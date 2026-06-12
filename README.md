# agents-toolkit

Local Codex workstation tooling for this machine.

This repository is the source tree for:

- the Windows-to-WSL Codex shim used by T3code-style app-server sessions;
- the Node proxy that translates Windows and WSL paths for Codex protocol
  traffic;
- the Windows skills updater wrappers for globally installed Codex skills;
- local Codex skills maintained on this workstation;
- operator docs for companion tooling, shell setup, and local archives.

Keep changes here first. Copy files into the active workstation install only
when the user asks to install, repair, or republish the local setup.

## Current Surface

| Area | Source of truth | Notes |
| --- | --- | --- |
| WSL Codex shim | [docs/wsl-shim.md](docs/wsl-shim.md) | Windows `.cmd`, WSL runner, Node proxy runtime, path translation, skills fallback, and T3code app-server behavior. |
| Skills updater | [docs/skills-updater.md](docs/skills-updater.md) | `sk-up` and `skills-updates` flags, state paths, locking, install/uninstall behavior, and verification. |
| Companion tools | [docs/codex-cli-tooling.md](docs/codex-cli-tooling.md) | Evo, RTK, ICM, CodSpeed, OMO/LazyCodex, OpenAI docs MCP, Discrawl, and adjacent CLIs. |
| Shell setup | [docs/shell-setup.md](docs/shell-setup.md) | Starship, fzf, zoxide, Atuin, PSReadLine, ble.sh, Tabby, Windows, WSL, and Arch parity. |
| Discrawl wiretap | [docs/discrawl-wiretap.md](docs/discrawl-wiretap.md) | Local Vesktop-cache archive, user systemd timer, limits, and verification. |
| Local skills | [skills/](skills/) | Repo-owned skill sources. In this repo, "work on skills" means this directory, not installed global skill copies. |
| Learning artifacts | [lessons/](lessons/) and [reference/](reference/) | Static HTML terminal-addon lesson and cheatsheet. |
| Plugin backups | [backups/](backups/) | Restore payloads from plugin backup operations. They are historical artifacts, not canonical docs to refresh. |

Root docs:

- [MISSION.md](MISSION.md): durable project mission and non-goals.
- [CONTEXT.md](CONTEXT.md): local vocabulary and architecture terms.
- [RESOURCES.md](RESOURCES.md): authoritative references used by this repo.
- [AGENTS.md](AGENTS.md): local agent instructions; it is workstation-local and
  excluded from normal VCS use by policy.
- [RTK.md](RTK.md): local shell-command wrapper rule for agents.

## Repository Layout

```text
bin/
  codex-wsl.cmd
  codex-wsl-proxy-runner.sh
  codex-wsl-proxy.js
  codex-wsl-proxy-runtime.js
  codex-wsl-path-translation.js
  codex-wsl-skills-fallback.js
  skills-updates.ps1
  skills-updates.cmd
  sk-up.cmd
docs/
  codex-cli-tooling.md
  discrawl-wiretap.md
  shell-setup.md
  skills-updater.md
  wsl-shim.md
skills/
  codex-goal-control/
  confidence-loop/
  discrawl-local/
  evo-end-to-end/
  icm/
  icm-recall/
  icm-store/
  tuck/
  yeet/
tests/
  skills-updates-install.ps1
```

## Common Tasks

### Verify the WSL proxy modules

```bash
node --test bin/codex-wsl-proxy-runtime.test.js
```

This covers protocol path translation, outbound Linux path conversion,
proxy-only environment cleanup, `skills/list` fallback schema, Windows path-list
handling, turn-id parsing, and timeout parsing.

### Probe the skills updater help path

From Windows PowerShell:

```powershell
bin\skills-updates.cmd --help
bin\sk-up.cmd -h
```

From WSL:

```bash
powershell.exe -NoProfile -ExecutionPolicy Bypass -File bin/skills-updates.ps1 --cmd-name skills-updates --help
powershell.exe -NoProfile -ExecutionPolicy Bypass -File bin/skills-updates.ps1 --cmd-name sk-up -h
```

For the source-install regression check:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests\skills-updates-install.ps1
```

From WSL, use the slash path form:

```bash
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests/skills-updates-install.ps1
```

### Check installed shim freshness

The intended Windows entry point is a native Windows symlink:

```text
C:\Users\nguco\bin\codex-wsl.cmd -> E:\dev\agents-toolkit\bin\codex-wsl.cmd
```

The WSL-side runtime files live under `~/.local/bin` when installed:

```bash
ls -l ~/.local/bin/codex-wsl-proxy*.js ~/.local/bin/codex-wsl-path-translation.js ~/.local/bin/codex-wsl-skills-fallback.js
```

If those installed files differ from `bin/`, repair by copying from this repo.
Do not mutate installed runtime files during ordinary source edits unless the
user explicitly asks for install or repair work.

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
  installs, Windows `%USERPROFILE%` paths, console codepages, Zed launching, and
  named mutexes.
- LazyCodex / OMO is kept as a Codex plugin. The WSL runner disables its
  telemetry, auto-update, and config-migration startup paths for T3code
  app-server sessions while leaving the plugin itself available.

## Safety Rules

- Preserve user changes and unrelated local state.
- Use Windows-native tooling to repair symlinks on Windows drives.
- Keep temporary agent notes out of VCS with `.git/info/exclude`, not
  `.gitignore`.
- Never record secrets, tokens, recovery codes, private personal data, or raw
  session exports in docs, examples, commits, or skill text.
- For behavior, workflow, script, or config changes, update the relevant human
  docs in the same turn.
