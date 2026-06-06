---
name: tuck
description: Tucks completed local changes into focused, reviewable git commits with mandatory read-only per-file review subagents before staging. Use when the user explicitly invokes $tuck; do not use for ordinary, small, or natural-language commit requests.
---

# Tuck

Use only after exact `$tuck`. Ordinary commit requests use the normal git flow.

## Flow

1. Run `git status --short`.
2. Inspect diffs before staging; treat unrelated modified or untracked files as user work.
3. Split commits by reviewer-facing intent. Do not mix formatting churn, unrelated docs, or separable behavior.
4. Before staging each commit, spawn one read-only reviewer per file in that commit.
5. Fix valid blocking findings, then stage only intended paths or hunks.
6. Re-check `git diff --cached`, commit with a concise imperative message, and run the smallest relevant check.
7. After all commits, ask whether to push and where. Never push or force-push without explicit confirmation.

## Review

For each intended commit, confirm scope with `git diff -- <paths>`. For intentional new files, confirm with `git ls-files --others --exclude-standard -- <paths>` and review full content.

Use one subagent per file, ideally in parallel. Prefer `explorer` when available. If subagents are unavailable, state the blocker and continue only with user approval.

Prompt:

```text
You are a read-only reviewer for one file in an intended git commit.
Commit scope: {COMMIT_SCOPE}
Assigned file: {FILE}
Review only this file's changed diff, plus nearby context needed to understand it.
Look for correctness bugs, regressions, missing tests, security/privacy risk,
broken docs, or maintainability problems that should block this commit.
Do not edit files. Do not spawn subagents. Return only actionable findings with
file paths and line references, or state that this file has no blocking issues.
```

Do not stage a commit until all reviewers for that commit return. Re-run review when a material fix changes the diff. For docs-only commits, review accuracy and clarity.

## Git Safety

Do not stage untracked files by default. Use `git add <specific paths>` or non-interactive patch staging. Do not use destructive commands such as `git reset --hard` or `git checkout --`.

For push, run `git branch --show-current` and `git remote -v`. Recommend the current branch when it exists and is not `main`; otherwise recommend a new branch.
