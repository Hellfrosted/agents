# Maintainer Guide

This guide captures the source-first workflow and local checks for maintaining
this repository.

## Source-First Workflow

Keep edits in this repository first. Copy files into active workstation install
locations only when explicitly doing install, repair, or republish work.

For behavior, workflow, script, or config changes:

1. Change the source file in this repository.
2. Run the smallest relevant check.
3. Update the focused docs page that owns the behavior.
4. Repair or republish installed files only when that is part of the task.

## Common Checks

### Documentation-only changes

Run a local Markdown link check over the changed docs and check for whitespace
errors before handing off documentation-only work:

```bash
node - README.md RESOURCES.md docs/README.md docs/maintainer-guide.md <<'NODE'
const fs = require('fs');
const path = require('path');
const files = process.argv.slice(2);
let failed = false;
for (const file of files) {
  const text = fs.readFileSync(file, 'utf8');
  const re = /\[[^\]]+\]\(([^)]+)\)/g;
  let match;
  while ((match = re.exec(text))) {
    const href = match[1];
    if (/^[a-z]+:/i.test(href) || href.startsWith('#')) continue;
    const target = href.split('#')[0];
    const resolved = path.resolve(path.dirname(file), target);
    if (!fs.existsSync(resolved)) {
      console.error(`${file}: missing link target ${href}`);
      failed = true;
    }
  }
}
process.exit(failed ? 1 : 0);
NODE

git diff --check
```

### Workstation hooks

Run the focused source tests before publishing changes to the active global hook:

```bash
PYTHONDONTWRITEBYTECODE=1 python3 -m unittest \
  hooks/test_wsl_command_guardrails.py
```

The hook contract is documented in [hook-guardrails.md](hook-guardrails.md).
The active global copy lives under `$CODEX_HOME/hooks/`. When repairing the
installed guardrail, copy `hooks/wsl_command_guardrails.py` to
`$CODEX_HOME/hooks/wsl_command_guardrails.py` and keep `$CODEX_HOME/hooks.json`
matching Codex shell tool names used by the current runtime. The source-side
reference snippet is `hooks/global-pretooluse-hooks.example.json`.

## Safety Rules

- Preserve user changes and unrelated local state.
- Use Windows-native tooling to repair symlinks on Windows drives.
- Keep temporary agent notes out of VCS with `.git/info/exclude`, not
  `.gitignore`.
- Never record secrets, tokens, recovery codes, private personal data, or raw
  session exports in docs, examples, commits, or skill text.
- For behavior, workflow, script, or config changes, update the relevant human
  docs in the same turn.
