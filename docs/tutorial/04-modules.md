# 模块详解

## Go CLI（`specflow` 二进制）

单二进制工具，无运行时依赖。负责所有结构化操作：初始化项目、任务生命周期管理、上下文构建、配置校验、版本管理等。通过 `go:embed` 将 skills、agents、runtime 模板、平台插件和 spec 模板编译进二进制。

CLI 包含 8 个内部包：

| 包 | 职责 |
|----|------|
| `config` | 项目配置读写与校验 |
| `taskstore` | 任务 CRUD、状态机、父子任务 |
| `session` | session 独占指针管理 |
| `vcs` | Git / JJ 适配层 |
| `fingerprint` | 文件指纹比对（update 冲突检测） |
| `context` | jsonl manifest 解析与上下文构建 |
| `installer` | init / update 安装逻辑 |
| `worktree` | 多 agent worktree 管理 |

## Skills（11 个 auto-trigger skill）

Skills 是 specflow 的"大脑"——每个 skill 承载一个特定阶段或能力的 know-how。它们不是你手动调用的函数，而是 AI 根据上下文自动触发的行为指南。

| Skill | 触发时机 | 核心职责 |
|-------|---------|---------|
| `specflow-brainstorm` | 用户同意建任务后 | 澄清需求、搜集证据、起草 prd.md 与 implement.md |
| `specflow-before-dev` | 编码前 | 读取 spec 库，将团队规范纳入工作记忆 |
| `specflow-check` | 实现完成后 | 对照 prd 验收标准逐条验证，自修复上限 3 轮 |
| `specflow-update-spec` | 有值得沉淀的知识时 | 把编码规范和踩坑经验写入 spec 库 |
| `specflow-break-loop` | 重复调试同一类 bug 后 | 5 维根因分析 + 预防机制设计 |
| `specflow-sync-requirements` | 验收通过后 | 将 prd 行为需求同步到 spec/requirements/ |
| `specflow-continue` | 手动调用 `/specflow:continue` | 读取任务状态，判断并推进到下一步 |
| `specflow-finish-work` | 手动调用 `/specflow:finish-work` | 归档任务 + 写 session journal |
| `specflow-spec-bootstrap` | 用户要求从代码库生成 spec | 分析代码库，自动生成项目专属 spec |
| `specflow-session-insight` | 用户引用过往对话 | 调用 specflow mem 检索跨会话历史对话 |
| `specflow-meta` | 用户要求理解或定制 specflow | 理解本地架构，引导定制入口 |

## Agents（3 个 native agent）

Agent 是实际执行任务的子进程。specflow 内置 3 个 native agent，定义在 `agents/` 目录，随 CLI 打包：

| Agent | 角色 | 权限 | 行为约束 |
|-------|------|------|---------|
| `specflow-implement` | 实现 agent | 可读写源码 | 严格遵循 implement.md 步骤；每完成一个行为切片立即更新状态标记；不执行 git 操作；需求不明确时停止报告 |
| `specflow-check` | 验收 agent | 可读写源码 | 以 prd 验收标准为唯一判据；发现问题可自行修复但修复后必须重新验证；不执行 git 操作；自修复上限 3 轮 |
| `specflow-research` | 研究 agent | 只读 | 不修改任何文件；研究结果写入 research/ 目录；不执行任何写操作的 git 命令 |

### Agent 三种来源

在 `agents.yaml` 中声明，支持三种 source：

- **native**：specflow 内置，随 CLI 打包，init/update 时自动同步到平台 agents 目录
- **platform**：复用平台已有 agent（如 OpenCode 自带的 code-reviewer），specflow 只注入上下文，不同步定义
- **custom**：用户自建 agent，定义放在 `.specflow/agents/`，通过 `specflow sync-agent` 同步到平台目录

agents.yaml 中每个 agent 声明包含以下字段：

| 字段 | 说明 |
|------|------|
| `source` | native \| platform \| custom |
| `jsonl_file` | 任务目录下的 jsonl 文件名，null 表示不注入文件上下文 |
| `require_task` | 是否要求活跃任务才能运行 |
| `readonly` | 是否只读 |
| `can_write` | 是否可写源码 |
| `constraints` | 行为约束列表（注入到 subagent prompt 中作为明确行为规则） |

## Plugins（3 个 JS 插件 + 1 个共享 lib）

插件是自动化的"手脚"——它们 hook OpenCode 的事件，在 AI 不知情的情况下自动注入上下文、同步状态、启动会话。每个插件逻辑很薄，只做 hook 拦截 + exec 调用 Go CLI + 结果注入。

| 插件 | Hook 事件 | 触发频率 | 职责 |
|------|----------|---------|------|
| `session-start.js` | `chat.message`（首消息去重） | 每会话一次 | 调用 `specflow get-context`，注入 session 上下文 |
| `inject-workflow-state.js` | `chat.message` | 每轮 | 解析 workflow.md，按 task status 注入面包屑（非破坏性） |
| `inject-subagent-context.js` | `tool.execute.before` | 每次 task 工具调用 | 拦截 subagent 派发，按 jsonl 构建上下文注入 prompt |

!!! warning "非破坏性注入"
    所有注入都采用非破坏性方式——注入内容作为独立的结构化块附加到上下文，**绝不修改用户消息原文**。这是 specflow 从 Trellis issue #367 中吸取的教训：把注入内容写进 user message 会导致 UI 污染，破坏对话可读性。

## Runtime 模板

安装到项目 `.specflow/` 的运行时文件，是 specflow 工作流的"配置"：

- **workflow.md**：工作流契约，定义阶段和面包屑标签块。改流程 = 改 markdown
- **config.yaml**：项目级配置（VCS 选择、journal 行数限制、mem 路径等）
- **agents.yaml**：agent 声明（source / jsonl_file / require_task / readonly / can_write / constraints）
- **spec/index.md**：规范库全局索引 + Pre-Development Checklist

## Spec 模板（9 个分类）

预置的规范库骨架，安装到 `.specflow/spec/`，提供编码规范的起点：

| 分类 | 说明 | 规范文件数 |
|------|------|-----------|
| `guides` | 思维指南（代码复用、跨层、跨平台），预填充直接可用 | 4 |
| `backend` | 后端开发规范骨架（目录结构、数据库、错误处理、质量、日志） | 6 |
| `frontend` | 前端开发规范骨架（目录结构、组件、Hook、状态管理、质量、类型安全） | 7 |
| `architecture` | 架构决策记录骨架（ADR 索引 + 模板） | 2 |
| `testing` | 测试规范骨架（测试约定、mock 策略、集成测试） | 4 |
| `security` | 安全规范骨架（认证授权、输入验证、密钥管理） | 4 |
| `api` | API 设计规范骨架（REST 约定、错误响应、版本管理） | 4 |
| `devops` | DevOps 规范骨架（CI/CD、部署、发布流程） | 4 |
| `git-conventions` | Git 工作流规范骨架（提交约定、分支策略） | 2 |

安装骨架后，可用 `specflow-spec-bootstrap` skill 从真实代码库分析并填充规范内容，无需手写。
