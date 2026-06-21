---
name: agent-loop-designer
description: Turns `$agent-loop-designer` plus a plain-language recurring task into a simple repeatable Codex loop, skill, automation prompt, Worktree-thread workflow, or plugin-backed workflow. Use when the user wants to stop remembering prompts for recurring Codex work, candidate fan-out, triage/discovery workflows, or repeatable skill/plugin workflows. Not for one-off skill/plugin creation or editing.
---

# Agent Loop Designer

## Simple Use

The user should only need to type:

```text
$agent-loop-designer <describe the thing I keep doing>
```

Treat everything after the skill name as the recurring workflow. Convert it into the simplest repeatable Codex loop that removes prompt memory from the user.

Loopify by default. Prefer a real loop surface whenever a task has any repeated
trigger, review, repair, follow-up, monitoring, handoff, or decision pattern.
Use a reusable prompt only when the user explicitly asks for prompt wording or
when no durable trigger, state, artifact, or next action exists after inspection.

## What To Produce

Default to a short answer with:

- The one-line command they will use next time.
- The result shape: what the loop returns when it runs.
- The smallest durable surface and concrete next step.

Do not put subagent plans, Worktree-thread plans, worker topology, or control-plane details in the user-facing result unless the user asks how the loop works. Workers are how the loop gets to the result, not the result.

If the user asks to actually make it repeatable, implement the durable surface rather than only explaining it.

Use [workflow-rules.md](references/workflow-rules.md) for surface selection, safety defaults, and packaging rules.
When consistency matters, produce and validate a loop spec with `$PLUGIN_ROOT/scripts/loop_spec.py`; see [loop-spec.md](references/loop-spec.md). Resolve `PLUGIN_ROOT` to the plugin root before running helper scripts; from this skill directory, use `PLUGIN_ROOT="$(cd ../.. && pwd)"`.
Before proposing or using CLIs, tools, Worktree threads, subagents, automations, plugins, MCP/connectors, or config changes, follow [docs-first.md](references/docs-first.md).
For parallel code-editing loops with symbol-level ownership, use [grit-coordination.md](references/grit-coordination.md) as the coordination backend. If Grit is not initialized in the target repo, the loop starts local initialization itself.
When an Agent Loop Designer run fails because the skill encoded a bad assumption, missed a guardrail, or produced a fragile instruction, apply [self-improvement.md](references/self-improvement.md) before the final response only when source edits are in scope. For read-only, diagnosis-only, or planning-only runs, report the proposed source fix instead.

If the request is about running `improve-codebase-architecture` and then automatically exploring or acting on all candidates, do not produce the old "pick one candidate" radar. Route it to `$architecture-cascade <optional repo area or problem>`.

## Required Loop Shape

Every loop must define:

- **Trigger**: prompt, skill invocation, standalone automation, or thread automation.
- **State**: markdown/control-plane files read or written by the loop.
- **Inputs**: repo files, docs, issues, logs, connectors, or external resources.
- **Workers**: thin main coordinator, visible threads for human-steerable exploration or isolated code-editing work, subagents for private/background context-heavy tasks, or main agent only for tiny loops.
- **Coordination**: none, manual Worktree threads, or Grit-backed symbol locks.
- **Artifact**: report, patch, issue list, PR draft, test result, reference doc, or decision record.
- **Decision point**: the user question that gates the next phase.
- **Next action**: stop, deepen design, implement, schedule, create worktree, or record a decision.
- **Safety rule**: what prevents unintended edits or excessive autonomy.
- **Failure learning**: what plugin or skill update will prevent repeat agent failure.

Keep the user's command memorable. Do not make them remember "control plane", "state model", or feature taxonomy. Use `references/control-plane-template.md` internally when a template helps.

Keep the main agent as a coordinator by default. Put exploration that needs human reading, intervention, steering, or decision-making in a visible thread so the user can inspect and interrupt it. Delegate private/background discovery, docs research, triage, test/log analysis, and review to subagents only when a compact packet is enough and no human intervention is expected. Use [subagents.md](references/subagents.md). Use worktree-backed threads for visible exploration or code-editing isolation, and Grit worktrees for symbol-locked code-editing; see [worktree-threads.md](references/worktree-threads.md) and [grit-coordination.md](references/grit-coordination.md).

Keep the Required Loop Shape and worker plan internal. The visible answer should read like the output contract of the running agent loop: report fields, changed artifacts, blocked items, verification, and next user decision.

## Workflow

1. If the text after `$agent-loop-designer` is missing or too vague to act on, ask one short question: what should Codex repeatedly do?
2. State the mission in one sentence.
3. Choose the smallest loop surface that can repeat the work: skill,
   automation prompt, Worktree-thread workflow, worktree policy, or
   plugin-backed skill. Use reusable prompt only as the fallback above.
4. Explain the choice in one sentence.
5. Define the failure-learning rule for the loop.
6. If the user asked to create it, edit the needed files, validate, and report
   the new one-line command. Reinstall plugin-backed loops only when the user
   explicitly asks to install or refresh the active plugin.

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

1. Start a visible architecture exploration thread to run the `improve-codebase-architecture` discovery process and return a compact candidate packet plus report location. The user can read and intervene in that thread while the main coordinator stays small.
2. Extract every candidate from the exploration thread packet into a queue with candidate id, files, recommendation strength, conflicts, and proposed owner. If a fact is missing, ask that thread for the smallest missing detail.
3. Create worktree-backed candidate threads. In WSL, create manual git worktrees under `<repo>/.codex-worktrees/architecture-cascade/...` by default and start local project threads under the same saved project, with each worker prompt naming its nested worktree cwd.
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
