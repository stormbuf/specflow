# 提交与分支策略

定义 Commit message 格式、分支命名、分支策略与 PR/MR 规范。

## Overview

<!-- 引导问题：回答以下问题后再填写具体章节 -->

<!--
- Commit message 格式是什么？是否使用 Conventional Commits？
- 分支命名约定是什么？
- 分支策略是什么（Git Flow / GitHub Flow / Trunk-based）？
- PR / MR 有什么规范？需要几个 reviewer？
- 代码审查的要求是什么？
- 是否使用 squash merge / rebase merge / merge commit？
- 是否有 commit 前 hook（lint / format / 密钥扫描）？
-->

（由团队填写：概述 Git 工作流的核心策略与原则。）

## Commit Message 格式

项目使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范。以下是格式说明：

### 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### type（必填）

| type | 含义 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat: 新增用户导出功能` |
| `fix` | Bug 修复 | `fix: 修复夜间模式样式错误` |
| `docs` | 文档变更 | `docs: 更新 README 安装说明` |
| `style` | 代码格式（不影响逻辑） | `style: 统一缩进为 2 空格` |
| `refactor` | 重构（非新功能、非修复） | `refactor: 提取表单验证逻辑` |
| `perf` | 性能优化 | `perf: 优化列表分页查询` |
| `test` | 测试相关 | `test: 补充用户模块单元测试` |
| `chore` | 构建 / 工具 / 依赖变更 | `chore: 升级 express 到 4.19` |
| `ci` | CI 配置变更 | `ci: 添加 GitHub Actions 工作流` |
| `build` | 构建系统变更 | `build: 修改 webpack 配置` |
| `revert` | 回退之前的 commit | `revert: feat: 新增用户导出功能` |

### scope（可选）

影响范围，如模块名或功能名：

```
feat(auth): 新增 OAuth2 登录
fix(api): 修复分页参数解析错误
```

### subject（必填）

- 使用中文描述，简洁明了
- 不超过 50 个字符
- 不加句号结尾
- 使用祈使句（如"新增"而非"新增了"）

### body（可选）

- 解释"为什么"做这个变更，而非"做了什么"（diff 已展示做了什么）
- 每行不超过 72 个字符
- 可以使用 ` - ` 列出多项要点

### footer（可选）

- BREAKING CHANGE：标注破坏性变更
- 关联 issue：`Closes #123` / `Refs #456`

### 完整示例

```
feat(auth): 新增 OAuth2 登录支持

用户现在可以通过 Google 和 GitHub 账号登录。
使用 passport-google-oauth20 和 passport-github 策略。

- 添加 OAuth2 回调路由
- 数据库新增 provider 和 provider_id 字段
- 前端添加第三方登录按钮

Closes #142
```

### Breaking Change 示例

```
feat(api)!: 重构用户 API 响应格式

将 name 字段拆分为 firstName 和 lastName。
客户端需要同步更新。

BREAKING CHANGE: GET /api/users 响应中的 name 字段移除，替换为 firstName 和 lastName
```

（由团队填写：补充团队对 commit message 的额外约定或调整。）

## 分支命名约定

（由团队填写：定义分支的命名规则。包括：

- 功能分支：`feat/<描述>` / `feature/<描述>`
- 修复分支：`fix/<描述>` / `bugfix/<描述>`
- 热修复分支：`hotfix/<描述>`
- 发布分支：`release/<版本号>`
- 其他约定（如包含 issue 编号：`feat/PROJ-123-user-export`）
- 命名风格（kebab-case / snake_case）
- 分支名长度限制
- 代码示例：分支创建命令
）

## 分支策略

（由团队填写：定义项目的分支模型。包括：

### 策略选型

- **Git Flow** — 适合有明确发布周期的项目（main / develop / feature / release / hotfix）
- **GitHub Flow** — 适合持续部署的项目（main + feature 分支）
- **Trunk-based Development** — 适合高频部署、强 CI 的团队（短生命周期分支 + main）

### 主分支
- main / master 分支的保护规则
- 谁有权限直接推送
- 是否要求 PR + 审查

### 功能分支
- 从哪个分支创建
- 生命周期限制（如不超过 3 天）
- 合并目标分支

### 发布分支（如使用 Git Flow）
- 何时创建
- 发布后的合并方向

### 热修复分支
- 从哪个分支创建
- 合并到哪些分支

（由团队填写：选择团队实际使用的策略并补充细节。）
）

## PR / MR 规范

（由团队填写：定义 Pull Request / Merge Request 的规范。包括：

### PR 标题
- 是否使用 Conventional Commits 格式
- 是否关联 issue 编号

### PR 描述模板

```markdown
## 变更说明
（简要描述本 PR 做了什么以及为什么）

## 变更类型
- [ ] 新功能（feat）
- [ ] Bug 修复（fix）
- [ ] 重构（refactor）
- [ ] 文档（docs）
- [ ] 其他

## 测试
- [ ] 已添加 / 更新测试
- [ ] 本地测试通过
- [ ] CI 通过

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 无 console.log / 调试代码残留
- [ ] 无硬编码密钥 / 敏感信息
- [ ] 已更新相关文档

## 关联 Issue
Closes #
```

### 审查要求
- 最少 reviewer 数量
- 谁可以审查（任意成员 / 指定 owner）
- 审查时限要求
- 审查重点（逻辑 / 安全 / 性能 / 规范）

### 合并方式
- Squash merge（推荐：保持 main 历史整洁）
- Rebase merge
- Merge commit
- 合并后是否删除源分支
）

## 常见错误

（由团队填写：列出团队在 Git 工作流上犯过的错误，例如：

- Commit message 只写"fix"或"update"，无法从历史中理解变更内容
- 一个 commit 包含多个不相关的变更，无法单独回退
- 功能分支生命周期过长（数周），合并时冲突严重
- 直接推送到 main 分支，绕过 PR 审查
- PR 描述为空或只写"同标题"，reviewer 无从了解上下文
- 分支命名随意（如 `test` / `tmp` / `wip`），无法识别用途
- 使用 merge commit 导致 main 历史混乱，难以阅读
- Commit 中包含调试代码 / 注释掉的代码 / .env 文件
- hotfix 只合并到 main 未合并到 develop，下次发版问题复现
- 大量 "WIP" commit 堆积在分支上，合并时未 squash
- Commit message 中的 type 与实际变更不符（如 refactor 标为 feat）
）
