# 实战示例

以下用一个"为 Web 应用添加 CSV 导出功能"的完整示例，展示 specflow 的端到端流程。

## 1. 描述需求

```
你：我们的订单管理页面需要加一个导出功能，
    用户可以点按钮把当前筛选的订单导出为 CSV 文件。

AI：这个需求涉及前端按钮、后端接口和数据格式化三个部分的改动，
    建议创建一个 specflow 任务来管理。是否创建？

你：好
```

## 2. 创建任务，进入规划

```bash
specflow task create --title "添加订单 CSV 导出功能" --intent "用户需要导出筛选后的订单为 CSV"
```

AI 加载 `specflow-brainstorm`，搜索代码库了解订单模块结构，起草 prd.md：

```markdown title="prd.md（节选）"
## 需求

WHEN 用户在订单管理页面点击"导出 CSV"按钮
THEN THE SYSTEM SHALL 生成包含当前筛选条件下所有订单的 CSV 文件并触发下载

IF 筛选结果为空
THEN THE SYSTEM SHALL 禁用导出按钮并提示"无数据可导出"

## 验收场景

```gherkin
Feature: 订单 CSV 导出
  Scenario: 正常导出
    Given 订单列表有 50 条数据
    And 用户设置了筛选条件"状态=已完成"
    When 用户点击"导出 CSV"按钮
    Then 系统生成 CSV 文件并触发下载
    And CSV 文件包含 30 条已完成订单
    And CSV 包含列：订单号、创建时间、金额、状态

  Scenario: 空数据导出
    Given 筛选结果为空
    Then 导出按钮被禁用
    And 显示"无数据可导出"提示
```
```

然后编写 implement.md，拆分为行为切片：

```markdown title="implement.md（节选）"
## 行为切片

- [ ] 1. 后端：新增 /api/orders/export 接口
  - 接收筛选参数，查询订单，格式化为 CSV
  - 返回 Content-Type: text/csv
- [ ] 2. 前端：添加导出按钮
  - 绑定点击事件，调用 export 接口
  - 空数据时禁用按钮
- [ ] 3. 集成验证
  - 测试正常导出和空数据场景
```

## 3. 整理上下文，激活任务

```bash
specflow add-context .specflow/changes/2026-07-01-add-csv-export-0 \
  specflow-implement ".specflow/spec/backend/api-conventions.md" "API 设计规范"

specflow add-context .specflow/changes/2026-07-01-add-csv-export-0 \
  specflow-implement "src/api/orders.ts" "现有订单接口"

specflow add-context .specflow/changes/2026-07-01-add-csv-export-0 \
  specflow-check ".specflow/changes/2026-07-01-add-csv-export-0/prd.md" "验收标准"

specflow task start
```

## 4. 执行实现 + 验收

```
AI：任务已激活。派发 specflow-implement 执行实现。
    [插件自动注入 prd.md + implement.md + api-conventions.md + orders.ts]

    → implement agent 按 implement.md 逐步实现
    → check agent 逐条验证 prd.md 验收标准
    → 验收通过，全部行为切片完成

    实现完成。我来同步需求并更新 spec。
    [加载 specflow-sync-requirements，同步 EARS 需求到 spec/requirements/orders.md]
    [加载 specflow-update-spec，把 CSV 导出的编码规范写入 spec]
```

## 5. 归档

```
你：/specflow:finish-work

AI：[检查工作区状态]
    [specflow task archive → 归档到 changes/archive/2026-07/]
    [specflow add-session → 写 journal 条目]
    [VCS auto-commit]

    任务已归档。CSV 导出功能开发完成。
```
