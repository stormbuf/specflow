# Mock 策略

定义 mock 的使用边界、最小化原则与外部依赖处理方式。

## Overview

<!-- 引导问题：回答以下问题后再填写具体章节 -->

<!--
- 什么应该 mock？什么不应该？
- mock 的最小化原则是什么？
- 外部依赖（数据库、API、文件系统）如何处理？
- mock 是在测试中内联定义还是使用独立 mock 文件？
- 如何确保 mock 行为与真实实现一致？
- 何时用 stub / spy / fake / dummy？
-->

（由团队填写：概述 mock 策略的核心原则。）

## Mock 的使用边界

（由团队填写：明确什么应该 mock、什么不应该。包括：

- 必须 mock 的场景：
  - 外部 HTTP API 调用
  - 发送邮件 / 短信 / 推送
  - 第三方支付 / 认证服务
  - 时间相关依赖（时钟）
- 不应该 mock 的场景：
  - 被测代码本身的内部逻辑
  - 纯函数 / 工具函数
  - 项目内的 repository / service 层（用真实实现或 fake）
  - 标准库函数
- 灰色地带的处理方式
）

## 最小化 Mock 原则

**mock 越少越好。** 每多 mock 一个东西，测试就离真实行为更远一步。

以下是最小化 mock 的指导原则：

1. **只 mock 你不拥有的东西** — 外部 API、第三方服务用 mock；项目内部的模块尽量用真实实现。
2. **mock 行为，不 mock 实现** — mock 对外的输入输出，不要 mock 内部的调用路径。如果重构实现，mock 不应该需要改。
3. **mock 粒度要粗** — mock 整个外部服务接口，不要 mock 单个方法再逐个拼装。
4. **fake 优于 mock** — 如果能提供一个轻量的 fake 实现（如内存数据库替代真实数据库），优先使用 fake。fake 可复用、可维护，mock 容易碎。

（由团队填写：补充团队对 mock 最小化的具体约定。）

## 外部依赖处理

（由团队填写：定义各类外部依赖在测试中的处理方式。包括：

### 数据库
- 单元测试：使用 fake / in-memory 还是 mock？
- 集成测试：使用真实数据库还是 testcontainer？
- 数据准备：fixture / factory / migration？
- 数据清理：truncate / transaction rollback / drop？

### HTTP API
- mock 方式：mock 库（nock / responses / httpmock）还是手写 mock server？
- mock 响应数据存放位置？
- 如何处理认证 / 重试逻辑的 mock？

### 文件系统
- 使用 tmpdir 还是 mock fs？
- fixture 文件的存放与管理？

### 消息队列 / 事件
- 如何 mock 消息发布 / 消费？
- 是否测试消息的顺序保证？
）

## Mock 的维护与一致性

（由团队填写：定义如何保持 mock 与真实实现的一致。包括：

- mock 数据的来源（手写 / 录制真实响应 / 从 OpenAPI 生成）
- 如何检测 mock 过期（接口变更后 mock 未更新）
- mock 的审查标准
- 是否使用 contract test 弥合 mock 与真实实现的差异
）

## 常见错误

（由团队填写：列出团队在 mock 策略上犯过的错误，例如：

- mock 了被测对象本身，测试变成"验证 mock 而非验证逻辑"
- mock 返回的数据结构与真实 API 不一致，测试通过但生产失败
- mock 过于细粒度，每次重构都需要修改大量 mock
- 只 mock happy path，未 mock 错误 / 超时 / 限流场景
- mock 中硬编码 URL / 参数，接口变更后 mock 仍"通过"
- 用 mock 替代了本应用 fake 的数据库，导致 SQL 逻辑未被执行
- 测试中 verify 了过多的内部调用细节，使测试脆弱
）
