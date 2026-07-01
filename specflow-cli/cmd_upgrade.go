package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/version"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级全局 CLI 二进制（从 GitHub Releases 下载）",
	RunE: func(cmd *cobra.Command, args []string) error {
		channel, _ := cmd.Flags().GetString("channel")
		force, _ := cmd.Flags().GetBool("force")

		fmt.Printf("当前版本: %s\n", version.Version)
		fmt.Printf("查询最新版本 (channel: %s)...\n", channel)

		// 查询 GitHub Releases
		releaseURL := "https://api.github.com/repos/stormbuf/specflow/releases/latest"
		if channel == "beta" {
			// beta 通道：列出所有 release，取第一个（包含 prerelease）
			releaseURL = "https://api.github.com/repos/stormbuf/specflow/releases?per_page=1"
		}

		req, _ := http.NewRequest("GET", releaseURL, nil)
		req.Header.Set("Accept", "application/vnd.github+json")
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("查询 GitHub Releases 失败: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("GitHub API 返回 %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)

		var release githubRelease
		if channel == "beta" {
			var releases []githubRelease
			if err := json.Unmarshal(body, &releases); err != nil || len(releases) == 0 {
				return fmt.Errorf("解析 release 失败")
			}
			release = releases[0]
		} else {
			if err := json.Unmarshal(body, &release); err != nil {
				return fmt.Errorf("解析 release 失败: %w", err)
			}
		}

		latestVersion := strings.TrimPrefix(release.TagName, "v")
		fmt.Printf("最新版本: %s\n", latestVersion)

		if latestVersion == version.Version && !force {
			fmt.Println("✅ 已是最新版本")
			return nil
		}

		// 匹配平台资产
		goos := runtime.GOOS
		goarch := runtime.GOARCH
		assetName := fmt.Sprintf("specflow_%s_%s", goos, goarch)

		var downloadURL string
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, assetName) {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}

		if downloadURL == "" {
			return fmt.Errorf("未找到匹配的二进制: %s (可用资产: %v)", assetName, assetNames(release.Assets))
		}

		fmt.Printf("下载: %s\n", downloadURL)

		// 下载 tar.gz
		resp2, err := http.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("下载失败: %w", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != 200 {
			return fmt.Errorf("下载失败: HTTP %d", resp2.StatusCode)
		}

		// 获取当前二进制路径
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("获取当前二进制路径失败: %w", err)
		}

		// 写入临时文件
		tmpFile := execPath + ".new"
		out, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("创建临时文件失败: %w", err)
		}

		// 写入下载内容（tar.gz 包含单个 specflow 二进制）
		if _, err := io.Copy(out, resp2.Body); err != nil {
			out.Close()
			os.Remove(tmpFile)
			return fmt.Errorf("写入失败: %w", err)
		}
		out.Close()

		// 解压 tar.gz
		// 简单处理：如果是 tar.gz，用 tar 命令解压
		if strings.HasSuffix(downloadURL, ".tar.gz") {
			tmpDir := execPath + ".tmpdir"
			os.MkdirAll(tmpDir, 0755)
			defer os.RemoveAll(tmpDir)

			if err := runCommand("tar", "xzf", tmpFile, "-C", tmpDir); err != nil {
				os.Remove(tmpFile)
				return fmt.Errorf("解压失败: %w（请手动安装: 下载 %s 解压后替换 %s）", err, downloadURL, execPath)
			}
			os.Remove(tmpFile)

			// 找到解压出来的 specflow 二进制
			entries, _ := os.ReadDir(tmpDir)
			binaryPath := ""
			for _, e := range entries {
				if e.Name() == "specflow" {
					binaryPath = fmt.Sprintf("%s/%s", tmpDir, e.Name())
					break
				}
			}
			if binaryPath == "" {
				return fmt.Errorf("解压后未找到 specflow 二进制")
			}

			// 复制到新文件
			if err := copyFile(binaryPath, tmpFile); err != nil {
				return fmt.Errorf("复制二进制失败: %w", err)
			}
		}

		// 替换旧二进制
		oldFile := execPath + ".old"
		os.Remove(oldFile)

		if err := os.Rename(execPath, oldFile); err != nil {
			os.Remove(tmpFile)
			return fmt.Errorf("备份旧版本失败: %w", err)
		}

		if err := os.Rename(tmpFile, execPath); err != nil {
			// 回滚
			os.Rename(oldFile, execPath)
			os.Remove(tmpFile)
			return fmt.Errorf("替换二进制失败: %w", err)
		}

		os.Chmod(execPath, 0755)
		os.Remove(oldFile)

		fmt.Printf("✅ 已升级到 %s\n", latestVersion)
		fmt.Println("请重启终端使新版本生效。")
		return nil
	},
}

func assetNames(assets []struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}) []string {
	names := make([]string, len(assets))
	for i, a := range assets {
		names[i] = a.Name
	}
	return names
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0755)
}

func init() {
	upgradeCmd.Flags().String("channel", "latest", "升级通道: latest|beta")
	upgradeCmd.Flags().Bool("force", false, "强制升级（即使版本相同）")
	rootCmd.AddCommand(upgradeCmd)
}
