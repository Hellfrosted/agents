# Evo End To End Reference

## Companion Skill Sources

| Skill | Use in this workflow | GitHub source if missing |
| --- | --- | --- |
| `$ask-matt` | Choose the Matt Pocock flow when a request is not yet clearly planning, triage, implementation, or codebase health work. | `https://github.com/mattpocock/skills.git` at `skills/engineering/ask-matt/SKILL.md` |
| `$grilling` | Clarify ambiguous goals, constraints, non-goals, success metrics, ownership, or forbidden changes. This is the internal interview skill used by the old grill wrappers. | `https://github.com/mattpocock/skills.git` at `skills/productivity/grilling/SKILL.md` |
| `$codebase-design` | Use Matt's deep-module vocabulary when evaluating module interfaces, seams, leverage, locality, and testability. | `https://github.com/mattpocock/skills.git` at `skills/engineering/codebase-design/SKILL.md` |
| `$domain-modeling` | Pair with `$grilling` for codebase planning when the workflow needs to add or change `CONTEXT.md`, ADRs, or domain terms. | `https://github.com/mattpocock/skills.git` at `skills/engineering/domain-modeling/SKILL.md` |
| `$improve-codebase-architecture` | Decompose architecture or testability work before choosing an Evo metric. | `https://github.com/mattpocock/skills.git` at `skills/engineering/improve-codebase-architecture/SKILL.md` |
| `$implement` | Execute a single PRD or issue after the plan has been split and the implementation scope is fresh. | `https://github.com/mattpocock/skills.git` at `skills/engineering/implement/SKILL.md` |
| `$tdd` | Drive implementation test-first when an Evo setup task needs behavior locked by public-interface tests. | `https://github.com/mattpocock/skills.git` at `skills/engineering/tdd/SKILL.md` |
| `evo:discover`, `evo:optimize`, `evo:finetuning`, `evo:infra-setup` | Evo plugin skills installed by `evo install <host>`. Use the bundle whose `evo_version` matches `evo --version`. | `https://github.com/evo-hq/evo.git` under `skills/`; npm mirrors may expose the same files under `npm/skills/`. |

The current `mattpocock/skills` package is installed through the Skills CLI and
may record `pluginName: "mattpocock-skills"` in the universal skill lockfile.
That metadata does not make it a Codex plugin; resolve missing Matt skills from
`https://github.com/mattpocock/skills.git`, not from the Codex plugin cache.

## Optimize Presets

Use presets only as fallbacks when the benchmark resource profile is unknown and
the user gave no exact values:

- **tiny**: `subagents=3 budget=5 stall=2`
- **small**: `subagents=3 budget=8 stall=3`
- **medium**: `subagents=4 budget=10 stall=4`
- **big**: `subagents=5 budget=14 stall=5`
- **huge**: `subagents=8 budget=20 stall=6`

Default to **medium** only when the benchmark is light, isolated, and no better
sizing signal is available. Reduce `subagents` to 1 for exclusive resources such
as a GPU, fixed port, shared database, or serialized fixture. Cap pool runs at
the pool slot count. Use **tiny** or **small** when the editable scope is narrow
or risky. Use **big** or **huge** only when the metric is stable, the baseline is
repeatable, and the approved scope can absorb broader exploration.

## Workspace And Backend Notes

- If using an existing Evo workspace from before v0.4.0, migrate host metadata
  with `evo host show`; if it prints `<not set>`, run `evo host set codex`.
- Local default backend: worktree.
- Faster local reuse: pool backend with a fixed workspace list.
- Remote backend: configure the provider first using Evo's `infra-setup`
  guidance for Modal, E2B, Daytona, AWS, Azure, SSH, manual, or custom
  providers.
- Runtime commands/env belong in `evo config runtime ...` and `evo env ...`, not
  hard-coded into benchmark scripts.
- Override individual calls with `evo run <exp_id> --timeout <seconds>` only
  when the configured per-experiment timeout is not appropriate.

## Evo v0.6.2 Notes

- Keep the Codex CLI and plugin bundle in lockstep. Report the update command
  first; run `evo update codex --version 0.6.2` only after explicit user
  approval. Add `--trust-hooks` only when the user explicitly approves trusting
  hooks without interactive review, then verify with `evo doctor codex`.
- The installed Evo plugin skills advertise `evo_version: 0.6.2`; stop and fix
  CLI/plugin drift before using `evo:discover`, `evo:optimize`,
  `evo:finetuning`, or `evo:infra-setup`.
