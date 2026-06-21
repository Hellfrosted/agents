#!/usr/bin/env python3
"""Block WSL shell footguns before Codex runs them."""

# allow: SIZE_OK - standalone Codex hook entrypoint copied into CODEX_HOME.

from __future__ import annotations

import json
import re
import shlex
import sys
from typing import NamedTuple


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
WINDOWS_PATH = re.compile(r"(^|[\s`])([A-Za-z]:[\\/][^\s'\"`;&|)]*)")
QUOTED_WINDOWS_PATH = re.compile(r"(^|[\s;&|()])(['\"])([A-Za-z]:[\\/][^'\"]*)\2(?=$|[\s;&|)])")
MNT_PATH = re.compile(r"^/mnt/[A-Za-z](/|$)")
TMP_ASSIGNMENT = re.compile(
    r"(^|[\s;&|()])(TMPDIR|TEMP|TMP)=\$?['\"]?/mnt/[A-Za-z](?:/[^'\"\s;&|)]*)?['\"]?(?=$|[\s;&|)])"
)
WINDOWS_PATH_BRIDGE_ALLOWLIST = {"wslpath", "explorer"}
SAFE_TEXT_COMMANDS = {"echo", "printf"}
SEARCH_TEXT_COMMANDS = {"rg", "grep"}
CODE_STRING_COMMANDS = {"python", "python3", "node", "perl", "ruby"}
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


SHELL_WORD_DELIMITERS = set(" \t\r\n;&|()<>")


class HeredocSpec(NamedTuple):
    delimiter: str
    strip_tabs: bool
    quoted: bool
    fd: str | None


class HeredocBody(NamedTuple):
    text: str
    executable_shell: bool
    expands: bool


class HeredocPartial(NamedTuple):
    delimiter: str
    strip_tabs: bool
    quoted: bool
    logical_prefix: str
    quote: str | None
    fd: str | None


def heredoc_specs(line: str) -> list[HeredocSpec]:
    specs, _ = heredoc_specs_and_partial(line)
    return specs


def heredoc_specs_and_partial(
    line: str, partial: HeredocPartial | None = None
) -> tuple[list[HeredocSpec], HeredocPartial | None]:
    specs: list[HeredocSpec] = []
    if partial is not None:
        delimiter, end, quoted, continued, quote = parse_heredoc_delimiter(
            line,
            0,
            initial=partial.delimiter,
            initial_quoted=partial.quoted,
            initial_quote=partial.quote,
        )
        if continued:
            logical_prefix = partial.logical_prefix + line[:-1]
            return specs, HeredocPartial(
                delimiter or partial.delimiter,
                partial.strip_tabs,
                quoted,
                logical_prefix,
                quote,
                partial.fd,
            )
        if delimiter:
            specs.append(HeredocSpec(delimiter, partial.strip_tabs, quoted, partial.fd))
        more_specs, next_partial = scan_heredoc_specs(line, end)
        specs.extend(more_specs)
        return specs, next_partial
    return scan_heredoc_specs(line)


def scan_heredoc_specs(line: str, index: int = 0) -> tuple[list[HeredocSpec], HeredocPartial | None]:
    specs: list[HeredocSpec] = []
    quote: str | None = None

    while index < len(line):
        char = line[index]
        if quote is None:
            if is_comment_start(line, index):
                break
            if line.startswith("$((", index):
                index = skip_arithmetic_expansion(line, index + 3)
                continue
            if line.startswith("((", index):
                index = skip_arithmetic_expansion(line, index + 2)
                continue
            if char == "\\":
                index += 2
                continue
            if char in {"'", '"'}:
                quote = char
                index += 1
                continue
            if line.startswith("<<<", index):
                index += 3
                continue
            if line.startswith("<<", index):
                fd = heredoc_fd_prefix(line, index)
                start = index + 2
                strip_tabs = start < len(line) and line[start] == "-"
                if strip_tabs:
                    start += 1
                while start < len(line) and line[start] in " \t":
                    start += 1
                delimiter, end, quoted, continued, quote = parse_heredoc_delimiter(line, start)
                if continued:
                    logical_prefix = line[:-1] if line.endswith("\\") else line[:end]
                    partial = HeredocPartial(delimiter or "", strip_tabs, quoted, logical_prefix, quote, fd)
                    return specs, partial
                if delimiter:
                    specs.append(HeredocSpec(delimiter, strip_tabs, quoted, fd))
                    index = end
                    continue
            index += 1
            continue

        if char == quote:
            quote = None
        elif quote == '"' and char == "\\":
            index += 1
        index += 1

    return specs, None


