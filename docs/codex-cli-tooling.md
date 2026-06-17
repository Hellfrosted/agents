# Codex companion tools

This page lists companion tools installed on this machine for Codex workflows,
plus repo-local wrappers that maintain those tools. It is an operator map, not
the source of truth for every wrapper flag. Use focused docs for detailed
contracts:

- [Codex WSL Shim](wsl-shim.md)
- [Skills Updater](skills-updater.md)
- [Discrawl Vesktop Wiretap](discrawl-wiretap.md)

## Main list

- Evo: experiment runner for optimization work.
- RTK: command wrapper agents use for shell commands.
- ICM: persistent memory for durable Codex context.
- OpenAI Developer Docs MCP: official OpenAI docs server for API and product
  docs lookup.
- CodSpeed: CLI plus hosted MCP for performance runs, comparisons, and
  flamegraph analysis.
- LazyCodex: installed as the `omo@sisyphuslabs` Codex plugin. The `omo` name
  here is the Codex marketplace/plugin identifier, not the full OpenCode OMO
  install.
- LazyCodex local MCP: LSP, AST grep, Context7, and grep.app search are provided
  by the LazyCodex plugin when their servers work on this host.
- LazyCodex CLI: local `lazycodex-ai` / `lazycodex` entrypoints for Codex
  plugin install and update helpers. Do not keep a WSL `omo` launcher; full OMO
  plus OpenCode belongs on Windows for this workstation.
- Discrawl: local Discord cache archive/search for Vesktop wiretap-only use.
- Skills updater: Windows wrappers for checking and updating globally installed
  Codex skills.
- Agent Browser: CLI-driven Chrome/Chromium automation for browser QA,
  screenshots, and page interaction from agent workflows.

## Codex app and plugin surfaces

These surfaces are part of the active workstation toolchain even when this repo
does not own their implementation:

- `google-calendar@openai-curated`: installed and enabled for connected Google
  Calendar scheduling, availability, and agenda work.
- `google-drive@openai-curated`: installed and enabled for connected Drive,
  Docs, Sheets, and Slides work.
- `chrome@openai-bundled`: installed and enabled for workflows that require the
  user's existing Chrome state.
- `browser@openai-bundled`: enabled in Codex config for in-app browser
  automation and local web target checks. In Windows + WSL Codex sessions, do
  not rely on `codex plugin list` alone: the Browser config can exist while the
  `@Browser` skill/tool surface is absent from the Linux-side thread. Verify the
  running thread's tool list or use `@Browser` directly in a fresh thread.
- `computer-use@openai-bundled`: installed and enabled for Windows desktop app
  control.
- `agent-loop-designer@personal`: installed and enabled for turning recurring
  Codex tasks into repeatable loops or worktree-thread workflows.
- `tabby` MCP: configured as a local HTTP MCP server at
  `http://172.27.48.1:3001/mcp`.
- `node_repl` MCP: app-provided JavaScript runtime used for Node-backed
  inspection and browser automation helpers.

Check these with:

```bash
codex plugin list
codex mcp list
```

Some plugin and MCP namespaces are assembled only when a Codex thread starts. If
config changed, start a fresh thread or reload Codex before treating a missing
namespace as a broken install. Browser is the main exception on this
workstation: because Codex Desktop is running the workspace through WSL, confirm
it from the actual thread tool/skill surface, not only from `codex plugin list`.

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

When the task is to run or inspect Windows-side shell behavior, prefer Tabby
MCP over brute-force WSL interop attempts. Use Tabby MCP to target existing
PowerShell or Command Prompt sessions, read terminal buffers, send interactive
input, stop running commands, split/focus panes, or use SFTP through an active
SSH session. This is especially useful when command success depends on the
loaded Windows profile, Tabby session state, an interactive prompt, or a
Windows-only executable that is awkward to launch reliably from WSL.

Keep direct WSL interop for simple, deterministic, non-interactive commands
that are already documented in this repo, such as path conversion, opening a
folder or URL, checking a Windows symlink, or invoking a known PowerShell test
script through WSL init.

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
  `evo update codex --version 0.5.3 --trust-hooks` for the current line, then
  verify with `evo doctor codex`.
- If Codex is pinned to a stale local `evo-hq` marketplace, remove that
  marketplace source and add `evo-hq/evo --ref v0.5.3` before reinstalling.
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

The repo includes Windows wrappers for installed global skill maintenance:

- `bin/sk-up.cmd`: short command names and flags.
- `bin/skills-updates.cmd`: long command names and flags.
- `bin/skills-updates.ps1`: implementation.

