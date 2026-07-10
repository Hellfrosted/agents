#!/usr/bin/env python3
"""Regression tests for the Codex adapter around ICM's PreToolUse hook."""

from __future__ import annotations

import json
import subprocess
import unittest
from pathlib import Path


HOOK = Path(__file__).parent / "icm_pre_json.py"


def payload(command: str) -> dict:
    return {
        "session_id": "test-session",
        "cwd": str(Path.cwd()),
        "hook_event_name": "PreToolUse",
        "tool_name": "Bash",
        "tool_use_id": "test-tool-use",
        "tool_input": {"command": command},
    }


def run_hook(command: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [str(HOOK)],
        input=json.dumps(payload(command)),
        text=True,
        capture_output=True,
        check=False,
    )


class IcmPreJsonTest(unittest.TestCase):
    def test_adds_required_updated_input_to_icm_allow(self) -> None:
        result = run_hook("icm status")
        self.assertEqual(result.returncode, 0, result.stderr)
        output = json.loads(result.stdout)
        hook_output = output["hookSpecificOutput"]
        self.assertEqual(hook_output["hookEventName"], "PreToolUse")
        self.assertEqual(hook_output["permissionDecision"], "allow")
        self.assertEqual(hook_output["updatedInput"], {"command": "icm status"})

    def test_non_icm_command_emits_nothing(self) -> None:
        result = run_hook("git status")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout, "")

    def test_malformed_payload_fails_open_without_invalid_output(self) -> None:
        result = subprocess.run(
            [str(HOOK)],
            input="not json",
            text=True,
            capture_output=True,
            check=False,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout, "")


if __name__ == "__main__":
    unittest.main()
