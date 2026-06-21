# Grit Coordination

Use Grit as the coordination backend when a loop has parallel code-editing workers whose ownership can be expressed as symbols, functions, classes, modules, or disjoint file areas in a Grit-supported language.

Do not use Grit for read-only loops, docs-only loops, generated-file work, migrations, broad formatting, or work whose ownership cannot be claimed before editing.

## Preflight

Run docs-first before relying on Grit. Check local help or upstream docs for the installed command surface, then record it in `docs_checked`.

Minimum preflight in the target repo:

```bash
grit --version
git rev-parse --show-toplevel
git status --short
git branch --show-current
```

If `grit` is missing and the user asked for Grit-backed coordination or plugin
setup, report the install command and ask for explicit approval before running
networked dependency installation. The upstream README command may fail when
the repository has multiple packages; use the package selector after approval:

```bash
cargo install --git https://github.com/rtk-ai/grit grit
```

## Initialization

If the repo does not have usable Grit state, initialize it without asking:

```bash
grit init
grit config set-local
grit symbols
```

Use the local backend by default. Do not configure Azure, S3, R2, MinIO, or other remote lock stores unless the user explicitly asks for distributed team coordination and provides the required non-secret parameters through an approved channel. When a durable loop uses a non-local backend, include `remote_backend_authorization` in `grit_coordination` with the explicit user request and the approved non-secret parameter source.

Keep `.grit/` local. If the repo has a `.git` directory and `.grit/` is not already excluded, add it to `.git/info/exclude`, not `.gitignore`. If `grit init` creates a new untracked `.gitignore` whose only purpose is ignoring `.grit`, migrate that entry to `.git/info/exclude` and remove the generated `.gitignore`. If `.gitignore` already existed, do not rewrite it unless the repo intentionally tracks local agent-state ignores.

Re-run `grit init` when `grit claim` reports a missing symbol or the loop has changed the symbol surface enough that the index is stale.

## Claiming

Each worker must have an agent id, intent, ownership, and output contract before it edits.

Prefer explicit symbol claims:

```bash
grit claim -a <agent-id> -i "<intent>" <path>::<symbol> <path>::<symbol>
```

Use dependency-aware claims when the worker will edit a caller and needs read locks on callees:

```bash
grit claim -a <agent-id> -i "<intent>" --with-deps <path>::<symbol>
```

Use `grit assign` only when the loop intentionally lets Grit choose a free symbol from a bounded file area. Use `--queue` only when the loop is allowed to wait for contested ownership.

When a claim succeeds, the worker's cwd is the Grit worktree, usually:

```text
.grit/worktrees/<agent-id>/
```

All worker reads, edits, tests, and git status checks happen inside that Grit worktree unless the coordinator explicitly requests read-only context from another thread.

## Thread Context

The coordinator owns integration context. It may inspect or message worker threads before making integration decisions, especially before `grit done`, after verification failures, when ownership overlaps, or when a blocked/stale lock needs explanation.

When Codex thread tools are available, the coordinator can read worker threads or send short context requests. Ask only for the missing integration fact: changed files, tests run, blocker, merge concern, or whether the worker is still active. Do not use cross-thread messages to bypass Grit locks or to ask a worker to edit outside its claim.

Worker prompts should include:

- agent id
- claimed symbols or assigned file area
- Grit worktree cwd
- allowed files/tests/docs
- verification command
- output contract
- rule to stop and report before editing outside the claim

## Done And Integration

`grit done` auto-commits, rebases, merges, and releases locks. Treat it as an integration action, not as routine cleanup.

Run `grit done -a <agent-id>` only when the current loop is allowed to create commits and merge worker work into the coordinator checkout. If the user did not explicitly allow commits or integration, stop after worker verification and report the Grit worktree path, changed files, tests, and integration choice.

Before `grit done`:

```bash
git status --short
grit status
```

If the main checkout is dirty, decide whether the changes belong to the loop. Do not run `grit done` over unrelated user changes. Ask one narrow question only when the dirty state blocks integration and cannot be isolated.

After `grit done`:

```bash
git status --short
grit status
```

Confirm the lock released and the work reached the intended checkout. If `grit done` fails, keep the agent branch/worktree intact, collect the error, and report the recovery choice.

Do not run `grit session pr`, push, or open a PR unless the user explicitly asked for a push/PR workflow.

## Cleanup

At loop end, run `grit status`. Release only locks owned by this loop. Use `grit gc` for stale expired locks when the stale ownership is proven and not actively held by another worker.

For durable loops, include `grit_coordination` in the loop spec with:

- backend: `local`
- remote backend authorization, only when backend is not `local`
- init policy
- claim strategy
- done policy
- thread context rule
- cleanup rule
