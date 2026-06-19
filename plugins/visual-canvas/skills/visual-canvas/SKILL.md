---
name: visual-canvas
description: Single public entrypoint for Visual Canvas. Use when the user asks for a polished visual HTML report, visual explanation, project recap, implementation plan, branch review, diagram-rich artifact, persistent visual profile, or reusable user-facing HTML policy.
---

# Visual Canvas

Create compact, project-backed visual HTML artifacts for developers. Visual
Canvas is a meta-skill: it owns the report pipeline and delegates taste,
image generation, browser checks, and specialist review to existing skills when
they are available.

## Route Internally

The user should only need to invoke `visual-canvas`. Do not ask them to invoke
mode-specific helper skills. Select the mode yourself from
`../../references/canvas/modes.md`.

If the user did not name a mode, pick the narrowest mode that satisfies the
request and continue. Do not ask unless the artifact purpose is genuinely
ambiguous.

## Required Shared References

Before producing an artifact, read:

- `../../references/canvas/modes.md`
- `../../references/canvas/report-pipeline.md`
- `../../references/canvas/artifact-contract.md`
- `../../references/canvas/profile-resolution.md`
- `../../references/canvas/asset-pipeline.md`
- `../../references/html/output-policy.md`
- `../../references/html/design-delegation.md`
- `../../references/html/visual-qa.md`

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
