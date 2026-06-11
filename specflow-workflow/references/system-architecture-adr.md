# System Architecture / ADR Stage

本文件定义系统架构与 ADR 阶段。系统架构只维护系统边界图和系统架构图，且所有 UML 必须使用 Mermaid；ADR 只记录不可逆或未来容易忘的长期决策。

路径锚定：本文件中的 `specflow/`、源码、测试、配置和部署路径均相对于 `{PROJECT_ROOT}/`；只有 `{SKILL_DIR}/assets/architecture.md` 和 `{SKILL_DIR}/assets/adr.md` 来自 skill 目录。

## 目标

- 创建或更新 `specflow/architecture.md` 中的两个 Mermaid UML：系统边界图、系统架构图。
- 为不可逆或未来容易忘的长期决策创建或更新 ADR。
- 阻止 Design、Tasks 和 Apply 在缺少系统架构或 ADR 约束时自由发挥。

## 触发

```text
IF {PROJECT_ROOT}/specflow/architecture.md 不存在:
  进入 System Architecture / ADR 阶段
ELSE IF 本次变更影响系统边界图或系统架构图:
  进入 System Architecture / ADR 阶段
ELSE IF 本次变更涉及 ADR 适用范围中的长期决策:
  进入 System Architecture / ADR 阶段
```

## 输入

- 用户原始需求和当前对话中已确认的技术偏好。
- 从 `{PROJECT_ROOT}/specflow/changes/` 子目录确定当前 change-id，读取以下产物（若存在）
- `{PROJECT_ROOT}/specflow/changes/<change-id>/proposal.md`，如果存在。
- `{PROJECT_ROOT}/specflow/changes/<change-id>/spec-delta.md`，如果存在。
- `{PROJECT_ROOT}/specflow/changes/<change-id>/design.md`，如果存在 — 了解 Design 阶段已识别的架构影响和 ADR 候选
- 已有 `specflow/architecture.md` 和 `specflow/adr/`，如果存在。
- 源码、manifest、测试配置、部署配置和项目文档事实。

## 输出

- `specflow/adr/NNNN-short-title.md`，当存在适用 ADR 的长期决策时
- `specflow/architecture.md`，当系统边界图或系统架构图需要创建或更新时

## 固定顺序

```text
先读取 design.md（如存在），提取 §9 中已识别的架构影响和 ADR 候选
将这些与 proposal、spec-delta、用户原话合并，形成问题队列
先处理 ADR 问题队列
IF 本次没有 ADR 适用范围内的新增、替代或废弃决策:
  记录“ADR 无变动”，跳过 ADR 写入
ELSE:
  按逐题讨论规则确认并写入 ADR
再处理系统架构问题队列
IF 系统边界图和系统架构图均无变动，且 architecture.md 已存在:
  记录“系统架构无变动”，跳过 architecture.md 写入
ELSE:
  按逐题讨论规则确认并写入 architecture.md
```

不得为了满足阶段顺序而硬改 ADR 或系统架构图。无新增决策、无图形变化时必须跳过对应产物写入。

## 系统架构范围

`specflow/architecture.md` 只包含两个 UML，后续架构变动也只更新这两个 UML：

- 系统边界图：当前系统、用户/角色、外部系统、输入输出、信任边界、明确不属于本系统的范围。
- 系统架构图：系统内部主要模块/服务/容器、数据存储、同步/异步调用、关键集成、部署边界和横切能力。

所有 UML 必须使用 Mermaid。不得使用 PlantUML、Graphviz、ASCII 图、图片或其他图形语法。

## ADR 适用范围

只为以下“不可逆”或“未来容易忘”的长期决策创建或更新 ADR：

- 技术栈选型：语言、运行时、框架、数据库、缓存、队列、测试框架、部署平台、云服务、关键 SDK。
- 核心架构约束：例如业务层不能直接访问数据库、前端不做权限判断、跨模块只能走指定入口。
- 可观测性：日志规范、监控选型、链路追踪方案、告警边界。
- 数据所有权与存储策略：数据归属、读写入口、持久化格式、一致性策略、迁移策略。
- 认证/安全方案：认证协议、授权边界、密钥管理、敏感数据处理。
- 关键集成方案：支付、邮件、外部 API、AI 服务或其他高耦合集成。
- AI 协作禁止项：例如禁止 AI 修改数据库 schema、禁止绕过测试、禁止改动生成文件源头之外的产物。

