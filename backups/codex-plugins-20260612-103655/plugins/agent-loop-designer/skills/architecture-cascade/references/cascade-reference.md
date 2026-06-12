# Architecture Cascade Reference

## Scope Boundary

Architecture cascade is turn-scoped. It may create worktree-backed worker
threads only for a current `$architecture-cascade` request or an equivalent
`$agent-loop-designer` request. Later `subagent`, `sub-agent`, or custom agent
requests are ordinary subagent requests.
For those ordinary role-specific subagent requests, use a non-full-history
spawn and put the role, specialty, constraints, and required context in the
message. In the current `spawn_agent` tool, omit `fork_context` or set
`fork_context: false`; on tool surfaces that use `fork_turns`, set
`fork_turns: "none"`. Do not combine a full-history fork with agent type,
model, or reasoning overrides.

## Report Location Guardrails

Do not rely on WSL `/tmp` when the user may reopen the report later.

- Copy or rewrite the final report to a durable non-repo path before opening or
  reporting it.
- For repos on mounted Windows drives, prefer a sibling durable temp folder on
  that drive, such as `/mnt/e/tmp/architecture-review-<timestamp>.html`.
- When giving a Windows browser URL for a mounted-drive report, use a Windows
  file URL such as `file:///E:/tmp/architecture-review-<timestamp>.html`.
- If OS temp is unavoidable, say that the report is volatile.

## Worktree Placement

Use Same-Project Nested Manual Worktree Mode. The worker prompt carries the
nested worktree cwd; the Codex thread target remains the main project.

- Create manual Git worktrees under
  `<repo>/.codex-worktrees/architecture-cascade/<run-id>/<candidate-id>`.
- Add `.codex-worktrees/` to `.git/info/exclude`; do not edit `.gitignore`.
- Use `git worktree add --detach <manual-path> <verified-ref>`.
- Start worker threads under the current saved project; the worker prompt names
  the nested worktree cwd for all reads, edits, `git status`, and validation.
- Do not create Codex-managed Worktree threads for architecture cascade.
- Do not create external-root worktrees for architecture cascade.
- Do not ask the user to add or open candidate worktree paths as separate Codex
  projects.
- If same-project thread creation cannot target the main saved project, stop and
  report the blocker. Do not create a new GUI project as a workaround.

This keeps every worker under the checkout the user already trusts in Codex.

## Starting State

Manual Git worktrees and Codex `create_thread` setup require verified existing
refs.

- `startingState: { type: "branch", branchName: "<ref>" }` starts from an existing ref; it does not create a branch.
- Never pass a new candidate branch name as `branchName` unless
  `git rev-parse --verify --quiet <ref>^{commit}` confirms it exists.
- Prefer the current existing branch ref after verifying it.
- If the current branch cannot be resolved, use another verified existing ref such as `main`; if none exists, stop and ask.
- Do not use `startingState: { type: "working-tree" }` for cascade fan-out.
  It replays the coordinator checkout's uncommitted diff into every Worktree.
- Candidate workers may create or switch branches inside their own worktree only after setup succeeds and the repo workflow permits it.

Minimum batch preflight:

```bash
git status --short
git branch --show-current
git rev-parse --verify --quiet "$(git branch --show-current)^{commit}"
```

## Worker Contracts

Candidate worker:

- Role: candidate owner.
- Worktree: one isolated manual git worktree under the main checkout.
- Ownership: exact candidate files, directly related tests, and candidate-specific docs.
- Write intent: explore first, then edit only when `actability` is `auto-safe`.
- Output: candidate id, actability, evidence, changed files, tests run,
  blockers, merge concerns.

Coordinator:

- Owns discovery, queue extraction, worker creation, status collection, and final summary.
- Does not edit candidate code directly while candidate workers are active.
- Does not merge worktree results automatically unless the user explicitly asks.

## Act Automatically Means

Act on every candidate that is bounded, non-conflicting, verified by its worker,
and permitted by the current session. Stop for product/domain decisions, ADR
conflicts, overlapping ownership, unavailable permissions, or no credible
verification path.

## Self-Improvement Rule

When repeatable coordinator failure is caused by this skill's instructions,
update the source plugin and active cached plugin before closing the turn.

- Patch the narrowest relevant skill or reference file.
- Base the update on observed command output, app error text, rejected API shape, missing durable artifact, or thread setup behavior.
- Add guardrails or preflights. Do not add incident logs, generated-by notes, or
  user-private raw logs.
- Validate source/cache parity when a source copy exists, search stale wording with `rg`, and run the smallest relevant bundled script check.
- Store a concise privacy-safe memory for durable failures after the patch
  succeeds.
