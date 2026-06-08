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
- CodSpeed: CLI plus hosted MCP for performance runs, comparisons, and
  flamegraph analysis.
- LazyCodex: installed as the `omo@sisyphuslabs` Codex plugin.
- LazyCodex local MCP: LSP, AST grep, Context7, and grep.app search are provided
  by the OMO plugin when their servers work on this host.
- OMO CLI: local `omo` entrypoint for LazyCodex-specific helpers.
- Discrawl: local Discord cache archive/search for Vesktop wiretap-only use.
- Skills updater: Windows wrappers for checking and updating globally installed
  Codex skills.

## Adjacent utilities

- `react-doctor`: React diagnostics and cleanup checks.
- `tokscale`: token counting and scale checks.
- `actionlint`: GitHub Actions workflow linting.
- `wslpath`: preferred WSL/Windows path conversion utility.
- `wslu`: optional WSL convenience tools, mainly `wslview`.

These are useful during agent work, but they are narrower than the main tools
above.

## WSL path and opener utilities

Use `wslpath` for path conversion in scripts and diagnostics. On this
workstation, `/usr/local/bin/wslpath` is first on `PATH` and converts the common
forms used by this repo:

```bash
wslpath -w /mnt/e/dev/agents-toolkit
wslpath -u 'E:\dev\agents-toolkit'
wslpath -u 'C:\Program Files\Git\cmd'
```

Keep WSLU installed only as optional convenience tooling. The upstream
`wslutilities/wslu` repository is archived, but Ubuntu still ships `wslu`; it is
acceptable for interactive helpers such as `wslview`. Do not add new repo or
shim logic that depends on WSLU when direct WSL interop works.

Prefer direct WSL interop for durable automation:

- `wslpath` for WSL/Windows path conversion.
- `explorer.exe .` for opening a folder from WSL.
- `powershell.exe /c start .` or `cmd.exe /c start` for opening files and URLs.
- `wsl.exe` for WSL management from Windows.

## Discrawl

Discrawl is installed in WSL as a local Discord archive/search tool. This
workstation uses it in Vesktop wiretap-only mode, with Discord bot tokens
disabled. See [Discrawl Vesktop Wiretap](discrawl-wiretap.md) for paths,
configuration, verification, and normal search commands.

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

Codex may also have a local `PreToolUse` shell hook under
`$CODEX_HOME/hooks/`. Despite historical RTK-oriented filenames, that hook is
not a blanket Python or RTK requirement. Its current purpose is to keep Python
package management on `uv`.

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
`pnpm dlx skills@latest`; global operations are forced to the universal
`.agents/skills` target. The script protects `.skill-lock.json` with a mutex
and preserves existing lockfile fields around those operations. Uninstalls also
remove the global installed skill directory, clear saved skips for the skill,
and remove the skill's lockfile entry under the same operation lock. If post-CLI
cleanup fails, the updater restores the pre-uninstall lockfile snapshot so
directory and lockfile state do not diverge.

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

## CodSpeed

CodSpeed CLI is installed at `/home/crunch/.local/bin/codspeed`. It is
authenticated as the default CLI profile.

The hosted CodSpeed MCP server is configured globally in Codex as `CodSpeed`:

```bash
codex mcp get CodSpeed
codex mcp login CodSpeed
```

Use `codspeed status` to verify CLI auth, repository linkage, and available
local executors. Local executor dependencies are installed on this WSL host:

- Simulation: CodSpeed's Valgrind fork is installed from source under
  `/home/crunch/.local`; `valgrind --version` should report
  `valgrind-3.26.0.codspeed3`.
- Walltime: Ubuntu `linux-tools-common`, `linux-tools-generic`, and
  `linux-perf` provide `/usr/bin/perf`.
- Memory: `codspeed setup` installs `codspeed-memtrack` automatically.

