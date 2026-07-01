// inject-workflow-state.js
//
// OpenCode 插件：每轮用户消息触发，按当前任务 status 匹配 workflow.md 中的
// `[workflow-state:<STATUS>]` 面包屑标签块，非破坏性注入到消息前。
//
// 对应接口契约 §2.2。
// Hook 事件：chat.message
//
// 非破坏性注入（规避 OpenCode #367）：
//   方案 A（优先）：向 output.parts 头部插入独立 text part，不修改用户原文。
//   方案 B（降级，已注释）：在用户原文 text part 前 prepend。

import { readFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import {
  isSpecflowProject,
  isSpecflowSubagent,
  exec,
  logError,
} from "../lib/specflow-context.js";

// 面包屑标签块正则，与接口契约 §3.3 一致。
// 匹配 [workflow-state:STATUS]...[/workflow-state:STATUS]，STATUS 仅允许 [A-Za-z0-9_-]+
const TAG_RE =
  /\[workflow-state:([A-Za-z0-9_-]+)\]\s*\n([\s\S]*?)\n\s*\[\/workflow-state:\1\]/g;

/**
 * 读取并解析 workflow.md 中的所有面包屑标签块。
 * @param {string} directory 项目根目录
 * @returns {Record<string, string>} { status: body }
 */
function loadBreadcrumbs(directory) {
  const templates = {};
  try {
    const workflowPath = join(directory, ".specflow", "workflow.md");
    if (!existsSync(workflowPath)) return templates;
    const text = readFileSync(workflowPath, "utf8");
    // 重置 lastIndex（全局正则复用安全）
    TAG_RE.lastIndex = 0;
    let match;
    while ((match = TAG_RE.exec(text)) !== null) {
      const status = match[1];
      const body = match[2];
      templates[status] = body;
    }
  } catch (e) {
    // 读取/解析失败不报错，buildBreadcrumb 会降级为固定文案
    logError(directory, `loadBreadcrumbs failed: ${e.message}`);
  }
  return templates;
}

/**
 * 构造面包屑文本。找不到匹配 status 的块时降级为固定英文文案。
 * @param {object|null} task specflow task current 返回的任务对象（null 表示无活跃任务）
 * @param {string} status 当前状态：no_task / planning / in_progress / completed
 * @param {Record<string, string>} templates loadBreadcrumbs 的返回值
 * @returns {string}
 */
function buildBreadcrumb(task, status, templates) {
  let body = templates[status];
  if (body === undefined) {
    body = "Refer to workflow.md for current step.";
  }
  const header =
    task === null || task === undefined
      ? `Status: ${status}`
      : `Task: ${task.task_id} (${status})`;
  return `<workflow-state>\n${header}\n${body}\n</workflow-state>`;
}

/**
 * 非破坏性注入（方案 A）：向 output.parts 头部插入独立 text part。
 * 不修改用户原文 text part。
 * @param {object} output hook 输出对象
 * @param {string} breadcrumb 面包屑文本
 */
function injectNonDestructive(output, breadcrumb, directory) {
  try {
    const parts = (output && output.parts) || [];
    // 方案 A（优先）：插入独立 part，不修改用户 message text part
    parts.unshift({
      type: "text",
      text: breadcrumb,
      metadata: { specflow: { workflowState: true } },
    });

    // 方案 B（降级）：若 OpenCode chat.message hook 仅支持修改 text part，
    // 取消下方注释，在用户原文前 prepend，用明确分隔标记，不压缩不替换原文。
    // const textPartIndex = parts.findIndex(
    //   (p) => p.type === "text" && typeof p.text === "string" && !p.metadata?.specflow
    // );
    // if (textPartIndex !== -1) {
    //   const original = parts[textPartIndex].text || "";
    //   parts[textPartIndex].text = `${breadcrumb}\n\n---\n\n${original}`;
    // }
  } catch (e) {
    // 注入失败不影响宿主
    logError(directory, `injectNonDestructive failed: ${e.message}`);
  }
}

export default async ({ directory }) => {
  return {
    "chat.message": async (input, output) => {
      try {
        // 1. 跳过条件检测
        if (isSpecflowSubagent(input)) return; // 跳过子 agent 消息
        if (process.env.SPECFLOW_HOOKS === "0") return; // 环境变量禁用
        if (!isSpecflowProject(directory)) return; // 非 specflow 项目

        // 2. exec 调用 Go CLI 获取当前任务状态
        const result = await exec(
          ["specflow", "task", "current", "--json"],
          { cwd: directory }
        );
        if (result.exitCode !== 0) return;

        let task = null;
        try {
          task = JSON.parse(result.stdout);
        } catch (e) {
          logError(directory, `inject-workflow-state JSON parse failed: ${e.message}`);
          task = null; // CLI 输出 null 或解析失败均视为无活跃任务
        }

        // 3. 解析 workflow.md 面包屑标签块
        const templates = loadBreadcrumbs(directory);
        const status = task ? task.status : "no_task";

        // 4. 构造面包屑并注入
        const breadcrumb = buildBreadcrumb(task, status, templates);
        injectNonDestructive(output, breadcrumb, directory);
      } catch (e) {
        // 插件任何异常都不应影响宿主
        logError(directory, `inject-workflow-state failed: ${e.message}`);
      }
    },
  };
};
