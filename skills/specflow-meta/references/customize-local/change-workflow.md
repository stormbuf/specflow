# 改工作流

用户想改 specflow 阶段、下一步动作提示、是否建任务、是否用 sub-agent、何时检查/收尾时，先编辑 `.specflow/workflow.md`。

## 先读这些文件

1. `.specflow/workflow.md`
2. 当前平台的入口文件（skills / agents / plugins）
3. 当前任务的 `task.json` 和 `prd.md`

## 常见需求与编辑点

| 需求 | 编辑点 |
| --- | --- |
| 改阶段名或阶段顺序 | `阶段索引` 和对应 Phase 段落。 |
| 改无任务时是否建任务 | `[workflow-state:no_task]` 状态块。 |
| 改规划阶段的下一步 | Phase 1 和 `[workflow-state:planning]`。 |
| 改 in_progress 阶段是否必须用 agent | Phase 2 和 `[workflow-state:in_progress]`。 |
| 改完成后的收尾流程 | Phase 3 和 `[workflow-state:completed]`。 |
| 改用户意图触发哪个 skill | `活跃任务路由` 部分。 |

## 修改步骤

1. 在 `.specflow/workflow.md` 中找到相关段落。
2. 改规则时，保留明确的触发条件和下一步动作。
3. 如果新增或重命名 skill/agent，同步 `.opencode/` 下对应文件。
4. 面包屑改动只需编辑 `[workflow-state:STATUS]` 块。插件是纯解析器——读取块内原文。保持开闭标签的 STATUS 字符串一致（`[workflow-state:foo]…[/workflow-state:foo]`），不匹配的对会被静默丢弃。
5. 让 AI 重新读 `.specflow/workflow.md`，不要沿用旧对话中的规则。

## 示例：放宽任务创建要求

改何时可跳过建任务，通常编辑 `[workflow-state:no_task]`：

```md
[workflow-state:no_task]
无活跃任务。先对用户请求分类：
- 简单对话 / 一次性问答 → 直接回答，不建任务
- inline 小任务（改几行代码） → 直接动手，不建任务
- 完整 specflow 任务 → 走任务流程
[/workflow-state:no_task]
```

如果正式 Phase 1 流程也需要改，同步 Phase 1 描述。

## 示例：改收尾流程

当前 Phase 3 收尾流程为：`specflow-sync-requirements` → spec update → commit（finish-work 统一处理）→ `/specflow:finish-work` 归档。`finish-work` 拒绝在脏工作区运行。

想改收尾顺序或增减步骤，编辑 Phase 3 描述和 `[workflow-state:completed]` 块。

## `/specflow:continue` 路由表

`/specflow:continue` 通过 `task.json.status` 结合任务产物状态决定恢复到哪个阶段步骤。映射固定在命令本身中；fork 添加自定义状态时，必须同时扩展 `workflow.md` 的状态块和此路由表。

| `status` | 产物状态 | 恢复到 |
| --- | --- | --- |
| `planning` | `prd.md` 缺失 | Phase 1（加载 `specflow-brainstorm`） |
| `planning` | 轻量任务，`prd.md` 已完成 | 请求 start 审查，运行 `specflow task start` |
| `planning` | 复杂任务缺 `design.md` 或 `implement.md` | 补全缺失的规划产物 |
| `in_progress` | 对话历史中无实现 | Phase 2（`specflow-implement`） |
| `in_progress` | 实现完成，未跑 `specflow-check` | Phase 2（`specflow-check`） |
| `in_progress` | check 通过 | Phase 3（spec update → commit） |
| `completed` | 任务仍在活跃树 | `/specflow:finish-work` 归档 |

## 注意

`.specflow/workflow.md` 是本地项目工作流，不是不可变模板。用户可按团队习惯调整。编辑后，平台入口文件可能仍持有旧描述，需一并检查。
