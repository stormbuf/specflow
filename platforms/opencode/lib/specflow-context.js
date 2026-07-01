// specflow-context.js
//
// specflow OpenCode 插件共享工具函数。
// 三个插件（inject-workflow-state / inject-subagent-context / session-start）通过
// `import { ... } from "../lib/specflow-context.js"` 复用这里的实现。
//
// 设计原则：
// - 零外部依赖，仅用 Node.js 内置模块（fs / path / child_process / util）
// - YAML 解析使用内置简易解析器（针对 agents.yaml 结构，非完整 YAML 实现）
// - 所有函数都做防御性处理，永远不抛异常给宿主

import { readFileSync, existsSync, appendFileSync } from "node:fs";
import { join } from "node:path";
import { execFile, execFileSync } from "node:child_process";
import { promisify } from "node:util";

const execFileP = promisify(execFile);

/**
 * 检测给定目录是否为 specflow 项目。
 * 判据：目录下存在 `.specflow/workflow.md`。
 * @param {string} directory 项目根目录绝对路径
 * @returns {boolean}
 */
export function isSpecflowProject(directory) {
  try {
    if (!directory || typeof directory !== "string") return false;
    return existsSync(join(directory, ".specflow", "workflow.md"));
  } catch {
    return false;
  }
}

/**
 * 检测当前消息是否来自 specflow 子 agent（避免递归注入）。
 * 判据：input.agent 以 "specflow-" 开头。
 * @param {object} input hook 事件输入对象
 * @returns {boolean}
 */
export function isSpecflowSubagent(input) {
  try {
    return Boolean(input && typeof input === "object" &&
      typeof input.agent === "string" && input.agent.startsWith("specflow-"));
  } catch {
    return false;
  }
}

/**
 * 读取并解析 `.specflow/agents.yaml`。
 * 优先调用 `specflow agents list --json`（支持完整 YAML 语法），
 * CLI 不可用时回退到内置简易解析器。
 * @param {string} directory 项目根目录绝对路径
 * @returns {{ agents: Record<string, object> }}
 */
export function loadAgentsConfig(directory) {
  // 优先使用 CLI 的 JSON 输出
  try {
    const stdout = execFileSync("specflow", ["agents", "list", "--json"], {
      cwd: directory,
      timeout: 5000,
      encoding: "utf8",
      stdio: ["pipe", "pipe", "pipe"],
    });
    const agents = JSON.parse(stdout);
    return { agents };
  } catch (_) {
    // CLI 不可用，回退到内置解析器
  }

  // 回退：内置简易 YAML 解析器
  try {
    const path = join(directory, ".specflow", "agents.yaml");
    if (!existsSync(path)) return { agents: {} };
    const text = readFileSync(path, "utf8");
    const parsed = parseSimpleYAML(text);
    if (!parsed || typeof parsed !== "object" || !parsed.agents) {
      return { agents: {} };
    }
    if (typeof parsed.agents !== "object" || Array.isArray(parsed.agents)) {
      return { agents: {} };
    }
    return { agents: parsed.agents };
  } catch {
    return { agents: {} };
  }
}

/**
 * exec 调用外部命令（主要用于调用 specflow CLI）。
 * @param {string[]|string} cmd 命令数组（[bin, ...args]）或空格分隔的字符串
 * @param {{ cwd?: string, env?: object, timeout?: number }} [options]
 * @returns {Promise<{ exitCode: number, stdout: string, stderr: string }>}
 */
