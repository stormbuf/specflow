# 代码复用思维指南

> **目的**：在写新代码之前停下来想一想——这个逻辑是不是已经存在了？

---

## 问题说明

**重复代码是不一致性 bug 的头号来源。**

当你复制粘贴或重写已有逻辑时：

- Bug 修复无法传播到所有副本
- 行为随时间逐渐分叉
- 代码库变得更难理解

---

## 编写新代码前的步骤

### Step 1：先搜索

```bash
# 搜索相似的函数名
grep -r "functionName" .

# 搜索相似逻辑的关键词
grep -r "keyword" .
```

### Step 2：问自己这些问题

| 问题 | 如果答案是"是"…… |
|------|-----------------|
| 是否已存在相似的函数？ | 使用或扩展它 |
| 这个模式是否在别处用过？ | 遵循已有模式 |
| 这是否可以成为共享 utility？ | 放到正确的位置创建 |
| 我是否在从另一个文件复制代码？ | **停下**——抽取到共享模块 |

---

## 常见重复模式

### 模式 1：复制粘贴函数

**Bad**：把一个校验函数复制到另一个文件

**Good**：抽取到共享 utilities，在需要的地方 import

### 模式 2：相似组件

**Bad**：创建一个与已有组件 80% 相似的新组件

**Good**：通过 props / variants 扩展已有组件

### 模式 3：重复常量

**Bad**：在多个文件中定义同一个常量

**Good**：单一数据源，到处 import

### 模式 4：重复的 payload 字段抽取

**Bad**：多个 consumer 各自在本地 cast 同一个 JSON / event 字段：

```typescript
const description = (ev as { description?: string }).description;
const context = (ev as { context?: ContextEntry[] }).context;
```

即使代码只有两行，这也是重复的 contract 逻辑。每个 consumer 现在都有自己的"什么是合法 payload"的定义。

**Good**：把 decoder、type guard 或 projection 放在数据 owner 旁边：

```typescript
if (isThreadEvent(ev)) {
  renderThreadEvent(ev);
}
```

**规则**：如果同一个 untyped payload 字段在 2 个以上位置被读取，在添加第三个 reader 之前创建共享的 type guard / normalizer / projection。

---

## 何时抽象 / 何时不抽象

**应该抽象**：

- 同样的代码出现 3 次以上
- 逻辑复杂到可能有 bug
- 多人可能需要这个功能

**不该抽象**：

- 只用了一次
- 简单的一行代码
- 抽象后比重复代码更复杂

---

## Reducer 应使用穷举结构

当 state 派生自 action 类型的值（`action`、`kind`、`status`、`phase`）时，优先使用单个 `switch` 的 reducer，而非分散的 `if/else` 更新。

```typescript
// BAD - action 特定的状态转换难以审计
if (action === "opened") { ... }
else if (action === "comment") { ... }
else if (action === "status") { ... }

// GOOD - 单个 reducer 拥有完整的转换表
switch (event.action) {
  case "opened":
    ...
    return;
  case "comment":
    ...
    return;
}
```

当 event log 是 source of truth 时，这一点尤为重要。reducer 是文档化的 replay model，display 代码和 command 不应重复这个 replay model 的片段。

---

## 批量修改后的检查

当你对多个文件做了类似修改后：

1. **复查**：是否遗漏了某些实例？
2. **搜索**：运行 grep 查找可能遗漏的位置
3. **考虑**：这是否应该被抽象？

---

## 提交前检查清单

- [ ] 已搜索是否存在相似的已有代码
- [ ] 没有应共享却被复制粘贴的逻辑
- [ ] 没有在共享 decoder 之外重复抽取 untyped payload 字段
- [ ] 常量只在一个地方定义
- [ ] 相似模式遵循相同结构
- [ ] Reducer / action 转换逻辑集中在一个 reducer 或 command dispatcher 中

---

## Gotcha：Python if/elif/else 穷举检查

**问题**：Python 的 `if/elif/else` 链没有编译期穷举检查。当你给 `Literal` 类型（如 `Platform`）添加新值时，已有的 `if/elif/else` 链会静默 fallthrough 到 `else` 分支，返回错误的默认值。

**症状**：新 platform 部分工作——某些方法返回 Claude 的默认值而非 platform 特定值，且不报错。

**示例**（`cli_adapter.py`）：

```python
# BAD: "gemini" fallthrough 到 else，返回 "claude"
@property
def cli_name(self) -> str:
    if self.platform == "opencode":
        return "opencode"
    else:
        return "claude"  # gemini 静默获得 "claude"！

# GOOD: 每个 platform 都有显式分支
@property
def cli_name(self) -> str:
    if self.platform == "opencode":
        return "opencode"
    elif self.platform == "gemini":
        return "gemini"
    else:
        return "claude"
```

**预防**：当你给 Python `Literal` 类型添加新值时，搜索所有基于该类型做 switch 的 `if/elif/else` 链，为每个新值添加显式分支。不要依赖 `else` 对新值是正确的。

---

**核心原则**：如果你要写新代码，先证明它不存在。
