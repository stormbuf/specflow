# 规约增量：<变更标题>

<!--
格式约定（硬约束）：
- 需求功能定义必须用 EARS 句式，关键字英文。至少一句 SHALL。可选句式：
    Ubiquitous : The <system> SHALL <response>.
    State-driven : WHILE <state> THE <system> SHALL <response>.
    Event-driven : WHEN <trigger> THE <system> SHALL <response>.
    Optional : WHERE <feature is included> THE <system> SHALL <response>.
    Unwanted : IF <condition>, THEN THE <system> SHALL <response>.
- 验收标准必须用 Gherkin ```gherkin Scenario 围栏块，关键字英文（Scenario/Given/When/Then/And/But），步骤正文中文。每条需求至少一个 Scenario。
- 主 spec 与 spec-delta 需求块同构，遵守同一套语法约束。
- 不再使用「系统必须…」散文式行为描述，也不再使用 #### 场景： 下的 - Given/- When/- Then 散列。
-->

## 目标主规约

- specflow/specs/<capability>.md <!-- 例：specflow/specs/order.md -->

## 新增需求

### 需求：<名称>
- 目标主规约：specflow/specs/<capability>.md  <!-- 仅一个目标时可省略；多个时每条需求必须显式标注 -->

<!-- EARS 句式行为描述，例：WHEN 用户提交订单 THE system SHALL 在 5 秒内向注册邮箱发送确认邮件。 -->

```gherkin
Scenario: <场景名称>
  Given <前置条件> <!-- 例：用户已登录且购物车中有商品 -->
  When <动作> <!-- 例：用户点击"提交订单"按钮 -->
  Then <可观察结果> <!-- 例：系统创建订单，订单状态为"待支付" -->
```

## 修改需求

<!-- 包含完整更新后的需求块（含 - 目标主规约 标注、EARS 行、Gherkin Scenario 块）。
从主 spec 复制整个需求块，修改以反映新行为。
例：
### 需求：订单列表查询
- 目标主规约：specflow/specs/order.md
（以下为从 specflow/specs/order.md 复制的完整块，已更新...）
The system SHALL ...
```gherkin
Scenario: ...
  Given ...
  When ...
  Then ...
```
... -->

## 删除需求

<!-- 每个删除需求必须包含"原因"和"迁移方案"。
例：
### 需求：旧版 PDF 导出
- 原因：依赖已废弃的 PDF 渲染库。
- 迁移方案：用户迁移至新版 CSV 导出。 -->

## 关键约束

<!-- 跨场景的系统级硬边界。只描述用户或外部系统可观察的行为约束，不写实现细节
  （索引名、类名、文件路径、缓存 key 等）。Design 阶段逐条映射到具体实现方案，
  Apply 阶段逐条核对。

  约束类型：数量限制 | 时间限制 | 权限约束 | 错误语义 | 数据边界 | 并发边界 | 其他
  无跨场景硬边界时写「无」。 -->

| 约束 | 类型 | 违反时行为 | 关联需求 |
|------|------|-----------|---------|
| 单次导出最多 5000 条 | 数量限制 | 返回错误码，提示"导出数量不能超过 5000 条" | 订单导出 |
| 时间范围不能超过 31 天 | 时间限制 | 返回错误码，提示"导出时间范围不能超过 31 天" | 订单导出 |

## 重命名需求

<!-- 每个重命名需求使用"原名称 / 新名称"。
例：
### 重命名：订单导出 / 订单数据导出 -->
