---
name: specflow-workflow
description: "[SLASH-ONLY] specflow-workflow。仅在用户显式调用 /specflow-workflow 时激活；用于稳定执行 Roadmap、Proposal、Spec Delta、System Architecture / ADR、Design、Tasks、Apply（含验证）和 Archive 全流程。"
metadata:
  version: "0.1.0"
  tags: "workflow,roadmap,spec,architecture,design,tasks,apply,archive"
---

# Specflow Workflow

> AI 使用本 skill 执行独立的变更生命周期：版本规划 -> 提议讨论 -> 规约变更 -> 系统架构 / ADR -> 技术方案 -> 任务拆解 -> 执行（含验证） -> 归档。
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
| `/specflow-workflow arch` | System Architecture / ADR |
| `/specflow-workflow architecture` | System Architecture / ADR |
| `/specflow-workflow adr` | System Architecture / ADR |
| `/specflow-workflow design` | Design |
| `/specflow-workflow tasks` | Tasks |
| `/specflow-workflow apply` | Apply |
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

已归档变更存放在：

```text
{PROJECT_ROOT}/specflow/archive/<change-id>/
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

`spec-delta.md` 是本次规约变化；`{PROJECT_ROOT}/specflow/specs/<capability>.md` 是归档后的主 spec。`System Architecture / ADR` 负责系统边界图、系统架构图和 ADR 的首次创建、后续更新、替代或废弃；`design.md` 是本次功能设计工作台，只消费或请求更新系统架构和 ADR。Archive 不创建或修改长期文档。

## 阶段路由

```text
IF 参数 = roadmap OR plan:
  执行 Roadmap 阶段，详见 {SKILL_DIR}/references/roadmap.md
ELSE IF 参数 = pr OR proposal:
  执行 Proposal 阶段，详见 {SKILL_DIR}/references/proposal.md
ELSE IF 参数 = sd OR spec:
  执行 Spec Delta 阶段，详见 {SKILL_DIR}/references/spec-delta.md
ELSE IF 参数 = arch OR architecture OR adr:
  执行 System Architecture / ADR 阶段，详见 {SKILL_DIR}/references/system-architecture-adr.md
ELSE IF 参数 = design:
  执行 Design 阶段，详见 {SKILL_DIR}/references/design.md
ELSE IF 参数 = tasks:
  执行 Tasks 阶段，详见 {SKILL_DIR}/references/tasks.md
ELSE IF 参数 = apply:
  执行 Apply 阶段，详见 {SKILL_DIR}/references/apply.md
ELSE IF 参数 = archive:
  执行 Archive 阶段，详见 {SKILL_DIR}/references/archive.md
ELSE:
  询问用户要进入哪个阶段
```

## 阶段顺序

1. Roadmap：维护规划台账和完成历史；Proposal 不强依赖 Roadmap。
2. Proposal：明确为什么做、范围、非范围、变更类型、是否需要规约/设计。
3. Spec Delta：描述本次对系统可观察行为的新增、修改、删除或重命名。
4. System Architecture / ADR：固定先 ADR、后系统架构；当存在适用 ADR 的长期决策、缺少 architecture.md、系统边界图或系统架构图需要变化时，逐题讨论并按需更新长期文档；无改动则跳过，不硬改。
5. Design：讨论并收敛本次实现方案，只消费系统架构和 ADR；发现缺失或需要长期决策时暂停并进入 System Architecture / ADR。
6. Tasks：将方案拆成可独立验证的执行项。
7. Apply：按 tasks.md 顺序执行并更新复选框；所有任务完成后自动验证并产出 verification.md。
8. Archive：落账（spec-delta → 主 spec）→ 归档（移入历史区）→ 记账（roadmap + 版本管理封存）。

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
   列出 {PROJECT_ROOT}/specflow/changes/ 下所有子目录作为候选 change-id，询问用户选择
   除非当前阶段明确要创建新变更
   注意：列出时确保能匹配子目录——部分 glob 的 * 不匹配子目录，会导致误判为"无现存变更"
ELSE IF 目标项目没有 {PROJECT_ROOT}/specflow/changes/<change-id>/ 且需要创建变更:
   使用当前日期生成 YYYY-MM-DD-short-slug-N 格式 change-id
   short-slug 使用 2-5 个英文 kebab-case 单词概括变更主题
   N 从 0 开始，只有同日期、同 short-slug 已存在时才递增（检查已有变更时同样注意上述子目录匹配陷阱）
   创建 {PROJECT_ROOT}/specflow/changes/<change-id>/
ELSE IF change-id 尚未确定:
   以上各分支均未命中意味着 change-id 仍为空。不得在不经用户确认和不读取 specflow/changes/ 目录的情况下跳过或猜测 change-id。必须明确询问用户要操作的 change-id（列出现有变更供选择，或允许创建新变更）。
ELSE IF 阶段需要读取上游产物但文件缺失:
  暂停并补齐缺失产物，不跳过阶段
ELSE:
  读取当前阶段需要的文件并继续
```

