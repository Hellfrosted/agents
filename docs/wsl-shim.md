# Codex WSL Shim

The WSL shim lets Windows tools start Codex while the real Codex process runs
inside WSL. It is optimized for T3code-style non-interactive sessions where
Codex is launched as an app server.

## Components

| File | Responsibility |
| --- | --- |
| `bin/codex-wsl.cmd` | Windows entry point. Finds `wsl.exe`, selects distro/user, preserves Windows cwd, passes selected environment names through WSL, and launches the WSL runner. |
| `bin/codex-wsl-proxy-runner.sh` | WSL startup wrapper. Normalizes `HOME`, `USER`, `CODEX_HOME`, and `PATH`; disables LazyCodex startup mutation paths; finds Node; launches the proxy. |
| `bin/codex-wsl-proxy.js` | Polyglot Bash/Node entry point that loads the runtime. |
| `bin/codex-wsl-proxy-runtime.js` | Child-process lifecycle, app-server idle reaping, active-turn tracking, stream forwarding, and `skills/list` fallback timing. |
| `bin/codex-wsl-path-translation.js` | Windows-to-WSL and WSL-to-Windows path conversion for known protocol path fields. |
| `bin/codex-wsl-skills-fallback.js` | Synthetic `skills/list` response when upstream Codex does not answer quickly enough. |
| `bin/codex-wsl-proxy-runtime.test.js` | Node tests for the runtime-adjacent policies. |

## Install Shape

The Windows entry point should be a native Windows symlink:

```text
C:\Users\nguco\bin\codex-wsl.cmd -> E:\dev\agents-toolkit\bin\codex-wsl.cmd
```

Repair that symlink with Windows tooling, not WSL `ln -s` on `/mnt/c`.

The WSL runner expects the proxy files under `~/.local/bin` unless overridden:

```bash
mkdir -p ~/.local/bin
cp bin/codex-wsl-proxy-runner.sh ~/.local/bin/
cp bin/codex-wsl-proxy.js ~/.local/bin/
cp bin/codex-wsl-proxy-runtime.js ~/.local/bin/
cp bin/codex-wsl-path-translation.js ~/.local/bin/
cp bin/codex-wsl-skills-fallback.js ~/.local/bin/
chmod +x ~/.local/bin/codex-wsl-proxy-runner.sh ~/.local/bin/codex-wsl-proxy.js
```

Keep repo source and installed runtime separate. Ordinary development edits stay
in `bin/`; copying to `~/.local/bin` is an explicit install or repair step.

## Defaults

`bin/codex-wsl.cmd` defaults to:

- WSL executable: `%SystemRoot%\System32\wsl.exe`
- distro: `Ubuntu`
- distro user: WSL default user
- proxy runner: `$HOME/.local/bin/codex-wsl-proxy-runner.sh`
- app-server idle timeout: `1800000` milliseconds
- outbound Linux path prefix: `\\wsl.localhost\<distro>\...` when the distro
  name is known

The WSL runner defaults to:

- `CODEX_HOME=$HOME/.codex`
- `PATH` prefixed with pnpm, Bun, user-local, and user bin directories
- LazyCodex / OMO startup mutation paths disabled for app-server sessions:
  - `OMO_CODEX_DISABLE_POSTHOG=1`
  - `OMO_CODEX_SEND_ANONYMOUS_TELEMETRY=0`
  - `OMO_DISABLE_POSTHOG=1`
  - `OMO_SEND_ANONYMOUS_TELEMETRY=0`
  - `LAZYCODEX_AUTO_UPDATE_DISABLED=1`
  - `OMO_CODEX_AUTO_UPDATE_DISABLED=1`
  - `LAZYCODEX_CONFIG_MIGRATION_DISABLED=1`
  - `OMO_CODEX_CONFIG_MIGRATION_DISABLED=1`

## Environment Overrides

Set these only when the local setup differs from the defaults:

