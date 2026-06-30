# Spec 系统

`.specflow/spec/` 是用户项目的工程规范库。specflow 不是让 AI 死记约定，而是在合适的时机注入相关 spec 或要求 AI 读取。

## 目录模型

常见单仓库结构：

```text
.specflow/spec/
├── index.md
├── backend/
│   ├── index.md
│   └── ...
├── frontend/
│   ├── index.md
│   └── ...
└── guides/
    ├── index.md
    └── ...
```

`index.md` 是每层的入口，列出 Pre-Development Checklist 和 Quality Check。具体规范放在同目录的其他 Markdown 文件中。

## 分类体系

specflow 内置 9 个 spec 分类模板，可通过 `specflow spec install` 安装：

| 分类 | 说明 |
| --- | --- |
| `guides` | 跨层思维指引（代码复用、跨层、跨平台） |
| `backend` | 后端规范（目录结构、错误处理、数据库、日志） |
| `frontend` | 前端规范（组件、状态管理、类型安全、hook） |
| `architecture` | 架构决策记录（ADR 模板） |
| `testing` | 测试规范（约定、集成模式、mock 策略） |
| `security` | 安全规范（认证、输入校验、密钥管理） |
| `api` | API 规范（REST 约定、版本管理、错误响应） |
| `devops` | DevOps 规范（CI/CD、部署、发布流程） |
| `git-conventions` | Git 约定（提交、分支） |

每个分类在 `spec-templates/` 下有 `.meta.yaml` 声明描述，供 `specflow spec install` 读取。

## Spec 模板安装

```bash
specflow spec install backend      # 安装 backend 分类模板
specflow spec install --all        # 安装所有分类模板
```

安装的模板是起点，不是契约。用户应根据实际项目代码调整、删除、拆分或新增 spec 文件。

## `specflow-spec-bootstrap` skill

`specflow-spec-bootstrap` skill 可从真实代码库出发，自动分析架构并生成或刷新 `.specflow/spec/` 规范。它拒绝占位符文本，要求每条重要规则指向真实文件或反复出现的本地 pattern。

## Spec 如何进入任务

任务进入实现前，规划阶段将相关 spec 写入 `implement.jsonl` / `check.jsonl`：

```jsonl
{"file": ".specflow/spec/backend/index.md", "reason": "后端约定"}
{"file": ".specflow/spec/testing/conventions.md", "reason": "测试期望"}
```

sub-agent 启动时由 `inject-subagent-context.js` 插件读取 jsonl 清单，加载引用的 spec 内容。

## Spec 应包含什么

Spec 应包含项目可执行的工程约定，而非通用最佳实践：

- 文件应放在哪里。
- 错误处理如何表达。
- API、hook、command 的输入输出契约。
- 禁止的 pattern。
- 需要测试的场景。
- 项目特有的坑和规避方法。

AI 在实现或调试中发现新规则时，应更新 `.specflow/spec/`，而不是只在对话中总结。`specflow-update-spec` skill 指导何时将经验提升为 spec。

## 本地定制点

| 需求 | 编辑位置 |
| --- | --- |
| 新增 spec 分类 | `.specflow/spec/<分类>/index.md` 和对应规范文件。 |
| 改 spec 内容 | 直接编辑 `.specflow/spec/` 下的 Markdown 文件。 |
| 安装新模板 | `specflow spec install <分类>`。 |
| 让任务读新 spec | 任务 `implement.jsonl` / `check.jsonl`。 |
| 改 spec 更新时机 | `.specflow/workflow.md` 和 `specflow-update-spec` skill。 |

## 边界

`.specflow/spec/` 是用户项目规范，不是 specflow 内置模板的永久副本。AI 应鼓励用户按实际项目代码更新 spec，而非把默认模板当作不可变文档。
