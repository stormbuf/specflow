# 前端组件规范

定义组件模式、props 约定、样式与无障碍标准。

## Overview

<!-- 引导问题：回答以下问题后再填写具体章节 -->

<!--
- 用什么组件模式？函数组件还是 class 组件？是否使用 composition over inheritance？
- props 如何定义？是否用 TypeScript interface / type？是否有命名约定？
- 如何处理组合？是否使用 children、render props、compound component？
- 无障碍标准是什么？WCAG 级别？哪些组件必须满足哪些 a11y 要求？
- 组件粒度如何划分？何时该拆分？
- 受控组件与非受控组件的使用边界？
- 组件状态放本地还是提升到父级？
-->

（由团队填写：概述组件设计的核心原则与约束。）

## 组件结构

（由团队填写：定义标准组件的文件结构。包括：

- 一个组件文件内的代码组织顺序（import → type → component → export）
- 组件行数上限与拆分时机
- 是否使用 `forwardRef`
- 组件是否拆分为子组件文件还是单文件内多个
）

## Props 约定

（由团队填写：定义 props 的命名与类型约定。包括：

- props 类型用 `interface` 还是 `type`
- 命名规则（如事件回调用 `onXxx`、布尔值用 `isXxx` / `hasXxx`）
- 可选 props 与默认值处理
- children 与 render props 的使用场景
- 是否禁止 spread props（`{...props}`）
）

## 样式模式

（由团队填写：定义组件样式的实现方式。包括：

- 方案选型（CSS Modules / styled-components / Tailwind / CSS-in-JS）
- 类名命名规则（BEM / kebab-case）
- 主题变量与设计 token 的使用
- 响应式断点约定
- 内联样式的使用边界
）

## 无障碍

（由团队填写：定义组件的无障碍要求。包括：

- 交互组件必须支持键盘操作
- 图片必须有 alt
- 表单控件必须有 label 关联
- 语义化 HTML 的使用（button / nav / main / section）
- ARIA 属性的使用场景
- 颜色对比度要求
）

## 常见错误

（由团队填写：列出团队在组件开发上犯过的错误，例如：

- 组件承担过多职责，成为"上帝组件"
- props 命名不一致，同一概念多种叫法
- 用 div + onClick 代替 button，丢失键盘可访问性
- 内联样式硬编码颜色，无法统一换肤
- 组件内部直接调 API，把数据获取与渲染耦合
- 没有处理 loading / error / empty 三态
）
