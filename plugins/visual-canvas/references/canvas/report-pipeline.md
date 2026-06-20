# Visual Canvas Report Pipeline

Visual Canvas runs as a compact artifact pipeline, not as an ad hoc HTML dump.

1. Classify intent and choose the public mode.
   Route general visual plans and recaps to the installed Agent-Native skills
   before creating a Visual Canvas artifact.
2. Resolve profiles in project-first order.
3. For artifact-producing modes, create or update `canvas.json` with title, mode, output path, profile
   summary, section outline, assets, and validation summary.
4. Decide slow assets early. If generated images are needed, start them in a
   separate worker immediately and continue deterministic work.
5. Build deterministic visuals: Mermaid, SVG, and charts.
6. Compose a descriptive HTML file named from the finalized report title.
7. Run the HTML policy checker and browser/visual QA where available.
8. Open or link the final HTML.

Default artifact detail is `compact`. Expanded sidecars are opt-in for debug or
resume-heavy workflows. The default artifact root is
`~/.agent/visual-canvas/projects/`; use repo-local output only on explicit
request.
