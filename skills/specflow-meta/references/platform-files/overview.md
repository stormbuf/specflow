# 平台文件总览

specflow 将同一套本地架构连接到不同 AI 工具。`.specflow/` 存共享运行时；平台目录存适配文件，定义每个 AI 工具如何进入 specflow。

本地 AI 修改 specflow 时，应先区分两类文件：

- **共享文件**：`.specflow/workflow.md`、`.specflow/changes/`、`.specflow/spec/`、`.specflow/agents.yaml`、`.specflow/config.yaml`。
- **平台文件**：当前仅 `.opencode/`（skills / plugins / lib / agents）。

平台文件不存业务状态。它们让 AI 工具读取 specflow 状态、调用 CLI、加载 skill/agent/hook。

## 平台文件分类

| 分类 | 路径 | 用途 |
| --- | --- | --- |
| plugins | `.opencode/plugins/` | 在 session start、用户消息、sub-agent 派发时注入上下文。 |
| lib | `.opencode/lib/` | 插件共享工具函数（零外部依赖）。 |
| agents | `.opencode/agents/` | 定义 `specflow-implement` / `specflow-check` / `specflow-research`。 |
| skills | `.opencode/skills/` | bundled skill 副本 + 项目本地 skill。 |

## specflow 与 Trellis 平台模型的差异

| 维度 | Trellis | specflow |
| --- | --- | --- |
| 平台数量 | 15+ | 仅 OpenCode |
| 平台入口 | settings/config + hooks + agents + skills + commands | plugins + agents + skills + lib |
| 上下文注入 | Python scripts + platform hooks | JS 插件 + exec 调用 Go CLI |
| 多平台扩展 | `trellis init --<platform>` | 预留 `platforms/<platform>/` 目录结构 |

## OpenCode 集成模式

specflow 当前使用 hook/extension driven 模式：

- `session-start.js`：每会话首条消息注入 specflow 概览。
- `inject-workflow-state.js`：每轮用户消息注入当前状态面包屑。
- `inject-subagent-context.js`：sub-agent 派发时注入 jsonl 上下文 + 行为约束。

三个插件均通过 exec 调用 specflow Go CLI 获取状态，通过非破坏性注入（向 `output.parts` 头部插入独立 text part）注入上下文。

## 本地修改顺序

用户要求为平台定制行为时，AI 按此顺序检查文件：

1. 读 `.specflow/workflow.md` 确认共享流程。
2. 读 `.opencode/plugins/` 看哪些插件已安装。
3. 读 `.opencode/agents/` 和 `.opencode/skills/`。
4. 修改离用户需求最近的本地文件。
5. 如果改动影响共享流程，同步 `.specflow/workflow.md` 或 `.specflow/spec/`。

不要只改平台文件而忘共享工作流。不要只改 `.specflow/workflow.md` 而忘平台入口可能仍持有旧描述。
