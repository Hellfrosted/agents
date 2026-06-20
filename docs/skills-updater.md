# Skills Updater

The skills updater maintains globally installed Codex skills under the user's
universal `.agents/skills` directory. It compares installed skill directories
with upstream repository content, opens diffs, installs changed or missing
skills, records temporary skips, and removes global skills cleanly. It does not
install Codex plugins; use `codex plugin ...` for plugin marketplaces and use
this updater for Skills CLI packages such as `mattpocock/skills`.

## Entry Points

| File | Use |
| --- | --- |
| `cmd/sk-up/` | Go source for the updater binary. |
| `bin/sk-up.cmd` | Windows launcher for adjacent `sk-up.exe`. |
| `bin/skills-updates.cmd` | Windows long-name alias for adjacent `sk-up.exe`. |
| `bin/build-sk-up-release.sh` | Cross-platform archive builder. |
| `tests/skills-updates-install.ps1` | Windows wrapper/source-install dry-run regression. |

Both `.cmd` wrappers switch the console to UTF-8, invoke the adjacent promoted
Go executable, then restore the previous codepage. They do not call PowerShell
as an updater fallback. Windows builds intentionally ship only `sk-up.exe`;
`skills-updates.cmd` sets `SK_UP_ENTRYPOINT=skills-updates` for the child
process so long-form help and structured `entrypoint` fields remain compatible
without placing a `skills-updates.exe` updater-looking binary on `PATH`.

## Portable Go Implementation

This page describes the promoted Go implementation. The implementation brief
and progress ledger remain in [sk-up-go-port.md](sk-up-go-port.md).

## Commands

Short form:

```bat
sk-up -l
sk-up -g
sk-up -d confidence-loop
sk-up -z confidence-loop evo-end-to-end
sk-up -i
sk-up -i confidence-loop
sk-up -I owner/repo
sk-up -s confidence-loop
sk-up -u confidence-loop
sk-up -S
sk-up -r confidence-loop
```

Long form:

```bat
skills-updates --list
skills-updates --global
skills-updates --diff confidence-loop
skills-updates --zed confidence-loop evo-end-to-end
skills-updates --install
skills-updates --install confidence-loop
skills-updates --install-source owner/repo
skills-updates --skip confidence-loop
skills-updates --unskip confidence-loop
skills-updates --skips
skills-updates --remove confidence-loop
```

Aliases:

- `--gui` is the same as `--zed`.
- `--uninstall` is the same as `--remove`.
- `--install-all` and `install-all` are explicit forms of targetless
  `--install`.

## Modes

- `--list` / `-l`: list installed global skills without checking upstream.
- `--global` / `-g`: fetch upstream repositories and print status for every
  global lockfile skill.
- `--diff` / `-d <skill>`: show a terminal diff for one changed skill.
- `--zed` / `-z [skill ...]`: open Zed diff views for changed skills.
- `--install` / `-i`: install every changed or missing unskipped skill.
- `--install` / `-i <skill ...>`: install named lockfile skills.
- `--install-source` / `-I <source ...>`: install source URLs, SSH remotes,
  `.git` URLs, or `owner/repo` packages explicitly.
- `--skip` / `-s <skill>`: save a skip for the current upstream tree hash.
- `--unskip` / `-u <skill>`: remove one saved skip.
- `--skips` / `-S`: list saved skips.
- `--remove` / `-r <skill ...>`: uninstall named global skills.

`-g` is only for status checks. Do not combine it with install, diff, or Zed
operations.

`--list` and `--global` answer different questions. `--list` reads installed
skill directories from `%AGENTS_HOME%\skills`; `--global` checks only skills
with lockfile metadata. If a skill appears in `--list` but not `--global`, treat
it as unmanaged drift until it is either intentionally documented as local-only
or reinstalled through the managed flow. If a lockfile entry has a `pluginName`,
that is Skills CLI package metadata, not a Codex plugin installation.

Last observed unmanaged drift on this workstation: `tdd` was installed under
`%USERPROFILE%\.agents\skills` but absent from `.skill-lock.json` before the
Matt Pocock package refresh.

## State Paths

The updater resolves global skill state from:

```text
%AGENTS_HOME%\.skill-lock.json
%AGENTS_HOME%\skills
```

When `AGENTS_HOME` is not set, it falls back to:

```text
%USERPROFILE%\.agents\.skill-lock.json
%USERPROFILE%\.agents\skills
```

Repository clones and skip state live in OS-native updater directories unless
overridden:

