---
name: architecture-cascade
description: Runs improve-codebase-architecture once, extracts every candidate, and fans out bounded candidates into worktree-backed worker threads only for that cascade turn. Use when the user wants architecture improvement candidates handled automatically instead of choosing one manually.
---

# Architecture Cascade

## Use

The user should only need:

```text
$architecture-cascade <optional repo area or problem>
```

Also use this when `$agent-loop-designer` is asked to run
`improve-codebase-architecture` and automatically explore or act on all
candidates.

## Non-Negotiables

- Worktree-backed worker threads are scoped to this cascade run only.
- Candidate worktrees must be created inside the main checkout under
  `.codex-worktrees/architecture-cascade/...`, and worker threads must target
  the same saved Codex project as the coordinator.
- Do not create or ask the user to open separate Codex GUI projects for
  candidate worktrees.
- Do not reinterpret later `subagent` or `sub-agent` requests as requests to
  create more threads. Those requests use normal subagent tooling unless the
  user explicitly invokes `$architecture-cascade` again.
- Do not create durable automations, recurring thread spawners, or custom
  subagent routing from this skill.
- The coordinator creates candidate workers only after discovery produced a
  concrete candidate queue.
- Candidate workers explore first, then edit only when `actability:
  auto-safe`.

## Required Docs

Before fan-out, read:

- `improve-codebase-architecture` skill instructions.
- Current Codex Worktree/thread docs or local thread tooling docs.
- `agent-loop-designer/references/worktree-threads.md`.
- Repo `AGENTS.md`, `CONTEXT.md`, `docs/adr/`, README/docs, and test/build
  guidance.

## Loop

1. Run one `improve-codebase-architecture` discovery pass.
2. Extract every candidate into a queue with id, files, problem, solution,
   strength, conflicts, likely tests, and proposed owner.
3. Create a practical batch of one-turn worktree-backed candidate workers.
4. Collect `implemented`, `rejected`, and `blocked` results.
5. Ask only for blocked candidates that need a domain decision, ADR reopening,
   unsafe permission change, or overlapping ownership resolution.
6. If the skill caused a repeatable setup failure, apply the self-improvement
   rule in [cascade-reference.md](references/cascade-reference.md).

## Output

Return an architecture action report:

- Implemented candidates: id, summary, changed files, checks, worker title/link.
- Rejected candidates: id, evidence, reason.
- Blocked candidates: id, exact user decision needed.
- Integration choices: completed worktrees ready for review or merge, with
  risks.
- Verification summary: checks passed, failed, or not run.

## Details

Use [cascade-reference.md](references/cascade-reference.md) for report location,
worktree placement, starting-state preflights, worker contracts, and the
self-improvement rule.
