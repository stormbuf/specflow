# Hooks 与插件机制

specflow 的 hooks 是 JS 插件，不是 Python 脚本。三个插件通过 OpenCode 的 hook 事件接入，exec 调用 specflow Go CLI 获取状态，非破坏性注入上下文。

## 三个 JS 插件

| 插件 | Hook 事件 | 用途 |
| --- | --- | --- |
| `session-start.js` | `chat.message`（per-session 去重） | 每会话首条消息注入 specflow 概览（开发者、VCS、活跃任务、spec 索引、journal）。 |
| `inject-workflow-state.js` | `chat.message` | 每轮用户消息注入当前任务状态对应的面包屑标签块。 |
| `inject-subagent-context.js` | `tool.execute.before` | 拦截 task 工具调用（sub-agent 派发），注入 jsonl 上下文 + 行为约束。 |

## 共享 lib

`.opencode/lib/specflow-context.js` 是三个插件的共享工具函数，零外部依赖，仅用 Node.js 内置模块（`fs` / `path` / `child_process` / `util`）。

提供四个核心函数：

| 函数 | 用途 |
| --- | --- |
| `isSpecflowProject(directory)` | 检测 `.specflow/workflow.md` 是否存在。 |
| `isSpecflowSubagent(input)` | 检测消息是否来自 specflow 子 agent（避免递归注入）。 |
| `loadAgentsConfig(directory)` | 读取并解析 `.specflow/agents.yaml`（内置简易 YAML 解析器）。 |
| `exec(cmd, options)` | exec 调用外部命令（主要用于调 specflow CLI）。 |

## session-start.js

每会话首条消息触发，内存去重（同一 sessionID 只注入一次）。exec 调用 `specflow get-context --json` 获取 session 上下文，构造 XML 风格标签块注入。

跳过条件：
- 是 specflow 子 agent 消息
- `SPECFLOW_HOOKS=0` 环境变量
- 非 specflow 项目
- 无 sessionID

## inject-workflow-state.js

每轮用户消息触发。exec 调用 `specflow task current --json` 获取当前任务，解析 `.specflow/workflow.md` 中的 `[workflow-state:<STATUS>]` 面包屑标签块，匹配当前 status 并注入。

正则匹配：`/\[workflow-state:([A-Za-z0-9_-]+)\]\s*\n([\s\S]*?)\n\s*\[\/workflow-state:\1\]/g`

找不到匹配块时降级为 `Refer to workflow.md for current step.`。

改面包屑只需改 `.specflow/workflow.md` 中的状态块，不需要改插件代码。

## inject-subagent-context.js

拦截 `tool.execute.before` 事件，只处理 `tool === "task"` 的调用。从 `output.args.subagent_type` 取 agent 类型名，在 `agents.yaml` 中查配置。未声明则跳过（扩展点）。

exec 调用 `specflow build-context <agent-type> --json` 构建上下文，拼装最终 prompt：

```text
<!-- specflow-hook-injected -->
# <agent-type> Agent Task

## 上下文
<jsonl manifest 引用的文件内容>

## 行为约束（必须遵守）
<constraints 列表>

---

## 你的任务
<原始 prompt>
```

原地修改 `output.args.prompt`。

## #367 规避设计

所有注入采用非破坏性方案（规避 OpenCode #367）：

- **方案 A（优先）**：向 `output.parts` 头部插入独立 text part，不修改用户原文。`metadata.specflow` 标记注入来源。
- **方案 B（降级，已注释）**：在用户原文 text part 前 prepend，用 `---` 分隔标记。

`session-start.js` 和 `inject-workflow-state.js` 使用方案 A。`inject-subagent-context.js` 直接修改 `args.prompt`（因为 sub-agent prompt 是新构造的，不存在原文保留问题）。

## 环境变量

| 变量 | 作用 |
| --- | --- |
| `SPECFLOW_HOOKS=0` | 禁用所有 specflow 插件（调试用） |

## 修改原则

1. **插件读取本地 `.specflow/`，不依赖上游源码路径**。
2. **插件任何异常都不应影响宿主**：所有插件在 catch 中静默处理错误。
3. **exec 调用有 30 秒超时**：避免 CLI 卡死阻塞宿主。
4. **改注入内容优先改 `.specflow/` 文件**：面包屑改 `workflow.md`，约束改 `agents.yaml`，spec 索引改 `spec/`。只有注入逻辑本身需要改时才编辑插件代码。

## 排障路径

用户说"AI 没读 specflow 状态"时：

1. 检查 `.opencode/plugins/` 下三个 JS 文件是否存在。
2. 检查 `SPECFLOW_HOOKS` 是否设为 `0`。
3. 手动运行 `specflow get-context --json` 和 `specflow task current --json` 确认 CLI 输出正常。
4. 检查 `.specflow/.runtime/sessions/` 下是否有 session 指针。
5. 检查 OpenCode 是否将 sessionID 传递给插件。
