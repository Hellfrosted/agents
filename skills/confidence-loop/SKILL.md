---
name: confidence-loop
description: Confidence audit with a 0-100 score. Use when the user invokes $confidence-loop/$confident-loop, asks for a numeric confidence score or 100% certainty, or asks to red-team/premortem a proposed strategy for loopholes and failure modes. Not for ordinary review, implementation, skill-trigger audits, or operational loops unless the user explicitly wants a confidence audit.
---

# Confidence Loop

Adversarially verify a strategy, plan, implementation, or answer. Treat "100% confident" as a request for evidence, not reassurance.

## Invocation

`$confident-loop` is an invocation alias for `$confidence-loop`. A bare
invocation runs in normal mode. Accept these modes and aliases:

| Mode | Aliases | Reviewer budget | Repair budget |
| --- | --- | --- | --- |
| single | solo | none | one pass |
| normal | default, c, b | one reviewer | one pass |
| hard | a | two to four reviewers | two passes |
| extreme | s, sr, ssr | up to six reviewers | three passes |

An explicit user budget may raise a mode's limits.

Do not trigger on named operational loops such as feedback loops, automation
loops, CI loops, or workflow loops unless the user explicitly asks for a
confidence/loophole audit of that loop. If the user forbids subagents, use
single mode.

When reviewers are unavailable, forbidden, or not authorized, run single mode
and label the reduced reviewer coverage. Do not imply that it was a
multi-reviewer audit.

Follow the shared protocol for
[`../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol`](../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol)
and the dedicated
[`approval and side-effect gates`](../shared-agent-protocol/SKILL.md#approval-read-only-and-side-effect-gates).
For repeated passes, also apply
[`../shared-agent-protocol/SKILL.md#bounded-loops`](../shared-agent-protocol/SKILL.md#bounded-loops).

## Standard Loop

1. State the strategy and success criteria. Complete when the object under
   review and the pass/fail criteria are explicit enough for another reviewer
   to judge.
2. List material assumptions, dependencies, edge cases, and failure modes.
   Complete when each plausible blocker has an evidence target or is marked as
   accepted uncertainty.
3. Mark each risk as disproven, accepted, blocked, or needing a fix. Complete
   when every listed material risk has exactly one status.
4. Revise the strategy to close real loopholes. Complete when every accepted
   fix is either applied or left behind an explicit approval gate.
5. If the task repeats review, repair, monitoring, handoff, or decision work,
   define it with the shared bounded-loop fields. Complete when every shared
   field is filled or explicitly unnecessary.
6. Gather ordinary read-only evidence and run the smallest relevant
   verification: source read, documentation research, command, test, search, or
   reasoning proof. Complete when the evidence and result are recorded, or the
   reason verification could not run is explicit.
7. Repeat within the selected budget until no material unresolved loopholes
   remain, or stop with the remaining uncertainty explicit.

Apply the shared approval gates to side effects. Ordinary read-only inspection
and research may proceed within the active request.

## Sub-Agent Review

When reviewers are authorized and in scope, use the selected mode's reviewer
budget. Use only as many as the uncertainty justifies, with distinct angles.

Use the shared delegation protocol for reviewer goal shape, read-only scope,
forking, evidence, and integration rules. Complete sub-agent review only when
each reviewer has returned evidence-backed objections or explicitly reports no
material objections.

Choose reviewer angles that match the problem. Common angles:

- **Skeptic**: false assumptions, loopholes, and failure modes.
- **Verifier**: smallest checks that would prove or disprove the strategy.
- **Domain reviewer**: code, docs, APIs, product constraints, or domain-specific gaps.
- **Implementability reviewer**: hidden execution, sequencing, ownership, or integration risks.
- **Security reviewer**: abuse cases, data exposure, authorization, or unsafe execution risks.
- **UX reviewer**: user workflow, copy, accessibility, and interaction gaps.
- **Maintenance reviewer**: future breakage, dependency, operational, or handoff risks.
- **Loop reviewer**: trigger/state/stop-condition gaps, runaway autonomy,
  missing human gates, weak budget/stall controls, and failure-learning gaps.

Give each reviewer the strategy, criteria, evidence, and open questions. They must not edit files, spawn agents, or decide the final answer.

Prompt shape:

```
TASK: act as a read-only confidence-loop reviewer. Angle: {reviewer angle}.
GOAL: {reviewer goal from the shared delegation protocol}
DELIVERABLE: material loopholes with evidence, verification checks, speculative objections labeled as such, and confidence: 0-100.
SCOPE: no file edits, no spawned agents, no final decision.
VERIFY: cite the evidence or reasoning used for every material objection.
Strategy: {strategy}
Criteria: {criteria}
Evidence: {facts}
Open questions: {unknowns}
```

Integrate results yourself. Accept evidence-backed issues, reject unsupported
ones, fix valid gaps only when edits are in scope, then verify.

## Codex Loop Audit

When a reviewed strategy should become repeatable, choose the durable Codex
surface that fits the host and request: a goal, skill, automation, visible task
workflow, plugin workflow, or compact reusable command. Apply the shared
bounded-loop and approval gates regardless of surface.

Treat PR-comment loops as one example only, not a default recommendation.
Choose the loop from the user's setup and evidence.

If the loop is not worth making durable, say why and leave it as a one-off
strategy.

## Confidence Score

Use a score from 0 to 100, not broad labels. Calibrate it this way:

- **100**: all material assumptions verified; relevant checks pass; no blocker remains.
- **90-99**: evidence strongly supports the answer, with only minor non-material uncertainty left.
- **70-89**: likely correct, but meaningful uncertainty, incomplete coverage, or unverified assumptions remain.
- **40-69**: plausible but materially uncertain; more evidence or revision is needed.
- **0-39**: weak, contradicted, missing core evidence, or verification failed.

If 100% is impossible, say why and what evidence would close the gap.

## Output

Keep it compact: Strategy, Codex loop shape when relevant, Loopholes found, Second opinions, Verification, Mode/reviewer count/pass count, Confidence score. Do not pad with hypothetical risks that do not apply.
