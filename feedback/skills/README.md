# Skill Feedback Events

This directory stores source-controlled feedback for repo-owned Codex Skills.
Use it for small, structured notes about real Skill behavior that should inform
later source-first improvements.

Do not store secrets, credentials, recovery codes, private personal data, or raw
chat/session exports here. If an event includes sensitive context, set
`private` to `true` and summarize only what future Skill text needs to learn.

## Event Files

Append JSON Lines events under `events/`. A practical naming convention is one
file per day or per thread:

```text
feedback/skills/events/2026-06-17.jsonl
feedback/skills/events/thread-<short-id>.jsonl
```

Each line is one JSON object:

```json
{"ts":"2026-06-17","skill":"task-brief","outcome":"miss","severity":"medium","actual":"treated global install like ordinary build work","expected":"flag machine-wide install as higher risk","source":"codex-thread","private":false}
```

Fields:

- `ts`: ISO date or timestamp.
- `skill`: Skill directory or name, such as `task-brief`.
- `outcome`: `miss`, `near-miss`, `friction`, `win`, or `note`.
- `severity`: `low`, `medium`, `high`, or `critical`.
- `actual`: what the Skill or agent did.
- `expected`: what should happen next time.
- `source`: short non-sensitive origin, such as `codex-thread`, `manual`, or
  `review`.
- `private`: boolean. When `true`, outer-loop review should summarize the
  lesson without copying sensitive detail into Skill source.

## Summaries

The `summaries/` directory stores rollups produced by the outer-loop review.
Summaries should group repeated feedback by Skill and failure mode, separate
one-off notes from repeated evidence, and link proposed source diffs to the
events that justify them.
