---
name: icm-store
description: Stores information in ICM persistent memory from Codex. Use when the user invokes `icm-store`, asks to store/save a note in ICM, or provides durable context that should be kept for future sessions.
---

# ICM Store

Store the user's provided content in ICM. Prefer the ICM MCP store tool. Use topic `note` unless the user specifies a better topic. CLI fallback:

```bash
rtk run 'icm store -t "note" -c "<content>"'
```

Never store secrets, tokens, passwords, recovery codes, private personal data, or raw session exports.
