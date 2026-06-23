# sk-up Go Implementation Reference

This document describes the promoted Go implementation of the skills updater.
It is no longer an in-progress port brief. Use it for the current CLI contract,
implementation boundaries, validation gates, and the consolidated rationale for
the historical port decisions.

For day-to-day operator usage, start with
[skills-updater.md](skills-updater.md). For source-first maintenance and common
checks, use [maintainer-guide.md](maintainer-guide.md).

## Current Outcome

`sk-up` and `skills-updates` are two entry point names for the same Go updater.
The implementation lives under `cmd/sk-up` and `internal/skup/...`. Windows
wrappers invoke an adjacent `sk-up.exe`; there is no supported PowerShell
updater fallback and Windows release archives intentionally do not ship
`skills-updates.exe`.

The updater maintains globally installed Codex skills under the universal
`.agents/skills` home. It can list installed skills, compare lockfile-backed
skills with upstream Git content, show diffs, open diffs, install updates,
install explicit sources, skip or unskip an upstream revision, list skips, and
remove installed skills.

## Source Material

- [CONTEXT.md](../CONTEXT.md): shared CLI vocabulary such as dry run,
  structured stdout, status names, and upstream hash skip.
- [docs/skills-updater.md](skills-updater.md): operator-facing behavior and
  examples.
- [docs/maintainer-guide.md](maintainer-guide.md): source-first workflow and
  verification guidance.
- [docs/adr/0001-port-sk-up-to-go.md](adr/0001-port-sk-up-to-go.md): why Go is
  the implementation language.
- [docs/adr/0002-break-sk-up-before-promotion.md](adr/0002-break-sk-up-before-promotion.md):
  compatibility decisions made before the Go updater became the contract.
- [docs/adr/0003-one-binary-two-skill-updater-names.md](adr/0003-one-binary-two-skill-updater-names.md):
  one binary with `sk-up` and `skills-updates` names.
- [docs/adr/0004-cross-platform-sk-up-state-paths.md](adr/0004-cross-platform-sk-up-state-paths.md):
  OS-native cache and state path rules.
- [docs/adr/0005-delegate-skill-installation.md](adr/0005-delegate-skill-installation.md):
  Skills CLI delegation and runner fallback.
- [docs/adr/0006-use-advisory-lock-files.md](adr/0006-use-advisory-lock-files.md):
  lockfile safety.
- [docs/adr/0007-require-git-not-tar.md](adr/0007-require-git-not-tar.md):
  external tool boundary.
- [docs/adr/0008-no-sk-up-config-file.md](adr/0008-no-sk-up-config-file.md):
  configuration surface.
- [docs/adr/0009-ship-sk-up-release-archives.md](adr/0009-ship-sk-up-release-archives.md):
  release archive packaging.
- [docs/adr/0010-use-private-go-sk-up-packages.md](adr/0010-use-private-go-sk-up-packages.md):
  Go source layout.
- [docs/adr/0011-validate-sk-up-before-promotion.md](adr/0011-validate-sk-up-before-promotion.md):
  validation bar that remains the regression floor.
- [docs/adr/0012-retire-powershell-sk-up.md](adr/0012-retire-powershell-sk-up.md):
  PowerShell retirement.

## Command Contract

Short form:

```text
sk-up -h
sk-up -l
sk-up -g
sk-up -d <skill>
sk-up -z [skill...]
sk-up -i [skill...]
sk-up -I <source...>
sk-up -s <skill>
sk-up -u <skill>
sk-up -S
sk-up -r <skill...>
```

Long form:

```text
skills-updates --help
skills-updates --list
skills-updates --global
skills-updates --diff <skill>
skills-updates --zed [skill...]
skills-updates --install [skill...]
skills-updates --install-source <source...>
skills-updates --skip <skill>
skills-updates --unskip <skill>
skills-updates --skips
skills-updates --remove <skill...>
```

Behavior map:

| Short form | Long form | Behavior |
| --- | --- | --- |
| `sk-up -l` | `skills-updates --list` | List installed global skill directories without Git. |
| `sk-up -g` | `skills-updates --global` | Compare lockfile-backed skills with upstream content. |
| `sk-up -d <skill>` | `skills-updates --diff <skill>` | Print a terminal diff for one changed skill. |
| `sk-up -z [skill...]` | `skills-updates --zed [skill...]` | Open changed skills in the configured diff tool. |
| `sk-up -i` | `skills-updates --install` | Install all changed or missing unskipped skills. |
| `sk-up -i <skill...>` | `skills-updates --install <skill...>` | Install named lockfile skills. |
| `sk-up -I <source...>` | `skills-updates --install-source <source...>` | Install package or repository sources explicitly. |
| `sk-up -s <skill>` | `skills-updates --skip <skill>` | Save a skip for the current upstream tree hash. |
| `sk-up -u <skill>` | `skills-updates --unskip <skill>` | Remove a saved skip. |
| `sk-up -S` | `skills-updates --skips` | List saved skips. |
| `sk-up -r <skill...>` | `skills-updates --remove <skill...>` | Remove named global skills. |

