# Codex companion tools

This page lists companion tools installed on this machine for Codex workflows,
plus repo-local wrappers that maintain those tools.

## Main list

- Evo: experiment runner for optimization work.
- RTK: command wrapper agents use for shell commands.
- ICM: persistent memory for durable Codex context.
- Codex Security: security review plugin from `openai-curated`.
- OpenAI Developer Docs MCP: official OpenAI docs server for API and product
  docs lookup.
- LazyCodex: installed as the `omo@sisyphuslabs` Codex plugin.
- Skills updater: Windows wrappers for checking and updating globally installed
  Codex skills.

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
- Keep the CLI and Codex plugin bundle in lockstep. Use
  `evo update codex --version 0.5.0` for the v0.5.0 line, then verify with
  `evo doctor codex`.
- If Codex is pinned to a stale local `evo-hq` marketplace, remove that
  marketplace source and add `evo-hq/evo --ref v0.5.0` before reinstalling.
- Do not install or upgrade it unless the user asks.
- Write a short experiment brief first: goal, metric, baseline, gate, editable
  scope, read-only context, forbidden changes, backend, runtime/env,
  per-experiment timeout, task skills, budget, stall rule, and merge rule.
- Get approval before Evo changes production behavior, APIs, persistence,
  auth/security, tests, packaging, dependencies, deployment, or user-visible
  behavior.
- For fine-tuning, post-training, reward design, or weight updates, use
  `$evo finetuning` before training-code edits. Keep held-out eval data out of
  training and run a smoke validation before spending the full budget.
- Before passing `subagents=N`, size the round from the binding resource:
  exclusive GPU/port/DB/shared mutable fixture means width 1 unless the harness
  isolates it; pool mode caps at slot count; remote mode caps at provider quota
  and cost.

Useful commands:

```bash
evo init --host codex
evo init --per-exp-timeout 600
evo host show
evo host set codex
evo config show
evo config get task-skills
evo config backend show
evo config runtime show
evo env show
evo run <exp_id> --check
evo run <exp_id> --timeout 600
evo wait --for ideators
evo abort <exp_id>
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

Codex also has a local `PreToolUse` shell hook at
`/home/crunch/.codex/hooks/rtk_pretooluse.py`. Despite the historical filename,
the hook is not a blanket Python or RTK requirement. Its current purpose is to
keep Python package management on `uv`.

The hook blocks direct `pip` package management, including:

```bash
pip install rich
pip3 install rich
python -m pip install rich
python3 -m pip install rich
rtk run 'python3 -m pip install rich'
bash -lc 'pip install rich'
```

Use these forms instead:

```bash
uv add rich
uv sync
uv run --with rich python -c 'import rich'
uv pip install --system rich
```

Plain Python execution is allowed. Skill validators, one-off scripts, and
commands such as `python3 script.py` or `python3 -c 'print(1)'` should not be
blocked by this hook. Keep it that way when editing the hook; the intended guard
is `uv` over `pip`, not "no Python".

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

## Skills Updater

The repo includes Windows wrappers for installed skill maintenance:

- `bin/sk-up.cmd`: short command names and flags.
- `bin/skills-updates.cmd`: long command names and flags.
- `bin/skills-updates.ps1`: implementation.

Common commands:

```bat
sk-up -l
sk-up -g
sk-up -d confidence-loop
sk-up -z confidence-loop evo-end-to-end
sk-up -i
sk-up -i confidence-loop
sk-up -i owner/repo
sk-up -s confidence-loop
sk-up -u confidence-loop
sk-up -S
sk-up -r confidence-loop
```

Use `-l` to list installed skills without checking upstream. Use `-g` to check
global skill status. Diff one skill with `-d`, or open one or more Zed diffs
with `-z`. Install all changed or missing skills with `-i`, install specific
lockfile skills with `-i <skill>`, or install a source URL/repo with
`-i <source>`. Remove a global skill with `-r`.

Skips are saved with `-s`, removed with `-u`, and listed with `-S`. They are
tied to the current upstream tree hash, so a new upstream tree makes the update
visible again.

The updater reads global skills from `%AGENTS_HOME%` when set, otherwise from
`%USERPROFILE%\.agents`. It caches upstream repositories and skip state under
`%LOCALAPPDATA%\skills-updates` when available, otherwise in a temp state
directory. Install and uninstall operations require `pnpm` and run
`pnpm dlx skills@latest`; the script protects `.skill-lock.json` with a mutex
and preserves existing lockfile fields around those operations.

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

Other T3code shim knobs:

- `CODEX_WSL_PROXY_IDLE_TIMEOUT_MS`: app-server idle timeout. The Windows shim
  defaults this to `1800000` milliseconds.
- `CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS`: timeout before the WSL proxy returns a
  fallback `skills/list` response. The proxy defaults this to `2000`
  milliseconds.
- `CODEX_WSL_PROXY_DEBUG_LOG`: WSL path for proxy debug logs.
- `CODEX_WSL_SHIM_DEBUG`: print the Windows shim launch arguments.

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
