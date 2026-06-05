# Roadmap Stage

本文件定义版本规划阶段。Roadmap 只维护 `specflow/roadmap.md` 工作台账，不替代 proposal、design 或 tasks。

路径锚定：本文件中的 `specflow/` 路径均相对于 `{PROJECT_ROOT}/`；只有 `{SKILL_DIR}/assets/roadmap.md` 来自 skill 目录。

## 目标

- 记录待做功能和纯技术变更。
- 维护优先级分区和正在进行项。
- 保存已完成历史。
- 保持规划条目简洁，不写实现细节。

## 输入

- 用户关于规划、优先级、开始做、完成项的指令。
- `specflow/roadmap.md`，如果存在。

## 输出

- `specflow/roadmap.md`

## 初始化

```text
IF specflow/roadmap.md 不存在:
  使用 {SKILL_DIR}/assets/roadmap.md 创建 {PROJECT_ROOT}/specflow/roadmap.md
ELSE IF ROADMAP_META 缺失:
  根据全文件现有 F/T 最大编号补齐 next_f / next_t
ELSE:
  读取并按现有分区更新，同时校验 ROADMAP_META
```

## 条目类型

```text
F = 业务功能项，包含该功能所需的技术实现、数据结构、流程改动和验证要求
T = 纯技术变更项，只在没有直接业务功能载体时使用

IF 技术工作服务于某个 F 项:
  归入该 F 项，不单独创建 T
ELSE IF 技术变更本身是规划目标:
  创建 T 项
```

F 项格式：

```markdown
- [ ] F{N}: {描述} [{新增|修改|删除|重设计}]
```

T 项格式：

```markdown
- [ ] T{N}: {描述} [技术变更] — {影响范围}
```

## 编号规则

```text
IF ROADMAP_META.next_f <= 全文件现有 F 最大编号:
  修正 next_f = 最大 F 编号 + 1
IF ROADMAP_META.next_t <= 全文件现有 T 最大编号:
  修正 next_t = 最大 T 编号 + 1
新增 F 项:
  N = ROADMAP_META.next_f
  写入 F{N}
  ROADMAP_META.next_f += 1
新增 T 项:
  N = ROADMAP_META.next_t
  写入 T{N}
  ROADMAP_META.next_t += 1
移动、完成或归档条目:
  保持原编号
删除未执行条目:
  不复用编号，除非用户明确要求重排
```

## 分区规则

```text
IF 用户要求规划新条目:
  写入 📋 下一批 (P0) 或 💡 远期 (P1/P2)
ELSE IF 用户要求开始做 X:
  将 X 移至 🔥 正在进行
ELSE IF 用户要求调整优先级:
  只移动分区或调整同区顺序
ELSE IF 用户要求完成 X:
  将 X 移至 已完成历史 顶部并标记 [x]
```

完成历史条目必须包含完成日期；如果来自 change，追加 `change: <change-id>` 和摘要。

## 边界

- Roadmap 不写 proposal、spec-delta、design 或 tasks 细节。
- 不擅自新增需求、范围、实现方式或目标版本。
- 已完成历史只读，除归档当前 change 或用户明确要求完成条目外不增删改。

## 完成条件

- `specflow/roadmap.md` 存在。
- ROADMAP_META 存在，且 next_f / next_t 大于全文件现有同类最大编号。
- 条目类型、编号和分区符合规则。
- 规划内容只来自用户指令或已确认上下文。
