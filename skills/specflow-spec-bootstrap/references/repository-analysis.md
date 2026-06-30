# 代码库分析

目标是在编写规则之前发现项目的真实架构。不要从通用 spec 模板出发填空。从代码出发，让 spec 结构自然形成。

## 分析顺序

1. 阅读现有 `.specflow/spec/` 目录树，记录哪些文件是模板、已过时、或已是项目特定的。
2. 检查 package manifest、构建脚本、workspace 配置和顶层文档，识别 package 和 runtime layer。
3. 用 GitNexus 查找执行流、module cluster、依赖 hub 和影响敏感区域。
4. 用 ABCoder 或语言原生工具获取精确的 signature、type、class 边界和实现示例。
5. 在将任何发现转化为 spec 规则之前，直接阅读代表性源码和测试文件。

## 需要捕获的信息

| 领域 | 问题 |
|------|------|
| Package 边界 | 每个 package 拥有什么？哪些 import 跨越了边界？ |
| Runtime layer | 哪些代码是 CLI、backend、frontend、worker、共享库、仅测试或工具？ |
| 核心抽象 | 哪些 type、service、store、command、route 或 adapter 定义了系统形态？ |
| 数据流 | 用户输入从哪里进入，如何验证，状态持久化在哪里？ |
| 错误处理 | 失败如何表示、记录、暴露和测试？ |
| 配置管理 | 默认值、环境配置、生成文件和模板放在哪里？ |
| 测试风格 | 哪些测试风格是新工作的可信参考？ |

## GitNexus 用法

先广后深，再检查具体 symbol：

```text
gitnexus_query({query: "CLI command execution flow"})
gitnexus_query({query: "template generation and migration"})
gitnexus_context({name: "SymbolName"})
gitnexus_cypher({query: "MATCH (n)-[r]->(m) RETURN n.name, type(r), m.name LIMIT 30"})
```

用 GitNexus 结果找到重要文件和流程。在检查相关源码文件之前，不要将 graph 输出作为最终权威引用。

## ABCoder 用法

当 spec 需要精确的代码结构时使用 ABCoder：

```text
list_repos()
get_repo_structure({repo_name: "package-name"})
get_file_structure({repo_name: "package-name", file_path: "src/example.ts"})
get_ast_node({repo_name: "package-name", node_ids: [{mod_path: "...", pkg_path: "...", name: "SymbolName"}]})
```

ABCoder 在文档化 constructor pattern、function signature、type contract 和 reference chain 时最有价值。

## 分析笔记

分析时保持简短笔记。笔记应包含：

- Package 或 layer 名称。
- 定义本地 pattern 的文件。
- Spec 应该教授的规则。
- 在旧代码、注释、测试或 migration path 中发现的反模式。
- 应该创建、删除、重命名或合并的 spec 文件。
