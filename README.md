# codex-wsl

Run Codex from Windows while the real process lives in WSL.

The Windows shim starts a small WSL-side proxy. The proxy translates Windows and
WSL paths in both directions, then hands off to the real `codex` binary. For
T3code-style non-interactive sessions, the proxy starts `codex app-server`
automatically and shuts it down after the configured idle period.

This workstation also installs LazyCodex as the `omo@sisyphuslabs` Codex
plugin. The WSL runner keeps that plugin enabled, but disables its telemetry,
auto-update, and config-migration startup paths for T3code app-server sessions.

## Layout

- `bin/codex-wsl.cmd` is the Windows entry point.
- `bin/codex-wsl-proxy-runner.sh` finds `node` inside WSL and starts the proxy.
- `bin/codex-wsl-proxy.js` translates paths and forwards traffic to Codex.
- `skills/` contains Codex skills that can be installed with the Skills CLI.
- `bin/skills-updates.ps1`, `bin/skills-updates.cmd`, and `bin/sk-up.cmd`
  check, diff, install, skip, and remove globally installed skills.
- `docs/codex-cli-tooling.md` documents companion tools used around Codex:
  Evo, RTK, ICM, Codex Security, OpenAI Developer Docs MCP, and adjacent
  utilities.
- `docs/discrawl-wiretap.md` documents the local Discrawl Vesktop wiretap setup
  and its user-level systemd auto-sync timer.

## Install

Copy the Windows shim somewhere on your Windows `PATH`, for example:

```bat
copy bin\codex-wsl.cmd %USERPROFILE%\bin\codex-wsl.cmd
```

Copy the WSL scripts into your WSL user:

```bash
mkdir -p ~/.local/bin
cp bin/codex-wsl-proxy-runner.sh ~/.local/bin/
cp bin/codex-wsl-proxy.js ~/.local/bin/
chmod +x ~/.local/bin/codex-wsl-proxy-runner.sh ~/.local/bin/codex-wsl-proxy.js
```

## Defaults

By default, `codex-wsl.cmd` uses the `Ubuntu` WSL distro and that distro's
default user. The runner uses `$HOME/.codex` for `CODEX_HOME`, prepends pnpm,
Bun, and local user bin directories to `PATH`, and expects the proxy at
`$HOME/.local/bin/codex-wsl-proxy-runner.sh`.

Set these environment variables only when your setup is different:

- `CODEX_WSL_EXE`: Windows path to `wsl.exe`.
- `CODEX_WSL_DISTRO`: WSL distro name, such as `Ubuntu-24.04`.
- `CODEX_WSL_USER`: WSL user to run as.
- `CODEX_WSL_PROXY`: WSL path to `codex-wsl-proxy-runner.sh`.
- `CODEX_WSL_PROXY_JS`: WSL path to `codex-wsl-proxy.js`.
- `CODEX_WSL_PROXY_DISTRO`: distro name to use when converting WSL paths to
  `\\wsl.localhost\...` paths. WSL usually sets this itself.
- `CODEX_WSL_PROXY_TARGET`: WSL path to the real `codex` binary.
- `CODEX_WSL_PROXY_IDLE_TIMEOUT_MS`: app-server idle timeout. The Windows shim
  defaults this to `1800000` milliseconds.
- `CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS`: timeout before the proxy synthesizes a
  fallback `skills/list` response. The proxy defaults this to `2000`
  milliseconds.
- `CODEX_WSL_PROXY_DEBUG_LOG`: WSL path for proxy debug logs.
- `CODEX_WSL_SHIM_DEBUG`: print the Windows shim launch arguments.
- `CODEX_SKILLS_DIRS` and `CODEX_SKILL_ROOTS`: extra skill roots passed through
  to WSL.
- `CODEX_HOME`: Codex home directory inside WSL.

LazyCodex startup defaults in the WSL runner:

- `OMO_CODEX_DISABLE_POSTHOG=1`
- `OMO_CODEX_SEND_ANONYMOUS_TELEMETRY=0`
- `OMO_DISABLE_POSTHOG=1`
- `OMO_SEND_ANONYMOUS_TELEMETRY=0`
- `LAZYCODEX_AUTO_UPDATE_DISABLED=1`
- `OMO_CODEX_AUTO_UPDATE_DISABLED=1`
- `LAZYCODEX_CONFIG_MIGRATION_DISABLED=1`
- `OMO_CODEX_CONFIG_MIGRATION_DISABLED=1`

## Use

Put the full path to `codex-wsl.cmd` into T3code. Arguments are passed through
to Codex after Windows paths are converted to WSL paths.

The WSL side needs `node` available on `PATH`.

The proxy translates arguments before spawning Codex, translates known JSON
protocol path fields in both directions, and converts WSL filesystem paths in
Codex output back to Windows paths when possible. Paths under `/mnt/<drive>/`
become drive-letter paths, and Linux paths become `\\wsl.localhost\<distro>\...`
when the distro name is known.

When the proxy receives no arguments and stdin is not a TTY, it runs
`codex app-server`. It tracks protocol activity and active turns so the idle
reaper does not stop the app server mid-turn.

If the upstream app server does not answer `skills/list` quickly enough, the
proxy returns a fallback skill list. It searches:

- `$HOME/.codex/skills`
- `$HOME/.codex/skills/.system`
- `$HOME/.agents/skills`
- repo-local `.codex/skills`
- repo-local `.agents/skills`
- roots from `CODEX_SKILLS_DIRS` and `CODEX_SKILL_ROOTS`
- nested `skills/` directories under `$HOME/.codex/plugins/cache`

## Skills Updater

The Windows wrappers run `bin/skills-updates.ps1` with UTF-8 console output.
Use `sk-up` for short flags or `skills-updates` for long flags:

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

Equivalent long-form examples:

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

The updater reads the global skill lockfile from `%AGENTS_HOME%\.skill-lock.json`
when `AGENTS_HOME` is set, otherwise from `%USERPROFILE%\.agents\.skill-lock.json`.
It compares installed skill folders with upstream repository content, not only
lockfile hashes. Repository clones and skip state live under
`%LOCALAPPDATA%\skills-updates` when available, otherwise under the temp
directory.

Installs and uninstalls require `pnpm`. The updater uses
`pnpm dlx skills@latest`, forces global operations to the universal
`.agents/skills` target, guards lockfile writes with a mutex, and preserves
existing lockfile entries around Skills CLI operations. Uninstalls also remove
the global installed skill directory, clear saved skips for the skill, and
remove the skill's lockfile entry under the same operation lock. If post-CLI
cleanup fails, the updater restores the pre-uninstall lockfile snapshot so
directory and lockfile state do not diverge. Saved skips are tied to the
current upstream tree hash and expire when upstream changes.
