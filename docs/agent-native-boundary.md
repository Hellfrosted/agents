# Agent-Native Boundary

BuilderIO's Agent-Native skills now own the general plan and recap workflows on
this workstation. Keep this repository focused on source that still has local
workstation value.

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
- `bin/skills-updates.*` and `bin/sk-up.cmd`: broad multi-source skill update,
  diff, skip, install, and uninstall wrappers for the universal `.agents`
  skill catalog.
- repo-owned local skills under `skills/`, such as ICM, Discrawl, Tuck, Yeet,
  Task Brief, Confidence Loop, and Evo End To End.
- `plugins/visual-canvas`: local and offline HTML report artifacts, reusable
  HTML output policy, visual profiles, static checks, and compact artifact
  structure.

## Visual Canvas Boundary

Visual Canvas should not duplicate Agent-Native plan or recap artifacts. When a
request is primarily a visual implementation plan, architecture plan, migration
plan, branch recap, diff recap, pull request recap, or agent-work recap, route
to the installed Agent-Native skill instead.

Keep Visual Canvas for requests that need a portable HTML artifact, a local
report outside the hosted plan app, profile guidance, policy checks against an
existing HTML file, or generated report assets whose structure should remain
repo-owned.

## Change Rule

When reducing overlap, prefer a source doc or routing change before deleting
code. Remove code only after current repo evidence shows it is unused or fully
delegated, and keep active global installs, plugin caches, lockfiles, and
machine configuration untouched unless the user explicitly asks for repair or
install work.
