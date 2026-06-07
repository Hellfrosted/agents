#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");
const { createGoalPanelSession } = require("./goal_panel_session");

const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 43873;
const SERVER_SCRIPT = path.join(__dirname, "codex_goal_panel_server.js");
const goalPanelSession = createGoalPanelSession();

function usage(exitCode = 1) {
  const text = [
    "Usage:",
    "  node scripts/codex_goal_panel_open.js [--thread <thread-id>] [--host 127.0.0.1] [--port 43873] [--json]",
    "",
    "Claims the current Codex thread for the local Goal panel and ensures the panel server is running.",
    "When --thread is omitted, CODEX_THREAD_ID is used.",
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
    portWasExplicit: false,
    threadId: process.env.CODEX_THREAD_ID || null,
    json: false,
  };

  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === "--help" || arg === "-h") usage(0);
    if (arg === "--json") {
      options.json = true;
      continue;
    }
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
      options.portWasExplicit = true;
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
    throw new Error("Refusing to use a non-local Goal panel host.");
  }
  if (!options.threadId) {
    throw new Error("Missing thread id. Launch from Codex or pass --thread <thread-id>.");
  }
  return options;
}

function startServer(options, claim) {
  fs.mkdirSync(goalPanelSession.paths.claimDir, { recursive: true });
  const logFd = fs.openSync(goalPanelSession.paths.serverLogFile, "a");
  const child = spawn(process.execPath, [SERVER_SCRIPT, "--host", options.host, "--port", String(options.port)], {
    detached: true,
    stdio: ["ignore", logFd, logFd],
    env: {
      ...process.env,
      CODEX_THREAD_ID: options.threadId,
      CODEX_GOAL_PANEL_WORKSPACE_ROOT: claim.workspaceRoot,
      CODEX_GOAL_PANEL_WORKSPACE_CWD: claim.workspaceCwd,
    },
  });
  goalPanelSession.writeServerState(options, child);
  child.unref();
}

async function main() {
  try {
    const options = parseArgs(process.argv.slice(2));
    const claim = goalPanelSession.claimThread(options.threadId);
    const server = await goalPanelSession.ensureServer(options, {
      startServer: (serverOptions) => startServer(serverOptions, claim),
    });
    const origin = `http://${options.host}:${options.port}`;
    const threadUrl = `${origin}/?threadId=${encodeURIComponent(options.threadId)}`;
    const result = {
      threadId: options.threadId,
      url: threadUrl,
      rootUrl: origin,
      threadUrl,
      server,
      claimFile: goalPanelSession.paths.claimFile,
      serverLogFile: goalPanelSession.paths.serverLogFile,
      workspaceRoot: claim.workspaceRoot,
      workspaceSource: claim.workspaceSource,
      claimedAt: claim.updatedAt,
    };
    if (options.json) {
      process.stdout.write(`${JSON.stringify(result, null, 2)}\n`);
      return;
    }
    process.stdout.write(`Codex Goal panel: ${result.threadUrl}\n`);
    process.stdout.write(`Root URL: ${result.rootUrl}\n`);
  } catch (error) {
    process.stderr.write(`${error.message || String(error)}\n`);
    process.exit(1);
  }
}

if (require.main === module) {
  main();
}
