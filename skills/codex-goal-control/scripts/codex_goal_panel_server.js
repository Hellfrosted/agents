#!/usr/bin/env node

const fs = require("fs");
const http = require("http");
const path = require("path");
const { URL } = require("url");
const { runGoalCommand } = require("./codex_goal");

const SKILL_ROOT = path.resolve(__dirname, "..");
const PANEL_ROOT = path.join(SKILL_ROOT, "assets", "panel");
const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 43873;
const DEFAULT_IDLE_SHUTDOWN_MS = 15 * 60 * 1000;
const DEFAULT_IDLE_SWEEP_MS = 30 * 1000;
const FORCE_EXIT_MS = 1000;
const CLAIM_DIR = path.join(require("os").tmpdir(), "codex-goal-panel");
const CLAIM_FILE = path.join(CLAIM_DIR, "current-thread.json");
const ICON_MAX_BYTES = 512 * 1024;
const ICON_CANDIDATES = [
  "favicon.ico",
  "favicon.svg",
  "favicon.png",
  "icon.svg",
  "icon.png",
  "public/favicon.ico",
  "public/favicon.svg",
  "public/favicon.png",
  "public/icon.svg",
  "public/icon.png",
  "public/apple-touch-icon.png",
  "app/favicon.ico",
  "app/icon.svg",
  "app/icon.png",
  "src/app/favicon.ico",
  "src/app/icon.svg",
  "src/app/icon.png",
  "assets/favicon.ico",
  "assets/favicon.svg",
  "assets/favicon.png",
  "assets/icon.svg",
  "assets/icon.png",
  "static/favicon.ico",
  "static/favicon.svg",
  "static/favicon.png",
];
const THREAD_BADGE_PALETTE = [
  "#06B6D4",
  "#22C55E",
  "#F59E0B",
  "#EC4899",
  "#8B5CF6",
  "#EF4444",
  "#14B8A6",
  "#EAB308",
];
const GOAL_STAGE_BADGES = {
  active: { fill: "#2EE86F", icon: "play" },
  paused: { fill: "#F6F3ED", icon: "pause" },
  budgetLimited: { fill: "#F59E0B", icon: "alert" },
  complete: { fill: "#FF2638", icon: "check" },
  none: { fill: "#545454", icon: "dash" },
  unknown: { fill: "#8A8A8A", icon: "unknown" },
};

function envPositiveInteger(name, fallback) {
  const value = Number(process.env[name]);
  return Number.isInteger(value) && value > 0 ? value : fallback;
}

function usage(exitCode = 1) {
  const text = [
    "Usage:",
    "  node scripts/codex_goal_panel_server.js [--thread <thread-id>] [--host 127.0.0.1] [--port 43873]",
    "",
    "Serves a local Codex Goal panel for the current thread.",
    "The server binds to 127.0.0.1 by default. When --thread is omitted, it uses the latest claimed thread.",
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

function readClaim() {
  try {
    const raw = fs.readFileSync(CLAIM_FILE, "utf8");
    const claim = JSON.parse(raw);
    const threadId = typeof claim.threadId === "string" ? claim.threadId.trim() : "";
    if (!threadId) return null;
    return {
      threadId,
      updatedAt: claim.updatedAt || null,
      source: claim.source || "claim",
      workspaceRoot:
        typeof claim.workspaceRoot === "string" && claim.workspaceRoot.trim()
          ? claim.workspaceRoot.trim()
          : null,
      workspaceCwd:
        typeof claim.workspaceCwd === "string" && claim.workspaceCwd.trim()
          ? claim.workspaceCwd.trim()
          : null,
      workspaceSource: claim.workspaceSource || null,
    };
  } catch {
    return null;
  }
}

function writeClaim(threadId, source = "api") {
  const cleanThreadId = String(threadId || "").trim();
  if (!cleanThreadId) throw new Error("Missing thread id for claim.");
  fs.mkdirSync(CLAIM_DIR, { recursive: true });
  const claim = {
    threadId: cleanThreadId,
    source,
    updatedAt: new Date().toISOString(),
  };
  fs.writeFileSync(CLAIM_FILE, `${JSON.stringify(claim, null, 2)}\n`, "utf8");
  return claim;
}

function resolveThreadId(url, options) {
  const explicit = url.searchParams.get("threadId");
  if (explicit) return { threadId: explicit, source: "query" };
  if (options.threadId) return { threadId: options.threadId, source: "server" };
  const claim = readClaim();
  if (claim) return { threadId: claim.threadId, source: claim.source, updatedAt: claim.updatedAt };
  return { threadId: null, source: "none" };
}

function requireThreadId(url, options) {
  const resolved = resolveThreadId(url, options);
  if (!resolved.threadId) {
    throw new Error("No Codex thread selected. Claim a thread first or pass ?threadId=<id>.");
  }
  return resolved.threadId;
}

function sendJson(res, statusCode, payload) {
  const body = JSON.stringify(payload, null, 2);
  res.writeHead(statusCode, {
    "content-type": "application/json; charset=utf-8",
    "cache-control": "no-store",
    "content-length": Buffer.byteLength(body),
  });
  res.end(body);
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let body = "";
    req.setEncoding("utf8");
    req.on("data", (chunk) => {
      body += chunk;
      if (body.length > 1024 * 1024) {
        reject(new Error("Request body too large"));
        req.destroy();
      }
    });
    req.on("end", () => {
      if (!body.trim()) {
        resolve({});
        return;
      }
      try {
        resolve(JSON.parse(body));
      } catch {
        reject(new Error("Request body must be JSON"));
      }
    });
    req.on("error", reject);
  });
}

