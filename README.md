# codex-wsl

Run Codex from Windows while the real process lives in WSL.

The Windows shim starts a small WSL-side proxy. The proxy translates Windows and
WSL paths in both directions, then hands off to the real `codex` binary.

## Layout

- `bin/codex-wsl.cmd` is the Windows entry point.
- `bin/codex-wsl-proxy-runner.sh` finds `node` inside WSL and starts the proxy.
- `bin/codex-wsl-proxy.js` translates paths and forwards traffic to Codex.
- `skills/` contains Codex skills that can be installed with the Skills CLI.
- `docs/codex-cli-tooling.md` documents companion tools used around Codex:
  Evo, RTK, ICM, Codex Security, OpenAI Developer Docs MCP, and adjacent
  utilities.

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

By default, `codex-wsl.cmd` uses your default WSL distro and that distro's
default user. The runner uses `$HOME/.codex` for `CODEX_HOME` and expects the
proxy at `$HOME/.local/bin/codex-wsl-proxy-runner.sh`.

Set these environment variables only when your setup is different:

- `CODEX_WSL_EXE`: Windows path to `wsl.exe`.
- `CODEX_WSL_DISTRO`: WSL distro name, such as `Ubuntu-24.04`.
- `CODEX_WSL_USER`: WSL user to run as.
- `CODEX_WSL_PROXY`: WSL path to `codex-wsl-proxy-runner.sh`.
- `CODEX_WSL_PROXY_JS`: WSL path to `codex-wsl-proxy.js`.
- `CODEX_WSL_PROXY_DISTRO`: distro name to use when converting WSL paths to
  `\\wsl.localhost\...` paths. WSL usually sets this itself.
- `CODEX_WSL_PROXY_TARGET`: WSL path to the real `codex` binary.
- `CODEX_HOME`: Codex home directory inside WSL.

## Use

Put the full path to `codex-wsl.cmd` into T3code. Arguments are passed through to Codex after
Windows paths are converted to WSL paths.

The WSL side needs `node` available on `PATH`.
