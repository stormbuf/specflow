package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/config"
	"github.com/stormbuf/specflow/internal/fingerprint"
	"github.com/stormbuf/specflow/internal/installer"
	"github.com/stormbuf/specflow/internal/version"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "把当前项目同步到本地 CLI 版本（含文件指纹冲突检测）",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		projectDir := getProjectDir()
		force, _ := cmd.Flags().GetBool("force")

		// 加载配置
		cfg, err := config.Load(specflowDir)
		if err != nil {
			return fmt.Errorf("读取 config.yaml 失败: %w", err)
		}

		// 加载旧指纹
		oldFP, err := fingerprint.Load(specflowDir)
		if err != nil {
			return fmt.Errorf("读取指纹失败: %w", err)
		}

		// 读取 install-map（用于确定安装目标路径）
		_, err = installer.LoadInstallMap(embeddedResources, cfg.Platform)
		if err != nil {
			return err
		}

		// 收集当前管理文件
		managedFiles, _ := fingerprint.CollectManagedFiles(projectDir)

		var updated []string
		var skipped []string
		var conflicts []map[string]string

		for _, relPath := range managedFiles {
			// 获取 CLI 内嵌的新版本内容
			newContent, err := getEmbeddedFile(relPath, cfg.Platform)
			if err != nil {
				// 文件不在 embed 中（可能是用户新增的），跳过
				skipped = append(skipped, relPath)
				continue
			}

			oldHash := oldFP.Files[relPath]
			curHash, _ := fingerprint.HashFile(filepath.Join(projectDir, relPath))
			newHash := fingerprint.HashBytes(newContent)

			result := fingerprint.ThreeWayCompare(oldHash, curHash, newHash)

			switch result {
			case fingerprint.MatchUserUnchanged:
			// 用户未修改 → 检查内容是否真的不同
			if curHash != newHash {
				writeFile(projectDir, relPath, newContent)
				updated = append(updated, relPath)
			} else {
				skipped = append(skipped, relPath)
			}

			case fingerprint.MatchCLIUnchanged:
				// CLI 未更新 → 保留用户版本
				skipped = append(skipped, relPath)

			case fingerprint.Conflict:
				// 冲突
				if force {
					writeFile(projectDir, relPath, newContent)
					updated = append(updated, relPath)
					conflicts = append(conflicts, map[string]string{
						"file":       relPath,
						"resolution": "overwrite (--force)",
					})
				} else {
					// 写入 update-candidates
					candidateDir := filepath.Join(specflowDir, ".update-candidates")
					os.MkdirAll(candidateDir, 0755)
					candidatePath := filepath.Join(candidateDir, strings.ReplaceAll(relPath, "/", "_"))
					os.WriteFile(candidatePath, newContent, 0644)
					skipped = append(skipped, relPath)
					conflicts = append(conflicts, map[string]string{
						"file":       relPath,
						"resolution": "merge (写入 .update-candidates/)",
					})
				}

			case fingerprint.NewFile:
				// 新文件（旧指纹中不存在）→ 直接写入
				writeFile(projectDir, relPath, newContent)
				updated = append(updated, relPath)
			}
		}

		// 刷新指纹
	newFP := &fingerprint.Fingerprints{
		SpecflowVersion: version.Version,
		Files:           make(map[string]string),
	}
	newFP.RecordAll(projectDir, managedFiles)
	newFP.Save(specflowDir)

	// 更新 AGENTS.md managed block
	agentsUpdated, _ := installer.UpdateAgentsMd(projectDir, embeddedResources)

	if useJSON(cmd) {
		data, _ := json.Marshal(map[string]interface{}{
			"updated":           updated,
			"skipped":           skipped,
			"conflicts":         conflicts,
			"agents_md_updated": agentsUpdated,
		})
		fmt.Println(string(data))
	} else {
		fmt.Printf("✅ 已更新 %d 个文件\n", len(updated))
		for _, f := range updated {
			fmt.Printf("  更新: %s\n", f)
		}
		fmt.Printf("⏭️  跳过 %d 个文件\n", len(skipped))
		for _, f := range skipped {
			fmt.Printf("  跳过: %s\n", f)
		}
		if len(conflicts) > 0 {
			fmt.Printf("⚠️  %d 个冲突:\n", len(conflicts))
			for _, c := range conflicts {
				fmt.Printf("  %s → %s\n", c["file"], c["resolution"])
			}
		}
		if agentsUpdated {
			fmt.Println("📝 AGENTS.md managed block 已更新")
		}
	}
		return nil
	},
}

// getEmbeddedFile 从 embed.FS 中读取文件
func getEmbeddedFile(relPath, platform string) ([]byte, error) {
	// 映射项目路径到 embed 路径
	var embedPath string
	switch {
	case strings.HasPrefix(relPath, ".specflow/workflow.md"):
		embedPath = "resources/specflow-runtime/workflow.md"
	case strings.HasPrefix(relPath, ".opencode/skills/"):
		relative := strings.TrimPrefix(relPath, ".opencode/skills/")
		embedPath = filepath.Join("resources/skills", relative)
	case strings.HasPrefix(relPath, ".opencode/plugins/"):
		relative := strings.TrimPrefix(relPath, ".opencode/plugins/")
		embedPath = filepath.Join(fmt.Sprintf("resources/platforms/%s/plugins", platform), relative)
	case strings.HasPrefix(relPath, ".opencode/agents/specflow-"):
		filename := filepath.Base(relPath)
		embedPath = filepath.Join("resources/agents", filename)
	default:
		return nil, fmt.Errorf("不在管理范围: %s", relPath)
	}
	return embeddedResources.ReadFile(filepath.ToSlash(embedPath))
}

func writeFile(projectDir, relPath string, content []byte) error {
	fullPath := filepath.Join(projectDir, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, content, 0644)
}

func init() {
	updateCmd.Flags().Bool("force", false, "跳过冲突询问，全部覆盖")
	rootCmd.AddCommand(updateCmd)
}

// 确保 fs 包被引用
var _ = fs.WalkDir
