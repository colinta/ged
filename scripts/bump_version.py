#!/usr/bin/env python3

from pathlib import Path
import sys

USAGE = "usage: just bump major|minor|bug"
VALID_PARTS = {"major", "minor", "bug"}
VERSION_FILE = Path(__file__).resolve().parent.parent / "VERSION"


def main() -> int:
    if len(sys.argv) != 2 or sys.argv[1] not in VALID_PARTS:
        print(USAGE, file=sys.stderr)
        return 1

    part = sys.argv[1]
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
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
