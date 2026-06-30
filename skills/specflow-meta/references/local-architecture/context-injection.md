# 上下文注入系统

specflow 上下文注入的目标是让 AI 在正确的时机读正确的文件，而不是依赖模型记忆。在用户项目中，注入由 OpenCode 插件（`.opencode/plugins/`）配合 specflow CLI 实现。

## 注入类型

| 类型 | 来源 | 用途 |
| --- | --- | --- |
| session context | `session-start.js` + `specflow get-context` | 新会话首条消息注入开发者、VCS、活跃任务、spec 索引、journal 概况。 |
| workflow context | `inject-workflow-state.js` + `specflow task current` | 每轮用户消息注入当前任务状态对应的面包屑。 |
| subagent context | `inject-subagent-context.js` + `specflow build-context` | sub-agent 派发时注入 jsonl 清单引用的文件内容 + 行为约束。 |

## session-start

`session-start.js` 在每会话首条消息触发（内存去重），exec 调用 `specflow get-context --json` 获取 session 上下文，以 XML 风格标签块注入：

```xml
<specflow-session-context>...</specflow-session-context>
<current-state>开发者 / VCS / 活跃任务</current-state>
<spec-indexes>spec 索引路径列表</spec-indexes>
<journal>最近会话日志</journal>
<ready>上下文已加载提示</ready>
```

如果用户感觉 AI 在新会话中不知道当前任务，先检查 `session-start.js` 插件是否已安装并运行。可通过 `SPECFLOW_HOOKS=0` 环境变量禁用所有插件。

## workflow-state

`inject-workflow-state.js` 每轮用户消息触发，exec 调用 `specflow task current --json` 获取当前任务状态，然后解析 `.specflow/workflow.md` 中的 `[workflow-state:<STATUS>]` 面包屑标签块，匹配当前 status 并注入。

找不到匹配 status 的块时，降级为固定文案 `Refer to workflow.md for current step.`。

用户想改"某状态下 AI 下一步该做什么"，编辑 `.specflow/workflow.md` 中对应的状态块即可，无需改插件代码。

## subagent context

`inject-subagent-context.js` 拦截 task 工具调用（subagent 派发），从 `.specflow/agents.yaml` 查 agent 配置，exec 调用 `specflow build-context <agent-type> --json` 构建上下文，拼装到 sub-agent prompt 前。

最终 prompt 结构：

```text
<!-- specflow-hook-injected -->
# <agent-type> Agent Task

## 上下文
<jsonl manifest 引用的文件内容，按 manifest 顺序排列>

## 行为约束（必须遵守）
<agents.yaml 中该 agent 的 constraints 列表>

---

## 你的任务
<调用方传入的原始 prompt>
```

agent 配置中 `jsonl_file` 为 `null` 的 agent（如 `specflow-research`）不注入文件上下文。配置未声明的 agent 类型不处理（扩展点）。

## #367 规避设计

specflow 插件采用非破坏性注入方案（规避 OpenCode #367）：

- **方案 A（优先）**：向 `output.parts` 头部插入独立 text part，不修改用户原文。
- **方案 B（降级，已注释）**：在用户原文 text part 前 prepend，用明确分隔标记。

`session-start.js` 和 `inject-workflow-state.js` 均使用方案 A。这确保用户原文不被修改，注入内容作为独立 part 存在。

## JSONL 读取规则

`implement.jsonl` / `check.jsonl` 每行一个 JSON 对象：

```jsonl
{"file": ".specflow/spec/backend/index.md", "reason": "后端规则"}
```

读取方（`context.go` 的 `ReadJSONL`）跳过无 `file` 字段的行（如 `_example` seed 行）。配置 JSONL 时只放 spec/research 文件，不预注册即将修改的代码文件。

`type: "directory"` 条目会读取该目录下所有 `.md` 文件。

## 活跃任务与 context key

活跃任务状态位于 `.specflow/.runtime/sessions/`，按 session 隔离。session key 格式为 `<platform>_<sanitized_session_id>`。插件从 hook 事件中获取 sessionID，通过 CLI 解析 context key。

如果 shell 命令看不到相同的 context key，`specflow task current` 可能报告无活跃任务。此时检查平台是否将 session 身份传递给 shell，而非手写全局 current-task 文件。

## 本地定制点

| 需求 | 编辑位置 |
| --- | --- |
| 改 session-start 注入内容 | `.opencode/plugins/session-start.js` 中的 `buildSessionContextBlock`。 |
| 改每轮面包屑规则 | `.specflow/workflow.md` 中的 `[workflow-state:STATUS]` 块。 |
| 改 sub-agent 上下文加载 | `.opencode/plugins/inject-subagent-context.js` 和 `.specflow/agents.yaml`。 |
| 改 jsonl 验证/展示 | `context.go` 中的 `ReadJSONL` / `BuildContext`（需改 CLI 源码）。 |
| 改活跃任务解析 | `session.go` 中的 session key 生成和指针读写（需改 CLI 源码）。 |

修改上下文注入时验证两件事：新会话能看到正确的任务，sub-agent 能看到正确的任务产物/spec/research。
