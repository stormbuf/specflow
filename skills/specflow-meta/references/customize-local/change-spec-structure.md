# 改 Spec 结构

用户想改 AI 遵循的工程约定、加新 spec 分层或调整分类时，编辑 `.specflow/spec/`。

## 先读这些文件

1. `.specflow/config.yaml`
2. `.specflow/spec/` 当前结构
3. `.specflow/workflow.md` 规划产物指引
4. 当前任务 `implement.jsonl` / `check.jsonl`

## 内置分类

specflow 内置 9 个 spec 分类模板：

| 分类 | 说明 |
| --- | --- |
| `guides` | 跨层思维指引 |
| `backend` | 后端规范 |
| `frontend` | 前端规范 |
| `architecture` | 架构决策记录 |
| `testing` | 测试规范 |
| `security` | 安全规范 |
| `api` | API 规范 |
| `devops` | DevOps 规范 |
| `git-conventions` | Git 约定 |

## 常见需求

| 需求 | 编辑位置 |
| --- | --- |
| 加 backend/frontend/docs/test spec 分层 | `.specflow/spec/<分类>/` |
| 加共享思维指引 | `.specflow/spec/guides/` |
| 安装新分类模板 | `specflow spec install <分类>` |
| 让任务读新 spec | 任务 `implement.jsonl` / `check.jsonl` |
| 从代码库自动生成 spec | `specflow-spec-bootstrap` skill |

## 加 spec 分层

手动创建：

```text
.specflow/spec/security/
├── index.md
└── auth.md
```

或通过模板安装：

```bash
specflow spec install security
```

`index.md` 应包含：

- 该分层适用哪些代码。
- Pre-Development Checklist。
- Quality Check。
- 指向具体规范文件的链接。

## 更新上下文

加 spec 不意味着每个任务自动读它。当前任务必须在 jsonl 中引用：

```bash
specflow task add-context <task> implement ".specflow/spec/security/index.md" "安全约定"
specflow task add-context <task> check ".specflow/spec/security/index.md" "安全审查规则"
```

## 从代码库生成 spec

`specflow-spec-bootstrap` skill 从真实代码库出发，分析架构、确定 spec 边界、编写有源码证据支撑的 spec 文档。拒绝占位符文本，要求每条规则指向真实文件或反复出现的本地 pattern。

## spec 内容应包含什么

项目可执行的工程约定，而非通用最佳实践：

- 文件应放在哪里。
- 错误处理如何表达。
- API 输入输出契约。
- 禁止的 pattern。
- 需要测试的场景。
- 项目特有的坑和规避方法。

## 注意

- spec 是用户项目约定，可按项目需求改。
- 不把临时任务信息放进 spec；临时信息放任务目录。
- 不把长期约定只放 agent 或 skill；保存到 spec。
- 改 spec 结构后，检查现有任务 jsonl 是否仍指向存在的文件。
