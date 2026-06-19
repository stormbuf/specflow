---
name: specflow-before-dev
description: "specflow 编码前置 skill。在派发 implement sub-agent 或动手写代码前激活，读取 spec 库与 Pre-Development Checklist，将规范纳入工作记忆。"
trigger: "即将开始编码（implement sub-agent 派发前）"
phase: in_progress
---

# specflow-before-dev

> 编码前置 skill。在写代码之前读取相关 spec，确保实现符合团队规范。

## 核心职责

在编写代码前读取 spec 库，按 Pre-Development Checklist 检查，将相关领域的规范纳入工作记忆。本 skill 不产出文件，只做信息加载。

## 执行步骤

### 1. 读取 spec 索引

- 读取 `.specflow/spec/index.md`，了解 spec 库的整体结构与 Pre-Development Checklist
- 确认当前任务 status 为 in_progress（非 planning / completed）

### 2. 按 Pre-Development Checklist 检查

逐项确认：

- [ ] 已读取当前任务的 prd.md 与 implement.md
- [ ] 已读取 implement.jsonl manifest 中声明的所有 spec 文件
- [ ] 已确认当前任务 status 为 in_progress
- [ ] 已确认要改的文件在 relatedFiles 或 jsonl manifest 范围内
- [ ] 不执行 git commit / push / merge（交给 finish-work）
- [ ] 遇到 prd 不明确的点，先停止并报告，不臆断

### 3. 读取相关领域的 spec

- 根据任务涉及的模块（backend / frontend / architecture 等），读取对应 spec 文件
- 若 spec 库中尚无相关领域 spec，记录为"无既有规范"，按通用最佳实践执行
- 将 spec 中与本次改动直接相关的条目提取到工作记忆

### 4. 记录需要注意的规范

- 在工作记忆中保留本次编码需要遵守的关键约束
- 若发现 spec 与 prd.md 存在冲突，停止并向调用方报告

## 产物

无文件产物。spec 内容纳入工作记忆，供后续实现与验收使用。

## 约束

- 编码前必须执行，不可跳过
- 只读 spec 文件，不修改任何文件
- 不执行 git 操作
- spec 与 prd 冲突时停止并报告，不自行裁决
