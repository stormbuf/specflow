package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"specflow/internal/installer"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 specflow 项目",
	Long:  "安装 specflow 三层结构到目标项目（skill 层 + 插件层 + 状态层）。",
	RunE: func(cmd *cobra.Command, args []string) error {
		developer, _ := cmd.Flags().GetString("user")
		opencode, _ := cmd.Flags().GetBool("opencode")
		pi, _ := cmd.Flags().GetBool("pi")
		platform, _ := cmd.Flags().GetString("platform")
		vcsType, _ := cmd.Flags().GetString("vcs")
		force, _ := cmd.Flags().GetBool("force")

		if developer == "" {
			return fmt.Errorf("必须指定 -u <developer>")
		}
		// 解析平台
		if platform == "" {
			if opencode {
				platform = "opencode"
			} else if pi {
				platform = "pi"
			} else {
				platform = "opencode"
			}
		}

		result, err := installer.Init(getProjectDir(), embeddedResources, installer.InitOptions{
			Developer: developer,
			Platform:  platform,
			VCS:       vcsType,
			Force:     force,
		})
		if err != nil {
			return err
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(result)
			fmt.Println(string(data))
		} else {
			fmt.Println("✅ specflow 已初始化")
			fmt.Printf("平台: %s\n", result.Platform)
			fmt.Printf("VCS: %s (自动检测)\n", result.VCS)
			fmt.Printf("开发者: %s\n", developer)
			fmt.Println("安装位置:")
			fmt.Printf("  - .specflow/ (运行时)\n")
			fmt.Printf("  - .opencode/skills/ (%d skills)\n", result.Skills)
			fmt.Printf("  - .opencode/plugins/ (%d plugins)\n", result.Plugins)
			fmt.Printf("  - .opencode/agents/ (%d native agents)\n", result.Agents)
			fmt.Println("文件指纹已记录")
			fmt.Println("\n请重启 AI Agent 以加载 specflow。")
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("user", "u", "", "开发者身份（必填）")
	initCmd.Flags().Bool("opencode", false, "目标平台 OpenCode")
	initCmd.Flags().Bool("pi", false, "目标平台 Pi")
	initCmd.Flags().String("platform", "", "目标平台（opencode | pi，与 --opencode/--pi 互斥）")
	initCmd.Flags().String("vcs", "", "版本管理工具: git | jj（自动检测）")
	initCmd.Flags().Bool("force", false, "覆盖已存在的 .specflow/ 目录")
	rootCmd.AddCommand(initCmd)
}
