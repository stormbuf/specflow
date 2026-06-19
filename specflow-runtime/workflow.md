# Specflow Workflow

本文件定义 specflow 工作流契约。插件 `inject-workflow-state.js` 会按当前任务 status 匹配下面的 `[workflow-state:<STATUS>]` 面包屑标签块，并注入到每轮用户消息前。

改面包屑 = 改 markdown，下一条消息即生效，无需重启或重新安装。

## 阶段索引

- Phase 1: 规划  (status=planning)
- Phase 2: 执行  (status=in_progress)
- Phase 3: 收尾  (status=completed)

[workflow-state:no_task]
无活跃任务。先对用户请求分类：
- 简单对话 / 一次性问答 → 直接回答，不建任务
- inline 小任务（改几行代码、小修小补）→ 直接动手，不建任务
- 完整 specflow 任务（新功能、重构、跨多文件改动、需要验收）→ 走任务流程

建任务必须经过 consent gate：先用一两句话说明为什么需要建任务，取得用户明确同意后再 `specflow task create`。不要在用户不知情的情况下创建任务。
[/workflow-state:no_task]

[workflow-state:planning]
当前处于规划阶段。建议流程：
1. 加载 `specflow-brainstorm` skill 进行需求探索与证据搜集
2. 编写 `prd.md`（需求文档）
3. 编写 `implement.md`（执行计划，按行为切片）
4. （可选）编写 `design.md`（架构设计）
5. curate jsonl manifest（`specflow add-context` 为 implement / check agent 挑选上下文文件）
6. `specflow task start` 进入 in_progress

产物必须落到任务目录文件中，不要只停留在对话里。
[/workflow-state:planning]

[workflow-state:in_progress]
当前处于执行阶段。建议流程：
1. 派发 `specflow-implement` sub-agent 按 `implement.md` 逐步实现
2. 派发 `specflow-check` sub-agent 验证实现并跑自修复循环
3. 加载 `specflow-sync-requirements` 同步 prd → requirements
4. spec update（`specflow-update-spec`）
5. commit（由 finish-work 统一处理，sub-agent 不单独 commit）
6. `/specflow:finish-work` 归档

派发 sub-agent 的 prompt 必须以 `Active task: <task-dir>` 开头。
上下文按 jsonl manifest 条目顺序加载，sub-agent 无需主动读取 prd.md / implement.md。
[/workflow-state:in_progress]

[workflow-state:completed]
任务已完成。运行 `/specflow:finish-work` 归档任务并写 journal。
归档后任务目录会被移动到 `.specflow/changes/archive/<YYYY-MM>/`，session 独占也会被释放。
[/workflow-state:completed]

## 活跃任务路由

根据当前任务状态与用户意图选择下一步动作：

- 规划阶段或需求不清 → `specflow-brainstorm`
- in_progress 阶段的实现 / 验收 → 派发 `specflow-implement` / `specflow-check` sub-agent
- 行为规约同步（prd → requirements）→ `specflow-sync-requirements`
- 重复调试同一类 bug → `specflow-break-loop`；spec 更新 → `specflow-update-spec`
- 推进当前任务下一步（不记得流程时）→ `/specflow:continue`
- 归档任务 → `/specflow:finish-work`
