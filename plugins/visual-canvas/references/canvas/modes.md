# Visual Canvas Modes

These modes are internal. The public user-facing skill is always
`visual-canvas`.

## Report

Use for general visual explanations and reports about systems, codebases,
decisions, data, or technical concepts.

Output bias:

- answer first
- diagrams for relationships and flows
- evidence tables where useful
- compact narrative sections

## Review

Use for branches, diffs, PRs, implementation plans, or technical proposals.

Output bias:

- findings and risks first
- change map
- validation gaps
- blocking decisions
- supporting diagrams only where they clarify the review

## Plan

Use for implementation plans, architecture proposals, migrations, and staged
technical work.

Output bias:

- target outcome
- current vs proposed shape
- staged execution
- dependencies and risks
- validation strategy
- open decisions

## Style Profile

Use when the user wants persistent style, report structure, diagram rules,
asset policy, or forbidden patterns.

This mode usually edits or proposes profile files. It does not need
`scaffold_canvas.py` unless the user also asks for a rendered profile report.

Profile locations:

```text
<repo>/.agent/visual-canvas.local.md
<repo>/.agent/visual-canvas.md
~/.agent/visual-canvas/profiles/default.md
<plugin>/references/canvas/default-profile.md
```

Profiles are Markdown with optional YAML frontmatter. Generated run metadata is
JSON.

## HTML Output Policy

Use when the user wants user-facing HTML guidance without a full report
project, or when another workflow needs Visual Canvas quality rules.

This mode can run directly against an existing HTML file. It does not need
`scaffold_canvas.py` unless the user also asks for a rendered policy report.

Read:

- `../html/output-policy.md`
- `../html/design-delegation.md`
- `../html/visual-qa.md`

If an HTML file exists on disk, run the checker from the Visual Canvas plugin
root. The plugin root is the directory containing `.codex-plugin/plugin.json`;
from this `SKILL.md`, it is two directories up.

```bash
python3 "$VISUAL_CANVAS_ROOT/scripts/check_html_policy.py" <file.html>
```
