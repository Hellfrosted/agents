#!/usr/bin/env node

const http = require("node:http");
const { createGoalPanel, isValidJsonContentType, verifyMutationRequest } = require("./goal_panel");
const { createGoalPanelSession } = require("./goal_panel_session");

const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 43873;
const FORCE_EXIT_MS = 1000;
const goalPanelSession = createGoalPanelSession();

function usage(exitCode = 1) {
  const text = [
    "Usage:",
    "  node scripts/codex_goal_panel_server.js [--thread <thread-id>] [--host 127.0.0.1] [--port 43873]",
    "",
    "Serves a local Codex Goal panel for the current thread.",
    "The server binds to 127.0.0.1 by default. Prefer --thread or CODEX_THREAD_ID for the target thread.",
    "",
  ].join("\n");
  const stream = exitCode === 0 ? process.stdout : process.stderr;
  stream.write(text);
  process.exit(exitCode);
}

function parseArgs(argv) {
  const options = {
    host: DEFAULT_HOST,
    port: DEFAULT_PORT,
    threadId: process.env.CODEX_THREAD_ID || null,
    workspaceRoot: process.env.CODEX_GOAL_PANEL_WORKSPACE_ROOT || process.cwd(),
  };

  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === "--help" || arg === "-h") usage(0);
    if (arg === "--host") {
      options.host = argv[i + 1] || "";
      i += 1;
      continue;
    }
    if (arg === "--port") {
      const parsed = Number(argv[i + 1]);
      if (!Number.isInteger(parsed) || parsed <= 0 || parsed > 65535) {
        throw new Error("--port must be a valid TCP port");
      }
      options.port = parsed;
      i += 1;
      continue;
    }
    if (arg === "--thread") {
      options.threadId = argv[i + 1] || "";
      i += 1;
      continue;
    }
    usage(1);
  }

  if (options.host !== "127.0.0.1" && options.host !== "localhost") {
    throw new Error("Refusing to bind Codex Goal panel outside localhost.");
  }
  return options;
}

function shutdownServer(state, reason) {
  if (state.shuttingDown) return;
  state.shuttingDown = true;
  if (state.sweepTimer) clearInterval(state.sweepTimer);

  const forceExit = setTimeout(() => process.exit(0), FORCE_EXIT_MS);
  forceExit.unref();

  if (!state.server) {
    process.exit(0);
    return;
  }

  state.server.close(() => {
    clearTimeout(forceExit);
    process.exit(0);
  });
  process.stderr.write(`Codex Goal panel server stopping: ${reason}\n`);
}

function startIdleShutdownTimer(panel) {
  panel.state.sweepTimer = setInterval(() => panel.runIdleSweep(), panel.state.idleSweepMs);
  panel.state.sweepTimer.unref();
}

function createServer(options) {
  const panel = createGoalPanel(options, {
    requestShutdown: (reason) => shutdownServer(panel.state, reason),
    session: goalPanelSession,
  });
  const server = http.createServer(panel.handleRequest);
  panel.state.server = server;
  startIdleShutdownTimer(panel);
  return server;
}

function main() {
  try {
    const options = parseArgs(process.argv.slice(2));
    const server = createServer(options);
    server.listen(options.port, options.host, () => {
      const rootUrl = `http://${options.host}:${options.port}/`;
      const threadUrl = options.threadId
        ? `${rootUrl}?threadId=${encodeURIComponent(options.threadId)}`
        : rootUrl;
      process.stdout.write(
        `${JSON.stringify(
          {
            listening: true,
            url: threadUrl,
            host: options.host,
            port: options.port,
            thread_id: options.threadId,
            claim_file: goalPanelSession.paths.claimFile,
          },
          null,
          2,
        )}\n`,
      );
    });
  } catch (error) {
    process.stderr.write(`${error.message || String(error)}\n`);
    process.exit(1);
  }
}

if (require.main === module) {
  main();
}

module.exports = {
  createServer,
  isValidJsonContentType,
  parseArgs,
  verifyMutationRequest,
};
