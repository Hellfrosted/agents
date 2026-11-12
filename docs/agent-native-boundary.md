# Visualization and Agent-Native Boundary

Codex's built-in visualization capability owns conversational visualizations,
interactive tools, charts, diagrams, maps, plots, and similar visual output.
BuilderIO's Agent-Native skills own the general plan and recap workflows on this
workstation. Keep this repository focused on source that still has local
workstation value.

## Built In To Codex

Use Codex's built-in visualization capability instead of maintaining or
publishing a local visualization plugin from this repository.

## Delegated To Agent-Native

Use the installed Agent-Native skills for:

- `/visual-plan`: implementation plans, architecture proposals, migrations, and
  staged technical work that should become a visual plan.
- `/visual-recap`: branch, commit, pull request, or diff recaps.
- `/agent-watchdog`, `/plan-arbiter`, `/plow-ahead`,
  `/efficient-frontier`, `/stay-within-limits`, `/quick-recap`, and
  `/read-the-damn-docs`: their named orchestration and instruction behaviors.

Those skills are installed under `$HOME/.codex/skills` and are updated through
their own Agent-Native update path. Do not copy their source into this
repository just to customize routine behavior.

## Kept Here

This repository still owns:

- `bin/codex-wsl*`: the Windows-to-WSL Codex app-server shim, protocol proxy,
  path translation, and skills fallback.
- `cmd/sk-up`, `internal/skup/...`, `bin/sk-up.cmd`, and
  `bin/skills-updates.cmd`: broad multi-source skill update, diff, skip,
  install, and uninstall behavior for the universal `.agents` skill catalog.
- repo-owned local skills under `skills/`, such as Discrawl, Tuck,
  Task Brief, Confidence Loop, and Evo End To End.

## Change Rule

Keep active global installs, plugin caches, lockfiles, and machine configuration
untouched unless the user explicitly asks for repair or install work.
