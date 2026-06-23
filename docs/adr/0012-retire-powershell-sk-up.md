# Retire the PowerShell skills updater

Status: implemented.

The PowerShell implementation is retired from the active path. The Windows
wrappers invoke the Go binary, and the PowerShell script is not shipped as a
supported fallback; keeping one supported implementation avoids drift in command
behavior, state handling, and structured output.

The retired PowerShell implementation is deleted from the active repo instead
of moved to an in-tree legacy location. Git history is the historical reference.
