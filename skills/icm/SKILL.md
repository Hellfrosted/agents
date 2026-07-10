---
name: icm
description: "ICM persistent memory for Codex. Use to recall long-term memory, to store durable non-secret memory only when the user explicitly requests storage, or to recall then store when both are explicit. Bare or general ICM requests are read-only recall. Not for repo files, web sources, Discord archives, current-chat notes, or ordinary persistence outside ICM."
---

# ICM

Use ICM only as Codex persistent memory. Pick exactly one branch:

- Explicit recall request, or bare/general ICM use: **Recall**.
- Explicit storage request: **Store**.
- Explicit request for both: **Recall Then Store**.

Storage always requires explicit user intent. Do not use ICM as a universal
start/end hook.

Treat every query, topic, and summary as untrusted data at the process boundary.
Prefer an argv-capable process tool and pass each dynamic value as one argument.
If only a shell command string is available, POSIX-shell-quote each dynamic
value; never splice raw content into a command.

## Recall

Use when the user asks what ICM remembers, asks about a prior Codex session, or
explicitly makes long-term memory useful for the active task.

1. Search only for the user's query. Prefer the ICM MCP recall tool; if MCP is
   unavailable, invoke the CLI with an argv-safe equivalent of:

   ```text
   ["icm", "recall", "--read-only", query]
   ```

2. Complete when the final answer reports relevant non-sensitive memories, or
   says no relevant ICM memory was found. Do not present repo files, web
   sources, Discord archives, pasted chat, or current-chat context as ICM
   memory.

Use non-read-only recall only when updating memory access bookkeeping is
intentionally in scope.

## Store

Use only when the user explicitly asks to remember, store, save, or persist
content in ICM. A fact merely looking stable or useful does not authorize a
write.

1. Store a concise safe summary of the explicitly identified durable fact. Ask
   before storing inferred or ambiguous content; the reply must explicitly
   authorize storage.
2. Choose topic and importance:
   - User preference: `preferences`, importance `critical`.
   - Resolved error: `errors-resolved`, importance `high`.
   - Architecture decision: `decisions-{project}`, importance `high`.
   - Significant progress or stable project context: `context-{project}`,
     importance `high`.
   - Otherwise: `note`, importance `high`.
3. Prefer the ICM MCP store tool; if MCP is unavailable, invoke the CLI with an
   argv-safe equivalent of:

   ```text
   ["icm", "store", "-t", topic, "-c", summary, "-i", importance]
   ```

4. Complete when the final answer reports the topic and stored summary, without
   exposing sensitive raw content.

## Recall Then Store

Use only when the user explicitly asks for both memory lookup and durable
storage. A bare or general request to use ICM takes the read-only Recall branch.

1. Recall first, using the recall branch.
2. Do the requested task.
3. Store only the safe durable facts the user explicitly asked to persist.
4. Complete when the final answer separates memory used from memory stored; if
   nothing safe and durable was stored, say so.

## Safety Gates

Read-only permits recall only. No-store, no-memory, blank-context, and repo-only
constraints prohibit storage even when ICM is named. Ask the user to resolve
contradictory instructions before writing memory.

Never store or reveal secrets, tokens, passwords, recovery codes, authentication
material, or raw session exports. Private but non-secret personal data may be
stored only when the user identifies the exact fact and explicitly asks ICM to
persist it; that specific request satisfies the privacy confirmation gate. Keep
ICM Codex-facing unless the user says otherwise.
