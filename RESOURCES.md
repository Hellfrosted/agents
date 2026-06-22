# agents-toolkit Resources

Use these references when changing this repo or validating its docs.

## Repo-Owned Contracts

- [README.md](README.md): public project overview and entry points.
- [docs/README.md](docs/README.md): maintained documentation index.
- [docs/maintainer-guide.md](docs/maintainer-guide.md): source-first workflow,
  common checks, install freshness checks, and design rationale.
- [docs/wsl-shim.md](docs/wsl-shim.md): Windows-to-WSL Codex shim contract.
- [docs/skills-updater.md](docs/skills-updater.md): skills updater contract.
- [docs/sk-up-go-port.md](docs/sk-up-go-port.md): portable Go skills updater
  implementation brief, CLI protocol, validation matrix, and progress ledger.
- [docs/agent-native-boundary.md](docs/agent-native-boundary.md):
  BuilderIO/Agent-Native ownership boundary for plan/recap skills and
  Visual Canvas.

## WSL And Windows Interop

- [Microsoft WSL documentation](https://learn.microsoft.com/windows/wsl/)
  Use for: `wsl.exe`, distro selection, `WSLENV`, Windows/WSL path boundaries,
  and WSL management from Windows.
- [`wslpath` manual behavior](https://manpages.ubuntu.com/manpages/noble/man1/wslpath.1.html)
  Use for: path conversion expectations when local direct conversion is not
  enough.
- [PowerShell documentation](https://learn.microsoft.com/powershell/)
  Use for: Windows wrapper regression harnesses and native process invocation
  from PowerShell. The skills updater implementation itself is Go.

## Codex

- [OpenAI Developer Docs](https://platform.openai.com/docs)
  Use through the configured OpenAI Developer Docs MCP when available.
- [Codex docs](https://developers.openai.com/codex)
  Use for current Codex product behavior when local docs or tools drift.

## Skills And Local Agent Workflows

- [Model Context Protocol](https://modelcontextprotocol.io/)
  Use for MCP lifecycle and tool server behavior.
- [Skills package source conventions](https://github.com/mattpocock/skills)
  Use only when maintaining skills that originate from that ecosystem.
- [BuilderIO Agent-Native skills](https://github.com/BuilderIO/skills)
  Use for the installed Agent-Native plan, recap, orchestration, and managed
  instruction-block workflows. Do not duplicate those skills in this repo.
- [LazyCodex / OMO plugin cache](https://github.com/sisyphuslabs)
  Use as external context only when the user asks for OMO/plugin work. Do not
  edit installed plugin/cache copies as part of ordinary agents-toolkit changes.
