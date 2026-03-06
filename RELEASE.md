# Release Process

## Prerequisites

- All tests pass: `go test ./...`
- README is up to date: `make readme`

## Steps

### 1. Tag the release

```bash
git tag X.Y.Z
git push origin main --tags
```

### 2. Get the tarball sha256

```bash
curl -sL https://github.com/colinta/ged/archive/refs/tags/X.Y.Z.tar.gz | shasum -a 256
```

### 3. Update the Homebrew formula

Edit `~/src/github.com/colinta/homebrew-ged/Formula/ged.rb`:

- Update the `url` line with the new tag
- Update the `sha256` line with the hash from step 2

```bash
cd ~/src/github.com/colinta/homebrew-ged
# edit Formula/ged.rb
git add -A && git commit -m "ged X.Y.Z"
git push origin main
```

### 4. Update the local template

Copy the updated formula back to this repo:

```bash
cp ~/src/github.com/colinta/homebrew-ged/Formula/ged.rb scripts/homebrew-formula.rb
```

## Installation

Users install with:

```bash
brew tap colinta/ged
brew install ged
```

Or without Go/Homebrew:

```bash
go install github.com/colinta/ged/cmd/ged@latest
```
