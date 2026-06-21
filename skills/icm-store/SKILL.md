---
name: icm-store
description: ICM store for explicit non-secret durable Codex memory. Use when the user invokes `icm-store`, explicitly asks to remember/store/save/persist content in ICM, or explicitly marks a standing preference, resolved error, architecture decision, or stable project fact as memory for future Codex sessions. Not for combined recall-plus-store ICM tasks; use icm. Not for repo docs, Discord archives, current-chat notes, read-only/no-memory/no-store tasks, ambiguous private facts, or ordinary save/persist requests outside ICM.
---

# ICM Store

Store a concise safe summary of the user's provided content in ICM. Prefer the
ICM MCP store tool. Use topic `note` unless the user specifies a better topic.
If MCP is unavailable, use the local ICM CLI directly:

```bash
icm store -t "note" -c "<content>"
```

Never store secrets, tokens, passwords, recovery codes, private personal data, or raw session exports. Refuse or ask before storing ambiguous private content. Report the topic and stored summary in the final response.

Read-only, no-store, no-memory, blank-context, and repo-only constraints
prohibit storage even when `icm-store` is invoked. Ask the user to resolve
contradictory instructions instead of writing memory.
