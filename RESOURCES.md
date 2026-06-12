# agents-toolkit Resources

Use these references when changing this repo or validating its docs.

## Repo-Owned Contracts

- [README.md](README.md): project map, current surface area, and common checks.
- [docs/wsl-shim.md](docs/wsl-shim.md): Windows-to-WSL Codex shim contract.
- [docs/skills-updater.md](docs/skills-updater.md): skills updater contract.
- [docs/codex-cli-tooling.md](docs/codex-cli-tooling.md): companion-tool
  operator notes.
- [docs/shell-setup.md](docs/shell-setup.md): shell setup and terminal-addon
  parity.
- [docs/discrawl-wiretap.md](docs/discrawl-wiretap.md): local Discrawl archive
  setup.
- [CONTEXT.md](CONTEXT.md): shared local terms.
- [AGENTS.md](AGENTS.md): workstation-local agent instructions for this repo.
- [RTK.md](RTK.md): local shell command wrapper rule.

## WSL And Windows Interop

- [Microsoft WSL documentation](https://learn.microsoft.com/windows/wsl/)
  Use for: `wsl.exe`, distro selection, `WSLENV`, Windows/WSL path boundaries,
  and WSL management from Windows.
- [`wslpath` manual behavior](https://manpages.ubuntu.com/manpages/noble/man1/wslpath.1.html)
  Use for: path conversion expectations when local direct conversion is not
  enough.
- [PowerShell documentation](https://learn.microsoft.com/powershell/)
  Use for: `.ps1` behavior, execution policy, `Start-Process`, profiles, and
  native process invocation.

## Codex And Companion Tools

- [OpenAI Developer Docs](https://platform.openai.com/docs)
  Use through the configured OpenAI Developer Docs MCP when available.
- [Codex docs](https://developers.openai.com/codex)
  Use for current Codex product behavior when local docs or tools drift.
- [Evo repository](https://github.com/evo-hq/evo)
  Use for Evo CLI/plugin release behavior.
- [CodSpeed docs](https://docs.codspeed.io/)
  Use for CodSpeed local executor and hosted run behavior.
- [Context7 MCP](https://context7.com/)
  Use through the OMO-provided MCP server when it initializes.

## Skills And Local Agent Workflows

- [Model Context Protocol](https://modelcontextprotocol.io/)
  Use for MCP lifecycle and tool server behavior.
- [Skills package source conventions](https://github.com/mattpocock/skills)
  Use only when maintaining skills that originate from that ecosystem.
- [LazyCodex / OMO plugin cache](https://github.com/sisyphuslabs)
  Use as external context only when the user asks for OMO/plugin work. Do not
  edit installed plugin/cache copies as part of ordinary agents-toolkit changes.

## Shell Addons

- [Atuin documentation](https://docs.atuin.sh/)
  Use for: history database, search UI, sync, key bindings, config, and
  `atuin init` flags.
- [fzf README](https://github.com/junegunn/fzf/blob/master/README.md)
  Use for: fuzzy finder behavior, shell key bindings, completion, and
  environment variables such as `FZF_CTRL_T_COMMAND`.
- [ble.sh README](https://github.com/akinomyoga/ble.sh)
  Use for: WSL Bash line editing, syntax highlighting, enhanced completion,
  auto suggestions, and fzf compatibility notes.
- [zoxide README](https://github.com/ajeetdsouza/zoxide)
  Use for: `zoxide init`, `z`, `zi`, directory ranking, shell hooks, and
  fzf-backed interactive directory selection.
- [Starship configuration docs](https://starship.rs/config/)
  Use for: prompt format, modules, palettes, `STARSHIP_CONFIG`, and cross-shell
  prompt behavior.
- [PSReadLine docs](https://learn.microsoft.com/powershell/module/psreadline/)
  Use for: PowerShell line editing, history predictions, list view, bell style,
  and edit mode.
- [PSFzf README](https://github.com/kelleyma49/PSFzf)
  Use for: how PowerShell wraps fzf and maps it into PSReadLine key chords.

## Local Learning Artifacts

- [lessons/0001-terminal-addon-ownership.html](lessons/0001-terminal-addon-ownership.html):
  interactive terminal-addon ownership lesson.
- [reference/0001-terminal-addons-cheatsheet.html](reference/0001-terminal-addons-cheatsheet.html):
  terminal-addon cheatsheet.
