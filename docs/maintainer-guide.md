# Maintainer Guide

This guide captures the source-first workflow and local checks for maintaining
this repository.

## Source-First Workflow

Keep edits in this repository first. Copy files into active workstation install
locations only when explicitly doing install, repair, or republish work.

For behavior, workflow, script, or config changes:

1. Change the source file in this repository.
2. Run the smallest relevant check.
3. Update the focused docs page that owns the behavior.
4. Repair or republish installed files only when that is part of the task.

## Common Checks

### Documentation-only changes

Run a local Markdown link check over the changed docs and check for whitespace
errors before handing off documentation-only work:

```bash
node - README.md RESOURCES.md docs/README.md docs/maintainer-guide.md <<'NODE'
const fs = require('fs');
const path = require('path');
const files = process.argv.slice(2);
let failed = false;
for (const file of files) {
  const text = fs.readFileSync(file, 'utf8');
  const re = /\[[^\]]+\]\(([^)]+)\)/g;
  let match;
  while ((match = re.exec(text))) {
    const href = match[1];
    if (/^[a-z]+:/i.test(href) || href.startsWith('#')) continue;
    const target = href.split('#')[0];
    const resolved = path.resolve(path.dirname(file), target);
    if (!fs.existsSync(resolved)) {
      console.error(`${file}: missing link target ${href}`);
      failed = true;
    }
  }
}
process.exit(failed ? 1 : 0);
NODE

git diff --check
```

### WSL proxy modules

```bash
node --test bin/codex-wsl-proxy-runtime.test.js
```

This covers protocol path translation, outbound Linux path conversion,
proxy-only environment cleanup, `skills/list` fallback schema, Windows path-list
handling, turn-id parsing, and timeout parsing.

### Workstation hooks

Run the focused source test before publishing changes to the active global hook:

```bash
python3 -m unittest hooks/test_wsl_command_guardrails.py
```

The active global copy lives under `$CODEX_HOME/hooks/`. When repairing the
installed guardrail, copy `hooks/wsl_command_guardrails.py` to
`$CODEX_HOME/hooks/wsl_command_guardrails.py` and keep `$CODEX_HOME/hooks.json`
matching Codex shell tool names used by the current runtime. The source-side
reference snippet is `hooks/global-pretooluse-hooks.example.json`.

### Skills updater help path

From Windows PowerShell:

```powershell
bin\skills-updates.cmd --help
bin\sk-up.cmd -h
```

On this workstation, run the skills-updater PowerShell checks from Windows
PowerShell, preferably through an existing Tabby PowerShell session when an
agent needs to drive the check.

Windows PowerShell can still target slash paths when needed:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File bin/skills-updates.ps1 --cmd-name skills-updates --help
powershell.exe -NoProfile -ExecutionPolicy Bypass -File bin/skills-updates.ps1 --cmd-name sk-up -h
```

For the source-install regression check:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests\skills-updates-install.ps1
```

The slash path form is also accepted from Windows PowerShell:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests/skills-updates-install.ps1
```

## Installed Shim Freshness

The intended Windows entry point is a native Windows symlink to
`bin/codex-wsl.cmd`. Repair that symlink with Windows-native tooling.

The WSL-side runtime files live under `~/.local/bin` when installed:

```bash
ls -l ~/.local/bin/codex-wsl-proxy*.js ~/.local/bin/codex-wsl-path-translation.js ~/.local/bin/codex-wsl-skills-fallback.js
```

If those installed files differ from `bin/`, repair by copying from this repo.
Do not mutate installed runtime files during ordinary source edits unless the
task explicitly asks for install or repair work.

## Current Design Rationale

- The Windows entry point stays small because Windows-specific work is limited
  to finding `wsl.exe`, selecting a distro/user, passing environment through
  `WSLENV`, and preserving the Windows current directory.
- The WSL runner owns Codex process startup because it can reliably normalize
  `HOME`, `USER`, `CODEX_HOME`, and WSL `PATH`.
- The Node proxy is split into runtime, path translation, and skills fallback
  modules so path-policy changes and app-server lifecycle changes are testable
  without exercising the Windows batch file.
- The skills updater is PowerShell-first because it manages Windows global skill
  installs, Windows `%USERPROFILE%` paths, console codepages, editor launching,
  and named mutexes.
- LazyCodex is kept as a Codex plugin in WSL. The WSL runner allows LazyCodex
  auto-update by default while disabling telemetry and config-migration startup
  paths for app-server sessions.

## Safety Rules

- Preserve user changes and unrelated local state.
- Use Windows-native tooling to repair symlinks on Windows drives.
- Keep temporary agent notes out of VCS with `.git/info/exclude`, not
  `.gitignore`.
- Never record secrets, tokens, recovery codes, private personal data, or raw
  session exports in docs, examples, commits, or skill text.
- For behavior, workflow, script, or config changes, update the relevant human
  docs in the same turn.
