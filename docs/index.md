# Specflow

**Spec 驱动的变更生命周期管理工具**

Specflow 让 AI 编码助手在变更全生命周期中自动获得正确的上下文、阶段感知和跨会话记忆。

??? quote "核心特性"
    - **上下文自动就位** — jsonl manifest 声明上下文文件与顺序，插件自动注入
    - **流程阶段感知** — workflow.md 面包屑每轮自动注入当前阶段
    - **跨会话记忆** — journal 记录会话日志，mem 检索历史对话
    - **自动生成规范** — specflow-spec-bootstrap skill 从代码库分析生成项目专属 spec
    - **配置驱动扩展** — 改流程 = 改 markdown，加 agent = 改 yaml
    - **版本管理中立** — 同时支持 Git 和 JJ
    - **状态注入不污染对话** — 非破坏性注入，绝不修改用户消息原文

## 快速开始

```bash
# 安装
brew install stormbuf/tap/specflow

# 初始化项目
specflow init -u 你的名字 --opencode

# 重启 AI Agent，开始使用
```

## 文档

- [什么是 Vibe Coding](tutorial/01-what-is-vibe-coding.md) — 从零理解 Spec 驱动的 AI 编码
- [Specflow 是什么](tutorial/02-what-is-specflow.md) — 痛点对比与核心能力
- [架构总览](tutorial/03-architecture.md) — 三层结构与设计哲学
- [模块详解](tutorial/04-modules.md) — CLI、Skills、Agents、Plugins、Runtime
- [工作流程](tutorial/05-workflow.md) — 规划 / 执行 / 收尾三阶段
- [快速上手](tutorial/06-quick-start.md) — 安装、初始化、第一个任务
- [Skill 详解](tutorial/07-skills.md) — 11 个 auto-trigger skill
- [使用场景](tutorial/08-use-cases.md) — 什么时候用 / 不用
- [实战示例](tutorial/09-example.md) — CSV 导出功能端到端流程
- [命令速查](tutorial/10-reference.md) — 全部 CLI 命令与配置文件

## 相关链接

- [GitHub 仓库](https://github.com/stormbuf/specflow)
- [Homebrew Tap](https://github.com/stormbuf/homebrew-tap)
- 架构设计参考 [Trellis](https://github.com/mindfold-ai/Trellis)（AGPL-3.0）

## 许可证

MIT License
