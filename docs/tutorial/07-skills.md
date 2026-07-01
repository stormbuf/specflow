# Skill 详解

## 规划阶段：specflow-brainstorm

用户同意建任务后激活。它的核心职责是把模糊的需求变成精确的 spec：

1. 理解用户需求意图，识别变更类型（新功能/重构/修复/优化）
2. 搜索代码库了解现状，记录关键发现和影响面
3. 起草 prd.md——用 EARS 语法描述行为需求，用 Gherkin 场景定义验收标准
4. 编写 implement.md——把实现拆分为可独立验证的行为切片
5. （可选）编写 design.md——复杂任务的技术设计

!!! example "EARS 语法"
    EARS（Easy Approach to Requirements Syntax）用结构化的句式描述需求。关键字使用英文：WHEN / IF / THEN / WHILE / 至少一句 SHALL。例如：

    ```
    WHEN 用户点击导出按钮 THEN THE SYSTEM SHALL 生成 CSV 文件并下载
    ```

## 编码前：specflow-before-dev

写代码之前自动激活。读取 `.specflow/spec/` 规范库，按 Pre-Development Checklist 检查，将相关领域的编码规范纳入 AI 的工作记忆。如果发现 spec 与 prd 存在冲突，停止并报告，不自行裁决。

## 验收：specflow-check

实现完成后激活。以 prd.md 的验收标准为 **唯一判据**，逐条验证。发现 finding 时可直接修复，自修复循环上限 3 轮。3 轮后仍未通过的，停止并报告未通过项清单。

## 经验沉淀：specflow-update-spec

完成任务后，如果发现值得持久化的编码规范或踩坑经验，自动触发。将经验固化为 spec 文件，更新 index.md 索引。spec 内容是 **团队规范，不是任务记录**——只写可复用的规则，不写一次性实现细节。

## 根因分析：specflow-break-loop

当反复调试同一类 bug 时触发。执行 5 维根因分析（What/Why/When/Where/How），区分表象与根因，设计预防措施，将结论沉淀到 spec 库。

## 跨会话记忆：specflow-session-insight

当用户问"上次怎么解的"、"之前讨论过吗"时触发。调用 `specflow mem` 检索历史对话日志，返回原始对话内容。这是一个能力型 skill，不是强制流程——AI 根据判断决定是否检索以及如何处理返回内容。

```bash
specflow mem search "导出功能"             # 按关键词检索
specflow mem search "导出" --phase brainstorm  # 只搜规划阶段
specflow mem search "导出" --limit 20         # 返回 20 条
specflow mem context "导出"                   # 检索并输出上下文片段
specflow mem list                          # 列出可检索的会话
```

## 工作流推进：/specflow:continue

当你不确定当前进度或下一步时，执行这个命令。AI 读取 task.json status 和 workflow.md 面包屑，判断当前 phase/step，告诉你下一步该做什么。

## 归档：/specflow:finish-work

完成任务后执行。归档任务（状态置为 completed、移动到 archive/、触发 VCS auto-commit、清 session 指针），写 session journal 条目。
