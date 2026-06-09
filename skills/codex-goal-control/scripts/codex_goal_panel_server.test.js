const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const {
  createGoalPanel,
  isValidJsonContentType,
  verifyMutationRequest,
} = require("./goal_panel");

function request(headers = {}) {
  return { headers, method: "GET", url: "/" };
}

function response() {
  let resolveEnd;
  const ended = new Promise((resolve) => {
    resolveEnd = resolve;
  });
  const res = {
    body: null,
    headers: null,
    statusCode: null,
    ended,
    writeHead(statusCode, headers = {}) {
      res.statusCode = statusCode;
      res.headers = headers;
    },
    end(body = "") {
      res.body = Buffer.isBuffer(body) ? body : Buffer.from(String(body), "utf8");
      resolveEnd(res);
    },
  };
  return res;
}

function panelSession(overrides = {}) {
  return {
    createLifecycleState: () => ({
      clients: new Map(),
      csrfToken: "csrf-token",
      idleShutdownMs: 1000,
      idleSweepMs: 250,
      lastActivityAt: 1000,
      server: null,
      sweepTimer: null,
      shuttingDown: false,
    }),
    forgetClient: () => ({ activeClients: 0 }),
    markClientSeen: () => ({ activeClients: 1 }),
    pruneClients: () => {},
    requireThreadId: () => "thread-1",
    resolveThreadId: () => ({ threadId: "thread-1", source: "server" }),
    resolveWorkspaceRoot: () => process.cwd(),
    writeClaim: (threadId) => ({ threadId }),
    ...overrides,
  };
}

test("Given a simple cross-origin POST When mutation guard runs Then it rejects the request", () => {
  const result = verifyMutationRequest(request({ "content-type": "text/plain" }), "token");

  assert.equal(result.statusCode, 403);
  assert.equal(result.payload.error, "Invalid CSRF token");
});

test("Given a valid token with text/plain When mutation guard runs Then it rejects the content type", () => {
  const result = verifyMutationRequest(
    request({
      "content-type": "text/plain",
      "x-codex-goal-csrf": "token",
    }),
    "token",
  );

  assert.equal(result.statusCode, 415);
  assert.equal(result.payload.error, "Request content-type must be application/json");
});

test("Given a valid token with JSON content When mutation guard runs Then it allows the request", () => {
  const result = verifyMutationRequest(
    request({
      "content-type": "application/json; charset=utf-8",
      "x-codex-goal-csrf": "token",
    }),
    "token",
  );

  assert.equal(result, null);
});

test("Given content type variants When parsed Then only JSON is accepted", () => {
  assert.equal(isValidJsonContentType("application/json"), true);
  assert.equal(isValidJsonContentType("application/json; charset=utf-8"), true);
  assert.equal(isValidJsonContentType("text/plain"), false);
  assert.equal(isValidJsonContentType(""), false);
});

test("Given the panel config API When handled Then thread config and CSRF token are returned", async () => {
  const panel = createGoalPanel(
    { host: "127.0.0.1", port: 43873, threadId: "thread-1", workspaceRoot: process.cwd() },
    { session: panelSession() },
  );
  const res = response();

  await panel.handleRequest({ headers: {}, method: "GET", url: "/api/config" }, res);
  await res.ended;

  assert.equal(res.statusCode, 200);
  assert.deepEqual(JSON.parse(res.body.toString("utf8")), {
    threadId: "thread-1",
    source: "server",
    csrfToken: "csrf-token",
  });
});

test("Given index HTML with icon links When served Then thread-aware icon URLs are rendered", async () => {
  const panelRoot = fs.mkdtempSync(path.join(os.tmpdir(), "codex-goal-panel-test-"));
  fs.writeFileSync(
    path.join(panelRoot, "index.html"),
    '<link href="/favicon.ico"><link href="/favicon.svg"><link href="/apple-touch-icon.png">',
    "utf8",
  );
  const panel = createGoalPanel(
    { host: "127.0.0.1", port: 43873, threadId: "thread-1", workspaceRoot: process.cwd() },
    { panelRoot, session: panelSession() },
  );
  const res = response();

  await panel.handleRequest({ headers: {}, method: "GET", url: "/" }, res);
  await res.ended;

  assert.equal(res.statusCode, 200);
  assert.equal(
    res.body.toString("utf8"),
    '<link href="/favicon.ico?threadId=thread-1"><link href="/favicon.svg?threadId=thread-1"><link href="/apple-touch-icon.png?threadId=thread-1">',
  );
});
