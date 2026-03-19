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

`ged` does **not** use strict semver reset behavior. `minor` and `bug` increments are independent.

Bump the version:

```bash
just bump bug    # or: major, minor
```

Then release that version:

```bash
just release
```

`just release` runs tests, regenerates docs, tags the current `VERSION`, creates the GitHub
release, and updates `homebrew-ged`.

If `../homebrew-ged` exists, it uses that repo. Otherwise it clones `github.com/colinta/homebrew-ged`
into `./homebrew-ged`, updates the formula, pushes it, and removes the temporary clone.

If you need to rerun only the Homebrew step after a release already exists, use:

```bash
just publish-homebrew
```
