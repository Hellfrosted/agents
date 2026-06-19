#!/usr/bin/env python3
# /// script
# requires-python = ">=3.11"
# ///

"""Repo-local structural checks for the Visual Canvas plugin source tree."""

from __future__ import annotations

import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent

REQUIRED_FILES = (
    ".codex-plugin/plugin.json",
    "skills/visual-canvas/SKILL.md",
    "references/canvas/artifact-contract.md",
    "references/canvas/asset-pipeline.md",
    "references/canvas/default-profile.md",
    "references/canvas/modes.md",
    "references/canvas/profile-resolution.md",
    "references/canvas/report-pipeline.md",
    "references/html/design-delegation.md",
    "references/html/output-policy.md",
    "references/html/visual-qa.md",
    "scripts/check_html_policy.py",
    "scripts/scaffold_canvas.py",
    "assets/templates/starter-report.html",
)


def fail(message: str) -> int:
    print(f"ERROR: {message}", file=sys.stderr)
    return 1


def main() -> int:
    missing = [path for path in REQUIRED_FILES if not (ROOT / path).is_file()]
    if missing:
        return fail("missing required files: " + ", ".join(missing))

    try:
        manifest = json.loads((ROOT / ".codex-plugin/plugin.json").read_text(encoding="utf-8"))
    except json.JSONDecodeError as error:
        return fail(f"plugin.json is invalid JSON: {error}")

    if manifest.get("name") != "visual-canvas":
        return fail("plugin.json name must be visual-canvas")
    if manifest.get("skills") != "./skills/":
        return fail("plugin.json skills must be ./skills/")

    skill_files = sorted(path.relative_to(ROOT).as_posix() for path in (ROOT / "skills").glob("*/SKILL.md"))
    if skill_files != ["skills/visual-canvas/SKILL.md"]:
        return fail("expected exactly one public skill: skills/visual-canvas/SKILL.md")

    generated = [
        path.relative_to(ROOT).as_posix()
        for path in ROOT.rglob("*")
        if path.name == "__pycache__" or path.suffix == ".pyc"
    ]
    if generated:
        return fail("generated Python cache files found: " + ", ".join(generated))

    print("OK: Visual Canvas plugin structure is valid.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
