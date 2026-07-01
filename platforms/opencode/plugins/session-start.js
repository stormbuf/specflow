// session-start.js
//
// OpenCode 插件：每会话首条消息触发（内存去重），exec 调用
// `specflow get-context` 获取 session 上下文并注入。
//
// 对应接口契约 §2.4。
// Hook 事件：chat.message（per-session 去重）
//
// 输出格式为 XML 风格标签块：
//   <specflow-session-context>...</specflow-session-context>
//   <current-state>...</current-state>
//   <spec-indexes>...</spec-indexes>
//   <journal>...</journal>
//   <ready>...</ready>

import {
  isSpecflowProject,
  isSpecflowSubagent,
  exec,
  logError,
} from "../lib/specflow-context.js";

/**
 * 根据 specflow get-context 的返回构造 session context 块。
 * @param {object} ctx get-context --json 的解析结果
 * @returns {string}
 */
function buildSessionContextBlock(ctx) {
  const developer = ctx.developer || "(未知)";
  const vcs = ctx.vcs || "(未知)";

  // 活跃任务行
  let activeTaskLine = "无";
  if (ctx.active_task && ctx.active_task.task_id) {
    const t = ctx.active_task;
    const title = t.title ? ` — ${t.title}` : "";
    activeTaskLine = `${t.task_id} (${t.status || "unknown"})${title}`;
  }

  // spec 索引列表
  const specIndexes = Array.isArray(ctx.spec_indexes)
    ? ctx.spec_indexes.map((p) => `- ${p}`).join("\n")
    : "(暂无)";

  // journal
  const journal = ctx.journal_latest
    ? `最近会话日志: ${ctx.journal_latest}`
    : "(暂无)";

  return `<specflow-session-context>
Specflow SessionStart 上下文。用于定位当前会话状态，按需加载细节。
</specflow-session-context>

<current-state>
开发者: ${developer}
VCS: ${vcs}
活跃任务: ${activeTaskLine}
</current-state>

<spec-indexes>
${specIndexes}
</spec-indexes>

<journal>
${journal}
</journal>

<ready>
上下文已加载。遵循 <current-state>。按需加载 workflow/spec/task 细节。
</ready>`;
}

export default async ({ directory, client }) => {
  // 内存去重：每个 session 只注入一次
  const processed = new Set();

  return {
    "chat.message": async (input, output) => {
      try {
        if (isSpecflowSubagent(input)) return;
        if (process.env.SPECFLOW_HOOKS === "0") return;
        if (!isSpecflowProject(directory)) return;

        const sessionID = input && input.sessionID;
        if (!sessionID) return; // 没有 sessionID 无法去重，跳过
        if (processed.has(sessionID)) return; // 本会话已注入

        // exec 调用 Go CLI 获取 session context
        const result = await exec(["specflow", "get-context", "--json"], {
          cwd: directory,
        });
        if (result.exitCode !== 0) return;

        let ctx;
        try {
          ctx = JSON.parse(result.stdout);
        } catch (e) {
          logError(directory, `session-start JSON parse failed: ${e.message}`);
          return; // 解析失败跳过
        }
        if (!ctx) return;

        const contextBlock = buildSessionContextBlock(ctx);

        // 非破坏性注入：向 output.parts 头部插入独立 text part
        const parts = (output && output.parts) || [];
        parts.unshift({
          type: "text",
          text: contextBlock,
          metadata: { specflow: { sessionStart: true } },
        });

        processed.add(sessionID);
      } catch (e) {
        // 插件任何异常都不应影响宿主
        logError(directory, `session-start failed: ${e.message}`);
      }
    },
  };
};
