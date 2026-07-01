---
name: specflow-continue
description: "specflow 工作流推进命令。手动调用 /specflow:continue 时激活，读取当前任务状态与 workflow.md 面包屑，判断当前 phase/step 并推进到下一步。"
type: command
trigger: "手动调用 /specflow:continue"
phase: any
---

# specflow-continue

> 工作流推进命令。当不确定当前进度或下一步时调用，自动判断并推进工作流。

## 核心职责

读取当前任务状态与 workflow.md 对应面包屑，判断当前 phase / step，推进到下一步。

## 执行步骤

### 1. 读取 task.json status

- 读取当前活跃任务的 `task.json`，获取 `status` 字段
- status 取值：`planning` / `in_progress` / `completed`
- 若无活跃任务，按以下方式处理：

```bash
specflow task current          # 检查当前任务
specflow task start            # 启动已有任务（指定 task-dir 或使用当前指针）
specflow task create --title "..." --intent "..."  # 或创建新任务
```

### 2. 读取 workflow.md 对应面包屑

- 根据 status 匹配 workflow.md 中的 `[workflow-state:<STATUS>]` 面包屑块
- 读取该状态下的建议流程

### 3. AI 判断当前 phase / step

结合 task.json status 与任务目录中已有产物，判断当前所处步骤：

```text
IF status = planning:
  检查 prd.md 是否存在
  检查 implement.md 是否存在
  检查 design.md 是否存在（可选）
  检查 jsonl manifest 是否已 curate
  判断下一步：补齐缺失产物 / task start

ELSE IF status = in_progress:
  检查 implement.md 中的切片完成状态
  检查验证结果是否已写入
  判断下一步：继续实现 / 派发 check / sync-requirements / update-spec

ELSE IF status = completed:
  提示用户执行 /specflow:finish-work 归档
```

### 4. 推进到下一步

- 根据判断结果，告知用户当前进度与下一步动作
- 若下一步需要激活其他 skill，说明应加载哪个 skill
- 若下一步是命令（如 `task start`、`/specflow:finish-work`），提示用户执行

## 产物

无文件产物。推进工作流状态，指导下一步动作。

## 约束

- 只读取和判断，不直接修改源码
- 不执行 git 操作
- 若任务状态与产物不一致（如 status=in_progress 但无 implement.md），报告异常
