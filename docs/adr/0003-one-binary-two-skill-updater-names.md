# One binary with two skill updater names

Status: implemented.

The portable skills updater ships as one Go binary that supports both `sk-up`
and `skills-updates` entry point names. Keeping one core executable avoids
divergent behavior while allowing `sk-up` to remain the terse daily-use
interface and `skills-updates` to remain the longer readable compatibility
interface. Unix installs can use a link or copied binary, and Windows wrappers
invoke the same core executable.
