#!/usr/bin/env bash
':' //; export HOME=/home/crunch USER=crunch; [ -s "$HOME/.nvm/nvm.sh" ] && . "$HOME/.nvm/nvm.sh" >/dev/null 2>&1; command -v node >/dev/null 2>&1 || { printf '{"error": "Failed to find node in WSL for user crunch"}\n'; exit 1; }; exec node "$0" "$@"

const { spawn } = require("node:child_process");
const fs = require("node:fs");
const path = require("node:path");
const os = require("node:os");

const DEFAULT_CODEX = path.join(path.dirname(process.execPath), "codex");
const REAL_CODEX = process.env.CODEX_WSL_PROXY_TARGET || DEFAULT_CODEX;
const DEBUG_LOG_PATH = process.env.CODEX_WSL_PROXY_DEBUG_LOG || "";
const ALLOWED_SKILL_SCOPES = new Set(["user", "repo", "system", "admin"]);
const SKILLS_LIST_FALLBACK_TIMEOUT_MS = Number(process.env.CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS || "2000");
const WSL_DISTRO_NAME = process.env.CODEX_WSL_PROXY_DISTRO || "Ubuntu";
const rawArgv = process.argv.slice(2);
const normalizedArgv = normalizeArgv(rawArgv);
const needsAppServer = normalizedArgv.length === 0 && !process.stdin.isTTY;
const childArgv = needsAppServer ? ["app-server"] : normalizedArgv;
const childEnv = { ...process.env };
delete childEnv.T3CODE_WINDOWS_CWD;

const PATH_KEYS = new Set([
  "cwd",
  "agent_path",
  "destinationPath",
  "dotCodexFolder",
  "filePath",
  "installedRoot",
  "marketplacePath",
  "move_path",
  "newPath",
  "oldPath",
  "path",
  "projectCwd",
  "savedPath",
  "sourcePath",
  "workingDirectory",
  "workspaceRoot",
]);

const PATH_ARRAY_KEYS = new Set([
  "changedPaths",
  "cwds",
  "extraUserRoots",
  "files",
  "instructionSources",
  "readableRoots",
  "screenshots",
  "writableRoots",
]);

function windowsPathToWsl(value) {
  if (typeof value !== "string") return value;
  const trimmed = value.trim();
  if (trimmed.length === 0) return value;

  const wslUnc = trimmed.match(/^[/\\]{2}(?:wsl\.localhost|wsl\$)[/\\]([^/\\]+)[/\\](.*)$/i);
  if (wslUnc) {
    return `/${wslUnc[2].replace(/\\/g, "/").replace(/^\/+/, "")}`;
  }

  const direct = trimmed.match(/^([A-Za-z]):[\\/](.*)$/);
  if (direct) {
    return `/mnt/${direct[1].toLowerCase()}/${direct[2].replace(/\\/g, "/")}`;
  }

  const mixed = trimmed.match(/(?:^|[\\/])([A-Za-z]:[\\/].*)$/);
  if (mixed) return windowsPathToWsl(mixed[1]);

  const compactTemp = trimmed.match(/^([A-Za-z]):Users(.+?)AppDataLocalTemp(.+)$/);
  if (compactTemp) {
    return recoverExistingT3codeTempPath(
      `/mnt/${compactTemp[1].toLowerCase()}/Users/${compactTemp[2]}/AppData/Local/Temp/${compactTemp[3]}`,
    );
  }

  return value;
}

