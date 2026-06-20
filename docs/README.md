# agents-toolkit Documentation

Focused docs for the scripts, skills, and operator workflows in this repository.

## Core Docs

| Topic | Doc | Use it for |
| --- | --- | --- |
| WSL shim | [wsl-shim.md](wsl-shim.md) | Windows-to-WSL Codex launch behavior, proxy runtime, path translation, fallback skills listing, and focused verification. |
| Skills updater | [skills-updater.md](skills-updater.md) | `skills-updates` and `sk-up` behavior, install and uninstall modes, state paths, locking, and PowerShell checks. |
| Agent-Native boundary | [agent-native-boundary.md](agent-native-boundary.md) | What BuilderIO/Agent-Native skills own now, what this repo still owns, and how Visual Canvas delegates plan/recap work. |
| Skill feedback loop | [skill-feedback-loop.md](skill-feedback-loop.md) | Repo-local feedback events, weekly review flow, and source-first skill updates. |
| Maintainer guide | [maintainer-guide.md](maintainer-guide.md) | Source-first workflow, common checks, install freshness checks, and design rationale. |

## Root Context

| Doc | Use it for |
| --- | --- |
| [../MISSION.md](../MISSION.md) | Durable project purpose, success criteria, constraints, and non-goals. |
| [../PRODUCT.md](../PRODUCT.md) | Product posture, users, design principles, and anti-references. |
| [../RESOURCES.md](../RESOURCES.md) | Repo-owned contracts and external references used when validating changes. |

## Repo Areas

| Path | Purpose |
| --- | --- |
| `bin/` | Source for launchers, proxy modules, path translation, skills fallback, and updater wrappers. |
| `plugins/visual-canvas/` | Source for portable local HTML report artifacts, visual profiles, and HTML output policy checks. Agent-Native owns general visual plan and recap workflows. |
| `skills/` | Repo-owned Codex skill sources. In this repository, "work on skills" means this directory unless a task says otherwise. |
| `feedback/skills/` | Captured feedback events and summaries for skill improvement. |
| `tests/` | Script-level checks that are practical to run from this repository. |

## Public Boundary

The root README stays intentionally short. Detailed local maintenance notes live
in this directory so the public project entry point explains the repository
without doubling as a workstation runbook.

Workstation setup, shell configuration, Codex automation definitions, AgentsView
notes, Discrawl wiretap notes, and restore archives now belong in
`/mnt/e/dev/dotfiles`.
