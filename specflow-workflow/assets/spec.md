# Feature: <capability>

<!--
主 spec（Main Spec）模板。
- 一个主 spec 文件对应一个 capability（一组「总是一起变更」的系统行为）。
- 由 Archive「落账」阶段首次创建或更新；不在 Proposal / Spec Delta / Design 阶段直接编写。
- 与 spec-delta 需求块同构，遵守同一套语法约束（EARS + Gherkin），便于落账原样合并。

格式约定（硬约束，与 spec-delta 一致）：
- 需求功能定义必须用 EARS 句式，关键字英文，至少一句 SHALL。
- 验收标准必须用 Gherkin ```gherkin Scenario 围栏块，关键字英文，步骤正文中文。每条需求至少一个 Scenario。
- 主 spec 内 `- 目标主规约：` 归属行可省略（主 spec 本身就是归属）；spec-delta 中必须保留。
-->

> <capability 简介：2-3 句话说明本能力覆盖的系统行为边界。>

## 需求：<名称>

<!-- EARS 句式行为描述，例：
The system SHALL allow users to query their own orders by status. -->

```gherkin
Scenario: <场景名称>
  Given <前置条件>
  When <动作>
  Then <可观察结果>
```

## 需求：<名称>

The <system> SHALL <response>.

```gherkin
Scenario: <场景名称>
  Given <前置条件>
  When <动作>
  Then <可观察结果>
```