def heredoc_delimiters(line: str) -> list[str]:
    return [spec.delimiter for spec in heredoc_specs(line)]


def heredoc_fd_prefix(line: str, operator_index: int) -> str | None:
    end = operator_index
    start = end
    while start > 0 and line[start - 1].isdigit():
        start -= 1
    if start == end:
        return None
    if start > 0 and line[start - 1] not in " \t;&|({":
        return None
    return line[start:end]


def is_comment_start(line: str, index: int) -> bool:
    return line[index] == "#" and (index == 0 or line[index - 1] in " \t;&|()")


def skip_arithmetic_expansion(line: str, index: int) -> int:
    depth = 1
    while index < len(line):
        if line.startswith("((", index):
            depth += 1
            index += 2
            continue
        if line.startswith("))", index):
            depth -= 1
            index += 2
            if depth == 0:
                return index
            continue
        if line[index] == "\\":
            index += 2
            continue
        index += 1
    return len(line)


def parse_heredoc_delimiter(
    line: str,
    start: int,
    initial: str = "",
    initial_quoted: bool = False,
    initial_quote: str | None = None,
) -> tuple[str | None, int, bool, bool, str | None]:
    delimiter: list[str] = list(initial)
    index = start
    quote = initial_quote
    quoted = initial_quoted

    while index < len(line):
        char = line[index]
        if quote is None:
            if char in SHELL_WORD_DELIMITERS:
                break
            if char == "$" and index + 1 < len(line) and line[index + 1] in {"'", '"'}:
                quoted = True
                quote = "ansi" if line[index + 1] == "'" else '"'
                index += 2
                continue
            if char in {"'", '"'}:
                quoted = True
                quote = char
                index += 1
                continue
            if char == "\\":
                if index + 1 < len(line):
                    quoted = True
                    delimiter.append(line[index + 1])
                    index += 2
                    continue
                return "".join(delimiter) or None, len(line), quoted, True, quote
            delimiter.append(char)
            index += 1
            continue

        if quote == "ansi":
            if char == "'":
                quote = None
                index += 1
                continue
            if char == "\\":
                value, index = read_ansi_c_escape(line, index)
                delimiter.append(value)
                continue
            delimiter.append(char)
            index += 1
            continue

        if char == quote:
            quote = None
        elif quote == '"' and char == "\\":
            if index + 1 >= len(line):
                return "".join(delimiter) or None, len(line), quoted, True, quote
            elif line[index + 1] in {'$', '`', '"', "\\"}:
                delimiter.append(line[index + 1])
                index += 2
                continue
            else:
                delimiter.append("\\")
        else:
            delimiter.append(char)
        index += 1

    return ("".join(delimiter) or None, index, quoted, False, quote)


def read_ansi_c_escape(line: str, index: int) -> tuple[str, int]:
    if index + 1 >= len(line):
        return "\\", index + 1
    char = line[index + 1]
    escapes = {
        "a": "\a",
        "b": "\b",
        "e": "\x1b",
        "E": "\x1b",
        "f": "\f",
        "n": "\n",
        "r": "\r",
        "t": "\t",
        "v": "\v",
        "\\": "\\",
        "'": "'",
        '"': '"',
        "?": "?",
    }
    if char in escapes:
        return escapes[char], index + 2
    if char == "x":
        digits = re.match(r"[0-9A-Fa-f]{1,2}", line[index + 2 :])
        if digits:
            return chr(int(digits.group(0), 16)), index + 2 + len(digits.group(0))
    if char == "u":
        digits = re.match(r"[0-9A-Fa-f]{1,4}", line[index + 2 :])
        if digits:
            return chr(int(digits.group(0), 16)), index + 2 + len(digits.group(0))
    if char == "U":
        digits = re.match(r"[0-9A-Fa-f]{1,8}", line[index + 2 :])
        if digits:
            return chr(int(digits.group(0), 16)), index + 2 + len(digits.group(0))
    if char in "01234567":
        digits = re.match(r"[0-7]{1,3}", line[index + 1 :])
        if digits:
            return chr(int(digits.group(0), 8)), index + 1 + len(digits.group(0))
    return char, index + 2


