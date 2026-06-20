# Ship sk-up release archives

The promoted Go skills updater ships canonical release archives with native
binaries for Linux, macOS, and Windows across amd64 and arm64. Each archive
contains the core executable, a `skills-updates` alias or wrapper, README/help
text, and license; package-manager formulas can follow later, but release
archives are the primary install unit.

Expected artifacts are:

- `sk-up-linux-amd64.tar.gz`
- `sk-up-linux-arm64.tar.gz`
- `sk-up-darwin-amd64.tar.gz`
- `sk-up-darwin-arm64.tar.gz`
- `sk-up-windows-amd64.zip`
- `sk-up-windows-arm64.zip`

Windows archives keep tiny `sk-up.cmd` and `skills-updates.cmd` wrappers around
`sk-up.exe`. They intentionally do not ship `skills-updates.exe`: Windows can
treat updater-looking executable names as installer/update programs and trigger
UAC before `--help` can run. The wrappers preserve current command names, can
set UTF-8 console behavior, and remain non-authoritative; all updater behavior
lives in the Go binary.