| Variable | Meaning |
| --- | --- |
| `CODEX_WSL_EXE` | Windows path to `wsl.exe`. |
| `CODEX_WSL_DISTRO` | WSL distro name, such as `Ubuntu-24.04`. |
| `CODEX_WSL_USER` | WSL user to run as. |
| `CODEX_WSL_PROXY` | WSL path to `codex-wsl-proxy-runner.sh`. |
| `CODEX_WSL_PROXY_JS` | WSL path to `codex-wsl-proxy.js`. |
| `CODEX_WSL_PROXY_DISTRO` | Distro name used for outbound `\\wsl.localhost\...` paths. |
| `CODEX_WSL_PROXY_TARGET` | WSL path to the real `codex` binary. |
| `CODEX_WSL_PROXY_IDLE_TIMEOUT_MS` | App-server idle timeout; `0` disables idle reaping. |
| `CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS` | Delay before synthetic `skills/list`; default `2000`. |
| `CODEX_WSL_PROXY_DEBUG_LOG` | WSL path for proxy debug logs. |
| `CODEX_WSL_SHIM_DEBUG` | Print Windows shim launch details. |
| `CODEX_SKILLS_DIRS` | Extra skill roots for fallback discovery. |
| `CODEX_SKILL_ROOTS` | Extra skill roots for fallback discovery. |
| `CODEX_HOME` | Codex home directory inside WSL after path conversion. |

The Windows shim also passes selected credential and proxy variables through WSL
when already present, including `OPENAI_API_KEY`, `OPENAI_BASE_URL`,
`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, `SSL_CERT_FILE`, `SSL_CERT_DIR`, and
`NODE_EXTRA_CA_CERTS`.

## Runtime Behavior

Arguments are converted from Windows paths to WSL paths before spawning Codex.
With no arguments and non-TTY stdin, the proxy starts:

```bash
codex app-server
```

The runtime then:

- forwards stdin/stdout/stderr to the child;
- converts known inbound protocol path fields to WSL paths;
- converts known outbound protocol path fields to Windows paths;
- tracks protocol activity and active turn ids;
- reaps app-server sessions only after the configured idle period and only when
  no active turn is known;
- forwards signals to the child and escalates to `SIGKILL` after five seconds;
- returns a synthetic `skills/list` response if upstream Codex misses the
  configured timeout.

## Path Translation Policy

Inbound examples:

```text
C:\Users\me\repo       -> /mnt/c/Users/me/repo
\\wsl.localhost\Ubuntu\home\me\repo -> /home/me/repo
```

Outbound examples:

```text
/mnt/e/dev/agents-toolkit -> E:\dev\agents-toolkit
/home/crunch/project      -> \\wsl.localhost\Ubuntu\home\crunch\project
```

Only known protocol path keys are transformed. Ordinary strings such as status
messages are left alone. Map keys under `fileChanges` and arrays under known
path-array fields are also converted.

## Skills Fallback

When the app server does not answer `skills/list` quickly enough, the proxy
builds a fallback response by scanning:

- `$HOME/.codex/skills`
- `$HOME/.codex/skills/.system`
- `$HOME/.agents/skills`
- repo-local `.codex/skills`
- repo-local `.agents/skills`
- per-cwd extra roots from the request
- roots from `CODEX_SKILLS_DIRS` and `CODEX_SKILL_ROOTS`
- nested `skills/` directories under `$HOME/.codex/plugins/cache`

The fallback reads `SKILL.md` frontmatter, infers a skill scope, validates the
response schema, sorts skills by name, and suppresses a later upstream response
for the same request id if the fallback already answered.

## Verification

Run the focused Node tests from the repo root:

```bash
node --test bin/codex-wsl-proxy-runtime.test.js
```

Useful install checks:

```bash
ls -l ~/.local/bin/codex-wsl-proxy*.js ~/.local/bin/codex-wsl-path-translation.js ~/.local/bin/codex-wsl-skills-fallback.js
test -x ~/.local/bin/codex-wsl-proxy-runner.sh
```

For Windows symlink checks, use Windows-native commands such as `cmd.exe /c dir`
or PowerShell `Get-Item`.
