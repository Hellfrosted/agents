# Worktree Thread Plan

Default fan-out model for repeatable loops that explore or edit multiple candidates.

Scope boundary: this applies only while a loop is being designed or run. Later
`subagent`, `sub-agent`, or custom agent requests use normal subagent tooling
unless the user explicitly invokes `$architecture-cascade`.

## Default

Use one coordinator plus one isolated worker thread per independent candidate. Create manual Git worktrees nested under the saved project checkout and start local project threads under that same saved project. Do not create separate Codex GUI projects unless the user asks.

Use subagents only when explicitly requested or for short read-only review. Keep this plan internal; user-facing output describes results, blockers, verification, and integration choices.

## Required Plan

- **Role**: what candidate or work item the thread owns.
- **Ownership**: exact files, directories, tests, or docs it may edit.
- **Write intent**: `none`, `artifact-only`, or `code-editing`.
- **Starting state**: verified existing git ref, or current working tree only when replaying uncommitted coordinator changes is required.
- **Project scope**: the current saved project id the worker must be created under.
- **Location strategy**: `same-project-nested-manual` unless the user explicitly asks otherwise.
- **Worker cwd**: nested manual worktree path when the thread target remains the saved project.
- **Output contract**: candidate status, patch summary, tests run, blockers, merge notes.
- **Integration rule**: report only, ask before merge, or wait for explicit integration.

## Same-Project Nested Manual Worktree Mode

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
root="$repo_root/.codex-worktrees/architecture-cascade"
run_id="$(date +%Y%m%d%H%M%S)"
candidate_path="$root/$run_id/<candidate-id>"
```

Add `.codex-worktrees/` to `.git/info/exclude` before creating worktrees. Keep
this local-only; do not edit `.gitignore`.

Helper:

```bash
python3 "$PLUGIN_ROOT/scripts/wsl_worktree_plan.py" --repo . --candidate c1 --candidate c2
```

Add `--create` only after reviewing paths. Create each checkout with a verified existing ref:

```bash
mkdir -p "$(dirname "$candidate_path")"
git worktree add --detach "$candidate_path" "$verified_ref"
```

Create the candidate thread under the current saved project id, and put the
nested worktree path in the worker prompt as the cwd/workdir:

```text
target: { type: "project", projectId: "<current saved project id>", environment: { type: "local" } }
worker_cwd: "<candidate_path>"
```

Ignore `CODEX_WSL_WORKTREE_ROOT` for architecture cascade. Do not use `--root`.
External roots create separate project/trust boundaries in the Codex GUI.

Use detached worktrees by default. Create named branches inside a candidate worktree only after the worker decides the work is auto-safe and the repo permits it.

## Starting State Guardrails

- `startingState: { type: "branch", branchName: "<ref>" }` means "start from this existing git ref." It does not create `<ref>`.
- Verify a ref before passing it as `branchName`:

```bash
git rev-parse --verify --quiet "<ref>^{commit}"
```

- Use the current branch only after `git branch --show-current` and verification. If missing or detached, use a verified known ref such as `main`, or stop and ask.
- Do not pass new names like `codex/my-candidate-branch` to `branchName`.
- Avoid `startingState: { type: "working-tree" }` for fan-out; it replays the coordinator checkout's uncommitted diff into every worktree.
- Use `working-tree` only when the loop deliberately needs uncommitted coordinator changes; inspect `git status --short`, state the replay requirement, and create one pilot Worktree before a batch.
- If the loop wants named branches, let the worker create or switch branches after Worktree setup succeeds and repo workflow allows it.

## Safety Defaults

- Do not merge or combine worktree changes automatically unless the user explicitly asks.
- Ask before creating a thread when ownership overlaps active user work or another candidate thread.
- Candidate threads explore first, then edit only when the work is bounded and testable.
- The coordinator keeps the queue and final status; workers own their isolated work.
- For a batch, preflight the starting ref once and reuse it. Do not retry failed setup by switching to `working-tree` unless uncommitted changes are required.
- Default to nested same-project manual worktrees. Do not use Codex-managed worktrees or external worktree roots for fan-out unless the user explicitly chooses a separate project.
