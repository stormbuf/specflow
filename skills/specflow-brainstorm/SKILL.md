---
name: specflow-brainstorm
description: "specflow 规划阶段 skill。当用户同意创建任务后激活，负责澄清需求、搜集证据、起草 prd.md 与 implement.md，按需产出 design.md。"
trigger: "用户同意创建任务，进入规划阶段（planning）"
phase: planning
---

# specflow-brainstorm

> 规划阶段的需求探索与产物起草 skill。产物必须落到任务目录文件中，不要只停留在对话里。

## 核心职责

澄清需求、搜集证据、起草 prd.md（需求文档）与 implement.md（执行计划），按需产出 design.md（技术设计）。

## Consent Gate

建任务必须经过 consent gate：

```text
IF 用户提出的需求属于完整 specflow 任务（新功能、重构、跨多文件改动、需要验收）:
  先用一两句话说明为什么需要建任务
  询问用户是否同意创建任务
  IF 用户明确同意:
    执行 specflow task create，进入 planning 阶段
  ELSE:
    不创建任务，按用户意图调整
ELSE IF 简单对话 / 一次性问答 / inline 小任务:
  直接处理，不建任务，不激活本 skill
```

不要在用户不知情的情况下创建任务。

## 执行步骤

### 1. 理解用户需求意图

- 复述用户需求的核心意图，确认双方理解一致
- 识别变更类型：新功能 / 重构 / 修复 / 性能优化 / 其他
- 明确范围与非范围，避免 scope creep

### 2. 搜索代码库 / 文档了解现状

- 搜索相关源码、配置、现有 spec 文件
- 阅读已有实现，理解当前架构与约束
- 记录关键发现（现有行为、依赖关系、潜在影响面）
- 若发现需求与现状存在矛盾或缺口，先与用户确认

### 3. 起草 prd.md（需求文档）

- 使用 `{SKILL_DIR}/assets/prd.md` 模板（单任务）或 `prd-parent.md`（父任务）
- 需求描述使用 EARS 语法（WHEN / IF / THEN / WHILE / 至少一句 SHALL）
- 验收标准使用 Gherkin Scenario（Feature / Scenario / Given / When / Then）
- 产物写入任务目录的 `prd.md`

### 4. 编写 implement.md（执行计划）

- 根据任务性质选择模板：
  - TDD 行为切片：`{SKILL_DIR}/assets/implement-tdd.md`
  - 度量驱动：`{SKILL_DIR}/assets/implement-metric.md`
- 将实现拆分为可独立验证的行为切片或步骤
- 每个切片标注状态（`[ ]` 待办 / `[x]` 完成）
- 产物写入任务目录的 `implement.md`

### 5.（可选）编写 design.md（技术设计）

当任务涉及架构决策、多模块协作或非平凡设计时，使用 `{SKILL_DIR}/assets/design.md` 模板编写技术设计。简单任务可跳过。

## 产物

| 产物 | 路径 | 必需 |
|---|---|---|
| prd.md | `<task-dir>/prd.md` | 是 |
| implement.md | `<task-dir>/implement.md` | 是 |
| design.md | `<task-dir>/design.md` | 否 |

## 模板

创建文档时按需使用以下模板：

- 单任务 PRD: `{SKILL_DIR}/assets/prd.md`
- 父任务 PRD: `{SKILL_DIR}/assets/prd-parent.md`
- TDD 执行计划: `{SKILL_DIR}/assets/implement-tdd.md`
- 度量驱动执行计划: `{SKILL_DIR}/assets/implement-metric.md`
- 技术设计: `{SKILL_DIR}/assets/design.md`

## 约束

- 必须经过 consent gate，先询问用户是否建任务
- EARS 关键字使用英文（WHEN / IF / THEN / WHILE / SHALL），至少一句 SHALL
- Gherkin 关键字使用英文（Feature / Scenario / Given / When / Then），步骤正文使用中文
- 产物必须落到任务目录文件中，不要只停留在对话里
- 需求不明确时停止并向用户确认，不臆断
- 不执行 git commit / push / merge
