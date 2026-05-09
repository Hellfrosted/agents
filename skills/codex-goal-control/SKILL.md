---
name: codex-goal-control
description: Open the local Codex goal panel and manage the current Codex thread goal through bundled helper scripts. Use when the user asks to open the goal panel, inspect/set/pause/resume/complete/clear a thread goal, manage goals without the CLI, or make the local goal web app target the current thread.
---

# Codex Goal Control

This skill is experimental. It depends on Codex goals and app-server goal methods, which may change across Codex versions.

Prerequisite:

```toml
[features]
goals = true
```

## Resolve Skill Root

Run bundled scripts from this skill's installed directory, not from `CODEX_HOME`.

`CODEX_HOME` is for Codex configuration/runtime state. Skill managers may install skills somewhere else, so do not append `/skills/codex-goal-control` to `CODEX_HOME` or any remembered skills root.

Use the loaded skill path as the authority. Set `SKILL_ROOT` to the directory that contains the loaded `codex-goal-control/SKILL.md` file:

```bash
SKILL_ROOT="<directory-containing-codex-goal-control-SKILL.md>"
```

Do not derive `SKILL_ROOT` from `HOME`, `CODEX_HOME`, the current working directory, the OS, or a remembered install path. If the loaded skill path is unavailable, locate the installed `codex-goal-control` skill first, then verify it before running helpers:

```bash
ls "$SKILL_ROOT/scripts/codex_goal.js"
```

## Ground Rules

1. Identify the requested operation before running commands: open panel, inspect, set, pause, resume, complete, or clear.
2. Target the current Codex thread with `CODEX_THREAD_ID`. Only pass `--thread` when the user explicitly gives a different thread id.
3. Do not infer the thread from a browser URL, claim file, cwd, old panel state, or previous output.
4. Do not invent missing objectives, budgets, target threads, or destructive intent.
5. After any mutation, read the goal back with the thread goal API and report what was proven.

If `CODEX_THREAD_ID` is missing and the user did not explicitly provide `--thread`, stop and say the runtime did not expose the current thread id.

## Open Panel

```bash
node "$SKILL_ROOT/scripts/codex_goal_panel_open.js" --json
```

Return `threadUrl` as the panel link. Run this again for each new Codex thread. Do not reuse an old localhost URL unless it includes the intended `threadId`.

If the helper returns `started-unverified-sandbox`, the sandbox could not verify the local listener. Use `serverLogFile` for diagnosis. Only start the panel server manually when `CODEX_THREAD_ID` is present, the user requested a panel, and the environment allows binding localhost:

```bash
node "$SKILL_ROOT/scripts/codex_goal_panel_server.js" --thread "$CODEX_THREAD_ID" --host 127.0.0.1 --port 43873
```

## Direct Goal Commands

Use direct commands when the user asks to manage the goal without opening the panel:

```bash
node "$SKILL_ROOT/scripts/codex_goal.js" get --json
node "$SKILL_ROOT/scripts/codex_goal.js" set "objective" --json
node "$SKILL_ROOT/scripts/codex_goal.js" set "objective" --budget 3000 --json
node "$SKILL_ROOT/scripts/codex_goal.js" pause --json
node "$SKILL_ROOT/scripts/codex_goal.js" resume --json
node "$SKILL_ROOT/scripts/codex_goal.js" complete --json
node "$SKILL_ROOT/scripts/codex_goal.js" clear --json
```

Before `set`, run `get --json`. If a goal already exists and the user did not clearly ask to replace it, report the existing objective/status and ask before overwriting it.

For `set`, pass `--budget` only when the user provides a concrete positive integer budget.

For `complete`, act only when the user explicitly asks or when the agent has actually achieved the objective and no required work remains.

For `clear`, act only when the user explicitly asks to clear/delete/remove the goal, or after confirmation.

After `set`, `pause`, `resume`, `complete`, or `clear`, run:

```bash
node "$SKILL_ROOT/scripts/codex_goal.js" get --json
```

For `clear`, a missing goal is the expected read-back proof.

## Proof Labels

- `thread_proven`: the Codex thread goal API returned the expected state for the target thread.
- `panel_proven`: a follow-up localhost panel `GET /api/goal` returned the expected state.
- `visual_proven`: the user or browser automation confirmed the visible browser render.

Do not present command success, server startup, a URL, or a stale claim as proof that the current thread's goal changed.
