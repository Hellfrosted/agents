const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const { createGoalPanelSession } = require("./goal_panel_session");

function tempSession() {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "goal-panel-session-test-"));
  const session = createGoalPanelSession({
    panelDir: tempDir,
    cwd: tempDir,
    now: () => new Date("2026-06-07T20:00:00.000Z"),
    execFileSync: () => {
      throw new Error("not a git repo");
    },
  });
  return { session, tempDir };
}

test("Given a thread id When claimed Then the session stores the thread and workspace", () => {
  const { session, tempDir } = tempSession();

  const claim = session.claimThread("thread-123");

  assert.equal(claim.threadId, "thread-123");
  assert.equal(claim.source, "codex-thread");
  assert.equal(claim.workspaceRoot, tempDir);
  assert.equal(claim.workspaceCwd, tempDir);
  assert.equal(claim.workspaceSource, "cwd");
  assert.equal(claim.updatedAt, "2026-06-07T20:00:00.000Z");
  assert.deepEqual(session.readClaim(), claim);
});

test("Given query, server option, and claim When resolving thread Then precedence is query before server before claim", () => {
  const { session } = tempSession();
  session.writeClaim("claim-thread", "api");

  assert.deepEqual(
    session.resolveThreadId(new URL("http://127.0.0.1/?threadId=query-thread"), { threadId: "server-thread" }),
    { threadId: "query-thread", source: "query" },
  );
  assert.deepEqual(
    session.resolveThreadId(new URL("http://127.0.0.1/"), { threadId: "server-thread" }),
    { threadId: "server-thread", source: "server" },
  );
  assert.equal(session.resolveThreadId(new URL("http://127.0.0.1/"), {}).threadId, "claim-thread");
});

test("Given a server process When state is written Then the session can read the listening address", () => {
  const { session } = tempSession();

  session.writeServerState({ host: "127.0.0.1", port: 43873 }, { pid: 12345 });

  assert.deepEqual(session.readServerState(), { host: "127.0.0.1", port: 43873 });
});

test("Given clients with stale activity When pruned Then only live session clients remain", () => {
  const { session } = tempSession();
  const state = session.createLifecycleState({
    idleShutdownMs: 1000,
    idleSweepMs: 250,
    csrfToken: "token",
    startedAt: 10_000,
  });

  session.markClientSeen(state, "live-client", 10_000);
  session.markClientSeen(state, "stale-client", 8_000);
  session.pruneClients(state, 10_000);

  assert.deepEqual([...state.clients.keys()], ["live-client"]);
});

test("Given an old default-port server and a ready saved port When ensuring server Then the saved port is reused", async () => {
  const { session } = tempSession();
  const options = { host: "127.0.0.1", port: 43873, portWasExplicit: false };
  session.writeServerState({ host: "127.0.0.1", port: 43874 }, { pid: 12345 });
  const seenOrigins = [];

  const result = await session.ensureServer(options, {
    isServerReady: async (origin) => {
      seenOrigins.push(origin);
      if (origin.endsWith(":43873")) return { ready: true, iconReady: false, error: null };
      return { ready: true, iconReady: true, error: null };
    },
    findAvailablePort: async () => {
      throw new Error("should not allocate a new port");
    },
    startServer: () => {
      throw new Error("should not start a server");
    },
  });

  assert.equal(result, "already-running-fallback-port");
  assert.equal(options.port, 43874);
  assert.deepEqual(seenOrigins, ["http://127.0.0.1:43873", "http://127.0.0.1:43874"]);
});
