---
type: quickstart
title: Agents Toolkit quickstart
description: Find the source and verification steps for repo-owned Codex hooks and skills.
tags: [codex, quickstart, maintenance]
verified:
  - by: openwiki/0.9.0
    at: 2026-11-18T04:17:50.689Z
sources:
  - id: openwiki-source-57066a52829deec25fa71dd7
    resource: repo://docs/maintainer-guide.md
  - id: openwiki-source-5aef912134f67c3d720302d2
    resource: repo://MISSION.md
generated: { by: "pi", at: "2026-11-18T04:17:50.689Z" }
---

# Agents Toolkit quickstart

This repository owns local Codex skill sources and the WSL shell guardrail source and tests. Workstation setup, shell configuration, automation definitions, and restore payloads are maintained separately in the dotfiles repository.

| If you need to… | Start here |
| --- | --- |
| Understand a blocked command, run regression tests, or verify an installed hook | [WSL shell guardrail](operations/wsl-shell-guardrail.md) |
| Pick or edit a skill, find the prompting lesson, or check what belongs upstream | [Skill workflows and lessons](workflows/skills-and-lessons.md) |
| Check maintenance and local repair rules | [Maintainer guide](../docs/maintainer-guide.md) |
| Understand scope or validate external product behavior | [Mission](../MISSION.md), [Product](../PRODUCT.md), and [Resources](../RESOURCES.md) |

Make changes to repository source first, run the smallest relevant check, and update the docs for the behavior you changed. Copy files to active installs only when installation, repair, or republishing is explicitly requested. Keep credentials, private personal data, and raw session exports out of documentation and examples.
