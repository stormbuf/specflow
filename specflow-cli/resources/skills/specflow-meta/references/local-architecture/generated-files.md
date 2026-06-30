# specflow init 生成的文件

`specflow init` 将 specflow 运行时写入用户项目。后续 `specflow update` 尝试更新 specflow 管理的模板文件，但通过 `.specflow/.fingerprints.json` 三路比对来判断哪些文件已被用户修改。

本页仅描述用户项目内可见、可编辑的文件。

## `.specflow/`

```text
.specflow/
├── workflow.md
├── config.yaml
├── agents.yaml
├── worktree.yaml
├── .developer
├── .fingerprints.json
├── .runtime/
│   └── sessions/
├── spec/
│   └── index.md
├── changes/
│   └── archive/
└── workspace/
    └── index.md
```

| 路径 | 通常可编辑？ | 说明 |
| --- | --- | --- |
| `.specflow/workflow.md` | 是 | 本地工作流文档和 AI 路由规则。 |
| `.specflow/config.yaml` | 是 | 项目配置：vcs、platform、mem、session 等。 |
| `.specflow/agents.yaml` | 是 | agent 声明，可扩展 custom agent。 |
| `.specflow/worktree.yaml` | 是 | 多 agent worktree 配置（复制文件、post_create、pre_merge）。 |
| `.specflow/spec/` | 是 | 项目 spec，供用户和 AI 定期更新。 |
| `.specflow/changes/` | 是 | 任务材料和产物，由任务工作流维护。 |
| `.specflow/workspace/` | 是 | 会话记录，通常由 `specflow finish-work` 写入。 |
| `.specflow/.runtime/` | 否 | 运行时状态，由 CLI 和插件自动写入。 |
| `.specflow/.developer` | 谨慎 | 当前 developer 身份。 |
| `.specflow/.fingerprints.json` | 否 | 指纹记录，不要手写业务规则。 |

## `.opencode/`

```text
.opencode/
├── skills/              # bundled skill 副本（specflow-meta 等）
├── plugins/             # 3 个 JS 插件
│   ├── session-start.js
│   ├── inject-workflow-state.js
│   └── inject-subagent-context.js
├── lib/                 # 插件共享工具
│   └── specflow-context.js
└── agents/              # native agent 定义
    ├── specflow-implement.md
    ├── specflow-check.md
    └── specflow-research.md
```

平台文件不存业务状态，只让 AI 工具读取 specflow 状态、调用 CLI、加载 skill/agent/hook。

## AGENTS.md managed block

`specflow init` 还会在项目根目录的 `AGENTS.md` 中写入 managed block：

```text
<!-- SPECFLOW:START -->
...
<!-- SPECFLOW:END -->
```

该 block 内的内容会被 `specflow update` 覆盖；block 外的内容由用户拥有，不受影响。

## 指纹三路比对机制

`.specflow/.fingerprints.json` 记录上次 specflow 写入模板文件时的内容 hash。`specflow update` 通过三路比对判断：

| 情况 | 比对结果 | update 行为 |
| --- | --- | --- |
| 用户未修改（当前 == 旧指纹） | `MatchUserUnchanged` | 可安全覆盖为新版本 |
| CLI 未更新（新版本 == 旧指纹） | `MatchCLIUnchanged` | 保留用户版本 |
| 用户修改了且 CLI 也更新了 | `Conflict` | 冲突，需用户决定 |
| 旧指纹中不存在 | `NewFile` | 新文件，直接写入 |

AI 定制本地文件时无需手动维护指纹。specflow update 识别为"用户已修改"是正常的。

## 本地定制边界

默认可编辑：

- `.specflow/workflow.md`
- `.specflow/config.yaml`
- `.specflow/agents.yaml`
- `.specflow/spec/**`
- `.opencode/agents/**`（custom agent 定义）
- 平台 hooks/plugins/skills/agents

默认不编辑：

- specflow Go 二进制
- specflow CLI 源码仓库
- `.specflow/.runtime/**` 下的具体状态文件
- `.specflow/.fingerprints.json` 的指纹内容

只有当用户明确想贡献上游时，才切换到 specflow CLI 源码视角。
