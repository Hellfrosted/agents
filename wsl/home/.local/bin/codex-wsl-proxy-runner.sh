#!/usr/bin/env bash
set -euo pipefail

export HOME=/home/crunch
export USER=crunch

NODE_ALIAS_FILE="$HOME/.nvm/alias/default"
if [ ! -r "$NODE_ALIAS_FILE" ]; then
  printf '{"error":"Failed to read NVM default alias at %s"}\n' "$NODE_ALIAS_FILE"
  exit 1
fi

NODE_VERSION="$(tr -d '\r\n' < "$NODE_ALIAS_FILE")"
NODE_BIN="$HOME/.nvm/versions/node/v$NODE_VERSION/bin/node"

if [ ! -x "$NODE_BIN" ]; then
  printf '{"error":"Failed to find node binary at %s"}\n' "$NODE_BIN"
  exit 1
fi

exec "$NODE_BIN" /home/crunch/.local/bin/codex-wsl-proxy.js "$@"
