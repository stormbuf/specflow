# specflow 本地架构总览

`specflow-meta` 面向已运行 `specflow init` 的用户项目。用户机器上通常只有 specflow 二进制和项目内生成的 specflow 文件，不一定有 specflow CLI 源码。

因此，AI 使用本 skill 时，默认定制目标是用户项目内的本地文件：

- `.specflow/`：工作流、配置、任务、spec、workspace 记忆、运行时状态。
- 平台目录：当前仅 `.opencode/`（skills / plugins / lib / agents）。

不要默认引导用户去 fork specflow CLI 仓库。只有当用户明确说想改 specflow 上游源码、发布二进制或贡献 PR 时，才切换到源码视角。

## 本地系统模型

specflow 在用户项目内提供三层：

1. **工作流层**：`.specflow/workflow.md` 定义阶段（planning / in_progress / completed）、skill 路由和面包屑标签块。是项目内工作流唯一事实源。
2. **持久层**：`.specflow/changes/`（任务）、`.specflow/spec/`（工程规范）、`.specflow/workspace/`（会话记忆）存储任务、规范和跨会话记录。
3. **平台集成层**：`.opencode/` 下的 plugins（hooks）、agents、skills、lib 连接 specflow 工作流与 OpenCode AI 工具。

三层均位于用户项目内，AI 可直接读取和修改。

## 核心路径

| 路径 | 用途 |
| --- | --- |
| `.specflow/workflow.md` | 工作流阶段、skill 路由、面包屑标签块。 |
| `.specflow/config.yaml` | 项目配置：vcs、platform、mem、session 等。 |
| `.specflow/agents.yaml` | agent 声明：source / jsonl_file / constraints。 |
| `.specflow/spec/` | 项目工程规范和思维指引。 |
| `.specflow/changes/` | 每个任务的 PRD、执行计划、jsonl 上下文。 |
| `.specflow/workspace/` | 按 developer 分目录的会话 journal。 |
| `.specflow/.runtime/` | session 级运行时状态（活跃任务指针）。 |
| `.specflow/.fingerprints.json` | 管理文件指纹，`specflow update` 三路比对依据。 |
| `.opencode/plugins/` | 3 个 JS 插件（session-start / inject-workflow-state / inject-subagent-context）。 |
| `.opencode/agents/` | native agent 定义（specflow-implement / specflow-check / specflow-research）。 |
| `.opencode/skills/` | bundled skill 副本（随 CLI 分发，update 会刷新）。 |
| `.opencode/lib/` | 插件共享工具函数（specflow-context.js）。 |

## specflow 与 Trellis 的关键差异

| 维度 | Trellis | specflow |
| --- | --- | --- |
| CLI 实现 | Python scripts | Go 单二进制 |
| 资源分发 | npm 包 + `node_modules` | `go:embed` 编译进二进制 |
| 平台支持 | 15+ 平台 | 仅 OpenCode |
| 管理文件追踪 | `.template-hashes.json` | `.fingerprints.json`（三路比对） |
| 多 agent 协作 | `trellis channel` 运行时 | 无此机制 |
| 任务目录 | `.trellis/tasks/` | `.specflow/changes/` |
| 命令调用 | `python3 ./.trellis/scripts/task.py` | `specflow task` |

## AI 定制原则

1. **先找本地事实源**：不要凭记忆编辑。先读 `.specflow/workflow.md`、`.specflow/config.yaml`、相关平台目录和任务文件。
2. **改用户项目，不改二进制**：修改项目内的生成文件，不修改 Go 二进制或 embed 资源。
3. **保持平台文件与 `.specflow/` 对齐**：工作流路由变了，检查平台 skill / agent 文件是否仍描述相同流程。
4. **项目规则放 `.specflow/spec/` 或本地 skill**：不要把团队约定放进 `specflow-meta`。
5. **保留用户修改**：文件已被本地修改时，从当前内容出发，不用默认模板覆盖。

## 如何使用本目录

- 想了解 init 后哪些文件存在，读 `generated-files.md`。
- 想改阶段、路由或下一步动作，读 `workflow.md`。
- 想改任务模型、jsonl 上下文或活跃任务行为，读 `task-system.md`。
- 想改编码规范注入，读 `spec-system.md`。
- 想了解 journal 和跨会话记忆，读 `workspace-memory.md`。
- 想改 hooks 或 sub-agent 上下文加载，读 `context-injection.md`。
