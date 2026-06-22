---
name: tuck
description: Reviewed local git commit workflow with bounded subagent review.
disable-model-invocation: true
---

# Tuck

Use after `$tuck` or a plain `tuck` request where `tuck` is being used as the
git commit workflow. Keep review out of the main context: use subagents, but
batch them so review stays token-aware. The primary output is local commit(s),
not a push.

If the user forbids subagents, do not run this Skill; ask whether they want an
ordinary commit workflow instead.

If the current tool surface requires explicit user authorization before
spawning subagents and the user has not authorized subagents or the tuck
workflow in the current request, ask before spawning reviewers.

When Codex subagents are authorized, follow the delegation, approval-gate,
evidence, and bounded-loop protocol in
[`../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol`](../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol)
and
[`../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates`](../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates).

## Flow

1. Run `git status --short`. Complete when every modified and untracked path is
   classified as intended, unrelated user work, or needing user clarification.
2. Inspect intended diffs before staging. Treat unrelated modified/untracked
   files as user work. Complete when each intended path has been read enough to
   explain the commit intent and risk.
3. Split commits by reviewer-facing intent; separate formatting, docs, and
   behavior when practical. Complete when each planned commit has a path list,
   intent, and risk focus.
4. Main agent only scopes commits with `git diff --stat -- <paths>` and minimal
   path checks; do not deep-review diffs in main context. Complete when the
   reviewer packet contains the relevant diff scope without loading unrelated
   history into the main context.
5. Before staging each commit, spawn batched read-only subagent reviewers.
   Complete when every reviewer returns blocking findings or explicitly reports
   no blocking findings.
6. Fix blocking findings, then stage only intended paths or hunks. Complete
   when all blocking findings are fixed, accepted behind an explicit user gate,
   or the tuck stops.
7. Re-check `git diff --cached`, run the smallest relevant check, then commit
   with a concise imperative message. Report any skipped check. Complete when
   the staged diff matches the planned commit and the command result is known.
8. Stop after the local commit summary and verification result. Do not offer
   push as part of the default tuck flow.

## Subagent Review

Always use at least one reviewer per commit. Prefer one batched reviewer for
small, familiar, docs-only, formatting-only, or obvious mechanical changes.
Split reviewers by risk area for larger or mixed commits.

Use the shared delegation protocol for reviewer goal shape, read-only scope,
forking, evidence, and integration rules. Paste the commit scope, assigned
paths, risk focus, and relevant diff context into the reviewer message.

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
GOAL: {commit-review goal from the shared delegation protocol}
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

If the same request asks to push after tuck, complete the reviewed local commit
first, then run `git branch --show-current` and `git remote -v` for the explicit
push gate. Recommend the current branch when it exists and is not `main`;
otherwise recommend a new branch. Never push or force-push without explicit
confirmation.
