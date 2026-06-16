# Discrawl Reference

## Local Setup

- Binary: `/home/crunch/.local/bin/discrawl`
- Config: `/home/crunch/.config/discrawl/config.toml`
- Database: `/home/crunch/.local/share/discrawl/discrawl.db`
- Vesktop cache: `/mnt/c/Users/nguco/AppData/Roaming/Vesktop/sessionData`
- Default sync source: `wiretap`
- Discord token source: `none`
- Desktop full cache: `true`
- Setup doc: `/mnt/e/dev/agents-toolkit/docs/discrawl-wiretap.md`

## Troubleshooting

For setup or troubleshooting, run:

```bash
discrawl --version
discrawl check-update
discrawl doctor --json
discrawl status --json
rg -n '^(token_source|source|path|full_cache|auto_update|media) =' /home/crunch/.config/discrawl/config.toml
```

The expected setup is wiretap-only with Vesktop's `sessionData` path. If the
cache path stops working, inspect the Vesktop data directories under
`/mnt/c/Users/nguco/AppData/Roaming/`.

## Limits

- The archive only contains what Vesktop has cached locally.
- Full-cache wiretap is the current local default and can be slower than a
  focused scan.
- Bot-visible complete history is intentionally not configured.
