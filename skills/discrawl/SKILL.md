---
name: discrawl
description: Discrawl local Discord archive search. Use for local Discord history questions, server/channel/DM cache lookup, Vesktop cache refresh/status, or bare "use discrawl" status. Not for official Discord/API docs, web/Slack search, pasted chat, live Discord actions, or ICM memory.
---

# Discrawl

Discrawl searches the workstation's local Discord/Vesktop archive. It reads the
desktop cache only; do not configure bot tokens, user tokens, selfbots, cloud
sync, or remote publishing for normal searches.

For paths, troubleshooting commands, and cache limits, see
[REFERENCE.md](REFERENCE.md).

When this workflow uses Codex subagents, follow the delegation, evidence, and
bounded-loop rules in
[`../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol`](../shared-agent-protocol/SKILL.md#codex-delegation-and-reviewer-protocol),
[`../shared-agent-protocol/SKILL.md#evidence-backed-output`](../shared-agent-protocol/SKILL.md#evidence-backed-output),
and
[`../shared-agent-protocol/SKILL.md#bounded-loops`](../shared-agent-protocol/SKILL.md#bounded-loops).

## Branches

### Search

Use when the user asks what Discord history says about a topic.

1. Refresh with `discrawl sync` unless the user asks for read-only,
   no-refresh, offline, or no-network behavior.
2. If the channel/server is ambiguous, check `discrawl --version`. On 0.11.4 or
   newer, resolve the name directly with `discrawl channels resolve "name"
   --json`, adding `--guild <id>` when known. On older versions, or when direct
   resolution is unavailable, filter and bound `discrawl --json channels list`
   locally before its output enters context, for example `discrawl --json
   channels list | jq --arg q "$term" '[.[] | select(.name | ascii_downcase |
   contains($q | ascii_downcase))][0:20]'`; never paste the unfiltered
   directory.
3. Run `discrawl search "query"` with useful `--guild`, `--channel`,
   `--author`, or `--limit` filters. Use `--dm` only when the user explicitly
   asks to search DMs.
4. Complete when the answer includes auditable evidence: command, filters,
   channel/guild name or id when available, timestamp, and enough context to
   support each claim. If nothing matches, report no results with the exact
   command and filters used.

### Channel Context

Use when the user asks what a channel or DM said near a time or topic.

1. Resolve a server channel with the version-gated direct resolver or bounded
   local-filter fallback from the Search branch. For a DM, use `discrawl dms
   --with "person" --last N` or resolve it implicitly with `discrawl messages
   --dm --channel "person" --last N`; the `--dm` scope is required.
2. Inspect a bounded window with `discrawl messages --channel "name" --last N`,
   `--days N`, `--since`, or `--before`. Keep `--dm` on every DM message query.
3. Complete when the summary names the channel, window, timestamps, and the
   message context relied on; no long private dumps.

### Status

Use for bare invocation or cache-health questions.

Run `discrawl status --json`; add `discrawl doctor --json` only for setup or
health checks. Complete by reporting message count, channel count, database
path, and any relevant health issue.

## Safety

- Treat `remote`, `cloud`, and `subscribe-cloud` as opt-in publishing or
  remote-read features.
- Treat `check-update` as optional network access.
- Never record token values in docs, examples, commits, or agent notes.
- Summarize only what is needed from private conversations.
