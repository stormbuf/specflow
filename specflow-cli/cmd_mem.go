package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/config"
)

var memCmd = &cobra.Command{
	Use:   "mem",
	Short: "跨会话对话检索",
}

var memListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出可检索的项目与会话",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)

		if !cfg.Mem.Enabled {
			if useJSON(cmd) {
				fmt.Println(`{"enabled": false}`)
			} else {
				fmt.Println("mem 未启用（config.yaml: mem.enabled = false）")
			}
			return nil
		}

		sessions := findMemSessions(cfg.Mem.LogPaths)
		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"enabled":  true,
				"sessions": sessions,
			})
			fmt.Println(string(data))
		} else {
			fmt.Printf("找到 %d 个会话日志:\n", len(sessions))
			for _, s := range sessions {
				fmt.Printf("  - %s\n", s)
			}
		}
		return nil
	},
}

var memSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "按关键词检索历史对话",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		query := args[0]
		phase, _ := cmd.Flags().GetString("phase")
		limit, _ := cmd.Flags().GetInt("limit")

		results := searchMemSessions(cfg.Mem.LogPaths, query, phase, limit)
		if useJSON(cmd) {
			data, _ := json.Marshal(results)
			fmt.Println(string(data))
		} else {
			if len(results) == 0 {
				fmt.Println("未找到匹配的对话片段")
			} else {
				for _, r := range results {
					fmt.Printf("--- %s ---\n", r["file"])
					fmt.Println(r["snippet"])
					fmt.Println()
				}
			}
		}
		return nil
	},
}

var memContextCmd = &cobra.Command{
	Use:   "context <query>",
	Short: "检索并输出上下文片段",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		query := args[0]
		phase, _ := cmd.Flags().GetString("phase")

		results := searchMemSessions(cfg.Mem.LogPaths, query, phase, 5)
		if len(results) == 0 {
			fmt.Println("未找到匹配的上下文")
			return nil
		}
		for _, r := range results {
			fmt.Printf("=== %s ===\n%s\n\n", r["file"], r["snippet"])
		}
		return nil
	},
}

func findMemSessions(logPaths []string) []string {
	var sessions []string
	for _, logPath := range logPaths {
		expanded := expandPath(logPath)
		entries, err := os.ReadDir(expanded)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				sessions = append(sessions, filepath.Join(expanded, entry.Name()))
			}
		}
	}
	return sessions
}

func searchMemSessions(logPaths []string, query string, phase string, limit int) []map[string]string {
	var results []map[string]string
	queryLower := strings.ToLower(query)

	for _, logPath := range logPaths {
		expanded := expandPath(logPath)
		filepath.Walk(expanded, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".jsonl") {
				return nil
			}
			if len(results) >= limit {
				return filepath.SkipDir
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			content := string(data)
			contentLower := strings.ToLower(content)
			idx := strings.Index(contentLower, queryLower)
			if idx == -1 {
				return nil
			}

			// 提取上下文片段（前后各 500 字符）
			start := idx - 500
			if start < 0 {
				start = 0
			}
			end := idx + len(query) + 500
			if end > len(content) {
				end = len(content)
			}
			snippet := content[start:end]

			results = append(results, map[string]string{
				"file":    path,
				"snippet": snippet,
			})
			return nil
		})
	}
	return results
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func init() {
	memSearchCmd.Flags().String("phase", "all", "检索阶段: brainstorm|implement|all")
	memSearchCmd.Flags().Int("limit", 10, "最大结果数")
	memContextCmd.Flags().String("phase", "all", "检索阶段: brainstorm|implement|all")

	memCmd.AddCommand(memListCmd, memSearchCmd, memContextCmd)
	rootCmd.AddCommand(memCmd)
}
