#!/usr/bin/env bash
set -euo pipefail

if [ -z "${HOME:-}" ]; then
  HOME="$(getent passwd "$(id -u)" | cut -d: -f6)"
  export HOME
fi

if [ -z "${USER:-}" ]; then
  USER="$(id -un)"
  export USER
fi

export CODEX_HOME="${CODEX_HOME:-$HOME/.codex}"

DEFAULT_PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
PATH_PREFIX="$HOME/.local/bin:$HOME/bin"
NODE_ALIAS_FILE="$HOME/.nvm/alias/default"

if [ -r "$NODE_ALIAS_FILE" ]; then
  NODE_VERSION="$(tr -d '\r\n' < "$NODE_ALIAS_FILE")"
  if [ -n "$NODE_VERSION" ]; then
    PATH_PREFIX="$PATH_PREFIX:$HOME/.nvm/versions/node/v$NODE_VERSION/bin"
  fi
fi

export PATH="$PATH_PREFIX:${PATH:-$DEFAULT_PATH}"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROXY_JS="${CODEX_WSL_PROXY_JS:-$SCRIPT_DIR/codex-wsl-proxy.js}"

if ! command -v node >/dev/null 2>&1 && [ -s "$HOME/.nvm/nvm.sh" ]; then
  . "$HOME/.nvm/nvm.sh" >/dev/null 2>&1
  if ! command -v node >/dev/null 2>&1 && command -v nvm >/dev/null 2>&1; then
    nvm use --silent default >/dev/null 2>&1 || true
  fi
fi

if ! command -v node >/dev/null 2>&1; then
  printf '{"error":"Failed to find node in WSL. Install node or make it available through nvm."}\n'
  exit 1
fi

if [ ! -r "$PROXY_JS" ]; then
  printf '{"error":"Failed to read Codex WSL proxy at %s"}\n' "$PROXY_JS"
  exit 1
fi

exec node "$PROXY_JS" "$@"
