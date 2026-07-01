package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormbuf/specflow/internal/installer"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "管理 spec 模板",
	Long:  "列出、安装 spec 模板到 .specflow/spec/ 目录。",
}

var specListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出可用的 spec 模板分类",
	RunE: func(cmd *cobra.Command, args []string) error {
		categories := installer.ListSpecCategories(embeddedResources)
		if len(categories) == 0 {
			fmt.Println("暂无可用 spec 模板")
			return nil
		}

		if useJSON(cmd) {
			fmt.Printf("[")
			for i, c := range categories {
				if i > 0 {
					fmt.Printf(",")
				}
				fmt.Printf(`{"name":"%s","description":"%s","files":%d}`, c.Name, c.Description, c.FileCount)
			}
			fmt.Println("]")
			return nil
		}

		fmt.Println("可用的 spec 模板分类：")
		fmt.Println()
		for i, c := range categories {
			fmt.Printf("  [%d] %-15s %s (%d 个文件)\n", i+1, c.Name, c.Description, c.FileCount)
		}
		fmt.Println()
		fmt.Println("使用 specflow spec install 安装选定的分类")
		return nil
	},
}

var specInstallCmd = &cobra.Command{
	Use:   "install [category...]",
	Short: "安装 spec 模板（交互式多选或指定分类）",
	Long:  "安装 spec 模板到 .specflow/spec/。不带参数时进入交互式多选；指定分类名时直接安装。",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		categories := installer.ListSpecCategories(embeddedResources)
		if len(categories) == 0 {
			return fmt.Errorf("暂无可用 spec 模板")
		}

		var selected []string

		if all {
			for _, c := range categories {
				selected = append(selected, c.Name)
			}
		} else if len(args) > 0 {
			// 验证分类名
			validMap := make(map[string]bool)
			for _, c := range categories {
				validMap[c.Name] = true
			}
			for _, arg := range args {
				if !validMap[arg] {
					return fmt.Errorf("未知的分类名: %s（使用 specflow spec list 查看可用分类）", arg)
				}
				selected = append(selected, arg)
			}
		} else {
			// 交互式多选
			var err error
			selected, err = interactiveSelect(categories)
			if err != nil {
				return err
			}
		}

		if len(selected) == 0 {
			fmt.Println("未选择任何分类，退出。")
			return nil
		}

		count, err := installer.InstallSpecTemplates(getProjectDir(), embeddedResources, selected)
		if err != nil {
			return err
		}

		fmt.Printf("✅ 已安装 %d 个 spec 模板文件到 .specflow/spec/\n", count)
		fmt.Printf("已安装分类: %s\n", strings.Join(selected, ", "))
		return nil
	},
}

// interactiveSelect 交互式多选
func interactiveSelect(categories []installer.SpecCategory) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("可用的 spec 模板分类：")
	fmt.Println()
	for i, c := range categories {
		fmt.Printf("  [%d] %-15s %s (%d 个文件)\n", i+1, c.Name, c.Description, c.FileCount)
	}
	fmt.Println()
	fmt.Print("选择要安装的分类（输入编号，逗号分隔；输入 all 全选；直接回车跳过）: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取输入失败: %w", err)
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, nil
	}

	if strings.ToLower(input) == "all" {
		var all []string
		for _, c := range categories {
			all = append(all, c.Name)
		}
		return all, nil
	}

	// 解析逗号分隔的编号
	var selected []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.Atoi(part)
		if err != nil || num < 1 || num > len(categories) {
			return nil, fmt.Errorf("无效的编号: %s", part)
		}
		selected = append(selected, categories[num-1].Name)
	}
	return selected, nil
}

func init() {
	specInstallCmd.Flags().Bool("all", false, "安装所有分类")
	specCmd.AddCommand(specListCmd)
	specCmd.AddCommand(specInstallCmd)
	rootCmd.AddCommand(specCmd)
}
