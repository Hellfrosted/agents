# Do not add a sk-up config file

The promoted Go skills updater will not support a config file. Configuration
comes from explicit CLI flags, environment variables, and platform defaults
only, which keeps automation behavior visible at the call site and avoids
long-term precedence, discovery, and migration rules for a separate config
surface.

Environment variables are limited to stable environment and integration
defaults: `SK_UP_AGENTS_HOME`, `SK_UP_CACHE_DIR`, `SK_UP_STATE_DIR`,
`SK_UP_SKILLS_COMMAND`, `SK_UP_DIFF_TOOL`, and `SK_UP_COLOR`. Targets,
structured-output modes, dry-run, install/remove mode, and skip actions stay
explicit in the command invocation.
