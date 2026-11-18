# agents-toolkit

Codex tooling source, scripts, hooks, and local skills.

This repository is public source for the Codex-adjacent tools owned here. It is
not the workstation dotfiles or setup repo; machine setup, shell configuration,
automation definitions, and restore payloads live in the dotfiles repository.

## What Is Here

- `hooks/`: source copies of local Codex workstation hooks and focused
  regression tests.
- `skills/`: repo-owned Codex skill sources.
- `docs/`: focused operator documentation for hooks and repo-owned skill workflows.

## Quick Checks

Run the workstation hook tests from the repo root in WSL.

```bash
python3 -B -m unittest hooks/test_wsl_command_guardrails.py
```

## Documentation

Start with the [OpenWiki quickstart](openwiki/quickstart.md) to find the right
source, workflow, or check. The [operator documentation index](docs/README.md)
links to the detailed maintenance and repair notes.

Common entry points:

- [Hook guardrails](docs/hook-guardrails.md): repo-owned Codex hook sources and
  verification.
- [Agent-Native boundary](docs/agent-native-boundary.md): what the installed
  BuilderIO/Agent-Native skills own versus what stays in this repo.
- [Maintainer guide](docs/maintainer-guide.md): source-first workflow,
  verification, and local repair notes.

## Development Posture

Keep source changes in this repo first. Copy files into an active workstation
install only as an explicit install or repair step.

When changing behavior, update the smallest relevant docs page and run the
smallest relevant check. Keep secrets, raw session exports, and private
credential material out of the repository.
