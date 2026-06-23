# Mission: agents-toolkit

This repository keeps the workstation's Codex-adjacent tooling understandable,
repairable, and reproducible.

## Purpose

Maintain the local source of truth for:

- Windows-to-WSL Codex launch behavior;
- T3code app-server proxy behavior;
- globally installed Codex skill maintenance wrappers;
- repo-owned Visual Canvas plugin source;
- repo-owned local skills;
- repo-owned Codex hook source and focused verification docs;
- promotion and wrapper checks for repo-owned workstation tooling.

The repo should answer what is true now, how to verify it, and where to make the
next source change. Historical notes belong only where they explain a current
design constraint or repair path.

## Success Looks Like

- A maintainer can identify which file owns a behavior before editing it.
- The active workstation install can be checked against this repo.
- Script and workflow changes have matching human-facing docs.
- Local skills describe current invocation and safety rules without relying on
  remembered installed paths.
- Docs are task-oriented enough to debug repo-owned tool failures without
  generic setup archaeology.

## Constraints

- Keep source edits in this repo first. Copy into active install locations only
  when the user asks for install, repair, or republish work.
- Preserve user changes and unrelated local state.
- Prefer focused docs over duplicated flag lists and stale snapshots.
- Keep machine-specific operational facts only when they are needed to run or
  repair this workstation.
- Do not store secrets, tokens, private personal data, raw chat exports, or
  credential material in docs, examples, commits, or skill text.

## Out Of Scope

- Redesigning the user's Codex workflow without a concrete request.
- Maintaining external plugin/cache contents as part of ordinary repo docs.
- Replacing installed global skills unless the user asks.
- Turning backup payloads into living documentation.
- Maintaining workstation dotfiles, shell setup, automation definitions, or
  restore payloads; those belong in the dotfiles repository.