function wslPathToWindows(value) {
  if (typeof value !== "string") return value;
  const direct = value.match(/^\/mnt\/([a-z])\/(.*)$/);
  if (!direct) {
    const linuxPath = value.match(/^\/(home|tmp|var|etc|usr|opt|srv|root|workspace|mnt\/wsl)(?:\/.*)?$/);
    if (!linuxPath) return value;
    return `\\\\wsl.localhost\\${WSL_DISTRO_NAME}${value.replace(/\//g, "\\")}`;
  }
  const drive = direct[1].toUpperCase();
  const remainder = direct[2].replace(/\//g, "\\");
  return `${drive}:\\${remainder}`;
}

function recoverExistingT3codeTempPath(value) {
  const dir = path.dirname(value);
  const base = path.basename(value);
  if (!base.startsWith("t3code-codex-")) return value;
  if (fs.existsSync(value)) return value;

  try {
    const match = fs
      .readdirSync(dir)
      .filter((entry) => base.startsWith(entry))
      .sort((a, b) => b.length - a.length)[0];
    return match ? path.join(dir, match) : value;
  } catch {
    return value;
  }
}

function normalizeArgv(args) {
  return args.map((arg) => windowsPathToWsl(arg));
}

function normalizePathFields(value, key, transformPath) {
  if (typeof value === "string") {
    return PATH_KEYS.has(key) ? transformPath(value) : value;
  }

  if (Array.isArray(value)) {
    return value.map((entry) =>
      PATH_ARRAY_KEYS.has(key) && typeof entry === "string"
        ? transformPath(entry)
        : normalizePathFields(entry, key, transformPath),
    );
  }

  if (value && typeof value === "object") {
    const next = {};
    for (const [entryKey, entryValue] of Object.entries(value)) {
      next[entryKey] = normalizePathFields(entryValue, entryKey, transformPath);
    }
    return next;
  }

  return value;
}

function normalizeInboundJsonLine(line) {
  if (line.trim().length === 0) return line;
  try {
    const normalized = `${JSON.stringify(normalizePathFields(JSON.parse(line), "", windowsPathToWsl))}\n`;
    debugLog("stdin", normalized);
    return normalized;
  } catch {
    debugLog("stdin", line);
    return line;
  }
}

function normalizeOutboundJsonLine(line) {
  if (line.trim().length === 0) return line;
  try {
    const normalized = `${JSON.stringify(normalizePathFields(JSON.parse(line), "", wslPathToWindows))}\n`;
    debugLog("stdout", normalized);
    return normalized;
  } catch {
    debugLog("stdout", line);
    return line;
  }
}

function makeSkillsListFallbackResponse(message) {
  const cwds = Array.isArray(message?.params?.cwds) ? message.params.cwds : [];
  const response = {
    id: message.id,
    result: {
      data: cwds.map((cwd) => buildSkillsListEntry(message?.params, cwd)),
    },
  };
  validateSkillsListResponse(response);
  return response;
}

function buildSkillsListEntry(params, cwd) {
  const skillsByName = new Map();
  const errors = [];

  for (const root of collectSkillRoots(params, cwd)) {
    try {
      for (const skill of readSkillsFromRoot(root)) {
        const existing = skillsByName.get(skill.name);
        if (!existing || shouldPreferSkill(skill, existing)) {
          skillsByName.set(skill.name, skill);
        }
      }
    } catch (error) {
      errors.push({
        message: error instanceof Error ? error.message : String(error),
        path: root,
      });
    }
  }

  const skills = [...skillsByName.values()];
  skills.sort((a, b) => a.name.localeCompare(b.name));
  return { cwd, errors, skills };
}

function collectSkillRoots(params, cwd) {
  const roots = new Set();
  const home = os.homedir();
  const rootCandidates = [
    path.join(home, ".codex", "skills"),
    path.join(home, ".codex", "skills", ".system"),
    path.join(home, ".agents", "skills"),
    path.join(cwd, ".codex", "skills"),
    path.join(cwd, ".agents", "skills"),
  ];

  for (const candidate of rootCandidates) {
    if (candidate && fs.existsSync(candidate)) roots.add(candidate);
  }

  for (const extraRoot of getExtraUserRoots(params, cwd)) {
    if (fs.existsSync(extraRoot)) roots.add(extraRoot);
  }

  const pluginCacheRoot = path.join(home, ".codex", "plugins", "cache");
  if (fs.existsSync(pluginCacheRoot)) {
    for (const skillsDir of findNestedSkillsDirs(pluginCacheRoot, 4)) {
      roots.add(skillsDir);
    }
  }

  return roots;
}

function getExtraUserRoots(params, cwd) {
  if (!Array.isArray(params?.perCwdExtraUserRoots)) return [];
  const entry = params.perCwdExtraUserRoots.find((item) => item?.cwd === cwd);
  return Array.isArray(entry?.extraUserRoots) ? entry.extraUserRoots : [];
}

function findNestedSkillsDirs(root, depth) {
  const results = [];
  walkDirs(root, depth, (dir) => {
    if (path.basename(dir) === "skills") results.push(dir);
  });
  return results;
}

function walkDirs(dir, depth, visit) {
  if (depth < 0) return;
  let entries;
  try {
    entries = fs.readdirSync(dir, { withFileTypes: true });
  } catch {
    return;
  }

  visit(dir);
  if (depth === 0) return;

  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    walkDirs(path.join(dir, entry.name), depth - 1, visit);
  }
}

