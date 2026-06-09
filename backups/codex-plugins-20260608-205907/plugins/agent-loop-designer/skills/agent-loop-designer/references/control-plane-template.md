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
- Safety rule:

## Codex Surfaces

- Skill:
- Automation:
- Worktree:
- Plugin/connectors:
- Subagents:

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
