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

## Install

```bash
go install github.com/colinta/ged/cmd/ged@latest
```

## Usage

```
HEADER

go run ./cmd/ged --help | sed 's/^//'

cat <<'FOOTER'
```

## Examples

```bash
# Simple substitution
echo "hello world" | ged 's/world/earth/'

# Global replace
echo "aaa" | ged 's/a/b/g'

# Filter lines (like grep)
cat log.txt | ged 'p/ERROR/'

# Filter with context (like grep -C)
cat log.txt | ged 'p/ERROR/context=2'

# Delete matching lines
cat data.txt | ged 'd/^#/'

# Chain multiple rules
echo -e "c\na\nb" | ged 's/$/!/' sort

# Conditional rules
echo -e "hello\nworld" | ged 'if/hello/' '{' upper '}'

# File editing
ged 's/old/new/g' --input=file.txt --write

# Show diff before editing
ged 's/old/new/g' --input=file.txt --diff --color

# Explain what rules do
ged --explain 'p/error/i/context=2' 's/ERROR/WARNING/' upper
```

## License

MIT
FOOTER
