#!/usr/bin/env bash
':' //; [ -n "${HOME:-}" ] || export HOME="$(getent passwd "$(id -u)" | cut -d: -f6)"; [ -n "${USER:-}" ] || export USER="$(id -un)"; export PATH="$HOME/.local/share/pnpm/bin:$HOME/.local/share/pnpm:$HOME/.bun/bin:$HOME/.local/bin:$HOME/bin:${PATH:-/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin}"; command -v node >/dev/null 2>&1 || { printf '{"error": "Failed to find node in WSL"}\n'; exit 1; }; exec node "$0" "$@"

const { startProxy } = require("./codex-wsl-proxy-runtime");

startProxy();