- If `evo doctor codex` reports Codex hooks still referencing
  `CLAUDE_PLUGIN_ROOT`, do not launch Evo workflows. Report the required repair
  command, `evo install codex --force`, for explicit user approval because it
  mutates the active Codex plugin/cache install.
- `evo:optimize` remains the canonical post-discover loop even when resource
  constraints force `subagents=1`. Load `$evo discover`'s sizing reference
  (`skills/discover/references/sizing-the-round.md` in the installed Evo plugin
  bundle) before choosing width, and use `evo wait` or bounded liveness checks
  instead of unbounded polling for long-running subagents, training, or
  evaluations.
- `autonomous` and `subagents-only` default to on in the installed
  `evo:optimize` skill unless the current user instruction or stored Evo
  defaults resolve them differently. Resolve and arm both modes before starting
  an optimize loop.

## Evo v0.5.2 Compatibility

- On the v0.5.2 line, keep the Codex CLI and plugin bundle in lockstep. Report
  the update command first; run `evo update codex --version 0.5.2` only after
  explicit user approval. Add `--trust-hooks` only when the user explicitly
  approves trusting hooks without interactive review, then verify with
  `evo doctor codex`.
- If Codex is pinned to a stale local marketplace such as
  `evo-hq-0.5.0-hookdrain`, remove that marketplace and add
  `evo-hq/evo --ref v0.5.2` before reinstalling.
- SDK packages such as `evo-hq-agent`, `@evo-hq/evo-agent`, and `@evo-hq/pi-evo`
  should match the 0.5.2 line when SDK instrumentation is used.

## Legacy Hook Recovery

- v0.5.1 moves the durable hook binary to `~/.evo/bin` and leaves a host-plugin
  fallback at the hook path, so host plugin cache rebuilds no longer delete the
  real binary and trigger `SessionStart hook (failed): exit 127`.
- `evo install codex` can trust Evo hooks during install; use
  `--no-trust-hooks` only when the user wants to review them manually through
  `/hooks`.
- `evo doctor codex` verifies hook trust and catches plugin updates that changed
  `hooks.json` enough to invalidate trust.

- v0.4.5 fixes Codex hook installation for Codex 0.130+ by registering and
  staging the plugin under the owner-name path Codex resolves (`evo@evo-hq`) and
  validating `evo-hook-drain` in `evo doctor codex`.
- Existing Codex installs with exit-127 hook failures do not self-heal with
  `evo update`. Recover with
  `uv tool install --force evo-hq-cli && evo install codex --force`.
- `evo install codex --force` stages `evo-hook-drain` into the Codex plugin
  cache and removes stale legacy registrations. It may leave hooks untrusted;
  trust them through `/hooks` or run `evo install codex --trust-hooks` only when
  the user explicitly approves skipping hook review.

## Backend Notes

- `evo init --host <claude-code|codex|cursor|opencode|openclaw|hermes|pi|generic>`
  is required for new workspaces. For this skill on Codex, use `codex`.
- New workspaces default to the `pareto_per_task` frontier strategy. Existing
  workspaces keep their configured strategy.
- Local execution has `worktree` and `pool` backends. Pool mode is useful when
  setup is expensive, but warm workspace state should stay out of commits.
- Pool mode defaults to `commit_strategy=tracked-only`; subagents must `git add`
  new source files and pass `--i-staged-new-files yes` to `evo run`.
- Remote experiments can run through Modal, E2B, Daytona, AWS, Azure, SSH,
  manual, or custom providers. Treat provider SDK installation, credentials, and
  cloud allocation as explicit user-approved setup.
- In remote mode, subagent briefs must state the experiment id and require
  `--exp-id <id>` on every `evo bash/read/write/edit/glob/grep` command.
- Backend provider credentials and benchmark runtime environment are separate.
  Configure benchmark variables with `evo env`, and do not copy secrets into
  worktrees or docs.
- `evo run <exp_id> --check` validates benchmark/gate wiring without committing,
  evaluating, or consuming retry budget.
- `evo:optimize` defaults to autonomous, subagents-only operation. The user can
  override either explicitly or via Evo defaults/config.
- `evo direct "<text>" --wait` expects an agent to acknowledge delivered
  directives with `evo ack <event_id>`.
- Use `evo gc` to clean worktrees, pool slots, and remote sandboxes.
- Use `evo config show`, `evo config backend show`,
  `evo config runtime show`, and `evo env show` to inspect setup before changing
  it.
