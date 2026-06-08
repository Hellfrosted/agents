const { spawn } = require("node:child_process");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

const { createPathTranslator } = require("./codex-wsl-path-translation");
const { createSkillsFallback } = require("./codex-wsl-skills-fallback");

function startProxy() {
  const debugLog = createDebugLogger(process.env.CODEX_WSL_PROXY_DEBUG_LOG || "");
  const translator = createPathTranslator({
    distroName: process.env.CODEX_WSL_PROXY_DISTRO || process.env.WSL_DISTRO_NAME || "",
    debugLog,
  });
  const skillsFallback = createSkillsFallback({
    windowsPathToWsl: translator.windowsPathToWsl,
  });
  const rawArgv = process.argv.slice(2);
  const normalizedArgv = rawArgv.map((arg) => translator.windowsPathToWsl(arg));
  const needsAppServer = normalizedArgv.length === 0 && !process.stdin.isTTY;
  const childArgv = needsAppServer ? ["app-server"] : normalizedArgv;
  const childEnv = buildChildEnv(translator.windowsPathToWsl);
  const child = spawn(resolveCodexTarget(), childArgv, {
    cwd: translator.windowsPathToWsl(process.env.T3CODE_WINDOWS_CWD) || os.homedir() || process.env.HOME || "/",
    env: childEnv,
    detached: true,
    stdio: ["pipe", "pipe", "pipe"],
  });

  const runtime = new ProxyRuntime({
    child,
    childArgv,
    debugLog,
    idleTimeoutMs: parseNonNegativeInteger(process.env.CODEX_WSL_PROXY_IDLE_TIMEOUT_MS, 0),
    normalizeInboundJsonLine: translator.normalizeInboundJsonLine,
    normalizeOutboundJsonLine: translator.normalizeOutboundJsonLine,
    skillsFallback,
    skillsFallbackTimeoutMs: parseNonNegativeInteger(process.env.CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS, 2000),
  });
  runtime.attach();
}

class ProxyRuntime {
  constructor(options) {
    Object.assign(this, options);
    this.activeTurnIds = new Set();
    this.pendingSkillsListRequests = new Map();
    this.stdinBuffer = "";
    this.stdoutBuffer = "";
    this.shuttingDown = false;
    this.childExited = false;
    this.lastActivityAt = Date.now();
  }

  attach() {
    this.attachLifecycle(); this.attachInput(); this.attachOutput();
    this.child.stderr.pipe(process.stderr);
  }

  attachLifecycle() {
    if (this.childArgv[0] === "app-server" && this.idleTimeoutMs > 0) {
      setInterval(() => this.reapIdleChild(), Math.min(this.idleTimeoutMs, 60_000)).unref();
    }

    process.once("SIGINT", () => this.shutdown("SIGINT"));
    process.once("SIGTERM", () => this.shutdown("SIGTERM"));
    process.once("SIGHUP", () => this.shutdown("SIGHUP"));
    this.child.stdin.on("error", () => {});
    this.child.on("error", (error) => failStart(error));
    this.child.on("exit", (code, signal) => this.exitFromChild(code, signal));
  }

  attachInput() {
    process.stdin.setEncoding("utf8");
    process.stdin.on("data", (chunk) => {
      this.recordActivity();
      this.stdinBuffer += chunk;
      this.flushInputLines();
    });
    process.stdin.on("end", () => {
      if (this.stdinBuffer.length > 0) this.forwardInputLine(this.stdinBuffer);
      this.child.stdin.end();
    });
    process.stdin.on("close", () => {
      if (!this.child.killed && this.childArgv[0] === "app-server") this.shutdown("SIGTERM");
    });
  }

  attachOutput() {
    this.child.stdout.setEncoding("utf8");
    this.child.stdout.on("data", (chunk) => {
      this.recordActivity();
      this.stdoutBuffer += chunk;
      this.flushOutputLines();
    });
    this.child.stdout.on("end", () => {
      if (this.stdoutBuffer.length > 0 && !this.handleChildJsonLine(this.stdoutBuffer)) {
        process.stdout.write(this.normalizeOutboundJsonLine(this.stdoutBuffer));
      }
    });
  }

  flushInputLines() {
    let newlineIndex;
    while ((newlineIndex = this.stdinBuffer.indexOf("\n")) !== -1) {
      const line = this.stdinBuffer.slice(0, newlineIndex + 1);
      this.stdinBuffer = this.stdinBuffer.slice(newlineIndex + 1);
      this.forwardInputLine(line);
    }
  }

  flushOutputLines() {
    let newlineIndex;
    while ((newlineIndex = this.stdoutBuffer.indexOf("\n")) !== -1) {
      const line = this.stdoutBuffer.slice(0, newlineIndex + 1);
      this.stdoutBuffer = this.stdoutBuffer.slice(newlineIndex + 1);
      if (!this.handleChildJsonLine(line)) process.stdout.write(this.normalizeOutboundJsonLine(line));
    }
  }

  forwardInputLine(line) {
    const normalizedLine = this.normalizeInboundJsonLine(line);
    try {
      const parsed = JSON.parse(normalizedLine);
      this.observeProtocolMessage(parsed);
      if (parsed?.method === "skills/list" && parsed?.id !== undefined) this.registerSkillsListFallback(parsed);
    } catch {
      // Non-JSON input still belongs to child stdin.
    }
    this.child.stdin.write(normalizedLine);
  }

