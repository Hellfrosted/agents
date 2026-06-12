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
- `failure_learning`

Valid `surface` values:

- `reusable-prompt`
- `skill`
- `automation-prompt`
- `worktree-policy`
- `plugin-backed-skill`

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

Each Worktree thread must define:

- `role`
- `ownership`
- `write_intent`: `none`, `artifact-only`, or `code-editing`
- `starting_state`
- `location_strategy`: `same-project-nested-manual` by default; other
  strategies require explicit user approval
- `output_contract`
- `integration_rule`

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
