# 改 Agent 行为

用户想改 `specflow-research`、`specflow-implement` 或 `specflow-check` 行为时，编辑项目内的 agent 定义文件和声明文件。

## 先读这些文件

1. `.opencode/agents/` 下的 native agent 定义（`specflow-implement.md` / `specflow-check.md` / `specflow-research.md`）
2. `.specflow/agents.yaml` — agent 声明
3. `.specflow/workflow.md` Phase 2 / research 路由
4. 当前任务 `prd.md`
5. 当前任务 `implement.jsonl` / `check.jsonl`

## 三种 agent source

`.specflow/agents.yaml` 中每个 agent 有 `source` 字段：

| source | 含义 | agent_file | 示例 |
| --- | --- | --- | --- |
| `native` | specflow 内置 agent，定义在 `agents/*.md`，安装到 `.opencode/agents/` | 不需要（自动安装） | `specflow-implement` / `specflow-check` / `specflow-research` |
| `platform` | 平台自带 agent，不安装到项目 | 不需要 | `opencode-builder` |
| `custom` | 用户自建 agent，必须提供 `agent_file` | 必须（`.specflow/agents/` 下的相对路径） | `specflow-oracle` |

## 常见路径

| 平台 | agent 路径 |
| --- | --- |
| OpenCode | `.opencode/agents/specflow-*.md` |
| custom | `.specflow/agents/<name>.md`（用户自建） |

## 常见需求

| 需求 | 改哪个 agent |
| --- | --- |
| research 必须写文件，不只回复对话 | `specflow-research` |
| 实现前必须读某些本地 spec | `specflow-implement` + `implement.jsonl` 配置规则 |
| 检查时必须跑特定命令 | `specflow-check` |
| agent 不得修改某些目录 | 对应 agent 的 write boundary 指令 |
| agent 输出格式固定 | 对应 agent 的产出说明 |
| 加自定义 agent | `.specflow/agents.yaml` 声明 + `.specflow/agents/<name>.md` 定义 |

## 修改原则

1. **保留角色边界**：research 调查并持久化；implement 写实现；check 审查并修复。
2. **不硬编码项目 spec 进 agent**：长期 spec 放 `.specflow/spec/`；agent 负责读取。
3. **让读取顺序明确**：jsonl manifest → spec/research → 原始 prompt。
4. **让写入边界明确**：哪些目录可写、哪些不可写。
5. **constraints 在 agents.yaml 声明**：行为约束通过 `inject-subagent-context.js` 注入，不需要改 agent 定义文件。

## 改 constraints

`agents.yaml` 中每个 agent 的 `constraints` 列表会被 `inject-subagent-context.js` 注入到 sub-agent prompt 的"行为约束"段。改 constraints 只需编辑 `agents.yaml`，不需要改 agent 定义文件或插件代码。

```yaml
specflow-implement:
  source: native
  constraints:
    - "严格遵循 implement.md 中的步骤顺序，不要跳步"
    - "不执行 git commit / push / merge"
    # 在此追加新约束
```

## 加 custom agent

在 `agents.yaml` 声明，并在 `.specflow/agents/` 下创建定义文件：

```yaml
specflow-oracle:
  source: custom
  agent_file: agents/oracle.md
  jsonl_file: oracle.jsonl
  require_task: true
  readonly: true
  can_write: false
  constraints:
    - "仅基于 spec/ 目录给出架构建议"
```

custom agent 的定义文件放在 `.specflow/agents/<name>.md`，不会被 `specflow update` 覆盖。

## 注意

- agent 定义文件中的 prelude（如"你的上下文已由 hook 注入"）不要删除，否则 agent 只凭对话上下文工作，绕过 specflow 核心机制。
- `inject-subagent-context.js` 只处理 `agents.yaml` 中声明的 agent 类型。未声明的 agent 不注入上下文（扩展点）。
- native agent 定义文件（`.opencode/agents/specflow-*.md`）被 `specflow update` 管理，本地修改会被指纹追踪。
