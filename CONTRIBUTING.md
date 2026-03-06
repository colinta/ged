# Contributing to ged

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- [just](https://github.com/casey/just) (task runner)
- [Node.js](https://nodejs.org/) (for man page generation via `npx marked-man`)

## Development

```bash
just build       # Build the binary to ./ged
just test        # Run all tests
just readme      # Regenerate README.md from --help output
just man         # Generate ged.1 man page (requires Node)
just docs        # Regenerate both README.md and ged.1
just clean       # Remove build artifacts
```

## Project Structure

```
cmd/ged/             CLI entry point, integration tests, YAML test data
internal/rule/       Rule implementations (LineRule, DocumentRule)
internal/parser/     Rule string parsing
internal/engine/     Processing pipeline
internal/diff/       LCS diff algorithm
scripts/             README generator, Homebrew formula template
```

## Adding a New Rule

1. Create `internal/rule/myrule.go` with the rule struct and `Apply`/`ApplyDocument` method
2. Add `Explain()` method in `internal/rule/explain.go`
3. Add parser support in `internal/parser/parser.go`
4. Write tests: rule unit tests + parser tests + YAML CLI tests in `cmd/ged/testdata/`
5. Update the `helpText` in `cmd/ged/main.go`
6. Run `just docs` to regenerate README and man page

## Documentation

The single source of truth for usage docs is the `helpText` constant in `cmd/ged/main.go`.

- `just readme` runs `ged --help` and wraps the output in markdown → `README.md`
- `just man` converts `README.md` to a man page via `npx marked-man` → `ged.1`

## Releasing

### 1. Verify

```bash
just test
just docs
```

### 2. Tag and push

```bash
git tag X.Y.Z
git push origin main --tags
```

### 3. Build release binaries and create GitHub release

```bash
mkdir -p /tmp/ged-release

# Build for all platforms
for os_arch in darwin-arm64 darwin-amd64 linux-arm64 linux-amd64; do
  GOOS=${os_arch%-*} GOARCH=${os_arch#*-} CGO_ENABLED=0 \
    go build -ldflags="-s -w" -o /tmp/ged-release/ged-${os_arch} ./cmd/ged
done

# Build man page
npx marked-man README.md --output /tmp/ged-release/ged.1

# Create tarballs
cd /tmp/ged-release
for os_arch in darwin-arm64 darwin-amd64 linux-arm64 linux-amd64; do
  tar czf ged-${os_arch}.tar.gz ged-${os_arch} ged.1
done

# Create GitHub release
cd /path/to/ged
gh release create X.Y.Z /tmp/ged-release/ged-*.tar.gz \
  --title "ged X.Y.Z" --notes "Release notes here"
```

### 4. Update the Homebrew tap

Update `Formula/ged.rb` in the [homebrew-ged](https://github.com/colinta/homebrew-ged) repo
with the new version number, download URLs, and sha256 hashes:

```bash
cd /tmp/ged-release
for f in ged-*.tar.gz; do echo "$f"; shasum -a 256 "$f"; done
```

Update the `url`, `sha256`, and `version` fields in the formula for each platform,
then push:

```bash
cd ~/src/github.com/colinta/homebrew-ged
# edit Formula/ged.rb
git add -A && git commit -m "ged X.Y.Z"
git push origin main
```

## Installation

### Homebrew (recommended — includes man page)

```bash
brew tap colinta/ged
brew install ged
```

### Manual install (if Homebrew is broken)

Download the binary for your platform from the
[latest release](https://github.com/colinta/ged/releases/latest), then:

```bash
tar xzf ged-darwin-arm64.tar.gz    # or your platform
cp ged-darwin-arm64 /usr/local/bin/ged
cp ged.1 /usr/local/share/man/man1/ged.1
```

### Go install (binary only, no man page)

```bash
go install github.com/colinta/ged/cmd/ged@latest
```
