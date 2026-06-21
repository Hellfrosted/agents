# Worktree Thread Plan

Use worktree-backed threads as the default visible workspace for exploration the user may read, interrupt, steer, or use for decisions, and as the default isolation model for repeatable loops that edit multiple candidates or need separate checkout state. Use subagents only for private/background context-heavy work where a compact packet is enough.

When parallel code-editing ownership can be expressed as Grit symbols, use `grit-coordination.md` first. In that mode, Grit creates `.grit/worktrees/<agent-id>/`, owns symbol locks, and controls `grit done` integration; this file still supplies the thread output contract and Codex thread context rules.

## Default

Prefer a coordinator thread plus visible exploration threads for human-steerable discovery/review and one isolated worktree-backed thread per independent code-editing candidate. Use subagents for private/background research that does not need human intervention. In WSL, create manual Git worktrees nested under the saved project checkout and start local project threads under that same saved project. Outside WSL, Codex-managed Worktree threads are usually fine.

Use thread messaging deliberately. The coordinator may read or message a Worktree thread to get missing integration context, merge concerns, verification status, or a blocker explanation. Ask for the smallest missing fact; do not import entire worker histories into the coordinator.

This plan is internal. User-facing loop output should describe results, changed artifacts, blockers, verification, and integration choices. Do not present the Worktree-thread topology as the final output unless the user asks how the loop works.

## Required Plan

For every worktree thread, specify:

- **Role**: what candidate or work item the thread owns.
- **Ownership**: exact files, directories, tests, or docs it may edit.
- **Write intent**: `none`, `artifact-only`, or `code-editing`.
- **Starting state**: verified existing git ref, or current working tree only when replaying uncommitted coordinator changes is explicitly required.
- **Project scope**: the current saved project id or workspace root the worker must be created under.
- **Location strategy**: `wsl-same-project-manual`, `wsl-external-manual`, or `codex-managed`.
- **Worker cwd**: nested manual worktree path when the thread target remains the saved project.
- **Output contract**: candidate status, patch summary, tests run, blockers, merge notes.
- **Integration rule**: whether the coordinator only reports, asks before merge, or waits for explicit integration.
- **Coordination backend**: `manual-worktree` or `grit`.
- **Visibility**: whether the thread is human-readable and may receive user intervention.

## WSL Same-Project Manual Worktree Mode

Use this mode when the agent is running in WSL and the current repo is a saved
Codex project. It is the default WSL fan-out mode.

Preflight:

```bash
test -n "${WSL_DISTRO_NAME:-}" || grep -qi microsoft /proc/version
git rev-parse --show-toplevel
git status --short
git branch --show-current
git rev-parse --verify --quiet "$(git branch --show-current)^{commit}"
```

Default root:

```bash
repo_root="$(git rev-parse --show-toplevel)"
root="${CODEX_WSL_WORKTREE_ROOT:-$repo_root/.codex-worktrees/architecture-cascade}"
run_id="$(date +%Y%m%d%H%M%S)"
candidate_path="$root/$run_id/<candidate-id>"
```

When using the default nested root, add `.codex-worktrees/` to
`.git/info/exclude` before creating worktrees. Keep this local-only; do not edit
`.gitignore` for agent worktrees.

When available, use the bundled helper to generate the same plan. Set `PLUGIN_ROOT`
to the agent-loop-designer plugin root first; from this reference directory, use
`PLUGIN_ROOT="$(cd ../../.. && pwd)"`.

```bash
python3 "$PLUGIN_ROOT/scripts/wsl_worktree_plan.py" --repo . --candidate c1 --candidate c2
```

Add `--create` only after reviewing the planned paths.

Create each candidate checkout with a verified existing ref:

```bash
mkdir -p "$(dirname "$candidate_path")"
git worktree add --detach "$candidate_path" "$verified_ref"
```

Then create the candidate thread as a local project thread under the current
saved project id, and put the nested worktree path in the worker prompt as the
cwd/workdir for all reads, edits, git status, and validation:

```text
target: { type: "project", projectId: "<current saved project id>", environment: { type: "local" } }
worker_cwd: "<candidate_path>"
```

If `CODEX_WSL_WORKTREE_ROOT` or `--root` is set, treat it as an explicit
external-root override. In that mode, create candidate threads with the manual
worktree path as the project id. If the thread tool rejects that path because it
is not a saved project, do not fall back to Codex-managed worktrees. Report the
created WSL paths and ask the user to add/open them as projects, or ask for
explicit approval to use managed worktrees on Windows storage.

Use detached worktrees by default. Create named branches inside a candidate
worktree only after the thread decides the work is auto-safe and the repo
workflow permits it.

## Starting State Guardrails

Treat Worktree setup state as an input ref, not a branch creation request.

- `startingState: { type: "branch", branchName: "<ref>" }` means "start from this existing git ref." It does not create `<ref>`.
- Verify a ref before passing it as `branchName`:

```bash
git rev-parse --verify --quiet "<ref>^{commit}"
```

- Use the current branch only after reading it with `git branch --show-current` and verifying it. If the current branch is missing or detached, use a verified known ref such as `main`, or stop and ask.
- Do not pass new names like `codex/my-candidate-branch` to `branchName` unless the branch already exists.
- Avoid `startingState: { type: "working-tree" }` for fan-out. It replays the coordinator checkout's uncommitted diff into every new worktree and can fail if the patch no longer applies.
- Use `working-tree` only when the loop deliberately needs uncommitted coordinator changes. In that case, inspect `git status --short`, mention the replay requirement, and create one pilot Worktree before queuing a batch.
- If the loop wants each candidate on a named branch, let the candidate thread create or switch branches after the Worktree is successfully created and after checking the repo workflow.

## Safety Defaults

- Create separate worktree threads for disjoint candidates.
- Do not merge or combine worktree changes automatically unless the user explicitly asks.
- Ask before creating a thread when ownership overlaps active user work or another candidate thread.
- Candidate threads explore first, then edit only when the work is bounded and testable.
- The coordinator keeps the queue and final status; candidate threads own their isolated work.
- The coordinator may read or message candidate threads to collect missing integration context before merging, running `grit done`, or reporting blockers.
- For a batch, preflight the starting ref once and reuse the verified ref for every candidate. Do not retry failed setup by switching to `working-tree` unless uncommitted changes are required.
- In WSL, default to nested same-project manual worktrees. Do not use Codex-managed worktrees on Windows storage for fan-out unless the user explicitly chooses that slower placement.
