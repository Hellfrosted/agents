---
name: tuck
description: Reviewed local git commit workflow with bounded subagent review.
---

# Tuck

Use after `$tuck` or a plain `tuck` request where `tuck` is being used as the
git commit workflow. Keep review out of the main context: use subagents, but
batch them so review stays token-aware. The primary output is local commit(s),
not a push.

If the user forbids subagents, do not run this Skill; ask whether they want an
ordinary commit workflow instead.

A `$tuck` or plain `tuck` request authorizes the required reviewers and local
commit. If the host requires a separate tool confirmation, obey that gate.

When Codex subagents are authorized, follow the delegation, approval-gate,
evidence, and bounded-loop protocol in
[`../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol`](../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol)
and
[`../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates`](../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates).

## Flow

1. Run `git status --short`, `git diff --cached`, and `git diff`; inspect
   untracked files that may be in scope. Classify every changed hunk as intended
   for a named commit, unrelated user work, or needing clarification. Existing
   staged hunks are part of the candidate commit unless the user explicitly
   chooses another disposition. Complete when every hunk is classified without
   altering unrelated user work.
2. Split intended hunks by reviewer-facing intent; separate formatting, docs,
   and behavior when practical. Complete when each planned commit has an exact
   hunk set, intent, and risk focus.
3. Stage only the intended hunks for one planned commit. Never use broad adds.
   Then capture the exact candidate patch with `git diff --cached --binary` and
   its stat. The candidate is the complete `HEAD`-to-index patch, including any
   hunks that were already staged before tuck began. If the index contains a
   hunk outside the planned commit, stop for user direction rather than
   unstaging it. Complete when the candidate patch contains every and only
   planned hunk.
4. Spawn batched read-only subagent reviewers with that exact candidate patch
   plus minimal nearby context. Complete when every reviewer reports blocking
   findings or explicitly reports no blocking findings for the whole candidate.
5. Fix blocking findings and restage only affected intended hunks. Any material
   candidate change invalidates the prior review: recapture the complete cached
   patch and rerun the relevant reviewer lanes. Complete when all blocking
   findings are fixed, explicitly accepted by the user, or tuck stops, and the
   final cached patch is fully reviewed.
6. Run the smallest relevant check and report any skipped check. Immediately
   before commit, confirm `git diff --cached --binary` still matches the final
   reviewed candidate exactly, and record `git rev-parse HEAD`. Complete when
   the candidate is unchanged and verification is known.
7. Commit with a concise imperative message. Completion requires a successful
   commit exit, a changed `HEAD`, and readback of the created commit with
   `git show --stat --oneline --decorate HEAD`. A failed command or unchanged
   `HEAD` is not completion.
8. Stop after the local commit summary and verification result. Do not offer a
   push as part of the default tuck flow.

## Subagent Review

Always use at least one reviewer per commit. Prefer one batched reviewer for
small, familiar, docs-only, formatting-only, or obvious mechanical changes.
Split reviewers by risk area for larger or mixed commits.

Use the shared delegation protocol for reviewer scope, forking, evidence, and
integration rules. Paste the complete candidate patch, commit intent, risk
focus, and any minimal nearby context into the reviewer message.

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
GOAL: determine whether the exact candidate patch is safe and complete for its stated intent.
DELIVERABLE: blocking findings only, with file paths and line references; if none, say: no blocking findings.
SCOPE: review only changed diff plus minimal nearby context; do not edit files or spawn agents.
VERIFY: every finding must be tied to a concrete changed line or omitted required check.
Commit intent: {COMMIT_INTENT}
Candidate patch: {COMPLETE_HEAD_TO_INDEX_PATCH}
Focus: {RISK_OR_QUESTION}
```

Do not commit until reviewers return. Re-run the affected reviewer lanes after
any material candidate change. If subagents are unavailable, stop and ask the
user.

## Git Safety

Do not stage untracked files by default. Use `git add <specific paths>` or
non-interactive patch staging. No destructive commands.

If the same request explicitly asks to push after tuck, that request is the push
authorization: do not ask for duplicate confirmation. Complete the reviewed
local commit first, then run `git branch --show-current` and `git remote -v`.
Ask only when the destination is ambiguous or the requested push would expand
scope. A normal push never authorizes force-push; force-push requires a separate
explicit request.
