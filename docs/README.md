# agents-toolkit Documentation

Focused docs for the scripts, skills, and operator workflows in this repository.

## Core Docs

| Topic | Doc | Use it for |
| --- | --- | --- |
| WSL shim | [wsl-shim.md](wsl-shim.md) | Windows-to-WSL Codex launch behavior, proxy runtime, path translation, fallback skills listing, and focused verification. |
| Skills updater | [skills-updater.md](skills-updater.md) | `skills-updates` and `sk-up` behavior, install and uninstall modes, state paths, locking, and PowerShell checks. |
| Skill feedback loop | [skill-feedback-loop.md](skill-feedback-loop.md) | Repo-local feedback events, weekly review flow, and source-first skill updates. |
| Companion tooling | [codex-cli-tooling.md](codex-cli-tooling.md) | Notes for adjacent Codex CLIs, MCP-backed docs, Evo, CodSpeed, LazyCodex, ICM, and local archive tools. |
| Shell setup | [shell-setup.md](shell-setup.md) | Starship, fzf, zoxide, Atuin, PSReadLine, ble.sh, Tabby, and cross-shell parity. |
| AgentsView | [agentsview.md](agentsview.md) | Local Codex/OpenCode session browser setup and verification. |
| Discrawl wiretap | [discrawl-wiretap.md](discrawl-wiretap.md) | Local Vesktop cache archive setup, limits, and verification. |
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
| `skills/` | Repo-owned Codex skill sources. In this repository, "work on skills" means this directory unless a task says otherwise. |
| `feedback/skills/` | Captured feedback events and summaries for skill improvement. |
| `tests/` | Script-level checks that are practical to run from this repository. |
| `docs/shell/` | Shell configuration files referenced by the shell setup docs. |

## Public Boundary

The root README stays intentionally short. Detailed local maintenance notes live
in this directory so the public project entry point explains the repository
without doubling as a workstation runbook.
