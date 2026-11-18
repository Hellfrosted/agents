---
type: operations
title: WSL shell guardrail
description: Explains how the Codex pre-tool hook checks WSL shell input and how maintainers verify or repair the installation.
tags: [codex, wsl, hooks]
verified:
  - by: openwiki/0.9.0
    at: 2026-11-18T04:17:50.689Z
sources:
  - id: openwiki-source-27b90d002c543e00a7b9af37
    resource: repo://docs/hook-guardrails.md
  - id: openwiki-source-a51a0c48587a1f718ac5a66a
    resource: repo://hooks/test_wsl_command_guardrails.py
  - id: openwiki-source-3f6511c1d5f18b8a34b2db6e
    resource: repo://hooks/wsl_command_guardrails.py
generated: { by: "pi", at: "2026-11-18T04:17:50.689Z" }
---

# WSL shell guardrail

[`wsl_command_guardrails.py`](../../hooks/wsl_command_guardrails.py) is a Codex `PreToolUse` hook for avoiding common Windows/WSL command mistakes. Its checks look at command shape, including Windows shell launches, temporary assignments on Windows mounts, and shell wrappers. Treat it as an operator guardrail, not as a sandbox.

The hook reads a JSON event from stdin. It checks shell commands before attempting path rewrites, then checks the rewritten commands again. A rejection emits a JSON `deny` decision with a reason. An accepted rewrite emits `allow` with `updatedInput`; the ordinary no-change path emits nothing. For parallel calls, the rewriter only changes recognized shell entries and preserves the other entries.

## Verify or repair

Run the [focused regression tests](../../hooks/test_wsl_command_guardrails.py) from the root in WSL. They start the hook as a subprocess and assert the output for allow, deny, and rewrite cases:

```bash
python3 -B -m unittest hooks/test_wsl_command_guardrails.py
```

Keep source edits in the repository. Only when installation or repair is explicitly requested, copy the hook into `$CODEX_HOME/hooks/` and verify both the installed file and the active matcher in `$CODEX_HOME/hooks.json`. The [example matcher](../../hooks/global-pretooluse-hooks.example.json) is an example, not a portable drop-in: check the path and actual tool names on the target runtime. See the [maintainer guide](../../docs/maintainer-guide.md) for local repair rules and [quickstart](../quickstart.md) for other tasks.
