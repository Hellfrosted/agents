# Discrawl Vesktop Wiretap

This workstation uses Discrawl as a local Discord cache reader for servers where
a bot cannot help. The setup is intentionally wiretap-only: it reads Vesktop's
local Electron cache, stores classifiable messages in a local SQLite archive,
and does not configure a Discord bot token.

## Installed Paths

- Binary: `/home/crunch/.local/bin/discrawl`
- Config: `/home/crunch/.config/discrawl/config.toml`
- Database: `/home/crunch/.local/share/discrawl/discrawl.db`
- Vesktop cache path:
  `/mnt/c/Users/nguco/AppData/Roaming/Vesktop/sessionData`

## Config Shape

The important settings are:

```toml
[discord]
token_source = "none"

[sync]
source = "wiretap"

[desktop]
path = "/mnt/c/Users/nguco/AppData/Roaming/Vesktop/sessionData"
full_cache = false

[share]
auto_update = false
media = false
```

`discrawl sync` therefore imports only the local Vesktop desktop cache. It does
not call the Discord API as a user, use a user token, run a selfbot, or publish
the archive.

## Normal Use

Refresh the local archive before searching:

```bash
discrawl sync
```

Inspect archive health:

```bash
discrawl check-update
discrawl doctor --json
discrawl status --json
discrawl channels list
```

Search cached server messages:

```bash
discrawl search "query"
discrawl search "query" --channel "channel-name"
discrawl messages --channel "channel-name" --last 50
```

DM search is available only for messages proven from local desktop cache:

```bash
discrawl search "query" --dm
discrawl messages --dm --last 50
```

## Limits

- Discrawl sees only messages Vesktop has cached locally.
- Servers and channels that have not been opened in Vesktop may not have useful
  cached history.
- `full_cache = false` keeps imports faster. Use `discrawl wiretap --full-cache`
  only for a deliberate slower archaeology pass.
- Do not add bot tokens unless the user explicitly wants Discord API sync; even
  then, never record token values in docs, examples, commits, or agent notes.
- Do not enable share/publish behavior unless the user explicitly wants to back
  up non-DM server archive data.
- Do not configure `remote`, `cloud`, or `subscribe-cloud` unless the user
  explicitly wants a Cloudflare-backed remote archive.

## Verification

The setup was verified on 2026-06-07:

```text
discrawl --version
0.10.0

discrawl check-update
discrawl: up to date (0.10.0)

discrawl doctor --json
discord_token: "discord token disabled by config"

discrawl sync --source wiretap
messages=26
guild_messages=26
dm_messages=0
dry_run=false
```