`--gui` is an alias for `--zed`; `--uninstall` is an alias for `--remove`;
`--install-all` and `install-all` are explicit forms of targetless install.

Source installs are explicit. Use `-I` or `--install-source` for source URLs,
SSH remotes, `.git` URLs, and `owner/repo` package shorthand. Named installs
use `-i <skill...>` and resolve the source from lockfile metadata.

## Flags And Environment

| Flag | Environment | Purpose |
| --- | --- | --- |
| `--agents-home <path>` | `SK_UP_AGENTS_HOME` | Override the global skills home. |
| `--cache-dir <path>` | `SK_UP_CACHE_DIR` | Override repository cache location. |
| `--state-dir <path>` | `SK_UP_STATE_DIR` | Override updater state location. |
| `--skills-command <cmd>` | `SK_UP_SKILLS_COMMAND` | Override delegated Skills CLI runner. |
| `--diff-tool <cmd>` | `SK_UP_DIFF_TOOL` | Override open-diff tool. |
| `--color auto|always|never` | `SK_UP_COLOR` | Control human output color. |
| `--json` | none | Emit one final JSON result on stdout. |
| `--jsonl` | none | Emit newline-delimited JSON events on stdout. |
| `--dry-run` | none | Plan and validate mutating work without changing state. |
| `--no-color` | none | Alias for `--color never`. |

There is no config file. Precedence is CLI flag, environment variable, then
platform default.

## State Paths

Skill install state remains under `AGENTS_HOME`:

```text
<AGENTS_HOME>/.skill-lock.json
<AGENTS_HOME>/skills
```

Fallback `AGENTS_HOME` is `$HOME/.agents` on Unix-like systems and
`%USERPROFILE%\.agents` on Windows.

Updater repository cache and skip state use OS-native cache/state directories:

```text
Linux cache: $XDG_CACHE_HOME/sk-up or $HOME/.cache/sk-up
Linux state: $XDG_STATE_HOME/sk-up or $HOME/.local/state/sk-up
macOS cache: $HOME/Library/Caches/sk-up
macOS state: $HOME/Library/Application Support/sk-up
Windows cache/state: %LOCALAPPDATA%\sk-up\cache and %LOCALAPPDATA%\sk-up\state
```

Use `--cache-dir`, `--state-dir`, `SK_UP_CACHE_DIR`, and `SK_UP_STATE_DIR` for
explicit overrides.

## Output Protocol

Human output is for operators. Structured output is the composable contract.
When `--json` or `--jsonl` is active, stdout contains only structured data.
Human progress, diagnostics, delegated runner output, warnings, and errors go
to stderr.

Stable result statuses are:

```text
ok
update
missing
skipped
error
```

Progress and operation words such as `fetch`, `clone`, `compare`, `install`,
`remove`, `repair`, and `plan` are event or action names, not result statuses.

Exit codes:

| Code | Meaning |
| --- | --- |
| `0` | Command completed successfully, including when updates are available. |
| `1` | Requested action failed or one or more targets errored. |
| `2` | Usage, configuration, dependency, invalid target, or invalid state path error. |
| `3` | Lock acquisition timeout. |
| `4` | Interrupted state was detected or repaired but the command cannot safely continue. |

## Implementation Layout

```text
cmd/sk-up/main.go
internal/skup/app
internal/skup/cli
internal/skup/compare
internal/skup/config
internal/skup/inventory
internal/skup/lockfile
internal/skup/output
internal/skup/plan
internal/skup/runner
internal/skup/state
internal/skup/status
```

The public contract is the CLI protocol, not an importable Go API. `cmd/sk-up`
is intentionally thin; it resolves the effective entry point name and calls the
app package. `internal/skup/cli` owns parsing and help text. `internal/skup/app`
routes commands and coordinates lockfile transactions. The compare, status,
state, runner, output, and config packages own the boundaries named in their
package paths.

## Core Behavior

Content comparison compares installed skill directory content with clean
upstream repository content and ignores CRLF line-ending differences. Lockfile
hashes alone are not authoritative update signals.

