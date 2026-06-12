# Docs-First Tool Use

Use this before a loop proposes or uses a CLI, script, tool, Worktree thread, subagent, custom agent, automation, plugin, MCP server, connector, or config change.

## Rule

Read the relevant docs or local help before designing or invoking the surface.

This is required for:

- Commands that create, edit, install, schedule, spawn, authenticate, change permissions, or touch config.
- Codex surfaces such as Worktree threads, subagents, custom agents, plugins, skills, automations, MCP/connectors, permissions, and worktrees.
- Third-party CLIs, libraries, framework APIs, or services where syntax or behavior may vary by version.
- Repo-local scripts when the script contract is not already obvious.

This is not required for routine shell inspection such as `pwd`, `ls`, `rg`, `sed`, `wc`, `git status`, or targeted file reads, unless using unusual flags or destructive behavior.

## What Counts

Prefer the closest reliable source:

- Local `--help`, `help`, man page, README, AGENTS.md, or script source for repo-local tools.
- Official product documentation for Codex, OpenAI, framework, cloud, package-manager, or service behavior.
- Installed package docs or type/schema definitions when local version matters.
- Existing project runbooks when they define the expected workflow.

For technical APIs or rapidly changing product behavior, prefer official docs over blogs or memory. If docs cannot be read, say what is unknown and choose the conservative path.

## Required Output

When docs affect the loop design, include a short record:

```text
Docs checked:
- Codex subagents docs: sandbox inheritance and approval behavior.
- codex plugin --help: reinstall command shape.
```

In a loop spec, put the same record in `docs_checked`.

## Subagent And Custom Agent Gate

Before creating or spawning subagents:

1. Read the current subagents docs or local custom-agent reference.
2. Confirm whether the child needs read-only, artifact-only, or code-editing access.
3. Confirm the active parent runtime permits that access.
4. Record the permission model in the subagent plan.
5. Record the spawn policy. Role-specific workers, reviewers, explorers, and
   custom-agent style assignments use non-full-history spawns with role/model
   intent in the message. In the current `spawn_agent` tool, omit
   `fork_context` or set `fork_context: false`; on tool surfaces that use
   `fork_turns`, set `fork_turns: "none"`. Full-history forks inherit the
   parent agent type, model, and reasoning effort, so do not combine them with
   `agent_type`, `model`, or `reasoning_effort` overrides.

If the child must write and the active parent turn is read-only, stop and report the mismatch instead of spawning a writer.