```text
Linux cache: $XDG_CACHE_HOME/sk-up or $HOME/.cache/sk-up
Linux state: $XDG_STATE_HOME/sk-up or $HOME/.local/state/sk-up
macOS cache: $HOME/Library/Caches/sk-up
macOS state: $HOME/Library/Application Support/sk-up
Windows cache/state: %LOCALAPPDATA%\sk-up\cache and %LOCALAPPDATA%\sk-up\state
```

Use `--cache-dir`, `--state-dir`, `SK_UP_CACHE_DIR`, and `SK_UP_STATE_DIR` for
explicit overrides.

## Comparison Model

For each lockfile skill, the updater:

1. resolves `sourceUrl` or legacy GitHub `source`;
2. groups skills by source repository;
3. updates or clones each source repository;
4. sparse-checks the relevant skill paths;
5. exports clean compare trees;
6. diffs installed directories against upstream content with CRLF ignored at
   line ends;
7. reports `OK`, `UPDATE`, `MISSING`, `SKIP`, or `ERROR`.

This compares real folder content, not just lockfile hashes.

Saved skips store the upstream tree hash. A new upstream tree makes the update
visible again.

## Install And Uninstall

Install operations delegate to the Skills CLI runner. The default runner order
is `pnpm dlx skills@latest`, `bunx skills@latest`,
`deno run -A npm:skills@latest`, then `npx -y skills@latest`. Override with
`--skills-command` or `SK_UP_SKILLS_COMMAND`.

Named install commands call the runner as:

```text
<skills-runner> add <source> -g -y --agent universal --skill <skill-name>
```

Source install commands use explicit `-I` / `--install-source` and call:

```text
<skills-runner> add <source> -g -y --agent universal
```

Named installs use lockfile sources. Source installs accept URL, SSH, `.git`, or
`owner/repo` arguments. Package installs can add multiple lockfile entries at
once; for example, `mattpocock/skills` records the active selected skills with
`pluginName: "mattpocock-skills"`.

Uninstall operations call:

```text
<skills-runner> remove -g -y --agent universal --skill <skill-name>
```

After the Skills CLI returns, uninstall also removes the installed skill
directory, clears any saved skip for that skill, and removes the lockfile entry.
Install and remove operations preserve existing lockfile entries and unrelated
lockfile fields around delegated Skills CLI mutations.

## Matt Pocock Skills Package Notes

The current `mattpocock/skills` package exposes `/ask-matt` as the router for
choosing a flow. Its active package manifest includes user-invoked skills such
as `ask-matt`, `grill-with-docs`, `triage`,
`improve-codebase-architecture`, `setup-matt-pocock-skills`, `to-prd`,
`to-issues`, `implement`, `prototype`, `grill-me`, `handoff`, `teach`, and
`writing-great-skills`, plus model-invoked skills such as `diagnosing-bugs`,
`tdd`, `domain-modeling`, `codebase-design`, and `grilling`.

Matt skills that are outside the active package manifest can still appear in a
lockfile if they were installed individually before the package refresh. Treat
entries under `skills/deprecated/`, `skills/in-progress/`, or `skills/personal/`
as explicit legacy installs when they do not include `pluginName:
"mattpocock-skills"`. Remove them only when the user asks to retire that skill:

```bat
sk-up -r request-refactor-plan
```

Do not treat the upstream `.claude-plugin/plugin.json` as a Codex plugin
manifest. For this workstation, Matt's package is managed through the Skills
CLI and the universal `.agents/skills` install.

## Lockfile Safety

Lockfile writes are guarded with advisory lock files next to
`.skill-lock.json`, not OS-specific named mutexes.

Skills CLI operations are wrapped in state transactions so existing lockfile
entries and unrelated lockfile fields survive add/remove commands.

## Verification

Run Go validation from WSL:

```bash
go test ./...
go test -race -shuffle=on -count=1 ./cmd/sk-up ./internal/skup/...
SK_UP_DIST_DIR="$(mktemp -d)" SK_UP_VERSION=test bin/build-sk-up-release.sh
```

Run the Windows wrapper regression from Windows PowerShell. From WSL, launch
Windows PowerShell through WSL init because the Windows drive mount can make
`.exe` files appear non-executable:

```bash
/init /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests/skills-updates-install.ps1
```

Source install regression:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests\skills-updates-install.ps1
```

That test builds temporary `sk-up.exe` beside copied wrappers, checks both help
entrypoint names, verifies no `skills-updates.exe` is created, runs
`-I owner/repo --dry-run --json`, and removes its temp directory.

Inventory drift check:

```powershell
bin\skills-updates.cmd --list
bin\skills-updates.cmd --global
```

The first command should include every installed global skill directory. The
second should include every lockfile-backed skill that the updater can compare
against upstream.
