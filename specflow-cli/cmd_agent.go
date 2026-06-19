package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"specflow/internal/config"
	"specflow/internal/installer"
)

var syncAgentCmd = &cobra.Command{
	Use:   "sync-agent <name>",
	Short: "同步 custom agent 到平台 agent 目录",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)
		agentName := args[0]

		count, err := installer.SyncAgent(specflowDir, getProjectDir(), agentName, cfg.Platform, embeddedResources)
		if err != nil {
			return err
		}
		if count == 0 {
			fmt.Printf("未找到需要同步的 custom agent: %s\n", agentName)
		} else {
			fmt.Printf("✅ 已同步 %d 个 custom agent\n", count)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncAgentCmd)
}
