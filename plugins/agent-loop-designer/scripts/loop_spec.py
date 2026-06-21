#!/usr/bin/env python3
"""Create, validate, and render Agent Loop Designer specs."""

# allow: SIZE_OK - standalone plugin helper script shipped as one file.

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any


REQUIRED_FIELDS = [
    "mission",
    "command",
    "surface",
    "trigger",
    "state",
    "inputs",
    "workers",
    "artifact",
    "decision_point",
    "next_action",
    "stop_condition",
    "loop_budget",
    "human_gate",
    "safety_rule",
    "failure_learning",
]

VALID_SURFACES = {
    "reusable-prompt",
    "skill",
    "automation-prompt",
    "worktree-policy",
    "plugin-backed-skill",
}

DOCS_REQUIRED_SURFACES = {"automation-prompt", "worktree-policy", "plugin-backed-skill"}
VALID_WRITE_INTENTS = {"none", "artifact-only", "code-editing"}
VALID_SANDBOXES = {"read-only", "workspace-write", "permission-profile", "danger-full-access"}
VALID_WORKTREE_LOCATION_STRATEGIES = {
    "wsl-manual",
    "wsl-same-project-manual",
    "wsl-external-manual",
    "codex-managed",
}
VALID_WORKTREE_COORDINATION_BACKENDS = {"manual-worktree", "grit"}

REQUIRED_SUBAGENT_FIELDS = [
    "role",
    "spawn_policy",
    "ownership",
    "write_intent",
    "sandbox_expectation",
    "approval_expectation",
    "output_contract",
    "coordination_rule",
    "context_budget",
]

REQUIRED_WORKTREE_THREAD_FIELDS = [
    "role",
    "ownership",
    "write_intent",
    "starting_state",
    "project_scope",
    "location_strategy",
    "worker_cwd",
    "output_contract",
    "integration_rule",
    "coordination_backend",
    "visibility",
]

REQUIRED_GRIT_COORDINATION_FIELDS = [
    "backend",
    "init_policy",
    "claim_strategy",
    "done_policy",
    "thread_context_rule",
    "cleanup_rule",
]

NEGATED_MENTION_PREFIX_RE = re.compile(
    r"(?:^|\s)(?:no|without|never|avoid|skip|exclude|disable|not|don't|dont|"
    r"do not|must not|cannot|can't|cant)\s+"
    r"(?:[\w-]+\s+){0,4}$"
)


def empty_spec(task: str) -> dict[str, Any]:
    return {
        "mission": task,
        "command": "",
        "surface": "",
        "trigger": "",
        "state": [],
        "inputs": [],
        "workers": [],
        "artifact": "",
        "decision_point": "",
        "next_action": "",
        "stop_condition": "",
        "loop_budget": "",
        "human_gate": "",
        "safety_rule": "",
        "failure_learning": {
            "trigger": "",
            "evidence": "",
            "update_target": "",
            "validation": "",
            "skip_when": "",
        },
        "tools": [],
        "docs_checked": [],
        "worktree_threads": [],
        "subagents": [],
    }


def load_spec(path: Path) -> dict[str, Any]:
    with path.open(encoding="utf-8") as handle:
        payload = json.load(handle)
    if not isinstance(payload, dict):
        raise ValueError("spec must be a JSON object")
    return payload


def non_empty(value: Any) -> bool:
    if isinstance(value, str):
        return bool(value.strip())
    if isinstance(value, list):
        return bool(value)
    if isinstance(value, dict):
        return bool(value)
    return value is not None


def list_mentions(values: Any, needles: tuple[str, ...]) -> bool:
    if not isinstance(values, list):
        return False
    return any(
        isinstance(value, str)
        and any(text_has_active_mention(value, needle) for needle in needles)
        for value in values
    )


def list_mentions_grit(values: Any) -> bool:
    if not isinstance(values, list):
        return False
    return any(
        isinstance(value, str) and text_has_grit_contract(value)
        for value in values
    )


def mention_pattern(needle: str) -> re.Pattern[str]:
    words = [word for word in re.split(r"[-\s]+", needle.lower()) if word]
    body = r"[-\s]+".join(re.escape(word) for word in words)
    return re.compile(rf"\b{body}s?\b")


def text_has_active_mention(text: str, needle: str) -> bool:
    normalized = text.lower()
    pattern = mention_pattern(needle)
    return any(
        not mention_is_negated(normalized, match.start())
        for match in pattern.finditer(normalized)
    )


def text_has_grit_contract(text: str) -> bool:
    normalized = text.lower()
    if re.search(r"\bgrit[-\s]+(?:assign|claim|claims|done|init|status|worktree|worktrees)\b", normalized):
        return True
    return text_has_active_mention(text, "grit")


