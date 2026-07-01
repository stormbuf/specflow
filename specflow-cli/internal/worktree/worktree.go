package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/stormbuf/specflow/internal/vcs"
)

// Config 对应 .specflow/worktree.yaml
type Config struct {
	WorktreeDir string   `yaml:"worktree_dir"`
	Copy        []string `yaml:"copy"`
	PostCreate  []string `yaml:"post_create"`
	PreMerge    []string `yaml:"pre_merge"`
}

// LoadConfig 读取 worktree.yaml
func LoadConfig(specflowDir string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(specflowDir, "worktree.yaml"))
	if err != nil {
		return nil, fmt.Errorf("读取 worktree.yaml 失败: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 worktree.yaml 失败: %w", err)
	}
	if cfg.WorktreeDir == "" {
		cfg.WorktreeDir = "../{project}-worktrees"
	}
	return &cfg, nil
}

// Info worktree 信息
type Info struct {
	Name    string
	Path    string
	Branch  string
	Current bool
}

// Create 创建 worktree
func Create(projectDir, specflowDir, name, baseBranch string) (*Info, error) {
	cfg, err := LoadConfig(specflowDir)
	if err != nil {
		return nil, err
	}

	vcsType := vcs.Detect(projectDir)
	if vcsType == "" {
		return nil, fmt.Errorf("未检测到 VCS")
	}

	// 解析 worktree 目录
	projectName := filepath.Base(projectDir)
	wtDir := strings.ReplaceAll(cfg.WorktreeDir, "{project}", projectName)
	if !filepath.IsAbs(wtDir) {
		wtDir = filepath.Join(projectDir, wtDir)
	}

	// 确保 worktree 根目录存在
	if err := os.MkdirAll(filepath.Dir(wtDir), 0755); err != nil {
		return nil, fmt.Errorf("创建 worktree 目录失败: %w", err)
	}

	wtPath := filepath.Join(wtDir, name)
	branchName := fmt.Sprintf("github.com/stormbuf/specflow/%s", name)

	// 创建 worktree
	switch vcsType {
	case "git":
		args := []string{"-C", projectDir, "worktree", "add", "-b", branchName, wtPath}
		if baseBranch != "" {
			args = append(args, baseBranch)
		}
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			return nil, fmt.Errorf("git worktree add 失败: %w\n%s", err, string(out))
		}
	case "jj":
		// jj 使用 git worktree（jj 仓库底层是 git）
		// 先检查是否是 colocated repo
		if _, err := os.Stat(filepath.Join(projectDir, ".git")); err == nil {
			args := []string{"-C", projectDir, "worktree", "add", "-b", branchName, wtPath}
			if baseBranch != "" {
				args = append(args, baseBranch)
			}
			if err := exec.Command("git", args...).Run(); err != nil {
				return nil, fmt.Errorf("git worktree add 失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("jj 仓库需 colocated（含 .git/）才能使用 worktree，当前为纯 jj 仓库")
		}
	default:
		return nil, fmt.Errorf("不支持的 VCS: %s", vcsType)
	}

	// 复制文件
	for _, src := range cfg.Copy {
		srcPath := filepath.Join(projectDir, src)
		data, err := os.ReadFile(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("⚠️  复制 %s 失败: %v\n", src, err)
			continue
		}
		destPath := filepath.Join(wtPath, src)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			fmt.Printf("⚠️  创建目录 %s 失败: %v\n", filepath.Dir(destPath), err)
			continue
		}
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			fmt.Printf("⚠️  写入 %s 失败: %v\n", destPath, err)
			continue
		}
	}

	// 执行 post_create 命令
	for _, cmd := range cfg.PostCreate {
		trimmed := strings.TrimSpace(cmd)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		fmt.Printf("执行: %s\n", trimmed)
		parts := strings.Fields(trimmed)
		c := exec.Command(parts[0], parts[1:]...)
		c.Dir = wtPath
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("⚠️  post_create 命令失败: %s\n", trimmed)
		}
	}

	return &Info{
		Name:   name,
		Path:   wtPath,
		Branch: branchName,
	}, nil
}

