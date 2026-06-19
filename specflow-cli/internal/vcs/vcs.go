package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Detect 自动检测项目使用的 VCS
// 优先 jj（jj 可同时管理 git 仓库），其次 git
func Detect(projectDir string) string {
	if _, err := os.Stat(filepath.Join(projectDir, ".jj")); err == nil {
		return "jj"
	}
	if _, err := os.Stat(filepath.Join(projectDir, ".git")); err == nil {
		return "git"
	}
	return ""
}

// HasUncommittedChanges 检查是否有未提交的改动
func HasUncommittedChanges(projectDir, vcsType string) (bool, error) {
	switch vcsType {
	case "git":
		out, err := exec.Command("git", "-C", projectDir, "status", "--porcelain").Output()
		if err != nil {
			return false, err
		}
		return strings.TrimSpace(string(out)) != "", nil
	case "jj":
		out, err := exec.Command("jj", "-R", projectDir, "status").Output()
		if err != nil {
			return false, err
		}
		return strings.TrimSpace(string(out)) != "", nil
	default:
		return false, fmt.Errorf("不支持的 VCS: %s", vcsType)
	}
}

// CurrentRev 获取当前版本号
func CurrentRev(projectDir, vcsType string) (string, error) {
	switch vcsType {
	case "git":
		out, err := exec.Command("git", "-C", projectDir, "rev-parse", "HEAD").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	case "jj":
		out, err := exec.Command("jj", "-R", projectDir, "log", "-r", "@-", "--no-graph", "-T", "commit_id").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	default:
		return "", fmt.Errorf("不支持的 VCS: %s", vcsType)
	}
}

// AutoCommit 自动提交
func AutoCommit(projectDir, vcsType, message string) error {
	switch vcsType {
	case "git":
		if err := exec.Command("git", "-C", projectDir, "add", "-A").Run(); err != nil {
			return err
		}
		return exec.Command("git", "-C", projectDir, "commit", "-m", message).Run()
	case "jj":
		if err := exec.Command("jj", "-R", projectDir, "describe", "-m", message).Run(); err != nil {
			return err
		}
		return exec.Command("jj", "-R", projectDir, "new").Run()
	default:
		return fmt.Errorf("不支持的 VCS: %s", vcsType)
	}
}