def mention_is_negated(text: str, start: int) -> bool:
    clause_prefix = re.split(r"[.;:]", text[:start])[-1]
    return bool(NEGATED_MENTION_PREFIX_RE.search(clause_prefix))


def spawn_policy_encodes_safe_fork(policy: str) -> bool:
    normalized = policy.lower()
    uses_fork_turns_none = bool(
        re.search(r"\bfork_turns\b\s*[:=]\s*[\"']?none\b", normalized)
    )
    uses_fork_context_false = bool(
        re.search(r"\bfork_context\b\s*[:=]\s*false\b", normalized)
    )
    omits_fork_context = bool(
        re.search(r"\bomit(?:ted)?\s+\bfork_context\b", normalized)
    )
    intentional_full_history = bool(
        re.search(r"\bfull[- ]history\b", normalized)
        and re.search(r"\binherit(?:s|ing)?\b", normalized)
        and re.search(r"\b(no|without|must not)\s+overrides?\b", normalized)
    )
    return (
        uses_fork_turns_none
        or uses_fork_context_false
        or omits_fork_context
        or intentional_full_history
    )


def validate_subagent(index: int, subagent: Any) -> list[str]:
    errors: list[str] = []
    prefix = f"subagents[{index}]"
    if not isinstance(subagent, dict):
        return [f"{prefix} must be an object"]

    for field in REQUIRED_SUBAGENT_FIELDS:
        if not non_empty(subagent.get(field)):
            errors.append(f"{prefix}.{field} is required")

    write_intent = subagent.get("write_intent")
    if write_intent and write_intent not in VALID_WRITE_INTENTS:
        errors.append(
            f"{prefix}.write_intent must be one of {', '.join(sorted(VALID_WRITE_INTENTS))}"
        )

    sandbox = subagent.get("sandbox_expectation")
    if sandbox and sandbox not in VALID_SANDBOXES:
        errors.append(
            f"{prefix}.sandbox_expectation must be one of {', '.join(sorted(VALID_SANDBOXES))}"
        )

    spawn_policy = str(subagent.get("spawn_policy", ""))
    if spawn_policy and not spawn_policy_encodes_safe_fork(spawn_policy):
        errors.append(
            f"{prefix}.spawn_policy must encode fork_context false/omitted, fork_turns none, or intentional full-history inheritance with no overrides"
        )

    if write_intent == "code-editing" and sandbox == "read-only":
        errors.append(f"{prefix} cannot do code-editing with read-only sandbox_expectation")

    if write_intent == "code-editing" and not str(subagent.get("ownership", "")).strip():
        errors.append(f"{prefix} needs explicit ownership for code-editing")

    return errors


def validate_worktree_thread(index: int, thread: Any) -> list[str]:
    errors: list[str] = []
    prefix = f"worktree_threads[{index}]"
    if not isinstance(thread, dict):
        return [f"{prefix} must be an object"]

    for field in REQUIRED_WORKTREE_THREAD_FIELDS:
        if not non_empty(thread.get(field)):
            errors.append(f"{prefix}.{field} is required")

    write_intent = thread.get("write_intent")
    if write_intent and write_intent not in VALID_WRITE_INTENTS:
        errors.append(
            f"{prefix}.write_intent must be one of {', '.join(sorted(VALID_WRITE_INTENTS))}"
        )

    location_strategy = thread.get("location_strategy")
    if location_strategy and location_strategy not in VALID_WORKTREE_LOCATION_STRATEGIES:
        errors.append(
            f"{prefix}.location_strategy must be one of {', '.join(sorted(VALID_WORKTREE_LOCATION_STRATEGIES))}"
        )

    coordination_backend = thread.get("coordination_backend")
    if coordination_backend and coordination_backend not in VALID_WORKTREE_COORDINATION_BACKENDS:
        errors.append(
            f"{prefix}.coordination_backend must be one of {', '.join(sorted(VALID_WORKTREE_COORDINATION_BACKENDS))}"
        )

    return errors


def validate_failure_learning(value: Any) -> list[str]:
    if not isinstance(value, dict):
        return ["failure_learning must be an object"]

    errors: list[str] = []
    for field in ["trigger", "evidence", "update_target", "validation", "skip_when"]:
        if not non_empty(value.get(field)):
            errors.append(f"failure_learning.{field} is required")
    return errors


def validate_grit_coordination(value: Any) -> list[str]:
    if not isinstance(value, dict):
        return ["grit_coordination must be an object"]

    errors: list[str] = []
    for field in REQUIRED_GRIT_COORDINATION_FIELDS:
        if not non_empty(value.get(field)):
            errors.append(f"grit_coordination.{field} is required")

    backend = value.get("backend")
    if (
        backend
        and backend != "local"
        and not non_empty(value.get("remote_backend_authorization"))
    ):
        errors.append(
            "grit_coordination.remote_backend_authorization is required when backend is not local"
        )

    return errors


