# Profile Resolution

Profile lookup order:

```text
<repo>/.agent/visual-canvas.local.md
<repo>/.agent/visual-canvas.md
~/.agent/visual-canvas/profiles/default.md
<plugin>/references/canvas/default-profile.md
```

Profiles are human-authored Markdown with optional YAML frontmatter. They may
describe visual style, report structure, diagram rules, asset policy, forbidden
patterns, and preferred specialist skills.

The run records profile sources and a short effective summary in `canvas.json`.
Expanded mode may also write profile sidecars, but compact mode does not.

Local project profiles are preferred by default. Commit a project profile only
when the style/report guidance is useful for other humans on the repo.
