---
name: specflow-spec-bootstrap
description: "specflow spec 引导 skill。从真实代码库出发，用单 agent 工作流创建或刷新 .specflow/spec/ 编码规范。适用于分析代码库、分解 spec 工作边界、编写有源码证据支撑的 spec 文档，拒绝占位符文本。"
trigger: "用户要求从代码库生成/初始化 spec 规范"
---

# specflow Spec Bootstrap

使用本 skill 从真实代码库创建或刷新 `.specflow/spec/` 规范。一个 agent 负责完整闭环：分析代码库、确定 spec 边界、编写文档、验证结果。工作流不依赖特定 host、CLI 或 agent 品牌。

## 工作流

1. 确认 specflow 已初始化，检查当前 `.specflow/spec/` 目录树。
2. 用可用的最佳工具分析代码库架构：GitNexus、ABCoder、语言原生工具、直接阅读源码。
3. 仅当代码库确实按 package 和 layer 组织时，才按 package 和 layer 分解 spec 工作。
4. 用项目中的具体 pattern、文件路径、示例和反模式填充或重塑 spec 文件。
5. 验证最终 spec 内部一致，不含模板占位符。

## 参考文档路由

| 需求 | 阅读 |
|------|------|
| 代码库架构分析 | [references/repository-analysis.md](references/repository-analysis.md) |
| Spec 工作分解与任务规划 | [references/spec-task-planning.md](references/spec-task-planning.md) |
| 编写高质量 specflow spec 文件 | [references/spec-writing.md](references/spec-writing.md) |
| GitNexus 和 ABCoder MCP 配置 | [references/mcp-setup.md](references/mcp-setup.md) |

## 核心原则

- 模板只是起点，不是契约。当代码库需要时，删除、重命名、拆分或新增 spec 文件。
- 优先使用有源码支撑的规则，而非通用建议。每条重要规则都应指向真实文件或反复出现的本地 pattern。
- 默认单 agent 执行。可选的 helper agent 是实现细节，不是必需或用户可见的依赖。
- 不要编写平台特定指令，除非目标项目已在该平台上标准化。
- 不要在 `.specflow/spec/` 中留下占位符文本、空标题或复制的样板代码。

## 完成标准

- `.specflow/spec/` 描述项目当前的真实状态。
- 每个相关的 package 或 layer 都有实用的编码指引和真实示例。
- 不适用的模板章节已删除。
- `index.md` 文件与最终 spec 文件集一致。
- 所有必需的配置或分析假设已记录在相关 spec 或任务说明中。
