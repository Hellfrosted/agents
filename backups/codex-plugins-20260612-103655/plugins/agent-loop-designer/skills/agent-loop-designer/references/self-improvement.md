# Self-Improvement

Use this when an Agent Loop Designer skill run fails because the plugin's own
instructions were incomplete, wrong, or too fragile.

## Failure Learning Contract

Every durable loop should define:

- `trigger`: which observed failure starts a plugin update.
- `evidence`: which command output, app error, API rejection, missing artifact, or behavior proves the failure.
- `update_target`: the narrowest source skill, reference, script, or validator rule to patch.
- `validation`: how to verify the source plugin, active cache, and any bundled script behavior.
- `skip_when`: when not to update the plugin.

## Patch Rule

If the failure came from this plugin's instructions, patch the plugin in the same
turn before final response:

1. Identify the bad assumption or missing guardrail.
2. Patch the source plugin under `/home/crunch/plugins/agent-loop-designer`.
3. Mirror the same patch into the active cache under `/home/crunch/.codex/plugins/cache/personal/agent-loop-designer/...`.
4. Keep the change operational: guardrail, preflight, contract, or validator rule.
5. Validate source/cache parity with `diff`, search stale wording with `rg`, and run the smallest relevant bundled script check.
6. Store a concise privacy-safe memory when the failure is durable.

Do not update the plugin for target-repo bugs, user-cancelled work, unavailable
external services, permission denials the user intentionally chose, or one-off
data issues. Report those normally and update the target repo's docs only when
that is the right durable home.
