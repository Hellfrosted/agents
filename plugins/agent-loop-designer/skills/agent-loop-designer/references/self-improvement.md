# Self-Improvement

Use this when an Agent Loop Designer skill run fails because the plugin's own
instructions were incomplete, wrong, or too fragile.

## Failure Learning Contract

Every durable loop should define:

- `trigger`: which observed failure starts a plugin update.
- `evidence`: which command output, app error, API rejection, missing artifact, or behavior proves the failure.
- `update_target`: the narrowest source skill, reference, script, or validator rule to patch.
- `validation`: how to verify the source plugin and any bundled script behavior.
- `skip_when`: when not to update the plugin.

## Patch Rule

If the failure came from this plugin's instructions and source edits are in
scope, patch the plugin in the same turn before final response:

1. Identify the bad assumption or missing guardrail.
2. Patch the source plugin under the current agents-toolkit checkout, usually
   `plugins/agent-loop-designer/`.
3. Keep the change operational: guardrail, preflight, contract, or validator
   rule.
4. Search stale wording with `rg` and run the smallest relevant bundled script
   check.

Do not mirror changes into active plugin caches, installed plugin directories,
or persistent memory unless the user explicitly asks for that separate
active-install or memory update.

For read-only, diagnosis-only, or planning-only runs, report the narrow source
patch that would fix the plugin instead of editing files.

Do not update the plugin for target-repo bugs, user-cancelled work, unavailable
external services, permission denials the user intentionally chose, or one-off
data issues. Report those normally and update the target repo's docs only when
that is the right durable home.
