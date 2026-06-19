---
name: architecture-cascade
description: Runs the improve-codebase-architecture process once, extracts all architecture candidates, starts separate worktree-backed candidate threads, and acts on every safe candidate. In WSL, defaults manual worktrees under the saved project checkout so candidate threads stay in the same Codex project group. Use when the user wants architecture improvement candidates handled automatically instead of picking one candidate manually.
---

# Architecture Cascade

## Use

The user should only need:

```text
$architecture-cascade <optional repo area or problem>
```

Also use this when `$agent-loop-designer` is asked to run `improve-codebase-architecture` once and then automatically explore or act on all candidates.

## Required Docs

Before spawning workers, read:

- `improve-codebase-architecture` skill instructions.
- `agent-loop-designer/references/subagents.md`, especially the boundary between private subagents and visible exploration threads.
- Codex Worktree/thread docs or local app thread tooling docs.
- `agent-loop-designer/references/worktree-threads.md`, especially WSL Same-Project Manual Worktree Mode.
- `agent-loop-designer/references/grit-coordination.md` when candidates can be owned by symbols in a Grit-supported language.
- Repo `AGENTS.md`, `CONTEXT.md`, `docs/adr/`, README/docs, and test/build guidance.

## Loop

1. Start a visible `architecture_explorer` thread to run one discovery pass using the `improve-codebase-architecture` process: domain vocabulary, ADRs, relevant docs, source, tests, and an HTML report. The exploration thread returns a compact candidate packet and report path. Apply the Report Location Guardrails before opening or sharing the report path.
2. Extract every candidate from the exploration thread packet into a queue with id, title, files, problem, solution, benefits, recommendation strength, ADR conflicts, likely tests, and proposed owner. Ask the explorer thread for one missing fact at a time instead of loading the full discovery context into the coordinator.
3. Create one worktree-backed candidate thread per candidate, capped to a practical batch size. When candidate ownership is symbol-addressable and Grit is available or can be installed for the requested workflow, initialize Grit if needed and use Grit claims/worktrees for the candidate workers. Each thread receives the candidate packet, repo guidance, docs checked, and an instruction to first verify actability before editing. Before creating any candidate thread, apply the Worktree Placement and Starting State Guardrails below.
4. Each candidate thread works in its own worktree. It maps implementation ownership, checks ADR/domain conflicts, identifies tests, and returns `actability: auto-safe | blocked | reject`.
5. Candidate threads for `auto-safe` candidates implement inside their own worktrees, run the smallest relevant verification, and report changed files, tests, residual risks, and any merge concerns.
6. The coordinator collects thread results. It may read or message candidate threads for missing integration context, merge concerns, verification status, or blocker details. If a candidate fails verification, stop that candidate, keep other independent worktree threads moving, and report the failure with evidence.
7. If the cascade itself fails because this skill missed a guardrail, encoded a bad tool assumption, created fragile output, or caused repeatable Worktree/thread setup failure, apply the Self-Improvement Rule below before the final response.
8. Report every candidate as implemented, rejected, or blocked. Ask the user only for blocked candidates that need a domain decision, ADR reopening, unsafe permission change, or overlapping ownership resolution.

## User-Facing Output

Return an architecture action report, not the execution topology:

- Implemented candidates: candidate id, summary, changed files, checks run, worktree/thread link or title.
- Rejected candidates: candidate id, evidence, and reason not to act.
- Blocked candidates: candidate id, exact user decision needed, and why automation stopped.
- Integration choices: completed worktrees ready for review or merge, with risks.
- Verification summary: checks passed, failed, or not run.

## Report Location Guardrails

Architecture cascade reports are user-facing artifacts. Do not rely on WSL `/tmp`
when the user may reopen the report after the current turn.

- The underlying architecture discovery may write the first report to the OS temp directory. For cascade runs, copy or rewrite the final report to a durable non-repo path before opening or reporting it.
- In WSL, `/tmp` can disappear after a distro restart, session cleanup, or temp cleanup. Browser URLs through `file://wsl$/...` can also be unreliable.
- If the repo is under a mounted Windows drive such as `/mnt/e/dev/...`, prefer a sibling durable temp folder on that drive, such as `/mnt/e/tmp/architecture-review-<timestamp>.html`.
- When giving the user a Windows browser URL for a mounted-drive report, use the Windows-style file URL, such as `file:///E:/tmp/architecture-review-<timestamp>.html`.
- If no durable mounted-drive path is available and OS temp must be used, tell the user that the report is volatile.

## Worktree Thread Contracts

## Worktree Placement Guardrails

Use WSL Same-Project Manual Worktree Mode when the coordinator is running in WSL
and the current repo is a saved Codex project. This keeps worker threads under
the same Codex project group without asking the app to create managed Worktrees.

- Do not ask the Codex app to create managed Worktree threads in that case. Managed app worktrees are created under `$CODEX_HOME/worktrees`, which can land on Windows storage and cause slow WSL I/O.
- Use `scripts/wsl_worktree_plan.py` to plan or create the manual worktrees when available.
- Create manual Git worktrees under the saved project checkout by default:
  `<repo>/.codex-worktrees/architecture-cascade/<run-id>/<candidate-id>`.
