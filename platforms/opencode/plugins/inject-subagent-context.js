// inject-subagent-context.js
//
// OpenCode 插件：拦截 task 工具调用（subagent 派发），exec 调用
// `specflow build-context` 取上下文，并拼装到 sub-agent prompt 前。
//
// 对应接口契约 §2.3。
// Hook 事件：tool.execute.before
//
// 关键点：
// - subagentType 直接取 output.args.subagent_type，不做前缀替换
// - 从 agents.yaml 查 agentConf，未声明则跳过（扩展点）
// - 上下文 + 行为约束 + 原始 prompt 拼装为最终 prompt

import {
  isSpecflowProject,
  loadAgentsConfig,
  exec,
  logError,
} from "../lib/specflow-context.js";

/**
 * 拼装最终 prompt（上下文 + 约束 + 原始 prompt）。
 * @param {string} agentType sub-agent 类型名（如 specflow-implement）
 * @param {string} originalPrompt 原始 prompt
 * @param {string} context specflow build-context 输出的上下文文本
 * @param {string[]} constraints 行为约束列表
 * @returns {string}
 */
function wrapPrompt(agentType, originalPrompt, context, constraints) {
  const constraintBlock =
    Array.isArray(constraints) && constraints.length > 0
      ? `\n\n## 行为约束（必须遵守）\n${constraints
          .map((c) => `- ${c}`)
          .join("\n")}`
      : "";

  return `<!-- specflow-hook-injected -->
# ${agentType} Agent Task

## 上下文

${context}${constraintBlock}

---

## 你的任务

${originalPrompt || ""}`;
}

export default async ({ directory }) => {
  return {
    "tool.execute.before": async (input, output) => {
      try {
        if (process.env.SPECFLOW_HOOKS === "0") return;
        if (!isSpecflowProject(directory)) return;
        if (!input || input.tool !== "task") return; // 只拦截 task 工具

        const args = output && output.args;
        if (!args) return;

        const subagentType = args.subagent_type; // 原名，如 specflow-implement
        if (!subagentType || typeof subagentType !== "string") return;

        // 直接查配置，不做前缀替换
        // 每次 hook 重新加载 agents.yaml（支持热更新）
        const config = loadAgentsConfig(directory);
        const agentConf = config.agents
          ? config.agents[subagentType]
          : undefined;
        if (!agentConf) return; // 配置未声明则不处理（扩展点）

        // exec 调用 Go CLI 构建上下文
        const result = await exec(
          ["specflow", "build-context", subagentType],
          { cwd: directory }
        );
        if (result.exitCode !== 0) return;

        const context = result.stdout || "";
        const originalPrompt = args.prompt || "";

        // 拼装最终 prompt（原地修改 output.args.prompt）
        args.prompt = wrapPrompt(
          subagentType,
          originalPrompt,
          context,
          agentConf.constraints
        );
      } catch (e) {
        // 插件任何异常都不应影响宿主
        logError(directory, `inject-subagent-context failed: ${e.message}`);
      }
    },
  };
};