  observeProtocolMessage(message) {
    if (!message || typeof message !== "object") return;
    this.recordActivity();
    if (this.childArgv[0] !== "app-server") return;
    if (message.method === "turn/started") {
      const turnId = readMessageTurnId(message);
      if (turnId) this.activeTurnIds.add(turnId);
      return;
    }
    if (message.method === "turn/completed") {
      const turnId = readMessageTurnId(message);
      if (turnId) this.activeTurnIds.delete(turnId);
    }
  }

  registerSkillsListFallback(message) {
    const timer = setTimeout(() => {
      const pending = this.pendingSkillsListRequests.get(message.id);
      if (!pending) return;
      pending.responded = true;
      const response = `${JSON.stringify(this.skillsFallback.makeResponse(message))}\n`;
      process.stdout.write(this.normalizeOutboundJsonLine(response));
    }, this.skillsFallbackTimeoutMs);
    this.pendingSkillsListRequests.set(message.id, { timer, responded: false });
  }

  handleChildJsonLine(line) {
    let message;
    try {
      message = JSON.parse(line);
    } catch {
      return false;
    }

    this.observeProtocolMessage(message);
    const pending = this.pendingSkillsListRequests.get(message?.id);
    if (!pending) return false;
    clearTimeout(pending.timer);
    this.pendingSkillsListRequests.delete(message.id);
    if (pending.responded) {
      this.debugLog("skills-fallback", `suppressed upstream response for id=${message.id}`);
      return true;
    }
    process.stdout.write(this.normalizeOutboundJsonLine(line));
    return true;
  }

  reapIdleChild() {
    if (this.shuttingDown || this.childExited || this.activeTurnIds.size > 0) return;
    const idleForMs = Date.now() - this.lastActivityAt;
    if (idleForMs < this.idleTimeoutMs) return;
    this.debugLog("idle-reaper", `app-server idle for ${idleForMs}ms; shutting down`);
    this.shutdown("SIGTERM");
  }

  shutdown(signal = "SIGTERM") {
    if (this.shuttingDown) return;
    this.shuttingDown = true;
    this.clearPendingSkillsListFallbacks();
    if (!this.child.killed) tryKill(this.child, signal);
    setTimeout(() => {
      if (!this.childExited) tryKill(this.child, "SIGKILL");
    }, 5000).unref();
  }

  exitFromChild(code, signal) {
    this.childExited = true;
    this.clearPendingSkillsListFallbacks();
    if (signal) {
      if (this.shuttingDown) process.exit(0);
      process.kill(process.pid, signal);
      return;
    }
    process.exit(code ?? 0);
  }

  recordActivity() {
    this.lastActivityAt = Date.now();
  }

  clearPendingSkillsListFallbacks() {
    for (const pending of this.pendingSkillsListRequests.values()) clearTimeout(pending.timer);
    this.pendingSkillsListRequests.clear();
  }
}

function buildChildEnv(windowsPathToWsl) {
  const childEnv = { ...process.env };
  delete childEnv.T3CODE_WINDOWS_CWD; delete childEnv.CODEX_WSL_PROXY_IDLE_TIMEOUT_MS;
  if (typeof childEnv.CODEX_HOME === "string") childEnv.CODEX_HOME = windowsPathToWsl(childEnv.CODEX_HOME);
  return childEnv;
}

function resolveCodexTarget() {
  return process.env.CODEX_WSL_PROXY_TARGET || resolveExecutableOnPath("codex") || path.join(path.dirname(process.execPath), "codex");
}

function failStart(error) {
  console.error(`codex-wsl-proxy: failed to start ${resolveCodexTarget()}: ${error.message}`);
  process.exit(127);
}

function resolveExecutableOnPath(name) {
  const pathEnv = process.env.PATH || "";
  for (const entry of pathEnv.split(path.delimiter)) {
    if (!entry) continue;
    const candidate = path.join(entry, name);
    try {
      fs.accessSync(candidate, fs.constants.X_OK);
      return candidate;
    } catch {}
  }
  return "";
}

function parseNonNegativeInteger(raw, fallback) {
  if (raw === undefined || raw === "") return fallback;
  const parsed = Number(raw);
  return Number.isFinite(parsed) && parsed >= 0 ? parsed : fallback;
}

function readMessageTurnId(message) {
  const params = message?.params;
  if (!params || typeof params !== "object") return undefined;
  if (typeof params.turnId === "string") return params.turnId;
  if (params.turn && typeof params.turn === "object" && typeof params.turn.id === "string") return params.turn.id;
  return undefined;
}

function createDebugLogger(debugLogPath) {
  return (stream, payload) => {
    if (!debugLogPath) return;
    try {
      fs.appendFileSync(debugLogPath, `[${new Date().toISOString()}] ${stream} ${String(payload).replace(/\n$/, "")}\n`, "utf8");
    } catch {}
  };
}

function tryKill(child, signal) {
  try {
    child.kill(signal);
  } catch {}
}

module.exports = { ProxyRuntime, buildChildEnv, parseNonNegativeInteger, readMessageTurnId, startProxy };
