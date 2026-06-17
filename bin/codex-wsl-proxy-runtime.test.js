const assert = require("node:assert/strict");
const test = require("node:test");

const { createPathTranslator } = require("./codex-wsl-path-translation");
const { createSkillsFallback, splitPathListEnv, validateSkillsListResponse } = require("./codex-wsl-skills-fallback");
const {
  AppServerSession,
  SkillsListFallbacks,
  buildChildEnv,
  parseNonNegativeInteger,
  readMessageTurnId,
} = require("./codex-wsl-proxy-runtime");

test("Given Windows paths When applying protocol path policy Then scalar fields, arrays, and keyed maps convert", () => {
  const translator = createPathTranslator({ distroName: "Ubuntu" });

  const line = translator.normalizeInboundJsonLine(
    `${JSON.stringify({
      method: "x",
      params: {
        cwd: "C:\\Users\\me\\repo",
        label: "C:\\not\\a\\path\\field",
        files: ["D:\\a\\b.txt"],
        fileChanges: {
          "E:\\old.txt": { path: "E:\\new.txt" },
        },
        events: [{ path: "F:\\event.txt", label: "F:\\not\\policy\\field" }],
      },
    })}\n`,
  );

  assert.deepEqual(JSON.parse(line).params, {
    cwd: "/mnt/c/Users/me/repo",
    label: "C:\\not\\a\\path\\field",
    files: ["/mnt/d/a/b.txt"],
    fileChanges: {
      "/mnt/e/old.txt": { path: "/mnt/e/new.txt" },
    },
    events: [{ path: "/mnt/f/event.txt", label: "F:\\not\\policy\\field" }],
  });
});

test("Given WSL paths When normalizing outbound protocol JSON Then Windows paths convert", () => {
  const translator = createPathTranslator({ distroName: "Ubuntu" });

  const line = translator.normalizeOutboundJsonLine(
    `${JSON.stringify({
      result: {
        path: "/mnt/c/Users/me/repo",
        cwd: "/home/me/project",
        message: "/home/me/project",
      },
    })}\n`,
  );

  assert.deepEqual(JSON.parse(line).result, {
    path: "C:\\Users\\me\\repo",
    cwd: "\\\\wsl.localhost\\Ubuntu\\home\\me\\project",
    message: "/home/me/project",
  });
});

test("Given child env When building env Then proxy-only vars drop and CODEX_HOME converts", () => {
  const originalEnv = process.env;
  process.env = {
    CODEX_HOME: "C:\\Users\\me\\.codex",
    CODEX_WSL_PROXY_IDLE_TIMEOUT_MS: "10",
    T3CODE_WINDOWS_CWD: "C:\\repo",
  };

  try {
    const translator = createPathTranslator();
    const childEnv = buildChildEnv(translator.windowsPathToWsl);

    assert.equal(childEnv.CODEX_HOME, "/mnt/c/Users/me/.codex");
    assert.equal(childEnv.CODEX_WSL_PROXY_IDLE_TIMEOUT_MS, undefined);
    assert.equal(childEnv.T3CODE_WINDOWS_CWD, undefined);
  } finally {
    process.env = originalEnv;
  }
});

test("Given skills list request When fallback builds response Then schema matches protocol", () => {
  const fallback = createSkillsFallback({ homeDir: "/tmp/no-such-home", windowsPathToWsl: (value) => value });

  const response = fallback.makeResponse({
    id: 42,
    params: {
      cwds: ["/tmp/no-such-cwd"],
    },
  });

  validateSkillsListResponse(response);
  assert.equal(response.id, 42);
  assert.deepEqual(response.result.data, [
    {
      cwd: "/tmp/no-such-cwd",
      errors: [],
      skills: [],
    },
  ]);
});

test("Given semicolon path list When splitting env Then Windows drive colons stay intact", () => {
  assert.deepEqual(splitPathListEnv("C:\\skills;D:\\more"), ["C:\\skills", "D:\\more"]);
  assert.deepEqual(splitPathListEnv("C:\\only"), ["C:\\only"]);
});

