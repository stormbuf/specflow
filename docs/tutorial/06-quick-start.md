# 快速上手

## 安装

```bash
# Homebrew（推荐）
brew install stormbuf/tap/specflow

# 或 go install
go install github.com/stormbuf/specflow/specflow-cli@latest

# 或从源码编译
cd specflow-cli && go build -o specflow .
```

## 初始化项目

```bash
specflow init -u 你的名字 --opencode
```

init 会自动检测版本管理系统（`.jj/` 优先，其次 `.git/`），安装三层结构，记录文件指纹，向项目 `AGENTS.md` 注入 managed block。安装完成后会提示选择 spec 模板（可跳过）。

| 选项 | 效果 |
|------|------|
| `--all-spec` | 安装所有 spec 模板 |
| `--with-spec guides,backend` | 指定分类安装（逗号分隔） |
| `--no-spec` | 跳过 spec 模板 |

init 还支持以下选项：

| 选项 | 效果 |
|------|------|
| `--vcs git\|jj` | 手动指定版本管理工具（默认自动检测） |
| `--force` | 覆盖已存在的 .specflow/ 目录 |
| `--platform opencode\|pi` | 指定目标平台（与 --opencode/--pi 互斥） |

安装完成后 **重启 AI Agent**，让插件生效。

## 第一个任务

在 OpenCode 中跟 AI 对话，描述你的需求。AI 会判断是否需要建任务：

```
你：我想给项目加一个 CSV 导出功能

AI：这个需求涉及多个文件的改动，建议创建一个 specflow 任务来管理。
    是否创建？

你：好

AI：[执行 specflow task create]
    [加载 specflow-brainstorm，编写 prd.md 和 implement.md]
    [整理 jsonl 上下文清单]
    [执行 specflow task start]

    任务已激活，进入执行阶段。我来派发 implement agent 实现功能。

    [派发 specflow-implement，插件自动注入上下文]
    [实现完成，派发 specflow-check 验收]
    [验收通过，同步需求，更新 spec]

    实现完成，可以归档了。执行 /specflow:finish-work 归档。

你：/specflow:finish-work

AI：[归档任务，写 journal，auto-commit]
    任务已归档。
```
