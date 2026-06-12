---
name: evo-end-to-end
description: Run a Codex planning-to-Evo workflow for evo-hq/evo v0.5.2+. Use when the user wants to turn a vague performance, architecture, refactor, flaky-test, slow-build, fine-tuning, post-training, or code-quality problem into an Evo-ready experiment brief, using companion skills such as $grill-me, $grill-with-docs, and $improve-codebase-architecture when needed.
---

# Evo End To End

Turn a fuzzy improvement request into an Evo-ready experiment, then run Evo only
after approval.

Target Evo release line: `evo-hq/evo` v0.5.2 or newer.

For companion skill sources, optimize presets, and release-specific notes, see
[REFERENCE.md](REFERENCE.md).

## Companion Skills

Names written as `$skill-name` are agent skills. Load and follow a named skill
when its trigger applies.

- Use `$grill-me` when the goal, constraints, non-goals, success metric, or
  forbidden changes are unclear after repo inspection.
- Use `$grill-with-docs` only for unclear terminology, ownership, `CONTEXT.md`,
  or ADR decisions.
- Use `$improve-codebase-architecture` when architecture or testability must be
  decomposed before choosing an Evo metric.
- Use `$evo finetuning` before writing or changing training code, reward design,
  model-weight update recipes, or post-training workflows.
- Use installed Evo plugin skills whose `evo_version` matches `evo --version`.

If a referenced `$xxx` skill is not installed, check the host's skill list and
then the user's skill lockfile for `sourceUrl` and `skillPath`. Cite the source
before continuing with the best available fallback.

## Flow

1. Inspect the repo for manifests, tests, docs, benchmarks, and likely editable
   scope before asking clarifying questions.
2. Verify Evo before running plugin skills:
   - `evo --version` must report `evo-hq-cli`, not the unrelated SLAM package.
   - The CLI version must match the loaded Evo plugin skill version.
   - If CLI and plugin drift, tell the user the lockstep updater command
     (`evo update <host> --version <version>`). Do not install or upgrade unless
     the user asks.
3. Draft an Evo brief with: goal, metric, baseline command/data, pass gate,
   editable scope, read-only context, forbidden changes, host, backend,
   runtime/env needs, per-experiment timeout, budget, stall rule, task skills,
   and merge rule.
4. Stop for approval before Evo edits production behavior, APIs, persistence,
   auth/security, tests, packaging, dependencies, deployment, user-visible
   behavior, dependency manifests, or remote/cloud infrastructure.
5. Run `$evo discover` with the approved brief. Optimize only after discovery
   records a baseline and the benchmark reviewer gate has passed.
6. Configure backend/runtime explicitly for pool, remote, or non-default
   runtimes. Keep benchmark variables in `evo env`, not in committed scripts or
   docs.
7. Set or confirm `--per-exp-timeout` / `evo config set per-exp-timeout
   <seconds>` for long benchmarks or training runs.
8. Run `evo run <exp_id> --check` when wiring risk is material and non-mutating
   validation is available.
9. Resolve `autonomous` and `subagents-only` the same way `$evo optimize` does,
   then arm that state with `evo autonomous on|off` and `evo subagents-only
   on|off`.
10. Size `subagents=<n>` from benchmark/backend resources first. Use width 1
    for exclusive GPUs, fixed ports, shared databases, singleton services, or
    timing-sensitive harnesses unless the harness isolates them.
    This is Evo CLI worker sizing, not Codex `spawn_agent` fork behavior. If
    this workflow also uses Codex subagents for planning, review, or triage,
    launch role-specific workers without full history. In the current
    `spawn_agent` tool, omit `fork_context` or set `fork_context: false`; on
    tool surfaces that use `fork_turns`, set `fork_turns: "none"`. Put the role
    and needed context in the message; do not override `agent_type`, `model`, or
    `reasoning_effort` on a full-history fork.
11. Run `$evo optimize subagents=<n> budget=<n> stall=<n>` within the approved
    scope.
12. Use `evo direct "<text>"` only to steer an already-running Evo session. If
    an agent receives an `[EVO DIRECTIVE id=...]` banner, it must run
    `evo ack <event_id>` before proceeding.
13. Manually review Evo output before merging behavior, API, persistence,
    security, packaging, deployment, or user-visible changes.

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