Use it to list installed skills, compare global skills against upstream
content, open diffs, install changed or missing skills, install source URLs,
save or remove skips, and uninstall global skills. See
[Skills Updater](skills-updater.md) for the full command table, state paths,
lockfile behavior, and verification.

For Matt Pocock's skills, use the updater as a Skills CLI package manager, not
as a Codex plugin manager:

```powershell
bin\sk-up.cmd -i mattpocock/skills
bin\sk-up.cmd -g
```

The current package uses `/ask-matt` as the flow router. Its lockfile entries
may include `pluginName: "mattpocock-skills"` because the Skills CLI imports
from the upstream `.claude-plugin/plugin.json`; that does not mean a Codex
plugin was installed. Deprecated Matt skills, such as `request-refactor-plan`,
remain only if they were installed separately and should be removed explicitly
when no longer wanted.

## OpenAI Developer Docs MCP

The OpenAI Developer Docs MCP server is configured as `openaiDeveloperDocs`. Use
it for current OpenAI API, SDK, platform, and product docs.

Prefer it over web search for OpenAI-specific questions. Fetch the exact doc
page before quoting or summarizing details.

## Agent Browser

`agent-browser` is installed from the global `pnpm` bin directory at
`/home/crunch/.local/share/pnpm/bin/agent-browser`. Use it for CLI-driven
Chrome/Chromium browser automation when a task needs a real browser surface,
including page inspection, interaction, screenshots, and lightweight QA.

There is no repo-owned `skills/agent-browser/` source in this repository. Treat
Agent Browser as an installed CLI workflow: load the version-matched CLI skill
text with `agent-browser skills get ...` so the instructions match the installed
CLI version.

The installed version verified on 2026-06-17 was `agent-browser 0.28.0`.
`agent-browser install` reported Chrome for Testing `149.0.7827.115` already
installed under `/home/crunch/.agent-browser/browsers`. `agent-browser doctor`
exited successfully with 9 pass, 0 warn, and 0 fail; its pass/warn counts can
vary when the headless launch smoke check is slow.

Before using it, load the version-matched CLI skill text:

```bash
agent-browser skills get core
agent-browser skills get core --full
```

Useful smoke checks:

```bash
agent-browser --version
agent-browser skills list
agent-browser install
agent-browser doctor
agent-browser open https://example.com
agent-browser get title
agent-browser get url
agent-browser snapshot -i -c
agent-browser screenshot example.png
agent-browser close --all
```

`agent-browser screenshot <path>` may still write to its managed temporary
screenshot directory and print the actual saved path. When preserving evidence,
copy the printed file into the desired evidence directory and record that path.

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

Current repo caveat: `Hellfrosted/agents` is not enabled as a CodSpeed
repository. Local CodSpeed runs from this repo upload to a local-runs project
unless the repository is enabled on CodSpeed.

## LazyCodex

LazyCodex is installed as the `omo@sisyphuslabs` Codex plugin. Keep it enabled
unless the user explicitly asks to remove it. In WSL, treat `omo@sisyphuslabs`
as the Codex marketplace/plugin name only; the full OMO plus OpenCode setup is
Windows-only on this workstation.

The local CLI entrypoints are `lazycodex-ai` and `lazycodex`. Use
`lazycodex-ai version` or `lazycodex-ai --help` as basic availability checks.
Do not keep a WSL `omo` wrapper, because stale wrappers can dispatch to the
full `oh-my-openagent` package.

### LazyCodex MCP runtime

The LazyCodex plugin declares the local MCP servers in its bundled
`.mcp.json`. Do not add direct `mcp_servers.lsp`, `mcp_servers.ast_grep`, or
direct custom Context7 MCP wiring to `/home/crunch/.codex/config.toml` unless
the user explicitly asks for a temporary diagnostic override. Prefer
LazyCodex-provided MCP servers over custom/local alternatives when those servers
work. The local plugin cache declaration launches the local stdio MCP servers
with `node` from the plugin root (`cwd = "."` in the plugin `.mcp.json`).

The AST grep integration requires the real `@ast-grep/cli` binary. On this
workstation it is installed globally with `pnpm`, exposing `sg` and `ast-grep`.
LazyCodex v4.11.0 no longer ships the old bundled `ast-grep-mcp` path; it
provisions AST grep through the shared ast-grep skill / `sg` resolver flow
instead.

Check the current Codex plugin and AST grep runtime with:

```bash
lazycodex-ai version
codex plugin list | rg 'omo@sisyphuslabs|VERSION|Marketplace `sisyphuslabs`' -C 2
sg --version
```

On 2026-06-17 the installed LazyCodex plugin cache path was
`/home/crunch/.codex/plugins/cache/sisyphuslabs/omo/4.11.0` and `sg --version`
reported `ast-grep 0.43.0`.

