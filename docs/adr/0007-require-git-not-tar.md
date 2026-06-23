# Require Git but not tar

Status: implemented.

The Go skills updater requires an external `git` executable for repository
fetching, sparse checkout, tree export, and diff-compatible behavior, but it
does not require an external `tar` executable. Git remains a core workflow
dependency and avoids reimplementing repository semantics, while Go can safely
handle archive extraction itself to reduce Windows portability friction.

Commands that do not need upstream comparison still work without Git,
including list, remove, unskip, skips, and named installs that can delegate from
lockfile source metadata. Commands that need upstream comparison, including
status, diff, open-diff, install changed/missing, and skip-current-update,
fail with a clear dependency error when Git is unavailable.
