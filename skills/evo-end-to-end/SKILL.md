---
name: evo-end-to-end
description: Run a Codex planning-to-Evo workflow for evo-hq/evo v0.5.0+. Use when the user wants to turn a vague performance, architecture, refactor, flaky-test, slow-build, fine-tuning, post-training, or code-quality problem into an Evo-ready experiment brief, using companion skills such as $grill-me, $grill-with-docs, and $improve-codebase-architecture when needed.
---

# Evo End To End

Turn a fuzzy improvement request into an Evo-ready experiment, then run Evo only after approval.

Target Evo release line: `evo-hq/evo` v0.5.0 or newer.

## Companion Skills

Names written as `$skill-name` are agent skills, not shell commands or prose labels. Load and follow the named skill when its trigger applies; if the skill is not installed, tell the user which source provides it before continuing with the best available fallback.

| Skill | Use in this workflow | GitHub source if missing |
| --- | --- | --- |
| `$grill-me` | Clarify an ambiguous goal, constraint set, non-goal list, success metric, or forbidden change. | `https://github.com/mattpocock/skills.git` at `skills/productivity/grill-me/SKILL.md` |
| `$grill-with-docs` | Challenge a plan against project terminology, `CONTEXT.md`, or ADRs, and update those docs as decisions crystallize. | `https://github.com/mattpocock/skills.git` at `skills/engineering/grill-with-docs/SKILL.md` |
| `$improve-codebase-architecture` | Decompose architecture or testability work before choosing an Evo metric. | `https://github.com/mattpocock/skills.git` at `skills/engineering/improve-codebase-architecture/SKILL.md` |
| `$evo discover`, `$evo optimize`, `$evo finetuning`, `$evo infra-setup` | Evo plugin skills installed by `evo install <host>`; use the installed evo plugin bundle whose `evo_version` matches `evo --version`. | `https://github.com/evo-hq/evo.git` under `plugins/evo/skills/` |

For any other `$xxx` reference encountered while executing this skill, first check the host's installed skill list. If absent, check the user's skill lockfile (usually `%USERPROFILE%\.agents\.skill-lock.json` on Windows or `/mnt/c/Users/<user>/.agents/.skill-lock.json` from WSL) for `sourceUrl` and `skillPath`, then cite that source to the user.

## Flow

1. Clarify only what local repo inspection cannot answer. Use `$grill-me` when the problem, constraints, non-goals, success metric, or forbidden changes are unclear.
2. Use `$grill-with-docs` only for unclear terminology, ownership, or ADR/CONTEXT decisions.
3. Use `$improve-codebase-architecture` only when architecture or testability must be decomposed before choosing a metric.
4. Inspect the repo for manifests, tests, docs, benchmarks, and likely editable scope.
5. Verify Evo is installable and version-aligned before running plugin skills:
   - Run `evo --version`; it must report `evo-hq-cli`, not the unrelated `evo` SLAM package.
   - Compare the CLI version to the installed evo plugin skill version. If a loaded evo skill is tagged `evo_version: 0.5.0`, `evo --version` must print exactly `evo-hq-cli 0.5.0`.
   - If the CLI and plugin drift, use Evo's lockstep updater (`evo update <host> --version <version>`). If Codex is pinned to a stale local marketplace, refresh that marketplace source before reinstalling the plugin.
   - Do not auto-install or upgrade the CLI unless the user explicitly asks.
6. Draft an Evo brief with: goal, metric, baseline command/data, pass gate, editable scope, read-only context, forbidden changes, host, backend, runtime/env needs, per-experiment timeout, budget, stall rule, task skills, merge rule.
7. Stop for approval before Evo edits production behavior, APIs, persistence, auth/security, tests, packaging, dependencies, deployment, user-visible behavior, dependency manifests, or remote/cloud infrastructure.
8. If the task involves fine-tuning, post-training, reward design, or weight updates, use `$evo finetuning` before writing or changing training code. Require literature research before the first train experiment, keep held-out eval data out of training, and set a smoke-run gate.
9. Run `$evo discover` with the approved brief; optimize only after discovery records a baseline and the benchmark reviewer gate has passed.
10. If using an existing Evo workspace from before v0.4.0, silently migrate host metadata with `evo host show`; if it prints `<not set>`, run `evo host set codex`.
11. For remote, pool, or non-default runtime setup, configure Evo explicitly before optimizing:
   - Local default: worktree backend.
   - Faster local reuse: pool backend with a fixed workspace list.
   - Remote: configure the provider first, using Evo's `infra-setup` guidance for Modal, E2B, Daytona, AWS, Azure, SSH, manual, or custom providers.
   - Runtime commands/env belong in `evo config runtime ...` and `evo env ...`, not hard-coded into benchmark scripts.
12. Set or confirm `--per-exp-timeout` / `evo config set per-exp-timeout <seconds>` for long benchmarks or training runs. Use `evo run <exp_id> --timeout <seconds>` only as a per-call override.
13. Run `evo run <exp_id> --check` when wiring risk is material and a non-mutating validation is available.
14. Before optimizing, resolve run behavior the same way `$evo optimize` does:
   - `autonomous` defaults on unless the user or stored defaults turn it off.
   - `subagents-only` defaults on unless the user or stored defaults turn it off.
   - Arm the resolved state with `evo autonomous on|off` and `evo subagents-only on|off`.
