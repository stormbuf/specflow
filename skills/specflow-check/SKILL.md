---
name: specflow-check
description: "specflow 验收 skill。实现完成后激活，对照 prd.md 验收标准逐条验证，可修复 finding 但禁止 git 操作，自修复上限 3 轮。"
trigger: "实现完成，需要验证"
phase: in_progress
---

# specflow-check

> 验收 skill。以 prd.md 验收标准为唯一判据，执行验证与自修复循环。

## 核心职责

对照 prd.md 验收标准验证实现，发现 finding 时在能力范围内直接修复（自修复循环上限 3 轮），无法修复的报告给调用方。

## 执行步骤

### 1. 读取 prd.md 验收标准

- 读取任务目录下的 prd.md，提取验收标准与 Gherkin 场景
- 确认验收标准清晰可验证；若不明确，停止并向调用方报告

### 2. 逐条验证

- 对每条验收标准执行验证（运行测试 / 检查输出 / 审查代码）
- 记录每条标准的验证结果（通过 / 未通过）
- 未通过项记录为 finding，包含：标准编号、现象、原因分析

### 3. 发现 finding 记录

- 将 finding 写入 implement.md 的验证段
- 每个 finding 标注严重程度（阻断 / 建议）与是否可修复

### 4. 可修复的 finding 直接修复（最多 3 轮）

```text
FOR round = 1 TO 3:
  修复当前未通过的 finding
  重新验证修复后的实现
  IF 全部通过:
    输出验收通过结论，结束
  ELSE:
    继续下一轮
IF 3 轮后仍有未通过项:
  停止，向调用方报告未通过项清单
```

### 5. 无法修复的报告

- 对超出修复能力范围的 finding（如 prd 矛盾、需回到规划阶段），停止并报告
- 报告内容包含：finding 描述、已尝试的修复、建议的下一步

## 产物

验证报告写入 implement.md 的验证段：

```markdown
## 验证结果

- 验收轮次: <N>
- 结论: <通过 / 未通过>
- 已验证项:
  - [x] <验收标准 1>
  - [x] <验收标准 2>
- 未通过项（若有）:
  - <finding 描述>
```

## 约束

- 以 prd.md 验收标准为唯一判据，不引入 prd 之外的期望
- 可修复 finding，但禁止 git commit / push / merge
- 自修复循环上限 3 轮
- 不修改 prd.md 的需求内容；验收标准以 prd.md 现状为准
- prd 验收标准不明确或存在矛盾时，停止并报告
