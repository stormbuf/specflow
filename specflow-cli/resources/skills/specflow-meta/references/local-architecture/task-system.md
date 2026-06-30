# 任务系统

specflow 任务系统完全存储在用户项目的 `.specflow/changes/` 下。每个任务是一个目录，包含需求、上下文、研究、状态和关系信息。

## 任务目录结构

```text
.specflow/changes/
├── 2026-04-28-example-task-1/
│   ├── task.json
│   ├── prd.md
│   ├── design.md
│   ├── implement.md
│   ├── implement.jsonl
│   ├── check.jsonl
│   └── research/
└── archive/
    └── 2026-04/
```

| 文件 | 用途 |
| --- | --- |
| `task.json` | 任务元数据：status、creator、assignee、vcs、baseRev、父子关系、relatedFiles、meta。 |
| `prd.md` | 需求、约束和验收标准。轻量任务可只有 PRD。 |
| `design.md` | 复杂任务的技术设计：边界、契约、数据流、兼容性、取舍。 |
| `implement.md` | 复杂任务的执行计划：有序检查清单、验证命令、回滚点。 |
| `implement.jsonl` | implement agent 必须先读的 spec/research 文件清单。 |
| `check.jsonl` | check agent 必须先读的 spec/research 文件清单。 |
| `research/` | 研究产物。复杂发现不应只存在于对话中。 |

## change-id 生成规则

`specflow task create` 自动生成 change-id，格式为 `YYYY-MM-DD-short-slug-N`：

- 日期取当天
- slug 从标题生成（小写、非字母数字转连字符、截断 50 字符）
- N 为同日同 slug 的递增序号（从 1 开始）

例如：`2026-04-28-add-auth-module-1`。

## task.json 字段

| 字段 | 含义 |
| --- | --- |
| `id` / `title` / `description` | 任务标识和描述。 |
| `status` | `planning` / `in_progress` / `completed`。 |
| `intent` | 变更意图简述。 |
| `creator` / `assignee` | 创建者和负责人。 |
| `vcs` / `baseRev` | 版本管理工具和基线版本。 |
| `children` / `parent` | 父子任务关系。 |
| `relatedFiles` | 相关文件列表。 |
| `meta` | 扩展字段，放项目私有数据。 |

## 父子任务

父子关系用于工作结构。父任务将多个独立可验证的交付物归组在同一需求集下，不是依赖调度器，不替代子任务自身的规划产物。

- 父任务拥有：源需求、子任务映射、跨子任务验收标准。
- 子任务可独立走规划/实现/检查/归档。

```bash
specflow task create "<child title>" --parent <parent-id>
```

`children` 是父任务上的历史列表。子任务归档后，父任务保留该子任务名，使 `[2/3 done]` 这类进度仍有意义。

## 活跃任务与 session 独占

活跃任务状态按 session 隔离，存储在：

```text
.specflow/.runtime/sessions/<platform>_<session-id>.json
```

`specflow task start` 将任务路径写入当前 session 的指针文件。不同 AI 窗口可指向不同任务，互不覆盖。

specflow 有 session 独占机制：若目标任务已被其他 session 指向，`task start` 会返回 `ExclusiveError`。归档任务或 `specflow task release` 可释放独占。`specflow doctor` 检测超过 `stale_threshold_hours` 的 stale 指针。

## JSONL 上下文

`implement.jsonl` 和 `check.jsonl` 是 sub-agent 先读的上下文清单。它们不替代 `implement.md`；`implement.md` 是人读的执行计划。

格式：

```jsonl
{"file": ".specflow/spec/backend/index.md", "reason": "后端约定"}
{"file": ".specflow/changes/2026-04-28-x/research/api.md", "reason": "API 研究"}
```

规则：

- 只放 spec 和 research 文件。
- 不放即将被修改的代码文件。
- seed 行（`_example`）无 `file` 字段，仅提示 AI 填入真实条目。

## 常用命令

```bash
specflow task create "<title>"              # 创建任务
specflow task start <task-id>               # 启动任务（planning -> in_progress）
specflow task current                       # 查看当前活跃任务
specflow task add-context <task> implement <file> <reason>  # 追加 jsonl 条目
specflow task list                          # 列出所有任务
specflow task finish                        # 完成任务（in_progress -> completed）
specflow task archive <task-id>             # 归档任务
```

## 本地定制点

| 需求 | 编辑位置 |
| --- | --- |
| 改任务默认字段 | `taskstore.go` 中的 `Create` 和 `Task` 结构体（需改 CLI 源码）。 |
| 改状态语义 | `.specflow/workflow.md` 和面包屑标签块。 |
| 加项目私有字段 | `task.json` 的 `meta` 字段。 |
| 改 jsonl 验证规则 | `context.go` 中的 `ReadJSONL`（需改 CLI 源码）。 |
| 改归档策略 | `taskstore.go` 中的 `Archive`（需改 CLI 源码）。 |

大部分行为定制通过 `.specflow/workflow.md` 和 `.specflow/config.yaml` 即可实现，不需要改 CLI 源码。
