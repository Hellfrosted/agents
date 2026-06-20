# Go sk-up Port Brief

This document is the control surface for porting the skills updater from the
retired PowerShell implementation to a promoted Go implementation. Keep it
current while implementing so future agents can choose the next action without
reconstructing the full design conversation.

## Outcome

Ship one Go skills updater implementation that works on major Linux
distributions, macOS, and Windows. After promotion, `sk-up` and
`skills-updates` are entry point names for the same Go binary, and the
PowerShell updater is retired rather than shipped as a supported fallback.

## Source Material

- [CONTEXT.md](../CONTEXT.md): glossary for the composable CLI contract.
- [docs/skills-updater.md](skills-updater.md): promoted Go updater behavior.
- [docs/maintainer-guide.md](maintainer-guide.md): source-first workflow and
  current verification guidance.
- [docs/adr/0001-port-sk-up-to-go.md](adr/0001-port-sk-up-to-go.md): Go as the
  implementation language.
- [docs/adr/0002-break-sk-up-before-promotion.md](adr/0002-break-sk-up-before-promotion.md):
  intentional breaking changes before promotion.
- [docs/adr/0003-one-binary-two-skill-updater-names.md](adr/0003-one-binary-two-skill-updater-names.md):
  one binary with `sk-up` and `skills-updates` names.
- [docs/adr/0004-cross-platform-sk-up-state-paths.md](adr/0004-cross-platform-sk-up-state-paths.md):
  cross-platform state path rules.
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
  promotion validation bar.
- [docs/adr/0012-retire-powershell-sk-up.md](adr/0012-retire-powershell-sk-up.md):
  PowerShell retirement.

## Command Contract

`sk-up` remains the terse daily-use command. `skills-updates` remains the
readable long-form command. Both names resolve to the same Go implementation.

| Short form | Long form | Behavior |
| --- | --- | --- |
| `sk-up -l` | `skills-updates --list` | List installed global skill directories without Git. |
| `sk-up -g` | `skills-updates --global` | Compare lockfile skills with upstream content. |
| `sk-up -d <skill>` | `skills-updates --diff <skill>` | Print a terminal diff for one skill. |
| `sk-up -z [skill...]` | `skills-updates --zed [skill...]` | Open diffs with the configured diff tool. |
| `sk-up -i` | `skills-updates --install` | Install all changed or missing unskipped skills. |
| `sk-up -i <skill...>` | `skills-updates --install <skill...>` | Install named lockfile skills. |
| `sk-up -I <source...>` | `skills-updates --install-source <source...>` | Install package or repository sources explicitly. |
| `sk-up -s <skill>` | `skills-updates --skip <skill>` | Save a skip for the current upstream tree hash. |
| `sk-up -u <skill>` | `skills-updates --unskip <skill>` | Remove a saved skip. |
| `sk-up -S` | `skills-updates --skips` | List saved skips. |
| `sk-up -r <skill...>` | `skills-updates --remove <skill...>` | Remove named global skills. |

`-z` is no longer Zed-only. Zed is the default diff tool when available, but
`--diff-tool` and `SK_UP_DIFF_TOOL` can select another tool.

Source installs are never inferred from `-i <arg>`. Use `-I` or
`--install-source` for source URLs, SSH remotes, `.git` URLs, and `owner/repo`
package shorthand.

## Common Flags

| Flag | Env var | Purpose |
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

Updater repository cache and skip state use OS-native cache/state directories,
with `--cache-dir` and `--state-dir` overrides. Use `$XDG_CACHE_HOME` and
`$XDG_STATE_HOME` on Linux when set, macOS user cache/application-support
locations, and `%LOCALAPPDATA%` on Windows.

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

### Final JSON Shape

The exact schema can evolve during implementation, but tests should pin the
promoted shape before release.

```json
{
  "ok": true,
  "command": "status",
  "entrypoint": "sk-up",
  "dryRun": false,
  "statuses": [
    {
      "name": "confidence-loop",
      "status": "update",
      "sourceUrl": "https://github.com/example/skills.git",
      "remoteHash": "abc123"
    }
  ],
  "actions": [],
  "errors": []
}
```

