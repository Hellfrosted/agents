const fs = require("node:fs");
const path = require("node:path");

const PROTOCOL_PATH_POLICY = {
  scalarFields: new Set([
    "cwd",
    "agent_path",
    "codexHome",
    "composerIcon",
    "destinationPath",
    "dotCodexFolder",
    "filePath",
    "grantRoot",
    "iconLarge",
    "iconSmall",
    "installedRoot",
    "localPluginPath",
    "logo",
    "managedDir",
    "marketplacePath",
    "move_path",
    "newPath",
    "oldPath",
    "path",
    "pluginPath",
    "projectCwd",
    "root",
    "savedPath",
    "source",
    "sourcePath",
    "windowsManagedDir",
    "workingDirectory",
    "working_directory",
    "workspaceRoot",
  ]),
  arrayFields: new Set([
    "changedPaths",
    "cwds",
    "extraLogFiles",
    "extraUserRoots",
    "files",
    "instructionSources",
    "preexisting_untracked_dirs",
    "preexisting_untracked_files",
    "read",
    "readableRoots",
    "readable_roots",
    "roots",
    "samplePaths",
    "screenshots",
    "sparsePaths",
    "upgradedRoots",
    "write",
    "writableRoots",
    "writable_roots",
  ]),
  keyedPathMaps: new Set(["fileChanges"]),

  shouldTranslateScalarField(key) {
    return this.scalarFields.has(key);
  },

  shouldTranslateArrayEntry(parentKey, entry) {
    return this.arrayFields.has(parentKey) && typeof entry === "string";
  },

  shouldTranslateMapKey(parentKey) {
    return this.keyedPathMaps.has(parentKey);
  },
};

function createPathTranslator({ distroName = "", debugLog = () => {} } = {}) {
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
      if (!linuxPath || !distroName) return value;
      return `\\\\wsl.localhost\\${distroName}${value.replace(/\//g, "\\")}`;
    }

    const drive = direct[1].toUpperCase();
    const remainder = direct[2].replace(/\//g, "\\");
    return `${drive}:\\${remainder}`;
  }

  function normalizeInboundJsonLine(line) {
    return normalizeJsonLine(line, "stdin", windowsPathToWsl);
  }

  function normalizeOutboundJsonLine(line) {
    return normalizeJsonLine(line, "stdout", wslPathToWindows);
  }

  function normalizeJsonLine(line, stream, transformPath) {
    if (line.trim().length === 0) return line;
    try {
      const normalized = `${JSON.stringify(normalizePathFields(JSON.parse(line), "", transformPath))}\n`;
      debugLog(stream, normalized);
      return normalized;
    } catch {
      debugLog(stream, line);
      return line;
    }
  }

  return {
    normalizeInboundJsonLine,
    normalizeOutboundJsonLine,
    windowsPathToWsl,
    wslPathToWindows,
  };
}

function normalizePathFields(value, key, transformPath) {
  if (typeof value === "string") {
    return PROTOCOL_PATH_POLICY.shouldTranslateScalarField(key) ? transformPath(value) : value;
  }

  if (Array.isArray(value)) {
    return value.map((entry) =>
      PROTOCOL_PATH_POLICY.shouldTranslateArrayEntry(key, entry)
        ? transformPath(entry)
        : normalizePathFields(entry, key, transformPath),
    );
  }

  if (value && typeof value === "object") {
    const next = {};
    for (const [entryKey, entryValue] of Object.entries(value)) {
      const normalizedKey = PROTOCOL_PATH_POLICY.shouldTranslateMapKey(key) ? transformPath(entryKey) : entryKey;
      next[normalizedKey] = normalizePathFields(entryValue, entryKey, transformPath);
    }
    return next;
  }

  return value;
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

module.exports = {
  createPathTranslator,
  normalizePathFields,
};
