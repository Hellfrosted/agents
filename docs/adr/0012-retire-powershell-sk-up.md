# Retire the PowerShell skills updater

After the Go skills updater is promoted, the current PowerShell implementation
is retired from the active path. The Windows wrappers invoke the Go binary, and
the PowerShell script is not shipped as a supported fallback; keeping one
supported implementation avoids drift in command behavior, state handling, and
structured output.

Promotion should delete the retired PowerShell implementation from the active
repo instead of moving it to an in-tree legacy location. Git history is the
historical reference.
