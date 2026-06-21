#!/usr/bin/env python3
"""Regression tests for Agent Loop Designer loop specs."""

from __future__ import annotations

import unittest

import loop_spec


SpecValue = str | list[str] | dict[str, str] | list[dict[str, str]]
Spec = dict[str, SpecValue]


def base_spec() -> Spec:
    return {
        "mission": "test loop",
        "command": "$test-loop",
        "surface": "skill",
        "trigger": "explicit test trigger",
        "state": ["repo checkout"],
        "inputs": ["plain input"],
        "workers": [],
        "artifact": "report",
        "decision_point": "stop when valid",
        "next_action": "report",
        "stop_condition": "validator catches contract drift",
        "loop_budget": "one pass",
        "human_gate": "none",
        "safety_rule": "read-only unless editing owned files",
        "failure_learning": {
            "trigger": "contract drift",
            "evidence": "validator output",
            "update_target": "loop_spec.py",
            "validation": "smoke tests",
            "skip_when": "none",
        },
        "tools": [],
        "docs_checked": [],
        "subagents": [],
        "worktree_threads": [],
    }


def worktree(coordination_backend: str) -> dict[str, str]:
    return {
        "role": "candidate",
        "ownership": "owned files",
        "write_intent": "code-editing",
        "starting_state": "verified HEAD",
        "project_scope": "current saved project",
        "location_strategy": "wsl-same-project-manual",
        "worker_cwd": ".codex-worktrees/run/candidate",
        "output_contract": "summary and tests",
        "integration_rule": "report only",
        "coordination_backend": coordination_backend,
        "visibility": "human-readable thread that can receive intervention",
    }


def grit_coordination(backend: str = "local") -> dict[str, str]:
    return {
        "backend": backend,
        "init_policy": "initialize grit state from approved settings",
        "claim_strategy": "claim explicit symbols before editing",
        "done_policy": "run grit done only when integration is authorized",
        "thread_context_rule": "ask worker threads for compact missing context",
        "cleanup_rule": "release only loop-owned locks",
    }


class LoopSpecValidationTest(unittest.TestCase):
    def assert_valid(self, spec: Spec) -> None:
        self.assertEqual(loop_spec.validate_spec(spec), [])

    def assert_error(self, spec: Spec, expected: str) -> None:
        errors = loop_spec.validate_spec(spec)
        self.assertIn(expected, errors)

    def test_active_sub_agent_requires_structured_subagents(self) -> None:
        spec = base_spec()
        spec["workers"] = ["one sub-agent reviewer"]

        self.assert_error(
            spec,
            "workers mention subagents/custom agents; add structured subagents",
        )

    def test_negated_sub_agent_does_not_require_structured_subagents(self) -> None:
        spec = base_spec()
        spec["workers"] = ["main agent only; no sub-agent delegation"]

        self.assert_valid(spec)

    def test_grit_backend_requires_grit_coordination(self) -> None:
        spec = base_spec()
        spec["workers"] = ["one Worktree thread edits code with Grit"]
        spec["docs_checked"] = ["loop-spec.md"]
        spec["worktree_threads"] = [worktree("grit")]

        self.assert_error(spec, "grit_coordination is required when the loop uses Grit")

    def test_negated_grit_does_not_require_grit_coordination(self) -> None:
        spec = base_spec()
        spec["workers"] = ["manual Worktree thread only; no Grit coordination"]
        spec["docs_checked"] = ["loop-spec.md"]
        spec["worktree_threads"] = [worktree("manual-worktree")]

        self.assert_valid(spec)

    def test_worktree_backed_thread_requires_structured_worktree_threads(self) -> None:
        spec = base_spec()
        spec["workers"] = ["worktree-backed candidate thread edits code"]

        self.assert_error(
            spec,
            "workers mention worktree threads; add structured worktree_threads",
        )

    def test_grit_done_guard_requires_grit_coordination(self) -> None:
        spec = base_spec()
        spec["workers"] = ["do not run grit done unless integration is authorized"]

        self.assert_error(spec, "grit_coordination is required when the loop uses Grit")

    def test_grit_claim_guard_requires_grit_coordination(self) -> None:
        spec = base_spec()
        spec["workers"] = ["worker must not edit outside Grit claims"]

        self.assert_error(spec, "grit_coordination is required when the loop uses Grit")

    def test_remote_grit_backend_requires_authorization(self) -> None:
        spec = base_spec()
        spec["workers"] = ["one Worktree thread edits code with Grit"]
        spec["docs_checked"] = ["loop-spec.md"]
        spec["worktree_threads"] = [worktree("grit")]
        spec["grit_coordination"] = grit_coordination("azure")

        self.assert_error(
            spec,
            "grit_coordination.remote_backend_authorization is required when backend is not local",
        )

    def test_authorized_remote_grit_backend_validates_and_renders(self) -> None:
        spec = base_spec()
        spec["workers"] = ["one Worktree thread edits code with Grit"]
        spec["docs_checked"] = ["loop-spec.md"]
        spec["worktree_threads"] = [worktree("grit")]
        coordination = grit_coordination("azure")
        coordination["remote_backend_authorization"] = (
            "User explicitly requested distributed team coordination; "
            "non-secret settings are in repo docs."
        )
        spec["grit_coordination"] = coordination

        self.assert_valid(spec)
        rendered = loop_spec.render_markdown(spec)
        self.assertIn("## Grit Coordination", rendered)
        self.assertIn("- Remote backend authorization: User explicitly requested", rendered)


if __name__ == "__main__":
    unittest.main()
