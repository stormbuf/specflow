# 触发模式

逐字用户措辞，帮助 AI 校准"何时该伸手查 `specflow mem`"的直觉。如果用户消息命中以下某个模式而你没有伸手，大概率漏了一个明显的回忆机会。

模式按**意图**分组，而非表面措辞。同一个意图会出现在不同语言和语体中。

## 1. 过往方案回忆

用户在问"我们（或我）上次是怎么解决这个的"。历史对话里有答案；代码库只显示结果，不显示推理过程。

**触发短语：**
- "How did we solve this last time?"
- "What did we end up doing about X?"
- "We dealt with this once already, didn't we?"
- "上次怎么解的?"
- "之前是怎么搞定 X 的?"
- "我记得以前修过类似的"

**伸手：** `specflow mem search "<症状关键词>" --limit 10`，然后 `specflow mem context "<关键词>"` 深入最接近的匹配。

## 2. 决策检索

用户引用了一个存在于旧对话中、但不在任何已提交文件里的决定。去 brainstorm 阶段的对话里找。

**触发短语：**
- "What was the decision on X?"
- "Did we decide to use Postgres or SQLite?"
- "The rationale for choosing X over Y was…?"
- "我们当时为啥选了 X 而不是 Y?"
- "关于 X 我们之前是怎么定的?"
- "之前讨论过 X 的方案吗?"

**伸手：** `specflow mem search "<决策关键词>" --phase brainstorm` 恢复讨论。

## 3. 跨会话续作

用户隔了一段时间回来，上下文是隐含的。

**触发短语：**
- "Where were we?"
- "Continue from last time."
- "Pick up where we left off."
- "继续上次的"
- "我们上次做到哪了"
- "接着昨天那个任务"

**伸手：** `specflow mem list` 找到最近的会话日志，然后 `specflow mem context "<任务关键词>"` 拉出上下文。

## 4. 熟悉 bug 调试

当前 bug 感觉像之前见过的。历史会话大概率记录了解决路径。

**触发短语：**
- "I feel like I've hit this before."
- "Doesn't this look like that bug from last month?"
- "Same kind of timeout I had in X."
- "这个错好像之前见过"
- "这个 bug 是不是上次那个?"
- "怎么又是这个 error?"

**伸手：** `specflow mem search "<错误信息片段>"`。用实际错误串中一个短的、有辨识度的 token 作为锚点。

## 5. 自我模式识别

用户在问自己是不是老犯同一类错或做同一类决策。

**触发短语：**
- "Do I always make this mistake?"
- "How often have I run into X?"
- "Is this a recurring thing for me?"
- "我每次都踩这个坑吗?"
- "我老犯这个错?"
- "这类问题之前出现过几次?"

**伸手：** `specflow mem search "<主题>" --limit 50` 扫描匹配的会话。可选地 `specflow mem context "<主题>"` 深入两三个做对比。

## 6. finish-work 回顾（按需）

用户明确想回顾这次任务 —— 不是强制步骤，只在被要求时触发。

**触发短语：**
- "Summarize what we did in this task."
- "What were the key decisions / surprises?"
- "Write up the lessons from this round."
- "总结一下这次的经验"
- "记一下这次踩的坑"
- "复盘下这个任务"

**伸手：** `specflow mem context "<任务关键词>" --phase brainstorm` 和 `specflow mem context "<任务关键词>" --phase implement` 分别恢复规划和执行阶段的对话。呈现摘要 —— 尽可能给出具体的 file:line 引用。是否将摘要写到某处（PRD、spec、notes）由用户决定；提议，不自动写。

## 不该伸手的反模式

- "这个函数是干什么的?" → 读文件。
- "这个测试为什么失败?" → 读测试输出和源文件。
- "我们代码库里 X 的正确 pattern 是什么?" → grep / 读 spec 文件。
- "Y 的最新 npm 版本是多少?" → 跑 `npm view`。
- "修一下这个 bug。" → 直接调试。只有在怀疑存在历史上下文时才伸手查 `mem`；否则只是噪音。

判断标准不变：一个资深同事在回答之前会不会问"我们之前不是聊过这个吗？"—— 会就查，不会就不查。

## 输出处理决策树

```
拿到 mem 返回内容后：

├─ 某段对话直接回答了用户的当前问题？
│   └─ YES → 内联引用在回复中，标注来源会话/阶段
│
├─ 翻出了一个本应写下但没写的承重决策？
│   └─ YES → 向用户提议更新 prd.md / design.md，确认后写入
│
├─ 发现属于当前任务记录但不适合放进 PRD？
│   └─ YES → 追加到 <task>/notes.md 或扩展现有 notes
│
├─ 发现的是项目级约定或坑，能帮到未来任务？
│   └─ YES → 运行 specflow-update-spec skill，提升为 spec
│
└─ 以上都不符合？
    └─ 仅内化，接下来几轮回答得更好就行，不写任何东西
```

specflow 不规定唯一去向。把每次回忆都塞进固定文件会让文件长成噪音。让场景决定。
