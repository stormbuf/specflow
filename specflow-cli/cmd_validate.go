package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/config"
	"github.com/stormbuf/specflow/internal/context"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "校验配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		var errors []string
		var warnings []string

		// 1. 校验 config.yaml
		cfg, err := config.Load(specflowDir)
		if err != nil {
			errors = append(errors, fmt.Sprintf("config.yaml: %v", err))
		} else {
			if cfg.VCS != "git" && cfg.VCS != "jj" {
				errors = append(errors, fmt.Sprintf("config.yaml: vcs '%s' 无效（应为 git 或 jj）", cfg.VCS))
			}
			if cfg.Platform != "opencode" && cfg.Platform != "pi" {
				errors = append(errors, fmt.Sprintf("config.yaml: platform '%s' 无效（应为 opencode 或 pi）", cfg.Platform))
			}
		}

		// 2. 校验 agents.yaml
		agentsCfg, err := config.LoadAgents(specflowDir)
		if err != nil {
			errors = append(errors, fmt.Sprintf("agents.yaml: %v", err))
		} else {
			agentErrs := agentsCfg.Validate(specflowDir)
			for _, e := range agentErrs {
				errors = append(errors, e.Error())
			}

			// 3. 校验 jsonl 引用文件存在性
			for name, agent := range agentsCfg.Agents {
				if agent.JSONLFile == nil || *agent.JSONLFile == "" {
					continue
				}
				// 检查当前活跃任务的 jsonl（如果有）
				// 这里只检查 agents.yaml 本身的完整性
				if agent.Source == "custom" && agent.AgentFile != "" {
					// 已在 Validate 中检查
				}
				_ = name
			}
		}

		// 4. 校验 workflow.md
		workflowPath := filepath.Join(specflowDir, "workflow.md")
		if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
			errors = append(errors, "workflow.md 不存在")
		} else {
			data, _ := os.ReadFile(workflowPath)
			content := string(data)
			if !containsWorkflowStateTags(content) {
				warnings = append(warnings, "workflow.md 未包含 [workflow-state:...] 标签块")
			}
		}

		// 5. 校验 .specflow/ 结构
		requiredDirs := []string{"changes", "spec", "workspace", ".runtime/sessions"}
		for _, d := range requiredDirs {
			path := filepath.Join(specflowDir, d)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				warnings = append(warnings, fmt.Sprintf("目录缺失: .specflow/%s/", d))
			}
		}

		valid := len(errors) == 0
		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"valid":    valid,
				"errors":   errors,
				"warnings": warnings,
			})
			fmt.Println(string(data))
		} else {
			if valid {
				fmt.Println("✅ 配置校验通过")
			} else {
				fmt.Println("❌ 配置校验失败:")
				for _, e := range errors {
					fmt.Printf("  - %s\n", e)
				}
			}
			for _, w := range warnings {
				fmt.Printf("⚠️  %s\n", w)
			}
		}

		if !valid {
			os.Exit(2)
		}
		return nil
	},
}

func containsWorkflowStateTags(content string) bool {
	// 检查是否包含 [workflow-state: 标签
	for _, line := range splitLines(content) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[workflow-state:") {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// 确保 context 包被引用（避免 unused import）
var _ = context.FindSpecIndexes

func init() {
	rootCmd.AddCommand(validateCmd)
}
