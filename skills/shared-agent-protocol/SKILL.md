---
name: shared-agent-protocol
disable-model-invocation: true
---

# Shared Agent Protocol

Reference only. Do not invoke this skill directly as a task workflow. Entrypoint
skills may link to a specific section when they need shared language without
adding model-visible trigger context.

Keep task ownership in the entrypoint skill. This file defines common protocol
concepts only; it must not grow into a general Codex handbook.

## Codex Delegation and Reviewer Protocol

Use delegation only when the active request or host surface authorizes it. If
the user forbids subagents, or the host requires authorization that has not
been granted, ask before spawning agents or use a non-delegated fallback when
the entrypoint skill allows it.

Keep the main agent responsible for scope, synthesis, fixes, and the final
decision. Give reviewers narrow lanes: read-only, no spawned agents, no file
edits, no final decision. Assign each reviewer an angle such as skepticism,
verification, security, implementation risk, maintenance, UX, or loop shape.

For Codex reviewer forks, prefer role-specific non-full-history forks. Provide
only the task, goal, scope, evidence, open questions, and relevant snippets. Do
not combine a full-history fork with role/model/reasoning overrides on tool
surfaces that reject that combination.

When a workflow requires a dedicated reviewer or lane goal, use a goal-writer
subagent first. The goal-writer is read-only, may not spawn agents, and returns
only the goal. The main agent then passes that goal to the reviewer or lane
worker.

Reviewer deliverables should be actionable and evidence-backed: blocking
findings, loopholes, verification gaps, or confidence notes tied to source
lines, commands, logs, docs, or explicit reasoning. Unsupported objections are
marked speculative and do not control the final decision.

## Approval, Read-Only, and Side-Effect Gates

Treat user scope, non-goals, allowed write sets, and read-only instructions as
hard gates. Read-only permits inspection and reasoning only; it excludes writes,
sync/refresh commands that mutate local state, memory storage, staging,
commits, pushes, installs, updates, publishing, and remote/cloud configuration.

Stop for explicit approval before destructive actions, credential access,
secret handling, dependency changes, machine-wide configuration, network
exposure, production or deployment changes, autonomous execution, automations,
commits, pushes, force operations, or tool-specific publishing/remote-sync
features.

When approval is missing, keep the action report-only: state the exact command,
target, risk, and expected postcondition instead of running it. If instructions
conflict, ask the smallest question that resolves the gate.

## Evidence-Backed Output

Tie material claims to evidence that another agent can audit: file paths and
line numbers, command outputs, test names, source URLs, timestamps, ids, or
quoted task constraints. Separate fact, inference, and uncertainty.

Use the smallest relevant verification for the boundary being changed or
claimed. If verification cannot run, say what was skipped and why. Do not pad
outputs with hypothetical risks that have no path to action.

For private or local archives, cite enough context to audit the finding without
dumping long conversations or sensitive material. For memory workflows, store
only explicit, durable, non-secret summaries when storage is in scope.

## Bounded Loops

Define repeated work as a bounded loop before running it:

- Trigger: what starts or wakes the loop.
- State: where progress, decisions, findings, and queue items live.
- Next action: what each pass does.
- Stop condition: concrete success, blocker, or exhaustion state.
- Human gate: actions that require explicit approval.
- Budget: max iterations, reviewers, wakeups, wall time, spend, or scope.
- Failure learning: durable skill, doc, memory, or process update that prevents
  the same miss from recurring, when such an update is in scope.

End the loop when the stop condition is met, a required approval is missing, or
the budget is exhausted. Report pass count, reviewer count when relevant,
checks run, remaining uncertainty, and the next gated action.

## Integration Pointers

Use exact section links from entrypoint skills instead of copying these rules:

- `confidence-loop`: reviewer lanes, evidence-backed loopholes, bounded loop
  shape, and explicit confidence uncertainty.
- `evo-end-to-end`: approval gates, Codex-vs-tool worker boundaries, bounded
  experiment loops, and report-only unsafe actions.
- `tuck`: commit-scoped read-only reviewers, no staging before review, and git
  side-effect gates.
- `task-brief`: compact gates, evidence criteria, and semantic boundaries for
  delegated work.
- `discrawl`: read-only/no-refresh gates, local archive evidence, and private
  conversation minimization.
- `icm`: recall-vs-store gates, no-memory/no-store constraints, and safe
  durable summaries.
