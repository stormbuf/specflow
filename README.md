# Specflow Workflow

本仓库维护一个独立的 `specflow-workflow` skill，用于稳定执行变更工作流，不依赖外部 schema CLI、Delta、archive 或 schema validate 机制。

`.opencode/skills/specflow-workflow` 以符号链接方式指向根目录 `specflow-workflow/`。`.opencode/skills/install-specflow-workflow` 提供项目级安装包装，用于把该 workflow 复制安装到其他 OpenCode 项目。

## 调用方式

该 skill 只支持显式调用，避免被 OpenCode 根据普通需求描述无意触发。

```text
/specflow-workflow roadmap   # Roadmap 阶段
/specflow-workflow plan      # Roadmap 阶段
/specflow-workflow pr        # Proposal 阶段
/specflow-workflow proposal  # Proposal 阶段
/specflow-workflow sd        # Spec Delta 阶段
/specflow-workflow spec      # Spec Delta 阶段
/specflow-workflow design    # Design 阶段
/specflow-workflow tasks     # Tasks 阶段
/specflow-workflow apply     # Apply 阶段（含验证）
/specflow-workflow archive   # Archive 阶段
```

单独输入 `/specflow-workflow` 时，应先询问用户选择阶段。

## 工作流

完整变更流程：

```text
版本规划 -> 提议讨论 -> 规约变更 -> 技术方案 -> 任务拆解 -> 执行（含验证） -> 归档
```

版本规划使用单文件台账：

```text
specflow/roadmap.md
```

`roadmap.md` 使用 `ROADMAP_META.next_f` 和 `ROADMAP_META.next_t` 分配 F/T 编号。新增条目后对应计数递增；删除未执行条目也不复用编号。若旧文件缺少 `ROADMAP_META`，Roadmap 阶段会根据现有 F/T 最大编号补齐。

每次变更使用独立目录，`<change-id>` 使用日期、短语义名和冲突顺序号：

```text
specflow/changes/<change-id>/
├── proposal.md
├── spec-delta.md
├── design.md
├── tasks.md
└── verification.md
```

格式：

```text
YYYY-MM-DD-short-slug-N
```

规则：

- `short-slug` 使用 2-5 个英文 kebab-case 单词概括变更主题。
- `N` 从 0 开始，只有同日期、同 `short-slug` 已存在时才递增。
- 同日期但不同 `short-slug` 的 change 各自从 `0` 开始。

示例中的 `2026-06-06-refine-workflow-skill-0` 是 change-id：

```text
2026-06-06-add-order-export-0
2026-06-06-refine-workflow-skill-0
2026-06-06-refine-workflow-skill-1
```

长期规约使用：

```text
specflow/specs/<capability>.md
```

已归档变更存放在：

```text
specflow/archive/<change-id>/
```

## 阶段定义

- `roadmap.md`：规划台账和完成历史，维护 `🔥 正在进行`、`📋 下一批 (P0)`、`💡 远期 (P1/P2)` 和 `已完成历史`。
- `proposal.md`：提议讨论结果，明确为什么做、范围、非范围、变更类型、是否需要规约和技术方案。
- `spec-delta.md`：本次规约变更，描述对系统可观察行为的新增、修改、删除或重命名。
- `design.md`：本次变更的技术方案，说明实现方式、接口、数据流、风险、迁移和验证策略。
- `tasks.md`：可独立验证的执行清单。
- `verification.md`：验证结果摘要，只记录 `passed | failed | partial`、检查项和必要 notes。

Proposal 不强依赖 Roadmap。若 `roadmap.md` 存在待实现项，Proposal 阶段会先询问是否从 Roadmap 选择本次执行项；选择后写入 `proposal.md` 的 `Roadmap 来源` 字段。Archive 阶段会按该来源将对应条目移入 Roadmap 的已完成历史。

## 项目级安装包装

项目级 skill `install-specflow-workflow` 用于把当前工作区的 `specflow-workflow` 复制安装到另一个 OpenCode 项目。

执行规则：

- 扫描 `~/project/` 下的一级目录。
- 使用 `question` 工具让用户选择目标项目。
- 在目标项目目录执行 `npx skills add "<resolved-source-skill-dir>" -a opencode --copy -y`。
- 不使用 `-g`，只做项目级安装。
- 安装后提醒用户重启目标项目中的 OpenCode。

## 归档规则

归档阶段只合并规约变更：

```text
spec-delta.md -> specflow/specs/<capability>.md
```

归档完成后，必须使用当前项目实际采用的版本管理工具封存结果。具体工具由执行该项目的 AI 根据项目规则和仓库事实判断；提交或变更描述必须包含 `change-id: <change-id>`。

示例：

```text
docs: archive workflow change

change-id: 2026-06-06-refine-workflow-skill-0
```

归档阶段不合并 `design.md`。整个 change 目录归档后移入 `specflow/archive/<change-id>/`，成为历史记录，用于解释当时为什么这样设计、排除了什么、识别了哪些风险。

如果变更产生长期架构影响或长期技术决策，只更新对应长期文档：

```text
architecture.md  # 当前架构事实、模块边界、数据流
adr/*.md         # 长期技术决策和约束
```

## 真相源边界

归档后，各类文档和代码的权威边界如下：

- 源码是实现真相源。
- 测试是可执行行为证据。
- `specflow/specs/` 是对外行为承诺。
- `architecture.md` 和 `adr/` 是长期结构与决策约束。
- `specflow/archive/<change-id>/design.md` 是历史方案记录，不是当前实现正确性的判定依据。

如果后续源码演进导致历史 `design.md` 过期，不回头维护历史 design：

- 行为变了：创建新的 `spec-delta.md`，归档后更新主 spec。
- 架构事实变了：更新 `architecture.md`。
- 长期技术决策变了：新增、更新或废弃 ADR。
- 实现细节变了：由源码和测试体现。
