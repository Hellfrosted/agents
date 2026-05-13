---
name: evo-end-to-end
description: Run a Codex planning-to-Evo workflow for evo-hq/evo v0.4.0+. Use when the user wants to start from a vague performance, architecture, refactor, flaky-test, slow-build, or code-quality problem; optionally use grill-me/grill-with-docs/improve-codebase-architecture; produce an Evo-ready experiment brief; then hand the brief to `$evo discover`, `$evo optimize`, and, when needed, Evo backend/runtime setup with safe scope, metric, gate, backend, host, budget, stall rule, and merge rules.
---

# Evo End To End

Turn a fuzzy improvement request into an Evo-ready experiment, then run Evo only after approval.

Target Evo release line: `evo-hq/evo` v0.4.0 or newer.

## Flow

1. Clarify only what local repo inspection cannot answer. Use `$grill-me` when the problem, constraints, non-goals, success metric, or forbidden changes are unclear.
2. Use `$grill-with-docs` only for unclear terminology, ownership, or ADR/CONTEXT decisions.
3. Use `$improve-codebase-architecture` only when architecture or testability must be decomposed before choosing a metric.
4. Inspect the repo for manifests, tests, docs, benchmarks, and likely editable scope.
5. Verify Evo is installable and version-aligned before running plugin skills:
   - `evo-version-check` when available from the installed plugin.
   - Otherwise `evo --version`; it must report `evo-hq-cli`, not the unrelated `evo` SLAM package.
   - Do not auto-install or upgrade the CLI unless the user explicitly asks.
6. Draft an Evo brief with: goal, metric, baseline command/data, pass gate, editable scope, read-only context, forbidden changes, host, backend, runtime/env needs, budget, stall rule, merge rule.
7. Stop for approval before Evo edits production behavior, APIs, persistence, auth/security, tests, packaging, dependencies, deployment, user-visible behavior, dependency manifests, or remote/cloud infrastructure.
8. Run `$evo discover` with the approved brief; optimize only after discovery records a baseline.
9. If using an existing Evo workspace from before v0.4.0, silently migrate host metadata with `evo host show`; if it prints `<not set>`, run `evo host set codex`.
10. For remote, pool, or non-default runtime setup, configure Evo explicitly before optimizing:
   - Local default: worktree backend.
   - Faster local reuse: pool backend with a fixed workspace list.
   - Remote: configure the provider first, using Evo's `infra-setup` guidance for Modal, E2B, Daytona, AWS, Azure, SSH, manual, or custom providers.
   - Runtime commands/env belong in `evo config runtime ...` and `evo env ...`, not hard-coded into benchmark scripts.
11. Run `evo run --check` when wiring risk is material and a non-mutating validation is available.
12. Run `$evo optimize subagents=<n> budget=<n> stall=<n>` within the approved scope. Choose a preset below unless the user gives exact values.
13. Use `evo direct "<text>"` only for mid-run steering of an already-running Evo session.
14. Manually review Evo output before merging behavior, API, persistence, security, packaging, deployment, or user-visible changes.

## Optimize Presets

Use these as starting points for `$evo optimize`:

- **tiny**: `subagents=3 budget=5 stall=2`
- **small**: `subagents=3 budget=8 stall=3`
- **medium**: `subagents=4 budget=10 stall=4`
- **big**: `subagents=5 budget=14 stall=5`
- **huge**: `subagents=8 budget=20 stall=6`

Default to **medium** for ordinary repo cleanup, focused architecture work, flaky-test diagnosis, or performance exploration. Use **tiny** or **small** when the editable scope is narrow or risky. Use **big** or **huge** only when the metric is stable, the baseline is repeatable, and the approved scope can absorb broader exploration.

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
Budget:
Stall rule:
Optimize preset:
Merge rule:
```

## Evo v0.4.0 Notes

- `evo init --host <claude-code|codex|opencode|openclaw|hermes|generic>` is required for new workspaces. For this skill on Codex, use `codex`.
- New workspaces default to the `pareto_per_task` frontier strategy instead of `argmax`. Existing workspaces keep their configured strategy.
- Local execution has two backends: `worktree` and `pool`. Pool mode is useful when setup is expensive, but it changes commit discipline because warm workspace state should stay out of commits.
- Remote experiments can run through Modal, E2B, Daytona, AWS, Azure, SSH, manual, or custom providers. Treat provider SDK installation, credentials, and cloud allocation as explicit user-approved setup.
- Backend provider credentials and benchmark runtime environment are separate concerns. Configure benchmark variables with `evo env`, and do not copy secrets into worktrees or docs.
- Use `evo gc` to clean worktrees, pool slots, and remote sandboxes across configured backends.
- Use `evo config show`, `evo config backend show`, `evo config runtime show`, and `evo env show` to inspect setup before changing it.
