---
name: evo-end-to-end
description: Run a Codex-only planning-to-Evo workflow. Use when the user wants to start from a vague performance, architecture, refactor, flaky-test, slow-build, or code-quality problem; optionally use grill-me/grill-with-docs/improve-codebase-architecture; produce an Evo-ready experiment brief; then hand the brief to `$evo discover` and `$evo optimize` with safe scope, metric, gate, budget, stall rule, merge rules, and tiny/small/medium/big/huge optimize presets.
---

# Evo End To End

Turn a fuzzy improvement request into an Evo-ready experiment, then run Evo only after approval.

## Flow

1. Clarify only what local repo inspection cannot answer. Use `$grill-me` when the problem, constraints, non-goals, success metric, or forbidden changes are unclear.
2. Use `$grill-with-docs` only for unclear terminology, ownership, or ADR/CONTEXT decisions.
3. Use `$improve-codebase-architecture` only when architecture or testability must be decomposed before choosing a metric.
4. Inspect the repo for manifests, tests, docs, benchmarks, and likely editable scope.
5. Draft an Evo brief with: goal, metric, baseline command/data, pass gate, editable scope, read-only context, forbidden changes, budget, stall rule, merge rule.
6. Stop for approval before Evo edits production behavior, APIs, persistence, auth/security, tests, packaging, dependencies, deployment, or user-visible behavior.
7. Run `$evo discover` with the approved brief; optimize only after discovery records a baseline.
8. Run `$evo optimize subagents=<n> budget=<n> stall=<n>` within the approved scope. Choose a preset below unless the user gives exact values.
9. Manually review Evo output before merging behavior, API, persistence, security, packaging, deployment, or user-visible changes.

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
Budget:
Stall rule:
Optimize preset:
Merge rule:
```
