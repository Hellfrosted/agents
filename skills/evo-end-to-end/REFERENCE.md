# Evo End To End Reference

## Authoritative Sources

For release-specific behavior, prefer the installed plugin bundle whose
`evo_version` matches `evo --version`, then the matching upstream tag:

- CLI: `plugins/evo/src/evo/cli.py`
- Host installers: `plugins/evo/src/evo/host_install/`
- Driving skills: `plugins/evo/skills/{discover,optimize,report,ship}/`
- Round sizing:
  `plugins/evo/skills/discover/references/sizing-the-round.md`
- CLI quick reference: `plugins/evo/skills/references/cli-quick-reference.md`

The npm mirror lives under `plugins/evo/npm/skills/`; do not use historical
`skills/` or `npm/skills/` repository-root paths for v0.7.

## Companion Skill Sources

| Skill | Use in this workflow | GitHub source if missing |
| --- | --- | --- |
| `$ask-matt` | Route an unclear engineering workflow. | `https://github.com/mattpocock/skills.git` at `skills/engineering/ask-matt/SKILL.md` |
| `$grilling` | Clarify ambiguous goals, constraints, success metrics, ownership, and forbidden changes. | `https://github.com/mattpocock/skills.git` at `skills/productivity/grilling/SKILL.md` |
| `$codebase-design` | Evaluate module interfaces, seams, leverage, locality, and testability. | `https://github.com/mattpocock/skills.git` at `skills/engineering/codebase-design/SKILL.md` |
| `$domain-modeling` | Pair domain vocabulary, `CONTEXT.md`, and ADR work with planning. | `https://github.com/mattpocock/skills.git` at `skills/engineering/domain-modeling/SKILL.md` |
| `$improve-codebase-architecture` | Decompose architecture or testability work before choosing an Evo metric. | `https://github.com/mattpocock/skills.git` at `skills/engineering/improve-codebase-architecture/SKILL.md` |
| Evo discover, optimize, report, ship, finetuning, infra-setup | Use the Evo bundle whose `evo_version` matches the CLI. | `https://github.com/evo-hq/evo.git` under `plugins/evo/skills/` |

The Matt Pocock package may record `pluginName: "mattpocock-skills"` in a
universal skill lockfile. Resolve missing skills from its source repository,
not from the Codex plugin cache. Do not invent or retain a companion skill name
that the current source tree no longer contains.

## Evo v0.7 Workflow Notes

- Keep the CLI and host plugin bundle in lockstep. Report
  `evo update <host> --version 0.7.0` when repair is needed; run it only after
  explicit user approval. Verify with `evo doctor <host>`.
- Host invocation differs. Claude Code uses `/evo:discover`,
  `/evo:optimize`, `/evo:report`, and `/evo:ship`. Codex advertises the plugin
  through `$evo`; select the corresponding driving skill rather than copying
  Claude slash-command syntax. Other hosts use their menu or natural language.
- `evo init` requires `--host` and a positive `--per-exp-timeout`. The timeout
  becomes the workspace default; override one call with `evo run <exp_id>
  --timeout N`, or change the workspace with `evo config set per-exp-timeout
  N`.
- `evo run <exp_id> --check` runs the configured benchmark and gates and writes
  artifacts to a fresh check-attempt directory. It does not commit, evaluate,
  or consume retry budget. Treat it as workspace-writing and potentially
  compute-expensive.
- `evo run <exp_id>` is terminal on success: it evaluates and commits the
  experiment. Do not follow it with `evo done`; that command is for manually
  recorded scores outside `evo run`.
- After optimization, use the report skill to summarize recorded state without
  spending compute. Use ship only after winner selection and approval; it
  distills the best valid experiment rather than landing the raw worktree diff.

## Sizing

Do not maintain independent tiny/small/medium/big preset tables here. Read the
matching plugin's `skills/discover/references/sizing-the-round.md` immediately
before choosing `subagents` and `budget`.

Its governing constraint is resource width:

- worktree and gitdir runs share local resources;
- pool width cannot exceed its slot count;
- remote width is bounded by provider quota, concurrency, and approved cost;
- serialize an exclusive GPU, fixed port, shared database/fixture, or unstable
  timing harness unless the benchmark isolates it;
- size budget from per-run cost and noise, then state the chosen width/budget
  and reason.

This is Evo worker sizing, not Codex subagent delegation.

## Workspace And Backend Notes

- Default local backend: `worktree`.
- `gitdir` relocates Git metadata away from `.git`. Use it where `.git` paths
  are forbidden, especially Claude Science.
- Claude Science defaults new workspaces to `gitdir` after `evo install
  claude-science`. The CLI cannot publish skills into Claude Science's catalog;
  the in-session `evo-autoresearch-claude-science-setup` skill from
  `evo-hq/evo-claude-science` does that through `host.skills`.
- Pool reuses fixed local workspaces. Respect the configured slot count and
  keep warm workspace state out of commits.
- Remote experiments can use Modal, E2B, Daytona, AWS, Azure, or SSH. Provider
  packages, credentials, and allocations require explicit approval.
- Runtime commands and benchmark variables belong in `evo config runtime ...`
  and `evo env ...`; do not copy secrets into worktrees, scripts, or docs.
- Inspect with `evo config show`, `evo config backend show`, `evo config runtime
  show`, and `evo env show` before changing configuration.
- `evo direct "<text>" --wait` delivers a directive; an agent receiving an
  `[EVO DIRECTIVE id=...]` banner must acknowledge it with `evo ack <event_id>`.
- `evo gc` mutates workspace state by cleaning worktrees, pool slots, and remote
  sandboxes; keep it behind explicit approval.

## Hook Repair

Keep one current recovery path instead of release-by-release hook archaeology:

1. Inspect with `evo doctor <host>`.
2. If the CLI and plugin are mismatched or the durable hook binary is missing,
   report `uv tool install --force evo-hq-cli && evo install <host> --force`
   for explicit approval.
3. For Codex, upstream `evo install codex` trusts Evo hooks by default. This
   repository's approval boundary is stricter: do not trust hooks implicitly.
   Use `--no-trust-hooks`, let the user review them through `/hooks`, and only
   accept automatic trust when the user explicitly authorizes it.
4. Re-run `evo doctor <host>` and record the result.

Do not preserve obsolete marketplace-removal recipes, `CLAUDE_PLUGIN_ROOT`
diagnostics, or pre-0.5 hook cache layouts in the active workflow. Consult the
matching historical tag only when repairing a deliberately pinned old install.

## Repository-Stricter Landing Override

Upstream Evo ship may propose creating a fresh branch, pushing, opening a pull
request, or merging. In this repository those actions remain report-only until
the user explicitly authorizes each one. Work on the current branch by default;
never create or switch branches merely because the upstream ship skill says to.
