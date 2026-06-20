# Port sk-up to Go

The portable `sk-up` replacement will be implemented in Go. The current
PowerShell-first implementation matches the original Windows workstation use
case, but the new goal is a composable CLI that works consistently across major
Linux distributions, macOS, and Windows; Go provides single-file native
binaries, practical process and filesystem APIs, straightforward JSON output,
and simpler cross-compilation than keeping PowerShell as the core runtime.
