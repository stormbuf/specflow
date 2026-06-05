# Tasks Stage

本文件定义任务拆解阶段。Tasks 将 proposal、spec-delta 和 design 转换为可独立验证的执行清单。

路径锚定：本文件中的 `specflow/`、源码、测试和项目文档路径均相对于 `{PROJECT_ROOT}/`；只有 `{SKILL_DIR}/assets/rules.md` 来自 skill 目录。

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

## 任务格式

必须使用复选框：

```markdown
- [ ] 1.1 <任务描述>
```

功能实现任务必须按“子项实现 -> 子项测试 -> 子项审查 -> 功能质量门禁”展开：

```markdown
## N. 功能：<功能名>

- [ ] N.1 实现 <可独立验证子项 A>
- [ ] N.2 测试 <可独立验证子项 A>
- [ ] N.3 审查 <可独立验证子项 A>
- [ ] N.4 实现 <可独立验证子项 B>
- [ ] N.5 测试 <可独立验证子项 B>
- [ ] N.6 审查 <可独立验证子项 B>
- [ ] N.7 执行 <功能名> 质量门禁
```

## 拆解规则

先读取 `{SKILL_DIR}/assets/rules.md`，将其中适用于本次变更的约束编排进任务清单：

```text
IF 项目缺少 architecture.md 和有效 ADR，且本次变更需要生成代码或确定实现结构，但 design.md 未标记 architecture_baseline_confirmed = yes:
  暂停并返回 Design 阶段完成架构基线讨论，不生成 tasks.md
IF rules.md 中约束要求执行前检查、测试、安全扫描、审查或质量门禁:
  按“子项后测试审查、功能后质量门禁”生成对应任务
ELSE IF rules.md 中约束不适用于本次变更:
  在对应验证任务中标注 [SKIP] 和原因
IF rules.md 与 proposal、spec-delta 或 design 冲突:
  暂停并列出冲突点，不生成 tasks.md
```

```text
IF 一个子项包含多个可观察行为:
  拆分
ELSE IF 一个子项混合重构和新功能:
  拆分
ELSE IF 一个子项混合修复和优化:
  拆分
ELSE IF API 签名变更和实现可独立验证:
  拆分
IF 一个功能包含多个可独立验证子项:
  每个子项后紧跟测试任务和审查任务
  功能末尾追加质量门禁任务
```

## 默认任务分组

1. 规约同步：确认或准备主 spec 更新。
2. 功能实现：按功能分组，并在每个可独立验证子项后紧跟适用测试和审查。
3. 功能质量门禁：每个功能的全部子项完成后，安排一次项目已定义的适用检查。
4. 验证：检查任务遗漏、例外说明和检查结果，并记录摘要。
5. 归档：合并 spec-delta 到主 spec。

## 多 Agent 节奏

```text
IF 平台支持多 agent 且任务风险非低:
  调度/代码/测试/审查按子项串行推进
ELSE:
  当前 agent 按相同职责顺序执行，不跳过测试或审查
```

同一子项测试或审查返修最多 3 轮；仍未通过时阻断用户。

## 完成条件

- `{PROJECT_ROOT}/specflow/changes/<change-id>/tasks.md` 存在。
- 如果项目缺少 architecture.md 和有效 ADR，且本次变更需要生成代码或确定实现结构，design.md 已标记 architecture_baseline_confirmed = yes。
- 已按当前 `{SKILL_DIR}/assets/rules.md` 编排适用的检查、测试、审查和质量门禁任务。
- 每个可独立验证子项后紧跟测试和审查任务。
- 每个功能末尾存在质量门禁任务。
- 任务覆盖 spec-delta 和 design 的关键约束。
- 任务顺序符合依赖关系。
- 验证和归档任务存在。
