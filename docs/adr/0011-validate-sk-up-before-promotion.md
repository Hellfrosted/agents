# Validate sk-up before promotion

Status: implemented as the regression floor.

The Go skills updater passed a promotion validation bar before it replaced the
PowerShell implementation. That bar remains the regression floor for current
maintenance: unit tests for CLI parsing, configuration and environment
resolution, runner fallback, path resolution, lockfile preservation, skip
semantics, JSON and JSONL output, and dry-run planning; integration tests with
fake Git and fake Skills CLI runners; golden help tests for both `sk-up` and
`skills-updates`; cross-platform build coverage; arm64 build smoke checks; and
manual workstation smoke tests for daily commands when the Windows surface is
being repaired or republished.
