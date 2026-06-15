# AGENTS.md

## 版本管理

- 本仓库使用 jj。使用 `jj log`、`jj diff`、`jj describe`、`jj new`；本地工作流不要使用 `git add`、`git commit`、`git stash` 或 `git checkout`。
- 用户要求用 jj 提交时，执行 `jj describe -m "..."` 后再执行 `jj new`，确保当前变更被封存，后续工作进入新变更。

## 项目形态

- 本仓库使用根目录 `specflow-workflow/` 作为稳定工作流 skill；将 `specflow-workflow/` 放到目标平台的 skills 目录中即可使用（如 `.agents/skills/`、`.claude/skills/`、`.cursor/skills/` 等）。
- 修改 `specflow-workflow/` 后，如需同步到其他项目，必须使用项目级 skill `install-specflow-workflow`，不要手写目标项目名或直接复制遗漏资源；该 skill 会扫描 `~/project/`、让用户选择目标项目，并执行项目级安装。

## 安装到其他项目

将当前工作区的 `specflow-workflow` 复制安装到另一个项目。

执行规则：

- 扫描 `~/project/` 下的一级目录。
- 使用 `question` 工具让用户选择目标项目。
- 将 `specflow-workflow/` 目录复制到目标项目的 `.agents/skills/` 目录下。
- 不使用 `-g`，只做项目级安装。
- 安装后提醒用户重启目标项目中的 AI Agent。

## 当前仓库事实

- 预期工作流工件仍在定义中。仓库出现可执行配置前，不要臆造 build、test、lint 或 codegen 命令。
- 当前仓库没有 manifest、源码树、CI 配置或测试运行器。
- `specflow/state.json` 和 `.agents/` 中未显式放行的内容是被忽略的本地状态，不应提交。
