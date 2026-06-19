---
name: confidence-loop
description: Stress-tests a strategy, plan, implementation approach, answer, or Codex loop design with sub-agent second opinions until remaining uncertainty is explicit and evidence-backed, then reports a 0-100 confidence score. Use when the user asks whether Codex is 100% confident, asks to find loopholes or failure modes, requests a confidence audit, says to run a loop until the strategy is factually solid, asks to harden a Codex loop, invokes $confident-loop or $confidence-loop, or invokes $confident-loop/$confidence-loop normal/hard/extreme/n/h/x.
---

# Confidence Loop

Adversarially verify a strategy, plan, implementation, or answer. Treat "100% confident" as a request for evidence, not reassurance.

## Invocation

If the user invokes bare `$confident-loop` or `$confidence-loop`, ask which mode to use before running:

- **normal** or **n**: default mode; use one sub-agent reviewer.
- **hard** or **h**: use two to four sub-agent reviewers.
- **extreme** or **x**: use as many sub-agent reviewers as the uncertainty, scope, and risk justify.

If the user invokes `$confident-loop normal`, `$confident-loop hard`, `$confident-loop extreme`, `$confident-loop n`, `$confident-loop h`, `$confident-loop x`, or the same forms with `$confidence-loop`, run immediately in that mode. Treat `default` as `normal`, `n` as `normal`, `h` as `hard`, and `x` as `extreme`.

Do not accept numeric shorthand as a mode. If the user appends a number such as
`1`, `2`, or `3`, treat it as ambiguous because agents may interpret it as a
requested sub-agent count. Ask whether they meant `n`, `h`, or `x`.

## Standard Loop

1. State the strategy and success criteria.
2. List material assumptions, dependencies, edge cases, and failure modes.
3. Mark each risk as disproven, accepted, blocked, or needing a fix.
4. Revise the strategy to close real loopholes.
5. If the task has repeated review, repair, follow-up, monitoring, handoff, or
   decision work, evaluate the Codex loop shape before finalizing:
   - Trigger: what starts or wakes the loop.
   - State: where progress, findings, decisions, and queue items live.
   - Next action: what the agent does after each pass.
   - Stop condition: what concrete state ends the loop.
   - Human gate: actions that require explicit user approval.
   - Loop budget: max iterations, wakeups, wall time, spend, or scope.
   - Failure learning: what durable skill/plugin/doc update prevents repeats.
6. Run the smallest relevant verification: source read, command, test, search, or reasoning proof.
7. Repeat until no material unresolved loopholes remain.

## Sub-Agent Review

Always use sub-agent reviewers. The mode controls only how many reviewers to use:

- **normal/default/n**: spawn one read-only reviewer.
- **hard/h**: spawn two to four read-only reviewers. Use only as many as the uncertainty justifies.
- **extreme/x**: spawn as many read-only reviewers as needed to cover the material uncertainty. Use focused batches with distinct angles until the remaining uncertainty is explicit and evidence-backed.

Give each reviewer a dedicated goal. The main agent must not draft that goal
itself. First spawn a dedicated goal-writer subagent that uses `$goalcraft` to
turn the task, criteria, and assigned angle into a reviewer goal, then returns
only that goal to the main agent. The main agent then passes the returned goal
to the reviewer. The goal must preserve the reviewer boundary: read-only, no
spawned agents, and no final decision.

When spawning reviewers in Codex, use non-full-history forks for role-specific
review. In the current `spawn_agent` tool, omit `fork_context` or set
`fork_context: false`; on tool surfaces that use `fork_turns`, set
`fork_turns: "none"`. Put the reviewer role, angle, constraints, and needed
context in the `message`. Do not combine a full-history fork with `agent_type`,
`model`, or `reasoning_effort` overrides; full-history forks inherit those
fields from the parent and will be rejected if overridden.

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
GOAL: {dedicated reviewer goal returned by the $goalcraft goal-writer subagent}
DELIVERABLE: material loopholes with evidence, verification checks, speculative objections labeled as such, and confidence: 0-100.
SCOPE: no file edits, no spawned agents, no final decision.
VERIFY: cite the evidence or reasoning used for every material objection.
Strategy: {strategy}
Criteria: {criteria}
Evidence: {facts}
Open questions: {unknowns}
```

Integrate results yourself. Accept evidence-backed issues, reject unsupported ones, fix valid gaps, then verify.

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

Keep it compact: Strategy, Codex loop shape when relevant, Loopholes found, Second opinions, Verification, Confidence score. Do not pad with hypothetical risks that do not apply.
