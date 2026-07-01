# 架构总览

Specflow 在你的项目里安装三层结构，各司其职：

```
你的项目
├── .opencode/                    ← 平台集成层
│   ├── skills/                   specflow 内置 skill（11 个）
│   ├── plugins/                  3 个 JS 插件（hook 自动化）
│   ├── lib/                      共享库
│   └── agents/                   native + custom agent 定义
│
├── .specflow/                    ← 工作流与状态层
│   ├── workflow.md               工作流契约（面包屑 + 路由）
│   ├── config.yaml               项目配置
│   ├── agents.yaml               agent 声明
│   ├── spec/                     规范库（自更新、团队共享）
│   │   ├── index.md              全局索引 + Pre-Dev Checklist
│   │   ├── <layer>/               分层规范（backend/frontend/...）
│   │   └── requirements/          行为规约（prd 派生）
│   ├── changes/                   变更任务工件
│   │   ├── <change-id>/           当前任务
│   │   │   ├── task.json         任务元数据 + 状态
│   │   │   ├── prd.md             需求文档
│   │   │   ├── implement.md       执行计划
│   │   │   ├── implement.jsonl    上下文清单
│   │   │   └── check.jsonl        验收上下文清单
│   │   └── archive/               归档任务
│   ├── workspace/                 跨会话 journal
│   └── .runtime/                  session 指针（gitignored）
│
└── AGENTS.md                     ← 项目规则注入点
```

| 层 | 载体 | 职责 | 进 Git |
|----|------|------|--------|
| Skill 层 | `.opencode/skills/` | auto-trigger skills + 命令式 skill，承载 know-how、模板、审查规则 | 是 |
| 插件层 | `.opencode/plugins/` | hook 自动化：上下文注入、工作流状态同步、会话启动 | 是 |
| 状态层 | `.specflow/` | 任务状态、工件、规范库、会话日志、运行时指针 | 部分 |

!!! tip "设计哲学"
    - **契约与实现分离**：workflow.md 是流程唯一事实源，插件只解析不存逻辑。
    - **状态驱动而非文件驱动**：阶段由 task.json 的 status 字段驱动面包屑切换。
    - **上下文窄而准**：上下文文件由 jsonl manifest 声明，按需加载，不搞大而全的固定契约。
