# Agent Loop Designer Workflow Rules

## Surface Selection

- Reusable prompt: use when the user only needs memorable wording.
- Skill: use when the workflow is manually triggered and repeatable.
- Automation prompt: use when the workflow should run on a schedule or heartbeat.
- Worktree policy: use when the loop may edit code, fan out candidates, or run beside active local work.
- Plugin-backed skill: use when the workflow should be installed, shared, or bundled with connectors/MCP.
- Worktree-backed threads: default for independent candidates, parallel exploration, or code-editing work. In WSL, prefer manual Linux-native git worktrees over Codex-managed worktrees on Windows storage.
- Subagents: use only when the user explicitly asks for subagents or the work is small and read-only inside the current thread.

For worktree-backed threads, see `worktree-threads.md` before assigning parallel candidate work.
For explicit subagent requests, see `subagents.md` before assigning write work.
For CLIs, tools, subagents, custom agents, automations, plugins, MCP/connectors, or config changes, see `docs-first.md` before proposing or using the surface.
For failures caused by this plugin's own instructions, see `self-improvement.md` before final response.

## Safety Defaults

- Preserve user changes.
- Prefer read-only exploration before editing.
- Use Worktree threads for background or parallel code changes.
- Do not commit, push, deploy, delete, or message external people unless the user explicitly asks.
- Do not store secrets or raw private logs in markdown state.
- Make the user decision point explicit before moving from discovery to implementation.
- Keep worker topology internal; user-facing output should be the loop-run result.
- When a loop is made durable, define how failures feed back into the skill or plugin without exposing control-plane details to the user.

## Consistency Layer

Use `scripts/loop_spec.py` when a loop will become a skill, automation, plugin-backed workflow, Worktree-thread workflow, or subagent workflow. Validate the spec before creating durable files. This catches missing decision points, unclear state, and ambiguous worker ownership.

If a spec names tools, Worktree threads, or subagents, include `docs_checked`. The validator rejects specs that do not record which docs, local help, or official references were read.
Every spec includes `failure_learning` so repeat agent failures have an explicit plugin-update path or an explicit skip rule.

## Automation Prompt Requirements

When drafting recurrence, specify:

- What to do each run.
- What counts as a finding.
- Where to write artifacts.
- Whether to run in a worktree.
- When to stop, archive, or ask the user.
- Which skill to invoke explicitly.

## Plugin Packaging Requirements

When creating or updating a plugin-backed loop:

- Keep `.codex-plugin/plugin.json` present.
- Put bundled skills under `skills/<skill-name>/SKILL.md`.
- Put optional detailed instructions under `references/`.
- Validate the plugin and the bundled skill.
- Use the plugin cachebuster/reinstall flow for existing local plugins.
- Add or preserve a self-improvement rule for failures caused by the plugin's own instructions.