`codspeed setup status` should show green checks for Valgrind, perf, and
memtrack. The CLI and MCP are usable for account auth, repository/local-runs
uploads, hosted run data, comparisons, and flamegraph queries when runs exist.

## LazyCodex

LazyCodex is installed as the `omo@sisyphuslabs` Codex plugin. Keep it enabled
unless the user explicitly asks to remove it.

The local CLI entrypoint is `omo`. It does not expose a `--version` flag, so
use `omo help` as the basic availability check.

### LazyCodex MCP runtime

The OMO plugin declares the local LazyCodex MCP servers in its bundled
`.mcp.json`. Do not add direct `mcp_servers.lsp`, `mcp_servers.ast_grep`, or
direct custom Context7 MCP wiring to `/home/crunch/.codex/config.toml` unless
the user explicitly asks for a temporary diagnostic override. Prefer
OMO-provided MCP servers over custom/local alternatives when the OMO servers
work. The local plugin cache declaration launches the local stdio MCP servers
through `/bin/bash -c` so it can prepend WSL user tool directories and base WSL
system paths before execing `/home/crunch/.local/share/pnpm/bin/node`; this is
needed when the Codex desktop app-server starts with a Windows-heavy PATH.

The AST grep server requires the real `@ast-grep/cli` binary. On this
workstation it is installed globally with `pnpm` and in the OMO plugin cache
`node_modules`, exposing `sg` and `ast-grep` for the plugin-provided server.

The plugin server can be checked without Codex by sending an MCP `initialize`,
`notifications/initialized`, and `tools/list` sequence to the CLI:

<!-- markdownlint-disable MD013 -->

```bash
printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"codex-check","version":"0"}}}' \
  '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' \
| /home/crunch/.local/share/pnpm/bin/node /home/crunch/.codex/plugins/cache/sisyphuslabs/omo/0.1.0/mcp/ast_grep/dist/cli.js mcp
```

<!-- markdownlint-enable MD013 -->

The response should include `search` and `replace`. The runtime also requires
the real `@ast-grep/cli` binary. OMO Context7 is a remote streamable HTTP MCP
server at `https://mcp.context7.com/mcp`; use the OMO declaration instead of the
old local Context7 compat proxy when the remote initializes successfully.

Keep `plugins."omo@sisyphuslabs".mcp_servers.git_bash` disabled in the WSL
Codex config until OMO `git_bash` supports this host. Its current server exits
on non-Windows hosts before returning tools, so it does not satisfy the
"prefer OMO when it works" rule here.

Keep the top-level `mcp_servers.grep_app` entry pointed at
`/home/crunch/.codex/mcp/grep-app-compat/grep-app-compat.mjs`, with
`plugins."omo@sisyphuslabs".mcp_servers.grep_app` disabled, until the hosted
`https://mcp.grep.app` server handles broad searches reliably. The compatibility
server preserves the `grep_app/searchGitHub` namespace, tries the hosted MCP
first, then falls back to authenticated `gh search code` when the hosted server
times out or returns HTML.

Keep redundant or noisy LazyCodex LSP entries disabled in
`/home/crunch/.codex/lsp-client.json`: Deno, ESLint, Pyright, Ty, Ruff, Svelte,
Astro, bash-ls, terraform-ls, and Prisma. The default diagnostic surface should
prefer the active non-duplicated servers: TypeScript/JavaScript, Oxlint, Biome,
BasedPyright, bash, Terraform, and Dockerfile plus the other single-language
servers that do not overlap.

After changing MCP dependencies, start a fresh Codex thread or restart/reload
Codex before checking whether namespaces are available, because tool assembly
happens at thread startup.

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
omo help
```

Some tools are exposed through MCP or plugins rather than plain CLI commands.
For those, verify them through the Codex config or the tool list in the running
session. `codex plugin list` confirms that OMO is installed and enabled, but a
fresh thread is still needed to confirm that the `lsp` and `ast_grep` MCP
namespaces mounted.
