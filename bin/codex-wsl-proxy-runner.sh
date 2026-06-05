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

# LazyCodex works as a Codex plugin, but its startup hooks can emit anonymous
# telemetry, spawn a detached self-updater, and migrate Codex config. The shim
# disables those network/background mutation paths for T3code's long-lived
# app-server sessions while leaving the plugin itself enabled.
export LAZYCODEX_AUTO_UPDATE_DISABLED="${LAZYCODEX_AUTO_UPDATE_DISABLED:-1}"
export OMO_CODEX_AUTO_UPDATE_DISABLED="${OMO_CODEX_AUTO_UPDATE_DISABLED:-1}"
export LAZYCODEX_CONFIG_MIGRATION_DISABLED="${LAZYCODEX_CONFIG_MIGRATION_DISABLED:-1}"
export OMO_CODEX_CONFIG_MIGRATION_DISABLED="${OMO_CODEX_CONFIG_MIGRATION_DISABLED:-1}"
export OMO_CODEX_DISABLE_POSTHOG="${OMO_CODEX_DISABLE_POSTHOG:-1}"
export OMO_CODEX_SEND_ANONYMOUS_TELEMETRY="${OMO_CODEX_SEND_ANONYMOUS_TELEMETRY:-0}"
export OMO_DISABLE_POSTHOG="${OMO_DISABLE_POSTHOG:-1}"
export OMO_SEND_ANONYMOUS_TELEMETRY="${OMO_SEND_ANONYMOUS_TELEMETRY:-0}"

DEFAULT_PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
PATH_PREFIX="$HOME/.local/share/pnpm/bin:$HOME/.local/share/pnpm:$HOME/.bun/bin:$HOME/.local/bin:$HOME/bin"

export PATH="$PATH_PREFIX:${PATH:-$DEFAULT_PATH}"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROXY_JS="${CODEX_WSL_PROXY_JS:-$SCRIPT_DIR/codex-wsl-proxy.js}"

if ! command -v node >/dev/null 2>&1; then
  printf '{"error":"Failed to find node in WSL. Install node or add it to PATH."}\n'
  exit 1
fi

if [ ! -r "$PROXY_JS" ]; then
  printf '{"error":"Failed to read Codex WSL proxy at %s"}\n' "$PROXY_JS"
  exit 1
fi

exec node "$PROXY_JS" "$@"
