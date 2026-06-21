---
name: icm
description: ICM task-level memory workflow for recall plus store across Codex sessions. Use when the user asks to use ICM generally or explicitly asks ICM to both recall past memory and store the result or decision in the same task. Not for one-off recall/store; use icm-recall or icm-store. Not for repo search, Discord/archive history, or no-memory/blank-context/read-only tasks.
---

# ICM

Use ICM as Codex persistent memory when durable context is likely to matter.
Do not treat this Skill as a universal start/end hook.

## Recall First

At task start, search relevant past context only when the active request makes
past-session memory useful. Prefer the ICM MCP recall tool. If MCP is
unavailable, use the local ICM CLI through `rtk run`:

```bash
rtk run 'icm recall --read-only "query"'
```

Use non-read-only recall only when updating memory access bookkeeping is
intentionally in scope.

## Store

Before replying, store privacy-safe durable facts only when the user explicitly
asks to remember/store them, or ask first when a fact seems durable but the user
did not request storage:

- Resolved error: topic `errors-resolved`, importance `high`.
- Architecture decision: topic `decisions-{project}`, importance `high`.
- User preference: topic `preferences`, importance `critical`.
- Significant progress or completion: topic `context-{project}`, importance `high`.
- Long run with meaningful reusable progress: save a concise progress summary.

Prefer ICM MCP store when available. If MCP is unavailable, use the local ICM
CLI through `rtk run`:

```bash
rtk run 'icm store -t "topic" -c "summary" -i high'
```

Read-only permits recall only. No-store, no-memory, blank-context, and
repo-only constraints prohibit storage even if ICM is named. Ask before storing
ambiguous personal or private facts. Never store secrets, tokens, passwords,
recovery codes, private personal data, or raw session exports. Keep ICM
Codex-facing unless the user says otherwise; do not add other harness setup or
touch install directories for ICM cleanup. If you store anything, report the
topic and safe summary in the final response.
