---
name: specflow-finish-work
description: "specflow 归档命令。手动调用 /specflow:finish-work 时激活，检查工作区干净后归档任务、写 journal、清 session 指针。工作区有未提交改动时拒绝执行。"
type: command
trigger: "手动调用 /specflow:finish-work"
phase: completed
---

# specflow-finish-work

> 归档命令。完成任务后调用，将任务移动到 archive/ 并写 journal。

## 核心职责

检查工作区无未提交改动，将任务状态置为 completed，移动到 archive/，触发 VCS auto-commit，写 journal 条目，清 session 指针。

## 执行步骤

### 1. 检查工作区无未提交改动

- 调用 `specflow` 检查 VCS 状态（`vcs.HasUncommittedChanges`）
- 若存在未提交改动：
  - 拒绝执行，提示用户先提交或通过 sub-agent 完成后再调用
  - 不自动提交未提交改动

### 2. status → completed

- 将任务 `task.json` 的 `status` 字段更新为 `completed`
- 确认所有验收标准已通过（若未通过，警告但允许归档）

### 3. 移动到 archive/

- 将任务目录从 `.specflow/changes/<task-dir>/` 移动到 `.specflow/changes/archive/<YYYY-MM>/`
- YYYY-MM 使用任务完成时的年月

### 4. VCS auto-commit

- 触发 VCS 自动提交，提交信息包含任务 ID 与标题
- 提交范围：归档操作产生的文件移动 + 任务产物

### 5. 写 journal 条目

- 在 journal 中追加一条记录，包含：
  - 任务 ID 与标题
  - 完成时间
  - 变更摘要（一句话）
  - 归档路径
- 若 journal 单文件超过 `max_journal_lines`（默认 2000），轮转到 `journal-(N+1).md`

### 6. 清 session 指针

- 清除当前 session 的活跃任务指针
- 释放 session 独占，使 session 可接受新任务

## 产物

- 归档后的任务目录（`.specflow/changes/archive/<YYYY-MM>/<task-dir>/`）
- journal 条目
- VCS 提交

## 约束

- 工作区有未提交改动时拒绝执行
- 不修改 prd.md / implement.md 的内容（只移动文件）
- 归档后任务目录不可恢复到活跃区（需重新创建任务）
- session 指针必须在归档成功后清除