The updater requires external `git` for upstream comparison commands but does
not require external `tar`. Commands that do not need upstream comparison still
work without Git, including list, remove, unskip, skips, and named installs
that can delegate from lockfile source metadata.

Install and remove delegate to the upstream Skills CLI. If no runner override
is set, runner resolution tries:

1. `pnpm dlx skills@latest`
2. `bunx skills@latest`
3. `deno run -A npm:skills@latest`
4. `npx -y skills@latest`

Runner execution uses tokenized process execution, not shell evaluation by
default.

Lockfile writes use advisory lock files next to `.skill-lock.json`. The updater
preserves unknown lockfile fields and unrelated skill entries across
install/remove transactions.

Skip entries are keyed by upstream tree hash. A skip hides one exact upstream
change until upstream changes again.

Dry-run mode plans and validates mutating operations without changing files,
lock state, skip state, installed skills, or external package-manager state.

Remove operations perform updater-owned cleanup after successful delegated
remove: installed skill directory, saved skip state, and lockfile entry.

## Release And Wrapper Shape

`bin/build-sk-up-release.sh` builds native archives for Linux, macOS, and
Windows on amd64 and arm64. The generated output belongs under `dist/sk-up/`
and is not committed.

Unix archives include the `sk-up` binary and a `skills-updates` hardlink inside
the archive work directory. Windows archives include:

```text
sk-up.exe
sk-up.cmd
skills-updates.cmd
```

They intentionally exclude `skills-updates.exe` to avoid Windows UAC heuristics
for updater-looking executable names. `skills-updates.cmd` sets
`SK_UP_ENTRYPOINT=skills-updates` and then invokes the adjacent `sk-up.exe`.

## Validation Matrix

Use the smallest relevant check for routine edits and the broader matrix for
promotion-shape or command-contract work.

```bash
go test ./...
go test -race -shuffle=on -count=1 ./cmd/sk-up ./internal/skup/...
go vet ./...
SK_UP_DIST_DIR="$(mktemp -d)" SK_UP_VERSION=test bin/build-sk-up-release.sh
bash tests/sk-up-promotion-check.sh
git diff --check
```

Focused help checks:

```bash
go run ./cmd/sk-up --help
SK_UP_ENTRYPOINT=skills-updates go run ./cmd/sk-up --help
```

Windows wrapper regression:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests\skills-updates-install.ps1
```

From WSL, use a Windows-capable shell surface when available. Do not treat a
failed WSL interop launch as proof that the Windows wrapper is broken; record
the capability gap and run the PowerShell check from a real Windows shell.

Do not install into the active workstation, publish releases, push branches, or
force-update remote state without explicit user approval.

## Consolidated Completion Audit

The Go updater is promoted in source. Current proof is carried by source files,
tests, and the promotion checks rather than by the old step-by-step port
ledger:

| Requirement | Current evidence |
| --- | --- |
| Promoted Go implementation | `cmd/sk-up`, `internal/skup/...`, `go.mod`, `go test ./...`. |
| Short and long entry points | `internal/skup/cli` help tests, `SK_UP_ENTRYPOINT=skills-updates`, `bin/sk-up.cmd`, `bin/skills-updates.cmd`. |
| No PowerShell updater fallback | `bin/skills-updates.ps1` is absent; `tests/sk-up-promotion-check.sh` fails if wrapper references return. |
| Structured stdout and diagnostics on stderr | `internal/skup/output` tests and app JSON/JSONL dry-run tests. |
| Dry-run mutation planning | `internal/skup/plan` and app dry-run tests. |
| OS-native cache/state paths | `internal/skup/config` tests. |
| Advisory locks and lockfile preservation | `internal/skup/lockfile` tests and app install/remove transaction tests. |
| Skills CLI delegation | `internal/skup/runner` tests and app install/remove fake runner tests. |
| Git-required comparison without external `tar` | `internal/skup/compare` tests and release/promotion checks. |
| Upstream-hash skips | `internal/skup/state` and app skip/unskip/skips tests. |
| Windows wrapper/package shape | `tests/sk-up-promotion-check.sh` and `tests/skills-updates-install.ps1`. |

The historical development ledger was intentionally consolidated here. Git
history and ADRs remain the detailed source for why decisions were made; this
file should stay focused on what is true now and how to verify it.

## Open Gates

Release publishing, active workstation install repair, and pushing branches are
outside ordinary source maintenance. Do those only when the user explicitly
asks for that operation and target.
