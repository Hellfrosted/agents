# Codex companion tools

This page lists companion tools installed on this machine for Codex workflows.
It is not a reference for the agent runtime or repo-local wrappers.

## Main list

- Evo: experiment runner for optimization work.
- RTK: command wrapper agents use for shell commands.
- ICM: persistent memory for durable Codex context.
- Codex Security: security review plugin from `openai-curated`.
- OpenAI Developer Docs MCP: official OpenAI docs server for API and product
  docs lookup.
- LazyCodex: installed as the `omo@sisyphuslabs` Codex plugin.

## Adjacent utilities

- `react-doctor`: React diagnostics and cleanup checks.
- `tokscale`: token counting and scale checks.
- `actionlint`: GitHub Actions workflow linting.

These are useful during agent work, but they are narrower than the main tools
above.

## Evo

Evo runs experiments for performance, architecture, flaky-test, slow-build, and
code-quality work.

Before running it:

- Check `evo --version`. It should report `evo-hq-cli`, not the unrelated SLAM
  package named `evo`.
- On Codex installs that hit exit-127 hook failures, recover with
  `uv tool install --force evo-hq-cli && evo install codex --force`; v0.4.5
  fixes the plugin cache path and doctor check, but a broken Codex install does
  not self-heal through `evo update`.
- Do not install or upgrade it unless the user asks.
- Write a short experiment brief first: goal, metric, baseline, gate, editable
  scope, read-only context, forbidden changes, backend, runtime/env, budget,
  stall rule, and merge rule.
- Get approval before Evo changes production behavior, APIs, persistence,
  auth/security, tests, packaging, dependencies, deployment, or user-visible
  behavior.

Useful commands:

```bash
evo init --host codex
evo host show
evo host set codex
evo config show
evo config backend show
evo config runtime show
evo env show
evo run --check
evo direct "<text>"
evo gc
```

Run discovery before optimization. Optimize only after the workspace has a
baseline and the editable scope is approved.

## RTK

RTK is the workstation command wrapper for agent-run shell commands. Use it for
searches, verification, package-manager commands, and one-off scripts.

Examples:

```bash
rtk rg --files
rtk pnpm test
rtk run 'icm recall "query"'
```

RTK is a wrapper convention. It does not replace understanding the command being
wrapped.

## ICM

ICM stores durable Codex context. Prefer the ICM MCP tools when they are
available. Use the CLI through `rtk` when MCP is not available.

Examples:

```bash
rtk run 'icm recall "query"'
rtk run 'icm store -t "topic" -c "summary" -i high'
```

Store resolved errors, architecture decisions, user preferences, and meaningful
project progress. Do not store secrets, tokens, passwords, recovery codes,
private personal data, or raw session exports.

## Codex Security

Codex Security is installed as the `codex-security@openai-curated` plugin. Use
it for security scans, diff reviews, validation, attack-path analysis, threat
models, and security fixes.

Keep its output focused on exploitable behavior. Security review should lead
with findings and include file and line references when possible.

## OpenAI Developer Docs MCP

The OpenAI Developer Docs MCP server is configured as `openaiDeveloperDocs`. Use
it for current OpenAI API, SDK, platform, and product docs.

Prefer it over web search for OpenAI-specific questions. Fetch the exact doc
page before quoting or summarizing details.

## LazyCodex

LazyCodex is installed as the `omo@sisyphuslabs` Codex plugin. Keep it enabled
unless the user explicitly asks to remove it.

For T3code sessions launched through `bin/codex-wsl.cmd`, the WSL runner
disables LazyCodex telemetry, auto-update, and config migration by default while
leaving the plugin available:

```bash
OMO_CODEX_DISABLE_POSTHOG=1
OMO_CODEX_SEND_ANONYMOUS_TELEMETRY=0
OMO_DISABLE_POSTHOG=1
OMO_SEND_ANONYMOUS_TELEMETRY=0
LAZYCODEX_AUTO_UPDATE_DISABLED=1
OMO_CODEX_AUTO_UPDATE_DISABLED=1
LAZYCODEX_CONFIG_MIGRATION_DISABLED=1
OMO_CODEX_CONFIG_MIGRATION_DISABLED=1
```

Serena is not part of the current installed Codex toolchain on this workstation;
it has already been uninstalled.

## Quick checks

```bash
evo --version
rtk --help
icm --help
codex plugin list
actionlint --version
react-doctor --version
tokscale --version
```

Some tools are exposed through MCP or plugins rather than plain CLI commands.
For those, verify them through the Codex config or the tool list in the running
session.
