# 加项目本地约定

用户往往不需要改 specflow 机制，只需要本地 AI 理解团队约定。此时优先用 `.specflow/spec/` 或项目本地 skill，不编辑 `specflow-meta`。

## 内容类型→位置表

| 内容类型 | 位置 |
| --- | --- |
| 代码必须遵循的规则 | `.specflow/spec/<分类>/` |
| 跨层思维方法 | `.specflow/spec/guides/` |
| 项目特定流程的 AI 能力 | `.opencode/skills/` 下的项目本地 skill |
| 一次性任务材料 | `.specflow/changes/<task>/` |
| 会话摘要 | `.specflow/workspace/<developer>/journal-N.md` |

## 创建项目本地 skill

用户想让 AI 知道"本项目如何定制 specflow"时，创建本地 skill：

```text
.opencode/skills/acme-specflow-local/
└── SKILL.md
```

示例：

```md
---
name: acme-specflow-local
description: "本仓库的 specflow 本地定制。改本项目的 specflow 工作流、本地 agent 或团队特定约定时使用。"
---

# Specflow Local

## 本地范围

本 skill 仅记录本仓库的 specflow 定制。

## 自定义工作流规则

- ...

## 本地 agent 改动

- ...
```

选不与 bundled skill 冲突的名字（加项目名前缀）。

## 写入 `.specflow/spec/`

如果内容是编码约定，写进 spec。示例：

```text
.specflow/spec/backend/error-handling.md
.specflow/spec/frontend/components.md
.specflow/spec/guides/cross-layer-thinking-guide.md
```

写完后更新对应 `index.md`，让 AI 能从入口找到新规则。

## 让当前任务使用新约定

写完 spec 后，加到当前任务上下文：

```bash
specflow task add-context <task> implement ".specflow/spec/backend/error-handling.md" "错误处理约定"
specflow task add-context <task> check ".specflow/spec/backend/error-handling.md" "审查错误处理"
```

## 不要把项目私有规则存进 `specflow-meta`

`specflow-meta` 是理解 specflow 架构和本地定制入口的公共 skill。项目私有内容放：

- `.specflow/spec/`
- 项目本地 skill
- 当前任务
- workspace journal

这能防止 specflow 内置 `specflow-meta` 未来更新时覆盖团队约定。
