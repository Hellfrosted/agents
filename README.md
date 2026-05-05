# codex-wsl

Windows and WSL launcher support files for running Codex from Windows into WSL.

## Layout

- `windows/bin/codex-wsl.cmd` mirrors `C:\Users\nguco\bin\codex-wsl.cmd`.
- `windows/.codex/skills/goal-maker/` mirrors `C:\Users\nguco\.codex\skills\goal-maker`.
- `wsl/home/.codex/skills/` mirrors `/home/crunch/.codex/skills`, excluding `.system`.
- `wsl/home/.local/bin/codex-wsl-proxy-runner.sh` mirrors `/home/crunch/.local/bin/codex-wsl-proxy-runner.sh`.
- `wsl/home/.local/bin/codex-wsl-proxy.js` mirrors `/home/crunch/.local/bin/codex-wsl-proxy.js`.
