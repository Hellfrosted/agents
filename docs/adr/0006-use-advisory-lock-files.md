# Use advisory lock files

The portable skills updater uses advisory lock files next to `.skill-lock.json`
instead of OS-specific named mutexes. Lock files give Linux, macOS, and Windows
the same coordination model, can include PID, host, and timestamp metadata for
diagnostics, and still allow the updater to keep backup and repair behavior for
interrupted lockfile writes.

The updater treats `.skill-lock.json` as externally owned shared state. It must
preserve unknown fields and unrelated skill entries across install/remove
transactions, making only targeted changes and restoring snapshots when a
transaction cannot complete safely.
