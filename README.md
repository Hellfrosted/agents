# agents-toolkit

Codex tooling source, scripts, hooks, and local skills.

This repository is public source for the Codex-adjacent tools owned here. It is
not the workstation dotfiles or setup repo; machine setup, shell configuration,
automation definitions, and restore payloads live in the dotfiles repository.

## What Is Here

- `bin/`: Windows-to-WSL Codex launch scripts, protocol proxy modules, and
  skills updater wrappers.
- `hooks/`: source copies of local Codex workstation hooks and focused
  regression tests.
- `plugins/`: repo-owned Codex plugin sources.
- `skills/`: repo-owned Codex skill sources.
- `docs/`: focused operator documentation for the shim, updater, hooks, and
  repo-owned skill workflows.
- `feedback/`: repo-local feedback capture for skill improvement work.
- `tests/`: focused regression checks for workstation scripts.

## Quick Checks

Run the Node proxy regression tests from the repo root:

```bash
node --test bin/codex-wsl-proxy-runtime.test.js
```

Probe the Windows skills updater wrappers from Windows PowerShell:

```powershell
bin\skills-updates.cmd --help
bin\sk-up.cmd -h
```

## Documentation

Start with [docs/README.md](docs/README.md) for the maintained docs index.

Common entry points:

- [WSL shim](docs/wsl-shim.md): Windows-to-WSL Codex launch and proxy behavior.
- [Skills updater](docs/skills-updater.md): global skill update wrappers.
- [Agent-Native boundary](docs/agent-native-boundary.md): what the installed
  BuilderIO/Agent-Native skills own versus what stays in this repo.
- [Skill feedback loop](docs/skill-feedback-loop.md): local feedback capture
  and improvement workflow.
- [Maintainer guide](docs/maintainer-guide.md): source-first workflow,
  verification, and local repair notes.

## Development Posture

Keep source changes in this repo first. Copy files into an active workstation
install only as an explicit install or repair step.

When changing behavior, update the smallest relevant docs page and run the
smallest relevant check. Keep secrets, raw session exports, and private
credential material out of the repository.
