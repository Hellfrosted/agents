#!/usr/bin/env node

const { spawn } = require("child_process");

function createCodexAppServer(options = {}) {
  const {
    codexHome = process.env.CODEX_HOME || null,
    clientInfo = {},
    capabilities = { experimentalApi: true },
  } = options;

  const env = { ...process.env };
  if (codexHome) env.CODEX_HOME = codexHome;

  const child = spawn("codex", ["app-server"], {
    stdio: ["pipe", "pipe", "pipe"],
    env,
  });

  let nextId = 1;
  const pending = new Map();
  const notifications = [];
  let stdoutBuffer = "";
  let stderrBuffer = "";

  child.stdout.setEncoding("utf8");
  child.stderr.setEncoding("utf8");

  child.stdout.on("data", (chunk) => {
    stdoutBuffer += chunk;
    while (stdoutBuffer.includes("\n")) {
      const newlineIndex = stdoutBuffer.indexOf("\n");
      const line = stdoutBuffer.slice(0, newlineIndex).trim();
      stdoutBuffer = stdoutBuffer.slice(newlineIndex + 1);
      if (!line) continue;

      let parsed;
      try {
        parsed = JSON.parse(line);
      } catch (error) {
        notifications.push({
          method: "_parse_error",
          params: { line, error: String(error) },
        });
        continue;
      }

      if (Object.prototype.hasOwnProperty.call(parsed, "id")) {
        const waiter = pending.get(parsed.id);
        if (waiter) {
          pending.delete(parsed.id);
          if (parsed.error) {
            waiter.reject(new Error(parsed.error.message || JSON.stringify(parsed.error)));
          } else {
            waiter.resolve(parsed.result);
          }
        }
        continue;
      }

      notifications.push(parsed);
    }
  });

  child.stderr.on("data", (chunk) => {
    stderrBuffer += chunk;
  });

  child.on("exit", (code, signal) => {
    const reason = new Error(`app-server exited code=${code} signal=${signal}`);
    for (const waiter of pending.values()) {
      waiter.reject(reason);
    }
    pending.clear();
  });

  function request(method, params) {
    const id = nextId;
    nextId += 1;
    const payload = { jsonrpc: "2.0", id, method, params };
    return new Promise((resolve, reject) => {
      pending.set(id, { resolve, reject });
      child.stdin.write(`${JSON.stringify(payload)}\n`);
    });
  }

  async function initialize() {
    return request("initialize", {
      clientInfo: {
        name: clientInfo.name || "codex-goal-control-app-server-client",
        title: clientInfo.title || "Codex Goal Control App Server Client",
        version: clientInfo.version || "0.0.1",
      },
      capabilities,
    });
  }

  async function waitForNotification(method, predicate, timeoutMs) {
    const startedAt = Date.now();
    for (;;) {
      for (const entry of notifications) {
        if (entry.method !== method) continue;
        if (!predicate || predicate(entry.params || {})) return entry;
      }
      if (Date.now() - startedAt > timeoutMs) {
        throw new Error(`Timed out waiting for notification ${method}`);
      }
      await new Promise((resolve) => setTimeout(resolve, 50));
    }
  }

  async function shutdown() {
    if (!child.killed) {
      child.kill("SIGTERM");
      await new Promise((resolve) => setTimeout(resolve, 100));
      if (!child.killed) child.kill("SIGKILL");
    }
  }

  return {
    child,
    request,
    initialize,
    waitForNotification,
    shutdown,
    notifications,
    getStderr: () => stderrBuffer.trim() || null,
  };
}

module.exports = {
  createCodexAppServer,
};