- Before creating nested worktrees, add `.codex-worktrees/` to
  `.git/info/exclude` if it is not already excluded. Keep this local-only; do
  not edit `.gitignore` for agent worktrees.
- Use `git worktree add --detach <manual-path> <verified-ref>` unless a branch is intentionally required by the repo workflow.
- Start candidate threads as local project threads using the current saved
  project id, with worker prompts explicitly naming the nested worktree path as
  the cwd/workdir for all reads, edits, git status, and validation.
- If `CODEX_WSL_WORKTREE_ROOT` or `--root` is set, treat that as an explicit
  external-root override. In that mode, create manual worktrees under that root
  and start candidate threads pointed at the manual worktree path. If the thread
  API rejects the manual path because it is not a saved project, stop after
  creating the worktrees and ask the user to add/open those paths as projects.
  Do not fall back to Windows-backed managed worktrees without explicit user
  approval.
- Use Codex-managed Worktree threads only when not in WSL or when the user
  explicitly chooses managed worktrees.

## Self-Improvement Rule

Architecture cascade is a self-improving skill. When the coordinator failure is
repeatable and caused by this skill's instructions, update the source plugin and
active cached plugin before closing the turn.

- Patch the narrowest relevant skill or reference file under the source plugin and mirror the same change into the active cache.
- Base the update on observed evidence: command output, app error text, rejected API shape, missing durable artifact, or thread setup behavior.
- Add guardrails, preflights, or output-location rules. Do not add broad apologies, incident logs, generated-by notes, or user-private raw logs.
- Validate source/cache parity with `diff`, search for stale wording with `rg`, and run the smallest relevant bundled script check.
- Store a concise privacy-safe memory for durable failures after the patch succeeds.
- Do not update the plugin for one-off repo failures, user-cancelled work, unavailable external services, or failures that belong in the target repo's own docs.

## Worktree Starting State Guardrails

Both manual Git worktrees and Codex `create_thread` Worktree setup are strict about
their starting refs.

- `startingState: { type: "branch", branchName: "<ref>" }` starts from an existing git ref. It does not create a new branch.
- Never pass a new candidate branch name such as `codex/arch-cascade-c1-...` as `branchName` unless `git rev-parse --verify --quiet <ref>^{commit}` confirms it already exists.
- For architecture-cascade fan-out, prefer the current existing branch ref, usually the output of `git branch --show-current`, after verifying it with `git rev-parse --verify --quiet`.
- If the current branch cannot be resolved, use another verified existing ref such as `main` only after confirming it exists. If no suitable ref exists, stop and ask instead of queuing Worktree threads.
- Do not use `startingState: { type: "working-tree" }` for architecture-cascade fan-out. It replays the coordinator checkout's uncommitted diff into every Worktree and can create failed pending chats when patch replay skips or conflicts.
- If uncommitted coordinator changes are intentionally required by a future loop, first state that requirement, inspect `git status --short`, and create one pilot Worktree thread before queuing a batch.
- Do not rely on Worktree setup to name the eventual implementation branch. Candidate threads may create or switch branches inside their own worktree only after setup succeeds and the repo workflow permits it.
- In WSL Same-Project Manual Worktree Mode, use the same verified ref in
  `git worktree add --detach`; do not use app `startingState`.

Minimum safe preflight before a batch:

```bash
git status --short
git branch --show-current
git rev-parse --verify --quiet "$(git branch --show-current)^{commit}"
```

When the checkout is dirty, decide explicitly whether those changes are relevant. If they are unrelated, still start from the verified branch ref, not from `working-tree`.

Candidate worktree threads:

- Role: candidate owner.
- Worktree: one isolated manual WSL git worktree or Codex-managed Worktree thread per candidate, chosen by the Worktree Placement Guardrails. In same-project manual mode, the thread target is the saved project and the worker cwd is the nested worktree path.
- Grit: for symbol-addressable candidates, claim the candidate symbols before editing and use the `.grit/worktrees/<agent-id>/` cwd. The coordinator runs `grit done` only when commits and integration are allowed by the current task.
- Ownership: exact candidate files, directly related tests, and candidate-specific docs.
- Write intent: explore first, then code-edit only when `actability` is `auto-safe`.
- Output: candidate id, actability, evidence, changed files, tests run, blockers, merge concerns.

Coordinator thread:

- Owns worker orchestration, queue extraction from the architecture explorer packet, thread creation, status collection, and final summary.
- Delegates architecture discovery to a visible exploration thread and avoids absorbing full discovery context unless a blocker cannot be resolved from compact packets or follow-up questions.
- Does not edit candidate code directly while candidate worktree threads are active.
- Does not merge worktree results automatically unless the user explicitly asks for merge/integration.

## Act Automatically Means

Act on every candidate that is bounded, non-conflicting, verified by its candidate thread, and permitted by the current session. Do not wait for the user to pick one.

Do not act automatically when the candidate needs a product/domain decision, contradicts an ADR without clear authorization, overlaps another active worktree thread's ownership, requires unavailable permissions, or lacks a credible verification path.
