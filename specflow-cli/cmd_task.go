package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/config"
	"github.com/stormbuf/specflow/internal/session"
	"github.com/stormbuf/specflow/internal/taskstore"
	"github.com/stormbuf/specflow/internal/vcs"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "任务管理",
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		intent, _ := cmd.Flags().GetString("intent")
		parent, _ := cmd.Flags().GetString("parent")

		if title == "" {
			return fmt.Errorf("必须指定 --title")
		}

		specflowDir := getSpecflowDir()
		cfg, err := config.Load(specflowDir)
		if err != nil {
			return fmt.Errorf("读取 config.yaml 失败: %w", err)
		}

		developer := readDeveloper(specflowDir)
		baseRev, _ := vcs.CurrentRev(getProjectDir(), cfg.VCS)

		task, err := taskstore.Create(specflowDir, taskstore.CreateOptions{
			Title:       title,
			Description: description,
			Intent:      intent,
			Developer:   developer,
			VCS:         cfg.VCS,
			BaseRev:     baseRev,
			Parent:      parent,
		})
		if err != nil {
			return err
		}

		// 设置 session 指针
		key := session.SessionKey(cfg.Platform, "cli")
		taskPath := filepath.Join(".specflow", "changes", task.ID)
		session.SetCurrentTask(specflowDir, key, cfg.Platform, "cli", taskPath)

		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"task_id":  task.ID,
				"task_dir": taskPath,
				"status":   task.Status,
				"parent":   task.Parent,
			})
			fmt.Println(string(data))
		} else {
			fmt.Printf("✅ 任务已创建: %s\n", task.ID)
			fmt.Printf("目录: %s\n", taskPath)
			if parent != "" {
				fmt.Printf("父任务: %s\n", parent)
			}
		}
		return nil
	},
}

var taskStartCmd = &cobra.Command{
	Use:   "start [task-dir]",
	Short: "开始任务（planning → in_progress）",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		key := session.SessionKey(cfg.Platform, "cli")

		var taskPath string
		if len(args) > 0 {
			taskPath = args[0]
		} else {
			tp, err := session.GetCurrentTask(specflowDir, key)
			if err != nil || tp == "" {
				return fmt.Errorf("无活跃任务，请指定 task-dir 或先 create")
			}
			taskPath = tp
		}

		// session 独占检查
		if err := session.CheckExclusive(specflowDir, key, taskPath); err != nil {
			if exclErr, ok := err.(*session.ExclusiveError); ok {
				if useJSON(cmd) {
					data, _ := json.Marshal(map[string]interface{}{
						"error":   "session_exclusive_conflict",
						"message": exclErr.Error(),
					})
					fmt.Println(string(data))
				} else {
					fmt.Printf("❌ %s\n", exclErr.Error())
				}
				os.Exit(3)
			}
			return err
		}

		taskID := filepath.Base(taskPath)
		task, err := taskstore.Load(specflowDir, taskID)
		if err != nil {
			return err
		}
		if err := task.Start(); err != nil {
			return err
		}
		taskDir := taskstore.TaskDir(specflowDir, taskID)
		if err := task.Save(taskDir); err != nil {
			return err
		}

		session.SetCurrentTask(specflowDir, key, cfg.Platform, "cli", taskPath)

		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"task_id":           task.ID,
				"status":            task.Status,
				"session_exclusive": true,
			})
			fmt.Println(string(data))
		} else {
			fmt.Printf("✅ 任务已开始: %s (in_progress)\n", task.ID)
		}
		return nil
	},
}

var taskFinishCmd = &cobra.Command{
	Use:   "finish",
	Short: "释放任务独占（清除 session 指针）",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		key := session.SessionKey(cfg.Platform, "cli")
		session.ClearCurrentTask(specflowDir, key)
		fmt.Println("✅ session 指针已清除，任务独占已释放")
		return nil
	},
}

var taskReleaseCmd = &cobra.Command{
	Use:   "release <task-id>",
	Short: "强制释放 stale session 指针",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		taskPath := filepath.Join(".specflow", "changes", args[0])
		count, err := session.Release(specflowDir, taskPath)
		if err != nil {
			return err
		}
		fmt.Printf("✅ 已释放 %d 个 stale session 指针\n", count)
		return nil
	},
}

