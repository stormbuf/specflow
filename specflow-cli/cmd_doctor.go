package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"specflow/internal/config"
	"specflow/internal/fingerprint"
	"specflow/internal/session"
)

type CheckResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "诊断 specflow 项目健康状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		specflowDir := getSpecflowDir()
		projectDir := getProjectDir()
		var checks []CheckResult
		var fixes []string

		// 1. 结构完整性
		checks = append(checks, checkStructure(specflowDir))

		// 2. config.yaml
		checks = append(checks, checkConfig(specflowDir))

		// 3. agents.yaml
		checks = append(checks, checkAgents(specflowDir))

		// 4. workflow.md
		checks = append(checks, checkWorkflow(specflowDir))

		// 5. native agent 同步状态
		nativeCheck := checkNativeAgentSync(specflowDir, projectDir)
		checks = append(checks, nativeCheck)
		if nativeCheck.Status == "warn" {
			fixes = append(fixes, "运行 specflow init --force 重新同步 native agent")
		}

		// 6. custom agent 同步状态
		customCheck := checkCustomAgentSync(specflowDir, projectDir)
		checks = append(checks, customCheck)
		if customCheck.Status == "warn" {
			fixes = append(fixes, "运行 specflow sync-agent all 同步 custom agent")
		}

		// 7. fingerprints 一致性
		fpCheck := checkFingerprints(specflowDir, projectDir)
		checks = append(checks, fpCheck)

		// 8. update-candidates
		candCheck := checkUpdateCandidates(specflowDir)
		checks = append(checks, candCheck)
		if candCheck.Status == "warn" {
			fixes = append(fixes, fmt.Sprintf("运行 diff 工具对比 %s 中的待合并文件", candCheck.Detail))
		}

		// 9. stale session 指针
		staleCheck := checkStaleSessions(specflowDir)
		checks = append(checks, staleCheck)
		if staleCheck.Status == "warn" {
			fixes = append(fixes, "运行 specflow task release <task-id> 清理 stale 指针")
		}

		if useJSON(cmd) {
			data, _ := json.Marshal(map[string]interface{}{
				"checks": checks,
				"fixes":  fixes,
			})
			fmt.Println(string(data))
		} else {
			for _, c := range checks {
				icon := "✅"
				switch c.Status {
				case "warn":
					icon = "⚠️ "
				case "fail":
					icon = "❌"
				}
				fmt.Printf("%s %s: %s\n", icon, c.Name, c.Detail)
			}
			if len(fixes) > 0 {
				fmt.Println("\n修复建议:")
				for _, f := range fixes {
					fmt.Printf("  - %s\n", f)
				}
			}
		}
		return nil
	},
}

func checkStructure(specflowDir string) CheckResult {
	required := []string{"workflow.md", "config.yaml", "agents.yaml", "spec/index.md"}
	for _, f := range required {
		if _, err := os.Stat(filepath.Join(specflowDir, f)); os.IsNotExist(err) {
			return CheckResult{"structure", "fail", fmt.Sprintf("缺失文件: %s", f)}
		}
	}
	dirs := []string{"changes", "spec", "workspace", ".runtime/sessions"}
	for _, d := range dirs {
		if _, err := os.Stat(filepath.Join(specflowDir, d)); os.IsNotExist(err) {
			return CheckResult{"structure", "warn", fmt.Sprintf("缺失目录: %s/", d)}
		}
	}
	return CheckResult{"structure", "pass", ""}
}

func checkConfig(specflowDir string) CheckResult {
	cfg, err := config.Load(specflowDir)
	if err != nil {
		return CheckResult{"config", "fail", err.Error()}
	}
	if cfg.VCS != "git" && cfg.VCS != "jj" {
		return CheckResult{"config", "fail", fmt.Sprintf("vcs '%s' 无效", cfg.VCS)}
	}
	return CheckResult{"config", "pass", fmt.Sprintf("vcs=%s platform=%s", cfg.VCS, cfg.Platform)}
}

func checkAgents(specflowDir string) CheckResult {
	ac, err := config.LoadAgents(specflowDir)
	if err != nil {
		return CheckResult{"agents", "fail", err.Error()}
	}
	errs := ac.Validate(specflowDir)
	if len(errs) > 0 {
		details := ""
		for _, e := range errs {
			details += e.Error() + "; "
		}
		return CheckResult{"agents", "fail", details}
	}
	return CheckResult{"agents", "pass", fmt.Sprintf("%d agents 声明", len(ac.Agents))}
}

