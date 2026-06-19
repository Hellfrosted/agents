#!/usr/bin/env python3
"""Regression tests for the Codex WSL command guardrail hook."""

from __future__ import annotations

import json
import shlex
import subprocess
import unittest
from pathlib import Path


HOOK = Path(__file__).with_name("wsl_command_guardrails.py")


def run_hook(command: str) -> subprocess.CompletedProcess[str]:
    payload = json.dumps(
        {
            "hook_event_name": "PreToolUse",
            "tool_name": "functions.shell_command",
            "tool_input": {"command": command},
        }
    )
    return subprocess.run([str(HOOK)], input=payload, text=True, capture_output=True, check=False)


def run_parallel_hook(command: str) -> subprocess.CompletedProcess[str]:
    payload = json.dumps(
        {
            "hook_event_name": "PreToolUse",
            "tool_name": "multi_tool_use.parallel",
            "tool_input": {
                "tool_uses": [
                    {
                        "recipient_name": "functions.shell_command",
                        "parameters": {"command": command},
                    }
                ]
            },
        }
    )
    return subprocess.run([str(HOOK)], input=payload, text=True, capture_output=True, check=False)


class WslCommandGuardrailsTest(unittest.TestCase):
    def assert_allowed(self, command: str) -> None:
        result = run_hook(command)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout, "")

    def assert_blocked(self, command: str, expected: str) -> None:
        result = run_hook(command)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("permissionDecision", result.stdout)
        self.assertIn("deny", result.stdout)
        self.assertIn(expected, result.stdout)

    def test_allows_wsl_native_git_and_bridge_helpers(self) -> None:
        self.assert_allowed("git status --short")
        self.assert_allowed(r"wslpath 'C:\\Users\\nguco\\file.txt'")
        self.assert_allowed("explorer.exe .")

    def test_blocks_windows_shells(self) -> None:
        self.assert_blocked("cmd.exe /d /c git status", "open Tabby and leave it running")
        self.assert_blocked("cmd /d /c dir", "Tabby MCP")
        self.assert_blocked("powershell.exe -NoProfile -Command Get-ChildItem", "Tabby MCP")
        self.assert_blocked("pwsh -NoProfile -File tools/repair.ps1", "Tabby MCP")

    def test_blocks_windows_terminal_launches(self) -> None:
        self.assert_blocked(
            "wt.exe new-tab powershell.exe -NoExit -Command git status",
            "open Tabby and leave it running",
        )

    def test_blocks_windows_git_executable(self) -> None:
        self.assert_blocked("'/mnt/c/Program Files/Git/cmd/git.exe' status", "Windows git.exe")
        self.assert_blocked(r"C:\\Git\\cmd\\git.exe status", "Windows git.exe")

    def test_blocks_ps1_as_executable(self) -> None:
        self.assert_blocked("tools/repair.ps1", "PS1")

    def test_blocks_nested_windows_shells(self) -> None:
        self.assert_blocked("bash -lc 'cmd.exe /c git status'", "Tabby MCP")
        self.assert_blocked("rtk run 'powershell.exe -NoProfile -Command git status'", "Tabby MCP")

    def test_blocks_default_rtk_wrapped_windows_commands(self) -> None:
        self.assert_blocked("rtk cmd.exe /d /c dir", "Tabby MCP")
        self.assert_blocked("rtk powershell.exe -NoProfile -Command Get-ChildItem", "Tabby MCP")
        self.assert_blocked("rtk wt.exe new-tab powershell.exe", "Windows Terminal")
        self.assert_blocked("rtk '/mnt/c/Program Files/Git/cmd/git.exe' status", "Windows git.exe")

    def test_blocks_parallel_payload_shell_commands(self) -> None:
        result = run_parallel_hook("cmd.exe /d /c dir")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("permissionDecision", result.stdout)
        self.assertIn("open Tabby and leave it running", result.stdout)

    def test_blocks_windows_shells_behind_launch_wrappers(self) -> None:
        self.assert_blocked("/init /mnt/c/Windows/System32/cmd.exe /d /c dir", "Tabby MCP")
        self.assert_blocked(
            "/init /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -NoProfile -Command Get-ChildItem",
            "Tabby MCP",
        )
        self.assert_blocked("nohup cmd.exe /d /c git status", "Tabby MCP")
        self.assert_blocked("setsid powershell.exe -NoProfile -Command git status", "Tabby MCP")
        self.assert_blocked("timeout 5 pwsh -NoProfile -Command git status", "Tabby MCP")
        self.assert_blocked("nice -n 5 wt.exe new-tab powershell.exe", "Windows Terminal")
        self.assert_blocked("stdbuf -oL cmd.exe /d /c dir", "Tabby MCP")
        self.assert_blocked("winpty cmd.exe /d /c dir", "Tabby MCP")

    def test_blocks_shell_compound_forms(self) -> None:
        self.assert_blocked("{ cmd.exe /d /c dir; }", "Tabby MCP")
        self.assert_blocked("if true; then cmd.exe /d /c dir; fi", "Tabby MCP")
        self.assert_blocked("while false; do powershell.exe -NoProfile -Command Get-ChildItem; done", "Tabby MCP")
        self.assert_blocked("for x in 1; do wt.exe new-tab powershell.exe; done", "Windows Terminal")
        self.assert_blocked("bash -lc '{ cmd.exe /d /c dir; }'", "Tabby MCP")

    def test_deep_nesting_denial_keeps_tabby_reminder(self) -> None:
        command = "cmd.exe /d /c dir"
        for _ in range(7):
            command = f"bash -lc {shlex.quote(command)}"
        self.assert_blocked(command, "open Tabby and leave it running")

    def test_blocks_raw_windows_paths_in_wsl_commands(self) -> None:
        self.assert_blocked(r"python C:\\Users\\nguco\\script.py", "Windows paths")

    def test_blocks_temp_on_windows_mounts(self) -> None:
        self.assert_blocked("TMPDIR=/mnt/c/tmp pnpm test", "TMPDIR")

    def test_blocks_linux_symlink_on_windows_mount(self) -> None:
        self.assert_blocked("ln -s /tmp/source /mnt/c/Users/nguco/bin/tool", "symlinks")


if __name__ == "__main__":
    unittest.main()