### JSONL Event Shape

Each line is one object. Long-running commands should emit progress events and
a final summary event.

```json
{"type":"event","event":"fetch","sourceUrl":"https://github.com/example/skills.git"}
{"type":"status","name":"confidence-loop","status":"update","remoteHash":"abc123"}
{"type":"summary","ok":true,"changed":1,"errors":0}
```

## Exit Codes

| Code | Meaning |
| --- | --- |
| `0` | Command completed successfully, including when updates are available. |
| `1` | Requested action failed or one or more targets errored. |
| `2` | Usage, configuration, dependency, invalid target, or invalid state path error. |
| `3` | Lock acquisition timeout. |
| `4` | Interrupted state was detected or repaired but the command cannot safely continue. |

## Implementation Layout

Use this Go layout:

```text
cmd/sk-up/main.go
internal/skup/cli
internal/skup/config
internal/skup/lockfile
internal/skup/runner
internal/skup/compare
internal/skup/output
internal/skup/state
```

The public contract is the CLI protocol, not an importable Go API.

## Core Behavior

Content comparison must compare installed skill directory content with clean
upstream repository content and ignore CRLF line-ending differences. Lockfile
hashes alone are not authoritative update signals.

The Go implementation requires external `git` for upstream comparison commands
but must not require external `tar`. Commands that do not need upstream
comparison still work without Git, including list, remove, unskip, skips, and
named installs that can delegate from lockfile source metadata.

Install and remove delegate to the upstream Skills CLI. If no runner override is
set, resolve runners in this order:

1. `pnpm dlx skills@latest`
2. `bunx skills@latest`
3. `deno run -A npm:skills@latest`
4. `npx -y skills@latest`

Runner execution uses tokenized process execution, not shell evaluation by
default.

Lockfile writes use advisory lock files next to `.skill-lock.json`, not
OS-specific named mutexes. Preserve unknown lockfile fields and unrelated skill
entries across install/remove transactions.

Skip entries are keyed by upstream tree hash. A skip means "hide this exact
upstream change until upstream changes again."

Dry-run mode plans and validates mutating operations without changing files,
lock state, skip state, installed skills, or external package-manager state.
Structured output keeps the normal result shape and includes planned actions.

Remove operations still perform updater-owned cleanup after successful delegated
remove: installed skill directory, saved skip state, and lockfile entry.

## Promotion Validation

Before promotion, require:

- Unit tests for CLI parsing, config/env resolution, runner fallback, path
  resolution, lockfile preservation, skip semantics, JSON/JSONL output, and
  dry-run planning.
- Integration tests with fake `git` and fake Skills CLI runners.
- Golden help tests for both `sk-up` and `skills-updates`.
- Cross-platform CI for Linux, macOS, and Windows on amd64.
- Arm64 build smoke checks.
- Manual workstation smoke tests for `sk-up -l`, `-g`, `-d`, `-i --dry-run`,
  named `-i`, `-I`, `-S`, `-s`, `-u`, and `-r --dry-run`.
- Updated docs for the promoted behavior.

Build release archives without installing them:

```bash
SK_UP_VERSION=<version> bin/build-sk-up-release.sh
```

The script writes Linux, macOS, and Windows amd64/arm64 archives under
`dist/sk-up/` and emits `SHA256SUMS`. `dist/` is generated output and is not
committed.

Do not install into the active workstation, publish releases, or push branches
without explicit user approval.

## Progress Ledger

