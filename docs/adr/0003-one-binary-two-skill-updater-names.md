# One binary with two skill updater names

The portable skills updater will ship as one Go binary that supports both
`sk-up` and `skills-updates` entry point names. Keeping one core executable
avoids divergent behavior while allowing `sk-up` to remain the terse daily-use
interface and `skills-updates` to remain the longer readable compatibility
interface; Unix installs can use a symlink or copied binary, and Windows can
keep wrappers that invoke the same core.