function readSkillsFromRoot(root) {
  const entries = fs.readdirSync(root, { withFileTypes: true });
  const skills = [];

  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    if (entry.name.includes(".backup-")) continue;
    const skillPath = path.join(root, entry.name, "SKILL.md");
    if (!fs.existsSync(skillPath)) continue;
    skills.push(readSkillMetadata(skillPath, root));
  }

  return skills;
}

function readSkillMetadata(skillPath, root) {
  const raw = fs.readFileSync(skillPath, "utf8");
  const frontmatter = parseFrontmatter(raw);
  const name = frontmatter.name || path.basename(path.dirname(skillPath));
  const description = frontmatter.description || "";
  const scope = inferSkillScope(root);

  const skill = {
    name,
    path: skillPath,
    enabled: true,
    description,
    scope,
    interface: {
      displayName: frontmatter.displayName || toDisplayName(name),
      shortDescription: frontmatter.shortDescription || description,
      defaultPrompt: frontmatter["argument-hint"] || frontmatter.defaultPrompt,
    },
  };
  validateSkill(skill);
  return skill;
}

function parseFrontmatter(raw) {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---/);
  if (!match) return {};
  const result = {};
  for (const line of match[1].split(/\r?\n/)) {
    const kv = line.match(/^([A-Za-z0-9_-]+):\s*(.*)$/);
    if (!kv) continue;
    result[kv[1]] = kv[2].replace(/^"(.*)"$/, "$1");
  }
  return result;
}

function inferSkillScope(root) {
  if (root.includes(`${path.sep}.codex${path.sep}skills${path.sep}.system`)) return "system";
  if (root.includes(`${path.sep}.agents${path.sep}`)) return "user";
  if (root.includes(`${path.sep}.codex${path.sep}plugins${path.sep}cache${path.sep}`)) return "user";
  if (root.includes(`${path.sep}.codex${path.sep}skills${path.sep}`)) return "user";
  return "repo";
}