def is_heredoc_terminator(content: str, spec: HeredocSpec) -> bool:
    candidate = content.lstrip("\t") if spec.strip_tabs else content
    return candidate == spec.delimiter


def heredoc_content_fragment(content: str, spec: HeredocSpec, continuation: str) -> tuple[str, str | None]:
    candidate = content.lstrip("\t") if spec.strip_tabs else content
    if not spec.quoted and continuation:
        candidate = continuation + candidate
    if not spec.quoted and candidate.endswith("\\"):
        return candidate, candidate[:-1]
    return candidate, None


def command_without_heredoc_bodies(command: str) -> str:
    lines = command.splitlines(keepends=True)
    output: list[str] = []
    pending: list[HeredocSpec] = []
    skipping: tuple[HeredocSpec, str] | None = None
    partial: HeredocPartial | None = None

    for line in lines:
        content = line.rstrip("\r\n")
        terminator = skipping
        if terminator is not None:
            spec, continuation = terminator
            candidate, next_continuation = heredoc_content_fragment(content, spec, continuation)
            if next_continuation is not None:
                skipping = (spec, next_continuation)
                continue
            if candidate == spec.delimiter:
                next_spec = pending.pop(0) if pending else None
                skipping = (next_spec, "") if next_spec else None
            else:
                skipping = (spec, "")
            continue

        output.append(line)
        specs, partial = heredoc_specs_and_partial(content, partial)
        pending.extend(specs)
        if partial is not None:
            continue
        if pending:
            skipping = (pending.pop(0), "")

    return "".join(output)


def transform_without_heredoc_bodies(command: str, transform) -> str:
    lines = command.splitlines(keepends=True)
    output: list[str] = []
    pending: list[tuple[HeredocSpec, bool]] = []
    skipping: tuple[HeredocSpec, bool, str] | None = None
    partial: HeredocPartial | None = None
    continued_command_prefix = ""

    for line in lines:
        content = line.rstrip("\r\n")
        terminator = skipping
        if terminator is not None:
            spec, executable_shell, continuation = terminator
            candidate, next_continuation = heredoc_content_fragment(content, spec, continuation)
            if next_continuation is not None:
                output.append(transform(line) if executable_shell else line)
                skipping = (spec, executable_shell, next_continuation)
                continue
            if candidate == spec.delimiter:
                output.append(line)
                next_pending = pending.pop(0) if pending else None
                skipping = (next_pending[0], next_pending[1], "") if next_pending else None
            else:
                output.append(transform(line) if executable_shell else line)
                skipping = (spec, executable_shell, "")
            continue

        output.append(transform(line))
        logical_line = continued_command_prefix + (
            partial.logical_prefix + content if partial is not None else content
        )
        specs, partial = heredoc_specs_and_partial(content, partial)
        shell_flags = heredoc_shell_flags(logical_line)
        pending.extend((spec, shell_flags[index] if index < len(shell_flags) else False) for index, spec in enumerate(specs))
        if partial is not None:
            continue
        if not specs and shell_line_continues(content):
            continued_command_prefix += content[:-1]
            continue
        continued_command_prefix = ""
        if pending:
            spec, shell_input = pending.pop(0)
            skipping = (spec, shell_input, "")

    return "".join(output)


def heredoc_bodies(command: str) -> list[HeredocBody]:
    lines = command.splitlines()
    bodies: list[HeredocBody] = []
    pending: list[tuple[HeredocSpec, bool]] = []
    skipping: tuple[HeredocSpec, bool, str, list[str]] | None = None
    partial: HeredocPartial | None = None
    continued_command_prefix = ""

    for line in lines:
        content = line.rstrip("\r\n")
        if skipping is not None:
            spec, executable_shell, continuation, body = skipping
            candidate, next_continuation = heredoc_content_fragment(content, spec, continuation)
            if next_continuation is not None:
                body.append(line)
                skipping = (spec, executable_shell, next_continuation, body)
                continue
            if candidate == spec.delimiter:
                bodies.append(HeredocBody("\n".join(body), executable_shell, not spec.quoted))
                next_pending = pending.pop(0) if pending else None
                skipping = (next_pending[0], next_pending[1], "", []) if next_pending else None
            else:
                body.append(line)
                skipping = (spec, executable_shell, "", body)
            continue

        logical_line = continued_command_prefix + (
            partial.logical_prefix + content if partial is not None else content
        )
        specs, partial = heredoc_specs_and_partial(content, partial)
        shell_flags = heredoc_shell_flags(logical_line)
        pending.extend((spec, shell_flags[index] if index < len(shell_flags) else False) for index, spec in enumerate(specs))
        if partial is not None:
            continue
        if not specs and shell_line_continues(content):
            continued_command_prefix += content[:-1]
            continue
        continued_command_prefix = ""
        if pending:
            spec, shell_input = pending.pop(0)
            skipping = (spec, shell_input, "", [])

    if skipping is not None:
        spec, executable_shell, _, body = skipping
        bodies.append(HeredocBody("\n".join(body), executable_shell, not spec.quoted))

    return bodies


