package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/worktree"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "管理多 agent worktree",
	Long:  "创建、列出、移除 worktree，用于多 agent 并行工作场景。",
}

var worktreeCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "创建 worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		base, _ := cmd.Flags().GetString("base")

		info, err := worktree.Create(getProjectDir(), getSpecflowDir(), name, base)
		if err != nil {
			return err
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(info)
			fmt.Println(string(data))
		} else {
			fmt.Printf("✅ worktree 已创建\n")
			fmt.Printf("名称: %s\n", info.Name)
			fmt.Printf("路径: %s\n", info.Path)
			fmt.Printf("分支: %s\n", info.Branch)
		}
		return nil
	},
}

var worktreeListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 worktree",
	RunE: func(cmd *cobra.Command, args []string) error {
		worktrees, err := worktree.List(getProjectDir())
		if err != nil {
			return err
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(worktrees)
			fmt.Println(string(data))
			return nil
		}

		if len(worktrees) == 0 {
			fmt.Println("暂无 worktree")
			return nil
		}

		fmt.Printf("%-20s %-40s %-25s %s\n", "NAME", "PATH", "BRANCH", "CURRENT")
		for _, wt := range worktrees {
			current := ""
			if wt.Current {
				current = "*"
			}
			fmt.Printf("%-20s %-40s %-25s %s\n", wt.Name, wt.Path, wt.Branch, current)
		}
		return nil
	},
}

var worktreeRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "移除 worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if err := worktree.Remove(getProjectDir(), name, force); err != nil {
			return err
		}

		fmt.Printf("✅ worktree %s 已移除\n", name)
		return nil
	},
}

var worktreeMergeCmd = &cobra.Command{
	Use:   "merge <name>",
	Short: "合并前验证检查（pre_merge）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if err := worktree.PreMerge(getProjectDir(), getSpecflowDir(), name); err != nil {
			return err
		}

		fmt.Printf("✅ pre_merge 检查全部通过，可以合并 %s\n", name)
		fmt.Println("合并命令:")
		fmt.Printf("  git merge specflow/%s\n", name)
		return nil
	},
}

func init() {
	worktreeCreateCmd.Flags().String("base", "", "基于哪个分支创建（默认当前分支）")
	worktreeRemoveCmd.Flags().Bool("force", false, "强制移除（忽略未提交改动）")
	worktreeCmd.AddCommand(worktreeCreateCmd)
	worktreeCmd.AddCommand(worktreeListCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
	worktreeCmd.AddCommand(worktreeMergeCmd)
	rootCmd.AddCommand(worktreeCmd)
}