def validate_spec(spec: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    for field in REQUIRED_FIELDS:
        if not non_empty(spec.get(field)):
            errors.append(f"{field} is required")

    command = spec.get("command")
    if isinstance(command, str) and command.strip() and not command.strip().startswith("$"):
        errors.append("command should start with `$`")

    surface = spec.get("surface")
    if surface and surface not in VALID_SURFACES:
        errors.append(f"surface must be one of {', '.join(sorted(VALID_SURFACES))}")

    for list_field in ["state", "inputs", "workers", "tools", "docs_checked"]:
        if list_field in spec and not isinstance(spec[list_field], list):
            errors.append(f"{list_field} must be an array")

    tools = spec.get("tools", [])
    docs_checked = spec.get("docs_checked", [])
    if isinstance(tools, list) and tools and not non_empty(docs_checked):
        errors.append("docs_checked is required when tools are specified")

    docs_required_reasons: list[str] = []
    if surface in DOCS_REQUIRED_SURFACES:
        docs_required_reasons.append(f"surface is {surface}")

    workers = spec.get("workers", [])
    workers_mention_worktree_thread = list_mentions(
        workers,
        (
            "worktree thread",
            "worktree threads",
            "worktree-thread",
            "worktree-backed",
            "codex worktree",
        ),
    )
    if workers_mention_worktree_thread:
        docs_required_reasons.append("workers mention worktree threads")
        if not spec.get("worktree_threads"):
            errors.append("workers mention worktree threads; add structured worktree_threads")

    workers_mention_subagent = list_mentions(
        workers, ("subagent", "sub-agent", "sub agent", "custom agent")
    )
    if workers_mention_subagent:
        docs_required_reasons.append("workers mention subagents/custom agents")
        if not spec.get("subagents"):
            errors.append("workers mention subagents/custom agents; add structured subagents")

    if list_mentions(
        spec.get("inputs", []),
        ("cli", "script", "tool", "mcp", "connector", "plugin", "automation", "config"),
    ):
        docs_required_reasons.append("inputs mention docs-first surfaces")

    if docs_required_reasons and not non_empty(docs_checked):
        errors.append(
            "docs_checked is required when " + "; ".join(docs_required_reasons)
        )

    if "failure_learning" in spec:
        errors.extend(validate_failure_learning(spec.get("failure_learning")))

    uses_grit_coordination = False

    worktree_threads = spec.get("worktree_threads", [])
    if worktree_threads is None:
        worktree_threads = []
    if not isinstance(worktree_threads, list):
        errors.append("worktree_threads must be an array")
    else:
        if worktree_threads and not non_empty(docs_checked):
            errors.append("docs_checked is required when worktree_threads are specified")
        for index, thread in enumerate(worktree_threads):
            errors.extend(validate_worktree_thread(index, thread))
            if isinstance(thread, dict) and thread.get("coordination_backend") == "grit":
                uses_grit_coordination = True

    subagents = spec.get("subagents", [])
    if subagents is None:
        subagents = []
    if not isinstance(subagents, list):
        errors.append("subagents must be an array")
    else:
        if subagents and not non_empty(docs_checked):
            errors.append("docs_checked is required when subagents are specified")
        for index, subagent in enumerate(subagents):
            errors.extend(validate_subagent(index, subagent))

    if list_mentions_grit(spec.get("tools", [])) or list_mentions_grit(workers):
        uses_grit_coordination = True

    if "grit_coordination" in spec:
        uses_grit_coordination = True

    if uses_grit_coordination and not non_empty(spec.get("grit_coordination")):
        errors.append("grit_coordination is required when the loop uses Grit")
    elif uses_grit_coordination:
        errors.extend(validate_grit_coordination(spec.get("grit_coordination")))

    return errors


def render_markdown(spec: dict[str, Any]) -> str:
    lines = [
        f"# {spec.get('command', 'Loop Spec')}",
        "",
        f"Mission: {spec.get('mission', '')}",
        "",
        f"- Surface: {spec.get('surface', '')}",
        f"- Trigger: {spec.get('trigger', '')}",
        f"- State: {', '.join(spec.get('state', []))}",
        f"- Inputs: {', '.join(spec.get('inputs', []))}",
        f"- Workers: {', '.join(spec.get('workers', []))}",
        f"- Tools: {', '.join(spec.get('tools', []))}",
        f"- Docs checked: {', '.join(spec.get('docs_checked', []))}",
        f"- Artifact: {spec.get('artifact', '')}",
        f"- Decision point: {spec.get('decision_point', '')}",
        f"- Next action: {spec.get('next_action', '')}",
        f"- Stop condition: {spec.get('stop_condition', '')}",
        f"- Loop budget: {spec.get('loop_budget', '')}",
        f"- Human gate: {spec.get('human_gate', '')}",
        f"- Safety rule: {spec.get('safety_rule', '')}",
    ]
    failure_learning = spec.get("failure_learning") or {}
    if failure_learning:
        lines.extend(
            [
                f"- Failure trigger: {failure_learning.get('trigger', '')}",
                f"- Failure evidence: {failure_learning.get('evidence', '')}",
                f"- Failure update target: {failure_learning.get('update_target', '')}",
                f"- Failure validation: {failure_learning.get('validation', '')}",
                f"- Failure skip rule: {failure_learning.get('skip_when', '')}",
            ]
        )
    subagents = spec.get("subagents") or []
    if subagents:
        lines.extend(["", "## Subagents"])
        for subagent in subagents:
            lines.extend(
                [
                    "",
                    f"### {subagent.get('role', '')}",
                    f"- Spawn policy: {subagent.get('spawn_policy', '')}",
                    f"- Ownership: {subagent.get('ownership', '')}",
                    f"- Write intent: {subagent.get('write_intent', '')}",
                    f"- Sandbox: {subagent.get('sandbox_expectation', '')}",
                    f"- Approval: {subagent.get('approval_expectation', '')}",
                    f"- Output: {subagent.get('output_contract', '')}",
                    f"- Coordination: {subagent.get('coordination_rule', '')}",
                    f"- Context budget: {subagent.get('context_budget', '')}",
                ]
            )
    worktree_threads = spec.get("worktree_threads") or []
    if worktree_threads:
        lines.extend(["", "## Worktree Threads"])
        for thread in worktree_threads:
            lines.extend(
                [
                    "",
                    f"### {thread.get('role', '')}",
                    f"- Ownership: {thread.get('ownership', '')}",
                    f"- Write intent: {thread.get('write_intent', '')}",
                    f"- Starting state: {thread.get('starting_state', '')}",
                    f"- Project scope: {thread.get('project_scope', '')}",
                    f"- Location strategy: {thread.get('location_strategy', '')}",
                    f"- Worker cwd: {thread.get('worker_cwd', '')}",
                    f"- Output: {thread.get('output_contract', '')}",
                    f"- Integration: {thread.get('integration_rule', '')}",
                    f"- Coordination backend: {thread.get('coordination_backend', '')}",
                    f"- Visibility: {thread.get('visibility', '')}",
                ]
            )
    grit_coordination = spec.get("grit_coordination") or {}
    if grit_coordination:
        lines.extend(
            [
                "",
                "## Grit Coordination",
                f"- Backend: {grit_coordination.get('backend', '')}",
                f"- Init policy: {grit_coordination.get('init_policy', '')}",
                f"- Claim strategy: {grit_coordination.get('claim_strategy', '')}",
                f"- Done policy: {grit_coordination.get('done_policy', '')}",
                f"- Thread context rule: {grit_coordination.get('thread_context_rule', '')}",
                f"- Cleanup rule: {grit_coordination.get('cleanup_rule', '')}",
                "- Remote backend authorization: "
                f"{grit_coordination.get('remote_backend_authorization', '')}",
            ]
        )
    return "\n".join(lines) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    subparsers = parser.add_subparsers(dest="command", required=True)

    template_parser = subparsers.add_parser("template", help="Print a blank loop spec JSON")
    template_parser.add_argument("--task", default="", help="Plain-language recurring task")

    validate_parser = subparsers.add_parser("validate", help="Validate a loop spec JSON file")
    validate_parser.add_argument("path", type=Path)

    render_parser = subparsers.add_parser("render", help="Render a loop spec JSON file as Markdown")
    render_parser.add_argument("path", type=Path)

    args = parser.parse_args()

    if args.command == "template":
        print(json.dumps(empty_spec(args.task), indent=2))
        return 0

    spec = load_spec(args.path)

    if args.command == "validate":
        errors = validate_spec(spec)
        if errors:
            for error in errors:
                print(f"ERROR: {error}", file=sys.stderr)
            return 1
        print("Loop spec is valid.")
        return 0

    if args.command == "render":
        errors = validate_spec(spec)
        if errors:
            for error in errors:
                print(f"ERROR: {error}", file=sys.stderr)
            return 1
        print(render_markdown(spec), end="")
        return 0

    raise AssertionError(args.command)


if __name__ == "__main__":
    raise SystemExit(main())
