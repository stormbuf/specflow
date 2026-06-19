package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Fingerprints 对应 .fingerprints.json
type Fingerprints struct {
	SpecflowVersion string            `json:"specflow_version"`
	Files           map[string]string `json:"files"`
}

// Load 读取 .fingerprints.json
func Load(specflowDir string) (*Fingerprints, error) {
	path := filepath.Join(specflowDir, ".fingerprints.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Fingerprints{Files: make(map[string]string)}, nil
		}
		return nil, err
	}
	var fp Fingerprints
	if err := json.Unmarshal(data, &fp); err != nil {
		return nil, err
	}
	if fp.Files == nil {
		fp.Files = make(map[string]string)
	}
	return &fp, nil
}

// Save 写入 .fingerprints.json
func (fp *Fingerprints) Save(specflowDir string) error {
	path := filepath.Join(specflowDir, ".fingerprints.json")
	data, err := json.MarshalIndent(fp, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// HashFile 计算文件内容的 sha256
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(h[:]), nil
}

// HashBytes 计算字节数组的 sha256
func HashBytes(data []byte) string {
	h := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(h[:])
}

// Record 记录管理文件的指纹
func (fp *Fingerprints) Record(projectDir, relPath string) error {
	fullPath := filepath.Join(projectDir, relPath)
	hash, err := HashFile(fullPath)
	if err != nil {
		return err
	}
	fp.Files[relPath] = hash
	return nil
}

// RecordAll 批量记录
func (fp *Fingerprints) RecordAll(projectDir string, relPaths []string) error {
	for _, p := range relPaths {
		if err := fp.Record(projectDir, p); err != nil {
			return fmt.Errorf("记录指纹 %s 失败: %w", p, err)
		}
	}
	return nil
}

// CompareResult 表示三路比对结果
type CompareResult int

const (
	// MatchUserUnchanged 用户未修改（当前 == 旧指纹）→ 可安全覆盖
	MatchUserUnchanged CompareResult = iota
	// MatchCLIUnchanged CLI 未更新（新版本 == 旧指纹）→ 保留用户版本
	MatchCLIUnchanged
	// Conflict 用户修改了且 CLI 也更新了 → 冲突
	Conflict
	// NewFile 新文件（旧指纹中不存在）
	NewFile
)

// ThreeWayCompare 三路比对
// oldHash: 上次记录的指纹
// curHash: 当前磁盘文件的 hash
// newHash: CLI 新版本内容的 hash
func ThreeWayCompare(oldHash, curHash, newHash string) CompareResult {
	if oldHash == "" {
		return NewFile
	}
	if curHash == oldHash {
		return MatchUserUnchanged
	}
	if newHash == oldHash {
		return MatchCLIUnchanged
	}
	return Conflict
}

// CollectManagedFiles 收集 specflow 管理的文件列表
// 包括: .specflow/workflow.md, .opencode/skills/*, .opencode/plugins/*, .opencode/agents/specflow-*
func CollectManagedFiles(projectDir string) ([]string, error) {
	var files []string

	// workflow.md
	managedPaths := []string{
		".specflow/workflow.md",
	}

	for _, p := range managedPaths {
		if _, err := os.Stat(filepath.Join(projectDir, p)); err == nil {
			files = append(files, p)
		}
	}

	// .opencode/skills/**/*
	skillsDir := filepath.Join(projectDir, ".opencode", "skills")
	err := filepath.WalkDir(skillsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(projectDir, path)
		files = append(files, rel)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// .opencode/plugins/*
	pluginsDir := filepath.Join(projectDir, ".opencode", "plugins")
	entries, err := os.ReadDir(pluginsDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".js" {
				files = append(files, filepath.Join(".opencode", "plugins", e.Name()))
			}
		}
	}

	// .opencode/agents/specflow-* (native agents only)
	agentsDir := filepath.Join(projectDir, ".opencode", "agents")
	entries, err = os.ReadDir(agentsDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasPrefix(e.Name(), "specflow-") {
				files = append(files, filepath.Join(".opencode", "agents", e.Name()))
			}
		}
	}

	return files, nil
}
