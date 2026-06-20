# Validate sk-up before promotion

The Go skills updater must pass a promotion validation bar before it replaces
the current PowerShell implementation. Required coverage includes unit tests for
CLI parsing, configuration and environment resolution, runner fallback, path
resolution, lockfile preservation, skip semantics, JSON and JSONL output, and
dry-run planning; integration tests with fake Git and fake Skills CLI runners;
golden help tests for both `sk-up` and `skills-updates`; cross-platform CI for
Linux, macOS, and Windows on amd64; arm64 build smoke checks; and manual
workstation smoke tests for the daily commands.
