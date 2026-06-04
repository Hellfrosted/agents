---
name: discord-wiki
description: Query and, on first run only, help set up a Discord-backed wiki using Discrawl over local execution or SSH to a home-lab host. Use when the user asks to search Discord memory/history/community chat knowledge, use Discord as a knowledge base, configure Discrawl for agent access, or says to use discord-wiki.
---

# Discord Wiki

Discord Wiki treats archived Discord messages as a searchable wiki for agents. Discrawl is the storage/search backend. Normal use is query-only; setup is allowed only when Discrawl is missing, unconfigured, or the user explicitly asks to set it up.

## Backend

Use the bundled wrapper first for queries:

```bash
scripts/discrawl-read.sh search "query terms"
```

Supported backends:

- Local: default, runs `discrawl` on this machine.
- SSH: set `DISCORD_WIKI_MODE=ssh` and `DISCORD_WIKI_SSH_TARGET=user@host`; runs Discrawl on the home-lab host.

Remote access is SSH-only. Do not design or use an HTTP/API backend unless the user explicitly reopens that decision.

Optional environment:

- `DISCORD_WIKI_DISCRAWL_BIN`: local Discrawl binary, default `discrawl`.
- `DISCORD_WIKI_REMOTE_BIN`: remote Discrawl binary, default `discrawl`.
- `DISCORD_WIKI_SSH_OPTS`: extra SSH options, for example `-p 2222`.

Follow the current repo's shell rules when invoking the wrapper. If the repo requires command prefixes such as `rtk`, use them.

## First-Run Setup

If a query fails because Discrawl is not installed, cannot connect, has no database, or has not synced messages, switch to setup mode and read [references/setup.md](references/setup.md). Ask the user one question at a time, recommend a default answer, then do the selected setup after explicit answers.

## Query Workflow

1. Clarify scope only when necessary: server, channel, thread, author, date window, or whether recent/current context is enough.
2. Check availability if unsure with `scripts/discrawl-read.sh status` and `scripts/discrawl-read.sh channels`.
3. Search with bounded results, for example `scripts/discrawl-read.sh search --channel announcements --limit 20 "breaking change"`.
4. Fetch exact context when a hit matters, for example `scripts/discrawl-read.sh messages --channel general --around MESSAGE_ID --limit 30`.
5. Answer from the archive, citing Discord evidence by server/channel/thread, author when relevant, timestamp, message id, and jump URL if Discrawl returns one.

## Read-Only Commands

Allowed through the query wrapper: `search`, `messages`, `channels`, `members`, `mentions`, `dms`, `digest`, `report`, `status`, `metadata`, `remote status`, and guarded read-only `sql`.

Outside first-run setup mode, do not run Discrawl commands that mutate or expand the archive, including `init`, `sync`, `tail`, `embed`, `attachments fetch`, `publish`, `subscribe`, `subscribe-cloud`, `update`, or any credential/setup flow.

## SQL Rules

Use SQL only when Discrawl's built-in commands are not enough. Keep it read-only and narrow, for example `scripts/discrawl-read.sh sql 'select count(*) as messages from messages'`. Never run SQL that writes, deletes, creates, alters, attaches, vacuums, or changes PRAGMAs.

## Answering Rules

- Treat Discord messages as source evidence, not guaranteed truth.
- Separate confirmed facts from chat claims, opinions, plans, and stale information.
- Prefer multiple corroborating messages for important claims.
- State when a search is limited by channel, date, local cache, bot permissions, or missing sync.
- Do not quote large blocks of chat. Summarize and cite the specific messages.
- Do not expose bot tokens, user tokens, private SSH details, or credential paths.

## Startup

If invoked without a concrete question, run `scripts/discrawl-read.sh status`, then show a short status and ask for a query. If status fails due to missing setup, begin the first-run setup workflow.