func checkWorkflow(specflowDir string) CheckResult {
	path := filepath.Join(specflowDir, "workflow.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return CheckResult{"workflow", "fail", "workflow.md 不存在"}
	}
	content := string(data)
	if !containsWorkflowStateTags(content) {
		return CheckResult{"workflow", "warn", "未包含 [workflow-state:...] 标签块"}
	}
	return CheckResult{"workflow", "pass", "标签块正常"}
}

func checkNativeAgentSync(specflowDir, projectDir string) CheckResult {
	cfg, _ := config.Load(specflowDir)
	agentsDir := filepath.Join(projectDir, ".opencode", "agents")
	if cfg.Platform == "pi" {
		agentsDir = filepath.Join(projectDir, ".pi", "agents")
	}

	ac, err := config.LoadAgents(specflowDir)
	if err != nil {
		return CheckResult{"native_agent_sync", "fail", err.Error()}
	}

	missing := []string{}
	for name, agent := range ac.Agents {
		if agent.Source != "native" {
			continue
		}
		path := filepath.Join(agentsDir, name+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return CheckResult{"native_agent_sync", "warn", fmt.Sprintf("未同步: %v", missing)}
	}
	return CheckResult{"native_agent_sync", "pass", "全部已同步"}
}

func checkCustomAgentSync(specflowDir, projectDir string) CheckResult {
	cfg, _ := config.Load(specflowDir)
	agentsDir := filepath.Join(projectDir, ".opencode", "agents")
	if cfg.Platform == "pi" {
		agentsDir = filepath.Join(projectDir, ".pi", "agents")
	}

	ac, err := config.LoadAgents(specflowDir)
	if err != nil {
		return CheckResult{"custom_agent_sync", "skip", "无法读取 agents.yaml"}
	}

	missing := []string{}
	for name, agent := range ac.Agents {
		if agent.Source != "custom" {
			continue
		}
		path := filepath.Join(agentsDir, name+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return CheckResult{"custom_agent_sync", "warn", fmt.Sprintf("未同步: %v", missing)}
	}
	count := 0
	for _, a := range ac.Agents {
		if a.Source == "custom" {
			count++
		}
	}
	if count == 0 {
		return CheckResult{"custom_agent_sync", "pass", "无 custom agent"}
	}
	return CheckResult{"custom_agent_sync", "pass", fmt.Sprintf("%d custom agent 已同步", count)}
}

func checkFingerprints(specflowDir, projectDir string) CheckResult {
	fp, err := fingerprint.Load(specflowDir)
	if err != nil {
		return CheckResult{"fingerprints", "warn", "无法读取 .fingerprints.json"}
	}
	if len(fp.Files) == 0 {
		return CheckResult{"fingerprints", "warn", "无指纹记录"}
	}

	mismatched := 0
	for relPath, oldHash := range fp.Files {
		curHash, err := fingerprint.HashFile(filepath.Join(projectDir, relPath))
		if err != nil {
			mismatched++
			continue
		}
		if curHash != oldHash {
			mismatched++
		}
	}
	if mismatched > 0 {
		return CheckResult{"fingerprints", "warn", fmt.Sprintf("%d 个文件被手动修改", mismatched)}
	}
	return CheckResult{"fingerprints", "pass", fmt.Sprintf("%d 个管理文件指纹一致", len(fp.Files))}
}

func checkUpdateCandidates(specflowDir string) CheckResult {
	dir := filepath.Join(specflowDir, ".update-candidates")
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		return CheckResult{"update_candidates", "pass", "无待处理文件"}
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	if count > 0 {
		return CheckResult{"update_candidates", "warn", fmt.Sprintf(".update-candidates/ 有 %d 个待合并文件", count)}
	}
	return CheckResult{"update_candidates", "pass", "无待处理文件"}
}

func checkStaleSessions(specflowDir string) CheckResult {
	cfg, _ := config.Load(specflowDir)
	threshold := 24
	if cfg.Session.StaleThresholdHours > 0 {
		threshold = cfg.Session.StaleThresholdHours
	}
	stale, err := session.FindStale(specflowDir, threshold)
	if err != nil {
		return CheckResult{"stale_sessions", "pass", "无 session 指针"}
	}
	if len(stale) > 0 {
		details := ""
		for _, s := range stale {
			details += fmt.Sprintf("%s (last_seen: %s); ", s.SessionID, s.LastSeenAt)
		}
		return CheckResult{"stale_sessions", "warn", fmt.Sprintf("发现 %d 个 stale 指针: %s", len(stale), details)}
	}
	return CheckResult{"stale_sessions", "pass", "无 stale 指针"}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
