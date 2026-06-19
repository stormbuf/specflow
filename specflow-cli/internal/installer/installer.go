package installer

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"specflow/internal/config"
	"specflow/internal/fingerprint"
	"specflow/internal/vcs"
)

// InstallMap 对应 install-map.yaml
type InstallMap struct {
	Platform       string            `yaml:"platform"`
	InstallTargets map[string]string `yaml:"install_targets"`
}

// LoadInstallMap 读取 install-map.yaml
func LoadInstallMap(embedFS embed.FS, platform string) (*InstallMap, error) {
	path := fmt.Sprintf("resources/platforms/%s/install-map.yaml", platform)
	data, err := embedFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 install-map.yaml 失败: %w", err)
	}
	var im InstallMap
	if err := yaml.Unmarshal(data, &im); err != nil {
		return nil, fmt.Errorf("解析 install-map.yaml 失败: %w", err)
	}
	return &im, nil
}

// InitOptions init 命令的选项
type InitOptions struct {
	Developer string
	Platform  string
	VCS       string
	Force     bool
}

// InitResult init 的结果
type InitResult struct {
	Platform string
	VCS      string
	Skills   int
	Plugins  int
	Agents   int
}

// Init 执行 specflow init
func Init(projectDir string, embedFS embed.FS, opts InitOptions) (*InitResult, error) {
	specflowDir := filepath.Join(projectDir, ".specflow")

	// 检查是否已存在
	if !opts.Force {
		if _, err := os.Stat(specflowDir); err == nil {
			return nil, fmt.Errorf(".specflow/ 已存在，使用 --force 覆盖")
		}
	}

	// 1. VCS 自动检测
	vcsType := opts.VCS
	if vcsType == "" {
		vcsType = vcs.Detect(projectDir)
		if vcsType == "" {
			return nil, fmt.Errorf("未检测到 VCS（.git/ 或 .jj/），请使用 --vcs 指定")
		}
	}

	// 2. 创建 .specflow/ 目录
	if err := os.MkdirAll(specflowDir, 0755); err != nil {
		return nil, fmt.Errorf("创建 .specflow/ 失败: %w", err)
	}

	// 3. 写入运行时模板
	if err := copyEmbeddedDir(embedFS, "resources/specflow-runtime", specflowDir); err != nil {
		return nil, fmt.Errorf("写入运行时模板失败: %w", err)
	}

	// 4. 写入 config.yaml（带实际配置）
	cfg := config.DefaultConfig()
	cfg.VCS = vcsType
	cfg.Platform = opts.Platform
	if err := cfg.Save(specflowDir); err != nil {
		return nil, fmt.Errorf("写入 config.yaml 失败: %w", err)
	}

	// 5. 写入 .developer
	if err := os.WriteFile(filepath.Join(specflowDir, ".developer"), []byte(opts.Developer), 0644); err != nil {
		return nil, fmt.Errorf("写入 .developer 失败: %w", err)
	}

	// 6. 写入 .vcs
	if err := os.WriteFile(filepath.Join(specflowDir, ".vcs"), []byte(vcsType), 0644); err != nil {
		return nil, fmt.Errorf("写入 .vcs 失败: %w", err)
	}

	// 7. 创建 .runtime/sessions/
	if err := os.MkdirAll(filepath.Join(specflowDir, ".runtime", "sessions"), 0755); err != nil {
		return nil, fmt.Errorf("创建 .runtime/ 失败: %w", err)
	}

	// 8. 创建 workspace/
	developerDir := filepath.Join(specflowDir, "workspace", opts.Developer)
	if err := os.MkdirAll(developerDir, 0755); err != nil {
		return nil, fmt.Errorf("创建 workspace/ 失败: %w", err)
	}

	// 9. 读取 install-map
	im, err := LoadInstallMap(embedFS, opts.Platform)
	if err != nil {
		return nil, err
	}

	// 10. 安装 skills
	skillsDir := filepath.Join(projectDir, im.InstallTargets["skills"])
	skillCount, err := copyEmbeddedDirWithCount(embedFS, "resources/skills", skillsDir)
	if err != nil {
		return nil, fmt.Errorf("安装 skills 失败: %w", err)
	}

	// 11. 安装 plugins
	pluginsDir := filepath.Join(projectDir, im.InstallTargets["plugins"])
	pluginCount, err := copyEmbeddedDirWithCount(embedFS, fmt.Sprintf("resources/platforms/%s/plugins", opts.Platform), pluginsDir)
	if err != nil {
		return nil, fmt.Errorf("安装 plugins 失败: %w", err)
	}

	// 12. 安装 native agents
	agentsDir := filepath.Join(projectDir, im.InstallTargets["agents"])
	agentCount, err := copyEmbeddedDirWithCount(embedFS, "resources/agents", agentsDir)
	if err != nil {
		return nil, fmt.Errorf("安装 agents 失败: %w", err)
	}

	// 13. 安装共享 lib
	libDir := filepath.Join(projectDir, filepath.Dir(strings.TrimRight(im.InstallTargets["plugins"], "/")), "lib")
	copyEmbeddedDir(embedFS, fmt.Sprintf("resources/platforms/%s/lib", opts.Platform), libDir)

	// 14. 记录文件指纹
	fp := &fingerprint.Fingerprints{
		SpecflowVersion: cfg.SpecflowVersion,
		Files:           make(map[string]string),
	}
	managedFiles, _ := fingerprint.CollectManagedFiles(projectDir)
	fp.RecordAll(projectDir, managedFiles)
	fp.Save(specflowDir)

	return &InitResult{
		Platform: opts.Platform,
		VCS:      vcsType,
		Skills:   skillCount,
		Plugins:  pluginCount,
		Agents:   agentCount,
	}, nil
}

// copyEmbeddedDir 递归复制 embed.FS 中的目录到目标路径
func copyEmbeddedDir(embedFS embed.FS, srcDir, destDir string) error {
	return fs.WalkDir(embedFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		data, err := embedFS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0644)
	})
}

// copyEmbeddedDirWithCount 复制并返回文件数量
func copyEmbeddedDirWithCount(embedFS embed.FS, srcDir, destDir string) (int, error) {
	count := 0
	return count, fs.WalkDir(embedFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		data, err := embedFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return err
		}
		count++
		return nil
	})
}

// SyncAgent 同步 custom agent 到平台目录
func SyncAgent(specflowDir, projectDir, agentName, platform string, embedFS embed.FS) (int, error) {
	im, err := LoadInstallMap(embedFS, platform)
	if err != nil {
		return 0, err
	}
	agentsDir := filepath.Join(projectDir, im.InstallTargets["agents"])

	agentsCfg, err := config.LoadAgents(specflowDir)
	if err != nil {
		return 0, err
	}

	count := 0
	for name, agent := range agentsCfg.Agents {
		if agentName != "all" && name != agentName {
			continue
		}
		if agent.Source != "custom" {
			continue
		}
		if agent.AgentFile == "" {
			continue
		}
		srcPath := filepath.Join(specflowDir, agent.AgentFile)
		data, err := os.ReadFile(srcPath)
		if err != nil {
			continue
		}
		destPath := filepath.Join(agentsDir, filepath.Base(agent.AgentFile))
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			continue
		}
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			continue
		}
		count++
	}
	return count, nil
}
