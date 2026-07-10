# Discrawl Reference

## Local Setup

- Binary: discover with `command -v discrawl`. Official distribution is a Go
  binary installed with Homebrew or built from source; local/manual installs
  commonly place it under the user's local bin directory.
- Config/database: read the resolved paths from `discrawl status --json`; do not
  assume a particular home directory. Existing installs may retain legacy
  `~/.discrawl` paths, while current installs use OS-native config/data paths.
- Vesktop cache: read the configured path from the config reported above. On
  Windows/WSL it is commonly below the Windows user's
  `AppData/Roaming/Vesktop/sessionData`, but the username and mount root vary.
- Default sync source: `wiretap`
- Discord token source: `none`
- Desktop full cache: `true`
- Embeddings: disabled by default
- Vector backend: `turbovec`
- TurboVec Python: read `CRAWLKIT_TURBOVEC_PYTHON`; if unset, report that the
  bridge path is not configured instead of guessing it.
- Workstation setup doc: check `$DOTFILES_ROOT` when set, then readable
  `docs/codex-workstation/discrawl-wiretap.md` candidates beneath the active
  workspace or known dotfiles checkout. Do not crawl entire mounted drives or
  assume one workstation path.

## Troubleshooting

For setup or troubleshooting, run:

```bash
discrawl --version
discrawl doctor --json
discrawl status --json
config_path="$(discrawl status --json | jq -r '.config_path // empty')"
test -n "$config_path" && test -r "$config_path" &&
  rg -n '^(token_source|source|path|full_cache|auto_update|media|enabled|provider|model|api_key_env|batch_size|vector_backend) =' "$config_path"
printf '%s\n' "$CRAWLKIT_TURBOVEC_PYTHON"
test -n "$CRAWLKIT_TURBOVEC_PYTHON" &&
  test -x "$CRAWLKIT_TURBOVEC_PYTHON" &&
  "$CRAWLKIT_TURBOVEC_PYTHON" -E -c 'import turbovec, numpy; print("turbovec_import=ok")'
```

Optional network-approved update check:

```bash
discrawl check-update
```

The expected setup is wiretap-only with Vesktop's `sessionData` path. If the
cache path stops working, inspect the configured Windows user's Vesktop data
directory rather than assuming a workstation username.

Discrawl `0.11.4` still sets `"vector": "not configured"` unconditionally in
the upstream doctor implementation, even when
`[search.embeddings].vector_backend` is `turbovec`; doctor does not probe the
vector bridge. Verify the bridge with the guarded Python import above or an
end-to-end semantic search against a test database. Recheck upstream source
before carrying this note into a later release line.

## Limits

- The archive only contains what Vesktop has cached locally.
- Full-cache wiretap is the current local default and can be slower than a
  focused scan.
- The `turbovec` backend only matters when embeddings are enabled and semantic
  or hybrid search is used. Real archive embeddings remain disabled unless the
  user explicitly opts in.
- Bot-visible complete history is intentionally not configured.