def shell_line_continues(content: str) -> bool:
    if not content.endswith("\\"):
        return False

    quote: str | None = None
    index = 0
    while index < len(content) - 1:
        char = content[index]
        if quote is None:
            if char == "#":
                break
            if char in {"'", '"'}:
                quote = char
                index += 1
                continue
            if char == "\\":
                index += 2
                continue
            index += 1
            continue

        if char == quote:
            quote = None
        elif quote == '"' and char == "\\":
            index += 1
        index += 1

    if quote == "'":
        return False

    trailing_backslashes = len(content) - len(content.rstrip("\\"))
    return trailing_backslashes % 2 == 1


def line_has_shell_heredoc_input(line: str) -> bool:
    return any(heredoc_shell_flags(line))


def heredoc_shell_flags(line: str) -> list[bool]:
    specs = heredoc_specs(line)
    if compound_shell_heredoc(line):
        return [heredoc_targets_stdin(spec) for spec in specs]
    try:
        tokens = shell_tokens(line)
    except ValueError:
        return [False] * len(specs)

    flags: list[bool] = []
    spec_index = 0
    for group in pipeline_groups(tokens):
        segment_specs: list[list[HeredocSpec]] = []
        for segment in group:
            count = segment.count("<<")
            segment_specs.append(specs[spec_index : spec_index + count])
            spec_index += count

        downstream_pipe_shell = False
        segment_shell_flags = [False] * len(group)
        for index in range(len(group) - 1, -1, -1):
            segment = group[index]
            segment_is_shell = segment_shell_sink(segment, segment_specs[index])
            segment_shell_flags[index] = segment_is_shell or downstream_pipe_shell
            has_stdin_heredoc = any(heredoc_targets_stdin(spec) for spec in segment_specs[index])
            downstream_pipe_shell = (segment_is_shell or downstream_pipe_shell) and not has_stdin_heredoc

        for index, _segment in enumerate(group):
            flags.extend(stdin_heredoc_flags(segment_specs[index], segment_shell_flags[index]))
    return flags


def pipeline_groups(tokens: list[str]) -> list[list[list[str]]]:
    groups: list[list[list[str]]] = []
    group: list[list[str]] = []
    segment: list[str] = []
    for token in tokens:
        if token == "|":
            if segment:
                group.append(segment)
                segment = []
            continue
        if token in COMMAND_SEPARATORS:
            if segment:
                group.append(segment)
                segment = []
            if group:
                groups.append(group)
                group = []
            continue
        segment.append(token)
    if segment:
        group.append(segment)
    if group:
        groups.append(group)
    return groups


def stdin_heredoc_flags(specs: list[HeredocSpec], shell_sink: bool) -> list[bool]:
    last_stdin_index = -1
    for index, spec in enumerate(specs):
        if heredoc_targets_stdin(spec):
            last_stdin_index = index
    return [shell_sink and index == last_stdin_index for index, _spec in enumerate(specs)]


def heredoc_targets_stdin(spec: HeredocSpec) -> bool:
    if spec.fd is None:
        return True
    try:
        return int(spec.fd, 10) == 0
    except ValueError:
        return False


def compound_shell_heredoc(line: str) -> bool:
    shell_pattern = r"(^|[\s;{(])(bash|sh|zsh)($|[\s;})])"
    shell_word = r"\b(bash|sh|zsh)\b"
    return bool(
        re.search(r"\{[^}\n]*" + shell_pattern + r"[^}\n]*}\s*\d*<<", line)
        or re.search(r"\([^)\n]*" + shell_pattern + r"[^)\n]*\)\s*\d*<<", line)
        or re.search(r"\bif\b.*\bthen\b.*" + shell_word + r".*\bfi\s*\d*<<", line)
        or re.search(r"\b(while|until)\b.*\bdo\b.*" + shell_word + r".*\bdone\s*\d*<<", line)
        or re.search(r"\bfor\b.*\bdo\b.*" + shell_word + r".*\bdone\s*\d*<<", line)
        or re.search(r"(\bfunction\s+\w+|\b\w+\s*\(\))\s*\{.*" + shell_word + r".*}\s*\d*<<", line)
    )


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


