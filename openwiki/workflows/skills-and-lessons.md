---
type: workflow
title: Skill workflows and lessons
description: Routes maintainers to repo-owned skill contracts, upstream workflows, and skill-authoring material.
tags: [codex, skills, maintenance]
verified:
  - by: openwiki/0.9.0
    at: 2026-11-18T04:17:50.689Z
sources:
  - id: openwiki-source-a9b55a65712732f6bcf808ba
    resource: repo://docs/agent-native-boundary.md
  - id: openwiki-source-7a21a87b4a02e711c924dcaa
    resource: repo://lessons/0001-prompting-for-codex-skills.html
  - id: openwiki-source-9c7448bd1a45150e267e7b46
    resource: repo://reference/skill-prompting-checklist.html
  - id: openwiki-source-923b5f957ed14536b821b731
    resource: repo://skills/shared-agent-protocol/SKILL.md
generated: { by: "pi", at: "2026-11-18T04:17:50.689Z" }
---

# Skill workflows and lessons

Start with the relevant `SKILL.md` for invocation rules and workflow limits:

| Task | Source |
| --- | --- |
| Search or inspect a local Discord archive | [Discrawl](../../skills/discrawl/SKILL.md) and its [local setup and troubleshooting reference](../../skills/discrawl/REFERENCE.md) |
| Turn rough intent into a Codex implementation ticket or handoff | [Task Brief](../../skills/task-brief/SKILL.md) |
| Audit confidence or stress-test a proposed strategy on request | [Confidence Loop](../../skills/confidence-loop/SKILL.md) |
| Review a diff for reuse, consistency, or unnecessary abstractions | [Code Refactor Review](../../skills/code-refactor-review/SKILL.md) |
| Plan an explicit Evo optimization campaign | [Evo End to End](../../skills/evo-end-to-end/SKILL.md) and its [operational reference](../../skills/evo-end-to-end/REFERENCE.md) |

[Shared Agent Protocol](../../skills/shared-agent-protocol/SKILL.md) is an internal reference for entrypoint skills, covering delegation, reviewer roles, approvals, and bounded loops.

## Edit a skill

Open the [interactive prompting lesson](../../lessons/0001-prompting-for-codex-skills.html) for examples and decision exercises, then use the [printable checklist](../../reference/skill-prompting-checklist.html) to check placement and whether repeated work has clear bounds. Edit the repository source first; don't overwrite active skills as part of routine maintenance.

## What belongs upstream

Use Codex's built-in visualization for conversational charts and diagrams. General visual plans and recaps belong to the installed BuilderIO Agent-Native skills. See the [ownership notes](../../docs/agent-native-boundary.md) for the full list of upstream commands and the install boundary, or the [quickstart](../quickstart.md) for other repository tasks.
