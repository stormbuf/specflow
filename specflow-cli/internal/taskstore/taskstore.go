package taskstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Task 对应 task.json
type Task struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Status       string   `json:"status"`
	Intent       string   `json:"intent"`
	Creator      string   `json:"creator"`
	Assignee     string   `json:"assignee"`
	CreatedAt    string   `json:"createdAt"`
	CompletedAt  *string  `json:"completedAt"`
	VCS          string   `json:"vcs"`
	BaseRev      *string  `json:"baseRev"`
	Children     []string `json:"children"`
	Parent       *string  `json:"parent"`
	RelatedFiles []string `json:"relatedFiles"`
	Meta         map[string]interface{} `json:"meta"`
}

const (
	StatusPlanning   = "planning"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

// ChangesDir 是任务目录名
const ChangesDir = "changes"

// changesDir 返回 .specflow/changes/ 的完整路径
func ChangesDirPath(specflowDir string) string {
	return filepath.Join(specflowDir, ChangesDir)
}

// TaskDir 返回任务的完整路径
func TaskDir(specflowDir, taskID string) string {
	return filepath.Join(ChangesDirPath(specflowDir), taskID)
}

// Load 读取 task.json
func Load(specflowDir, taskID string) (*Task, error) {
	path := filepath.Join(TaskDir(specflowDir, taskID), "task.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 task.json 失败: %w", err)
	}
	var t Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("解析 task.json 失败: %w", err)
	}
	return &t, nil
}

// LoadByDir 通过目录路径读取 task.json
func LoadByDir(taskDir string) (*Task, error) {
	path := filepath.Join(taskDir, "task.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 task.json 失败: %w", err)
	}
	var t Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("解析 task.json 失败: %w", err)
	}
	return &t, nil
}

// Save 写入 task.json
func (t *Task) Save(taskDir string) error {
	path := filepath.Join(taskDir, "task.json")
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 task.json 失败: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// SaveToSpecflowDir 写入 task.json（通过 specflowDir 和 taskID）
func (t *Task) SaveToSpecflowDir(specflowDir string) error {
	return t.Save(TaskDir(specflowDir, t.ID))
}

// CreateOptions 创建任务的选项
type CreateOptions struct {
	Title       string
	Description string
	Intent      string
	Developer   string
	VCS         string
	BaseRev     string
	Parent      string // 父任务 ID
}

// GenerateChangeID 生成 change-id: YYYY-MM-DD-short-slug-N
func GenerateChangeID(specflowDir, title string) string {
	date := time.Now().Format("2006-01-02")
	slug := slugify(title)
	if slug == "" {
		slug = "task"
	}

	// 查找同日同 slug 的最大序号
	changesDir := ChangesDirPath(specflowDir)
	entries, _ := os.ReadDir(changesDir)
	prefix := fmt.Sprintf("%s-%s-", date, slug)
	maxN := 0
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			// 提取末尾的数字
			suffix := strings.TrimPrefix(name, prefix)
			var n int
			if _, err := fmt.Sscanf(suffix, "%d", &n); err == nil && n > maxN {
				maxN = n
			}
		}
	}
	return fmt.Sprintf("%s-%s-%d", date, slug, maxN+1)
}

var nonAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlnumRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 50 {
		s = s[:50]
	}
	return s
}

// Create 创建新任务
func Create(specflowDir string, opts CreateOptions) (*Task, error) {
	taskID := GenerateChangeID(specflowDir, opts.Title)
	taskDir := TaskDir(specflowDir, taskID)

	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return nil, fmt.Errorf("创建任务目录失败: %w", err)
	}

	now := time.Now().Format("2006-01-02")
	var baseRev *string
	if opts.BaseRev != "" {
		baseRev = &opts.BaseRev
	}

	var parent *string
	if opts.Parent != "" {
		parent = &opts.Parent
	}

	t := &Task{
		ID:           taskID,
		Title:        opts.Title,
		Description:  opts.Description,
		Status:       StatusPlanning,
		Intent:       opts.Intent,
		Creator:      opts.Developer,
		Assignee:     opts.Developer,
		CreatedAt:    now,
		CompletedAt:  nil,
		VCS:          opts.VCS,
		BaseRev:      baseRev,
		Children:     []string{},
		Parent:       parent,
		RelatedFiles: []string{},
		Meta:         map[string]interface{}{},
	}

	if err := t.Save(taskDir); err != nil {
		return nil, err
	}

	// 处理父子关系
	if opts.Parent != "" {
		parentTask, err := Load(specflowDir, opts.Parent)
		if err != nil {
			return t, fmt.Errorf("加载父任务失败: %w", err)
		}
		parentTask.Children = append(parentTask.Children, taskID)
		if err := parentTask.SaveToSpecflowDir(specflowDir); err != nil {
			return t, fmt.Errorf("更新父任务 children 失败: %w", err)
		}
	}

	// seed 空的 jsonl 和模板文件
	seedFiles(taskDir)

	return t, nil
}

