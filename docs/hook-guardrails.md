# Hook Guardrails

This repository owns source copies of local Codex hook behavior that is useful
to audit, test, and repair from source. Active global hook installs are not
mutated during ordinary repo work; copy or repair them only when the user asks
for install or repair.

## Hook Surfaces

| File | Current role |
| --- | --- |
| `hooks/wsl_command_guardrails.py` | Pre-tool-use shell guardrail for WSL Codex threads. It blocks command shapes that commonly mis-handle Windows shells, Windows paths, WSL mount temp dirs, and heredoc/script execution. |
| `hooks/test_wsl_command_guardrails.py` | Focused unittest coverage for the guardrail parser and blocking decisions. |
| `hooks/global-pretooluse-hooks.example.json` | Source-side example for wiring the guardrail into a global Codex hook matcher. |
| `.codex/hooks.json` | Repo-local hook config. It runs the Impeccable post-edit adapter for `Edit`, `Write`, and `apply_patch` tool events. |
| `bin/impeccable-hook.mjs` | Adapter that delegates to a local Impeccable skill hook only when that skill source exists in `.agents/skills/impeccable`. |

## WSL Command Guardrails

The guardrail protects this workstation's WSL Codex sessions from high-friction
shell mistakes:

- direct Windows shell launches when Tabby or a documented bridge is the safer
  surface;
- unconverted Windows paths passed into WSL commands;
- temporary file or directory assignments under `/mnt/<drive>`;
- heredoc and shell-wrapper cases that hide blocked commands;
- command strings that mix shell separators with risky Windows interop.

The guardrail is intentionally conservative about command shape, not command
intent. If a legitimate repair needs Windows-side behavior, prefer the
documented Windows/Tabby surface or use a small non-interactive bridge command
that is already covered by local guidance.

## Install Shape

The source file is tracked here:

```text
hooks/wsl_command_guardrails.py
```

The active global copy, when repaired, lives under:

```text
$CODEX_HOME/hooks/wsl_command_guardrails.py
```

The source-side hook snippet is:

```text
hooks/global-pretooluse-hooks.example.json
```

Keep the active global hook matcher aligned with the shell tool names used by
the current Codex runtime. Do not edit `$CODEX_HOME/hooks/` as part of ordinary
source changes unless the task is explicitly an install or repair.

## Repo-Local Impeccable Hook

`.codex/hooks.json` wires post-edit UI checks to `bin/impeccable-hook.mjs`. The
adapter exits successfully when `.agents/skills/impeccable/scripts/hook.mjs` is
not present, so the repo config is safe on machines that do not have that local
skill source.

This hook is source hygiene for UI changes in this repo; it is not a replacement
for Visual Canvas policy checks or browser/visual QA when a generated artifact
has meaningful visual surface area.

## Verification

Run the hook regression from the repo root:

```bash
python3 -m unittest hooks/test_wsl_command_guardrails.py
```

For hook config or adapter edits, also run:

```bash
node --check bin/impeccable-hook.mjs
git diff --check
```

When repairing the active global guardrail, verify the installed file path and
the active hook matcher separately after copying from source.
