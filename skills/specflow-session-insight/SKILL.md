---
name: specflow-session-insight
description: "通过 `specflow mem` CLI 检索跨会话历史对话。当用户问'上次怎么解的'、'之前讨论过吗'、'我们当时怎么定的'、'想起一段对话'，或在新 brainstorm 可能与过往工作重叠、调试熟悉的 bug、跨会话续作任务、finish-work 回顾时触发。返回原始历史对话，由 AI 当场决定是否更新 spec、追加 task notes、内联引用或仅内化。"
trigger: "用户引用过往对话、跨会话续作、brainstorm 重复风险检测、finish-work 回顾、模式识别"
---

# specflow Session Insight

本 skill 教会 AI **如何调用 `specflow mem`** —— specflow 内置的跨会话记忆检索工具 —— 以及**何时应该主动伸手去查**。

这是一个**能力型 skill，不是工作流**。没有固定的输出文件，没有强制的写回步骤，没有"每次 finish-work 都必须运行"的规则。`mem` 返回的内容如何处理，完全由对话当下的判断决定。skill 的存在只是让 AI 知道这个能力可用，并能自主决策。

## `specflow mem` 是什么

一个本地 CLI，索引用户过往的对话日志（JSONL 文件，路径由 `.specflow/config.yaml` 的 `mem.log_paths` 配置，默认为 `~/.opencode/sessions/`），支持列出会话、按关键词检索、按阶段过滤、输出上下文片段。

`mem` 中的所有内容都不上传，全部在本地读取。

## 何时该伸手

判断标准是："一个资深同事会不会问'我们之前不是聊过这个吗？'"—— 如果会，就是该查的时刻。具体场景：

1. **brainstorm 重复风险。** 新任务触碰到用户之前做过的领域，你想在重新问用户之前先确认是否已有决策。
2. **熟悉 bug 调试。** 当前 bug 模式感觉像是用户之前报过/修过的。拉出相关历史会话可以省掉一整轮调试循环。
3. **跨会话续作。** 用户隔了一段时间回来，说"我们上次做到哪了" / "继续上次的"但没给具体上下文。
4. **决策检索。** 用户提到"我们当时对 X 的决定"，但这个决定存在于旧 brainstorm 对话里，没写进任何 `prd.md` / `spec/`。
5. **finish-work 回顾。** 用户明确要求复盘这次任务中决定了什么 / 踩了什么坑 / 有什么意外 —— 不是每次 finish-work 都强制执行，只在被要求时触发。
6. **跨会话模式识别。** 用户问"我是不是老在 X 上犯同样的错" / "我每次都踩这个坑吗" —— 跨会话检索能回答这个问题。

如果以上都不符合，不要调 `mem`。它是工具，不是仪式。

## 何时不该伸手

- 相关上下文已在当前轮次、`prd.md`、`design.md`、最近的 VCS log 或打开的文件中。`mem` 用于那些已经脱离即时可达范围的东西。
- 用户问的是代码里的事实，不是过去对话里的事实。`git log -p` / `grep` / 直接读文件更快也更权威。
- 你在 sub-agent（`specflow-implement` / `specflow-check`）中，dispatch prompt 已包含精选的上下文。再叠加 `mem` 通常只是添乱。
- 用户明确说了"别翻历史，直接回答我的问题"。

## 拿到 `mem` 返回内容后怎么处理

把输出当作**原料**，不是交付物。拿到之后，根据当前对话判断：

- **内联引用在回复中** —— 如果某段历史对话正好回答了用户的当前问题，直接引用，并标注来源会话/阶段以便用户核实。
- **更新 `<task>/prd.md` 或 `<task>/design.md`** —— 如果 `mem` 翻出了一个本应写下但没写的承重决策。先向用户提议编辑，再动手。
- **追加到 task-local notes 文件**（如 `<task>/notes.md` 或扩展现有文件）—— 如果发现属于当前任务记录但不适合放进 PRD 的内容。
- **更新 `.specflow/spec/`** —— 如果发现的是项目级约定或坑，能帮到未来任务。为此运行 `specflow-update-spec` skill —— `session-insight` 止步于发现。
- **仅内化** —— 接下来几轮回答得更好就行，不写任何东西。对于一次性回忆，这往往是最对的选择。

specflow 不规定唯一去向。把每次回忆都塞进固定文件会让文件长成噪音。让场景决定。

## 怎么调用

完整 CLI 参考见 [references/cli-quick-reference.md](references/cli-quick-reference.md)。80% 的场景是以下之一：

```bash
# 按关键词检索历史对话（默认搜索所有阶段）
specflow mem search "<关键词>"

# 按阶段过滤，只搜 brainstorm 阶段的对话
specflow mem search "<关键词>" --phase brainstorm

# 检索并输出上下文片段（比 search 返回更多周边上下文）
specflow mem context "<关键词>"

# 列出可检索的会话日志
specflow mem list
```

阶段过滤（`--phase brainstorm|implement|all`）按 specflow 任务生命周期切分会话。`brainstorm` = 规划讨论阶段，`implement` = 执行实现阶段。默认 `all`。

## 触发模式

[references/triggering-patterns.md](references/triggering-patterns.md) 列出了更多逐字用户措辞（中英文），帮助你校准"该伸手"的直觉。

## 不在范围内

- `mem` 不编辑代码、不更新文件。任何写回操作都是你当场的自主决策。
- `mem` 对日志文件只读。不推送、不同步到远端。
- 本 skill 不替代 `specflow-update-spec`（后者是将发现提升为项目级指引的正确工具），也不替代平台原生的 task / spec 工作流。
