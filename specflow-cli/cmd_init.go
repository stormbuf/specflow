package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"specflow/internal/installer"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 specflow 项目",
	Long:  "安装 specflow 三层结构到目标项目（skill 层 + 插件层 + 状态层），可选安装 spec 模板。",
	RunE: func(cmd *cobra.Command, args []string) error {
		developer, _ := cmd.Flags().GetString("user")
		opencode, _ := cmd.Flags().GetBool("opencode")
		pi, _ := cmd.Flags().GetBool("pi")
		platform, _ := cmd.Flags().GetString("platform")
		vcsType, _ := cmd.Flags().GetString("vcs")
		force, _ := cmd.Flags().GetBool("force")
		noSpec, _ := cmd.Flags().GetBool("no-spec")
		allSpec, _ := cmd.Flags().GetBool("all-spec")
		withSpec, _ := cmd.Flags().GetString("with-spec")

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

		// spec 模板安装
		specCount := 0
		specCategories := []string{}
		if !noSpec {
			categories := installer.ListSpecCategories(embeddedResources)
			if len(categories) > 0 {
				if allSpec {
					for _, c := range categories {
						specCategories = append(specCategories, c.Name)
					}
				} else if withSpec != "" {
					// 验证分类名
					validMap := make(map[string]bool)
					for _, c := range categories {
						validMap[c.Name] = true
					}
					for _, name := range strings.Split(withSpec, ",") {
						name = strings.TrimSpace(name)
						if !validMap[name] {
							return fmt.Errorf("未知的 spec 分类: %s（使用 specflow spec list 查看可用分类）", name)
						}
						specCategories = append(specCategories, name)
					}
				} else if !useJSON(cmd) {
					// 交互式选择
					fmt.Println("\n--- Spec 模板安装 ---")
					selected, err := interactiveSelect(categories)
					if err != nil {
						fmt.Printf("⚠️  spec 模板选择失败: %v（跳过）\n", err)
					} else {
						specCategories = selected
					}
				}

				if len(specCategories) > 0 {
					specCount, err = installer.InstallSpecTemplates(getProjectDir(), embeddedResources, specCategories)
					if err != nil {
						fmt.Printf("⚠️  spec 模板安装失败: %v\n", err)
					}
				}
			}
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(struct {
				*installer.InitResult
				SpecTemplates int      `json:"spec_templates"`
				SpecCategories []string `json:"spec_categories"`
			}{
				InitResult:     result,
				SpecTemplates:  specCount,
				SpecCategories: specCategories,
			})
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
			if specCount > 0 {
				fmt.Printf("  - .specflow/spec/ (%d spec 模板文件, 分类: %s)\n", specCount, strings.Join(specCategories, ", "))
			}
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
	initCmd.Flags().Bool("no-spec", false, "跳过 spec 模板安装")
	initCmd.Flags().Bool("all-spec", false, "安装所有 spec 模板")
	initCmd.Flags().String("with-spec", "", "指定 spec 模板分类（逗号分隔，如 guides,backend）")
	rootCmd.AddCommand(initCmd)
}
