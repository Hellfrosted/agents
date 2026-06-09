const assert = require("node:assert/strict");
const test = require("node:test");

const { createPathTranslator } = require("./codex-wsl-path-translation");
const { createSkillsFallback, splitPathListEnv, validateSkillsListResponse } = require("./codex-wsl-skills-fallback");
const { buildChildEnv, parseNonNegativeInteger, readMessageTurnId } = require("./codex-wsl-proxy-runtime");

test("Given Windows paths When normalizing protocol JSON Then only path fields convert", () => {
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

test("Given integer config When parsing timeout Then invalid values use fallback", () => {
  assert.equal(parseNonNegativeInteger("123", 5), 123);
  assert.equal(parseNonNegativeInteger("-1", 5), 5);
  assert.equal(parseNonNegativeInteger("no", 5), 5);
  assert.equal(parseNonNegativeInteger("", 5), 5);
});
