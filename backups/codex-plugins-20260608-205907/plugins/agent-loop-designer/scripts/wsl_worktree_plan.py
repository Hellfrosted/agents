#!/usr/bin/env python3
"""Plan or create WSL-native manual Git worktrees for candidate fan-out."""

from __future__ import annotations

import argparse
import json
import os
import re
import subprocess
import sys
from datetime import datetime
from pathlib import Path
from typing import Any


def run_git(repo: Path, args: list[str]) -> str:
    completed = subprocess.run(
        ["git", "-C", str(repo), *args],
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    return completed.stdout.strip()


def is_wsl() -> bool:
    if os.environ.get("WSL_DISTRO_NAME"):
        return True
    try:
        return "microsoft" in Path("/proc/version").read_text(encoding="utf-8").lower()
    except OSError:
        return False


def slug(value: str) -> str:
    cleaned = re.sub(r"[^A-Za-z0-9._-]+", "-", value.strip()).strip("-._")
    return cleaned.lower() or "candidate"


def default_root() -> Path:
    return Path(os.environ.get("CODEX_WSL_WORKTREE_ROOT", Path.home() / "codex-worktrees"))


def resolve_ref(repo: Path, requested_ref: str | None) -> str:
    ref = requested_ref or run_git(repo, ["branch", "--show-current"])
    if not ref:
        raise ValueError("no ref supplied and current checkout is detached")
    run_git(repo, ["rev-parse", "--verify", "--quiet", f"{ref}^{{commit}}"])
    return ref


def build_plan(args: argparse.Namespace) -> dict[str, Any]:
    repo = Path(args.repo).resolve()
    repo_root = Path(run_git(repo, ["rev-parse", "--show-toplevel"]))
    verified_ref = resolve_ref(repo_root, args.ref)
    run_id = args.run_id or datetime.now().strftime("%Y%m%d%H%M%S")
    root = Path(args.root).expanduser() if args.root else default_root()

    if is_wsl() and str(root).startswith("/mnt/") and not args.allow_windows_mount:
        raise ValueError(
            f"refusing Windows-mounted worktree root in WSL: {root}; "
            "set CODEX_WSL_WORKTREE_ROOT under /home or pass --allow-windows-mount"
        )

    candidates = args.candidate or ["candidate"]
    worktrees = []
    for candidate in candidates:
        candidate_id = slug(candidate)
        path = root / repo_root.name / run_id / candidate_id
        worktrees.append(
            {
                "candidate": candidate,
                "candidate_id": candidate_id,
                "path": str(path),
                "create_command": [
                    "git",
                    "-C",
                    str(repo_root),
                    "worktree",
                    "add",
                    "--detach",
                    str(path),
                    verified_ref,
                ],
                "thread_target": {
                    "type": "project",
                    "projectId": str(path),
                    "environment": {"type": "local"},
                },
            }
        )

    return {
        "wsl": is_wsl(),
        "repo_root": str(repo_root),
        "verified_ref": verified_ref,
        "run_id": run_id,
        "root": str(root),
        "worktrees": worktrees,
    }


def create_worktrees(plan: dict[str, Any]) -> None:
    for worktree in plan["worktrees"]:
        path = Path(worktree["path"])
        path.parent.mkdir(parents=True, exist_ok=True)
        subprocess.run(worktree["create_command"], check=True)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", default=".", help="Repository path; defaults to cwd")
    parser.add_argument("--ref", help="Existing git ref to use; defaults to current branch")
    parser.add_argument("--root", help="Manual worktree root; defaults to CODEX_WSL_WORKTREE_ROOT or ~/codex-worktrees")
    parser.add_argument("--run-id", help="Stable run id; defaults to current timestamp")
    parser.add_argument("--candidate", action="append", help="Candidate id/title; repeat for multiple candidates")
    parser.add_argument("--create", action="store_true", help="Create the planned detached git worktrees")
    parser.add_argument("--allow-windows-mount", action="store_true", help="Allow /mnt/* worktree roots in WSL")
    args = parser.parse_args()

    try:
        plan = build_plan(args)
        if args.create:
            create_worktrees(plan)
        print(json.dumps(plan, indent=2))
        return 0
    except (OSError, subprocess.CalledProcessError, ValueError) as error:
        print(f"ERROR: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
