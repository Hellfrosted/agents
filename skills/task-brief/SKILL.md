---
name: task-brief
description: Task brief generator for compact Codex execution briefs. Use for implementation tickets, delegation prompts, future-work handoffs, or rough intent that needs concrete scope and evidence. Not for ordinary implementation/review, confidence audits, conversation-compaction handoffs, Evo briefs, PRDs, image prompts, generic writing prompts, or non-Codex prompt engineering.
---

# Task Brief

Write a compact brief that fixes scope, boundaries, and proof. Ask only when
missing detail materially changes outcome, risk, credentials, deployment,
hardware, destructive actions, or ownership.

## Flow

1. Choose the branch: build/fix, research, review, operate, or future-work
   handoff. Complete when the goal names one observable outcome.
2. Pin context and semantics. Complete when the brief names the relevant
   repo/path/link/current state and defines any domain term that must not drift.
3. Set boundaries. Complete when `Do` states allowed actions and `Do not`
   states safety limits or non-goals.
4. Set proof and mode. Complete when `Evidence` says what verifies done and
   `Mode` is light or heavy with the trigger named.
5. Pin the deliverable. Complete when `Output` names its format and destination
   or audience.
6. Emit the canonical shape only; add branch-specific detail inside its fields.
   Complete when every canonical field has a concrete value or `none`.

## Canonical Shape

```text
Goal: <one concrete outcome>
Context: <paths, repo, links, current state>
Semantic target: <domain term that must not drift, or none>
Do: <allowed actions>
Do not: <safety boundaries and non-goals>
Evidence: <what proves done>
Output: <file/report/commit/PR/table/etc>
Mode: <light or heavy, with trigger if heavy>
```

## Branch Notes

- Build/fix: include the observable failure or workflow, smallest acceptable
  change, focused validation, and real-surface QA when user-facing. For a
  behavior fix, `Evidence` names the pre-fix RED reproduction or failing test
  when feasible, the post-fix GREEN command, and any UI/API proof.
- Research: require primary sources for unstable or high-stakes facts, separate
  fact from inference, and cite files/URLs with dates when relevant. Evidence is
  source-backed findings; no mutation unless explicitly requested.
- Review: make the requested output findings-first, with each material finding
  tied to source lines or equivalent evidence and with verification gaps stated.
  Include edits only when the user requests review-and-fix.
- Operate: require current-state inspection, dry-run/read-only checks when
  available, before/after evidence, and explicit approval for destructive,
  credential, exposure, machine-wide, or production-risk actions.
- Future-work handoff: include current state, exact next action, known blockers,
  allowed write set, and resume evidence. Conversation-compaction handoffs
  depend on the external `$handoff` skill. If it is unavailable, say so and use
  this canonical shape as a compact fallback.

## Mode

Default to light for ordinary prompt shaping and implementation briefs. Use
heavy only when the brief drives security, money, deployment, legal/medical
risk, production operations, credentials, irreversible data, concurrency,
external APIs, cross-module contracts, or architecture decisions.

Heavy mode keeps the canonical shape but adds a concrete delta:

- `Do` names sequencing, owners, and rollback or recovery actions.
- `Do not` names approval gates and the exact boundary each gate protects.
- `Evidence` covers boundary checks, failure-path proof, and rollback proof.
- `Output` includes an audit trail and unresolved risks or decisions.

Heavy mode is complete only when each applicable delta is present or explicitly
marked unnecessary.

## Example

```text
Goal: identify client-only ATM10 mods that are also server-capable.
Context: client mods path: <path>; official list: <path or URL>.
Semantic target: "sideness" means upstream mod loader side support, not install location.
Do: compare filenames, resolve mod identity, classify mod sideness from upstream docs/source.
Do not: infer sidedness from where the file is installed.
Evidence: table with source link or local source evidence per mod; uncertain items separated.
Output: markdown report.
Mode: light unless adding/removing server mods or touching deployment config.
```
