# Visual QA

Run the strongest available checks for the current harness.

Minimum:

- Run `$VISUAL_CANVAS_ROOT/scripts/check_html_policy.py` on the final HTML.
  `VISUAL_CANVAS_ROOT` is the directory containing `.codex-plugin/plugin.json`
  for this installed plugin. Treat the checker as a lint smoke check, not as
  full visual QA.
- Open the HTML in the available browser or provide a durable local path.
- Inspect at least desktop and mobile widths when the harness supports it.

Check for:

- blank or broken render
- missing assets
- unreadable contrast
- text overflow or overlap
- mobile horizontal scrolling caused by layout mistakes
- unstyled default links or controls
- unsupported motion without reduced-motion handling
- Mermaid diagrams that render too small to read

When browser automation is unavailable, say so and report the static checker
result plus the file path.
