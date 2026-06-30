# 跨层思维指南

> **目的**：在实现跨层功能之前，先想清楚数据如何跨层流动。

---

## 问题说明

**大多数 bug 发生在层边界，而非层内部。**

常见的跨层 bug：

- API 返回格式 A，frontend 期望格式 B
- Database 存储 X，Service 转换成 Y，但丢失了数据
- 多个层用不同方式实现了同一逻辑

---

## 实现跨层功能前的步骤

### Step 1：映射数据流

画出数据如何流动：

```
Source → Transform → Store → Retrieve → Transform → Display
```

对每一个箭头，问自己：

- 数据是什么格式？
- 可能出什么问题？
- 谁负责 validation？

### Step 2：识别边界

| 边界 | 常见问题 |
|------|----------|
| API ↔ Service | 类型不匹配、字段缺失 |
| Service ↔ Database | 格式转换、null 处理 |
| Backend ↔ Frontend | 序列化、日期格式 |
| Component ↔ Component | props 结构变化 |

### Step 3：定义契约

对每个边界：

- 输入的精确格式是什么？
- 输出的精确格式是什么？
- 可能发生哪些错误？

---

## 常见跨层错误

### 错误 1：隐式格式假设

**Bad**：假设日期格式而不检查

**Good**：在边界处做显式格式转换

### 错误 2：分散的校验

**Bad**：在多个层重复校验同一件事

**Good**：在入口点校验一次

### 错误 3：抽象泄漏

**Bad**：Component 知道 database schema

**Good**：每一层只了解它的直接邻居

### 错误 4：每个 consumer 都在解析同一个 payload

**Bad**：一个 command 读取 JSONL events 并在行内 cast 字段：

```typescript
const thread = (ev as { thread?: string }).thread;
const labels = (ev as { labels?: string[] }).labels;
```

这看起来是局部的，但意味着每个 consumer 都拥有 event contract 的一个私有版本。下次字段变更会更新一个 command 而遗漏另一个。

**Good**：在 event 边界解码一次，然后导出 typed projection：

```typescript
if (!isThreadEvent(ev)) return false;
return ev.thread === filter.thread;
```

**规则**：对于 append-only log、JSON stream、RPC payload 或配置文件，创建唯一的 owner 负责：

- event / payload 类型定义
- 从 `unknown` 到 typed 的 type guard 和 normalization
- UI command 使用的 metadata projection
- 从 source of truth replay state 的 reducer

Rendering 代码可以格式化字段，但不得重新定义 payload contract。

---

## 检查清单

### 实现前

- [ ] 已映射完整的数据流
- [ ] 已识别所有层边界
- [ ] 已定义每个边界的格式
- [ ] 已确定 validation 发生的位置

### 实现后

- [ ] 已用 edge case 测试（null、空值、非法值）
- [ ] 已验证每个边界的 error handling
- [ ] 已检查数据 round-trip 后完好无损
- [ ] 已检查 consumer 使用共享 decoder / projection，而非在本地 cast payload 字段
- [ ] 已检查派生 state 指回 source event identifier（`seq`、`id`、`version`），而非发明第二个 cursor

---

## Event Log / Projection 边界

Append-only log 是跨层契约。一个 event 会流经：

```
CLI input → event writer → events.jsonl → reader → filter → reducer → display
```

### 检查清单：添加新 Event Kind 或字段后

- [ ] 将 event kind 添加到中央 event taxonomy
- [ ] 在 event 层添加 typed event variant 或 type guard
- [ ] 为来自用户输入或 JSON 的 array / object 字段添加 normalization helper
- [ ] `seq` / `id` 的赋值只在 event writer 中进行
- [ ] filter 和 reducer 消费 typed event guard，而非本地 cast
- [ ] display 代码消费 reducer 输出或 typed event，而非 raw JSON
- [ ] 添加至少一个回归测试，证明 history replay 和 live filtering 使用同一 filter model

---

## 何时创建流程文档

以下情况建议创建详细的流程文档：

- 功能跨越 3 个及以上层
- 涉及多个团队
- 数据格式复杂
- 该功能曾经引发过 bug

---

**核心原则**：bug 不在层里，bug 在层与层之间。
