# 工作区记忆系统

`.specflow/workspace/` 存储跨会话记忆。目的是让 AI 和人理解不同窗口、不同日期之前发生了什么。

## 目录结构

```text
.specflow/workspace/
├── index.md
└── <developer>/
    ├── index.md
    ├── journal-1.md
    └── journal-2.md
```

| 文件 | 用途 |
| --- | --- |
| `.specflow/.developer` | 当前 developer 身份。 |
| `.specflow/workspace/index.md` | 全局 workspace 概览。 |
| `.specflow/workspace/<developer>/index.md` | 该 developer 的会话索引。 |
| `.specflow/workspace/<developer>/journal-N.md` | 会话 journal。 |

## Developer 身份

`specflow init` 时通过 `--developer <name>` 设置。也可后续修改 `.specflow/.developer`。AI 不应随意更改 developer 身份；如果身份不对，先确认当前使用项目的是谁。

## Journal

`journal-N.md` 记录每次会话完成或部分完成的工作。默认每个 journal 约 2000 行，超限后轮转到 `journal-(N+1).md`。行数上限由 `.specflow/config.yaml` 的 `max_journal_lines` 控制。

journal 通常由 `specflow finish-work` 命令写入，记录会话标题、摘要和 commit 信息。

## `specflow mem` 跨会话检索

`specflow mem` 是本地 CLI，索引用户过往的对话日志（JSONL 文件），支持列出会话、按关键词检索、按阶段过滤、输出上下文片段。

```bash
specflow mem search "<关键词>"                  # 按关键词检索
specflow mem search "<关键词>" --phase brainstorm  # 按阶段过滤
specflow mem context "<关键词>"                 # 输出上下文片段
specflow mem list                               # 列出可检索会话
```

`mem` 中的所有内容都不上传，全部在本地读取。日志路径由 `.specflow/config.yaml` 的 `mem.log_paths` 配置，默认为 `~/.opencode/sessions/`。

`specflow-session-insight` skill 教导 AI 何时应该主动伸手去查 `mem`。判断标准是"一个资深同事会不会问'我们之前不是聊过这个吗？'"——如果会，就是该查的时刻。

## workspace 记忆与任务的关系

| 系统 | 存什么 |
| --- | --- |
| `.specflow/changes/` | 特定任务的需求、设计、研究、状态。 |
| `.specflow/workspace/` | 跨任务、跨会话的工作记录。 |
| `.specflow/spec/` | 作为长期约定保留的工程知识。 |

如果信息只对当前任务有用，放任务目录。如果信息描述当前会话发生了什么，放 workspace journal。如果信息应该在未来每次写代码时遵循，放 spec。

## 本地定制点

| 需求 | 编辑位置 |
| --- | --- |
| 改 journal 最大行数 | `max_journal_lines` in `.specflow/config.yaml`。 |
| 改 mem 日志路径 | `mem.log_paths` in `.specflow/config.yaml`。 |
| 启用/禁用 mem | `mem.enabled` in `.specflow/config.yaml`。 |
| 改 stale session 阈值 | `session.stale_threshold_hours` in `.specflow/config.yaml`。 |

## AI 使用规则

AI 不应把 workspace 当作唯一事实源。恢复任务时，先读当前任务，再用 workspace 补充背景。任务完成后，在 workspace 记录重要的过程笔记；如果涌现了长期规则，更新 spec。
