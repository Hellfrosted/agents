const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

const ALLOWED_SKILL_SCOPES = new Set(["user", "repo", "system", "admin"]);

function createSkillsFallback({ homeDir = os.homedir(), windowsPathToWsl }) {
  return {
    makeResponse(message) {
      const cwds = Array.isArray(message?.params?.cwds) ? message.params.cwds : [];
      const response = {
        id: message.id,
        result: {
          data: cwds.map((cwd) => buildSkillsListEntry(message?.params, cwd, homeDir, windowsPathToWsl)),
        },
      };
      validateSkillsListResponse(response);
      return response;
    },
  };
}

function buildSkillsListEntry(params, cwd, homeDir, windowsPathToWsl) {
  const skillsByName = new Map();
  const errors = [];

  for (const root of collectSkillRoots(params, cwd, homeDir, windowsPathToWsl)) {
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

function collectSkillRoots(params, cwd, home, windowsPathToWsl) {
  const roots = new Set();
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

  for (const extraRoot of getEnvSkillRoots(windowsPathToWsl)) {
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

function getEnvSkillRoots(windowsPathToWsl) {
  return [
    ...splitPathListEnv(process.env.CODEX_SKILLS_DIRS),
    ...splitPathListEnv(process.env.CODEX_SKILL_ROOTS),
  ]
    .map((root) => windowsPathToWsl(root))
    .filter((root) => typeof root === "string" && root.trim().length > 0);
}

function splitPathListEnv(value) {
  if (typeof value !== "string" || value.trim().length === 0) return [];
  const trimmed = value.trim();
  if (trimmed.includes(";")) return trimmed.split(";").map((entry) => entry.trim()).filter(Boolean);
  if (/^[A-Za-z]:[\\/]/.test(trimmed)) return [trimmed];
  return trimmed.split(path.delimiter).map((entry) => entry.trim()).filter(Boolean);
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
    if (entry.isDirectory()) walkDirs(path.join(dir, entry.name), depth - 1, visit);
  }
}

function readSkillsFromRoot(root) {
  const entries = fs.readdirSync(root, { withFileTypes: true });
  const skills = [];

  for (const entry of entries) {
    if (!entry.isDirectory() || entry.name.includes(".backup-")) continue;
    const skillPath = path.join(root, entry.name, "SKILL.md");
    if (fs.existsSync(skillPath)) skills.push(readSkillMetadata(skillPath, root));
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
    if (kv) result[kv[1]] = kv[2].replace(/^"(.*)"$/, "$1");
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
  return name.replace(/[-_]+/g, " ").replace(/\b\w/g, (char) => char.toUpperCase());
}

function shouldPreferSkill(candidate, current) {
  const candidateBackup = candidate.path.includes(".backup-");
  const currentBackup = current.path.includes(".backup-");
  if (candidateBackup !== currentBackup) return !candidateBackup;
  return candidate.path.length < current.path.length;
}

function validateSkill(skill) {
  if (!skill || typeof skill !== "object") throw new Error("Invalid skill entry: expected object");
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

module.exports = {
  createSkillsFallback,
  parseFrontmatter,
  splitPathListEnv,
  validateSkillsListResponse,
};
