# Discord Wiki Setup

Use this only when Discrawl is missing, unreachable, unconfigured, empty, or the user explicitly asks to set up Discord Wiki.

## Reliability Model

The desired system is not live scraping during agent use. The skill should query a database that a home-lab indexer has already built.

Target flow:

```text
user account / Discord Desktop cache identifies useful servers and channels
  -> home-lab indexing job records those targets
  -> async Discrawl jobs backfill history slowly
  -> scheduled refresh keeps the DB warm
  -> agents query the already-built SQLite/FTS database over SSH
```

First setup should prefer registering an indexing target and starting or checking the async backfill, not blocking an agent session until all history is indexed.

Important constraint: "index every message ever sent" requires an indexing source that can page channel history for that server/channel. Wiretap can identify and seed targets from cached account activity, but cache data is not complete history. For complete server backfill, use a source with actual history access for that server, such as bot sync where feasible, or another explicitly chosen collector.

Non-admin compatibility is a hard requirement. The default design must work for servers where the user is only a normal member and cannot add a bot, change server settings, or get admin exports. Index only channels/messages the user's account can legitimately see. Do not suggest bypassing private channels, hidden history, role restrictions, or server permissions.

## Setup Interview

Ask one question at a time. Include the recommended answer in each question.

Decision order:

1. Backend location: recommend SSH home-lab if the user has one; otherwise local.
2. SSH target, if remote: ask for `user@host` and optional SSH options.
3. Target discovery: recommend user-account/cache-based discovery to identify which servers/channels should be indexed.
4. Backfill source: ask how the home-lab indexer should obtain history for the selected targets without requiring admin rights; recommend bot sync only where the user has permission to add one or the server provides an approved export path.
5. Storage scope: recommend text, metadata, channels, threads, members, FTS; no media and no embeddings for the first pass.
6. Async mode: recommend a home-lab scheduled job or service that backfills slowly and refreshes regularly.
7. Query readiness: ask whether the skill should return partial results while backfill is still running.

## Credential Rule

Keep credentials user-controlled. Never ask the user to paste bot tokens, user tokens, SSH keys, passwords, recovery codes, or private connection details into chat. Tell them to create/store credentials directly on the target machine, then run or help run Discrawl commands there.

## Target Discovery

Use wiretap as a discovery/seed path, not as the reliability guarantee:

```bash
discrawl init
discrawl doctor
discrawl wiretap --dry-run
discrawl sync --source wiretap
discrawl search "test query"
```

Wiretap imports classifiable messages from Discord Desktop cache and avoids per-server bot setup. It is good for identifying active servers/channels and seeding the DB. It only sees what the Desktop client has cached and can classify; it is not guaranteed to be complete server or DM history.

For non-admin servers, discovery and indexing should be based on member-visible account activity, cache-derived targets, and any approved export method available to a normal member. If the chosen collector cannot read older channel history, label coverage as partial rather than claiming full backfill.

## Headless Server Mode

The home lab should own the durable DB and async jobs. Wiretap on a headless server is functional only if the Discord Desktop cache is available on that server. Recommended options:

- Sync or mount the Discord Desktop cache directory onto the home-lab host, then run wiretap with `--path` for target discovery and seed data.
- Run target discovery locally on the workstation, then copy only target metadata to the home lab.
- Use bot sync or another explicit complete-history collector only for selected servers/channels where full backfill matters and the user has permission to use that source.

Do not treat wiretap cache import as "everything ever sent." It is a useful indicator for what to index; complete history needs a separate backfill source.

## Allowed Setup Actions

After explicit user answers, acceptable setup actions are:

- Check Discrawl availability: `discrawl doctor`, `discrawl status`, `discrawl metadata --json`.
- Install or update Discrawl only if the user explicitly approves the install method.
- Run `discrawl init` on the selected backend.
- Help the user locate or mount the Discord Desktop cache path without exposing secrets.
- Help the user configure a bot token on the target machine only if they explicitly choose bot sync and have permission to add/use a bot for that server.
- Register indexing targets and run the selected first sync/backfill command.
- Help create or check a scheduled home-lab indexing job after user approval.
- Verify with `discrawl status`, `discrawl channels`, and a small `discrawl search`.

Do not enable media downloads, embeddings, Git snapshots, cloud subscribe, or publishing during first setup unless the user explicitly asks for those features.

## CLI Client Fallback

Do not replace Discrawl with a Discord TUI/CLI client for the stored knowledge-base path. Tools such as Discordo or discord-cli can be useful for live listening or narrow recent-history experiments, but they do not replace Discrawl's SQLite archive, FTS index, channel/thread/member model, status checks, or read-only SQL surface.

If Discrawl cannot be used, offer a separate live-only design instead of silently changing this skill's backend. That design should be explicit about losing durable storage, broad search, and historical indexing.

## Recommended First Pass

Prefer:

- SSH home-lab backend.
- Wiretap or user-account/cache-derived data to identify what should be indexed.
- Async home-lab indexing/backfill jobs so searches use an already-built DB.
- Text and metadata only.
- FTS search enabled by Discrawl.
- No media downloads.
- No embeddings.
- Non-admin/member-visible server support by default.
- Complete-history backfill only for selected servers/channels where the user chooses a source that can actually read full history without bypassing permissions.

The goal is a working, cheap, searchable Discord knowledge base before adding semantic search or attachment storage.
