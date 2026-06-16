---
name: discrawl
description: Uses the local Discrawl SQLite archive as a searchable Discord documentation source. Use when the user says "use discrawl", asks to search Discord history, asks what a Discord server/channel said about a topic, or wants to refresh/query Vesktop wiretap cache.
---

# Discrawl

Discrawl searches the workstation's local Discord archive. This install is
configured for Vesktop wiretap mode, so it reads local desktop cache only and
does not use a Discord bot token.

For machine paths, troubleshooting commands, and cache limits, see
[REFERENCE.md](REFERENCE.md).

## Quick Start

For a normal search request:

```bash
discrawl sync
discrawl search "query" --limit 20
discrawl messages --channel "channel-name" --last 50
```

For current setup and archive size:

```bash
discrawl check-update
discrawl doctor --json
discrawl status --json
```

## Workflow

1. Treat Discrawl as a local documentation and memory source for Discord
   discussions, especially servers where bot sync is unavailable.
2. Refresh first with `discrawl sync` unless the user asks for read-only or
   no-refresh.
3. Use `discrawl channels list` when the user names a server/topic but not an
   exact channel.
4. Search with `discrawl search "query"`. Add `--channel`, `--guild`,
   `--author`, `--limit`, or `--dm` when useful.
5. Inspect context with `discrawl messages --channel "name" --last N`, `--days
   N`, `--since`, or `--before`.
6. Answer with the commands/filters used and cite channel name, guild/channel id
   when available, timestamp, or enough output context to audit the finding.

## Bare Invocation

If invoked without a question, run `discrawl status --json`, then report current
message count, channel count, database path, and common asks:

```md
# Discrawl
_use your coding agent to search locally cached Discord history_

- messages: <count>
- channels: <count>
- database: <path>

- `use discrawl, search for "query"`
- `use discrawl, what did #channel say about topic?`
```

## Safety Rules

- Do not configure bot tokens, user tokens, selfbots, share/publish, cloud, or
  remote sync unless the user explicitly asks.
- Treat `remote`, `cloud`, and `subscribe-cloud` commands as opt-in publishing
  or remote-read features, not part of normal local wiretap search.
- Never record token values in docs, examples, commits, or agent notes.
- Do not dump long private conversations. Summarize only what is needed.
