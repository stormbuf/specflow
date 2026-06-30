# 前端 Hook 规范

定义自定义 hook 的结构、数据获取与副作用管理约定。

## Overview

<!-- 引导问题：回答以下问题后再填写具体章节 -->

<!--
- 自定义 hook 模式是什么？命名规则、返回值结构、职责边界？
- 数据获取模式？是否用 React Query / SWR / 自封装 hook？缓存策略？
- 副作用如何管理？useEffect 的使用边界？何时该用 useMemo / useCallback？
- 依赖数组如何处理？是否有 lint 规则强制 exhaustive-deps？
- hook 之间的组合与复用约定？
- hook 中的错误处理与 loading 状态？
- hook 与状态管理库的关系？
-->

（由团队填写：概述 hook 设计的核心原则与约束。）

## Hook 结构

（由团队填写：定义标准 hook 的结构。包括：

- 命名规则（必须 `use` 开头）
- 返回值结构（对象 vs 元组，何时用哪种）
- 参数设计（配置对象 vs 位置参数）
- 单个 hook 的职责上限
- hook 内部状态与外部状态的边界
）

## 数据获取

（由团队填写：定义数据获取的标准模式。包括：

- 是否统一使用 React Query / SWR 等库
- 自封装数据获取 hook 的命名与结构（如 `useUserList`）
- loading / error / data 三态的处理
- 缓存与失效策略
- 请求去重与竞态处理
- 是否禁止在 useEffect 里直接 fetch
）

## 副作用管理

（由团队填写：定义 useEffect 等副作用 hook 的使用规范。包括：

- useEffect 的使用边界（哪些场景该用，哪些不该用）
- cleanup 函数的要求
- 依赖数组的填写规则
- useEffect 内部的异步处理模式
- 何时用 useLayoutEffect
- 事件监听的注册与销毁约定
）

## 常见错误

（由团队填写：列出团队在 hook 使用上犯过的错误，例如：

- useEffect 依赖数组遗漏，导致闭包陷阱
- 在 useEffect 里直接 fetch 没有处理竞态
- 用 useMemo / useCallback 滥用优化，反而降低可读性
- hook 内部状态没有清理，组件卸载后还在 setState
- 把业务逻辑堆进 hook，导致 hook 过重
- 在条件语句或循环里调用 hook，违反 Rules of Hooks
）
