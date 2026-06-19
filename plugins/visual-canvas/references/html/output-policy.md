# HTML Output Policy

This is the reusable replacement for global user-facing HTML guidance.

## Baseline

- Default to dark mode unless the user asks for light mode or a project style
  clearly overrides it.
- Use readable contrast. Body text should meet 4.5:1 contrast; large text
  should meet 3:1.
- Respect `prefers-reduced-motion` for every animation.
- Use real visual assets when a site, game, report, or visual explanation needs
  them. Assets may be generated images, SVG, charts, screenshots, or other
  meaningful media.
- Do not make a landing page when the user asked for a tool, report, game, or
  app. Make the actual usable artifact first.

## Layout

- Use full-width sections or unframed layouts. Do not put cards inside cards.
- Cards are for repeated items, compact panels, and genuinely framed tools; do
  not use them as default page-section wrappers.
- Fixed-format elements need stable dimensions: use `aspect-ratio`, explicit
  grid tracks, min/max constraints, or container-relative sizing.
- Text must not overflow or overlap its container at mobile or desktop sizes.
- Do not scale font size directly with viewport width.
- Keep letter spacing at `0` unless typography truly needs a small adjustment.

## Visual Tone

- Avoid one-note palettes dominated by one hue family.
- Avoid default purple/indigo gradients, beige/cream monoculture, dark
  blue/slate monoculture, and brown/orange/espresso monoculture.
- Avoid gradient orbs, bokeh blobs, decorative diagonal stripes, and generic
  stock-like atmospheric imagery.
- Use icons for familiar controls where appropriate. Prefer an installed icon
  library when the project already has one.

## HTML Artifact Rules

- Embed a compact metadata JSON block in final HTML when the artifact has a
  project directory.
- Use semantic HTML for tables, headings, navigation, and sections.
- If using Mermaid, include zoom/pan or a readable presentation strategy for
  complex diagrams.
- If using generated images, keep prompts and saved asset paths in `canvas.json`
  or the relevant asset metadata.

The static checker only enforces a small set of high-confidence failures. It is
not a substitute for browser inspection, contrast review, content review, or
specialist design/taste skills.
