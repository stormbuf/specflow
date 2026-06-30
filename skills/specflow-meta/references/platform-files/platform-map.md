# 平台文件路径

本页列出 specflow 在用户项目中的常见文件位置。当前仅支持 OpenCode 平台。

## 当前矩阵

| 平台 | CLI flag | 主目录 | Skill 目录 | Agent 目录 | 插件目录 |
| --- | --- | --- | --- | --- | --- |
| OpenCode | `--platform opencode` | `.opencode/` | `.opencode/skills/` | `.opencode/agents/` | `.opencode/plugins/` + `.opencode/lib/` |

`specflow init` 默认平台为 `opencode`，也可通过 `--platform` 指定。`.specflow/config.yaml` 的 `platform` 字段记录当前平台。

## OpenCode 文件清单

```text
.opencode/
├── skills/                  # bundled skill 副本
│   ├── specflow-meta/
│   ├── specflow-brainstorm/
│   ├── specflow-check/
│   ├── specflow-implement/  # (作为 skill 副本，非 agent)
│   ├── specflow-research/   # (作为 skill 副本，非 agent)
│   ├── specflow-before-dev/
│   ├── specflow-break-loop/
│   ├── specflow-continue/
│   ├── specflow-finish-work/
│   ├── specflow-sync-requirements/
│   ├── specflow-update-spec/
│   ├── specflow-spec-bootstrap/
│   └── specflow-session-insight/
├── plugins/                 # 3 个 JS 插件
│   ├── session-start.js
│   ├── inject-workflow-state.js
│   └── inject-subagent-context.js
├── lib/                     # 插件共享工具
│   └── specflow-context.js
└── agents/                  # native agent 定义
    ├── specflow-implement.md
    ├── specflow-check.md
    └── specflow-research.md
```

## install-map.yaml

`.opencode/` 的安装映射由 `platforms/opencode/install-map.yaml` 定义：

```yaml
platform: opencode
install_targets:
  skills: ".opencode/skills/"
  plugins: ".opencode/plugins/"
  agents: ".opencode/agents/"
```

`specflow init` 和 `specflow update` 读取此文件决定内嵌资源安装到哪些目录。

## 多平台扩展点

specflow 预留了多平台扩展点。`platforms/` 目录结构为：

```text
platforms/
└── <platform>/
    ├── install-map.yaml     # 该平台的安装映射
    ├── plugins/             # 该平台的插件
    ├── lib/                 # 该平台的共享库
    └── agents/              # 该平台的 native agent 定义
```

当前仅有 `platforms/opencode/`。未来新增平台时，在 `platforms/` 下新建对应目录，specflow CLI 的 `installer` 包会自动发现并安装。

## 平台文件与共享文件的关系

| 共享文件 | 对应平台文件 |
| --- | --- |
| `.specflow/workflow.md`（面包屑标签块） | `inject-workflow-state.js`（解析并注入） |
| `.specflow/agents.yaml`（agent 声明） | `inject-subagent-context.js`（查配置注入约束） |
| `.specflow/changes/<task>/*.jsonl`（上下文清单） | `inject-subagent-context.js`（exec `build-context` 加载） |
| `.specflow/spec/`（spec 索引） | `session-start.js`（exec `get-context` 列出索引） |

## 修改平台文件时的决策规则

1. 用户指定平台：只改该平台目录，除非共享 workflow/spec 也要改。
2. 用户说"我的 AI"：检查项目中实际存在的配置目录，推断当前平台。
3. 用户要项目规则：优先 `.specflow/spec/` 或项目本地 skill。
4. 用户要 specflow 行为：改 `.specflow/workflow.md` + 平台插件/agents/skills。

## 路径不一致时

平台生态会变化，用户项目可能已有自定义。如果本表与本地文件不一致，以用户项目内实际文件为权威：

- 检查 `install-map.yaml` 注册的安装目标。
- 检查插件实际读取的路径。
- 以 agent 文件中当前写入的读取规则判断行为。

不要因为某个文件不在本路径表中就删除它。
