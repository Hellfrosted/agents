# Agent Loop Designer Workflow Rules

## Surface Selection

Loop-first rule: design a loop whenever there is any repeated trigger, review,
repair, follow-up, monitoring, handoff, or decision pattern. Prefer the smallest
durable loop surface that captures trigger, state, artifact, stop condition, and
safety gate. Do not downgrade a repeatable workflow to a reusable prompt just
because it is simple.

- Reusable prompt: use only when the user explicitly asks for wording or when
  no durable trigger, state, artifact, or next action exists after inspection.
- Skill: use when the workflow is manually triggered and repeatable.
- Automation prompt: use when the workflow should run on a schedule or heartbeat.
- Worktree policy: use when the loop may edit code, fan out candidates, or run beside active local work.
- Plugin-backed skill: use when the workflow should be installed, shared, or bundled with connectors/MCP.
- Visible exploration threads: default for exploration that the user may read, interrupt, steer, or use for decisions.
- Subagents: default for private/background read-only or artifact-only work such as docs research, triage, test/log analysis, review, and summarization when a compact packet is enough.
- Worktree-backed threads: default for independent candidates that need visible exploration, code edits, or isolated checkout state. In WSL, prefer same-project manual git worktrees nested under the saved project checkout over Codex-managed worktrees on Windows storage.
- Grit-backed coordination: use when parallel code-editing work has symbol-level ownership in a Grit-supported language. Grit supplies the worktrees, locks, and integration protocol; see `grit-coordination.md`.

For private subagent delegation, see `subagents.md` before assigning context-heavy work.
For worktree-backed threads, see `worktree-threads.md` before assigning parallel candidate work.
For Grit-backed coordination, see `grit-coordination.md` before claiming symbols or running `grit done`.
For CLIs, tools, subagents, custom agents, automations, plugins, MCP/connectors, or config changes, see `docs-first.md` before proposing or using the surface.
For failures caused by this plugin's own instructions, see `self-improvement.md` before final response.

## Safety Defaults

- Preserve user changes.
- Every loop has a stop condition, a maximum scope or iteration limit, and a
  human gate for destructive, credential, push/PR, deploy, external-message,
  install, marketplace, or plugin-cache changes.
- Keep the main agent as a thin coordinator for non-trivial loops.
- Use visible threads when the user may need to read, intervene, steer, or decide from exploration.
- Delegate private/background context-heavy work to subagents and require compact, structured packets back.
- Use Worktree threads for visible exploration, background worker threads, or parallel code changes.
- Let the coordinator query subagents or Worktree threads for missing context before integration, blocker decisions, or final reporting.
- Do not commit, push, deploy, delete, or message external people unless the user explicitly asks.
- Do not run `grit done` unless the loop is allowed to commit and integrate worker work. Do not run `grit session pr` unless the user explicitly asks for push/PR behavior.
- Do not store secrets or raw private logs in markdown state.
- Make the user decision point explicit before moving from discovery to implementation.
- Keep worker topology internal; user-facing output should be the loop-run result.
- When a loop is made durable, define how failures feed back into the skill or plugin without exposing control-plane details to the user.
- Treat examples such as PR-comment loops as patterns, not defaults. Choose the
  actual loop from the user's setup and current workflow evidence.

## Consistency Layer

Use `$PLUGIN_ROOT/scripts/loop_spec.py` when a loop will become a skill, automation, plugin-backed workflow, Worktree-thread workflow, or subagent workflow. Set `PLUGIN_ROOT` to the agent-loop-designer plugin root first; from this reference directory, use `PLUGIN_ROOT="$(cd ../../.. && pwd)"`. Validate the spec before creating durable files. This catches missing decision points, unclear state, and ambiguous worker ownership.

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
- Use the plugin cachebuster/reinstall flow for existing local plugins only
  when the user explicitly asks to install or refresh the active plugin.
- Add or preserve a self-improvement rule for failures caused by the plugin's own instructions.
