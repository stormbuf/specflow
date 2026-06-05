# Apply Stage

本文件定义执行阶段。Apply 只执行已确认的 tasks.md，不重新发明 proposal、spec-delta 或 design。

路径锚定：本文件中的 `specflow/`、源码、测试和项目命令路径均相对于 `{PROJECT_ROOT}/`；只有 `{SKILL_DIR}/assets/rules.md` 来自 skill 目录。

## 目标

- 按顺序完成 tasks.md。
- 每完成一项立即更新复选框。
- 只修改当前任务需要的文件。
- 执行适用测试、审查和质量门禁。

## 输入

- `{SKILL_DIR}/assets/rules.md`，每次进入 Apply 阶段和每个任务开始前重新读取
- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/tasks.md`
- `specflow/changes/<change-id>/spec-delta.md`，如果存在
- `specflow/changes/<change-id>/design.md`，如果存在
- 相关主 spec、源码和测试

## 执行前检查

```text
IF 当前 change 目录中的 tasks.md 缺失:
  暂停并进入 Tasks 阶段
ELSE IF proposal 的“是否需要规约”为 yes 但当前 change 目录中的 spec-delta.md 缺失:
  暂停并进入 Spec Delta 阶段
ELSE IF proposal 的“是否需要技术方案”为 yes 但当前 change 目录中的 design.md 缺失:
  暂停并进入 Design 阶段
ELSE IF 项目缺少 architecture.md 和有效 ADR，且本次变更需要生成代码或确定实现结构，但 design.md 未标记 architecture_baseline_confirmed = yes:
  暂停并进入 Design 阶段完成架构基线讨论
ELSE:
  阅读 rules.md 和所有相关阶段产物，并执行当前未完成任务
```

## 规则编排

Apply 不缓存规则。每个任务开始前读取 `{SKILL_DIR}/assets/rules.md`，将当前规则转换为本轮执行检查项：

```text
IF rules.md 要求改动前列文件清单、影响面或确认点:
  先输出并确认对应清单，再编辑文件
IF 当前任务是可独立验证子项实现:
  在勾选该子项前执行适用测试和代码审查，或记录 [SKIP] 原因
IF 当前任务是功能质量门禁:
  在该功能所有子项完成后执行项目已定义的适用检查，或记录 [SKIP] 原因
IF rules.md 新增约束导致当前 tasks.md 不完整:
  暂停 Apply，回到 Tasks 阶段补齐任务
IF rules.md 与已确认 design 或 tasks 冲突:
  暂停，列出冲突并请求用户决策
```

## 执行纪律

- 严格按任务顺序推进，除非依赖关系允许并行。
- 不顺便重构、优化或清理无关代码。
- 自己引入的无用 import、变量、函数或文件必须清理。
- 发现需求歧义、设计缺口或测试不可执行时必须记录并阻断。
- 需要偏离设计时必须暂停，说明原因和影响，获得确认后继续。

## 测试策略

```text
IF 代码包含业务逻辑、数据转换、状态判断、边界处理或错误处理:
  添加或更新单元测试
IF 涉及序列化、反序列化、格式转换、不变量或集合操作:
  判断是否需要属性测试
IF 修改对外 API、CLI、事件、配置、模块边界或数据流:
  添加或更新集成测试
IF 修复 bug:
  先提供复现测试或等价复现步骤
```

## 完成条件

- 当前任务的代码和文档改动完成。
- 如果项目缺少 architecture.md 和有效 ADR，且本次变更需要生成代码或确定实现结构，design.md 已标记 architecture_baseline_confirmed = yes。
- 已按当前 `{SKILL_DIR}/assets/rules.md` 完成本轮适用检查，或记录 `[SKIP]` 原因。
- 适用测试已运行；无法运行时记录原因。
- `{PROJECT_ROOT}/specflow/changes/<change-id>/tasks.md` 对应复选框已更新。
