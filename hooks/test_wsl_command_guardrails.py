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


def run_parallel_cmd_hook(command: str) -> subprocess.CompletedProcess[str]:
    payload = json.dumps(
        {
            "hook_event_name": "PreToolUse",
            "tool_name": "multi_tool_use.parallel",
            "tool_input": {
                "tool_uses": [
                    {
                        "recipient_name": "exec_command",
                        "parameters": {"cmd": command},
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

    def assert_rewritten(self, command: str, expected: str) -> None:
        result = run_hook(command)
        self.assertEqual(result.returncode, 0, result.stderr)
        output = json.loads(result.stdout)
        hook_output = output["hookSpecificOutput"]
        self.assertEqual(hook_output["hookEventName"], "PreToolUse")
        self.assertEqual(hook_output["permissionDecision"], "allow")
        self.assertEqual(hook_output["updatedInput"]["command"], expected)

    def test_allows_wsl_native_git_and_bridge_helpers(self) -> None:
        self.assert_allowed("git status --short")
        self.assert_allowed(r"wslpath 'C:\\Users\\nguco\\file.txt'")
        self.assert_allowed("explorer.exe .")

    def test_allows_search_text_that_mentions_windows_paths(self) -> None:
        self.assert_allowed(r"rg -n 'C:\\Users\\nguco' /home/crunch/.codex/hooks")

    def test_blocks_windows_shells(self) -> None:
        self.assert_blocked("cmd.exe /d /c git status", "open Tabby and leave it running")
        self.assert_blocked(r"C:\\Windows\\System32\\cmd.exe /d /c dir", "Tabby MCP")
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
        self.assert_blocked("'/mnt/c/Program Files/Git/cmd/git' status", "Windows git.exe")
        self.assert_blocked(r"C:\\Git\\cmd\\git.exe status", "Windows git.exe")
        self.assert_blocked(r"'C:\\Program Files\\Git\\bin\\bash.exe' --version", "Windows executable")

    def test_blocks_ps1_as_executable(self) -> None:
        self.assert_blocked("tools/repair.ps1", "PS1")

    def test_blocks_nested_windows_shells(self) -> None:
        self.assert_blocked("bash -lc 'cmd.exe /c git status'", "Tabby MCP")
        self.assert_blocked("env -S 'cmd.exe /d /c dir'", "Tabby MCP")
        self.assert_blocked(
            "env --split-string='powershell.exe -NoProfile -Command Get-ChildItem'",
            "Tabby MCP",
        )
        self.assert_blocked("env -S'git.exe status'", "Windows git.exe")
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

    def test_blocks_shell_heredoc_bodies(self) -> None:
        self.assert_blocked("bash <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("bash <<'EOF'\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_blocked("bash <<'EOF'\n EOF\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("bash <<'EOF-MARKER'\ncmd.exe /d /c dir\nEOF-MARKER\n", "Tabby MCP")
        self.assert_blocked("cat <<'EOF' | bash\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("<<'EOF' bash\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("0<<'EOF' bash\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("<<'EOF' env bash\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("00<<'EOF' bash\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("bash -s 1 <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("env <<'EOF' bash\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("{ bash; } <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("if true; then bash; fi <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("env -S'bash -s' <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("env -vS'bash -s' <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("env -S'bash\\_-s' <<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("bash \\\n<<'EOF'\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_allowed("5<<'EOF' bash\ncmd.exe /d /c dir\nEOF\n")
        self.assert_allowed("5 <<'EOF' bash\ncmd.exe /d /c dir\nEOF\n")
        self.assert_allowed("bash <<'DATA' <<'SCRIPT'\ncmd.exe /d /c dir\nDATA\necho safe\nSCRIPT\n")
        self.assert_allowed("cat <<'DATA'; bash <<'SCRIPT'\ncmd.exe /d /c dir\nDATA\necho safe\nSCRIPT\n")
        self.assert_allowed("cat <<'DATA' | bash <<'SCRIPT'\ncmd.exe /d /c dir\nDATA\necho safe\nSCRIPT\n")
        self.assert_blocked("bash <<'EOF'\npowershell.exe -NoProfile -Command Get-ChildItem\nEOF\n", "Tabby MCP")
        self.assert_blocked("bash <<'EOF'\nwt.exe new-tab powershell.exe\nEOF\n", "Windows Terminal")
        self.assert_blocked("bash <<'EOF'\ntools/repair.ps1\nEOF\n", "PS1")
        self.assert_blocked("bash <<'EOF'\nTMPDIR=/mnt/c/tmp pnpm test\nEOF\n", "TMPDIR")

    def test_blocks_unquoted_heredoc_command_substitutions(self) -> None:
        self.assert_blocked("cat <<EOF\n$(cmd.exe /d /c dir)\nEOF\n", "Tabby MCP")
        self.assert_blocked("cat <<EOF\n$(( $(cmd.exe /d /c dir) + 1 ))\nEOF\n", "Tabby MCP")
        self.assert_blocked("cat <<EOF\n$\\\n(cmd.exe /d /c dir)\nEOF\n", "Tabby MCP")
        self.assert_blocked("cat <<EOF\n${ cmd.exe /d /c dir; }\nEOF\n", "Tabby MCP")

    def test_quoted_heredoc_operator_text_does_not_hide_later_commands(self) -> None:
        self.assert_blocked('echo "<<EOF"\ncmd.exe /d /c dir\nEOF\n', "Tabby MCP")

    def test_non_heredoc_shift_forms_do_not_hide_later_commands(self) -> None:
        self.assert_blocked("cat <<<foo\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_blocked("# <<EOF\ncmd.exe /d /c dir\nEOF\n", "Tabby MCP")
        self.assert_blocked("echo $((1<<2))\nTMPDIR=/mnt/c/tmp pnpm test\n", "TMPDIR")

    def test_tab_stripped_heredoc_terminator_does_not_hide_later_commands(self) -> None:
        self.assert_blocked("cat <<-EOF\nsafe\n\tEOF\ncmd.exe /d /c dir\n", "Tabby MCP")

    def test_shell_quoted_delimiter_does_not_hide_later_commands(self) -> None:
        self.assert_blocked("cat <<$'EOF'\nsafe\nEOF\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_blocked('cat <<$"EOF"\nsafe\nEOF\ncmd.exe /d /c dir\n', "Tabby MCP")
        self.assert_blocked("cat <<$'EO\\x46'\nsafe\nEOF\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_blocked("cat <<$'EO\\u0046'\nsafe\nEOF\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_blocked('cat <<"EO\\F"\nsafe\nEO\\F\ncmd.exe /d /c dir\n', "Tabby MCP")
        self.assert_blocked('cat <<"EO\\\nF"\nsafe\nEOF\ncmd.exe /d /c dir\n', "Tabby MCP")

    def test_unquoted_heredoc_line_continuation_can_form_terminator(self) -> None:
        self.assert_blocked("cat <<EOF\nEO\\\nF\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_blocked("cat <<EO\\\nF\nsafe\nEOF\ncmd.exe /d /c dir\n", "Tabby MCP")
        self.assert_allowed("cat <<EO\\\nF\ncmd.exe /d /c dir\nEOF\n")
        self.assert_blocked("cat <<EOF> /tmp/out\nsafe\nEOF\ncmd.exe /d /c dir\n", "Tabby MCP")

    def test_deep_nesting_denial_keeps_tabby_reminder(self) -> None:
        command = "cmd.exe /d /c dir"
        for _ in range(7):
            command = f"bash -lc {shlex.quote(command)}"
        self.assert_blocked(command, "open Tabby and leave it running")

    def test_rewrites_raw_windows_paths_in_wsl_commands(self) -> None:
        self.assert_rewritten(
            r"python C:\\Users\\nguco\\script.py",
            "python /mnt/c/Users/nguco/script.py",
        )
        self.assert_rewritten(
            r"python 'C:\\Users\\Nate Foo\\script.py'",
            "python '/mnt/c/Users/Nate Foo/script.py'",
        )

    def test_preserves_heredoc_body_during_path_rewrite(self) -> None:
        command = "python C:\\Users\\nguco\\script.py <<'EOF'\nliteral C:\\Users\\secret\\value\nEOF\n"
        expected = "python /mnt/c/Users/nguco/script.py <<'EOF'\nliteral C:\\Users\\secret\\value\nEOF\n"

        self.assert_rewritten(command, expected)

    def test_preserves_indented_heredoc_body_during_path_rewrite(self) -> None:
        command = "python C:\\Users\\nguco\\script.py <<'EOF'\n EOF\nliteral C:\\Users\\secret\\value\nEOF\n"
        expected = "python /mnt/c/Users/nguco/script.py <<'EOF'\n EOF\nliteral C:\\Users\\secret\\value\nEOF\n"

        self.assert_rewritten(command, expected)

    def test_rewrites_shell_heredoc_body_paths(self) -> None:
        command = 'bash <<\'EOF\'\npython "C:\\Users\\nguco\\script.py"\nEOF\n'
        expected = 'bash <<\'EOF\'\npython "/mnt/c/Users/nguco/script.py"\nEOF\n'

        self.assert_rewritten(command, expected)

    def test_preserves_non_shell_data_heredoc_body_paths(self) -> None:
        self.assert_allowed('grep bash <<\'EOF\'\nC:\\Users\\nguco\\note.txt\nEOF\n')

    def test_parallel_cmd_rewrite_preserves_cmd_key(self) -> None:
        result = run_parallel_cmd_hook(r"python C:\\Users\\nguco\\script.py")
        self.assertEqual(result.returncode, 0, result.stderr)
        output = json.loads(result.stdout)
        parameters = output["hookSpecificOutput"]["updatedInput"]["tool_uses"][0]["parameters"]
        self.assertEqual(parameters, {"cmd": "python /mnt/c/Users/nguco/script.py"})

    def test_allows_wsl_style_mounted_windows_paths(self) -> None:
        self.assert_allowed("python /mnt/c/Users/nguco/script.py")

    def test_blocks_temp_on_windows_mounts(self) -> None:
        self.assert_blocked("TMPDIR=/mnt/c/tmp pnpm test", "TMPDIR")
        self.assert_blocked("TMPDIR='/mnt/c/tmp' pnpm test", "TMPDIR")
        self.assert_blocked("export TMPDIR='/mnt/c/tmp'; pnpm test", "TMPDIR")
        self.assert_blocked("TMPDIR=$'/mnt/c/tmp' pnpm test", "TMPDIR")
        self.assert_blocked('TMPDIR=/mnt/"c"/tmp pnpm test', "TMPDIR")
        self.assert_blocked("export TMPDIR=$'/mnt/c/tmp'; pnpm test", "TMPDIR")
        self.assert_blocked("TMPDIR=$'\\057mnt\\057c\\057tmp' pnpm test", "TMPDIR")
        self.assert_blocked("declare -x TMPDIR=/mnt/c/tmp; pnpm test", "TMPDIR")
        self.assert_blocked("env -u PATH TMPDIR=/mnt/c/tmp pnpm test", "TMPDIR")
        self.assert_blocked("env -S'TMPDIR=/mnt/c/tmp pnpm test'", "TMPDIR")
        self.assert_blocked("env -vS'TMPDIR=/mnt/c/tmp pnpm test'", "TMPDIR")

    def test_allows_temp_assignment_text_arguments(self) -> None:
        self.assert_allowed("echo 'TMPDIR=/mnt/c/tmp'")

    def test_preserves_code_strings_during_path_rewrite(self) -> None:
        self.assert_allowed(r'python -c "print(\'C:\\Users\\nguco\')"')
        self.assert_allowed(r"python -c 'print(\"C:\\Users\\nguco\")'")
        self.assert_blocked("bash -lc 'echo safe\ncmd.exe /d /c dir'", "Tabby MCP")

    def test_rewrites_search_file_operands(self) -> None:
        self.assert_rewritten(
            r"grep needle C:\\Users\\nguco\\file.txt",
            "grep needle /mnt/c/Users/nguco/file.txt",
        )
        self.assert_rewritten(
            r"rg needle C:\\Users\\nguco",
            "rg needle /mnt/c/Users/nguco",
        )
        self.assert_rewritten(
            r"grep 'C:\\Users\\needle' C:\\Users\\nguco\\file.txt",
            r"grep 'C:\\Users\\needle' /mnt/c/Users/nguco/file.txt",
        )
        self.assert_rewritten(
            r"grep 'C:\\Users\\needle' 'C:\\Users\\nguco\\file.txt'",
            r"grep 'C:\\Users\\needle' '/mnt/c/Users/nguco/file.txt'",
        )
        self.assert_rewritten(
            r"grep -e C:\\Users\\needle C:\\Users\\nguco\\file.txt",
            r"grep -e C:\\Users\\needle /mnt/c/Users/nguco/file.txt",
        )

    def test_blocks_linux_symlink_on_windows_mount(self) -> None:
        self.assert_blocked("ln -s /tmp/source /mnt/c/Users/nguco/bin/tool", "symlinks")


if __name__ == "__main__":
    unittest.main()
