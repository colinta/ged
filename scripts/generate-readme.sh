#!/bin/bash
# Generate README.md from ged --help output.
# Single source of truth: the helpText constant in cmd/ged/main.go.
#
# Usage: ./scripts/generate-readme.sh > README.md

set -euo pipefail

cd "$(dirname "$0")/.."

cat <<'HEADER'
# ged

A streaming text editor for pipelines. Like `sed`, but with modern regex, intuitive syntax, and composable rules.

## Installation

### Homebrew (recommended — includes man page)

```bash
brew tap colinta/ged
brew install ged
```

### Go install (binary only, no man page)

```bash
go install github.com/colinta/ged/cmd/ged@latest
```

### Manual install

Download the binary for your platform from the
[latest release](https://github.com/colinta/ged/releases/latest), then:

```bash
tar xzf ged-darwin-arm64.tar.gz    # or your platform
cp ged-darwin-arm64 /usr/local/bin/ged
cp ged.1 /usr/local/share/man/man1/ged.1
```

## Usage

```
HEADER

go run ./cmd/ged help | sed 's/^//'

cat <<'FOOTER'
```

## License

MIT
FOOTER
