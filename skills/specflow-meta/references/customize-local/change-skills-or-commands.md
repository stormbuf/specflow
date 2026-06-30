# 改 Skills / Commands

用户想改 AI 入口点、自动触发规则或显式命令行为时，编辑 `.opencode/skills/` 下的 skill 文件。

编辑前，先分类要碰的 skill：

- **Bundled upstream skill** — `specflow-meta`、`specflow-spec-bootstrap`、`specflow-session-insight` 等。源码在 specflow CLI 仓库的 `skills/` 目录下，通过 `go:embed` 编译进二进制，`specflow init` / `specflow update` 分发到 `.opencode/skills/`。本地修改被 `.fingerprints.json` 追踪，下次 update 会标记为"用户已修改"。
- **项目本地 skill** — `.opencode/skills/` 下其他任何 skill。用户拥有，`specflow update` 不刷新。

## 先读这些文件

1. `.specflow/workflow.md`
2. `.opencode/skills/` 目录
3. 相关 agent 或插件文件
4. `.specflow/spec/` 中是否已有项目规则
5. `.specflow/.fingerprints.json` — 确认要编辑的 skill 是否是 upstream-owned（有条目）还是项目本地（无条目）

## 选择哪种入口类型

| 目标 | 建议 |
| --- | --- |
| AI 应自动知道一个能力 | 加或改 skill。 |
| 团队项目约定 | 优先 `.specflow/spec/` 或项目本地 skill —— 绝不放 bundled skill 目录。 |
| 微调 bundled skill 适配本项目 | 创建一个不同名字的项目本地 skill 覆盖意图，或编辑 `.specflow/spec/`。在 bundled skill 目录内的编辑只能存活到下次 `specflow update`。 |
| 贡献回上游 | 编辑 specflow CLI 仓库 `skills/` 下的源文件，不是部署副本。 |
| 改 specflow 流程语义 | 同步 `.specflow/workflow.md`。 |

## 修改 skill

skill 通常长这样：

```text
<skill-name>/
├── SKILL.md
└── references/
```

`SKILL.md` 应短小，负责触发/路由。长内容放 `references/`，让 AI 按需读取。

frontmatter 的 `description` 应明确何时使用：

```yaml
description: "在定制本项目部署流程和发布检查清单时使用。"
```

不要写"有用的项目 skill"这类模糊描述，会导致错误触发。

### Bundled vs 项目本地

| 方面 | Bundled（`specflow-meta` 等） | 项目本地 |
| --- | --- | --- |
| 事实源 | specflow CLI 仓库 `skills/<name>/` | 用户项目内 |
| 分发 | `go:embed` 编译，`specflow init` / `update` 分发 | 用户创建，不被移动 |
| 指纹追踪 | 每个文件记录在 `.fingerprints.json`；update 时冲突提示 | 不追踪 |
| 本地编辑 | 允许，但下次 update 标记"用户已修改" | 自由编辑 |
| 正确定制方式 | 新建不同名字的项目本地 skill 补充或覆盖 | 直接编辑文件 |

## 加项目本地 skill

用户想记录团队私有定制时，创建项目本地 skill —— 绝不把项目私有内容放进 bundled skill 目录：

```text
.opencode/skills/acme-specflow-deploy/
└── SKILL.md
```

选不与 bundled 集冲突的名字：

- `specflow-meta`
- `specflow-spec-bootstrap`
- `specflow-session-insight`
- `specflow-brainstorm` / `specflow-check` / `specflow-implement` / `specflow-research`
- `specflow-before-dev` / `specflow-break-loop` / `specflow-continue` / `specflow-finish-work`
- `specflow-sync-requirements` / `specflow-update-spec`

重名会导致 `specflow update` 覆盖项目本地副本。常用约定是加项目名前缀：`acme-specflow-deploy`、`acme-specflow-onboarding`。

## 注意

- 不把长期工程约定藏在 skill 里；写进 `.specflow/spec/`。
- 不手编 bundled skill 目录下的文件期望持久化 —— `specflow update` 会覆盖。要么贡献上游，要么加项目本地 skill 补充。
- `specflow update` 报告 bundled skill 文件"用户已修改"冲突后，选"保留"仅当你接受手动维护分歧；否则接受覆盖，将意图重新实现为项目本地 skill。
