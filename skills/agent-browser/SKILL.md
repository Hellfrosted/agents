---
name: agent-browser
description: Uses the installed Agent Browser CLI for browser automation, visual QA, screenshots, page inspection, and Electron/Slack/browser-adjacent workflows. Use when the user asks to use agent-browser, automate a browser page, test a web app through a browser, take screenshots, inspect rendered UI, or load Agent Browser's version-matched workflow skills.
---

# Agent Browser

Agent Browser is installed as a global CLI on this workstation. This repo-owned
skill is intentionally a small launcher: load the version-matched CLI skill
content before running browser automation so commands stay aligned with the
installed Agent Browser release.

## Start Here

Check the install and load the current core workflow:

```bash
agent-browser --version
agent-browser skills get core
agent-browser skills get core --full
```

Use specialized CLI-provided skill text when the task is not a normal browser
page:

```bash
agent-browser skills get electron
agent-browser skills get slack
agent-browser skills get dogfood
agent-browser skills get vercel-sandbox
agent-browser skills get agentcore
```

Run `agent-browser skills list` to see every workflow available in the installed
version.

## Common Checks

Use these for quick install and browser-surface smoke tests:

```bash
agent-browser install
agent-browser doctor
agent-browser open https://example.com
agent-browser get title
agent-browser get url
agent-browser snapshot -i -c
agent-browser screenshot example.png
agent-browser close --all
```

`agent-browser screenshot <path>` may print a managed screenshot path instead of
writing exactly to the requested path. When preserving evidence, record or copy
the printed path.

## Safety

- Do not paste secrets into browser sessions.
- Prefer a fresh session for unauthenticated QA.
- Close sessions with `agent-browser close --all` after test runs unless the
  user asks to keep them open.
- For Codex in-app browser work, prefer `@Browser` when it is available in the
  current thread. Use Agent Browser as the CLI fallback when the Browser plugin
  or required Node REPL tool surface is not exposed.
