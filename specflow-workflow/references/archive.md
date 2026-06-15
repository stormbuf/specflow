# Archive Stage

本文件定义归档阶段。Archive 包含三个子阶段，顺序执行：**落账 → 归档 → 记账**。
前一阶段阻断则整条链路中断。

路径锚定：本文件中的 `specflow/`、版本管理和项目规则路径均相对于 `{PROJECT_ROOT}/`。

---

## 1. 落账（Settlement）

> spec-delta 的改动应用到主 spec，主 spec 成为系统当前行为的唯一真相。
> 没有主 spec 则新建，移除功能则删除主 spec。

### 输入

- `specflow/changes/<change-id>/spec-delta.md`
- 目标 `specflow/specs/<capability>.md`（capability 名称取自 spec-delta 中"目标主规约"节）

### 规则

```text
IF specflow/specs/<capability>.md 不存在:
   根据 spec-delta 中"目标主规约"节确定 capability 名称
   用 {SKILL_DIR}/assets/spec.md 模板创建该文件（# Feature: <capability> 顶层 + 同构需求块）
   用"新增需求"填充需求块
   IF spec-delta 中存在修改需求 OR 删除需求 OR 重命名需求:
       阻断 — 主 spec 缺失但 spec-delta 声称要修改/删除/重命名，数据不一致
       暂停，询问用户：
         a) spec-delta 写错了 → 回到 spec-delta 阶段修正
         b) 主 spec 被误删了 → 恢复主 spec 后继续合并
ELSE:
   从 spec-delta 中筛选"目标主规约"标注为当前 capability 的需求（若 spec-delta 仅声明一个 capability，所有未标注的需求默认归属于它）
   IF 无归属当前 capability 的需求:
     跳过（该 capability 本次无变更）
   FOR EACH 筛选后的新增需求:
     添加到目标主 spec，原样保留 EARS 行与 ```gherkin Scenario 块
   FOR EACH 筛选后的修改需求:
     用完整更新后的需求块替换主 spec 中同名需求，保留 EARS 句式与 Scenario 块完整
   FOR EACH 筛选后的删除需求:
     从主 spec 移除需求，需确保 spec-delta 记录了原因/迁移方案
   FOR EACH 筛选后的重命名需求:
     在主 spec 中更新需求名称，保持 EARS 需求句式与 Gherkin Scenario 块一致
   若执行完删除后主 spec 无剩余需求:
     删除该主 spec 文件
```

> spec-delta 与主 spec 需求块同构（EARS 行 + ```gherkin Scenario 块）；落账原样合并，不重写需求措辞或 Scenario 结构。
> 主 spec 内 `- 目标主规约：` 归属行可省略（主 spec 本身即归属）；spec-delta 中保留。
> 若 spec-delta 中"目标主规约"节记录多个 capability，每个 capability 对应一个独立的 `# Feature:` 主 spec 文件，分别执行以上规则。

### 阻断

```text
IF spec-delta 中"目标主规约"节缺失或无法定位 capability 名称:
   暂停询问
IF spec-delta 无法无歧义合并到主 spec:
   暂停，说明冲突并请求用户确认
IF 主 spec 不存在且 spec-delta 含修改/删除/重命名需求:
   暂停，数据不一致，询问用户修正方向
```

### 完成条件

- 主 spec 反映本次变更后的系统行为
- 若功能被完全移除，对应主 spec 文件已删除

---

## 2. 归档（Archiving）

> 变更目录从活跃区移入历史区，封存为只读记录。

### 输入

- `specflow/changes/<change-id>/`（完整的变更工作目录）

### 规则

```text
IF specflow/archive/ 不存在:
   创建 specflow/archive/
移动 specflow/changes/<change-id>/ → specflow/archive/<change-id>/
```

### 阻断

```text
IF 当前 change 目录中 verification.md 缺失:
   暂停，重新执行 Apply 阶段
IF verification Result = failed:
   暂停，不归档
IF verification Result = partial:
   暂停，告知用户 verification.md 中已记录的风险，询问是否继续归档
IF 当前 change 目录中 tasks.md 存在未完成任务且没有 Notes 说明例外:
   暂停，完成任务或记录例外
```

### 完成条件

- 变更目录已移入 `specflow/archive/<change-id>/`
- 该 change 不再出现在 `specflow/changes/` 活跃变更列表中
- 归档目录内容为只读历史记录，后续不得修改

---

## 3. 记账（Bookkeeping）

> 关联的 roadmap 条目标记完成，版本管理提交封存（附带 change-id）。

### 输入

- `specflow/changes/<change-id>/proposal.md`（Roadmap 来源字段）
- `specflow/roadmap.md`（若 proposal.md 记录了 Roadmap 来源）

### 规则

#### 3a. Roadmap 同步

```text
IF proposal.md 的 Roadmap 来源 = 无:
   不更新 specflow/roadmap.md
IF specflow/roadmap.md 缺失:
   暂停，询问用户是否创建 roadmap 或跳过同步
ELSE:
   将 Roadmap 来源中的对应条目移入"已完成历史"顶部
   标记 [x]，追加完成日期、change-id 和归档摘要
   只更新本 change 明确记录的来源项；不得整理、补写或改写其他 roadmap 条目
   已完成历史除本次归档追加外只读
```

#### 3b. 版本管理封存

```text
使用当前项目实际采用的版本管理工具封存归档结果
提交或变更描述必须包含 `change-id: <change-id>`
示例：
  docs: archive workflow change

  change-id: 2026-06-06-refine-workflow-skill-0
```

### 阻断

```text
IF proposal.md 记录 Roadmap 来源但无法在 specflow/roadmap.md 定位对应条目:
   暂停，说明缺失项，询问用户跳过同步或手动指定条目
IF 无法判断当前项目使用的版本管理工具:
   暂停，读取项目规则或询问用户
```

### 完成条件

- 若有关联 Roadmap 来源，`specflow/roadmap.md` 对应条目已移入完成历史
- 版本管理已封存，提交描述包含 `change-id: <change-id>`

---

## 总完成条件

1. 主 spec 反映本次变更后的系统行为（落账）
2. 变更目录已移入 `specflow/archive/<change-id>/`（归档）
3. 若有关联，roadmap 条目已标记完成（记账）
4. 版本管理已封存，附 change-id（记账）
5. 最终回复说明归档结果、更新的主 spec、验证摘要

## 不做什么

- 不创建或修改长期文档（architecture.md / ADR）
- 不合并 `design.md`（design 随 change 目录归档，作为历史方案记录）
- 不使用外部 archive 机制
- 不整理/补写/改写其他 roadmap 条目
