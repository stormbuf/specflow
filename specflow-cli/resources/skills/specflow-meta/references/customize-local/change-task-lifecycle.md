# 改任务生命周期

任务生命周期包括创建、启动、上下文配置、完成、归档、父子任务和 session 独占。默认定制目标是 `.specflow/changes/` 和 `.specflow/config.yaml`。

## 先读这些文件

1. `.specflow/workflow.md`
2. `.specflow/config.yaml`
3. `.specflow/agents.yaml`
4. 当前任务的 `.specflow/changes/<task>/task.json`

## 常见需求与编辑点

| 需求 | 编辑点 |
| --- | --- |
| 改任务创建 consent gate 策略 | `.specflow/workflow.md` 的 `[workflow-state:no_task]` 块。 |
| 改 session 独占行为 | `.specflow/config.yaml` 的 `session.stale_threshold_hours`。 |
| 改 journal 行数上限 | `.specflow/config.yaml` 的 `max_journal_lines`。 |
| 加项目私有字段 | `task.json` 的 `meta` 字段。 |
| 改任务状态语义 | `.specflow/workflow.md` 和面包屑标签块。 |
| 改归档目录结构 | `taskstore.go` 中的 `Archive`（需改 CLI 源码）。 |

## session 独占

specflow 有 session 独占机制：`specflow task start` 检查目标任务是否已被其他 session 指向。若是，返回 `ExclusiveError`，提示先 finish 或 release。

```bash
specflow task release <task-id>    # 强制释放指向任务的所有 session 指针
specflow doctor                    # 检测 stale session 指针
```

session 指针存储在 `.specflow/.runtime/sessions/<platform>_<session-id>.json`。不要手编这些文件来"修复"业务状态。

## 改任务字段

如果用户想加项目本地字段，优先放 `task.json` 的 `meta` 下，避免破坏 CLI 对标准字段的假设：

```json
"meta": {
  "linearIssue": "ENG-123",
  "risk": "high"
}
```

如果确实需要改标准字段，需检查 CLI 中所有读取 `task.json` 的代码路径（`taskstore.go`、`context.go` 等）。

## 改活跃任务

活跃任务是 session 级状态，存储在 `.specflow/.runtime/sessions/`。不要回退到全局 `.current-task` 模型。`specflow task create` 创建任务后，需在 AI session 中运行 `specflow task start` 才能设为活跃。

如果用户想改活跃任务行为，编辑：

- `.specflow/workflow.md` 中活跃任务路由描述
- `.specflow/config.yaml` 中 session 相关配置
- 平台插件中 session 身份传递逻辑

## 修改步骤

1. 用 `specflow task current` 确认当前任务。
2. 读当前任务的 `task.json`，确认 status 和字段。
3. 配置需求，先改 `.specflow/config.yaml`。
4. 流程需求，改 `.specflow/workflow.md`。
5. agent 行为需求，改 `.specflow/agents.yaml` 和 `.opencode/agents/` 下的定义。
6. 如果 AI 流程变了，同步 `.specflow/workflow.md`。

## 常用命令

```bash
specflow task create "<title>"                        # 创建任务
specflow task create "<title>" --parent <parent-id>   # 创建子任务
specflow task start <task-id>                         # 启动任务（设为活跃）
specflow task current                                 # 查看当前活跃任务
specflow task list                                    # 列出所有任务
specflow task finish                                  # 完成任务
specflow task archive <task-id>                       # 归档任务
specflow task release <task-id>                       # 释放 session 独占
specflow task add-subtask <parent> <child>            # 关联父子
specflow task remove-subtask <parent> <child>         # 解除父子
```

## 不要做

- 不直接编辑 `.specflow/.runtime/sessions/` 来"修复"业务状态。
- 不把项目私有字段硬编码进 CLI 源码；优先 `meta`。
- 不默认要求用户 fork specflow CLI。
