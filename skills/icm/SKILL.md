---
name: icm
description: Provides the ICM (Infinite Context Memory) persistent-memory rule for Codex. Use when persistent memory should be consulted or maintained for a task, when the user asks to use ICM generally, or when durable context such as user preferences, resolved errors, architecture decisions, or significant project progress should persist across Codex sessions.
---

# ICM

Use ICM (Infinite Context Memory) as the persistent memory layer when durable context is likely to matter.

## Recall

At the start of each task, search for relevant past context.

```bash
rtk run 'icm recall "query"'
```

Prefer the ICM MCP recall tool when it is available in the Codex session; it accesses the same memory store.

## Store

Store when any of these happens.

1. Error resolved: `rtk run 'icm store -t errors-resolved -c "description" -i high'`
2. Architecture decision: `rtk run 'icm store -t decisions-{project} -c "description" -i high'`
3. User preference discovered: `rtk run 'icm store -t preferences -c "description" -i critical'`
4. Significant task completed: `rtk run 'icm store -t context-{project} -c "summary" -i high'`
5. Conversation exceeds about 20 tool calls without a store: store a progress summary.

Do this before responding to the user when the information is durable, useful later, and privacy-safe.

## Guardrails

Never store secrets, tokens, passwords, recovery codes, private personal data, or raw session exports.

When cleaning integrations, keep ICM Codex-facing unless the user says otherwise. Do not add ICM setup for other agent harnesses, and do not touch agent install directories for this purpose.
