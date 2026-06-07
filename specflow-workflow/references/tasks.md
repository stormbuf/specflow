# Tasks Stage

本文件定义任务拆解阶段。Tasks 将 proposal、spec-delta 和 design 转换为可独立验证的执行清单。

路径锚定：本文件中的 `specflow/`、源码、测试和项目文档路径均相对于 `{PROJECT_ROOT}/`；只有 `{SKILL_DIR}/assets/rules.md` 和模板文件来自 skill 目录。

## 目标

- 每个任务只解决一个行为或契约问题。
- 任务按依赖顺序排列。
- 每个任务有可检查完成条件。
- 不混合无关目标。

## 输入

- `{SKILL_DIR}/assets/rules.md`，每次进入 Tasks 阶段时重新读取
- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/spec-delta.md`，如果存在
- `specflow/changes/<change-id>/design.md`，如果存在

## 输出

- `specflow/changes/<change-id>/tasks.md`

## 变更类型与分组选型

proposal 变更类型直接对应任务分组，无需映射：

| proposal 变更类型 | 模板 |
|---|---|
| `functional` | `{SKILL_DIR}/assets/tasks-functional.md` |
| `nonfunctional` | `{SKILL_DIR}/assets/tasks-nonfunctional.md` |
| `infrastructure` | `{SKILL_DIR}/assets/tasks-infra.md` |
| `lightweight` | `{SKILL_DIR}/assets/tasks-lightweight.md` |

## 拆解规则

先读取 `{SKILL_DIR}/assets/rules.md`，然后按以下流程执行：

```text
读取 proposal.md，确定变更类型和对应的分组模板

IF design.md §9 包含非"无"的架构变更描述:
  读取 architecture.md，对照描述核实变更是否已反映在文档中
  IF 未反映: 告知用户并暂停，请先执行 /specflow-workflow arch
IF design.md §9 包含非"无"的 ADR 候选描述:
  读取 adr/ 目录，对照描述核实对应 ADR 是否存在且状态为 accepted
  IF 缺失或未确认: 告知用户并暂停，请先执行 /specflow-workflow arch
IF 项目缺少 architecture.md，且本次变更需要生成代码或确定实现结构:
  告知用户并暂停，请先执行 /specflow-workflow arch

加载对应的分组模板，按模板骨架展开任务

读取 rules.md 中检查项分类的规则：

IF 规则在当前变更中触发（Given 条件满足）:
  按规则声明的粒度和时机，将对应检查编入模板骨架的适当位置
ELSE:
  在验证任务中标注 [SKIP] 及原因

IF rules.md 与 proposal、spec-delta 或 design 冲突:
  列出冲突点，告知用户，不生成 tasks.md
```

模板骨架定义了任务排列顺序，rules 检查项通过 Given 条件自行判断是否适用、落在骨架的哪个位置。拆解时不硬编码规则类型。

> 多 Agent 策略是执行策略而非检查项。该分类无 Given/When/Then 结构，Tasks 阶段不会在其上编排任务项。Apply 阶段独立读取该分类进行 agent 调度，详见 rules.md §多 Agent 策略说明。

## 默认任务分组

### 功能性变更

1. 规约同步 — 确认主 spec 更新
2. 功能实现 — 每个子项内四步按顺序排列：实现 → 安全扫描 → 测试（单测/属性测试）→ 审查。子项之间独立，不得跨子项按类型归堆（禁止「所有实现→所有测试→所有审查」的排列方式）。功能末尾：集成测试。
3. 质量门禁 — lint、类型检查、全量测试、构建
4. 验证 — 检查遗漏、汇总结果

### 非功能变更

1. 影响面与方案 — 受影响模块/接口/配置；策略及选型确认；涉及性能/安全时采集基线
2. 分步实现 — 按影响范围拆分步骤，每步后跟专项验证
3. 整体审查 — 代码审查
4. 专项验证 — 回归测试 + 性能对比/安全扫描/故障注入
5. 质量门禁 — lint、类型检查、全量测试、构建
6. 验证 — 检查遗漏、汇总专项结果

### 基础设施变更

1. 准备工作 — 选型确认、迁移策略、回滚方案、采集基线
2. 实现 — 适配代码、迁移脚本、配置更新
3. 整体审查 — 代码审查
4. 迁移/切换验证 — 数据一致性、服务连通性、功能可用性
5. 回归 + 基线对比 — 运行已有测试、对比关键指标
6. 质量门禁 — lint、类型检查、全量测试、构建
7. 验证 — 检查遗漏、汇总结果

### 轻量变更

1. 执行 — 完成单一任务
2. 验证 — 确认完成、汇总结果

## 完成条件

- `{PROJECT_ROOT}/specflow/changes/<change-id>/tasks.md` 存在。
- architecture.md 已存在且反映本次变更的架构需求（与 design.md §9 一致）；涉及 ADR 时对应 ADR 已存在且状态为 accepted。
- 已按当前 `{SKILL_DIR}/assets/rules.md` 编排适用的检查、测试、审查和质量门禁任务。
- 任务覆盖 spec-delta 和 design 的关键约束。
- 任务顺序符合依赖关系。
- 验证任务存在。
