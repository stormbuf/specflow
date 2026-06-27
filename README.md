# Specflow

Spec 驱动的变更生命周期管理工具。通过 Go CLI + 平台插件 + auto-trigger skills，让 AI 编码助手在变更全生命周期中自动获得正确的上下文、阶段感知和跨会话记忆。

## 核心特性

- **上下文自动就位** — 每次派发子 agent 时，相关任务工件与规范自动注入，无需手动指定文件
- **流程阶段感知** — AI 每轮对话自动知道"现在在第几步、下一步做什么"，流程不再靠人盯
- **跨会话记忆** — 项目规范、踩坑经验、会话日志持久留存，新会话不再从零开始
- **配置驱动扩展** — 新增 agent 类型或接入第三方 agent，只需改配置文件，不改代码
- **版本管理中立** — 不管用 git 还是 jj，开箱即用
- **状态注入不污染对话** — 工作流上下文作为独立消息注入，你的原始对话内容保持干净

## 架构

三层职责，各司其职：

| 层 | 载体 | 职责 |
|----|------|------|
| skill 层 | `.opencode/skills/` | auto-trigger skills + 命令式 skill |
| 插件层 | `.opencode/plugins/` | 上下文注入、工作流状态同步、会话启动自动化 |
| 状态层 | `.specflow/` | 任务状态、工件、规范库、会话日志、运行时指针 |

## 安装

### 从源码编译

```bash
cd specflow-cli
go build -o specflow .
```

### Homebrew（规划中）

```bash
brew install ./specflow.rb
```

### 初始化项目

```bash
specflow init -u <developer> --opencode
```

init 会自动检测版本管理系统（`.jj/` 优先，其次 `.git/`），安装三层结构，记录文件指纹。安装完成后重启 AI Agent。

## 使用

### 典型工作流

```bash
# 1. 用户描述需求，AI 经确认后建任务
specflow task create --title "添加导出功能" --intent "用户需要导出数据为 CSV"

# 2. AI 加载 specflow-brainstorm 编写 prd.md / implement.md
# 3. 整理 jsonl 上下文清单
specflow add-context .specflow/changes/<change-id> specflow-implement ".specflow/spec/backend/coding-style.md" "编码规范"

# 4. 开始任务
specflow task start

# 5. AI 派发 specflow-implement 子 agent（上下文由插件自动注入）
# 6. AI 派发 specflow-check 子 agent 验证
# 7. 同步行为规约、更新规范库

# 8. 归档
/specflow:finish-work
```

### 命令速查

| 命令 | 用途 |
|------|------|
| `specflow init` | 初始化项目 |
| `specflow task create/start/finish/archive` | 任务生命周期 |
| `specflow task current/list` | 查看任务状态 |
| `specflow task release` | 清理 stale session 指针 |
| `specflow get-context` | 聚合 session 上下文 |
| `specflow build-context` | 构建子 agent 上下文 |
| `specflow add-context` | 管理 jsonl 上下文清单 |
| `specflow sync-agent` | 同步 custom agent |
| `specflow add-session` | 记录会话日志 |
| `specflow validate` | 校验配置 |
| `specflow doctor` | 诊断项目健康 |
| `specflow mem search` | 跨会话对话检索 |
| `specflow update` | 同步项目到 CLI 版本 |
| `specflow upgrade` | 升级 CLI 二进制 |
| `/specflow:continue` | 任务内推进下一步 |
| `/specflow:finish-work` | 归档任务 + 写会话日志 |

### Agent 接入

在 `.specflow/agents.yaml` 中声明 agent，支持三种来源：

```yaml
agents:
  specflow-implement:          # native: specflow 内置，随 CLI 打包
    source: native
    jsonl_file: implement.jsonl
    constraints:
      - "禁止 git commit / push / merge"

  code-reviewer:               # platform: 复用平台已有 agent，不同步
    source: platform
    jsonl_file: code-reviewer.jsonl

  specflow-oracle:             # custom: 自建，sync-agent 同步到平台
    source: custom
    jsonl_file: oracle.jsonl
    agent_file: agents/oracle.md
```

### 版本管理

两步升级体系：

- `specflow upgrade` — 升级全局 CLI 二进制
- `specflow update` — 同步项目到本地 CLI 版本，含文件指纹冲突检测

## 文档

- [重构设计思路](docs/重构设计思路.md) — 背景、调研、架构决策、方案设计
- [接口契约](docs/接口契约.md) — 五个模块的详细接口定义

## 许可证

MIT
