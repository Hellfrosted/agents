# agents-toolkit

Codex workstation tooling, scripts, and local skills.

This repository is public source for a mostly personal setup. It is not a
general-purpose installer, but the pieces are kept documented and testable so
they can be reused, audited, or repaired without relying on memory.

## What Is Here

- `bin/`: Windows-to-WSL Codex launch scripts, protocol proxy modules, and
  skills updater wrappers.
- `skills/`: repo-owned Codex skill sources.
- `docs/`: focused operator documentation for the shim, updater, shell setup,
  companion tools, and local archives.
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
