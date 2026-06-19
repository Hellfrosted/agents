#!/usr/bin/env python3
"""Block WSL shell footguns before Codex runs them."""

from __future__ import annotations

import json
import os
import re
import shlex
import sys


COMMAND_SEPARATORS = {
    ";",
    "&",
    "&&",
    "|",
    "||",
    "(",
    ")",
    "{",
    "}",
    "\n",
    "then",
    "else",
    "elif",
    "do",
    "fi",
    "done",
}
SHELLS = {"bash", "sh", "zsh"}
SHELL_TOOL_NAMES = {
    "Bash",
    "exec_command",
    "functions.exec_command",
    "shell_command",
    "functions.shell_command",
}
WINDOWS_SHELL_EXECUTABLES = {"cmd", "powershell", "pwsh"}
WINDOWS_TERMINAL_EXECUTABLES = {"wt", "windowsterminal", "windows-terminal"}
TABBY_NUDGE = "If Tabby MCP is not accessible, remind the user to open Tabby and leave it running in the background."
WINDOWS_PATH = re.compile(r"(^|[\s'\"`])([A-Za-z]:[\\/][^\s'\"`;&|)]*)")
MNT_PATH = re.compile(r"^/mnt/[A-Za-z](/|$)")
TMP_ASSIGNMENT = re.compile(r"(^|[\s;&|()])(TMPDIR|TEMP|TMP)=/mnt/[A-Za-z](/|[\s;&|)]|$)")
WINDOWS_PATH_BRIDGE_ALLOWLIST = {"wslpath", "explorer"}
WRAPPERS = {"command", "exec", "noglob", "time"}
NO_ARG_LAUNCH_WRAPPERS = {"/init", "nohup", "setsid", "winpty"}


def shell_commands(payload: dict) -> list[str]:
    tool_input = payload.get("tool_input") or {}
    if not isinstance(tool_input, dict):
        return []

    tool_name = payload.get("tool_name")
    if tool_name in SHELL_TOOL_NAMES:
        command = tool_input.get("command") or tool_input.get("cmd")
        return [command] if isinstance(command, str) and command.strip() else []

    commands: list[str] = []
    tool_uses = tool_input.get("tool_uses")
    if isinstance(tool_uses, list):
        for tool_use in tool_uses:
            if not isinstance(tool_use, dict):
                continue
            nested_name = tool_use.get("recipient_name") or tool_use.get("tool_name")
            if nested_name not in SHELL_TOOL_NAMES:
                continue
            parameters = tool_use.get("parameters") or {}
            if not isinstance(parameters, dict):
                continue
            command = parameters.get("command") or parameters.get("cmd")
            if isinstance(command, str) and command.strip():
                commands.append(command)
    return commands


def shell_tokens(command: str) -> list[str]:
    lexer = shlex.shlex(command, posix=True, punctuation_chars=True)
    lexer.whitespace_split = True
    return list(lexer)


def split_segments(tokens: list[str]) -> list[list[str]]:
    segments: list[list[str]] = []
    segment: list[str] = []
    for token in tokens:
        if token in COMMAND_SEPARATORS:
            if segment:
                segments.append(segment)
                segment = []
            continue
        segment.append(token)
    if segment:
        segments.append(segment)
    return segments


def executable_name(token: str) -> str:
    name = re.split(r"[\\/]", token)[-1].lower()
    return name[:-4] if name.endswith(".exe") else name


def is_env_assignment(token: str) -> bool:
    if "=" not in token or token.startswith(("=", "-", "./", "../", "/")):
        return False
    key = token.split("=", 1)[0]
    return bool(key) and not key[0].isdigit() and key.replace("_", "").isalnum()


def normalize_segment(segment: list[str]) -> list[str]:
    normalized = list(segment)

    while normalized:
        head = normalized[0]
        if head in WRAPPERS:
            normalized.pop(0)
            continue
        if head in NO_ARG_LAUNCH_WRAPPERS:
            normalized.pop(0)
            while normalized and normalized[0].startswith("-"):
                normalized.pop(0)
            continue
        if head == "timeout":
            normalized.pop(0)
            while normalized and normalized[0].startswith("-"):
                option = normalized.pop(0)
                if option in {"-k", "--kill-after", "-s", "--signal"} and normalized:
                    normalized.pop(0)
            if normalized:
                normalized.pop(0)
            continue
        if head == "nice":
            normalized.pop(0)
            while normalized and normalized[0].startswith("-"):
                option = normalized.pop(0)
                if option == "-n" and normalized:
                    normalized.pop(0)
            continue
        if head == "stdbuf":
            normalized.pop(0)
            while normalized and normalized[0].startswith("-"):
                option = normalized.pop(0)
                if option in {"-i", "-o", "-e"} and normalized:
                    normalized.pop(0)
            continue
        if is_env_assignment(head):
            normalized.pop(0)
            continue
        if head == "env":
            normalized.pop(0)
            while normalized and (normalized[0].startswith("-") or is_env_assignment(normalized[0])):
                normalized.pop(0)
            continue
        if head == "sudo":
            normalized.pop(0)
            while len(normalized) >= 2 and normalized[0] in {"-u", "-g", "-h", "-p"}:
                del normalized[:2]
            while normalized and normalized[0].startswith("-"):
                normalized.pop(0)
            continue
        break

    return normalized


def segment_executable_token(segment: list[str]) -> str:
    for token in normalize_segment(segment):
        if is_env_assignment(token):
            continue
        return token
    return ""


def segment_executable(segment: list[str]) -> str:
    token = segment_executable_token(segment)
    if not token:
        return ""
    return executable_name(token)