def contains_windows_tmp_assignment(command: str) -> bool:
    try:
        tokens = shell_tokens(command)
    except ValueError:
        return False
    return any(segment_contains_windows_tmp_assignment(segment) for segment in split_segments(tokens))


def segment_contains_windows_tmp_assignment(segment: list[str]) -> bool:
    index = 0
    while index < len(segment) and is_env_assignment(segment[index]):
        if is_windows_tmp_assignment_token(segment[index]):
            return True
        index += 1

    if index >= len(segment):
        return False

    executable = executable_name(segment[index])
    if executable in {"export", "declare", "typeset", "local", "readonly"}:
        return any(
            is_windows_tmp_assignment_token(token)
            for token in assignment_builtin_operands(segment[index + 1 :])
        )
    if executable == "env":
        split_command = env_split_string(segment[index:])
        if split_command and contains_windows_tmp_assignment(split_command):
            return True
        index += 1
        while index < len(segment) and segment[index].startswith("-"):
            option = segment[index]
            index += 1
            if option in {"-u", "--unset", "-C", "--chdir", "-S", "--split-string"} and index < len(segment):
                index += 1
        while index < len(segment) and is_env_assignment(segment[index]):
            if is_windows_tmp_assignment_token(segment[index]):
                return True
            index += 1
    return False


def env_split_string(segment: list[str]) -> str | None:
    if not segment or executable_name(segment[0]) != "env":
        return None
    index = 1
    while index < len(segment):
        token = segment[index]
        if token == "--":
            return None
        if token in {"-S", "--split-string"}:
            return normalize_env_split_string(segment[index + 1] if index + 1 < len(segment) else "")
        if token.startswith("-S") and len(token) > 2:
            return normalize_env_split_string(token[2:])
        if token.startswith("-") and not token.startswith("--") and "S" in token[1:]:
            after = token.split("S", 1)[1]
            return normalize_env_split_string(after if after else (segment[index + 1] if index + 1 < len(segment) else ""))
        prefix = "--split-string="
        if token.startswith(prefix):
            return normalize_env_split_string(token[len(prefix):])
        if token in {"-u", "--unset", "-C", "--chdir"}:
            index += 2
            continue
        if token.startswith("-"):
            index += 1
            continue
        return None
    return None


def normalize_env_split_string(value: str) -> str:
    return value.replace("\\_", " ")


def assignment_builtin_operands(tokens: list[str]) -> list[str]:
    operands: list[str] = []
    index = 0
    while index < len(tokens):
        token = tokens[index]
        if token == "--":
            operands.extend(tokens[index + 1 :])
            break
        if token.startswith("-"):
            index += 1
            continue
        operands.append(token)
        index += 1
    return operands


def is_windows_tmp_assignment_token(token: str) -> bool:
    for name in ("TMPDIR", "TEMP", "TMP"):
        prefix = f"{name}="
        if not token.startswith(prefix):
            continue
        value = token[len(prefix):]
        if value.startswith("$/"):
            value = value[1:]
        elif value.startswith("$\\"):
            value = decode_ansi_c_fragment(value[1:])
        return bool(MNT_PATH.match(value))
    return False


def decode_ansi_c_fragment(value: str) -> str:
    decoded: list[str] = []
    index = 0
    while index < len(value):
        if value[index] == "\\":
            char, index = read_ansi_c_escape(value, index)
            decoded.append(char)
            continue
        decoded.append(value[index])
        index += 1
    return "".join(decoded)


def shell_expansion_commands(text: str) -> list[str]:
    text = normalize_shell_line_continuations(text)
    commands: list[str] = []
    index = 0
    while index < len(text):
        char = text[index]
        if char == "\\":
            index += 2
            continue
        if text.startswith("$((", index):
            index += 3
            continue
        if text.startswith("$(", index):
            command, end = read_parenthesized_command(text, index + 2)
            if command is not None:
                commands.append(command)
                index = end
                continue
        if text.startswith("${", index):
            command, end = read_braced_command(text, index + 2)
            if command is not None:
                commands.append(command)
                index = end
                continue
        if char == "`":
            command, end = read_backtick_command(text, index + 1)
            if command is not None:
                commands.append(command)
                index = end
                continue
        index += 1
    return commands


