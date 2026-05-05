#!/usr/bin/env node

const { createCodexAppServer } = require("./codex_app_server_client");

function usage(exitCode = 1) {
  const text = [
    "Usage:",
    "  node scripts/codex_goal.js get [--thread <thread-id>] [--json]",
    "  node scripts/codex_goal.js set <objective> [--budget <tokens>] [--thread <thread-id>] [--json]",
    "  node scripts/codex_goal.js pause [--thread <thread-id>] [--json]",
    "  node scripts/codex_goal.js resume [--thread <thread-id>] [--json]",
    "  node scripts/codex_goal.js complete [--thread <thread-id>] [--json]",
    "  node scripts/codex_goal.js clear [--thread <thread-id>] [--json]",
    "",
    "Defaults:",
    "  --thread defaults to CODEX_THREAD_ID, so commands run against the current Codex app thread when launched from Codex.",
    "",
  ].join("\n");
  const stream = exitCode === 0 ? process.stdout : process.stderr;
  stream.write(text);
  process.exit(exitCode);
}

function parseArgs(argv) {
  const [command, ...rest] = argv;
  if (!command || command === "help" || command === "--help" || command === "-h") usage(command ? 0 : 1);

  const options = {
    command,
    threadId: process.env.CODEX_THREAD_ID || null,
    objective: null,
    tokenBudget: undefined,
    json: false,
  };

  const positional = [];
  for (let i = 0; i < rest.length; i += 1) {
    const arg = rest[i];
    if (arg === "--json") {
      options.json = true;
      continue;
    }
    if (arg === "--thread") {
      options.threadId = rest[i + 1] || "";
      i += 1;
      continue;
    }
    if (arg === "--budget") {
      const raw = rest[i + 1];
      const parsed = Number(raw);
      if (!Number.isInteger(parsed) || parsed <= 0) {
        throw new Error("--budget must be a positive integer");
      }
      options.tokenBudget = parsed;
      i += 1;
      continue;
    }
    if (arg === "--no-budget") {
      options.tokenBudget = null;
      continue;
    }
    positional.push(arg);
  }

  if (!options.threadId) {
    throw new Error("Missing thread id. Pass --thread <thread-id> or launch from Codex with CODEX_THREAD_ID.");
  }

  if (command === "set") {
    options.objective = positional.join(" ").trim();
    if (!options.objective) throw new Error("Missing goal objective for set.");
  } else if (positional.length > 0) {
    throw new Error(`Unexpected arguments for ${command}: ${positional.join(" ")}`);
  }

  if (!["get", "set", "pause", "resume", "complete", "clear"].includes(command)) {
    usage(1);
  }

  return options;
}

function formatGoal(goal) {
  if (!goal) return ["Goal: none"];
  const lines = [
    `Goal: ${goal.objective}`,
    `Status: ${goal.status}`,
    `Tokens: ${goal.tokensUsed}${goal.tokenBudget == null ? "" : ` / ${goal.tokenBudget}`}`,
    `Time: ${goal.timeUsedSeconds}s`,
    `Thread: ${goal.threadId}`,
  ];
  return lines;
}

async function runGoalCommand(options) {
  const client = createCodexAppServer({
    clientInfo: {
      name: "codex-goal-control",
      title: "Codex Goal Control",
      version: "0.0.1",
    },
  });

  try {
    await client.initialize();
    if (options.command === "get") {
      const result = await client.request("thread/goal/get", { threadId: options.threadId });
      return { command: options.command, threadId: options.threadId, goal: result.goal || null };
    }

    if (options.command === "clear") {
      const result = await client.request("thread/goal/clear", { threadId: options.threadId });
      return { command: options.command, threadId: options.threadId, ...result };
    }

    const statusByCommand = {
      set: "active",
      pause: "paused",
      resume: "active",
      complete: "complete",
    };
    const params = {
      threadId: options.threadId,
      status: statusByCommand[options.command],
    };
    if (options.command === "set") {
      params.objective = options.objective;
      if (options.tokenBudget !== undefined) params.tokenBudget = options.tokenBudget;
    }

    const result = await client.request("thread/goal/set", params);
    return { command: options.command, threadId: options.threadId, goal: result.goal };
  } finally {
    await client.shutdown();
  }
}

async function main() {
  try {
    const options = parseArgs(process.argv.slice(2));
    const result = await runGoalCommand(options);
    if (options.json) {
      process.stdout.write(`${JSON.stringify(result, null, 2)}\n`);
      return;
    }
    if (Object.prototype.hasOwnProperty.call(result, "cleared")) {
      process.stdout.write(`Goal cleared: ${result.cleared ? "yes" : "no"}\n`);
      return;
    }
    process.stdout.write(`${formatGoal(result.goal).join("\n")}\n`);
  } catch (error) {
    process.stderr.write(`${error.message || String(error)}\n`);
    process.exit(1);
  }
}

if (require.main === module) {
  main();
}

module.exports = {
  runGoalCommand,
};
