const fs = require("node:fs");
const crypto = require("node:crypto");
const os = require("node:os");
const path = require("node:path");
const { execFileSync: defaultExecFileSync } = require("node:child_process");
const { findAvailablePort, isServerReady } = require("./goal_panel_server_probe");

const DEFAULT_PANEL_DIR = path.join(os.tmpdir(), "codex-goal-panel");
const DEFAULT_IDLE_SHUTDOWN_MS = 15 * 60 * 1000;
const DEFAULT_IDLE_SWEEP_MS = 30 * 1000;

function envPositiveInteger(name, fallback) {
  const value = Number(process.env[name]);
  return Number.isInteger(value) && value > 0 ? value : fallback;
}

function createGoalPanelSession(options = {}) {
  const panelDir = options.panelDir || DEFAULT_PANEL_DIR;
  const cwd = options.cwd || process.cwd();
  const now = options.now || (() => new Date());
  const execFileSync = options.execFileSync || defaultExecFileSync;
  const paths = {
    claimDir: panelDir,
    claimFile: path.join(panelDir, "current-thread.json"),
    serverLogFile: path.join(panelDir, "server.log"),
    serverStateFile: path.join(panelDir, "server-state.json"),
  };

  function resolveWorkspace() {
    try {
      const root = execFileSync("git", ["-C", cwd, "rev-parse", "--show-toplevel"], {
        encoding: "utf8",
        stdio: ["ignore", "pipe", "ignore"],
        timeout: 1000,
      }).trim();
      if (root) return { root, cwd, source: "git" };
    } catch {
      // The panel can be opened from scratch directories outside a Git repo.
    }
    return { root: cwd, cwd, source: "cwd" };
  }

  function readClaim() {
    try {
      const raw = fs.readFileSync(paths.claimFile, "utf8");
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

  function writeClaim(threadId, source = "api", workspace = null) {
    const cleanThreadId = String(threadId || "").trim();
    if (!cleanThreadId) throw new Error("Missing thread id for claim.");
    fs.mkdirSync(paths.claimDir, { recursive: true });
    const claim = {
      threadId: cleanThreadId,
      source,
      updatedAt: now().toISOString(),
    };
    if (workspace) {
      claim.workspaceRoot = workspace.root;
      claim.workspaceCwd = workspace.cwd;
      claim.workspaceSource = workspace.source;
    }
    fs.writeFileSync(paths.claimFile, `${JSON.stringify(claim, null, 2)}\n`, "utf8");
    return claim;
  }

  function claimThread(threadId) {
    return writeClaim(threadId, "codex-thread", resolveWorkspace());
  }

  function resolveThreadId(url, serverOptions) {
    const explicit = url.searchParams.get("threadId");
    if (explicit) return { threadId: explicit, source: "query" };
    if (serverOptions.threadId) return { threadId: serverOptions.threadId, source: "server" };
    const claim = readClaim();
    if (claim) return { threadId: claim.threadId, source: claim.source, updatedAt: claim.updatedAt };
    return { threadId: null, source: "none" };
  }

  function requireThreadId(url, serverOptions) {
    const resolved = resolveThreadId(url, serverOptions);
    if (!resolved.threadId) {
      throw new Error("No Codex thread selected. Claim a thread first or pass ?threadId=<id>.");
    }
    return resolved.threadId;
  }

  function resolveWorkspaceRoot(serverOptions) {
    const claim = readClaim();
    const root = serverOptions.workspaceRoot || claim?.workspaceRoot || cwd;
    return path.resolve(root);
  }

  function readServerState() {
    try {
      const state = JSON.parse(fs.readFileSync(paths.serverStateFile, "utf8"));
      const port = Number(state.port);
      const host = typeof state.host === "string" ? state.host : "";
      if (!host || !Number.isInteger(port) || port <= 0 || port > 65535) return null;
      return { host, port };
    } catch {
      return null;
    }
  }

  function writeServerState(serverOptions, child) {
    const state = {
      host: serverOptions.host,
      port: serverOptions.port,
      pid: child.pid,
      startedAt: now().toISOString(),
    };
    fs.writeFileSync(paths.serverStateFile, `${JSON.stringify(state, null, 2)}\n`, "utf8");
    return state;
  }

  function createLifecycleState(overrides = {}) {
    const startedAt = overrides.startedAt || Date.now();
    return {
      clients: new Map(),
      csrfToken: overrides.csrfToken || crypto.randomBytes(32).toString("base64url"),
      idleShutdownMs:
        overrides.idleShutdownMs ||
        envPositiveInteger("CODEX_GOAL_PANEL_IDLE_SHUTDOWN_MS", DEFAULT_IDLE_SHUTDOWN_MS),
      idleSweepMs:
        overrides.idleSweepMs ||
        envPositiveInteger("CODEX_GOAL_PANEL_IDLE_SWEEP_MS", DEFAULT_IDLE_SWEEP_MS),
      lastActivityAt: startedAt,
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

  function markClientSeen(state, clientId, now = Date.now()) {
    const cleanClientId = String(clientId || "").trim();
    if (!cleanClientId) throw new Error("Missing dashboard client id.");
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

  async function ensureServer(serverOptions, adapters = {}) {
    const checkReady = adapters.isServerReady || isServerReady;
    const findPort = adapters.findAvailablePort || findAvailablePort;
    const startServer = adapters.startServer;
    let origin = `http://${serverOptions.host}:${serverOptions.port}`;
    const initial = await checkReady(origin);
    if (initial.ready && initial.iconReady) return "already-running";
    if (initial.ready && !initial.iconReady) {
      if (serverOptions.portWasExplicit) {
        throw new Error(
          `Goal panel server at ${origin} is running an older script without favicon support. Stop it or choose another --port.`,
        );
      }
      const state = readServerState();
      if (state && (state.host !== serverOptions.host || state.port !== serverOptions.port)) {
        const stateOrigin = `http://${state.host}:${state.port}`;
        const stateReady = await checkReady(stateOrigin);
        if (stateReady.ready && stateReady.iconReady) {
          serverOptions.host = state.host;
          serverOptions.port = state.port;
          return "already-running-fallback-port";
        }
      }
      serverOptions.port = await findPort(serverOptions.host, serverOptions.port + 1);
      origin = `http://${serverOptions.host}:${serverOptions.port}`;
    }
    startServer(serverOptions);
    for (let i = 0; i < 20; i += 1) {
      const next = await checkReady(origin);
      if (next.ready) return "started";
      if (next.error && String(next.error.message || next.error).includes("EPERM")) {
        return "started-unverified-sandbox";
      }
      await new Promise((resolve) => setTimeout(resolve, 100));
    }
    throw new Error("Goal panel server did not become ready.");
  }

  return {
    paths,
    claimThread,
    createLifecycleState,
    ensureServer,
    findAvailablePort,
    forgetClient,
    isServerReady,
    markClientSeen,
    pruneClients,
    readClaim,
    readServerState,
    requireThreadId,
    resolveThreadId,
    resolveWorkspace,
    resolveWorkspaceRoot,
    writeClaim,
    writeServerState,
  };
}

module.exports = {
  createGoalPanelSession,
};