func seedFiles(taskDir string) {
	// seed prd.md
	prdPath := filepath.Join(taskDir, "prd.md")
	if _, err := os.Stat(prdPath); os.IsNotExist(err) {
		os.WriteFile(prdPath, []byte("# 需求文档\n\n## 意图\n\nTODO: 简述变更意图\n\n## 需求\n\nTODO: EARS / Gherkin 格式需求\n"), 0644)
	}

	// seed implement.md
	implPath := filepath.Join(taskDir, "implement.md")
	if _, err := os.Stat(implPath); os.IsNotExist(err) {
		os.WriteFile(implPath, []byte("# 执行计划\n\nTODO: TDD 行为切片 / 度量驱动\n"), 0644)
	}

	// seed implement.jsonl
	seedJSONL(filepath.Join(taskDir, "implement.jsonl"))
	// seed check.jsonl
	seedJSONL(filepath.Join(taskDir, "check.jsonl"))
}

const seedExample = `{"_example":"Fill with {\"file\": \"<path>\", \"reason\": \"<why>\"}. Put spec/research files only — no code paths. Delete this line once real entries are added."}` + "\n"

func seedJSONL(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte(seedExample), 0644)
	}
}

// Start 将任务状态从 planning 改为 in_progress
func (t *Task) Start() error {
	if t.Status != StatusPlanning {
		return fmt.Errorf("任务状态为 %s，无法 start（仅 planning 可 start）", t.Status)
	}
	t.Status = StatusInProgress
	return nil
}

// Complete 将任务状态改为 completed
func (t *Task) Complete() {
	t.Status = StatusCompleted
	now := time.Now().Format("2006-01-02")
	t.CompletedAt = &now
}

// Archive 将任务目录移动到 archive/<YYYY-MM>/
func Archive(specflowDir, taskID string) error {
	taskDir := TaskDir(specflowDir, taskID)
	archiveDir := filepath.Join(ChangesDirPath(specflowDir), "archive", time.Now().Format("2006-01"))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("创建归档目录失败: %w", err)
	}
	dest := filepath.Join(archiveDir, taskID)
	return os.Rename(taskDir, dest)
}

// ListAll 列出所有任务（不含归档）
func ListAll(specflowDir string) ([]*Task, error) {
	changesDir := ChangesDirPath(specflowDir)
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, fmt.Errorf("读取 changes/ 目录失败: %w", err)
	}

	var tasks []*Task
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "archive" {
			continue
		}
		t, err := LoadByDir(filepath.Join(changesDir, entry.Name()))
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	// 按创建日期排序
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt > tasks[j].CreatedAt
	})
	return tasks, nil
}

// AddSubtask 关联父子任务
func AddSubtask(specflowDir, parentID, childID string) error {
	parent, err := Load(specflowDir, parentID)
	if err != nil {
		return fmt.Errorf("加载父任务失败: %w", err)
	}
	// 先验证子任务存在
	child, err := Load(specflowDir, childID)
	if err != nil {
		return fmt.Errorf("加载子任务失败: %w", err)
	}
	// 检查是否已存在
	for _, c := range parent.Children {
		if c == childID {
			return fmt.Errorf("子任务 %s 已存在于父任务 %s 中", childID, parentID)
		}
	}
	parent.Children = append(parent.Children, childID)
	if err := parent.SaveToSpecflowDir(specflowDir); err != nil {
		return err
	}
	child.Parent = &parentID
	return child.SaveToSpecflowDir(specflowDir)
}

// RemoveSubtask 解除父子关联
func RemoveSubtask(specflowDir, parentID, childID string) error {
	parent, err := Load(specflowDir, parentID)
	if err != nil {
		return fmt.Errorf("加载父任务失败: %w", err)
	}
	newChildren := []string{}
	for _, c := range parent.Children {
		if c != childID {
			newChildren = append(newChildren, c)
		}
	}
	parent.Children = newChildren
	if err := parent.SaveToSpecflowDir(specflowDir); err != nil {
		return err
	}

	child, err := Load(specflowDir, childID)
	if err != nil {
		return fmt.Errorf("加载子任务失败: %w", err)
	}
	child.Parent = nil
	return child.SaveToSpecflowDir(specflowDir)
}
