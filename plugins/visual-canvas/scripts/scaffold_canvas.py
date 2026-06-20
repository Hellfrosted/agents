#!/usr/bin/env python3
# /// script
# requires-python = ">=3.11"
# ///
# ─── How to run ───
# python3 plugins/visual-canvas/scripts/scaffold_canvas.py "Report Title" --mode visual-report

"""Create a compact Visual Canvas artifact project directory.

This script is for Visual Canvas local HTML artifact-producing modes:
visual-report, visual-review, and visual-plan. Agent-Native owns general visual
plan and recap artifacts; use visual-plan here only for portable local HTML
planning artifacts. Profile editing and HTML policy checks usually do not need
a project scaffold.
"""

from __future__ import annotations

import html
import json
import re
import sys
from dataclasses import asdict, dataclass
from datetime import UTC, datetime
from pathlib import Path
from typing import Literal

Mode = Literal["visual-report", "visual-review", "visual-plan"]


@dataclass(frozen=True, slots=True)
class Canvas:
    schemaVersion: int
    id: str
    mode: Mode
    title: str
    createdAt: str
    artifactDetail: Literal["compact", "expanded"]
    source: dict[str, str]
    paths: dict[str, str]
    profileSources: list[str]
    profileSummary: str
    sections: list[dict[str, str]]
    assets: list[dict[str, str]]
    validation: dict[str, str]


@dataclass(frozen=True, slots=True)
class Args:
    title: str
    mode: Mode
    out_dir: Path | None
    with_html: bool


def slugify(value: str) -> str:
    lowered = value.lower()
    slug = re.sub(r"[^a-z0-9]+", "-", lowered).strip("-")
    return slug or "visual-canvas"


def build_canvas(title: str, mode: Mode, html_name: str) -> Canvas:
    now = datetime.now(UTC).replace(microsecond=0)
    date_prefix = now.strftime("%Y%m%d")
    return Canvas(
        schemaVersion=1,
        id=f"{date_prefix}-{slugify(title)}",
        mode=mode,
        title=title,
        createdAt=now.isoformat().replace("+00:00", "Z"),
        artifactDetail="compact",
        source={},
        paths={"html": html_name},
        profileSources=[],
        profileSummary="",
        sections=[],
        assets=[],
        validation={"status": "not-run"},
    )


def parse_mode(value: str) -> Mode:
    match value:
        case "visual-report" | "visual-review" | "visual-plan":
            return value
        case unreachable:
            raise SystemExit(f"Invalid --mode: {unreachable}")


def parse_args(argv: list[str]) -> Args:
    title_parts: list[str] = []
    mode: Mode = "visual-report"
    out_dir: Path | None = None
    with_html = False
    index = 1

    while index < len(argv):
        item = argv[index]
        match item:
            case "--mode":
                index += 1
                if index >= len(argv):
                    raise SystemExit("--mode requires a value")
                mode = parse_mode(argv[index])
            case "--out-dir":
                index += 1
                if index >= len(argv):
                    raise SystemExit("--out-dir requires a value")
                out_dir = Path(argv[index])
            case "--with-html":
                with_html = True
            case "--help" | "-h":
                print("Usage: scaffold_canvas.py <title> [--mode MODE] [--out-dir DIR] [--with-html]")
                raise SystemExit(0)
            case _:
                if item.startswith("-"):
                    raise SystemExit(f"Unknown option: {item}")
                title_parts.append(item)
        index += 1

    title = " ".join(title_parts).strip()
    if not title:
        raise SystemExit("Usage: scaffold_canvas.py <title> [--mode MODE] [--out-dir DIR] [--with-html]")
    return Args(title=title, mode=mode, out_dir=out_dir, with_html=with_html)


def starter_html(title: str, canvas_id: str) -> str:
    template_path = Path(__file__).resolve().parent.parent / "assets" / "templates" / "starter-report.html"
    metadata = json.dumps({"schemaVersion": 1, "canvasId": canvas_id}, separators=(",", ":"))
    return (
        template_path.read_text(encoding="utf-8")
        .replace("{{ title }}", html.escape(title, quote=True))
        .replace("{{ metadata_json }}", html.escape(metadata, quote=False))
    )


def main() -> int:
    args = parse_args(sys.argv)
    slug = slugify(args.title)
    out_dir = args.out_dir or Path.home() / ".agent" / "visual-canvas" / "projects" / slug
    html_name = f"{slug}.html"
    canvas = build_canvas(args.title, args.mode, html_name)

    out_dir.mkdir(parents=True, exist_ok=True)
    (out_dir / "canvas.json").write_text(json.dumps(asdict(canvas), indent=2) + "\n", encoding="utf-8")
    if args.with_html:
        (out_dir / html_name).write_text(starter_html(args.title, canvas.id), encoding="utf-8")
    print(out_dir)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