export async function exec(cmd, options = {}) {
  let bin;
  let args;
  if (Array.isArray(cmd)) {
    bin = cmd[0];
    args = cmd.slice(1);
  } else {
    const parts = String(cmd || "").split(/\s+/).filter(Boolean);
    bin = parts[0];
    args = parts.slice(1);
  }

  if (!bin) {
    return { exitCode: 1, stdout: "", stderr: "exec: empty command" };
  }

  try {
    const { stdout, stderr } = await execFileP(bin, args, {
      maxBuffer: 10 * 1024 * 1024,
      timeout: typeof options.timeout === "number" ? options.timeout : 30000,
      cwd: options.cwd,
      env: options.env ? { ...process.env, ...options.env } : process.env,
    });
    return { exitCode: 0, stdout: stdout || "", stderr: stderr || "" };
  } catch (err) {
    // 进程非零退出或无法 spawn 时，execFile 的 promise 会 reject。
    // - 进程已运行但非零退出：err 带 code(number) / stdout / stderr
    // - 无法 spawn（如二进制不存在）：err.code 为 'ENOENT' 等字符串，stderr 可能为空串
    // 当 stderr 为空时回退到 err.message，确保调用方总能拿到可读的错误信息。
    const exitCode = typeof err.code === "number" ? err.code : 1;
    const stdout = typeof err.stdout === "string" ? err.stdout : "";
    const stderr =
      (typeof err.stderr === "string" && err.stderr.length > 0
        ? err.stderr
        : "") || (err.message || "");
    return { exitCode, stdout, stderr };
  }
}

/**
 * 将错误日志写入 `.specflow/logs/plugins.log`。
 * 仅当 `.specflow/logs` 目录已存在时才写入（不主动创建目录）。
 * 写入失败时静默忽略，不影响宿主。
 * @param {string} directory 项目根目录绝对路径
 * @param {string} message 日志消息
 */
export function logError(directory, message) {
  try {
    const logsDir = join(directory, ".specflow", "logs");
    if (!existsSync(logsDir)) return; // Don't create dirs, just skip if not a specflow project
    const timestamp = new Date().toISOString();
    const logLine = `[${timestamp}] ${message}\n`;
    appendFileSync(join(logsDir, "plugins.log"), logLine);
  } catch (_) { /* swallow */ }
}

// ---------------------------------------------------------------------------
// 内置简易 YAML 解析器
//
// 仅针对 agents.yaml 的结构实现：嵌套 map（按缩进）、标量值、列表、引号字符串、
// bool / null / 数字。非完整 YAML 1.2 实现，不处理 anchor / alias / 多行字符串 /
// flow style 等高级特性。解析失败时返回 {}。
// ---------------------------------------------------------------------------

/**
 * @param {string} text
 * @returns {object}
 */
function parseSimpleYAML(text) {
  const lines = preprocessLines(text);
  if (lines.length === 0) return {};
  const { value } = parseNode(lines, 0, lines[0].indent);
  return value || {};
}

/**
 * 预处理：去掉空行与注释行，计算每行缩进，去掉行内注释。
 * @returns {Array<{ indent: number, content: string }>}
 */
function preprocessLines(text) {
  const result = [];
  const rawLines = String(text).split(/\r?\n/);
  for (const rawLine of rawLines) {
    if (rawLine.trim() === "") continue;
    const trimmed = rawLine.trim();
    if (trimmed.startsWith("#")) continue;
    const indent = rawLine.length - rawLine.replace(/^\s+/, "").length;
    const content = stripInlineComment(trimmed);
    if (content === "") continue;
    result.push({ indent, content });
  }
  return result;
}

/**
 * 去掉行内注释（# 前需有空白），尊重引号。
 */
function stripInlineComment(s) {
  let inSingle = false;
  let inDouble = false;
  for (let i = 0; i < s.length; i++) {
    const ch = s[i];
    if (ch === "'" && !inDouble) inSingle = !inSingle;
    else if (ch === '"' && !inSingle) inDouble = !inDouble;
    else if (ch === "#" && !inSingle && !inDouble) {
      if (i === 0 || /\s/.test(s[i - 1])) {
        return s.slice(0, i).trim();
      }
    }
  }
  return s.trim();
}

/**
 * 解析一个节点（map 或 list），从 startIdx 开始，所有同级行 indent === indent。
 * @returns {{ value: any, nextIdx: number }}
 */