function createLifecycleState() {
  return {
    clients: new Map(),
    idleShutdownMs: envPositiveInteger("CODEX_GOAL_PANEL_IDLE_SHUTDOWN_MS", DEFAULT_IDLE_SHUTDOWN_MS),
    idleSweepMs: envPositiveInteger("CODEX_GOAL_PANEL_IDLE_SWEEP_MS", DEFAULT_IDLE_SWEEP_MS),
    lastActivityAt: Date.now(),
    server: null,
    sweepTimer: null,
    shuttingDown: false,
  };
}

function pruneClients(state, now = Date.now()) {
  for (const [clientId, seenAt] of state.clients) {
    if (now - seenAt > state.idleShutdownMs) state.clients.delete(clientId);
  }
}

function markClientSeen(state, clientId) {
  const cleanClientId = String(clientId || "").trim();
  if (!cleanClientId) throw new Error("Missing dashboard client id.");
  const now = Date.now();
  state.clients.set(cleanClientId, now);
  state.lastActivityAt = now;
  return {
    clientId: cleanClientId,
    activeClients: state.clients.size,
    idleShutdownSeconds: Math.round(state.idleShutdownMs / 1000),
  };
}

function forgetClient(state, clientId) {
  const cleanClientId = String(clientId || "").trim();
  if (cleanClientId) state.clients.delete(cleanClientId);
  if (state.clients.size === 0) state.lastActivityAt = Date.now();
  return { clientId: cleanClientId || null, activeClients: state.clients.size };
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

function startIdleShutdownTimer(state) {
  state.sweepTimer = setInterval(() => {
    const now = Date.now();
    pruneClients(state, now);
    if (state.clients.size === 0 && now - state.lastActivityAt >= state.idleShutdownMs) {
      shutdownServer(state, "idle-timeout");
    }
  }, state.idleSweepMs);
  state.sweepTimer.unref();
}

function contentTypeFor(filePath) {
  if (filePath.endsWith(".html")) return "text/html; charset=utf-8";
  if (filePath.endsWith(".css")) return "text/css; charset=utf-8";
  if (filePath.endsWith(".js")) return "application/javascript; charset=utf-8";
  if (filePath.endsWith(".json")) return "application/json; charset=utf-8";
  if (filePath.endsWith(".svg")) return "image/svg+xml; charset=utf-8";
  if (filePath.endsWith(".png")) return "image/png";
  if (filePath.endsWith(".ico")) return "image/x-icon";
  if (filePath.endsWith(".jpg") || filePath.endsWith(".jpeg")) return "image/jpeg";
  if (filePath.endsWith(".webp")) return "image/webp";
  return "application/octet-stream";
}

function resolveWorkspaceRoot(options) {
  const claim = readClaim();
  const root = claim?.workspaceRoot || options.workspaceRoot || process.cwd();
  return path.resolve(root);
}

function isInsidePath(childPath, parentPath) {
  const relative = path.relative(parentPath, childPath);
  return Boolean(relative && !relative.startsWith("..") && !path.isAbsolute(relative));
}

function findRepoIcon(workspaceRoot) {
  let realWorkspaceRoot = null;
  try {
    realWorkspaceRoot = fs.realpathSync(workspaceRoot);
  } catch {
    return null;
  }

  for (const candidate of ICON_CANDIDATES) {
    const filePath = path.resolve(realWorkspaceRoot, candidate);
    if (!isInsidePath(filePath, realWorkspaceRoot)) continue;
    let stat = null;
    try {
      stat = fs.statSync(filePath);
    } catch {
      continue;
    }
    if (!stat.isFile()) continue;
    try {
      const realFilePath = fs.realpathSync(filePath);
      if (!isInsidePath(realFilePath, realWorkspaceRoot)) continue;
      return realFilePath;
    } catch {
      continue;
    }
  }
  return null;
}

function escapeXml(value) {
  return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function hashString(value) {
  let hash = 0;
  for (const char of String(value || "")) {
    hash = (hash * 31 + char.charCodeAt(0)) >>> 0;
  }
  return hash;
}

function shortThreadLabel(threadId) {
  const clean = String(threadId || "").replace(/[^a-z0-9]/gi, "").toUpperCase();
  if (!clean) return "--";
  return clean.slice(-2).padStart(2, "0");
}

function threadBadgeFor(threadId) {
  const hash = hashString(threadId);
  return {
    label: shortThreadLabel(threadId),
    fill: THREAD_BADGE_PALETTE[hash % THREAD_BADGE_PALETTE.length],
  };
}

function normalizeGoalStatus(status) {
  const clean = String(status || "").trim();
  if (clean === "active" || clean === "paused" || clean === "budgetLimited" || clean === "complete") {
    return clean;
  }
  if (clean === "none") return "none";
  return "unknown";
}

async function goalStatusForIcon(url, options, threadId) {
  const explicitStatus = url.searchParams.get("goalStatus");
  if (explicitStatus) return normalizeGoalStatus(explicitStatus);
  if (!threadId) return "none";

  try {
    const result = await runGoalCommand({ command: "get", threadId, json: true });
    return result.goal ? normalizeGoalStatus(result.goal.status) : "none";
  } catch {
    return "unknown";
  }
}

function repoIconData(iconPath) {
  if (!iconPath) return null;
  let stat = null;
  try {
    stat = fs.statSync(iconPath);
  } catch {
    return null;
  }
  if (!stat.isFile() || stat.size > ICON_MAX_BYTES) return null;

  try {
    const mime = contentTypeFor(iconPath).split(";")[0];
    const data = fs.readFileSync(iconPath);
    return {
      href: `data:${mime};base64,${data.toString("base64")}`,
      source: iconPath,
    };
  } catch {
    return null;
  }
}

function fallbackIconMarkup(workspaceRoot) {
  const label = path.basename(workspaceRoot || "Goal") || "Goal";
  const initial = (label.match(/[a-z0-9]/i)?.[0] || "G").toUpperCase();
  return [
    "<rect width=\"64\" height=\"64\" rx=\"14\" fill=\"#111827\"/>",
    "<path d=\"M14 22L24 32L14 42\" fill=\"none\" stroke=\"#34D399\" stroke-width=\"5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/>",
    "<path d=\"M27 44H36\" fill=\"none\" stroke=\"#34D399\" stroke-width=\"5\" stroke-linecap=\"round\"/>",
    `<text x="42" y="38" text-anchor="middle" font-family="Arial, sans-serif" font-size="26" font-weight="700" fill="#A78BFA">${escapeXml(initial)}</text>`,
  ].join("");
}

function goalStageIconMarkup(stageBadge) {
  const stroke = "#050505";
  const fill = "#050505";
  if (stageBadge.icon === "play") {
    return `<path d="M45 40L45 56L57 48Z" fill="${fill}"/>`;
  }
  if (stageBadge.icon === "pause") {
    return `<path d="M43 40V56M53 40V56" fill="none" stroke="${stroke}" stroke-width="5" stroke-linecap="round"/>`;
  }
  if (stageBadge.icon === "alert") {
    return [
      `<path d="M48 39V50" fill="none" stroke="${stroke}" stroke-width="5" stroke-linecap="round"/>`,
      `<circle cx="48" cy="56" r="2.7" fill="${fill}"/>`,
    ].join("");
  }
  if (stageBadge.icon === "check") {
    return `<path d="M40 48L46 54L57 42" fill="none" stroke="${stroke}" stroke-width="5" stroke-linecap="round" stroke-linejoin="round"/>`;
  }
  if (stageBadge.icon === "dash") {
    return `<path d="M40 48H56" fill="none" stroke="${stroke}" stroke-width="5" stroke-linecap="round"/>`;
  }
  return [
    `<path d="M43 43C44 40 47 39 50 40C53 41 54 44 52 47C51 49 48 50 48 53" fill="none" stroke="${stroke}" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>`,
    `<circle cx="48" cy="57" r="2.3" fill="${fill}"/>`,
  ].join("");
}

function panelIcon(workspaceRoot, threadId, goalStatus, repoIcon) {
  const threadBadge = threadBadgeFor(threadId);
  const stageBadge = GOAL_STAGE_BADGES[normalizeGoalStatus(goalStatus)] || GOAL_STAGE_BADGES.unknown;
  const baseIcon = repoIcon
    ? `<image href="${repoIcon.href}" x="0" y="0" width="64" height="64" preserveAspectRatio="xMidYMid slice" clip-path="url(#icon-clip)"/>`
    : fallbackIconMarkup(workspaceRoot);

  return Buffer.from(
    [
      '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">',
      "<defs><clipPath id=\"icon-clip\"><rect width=\"64\" height=\"64\" rx=\"14\"/></clipPath></defs>",
      baseIcon,
      "<rect width=\"64\" height=\"64\" rx=\"14\" fill=\"none\" stroke=\"rgba(255,255,255,0.22)\"/>",
      "<circle cx=\"16\" cy=\"16\" r=\"15\" fill=\"#050505\" stroke=\"#F6F3ED\" stroke-width=\"3\"/>",
      `<circle cx="16" cy="16" r="11.5" fill="${threadBadge.fill}"/>`,
      `<text x="16" y="20.2" text-anchor="middle" font-family="Arial, sans-serif" font-size="10.5" font-weight="900" fill="#F6F3ED" stroke="#050505" stroke-width="1.9" paint-order="stroke">${escapeXml(threadBadge.label)}</text>`,
      "<circle cx=\"48\" cy=\"48\" r=\"15\" fill=\"#050505\" stroke=\"#F6F3ED\" stroke-width=\"3\"/>",
      `<circle cx="48" cy="48" r="11.5" fill="${stageBadge.fill}"/>`,
      goalStageIconMarkup(stageBadge),
      "</svg>",
    ].join(""),
    "utf8",
  );
}

function pageHtml(data, url, options) {
  const resolved = resolveThreadId(url, options);
  if (!resolved.threadId) return data;

  const iconQuery = `?threadId=${encodeURIComponent(resolved.threadId)}`;
  return Buffer.from(
    data
      .toString("utf8")
      .replaceAll('href="/favicon.ico"', `href="/favicon.ico${iconQuery}"`)
      .replaceAll('href="/favicon.svg"', `href="/favicon.svg${iconQuery}"`)
      .replaceAll('href="/apple-touch-icon.png"', `href="/apple-touch-icon.png${iconQuery}"`),
    "utf8",
  );
}

async function serveRepoIcon(req, res, options) {
  const url = new URL(req.url, `http://${options.host}:${options.port}`);
  const workspaceRoot = resolveWorkspaceRoot(options);
  const iconPath = findRepoIcon(workspaceRoot);
  const resolved = resolveThreadId(url, options);
  const threadId = resolved.threadId || "";
  const goalStatus = await goalStatusForIcon(url, options, threadId);
  const repoIcon = repoIconData(iconPath);
  const data = panelIcon(workspaceRoot, threadId, goalStatus, repoIcon);
  const headers = {
    "content-type": "image/svg+xml; charset=utf-8",
    "cache-control": "no-store",
    "content-length": data.length,
    "x-codex-goal-thread": threadId || "none",
    "x-codex-goal-status": goalStatus,
  };
  if (repoIcon?.source) headers["x-codex-goal-icon-source"] = repoIcon.source;
  res.writeHead(200, headers);
  res.end(data);
}

async function serveStatic(req, res, options) {
  const url = new URL(req.url, "http://localhost");
  if (
    url.pathname === "/favicon.ico" ||
    url.pathname === "/favicon.svg" ||
    url.pathname === "/apple-touch-icon.png"
  ) {
    await serveRepoIcon(req, res, options);
    return;
  }
  const rawPath = url.pathname === "/" ? "/index.html" : url.pathname;
  const normalized = path.normalize(rawPath).replace(/^(\.\.[/\\])+/, "");
  const filePath = path.join(PANEL_ROOT, normalized);
  if (!filePath.startsWith(PANEL_ROOT)) {
    res.writeHead(403);
    res.end("Forbidden");
    return;
  }
  fs.readFile(filePath, (error, data) => {
    if (error) {
      res.writeHead(404);
      res.end("Not found");
      return;
    }
    const body = filePath.endsWith("index.html") ? pageHtml(data, url, options) : data;
    res.writeHead(200, {
      "content-type": contentTypeFor(filePath),
      "cache-control": "no-store",
      "content-length": body.length,
    });
    res.end(body);
  });
}

async function handleApi(req, res, options, state) {
  const url = new URL(req.url, `http://${options.host}:${options.port}`);

  if (url.pathname === "/api/config" && req.method === "GET") {
    sendJson(res, 200, resolveThreadId(url, options));
    return;
  }

  if (url.pathname === "/api/claim" && req.method === "POST") {
    const body = await readBody(req);
    sendJson(res, 200, writeClaim(body.threadId, "api"));
    return;
  }

  if (url.pathname === "/api/session/heartbeat" && req.method === "POST") {
    const body = await readBody(req);
    sendJson(res, 200, markClientSeen(state, body.clientId));
    return;
  }

  if (url.pathname === "/api/session/close" && (req.method === "POST" || req.method === "DELETE")) {
    const body = await readBody(req);
    sendJson(res, 200, forgetClient(state, body.clientId));
    return;
  }

  if (url.pathname === "/api/server/stop" && req.method === "POST") {
    const body = await readBody(req);
    forgetClient(state, body.clientId);
    sendJson(res, 200, { stopping: true, reason: "api-stop" });
    setTimeout(() => shutdownServer(state, "api-stop"), 25).unref();
    return;
  }

  if (url.pathname === "/api/goal" && req.method === "GET") {
    const threadId = requireThreadId(url, options);
    sendJson(res, 200, await runGoalCommand({ command: "get", threadId, json: true }));
    return;
  }

  if (url.pathname === "/api/goal" && req.method === "DELETE") {
    const threadId = requireThreadId(url, options);
    sendJson(res, 200, await runGoalCommand({ command: "clear", threadId, json: true }));
    return;
  }

  if (url.pathname === "/api/goal" && req.method === "POST") {
    const threadId = requireThreadId(url, options);
    const body = await readBody(req);
    const objective = String(body.objective || "").trim();
    if (!objective) throw new Error("Missing objective");
    const tokenBudget =
      body.tokenBudget === null || body.tokenBudget === undefined || body.tokenBudget === ""
        ? undefined
        : Number(body.tokenBudget);
    if (tokenBudget !== undefined && (!Number.isInteger(tokenBudget) || tokenBudget <= 0)) {
      throw new Error("tokenBudget must be a positive integer");
    }
    sendJson(
      res,
      200,
      await runGoalCommand({
        command: "set",
        threadId,
        objective,
        tokenBudget,
        json: true,
      }),
    );
    return;
  }

  if (url.pathname === "/api/goal/status" && req.method === "POST") {
    const threadId = requireThreadId(url, options);
    const body = await readBody(req);
    const statusCommand = {
      active: "resume",
      paused: "pause",
      complete: "complete",
    }[String(body.status || "")];
    if (!statusCommand) throw new Error("status must be active, paused, or complete");
    sendJson(res, 200, await runGoalCommand({ command: statusCommand, threadId, json: true }));
    return;
  }

  sendJson(res, 404, { error: "not found" });
}

function createServer(options) {
  const state = createLifecycleState();
  const server = http.createServer(async (req, res) => {
    try {
      if (req.url.startsWith("/api/")) {
        await handleApi(req, res, options, state);
        return;
      }
      await serveStatic(req, res, options);
    } catch (error) {
      sendJson(res, 500, { error: error.message || String(error) });
    }
  });
  state.server = server;
  startIdleShutdownTimer(state);
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
            claim_file: CLAIM_FILE,
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