## 阻断条件

各阶段的详细阻断条件定义在对应 `references/` 文件中，以下为编排级规则：

```text
IF 任一阶段产物缺失:
  不跳过，返回补齐缺失产物
ELSE IF 任一阶段执行中触发其 reference 定义的阻断:
  暂停，处理阻断后恢复
ELSE IF 阶段产物间发现冲突:
  暂停，修正冲突文档后继续
```

## 模板

创建文档时按需使用：

- Proposal: `{SKILL_DIR}/assets/proposal.md`
- Roadmap: `{SKILL_DIR}/assets/roadmap.md`
- Spec Delta: `{SKILL_DIR}/assets/spec-delta.md`
- Design: `{SKILL_DIR}/assets/design.md`
- Architecture: `{SKILL_DIR}/assets/architecture.md`
- ADR: `{SKILL_DIR}/assets/adr.md`
- Tasks: `{SKILL_DIR}/assets/tasks.md`（索引）；按变更类型从以下选择：
  - `functional`（功能性）→ `{SKILL_DIR}/assets/tasks-functional.md`
  - `nonfunctional`（非功能）→ `{SKILL_DIR}/assets/tasks-nonfunctional.md`
  - `infrastructure`（基础设施）→ `{SKILL_DIR}/assets/tasks-infra.md`
  - `lightweight`（轻量）→ `{SKILL_DIR}/assets/tasks-lightweight.md`
- Verification: `{SKILL_DIR}/assets/verification.md`

## 关键约束

1. 不使用外部 schema CLI、Delta、archive 或 schema validate 机制。
2. 只在用户显式调用 `/specflow-workflow` 时激活，禁止自动匹配普通需求描述。
3. `spec-delta.md` 不是主 spec；归档后必须合并到 `{PROJECT_ROOT}/specflow/specs/<capability>.md`。
4. 系统架构和 ADR 必须来自 System Architecture / ADR 阶段的逐题用户确认或已有项目事实；Design 不临场创建或修改长期文档。
5. `roadmap.md` 只维护规划台账和完成历史，不承载 proposal、design 或 tasks 细节。
6. 执行阶段只按 `tasks.md` 推进，不重新设计需求或方案。
7. Tasks 和 Apply 阶段每次运行都必须重新读取 `{SKILL_DIR}/assets/rules.md`，动态编排适用检查项，不复制或缓存规则正文。
8. Apply 末尾必须执行验证并产出 verification.md；中断后重新运行 Apply 应跳过已完成任务并完成验证。
9. 所有 UML 必须使用 Mermaid；系统架构只维护系统边界图和系统架构图。
10. System Architecture / ADR 阶段固定先 ADR、后系统架构；系统架构和 ADR 必须逐个问题确认，每个问题给出调研后的 2-3 个方案、推荐方案，并保留用户自述选项；无改动则跳过，不硬改。
