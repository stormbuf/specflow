package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Session 对应 .runtime/sessions/<session-key>.json
type Session struct {
	Platform    string `json:"platform"`
	SessionID   string `json:"session_id"`
	CurrentTask string `json:"current_task"`
	LastSeenAt  string `json:"last_seen_at"`
}

// SessionsDir 返回 .runtime/sessions/ 的完整路径
func SessionsDir(specflowDir string) string {
	return filepath.Join(specflowDir, ".runtime", "sessions")
}

// sessionKey 生成 session-key: <platform>_<sanitized_session_id>
func SessionKey(platform, sessionID string) string {
	safe := sanitizeKey(sessionID)
	return fmt.Sprintf("%s_%s", platform, safe)
}

var sanitizeRe = regexp.MustCompile(`[^A-Za-z0-9._-]`)

func sanitizeKey(s string) string {
	s = sanitizeRe.ReplaceAllString(s, "_")
	if len(s) > 160 {
		s = s[:160]
	}
	return s
}

// sessionPath 返回 session 指针文件路径
func sessionPath(specflowDir, key string) string {
	return filepath.Join(SessionsDir(specflowDir), key+".json")
}

// Load 读取 session 指针
func Load(specflowDir, key string) (*Session, error) {
	path := sessionPath(specflowDir, key)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Save 写入 session 指针
func (s *Session) Save(specflowDir, key string) error {
	dir := SessionsDir(specflowDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建 sessions 目录失败: %w", err)
	}
	s.LastSeenAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	path := sessionPath(specflowDir, key)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SetCurrentTask 设置当前活跃任务
func SetCurrentTask(specflowDir, key, platform, sessionID, taskPath string) error {
	s := &Session{
		Platform:    platform,
		SessionID:   sessionID,
		CurrentTask: taskPath,
	}
	return s.Save(specflowDir, key)
}

// ClearCurrentTask 清除当前活跃任务
func ClearCurrentTask(specflowDir, key string) error {
	path := sessionPath(specflowDir, key)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s.CurrentTask = ""
	return s.Save(specflowDir, key)
}

// ExclusiveError 表示 session 独占冲突
type ExclusiveError struct {
	TaskPath   string
	HolderKey  string
	HolderSessionID string
}

func (e *ExclusiveError) Error() string {
	return fmt.Sprintf("任务 %s 正被 session %s 占用，请先 finish 或 task release", e.TaskPath, e.HolderSessionID)
}

// CheckExclusive 检查 session 独占。
// 若目标任务已被其他 session 指向，返回 *ExclusiveError。
func CheckExclusive(specflowDir, currentKey, taskPath string) error {
	dir := SessionsDir(specflowDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		key := strings.TrimSuffix(entry.Name(), ".json")
		if key == currentKey {
			continue
		}
		s, err := Load(specflowDir, key)
		if err != nil {
			continue
		}
		if s.CurrentTask == taskPath {
			return &ExclusiveError{
				TaskPath:   taskPath,
				HolderKey:  key,
				HolderSessionID: s.SessionID,
			}
		}
	}
	return nil
}

// Release 强制清除指向指定任务的所有 session 指针
func Release(specflowDir, taskPath string) (int, error) {
	dir := SessionsDir(specflowDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		key := strings.TrimSuffix(entry.Name(), ".json")
		s, err := Load(specflowDir, key)
		if err != nil {
			continue
		}
		if s.CurrentTask == taskPath {
			s.CurrentTask = ""
			if err := s.Save(specflowDir, key); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

// GetCurrentTask 从 session 指针读取当前活跃任务
func GetCurrentTask(specflowDir, key string) (string, error) {
	s, err := Load(specflowDir, key)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return s.CurrentTask, nil
}

// StaleSession 表示一个 stale session
type StaleSession struct {
	Key        string
	SessionID  string
	TaskPath   string
	LastSeenAt string
}

// FindStale 查找超过阈值的 stale session 指针
func FindStale(specflowDir string, thresholdHours int) ([]StaleSession, error) {
	dir := SessionsDir(specflowDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	threshold := time.Duration(thresholdHours) * time.Hour
	now := time.Now()
	var stale []StaleSession

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		key := strings.TrimSuffix(entry.Name(), ".json")
		s, err := Load(specflowDir, key)
		if err != nil {
			continue
		}
		if s.CurrentTask == "" {
			continue
		}
		lastSeen, err := time.Parse("2006-01-02T15:04:05Z", s.LastSeenAt)
		if err != nil {
			continue
		}
		if now.Sub(lastSeen) > threshold {
			stale = append(stale, StaleSession{
				Key:        key,
				SessionID:  s.SessionID,
				TaskPath:   s.CurrentTask,
				LastSeenAt: s.LastSeenAt,
			})
		}
	}
	return stale, nil
}
