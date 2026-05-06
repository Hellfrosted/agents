---
name: tuck
description: Tucks completed local changes into focused, reviewable git commits with mandatory read-only per-file review subagents before staging. Use when the user explicitly invokes $tuck; do not use for ordinary, small, or natural-language commit requests.
---

# Tuck

## Quick Start

Use this only after an exact `$tuck` invocation. For any ordinary commit request, use the normal git commit flow instead.

1. Run `git status --short`.
2. Inspect the diff before staging.
3. Split changes by reviewer-facing intent, not by implementation accident.
4. Before each commit, the main agent spawns one read-only review subagent per file in that commit.
5. Commit only related tracked files or intentional new files.
6. After committing, ask whether to push and where.

### 1. Protect Worktree

- Treat unrelated modified or untracked files as user work unless the user explicitly says otherwise.
- Do not stage untracked files by default.
- Do not use destructive commands such as `git reset --hard` or `git checkout --`.
- If a file mixes related and unrelated changes, inspect it and stage only the relevant hunks.

### 2. Choose Chunks

- Prefer shared primitive or helper extraction before feature usage.
- Prefer data/model/API contract changes before UI or integration changes.
- Prefer bug fix and regression test together when the test directly proves the fix.
- Prefer documentation-only updates separate from code unless they explain that exact code change.

- Avoid "all files I touched" as one commit when there are separable intents.
- Avoid formatting churn mixed with behavioral changes.
- Avoid a knowingly broken commit unless the user explicitly requested a work-in-progress stack.

### 3. Review Files

Before staging each commit chunk, spawn read-only review subagents for the exact files intended for that commit.

- Use `git diff -- <paths>` to confirm the tracked-file diff scope first.
- For intended new untracked files, use `git ls-files --others --exclude-standard -- <paths>` to confirm they are untracked, then review the whole file with `git diff --no-index /dev/null -- <file>` or equivalent full-file context before staging.
- Launch one separate subagent per changed file, ideally in parallel.
- Give each review subagent exactly one assigned file path; it reviews only that file's changed diff plus enough nearby context to understand it.
- Do not ask a review subagent to spawn more subagents. The main agent is responsible for fan-out.
- Reviews must not mutate files. Implement valid fixes yourself, then re-run the relevant review if the diff materially changes.
- For docs-only commits, still run per-file review, focused on accuracy, clarity, and whether docs match the code or workflow they describe.

Use this subagent prompt shape once per file, replacing `{FILE}` and `{COMMIT_SCOPE}`:

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

- Spawn review subagents with the normal subagent tool, not shell `codex exec`. Use the default inherited model unless the user explicitly requests another model or the task clearly requires one.
- Prefer `explorer` subagents when available because the task is read-only code review. If only generic subagents are available, still keep the prompt read-only and single-file.
- Wait for the review subagents before staging. When there are many files, wait in batches as needed, but do not stage a chunk until every assigned file in that chunk has returned.
- Fix valid blocking findings before committing.
- Ignore non-actionable or speculative findings, but mention any meaningful residual risk in the final commit summary.
- If subagents are unavailable, state the exact blocker and continue only if the user approves committing without the per-file review.

### 4. Stage

- Use `git add <specific paths>` for whole-file chunks.
- Use non-interactive patch staging when one file needs to be split.
- Re-check `git diff --cached` before every commit.
- Keep generated lockfiles, local docs, scratch files, and tool artifacts out unless they are intentional.

### 5. Commit

- Use concise imperative commit messages.
- Mention the primary user-facing or reviewer-facing change.
- Do not mention internal agent process.

### 6. Verify

- Run the smallest relevant checks before or after the final commit.
- If checks cannot run, state the exact blocker.
- If checks fail, fix or ask before committing the known failure.

## Push

After all requested commits are created, do not push automatically. Ask whether to push to the current branch, a new branch, or `main`. Recommend the current branch when it exists and is not `main`; recommend a new branch when currently on `main`. Only push to `main` after explicit confirmation.

- Run `git branch --show-current` and `git remote -v`.
- Use normal `git push`; use `gh` only for GitHub-specific actions such as opening a PR.
- Never force-push unless explicitly requested.
- If creating a new branch, ask for the branch name unless the user already gave one.
