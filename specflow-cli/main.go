package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "specflow",
	Short: "Specflow CLI — 变更生命周期管理工具",
	Long:  "Specflow 是一个基于 spec + journal + workflow 的变更生命周期管理 CLI。",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolP("json", "j", false, "输出 JSON 格式")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细日志输出")
}

// getProjectDir 获取项目根目录（当前工作目录）
func getProjectDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

// getSpecflowDir 获取 .specflow/ 目录路径
func getSpecflowDir() string {
	return fmt.Sprintf("%s/.specflow", getProjectDir())
}

// printJSON 标记是否使用 JSON 输出
func useJSON(cmd *cobra.Command) bool {
	jsonFlag, _ := cmd.Root().PersistentFlags().GetBool("json")
	return jsonFlag
}
