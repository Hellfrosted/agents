# Use private Go packages for sk-up

Status: implemented.

The Go skills updater source lives under `cmd/sk-up` for the executable and
`internal/skup/...` for implementation packages. This follows common Go layout,
keeps internal boundaries private while the CLI protocol is the public contract,
and supports one binary that adapts behavior for the `sk-up` and
`skills-updates` entry point names.
