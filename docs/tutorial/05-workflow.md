# 工作流程

Specflow 把一次完整的变更拆成三个阶段，由 `task.json` 的 `status` 字段驱动：

| 阶段 | status 值 | 做什么 |
|------|----------|--------|
| 规划 | `planning` | 澄清需求，编写 prd.md 和 implement.md，整理 jsonl 上下文清单，激活任务 |
| 执行 | `in_progress` | 派发 implement agent 实现，派发 check agent 验收，自修复循环 |
| 收尾 | `completed` | 同步行为规约，沉淀经验到 spec，提交代码，归档任务，写 journal |

## Phase 1：规划

当用户提出一个需要完整任务的需求时（新功能、重构、跨多文件改动），AI 不会直接动手，而是先经过 **consent gate**——用一两句话说明为什么需要建任务，取得用户同意后才创建任务。

```bash
# AI 判断需要建任务，先问用户
"这个需求涉及多个模块的改动，建议创建一个 specflow 任务来管理。是否创建？"

# 用户同意后
specflow task create --title "添加 CSV 导出功能" --intent "用户需要导出数据为 CSV"
```

任务创建后，AI 加载 `specflow-brainstorm` skill 进行需求探索，产出两个核心文件：

- **prd.md**：需求文档，用 EARS 语法描述行为需求（WHEN/IF/THEN/SHALL），用 Gherkin 场景定义验收标准
- **implement.md**：执行计划，把实现拆分为可独立验证的行为切片

然后 AI 整理 jsonl manifest——声明 implement agent 和 check agent 各自需要哪些上下文文件：

```bash
specflow add-context .specflow/changes/<change-id> specflow-implement ".specflow/spec/backend/coding-style.md" "编码规范"
specflow add-context .specflow/changes/<change-id> specflow-implement ".specflow/changes/<change-id>/prd.md" "需求文档"
```

最后激活任务，进入执行阶段：

```bash
specflow task start
```

## Phase 2：执行

AI 派发 `specflow-implement` 子 agent 按 implement.md 逐步实现。插件 `inject-subagent-context.js` 自动拦截派发，按 jsonl manifest 构建上下文并注入 prompt——子 agent 不需要主动读取任何文件，上下文已在它的 prompt 中。

实现完成后，AI 派发 `specflow-check` 子 agent 验收。check agent 以 prd.md 的验收标准为唯一判据，逐条验证。发现问题时可直接修复（自修复上限 3 轮）。

!!! info "面包屑机制"
    整个执行过程中，插件 `inject-workflow-state.js` 每轮解析 workflow.md，按当前 task status 匹配 `[workflow-state:STATUS]` 面包屑标签块并注入。AI 每轮对话都自动知道"现在在执行阶段，下一步该做什么"——不需要你提醒。

## Phase 3：收尾

验收通过后，进入收尾阶段：

1. **行为规约同步**：加载 `specflow-sync-requirements` skill，将 prd 的 EARS 需求同步到 `spec/requirements/`
2. **经验沉淀**：加载 `specflow-update-spec` skill，把值得复用的编码规范和踩坑经验写入 spec 库
3. **提交**：通过 `/specflow:finish-work` 归档任务并写 session journal

```
/specflow:finish-work
→ 归档任务到 changes/archive/<YYYY-MM>/
→ 写 session journal 条目
→ VCS auto-commit
→ 清 session 独占指针
```
