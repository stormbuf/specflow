package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"specflow/internal/config"
	"specflow/internal/context"
	"specflow/internal/session"
	"specflow/internal/taskstore"
)

var getContextCmd = &cobra.Command{
	Use:   "get-context",
	Short: "聚合 session 上下文（供 session-start.js 调用）",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)

		developer := readDeveloper(specflowDir)
		key := session.SessionKey(cfg.Platform, "cli")

		taskPath, _ := session.GetCurrentTask(specflowDir, key)
		var activeTask map[string]interface{}
		if taskPath != "" {
			taskID := filepath.Base(taskPath)
			task, err := taskstore.Load(specflowDir, taskID)
			if err == nil {
				activeTask = map[string]interface{}{
					"task_dir": taskPath,
					"task_id":  task.ID,
					"status":   task.Status,
					"title":    task.Title,
				}
			}
		}

		specIndexes := context.FindSpecIndexes(specflowDir)
		var indexPaths []string
		for _, idx := range specIndexes {
			indexPaths = append(indexPaths, idx.Path)
		}

		workflowOverview := context.WorkflowOverview(specflowDir)

		// 查找最新 journal
		journalLatest := ""
		journalDir := filepath.Join(specflowDir, "workspace", developer)
		if entries, err := readDirSorted(journalDir); err == nil && len(entries) > 0 {
			journalLatest = filepath.Join(".specflow", "workspace", developer, entries[0])
		}

		result := map[string]interface{}{
			"developer":         developer,
			"vcs":               cfg.VCS,
			"active_task":       activeTask,
			"spec_indexes":      indexPaths,
			"journal_latest":    journalLatest,
			"workflow_overview": workflowOverview,
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(result)
			fmt.Println(string(data))
		} else {
			fmt.Printf("开发者: %s\n", developer)
			fmt.Printf("VCS: %s\n", cfg.VCS)
			if activeTask != nil {
				fmt.Printf("活跃任务: %s (%s) — %s\n", activeTask["task_id"], activeTask["status"], activeTask["title"])
			} else {
				fmt.Println("活跃任务: 无")
			}
		}
		return nil
	},
}

var buildContextCmd = &cobra.Command{
	Use:   "build-context <agent-name>",
	Short: "构建 subagent 上下文（供 inject-subagent-context.js 调用）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		agentName := args[0]

		// 查 agents.yaml
		agentsCfg, err := config.LoadAgents(specflowDir)
		if err != nil {
			return fmt.Errorf("读取 agents.yaml 失败: %w", err)
		}
		agentConf, ok := agentsCfg.Agents[agentName]
		if !ok {
			return fmt.Errorf("agent '%s' 未在 agents.yaml 中声明", agentName)
		}

		// 解析任务目录
		key := session.SessionKey(cfg.Platform, "cli")
		taskPath, _ := session.GetCurrentTask(specflowDir, key)
		if taskPath == "" && agentConf.RequireTask {
			return fmt.Errorf("无活跃任务，agent '%s' 要求活跃任务", agentName)
		}

		taskDir := ""
		if taskPath != "" {
			taskDir = filepath.Join(getProjectDir(), taskPath)
		}

		// 构建上下文
		var contextStr string
		if agentConf.JSONLFile != nil && *agentConf.JSONLFile != "" && taskDir != "" {
			contextStr, err = context.BuildContext(getProjectDir(), taskDir, *agentConf.JSONLFile)
			if err != nil {
				return err
			}
		}

		// 输出上下文文本（非 JSON，供插件直接注入）
		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"agent":       agentName,
				"context":     contextStr,
				"constraints": agentConf.Constraints,
			})
			fmt.Println(string(data))
		} else {
			fmt.Print(contextStr)
		}
		return nil
	},
}

var addContextCmd = &cobra.Command{
	Use:   "add-context <task-dir> <agent-name> <file-path> <reason>",
	Short: "向 jsonl 追加上下文条目",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		taskDir := args[0]
		agentName := args[1]
		filePath := args[2]
		reason := args[3]

		// 查 agents.yaml 获取 jsonl 文件名
		agentsCfg, err := config.LoadAgents(specflowDir)
		if err != nil {
			return err
		}
		agentConf, ok := agentsCfg.Agents[agentName]
		if !ok {
			return fmt.Errorf("agent '%s' 未声明", agentName)
		}
		if agentConf.JSONLFile == nil || *agentConf.JSONLFile == "" {
			return fmt.Errorf("agent '%s' 未指定 jsonl_file", agentName)
		}

		fullTaskDir := filepath.Join(getProjectDir(), taskDir)
		return context.AddContext(fullTaskDir, *agentConf.JSONLFile, filePath, reason)
	},
}

func init() {
	rootCmd.AddCommand(getContextCmd, buildContextCmd, addContextCmd)
}
