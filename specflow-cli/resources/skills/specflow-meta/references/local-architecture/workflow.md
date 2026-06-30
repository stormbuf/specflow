# 工作流系统

`.specflow/workflow.md` 是 specflow 工作流在用户项目内的唯一事实源。AI 不需要 specflow 源码就能理解当前项目如何推进任务；读这一个文件就够了。

## 文件职责

`.specflow/workflow.md` 有三大职责：

1. **定义工作流阶段**：规划（planning）、执行（in_progress）、收尾（completed）。
2. **定义 skill 路由**：用户表达某种意图时，AI 应该用哪个 skill 或 agent。
3. **提供面包屑标签块**：插件按当前任务状态选取对应块注入对话。

## 当前阶段模型

```text
Phase 1: 规划    (status=planning)    -> 明确要构建什么，产出 prd.md 和必要的研究
Phase 2: 执行    (status=in_progress) -> 按 PRD 和 spec 实现，然后 check
Phase 3: 收尾    (status=completed)   -> 最终验证、保留经验、归档
```

每个阶段包含编号步骤，如"加载 specflow-brainstorm skill"。这些编号不是 `task.json` 中的运行时字段，而是供 AI 和人阅读的工作流结构。

## 面包屑标签块

`workflow.md` 底部包含如下状态块：

```text
[workflow-state:no_task]
...
[/workflow-state:no_task]
```

插件 `inject-workflow-state.js` 根据当前任务 status 匹配对应块并注入每轮用户消息。常见状态：

| 状态 | 含义 |
| --- | --- |
| `no_task` | 当前会话无活跃任务。 |
| `planning` | 任务处于需求、研究或上下文配置阶段。 |
| `in_progress` | 任务已进入实现和检查阶段。 |
| `completed` | 任务已完成，等待归档。 |

用户想改"某状态下 AI 下一步该做什么"时，直接编辑对应状态块即可。插件是纯解析器 —— 读取块内原文，不嵌入兜底文本。

## Skill 路由

`workflow.md` 的"活跃任务路由"部分按意图分发：

- 规划阶段或需求不清 → `specflow-brainstorm`
- in_progress 阶段的实现/验收 → 派发 `specflow-implement` / `specflow-check` sub-agent
- 行为规约同步 → `specflow-sync-requirements`
- 重复调试 → `specflow-break-loop`
- spec 更新 → `specflow-update-spec`
- 推进当前任务 → `/specflow:continue`
- 归档任务 → `/specflow:finish-work`

改本地 AI 行为时，先更新 `workflow.md` 中的路由描述，再检查对应的平台 skill / agent 文件是否需要同步。

## 本地修改模式

| 目标 | 编辑点 |
| --- | --- |
| 新增一个阶段 | 更新阶段索引、阶段正文、路由和状态块。 |
| 改任务创建策略 | 更新 `[workflow-state:no_task]` 块和 Phase 1 描述。 |
| 改默认实现/检查路径 | 更新 Phase 2 和 skill 路由。 |
| 改收尾流程 | 更新 Phase 3 和 `[workflow-state:completed]` 块。 |

编辑后，让 AI 重新读 `.specflow/workflow.md`；不要假设旧对话中的流程仍然有效。

## 与平台文件的关系

`workflow.md` 是本地工作流的语义中心，但平台入口文件（skills / agents / plugins）也可能包含流程描述。只改 `workflow.md` 而不看平台文件，可能导致平台入口仍持有旧语言。用户想改"AI 实际做什么"时，同时检查相关平台目录。

## 路由一致性

`/specflow:continue` 通过 `task.json.status` 结合任务目录内的产物状态决定恢复到哪个阶段步骤。如果 fork 添加了自定义状态，必须同时扩展 `workflow.md` 的状态块和 continue 命令的路由表，否则会落入默认分支，用户不会停在预期步骤。
