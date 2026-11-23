---
name: council
description: >-
  Kick off agentic work with independent reviews from multiple model families, then
  verify, repair, and score the result. Use when the user says "council", "sol loop",
  "$confidence-loop", or asks for a multi-model review, red-team, premortem,
  loophole audit, numeric confidence score, or work with independent reviewers.
---

# Council

One magic word for the whole loop: scope the work, seat independent reviewers
across model families, run bounded rounds, verify and repair, then score and hand off.

Outcome: a scored, evidence-backed decision on whether the current change is ready,
with one finding ledger and an explicit confidence score. Proof: every material
finding has evidence, every risk has exactly one status, and each score comes from a
seat that actually ran. When no seat can run, say so and report the single-family
result rather than implying a council happened.

Use this when independent cross-family judgement or a quantified confidence decision
is needed. For a single-model reuse and slop pass, use `code-refactor-review` instead.

Aliases: `sol loop` and `confidence loop` / `$confidence-loop` mean this skill.
If the user names a single model, or no multi-model route exists, run single-family
mode and label the run `single family`.

## 0. Gates

- Read-only inspection and research may proceed within the active request.
- Edits are allowed only inside the scoped change. Never edit a scoped file while a
  review is running.
- No merge, force push, PR comment, or base-branch change unless the user explicitly
  authorizes it. A round ending is not approval to ship.
- Preserve unrelated work. Do not stash, discard, or include it.

## 1. Scope and proof

Run `git status --short`. Limit the review to change paths that separate cleanly;
otherwise stop and ask. After `git fetch origin`, pin the base:

- PR: `gh pr view <n> --json baseRefName,headRefOid`; confirm `HEAD` is the PR head,
  then `git merge-base HEAD origin/<baseRefName>`.
- Local: `git merge-base HEAD origin/<default-branch>`.
- No remote: a base the user names, or `HEAD` for uncommitted work only.

Use one diff command every round. Include untracked files. State intent and
out-of-scope in 1-2 sentences. Run covering checks before seating reviewers.
Record missing or pre-failing checks in the brief. Fix known in-scope issues and
relevant PR comments first. If the user is considering a rewrite, review as-is.

## 2. Seat the council

In T3 Code, call `orchestrator_capabilities` on the target machine. Require
`canRunChildTask` and, for cross-provider runs, `canRunCrossProviderChildTask`.
Use exact catalog IDs. Cap every seat's thinking effort at `medium`, including
re-reviews; if a route cannot honor the cap, report the seat as unavailable.
Outside T3, use the host's subagents.

Seat one strong model per family that adds a useful lens. Model names below are
hints, not requirements: use the strongest model currently available in that family.

| Family | Lens |
| --- | --- |
| OpenAI | correctness, specification, regressions |
| Anthropic | architecture, maintainability |
| Google | integration assumptions, edge cases |
| xAI | challenge assumptions, overlooked failure cases |
| DeepSeek | implementation detail, failure paths |
| Moonshot | long-context consistency, missing evidence |
| Zhipu | tests, portability |

Skip duplicate routes and opaque aliases (for example Copilot `auto`/`default`).
Report unavailable seats, never silently replace a named member, and label any
incomplete run `partial`.

## 3. Run bounded rounds

- Default: 3 rounds. Extend to 5 only for security, permissions, auth, data loss,
  or state/concurrency while a blocker or `should` stays open with evidence, or when
  the user asks.
- Round 1 is full. Later rounds cover changed code and open findings only.
- Stop early when a round returns only nits. Declined nits and comment, doc, or
  format-only fixes need no new round.
- Keep one ledger across all rounds and seats: `id`, `severity`, `source seat`,
  `claim`, `evidence`, `outcome`. Never guess severity, evidence, or usage.
- Every seat is read-only. No seat may edit files, commit, push, comment on PRs,
  or start its own reviewer. The main agent owns synthesis, fixes, and the decision.

