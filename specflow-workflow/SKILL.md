---
name: specflow-workflow
description: "[SLASH-ONLY] specflow-workflow。仅在用户显式调用 /specflow-workflow 时激活；用于稳定执行 Roadmap、Proposal、Spec Delta、Design、Tasks、Apply、Verify 和 Archive 全流程。"
metadata:
  version: "0.1.0"
  tags: "workflow,roadmap,spec,design,tasks,apply,archive"
---

# Specflow Workflow

> AI 使用本 skill 执行独立的变更生命周期：版本规划 -> 提议讨论 -> 规约变更 -> 技术方案 -> 任务拆解 -> 执行 -> 验证 -> 归档。
> 最高信条：阶段产物必须前后承接；执行阶段不得重新发明 proposal、spec-delta 或 design。
>
> 路径约定：`{SKILL_DIR}/` = 本 skill 所在目录，仅用于读取本 skill 自带的 `references/` 和 `assets/`；`{PROJECT_ROOT}/` = 用户调用 slash command 时的目标项目根目录。所有 `specflow/`、源码、测试、配置、版本管理路径均相对于 `{PROJECT_ROOT}/`，不得相对于 `{SKILL_DIR}/` 解析。

## 调用方式

本 skill 只支持显式调用，禁止根据普通需求描述自动触发。

```text
IF 本 skill 已被 slash command 显式加载:
  视为用户已显式调用 /specflow-workflow
  按参数路由到对应阶段
ELSE:
  不根据普通需求描述自动激活本 skill
```

不得因当前可见消息缺少字面量 `/specflow-workflow` 而拒绝执行；slash command 解析后的参数同样有效。

| 命令 | 阶段 |
|---|---|
| `/specflow-workflow roadmap` | Roadmap |
| `/specflow-workflow plan` | Roadmap |
| `/specflow-workflow pr` | Proposal |
| `/specflow-workflow proposal` | Proposal |
| `/specflow-workflow sd` | Spec Delta |
| `/specflow-workflow spec` | Spec Delta |
| `/specflow-workflow design` | Design |
| `/specflow-workflow tasks` | Tasks |
| `/specflow-workflow apply` | Apply |
| `/specflow-workflow verify` | Verify |
| `/specflow-workflow archive` | Archive |
| `/specflow-workflow` | 询问用户选择阶段 |

## 工作区约定

目标项目内每次变更使用独立目录，`<change-id>` 使用日期、短语义名和冲突顺序号：

```text
{PROJECT_ROOT}/specflow/changes/<change-id>/
├── proposal.md
├── spec-delta.md
├── design.md
├── tasks.md
└── verification.md
```

长期架构与决策文档使用：

```text
{PROJECT_ROOT}/specflow/roadmap.md
{PROJECT_ROOT}/specflow/architecture.md
{PROJECT_ROOT}/specflow/adr/NNNN-short-title.md
```

格式：

```text
YYYY-MM-DD-short-slug-N
```

规则：

- `short-slug` 使用 2-5 个英文 kebab-case 单词概括变更主题。
- `N` 从 0 开始，只有同日期、同 `short-slug` 已存在时才递增。
- 同日期但不同 `short-slug` 的 change 各自从 `0` 开始。

示例：

```text
2026-06-06-add-order-export-0
2026-06-06-refine-workflow-skill-0
2026-06-06-refine-workflow-skill-1
```

长期规约使用：

```text
specflow/specs/<capability>.md
```

`spec-delta.md` 是本次规约变化；`{PROJECT_ROOT}/specflow/specs/<capability>.md` 是归档后的主 spec。`design.md` 是本次设计讨论和方案收敛记录；归档时仅将已确认的长期影响提炼到 `{PROJECT_ROOT}/specflow/architecture.md` 或 `{PROJECT_ROOT}/specflow/adr/`。

## 阶段路由

```text
IF 参数 = roadmap OR plan:
  执行 Roadmap 阶段，详见 {SKILL_DIR}/references/roadmap.md
ELSE IF 参数 = pr OR proposal:
  执行 Proposal 阶段，详见 {SKILL_DIR}/references/proposal.md
ELSE IF 参数 = sd OR spec:
  执行 Spec Delta 阶段，详见 {SKILL_DIR}/references/spec-delta.md
ELSE IF 参数 = design:
  执行 Design 阶段，详见 {SKILL_DIR}/references/design.md
ELSE IF 参数 = tasks:
  执行 Tasks 阶段，详见 {SKILL_DIR}/references/tasks.md
ELSE IF 参数 = apply:
  执行 Apply 阶段，详见 {SKILL_DIR}/references/apply.md
ELSE IF 参数 = verify:
  执行 Verify 阶段，详见 {SKILL_DIR}/references/verify.md
ELSE IF 参数 = archive:
  执行 Archive 阶段，详见 {SKILL_DIR}/references/archive.md
ELSE:
  询问用户要进入哪个阶段
```

## 阶段顺序

