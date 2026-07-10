#!/usr/bin/env python3

"""Behavioral tests for the Visual Canvas checker and scaffolder."""

from __future__ import annotations

import importlib.util
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
SCRIPTS = ROOT / "scripts"


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"Unable to load {path}")
    module = importlib.util.module_from_spec(spec)
    sys.modules[name] = module
    spec.loader.exec_module(module)
    return module


checker = load_module("visual_canvas_checker", SCRIPTS / "check_html_policy.py")
scaffolder = load_module("visual_canvas_scaffolder", SCRIPTS / "scaffold_canvas.py")


class HtmlPolicyTests(unittest.TestCase):
    def check_text(self, text: str) -> set[str]:
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "report.html"
            path.write_text(text, encoding="utf-8")
            return {finding.code for finding in checker.check_html(path)}

    def test_starter_template_passes(self) -> None:
        template = ROOT / "assets" / "templates" / "starter-report.html"
        self.assertEqual(checker.check_html(template), [])

    def test_detects_missing_viewport_and_motion_guard(self) -> None:
        codes = self.check_text("<style>.x { transition: color 1s; }</style>")
        self.assertEqual(codes, {"missing-viewport", "missing-reduced-motion"})

    def test_detects_unreadable_mermaid_and_policy_patterns(self) -> None:
        codes = self.check_text(
            '<meta name="viewport" content="width=device-width">'
            '<div class="mermaid orb" style="font-size: 2vw"></div>'
        )
        self.assertEqual(
            codes,
            {"decorative-orb", "mermaid-no-readable-controls", "viewport-scaled-type"},
        )


class ScaffoldTests(unittest.TestCase):
    def test_slugify_has_a_stable_fallback(self) -> None:
        self.assertEqual(scaffolder.slugify("A Useful Report"), "a-useful-report")
        self.assertEqual(scaffolder.slugify("!!!"), "visual-canvas")

    def test_cli_writes_matching_canvas_and_escaped_html(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            out_dir = Path(directory) / "canvas"
            subprocess.run(
                [
                    sys.executable,
                    str(SCRIPTS / "scaffold_canvas.py"),
                    "Report <One>",
                    "--mode",
                    "visual-review",
                    "--out-dir",
                    str(out_dir),
                    "--with-html",
                ],
                check=True,
                capture_output=True,
                text=True,
            )

            canvas = json.loads((out_dir / "canvas.json").read_text(encoding="utf-8"))
            html = (out_dir / canvas["paths"]["html"]).read_text(encoding="utf-8")
            self.assertEqual(canvas["mode"], "visual-review")
            self.assertEqual(canvas["title"], "Report <One>")
            self.assertIn("Report &lt;One&gt;", html)
            metadata = json.dumps(
                {"schemaVersion": 1, "canvasId": canvas["id"]},
                separators=(",", ":"),
            )
            self.assertIn(metadata, html)


if __name__ == "__main__":
    unittest.main()
