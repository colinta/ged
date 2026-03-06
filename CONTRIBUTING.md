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

### 3. Get the tarball sha256

```bash
curl -sL https://github.com/colinta/ged/archive/refs/tags/X.Y.Z.tar.gz | shasum -a 256
```

### 4. Update the Homebrew tap

Edit `Formula/ged.rb` in the [homebrew-ged](https://github.com/colinta/homebrew-ged) repo:

- Update the `url` with the new tag
- Update the `sha256` with the hash from step 3

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

### Go install (binary only, no man page)

```bash
go install github.com/colinta/ged/cmd/ged@latest
```
