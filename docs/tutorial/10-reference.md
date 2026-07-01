# 命令速查

## 全局标志

所有命令都支持以下持久标志：

| 标志 | 说明 |
|------|------|
| `--json` / `-j` | 输出 JSON 格式 |
| `--verbose` / `-v` | 详细日志输出 |

## CLI 命令

### 初始化

| 命令 | 用途 |
|------|------|
| `specflow init` | 初始化项目（安装三层结构） |

init 标志：`--user` / `-u`（必填）、`--opencode`、`--pi`、`--platform`、`--vcs`、`--force`、`--no-spec`、`--all-spec`、`--with-spec`

### 任务管理

| 命令 | 用途 |
|------|------|
| `specflow task create` | 创建任务 |
| `specflow task start` | 激活任务（检查 session 独占） |
| `specflow task finish` | 释放 session 独占（不改 status） |
| `specflow task archive` | 归档任务（status → completed + 移动 + auto-commit） |
| `specflow task current` | 查看当前活跃任务 |
| `specflow task list` | 列出任务 |
| `specflow task release` | 强制释放 stale session 指针 |
| `specflow task add-subtask` | 链接父子任务 |
| `specflow task remove-subtask` | 解除父子任务关联 |

task create 标志：`--title`、`--description`、`--intent`、`--parent`

task add-subtask 参数：`<parent-dir> <child-dir>`（2 个位置参数）

task remove-subtask 参数：`<parent-dir> <child-dir>`（2 个位置参数）

### 上下文管理

| 命令 | 用途 |
|------|------|
| `specflow get-context` | 聚合 session 上下文 |
| `specflow build-context <agent-name>` | 按 jsonl 构建子 agent 上下文 |
| `specflow add-context <task-dir> <agent-name> <file-path> <reason>` | 追加 jsonl 上下文条目 |

add-context 需要 4 个位置参数：任务目录、agent 名称、文件路径、理由说明。

### Agent 管理

| 命令 | 用途 |
|------|------|
| `specflow sync-agent <name>` | 同步 custom agent 到平台目录 |
| `specflow agents list` | 列出已声明的 agent |

### Spec 模板

| 命令 | 用途 |
|------|------|
| `specflow spec list` | 列出可用 spec 模板分类 |
| `specflow spec install [category...]` | 安装 spec 模板 |

spec install 标志：`--all`（安装所有分类）

### 跨会话检索

| 命令 | 用途 |
|------|------|
| `specflow mem list` | 列出可检索的项目与会话 |
| `specflow mem search <query>` | 按关键词检索历史对话 |
| `specflow mem context <query>` | 检索并输出上下文片段 |

mem search 标志：`--phase`（brainstorm \| implement \| all，默认 all）、`--limit`（默认 10）

mem context 标志：`--phase`（默认 all）

### 会话日志

| 命令 | 用途 |
|------|------|
| `specflow add-session` | 向 journal 追加 session 条目 |

add-session 标志：`--title`、`--summary`、`--task`

### 诊断

| 命令 | 用途 |
|------|------|
| `specflow validate` | 校验配置文件完整性 |
| `specflow doctor` | 诊断项目健康状态（9 项检查） |

doctor 检查项：structure、config、agents、workflow、native_agent_sync、custom_agent_sync、fingerprints、update_candidates、stale_sessions

### 版本管理

| 命令 | 用途 |
|------|------|
| `specflow update` | 同步项目到本地 CLI 版本 |
| `specflow upgrade` | 升级全局 CLI 二进制 |

update 标志：`--force`（跳过冲突询问，全部覆盖）

upgrade 标志：`--channel`（latest \| beta，默认 latest）、`--force`（强制升级）

### Worktree 管理

| 命令 | 用途 |
|------|------|
| `specflow worktree create <name>` | 创建 worktree |
| `specflow worktree list` | 列出 worktree |
| `specflow worktree remove <name>` | 删除 worktree |
| `specflow worktree merge <name>` | 合并 worktree |

worktree create 标志：`--base`（基于哪个分支创建）

worktree remove 标志：`--force`（强制移除）

## Slash 命令

| 命令 | 用途 |
|------|------|
| `/specflow:continue` | 任务内推进下一步（AI 自动判断当前阶段） |
| `/specflow:finish-work` | 归档任务 + 写会话日志 |

## 配置文件

| 文件 | 位置 | 用途 |
|------|------|------|
| `workflow.md` | .specflow/ | 工作流契约（面包屑 + 路由） |
| `config.yaml` | .specflow/ | 项目配置（VCS / journal / mem） |
| `agents.yaml` | .specflow/ | Agent 声明 |
| `.fingerprints.json` | .specflow/ | 文件指纹（update 冲突检测） |
| `.developer` | .specflow/ | 开发者身份（gitignored） |
| `.vcs` | .specflow/ | 版本管理选择（gitignored） |

## 升级流程

两步升级体系：

1. `specflow upgrade` — 升级全局 CLI 二进制（跟随 latest / beta channel）
2. `specflow update` — 同步项目到本地 CLI 版本（含文件指纹三路比对冲突检测）

完整升级 = 先 `upgrade`（CLI）再 `update`（项目）。update 时如果检测到用户修改过管理文件且 CLI 也有更新，会交互询问：覆盖 / 合并 / 放弃。