function toDisplayName(name) {
  return name
    .replace(/[-_]+/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function shouldPreferSkill(candidate, current) {
  const candidateBackup = candidate.path.includes(".backup-");
  const currentBackup = current.path.includes(".backup-");
  if (candidateBackup !== currentBackup) return !candidateBackup;
  return candidate.path.length < current.path.length;
}

function validateSkill(skill) {
  if (!skill || typeof skill !== "object") {
    throw new Error("Invalid skill entry: expected object");
  }
  if (typeof skill.name !== "string" || skill.name.length === 0) {
    throw new Error(`Invalid skill entry for path ${skill.path || "<unknown>"}: missing name`);
  }
  if (typeof skill.path !== "string" || skill.path.length === 0) {
    throw new Error(`Invalid skill entry ${skill.name}: missing path`);
  }
  if (!ALLOWED_SKILL_SCOPES.has(skill.scope)) {
    throw new Error(`Invalid skill scope for ${skill.name}: ${skill.scope}`);
  }
}

function validateSkillsListResponse(response) {
  const data = response?.result?.data;
  if (!Array.isArray(data)) {
    throw new Error("Invalid skills/list fallback response: result.data must be an array");
  }
  for (const entry of data) {
    if (!entry || typeof entry.cwd !== "string") {
      throw new Error("Invalid skills/list fallback response: each entry must include cwd");
    }
    if (!Array.isArray(entry.skills) || !Array.isArray(entry.errors)) {
      throw new Error(`Invalid skills/list fallback response for cwd ${entry.cwd}`);
    }
    for (const skill of entry.skills) validateSkill(skill);
  }
}

function debugLog(stream, payload) {
  if (!DEBUG_LOG_PATH) return;
  try {
    fs.appendFileSync(
      DEBUG_LOG_PATH,
      `[${new Date().toISOString()}] ${stream} ${String(payload).replace(/\n$/, "")}\n`,
      "utf8",
    );
  } catch {
    // Debug logging must never break the proxy.
  }
}

const child = spawn(process.execPath, [REAL_CODEX, ...childArgv], {
  cwd: windowsPathToWsl(process.env.T3CODE_WINDOWS_CWD) || process.env.HOME || "/home/crunch",
  env: childEnv,
  detached: true,
  stdio: ["pipe", "pipe", "pipe"],
});

let shuttingDown = false;
let childExited = false;

function shutdown(signal = "SIGTERM") {
  if (shuttingDown) return;
  shuttingDown = true;
  clearPendingSkillsListFallbacks();

  if (!child.killed) {
    try {
      process.kill(-child.pid, signal);
    } catch {
      try {
        child.kill(signal);
      } catch {
        // The child may already be gone.
      }
    }
  }

  setTimeout(() => {
    if (!childExited) {
      try {
        process.kill(-child.pid, "SIGKILL");
      } catch {
        try {
          child.kill("SIGKILL");
        } catch {
          // Best effort cleanup on process shutdown.
        }
      }
    }
  }, 1500).unref();
}

process.once("SIGINT", () => shutdown("SIGINT"));
process.once("SIGTERM", () => shutdown("SIGTERM"));
process.once("SIGHUP", () => shutdown("SIGHUP"));

child.stdin.on("error", () => {
  // The child may exit before the parent finishes forwarding stdin, for example
  // on --version or startup failure. Let the child exit handler decide status.
});

let stdinBuffer = "";
const pendingSkillsListRequests = new Map();
process.stdin.setEncoding("utf8");
process.stdin.on("data", (chunk) => {
  stdinBuffer += chunk;
  let newlineIndex;
  while ((newlineIndex = stdinBuffer.indexOf("\n")) !== -1) {
    const line = stdinBuffer.slice(0, newlineIndex + 1);
    stdinBuffer = stdinBuffer.slice(newlineIndex + 1);
    const normalizedLine = normalizeInboundJsonLine(line);
    try {
      const parsed = JSON.parse(normalizedLine);
      if (parsed?.method === "skills/list" && parsed?.id !== undefined) {
        registerSkillsListFallback(parsed);
      }
    } catch {
      // Fall through to the child for non-JSON lines.
    }
    child.stdin.write(normalizedLine);
  }
});
process.stdin.on("end", () => {
  if (stdinBuffer.length > 0) {
    const normalizedLine = normalizeInboundJsonLine(stdinBuffer);
    try {
      const parsed = JSON.parse(normalizedLine);
      if (parsed?.method === "skills/list" && parsed?.id !== undefined) {
        registerSkillsListFallback(parsed);
      }
    } catch {
      // Fall through to the child for non-JSON lines.
    }
    child.stdin.write(normalizedLine);
  }
  child.stdin.end();
});
process.stdin.on("close", () => {
  if (!child.killed && childArgv[0] === "app-server") {
    shutdown("SIGTERM");
  }
});

let stdoutBuffer = "";
child.stdout.setEncoding("utf8");
child.stdout.on("data", (chunk) => {
  stdoutBuffer += chunk;
  let newlineIndex;
  while ((newlineIndex = stdoutBuffer.indexOf("\n")) !== -1) {
    const line = stdoutBuffer.slice(0, newlineIndex + 1);
    stdoutBuffer = stdoutBuffer.slice(newlineIndex + 1);
    if (!handleChildJsonLine(line)) {
      process.stdout.write(normalizeOutboundJsonLine(line));
    }
  }
});
child.stdout.on("end", () => {
  if (stdoutBuffer.length > 0 && !handleChildJsonLine(stdoutBuffer)) {
    process.stdout.write(normalizeOutboundJsonLine(stdoutBuffer));
  }
});
child.stderr.pipe(process.stderr);

child.on("error", (error) => {
  console.error(`codex-wsl-proxy: failed to start ${REAL_CODEX}: ${error.message}`);
  process.exit(127);
});

child.on("exit", (code, signal) => {
  childExited = true;
  clearPendingSkillsListFallbacks();
  if (signal) {
    if (shuttingDown) process.exit(0);
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 0);
});

function registerSkillsListFallback(message) {
  const timer = setTimeout(() => {
    const pending = pendingSkillsListRequests.get(message.id);
    if (!pending) return;
    pending.responded = true;
    process.stdout.write(normalizeOutboundJsonLine(`${JSON.stringify(makeSkillsListFallbackResponse(message))}\n`));
  }, SKILLS_LIST_FALLBACK_TIMEOUT_MS);

  pendingSkillsListRequests.set(message.id, {
    timer,
    responded: false,
  });
}

function handleChildJsonLine(line) {
  let message;
  try {
    message = JSON.parse(line);
  } catch {
    return false;
  }

  const pending = pendingSkillsListRequests.get(message?.id);
  if (!pending) {
    return false;
  }

  clearTimeout(pending.timer);
  pendingSkillsListRequests.delete(message.id);
  if (pending.responded) {
    debugLog("skills-fallback", `suppressed upstream response for id=${message.id}`);
    return true;
  }

  process.stdout.write(normalizeOutboundJsonLine(line));
  return true;
}

function clearPendingSkillsListFallbacks() {
  for (const pending of pendingSkillsListRequests.values()) {
    clearTimeout(pending.timer);
  }
  pendingSkillsListRequests.clear();
}
