---
name: confidence-loop
description: Confidence audit with a 0-100 score. Use when the user invokes $confidence-loop/$confident-loop, asks for a numeric confidence score or 100% certainty, or asks to red-team/premortem a proposed strategy for loopholes and failure modes. Not for ordinary review, implementation, skill-trigger audits, or operational loops unless the user explicitly wants a confidence audit.
---

# Confidence Loop

Adversarially verify a strategy, plan, implementation, or answer. Treat "100% confident" as a request for evidence, not reassurance.

## Invocation

If the user invokes bare `$confident-loop` or `$confidence-loop`, ask which mode to use before running:

- **normal**, **c**, or **b**: default mode; use one sub-agent reviewer.
- **hard** or **a**: use two to four sub-agent reviewers.
- **extreme**, **s**, **sr**, or **ssr**: use up to six sub-agent reviewers unless the user gives a larger budget.

If the user invokes `$confident-loop normal`, `$confident-loop hard`, `$confident-loop extreme`, `$confident-loop c`, `$confident-loop b`, `$confident-loop a`, `$confident-loop s`, `$confident-loop sr`, `$confident-loop ssr`, or the same forms with `$confidence-loop`, run immediately in that mode. Treat `default`, `c`, and `b` as `normal`; `a` as `hard`; and `s`, `sr`, and `ssr` as `extreme`.

Do not trigger on named operational loops such as feedback loops, automation
loops, CI loops, or workflow loops unless the user explicitly asks for a
confidence/loophole audit of that loop. If the user forbids subagents, do not
run this Skill; perform an ordinary non-skill review if possible or report the
conflict.

If the current tool surface requires explicit user authorization before
spawning subagents and the user has not authorized subagents, delegation, or
parallel review in the current request, ask before spawning reviewers or perform
a single-agent confidence audit and label it as such.

When Codex subagents are authorized and needed, follow the delegation,
approval-gate, evidence, and bounded-loop protocol in
[`../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol`](../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol)
and
[`../shared-agent-protocol/SKILL.md#bounded-loops`](../shared-agent-protocol/SKILL.md#bounded-loops).

Default loop budgets:

- **normal/default/c/b**: one reviewer and one repair pass.
- **hard/a**: two to four reviewers and at most two repair passes.
- **extreme/s/sr/ssr**: at most six reviewers and three repair passes unless
  the user gives an explicit larger budget.

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
5. If the task has repeated review, repair, follow-up, monitoring, handoff, or
   decision work, evaluate the Codex loop shape before finalizing. Complete
   when all loop fields below are filled or explicitly unnecessary:
   - Trigger: what starts or wakes the loop.
   - State: where progress, findings, decisions, and queue items live.
   - Next action: what the agent does after each pass.
   - Stop condition: what concrete state ends the loop.
   - Human gate: actions that require explicit user approval.
   - Loop budget: max iterations, wakeups, wall time, spend, or scope.
   - Failure learning: what durable skill/plugin/doc update prevents repeats.
6. Run the smallest relevant verification: source read, command, test, search,
   or reasoning proof. Complete when the verification result is recorded, or
   the reason it could not run is explicit.
7. Repeat within the selected budget until no material unresolved loopholes
   remain, or stop with the remaining uncertainty explicit.

Only make file edits, network calls, installs, credential access, destructive
commands, automations, commits, or pushes when the active user request already
authorizes that action. Otherwise report the needed fix or approval gate.

## Sub-Agent Review

When subagents are authorized and in scope, use sub-agent reviewers. The mode
controls only how many reviewers to use:

- **normal/default/c/b**: spawn one read-only reviewer.
- **hard/a**: spawn two to four read-only reviewers. Use only as many as the uncertainty justifies.
- **extreme/s/sr/ssr**: spawn up to six read-only reviewers unless the user gives a larger budget. Use focused batches with distinct angles until the remaining uncertainty is explicit and evidence-backed.

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

When a reviewed strategy can become a repeatable Codex loop, prefer loopifying it
instead of leaving the user with a one-off prompt. Keep the loop bounded and
human-gated. Do not use Codex Goal Control or `/goal` mechanics as the loop
primitive; define the loop in ordinary Codex surfaces such as a skill,
automation prompt, visible thread workflow, worktree-thread workflow, plugin
workflow, or compact reusable command.

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
