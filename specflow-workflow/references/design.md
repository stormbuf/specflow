# Design Stage

本文件定义技术方案阶段。Design 是技术讨论、影响发现和方案收敛阶段；它说明如何实现 spec-delta，不重新定义用户需求。

路径锚定：本文件中的 `specflow/`、源码、测试、配置和部署路径均相对于 `{PROJECT_ROOT}/`；不得在 `{SKILL_DIR}/` 下查找 proposal、spec-delta、design 或长期项目文档。

## 目录

- 目标：行 20-27
- 输入 / 输出：行 29-37
- 必要章节：行 39-55
- 架构基线分支：行 57-73
- 讨论循环：行 75-91
- 架构影响判定：行 93-107
- 影响升级：行 109-123
- 长期文档边界：行 125-136
- 图表：行 138-151
- 阻断 / 完成条件：行 153-179

## 目标

- 说明目标和非目标。
- 复核规约影响。
- 发现本次实现是否升级为长期架构影响或长期技术决策。
- 定义接口、数据模型、事件、配置或跨模块契约。
- 记录关键流程、风险、迁移和验证策略。
- 标记是否需要在归档时更新 `specflow/architecture.md` 或 `specflow/adr/`。

## 输入

- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/spec-delta.md`，如果“是否需要规约”为 yes
- 相关主 spec、`specflow/architecture.md`、`specflow/adr/`、源码、测试、配置和部署事实

## 输出

- `specflow/changes/<change-id>/design.md`

## 必要章节

- 背景
- 目标 / 非目标
- 规约影响
- 实现思路
- 方案选项
- 接口 / 数据模型
- 流程
- 架构影响
- 长期文档影响
- ADR 候选
- 技术决策
- 风险 / 权衡
- 迁移计划
- 验证策略
- 开放问题

## 架构基线分支

```text
IF {PROJECT_ROOT}/specflow/architecture.md 不存在 且 {PROJECT_ROOT}/specflow/adr/ 不存在或没有有效 ADR 且本次变更需要生成代码或确定实现结构:
  进入完整技术决策和架构讨论
  至少覆盖语言 / 运行时 / 项目类型、目录结构与模块边界、核心产物存储方式、接口 / 数据模型 / 状态流转、测试策略与质量门禁、依赖与基础设施边界、ADR 候选
  FOR EACH 长期技术取舍:
    给出 2-3 个可行方案
    比较复杂度、风险、维护成本、迁移成本和验证方式
    给出推荐方案和理由
    询问用户确认并暂停
  将已确认的最小架构基线写入 design.md 的“架构影响”“长期文档影响”“ADR 候选”“技术决策”章节
  标记 architecture_baseline_confirmed = yes
  未确认前不得进入 Tasks 或 Apply
```

“需要生成代码或确定实现结构”包括新增源码、选择技术栈、创建目录结构、定义持久化格式、定义对外接口、确定跨模块数据流、引入依赖或基础设施。

## 讨论循环

```text
读取当前 change 目录中的 proposal、spec-delta，以及主 spec、architecture、ADR、源码和配置事实
识别本次实现问题
识别可能升级为长期架构影响或 ADR 的问题
FOR EACH 关键取舍:
  给出 2-3 个可行方案
  比较复杂度、风险、维护成本、迁移成本和验证方式
  IF 决策会影响长期结构或长期技术方向:
    询问用户确认关键取舍并暂停
  ELSE:
    记录推荐方案和理由
将确认结果写入 design.md
IF 无开放问题:
  进入 Tasks 阶段
```

## 架构影响判定

```text
IF 技术内容只是本次实现的内部组织方式:
  写入“实现思路”或“接口 / 数据模型”
  不写入“架构影响”
  architecture_update = no，除非存在其他长期影响
ELSE IF 技术内容成为稳定模块边界、唯一读写入口、跨能力依赖点、对外接口、长期数据流、安全边界、部署边界或后续实现约束:
  写入“架构影响”
  判断是否需要同步 architecture.md 或 ADR
```

内部 service、class、function、helper、目录名或文件名本身不自动构成架构影响。只有当它被确认成长期边界、稳定依赖入口或跨变更约束时，才升级为架构影响。

示例：`BookBibleService` 仅用于本次组装 `BOOK_BIBLE` 时，属于实现思路；如果它被确认为后续创作能力访问 `BOOK_BIBLE` 的唯一入口，则属于架构影响。

## 影响升级

```text
FOR EACH design 过程中发现的技术内容:
  IF 只影响本次实现:
    记录在 design.md
  IF 改变当前系统长期结构、模块边界、数据流、安全边界、部署拓扑或基础设施形态:
    在“长期文档影响”中标记 architecture_update = yes
    说明需要更新 specflow/architecture.md 的内容
  IF 形成长期技术决策且存在取舍成本:
    在“长期文档影响”中标记 adr_needed = yes
    在“ADR 候选”中记录标题、决策点、备选方案、推荐方案和后果
```

长期技术决策包括架构边界、基础设施、技术栈、数据库、队列、缓存、对象存储、部署发布、安全模型、可观测性、构建与工程约束。

## 长期文档边界

```text
IF 是本次怎么做:
  写入 design.md
ELSE IF 是当前系统长期结构或现状:
  在 design.md 标记归档时同步 specflow/architecture.md
ELSE IF 是长期技术决策及其理由:
  在 design.md 标记归档时创建或更新 specflow/adr/NNNN-short-title.md
```

不得绕过 design.md 直接生成 `specflow/architecture.md` 或 ADR。长期文档必须来自 Design 阶段确认过的内容或已有项目事实。

## 图表

```text
IF 跨模块或跨服务调用复杂:
  使用 Mermaid component 或 sequence diagram
ELSE IF 异步流程或消息流复杂:
  使用 Mermaid sequence diagram
ELSE IF 状态流转复杂:
  使用 Mermaid state diagram
ELSE IF 数据模型复杂:
  使用字段表或 ER diagram
ELSE:
  使用文字和表格即可
```

## 阻断

```text
IF proposal 的“是否需要技术方案”为 yes 但没有当前 change 目录中的 design.md:
  暂停并在 {PROJECT_ROOT}/specflow/changes/<change-id>/ 创建 design.md
ELSE IF 缺少 architecture.md 和有效 ADR 且本次变更需要生成代码或确定实现结构，但 design.md 未标记 architecture_baseline_confirmed = yes:
  暂停并进入架构基线分支，与用户完成技术决策和架构讨论
ELSE IF design 与 spec-delta 描述的行为冲突:
  暂停并修正冲突
ELSE IF design 引入长期架构约束但未记录架构影响:
  补齐架构影响
ELSE IF design 标记 architecture_update = yes 但没有说明同步内容:
  暂停并补齐同步内容
ELSE IF design 标记 adr_needed = yes 但没有 ADR 候选:
  暂停并补齐 ADR 候选
ELSE IF design 包含未决关键问题:
  询问用户并暂停
```

## 完成条件

- 当前 change 目录中的 design.md 覆盖 spec-delta 的所有关键行为。
- 如果项目缺少 architecture.md 和有效 ADR，且本次变更需要生成代码或确定实现结构，design.md 已记录并确认最小架构基线。
- 接口、数据、错误路径和失败语义明确。
- 验证策略可执行；工具链缺失时写明人工检查方式。
- 架构影响、长期文档影响和技术决策有明确结论。
- `architecture_update` 与 `adr_needed` 均有 yes/no 结论；yes 时有同步内容或 ADR 候选。
