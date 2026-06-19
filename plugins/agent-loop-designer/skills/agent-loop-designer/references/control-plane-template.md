# Agent Loop Control Plane Template

Use this template when designing a repeatable Codex workflow.

## Mission

One sentence: what should Codex repeatedly help with, and why?

## Loop

- Trigger:
- State:
- Inputs:
- Workers:
- Artifact:
- Decision point:
- Next action:
- Stop condition:
- Loop budget:
- Human gate:
- Safety rule:

## Codex Surfaces

- Skill:
- Automation:
- Worktree:
- Plugin/connectors:
- Subagents:
  - Spawn policy: use non-full-history spawns for role-specific workers,
    reviewers, explorers, or custom-agent style assignments. In the current
    `spawn_agent` tool, omit `fork_context` or set `fork_context: false`; on
    tool surfaces that use `fork_turns`, set `fork_turns: "none"`.
    Full-history forks inherit parent agent type/model/reasoning and must not
    carry those overrides.

## First Manual Run

1. Read durable state.
2. Explore with bounded scope.
3. Produce a reviewable artifact.
4. Ask one decision question.
5. Stop unless the user chooses the next phase.

## Automation Upgrade

Schedule:

Prompt:

Finding criteria:

Stop/archive criteria:

Worktree policy:

## Risks

- What could be changed unintentionally?
- What should require explicit approval?
- What data should never be written to markdown or logs?