1. Roadmap：维护规划台账和完成历史；Proposal 不强依赖 Roadmap。
2. Proposal：明确为什么做、范围、非范围、变更类型、是否需要规约/设计。
3. Spec Delta：描述本次对系统可观察行为的新增、修改、删除或重命名。
4. Design：讨论并收敛实现方案，发现长期架构影响或 ADR 候选，记录接口、数据流、风险和验证策略。
5. Tasks：将方案拆成可独立验证的执行项。
6. Apply：按 tasks.md 顺序执行并更新复选框。
7. Verify：检查 tasks 是否有遗漏项没做，并记录归档前验证摘要。
8. Archive：将 spec-delta 合并回主 spec，冻结本次变更记录。

## 前置条件

```text
先确定 PROJECT_ROOT:
  优先使用 slash command 被调用时的工作目录
  如果用户提供绝对 change 路径，则从该路径反推 PROJECT_ROOT
  如果无法确定 PROJECT_ROOT，暂停并询问用户
IF 用户提供 change-id:
  只在 {PROJECT_ROOT}/specflow/changes/<change-id>/ 查找
  如果目录存在，读取该目录下阶段产物并继续
  如果目录不存在，暂停并报告已检查的绝对路径，不创建同名变更
ELSE IF {PROJECT_ROOT}/specflow/changes/ 下已有变更且用户未指定 change-id:
  列出现有 change-id 并询问用户选择，除非当前阶段明确要创建新变更
ELSE IF 目标项目没有 {PROJECT_ROOT}/specflow/changes/<change-id>/ 且需要创建变更:
  使用当前日期生成 YYYY-MM-DD-short-slug-N 格式 change-id
  short-slug 使用 2-5 个英文 kebab-case 单词概括变更主题
  N 从 0 开始，只有同日期、同 short-slug 已存在时才递增
  创建 {PROJECT_ROOT}/specflow/changes/<change-id>/
ELSE IF 阶段需要读取上游产物但文件缺失:
  暂停并补齐缺失产物，不跳过阶段
ELSE:
  读取当前阶段需要的文件并继续
```

## 阻断条件

```text
IF proposal 仍有关键开放问题:
  暂停，使用苏格拉底追问询问用户
ELSE IF proposal 包含没有用户原话、当前对话、主 spec、架构文档、ADR、源码、测试或项目文档依据的新业务动机、范围、非范围或影响面:
  暂停，移入开放问题并使用苏格拉底追问询问用户
ELSE IF spec-delta 改变可观察行为但没有目标 capability:
  暂停，补齐目标 {PROJECT_ROOT}/specflow/specs/<capability>.md 信息
ELSE IF spec-delta 包含没有 proposal、主 spec、源码、测试或项目文档依据的新业务约束:
  暂停，移入开放问题并使用苏格拉底追问询问用户
ELSE IF design 与 spec-delta 冲突:
  暂停，更新 design 或 spec-delta 后继续
ELSE IF tasks 未覆盖 spec-delta 或 design 的关键约束:
  暂停，补齐 tasks
ELSE IF apply 需要偏离已确认设计:
  暂停，说明原因并请求确认
ELSE IF verify 失败且无法自行修复:
  暂停，报告失败项和下一步选择
ELSE IF archive 前 tasks 未完成且未记录例外:
  暂停，完成任务或记录未完成原因
```

## 模板

创建文档时按需使用：

- Proposal: `{SKILL_DIR}/assets/proposal.md`
- Roadmap: `{SKILL_DIR}/assets/roadmap.md`
- Spec Delta: `{SKILL_DIR}/assets/spec-delta.md`
- Design: `{SKILL_DIR}/assets/design.md`
- Architecture: `{SKILL_DIR}/assets/architecture.md`
- ADR: `{SKILL_DIR}/assets/adr.md`
- Tasks: `{SKILL_DIR}/assets/tasks.md`
- Verification: `{SKILL_DIR}/assets/verification.md`

## 关键约束

1. 不使用外部 schema CLI、Delta、archive 或 schema validate 机制。
2. 只在用户显式调用 `/specflow-workflow` 时激活，禁止自动匹配普通需求描述。
3. `spec-delta.md` 不是主 spec；归档后必须合并到 `{PROJECT_ROOT}/specflow/specs/<capability>.md`。
4. `design.md` 是本次设计工作台；长期架构和 ADR 必须来自 Design 阶段确认过的内容或已有项目事实。
5. `roadmap.md` 只维护规划台账和完成历史，不承载 proposal、design 或 tasks 细节。
6. 执行阶段只按 `tasks.md` 推进，不重新设计需求或方案。
7. Tasks 和 Apply 阶段每次运行都必须重新读取 `{SKILL_DIR}/assets/rules.md`，动态编排适用检查项，不复制或缓存规则正文。
8. 验证主要检查 tasks 遗漏和例外；只记录结果摘要：passed / failed / partial、检查项和必要 notes。
9. 发现阶段产物冲突时先暂停修正文档，再继续执行。