// List 列出所有 worktree
func List(projectDir string) ([]Info, error) {
	vcsType := vcs.Detect(projectDir)
	if vcsType == "" {
		return nil, fmt.Errorf("未检测到 VCS")
	}

	var lines []string
	switch vcsType {
	case "git":
		out, err := exec.Command("git", "-C", projectDir, "worktree", "list", "--porcelain").Output()
		if err != nil {
			return nil, fmt.Errorf("git worktree list 失败: %w", err)
		}
		lines = strings.Split(string(out), "\n")
	case "jj":
		if _, err := os.Stat(filepath.Join(projectDir, ".git")); err != nil {
			return nil, fmt.Errorf("jj 仓库需 colocated 才能列出 worktree")
		}
		out, err := exec.Command("git", "-C", projectDir, "worktree", "list", "--porcelain").Output()
		if err != nil {
			return nil, fmt.Errorf("git worktree list 失败: %w", err)
		}
		lines = strings.Split(string(out), "\n")
	default:
		return nil, fmt.Errorf("不支持的 VCS: %s", vcsType)
	}

	var worktrees []Info
	var current Info
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Info{}
			}
			continue
		}
		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
			if current.Path == projectDir {
				current.Current = true
			}
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
		}
	}
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	// 提取 name
	for i := range worktrees {
		if worktrees[i].Branch != "" {
			worktrees[i].Name = strings.TrimPrefix(worktrees[i].Branch, "github.com/stormbuf/specflow/")
		}
	}

	return worktrees, nil
}

// Remove 移除 worktree
func Remove(projectDir, name string, force bool) error {
	vcsType := vcs.Detect(projectDir)
	if vcsType == "" {
		return fmt.Errorf("未检测到 VCS")
	}

	branchName := fmt.Sprintf("github.com/stormbuf/specflow/%s", name)

	switch vcsType {
	case "git":
		// 先找 worktree path
		worktrees, err := List(projectDir)
		if err != nil {
			return err
		}
		var wtPath string
		for _, wt := range worktrees {
			if wt.Name == name || wt.Branch == branchName {
				wtPath = wt.Path
				break
			}
		}
		if wtPath == "" {
			return fmt.Errorf("未找到 worktree: %s", name)
		}

		args := []string{"-C", projectDir, "worktree", "remove", wtPath}
		if force {
			args = append(args, "--force")
		}
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			return fmt.Errorf("git worktree remove 失败: %w\n%s\n提示: 如有未跟踪文件，使用 --force 强制移除", err, string(out))
		}

		// 删除分支
		branchArgs := []string{"-C", projectDir, "branch", "-d", branchName}
		if force {
			branchArgs = []string{"-C", projectDir, "branch", "-D", branchName}
		}
		exec.Command("git", branchArgs...).Run()
		return nil
	case "jj":
		if _, err := os.Stat(filepath.Join(projectDir, ".git")); err == nil {
			// colocated, 用 git worktree
			worktrees, err := List(projectDir)
			if err != nil {
				return err
			}
			var wtPath string
			for _, wt := range worktrees {
				if wt.Name == name || wt.Branch == branchName {
					wtPath = wt.Path
					break
				}
			}
			if wtPath == "" {
				return fmt.Errorf("未找到 worktree: %s", name)
			}
			args := []string{"-C", projectDir, "worktree", "remove", wtPath}
			if force {
				args = append(args, "--force")
			}
			if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
				return fmt.Errorf("git worktree remove 失败: %w\n%s\n提示: 如有未跟踪文件，使用 --force 强制移除", err, string(out))
			}
			exec.Command("git", "-C", projectDir, "branch", "-D", branchName).Run()
			return nil
		}
		return fmt.Errorf("jj 仓库需 colocated 才能移除 worktree")
	default:
		return fmt.Errorf("不支持的 VCS: %s", vcsType)
	}
}

// PreMerge 在合并前执行验证检查
func PreMerge(projectDir, specflowDir, name string) error {
	cfg, err := LoadConfig(specflowDir)
	if err != nil {
		return err
	}

	worktrees, err := List(projectDir)
	if err != nil {
		return err
	}

	var wtPath string
	for _, wt := range worktrees {
		if wt.Name == name {
			wtPath = wt.Path
			break
		}
	}
	if wtPath == "" {
		return fmt.Errorf("未找到 worktree: %s", name)
	}

	failed := false
	for _, cmd := range cfg.PreMerge {
		trimmed := strings.TrimSpace(cmd)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		fmt.Printf("检查: %s\n", trimmed)
		parts := strings.Fields(trimmed)
		c := exec.Command(parts[0], parts[1:]...)
		c.Dir = wtPath
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("❌ 失败: %s\n", trimmed)
			failed = true
		} else {
			fmt.Printf("✅ 通过: %s\n", trimmed)
		}
	}

	if failed {
		return fmt.Errorf("pre_merge 检查未通过，请修复后再合并")
	}
	return nil
}
