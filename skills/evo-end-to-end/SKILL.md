---
name: evo-end-to-end
description: "Evo campaign planner for evo-hq/evo. Use when the user explicitly names Evo/evo-hq for optimization or autoresearch, invokes $evo-end-to-end, or asks for an Evo-ready experiment brief. Not for ordinary optimization, refactors, testing, training, or unrelated uses of Evo."
---

# Evo End To End

Turn a fuzzy improvement request into an Evo-ready experiment, then run Evo only
after approval.

Target Evo release line: `evo-hq/evo` v0.5.2 or newer.

For companion skill sources, optimize presets, and release-specific notes, see
[REFERENCE.md](REFERENCE.md).

Approval boundary: commands in this Skill or `REFERENCE.md` that update,
install, force reinstall, trust hooks, change remote providers, change
credentials, garbage-collect workspaces, or arm autonomous execution are
report-only until the user explicitly approves that action and target.

When this workflow uses Codex subagents for planning, review, or triage, follow
the delegation, approval-gate, evidence, and bounded-loop protocol in
[`../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol`](../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol),
[`../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates`](../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates),
and
[`../shared-agent-protocol/SKILL.md#bounded-loops`](../shared-agent-protocol/SKILL.md#bounded-loops).

## Companion Skills

Names written as [`$skill-name`](codex://skills) or plugin IDs such as
`evo:discover` are agent skills. Load and follow a named skill when its trigger
applies.

- Use [`$grilling`](codex://skills) when the goal, constraints, non-goals, success metric,
  ownership, or forbidden changes are unclear after repo inspection.
- Add [`$domain-modeling`](codex://skills) alongside [`$grilling`](codex://skills) when unclear terminology,
  `CONTEXT.md`, ADR decisions, or domain vocabulary need to be sharpened or
  updated while planning.
- Use [`$improve-codebase-architecture`](codex://skills) when architecture or testability must be
  decomposed before choosing an Evo metric.
- Use [`evo:finetuning`](codex://skills) before writing or changing training code, reward design,
  model-weight update recipes, or post-training workflows.
- Use installed Evo plugin skills whose `evo_version` matches `evo --version`.

If a referenced [`$xxx`](codex://skills) skill is not installed, check the host's skill list and
then the user's skill lockfile for `sourceUrl` and `skillPath`. Cite the source
before continuing with the best available fallback or hardstop as needed.

## Flow

1. Inspect the repo for manifests, tests, docs, benchmarks, and likely editable
   scope before asking clarifying questions. Complete when the likely baseline,
   validation commands, and forbidden areas are known or listed as unknowns.
2. Verify Evo before running plugin skills:
   - `evo --version` must report `evo-hq-cli`, not the unrelated SLAM package.
   - The CLI version must match the loaded Evo plugin skill version.
   - If CLI and plugin drift, tell the user the lockstep updater command
     (`evo update <host> --version <version>`). Do not install or upgrade unless
     the user asks.
   Complete when the CLI identity/version and plugin skill version are recorded,
   or the missing check is a blocker.
3. Draft an Evo brief with: goal, metric, baseline command/data, pass gate,
   editable scope, read-only context, forbidden changes, host, backend,
   runtime/env needs, per-experiment timeout, budget, stall rule, task skills,
   and merge rule. Complete when every brief field is filled or explicitly
   marked unknown.
4. Stop for approval before Evo edits production behavior, APIs, persistence,
   auth/security, tests, packaging, dependencies, deployment, user-visible
   behavior, dependency manifests, or remote/cloud infrastructure. Complete
   only after the user approves the affected scope, or stop with a report-only
   brief.
5. Run [`evo:discover`](codex://skills) with the approved brief. Optimize only after discovery
   records a baseline and the benchmark reviewer gate has passed. Complete
   when the discovery experiment id, baseline proof, and reviewer gate result
   are recorded.
6. Configure backend/runtime explicitly for pool, remote, or non-default
   runtimes. Keep benchmark variables in `evo env`, not in committed scripts or
   docs. Complete when backend, runtime, and env placement are recorded.
7. Set or confirm `--per-exp-timeout` / `evo config set per-exp-timeout
   <seconds>` for long benchmarks or training runs. Complete when the timeout
   value is recorded or the run is confirmed short enough not to need one.
8. Run `evo run <exp_id> --check` when wiring risk is material and non-mutating
   validation is available. Complete when the check result is recorded, or the
   unavailable check is named.
9. Resolve `autonomous` and `subagents-only` the same way [`evo:optimize`](codex://skills) does,
   then arm that state with `evo autonomous on|off` and `evo subagents-only
   on|off` only after the user has approved the values. Complete when both
   values and their approval state are recorded.
10. Size `subagents=<n>` from benchmark/backend resources first. Use width 1
    for exclusive GPUs, fixed ports, shared databases, singleton services, or
    timing-sensitive harnesses unless the harness isolates them.
    This is Evo CLI worker sizing, not Codex subagent fork behavior. Complete
    when `subagents=<n>` is justified by resource isolation and benchmark risk.
11. Run [`evo:optimize`](codex://skills) with `subagents=<n> budget=<n> stall=<n>` within the approved
    scope. Complete when the optimize experiment id and approved limits are
    recorded.
12. Use `evo direct "<text>"` only to steer an already-running Evo session. If
    an agent receives an `[EVO DIRECTIVE id=...]` banner, it must run
    `evo ack <event_id>` before proceeding. Complete when each directive is
    acknowledged before action.
13. Manually review Evo output before merging behavior, API, persistence,
    security, packaging, deployment, or user-visible changes. Complete when
    material output is accepted, rejected, or left behind an explicit approval
    gate.

## Output

Report the approved brief plus the evidence that exists so far:

- Evo CLI and plugin versions checked.
- Discovery or optimize experiment id, when created.
- Baseline command/data and pass gate proof.
- Backend, runtime/env, budget, stall rule, autonomous value, and
  subagents-only value.
- `--check`, discover, optimize, and manual-review results.
- Remaining approvals before merge, install, cleanup, push, or deploy.

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
