#!/usr/bin/env node

const fs = require("fs");
const http = require("http");
const net = require("net");
const os = require("os");
const path = require("path");
const { execFileSync, spawn } = require("child_process");

const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 43873;
const CLAIM_DIR = path.join(os.tmpdir(), "codex-goal-panel");
const CLAIM_FILE = path.join(CLAIM_DIR, "current-thread.json");
const SERVER_LOG_FILE = path.join(CLAIM_DIR, "server.log");
const SERVER_STATE_FILE = path.join(CLAIM_DIR, "server-state.json");
const SERVER_SCRIPT = path.join(__dirname, "codex_goal_panel_server.js");

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

function resolveWorkspace() {
  const cwd = process.cwd();
  try {
    const root = execFileSync("git", ["-C", cwd, "rev-parse", "--show-toplevel"], {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "ignore"],
      timeout: 1000,
    }).trim();
    if (root) return { root, cwd, source: "git" };
  } catch {
    // This helper is often run from scratch directories that are not Git repos.
  }
  return { root: cwd, cwd, source: "cwd" };
}

function claimThread(threadId) {
  fs.mkdirSync(CLAIM_DIR, { recursive: true });
  const workspace = resolveWorkspace();
  const claim = {
    threadId,
    source: "codex-thread",
    workspaceRoot: workspace.root,
    workspaceCwd: workspace.cwd,
    workspaceSource: workspace.source,
    updatedAt: new Date().toISOString(),
  };
  fs.writeFileSync(CLAIM_FILE, `${JSON.stringify(claim, null, 2)}\n`, "utf8");
  return claim;
}

function requestJson(url) {
  return new Promise((resolve, reject) => {
    const req = http.get(url, { timeout: 1000 }, (res) => {
      let body = "";
      res.setEncoding("utf8");
      res.on("data", (chunk) => {
        body += chunk;
      });
      res.on("end", () => {
        try {
          resolve({ ok: res.statusCode >= 200 && res.statusCode < 300, data: JSON.parse(body) });
        } catch (error) {
          reject(error);
        }
      });
    });
    req.on("timeout", () => {
      req.destroy(new Error("request timeout"));
    });
    req.on("error", reject);
  });
}

function requestOk(url) {
  return new Promise((resolve) => {
    const req = http.get(url, { timeout: 1000 }, (res) => {
      const ok = res.statusCode >= 200 && res.statusCode < 300;
      res.resume();
      res.on("end", () => resolve(ok));
    });
    req.on("timeout", () => {
      req.destroy(new Error("request timeout"));
    });
    req.on("error", () => resolve(false));
  });
}

async function isServerReady(origin) {
  try {
    const result = await requestJson(`${origin}/api/config`);
    if (!result.ok || !result.data) return { ready: false, iconReady: false, error: null };
    const iconReady = await requestOk(`${origin}/favicon.ico`);
    return { ready: true, iconReady, error: null };
  } catch (error) {
    return { ready: false, iconReady: false, error };
  }
}

function readServerState() {
  try {
    const state = JSON.parse(fs.readFileSync(SERVER_STATE_FILE, "utf8"));
    const port = Number(state.port);
    const host = typeof state.host === "string" ? state.host : "";
    if (!host || !Number.isInteger(port) || port <= 0 || port > 65535) return null;
    return { host, port };
  } catch {
    return null;
  }
}

function writeServerState(options, child) {
  const state = {
    host: options.host,
    port: options.port,
    pid: child.pid,
    startedAt: new Date().toISOString(),
  };
  fs.writeFileSync(SERVER_STATE_FILE, `${JSON.stringify(state, null, 2)}\n`, "utf8");
}

function canBind(host, port) {
  return new Promise((resolve) => {
    const server = net.createServer();
    server.once("error", () => resolve(false));
    server.once("listening", () => {
      server.close(() => resolve(true));
    });
    server.listen(port, host);
  });
}

async function findAvailablePort(host, startPort) {
  for (let port = startPort; port <= 65535; port += 1) {
    if (await canBind(host, port)) return port;
  }
  throw new Error("No available local port found for Goal panel server.");
}

function startServer(options, claim) {
  fs.mkdirSync(CLAIM_DIR, { recursive: true });
  const logFd = fs.openSync(SERVER_LOG_FILE, "a");
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
  writeServerState(options, child);
  child.unref();
}

async function ensureServer(options, claim) {
  let origin = `http://${options.host}:${options.port}`;
  const initial = await isServerReady(origin);
  if (initial.ready && initial.iconReady) return "already-running";
  if (initial.ready && !initial.iconReady) {
    if (options.portWasExplicit) {
      throw new Error(
        `Goal panel server at ${origin} is running an older script without favicon support. Stop it or choose another --port.`,
      );
    }
    const state = readServerState();
    if (state && (state.host !== options.host || state.port !== options.port)) {
      const stateOrigin = `http://${state.host}:${state.port}`;
      const stateReady = await isServerReady(stateOrigin);
      if (stateReady.ready && stateReady.iconReady) {
        options.host = state.host;
        options.port = state.port;
        return "already-running-fallback-port";
      }
    }
    options.port = await findAvailablePort(options.host, options.port + 1);
    origin = `http://${options.host}:${options.port}`;
  }
  startServer(options, claim);
  for (let i = 0; i < 20; i += 1) {
    const next = await isServerReady(origin);
    if (next.ready) return "started";
    if (next.error && String(next.error.message || next.error).includes("EPERM")) {
      return "started-unverified-sandbox";
    }
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  throw new Error("Goal panel server did not become ready.");
}

async function main() {
  try {
    const options = parseArgs(process.argv.slice(2));
    const claim = claimThread(options.threadId);
    const server = await ensureServer(options, claim);
    const origin = `http://${options.host}:${options.port}`;
    const threadUrl = `${origin}/?threadId=${encodeURIComponent(options.threadId)}`;
    const result = {
      threadId: options.threadId,
      url: threadUrl,
      rootUrl: origin,
      threadUrl,
      server,
      claimFile: CLAIM_FILE,
      serverLogFile: SERVER_LOG_FILE,
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
