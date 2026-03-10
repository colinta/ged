class Ged < Formula
  desc "Streaming text editor for pipelines — modern sed alternative"
  homepage "https://github.com/colinta/ged"
  version "1.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/colinta/ged/releases/download/1.1.0/ged-darwin-arm64.tar.gz"
      sha256 "63306dd9e025f2c80c2473db4d47bc44daf1414d5ca36e98c861b59b0c22e3cb"

      def install
        bin.install "ged-darwin-arm64" => "ged"
        man1.install "ged.1"
      end
    else
      url "https://github.com/colinta/ged/releases/download/1.1.0/ged-darwin-amd64.tar.gz"
      sha256 "dfcce7201c2a5a8e98ee5e6f6567364a6220e106e21c5d554c97981382059ef8"

      def install
        bin.install "ged-darwin-amd64" => "ged"
        man1.install "ged.1"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/colinta/ged/releases/download/1.1.0/ged-linux-arm64.tar.gz"
      sha256 "5696a5f4e22d7f46d72fcc1a767c46409740141334a6361aa42e0f01a89a1822"

      def install
        bin.install "ged-linux-arm64" => "ged"
        man1.install "ged.1"
      end
    else
      url "https://github.com/colinta/ged/releases/download/1.1.0/ged-linux-amd64.tar.gz"
      sha256 "898286addec513a9f2f18cef8b5942e41dd740b163f54729460d20886b4d5723"

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
