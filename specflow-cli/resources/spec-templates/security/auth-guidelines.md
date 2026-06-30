# 认证与授权

定义认证机制、授权模型、密码存储与 Session/Token 管理的统一约定。

## Overview

<!-- 引导问题：回答以下问题后再填写具体章节 -->

<!--
- 认证机制是什么（JWT / Session / OAuth / SSO）？
- 授权模型是什么（RBAC / ABAC / 自定义）？
- 密码如何存储（bcrypt / argon2 / scrypt）？盐值策略？
- Session / Token 如何管理？过期策略？刷新机制？
- 权限检查在哪一层执行（middleware / handler / service）？
- 是否有第三方认证集成（OAuth2 / OIDC / SAML）？
-->

（由团队填写：概述认证与授权的整体架构与核心原则。）

## 认证机制

（由团队填写：定义项目的认证方式。包括：

- 认证方式选型（JWT / Session Cookie / OAuth2 / OIDC）
- Token / Session 的生成、签名与验证流程
- Token 的存储位置（HttpOnly Cookie / localStorage / 内存）
- Token 刷新机制（refresh token 轮换策略）
- 多设备登录 / 登出的处理
- 第三方登录集成（Google / GitHub / 企业 SSO）
- 代码示例：认证 middleware 的标准实现
）

## 授权模型

（由团队填写：定义项目的授权方式。包括：

- 授权模型选型（RBAC / ABAC / ACL / 自定义）
- 角色 / 权限的定义与分配方式
- 权限检查的执行位置（middleware / handler / service / 数据库层）
- 资源级权限控制（如"只能编辑自己的文章"）
- 权限的缓存策略
- 权限变更的生效时机
- 代码示例：权限检查的标准模式
）

## 密码存储

（由团队填写：定义密码的存储与验证策略。包括：

- 哈希算法（bcrypt / argon2 / scrypt）及参数配置
- 盐值策略（自动盐 / 手动盐）
- 密码复杂度要求
- 密码重置流程
- 登录失败锁定 / 限流策略
- 代码示例：密码哈希与验证
）

## Session / Token 管理

（由团队填写：定义 Session 和 Token 的生命周期管理。包括：

- 过期时间设置（access token / session / refresh token）
- 续期 / 刷新机制
- 吊销 / 注销机制（黑名单 / 版本号）
- 并发会话控制
- CSRF 防护（如使用 Cookie 存储 Token）
- 代码示例：Token 生命周期管理
）

## 常见错误

（由团队填写：列出团队在认证授权上犯过的错误，例如：

- JWT secret 硬编码在代码中或使用弱密钥
- Token 不过期或过期时间过长
- 权限检查只在前端做，后端未校验
- 密码用 MD5 / SHA1 哈希，未使用慢哈希算法
- 登录接口无限流，可被暴力破解
- refresh token 不轮换，泄露后可永久使用
- Session 固定攻击未防护（登录后未重新生成 Session ID）
- 权限检查遗漏某些 API 端点，导致越权访问
- Token 存储在 localStorage，存在 XSS 窃取风险
）
