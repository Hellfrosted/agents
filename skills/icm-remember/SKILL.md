---
name: icm-remember
description: Stores information in ICM persistent memory from Codex. Use when the user invokes `icm-remember`, asks to remember something, asks to store/save a note in ICM, or provides durable context that should be kept for future sessions.
---

# ICM Remember

Store the user's provided content in ICM memory.

If using the CLI, replace `<content>` with the user's content and run:

```bash
rtk run 'icm store -t "note" -c "<content>"'
```

When the ICM MCP store tool is available, use it instead of the CLI. Preserve the same default behavior: store the user's provided content under the `note` topic unless the user specifies a more precise topic.

Never store secrets, tokens, passwords, recovery codes, private personal data, or raw session exports.
