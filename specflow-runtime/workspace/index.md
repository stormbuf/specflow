# Workspace 会话日志索引

> 记录所有开发者在 AI Agent 协作过程中的跨会话工作记录

---

## 概述

本目录用于追踪所有在本项目中与 AI Agent 协作的开发者工作记录。每个开发者拥有独立的子目录，按 journal 文件顺序沉淀每次 session 的工作内容，便于跨会话回溯与上下文衔接。

### 目录结构

```
workspace/
├── index.md              # 本文件 —— 全局索引
└── {developer}/          # 各开发者的独立目录
    ├── index.md          # 个人索引（含 session 历史）
    └── journal-N.md      # journal 文件（顺序编号：1, 2, 3...）
```

---

## 活跃开发者

| Developer | 最后活跃 | Session 数 | 当前文件 |
|-----------|----------|-----------|----------|
| （暂无）  | -        | -         | -        |

---

## 新开发者引导

### 首次加入

运行初始化命令创建开发者身份与目录：

```bash
specflow init -u <name>
```

该命令会完成以下操作：

1. 创建开发者身份文件（`.specflow/.developer`，已 gitignore）
2. 创建对应的 `workspace/{developer}/` 目录
3. 创建个人 index 与初始 journal 文件

### 老开发者回归

1. 查看当前开发者名称：
   ```bash
   cat .specflow/.developer
   ```

2. 查看个人索引：
   ```bash
   cat .specflow/workspace/$(cat .specflow/.developer)/index.md
   ```

---

## 规则说明

### Journal 文件规则

- 每个 journal 文件**最多 2000 行**（行数上限可通过 `config.yaml` 的 `max_journal_lines` 配置）
- 当当前文件达到上限时，轮转创建 `journal-{N+1}.md` 继续记录
- 创建新 journal 文件时需同步更新个人 `index.md`

### Session 记录说明

每个 session 的记录由 `/specflow:finish-work` 命令在任务归档时**自动生成**，无需手动编写。记录内容覆盖本次 session 的关键信息，确保下一次会话可无缝衔接。

---

## Session 记录模板

session 记录采用以下格式（由 `/specflow:finish-work` 自动填充）：

```markdown
## Session {N}: {标题}

**日期**: YYYY-MM-DD
**任务**: {task-name}
**分支**: `{branch-name}`

### Summary

{一句话概述本次 session 的工作内容}

### Branch

- 工作分支：`{branch-name}`
- 基线分支：`{base-branch}`

### Main Changes

- {改动 1}
- {改动 2}

### Git Commits

| Hash | Message |
|------|---------|
| `abc1234` | {commit message} |

### Testing

- [OK] {测试结果}

### Status

[OK] **已完成** / # **进行中** / [P] **受阻**

### Next Steps

- {下一步 1}
- {下一步 2}
```
