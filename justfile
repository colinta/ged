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
dist: clean build docs
    mkdir -p dist
    GOOS=darwin  GOARCH=arm64 go build -o dist/ged-darwin-arm64 ./cmd/ged/
    GOOS=darwin  GOARCH=amd64 go build -o dist/ged-darwin-amd64 ./cmd/ged/
    GOOS=linux   GOARCH=arm64 go build -o dist/ged-linux-arm64  ./cmd/ged/
    GOOS=linux   GOARCH=amd64 go build -o dist/ged-linux-amd64  ./cmd/ged/
    cd dist && for bin in ged-*; do tar czf "${bin}.tar.gz" "$bin" -C .. ged.1; done

# Publish a new version: just release 1.2.0
release version: test dist
    #!/usr/bin/env bash
    set -euo pipefail

    version="{{version}}"

    # Sanity checks
    if git tag | grep -qx "$version"; then
        echo "Error: tag $version already exists" >&2
        exit 1
    fi
    if [ ! -d "{{homebrew_repo}}" ]; then
        echo "Error: homebrew repo not found at {{homebrew_repo}}" >&2
        exit 1
    fi

    # Compute SHA256s
    sha_darwin_arm64=$(shasum -a 256 dist/ged-darwin-arm64.tar.gz | awk '{print $1}')
    sha_darwin_amd64=$(shasum -a 256 dist/ged-darwin-amd64.tar.gz | awk '{print $1}')
    sha_linux_arm64=$(shasum -a 256 dist/ged-linux-arm64.tar.gz | awk '{print $1}')
    sha_linux_amd64=$(shasum -a 256 dist/ged-linux-amd64.tar.gz | awk '{print $1}')

    # Generate homebrew formula
    cat > scripts/homebrew-formula.rb <<EOF
    class Ged < Formula
      desc "Streaming text editor for pipelines — modern sed alternative"
      homepage "https://github.com/colinta/ged"
      version "$version"
      license "MIT"

      on_macos do
        if Hardware::CPU.arm?
          url "https://github.com/colinta/ged/releases/download/$version/ged-darwin-arm64.tar.gz"
          sha256 "$sha_darwin_arm64"

          def install
            bin.install "ged-darwin-arm64" => "ged"
            man1.install "ged.1"
          end
        else
          url "https://github.com/colinta/ged/releases/download/$version/ged-darwin-amd64.tar.gz"
          sha256 "$sha_darwin_amd64"

          def install
            bin.install "ged-darwin-amd64" => "ged"
            man1.install "ged.1"
          end
        end
      end

      on_linux do
        if Hardware::CPU.arm?
          url "https://github.com/colinta/ged/releases/download/$version/ged-linux-arm64.tar.gz"
          sha256 "$sha_linux_arm64"

          def install
            bin.install "ged-linux-arm64" => "ged"
            man1.install "ged.1"
          end
        else
          url "https://github.com/colinta/ged/releases/download/$version/ged-linux-amd64.tar.gz"
          sha256 "$sha_linux_amd64"

          def install
            bin.install "ged-linux-amd64" => "ged"
            man1.install "ged.1"
          end
        end
      end

      test do
        assert_match "hello earth",
          pipe_output("#{bin}/ged 's/world/earth/'", "hello world\n").strip
      end
    end
    EOF

    # Strip leading whitespace from heredoc (4 spaces)
    sed -i '' 's/^    //' scripts/homebrew-formula.rb

    # Copy to homebrew repo
    cp scripts/homebrew-formula.rb "{{homebrew_repo}}/Formula/ged.rb"

    # Commit and tag ged repo
    git add -A
    git commit -m "Release $version"
    git tag "$version"

    # Push ged repo and create GitHub release
    git push origin main
    git push origin "$version"
    gh release create "$version" \
        dist/ged-darwin-arm64.tar.gz \
        dist/ged-darwin-amd64.tar.gz \
        dist/ged-linux-arm64.tar.gz \
        dist/ged-linux-amd64.tar.gz \
        --title "$version" \
        --generate-notes

    # Commit and push homebrew repo
    cd "{{homebrew_repo}}"
    git add -A
    git commit -m "Update to $version"
    git push origin main

    echo ""
    echo "✅ Released $version"
