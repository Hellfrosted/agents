# Port sk-up to Go

Status: implemented.

The portable `sk-up` replacement is implemented in Go. The original
PowerShell-first implementation matched the first Windows workstation use case,
but the current contract is a composable CLI that works consistently across
major Linux distributions, macOS, and Windows. Go provides single-file native
binaries, practical process and filesystem APIs, straightforward JSON output,
and simpler cross-compilation than keeping PowerShell as the core runtime.
