#!/usr/bin/env python3
"""Adapt ICM's PreToolUse allow response to the current Codex hook contract."""

from __future__ import annotations

import json
import os
import subprocess
import sys


DEFAULT_ICM_BIN = "/home/crunch/.local/bin/icm"


def main() -> int:
    raw_payload = sys.stdin.read()
    try:
        payload = json.loads(raw_payload)
    except json.JSONDecodeError:
        return 0

    if payload.get("hook_event_name") != "PreToolUse":
        return 0

    try:
        result = subprocess.run(
            [os.environ.get("ICM_BIN", DEFAULT_ICM_BIN), "hook", "pre"],
            input=raw_payload,
            text=True,
            capture_output=True,
            check=False,
            timeout=4,
        )
    except (OSError, subprocess.TimeoutExpired):
        return 0

    if result.returncode != 0 or not result.stdout.strip():
        return 0

    try:
        response = json.loads(result.stdout)
    except json.JSONDecodeError:
        return 0

    hook_output = response.get("hookSpecificOutput")
    if not isinstance(hook_output, dict):
        return 0

    if hook_output.get("permissionDecision") == "allow" and "updatedInput" not in hook_output:
        tool_input = payload.get("tool_input")
        if not isinstance(tool_input, dict) or not isinstance(tool_input.get("command"), str):
            return 0
        hook_output["updatedInput"] = tool_input

    print(json.dumps(response, separators=(",", ":")))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
