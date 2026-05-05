---
name: icm-recall
description: Searches ICM persistent memory from Codex. Use when the user invokes `icm-recall`, asks to recall or search ICM memory, asks what ICM remembers, or provides a query that should be looked up in long-term memory.
---

# ICM Recall

Search ICM memory for the user's query.

If using the CLI, replace `<query>` with the user's query and run:

```bash
rtk run 'icm recall "<query>"'
```

When the ICM MCP recall tool is available, use it instead of the CLI. Preserve the same behavior: search for the user's provided query and report the relevant memory results concisely.
