# 跨平台思维指南

> **目的**：在平台特定的假设变成 bug 之前发现它们。

---

## 为什么这很重要

**大多数跨平台 bug 来自隐式假设**：

- 假设 shebang 有效 → 在 Windows 上失效
- 假设 `/` 路径分隔符 → 在 Windows 上失效
- 假设 `\n` 行尾 → 行为不一致
- 假设命令可用 → `grep` vs. `findstr`

---

## 平台差异检查清单

### 1. 脚本执行

| 假设 | macOS/Linux | Windows |
|------|-------------|---------|
| Shebang (`#!/usr/bin/env python3`) | 有效 | 被忽略 |
| 直接执行 (`./script.py`) | 有效 | 失败 |
| `python3` 命令 | 总是可用 | 可能需要 `python` |
| `python` 命令 | 可能是 Python 2 | 通常是 Python 3 |

**规则 1**：对于面向用户的文档和错误信息，显式声明平台规则（Windows 上用 `python`，其他平台用 `python3`），或通过平台感知 helper 渲染命令。

```python
# BAD - 假设 shebang 有效
print("Usage: ./script.py <args>")

# GOOD - 平台感知的措辞
print("Usage: python on Windows, python3 elsewhere")
```

**规则 2**：运行时从 JavaScript 调用 Python 时，动态检测平台：

```javascript
import { platform } from "os"

const PYTHON_CMD = platform() === "win32" ? "python" : "python3"
execSync(`${PYTHON_CMD} "${scriptPath}"`, { ... })
```

**规则 3**：当 Python 从 Python 调用时，使用 `sys.executable`：

```python
import sys
import subprocess

# BAD - 硬编码命令
subprocess.run(["python3", "other_script.py"])

# GOOD - 使用当前解释器
subprocess.run([sys.executable, "other_script.py"])
```

---

### 2. 路径处理

| 假设 | macOS/Linux | Windows |
|------|-------------|---------|
| `/` 分隔符 | 有效 | 有时有效 |
| `\` 分隔符 | 转义字符 | 原生 |
| `pathlib.Path` | 有效 | 有效 |

**规则（Python）**：所有路径操作使用 `pathlib.Path`。

```python
# BAD - 字符串拼接
path = base + "/" + filename

# GOOD - pathlib
from pathlib import Path
path = Path(base) / filename
```

**规则（TypeScript）**：路径字符串作为 logical key（Map key、JSON 字段、hash 字典 key）跨 OS 持久化时，归一化为 POSIX；直接传给 `fs.*` 时保持 OS-native。

```typescript
// BAD - logical key 携带 OS-native 分隔符
files.set(path.join(".opencode", entry), readFile(entry));  // Windows 上是 \

// GOOD - 在边界处归一化
files.set(toPosix(path.join(".opencode", entry)), readFile(entry));
```

**常见违规点**：`path.relative(cwd, fullPath)` 在 Windows 上产生 `\`。如果随后用作 hash 字典查找 key，先做 `toPosix`，否则在 Windows 上查找会失败。

---

### 3. 行尾

| 格式 | macOS/Linux | Windows | Git |
|------|-------------|---------|-----|
| `\n` (LF) | 原生 | 部分工具 | 归一化 |
| `\r\n` (CRLF) | 多余字符 | 原生 | 转换 |

**规则 1**：使用 `.gitattributes` 强制一致的行尾。

```gitattributes
* text=auto eol=lf
*.sh text eol=lf
*.py text eol=lf
```

**规则 2**：跨平台 hashing 或比较**内容**时，在计算 hash 前归一化行尾。`sha256(LF)` ≠ `sha256(CRLF)`。

```typescript
// BAD - autocrlf=true 的 Windows 用户得到不同的 hash
export function computeHash(content: string): string {
  return createHash("sha256").update(content, "utf-8").digest("hex");
}

