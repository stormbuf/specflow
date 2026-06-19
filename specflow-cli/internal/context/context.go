package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// JSONLEntry 对应 jsonl 中的一行
type JSONLEntry struct {
	File   string `json:"file"`
	Reason string `json:"reason"`
	Type   string `json:"type,omitempty"`
}

// ReadJSONL 读取 jsonl 文件，返回条目列表
func ReadJSONL(path string) ([]JSONLEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []JSONLEntry
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry JSONLEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue // 跳过无法解析的行
		}
		if entry.File == "" {
			continue // 跳过无 file 字段的行（如 _example）
		}
		_ = i
		entries = append(entries, entry)
	}
	return entries, nil
}

// BuildContext 根据 jsonl 条目加载文件内容，拼装上下文
func BuildContext(projectDir, taskDir, jsonlFile string) (string, error) {
	jsonlPath := filepath.Join(taskDir, jsonlFile)
	entries, err := ReadJSONL(jsonlPath)
	if err != nil {
		return "", fmt.Errorf("读取 jsonl 失败: %w", err)
	}

	var parts []string
	for _, entry := range entries {
		if entry.Type == "directory" {
			// 读取目录下所有 .md 文件
			dirPath := filepath.Join(projectDir, entry.File)
			dirEntries, err := os.ReadDir(dirPath)
			if err != nil {
				continue
			}
			for _, de := range dirEntries {
				if de.IsDir() || !strings.HasSuffix(de.Name(), ".md") {
					continue
				}
				filePath := filepath.Join(entry.File, de.Name())
				content := readFileContent(projectDir, filePath)
				if content != "" {
					parts = append(parts, fmt.Sprintf("=== %s ===\n%s", filePath, content))
				}
			}
		} else {
			content := readFileContent(projectDir, entry.File)
			if content != "" {
				label := entry.File
				if entry.Reason != "" {
					label = fmt.Sprintf("%s (%s)", entry.File, entry.Reason)
				}
				parts = append(parts, fmt.Sprintf("=== %s ===\n%s", label, content))
			}
		}
	}
	return strings.Join(parts, "\n\n"), nil
}

func readFileContent(projectDir, relPath string) string {
	fullPath := filepath.Join(projectDir, relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}
	return string(data)
}

// AddContext 向 jsonl 文件追加一行
func AddContext(taskDir, jsonlFile, filePath, reason string) error {
	path := filepath.Join(taskDir, jsonlFile)
	entry := JSONLEntry{
		File:   filePath,
		Reason: reason,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	line := string(data) + "\n"

	// 检查文件是否存在，如果只有 _example 行则追加
	existing, err := os.ReadFile(path)
	if err == nil && strings.Contains(string(existing), "_example") {
		// 替换 _example 行
		lines := strings.Split(string(existing), "\n")
		var newLines []string
		for _, l := range lines {
			if strings.Contains(l, "_example") {
				continue
			}
			if strings.TrimSpace(l) != "" {
				newLines = append(newLines, l)
			}
		}
		newLines = append(newLines, strings.TrimSpace(line))
		return os.WriteFile(path, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
	}

	// 追加
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

// SpecIndex 表示 spec 索引路径
type SpecIndex struct {
	Path string `json:"path"`
}

// FindSpecIndexes 扫描 spec 目录下的 index.md 文件
func FindSpecIndexes(specflowDir string) []SpecIndex {
	var indexes []SpecIndex
	specDir := filepath.Join(specflowDir, "spec")

	// 顶层 index.md
	topIndex := filepath.Join(specDir, "index.md")
	if _, err := os.Stat(topIndex); err == nil {
		indexes = append(indexes, SpecIndex{Path: filepath.Join(".specflow", "spec", "index.md")})
	}

	// 子目录 index.md
	entries, err := os.ReadDir(specDir)
	if err != nil {
		return indexes
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subIndex := filepath.Join(specDir, entry.Name(), "index.md")
		if _, err := os.Stat(subIndex); err == nil {
			indexes = append(indexes, SpecIndex{
				Path: filepath.Join(".specflow", "spec", entry.Name(), "index.md"),
			})
		}
	}
	return indexes
}

// WorkflowOverview 提取 workflow.md 中的阶段索引部分
func WorkflowOverview(specflowDir string) string {
	path := filepath.Join(specflowDir, "workflow.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	content := string(data)

	// 提取 ## 阶段索引 到下一个 ## 之间的内容
	start := strings.Index(content, "## 阶段索引")
	if start == -1 {
		return ""
	}
	rest := content[start:]
	end := strings.Index(rest[3:], "\n## ")
	if end == -1 {
		return rest
	}
	return rest[:end+3]
}
