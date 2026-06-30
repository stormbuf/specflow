# AGENTS.md

## 构建与开发

- **编译 CLI**：`cd specflow-cli && go build -o specflow .`
- **Go module**：模块名 `specflow`（本地模块，非远程 import path）
- **依赖**：`cobra`、`yaml.v3`（见 `specflow-cli/go.mod`）
- **编译产物** `specflow-cli/specflow` 已在 `.gitignore` 中排除，不要提交
- **embed 资源同步**：`specflow-cli/resources/` 是 `skills/`、`agents/`、`specflow-runtime/`、`platforms/`、`spec-templates/` 的内嵌副本，通过 `//go:embed resources` 编译进二进制。修改源目录后必须同步：`rsync -av skills/ specflow-cli/resources/skills/`（其他目录同理），否则编译产物不含最新改动

## 项目结构

| 目录 | 说明 |
|------|------|
| `specflow-cli/` | Go CLI 源码，含 8 个 internal 包（config/taskstore/session/vcs/fingerprint/context/installer/worktree） |
| `skills/` | 11 个 auto-trigger skill（含 specflow-spec-bootstrap、specflow-session-insight、specflow-meta） + 5 个模板 |
| `agents/` | 3 个 native agent 定义（specflow-implement / specflow-check / specflow-research） |
| `specflow-runtime/` | 运行时模板（workflow.md / config.yaml / agents.yaml / spec seed / workspace index / worktree.yaml / agents.md.tmpl） |
| `spec-templates/` | spec 模板源文件（9 个分类：guides / backend / frontend / architecture / testing / security / api / devops / git-conventions），每个分类含 .meta.yaml 声明描述，供 `spec install` 读取 |
| `platforms/opencode/` | OpenCode 平台适配（3 个 JS 插件 + 1 个共享 lib + install-map.yaml） |
| `specflow-workflow/` | **旧版单体 skill，已废弃，不要修改** |

## 版本管理

- 本仓库使用 jj。使用 `jj log`、`jj diff`、`jj describe`、`jj new`；**不要使用** `git add`、`git commit`、`git stash`、`git checkout`。
- 用户要求提交时，执行 `jj describe -m "..."` 后再执行 `jj new`，确保当前变更被封存，后续工作进入新变更。
- specflow CLI 同时支持 git 和 jj，通过 `config.yaml:vcs` 配置或 `specflow init` 时自动扫描检测（优先 jj）。

## 边界

- ✅ 可以修改 `specflow-cli/`、`skills/`、`agents/`、`specflow-runtime/`、`spec-templates/`、`platforms/` 下的源文件
- ✅ 修改源目录后同步 `specflow-cli/resources/` embed 副本
- ⚠️ 修改 `skills/` 或 `agents/` 后需同步对应的 `specflow-cli/resources/` 副本
- ❌ 不要修改 `specflow-workflow/`（已废弃）
- ❌ 不要提交 `specflow-cli/specflow` 二进制文件
- ❌ 不要提交 `.specflow/state.json` 和 `.opencode/` 中未显式放行的本地状态

## 当前仓库事实

- Go CLI 已实现并通过端到端验证。
- 11 个 skill 已实现（含 specflow-spec-bootstrap 从代码库自动生成 spec、specflow-session-insight 跨会话记忆检索、specflow-meta 架构理解与定制），目前是 placeholder 级别，后续根据实际使用迭代。
- OpenCode 插件已实现（3 个 JS 插件 + 1 个共享 lib），通过 `node --check` 语法校验。
- `specflow-cli/resources/` 是 embed 副本，修改源目录后需手动同步，否则编译产物不含最新改动。
