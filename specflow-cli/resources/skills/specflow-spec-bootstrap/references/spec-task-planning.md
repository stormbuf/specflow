# Spec 工作分解与任务规划

默认使用单 agent 执行模型。Agent 可以创建 specflow task 用于可追溯性，但本 skill 不要求特定平台、CLI 或并行 worker 模型。

## 分解原则

围绕真实的所有权边界创建 spec 工作单元：

- 一个 package：当该 package 有自己的约定时。
- 一个 layer：当同一个 package 有不同的 frontend、backend、CLI、worker 或共享库规则时。
- 一个跨切面 guide：当一个 pattern 跨越多个 package 且不属于任何一个 layer 时。

避免人为分解。小型库通常只需要一次集中的 spec 工作，不需要多个 task。

## Task 形状

当 specflow task 有用时，编写简洁的 PRD，包含以下章节：

```markdown
# 填充 <package-or-layer> specflow Spec

## 目标
为 <scope> 编写项目特定的 `.specflow/spec/` 指引。

## 范围
- Spec 目录：
- 需检查的源码目录：
- 需检查的测试：
- 不在范围内：

## 架构上下文
总结代码库分析中的具体发现。

## 需创建或更新的文件
- `.specflow/spec/.../index.md`
- `.specflow/spec/.../<topic>.md`

## 规则
- 根据 real codebase 调整 spec 文件集。
- 使用真实源码示例并标注文件路径。
- 删除不适用的模板章节。
- 除非 task 明确要求，不修改产品源码。

## 验收标准
- [ ] Spec 包含来自代码库的具体示例和反模式。
- [ ] 无占位符文本残留。
- [ ] Index 文件与最终 spec 文件一致。
- [ ] 规则有源码文件、测试或项目文档支撑。
```

## 可选 Helper Agent

如果 host 支持 subagent，helper 可以检查独立的 package 或执行验证。它们是可选的。主 agent 仍然负责集成和最终质量。

Helper task 必须有清晰的职责边界：

- 只读研究 task：可以检查分配范围内所需的任何源码。
- 写 task：应负责不重叠的 spec 目录。
- 验证 task：应检查占位符清除、链接有效性和一致性。

不要在 skill 中编码 helper agent 名称、vendor 特定命令或平台特定路由。只将必需的工作和验收标准放入 task。