LazyCodex Context7 is a remote streamable HTTP MCP server at
`https://mcp.context7.com/mcp`; use the plugin declaration instead of the old
local Context7 compat proxy when the remote initializes successfully.

Keep `plugins."omo@sisyphuslabs".mcp_servers.git_bash` disabled in the WSL
Codex config until LazyCodex `git_bash` supports this host. Its current server
exits on non-Windows hosts before returning tools, so it does not satisfy the
"prefer LazyCodex when it works" rule here.

Keep the top-level `mcp_servers.grep_app` entry pointed at
`/home/crunch/.codex/mcp/grep-app-compat/grep-app-compat.mjs`, with
`plugins."omo@sisyphuslabs".mcp_servers.grep_app` disabled, until the hosted
`https://mcp.grep.app` server handles broad searches reliably. The compatibility
server preserves the `grep_app/searchGitHub` namespace, tries the hosted MCP
first, then falls back to authenticated `gh search code` when the hosted server
times out or returns HTML.

Keep redundant or noisy LazyCodex LSP entries disabled in
`/home/crunch/.codex/lsp-client.json`: ESLint, Pyright, Ty, Ruff, Svelte,
Astro, bash-ls, terraform-ls, and Prisma. The default diagnostic surface should
prefer the active non-duplicated servers: TypeScript/JavaScript, Oxlint, Biome,
BasedPyright, bash, Terraform, and Dockerfile plus the other single-language
servers that do not overlap.

After changing MCP dependencies, start a fresh Codex thread or restart/reload
Codex before checking whether namespaces are available, because tool assembly
happens at thread startup.

For T3code sessions launched through `bin/codex-wsl.cmd`, the WSL runner allows
LazyCodex auto-update by default so Codex plugin updates do not require manual
remembering. It still disables telemetry and config migration by default while
leaving the plugin available:

```bash
OMO_CODEX_DISABLE_POSTHOG=1
OMO_CODEX_SEND_ANONYMOUS_TELEMETRY=0
OMO_DISABLE_POSTHOG=1
OMO_SEND_ANONYMOUS_TELEMETRY=0
LAZYCODEX_CONFIG_MIGRATION_DISABLED=1
OMO_CODEX_CONFIG_MIGRATION_DISABLED=1
```

Set `LAZYCODEX_AUTO_UPDATE_DISABLED=1` or `OMO_CODEX_AUTO_UPDATE_DISABLED=1`
explicitly only when diagnosing update-related startup issues.

Other T3code shim knobs:

- `CODEX_WSL_PROXY_IDLE_TIMEOUT_MS`: app-server idle timeout. The Windows shim
  defaults this to `1800000` milliseconds.
- `CODEX_WSL_PROXY_SKILLS_TIMEOUT_MS`: timeout before the WSL proxy returns a
  fallback `skills/list` response. The proxy defaults this to `2000`
  milliseconds.
- `CODEX_WSL_PROXY_DEBUG_LOG`: WSL path for proxy debug logs.
- `CODEX_WSL_SHIM_DEBUG`: print the Windows shim launch arguments.

See [Codex WSL Shim](wsl-shim.md) for the full shim contract, install shape,
path translation policy, and runtime verification.

Serena is not part of the current installed Codex toolchain on this workstation;
it has already been uninstalled.

## Quick checks

```bash
evo --version
rtk --help
icm --help
codex plugin list
codex mcp list
codspeed status
codspeed setup status
agent-browser --version
agent-browser doctor
actionlint --version
react-doctor --version
tokscale --version
lazycodex-ai --help
lazycodex-ai version
```

Some tools are exposed through MCP or plugins rather than plain CLI commands.
For those, verify them through the Codex config or the tool list in the running
session. `codex plugin list` confirms that the LazyCodex plugin is installed
and enabled, but a fresh thread is still needed to confirm that the `lsp` and
`ast_grep` MCP namespaces mounted.

For Browser use, check the current thread directly. Official Codex docs say the
in-app browser is controlled through the Browser plugin and `@Browser`
(`https://developers.openai.com/codex/app/browser`), and Codex plugins can
bundle skills, apps, and MCP servers
(`https://developers.openai.com/codex/plugins`). Public OpenAI Codex issue
reports also show Windows + WSL sessions where Browser feature/config state is
enabled but the Browser use skill or required Node REPL tool surface is missing
from the actual session (`https://github.com/openai/codex/issues/19365`,
`https://github.com/openai/codex/issues/21440`). When that happens, use the
installed `agent-browser` CLI as the fallback for browser QA and load its current
instructions with `agent-browser skills get core`.