15. Before choosing `subagents=<n>`, read Evo's sizing guidance and size width from the binding resource: exclusive GPU/port/DB implies width 1 unless the harness isolates it; pool caps at slot count; remote caps at quota/cost.
16. Run `$evo optimize subagents=<n> budget=<n> stall=<n>` within the approved scope. Size the round from benchmark/backend resources first; use the presets below only as fallbacks or user-facing shorthand.
17. Use `evo direct "<text>"` only for mid-run steering of an already-running Evo session. If an agent receives an `[EVO DIRECTIVE id=...]` banner, it must run `evo ack <event_id>` before proceeding.
18. Manually review Evo output before merging behavior, API, persistence, security, packaging, deployment, or user-visible changes.

## Optimize Presets

Use these as fallbacks for `$evo optimize` when the benchmark resource profile is unknown and the user gave no exact values:

- **tiny**: `subagents=3 budget=5 stall=2`
- **small**: `subagents=3 budget=8 stall=3`
- **medium**: `subagents=4 budget=10 stall=4`
- **big**: `subagents=5 budget=14 stall=5`
- **huge**: `subagents=8 budget=20 stall=6`

Default to **medium** only when the benchmark is light, isolated, and no better sizing signal is available. Reduce `subagents` to 1 for exclusive resources such as a GPU, fixed port, shared database, or serialized fixture. Cap pool runs at the pool slot count. Use **tiny** or **small** when the editable scope is narrow or risky. Use **big** or **huge** only when the metric is stable, the baseline is repeatable, and the approved scope can absorb broader exploration.

## Brief Template

```text
Goal:
Metric:
Baseline:
Gate:
Editable scope:
Read-only context:
Forbidden changes:
Host: codex
Backend: worktree | pool | remote:<provider>
Runtime/env:
Per-exp timeout:
Budget:
Stall rule:
Autonomous: on | off
Subagents-only: on | off
Task skills:
Optimize preset:
Merge rule:
```

## Evo v0.5.0 Notes

- v0.5.0 adds a first-class `$evo finetuning` path for SFT, LoRA, preference optimization, RFT, and RL training moves. Use it before writing or changing `train.py`, reward code, or model-weight update recipes.
- Training runs must keep the held-out benchmark out of training data, perform literature research before the first training experiment, and use smoke-run validation before spending full budget.
- New Evo workspaces should set a realistic per-experiment timeout at init with `--per-exp-timeout <seconds>` or later with `evo config set per-exp-timeout <seconds>`. Override individual calls with `evo run <exp_id> --timeout <seconds>`.
- `task-skills` config lets discovery record task-category skills, such as `finetuning`, that subagents should load on demand. Inspect it with `evo config get task-skills`.
- `$evo optimize` requires resource-bound round sizing. Read Evo's `sizing-the-round.md` before passing a concrete `subagents=N`; use width 1 for exclusive GPUs, ports, singleton services, shared mutable fixtures, or unknown timing-sensitive benchmarks.
- `evo wait` now has process, log, GPU, and ideator selectors for waiting on long-running work without burning context.
- `evo abort` stops the experiment subprocess tree cross-platform, including detached benchmark or training children that would otherwise survive the driver.
- The dashboard now supports live log tailing, per-experiment annotations, and `EVO_DASHBOARD_HOST` for binding on cloud or Modal hosts.
- The installed `evo-hq-agent`, `@evo-hq/evo-agent`, and `@evo-hq/pi-evo` SDK packages should match the 0.5.0 line when SDK instrumentation is used.

## Evo v0.4.5 Notes

- v0.4.5 fixes Codex hook installation for Codex 0.130+ by registering and staging the plugin under the owner-name path Codex resolves (`evo@evo-hq`) and by validating the resolved `evo-hook-drain` binary in `evo doctor codex`.
- Existing Codex installs with exit-127 hook failures do not self-heal with `evo update`. Recover with `uv tool install --force evo-hq-cli && evo install codex --force`.
- `evo install codex --force` stages `evo-hook-drain` into the Codex plugin cache and removes stale legacy registrations. It may leave hooks untrusted; trust them through `/hooks` or run `evo install codex --trust-hooks` only when the user explicitly approves skipping the hook review.

## Evo v0.4.4 Notes

- `evo init --host <claude-code|codex|cursor|opencode|openclaw|hermes|pi|generic>` is required for new workspaces. For this skill on Codex, use `codex`.
- New workspaces default to the `pareto_per_task` frontier strategy instead of `argmax`. Existing workspaces keep their configured strategy.
- Local execution has two backends: `worktree` and `pool`. Pool mode is useful when setup is expensive, but it changes commit discipline because warm workspace state should stay out of commits.
- Pool mode defaults to `commit_strategy=tracked-only`; subagents must `git add` new source files and pass `--i-staged-new-files yes` to `evo run`.
- Remote experiments can run through Modal, E2B, Daytona, AWS, Azure, SSH, manual, or custom providers. Treat provider SDK installation, credentials, and cloud allocation as explicit user-approved setup.
- In remote mode, subagent briefs must state the experiment id explicitly and require `--exp-id <id>` on every `evo bash/read/write/edit/glob/grep` command.
- Backend provider credentials and benchmark runtime environment are separate concerns. Configure benchmark variables with `evo env`, and do not copy secrets into worktrees or docs.
- `evo run <exp_id> --check` validates benchmark/gate wiring without committing, evaluating, or consuming retry budget.
- `$evo optimize` defaults to autonomous, subagents-only operation. The user can override either explicitly, or via `evo config get default-autonomous`, `evo defaults get autonomous`, `evo config get default-subagents-only`, and `evo defaults get subagents-only`.
- `evo direct "<text>" --wait` expects an agent to acknowledge delivered directives with `evo ack <event_id>`.
- Use `evo gc` to clean worktrees, pool slots, and remote sandboxes across configured backends.
- Use `evo config show`, `evo config backend show`, `evo config runtime show`, and `evo env show` to inspect setup before changing it.
