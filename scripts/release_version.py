#!/usr/bin/env python3

from __future__ import annotations

import os
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

import publish_homebrew

PROJECT_ROOT = Path(__file__).resolve().parent.parent
VERSION_FILE = PROJECT_ROOT / "VERSION"
TARGETS = [
    "darwin-arm64",
    "darwin-amd64",
    "linux-arm64",
    "linux-amd64",
]


def run(command: list[str], *, cwd: Path | None = None, env: dict[str, str] | None = None) -> None:
    subprocess.run(command, cwd=cwd, env=env, check=True)


def read_version() -> str:
    return VERSION_FILE.read_text().strip()


def ensure_no_staged_changes() -> None:
    result = subprocess.run(
        ["git", "diff", "--cached", "--quiet"],
        cwd=PROJECT_ROOT,
    )
    if result.returncode != 0:
        raise RuntimeError("refusing to release with staged changes; commit or unstage them first")


def build_release_assets(release_dir: Path, version: str) -> list[str]:
    for target in TARGETS:
        goos, goarch = target.split("-", 1)
        env = os.environ.copy()
        env["GOOS"] = goos
        env["GOARCH"] = goarch
        env["CGO_ENABLED"] = "0"
        output_path = release_dir / f"ged-{target}"
        run(
            [
                "go",
                "build",
                "-ldflags",
                f"-s -w -X main.version={version}",
                "-o",
                str(output_path),
                "./cmd/ged",
            ],
            cwd=PROJECT_ROOT,
            env=env,
        )

    run(["npx", "marked-man", "README.md", "--output", str(release_dir / "ged.1")], cwd=PROJECT_ROOT)

    asset_paths: list[str] = []
    for target in TARGETS:
        archive_path = release_dir / f"ged-{target}.tar.gz"
        run(
            ["tar", "czf", str(archive_path), f"ged-{target}", "ged.1"],
            cwd=release_dir,
        )
        asset_paths.append(str(archive_path))

    return asset_paths


def release(version: str) -> None:
    ensure_no_staged_changes()

    release_dir = Path(tempfile.mkdtemp(prefix="ged-release-"))
    try:
        run(["git", "tag", version], cwd=PROJECT_ROOT)
        run(["git", "push", "origin", "main", "--tags"], cwd=PROJECT_ROOT)

        asset_paths = build_release_assets(release_dir, version)
        run(
            [
                "gh",
                "release",
                "create",
                version,
                *asset_paths,
                "--title",
                f"ged {version}",
                "--generate-notes",
            ],
            cwd=PROJECT_ROOT,
        )

        publish_homebrew.publish_homebrew(version)
    finally:
        shutil.rmtree(release_dir)


def main() -> int:
    version = read_version()
    try:
        release(version)
    except subprocess.CalledProcessError as error:
        print(f"command failed: {error.cmd}", file=sys.stderr)
        return error.returncode or 1
    except Exception as error:  # noqa: BLE001
        print(str(error), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
