# Architecture Cascade Example

This is the starter pattern for a loop similar to `improve-codebase-architecture`.

## Mission

Find codebase architecture deepening opportunities once, then explore and act on all candidates instead of stopping for a single user selection.

## Loop

- Trigger: `$architecture-cascade` or `$agent-loop-designer` followed by "run improve-codebase-architecture and act on all candidates in worktree-backed threads".
- State: `CONTEXT.md`, `docs/adr/`, README/docs, prior architecture reports if available.
- Inputs: source tree, tests, recent diffs, domain docs, build/test commands.
- Workers: main coordinator thread; one Codex Worktree thread per candidate.
- Artifact: durable non-repo HTML report plus an architecture candidate queue with status for every candidate.
- Decision point: only blocked candidates ask for user input.
- Next action: implement safe candidates, verify them, update `CONTEXT.md` or offer ADRs when decisions crystallize.
- Safety rule: candidate edits happen in isolated worktree-backed threads after each thread confirms files, tests, conflicts, and ownership. In WSL, use manual Linux-native git worktrees unless the user explicitly chooses managed worktrees.

## Automation Upgrade

Run weekly in a coordinator worktree. Create candidate worktree-backed threads for bounded candidates; in WSL, place manual git worktrees under `$HOME/codex-worktrees`. Report blocked candidates with the exact decision needed.
