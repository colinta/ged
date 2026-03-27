#!/usr/bin/env python3

from __future__ import annotations

from pathlib import Path
import subprocess
import sys

USAGE = "usage: just bump major|minor|bug"
VALID_PARTS = {"major", "minor", "bug"}
PROJECT_ROOT = Path(__file__).resolve().parent.parent
VERSION_FILE = PROJECT_ROOT / "VERSION"
RELEASE_WARNING = "you should commit before running 'just release'"


def run(command: list[str]) -> None:
    subprocess.run(command, cwd=PROJECT_ROOT, check=True)


def has_staged_changes() -> bool:
    result = subprocess.run(
        ["git", "diff", "--cached", "--quiet"],
        cwd=PROJECT_ROOT,
    )
    return result.returncode != 0


def working_tree_is_clean() -> bool:
    result = subprocess.run(
        ["git", "status", "--porcelain"],
        cwd=PROJECT_ROOT,
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip() == ""


def commit_version_bump(part: str) -> str:
    run(["git", "commit", "-m", f"{part} version bump", "VERSION"])
    result = subprocess.run(
        ["git", "rev-parse", "HEAD"],
        cwd=PROJECT_ROOT,
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip()


def tag_version(version: str, commit: str) -> None:
    run(["git", "tag", version, commit])


def main() -> int:
    if len(sys.argv) != 2 or sys.argv[1] not in VALID_PARTS:
        print(USAGE, file=sys.stderr)
        return 1

    part = sys.argv[1]
    had_staged_changes = has_staged_changes()
    version = VERSION_FILE.read_text().strip()

    try:
        major_text, minor_text, bug_text = version.split(".")
        major = int(major_text)
        minor = int(minor_text)
        bug = int(bug_text)
    except ValueError:
        print(f"invalid VERSION format: {version!r}", file=sys.stderr)
        return 1

    if part == "major":
        major += 1
    elif part == "minor":
        minor += 1
    else:
        bug += 1

    new_version = f"{major}.{minor}.{bug}\n"
    VERSION_FILE.write_text(new_version)
    print(f"VERSION -> {new_version.strip()}")

    if had_staged_changes:
        print(RELEASE_WARNING)
        return 0

    commit_sha = commit_version_bump(part)
    tag_version(new_version.strip(), commit_sha)

    if not working_tree_is_clean():
        print(RELEASE_WARNING)

    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except subprocess.CalledProcessError as error:
        print(f"command failed: {error.cmd}", file=sys.stderr)
        raise SystemExit(error.returncode or 1)
