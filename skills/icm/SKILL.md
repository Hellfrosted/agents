---
name: icm
description: "ICM persistent memory for Codex. Use for three branches only: recall long-term memory, store explicit durable non-secret memory, or recall-then-store when the user asks for both. Not for repo files, web sources, Discord archives, current-chat notes, ordinary persistence outside ICM, or read-only/no-memory tasks."
---

# ICM

Use ICM only as Codex persistent memory. Pick exactly one branch: recall,
store, or recall-then-store. Do not use ICM as a universal start/end hook.

## Recall

Use when the user asks what ICM remembers, asks about a prior Codex session, or
explicitly makes long-term memory useful for the active task.

1. Search only for the user's query. Prefer the ICM MCP recall tool; if MCP is
   unavailable, run:

   ```bash
   icm recall --read-only "<query>"
   ```

2. Complete when the final answer reports relevant non-sensitive memories, or
   says no relevant ICM memory was found. Do not present repo files, web
   sources, Discord archives, pasted chat, or current-chat context as ICM
   memory.

Use non-read-only recall only when updating memory access bookkeeping is
intentionally in scope.

## Store

Use when the user explicitly asks to remember, store, save, or persist content
in ICM, or marks a stable preference, resolved error, architecture decision, or
project fact as future Codex memory.

1. Store a concise safe summary of the user-provided durable fact. Ask before
   storing inferred, ambiguous, private, or potentially secret material.
2. Choose topic and importance:
   - User preference: `preferences`, importance `critical`.
   - Resolved error: `errors-resolved`, importance `high`.
   - Architecture decision: `decisions-{project}`, importance `high`.
   - Significant progress or stable project context: `context-{project}`,
     importance `high`.
   - Otherwise: `note`, importance `high`.
3. Prefer the ICM MCP store tool; if MCP is unavailable, run:

   ```bash
   icm store -t "topic" -c "summary" -i high
   ```

4. Complete when the final answer reports the topic and stored summary, without
   exposing sensitive raw content.

## Recall Then Store

Use when the user asks to use ICM generally for a task, or explicitly asks for
both memory lookup and durable storage.

1. Recall first, using the recall branch.
2. Do the requested task.
3. Store only explicit safe durable facts, or ask before storing inferred facts.
4. Complete when the final answer separates memory used from memory stored; if
   nothing safe and durable was stored, say so.

## Safety Gates

Read-only permits recall only. No-store, no-memory, blank-context, and
repo-only constraints prohibit storage even when ICM is named. Ask the user to
resolve contradictory instructions before writing memory.

Never store or reveal secrets, tokens, passwords, recovery codes, private
personal data, or raw session exports. Keep ICM Codex-facing unless the user
says otherwise.
