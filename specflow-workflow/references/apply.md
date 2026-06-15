# Apply Stage

本文件定义执行阶段。Apply 是忠实执行器：只执行已确认的 tasks.md，不重新发明 proposal、spec-delta 或 design，不独立推导调度逻辑。所有 tasks 完成后自动验证并产出 verification.md。

路径锚定：本文件中的 `specflow/`、源码、测试和项目命令路径均相对于 `{PROJECT_ROOT}/`。

## 目标

- 按顺序完成 tasks.md 中的未完成任务。
- 每完成一项立即更新复选框（中断恢复的基础）。
- 只修改当前任务需要的文件。
- 执行适用测试、审查和质量门禁。
- 所有任务完成后自动验证并产出 verification.md。

## 输入

- `specflow/changes/<change-id>/tasks.md` — Apply 推进的核心清单，含勾选框状态、Agent 标注、covers 字段、[SKIP] 标记
- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/spec-delta.md`，如果存在
- `specflow/changes/<change-id>/design.md`，如果存在
- 相关主 spec、源码和测试

## 前置

进入本阶段后，第一件事是读取 tasks.md 并一次性提取执行契约到工作记忆：

```text
读取 tasks.md 全文
提取执行契约到工作记忆：
  - 任务顺序与分组
  - 每个任务项的 agent 标注
  - 子项内串行约束（顺序由各族拓扑定义，已在 tasks.md 落地）
  - 各 section 的启用 / [SKIP] 状态
后续任务执行引用工作记忆，不重读外部规则文件。
```

## 中断恢复

Apply 通过 tasks.md 复选框追踪进度，中断后重新运行可恢复。

```text
读取 tasks.md 所有复选框
IF 存在未完成复选框:
  从第一个未完成任务继续执行（含实现任务和验证任务）
ELSE IF verification.md 缺失:
  所有任务已标记完成但验证记录缺失 → 执行验证步骤并写入
ELSE:
  所有任务已完成。读取 verification.md 报告状态。不重做任何事。
```

验证任务的复选框勾选表示已验证完成。只有当 verification.md 与复选框状态不一致时才修复数据，不无故重做。

## 执行前检查

```text
IF 当前 change 目录中的 tasks.md 缺失:
  暂停并进入 Tasks 阶段
ELSE:
  按 SKILL.md「阶段进入条件」全局门禁做一致性核对（不重复 Tasks 已完成的校验）：
    IF proposal/spec-delta/design/architecture 与 tasks.md 声明的任务族不一致:
      暂停，修正不一致后继续
    ELSE:
      按中断恢复逻辑执行
```

## 执行纪律

### 工作方式

- 严格按 tasks.md 顺序推进，不跳过、不乱序。
- 同一任务最多返修三轮；三轮后仍有问题则阻断，告知用户已尝试的方案和下一步选择。
- 不获确认的情况下不继续执行；不自行假设需求、补充需求。

### Agent 调度

调度依据是 tasks.md 中已编排的 agent 标注和串行约束。Tasks 阶段已将 Agent 角色、串行顺序、降级规则编进 tasks.md；Apply 只执行，不推导。

**执行规则：**

1. 按 tasks.md 顺序遍历未完成任务，读取任务项的 agent 标注。
2. 无依赖子项组可并行委托；组内任务按顺序串行（顺序由各族拓扑定义，已在 tasks.md 落地）。
3. 标注匹配 + 平台支持 → 委托对应 agent；否则编排器自行执行。
4. 测试任务：agent 产出后立即执行验证。失败时按常识分流（测试代码问题→测试 agent 修复；实现 bug→退回代码生成 agent 修复→重跑安全扫描→重跑测试）。修复后重跑，最多三轮，超出阻断。
5. 审查任务：审查发现问题 → 退回代码生成 agent 修复 → 重跑当前子项的安全扫描→测试→审查。最多三轮，超出阻断。
6. 非测试非审查任务（含安全扫描）：验证产出满足 tasks.md 约束，不满足则修正。
7. 每轮完成后更新 tasks.md 复选框；超过三轮上限时按「工作方式」阻断处理。

**降级：** 按 agent 粒度降级——每个任务项对应 agent 可用则委托，不可用则编排器自行执行。不因单个 agent 不可用而全局降级。

### 沟通边界

- 原子修改（<20 行，单文件）→ 事后摘要告知。
- 中等及以上变更（多文件 / 改对外接口 / 改数据结构）→ 先出方案等确认。
- 发现需求歧义、设计缺口或测试不可执行 → 记录并阻断。
- 需要偏离设计 → 暂停，说明原因和影响，获得确认后继续。

### 变更边界

- 不顺便重构、优化或清理无关代码。
- 自己引入的无用 import、变量、函数或文件必须清理。

## 验证

所有 tasks 完成后，Apply 必须执行以下验证步骤并产出 verification.md。

### 遗漏检查

```text
IF tasks.md 存在未完成复选框且没有备注说明例外:
  Result = failed，需要完成缺失项或记录例外
ELSE IF tasks.md 子项 covers 字段未覆盖 spec-delta 全部需求 或 design 全部关键约束:
  Result = failed，需要补齐遗漏项
ELSE IF tasks.md 中验证项缺少结果记录:
  补齐检查结果；无法补齐时记录原因
ELSE:
  汇总验证摘要
```

### 检查命令

```text
IF 已执行适用检查且结果可追溯:
  记录已有结果
ELSE IF 项目存在聚合质量门禁命令:
  运行聚合命令并记录结果
ELSE IF 项目存在明确 lint、类型检查、测试或构建命令:
  运行适用命令并记录结果
ELSE:
  记录未发现可执行命令，不臆造命令
```

不得臆造不存在的 build、test、lint 或 codegen 命令。

### 产出

写入 `{PROJECT_ROOT}/specflow/changes/<change-id>/verification.md`（使用 `{SKILL_DIR}/assets/verification.md` 模板）。

### Result

- `passed`：tasks 无遗漏，所有适用检查通过。
- `failed`：存在未完成任务、遗漏任务、失败检查或未说明例外。
- `partial`：tasks 已完成，但部分检查无法执行或工具链缺失，并已记录风险。

## 完成条件

- 当前任务的代码和文档改动完成。
- architecture/ADR gate 由 SKILL.md 全局门禁统一校验，本文件不重复。
- 已按 tasks.md 启用的 section 完成本轮适用检查，或记录 `[SKIP]` 原因。
- 适用测试已运行；无法运行时记录原因。
- `{PROJECT_ROOT}/specflow/changes/<change-id>/tasks.md` 对应复选框已更新。
- `{PROJECT_ROOT}/specflow/changes/<change-id>/verification.md` 存在，Result 不为空。
- 已按 tasks.md agent 标注执行委托（如适用）；降级时按 agent 粒度处理，不全局降级。