// GOOD - hashing 前归一化
export function computeHash(content: string): string {
  const normalized = content.replace(/\r\n/g, "\n");
  return createHash("sha256").update(normalized, "utf-8").digest("hex");
}
```

---

### 4. 环境变量

| 变量 | macOS/Linux | Windows |
|------|-------------|---------|
| `HOME` | 已设置 | 使用 `USERPROFILE` |
| `PATH` 分隔符 | `:` | `;` |
| 大小写敏感 | 敏感 | 不敏感 |

**规则 1**：使用 `pathlib.Path.home()` 而非环境变量。

```python
# BAD
home = os.environ.get("HOME")

# GOOD
home = Path.home()
```

**规则 2**：向 shell 命令注入环境变量时，为实际解析命令的 shell 生成前缀。Windows 上的 "Bash" 可能通过 PowerShell、Git Bash、MSYS2 等执行。

```javascript
// BAD - 当 host shell 是 PowerShell 时失效
command = `export VAR=${shellQuote(value)}; ${command}`;

// GOOD - shell 方言感知
const prefix = process.platform === "win32" && !isWindowsPosixShell(process.env)
  ? `$env:VAR = ${powershellQuote(value)}; `
  : `export VAR=${shellQuote(value)}; `;
command = `${prefix}${command}`;
```

---

### 5. 命令可用性

| 命令 | macOS/Linux | Windows |
|------|-------------|---------|
| `grep` | 内置 | 不可用 |
| `find` | 内置 | 语法不同 |
| `cat` | 内置 | 使用 `type` |
| `tail -f` | 内置 | 不可用 |

**规则**：尽可能使用 Python 标准库而非 shell 命令。

```python
# BAD - tail -f 在 Windows 上不可用
subprocess.run(["tail", "-f", log_file])

# GOOD - 跨平台实现
def tail_follow(file_path: Path) -> None:
    """Follow a file like 'tail -f', cross-platform compatible."""
    with open(file_path, "r", encoding="utf-8", errors="replace") as f:
        f.seek(0, 2)
        while True:
            line = f.readline()
            if line:
                print(line, end="", flush=True)
            else:
                time.sleep(0.1)
```

---

### 6. 文件编码

| 默认编码 | macOS/Linux | Windows |
|----------|-------------|---------|
| 终端 | UTF-8 | 常为 CP1252 或 GBK |
| 文件 I/O | UTF-8 | 系统 locale |
| Git 输出 | UTF-8 | 可能变化 |

**规则**：始终显式指定 `encoding="utf-8"` 并使用 `errors="replace"`。

```python
# BAD - 依赖系统默认
with open(file, "r") as f:
    content = f.read()

# GOOD - 显式编码 + 错误处理
with open(file, "r", encoding="utf-8", errors="replace") as f:
    content = f.read()
```

**Git 命令**：强制 UTF-8 输出编码：

```python
git_args = ["git", "-c", "i18n.logOutputEncoding=UTF-8"] + args
result = subprocess.run(
    git_args, capture_output=True, text=True,
    encoding="utf-8", errors="replace"
)
```

---

## 提交前检查清单

- [ ] 面向用户的 Python 调用是平台感知的（Windows 用 `python`，其他用 `python3`）
- [ ] Python 从 Python 调用 subprocess 时使用 `sys.executable`
- [ ] 所有路径使用 `pathlib.Path`
- [ ] 没有硬编码路径分隔符（`/` 或 `\`）
- [ ] 作为 logical / 持久化 key 的路径字符串归一化为 POSIX；`fs.*` 调用保持 OS-native
- [ ] 跨 OS 内容 hash 在 hashing 前归一化行尾（`\r\n` → `\n`）
- [ ] 没有缺少 fallback 的平台特定命令（如 `tail -f`）
- [ ] 所有文件 I/O 指定 `encoding="utf-8"` 和 `errors="replace"`
- [ ] 所有 subprocess 调用指定 `encoding="utf-8"` 和 `errors="replace"`
- [ ] Git 命令使用 `-c i18n.logOutputEncoding=UTF-8`
- [ ] 已运行搜索找到所有受影响的位置

---

**核心原则**：如果不是显式的，那就是假设。而假设会出错。
