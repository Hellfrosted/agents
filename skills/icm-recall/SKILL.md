---
name: icm-recall
description: ICM recall for past-session or project memory. Use when the user invokes `icm-recall`, asks what ICM remembers, asks what we decided or learned in a prior Codex session, or asks to search/recall long-term Codex memory. Not for combined recall-plus-store ICM tasks; use icm. Not for repo files, web sources, Discord archives, or current-chat context.
---

# ICM Recall

Search ICM for the user's query and report only relevant memories concisely.
Prefer the ICM MCP recall tool. If MCP is unavailable, use the local ICM CLI
directly:

```bash
icm recall --read-only "<query>"
```

Use non-read-only recall only when updating memory access bookkeeping is
intentionally in scope.

Do not reveal secrets, credentials, private personal data, or raw session
exports from memory. Summarize only non-sensitive facts relevant to the current
request.
