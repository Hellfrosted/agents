# Break sk-up behavior before promotion

Status: implemented before Go promotion.

The Go `sk-up` port intentionally resolved command behavior before the tool was
promoted. That was the last low-cost window to replace PowerShell-era workflow
compromises with a coherent composable CLI contract. Compatibility is now a
user-facing commitment. There is no throwaway v1: the promoted release is the
compatibility contract users build on.

The short `sk-up` flags remain part of the intended interface because `sk-up`
itself is already the shorthand command. The longer `skills-updates` entry point
can keep readable long-form compatibility for users and scripts that prefer it.

The existing `-z` muscle memory remains, but it is an alias for opening diffs
with a configured diff tool rather than a Zed-only contract. Zed is the default
when available, and users can set `--diff-tool` or `SK_UP_DIFF_TOOL` for other
editors.

Source installs require explicit syntax instead of being inferred from
`sk-up -i <arg>`. Named installs keep `sk-up -i <skill...>`, while package or
repository installs use `sk-up -I <source...>` or `skills-updates
--install-source <source...>` so automation does not depend on URL/name
heuristics.

No-target install remains intentional shorthand: `sk-up -i` installs all
changed or missing unskipped skills, with `skills-updates --install` and
`skills-updates --install-all` as long-form equivalents.
