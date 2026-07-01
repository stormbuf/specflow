package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/config"
	"github.com/stormbuf/specflow/internal/installer"
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

var agentsListCmd = &cobra.Command{
	Use:   "agents list",
	Short: "列出 agents.yaml 中配置的所有 agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		agentsCfg, err := config.LoadAgents(specflowDir)
		if err != nil {
			return fmt.Errorf("读取 agents.yaml 失败: %w", err)
		}

		if useJSON(cmd) {
			data, _ := json.MarshalIndent(agentsCfg.Agents, "", "  ")
			fmt.Println(string(data))
		} else {
			if len(agentsCfg.Agents) == 0 {
				fmt.Println("未配置任何 agent")
				return nil
			}
			fmt.Printf("%-25s %-12s %-20s %s\n", "NAME", "SOURCE", "JSONL", "CONSTRAINTS")
			for name, agent := range agentsCfg.Agents {
				jsonl := ""
				if agent.JSONLFile != nil {
					jsonl = *agent.JSONLFile
				}
				fmt.Printf("%-25s %-12s %-20s %d rules\n", name, agent.Source, jsonl, len(agent.Constraints))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncAgentCmd, agentsListCmd)
}
