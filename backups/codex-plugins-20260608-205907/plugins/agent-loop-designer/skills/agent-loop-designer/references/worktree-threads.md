# Worktree Thread Plan

Use worktree-backed threads as the default fan-out model for repeatable loops that explore or edit multiple candidates.

## Default

Prefer a coordinator thread plus one isolated worktree-backed thread per independent candidate. In WSL, create manual Linux-native Git worktrees and start local project threads there. Outside WSL, Codex-managed Worktree threads are usually fine.

Use subagents only when the user explicitly asks for subagents, or when the task is a short read-only review that should stay inside the current thread.

This plan is internal. User-facing loop output should describe results, changed artifacts, blockers, verification, and integration choices. Do not present the Worktree-thread topology as the final output unless the user asks how the loop works.

## Required Plan

For every worktree thread, specify:

- **Role**: what candidate or work item the thread owns.
- **Ownership**: exact files, directories, tests, or docs it may edit.
- **Write intent**: `none`, `artifact-only`, or `code-editing`.
- **Starting state**: verified existing git ref, or current working tree only when replaying uncommitted coordinator changes is explicitly required.
- **Location strategy**: `wsl-manual` or `codex-managed`.
- **Output contract**: candidate status, patch summary, tests run, blockers, merge notes.
- **Integration rule**: whether the coordinator only reports, asks before merge, or waits for explicit integration.

## WSL Manual Worktree Mode

Use this mode when the agent is running in WSL and a managed app worktree would
land on Windows storage such as `/mnt/c/Users/.../.codex/worktrees`.

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
root="${CODEX_WSL_WORKTREE_ROOT:-$HOME/codex-worktrees}"
repo="$(basename "$(git rev-parse --show-toplevel)")"
run_id="$(date +%Y%m%d%H%M%S)"
candidate_path="$root/$repo/$run_id/<candidate-id>"
```

When available, use the bundled helper to generate the same plan:

```bash
python3 "$PLUGIN_ROOT/scripts/wsl_worktree_plan.py" --repo . --candidate c1 --candidate c2
```

Add `--create` only after reviewing the planned paths.

Create each candidate checkout with a verified existing ref:

```bash
mkdir -p "$(dirname "$candidate_path")"
git worktree add --detach "$candidate_path" "$verified_ref"
```

Then create the candidate thread as a local project thread whose project id or
workspace root is the manual worktree path:

```text
target: { type: "project", projectId: "<candidate_path>", environment: { type: "local" } }
```

If the thread tool rejects that path because it is not a saved project, do not
fall back to Codex-managed worktrees. Report the created WSL paths and ask the
user to add/open them as projects, or ask for explicit approval to use managed
worktrees on Windows storage.

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
- For a batch, preflight the starting ref once and reuse the verified ref for every candidate. Do not retry failed setup by switching to `working-tree` unless uncommitted changes are required.
- In WSL, do not use Codex-managed worktrees on Windows storage for fan-out unless the user explicitly chooses that slower placement.