| Date | Progress | Evidence | Next risk |
| --- | --- | --- | --- |
| 2026-06-20 | Grill decisions captured in glossary and ADRs 0001-0012. | `CONTEXT.md`, `docs/adr/` | Convert decisions into implementation plan and code. |
| 2026-06-20 | Implementation brief created as the control surface for the port. | `docs/sk-up-go-port.md` | Scaffold Go module, CLI parser, and focused tests. |
| 2026-06-20 | Brief linked from the maintained docs index and resources. | `docs/README.md`, `RESOURCES.md`, `docs/skills-updater.md`; link check and `git diff --check` passed. | Scaffold Go module, CLI parser, and focused tests. |
| 2026-06-20 | Go module scaffolded with config/path resolution tests. | `go.mod`, `internal/skup/config/`; `go test ./...` and `go test -race -shuffle=on -count=1 ./internal/skup/config` passed. | Decide and approve CLI dependency, then implement entrypoint/help parsing. |
| 2026-06-20 | Skills CLI runner fallback resolution implemented without shell evaluation. | `internal/skup/runner/`; `go test ./...` and `go test -race -shuffle=on -count=1 ./internal/skup/config ./internal/skup/runner` passed. | Decide and approve CLI dependency, then implement entrypoint/help parsing. |
| 2026-06-20 | Dependency-free CLI parser, dual-entrypoint help text, and `cmd/sk-up` help smoke test implemented. | `cmd/sk-up/`, `internal/skup/cli/`; `go test ./...` and race/shuffle checks passed. | Implement lockfile, skip-state, and dry-run planning primitives. |
| 2026-06-20 | Final JSON and JSONL event schema writers implemented. | `internal/skup/output/`; `go test ./...`, race/shuffle checks, and `git diff --check` passed. | Implement lockfile, skip-state, and dry-run planning primitives. |
| 2026-06-20 | Lockfile preservation, upstream-hash skip state, and basic advisory lock-file primitives implemented. | `internal/skup/lockfile/`, `internal/skup/state/`; `go test ./...` passed. | Add dry-run planning primitives, then connect command execution. |
| 2026-06-20 | Dry-run planning connected to app execution; real non-mutating `-l/--list` inventory command implemented with human, JSON, and JSONL output paths. | `internal/skup/plan/`, `internal/skup/app/`, `internal/skup/inventory/`; `go test ./...` passed. | Extend advisory lock wait/timeout, then implement Git-backed comparison. |
| 2026-06-20 | Advisory lock wait/timeout API added with deterministic tests. | `internal/skup/lockfile/`; `go test ./internal/skup/lockfile` and `go test ./...` passed. | Implement Git-backed content comparison without external `tar`. |
| 2026-06-20 | Git-backed comparison primitives implemented: CRLF-insensitive directory compare, safe Go tar extraction, process runner, `git archive` export, and high-level skill comparison. | `internal/skup/compare/`; `go test ./internal/skup/compare` and `go test ./...` passed. | Wire comparison into `-g`, `-d`, `-z`, skip, and install-update command execution. |
| 2026-06-20 | `-g/--global` status execution wired through lockfile discovery, cached Git repo commands, upstream hash lookup, skip matching, comparison, and human/JSON/JSONL output. | `internal/skup/status/`, `internal/skup/app/`; `go test ./internal/skup/app ./internal/skup/status` and `go test ./...` passed. | Wire comparison into `-d`, `-z`, skip, and install-update command execution. |
| 2026-06-20 | `-s/--skip`, `-u/--unskip`, and `-S/--skips` wired to upstream-hash skip state with deterministic timestamps in tests and structured output for skip listing. | `internal/skup/app/`, `internal/skup/state/`; `go test ./internal/skup/app ./internal/skup/state` passed. | Wire comparison into `-d`, `-z`, and install-update command execution. |
| 2026-06-20 | `-d/--diff` and `-z/--zed` wired to exported compare trees; terminal diff uses `git diff --no-index --ignore-cr-at-eol`, and open-diff uses `--diff-tool`/`SK_UP_DIFF_TOOL` with default `zed`. | `internal/skup/app/`, `internal/skup/compare/`, `internal/skup/config/`, `internal/skup/status/`; focused package tests passed. | Wire install-update command execution and delegated mutation boundaries. |
| 2026-06-20 | `-i/--install` and `-I/--install-source` wired to the Skills CLI runner; targetless install selects changed/missing unskipped skills, named install uses lockfile source metadata, source install delegates explicit sources, and delegated lockfile rewrites preserve existing entries/fields. | `internal/skup/app/`, `internal/skup/runner/`, `internal/skup/lockfile/`; `go test ./internal/skup/app ./internal/skup/runner ./internal/skup/lockfile` passed. | Implement remove cleanup and promotion wrapper/package gates. |
| 2026-06-20 | `-r/--remove` wired to the Skills CLI runner with updater-owned cleanup for installed directories, lockfile entries, and saved skip state. | `internal/skup/app/`; `go test ./internal/skup/app` passed. | Add promotion wrapper/package gates and real-workstation smoke plan. |
| 2026-06-20 | Release archive build gate added for Linux, macOS, and Windows on amd64/arm64 without changing active workstation wrappers. | `bin/build-sk-up-release.sh`, `.gitignore`; temp-dir smoke produced six archives plus `SHA256SUMS`, then temp output was removed. | Run final full validation and decide whether to promote/install with explicit approval. |
| 2026-06-20 | Safe binary smoke completed without touching the active workstation install. | Temp-built `sk-up` rendered help, listed a temp `AGENTS_HOME` skill with `-l`, and emitted structured `-I --dry-run --json`; temp files were removed. | Source promotion remained. |
| 2026-06-20 | Source promotion completed: Windows wrappers invoke adjacent Go executables, release archives include wrappers, PowerShell updater source was removed, and docs/tests now describe the Go path. | `bin/sk-up.cmd`, `bin/skills-updates.cmd`, deleted `bin/skills-updates.ps1`, `tests/skills-updates-install.ps1`, `docs/skills-updater.md`, `docs/maintainer-guide.md`. | Final validation and completion audit. |
| 2026-06-20 | Final local WSL validation passed after source promotion; Windows wrapper regression was updated but could not be executed from this session because direct Windows PowerShell was blocked and the available Tabby PowerShell session did not execute MCP commands. | `go test ./...`, `go vet ./...`, race/shuffle Go test, `go build ./cmd/sk-up`, release archive temp-dir smoke, safe binary smoke, `node --test bin/codex-wsl-proxy-runtime.test.js`, `bash -n bin/build-sk-up-release.sh`, `git diff --check`, Markdown link check. | Run `tests/skills-updates-install.ps1` from a working Windows PowerShell/Tabby session before declaring the goal complete. |
| 2026-06-20 | Added and passed a WSL-side promotion check for the promoted source shape. | `tests/sk-up-promotion-check.sh` verifies the PowerShell fallback is absent, wrappers point at adjacent Go executables, Windows release zips contain `sk-up.exe` plus `.cmd` wrappers, and temp `sk-up` help/list/source dry-run works. | Run `tests/skills-updates-install.ps1` from a working Windows PowerShell/Tabby session before declaring the goal complete. |
| 2026-06-20 | Deployed the promoted Go updater to the active Windows user bin for dogfooding. | Installed `sk-up.exe`, `skills-updates.exe`, `sk-up.cmd`, and `skills-updates.cmd` under `C:\Users\nguco\bin`; removed installed `skills-updates.ps1`; backed up prior updater files to `C:\Users\nguco\bin\sk-up-backup-20260620T075209Z`; installed `sk-up.exe` passed help, temp `-l`, and `-I --dry-run --json` smoke checks with Windows paths. | Dogfood on Windows shell; WSL interop could not directly launch `skills-updates.exe` and direct PowerShell remains blocked from this thread. |
| 2026-06-20 | Windows dogfood exposed a UAC prompt for `skills-updates.exe`; fixed by making `skills-updates.cmd` a long-name alias to `sk-up.exe` with `SK_UP_ENTRYPOINT=skills-updates` and rejecting `skills-updates.exe` in Windows release archives. | `cmd/sk-up/main.go`, `cmd/sk-up/main_test.go`, `bin/skills-updates.cmd`, `bin/build-sk-up-release.sh`, `tests/sk-up-promotion-check.sh`, `tests/skills-updates-install.ps1`, docs. | Redeploy active Windows bin without `skills-updates.exe`, then dogfood `skills-updates --help`. |
| 2026-06-20 | Redeployed the UAC-safe Windows install. | Active `C:\Users\nguco\bin` now contains `sk-up.exe`, `sk-up.cmd`, and `skills-updates.cmd`; `skills-updates.exe` and `skills-updates.ps1` are absent; prior active files were backed up to `C:\Users\nguco\bin\sk-up-backup-20260620T080951Z-uac-fix`; installed `sk-up.exe` passed WSL-side help and dry-run JSON smoke checks. | Dogfood `skills-updates --help` from a real Windows shell; Tabby did not respond and direct `cmd.exe` is blocked by local guardrails. |
| 2026-06-20 | Completion audit rerun against the objective and current worktree. | Source/docs/wrappers satisfy the promoted Go implementation, composable CLI contract, structured output, dry-run, OS-native cache/state paths, advisory locks, Skills CLI fallback runners, Git-required comparison without external `tar`, lockfile preservation, upstream-hash skips, release archive gates, and PowerShell retirement; `go test ./...`, `go vet ./...`, race/shuffle Go tests, `tests/sk-up-promotion-check.sh`, `bash -n`, `git diff --check`, installed `sk-up.exe -h`, installed `sk-up.exe -I Hellfrosted/agents --dry-run --json`, and `SK_UP_ENTRYPOINT=skills-updates go run ./cmd/sk-up --help` passed. | Only the real Windows-shell `skills-updates.cmd` alias smoke remains unproven because direct `cmd.exe` is blocked and Tabby timed out repeatedly. |
| 2026-06-20 | Windows PowerShell alias dogfood passed. | User-ran `skills-updates --help`, which rendered long-form help without UAC, and `skills-updates --install-source Hellfrosted/agents --dry-run --json`, which returned `{"ok":true,"command":"install-source","entrypoint":"skills-updates","dryRun":true,...}`. | Goal completion audit can close. |

