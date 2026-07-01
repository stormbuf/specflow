package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/stormbuf/specflow/internal/version"
)

// Config 对应 .specflow/config.yaml
type Config struct {
	SpecflowVersion string         `yaml:"specflow_version"`
	VCS             string         `yaml:"vcs"`
	Platform        string         `yaml:"platform"`
	MaxJournalLines int            `yaml:"max_journal_lines"`
	Mem             MemConfig      `yaml:"mem"`
	Session         SessionConfig  `yaml:"session"`
}

type MemConfig struct {
	Enabled  bool     `yaml:"enabled"`
	LogPaths []string `yaml:"log_paths"`
}

type SessionConfig struct {
	StaleThresholdHours int `yaml:"stale_threshold_hours"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		SpecflowVersion: version.Version,
		VCS:             "jj",
		Platform:        "opencode",
		MaxJournalLines: 2000,
		Mem: MemConfig{
			Enabled:  true,
			LogPaths: []string{"~/.opencode/sessions/"},
		},
		Session: SessionConfig{
			StaleThresholdHours: 24,
		},
	}
}

// Load 从 .specflow/config.yaml 读取配置
func Load(specflowDir string) (*Config, error) {
	path := filepath.Join(specflowDir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 config.yaml 失败: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 config.yaml 失败: %w", err)
	}
	if cfg.MaxJournalLines == 0 {
		cfg.MaxJournalLines = 2000
	}
	if cfg.Session.StaleThresholdHours == 0 {
		cfg.Session.StaleThresholdHours = 24
	}
	return &cfg, nil
}

// Save 写入 .specflow/config.yaml
func (c *Config) Save(specflowDir string) error {
	path := filepath.Join(specflowDir, "config.yaml")
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("序列化 config.yaml 失败: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// AgentConfig 对应 agents.yaml 中单个 agent 的配置
type AgentConfig struct {
	Source      string   `yaml:"source"`
	JSONLFile   *string  `yaml:"jsonl_file"`
	AgentFile   string   `yaml:"agent_file"`
	RequireTask bool     `yaml:"require_task"`
	Readonly    bool     `yaml:"readonly"`
	CanWrite    bool     `yaml:"can_write"`
	Constraints []string `yaml:"constraints"`
}

// AgentsConfig 对应 .specflow/agents.yaml
type AgentsConfig struct {
	Agents map[string]AgentConfig `yaml:"agents"`
}

// LoadAgents 从 .specflow/agents.yaml 读取
func LoadAgents(specflowDir string) (*AgentsConfig, error) {
	path := filepath.Join(specflowDir, "agents.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 agents.yaml 失败: %w", err)
	}
	var ac AgentsConfig
	if err := yaml.Unmarshal(data, &ac); err != nil {
		return nil, fmt.Errorf("解析 agents.yaml 失败: %w", err)
	}
	return &ac, nil
}

// Validate 校验 agents.yaml
func (ac *AgentsConfig) Validate(specflowDir string) []error {
	var errs []error
	for name, agent := range ac.Agents {
		if agent.Source == "" {
			errs = append(errs, fmt.Errorf("agents.yaml: %s.source 未指定", name))
			continue
		}
		switch agent.Source {
		case "native", "platform":
			// 不需要 agent_file
		case "custom":
			if agent.AgentFile == "" {
				errs = append(errs, fmt.Errorf("agents.yaml: %s.source=custom 但未指定 agent_file", name))
			} else {
				fullPath := filepath.Join(specflowDir, agent.AgentFile)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					errs = append(errs, fmt.Errorf("agents.yaml: %s.agent_file '%s' 不存在", name, agent.AgentFile))
				}
			}
		default:
			errs = append(errs, fmt.Errorf("agents.yaml: %s.source '%s' 无效（应为 native/platform/custom）", name, agent.Source))
		}
	}
	return errs
}
