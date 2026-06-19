---
name: specflow-sync-requirements
description: "specflow 需求同步 skill。验收通过后将 prd.md 的 EARS/Gherkin 行为需求同步到 spec/requirements/，按领域组织并更新索引与 Pre-Development Checklist。"
trigger: "Phase 3.1，验收通过后"
phase: in_progress
---

# specflow-sync-requirements

> 需求同步 skill。将 prd.md 的行为需求派生为 spec/requirements/ 下的领域需求文件。

## 核心职责

将 prd.md 中的 EARS / Gherkin 行为需求同步到 `.specflow/spec/requirements/<domain>.md`，更新 requirements 索引与 Pre-Development Checklist。

## 执行步骤

### 1. 读取 prd.md 中的 EARS / Gherkin 需求

- 读取任务目录下的 prd.md
- 提取所有 EARS 需求句（WHEN / IF / THEN / WHILE / SHALL）
- 提取所有 Gherkin 场景（Feature / Scenario / Given / When / Then）

### 2. 确定领域

- 根据需求涉及的功能模块确定领域（如 auth、order、payment 等）
- 若需求跨多个领域，按领域拆分

### 3. 生成 / 更新 .specflow/spec/requirements/<domain>.md

- 已有领域文件：追加或修订需求条目，保持原有结构
- 无领域文件：新建，使用以下结构：

```markdown
# <Domain> 需求

## 行为需求（EARS）

- WHEN <条件> THE SYSTEM SHALL <行为>
- IF <条件> THEN THE SYSTEM SHALL <行为>

## 验收场景（Gherkin）

\`\`\`gherkin
Feature: <功能>
  Scenario: <场景>
    Given <前置>
    When <动作>
    Then <预期>
\`\`\`
```

- 需求条目原样同步，不改写语义

### 4. 更新 requirements/index.md

- 在 `.specflow/spec/requirements/index.md` 追加或更新领域条目
- 格式：`- [domain.md](domain.md) — 一句话摘要`

### 5. 更新 Pre-Development Checklist

- 若新增了领域需求，在 `.specflow/spec/index.md` 的 Pre-Development Checklist 中补充对应检查项
- 确保后续开发能通过 checklist 发现需要读取的 requirements

## 产物

- `.specflow/spec/requirements/<domain>.md`
- 更新后的 `.specflow/spec/requirements/index.md`
- 更新后的 `.specflow/spec/index.md`（Pre-Development Checklist）

## 约束

- requirements 是 prd 的派生产物，只读，按领域组织
- 需求条目原样同步，不改写语义
- 不修改 prd.md
- 不执行 git 操作
