# Evo End To End Reference

## Companion Skill Sources

| Skill | Use in this workflow | GitHub source if missing |
| --- | --- | --- |
| `$grill-me` | Clarify ambiguous goals, constraints, non-goals, success metrics, or forbidden changes. | `https://github.com/mattpocock/skills.git` at `skills/productivity/grill-me/SKILL.md` |
| `$grill-with-docs` | Challenge a plan against project terminology, `CONTEXT.md`, or ADRs, and update docs as decisions crystallize. | `https://github.com/mattpocock/skills.git` at `skills/engineering/grill-with-docs/SKILL.md` |
| `$improve-codebase-architecture` | Decompose architecture or testability work before choosing an Evo metric. | `https://github.com/mattpocock/skills.git` at `skills/engineering/improve-codebase-architecture/SKILL.md` |
| `$evo discover`, `$evo optimize`, `$evo finetuning`, `$evo infra-setup` | Evo plugin skills installed by `evo install <host>`. Use the bundle whose `evo_version` matches `evo --version`. | `https://github.com/evo-hq/evo.git` under `plugins/evo/skills/` |

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

## Evo v0.5.2 Notes

- v0.5.2 improves the optimize meta controller: meta ticks keep a journal,
  appended prompt directives accumulate instead of overwriting each other, and
  the meta can harden verifier prompts during a run.
- Meta and hard implement/revise agents follow the session model instead of a
  stale pinned model; easy briefs still route to the lighter model.
- Keep the Codex CLI and plugin bundle in lockstep with
  `evo update codex --version 0.5.2 --trust-hooks`, then verify with
  `evo doctor codex`.
- If Codex is pinned to a stale local marketplace such as
  `evo-hq-0.5.0-hookdrain`, remove that marketplace and add
  `evo-hq/evo --ref v0.5.2` before reinstalling.
- SDK packages such as `evo-hq-agent`, `@evo-hq/evo-agent`, and `@evo-hq/pi-evo`
  should match the 0.5.2 line when SDK instrumentation is used.

## Evo v0.5.1 Notes

- v0.5.1 moves the durable hook binary to `~/.evo/bin` and leaves a host-plugin
  fallback at the hook path, so host plugin cache rebuilds no longer delete the
  real binary and trigger `SessionStart hook (failed): exit 127`.
- `evo install codex` can trust Evo hooks during install; use
  `--no-trust-hooks` only when the user wants to review them manually through
  `/hooks`.
- `evo doctor codex` verifies hook trust and catches plugin updates that changed
  `hooks.json` enough to invalidate trust.

## Evo v0.5.0 Notes

- v0.5.0 adds `$evo finetuning` for SFT, LoRA, preference optimization, RFT, and
  RL training moves.
- Training runs must keep the held-out benchmark out of training data, perform
  literature research before the first train experiment, and use smoke-run
  validation before full budget.
- New workspaces should set a realistic per-experiment timeout at init with
  `--per-exp-timeout <seconds>` or later with
  `evo config set per-exp-timeout <seconds>`.
- `task-skills` records task-category skills, such as `finetuning`, that
  subagents should load on demand. Inspect it with `evo config get task-skills`.
- `$evo optimize` requires resource-bound round sizing. Read Evo's
  `sizing-the-round.md` before passing a concrete `subagents=N`.
- `evo wait` has process, log, GPU, and ideator selectors for long-running work.
- `evo abort` stops the experiment subprocess tree cross-platform, including
  detached benchmark or training children.
- The dashboard supports live log tailing, per-experiment annotations, and
  `EVO_DASHBOARD_HOST` for binding on cloud or Modal hosts.

## Evo v0.4.5 Notes

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

## Evo v0.4.4 Notes

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
- `$evo optimize` defaults to autonomous, subagents-only operation. The user can
  override either explicitly or via Evo defaults/config.
- `evo direct "<text>" --wait` expects an agent to acknowledge delivered
  directives with `evo ack <event_id>`.
- Use `evo gc` to clean worktrees, pool slots, and remote sandboxes.
- Use `evo config show`, `evo config backend show`,
  `evo config runtime show`, and `evo env show` to inspect setup before changing
  it.
