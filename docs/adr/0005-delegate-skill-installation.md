# Delegate skill installation

The Go skills updater delegates add and remove operations to the upstream Skills
CLI instead of reimplementing package installation. The updater owns discovery,
comparison, skip handling, dry-run planning, lockfile safety, and output
protocols; the Skills CLI remains responsible for install semantics, isolated
behind a runner interface so the boundary is testable and replaceable.

The default runner is `pnpm dlx skills@latest`, but users and tests can override
it with `--skills-command` or `SK_UP_SKILLS_COMMAND`. When no override is set,
the updater should fall back across available runners in this order:

1. `pnpm dlx skills@latest`
2. `bunx skills@latest`
3. `deno run -A npm:skills@latest`
4. `npx -y skills@latest`

Runner execution should use tokenized process execution rather than shell
evaluation by default.

Remove operations keep defensive updater-owned cleanup after a successful
delegated remove: delete the installed skill directory, clear saved skip state,
and remove the lockfile entry. Cleanup must be transactional where lockfile
state is involved and visible in dry-run planned actions.
