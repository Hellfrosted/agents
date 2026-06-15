---
name: tuck
description: Tucks local changes into focused local commits with token-aware subagent review. Commit-only workflow; do not push unless the user explicitly asks after tuck completes. Use when the user says $tuck or tuck, not for ordinary commit requests.
---

# Tuck

Use after `$tuck` or a plain `tuck` request. Keep review out of the main context:
use subagents, but batch them so review stays token-aware. The primary output is
local commit(s), not a push.

## Flow

1. Run `git status --short`.
2. Inspect intended diffs before staging. Treat unrelated modified/untracked files as user work.
3. Split commits by reviewer-facing intent; separate formatting, docs, and behavior when practical.
4. Main agent only scopes commits with `git diff --stat -- <paths>` and minimal
   path checks; do not deep-review diffs in main context.
5. Before staging each commit, spawn batched read-only subagent reviewers.
6. Fix blocking findings, then stage only intended paths or hunks.
7. Re-check `git diff --cached`, commit with a concise imperative message, then run the smallest relevant check.
8. Stop after the local commit summary and verification result. Do not offer push as part of the default tuck flow.

## Subagent Review

Always use at least one reviewer per commit. Prefer one batched reviewer for
small, familiar, docs-only, formatting-only, or obvious mechanical changes.
Split reviewers by risk area for larger or mixed commits.

Give each reviewer a dedicated commit-review goal. The main agent must not
draft that goal itself. First spawn a dedicated goal-writer subagent that uses
`$goalcraft` to turn the commit scope, assigned paths, and risk focus into a
reviewer goal, then returns only that goal to the main agent. The main agent
then passes the returned goal to the reviewer. The goal must keep the reviewer
read-only and commit-scoped.

When spawning Codex reviewers, use non-full-history forks for role-specific
review. In the current `spawn_agent` tool, omit `fork_context` or set
`fork_context: false`; on tool surfaces that use `fork_turns`, set
`fork_turns: "none"`. Paste the commit scope, assigned paths, risk focus, and
relevant diff context into the `message`. Do not combine a full-history fork
with `agent_type`, `model`, or `reasoning_effort` overrides; full-history forks
inherit those fields from the parent and will be rejected if overridden.

Add targeted reviewers for:

- Security, credentials, auth, persistence, paths, process execution,
  install/update logic, or public APIs.
- Large, unfamiliar, highly coupled, or generated diffs.
- A suspected blocker needing independent scrutiny.
- Explicit user request for deep/exhaustive review.

Use at most 3 reviewers per commit unless the user asks for exhaustive review.
Use per-file review only when one file is the risk boundary.

Short reviewer prompt:

```text
TASK: act as a read-only tuck reviewer.
GOAL: {dedicated commit-review goal returned by the $goalcraft goal-writer subagent}
DELIVERABLE: blocking findings only, with file paths and line references; if none, say: no blocking findings.
SCOPE: review only changed diff plus minimal nearby context; do not edit files or spawn agents.
VERIFY: every finding must be tied to a concrete changed line or omitted required check.
Commit scope: {COMMIT_SCOPE}
Assigned paths: {PATHS}
Focus: {RISK_OR_QUESTION}
```

Do not stage until reviewers return. Re-run review only when a material fix
changes reviewed risk. If subagents are unavailable, stop and ask the user.

## Git Safety

Do not stage untracked files by default. Use `git add <specific paths>` or non-interactive patch staging. No destructive commands.

If the user separately asks to push after tuck, run `git branch --show-current`
and `git remote -v`. Recommend the current branch when it exists and is not
`main`; otherwise recommend a new branch. Never push or force-push without
explicit confirmation.
