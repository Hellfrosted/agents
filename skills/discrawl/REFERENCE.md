# Discrawl Reference

## Local Setup

- Binary: `/home/crunch/.local/bin/discrawl`
- Config: `/home/crunch/.config/discrawl/config.toml`
- Database: `/home/crunch/.local/share/discrawl/discrawl.db`
- Vesktop cache: `/mnt/c/Users/nguco/AppData/Roaming/Vesktop/sessionData`
- Default sync source: `wiretap`
- Discord token source: `none`
- Desktop full cache: `true`
- Embeddings: disabled by default
- Vector backend: `turbovec`
- TurboVec Python: `/home/crunch/.local/share/discrawl/turbovec-venv/bin/python`
- Setup doc: `/mnt/e/dev/dotfiles/docs/codex-workstation/discrawl-wiretap.md`

## Troubleshooting

For setup or troubleshooting, run:

```bash
discrawl --version
discrawl check-update
discrawl doctor --json
discrawl status --json
rg -n '^(token_source|source|path|full_cache|auto_update|media|enabled|provider|model|api_key_env|batch_size|vector_backend) =' /home/crunch/.config/discrawl/config.toml
printf '%s\n' "$CRAWLKIT_TURBOVEC_PYTHON"
"$CRAWLKIT_TURBOVEC_PYTHON" -E -c 'import turbovec, numpy; print("turbovec_import=ok")'
```

The expected setup is wiretap-only with Vesktop's `sessionData` path. If the
cache path stops working, inspect the Vesktop data directories under
`/mnt/c/Users/nguco/AppData/Roaming/`.

Discrawl `0.11.0` reports `"vector": "not configured"` from
`discrawl doctor --json` even when `[search.embeddings].vector_backend` is set
to `turbovec`; the doctor command does not probe the vector bridge. Verify the
bridge with the Python import command above or with an end-to-end semantic
search against a test database.

## Limits

- The archive only contains what Vesktop has cached locally.
- Full-cache wiretap is the current local default and can be slower than a
  focused scan.
- The `turbovec` backend only matters when embeddings are enabled and semantic
  or hybrid search is used. Real archive embeddings remain disabled unless the
  user explicitly opts in.
- Bot-visible complete history is intentionally not configured.
