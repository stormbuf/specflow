# Spec 写作标准

specflow spec 是面向未来 agent 的编码指引。它应该解释如何在这个代码库中工作，而不是一个通用项目可能如何组织。

## 从证据出发

每条重要规则都应有以下之一作为支撑：

- 展示首选 pattern 的源码文件。
- 展示预期行为的测试文件。
- 定义约定的项目文档。
- 跨多个文件反复出现的 pattern。

仅在能令规则更清晰时使用简短代码片段。优先链接文件路径并指明 symbol 或行为。

## 文件结构

保持 spec 树与项目对齐：

- 以 `index.md` 作为 spec 目录的导航文件。
- 当开发者会独立查找某些主题时，拆分为单独文件。
- 当单独文件会重复相同规则时，合并主题。
- 删除不适用的模板文件。
- 为模板遗漏的重要本地 pattern 新增文件。

## 内容标准

好的 spec 章节应包含：

- 规则的适用场景。
- 应遵循的本地 pattern。
- 证明该 pattern 的源码或测试文件。
- 常见错误或反模式。
- 当验证命令具体且可靠时，提供验证命令或检查方式。

应避免：

- 占位符描述。
- 通用框架建议。
- 只在某个 agent host 中有效的工具指令。
- 大段复制的代码块。
- 基于偶然实现细节的规则。

## 示例形状

```markdown
## Command Handler

Command handler 应将参数解析、验证和副作用分离。本地 pattern 是：

- 在 command 边界解析 CLI flag。
- 在调用核心逻辑之前，将原始输入转换为 typed task option。
- 将文件系统写入保持在 command 或 service layer，不要放在 template helper 中。

参考文件：
- `packages/cli/src/commands/example.ts`
- `packages/cli/test/commands/example.test.ts`

避免将原始 `process.argv` 或未验证的 config 对象传入共享 helper。
```

## 最终检查

完成前执行：

```bash
grep -R "To be filled\\|TODO: fill\\|placeholder" .specflow/spec
```

同时检查链接、index 文件，以及是否有 spec 仍在描述模板而非当前代码库。
