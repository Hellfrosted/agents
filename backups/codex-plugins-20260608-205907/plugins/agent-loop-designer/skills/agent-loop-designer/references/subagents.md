# Subagent Design Checklist

Use this when a loop uses subagents or custom agents.

Subagents are not the default fan-out model for this plugin. Prefer worktree-backed threads for independent candidates and any code-editing work. In WSL, that means manual Linux-native git worktrees by default. Use this file only when the user explicitly asks for subagents or when a small read-only check should stay inside the current thread.

Read the current Codex subagents and configuration docs before creating or recommending custom agent files, permission profiles, sandbox settings, or subagent prompts. Record the docs in `docs_checked` when using the loop spec.

## Why This Exists

Subagents are best for parallel, read-heavy work: exploration, review, triage, test/log analysis, and summarization. Write-heavy subagents need tighter design because parallel edits create conflicts and permission mismatches.

The common failure mode is assigning implementation work to an agent that effectively runs read-only. That can happen when the custom agent file, parent session runtime overrides, sandbox mode, approval policy, or permission profile do not agree.

## Two Config Interactions To Check

1. **Parent runtime overrides vs custom agent defaults**
   - Custom agent files can set defaults such as `sandbox_mode`, `model`, `model_reasoning_effort`, MCP servers, and skills.
   - Subagents inherit the parent session when fields are omitted.
   - Codex reapplies the parent turn's live runtime overrides when spawning a child, including `/permissions` changes, `--yolo`, sandbox choices, and approval choices.
   - Result: a writer custom agent may still fail if the active parent turn is effectively read-only or cannot surface/allow the needed approval.

2. **Permission profiles vs legacy sandbox keys**
   - Codex has legacy sandbox keys such as `sandbox_mode` and `sandbox_workspace_write.*`.
   - Codex also has permission profiles through `default_permissions` and `[permissions.<name>]`.
   - Do not design a loop that mixes `default_permissions` with `sandbox_mode` or `[sandbox_workspace_write]` in the same intended config surface. Pick one model for the loop.

## Required Subagent Plan

For every subagent, specify:

- **Role**: what this agent is responsible for.
- **Ownership**: exact files, directories, or read-only area.
- **Write intent**: `none`, `artifact-only`, or `code-editing`.
- **Sandbox expectation**: `read-only`, `workspace-write`, or explicit permission profile.
- **Approval expectation**: whether the task can proceed without fresh approval.
- **Output contract**: summary, findings, patch, test result, or JSON.
- **Coordination rule**: whether the parent waits, messages follow-up, or closes the agent.

## Safe Defaults

- Use read-only subagents for exploration and review.
- Use one writer at a time for code edits unless ownership is completely disjoint.
- If a child must write, ensure the parent turn is not in read-only mode and can handle the required approval flow.
- In non-interactive or automation runs, avoid child tasks that need fresh approval. They will fail if approval cannot surface.
- Do not rely on a child `sandbox_mode = "workspace-write"` if the parent turn has a stricter live runtime override.
- Prefer explicit role names: `explorer`, `reviewer`, `test_runner`, `docs_researcher`, `worker`.

## Example Read-Only Agent

```toml
name = "architecture_explorer"
description = "Read-only explorer for architecture friction and domain vocabulary."
model_reasoning_effort = "medium"
sandbox_mode = "read-only"
developer_instructions = """
Map relevant code paths and return concise findings with file references.
Do not edit files.
"""
```

## Example Writer Agent

```toml
name = "targeted_worker"
description = "Implementation worker for one explicitly assigned file/module owner."
model_reasoning_effort = "medium"
sandbox_mode = "workspace-write"
developer_instructions = """
Edit only the files assigned by the parent.
Do not revert other agents' or user changes.
Stop and report if permissions prevent the requested write.
"""
```

Before spawning this writer, the parent loop should confirm the active session allows workspace writes or uses a permission profile that grants the needed write roots.

## Prompt Pattern

```text
Use subagents for the read-heavy parts only:
- architecture_explorer: read-only, map relevant modules and return file references.
- reviewer: read-only, find correctness and test risks.

Wait for both, summarize their findings, then ask me before spawning any writer.
```
