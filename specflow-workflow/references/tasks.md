# Tasks Stage

本文件定义任务拆解阶段。Tasks 将 proposal、spec-delta 和 design 转换为可独立验证的执行清单。

路径锚定：本文件中的 `specflow/`、源码、测试和项目文档路径均相对于 `{PROJECT_ROOT}/`；模板文件 `{SKILL_DIR}/assets/tasks-template.md` 来自 skill 目录。

## 目标

- 每个任务只解决一个行为或契约问题。
- 任务按依赖顺序排列。
- 每个任务有可检查完成条件。
- 不混合无关目标。

## 输入

- `{SKILL_DIR}/assets/tasks-template.md` — 统一任务模板，含上游产物依赖表、任务族插槽、Agent 角色表和 section 启用条件。`{SKILL_DIR}` 指本 skill 安装目录。
- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/spec-delta.md`，如果存在
- `specflow/changes/<change-id>/design.md`，如果存在

## 输出

- `specflow/changes/<change-id>/tasks.md`

## 拆解规则

Tasks 阶段是显式编排权威。Agent 角色标注、串行约束、启用条件在拆解时一次性编进 tasks.md；Apply 阶段只做"读 tasks.md → 委托/降级 → 勾选 → 验证"，不再独立推导调度逻辑。

按以下确定性流程拆解：

```text
1. 读取 proposal.md → 提取变更类型字段（functional / nonfunctional / infrastructure / lightweight）

2. 读取 tasks-template.md「上游产物依赖表」→ 该类型必需哪些上游产物

3. 阻断检查（引用 SKILL.md「阶段进入条件」全局门禁）：
   IF 必需产物缺失:
     暂停并补齐
   IF design.md §9 含架构/ADR 变更但 architecture.md 未反映或 ADR 未 accepted:
     暂停，请先执行 /specflow-workflow arch

4. 查表选任务族：
   - functional → 功能族
   - nonfunctional → 优化族
   - infrastructure → 迁移族
   - lightweight → 轻量族

5. 读取该族声明的拆解依据（design 指定章节）→ 切分子项：
   - 功能族：design §1（模块定位）+ §4（接口设计）；每个对外接口/独立模块 = 一个子项
   - 优化族：design §2（优化方案设计）的各优化点；每个独立优化点 = 一个子项
   - 迁移族：design §2（迁移流程设计）的各阶段 + §3（数据迁移设计）的各迁移点；每个迁移阶段/受影响模块 = 一个子项
   - 轻量族：proposal 范围直接映射为单一任务

6. 装配 tasks.md：
   - 公共前置（规约同步，仅 spec-delta 存在时启用）
   - 选中的任务族区块（按该族拓扑，子项内步骤顺序固定，Agent 标注从模板复制）
   - 公共收尾（质量门禁 + 验证）

7. 每个子项填 covers 字段（强制）：
   - 功能族：标注覆盖的 spec-delta Gherkin Scenario（如 covers: spec-delta §用户注册.已登录用户注册）
   - 优化族：标注 design §2 优化点
   - 迁移族：标注 design §2/§3 迁移点
   - 轻量族：无需 covers

8. 按各 section 的「启用条件清单」决定是否启用：
   IF 满足启用条件:
     保留该 section
   ELSE:
     标 [SKIP] 并写明未命中哪条启用条件
```

**多 Agent 策略已在模板落地。** 模板的「实现 section 定义」含子项组织原则（串行/并行/降级）和 Agent 角色表。拆解时将 Agent 标注从模板复制到每个任务项，不需要读取独立的规则文件。Apply 阶段按 tasks.md 中已编排的字段执行。

## 审查

tasks.md 写入后，按 SKILL.md「阶段产出物审查」启动独立审查-修复-审查循环（最多三轮），委托审查 agent 检查：

- **覆盖完整性**：功能族逐条核对 spec-delta 每条 Gherkin Scenario 是否被某子项 covers；其他族逐条核对 design 指定章节是否被 covers。
- **子项 covers 字段完整性**：每个子项是否填写 covers，指向是否准确。
- **section 启用决策正确性**：每个 [SKIP] 是否真的不满足启用条件；未 SKIP 的 section 是否确实满足启用条件。
- **任务顺序符合依赖关系**：子项内步骤严格串行，顺序符合各族拓扑。
- **验证任务存在**。

审查通过后方可进入 Apply 阶段。三轮后仍有问题标注为已知问题继续。

## 完成条件

- `{PROJECT_ROOT}/specflow/changes/<change-id>/tasks.md` 存在。
- 按上游产物依赖表，该变更类型的必需产物齐备（architecture/ADR gate 由 SKILL.md 全局门禁统一校验，本文件不重复）。
- 每个子项含 covers 字段，指向 spec-delta Scenario 或 design 章节。
- section 启用决策有据可查：启用的 section 满足启用条件，[SKIP] 的 section 写明未命中哪条。
- 任务覆盖 spec-delta 和 design 的关键约束（由 covers 字段机械核对）。
- 任务顺序符合依赖关系。
- 验证任务存在。
