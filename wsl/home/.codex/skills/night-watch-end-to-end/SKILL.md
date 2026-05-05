---
name: night-watch-end-to-end
description: Run a Codex-only planning-to-handoff workflow for Night Watch. Use when the user wants to start from a problem, use grill-me/improve-codebase-architecture/to-prd/to-issues, and prepare AFK issues for Night Watch without invoking the native Night Watch CLI.
---

# Night Watch End To End

## Best Full Sequence

Use this sequence for normal Codex-to-Night-Watch work:

1. `$grill-me` - clarify the problem, constraints, risks, non-goals, success criteria, and what would make the work safe to run unattended.
2. Optional: `$grill-with-docs` - use when terminology, business concepts, ownership boundaries, or architectural language are unclear. Update `CONTEXT.md`/ADRs only when decisions crystallize.
3. Optional: `$improve-codebase-architecture` - use when the work is architectural, refactor-heavy, testability-related, or needs repo-informed decomposition before PRD/issues.
4. Optional: draft a staged plan before the PRD when the work needs staged commits, a larger refactor plan, or a sequence where every step must leave the repo working.
5. `$to-prd` - convert the approved plan into a parent PRD GitHub issue.
6. `$to-issues` - split the PRD into thin vertical slice issues. Mark each issue as `AFK` or `HITL`.
7. Night Watch handoff - prepare every unblocked `AFK` issue for the Night Watch board and mark it for `Ready` handling using available non-CLI tooling.
8. Leave `HITL`, blocked, unclear, production-risky, API-changing, dependency-changing, packaging-changing, migration-heavy, or security-sensitive issues out of `Ready` unless the user explicitly approves that issue for autonomous execution.
9. Do not invoke the native Night Watch CLI. If immediate execution is requested, use the available non-CLI Night Watch integration; otherwise leave prepared work for scheduled pickup.
10. Review Night Watch PRs manually before merging if they touch production behavior, APIs, persistence, auth/security, tests, packaging, dependencies, deployment, or user-visible behavior.

Good starting prompt:

```text
Use night-watch-end-to-end for <repo/path>.

Problem: I want to improve <problem>.

Grill me until we have a clear plan. Then, after I approve it, create the PRD, split it into issues, and prepare unblocked AFK issues for Night Watch Ready handling.
```

Fast fully automated prompt:

```text
Use night-watch-end-to-end for <repo/path>.

Problem: <problem>.

Inspect the repo, ask only blocking questions, draft the PRD/issues, and stop before marking work Ready if risk is unclear.
```
