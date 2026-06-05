# Archive Stage

本文件定义归档阶段。Archive 将 spec-delta 合并回主 spec，根据 design 同步长期架构文档和 ADR，并冻结本次变更记录。

路径锚定：本文件中的 `specflow/`、版本管理和项目规则路径均相对于 `{PROJECT_ROOT}/`；只有 `{SKILL_DIR}/assets/architecture.md` 和 `{SKILL_DIR}/assets/adr.md` 来自 skill 目录。

## 目标

- 将 `specflow/changes/<change-id>/spec-delta.md` 合并到 `specflow/specs/<capability>.md`。
- 当 design 标记长期架构影响时，更新或创建 `specflow/architecture.md`。
- 当 design 标记长期技术决策时，创建或更新 `specflow/adr/NNNN-short-title.md`。
- 保留 `specflow/changes/<change-id>/` 作为历史记录。
- 使用当前项目实际采用的版本管理工具封存归档结果，提交或变更描述必须包含 `change-id: <change-id>`。
- 不使用外部 archive 机制。

## 输入

- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/spec-delta.md`
- `specflow/changes/<change-id>/tasks.md`
- `specflow/changes/<change-id>/verification.md`
- `specflow/changes/<change-id>/design.md`，如果存在
- 目标 `specflow/specs/<capability>.md`
- `specflow/roadmap.md`，如果 proposal.md 记录了 Roadmap 来源

## 长期文档同步

```text
IF design.md 中 architecture_update = yes:
  IF specflow/architecture.md 不存在:
    使用 {SKILL_DIR}/assets/architecture.md 创建 {PROJECT_ROOT}/specflow/architecture.md
  ELSE:
    更新 specflow/architecture.md 的对应章节
  只写入 design 已确认内容和已有项目事实
IF design.md 中 adr_needed = yes:
  IF ADR 候选是新决策:
    使用 {SKILL_DIR}/assets/adr.md 创建 {PROJECT_ROOT}/specflow/adr/NNNN-short-title.md
  ELSE IF ADR 候选修订既有决策:
    更新既有 ADR，或创建新 ADR 并将旧 ADR 状态改为 superseded
  只写入 design 的 ADR 候选和已确认取舍
```

`specflow/architecture.md` 记录当前系统长期结构和现状。`specflow/adr/` 记录长期技术决策及其理由。不得把未经 Design 阶段确认的推断写入长期文档。

职责边界：Design 负责发现、讨论和标记长期影响；Archive 负责创建或维护 `specflow/architecture.md` 与 `specflow/adr/`。

## Roadmap 同步

```text
IF proposal.md 的 Roadmap 来源 = 无:
  不更新 specflow/roadmap.md
ELSE IF specflow/roadmap.md 缺失:
  暂停，询问用户是否创建 roadmap 或跳过同步
ELSE:
  将 Roadmap 来源中的对应条目移入 已完成历史 顶部
  标记为 [x]
  追加完成日期、change-id 和归档摘要
```

只更新本 change 明确记录的来源项；不得整理、补写或改写其他 roadmap 条目。已完成历史除本次归档追加外只读。

## 合并规则

```text
FOR EACH 新增需求:
  添加到目标主 spec
FOR EACH 修改需求:
  用完整更新后的需求块替换主 spec 中同名需求
FOR EACH 删除需求:
  从主 spec 移除需求，并确保 spec-delta 记录原因 / 迁移方案
FOR EACH 重命名需求:
  在主 spec 中更新需求名称，并保持场景内容一致或按 patch 修改
```

## 阻断

```text
IF 当前 change 目录中的 verification.md 缺失:
  暂停并进入 Verify 阶段
ELSE IF verification Result = failed:
  暂停，不归档
ELSE IF 当前 change 目录中的 tasks.md 存在未完成任务且没有 Notes 说明例外:
  暂停，完成任务或记录例外
ELSE IF spec-delta 无法无歧义合并到主 spec:
  暂停，说明冲突并请求用户确认
ELSE IF design.md 标记 architecture_update = yes 但没有可同步内容:
  暂停并返回 Design 阶段补齐
ELSE IF design.md 标记 adr_needed = yes 但没有 ADR 候选:
  暂停并返回 Design 阶段补齐
ELSE IF 需要创建 ADR 但无法确定下一个 ADR 编号:
  暂停，检查 specflow/adr/ 后确定编号
ELSE IF proposal.md 记录 Roadmap 来源但无法在 specflow/roadmap.md 定位对应条目:
  暂停，说明缺失项并询问用户是否跳过同步或手动指定条目
ELSE IF 无法判断当前项目使用的版本管理工具:
  暂停，读取项目规则或询问用户
```

## 版本管理

```text
IF 项目规则明确指定版本管理工具:
  使用项目规则指定的工具
ELSE IF 仓库事实能判断版本管理工具:
  使用仓库事实对应的工具
ELSE:
  询问用户
```

提交或变更描述必须显式说明 `change-id`。示例：

```text
docs: archive workflow change

change-id: 2026-06-06-refine-workflow-skill-0
```

## 完成条件

- 主 spec 反映本次变更后的当前系统行为。
- `specflow/architecture.md` 反映本次变更后的长期架构影响，若 design 标记需要同步。
- `specflow/adr/` 记录本次确认的长期技术决策，若 design 标记需要 ADR。
- 变更目录保留为历史记录，不再修改。
- 如果 proposal.md 记录 Roadmap 来源，`specflow/roadmap.md` 已将对应项移入已完成历史。
- 归档结果已用当前项目版本管理工具封存，提交或变更描述包含 `change-id: <change-id>`。
- 最终回复说明归档结果、更新的主 spec、验证摘要。
