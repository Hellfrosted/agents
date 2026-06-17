# Skill Feedback Loop

The Skill feedback loop lets this repo learn from repeated Skill failures
without changing active workstation behavior automatically.

The loop is source-first:

1. Record small structured feedback events under `feedback/skills/events/`.
2. Review events periodically with `$improve-skills`.
3. Let the review write summaries and propose diffs to repo-owned Skill source.
4. Stop for human review before any proposed Skill-source change is installed,
   copied, committed, pushed, or applied to active workstation behavior.

## What This Covers

This workflow is for repo-owned Skills under `skills/` in
`E:\dev\agents-toolkit`. It does not update globally installed Skills under
`%USERPROFILE%\.agents\skills` and does not edit plugin caches such as
`~/.codex/plugins/cache`.

Use it for durable behavior feedback such as:

- A Skill missed a recurring risk.
- A Skill asked unnecessary questions when repo evidence was enough.
- A Skill produced friction that should be avoided next time.
- A repeated positive behavior should be preserved.

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
failures, propose source-first diffs only, verify changed paths are limited to
allowed write surfaces, and stop for human review. Do not install, copy, commit,
push, or edit plugin/cache/global Skill directories.
```

The review should:

- Read `AGENTS.md`.
- Inspect `feedback/skills/events/*.jsonl` and
  `feedback/skills/summaries/*.md`.
- Group feedback by Skill and repeated failure type.
- Treat one-off low or medium feedback as a summary note.
- Require repeated similar evidence, ideally two or three events, before
  changing permanent Skill instructions.
- Allow a single high or critical event to justify a Skill edit only when the
  expected behavior is clear and likely to recur.
- Write rollups under `feedback/skills/summaries/`.
- Edit only repo-owned Skill source when evidence justifies it.
- Run `git diff --name-only` before the final response and verify changed paths
  stay within the allowed surfaces.

## Allowed Write Surfaces

For `$improve-skills`, allowed write surfaces are:

```text
skills/<target>/SKILL.md
skills/<target>/REFERENCE.md
feedback/skills/summaries/*.md
docs/skill-feedback-loop.md
README.md, only when needed to link canonical docs
```

Feedback event files may be added when the user explicitly asks to record new
feedback.

Forbidden surfaces include:

```text
%USERPROFILE%\.agents\skills
$HOME/.agents/skills
~/.codex/plugins/cache
bin/
backups/
unrelated docs or code
```

## Weekly Automation Prompt

When creating a recurring Codex automation, use the automation tool rather than
raw automation text. The prompt body should be:

```text
Weekly Skill feedback review:

In E:\dev\agents-toolkit, use $improve-skills to review feedback/skills.
If $improve-skills is not exposed as an active Skill, read
skills/improve-skills/SKILL.md directly and follow it. Find repeated Skill
failures, propose source-first diffs only, verify changed paths are limited to
allowed write surfaces, and stop for human review. Do not install, copy, commit,
push, or edit plugin/cache/global Skill directories.
```

The automation should stop after proposing source diffs and summaries. It should
not apply those changes to active installed Skills until the user explicitly asks
for install or repair work.

## Review Output

Each outer-loop run should report:

- Event files reviewed.
- Repeated failures found and supporting event counts.
- Skill source files changed, or why no edit was justified.
- Summary files written.
- Validation commands run.
- Deferred one-off notes that need more evidence.
