#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  discrawl-read.sh <read-command> [args...]

Backends:
  local: default, runs $DISCORD_WIKI_DISCRAWL_BIN or discrawl
  ssh:   set DISCORD_WIKI_MODE=ssh and DISCORD_WIKI_SSH_TARGET=user@host

Allowed commands:
  search messages channels members mentions dms digest report status metadata remote sql
USAGE
}

if [[ $# -lt 1 ]]; then
  usage
  exit 2
fi

cmd="$1"
shift

case "$cmd" in
  search|messages|channels|members|mentions|dms|digest|report|status|metadata|remote|sql)
    ;;
  *)
    echo "discord-wiki: refused non-read Discrawl command: $cmd" >&2
    usage >&2
    exit 2
    ;;
esac

if [[ "$cmd" == "sql" ]]; then
  query="${1:-}"
  if [[ -z "$query" ]]; then
    echo "discord-wiki: sql requires a query string" >&2
    exit 2
  fi

  normalized="$(printf '%s' "$query" | tr '\n' ' ' | sed -E 's/^[[:space:]]+//; s/[[:space:]]+$//')"
  lower="$(printf '%s' "$normalized" | tr '[:upper:]' '[:lower:]')"

  if ! [[ "$lower" =~ ^(select|with|pragma[[:space:]]+(user_version|table_info|index_list|index_info|database_list|foreign_key_list|integrity_check|quick_check)\b) ]]; then
    echo "discord-wiki: refused SQL that does not start with SELECT, WITH, or a read-only PRAGMA" >&2
    exit 2
  fi

  if [[ "$lower" =~ (^|[^a-z])(insert|update|delete|drop|alter|create|replace|attach|detach|vacuum|reindex|analyze|begin|commit|rollback|savepoint|release)([^a-z]|$) ]]; then
    echo "discord-wiki: refused SQL containing a write or database-control keyword" >&2
    exit 2
  fi
fi

if [[ "$cmd" == "remote" ]]; then
  subcmd="${1:-status}"
  if [[ "$subcmd" != "status" ]]; then
    echo "discord-wiki: refused non-status Discrawl remote command: remote $subcmd" >&2
    exit 2
  fi
fi

mode="${DISCORD_WIKI_MODE:-local}"

case "$mode" in
  local|"")
    bin="${DISCORD_WIKI_DISCRAWL_BIN:-discrawl}"
    exec "$bin" "$cmd" "$@"
    ;;
  ssh)
    target="${DISCORD_WIKI_SSH_TARGET:-}"
    if [[ -z "$target" ]]; then
      echo "discord-wiki: DISCORD_WIKI_SSH_TARGET is required when DISCORD_WIKI_MODE=ssh" >&2
      exit 2
    fi
    remote_bin="${DISCORD_WIKI_REMOTE_BIN:-discrawl}"
    # shellcheck disable=SC2206
    ssh_opts=(${DISCORD_WIKI_SSH_OPTS:-})
    printf -v remote_command '%q ' "$remote_bin" "$cmd" "$@"
    exec ssh "${ssh_opts[@]}" "$target" bash -lc "$remote_command"
    ;;
  *)
    echo "discord-wiki: unknown DISCORD_WIKI_MODE: $mode" >&2
    exit 2
    ;;
esac
