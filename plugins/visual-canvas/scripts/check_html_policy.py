#!/usr/bin/env python3
# /// script
# requires-python = ">=3.11"
# ///
# ─── How to run ───
# python3 plugins/visual-canvas/scripts/check_html_policy.py path/to/file.html

"""Static checks for Visual Canvas user-facing HTML artifacts."""

from __future__ import annotations

import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Final


@dataclass(frozen=True, slots=True)
class Finding:
    code: str
    message: str


@dataclass(frozen=True, slots=True)
class Check:
    code: str
    pattern: re.Pattern[str]
    message: str


CHECKS: Final[tuple[Check, ...]] = (
    Check(
        "gradient-text",
        re.compile(r"background-clip\s*:\s*text|-webkit-background-clip\s*:\s*text", re.I),
        "Avoid gradient text in generated HTML artifacts.",
    ),
    Check(
        "decorative-orb",
        re.compile(r"(?:class|id)\s*=\s*[\"'][^\"']*\b(?:orb|blob|bokeh)\b|[.#][\w-]*(?:orb|blob|bokeh)[\w-]*", re.I),
        "Avoid decorative orb/blob/bokeh naming and patterns.",
    ),
    Check(
        "tailwind-purple",
        re.compile(r"#(?:8b5cf6|7c3aed|a78bfa|d946ef)\b", re.I),
        "Avoid default violet/fuchsia accent colors.",
    ),
    Check(
        "oversized-radius",
        re.compile(r"border-radius\s*:\s*(?:3[2-9]|[4-9]\d|\d{3,})px", re.I),
        "Large card radii usually read as AI-generated; justify or reduce.",
    ),
    Check(
        "huge-shadow",
        re.compile(r"box-shadow\s*:[^;]*(?:1[6-9]|[2-9]\d)px", re.I),
        "Wide soft shadows are often ghost-card decoration.",
    ),
    Check(
        "stripe-background",
        re.compile(r"repeating-linear-gradient\s*\(", re.I),
        "Avoid decorative repeating stripe backgrounds.",
    ),
    Check(
        "viewport-scaled-type",
        re.compile(r"font-size\s*:[^;]*(?:vw|vh|vmin|vmax)", re.I),
        "Do not scale font size directly with viewport units.",
    ),
)


def check_html(path: Path) -> list[Finding]:
    text = path.read_text(encoding="utf-8")
    findings: list[Finding] = []

    if "<meta name=\"viewport\"" not in text and "<meta name='viewport'" not in text:
        findings.append(Finding("missing-viewport", "Missing responsive viewport meta tag."))

    has_animation = re.search(r"@keyframes|transition\s*:|animation\s*:", text, re.I) is not None
    has_reduced_motion = "prefers-reduced-motion" in text
    if has_animation and not has_reduced_motion:
        findings.append(Finding("missing-reduced-motion", "Animations need prefers-reduced-motion handling."))

    has_mermaid = "mermaid" in text.lower()
    has_zoom = re.search(r"\b(zoom|pan|diagram-shell|mermaid-viewport)\b", text, re.I) is not None
    if has_mermaid and not has_zoom:
        findings.append(Finding("mermaid-no-readable-controls", "Mermaid diagrams need a readable zoom/pan strategy."))

    for check in CHECKS:
        if check.pattern.search(text):
            findings.append(Finding(check.code, check.message))

    return findings


def main(argv: list[str]) -> int:
    if len(argv) != 2:
        print("Usage: check_html_policy.py <file.html>", file=sys.stderr)
        return 2

    path = Path(argv[1])
    if not path.is_file():
        print(f"Not a file: {path}", file=sys.stderr)
        return 2

    findings = check_html(path)
    if not findings:
        print("OK: no static HTML policy findings.")
        return 0

    for finding in findings:
        print(f"{finding.code}: {finding.message}")
    return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
