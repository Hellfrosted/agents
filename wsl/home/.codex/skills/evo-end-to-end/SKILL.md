---
name: evo-end-to-end
description: Run a Codex-only planning-to-Evo workflow. Use when the user wants to start from a vague performance, architecture, refactor, flaky-test, slow-build, or code-quality problem; optionally use grill-me/grill-with-docs/improve-codebase-architecture; produce an Evo-ready experiment brief; then hand the brief to `$evo discover` and `$evo optimize` with safe scope, metric, gate, budget, and merge rules.
---

# Evo End To End

## Best Full Sequence

Use this sequence for normal Codex-to-Evo work:

1. `$grill-me` - clarify the problem, constraints, risks, non-goals, measurable success, forbidden changes, and what would make the run safe for autonomous optimization.
2. Optional: `$grill-with-docs` - use when terminology, ownership boundaries, architectural language, or documented decisions are unclear. Update `CONTEXT.md`/ADRs only when decisions crystallize.
3. Optional: `$improve-codebase-architecture` - use when the problem is architectural, refactor-heavy, testability-related, or needs repo-informed decomposition before an Evo metric can be chosen.
4. Optional: draft a staged plan before the Evo brief when the improvement needs multiple dependent steps or every step must leave the repo working.
5. Repo inspection - prefer local evidence over user questions when commands, manifests, tests, docs, or code can answer safely.
6. Draft the Evo experiment brief - define goal, metric, baseline, gate, editable scope, read-only context, forbidden changes, budget, and merge rule.
7. Get user approval before running Evo if the run may edit production behavior, APIs, persistence, auth/security, tests, packaging, dependencies, deployment, or user-visible behavior.
8. `$evo discover` - hand Evo the approved brief so it can inspect the repo, instrument the benchmark, and record baseline.
9. `$evo optimize subagents=<n> budget=<n> stall=<n>` - run the optimization loop only after discovery/baseline is clear.
10. Review Evo results manually before merging if they touch production behavior, APIs, persistence, auth/security, tests, packaging, dependencies, deployment, or user-visible behavior.

Good starting prompt:

```text
Use $evo-end-to-end for <repo/path>.

Problem: I want to improve <problem>.

Grill me until we have a clear optimization plan. Then, after I approve it, create an Evo-ready brief with metric, baseline, gate, scope, forbidden changes, budget, and merge rule. Run Evo only if I say "start Evo".
```

Fast guided prompt:

```text
Use $evo-end-to-end for <repo/path>.

Problem: <problem>.

Inspect the repo, ask only blocking questions, draft an Evo-ready brief, and stop for approval before running Evo.
```
