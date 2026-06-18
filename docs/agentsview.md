# AgentsView

This workstation uses AgentsView as a local browser for Codex and OpenCode
session history. It indexes local session stores into a SQLite database and
serves the UI on a loopback-only URL.

## Installed Paths

- Binary: `/home/crunch/.local/bin/agentsview`
- Config: `/home/crunch/.agentsview/config.toml`
- Database: `/home/crunch/.agentsview/sessions.db`
- User systemd units:
  - `/home/crunch/.config/systemd/user/agentsview.socket`
  - `/home/crunch/.config/systemd/user/agentsview.service`
  - `/home/crunch/.config/systemd/user/agentsview-backend.service`

Last observed from this repo on 2026-06-18: AgentsView `v0.33.1`, database
`160M`, with `1723` Codex sessions and `86` OpenCode sessions indexed.

## Indexed Roots

AgentsView reads these Codex roots:

```toml
codex_sessions_dirs = [
  "/home/crunch/.codex/sessions",
  "/home/crunch/.codex/archived_sessions",
  "/mnt/c/Users/nguco/.codex/sessions",
]
```

It reads these OpenCode roots:

```toml
opencode_dirs = [
  "/mnt/c/Users/nguco/.local/share/opencode",
  "/mnt/c/Users/nguco/.local/share/OpenCode",
]
```

`/mnt/c/Users/nguco/AppData/Roaming/opencode` and
`/mnt/c/Users/nguco/AppData/Roaming/OpenCode` are WebView/cache roots on this
machine, not the OpenCode session store.

## Lazy Local URL

Open AgentsView at:

```text
http://127.0.0.1:18080
```

The URL is socket-activated by user-level systemd:

- `agentsview.socket` listens on `127.0.0.1:18080`.
- `agentsview.service` runs `systemd-socket-proxyd`.
- `agentsview-backend.service` starts AgentsView on `127.0.0.1:18081`.

Port `8080` is not used because it is already occupied by `tunnel-client`.

This setup starts AgentsView on the first URL hit only when WSL user systemd is
already running. A browser URL alone cannot cold-start the WSL distro.

## Normal Commands

Refresh the index:

```bash
agentsview sync
```

Inspect indexed projects:

```bash
agentsview projects
```

Check socket and service state:

```bash
systemctl --user status agentsview.socket agentsview.service agentsview-backend.service
```

Stop the running UI while leaving lazy activation enabled:

```bash
systemctl --user stop agentsview.service agentsview-backend.service
```

Disable the lazy URL:

```bash
systemctl --user disable --now agentsview.socket
```

Troubleshoot recent service logs:

```bash
journalctl --user -u agentsview.service -u agentsview-backend.service -n 80 --no-pager
```

Verify the URL:

```bash
curl -I http://127.0.0.1:18080/
```

## Limits

- Keep the URL loopback-only. Do not expose the session browser outside this
  workstation unless the user explicitly asks for a sharing design.
- Do not record raw session contents, tokens, or private data in repo docs.
- Docker is not the chosen path here because WSL can read every needed root
  directly and Docker would add bind-mount configuration without improving the
  local workflow.
- Portless is not used because the browser needs a stable local URL, and a
  user systemd socket already provides lazy startup with less moving state.
