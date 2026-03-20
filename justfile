version := `cat VERSION`

build:
    go build -ldflags "-X main.version={{version}}" -o ged ./cmd/ged/

install:
    go install -ldflags "-X main.version={{version}}" ./cmd/ged/

test:
    go test ./...

bump part:
    python3 scripts/bump_version.py {{part}}

release: test docs
    python3 scripts/release_version.py

publish-homebrew:
    python3 scripts/publish_homebrew.py

readme:
    ./scripts/generate-readme.sh > README.md

man: readme
    npx marked-man README.md > ged.1

docs: readme man

clean:
    rm -rf ged ged.1 dist

# Build release binaries for all platforms
dist: clean docs
    mkdir -p dist
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version={{version}}" -o dist/ged-darwin-arm64 ./cmd/ged/
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version={{version}}" -o dist/ged-darwin-amd64 ./cmd/ged/
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version={{version}}" -o dist/ged-linux-arm64 ./cmd/ged/
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version={{version}}" -o dist/ged-linux-amd64 ./cmd/ged/
    cd dist && for bin in ged-*; do tar czf "${bin}.tar.gz" "$bin" -C .. ged.1; done
