# 改上下文加载

上下文加载决定 AI 何时读取工作流、任务、spec、研究、workspace 和 git 状态。用户说"AI 不知道当前任务"、"agent 没读 spec"或"上下文太多/太少"时，读本页。

## 先读这些文件

1. `.specflow/workflow.md`
2. `.opencode/plugins/session-start.js`
3. `.opencode/plugins/inject-workflow-state.js`
4. `.opencode/plugins/inject-subagent-context.js`
5. `.opencode/lib/specflow-context.js`
6. 当前任务的 `implement.jsonl` / `check.jsonl`
7. `.specflow/agents.yaml`

## 上下文来源

| 来源 | 用途 |
| --- | --- |
| `.specflow/workflow.md` | 工作流和下一步动作提示。 |
| `.specflow/changes/<task>/prd.md` | 当前任务需求。 |
| `.specflow/changes/<task>/design.md` | 复杂任务技术设计。 |
| `.specflow/changes/<task>/implement.md` | 复杂任务执行计划。 |
| `.specflow/changes/<task>/implement.jsonl` | 实现前要读的 spec/research 清单。 |
| `.specflow/changes/<task>/check.jsonl` | 检查时要读的 spec/research 清单。 |
| `.specflow/spec/` | 项目 spec。 |
| `.specflow/workspace/` | 会话记录。 |
| git status | 当前工作区改动。 |

## 常见需求与编辑点

| 需求 | 编辑点 |
| --- | --- |
| 新会话注入更多/更少信息 | `session-start.js` 中的 `buildSessionContextBlock`。 |
| 改每轮用户输入的提示 | `[workflow-state:STATUS]` 块 in `.specflow/workflow.md`。插件是纯解析器。 |
| agent 没读 spec | 任务 jsonl + `.specflow/agents.yaml` 中的 `jsonl_file` + `inject-subagent-context.js`。 |
| 活跃任务丢失 | `session.go` 中的 session key 解析 + 平台 session 身份传递。 |
| 改 jsonl 验证规则 | `context.go` 中的 `ReadJSONL`（需改 CLI 源码）。 |
| 改 agent 行为约束 | `.specflow/agents.yaml` 中对应 agent 的 `constraints` 列表。 |

## JSONL 规则

`implement.jsonl` / `check.jsonl` 是关键上下文加载接口：

```jsonl
{"file": ".specflow/spec/backend/index.md", "reason": "后端约定"}
{"file": ".specflow/changes/2026-04-28-x/research/api.md", "reason": "API 研究"}
```

只放 spec/research 文件。不放即将修改的代码文件；agent 在实现过程中自己读代码文件。

`type: "directory"` 条目会读取该目录下所有 `.md` 文件。

## 改 session 上下文

用户想让每个新会话看到更多项目状态时，编辑：

- `.opencode/plugins/session-start.js` 中的 `buildSessionContextBlock` 函数
- 或 `specflow get-context` 的输出逻辑（需改 CLI 源码）

上下文不能无限增长。优先注入索引和路径，让 AI 按需读详细文件。

## 改 sub-agent 上下文

specflow 使用 hook push 模式：`inject-subagent-context.js` 拦截 task 工具调用，从 `agents.yaml` 查配置，exec 调用 `specflow build-context` 构建上下文，拼装到 prompt 前。

确保 agent 最终读到：

1. jsonl manifest 引用的 spec/research 文件内容
2. `agents.yaml` 中声明的行为约束
3. 调用方传入的原始 prompt

`agents.yaml` 中 `jsonl_file` 为 `null` 的 agent 不注入文件上下文。

## 环境变量

`SPECFLOW_HOOKS=0` 可禁用所有 specflow 插件，用于调试。

## 排障顺序

```bash
specflow task current                    # 确认活跃任务
specflow get-context --json              # 检查 session context 输出
specflow build-context specflow-implement --json  # 检查 sub-agent 上下文
```

确认任务和 jsonl 正确后再编辑插件/agent。
