# Loop Spec

Use the loop spec when consistency matters. It is the structured contract behind the simple user command.

Generate a blank spec:

```bash
python3 "$PLUGIN_ROOT/scripts/loop_spec.py" template --task "recurring task"
```

Validate a spec:

```bash
python3 "$PLUGIN_ROOT/scripts/loop_spec.py" validate loop-spec.json
```

Render a spec:

```bash
python3 "$PLUGIN_ROOT/scripts/loop_spec.py" render loop-spec.json
```

Required fields:

- `mission`
- `command`
- `surface`
- `trigger`
- `state`
- `inputs`
- `workers`
- `tools` when the loop depends on CLIs, MCP tools, connector tools, app tools, or scripts
- `docs_checked` when `tools`, `worktree_threads`, or `subagents` are present, or when the surface is `automation-prompt`, `worktree-policy`, or `plugin-backed-skill`
- `artifact`
- `decision_point`
- `next_action`
- `safety_rule`
- `stop_condition`
- `loop_budget`
- `human_gate`
- `failure_learning`
- `grit_coordination` when the loop uses Grit claims, Grit worktrees, or `grit done`

Valid `surface` values:

- `reusable-prompt`
- `skill`
- `automation-prompt`
- `worktree-policy`
- `plugin-backed-skill`

Use `reusable-prompt` only when the user explicitly asks for prompt wording or
when no durable trigger, state, artifact, or next action exists after
inspection. If a task can repeat, choose a loop-capable surface instead.

Loop control fields:

- `stop_condition`: observable state where the loop must stop.
- `loop_budget`: maximum iterations, wakeups, wall time, spend, or scope.
- `human_gate`: actions that require explicit user approval before continuing.

Each subagent must define:

- `role`
- `spawn_policy`: omit `fork_context` or set `fork_context: false` in the
  current `spawn_agent` tool for role-specific agents; use
  `fork_turns: "none"` on tool surfaces that expose `fork_turns`; full history
  only when inheriting the parent agent type/model/reasoning is intended
- `ownership`
- `write_intent`: `none`, `artifact-only`, or `code-editing`
- `sandbox_expectation`: `read-only`, `workspace-write`, `permission-profile`, or `danger-full-access`
- `approval_expectation`
- `output_contract`
- `coordination_rule`
- `context_budget`: what raw context the subagent may absorb and what compact packet it must return

Each Worktree thread must define:

- `role`
- `ownership`
- `write_intent`: `none`, `artifact-only`, or `code-editing`
- `starting_state`
- `location_strategy`: `wsl-manual` or `codex-managed`
- `visibility`: whether the thread is human-readable and can receive user intervention
- `output_contract`
- `integration_rule`

When `grit_coordination` is present, define:

- `backend`: `local` unless the user explicitly requested Azure/S3-compatible team coordination.
- `init_policy`: how the loop detects missing or stale Grit state and runs `grit init`, `grit config set-local`, and `grit symbols`.
- `claim_strategy`: explicit symbols, dependency-aware claims, bounded `grit assign`, or queued claims.
- `done_policy`: when `grit done -a <agent-id>` is allowed, and when the coordinator reports worktrees instead of integrating.
- `thread_context_rule`: how the coordinator reads or messages worker threads for missing integration context.
- `cleanup_rule`: how the loop checks `grit status`, releases only loop-owned locks, and handles stale locks.

`failure_learning` must define:

- `trigger`: observed failure that starts a plugin update.
- `evidence`: output or behavior that proves the failure.
- `update_target`: narrowest source plugin file to patch.
- `validation`: source/cache and script checks to run.
- `skip_when`: cases where the failure should not update the plugin.

The validator rejects subagents whose `spawn_policy` does not explicitly encode
a non-full-history role-specific launch rule for either the current
`fork_context` surface or a `fork_turns` surface. It also rejects `code-editing`
subagents with a `read-only` sandbox expectation.
The validator also rejects specs that omit `failure_learning`, name `tools`, `worktree_threads`, or `subagents`, use docs-sensitive surfaces, or mention Worktree threads/subagents/custom agents in `workers` without `docs_checked`. If `workers` mentions Worktree threads or subagents/custom agents, define the matching structured entries too.
The validator allows extra structured fields such as `visibility`, `context_budget`, and `grit_coordination`; the loop author is responsible for filling those contracts when they are named in `tools`, `workers`, or `inputs`.