def shell_command_string(segment: list[str]) -> str | None:
    args = normalize_segment(segment)
    if not args or executable_name(args[0]) not in SHELLS:
        return None

    index = 1
    while index < len(args):
        option = args[index]
        if option == "--":
            index += 1
            continue
        if not option.startswith("-"):
            return None
        if option in {"-o", "+o", "-O", "+O"}:
            index += 2
            continue
        if "c" in option[1:]:
            return args[index + 1] if index + 1 < len(args) else ""
        if "o" in option[1:] or "O" in option[1:]:
            index += 2
            continue
        index += 1

    return None


def is_windows_path_token(token: str) -> bool:
    return bool(re.match(r"^[A-Za-z]:[\\/]", token) or token.startswith("/mnt/c/") or token.startswith("/mnt/C/"))


def raw_executable_name(token: str) -> str:
    return re.split(r"[\\/]", token)[-1].lower()


def is_windows_git_executable(token: str) -> bool:
    raw_name = raw_executable_name(token)
    return raw_name == "git.exe" or (raw_name == "git" and is_windows_path_token(token))


def is_ps1_executable(token: str) -> bool:
    return raw_executable_name(token).endswith(".ps1")


def rtk_nested_command(segment: list[str]) -> str | None:
    args = normalize_segment(segment)
    if not args or executable_name(args[0]) != "rtk":
        return None

    rtk_args = args[1:]
    while rtk_args and rtk_args[0].startswith("-"):
        rtk_args.pop(0)
    if not rtk_args:
        return None
    if rtk_args[0] == "run" and len(rtk_args) >= 2:
        return rtk_args[1]
    if rtk_args[0] == "proxy" and len(rtk_args) >= 2:
        return shlex.join(rtk_args[1:])
    return shlex.join(rtk_args)


def contains_ln_s_on_windows_mount(segment: list[str]) -> bool:
    if segment_executable(segment) != "ln":
        return False
    args = [token for token in segment[1:] if not is_env_assignment(token)]
    has_symlink_flag = any(token == "-s" or (token.startswith("-") and "s" in token[1:]) for token in args)
    return has_symlink_flag and any(MNT_PATH.match(token) for token in args if not token.startswith("-"))


def allowed_windows_path_bridge(command: str) -> bool:
    try:
        segments = split_segments(shell_tokens(command))
    except ValueError:
        return False
    return any(segment_executable(segment) in WINDOWS_PATH_BRIDGE_ALLOWLIST for segment in segments)


def blocked_segment_reason(segment: list[str], depth: int) -> str | None:
    nested_shell = shell_command_string(segment)
    if nested_shell is not None:
        return blocked_reason(nested_shell, depth + 1)

    nested_rtk = rtk_nested_command(segment)
    if nested_rtk:
        return blocked_reason(nested_rtk, depth + 1)

    executable_token = segment_executable_token(segment)
    executable = executable_name(executable_token) if executable_token else ""

    if executable in WINDOWS_TERMINAL_EXECUTABLES:
        return (
            "Blocked: do not launch Windows Terminal from Codex shell commands. "
            "Use the Tabby MCP terminal/session tools for Windows-side terminal work. "
            f"{TABBY_NUDGE}"
        )

    if executable in WINDOWS_SHELL_EXECUTABLES:
        return (
            "Blocked: do not run Windows CMD or PowerShell from Codex shell commands. "
            "Use Tabby MCP for Windows-side CMD, PowerShell, PS1, or Windows Git work. "
            f"{TABBY_NUDGE}"
        )

    if executable_token and is_windows_git_executable(executable_token):
        return (
            "Blocked: do not run Windows git.exe from Codex shell commands. "
            "Use WSL-native git here, or Tabby MCP for Windows-side Git work. "
            f"{TABBY_NUDGE}"
        )

    if executable_token and is_ps1_executable(executable_token):
        return (
            "Blocked: do not run PS1 scripts from Codex shell commands. "
            "Use Tabby MCP for Windows-side PowerShell work. "
            f"{TABBY_NUDGE}"
        )

    return None


def blocked_reason(command: str, depth: int = 0) -> str | None:
    if depth > 5:
        return f"Blocked: deeply nested shell command could not be safely inspected. {TABBY_NUDGE}"

    if TMP_ASSIGNMENT.search(command):
        return "Blocked: keep TMPDIR, TEMP, and TMP on Linux storage such as /tmp, not /mnt/*."

    try:
        tokens = shell_tokens(command)
    except ValueError:
        return None

    for segment in split_segments(tokens):
        reason = blocked_segment_reason(segment, depth)
        if reason:
            return reason

    if WINDOWS_PATH.search(command) and not allowed_windows_path_bridge(command):
        return "Blocked: translate Windows paths with wslpath before using them in WSL shell commands."

    for segment in split_segments(tokens):
        if contains_ln_s_on_windows_mount(segment):
            return (
                "Blocked: create symlinks on Windows drives through Tabby MCP using "
                "Windows-native mklink or New-Item in an existing Windows-side session. "
                f"{TABBY_NUDGE}"
            )

    return None


def deny(reason: str) -> None:
    print(
        json.dumps(
            {
                "hookSpecificOutput": {
                    "hookEventName": "PreToolUse",
                    "permissionDecision": "deny",
                    "permissionDecisionReason": reason,
                }
            }
        )
    )


def main() -> int:
    try:
        payload = json.load(sys.stdin)
    except json.JSONDecodeError:
        return 0

    if payload.get("hook_event_name") != "PreToolUse":
        return 0

    for command in shell_commands(payload):
        reason = blocked_reason(command)
        if reason:
            deny(reason)
            return 0

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
