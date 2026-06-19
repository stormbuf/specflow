package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"specflow/internal/config"
)

var addSessionCmd = &cobra.Command{
	Use:   "add-session",
	Short: "向 journal 追加 session 条目",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		cfg, _ := config.Load(specflowDir)

		title, _ := cmd.Flags().GetString("title")
		summary, _ := cmd.Flags().GetString("summary")
		taskDir, _ := cmd.Flags().GetString("task")

		if title == "" {
			return fmt.Errorf("必须指定 --title")
		}

		developer := readDeveloper(specflowDir)
		journalDir := filepath.Join(specflowDir, "workspace", developer)
		if err := os.MkdirAll(journalDir, 0755); err != nil {
			return err
		}

		// 查找当前 journal 文件
		journalFile := findLatestJournal(journalDir)
		if journalFile == "" {
			journalFile = "journal-1.md"
		}

		journalPath := filepath.Join(journalDir, journalFile)

		// 检查行数，超限则轮转
		if shouldRotate(journalPath, cfg.MaxJournalLines) {
			num := extractJournalNum(journalFile)
			journalFile = fmt.Sprintf("journal-%d.md", num+1)
			journalPath = filepath.Join(journalDir, journalFile)
			// 写入 header
			os.WriteFile(journalPath, []byte(fmt.Sprintf("# Journal %d\n\n", num+1)), 0644)
		}

		// 构建条目
		now := time.Now()
		var entry strings.Builder
		entry.WriteString(fmt.Sprintf("\n## %s\n", title))
		if taskDir != "" {
			taskID := filepath.Base(taskDir)
			entry.WriteString(fmt.Sprintf("**任务**: %s\n", taskID))
		}
		entry.WriteString(fmt.Sprintf("**日期**: %s\n", now.Format("2006-01-02 15:04")))
		entry.WriteString(fmt.Sprintf("**摘要**: %s\n", summary))
		entry.WriteString("\n")

		// 追加
		f, err := os.OpenFile(journalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(entry.String())
		if err != nil {
			return err
		}

		fmt.Printf("✅ session 已记录到 %s\n", journalFile)
		return nil
	},
}

func findLatestJournal(dir string) string {
	names, err := readDirSorted(dir)
	if err != nil || len(names) == 0 {
		return ""
	}
	// 找到 journal-N.md 格式的文件
	for _, name := range names {
		if strings.HasPrefix(name, "journal-") && strings.HasSuffix(name, ".md") {
			return name
		}
	}
	return ""
}

func shouldRotate(path string, maxLines int) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	if maxLines <= 0 {
		maxLines = 2000
	}
	lines := strings.Count(string(data), "\n")
	return lines >= maxLines
}

func extractJournalNum(filename string) int {
	// journal-N.md → N
	name := strings.TrimPrefix(filename, "journal-")
	name = strings.TrimSuffix(name, ".md")
	var num int
	fmt.Sscanf(name, "%d", &num)
	return num
}

func init() {
	addSessionCmd.Flags().String("title", "", "session 标题")
	addSessionCmd.Flags().String("summary", "", "session 摘要")
	addSessionCmd.Flags().String("task", "", "关联任务目录")
	rootCmd.AddCommand(addSessionCmd)
}
