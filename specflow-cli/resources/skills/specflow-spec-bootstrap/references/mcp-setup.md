# MCP 工具配置

GitNexus 和 ABCoder 是引导 specflow spec 时的推荐工具，因为它们向 agent 暴露架构和 AST 上下文。它们是工具选择，不是平台要求。通过你的 agent host 提供的 MCP 机制进行配置。

## GitNexus

GitNexus 从代码库构建代码知识图谱。用于 module 边界、执行流、依赖关系、影响范围和图查询。

### 安装与索引

```bash
# 在代码库根目录运行。
npx gitnexus analyze

# 检查索引状态。
npx gitnexus status

# 代码变更后当索引过期时重新索引。
npx gitnexus analyze
```

索引写入 `.gitnexus/`。仅当项目已使用 embedding 时才保留 embedding；否则普通索引足以用于 spec 引导。

### MCP Server 命令

在 host 的 MCP 配置中使用此 server 命令：

```bash
npx -y gitnexus mcp
```

### 常用工具

| 工具 | 用途 |
|------|------|
| `gitnexus_query` | 按概念查找执行流和功能区域 |
| `gitnexus_context` | 检查 symbol 的 caller、callee、reference 和 process 参与情况 |
| `gitnexus_impact` | 在修改 symbol 前了解影响范围 |
| `gitnexus_detect_changes` | 完成前检查已变更的 symbol 和受影响的流程 |
| `gitnexus_cypher` | 运行直接图查询 |
| `gitnexus_list_repos` | 列出已索引的代码库 |

## ABCoder

ABCoder 将代码解析为 UniAST，提供精确的 package、file 和 node 级结构。用于 signature、type 形状、实现、依赖和反向引用。

### 安装

```bash
go install github.com/cloudwego/abcoder@latest
abcoder --help
```

### 解析代码库

```bash
abcoder parse /absolute/path/to/package \
  --lang typescript \
  --name package-name \
  --output ~/abcoder-asts
```

对于 monorepo，用稳定的 `--name` 解析每个 package，以便 task 说明可以引用相同的 repository 名称。

### MCP Server 命令

在 host 的 MCP 配置中使用此 server 命令：

```bash
abcoder mcp ~/abcoder-asts
```

### 常用工具

| 工具 | Layer | 用途 |
|------|-------|------|
| `list_repos` | 1 | 列出已解析的代码库 |
| `get_repo_structure` | 2 | 检查 package 和 file |
| `get_package_structure` | 3 | 检查 package 内的 node |
| `get_file_structure` | 3 | 检查文件中的 function、class、type 和 signature |
| `get_ast_node` | 4 | 获取代码、依赖、引用和实现 |

## 验证

配置完成后，从 agent host 验证两个 MCP server 均可见。然后在开始 spec 编写之前，对每个 server 运行一次简单查询。

```bash
ls .gitnexus/meta.json
ls ~/abcoder-asts/*.json
```
