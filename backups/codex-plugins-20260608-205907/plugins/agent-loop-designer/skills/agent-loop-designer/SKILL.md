---
name: agent-loop-designer
description: Turns `$agent-loop-designer` plus a plain-language recurring task into a simple repeatable Codex loop, skill, automation prompt, Worktree-thread workflow, or plugin-backed workflow. Use when the user wants to stop remembering prompts for recurring Codex work, candidate fan-out, triage/discovery workflows, or skill/plugin design.
---

# Agent Loop Designer

## Simple Use

The user should only need to type:

```text
$agent-loop-designer <describe the thing I keep doing>
```

Treat everything after the skill name as the recurring workflow. Convert it into the simplest repeatable Codex loop that removes prompt memory from the user.

## What To Produce

Default to a short answer with:

- The one-line command they will use next time.
- The result shape: what the loop returns when it runs.
- The smallest durable surface and concrete next step.

Do not put Worktree-thread plans, worker topology, or control-plane details in the user-facing result unless the user asks how the loop works. Worktree threads are how the loop gets to the result, not the result.

If the user asks to actually make it repeatable, implement the durable surface rather than only explaining it.

Use [workflow-rules.md](references/workflow-rules.md) for surface selection, safety defaults, and packaging rules.
When consistency matters, produce and validate a loop spec with `scripts/loop_spec.py`; see [loop-spec.md](references/loop-spec.md).
Before proposing or using CLIs, tools, Worktree threads, subagents, automations, plugins, MCP/connectors, or config changes, follow [docs-first.md](references/docs-first.md).
When an Agent Loop Designer run fails because the skill encoded a bad assumption, missed a guardrail, or produced a fragile instruction, apply [self-improvement.md](references/self-improvement.md) before the final response.

If the request is about running `improve-codebase-architecture` and then automatically exploring or acting on all candidates, do not produce the old "pick one candidate" radar. Route it to `$architecture-cascade <optional repo area or problem>`.

## Required Loop Shape

Every loop must define:

- **Trigger**: prompt, skill invocation, standalone automation, or thread automation.
- **State**: markdown/control-plane files read or written by the loop.
- **Inputs**: repo files, docs, issues, logs, connectors, or external resources.
- **Workers**: main coordinator, Worktree threads, or main agent only.
- **Artifact**: report, patch, issue list, PR draft, test result, reference doc, or decision record.
- **Decision point**: the user question that gates the next phase.
- **Next action**: stop, deepen design, implement, schedule, create worktree, or record a decision.
- **Safety rule**: what prevents unintended edits or excessive autonomy.
- **Failure learning**: what plugin or skill update will prevent repeat agent failure.

Keep the user's command memorable. Do not make them remember "control plane", "state model", or feature taxonomy. Use `references/control-plane-template.md` internally when a template helps.

For fan-out or code-editing loops, use worktree-backed threads by default and include role, ownership, write intent, starting state, location strategy, output contract, and integration rule. Use [worktree-threads.md](references/worktree-threads.md). Use subagents only when explicitly requested or for small read-only checks; then use [subagents.md](references/subagents.md).

Keep the Required Loop Shape and Worktree-thread plan internal. The visible answer should read like the output contract of the running agent loop: report fields, changed artifacts, blocked items, verification, and next user decision.

## Workflow

1. If the text after `$agent-loop-designer` is missing or too vague to act on, ask one short question: what should Codex repeatedly do?
2. State the mission in one sentence.
3. Choose one durable surface: reusable prompt, skill, automation prompt, Worktree-thread workflow, worktree policy, or plugin-backed skill.
4. Explain the choice in one sentence.
5. Define the failure-learning rule for the loop.
6. If the user asked to create it, edit the needed files, validate, reinstall when plugin-backed, and report the new one-line command.

## Quiet State

Create or propose markdown state only when it carries durable information:

- `CONTEXT.md` for project/domain vocabulary.
- `docs/adr/` for durable technical decisions.
- `STATE.md` or `QUEUE.md` for active loop state.
- `learning-records/` for teaching-oriented workflows.
- A non-repo HTML report for discovery output; on WSL, prefer a durable mounted-drive temp path over `/tmp` when the user may reopen it.

Keep local-only agent notes out of VCS unless the user asks for repo documentation.

## Architecture Cascade Pattern

For `improve-codebase-architecture` style loops where the user wants all candidates handled:

1. Run one discovery pass using the `improve-codebase-architecture` process.
2. Extract every candidate into a queue with candidate id, files, recommendation strength, conflicts, and proposed owner.
3. Create worktree-backed candidate threads. In WSL, prefer manual WSL-native git worktrees under `$HOME/codex-worktrees/...` and start local project threads there instead of using Codex-managed worktrees under Windows `$CODEX_HOME/worktrees`.
4. Act automatically in each candidate worktree when the candidate is bounded, non-conflicting, and has clear verification.
5. Ask only for blocked candidates: ADR conflict, overlapping ownership, unsafe permission mismatch, unclear tests, or product/domain decision.

## Example Prompt

```text
$agent-loop-designer run improve-codebase-architecture once, then explore and act on all candidates in worktree threads
```

Expected response shape:

```text
Use: $architecture-cascade

When it runs, it returns an architecture action report: implemented candidates with changed files and checks, rejected candidates with evidence, blocked candidates with the exact decision needed, and integration choices for completed worktrees.
```
