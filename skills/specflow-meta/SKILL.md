---
name: specflow-meta
description: "理解并定制 specflow 在用户项目内的本地架构。当需要修改 .specflow/ 目录、平台 hooks/plugins/agents/skills、工作流定义、任务系统、spec 结构、上下文注入机制、跨会话记忆（specflow mem）、或 bundled skill 自动分发机制时触发。适用于已运行 specflow init 的项目，帮 AI 理解 specflow 本地架构和定制入口。"
trigger: "用户要求理解或定制 specflow 本地架构、改工作流、改任务生命周期、改上下文加载、改 agent 行为、改 skills、改 spec 结构、加项目本地约定"
---

# specflow Meta

本 skill 面向已在项目中运行过 `specflow init` 的用户。读完之后，AI 应当理解 specflow 在该项目内的架构、运作模型和定制入口，然后根据用户请求修改生成的 `.specflow/` 和平台目录文件。

specflow 是一个 Go 单二进制 CLI，通过 `go:embed` 将 skills、agents、runtime 模板、平台插件和 spec 模板编译进可执行文件。与 npm 分发模式不同，specflow 不依赖 `node_modules`，所有管理文件通过 `specflow init` 写入项目，通过 `specflow update` 做三路指纹比对增量同步。当前仅支持 OpenCode 平台，但 `platforms/` 目录结构预留了多平台扩展点。

默认操作范围是用户项目内的本地文件：

- `.specflow/`：工作流、配置、任务、spec、workspace 记忆、运行时状态。
- 平台目录：当前仅 `.opencode/`（skills / plugins / lib / agents）。
- `.specflow/.fingerprints.json`：管理文件指纹，`specflow update` 据此判断哪些文件已被用户修改。

不要假设用户拥有 specflow CLI 源码仓库。不要默认引导用户去 fork specflow CLI。

## specflow 三层架构

specflow 在用户项目内提供三层：

1. **工作流层**：`.specflow/workflow.md` 定义阶段、skill 路由和面包屑标签块。
2. **持久层**：`.specflow/changes/`（任务）、`.specflow/spec/`（工程规范）、`.specflow/workspace/`（会话记忆）。
3. **平台集成层**：`.opencode/` 下的 hooks/plugins/agents/skills 连接 specflow 工作流与 AI 工具。

三层均位于用户项目内，AI 可直接读取和修改。

## How To Use

1. 先读 `references/local-architecture/overview.md` 建立 specflow 本地系统模型。
2. 如果请求涉及特定 AI 工具平台，读 `references/platform-files/platform-map.md` 和相关平台文件说明。
3. 如果用户想改行为，读 `references/customize-local/overview.md`，再打开对应定制主题。
4. 编辑前，先读用户项目内的实际文件，以本地内容为权威。
5. 改完后确认语义同步：共享流程变了，检查平台入口文件是否需要同步。

## References

### Local Architecture

- `references/local-architecture/overview.md`：specflow 本地三层架构总览与定制原则。
- `references/local-architecture/generated-files.md`：`specflow init` 生成的文件清单与定制边界。
- `references/local-architecture/workflow.md`：阶段、路由、面包屑标签块系统。
- `references/local-architecture/task-system.md`：`.specflow/changes/` 任务结构、父子任务、session 独占。
- `references/local-architecture/spec-system.md`：`.specflow/spec/` 的组织方式、注入机制与模板安装。
- `references/local-architecture/workspace-memory.md`：`.specflow/workspace/` journal 与 `specflow mem` 跨会话检索。
- `references/local-architecture/context-injection.md`：session context、workflow context、subagent context 注入与 #367 规避。

### Customize Local

- `references/customize-local/overview.md`：选择正确的本地定制入口。
- `references/customize-local/change-workflow.md`：改阶段、路由、面包屑标签块。
- `references/customize-local/change-task-lifecycle.md`：改任务创建、状态、归档、父子关系。
- `references/customize-local/change-context-loading.md`：改 jsonl 清单、agents.yaml、上下文注入。
- `references/customize-local/change-agents.md`：改 native/platform/custom agent 行为。
- `references/customize-local/change-skills-or-commands.md`：区分 bundled skill 与项目本地 skill。
- `references/customize-local/change-spec-structure.md`：调整 `.specflow/spec/` 结构与模板安装。
- `references/customize-local/add-project-local-conventions.md`：把团队规则放进 spec 或本地 skill。

### Platform Files

- `references/platform-files/overview.md`：共享文件与平台文件的关系。
- `references/platform-files/platform-map.md`：当前 OpenCode 平台文件路径与多平台扩展点。
- `references/platform-files/hooks-and-plugins.md`：3 个 JS 插件与共享 lib 机制。

## Current Rules

- `.specflow/workflow.md` 是本地工作流唯一事实源。插件 `inject-workflow-state.js` 按当前任务 status 匹配 `[workflow-state:<STATUS>]` 面包屑标签块并注入每轮用户消息。改面包屑只需改 markdown，下一条消息即生效。
- `.specflow/config.yaml` 是项目级配置入口。包含 `vcs`（git|jj）、`platform`（opencode）、`max_journal_lines`、`mem` 配置（`enabled` / `log_paths`）、`session.stale_threshold_hours`。
- `.specflow/agents.yaml` 声明 agent 的 source（native|platform|custom）、jsonl_file、constraints 等字段。`specflow validate` 校验声明的合法性。
- `.specflow/spec/` 存项目工程规范，按分类组织（guides/backend/frontend/architecture/testing/security/api/devops/git-conventions），每层有 `index.md`。可通过 `specflow spec install` 安装模板。
- `.specflow/changes/` 存任务（task.json/prd.md/implement.md/design.md/implement.jsonl/check.jsonl）。change-id 格式为 `YYYY-MM-DD-short-slug-N`。任务有 session 独占机制，归档后移动到 `changes/archive/<YYYY-MM>/`。
- `.specflow/workspace/` 存跨会话 journal（按 developer 分目录）。journal 由 `specflow finish-work` 写入，记录完成的会话工作。
- `.specflow/.fingerprints.json` 跟踪 specflow 管理文件的内容指纹。`specflow update` 通过三路比对（旧指纹 vs 当前磁盘 vs CLI 新版本）决定覆盖/保留/冲突。
- `.specflow/.runtime/sessions/` 存 session 级别的活跃任务指针，按 `<platform>_<session-id>.json` 隔离，不同 AI 窗口可指向不同任务。

## Do Not

- 不要把 specflow CLI 源码作为本地定制的默认目标。
- 不要修改 Go 二进制文件或 `go:embed` 资源来实现项目需求；改项目内的生成文件。
- 不要用默认模板覆盖用户已修改的本地文件；先查 `.fingerprints.json`，优先生成 `.new` sidecar 而非破坏性覆盖。
- 不要把团队私有项目规则放进 bundled skill（`specflow-meta` 等）；项目规则放 `.specflow/spec/`、项目本地 skill、当前任务或 workspace journal —— `specflow update` 会覆盖 bundled skill 目录内的内容。
- 不要手编 `.specflow/.runtime/sessions/` 下的 session 指针文件来"修复"业务状态。
- 不要把已移除或从未存在的机制描述为当前 specflow 行为；先对照本地 `.specflow/config.yaml` 和 `specflow --help` 确认。
