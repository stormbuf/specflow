# 本地定制总览

本目录面向在用户项目中工作的本地 AI，该项目已通过 specflow 二进制安装并运行过 `specflow init`。AI 应修改项目内的 `.specflow/` 和 `.opencode/` 目录，不修改 specflow CLI 源码。

## 先确定用户想改什么

| 用户措辞 | 先读 |
| --- | --- |
| "改 specflow 流程 / 阶段 / 下一步提示" | `change-workflow.md` |
| "改任务创建、状态、归档" | `change-task-lifecycle.md` |
| "AI 没读上下文 / 改注入内容" | `change-context-loading.md` |
| "改 implement/check/research agent 行为" | `change-agents.md` |
| "加一个 skill / command" | `change-skills-or-commands.md` |
| "调整项目 spec 结构" | `change-spec-structure.md` |
| "加团队约定和本地笔记" | `add-project-local-conventions.md` |

## 意图→文件对照表

| 定制意图 | 编辑文件 |
| --- | --- |
| 改工作流阶段、面包屑、skill 路由 | `.specflow/workflow.md` |
| 改项目配置（vcs、mem、session） | `.specflow/config.yaml` |
| 改 agent 声明（source、jsonl_file、constraints） | `.specflow/agents.yaml` |
| 改 agent 行为定义 | `.opencode/agents/specflow-*.md` |
| 改 sub-agent 上下文注入 | `.opencode/plugins/inject-subagent-context.js` |
| 改 session-start 注入 | `.opencode/plugins/session-start.js` |
| 改每轮面包屑注入 | `.specflow/workflow.md`（插件是纯解析器） |
| 改 spec 内容或结构 | `.specflow/spec/` |
| 改多 agent worktree 配置 | `.specflow/worktree.yaml` |
| 加项目本地 skill | `.opencode/skills/<新名字>/SKILL.md` |
| 加 custom agent | `.specflow/agents.yaml` + `.specflow/agents/<name>.md` |

## 通用操作顺序

1. **确认平台和目录**：检查 `.opencode/` 是否存在。
2. **确认当前活跃任务**：运行 `specflow task current`。
3. **读本地事实源**：优先 `.specflow/workflow.md`、`.specflow/config.yaml`、相关平台文件。
4. **窄幅修改**：只编辑与用户请求相关的文件。
5. **同步语义**：共享流程变了，检查平台入口是否需要同步；平台入口变了，检查 `workflow.md` 是否仍一致。

## 本地文件优先级

| 层 | 文件 |
| --- | --- |
| 工作流 | `.specflow/workflow.md` |
| 项目配置 | `.specflow/config.yaml` |
| Agent 声明 | `.specflow/agents.yaml` |
| 任务产物 | `.specflow/changes/<task>/` |
| 项目 spec | `.specflow/spec/` |
| 平台集成 | `.opencode/` |
| Agent 定义 | `.opencode/agents/` |
| 插件 | `.opencode/plugins/` |
| Skill | `.opencode/skills/` |

## 默认不做的事

- 不修改 specflow Go 二进制。
- 不假设用户有 specflow CLI 源码仓库。
- 不用默认模板覆盖用户已修改的本地文件。
- 不把团队项目规则放进 `specflow-meta`；项目规则放 `.specflow/spec/` 或本地 skill。

## 何时检查上游源码

只有当用户明确表达以下目标时，才切换到上游源码视角：

- "我想给 specflow 提 PR"
- "我想改 specflow init/update 的生成逻辑"
- "我想 fork specflow"

否则，默认修改用户项目内的 specflow 文件。