Launch seats in parallel with the same request, scope, evidence, and base diff.
Give each seat its lens. Ask for prioritized findings with evidence and a score
out of 10. Do not edit scoped files while seats run. A timeout, error, or silence
is not approval; on failure, stop the seat, verify it stopped, then restart it with
the ledger. Do not poll or wait for non-blocking members.

Round 1 brief (child seats get no history):

```
Review only. Do not edit, commit, push, comment on PRs, or start subagents.
Repo: <abs path>  Base: <sha>  Intent: <intent>  Out of scope: <scope>  Checks: <results>
See: git diff <base> [-- <paths>]; git status --short (read untracked).
Read surrounding code. Find bugs, broken callers, mismatch with intent.
Skip style unless it hides a bug. Hold to the smallest change. Flag code,
tests, or guards the intent does not need. Other paths, rare cases, small
perf, extra compat: nits.
blocker = intended case wrong. should = real bug callers will hit, in intent. nit = rest.
Reply:
VERDICT: APPROVE or CHANGES
S1 [blocker|should|nit] <file>:<line> <claim>. Fix: <suggestion> (stable ids)
SCORE: <0-10>
```

## 4. Verify, repair, and re-score

Check every finding yourself before acting. Re-label severity if the evidence
supports it; another path, a rare case, or a small cost is a nit whatever label
it came with.

- Agree: apply the smallest fix. Prefer delete or narrow over a new lock, cache,
  helper, or option. Fix the case that happened. Add one focused test that fails
  without the fix only if intent needs it.
- Disagree: state why with evidence - file:line, test, command output.
- Out of scope or nit: decline in one line and list as follow-up.

For each material risk, mark exactly one status: `disproven`, `accepted`,
`blocked`, or `fixed`. Verify with the smallest relevant check: source read,
command, test, search, documentation, or reasoning proof. Record the result.
Leave any accepted fix outside the current scope behind an explicit approval gate.

Update the ledger, re-run checks, and run `git diff --shortstat <base>`. If review
grew the diff, cut it before the next round. Do not agree just to end the loop.

Start another round only if budget remains and a blocker or `should` is open,
behavior code changed, or approval does not cover the current behavior diff. Send
the round 1 brief, the ledger (`fixed:` / `disagree:` / `declined:` with evidence),
the same base diff, and the change since last round. Ask seats to cover open ids
and changed code only. Mark each old id `resolved` or `still open: <reason>`.

## 5. Stop and score

Agreed when the current behavior diff is APPROVE, no blocker or `should` is open,
and you have no concern of your own. Stuck when an item survives two replies or
the budget ends: show both sides in 2-3 lines each and let the user decide.
Budget end is not approval.

Verify findings, then summarize agreement, disagreement, and next steps. Give a
calibrated confidence score: 100 = all material assumptions verified and checks
pass; 90-99 = minor non-material uncertainty; 70-89 = meaningful uncertainty or
unverified assumptions; 40-69 = materially uncertain; 0-39 = weak or contradicted.
If 100 is impossible, say what evidence would close the gap.

Report in this format:

```markdown
## Council review

Mode: <full|single family>   Rounds: <N>   Seats: <families run / unavailable>

| Seat | Round 1 | Round 2 |
| --- | --- | --- |
| <family / model> | 8/10 | 10/10 |

Finding ledger: <fixed / disagreed / declined, with ids>
Verification: <checks and results>   Remaining uncertainty: <explicit>
Confidence: <0-100>   Decision: <agreed | stuck - user decides>
```

Use `no score` for a completed review without a score and `unavailable` for a
seat that could not run. If the work ends in a PR, put this section in the PR
description above the model and harness line; otherwise put it in the final reply.

## 6. Hand off

Do not push during the loop. Push only if the user asked to push, ship, or update
the PR; use `--force-with-lease` after a rebase and never force-push without
authorization. Do not merge unless the user explicitly asks for a merge.

Report: rounds, seats and families, what changed, diff size before and after,
fixed and declined findings with reasons, verification results, confidence score,
and whether anything was pushed.
