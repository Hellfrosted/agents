# Artifact Contract

Default generated artifact location is outside the repo:

```text
~/.agent/visual-canvas/projects/<slug>/
```

Use a repo-local output directory only when the user explicitly wants the
artifact to become project documentation or durable repo evidence.

Default project layout:

```text
<project-dir>/
  canvas.json
  <descriptive-title>.html
  visuals/
```

`visuals/` is created only when external assets are needed. Inline CSS, inline
SVG, inline Mermaid source, and inline metadata are acceptable when they keep
the artifact portable and readable.

Only artifact-producing modes create a project directory. `Report` and local
HTML `Review` normally produce artifacts. Agent-Native owns the general
`/visual-plan` and `/visual-recap` workflows; use Visual Canvas planning or
review artifacts only when the user specifically needs a portable local HTML
report. `Style Profile` edits profile guidance, and `HTML Output Policy` can
run as policy/check guidance against an existing HTML file without creating a
Visual Canvas project.

`canvas.json` is the single compact run contract. It should include:

- `schemaVersion`
- `id`
- `mode`
- `title`
- `createdAt`
- `source`
- `paths.html`
- `profileSources`
- `profileSummary`
- `sections`
- `assets`
- `validation`
- `artifactDetail`

Do not duplicate full report prose in JSON. The HTML is the prose artifact;
`canvas.json` is the project index and compact execution record.

Use expanded sidecars only when `artifactDetail` is `expanded`.
