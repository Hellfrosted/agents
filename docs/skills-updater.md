# Skills Updater

The skills updater maintains globally installed Codex skills under the user's
universal `.agents/skills` directory. It compares installed skill directories
with upstream repository content, opens diffs, installs changed or missing
skills, records temporary skips, and removes global skills cleanly.

## Entry Points

| File | Use |
| --- | --- |
| `bin/sk-up.cmd` | Short Windows command names and flags. |
| `bin/skills-updates.cmd` | Long Windows command names and flags. |
| `bin/skills-updates.ps1` | PowerShell implementation. |
| `tests/skills-updates-install.ps1` | Regression test for source URL installs. |

Both `.cmd` wrappers switch the console to UTF-8, run the PowerShell script with
`-NoProfile -ExecutionPolicy Bypass`, pass `--cmd-name` so help text matches the
invoked wrapper, then restore the previous codepage.

## Commands

Short form:

```bat
sk-up -l
sk-up -g
sk-up -d confidence-loop
sk-up -z confidence-loop evo-end-to-end
sk-up -i
sk-up -i confidence-loop
sk-up -i owner/repo
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
skills-updates --install owner/repo
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
- `--install` / `-i <source>`: install a source URL or `owner/repo` package.
- `--skip` / `-s <skill>`: save a skip for the current upstream tree hash.
- `--unskip` / `-u <skill>`: remove one saved skip.
- `--skips` / `-S`: list saved skips.
- `--remove` / `-r <skill ...>`: uninstall named global skills.

`-g` is only for status checks. Do not combine it with install, diff, or Zed
operations.

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

Repository clones and skip state live under:

```text
%LOCALAPPDATA%\skills-updates
```

When `LOCALAPPDATA` is unavailable, the script uses a temp directory named
`skills-updates-state`.

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

Install operations require `pnpm` and call:

```text
pnpm dlx skills@latest add <source> -g -y --agent universal --skill <skill-name>
pnpm dlx skills@latest add <source> -g -y --agent universal
```

Named installs use lockfile sources. Source installs accept URL, SSH, `.git`,
or `owner/repo` arguments and do not mix with named lockfile installs in the
same command.

Uninstall operations require `pnpm` and call:

```text
pnpm dlx skills@latest remove -g -y --agent universal --skill <skill-name>
```

After the Skills CLI returns, uninstall also removes the installed skill
directory, clears any saved skip for that skill, and removes the lockfile entry.
If post-CLI cleanup fails, the script restores the pre-uninstall lockfile
snapshot so directory and lockfile state do not silently diverge.

## Lockfile Safety

Lockfile writes are guarded with a named mutex derived from the lockfile path.
The script also writes an adjacent backup before raw lockfile replacement and
repairs from that backup if it detects an interrupted prior write.

Skills CLI operations are wrapped in state transactions so existing lockfile
entries and unrelated lockfile fields survive add/remove commands.

## Verification

Help-path smoke checks:

```powershell
bin\skills-updates.cmd --help
bin\sk-up.cmd -h
```

Source install regression:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File tests\skills-updates-install.ps1
```

That test creates a temporary `AGENTS_HOME`, places a fake `pnpm` first on
`PATH`, runs a source URL install, asserts the expected `pnpm dlx skills@latest
add <source>` invocation, and removes its temp directory.
