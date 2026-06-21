---
name: improve-skills
description: Skill improvement loop for repo-owned agents-toolkit Skills. Use when the user asks to process Skill feedback, evaluate or tune Skill activation/behavior, red-team Skill triggers, inspect cross-Skill collisions, or edit repo-owned Skill source. Not for creating/installing Skills or drafting briefs about Skill work.
---

# Improve Skills

Review repo-owned Skill feedback and eval results, then make source-first Skill
changes when the user asks for edits. This Skill is an outer loop: it improves
Skill text from observed behavior, but does not install, copy, commit, push, or
schedule changes unless the user explicitly asks for that separate operation.

## Scope

Work in `/mnt/e/dev/agents-toolkit` unless the user gives a more specific path
inside this repo. Read `AGENTS.md` before editing.

Allowed write surfaces:

- `skills/<target>/SKILL.md`
- `skills/<target>/REFERENCE.md`
- `skills/<target>/agents/openai.yaml`, only when metadata would otherwise be
  stale after a Skill change
- `plugins/<plugin>/skills/**/SKILL.md`
- `plugins/<plugin>/skills/**/references/*.md`
- `plugins/<plugin>/references/**/*.md`
- `feedback/skills/README.md`
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
- Target Skill source under `skills/<target>/` and
  `plugins/<plugin>/skills/<target>/`
- Plugin-root references named by target Skills, such as
  `plugins/<plugin>/references/...`
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
8. When writes are in scope, update or create a concise summary in
   `feedback/skills/summaries/` that lists reviewed event files, repeated
   lessons, proposed edits, and deferred one-off notes. For read-only runs,
   report the proposed summary content instead.
9. If no Skill edit is justified, write the summary only when writes are in
   scope; otherwise report the proposed summary and explain why no edit is
   justified.

## Activation Eval Loop

Use this loop when the user asks to tune triggers, run evals, or check unwanted
effects even if there are no feedback event files:

1. Inventory every target `skills/<target>/SKILL.md` and
   `plugins/<plugin>/skills/<target>/SKILL.md`; read `agents/openai.yaml` and
   linked reference files when present.
2. For each target, write a trigger matrix in working notes: prompts that should
   trigger, prompts that should not trigger, likely over-trigger risks, likely
   under-trigger risks, intended behavior, and unwanted effects.
3. Launch blank-slate subagents for independent evals when available and in
   scope. If the current tool surface requires explicit user authorization,
   treat user requests for subagents, blank-slate evals, delegation, or parallel
   evals as authorization; otherwise ask before spawning. If the user asks for
   read-only/no-subagent work or subagents are unavailable, continue with the
   single-agent matrix and report the confidence limit. Use
   non-full-history forks (`fork_context: false` or `fork_turns: "none"`), keep
   prompts task-like, and ask for read-only findings. Do not pass suspected
   fixes or prior conclusions.
4. Cover at least these angles when scope is broad: activation precision,
   behavior adherence/safety gates, and cross-Skill collisions.
5. When edits are in scope, patch the smallest frontmatter description or body
   text that fixes an evidence-backed failure. Put all "when to use" trigger
   language in frontmatter, because the body only loads after activation. For
   read-only runs, report the proposed wording instead.
6. When subagents are in scope, re-run focused blank-slate evals for any Skill
   whose trigger or high-risk behavior changed. Continue until no P0/P1
   activation or safety issue remains, or until the remaining uncertainty,
   budget limit, or human gate is explicit.

Avoid false precision. A good description says exactly what should activate the
Skill, names the common synonyms users actually type, and also excludes nearby
workflows that should stay ordinary agent behavior or another Skill.

## Validation

Before the final response:

1. Compare changed paths against the pre-run dirty baseline. If no baseline was
   captured, use `git status --short` and the known task edits to separate
   pre-existing user work from this run's changes.
2. Verify every path changed by this run is in the allowed write surfaces above,
   except feedback event files the user explicitly asked to add. Report
   pre-existing out-of-scope dirty paths separately; do not treat them as this
   run's validation failure.
3. Re-read every Skill file you changed and confirm the frontmatter remains
   valid, self-contained, and neither too broad nor too narrow for the eval
   cases.
4. Run the smallest available validation for Skill metadata/frontmatter.
5. Stop before install, copy to active Skill paths, commit, push, or
   create/update automations unless explicitly asked.

## Output

Report:

- Feedback files reviewed.
- Eval passes run, including subagent angles and whether they were blank-slate.
- Repeated failures found and how many events support each one.
- Skill source files changed, or why no edit was justified.
- Summary file written.
- Validation command results.
- Any deferred one-off notes that need more evidence.
