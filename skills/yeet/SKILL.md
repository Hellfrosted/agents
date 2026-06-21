---
name: yeet
description: Yeet git shortcut for committing and pushing this turn's changes. Use only when the user invokes $yeet or uses "yeet" as the git commit-and-push action. Not for ordinary commit/push/PR/deploy/discard/delete/stash requests or reviewed/local commit workflows; use tuck first when review is requested.
---

# Yeet

Commit and push only the changes intended for this turn.

IMPORTANT: just do so for this turn. This doesn't mean you should commit and
push changes in future turns.

## Flow

1. Run `git status --short` and inspect the relevant diff before staging.
2. Stage only paths or hunks that belong to this turn. Ask before including
   unrelated modified files or untracked files.
3. Run the smallest relevant validation before committing when practical.
4. Commit with a concise imperative message.
5. Run `git branch --show-current` and `git remote -v`.
6. Push the current branch only when it is clear and not `main`. If the current
   branch is `main`, detached, has no upstream, has multiple plausible remotes,
   or the push target is ambiguous, ask before pushing.
7. Report the commit hash, branch, remote, and validation result.

Do not force-push, delete/discard/stash changes, create a PR, change branches,
or include secrets/credentials unless the user explicitly asks for that exact
action and the required safety checks pass.

If the user asks for review before commit, use Tuck first and push only through
the explicit push gate after the reviewed local commit exists.

Do not add yourself as a co-author.