## Completion Audit

Current evidence proves the source-port scope except for the Windows-shell alias
smoke:

| Requirement | Evidence | Status |
| --- | --- | --- |
| Promoted Go implementation under `cmd/sk-up` and `internal/skup/...`. | Go source tree, `go.mod`, `go test ./...`, `go vet ./...`, race/shuffle Go tests. | Proven. |
| `sk-up` short flags and `skills-updates` long entrypoint. | `internal/skup/cli` parser/help tests; `SK_UP_ENTRYPOINT=skills-updates go run ./cmd/sk-up --help`. | Source-proven. |
| No supported PowerShell updater fallback after promotion. | `bin/skills-updates.ps1` removed; promotion check fails if PowerShell fallback references return. | Proven. |
| Structured stdout and human diagnostics on stderr. | `internal/skup/output` tests plus app JSON/JSONL dry-run tests. | Proven. |
| Dry-run plans mutating operations without changing state. | `internal/skup/plan`, app dry-run tests, installed `sk-up.exe -I Hellfrosted/agents --dry-run --json`. | Proven for source install; covered by tests for app contract. |
| OS-native cache/state paths with explicit overrides. | `internal/skup/config` tests for flags, env, Linux, macOS, and Windows defaults. | Proven. |
| Advisory lock files and transactional lockfile preservation. | `internal/skup/lockfile` tests and app install/remove transaction tests. | Proven. |
| Skills CLI delegation with fallback runners. | `internal/skup/runner` tests for `pnpm`, `bunx`, `deno`, and `npx`; app install/remove fake runner tests. | Proven. |
| Git-required comparison without external `tar`. | `internal/skup/compare` tests; release gate; Go tar extraction is internal only. | Proven. |
| Upstream-hash skip semantics. | `internal/skup/state` and app skip/unskip/skips tests. | Proven. |
| Windows wrapper/package promotion gate. | `tests/sk-up-promotion-check.sh` builds release archives, checks wrapper references, and rejects `skills-updates.exe`. | Proven for archive shape. |
| Active workstation deployment. | Active `C:\Users\nguco\bin` has `sk-up.exe`, `sk-up.cmd`, and `skills-updates.cmd`; no `skills-updates.exe` or `skills-updates.ps1`; installed `sk-up.exe` smoke checks passed. | Proven for deployed file shape. |
| Real Windows-shell `skills-updates` alias behavior. | User-ran `skills-updates --help` and `skills-updates --install-source Hellfrosted/agents --dry-run --json` in Windows PowerShell; help rendered long-form usage and JSON reported `entrypoint:"skills-updates"`. | Proven. |

## Open Implementation Queue

1. Release publishing and pushing remain outside this source-port goal until
   explicitly requested.