var taskArchiveCmd = &cobra.Command{
	Use:   "archive [task-dir]",
	Short: "归档任务（completed + 移到 archive/）",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		key := session.SessionKey(cfg.Platform, "cli")

		var taskPath string
		if len(args) > 0 {
			taskPath = args[0]
		} else {
			tp, _ := session.GetCurrentTask(specflowDir, key)
			if tp == "" {
				return fmt.Errorf("无活跃任务")
			}
			taskPath = tp
		}

		taskID := filepath.Base(taskPath)
		task, err := taskstore.Load(specflowDir, taskID)
		if err != nil {
			return err
		}
		task.Complete()
		taskDir := taskstore.TaskDir(specflowDir, taskID)
		task.Save(taskDir)

		if err := taskstore.Archive(specflowDir, taskID); err != nil {
			return err
		}

		// VCS auto-commit
		vcs.AutoCommit(getProjectDir(), cfg.VCS, fmt.Sprintf("chore(task): archive %s", taskID))

		// 清 session
		session.ClearCurrentTask(specflowDir, key)

		archiveMonth := ""
		if task.CompletedAt != nil {
			archiveMonth = (*task.CompletedAt)[:7] // YYYY-MM
		}
		fmt.Printf("✅ 任务已归档: archive/%s/%s\n", archiveMonth, taskID)
		return nil
	},
}

var taskCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "显示当前活跃任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		key := session.SessionKey(cfg.Platform, "cli")

		taskPath, _ := session.GetCurrentTask(specflowDir, key)
		if taskPath == "" {
			if useJSON(cmd) {
				fmt.Println("null")
			} else {
				fmt.Println("无活跃任务")
			}
			return nil
		}

		taskID := filepath.Base(taskPath)
		task, err := taskstore.Load(specflowDir, taskID)
		if err != nil {
			return err
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"task_dir": taskPath,
				"task_id":  task.ID,
				"status":   task.Status,
				"title":    task.Title,
			})
			fmt.Println(string(data))
		} else {
			fmt.Printf("任务: %s\n", task.ID)
			fmt.Printf("标题: %s\n", task.Title)
			fmt.Printf("状态: %s\n", task.Status)
		}
		return nil
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		tasks, err := taskstore.ListAll(specflowDir)
		if err != nil {
			return err
		}

		// 过滤父任务（parent == nil）
		var parents []*taskstore.Task
		taskMap := make(map[string]*taskstore.Task)
		for _, t := range tasks {
			taskMap[t.ID] = t
			if t.Parent == nil {
				parents = append(parents, t)
			}
		}

		if useJSON(cmd) {
			var result []map[string]interface{}
			for _, p := range parents {
				item := map[string]interface{}{
					"id":       p.ID,
					"title":    p.Title,
					"status":   p.Status,
					"children": []map[string]interface{}{},
				}
				var children []map[string]interface{}
				for _, childID := range p.Children {
					if c, ok := taskMap[childID]; ok {
						children = append(children, map[string]interface{}{
							"id":     c.ID,
							"title":  c.Title,
							"status": c.Status,
						})
					}
				}
				item["children"] = children
				result = append(result, item)
			}
			data, _ := json.Marshal(result)
			fmt.Println(string(data))
		} else {
			for _, p := range parents {
				completed := 0
				for _, childID := range p.Children {
					if c, ok := taskMap[childID]; ok && c.Status == taskstore.StatusCompleted {
						completed++
					}
				}
				fmt.Printf("● %s — %s [%s]", p.ID, p.Title, p.Status)
				if len(p.Children) > 0 {
					fmt.Printf(" [%d/%d]", completed, len(p.Children))
				}
				fmt.Println()
				for _, childID := range p.Children {
					if c, ok := taskMap[childID]; ok {
						fmt.Printf("  ├─ %s — %s [%s]\n", c.ID, c.Title, c.Status)
					}
				}
			}
		}
		return nil
	},
}

var taskAddSubtaskCmd = &cobra.Command{
	Use:   "add-subtask <parent-dir> <child-dir>",
	Short: "关联父子任务",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		parentID := filepath.Base(args[0])
		childID := filepath.Base(args[1])
		return taskstore.AddSubtask(specflowDir, parentID, childID)
	},
}

var taskRemoveSubtaskCmd = &cobra.Command{
	Use:   "remove-subtask <parent-dir> <child-dir>",
	Short: "解除父子关联",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		parentID := filepath.Base(args[0])
		childID := filepath.Base(args[1])
		return taskstore.RemoveSubtask(specflowDir, parentID, childID)
	},
}

func readDeveloper(specflowDir string) string {
	data, err := os.ReadFile(filepath.Join(specflowDir, ".developer"))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func init() {
	taskCreateCmd.Flags().String("title", "", "任务标题")
	taskCreateCmd.Flags().String("description", "", "任务描述")
	taskCreateCmd.Flags().String("intent", "", "变更意图")
	taskCreateCmd.Flags().String("parent", "", "父任务 ID")

	taskCmd.AddCommand(taskCreateCmd, taskStartCmd, taskFinishCmd, taskReleaseCmd,
		taskArchiveCmd, taskCurrentCmd, taskListCmd, taskAddSubtaskCmd, taskRemoveSubtaskCmd)
	rootCmd.AddCommand(taskCmd)
}