当前实现细节、可轻易回滚的局部代码组织、一次性脚本、普通类名/函数名不创建 ADR。

## 逐题讨论规则

```text
先识别 ADR 问题队列，再识别系统架构问题队列；每个队列内部按阻塞顺序排序
FOR EACH 问题:
  只讨论当前一个问题，不一次性展示后续全部问题
  基于当前项目事实和必要外部调研给出 2-3 个可行方案
  每个方案必须解释成立理由、代价、风险和验证方式
  标明推荐方案和推荐理由
  必须保留“用户自述选项”，允许用户输入自己的方案
  使用 question 工具请求用户确认当前一个问题，并暂停
  IF 用户没有明确确认当前问题:
    不得进入下一个问题
    不得创建或修改 architecture.md
    不得创建或修改 ADR
    不得进入 Design、Tasks 或 Apply
  ELSE:
    记录用户确认原文或确认摘要
    进入下一个问题
```

AI 的推荐方案、推断、草案或沉默通过均不等同于用户确认。不得用一次性总确认替代系统架构或 ADR 的逐题确认。

## 技术栈 ADR 评估

技术栈 ADR 的每个候选方案必须逐项比较：

- 业务与质量驱动匹配。
- 硬约束：已有代码、团队技能、运行环境、许可证、预算、部署限制、合规或供应商限制。
- 生态成熟度：社区活跃度、文档质量、维护频率、长期支持、插件/扩展、已知风险。
- 集成边界：与模块、数据层、构建系统、测试系统、部署系统和运维工具的耦合。
- 运行与运维：资源占用、部署复杂度、监控、日志、调试、故障恢复、升级路径。
- 安全与合规：认证授权、权限边界、数据保护、依赖漏洞、供应链风险、许可证风险。
- 迁移与退出：引入成本、迁移步骤、回滚方案、锁定风险、替换成本、数据迁移影响。
- 验证方式：实验、原型、benchmark、测试、代码审查规则或验收标准。

不得只用“主流”“简单”“性能好”“生态好”等泛化理由确认技术栈 ADR。每个候选方案必须有明确放弃理由；被选方案必须记录负面后果和未来触发重新评估的条件。

## 写入规则

```text
IF 当前问题属于 ADR 适用范围且已逐题确认:
  使用 {SKILL_DIR}/assets/adr.md 创建或更新 {PROJECT_ROOT}/specflow/adr/NNNN-short-title.md
  若替代已有 ADR：新 ADR 的 status 为 accepted，被替代的 ADR status 改为 superseded，
  并在"替代/被替代"中互相引用
IF 所有系统架构问题均已逐题确认且系统边界图或系统架构图需要创建或更新:
  使用 {SKILL_DIR}/assets/architecture.md 创建或更新 {PROJECT_ROOT}/specflow/architecture.md
  若更新已有 architecture.md：
  - 直接编辑文件中对应的 Mermaid 图，增加/修改/删除节点和边
  - 在"用户确认记录"中追加本次确认摘要
  - 在"最近更新"中追加一条：`- <change-id> — <变更简述>`
IF ADR 无变动或系统架构无变动:
  跳过对应长期文档写入，不硬改文件
IF 存在未确认问题:
  暂停，不写入长期文档
```

## 审查

architecture.md 或 ADR 写入后，按 SKILL.md「阶段产出物审查」启动独立审查-修复-审查循环（最多三轮），委托审查 agent 检查：
- Mermaid UML 语法正确性、节点和边的完整性
- ADR 逐题确认记录是否可追溯
- 各方案达成理由与放弃理由是否充分
- 技术栈 ADR 是否逐项评估（业务匹配、硬约束、生态、集成、运维、安全、迁移）

审查通过后方可进入下一阶段。三轮后仍有问题标注为已知问题继续。

## 完成条件

- `specflow/architecture.md` 只包含系统边界图、系统架构图及必要元信息。
- 所有 UML 均为 Mermaid。
- 需要 ADR 的长期决策均已记录。
- 系统架构和 ADR 的用户确认记录均可追溯。
- 后续 Design、Tasks 和 Apply 可直接消费系统架构和 ADR，不需要临场选择全局技术方向。
