---
name: specflow-finish-work
description: "specflow 归档命令。手动调用 /specflow:finish-work 时激活，归档当前任务并写 session journal。归档前建议检查工作区是否有未提交改动。"
type: command
trigger: "手动调用 /specflow:finish-work"
phase: completed
---

# specflow-finish-work

> 归档命令。完成任务后调用，将任务归档并记录 session journal。

## 核心职责

归档当前任务（状态置为 completed、移动到 archive/、触发 VCS auto-commit、清 session 指针），并写 session journal 条目。

## 执行步骤

### 1. 检查工作区状态

- 运行 `specflow doctor` 检查项目健康状态
- 提醒用户确认工作区无未提交改动：若有未提交改动，建议先提交后再归档
- `specflow task archive` 只提交归档操作本身（文件移动 + 任务产物），不会自动提交已有的未提交改动
- 确认所有验收标准已通过（若未通过，警告但允许归档）

### 2. 归档任务

```bash
specflow task archive [task-dir]
```

- 若未指定 `<task-dir>`，默认归档当前活跃任务（即 `specflow task current` 所指向的任务）
- 该命令依次执行：status → completed → 移动到 `archive/<YYYY-MM>/` → VCS auto-commit → 清 session 指针
- 归档路径为 `.specflow/changes/archive/<YYYY-MM>/<task-dir>/`

### 3. 写 session journal

```bash
specflow add-session --title "任务标题" --summary "变更摘要（一句话）" --task <task-dir>
```

- `--title`（必需）：session 标题
- `--summary`：变更摘要
- `--task`：关联的任务目录
- 条目追加到 journal 文件中，记录任务 ID、完成时间与变更摘要

### 4. 确认归档结果

- 确认任务已移动到 `archive/<YYYY-MM>/`
- 确认 session 指针已清除（任务独占已释放，可接受新任务）

## 产物

- 归档后的任务目录（`.specflow/changes/archive/<YYYY-MM>/<task-dir>/`）
- session journal 条目
- VCS 提交（归档操作 + 任务产物）

## 约束

- `specflow task archive` 不检查未提交改动，归档前需人工确认工作区干净
- 不修改 prd.md / implement.md 的内容（只移动文件）
- 归档后任务目录不可恢复到活跃区（需重新创建任务）
- session 指针在归档成功后自动清除
