# Homebrew formula template for ged.
# To use: create a repo github.com/colinta/homebrew-ged with this file as Formula/ged.rb
# Then: brew tap colinta/ged && brew install ged
#
# Before publishing, update:
#   1. url — point to the release tag tarball
#   2. sha256 — run: curl -sL <url> | shasum -a 256
#
class Ged < Formula
  desc "Streaming text editor for pipelines — modern sed alternative"
  homepage "https://github.com/colinta/ged"
  # Update URL and sha256 for each release:
  url "https://github.com/colinta/ged/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "go" => :build
  depends_on "node" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/ged"
    system "npx", "marked-man", "README.md", "--output", "ged.1"
    man1.install "ged.1"
  end

  test do
    assert_match "hello earth",
      pipe_output("#{bin}/ged 's/world/earth/'", "hello world\n").strip
  end
end
