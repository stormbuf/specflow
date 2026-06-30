# REST 约定

定义 URL 命名、HTTP 方法语义、请求/响应格式与分页排序的统一约定。

## Overview

<!-- 引导问题：回答以下问题后再填写具体章节 -->

<!--
- URL 命名约定是什么（复数 / 单数、嵌套层级、kebab-case）？
- HTTP 方法的语义如何使用（GET / POST / PUT / PATCH / DELETE）？
- 请求 / 响应格式是什么？JSON 字段命名用 camelCase 还是 snake_case？
- 分页、过滤、排序的参数约定是什么？
- 状态码的使用规范是什么？
- 是否有统一的请求 / 响应包装格式？
-->

（由团队填写：概述 REST API 设计的核心原则。）

## URL 命名约定

（由团队填写：定义 URL 的命名规则。包括：

- 资源命名：复数还是单数（如 `/users` vs `/user`）
- 路径风格：kebab-case / snake_case / camelCase
- 嵌套资源的层级限制（如 `/users/:id/posts/:postId` 最多几层）
- 子资源 vs 独立资源的选择标准
- 查询参数的命名风格
- ID 格式（数字 / UUID / slug）
- 代码示例：推荐 vs 不推荐的 URL 设计
）

## HTTP 方法语义

（由团队填写：定义 HTTP 方法的使用规范。包括：

- GET：读取资源（必须幂等、无副作用）
- POST：创建资源 / 触发动作
- PUT：整体替换资源（幂等）
- PATCH：部分更新资源
- DELETE：删除资源（幂等）
- 是否使用 POST 执行非 CRUD 动作（如 `/users/:id/activate`）
- 幂等性要求与实现方式
- 代码示例：各方法的标准用法
）

## 请求与响应格式

（由团队填写：定义请求和响应的数据格式。包括：

- Content-Type 约定（application/json）
- JSON 字段命名风格（camelCase / snake_case）
- 时间格式（ISO 8601 / Unix timestamp）
- 空值表示（null / 省略字段 / 空字符串）
- 枚举值格式（字符串 / 数字）
- 布尔值字段命名（isXxx / hasXxx / xxxEnabled）
- 是否有统一的响应包装格式（如 `{ code, message, data }`）
- 代码示例：标准请求与响应结构
）

## 分页、过滤与排序

（由团队填写：定义列表资源的查询参数约定。包括：

### 分页
- 分页方式（offset/limit / cursor / page/pageSize）
- 默认值与最大值限制
- 分页响应中的元数据（total / hasNext / cursor）

### 过滤
- 过滤参数命名规则（如 `?status=active`）
- 多条件过滤（AND / OR）
- 范围查询（如 `?created_after=2024-01-01`）
- 模糊搜索参数命名

### 排序
- 排序参数格式（如 `?sort=-created_at,name`）
- 默认排序
- 多字段排序
- 代码示例：分页 / 过滤 / 排序的标准用法
）

## 状态码使用

（由团队填写：定义 HTTP 状态码的使用规范。包括：

- 2xx：200 OK / 201 Created / 204 No Content 的使用场景
- 3xx：是否使用重定向
- 4xx：400 / 401 / 403 / 404 / 409 / 422 / 429 的使用场景
- 5xx：500 / 502 / 503 的使用场景
- 是否使用非标准状态码
- 状态码与业务逻辑的映射规则
- 代码示例：各场景的状态码选择
）

## 常见错误

（由团队填写：列出团队在 REST 设计上犯过的错误，例如：

- URL 用动词而非名词（如 `/getUser` 而非 `GET /users/:id`）
- GET 请求修改数据，违反幂等性
- 用 POST 执行所有操作，未正确使用 PUT / PATCH / DELETE
- 状态码滥用（如用 200 返回错误信息）
- 分页参数不统一，不同接口用不同分页方式
- JSON 字段命名混乱，同一 API 中混用 camelCase 和 snake_case
- 嵌套层级过深（如 `/users/:id/posts/:postId/comments/:commentId/replies/:replyId`）
- 时间格式不统一，有的用 ISO 8601 有的用 timestamp
- 列表接口不返回总数，前端无法显示分页信息
）
