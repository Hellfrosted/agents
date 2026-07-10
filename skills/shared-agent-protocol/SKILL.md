---
name: shared-agent-protocol
description: "Internal shared reference for delegation, reviewer, approval, and bounded-loop protocol used by entrypoint skills."
---

# Shared Agent Protocol

Reference only. Do not invoke this skill directly as a task workflow. Entrypoint
skills may link to a specific section when they need shared language without
adding model-visible trigger context.

Keep task ownership in the entrypoint skill. This file defines common protocol
concepts only; it must not grow into a general Codex handbook.

## Codex Delegation and Reviewer Protocol

Use delegation only when the active request, invoked entrypoint, or host surface
authorizes it. Invoking an entrypoint whose documented workflow requires
reviewers authorizes those reviewers unless the user forbids subagents or the
host requires a separate confirmation. When authorization is absent, ask before
spawning agents or use a non-delegated fallback when the entrypoint allows it.

Keep the main agent responsible for scope, synthesis, fixes, and the final
decision. Give reviewers narrow lanes: read-only, no spawned agents, no file
edits, no final decision. Assign each reviewer an angle such as skepticism,
verification, security, implementation risk, maintenance, UX, or loop shape.

For Codex reviewer forks, prefer role-specific non-full-history forks. Provide
only the task, goal, scope, evidence, open questions, and relevant snippets. Do
not combine a full-history fork with role/model/reasoning overrides on tool
surfaces that reject that combination.

Reviewer deliverables should be actionable and evidence-backed: blocking
findings, loopholes, verification gaps, or confidence notes tied to source
lines, commands, logs, docs, or explicit reasoning. Unsupported objections are
marked speculative and do not control the final decision.

## Approval, Read-Only, and Side-Effect Gates

Treat user scope, non-goals, allowed write sets, and read-only instructions as
hard gates. Read-only permits inspection and reasoning only; it excludes writes,
sync/refresh commands that mutate local state, memory storage, staging,
commits, pushes, installs, updates, publishing, and remote/cloud configuration.

An explicit request for an action, or invocation of an entrypoint whose stated
primary result necessarily includes that action, satisfies this protocol's
approval gate for ordinary in-scope execution. Do not ask for the same approval
again unless the host/tool requires it, the target is ambiguous, or the action
expands beyond the request.

Require specific authorization for destructive actions, credential or secret
access, dependency changes, machine-wide configuration, network exposure,
production or deployment changes, autonomous execution, automations, commits,
pushes, force operations, and tool-specific publishing or remote sync. An
entrypoint authorizes only the side effects its contract names; for example, a
local-commit workflow does not imply push authorization.

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

- [`confidence-loop`](../confidence-loop/SKILL.md): reviewer lanes,
  evidence-backed loopholes, and bounded loops.
- [`evo-end-to-end`](../evo-end-to-end/SKILL.md): approval gates,
  Codex-versus-tool worker boundaries, and bounded experiment loops.
- [`discrawl`](../discrawl/SKILL.md): bounded private-archive discovery,
  evidence, and delegated search lanes.
- [`tuck`](../tuck/SKILL.md): commit-scoped reviewers and git side-effect
  gates.