function parseNode(lines, startIdx, indent) {
  if (startIdx >= lines.length) return { value: null, nextIdx: startIdx };
  const first = lines[startIdx];
  if (first.content.startsWith("-")) {
    return parseList(lines, startIdx, indent);
  }
  return parseMap(lines, startIdx, indent);
}

function parseMap(lines, startIdx, indent) {
  const result = {};
  let i = startIdx;
  while (i < lines.length) {
    const line = lines[i];
    if (line.indent < indent) break;
    if (line.indent > indent) { i++; continue; }
    if (line.content.startsWith("-")) break; // 不是 map 行

    const entry = parseMapEntry(line.content);
    if (entry.hasInlineValue) {
      result[entry.key] = entry.value;
      i++;
    } else {
      // 嵌套块：找下一行的缩进
      const childIndent = findChildIndent(lines, i + 1, indent);
      if (childIndent === null) {
        result[entry.key] = null;
        i++;
      } else {
        const { value, nextIdx } = parseNode(lines, i + 1, childIndent);
        result[entry.key] = value;
        i = nextIdx;
      }
    }
  }
  return { value: result, nextIdx: i };
}

function parseList(lines, startIdx, indent) {
  const result = [];
  let i = startIdx;
  while (i < lines.length) {
    const line = lines[i];
    if (line.indent < indent) break;
    if (line.indent > indent) { i++; continue; }
    if (!line.content.startsWith("-")) break;

    const itemContent = line.content.replace(/^-\s*/, "").trim();
    if (itemContent === "") {
      // 嵌套块 item
      const childIndent = findChildIndent(lines, i + 1, indent);
      if (childIndent === null) {
        result.push(null);
        i++;
      } else {
        const { value, nextIdx } = parseNode(lines, i + 1, childIndent);
        result.push(value);
        i = nextIdx;
      }
    } else {
      result.push(parseScalar(itemContent));
      i++;
    }
  }
  return { value: result, nextIdx: i };
}

/**
 * 找 parentIndent 之后的第一个子块的缩进。
 * @returns {number|null}
 */
function findChildIndent(lines, startIdx, parentIndent) {
  if (startIdx >= lines.length) return null;
  const first = lines[startIdx];
  if (first.indent <= parentIndent) return null;
  return first.indent;
}

function parseMapEntry(content) {
  const colonIdx = findColon(content);
  if (colonIdx === -1) {
    return { key: content, value: null, hasInlineValue: true };
  }
  const key = content.slice(0, colonIdx).trim();
  const rawValue = content.slice(colonIdx + 1).trim();
  if (rawValue === "") {
    return { key, value: null, hasInlineValue: false };
  }
  return { key, value: parseScalar(rawValue), hasInlineValue: true };
}

/**
 * 找第一个不在引号内的冒号。
 */
function findColon(s) {
  let inSingle = false;
  let inDouble = false;
  for (let i = 0; i < s.length; i++) {
    const ch = s[i];
    if (ch === "'" && !inDouble) inSingle = !inSingle;
    else if (ch === '"' && !inSingle) inDouble = !inDouble;
    else if (ch === ":" && !inSingle && !inDouble) return i;
  }
  return -1;
}

/**
 * 解析标量值：null / bool / 数字 / 引号字符串 / 普通字符串。
 */
function parseScalar(s) {
  const v = s.trim();
  if (v === "null" || v === "~" || v === "") return null;
  if (v === "true") return true;
  if (v === "false") return false;
  if (v === "[]") return [];
  if (v === "{}") return {};
  if ((v.startsWith('"') && v.endsWith('"') && v.length >= 2) ||
      (v.startsWith("'") && v.endsWith("'") && v.length >= 2)) {
    return v.slice(1, -1);
  }
  if (/^-?\d+$/.test(v)) return parseInt(v, 10);
  if (/^-?\d+\.\d+$/.test(v)) return parseFloat(v);
  return v;
}

export default {
  isSpecflowProject,
  isSpecflowSubagent,
  loadAgentsConfig,
  exec,
  logError,
};