def normalize_shell_line_continuations(text: str) -> str:
    return re.sub(r"\\\r?\n", "", text)


def read_parenthesized_command(text: str, index: int) -> tuple[str | None, int]:
    start = index
    depth = 1
    quote: str | None = None
    while index < len(text):
        char = text[index]
        if quote is None:
            if char == "\\":
                index += 2
                continue
            if char in {"'", '"'}:
                quote = char
                index += 1
                continue
            if char == "(":
                depth += 1
            elif char == ")":
                depth -= 1
                if depth == 0:
                    return text[start:index], index + 1
            index += 1
            continue

        if char == quote:
            quote = None
        elif char == "\\":
            index += 1
        index += 1
    return None, len(text)


def read_braced_command(text: str, index: int) -> tuple[str | None, int]:
    start = index
    quote: str | None = None
    while index < len(text):
        char = text[index]
        if quote is None:
            if char == "\\":
                index += 2
                continue
            if char in {"'", '"'}:
                quote = char
                index += 1
                continue
            if char == "}":
                return text[start:index], index + 1
            index += 1
            continue

        if char == quote:
            quote = None
        elif char == "\\":
            index += 1
        index += 1
    return None, len(text)


def read_backtick_command(text: str, index: int) -> tuple[str | None, int]:
    start = index
    while index < len(text):
        char = text[index]
        if char == "\\":
            index += 2
            continue
        if char == "`":
            return text[start:index], index + 1
        index += 1
    return None, len(text)


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


def segment_shell_sink(segment: list[str], specs: list[HeredocSpec] | None = None) -> bool:
    if segment_executable(segment) in SHELLS:
        return True
    without_redirections = segment_without_heredoc_redirections(segment, specs)
    if segment_executable(without_redirections) in SHELLS:
        return True
    split_command = env_split_string(without_redirections)
    if split_command:
        try:
            split_tokens = shell_tokens(split_command)
        except ValueError:
            split_tokens = []
        if split_tokens and segment_shell_sink(split_tokens):
            return True
    if leading_heredoc_redirected_executable(segment, specs) in SHELLS:
        return True
    return False


def segment_without_heredoc_redirections(
    segment: list[str], specs: list[HeredocSpec] | None = None
) -> list[str]:
    output: list[str] = []
    index = 0
    spec_index = 0
    while index < len(segment):
        token = segment[index]
        if token.isdigit() and index + 1 < len(segment) and segment[index + 1] == "<<":
            if specs is not None and (
                spec_index >= len(specs) or specs[spec_index].fd != token
            ):
                output.append(token)
                index += 1
                continue
            index += 3
            spec_index += 1
            continue
        if token == "<<":
            index += 2
            spec_index += 1
            continue
        output.append(token)
        index += 1
    return output


def leading_heredoc_redirected_executable(
    segment: list[str], specs: list[HeredocSpec] | None = None
) -> str:
    args = segment_without_heredoc_redirections(segment, specs)
    if args:
        return segment_executable(args)
    args = list(segment)
    index = 0
    saw_heredoc = False

    while index < len(args):
        token = args[index]
        if token.isdigit() and index + 1 < len(args) and args[index + 1] == "<<":
            index += 1
            token = args[index]
        if token != "<<":
            break
        saw_heredoc = True
        index += 1
        if index < len(args) and args[index] == "-":
            index += 1
        if index < len(args):
            index += 1

    if not saw_heredoc or index >= len(args):
        return ""
    return segment_executable(args[index:])


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


def is_raw_windows_path_token(token: str) -> bool:
    return bool(re.match(r"^[A-Za-z]:[\\/]", token))


def raw_executable_name(token: str) -> str:
    return re.split(r"[\\/]", token)[-1].lower()


def is_windows_git_executable(token: str) -> bool:
    raw_name = raw_executable_name(token)
    return raw_name == "git.exe" or (
        raw_name == "git" and (is_raw_windows_path_token(token) or bool(MNT_PATH.match(token)))
    )


