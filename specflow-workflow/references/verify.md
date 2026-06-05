# Verify Stage

本文件定义验证阶段。Verify 主要检查 tasks 是否有遗漏项没做，并记录归档前验证摘要；它不是独立重做实现或重新设计方案的阶段。

路径锚定：本文件中的 `specflow/`、manifest、测试配置、命令和 diff 路径均相对于 `{PROJECT_ROOT}/`。

## 目标

- 对照 `tasks.md` 检查是否存在未完成任务、遗漏任务或未说明例外。
- 汇总 Apply 阶段已完成的测试、审查和质量门禁结果。
- 发现必要检查缺失时补跑；无法补跑时记录原因和风险。
- 明确 passed / failed / partial。

## 输入

- `specflow/changes/<change-id>/tasks.md`
- `specflow/changes/<change-id>/proposal.md`
- `specflow/changes/<change-id>/spec-delta.md`，如果存在
- `specflow/changes/<change-id>/design.md`，如果存在
- 项目 manifest、测试配置、已有命令说明
- 本次变更 diff

## 输出

- `specflow/changes/<change-id>/verification.md`

## 遗漏检查

```text
IF 当前 change 目录中的 tasks.md 存在未完成复选框且没有备注说明例外:
  Result = failed，并返回 Apply 或 Tasks 阶段处理
ELSE IF tasks.md 未覆盖 proposal、spec-delta 或 design 的关键约束:
  Result = failed，并返回 Tasks 阶段补齐
ELSE IF tasks.md 中验证项缺少结果记录:
  补齐检查结果；无法补齐时记录原因
ELSE:
  汇总验证摘要
```

## 检查命令

```text
IF Apply 阶段已运行适用检查且结果可追溯:
  记录已有结果
ELSE IF 项目存在聚合质量门禁命令:
  运行聚合命令并记录结果
ELSE IF 项目存在明确 lint、类型检查、测试或构建命令:
  运行适用命令并记录结果
ELSE:
  记录未发现可执行命令，不臆造命令
```

不得臆造不存在的 build、test、lint 或 codegen 命令。

## Result

- `passed`：tasks 无遗漏，所有适用检查通过。
- `failed`：存在未完成任务、遗漏任务、失败检查或未说明例外。
- `partial`：tasks 已完成，但部分检查无法执行或工具链缺失，并已记录风险。

## 完成条件

- `{PROJECT_ROOT}/specflow/changes/<change-id>/verification.md` 存在。
- Result、Checks、Notes 已填写。
- tasks 的未完成项、遗漏项和例外已明确处理。
- 无法执行的检查已说明原因和风险。
