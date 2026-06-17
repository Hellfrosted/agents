---
name: improve-skills
description: Reviews agents-toolkit Skill feedback, identifies repeated failures, and proposes source-first diffs to repo-owned Skill files while stopping for human review. Use when asked to improve repo-owned Skills from feedback or run the Skill feedback loop.
---

# Improve Skills

Review repo-owned Skill feedback and propose source-first Skill changes. This
Skill is an outer loop: it improves Skill text from observed behavior, but does
not install, copy, commit, push, or schedule changes unless the user explicitly
asks for that separate operation.

## Scope

Work in `/mnt/e/dev/agents-toolkit` unless the user gives a more specific path
inside this repo. Read `AGENTS.md` before editing.

Allowed write surfaces:

- `skills/<target>/SKILL.md`
- `skills/<target>/REFERENCE.md`
- `feedback/skills/summaries/*.md`
- `docs/skill-feedback-loop.md`
- `README.md`, only when needed to link canonical docs

Forbidden write surfaces:

- Installed global Skill directories, including `%USERPROFILE%\.agents\skills`
  and `$HOME/.agents/skills`.
- `~/.codex/plugins/cache` and other plugin/cache directories.
- OMO plugin/cache files.
- Unrelated `bin/` runtime files.
- Unrelated project docs, code, backups, and generated artifacts.

## Inputs

Inspect:

- `feedback/skills/events/*.jsonl`
- `feedback/skills/summaries/*.md`
- Target Skill source under `skills/<target>/`
- Relevant repo docs only when needed to verify source-first workflow or
  existing terminology.

Never copy secrets, credentials, private personal data, raw logs, or raw chat
exports into Skill text. For events with `"private": true`, carry forward only
the safe lesson needed for future behavior.

## Review Process

1. Parse feedback events as append-only observations. If an event is malformed,
   record that in the summary and continue with valid events.
2. Group events by `skill`, then by repeated failure type or expected future
   behavior.
3. Treat `miss`, `near-miss`, and `friction` as improvement candidates. Treat
   `win` and `note` as context unless they clarify an existing candidate.
4. Require repeated similar evidence, ideally two or three events, before
   changing permanent Skill instructions.
5. A single `high` or `critical` severity event may justify a Skill edit when
   the expected behavior is clear, source-backed, and likely to recur.
6. Treat one-off low or medium feedback as a summary note, not a permanent Skill
   rule.
7. Prefer the smallest Skill text change that prevents the repeated failure.
   Do not add broad policy, speculative future cases, or duplicate global
   instructions.
8. Update or create a concise summary in `feedback/skills/summaries/` that
   lists reviewed event files, repeated lessons, proposed edits, and deferred
   one-off notes.
9. If no Skill edit is justified, only write the summary and explain why.

## Validation

Before the final response:

1. Run `git diff --name-only`.
2. Verify every changed path is in the allowed write surfaces above, except
   feedback event files the user explicitly asked to add.
3. Re-read every Skill file you changed and confirm the frontmatter remains
   valid and self-contained.
4. Stop for human review. Do not install, copy to active Skill paths, commit,
   push, or create/update automations unless explicitly asked.

## Output

Report:

- Feedback files reviewed.
- Repeated failures found and how many events support each one.
- Skill source files changed, or why no edit was justified.
- Summary file written.
- Validation command results.
- Any deferred one-off notes that need more evidence.