def is_windows_executable_path(token: str) -> bool:
    return raw_executable_name(token).endswith(".exe") and (
        is_raw_windows_path_token(token) or bool(MNT_PATH.match(token))
    )


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


def raw_windows_path_to_wsl(path: str) -> str:
    drive = path[0].lower()
    rest = re.sub(r"[\\/]+", "/", path[2:]).lstrip("/")
    return f"/mnt/{drive}/{rest}" if rest else f"/mnt/{drive}"


def rewrite_raw_windows_paths(segment: str) -> str:
    def replace_quoted(match: re.Match[str]) -> str:
        prefix = match.group(1)
        quote = match.group(2)
        path = match.group(3)
        return f"{prefix}{quote}{raw_windows_path_to_wsl(path)}{quote}"

    def replace(match: re.Match[str]) -> str:
        prefix = match.group(1)
        path = match.group(2)
        return f"{prefix}{raw_windows_path_to_wsl(path)}"

    return WINDOWS_PATH.sub(replace, QUOTED_WINDOWS_PATH.sub(replace_quoted, segment))


def rewrite_windows_paths_for_wsl(command: str) -> str:
    return transform_without_heredoc_bodies(command, rewrite_windows_paths_for_wsl_line)


def rewrite_windows_paths_for_wsl_line(line: str) -> str:
    parts = re.split(r"([;&|()\n]+)", line)
    rewritten_parts: list[str] = []
    for part in parts:
        if not part or re.fullmatch(r"[;&|()\n]+", part):
            rewritten_parts.append(part)
            continue
        if not WINDOWS_PATH.search(part) and not QUOTED_WINDOWS_PATH.search(part):
            rewritten_parts.append(part)
            continue
        try:
            executable = segment_executable(shell_tokens(part))
        except ValueError:
            rewritten_parts.append(part)
            continue
        if executable in WINDOWS_PATH_BRIDGE_ALLOWLIST | SAFE_TEXT_COMMANDS:
            rewritten_parts.append(part)
            continue
        if executable in CODE_STRING_COMMANDS and has_code_string_option(part):
            rewritten_parts.append(part)
            continue
        if executable in SEARCH_TEXT_COMMANDS:
            rewritten_parts.append(rewrite_search_windows_paths(part))
            continue
        rewritten_parts.append(rewrite_raw_windows_paths(part))
    return "".join(rewritten_parts)


def has_code_string_option(part: str) -> bool:
    try:
        tokens = shell_tokens(part)
    except ValueError:
        return False
    return any(token in {"-c", "-e"} or (token.startswith("-") and any(flag in token[1:] for flag in {"c", "e"})) for token in tokens[1:])


def rewrite_search_windows_paths(part: str) -> str:
    words = shell_word_spans(part)
    pattern_index = search_pattern_word_index([word for word, _, _ in words])
    if pattern_index is None:
        return rewrite_raw_windows_paths(part)
    start, end = words[pattern_index][1], words[pattern_index][2]
    return rewrite_raw_windows_paths(part[:start]) + part[start:end] + rewrite_raw_windows_paths(part[end:])


def search_pattern_word_index(tokens: list[str]) -> int | None:
    index = 1
    while index < len(tokens):
        token = tokens[index]
        if token == "--":
            return index + 1 if index + 1 < len(tokens) else None
        if token in {"-e", "--regexp"}:
            return index + 1 if index + 1 < len(tokens) else None
        if token.startswith("--regexp="):
            return index
        if token.startswith("-"):
            index += 2 if token in {"-e", "-f", "--regexp", "--file", "-g", "--glob", "-t", "--type"} else 1
            continue
        return index
    return None


def shell_word_spans(command: str) -> list[tuple[str, int, int]]:
    words: list[tuple[str, int, int]] = []
    index = 0
    while index < len(command):
        while index < len(command) and command[index].isspace():
            index += 1
        if index >= len(command):
            break
        start = index
        value: list[str] = []
        quote: str | None = None
        while index < len(command):
            char = command[index]
            if quote is None and char.isspace():
                break
            if quote is None and char in {"'", '"'}:
                quote = char
                index += 1
                continue
            if quote is not None and char == quote:
                quote = None
                index += 1
                continue
            if char == "\\" and index + 1 < len(command):
                value.append(command[index + 1])
                index += 2
                continue
            value.append(char)
            index += 1
        words.append(("".join(value), start, index))
    return words


