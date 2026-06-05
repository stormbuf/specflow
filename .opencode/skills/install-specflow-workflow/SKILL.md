---
name: install-specflow-workflow
description: "Install specflow-workflow into another OpenCode project. Use when the user asks to install, distribute, or copy specflow workflow to a selected project."
metadata:
  version: "0.1.0"
---

# Install Specflow Workflow

> AI 使用本 skill 将当前工作区的 `specflow-workflow` 安装到另一个 OpenCode 项目的项目级 skills 中。
> 最高信条：目标项目必须由用户选择；不得硬编码项目名或路径。

## 路径约定

- 当前工作区根目录：OpenCode 当前 workspace root。
- specflow skill 源目录：当前工作区根目录下的 `specflow-workflow/`。
- 默认项目根目录：`~/project/`。
- 目标项目：由扫描结果和用户选择决定

## 触发

当用户要求把 specflow、specflow-workflow 或当前 workflow skill 安装、分发、复制到其他 OpenCode 项目时执行。

## 执行流程

```text
1. 解析当前工作区根目录下的 specflow-workflow/SKILL.md，得到源 skill 目录绝对路径。
2. 扫描 ~/project/ 下的一级目录。
3. 过滤掉当前工作区目录和明显非项目目录。
4. 使用 question 工具列出候选项目，让用户选择一个目标项目。
5. 用户选择后，在目标项目目录执行安装命令。
6. 安装完成后提示用户重启目标项目中的 OpenCode 会话。
```

## 项目扫描规则

```text
候选目录 = ~/project/* 中的一级目录
排除目录:
  - 当前工作区目录
  - 以 . 开头的目录
  - 不可读目录
优先展示:
  - 包含 .opencode/ 的目录
  - 包含 .git/ 或 .jj/ 的目录
  - 包含 package.json、pyproject.toml、Cargo.toml、go.mod、README.md 的目录
```

若候选项目超过 20 个，只展示最可能的 20 个，并允许用户自定义路径。

## 安装命令

在用户选择的目标项目目录执行：

```bash
npx skills add "<resolved-source-skill-dir>" -a opencode --copy -y
```

使用 `--copy`，让目标项目获得独立副本，避免当前 specflow 仓库后续修改影响目标项目。

## 验证

安装后检查目标项目中是否出现项目级 skill：

```text
优先检查:
  <target>/.opencode/skills/specflow-workflow/SKILL.md
如果 CLI 使用其他项目级目录:
  运行 npx skills list -a opencode，并确认 specflow-workflow 出现
```

## 阻断

```text
IF 当前工作区根目录/specflow-workflow/SKILL.md 不存在:
  暂停，报告源 skill 缺失
ELSE IF 未扫描到候选项目:
  询问用户输入目标项目绝对路径
ELSE IF 用户未选择目标项目:
  暂停，不执行安装
ELSE IF 目标路径不存在或不是目录:
  暂停，报告路径无效
ELSE IF 安装命令失败:
  报告命令、退出码和关键错误输出
```

## 完成条件

- 用户已明确选择目标项目。
- 安装命令在目标项目目录执行完成。
- 已验证 `specflow-workflow` 出现在目标项目项目级 OpenCode skills 中，或报告无法验证的原因。
- 最终回复包含目标项目路径、安装方式和重启提醒。

## 关键约束

1. 不硬编码任何目标项目名。
2. 执行安装前必须使用 question 工具让用户选择目标项目。
3. 不使用 `-g`，只做项目级安装。
4. 默认使用 `--copy`，不创建跨项目共享 symlink。
5. 安装后提醒用户重启目标项目中的 OpenCode。
