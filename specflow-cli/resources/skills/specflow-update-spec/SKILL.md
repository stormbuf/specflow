---
name: specflow-update-spec
description: "specflow 经验沉淀 skill。完成工作后将值得持久化的编码规范、踩坑经验写入 spec 库，保持 spec 与 index.md 索引同步。"
trigger: "完成工作后，有值得持久化的知识"
phase: in_progress
---

# specflow-update-spec

> 经验沉淀 skill。将任务中的编码规范与踩坑经验固化为团队 spec，而非任务记录。

## 核心职责

回顾本次变更中的编码规范与踩坑经验，确定归属领域，更新或新建 spec 文件，同步 index.md 索引。

## 执行步骤

### 1. 回顾本次变更中的编码规范 / 踩坑

- 回顾 prd.md、implement.md、验证结果与调试过程
- 识别值得沉淀的知识：
  - 编码规范（命名、结构、错误处理等）
  - 踩坑经验（易错点、陷阱、调试教训）
  - 架构约束（模块边界、依赖规则等）
- 过滤掉一次性任务细节，只保留可复用的规范

### 2. 确定写入哪个 spec 领域

- 对照 `.specflow/spec/index.md` 的目录结构
- 判断知识归属：backend / frontend / architecture / 其他
- 若跨多个领域，拆分到各自 spec 文件

### 3. 更新或新建 spec 文件

- 已有 spec 文件：追加或修订相关条目，保持原有结构
- 无对应 spec 文件：在领域目录下新建，文件名使用 kebab-case
- spec 条目使用祈使句，描述"应该做什么"而非"本次做了什么"

### 4. 更新 index.md 索引

- 在 `.specflow/spec/index.md` 对应领域 section 追加一行：
  `- [相对路径](相对路径) — 一句话摘要`
- 若新建了领域目录，在 index.md 中补充目录结构说明

## 产物

- 更新或新建的 spec 文件（`.specflow/spec/<domain>/<name>.md`）
- 更新后的 `.specflow/spec/index.md`

## 约束

- spec 内容是团队规范，不是任务记录
- 只写可复用的规范，不写一次性实现细节
- 不修改 prd.md / implement.md
- 不执行 git 操作
