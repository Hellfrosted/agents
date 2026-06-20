# Use OS-native updater state paths

The portable skills updater keeps installed skills and `.skill-lock.json` under
`AGENTS_HOME`, falling back to the user's `.agents` directory, but stores updater
cache and skip state in OS-native cache/state locations with explicit
overrides. This preserves the universal skills install contract while avoiding a
Windows-shaped `%LOCALAPPDATA%\skills-updates` model on Linux and macOS.

`--agents-home`, `--cache-dir`, and `--state-dir` should exist for automation,
with matching environment variables.
