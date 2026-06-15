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
    列出时必须遍历目录内容获取子目录名；不要用通配符匹配子目录（通配符可能只匹配文件，导致误判为"无现存变更"）。
ELSE IF 目标项目没有 {PROJECT_ROOT}/specflow/changes/<change-id>/ 且需要创建变更:
   使用当前日期生成 YYYY-MM-DD-short-slug-N 格式 change-id
   short-slug 使用 2-5 个英文 kebab-case 单词概括变更主题
     N 从 0 开始，只有同日期、同 short-slug 已存在时才递增（检查已有变更时同样不要用通配符，需遍历目录内容获取子目录名）
   创建 {PROJECT_ROOT}/specflow/changes/<change-id>/
ELSE IF change-id 尚未确定:
   以上各分支均未命中意味着 change-id 仍为空。不得在不经用户确认和不读取 specflow/changes/ 目录的情况下跳过或猜测 change-id。必须明确询问用户要操作的 change-id（列出现有变更供选择，或允许创建新变更）。
ELSE IF 阶段需要读取上游产物但文件缺失:
  暂停并补齐缺失产物，不跳过阶段
ELSE:
  读取当前阶段需要的文件并继续
```

## 阶段产出物审查

Proposal、Spec Delta、System Architecture / ADR、Design、Tasks 五个阶段在输出产物文档后，须启动独立审查-修复-审查循环，最多三轮。

### 审查范围

| 阶段 | 审查产物 | 审查重点 |
|---|---|---|
| Proposal | `proposal.md` | 范围/非范围是否明确、变更类型结论及理由是否充分、是否有无事实来源的断言、阻断性开放问题是否已解决 |
| Spec Delta | `spec-delta.md` | EARS 句式合规、Gherkin Scenario Given/When/Then 步骤可验证性、与 proposal 范围一致性、归属标注完整性、开放问题是否已分类 |
| System Architecture / ADR | `architecture.md` / `adr/*.md` | Mermaid UML 正确性、ADR 逐题确认记录、方案达成与放弃理由、技术栈评估完整性 |
| Design | `design.md` | 族插槽选择与 proposal 变更类型一致、各族维度覆盖完整性、§0 上下文完整、§9 架构影响与 architecture.md / ADR 一致、§9 分族检查清单逐项检查 |
| Tasks | `tasks.md` | 覆盖 spec-delta 和 design 约束（由子项 covers 字段核对）、任务顺序依赖正确、section 启用决策正确、验证任务存在 |

### 审查流程

```text
1. 审查 — 委托审查 agent 读取产物文档，逐项检查审查重点，输出问题清单（如有）
2. 修复 — 按问题清单逐项修复产物文档
3. 再审查 — 审查 agent 再次检查，确认问题已修复或发现新问题
4. 最多迭代 3 轮
IF 3 轮内所有问题修复:
   确认产物质量通过，进入下一阶段
ELSE IF 3 轮后仍存在未解决问题:
   将剩余问题标注为已知问题，写入产物文档末尾的「已知问题」节
   继续进入下一阶段，不无限循环
ELSE IF 审查无问题:
   直接通过
```

### 审查委托

审查时必须委托审查 agent，不得自行审查。审查 agent 的 prompt 须包含：
- 产物文件路径
- 当前阶段名称
- 审查重点清单
- 上一轮问题清单（第 2、3 轮时）

审查 agent 返回结果须区分：无问题 / 有问题（列出逐条问题及严重程度：阻断/建议）。

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

## 阶段进入条件

architecture/ADR gate 等跨阶段前置校验统一在此定义一次，各 `references/` 文件引用本表，不重复：

| 阶段 | 增量进入条件 |
|---|---|
| Tasks | proposal 必需；按 `tasks-template.md`「上游产物依赖表」核对必需产物齐备；若 design.md §9 含架构/ADR 变更，`architecture.md` 已反映且对应 ADR 状态为 accepted；项目缺少 `architecture.md` 且需生成代码或确定实现结构时暂停 |
| Apply | `tasks.md` 存在；不重复 Tasks 的上游校验，仅做一致性核对（proposal/spec-delta/design/architecture 与 tasks.md 声明的任务族一致） |

## 执行纪律

以下原则在 Design、Tasks、Apply 各阶段均适用，作为脑内规则实时遵循（不在 tasks.md 中生成任务项）：

- **思考先行**：写代码/拆任务前先列改动文件清单 + 意图。不在不确定时直接动手。
- **契约显式化**：对外接口、事件、配置、数据结构的语义必须显式；状态值、枚举值、错误语义和边界行为必须有单一真相源；跨模块共享契约必须集中定义；实现必须遵循 design.md 中已确认的接口、数据流、状态流转、UML 图和约束。
- **单一职责**：每个模块/函数职责可用一句话描述，且描述不含"和"字；不为单次使用创建抽象；不引入未要求的"灵活性"或"可配置性"。
- **影响面扫描**：改动前列受影响文件、接口、数据、配置或流程。确认影响面可控后再开始改。
- **错误处理策略**：所有外部输入、外部依赖、IO、网络、并发、状态转换和持久化操作必须考虑失败路径；错误语义显式；不吞掉错误；错误日志不得泄露密钥/令牌/隐私数据；错误处理必须有对应验证。

## 模板

创建文档时按需使用：

- Proposal: `{SKILL_DIR}/assets/proposal.md`
- Roadmap: `{SKILL_DIR}/assets/roadmap.md`
- Spec Delta: `{SKILL_DIR}/assets/spec-delta.md`
- 主 spec（Main Spec）: `{SKILL_DIR}/assets/spec.md`（Archive「落账」首次创建主 spec 时使用）
- Design: `{SKILL_DIR}/assets/design.md`
- Architecture: `{SKILL_DIR}/assets/architecture.md`
- ADR: `{SKILL_DIR}/assets/adr.md`
- Tasks: `{SKILL_DIR}/assets/tasks.md`（索引）→ `{SKILL_DIR}/assets/tasks-template.md`（统一模板，含上游依赖表 + 任务族插槽 + Agent 角色表 + section 启用条件）
- Verification: `{SKILL_DIR}/assets/verification.md`

## 关键约束

1. 不使用外部 schema CLI、Delta、archive 或 schema validate 机制。
2. 只在用户显式调用 `/specflow-workflow` 时激活，禁止自动匹配普通需求描述。
3. `spec-delta.md` 不是主 spec；归档后必须合并到 `{PROJECT_ROOT}/specflow/specs/<capability>.md`。spec-delta 与主 spec 需求块同构，落账时原样合并。
4. 系统架构和 ADR 必须来自 System Architecture / ADR 阶段的逐题用户确认或已有项目事实；Design 不临场创建或修改长期文档。
5. `roadmap.md` 只维护规划台账和完成历史，不承载 proposal、design 或 tasks 细节。
6. 执行阶段只按 `tasks.md` 推进，不重新设计需求或方案。
7. Tasks 阶段是显式编排权威：Agent 角色、串行约束、section 启用条件在拆解时一次性编进 `tasks.md`。Apply 阶段启动时一次性从 `tasks.md` 提取执行契约到工作记忆，不重读外部规则文件，不独立推导调度逻辑。
8. Apply 末尾必须执行验证并产出 verification.md；中断后重新运行 Apply 应跳过已完成任务并完成验证。
9. 所有 UML 必须使用 Mermaid；系统架构只维护系统边界图和系统架构图。
10. System Architecture / ADR 阶段固定先 ADR、后系统架构；系统架构和 ADR 必须逐个问题确认，每个问题给出调研后的 2-3 个方案、推荐方案，并保留用户自述选项；无改动则跳过，不硬改。
11. Proposal、Spec Delta、System Architecture / ADR、Design、Tasks 阶段产出文档后，必须通过独立审查-修复-审查循环（最多三轮），不得未经审查直接进入下一阶段。
12. spec-delta 与主 spec 的需求行为描述必须用 EARS 语法（关键字英文，至少一句 SHALL）；验收标准必须用 Gherkin `Scenario` 围栏块（```gherkin，关键字英文，步骤正文中文）；二者需求块同构。不再使用「系统必须…」散文式描述或 `- Given/- When/- Then` 散列。