test("Given protocol messages When reading turn id Then both known shapes work", () => {
  assert.equal(readMessageTurnId({ params: { turnId: "a" } }), "a");
  assert.equal(readMessageTurnId({ params: { turn: { id: "b" } } }), "b");
  assert.equal(readMessageTurnId({ params: {} }), undefined);
});

test("Given app-server session When idle with active turn Then reaper waits for completion", () => {
  let now = 0;
  const signals = [];
  const session = new AppServerSession({
    debugLog: () => {},
    enabled: true,
    idleTimeoutMs: 10,
    now: () => now,
    requestShutdown: (signal) => signals.push(signal),
  });

  session.observeProtocolMessage({ method: "turn/started", params: { turnId: "turn-1" } });
  now = 20;
  session.reapIdleChild();

  assert.deepEqual(signals, []);

  session.observeProtocolMessage({ method: "turn/completed", params: { turnId: "turn-1" } });
  now = 31;
  session.reapIdleChild();

  assert.deepEqual(signals, ["SIGTERM"]);
});

test("Given non app-server session When stdin closes Then child keeps running", () => {
  const appServer = new AppServerSession({
    debugLog: () => {},
    enabled: true,
    idleTimeoutMs: 10,
    requestShutdown: () => {},
  });
  const command = new AppServerSession({
    debugLog: () => {},
    enabled: false,
    idleTimeoutMs: 10,
    requestShutdown: () => {},
  });

  assert.equal(appServer.shouldShutdownOnInputClose({ killed: false }), true);
  assert.equal(appServer.shouldShutdownOnInputClose({ killed: true }), false);
  assert.equal(command.shouldShutdownOnInputClose({ killed: false }), false);
});

test("Given skills fallback answered When upstream replies later Then upstream response is suppressed", async () => {
  const lines = [];
  const logs = [];
  const fallback = new SkillsListFallbacks({
    debugLog: (stream, payload) => logs.push([stream, payload]),
    normalizeOutboundJsonLine: (line) => line,
    skillsFallback: {
      makeResponse(message) {
        return { id: message.id, result: { data: [] } };
      },
    },
    timeoutMs: 0,
    writeOutput: (line) => lines.push(line),
  });

  fallback.observeInboundMessage({ id: 7, method: "skills/list" });
  await new Promise((resolve) => setTimeout(resolve, 0));

  assert.equal(fallback.handleUpstreamMessage({ id: 7 }, `${JSON.stringify({ id: 7, result: "upstream" })}\n`), true);
  assert.deepEqual(lines.map((line) => JSON.parse(line)), [{ id: 7, result: { data: [] } }]);
  assert.deepEqual(logs, [["skills-fallback", "suppressed upstream response for id=7"]]);
});

test("Given skills fallback pending When upstream replies first Then upstream response is forwarded", () => {
  const lines = [];
  const fallback = new SkillsListFallbacks({
    debugLog: () => {},
    normalizeOutboundJsonLine: (line) => `normalized:${line}`,
    skillsFallback: {
      makeResponse(message) {
        return { id: message.id, result: { data: [] } };
      },
    },
    timeoutMs: 1000,
    writeOutput: (line) => lines.push(line),
  });
  const upstreamLine = `${JSON.stringify({ id: 8, result: "upstream" })}\n`;

  fallback.observeInboundMessage({ id: 8, method: "skills/list" });

  assert.equal(fallback.handleUpstreamMessage({ id: 8 }, upstreamLine), true);
  assert.deepEqual(lines, [`normalized:${upstreamLine}`]);
  fallback.clear();
});

test("Given integer config When parsing timeout Then invalid values use fallback", () => {
  assert.equal(parseNonNegativeInteger("123", 5), 123);
  assert.equal(parseNonNegativeInteger("-1", 5), 5);
  assert.equal(parseNonNegativeInteger("no", 5), 5);
  assert.equal(parseNonNegativeInteger("", 5), 5);
});
