# Skill Feedback Loop

The Skill feedback loop lets this repo learn from repeated Skill failures and
activation evals without changing active workstation behavior automatically.

The loop is source-first:

1. Record small structured feedback events under `feedback/skills/events/`.
2. Review events periodically with `$improve-skills`.
3. Let the review write summaries and, when requested, edit repo-owned Skill
   source.
4. Stop for human review before any Skill-source change is installed, copied,
   committed, pushed, or applied to active workstation behavior.

## What This Covers

This workflow is for repo-owned Skills under `skills/` in
`E:\dev\agents-toolkit`. It does not update globally installed Skills under
`%USERPROFILE%\.agents\skills` and does not edit plugin caches such as
`~/.codex/plugins/cache`.

Use it for durable behavior feedback and trigger-shaping work such as:

- A Skill missed a recurring risk.
- A Skill asked unnecessary questions when repo evidence was enough.
- A Skill produced friction that should be avoided next time.
- A repeated positive behavior should be preserved.
- A Skill description over-triggers, under-triggers, or collides with a nearby
  Skill.
- A blank-slate subagent eval exposes unintended behavior or weak guardrails.

Do not store secrets, credentials, private personal data, recovery codes, raw
logs, or raw chat exports in feedback files.

## Recording Feedback

Append JSON Lines events to `feedback/skills/events/`. Use one JSON object per
line:

```json
{"ts":"2026-06-17","skill":"task-brief","outcome":"miss","severity":"medium","actual":"treated global install like ordinary build work","expected":"flag machine-wide install as higher risk","source":"codex-thread","private":false}
```

Required fields:

- `ts`: ISO date or timestamp.
- `skill`: Skill directory or name, such as `task-brief`.
- `outcome`: `miss`, `near-miss`, `friction`, `win`, or `note`.
- `severity`: `low`, `medium`, `high`, or `critical`.
- `actual`: what the Skill or agent did.
- `expected`: what should happen next time.
- `source`: short non-sensitive origin, such as `codex-thread`, `manual`, or
  `review`.
- `private`: boolean. Use `true` when the outer loop should summarize the
  lesson without copying sensitive detail into Skill source.

Practical file names:

```text
feedback/skills/events/2026-06-17.jsonl
feedback/skills/events/thread-<short-id>.jsonl
```

## Running The Outer Loop

Run the review from this repo:

```text
In E:\dev\agents-toolkit, use $improve-skills to review feedback/skills.
If $improve-skills is not exposed as an active Skill, read
skills/improve-skills/SKILL.md directly and follow it. Find repeated Skill
failures or activation issues, make source-first edits only when requested,
verify changed paths are limited to allowed write surfaces, and stop before
installing or applying changes to active workstation behavior. Do not install,
copy, commit, push, or edit plugin/cache/global Skill directories.
```

The review should:

- Read `AGENTS.md`.
- Inspect `feedback/skills/events/*.jsonl` and
  `feedback/skills/summaries/*.md`.
- Group feedback by Skill and repeated failure type.
- For trigger-shaping requests, build activation matrices and use blank-slate
  subagent evals when available.
- Treat one-off low or medium feedback as a summary note.
- Require repeated similar evidence, ideally two or three events, before
  changing permanent Skill instructions.
- Allow a single high or critical event to justify a Skill edit only when the
  expected behavior is clear and likely to recur.
- Write rollups under `feedback/skills/summaries/`.
- Edit only repo-owned Skill source when evidence justifies it.
- Compare changed paths against the pre-run dirty baseline before the final
  response. Verify paths changed by the loop stay within the allowed surfaces,
  and report pre-existing out-of-scope dirty paths separately.

## Allowed Write Surfaces

For `$improve-skills`, `skills/improve-skills/SKILL.md` is the canonical source
for allowed and forbidden write surfaces. Feedback event files may be added when
the user explicitly asks to record new feedback.

Do not install, copy to global Skill directories, edit plugin/cache files,
commit, push, or update automations unless the user explicitly asks for that
separate operation.

## Weekly Automation Prompt

When creating a recurring Codex automation, use the automation tool rather than
raw automation text. The prompt body should be:

```text
Weekly Skill feedback review:

In E:\dev\agents-toolkit, use $improve-skills to review feedback/skills.
If $improve-skills is not exposed as an active Skill, read
skills/improve-skills/SKILL.md directly and follow it. Find repeated Skill
failures or activation issues, make source-first edits only when requested,
verify changed paths are limited to allowed write surfaces, and stop before
installing or applying changes to active workstation behavior. Do not install,
copy, commit, push, or edit plugin/cache/global Skill directories.
```

The automation should stop after summaries and any requested source edits. It
should not apply those changes to active installed Skills until the user
explicitly asks for install or repair work.

## Review Output

Each outer-loop run should use the output shape in
`skills/improve-skills/SKILL.md`.
