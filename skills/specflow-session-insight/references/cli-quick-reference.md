# `specflow mem` CLI 速查

三个子命令的完整用法。本文件是权威参考 —— `specflow mem --help` 在运行时打印相同内容，如发现与本文件不一致以运行时 help 为准。

## 子命令

| 命令 | 用途 |
|------|------|
| `list` | 列出可检索的项目与会话日志。无参数时为默认子命令。 |
| `search <query>` | 按关键词检索历史对话，返回匹配片段（关键词前后各约 500 字符）。 |
| `context <query>` | 检索并输出上下文片段，格式更宽，适合理解前后文脉络。 |

## Flags

| Flag | 子命令 | 含义 |
|------|--------|------|
| `--phase brainstorm\|implement\|all` | search / context | 按阶段过滤对话。`brainstorm` = 规划讨论阶段，`implement` = 执行实现阶段。默认 `all`。 |
| `--limit N` | search | 最大结果数。默认 `10`。 |
| `--json` / `-j` | 全部（persistent flag） | 输出机器可解析的 JSON 格式，而非人类可读格式。 |

## 配置

`mem` 的行为由 `.specflow/config.yaml` 控制：

```yaml
mem:
  enabled: true              # 是否启用 mem 跨会话检索
  log_paths:                 # 对话日志搜索路径（自动检测，可手动覆盖）
    - ~/.opencode/sessions/
```

- `mem.enabled = false` 时，`list` 命令会输出"mem 未启用"提示。
- `mem.log_paths` 可配置多个路径，`mem` 会逐一扫描其中的 `.jsonl` 文件。

## 常用命令示例

```bash
# 列出可检索的会话日志
specflow mem list

# 按关键词检索，返回最多 20 条匹配
specflow mem search "deadlock" --limit 20

# 只搜 brainstorm 阶段的对话（适合找过去的决策讨论）
specflow mem search "数据库选型" --phase brainstorm

# 检索并输出更多上下文（适合需要理解前后文脉络的场景）
specflow mem context "lock contention" --phase implement

# 输出 JSON 格式（适合管道处理或后续脚本消费）
specflow mem search "timeout" --json
```

## 输出格式

- **默认人类可读输出**（不加 `--json`）：带文件标识和片段标记。`search` 用 `---` 分隔每条结果，`context` 用 `===` 分隔并空行。适合直接阅读。
- **`--json` 输出**：稳定 schema，适合程序解析。当你需要将 `mem` 输出喂给后续步骤时（如汇总成经验总结），优先使用 `--json`。

## 隐私

- `mem` 的所有读取都在本地完成，不上传任何对话数据。
- `mem` 对日志文件只读，不修改、不删除、不推送。

## 注意事项

- `--phase` 切分依赖于会话日志中记录的任务阶段标记。如果用户在 AI 对话循环之外操作（如从另一个终端运行 specflow 命令），该会话可能没有阶段边界。`--phase all` 是安全兜底。
- `mem` 直接索引平台 JSONL 文件。如果用户清理了会话存储，`mem` 无法恢复磁盘上已不存在的内容。
- 如果需要了解更多运行时细节，在用户 shell 中运行 `specflow mem --help`。运行时 help 是权威来源，在快速迭代期间可能领先于本文件。
