# agents-toolkit Documentation

Focused docs for the scripts, skills, and operator workflows in this repository.

## Core Docs

| Topic | Doc | Use it for |
| --- | --- | --- |
| WSL shim | [wsl-shim.md](wsl-shim.md) | Windows-to-WSL Codex launch behavior, proxy runtime, path translation, fallback skills listing, and focused verification. |
| Skills updater | [skills-updater.md](skills-updater.md) | `skills-updates` and `sk-up` behavior, install and uninstall modes, state paths, locking, and Go wrapper checks. |
| sk-up implementation | [sk-up-go-port.md](sk-up-go-port.md) | Current Go updater implementation reference, command contract, validation matrix, and consolidated completion audit. |
| Hook guardrails | [hook-guardrails.md](hook-guardrails.md) | WSL shell footgun guardrails, repo-local Impeccable hook adapter, install shape, and focused verification. |
| Agent-Native boundary | [agent-native-boundary.md](agent-native-boundary.md) | What BuilderIO/Agent-Native skills own now, what this repo still owns, and how Visual Canvas delegates plan/recap work. |
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
| `hooks/` | Source and tests for repo-owned Codex workstation hook guardrails. Active global hook installs are repaired only on request. |
| `plugins/visual-canvas/` | Source for portable local HTML report artifacts, visual profiles, and HTML output policy checks. Agent-Native owns general visual plan and recap workflows. |
| `skills/` | Repo-owned Codex skill sources. In this repository, "work on skills" means this directory unless a task says otherwise. |
| `tests/` | Script-level checks for updater promotion shape and Windows wrapper behavior. |

## Public Boundary

The root README stays intentionally short. Detailed local maintenance notes live
in this directory so the public project entry point explains the repository
without doubling as a workstation runbook.

Workstation setup, shell configuration, Codex automation definitions, AgentsView
notes, Discrawl wiretap notes, and restore archives now belong in
`/mnt/e/dev/dotfiles`.
