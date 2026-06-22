---
name: icm
description: ICM persistent memory for Codex. Use when the user invokes `icm`, asks what ICM remembers, asks to recall/search long-term memory, asks to remember/store/save durable non-secret context, or asks to recall and store in one task. Not for repo files, web sources, Discord archives, current-chat notes, ordinary save/persist outside ICM, or no-memory/no-store/read-only tasks.
---

# ICM

Use ICM as Codex persistent memory. Pick exactly the branch the request needs:
recall, store, or recall-then-store. Do not treat ICM as a universal start/end
hook.

## Recall

Use when the user asks what ICM remembers, asks what was decided or learned in a
prior Codex session, asks to search long-term memory, or when the active task
explicitly makes past-session memory useful.

Search only for the user's query and report relevant non-sensitive memories
concisely. Prefer the ICM MCP recall tool. If MCP is unavailable, use:

```bash
icm recall --read-only "<query>"
```

Use non-read-only recall only when updating memory access bookkeeping is
intentionally in scope. Do not use ICM for repo files, web sources, Discord
archives, pasted chat, or current-chat context.

## Store

Use when the user explicitly asks to remember, store, save, or persist content
in ICM, or explicitly marks a standing preference, resolved error, architecture
decision, or stable project fact as future Codex memory.

Store a concise safe summary of the user-provided content. Use topic `note`
unless a better topic is clear:

- Resolved error: `errors-resolved`, importance `high`.
- Architecture decision: `decisions-{project}`, importance `high`.
- User preference: `preferences`, importance `critical`.
- Significant progress or completion: `context-{project}`, importance `high`.

Prefer the ICM MCP store tool. If MCP is unavailable, use:

```bash
icm store -t "topic" -c "summary" -i high
```

Report the topic and stored summary in the final response.

## Recall Then Store

Use when the user asks to use ICM generally for a task or explicitly asks for
both memory lookup and durable storage. Recall first, complete the task, then
store only explicit safe durable facts or ask before storing inferred durable
facts.

## Safety Gates

Read-only permits recall only. No-store, no-memory, blank-context, and repo-only
constraints prohibit storage even when ICM is named. Ask the user to resolve
contradictory instructions before writing memory.

Never store or reveal secrets, tokens, passwords, recovery codes, private
personal data, or raw session exports. Refuse or ask before storing ambiguous
private content. Keep ICM Codex-facing unless the user says otherwise; do not
add harness setup or touch install directories for ICM cleanup.
