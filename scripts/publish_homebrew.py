#!/usr/bin/env python3

from __future__ import annotations

import hashlib
import shutil
import subprocess
import sys
import urllib.request
from dataclasses import dataclass
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parent.parent
VERSION_FILE = PROJECT_ROOT / "VERSION"
SIBLING_HOMEBREW_REPO = PROJECT_ROOT.parent / "homebrew-ged"
LOCAL_HOMEBREW_REPO = PROJECT_ROOT / "homebrew-ged"
HOMEBREW_REPO_URL = "https://github.com/colinta/homebrew-ged"
FORMULA_PATH = Path("Formula/ged.rb")
TARGETS = [
    "darwin-arm64",
    "darwin-amd64",
    "linux-arm64",
    "linux-amd64",
]


@dataclass
class HomebrewRepo:
    path: Path
    created_here: bool


def run(command: list[str], *, cwd: Path | None = None) -> None:
    subprocess.run(command, cwd=cwd, check=True)


def read_version() -> str:
    return VERSION_FILE.read_text().strip()


def asset_url(version: str, target: str) -> str:
    return f"https://github.com/colinta/ged/releases/download/{version}/ged-{target}.tar.gz"


def sha256_for_url(url: str) -> str:
    hasher = hashlib.sha256()
    with urllib.request.urlopen(url) as response:
        while True:
            chunk = response.read(1024 * 1024)
            if not chunk:
                break
            hasher.update(chunk)
    return hasher.hexdigest()


def locate_homebrew_repo() -> HomebrewRepo:
    if SIBLING_HOMEBREW_REPO.exists():
        print(f"Using sibling homebrew repo: {SIBLING_HOMEBREW_REPO}")
        return HomebrewRepo(SIBLING_HOMEBREW_REPO, created_here=False)

    if LOCAL_HOMEBREW_REPO.exists():
        print(f"Using local homebrew repo: {LOCAL_HOMEBREW_REPO}")
        return HomebrewRepo(LOCAL_HOMEBREW_REPO, created_here=False)

    print(f"Cloning {HOMEBREW_REPO_URL} to {LOCAL_HOMEBREW_REPO}")
    run(["git", "clone", HOMEBREW_REPO_URL, str(LOCAL_HOMEBREW_REPO)])
    return HomebrewRepo(LOCAL_HOMEBREW_REPO, created_here=True)


def render_formula(version: str, shas: dict[str, str]) -> str:
    return f'''class Ged < Formula
  desc "Streaming text editor for pipelines — modern sed alternative"
  homepage "https://github.com/colinta/ged"
  version "{version}"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "{asset_url(version, "darwin-arm64")}"
      sha256 "{shas["darwin-arm64"]}"

      def install
        bin.install "ged-darwin-arm64" => "ged"
        man1.install "ged.1"
      end
    else
      url "{asset_url(version, "darwin-amd64")}"
      sha256 "{shas["darwin-amd64"]}"

      def install
        bin.install "ged-darwin-amd64" => "ged"
        man1.install "ged.1"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "{asset_url(version, "linux-arm64")}"
      sha256 "{shas["linux-arm64"]}"

      def install
        bin.install "ged-linux-arm64" => "ged"
        man1.install "ged.1"
      end
    else
      url "{asset_url(version, "linux-amd64")}"
      sha256 "{shas["linux-amd64"]}"

      def install
        bin.install "ged-linux-amd64" => "ged"
        man1.install "ged.1"
      end
    end
  end

  test do
    assert_match "hello earth",
      pipe_output("#{{bin}}/ged 's/world/earth/'", "hello world\n").strip
  end
end
'''


def update_formula(repo: Path, version: str) -> None:
    shas: dict[str, str] = {}
    for target in TARGETS:
        url = asset_url(version, target)
        print(f"Fetching {url}")
        shas[target] = sha256_for_url(url)

    formula = render_formula(version, shas)
    formula_path = repo / FORMULA_PATH
    formula_path.write_text(formula)
    print(f"Updated {formula_path}")


def commit_and_push(repo: Path, version: str) -> None:
    run(["git", "add", str(FORMULA_PATH)], cwd=repo)

    diff = subprocess.run(
        ["git", "diff", "--cached", "--quiet", "--", str(FORMULA_PATH)],
        cwd=repo,
    )
    if diff.returncode == 0:
        print("No Homebrew formula changes to commit")
        return

    run(["git", "commit", "-m", f"ged {version}"], cwd=repo)
    run(["git", "push", "origin", "main"], cwd=repo)


def publish_homebrew(version: str) -> None:
    homebrew_repo = locate_homebrew_repo()
    try:
        run(["git", "pull"], cwd=homebrew_repo.path)
        update_formula(homebrew_repo.path, version)
        commit_and_push(homebrew_repo.path, version)
    finally:
        if homebrew_repo.created_here:
            shutil.rmtree(homebrew_repo.path)
            print(f"Removed temporary clone: {homebrew_repo.path}")


def main() -> int:
    version = read_version()
    try:
        publish_homebrew(version)
    except subprocess.CalledProcessError as error:
        print(f"command failed: {error.cmd}", file=sys.stderr)
        return error.returncode or 1
    except Exception as error:  # noqa: BLE001
        print(str(error), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
