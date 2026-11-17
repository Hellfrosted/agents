---
name: evo-end-to-end
description: "Evo campaign planner for evo-hq/evo. Use when the user explicitly names Evo/evo-hq for optimization or autoresearch, invokes $evo-end-to-end, or asks for an Evo-ready experiment brief. Not for ordinary optimization, refactors, testing, training, or unrelated uses of Evo."
---

# Evo End To End

Turn an optimization request into a bounded Evo campaign. For a brief-only
request, stop after drafting the brief.

## Inspect and brief

Inspect the repo, existing benchmarks, and validation entry points. Use
`$grilling` when a consequential scope or metric decision remains unclear.
Use the installed Evo skills for discovery, optimization, reporting, and shipping.
Check that `evo --version` reports `evo-hq-cli` and matches their `evo_version`.
For drift or missing hooks, see [installation diagnostics](REFERENCE.md#installation-diagnostics).
On Codex, follow the [shared agent protocol](../shared-agent-protocol/SKILL.md)
when delegating, reviewing evidence, or running bounded loops.

Record the brief before allocating experiments:

- Goal, metric and direction, baseline command, and pass/fail gate.
- Editable scope, protected behavior, and changes outside scope.
- Host, backend, runtime, benchmark environment, and positive per-experiment
  timeout. Identify hosted services and expected cost where applicable.
- Round width, per-worker budget, stall rounds, and any campaign-wide cap.
  Size width and depth using the matching Evo plugin's round-sizing reference.
- Autonomous and subagents-only settings, task skills, and winner criteria.
  For noisy scores, define matched conditions, repeat count, summary statistic,
  and minimum worthwhile improvement before promoting a winner.

Use the user's existing authorization within scope. Obtain approval for any
additional installation, hook trust, infrastructure or credential change,
workspace garbage collection, or autonomous execution. Confirm changes to
protected production behavior, APIs, data, auth, tests, dependencies,
packaging, deployment, or user-visible behavior before widening scope.

## Discover and validate

Invoke the installed discover skill with the brief. Set the timeout at
initialization and configure the intended backend, runtime, and `evo env`
before allocating or scoring the baseline. Require a passing benchmark review
before scoring. When wiring risk is material, run the discovery skill's
`--check` before the scored baseline and retain its artifact path. Record the
baseline experiment ID, score, and parity gate result.

When Evo starts its dashboard on a tailnet host, automatically publish it
through [Tailscale Serve](REFERENCE.md#dashboard-access) and share the verified
HTTPS URL. Record its cleanup details in campaign state.

## Optimize

Use the installed optimize skill within the approved scope and limits.
Resolve and record autonomous and subagents-only settings before starting.
Apply the approved settings and keep unattended execution within the agreed
limits. Reconfirm a finalist against the brief's promotion criteria. For a
hosted CI campaign, see [measurement and parity](REFERENCE.md#hosted-ci).

## Report and clean up

Use the installed report skill to summarize recorded evidence. Include the
baseline and finalist experiment IDs and scores, gate results, and remaining
decisions. Keep the dashboard route across waits while monitoring continues.
When monitoring ends, including after cancellation or terminal failure,
[verify and remove the campaign's route](REFERENCE.md#dashboard-access).

## Ship if approved

Use the installed ship skill after the user selects a qualifying winner and
approves landing it. Review the distilled diff against the brief and parity
evidence. Follow the target repository's rules for branches, pushes, pull
requests, and merge authorization.
