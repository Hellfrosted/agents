---
name: visual-canvas
description: Visual Canvas for portable local HTML. Use when the user explicitly wants a portable local/offline HTML report, Visual Canvas profile, reusable HTML output policy, or HTML quality check. Route general visual plans and recaps to visual-plan or visual-recap.
---

# Visual Canvas

Create compact, local HTML artifacts and HTML policy outputs for developers.
Visual Canvas is a meta-skill: it owns routing, artifact structure, profile
resolution, asset orchestration, validation evidence, and final delivery.

## Route First

First read `../../references/canvas/modes.md` and choose the narrowest mode
that satisfies the request. If the artifact purpose is genuinely ambiguous,
ask one clarifying question; otherwise continue.

Do not use Visual Canvas for general Agent-Native plan or recap lanes. Route
implementation plans, architecture plans, migration plans, branch recaps, pull
request recaps, diff recaps, and agent-work recaps to the installed
`/visual-plan` or `/visual-recap` skill unless the user specifically asks for a
portable local HTML report.

The user should only need to invoke `visual-canvas`. Select any mode-specific
helper behavior yourself.

## Progressive References

After choosing a mode, load only the references needed by that branch:

- Report or local HTML Review: read
  `../../references/canvas/report-pipeline.md`,
  `../../references/canvas/artifact-contract.md`,
  `../../references/canvas/profile-resolution.md`,
  `../../references/canvas/asset-pipeline.md`,
  `../../references/html/output-policy.md`,
  `../../references/html/design-delegation.md`, and
  `../../references/html/visual-qa.md`.
- Local HTML Plan: read the same artifact-producing references as Report, but
  only after confirming `/visual-plan` is not the right route.
- Style Profile: read `../../references/canvas/profile-resolution.md`. Read
  artifact references only if the user also wants a rendered profile report.
- HTML Output Policy: read `../../references/html/output-policy.md`,
  `../../references/html/design-delegation.md`, and
  `../../references/html/visual-qa.md`. Use the checker directly against an
  existing HTML file when one is provided.

## Default Artifact Shape

Write compact artifacts by default:

```text
<project-dir>/
  canvas.json
  <descriptive-report-title>.html
  visuals/                 # only when external assets are needed
```

Use expanded sidecars only when a profile or user request explicitly asks for
debug/resume detail.

## Hard Boundaries

- Do not build a custom viewer or workbench app. Generate HTML and open/link it
  through the current harness or the user's browser.
- Do not duplicate full report prose in JSON. `canvas.json` holds metadata,
  structure, asset records, and validation summary; the HTML holds the prose.
- Do not reflexively generate images. If imagegen assets are needed, decide
  early and dispatch them in parallel.

## Completion

Finish with the artifact or policy result, plus validation evidence. For
artifact-producing modes, include the HTML path, `canvas.json` path when
created, checker/browser results, and any unavailable-check note. For HTML
policy or profile-only modes, include the policy/profile result and the checks
that ran, or state why a check was unavailable.