def rewritten_shell_input(tool_input: dict) -> dict | None:
    command_key = "command" if "command" in tool_input else "cmd"
    command = tool_input.get(command_key)
    if not isinstance(command, str) or not command.strip():
        return None

    rewritten = rewrite_windows_paths_for_wsl(command)
    if rewritten == command:
        return None

    updated = dict(tool_input)
    updated[command_key] = rewritten
    return updated


def rewritten_parallel_input(tool_input: dict) -> dict | None:
    tool_uses = tool_input.get("tool_uses")
    if not isinstance(tool_uses, list):
        return None

    changed = False
    updated_tool_uses: list = []
    for tool_use in tool_uses:
        if not isinstance(tool_use, dict):
            updated_tool_uses.append(tool_use)
            continue
        nested_name = tool_use.get("recipient_name") or tool_use.get("tool_name")
        parameters = tool_use.get("parameters") or {}
        if nested_name not in SHELL_TOOL_NAMES or not isinstance(parameters, dict):
            updated_tool_uses.append(tool_use)
            continue
        updated_parameters = rewritten_shell_input(parameters)
        if updated_parameters is None:
            updated_tool_uses.append(tool_use)
            continue
        updated_tool_use = dict(tool_use)
        updated_tool_use["parameters"] = updated_parameters
        updated_tool_uses.append(updated_tool_use)
        changed = True

    if not changed:
        return None

    updated_input = dict(tool_input)
    updated_input["tool_uses"] = updated_tool_uses
    return updated_input


def rewritten_tool_input(payload: dict) -> dict | None:
    tool_input = payload.get("tool_input") or {}
    if not isinstance(tool_input, dict):
        return None

    tool_name = payload.get("tool_name")
    if tool_name in SHELL_TOOL_NAMES:
        return rewritten_shell_input(tool_input)
    if tool_name == "multi_tool_use.parallel":
        return rewritten_parallel_input(tool_input)
    return None


def blocked_segment_reason(segment: list[str], depth: int) -> str | None:
    split_command = env_split_string(segment)
    if split_command:
        reason = blocked_reason(split_command, depth + 1)
        if reason:
            return reason

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

    if executable_token and is_windows_executable_path(executable_token):
        return (
            "Blocked: do not run Windows executable paths from Codex shell commands. "
            "Use Tabby MCP for Windows-side executables. "
            f"{TABBY_NUDGE}"
        )

    return None


def blocked_reason(command: str, depth: int = 0) -> str | None:
    if depth > 5:
        return f"Blocked: deeply nested shell command could not be safely inspected. {TABBY_NUDGE}"

    for body in heredoc_bodies(command):
        if body.executable_shell:
            reason = blocked_reason(body.text, depth + 1)
            if reason:
                return reason
            continue
        if body.expands:
            for nested_command in shell_expansion_commands(body.text):
                reason = blocked_reason(nested_command, depth + 1)
                if reason:
                    return reason

    command = command_without_heredoc_bodies(command)
    if contains_windows_tmp_assignment(command):
        return "Blocked: keep TMPDIR, TEMP, and TMP on Linux storage such as /tmp, not /mnt/*."

    try:
        tokens = shell_tokens(command)
    except ValueError:
        tokens = []

    for segment in split_segments(tokens):
        reason = blocked_segment_reason(segment, depth)
        if reason:
            return reason
        if contains_ln_s_on_windows_mount(segment):
            return (
                "Blocked: create symlinks on Windows drives through Tabby MCP using "
                "Windows-native mklink or New-Item in an existing Windows-side session. "
                f"{TABBY_NUDGE}"
            )

    if "\n" in command:
        for line in command.splitlines():
            if not line.strip():
                continue
            reason = blocked_reason(line, depth + 1)
            if reason:
                return reason
        return None

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


def allow_rewrite(updated_input: dict) -> None:
    print(
        json.dumps(
            {
                "hookSpecificOutput": {
                    "hookEventName": "PreToolUse",
                    "permissionDecision": "allow",
                    "updatedInput": updated_input,
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

    updated_input = rewritten_tool_input(payload)
    if updated_input is not None:
        updated_payload = dict(payload)
        updated_payload["tool_input"] = updated_input
        for command in shell_commands(updated_payload):
            reason = blocked_reason(command)
            if reason:
                deny(reason)
                return 0
        allow_rewrite(updated_input)
        return 0

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
