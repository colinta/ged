class Ged < Formula
  desc "Streaming text editor for pipelines — modern sed alternative"
  homepage "https://github.com/colinta/ged"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/colinta/ged/releases/download/1.0.0/ged-darwin-arm64.tar.gz"
      sha256 "82ee254cdfe4f0d9e928d71f0e760be6af9aacef944bea2596735899ad42e035"

      def install
        bin.install "ged-darwin-arm64" => "ged"
        man1.install "ged.1"
      end
    else
      url "https://github.com/colinta/ged/releases/download/1.0.0/ged-darwin-amd64.tar.gz"
      sha256 "d2e89761a3e07cedffe5db451976dc6917a83b62ef672a43777a21b01bad8dd6"

      def install
        bin.install "ged-darwin-amd64" => "ged"
        man1.install "ged.1"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/colinta/ged/releases/download/1.0.0/ged-linux-arm64.tar.gz"
      sha256 "a1067fb096e2ddc18d31a807bce2c56e3752b91ed2796b56391486223c21e23c"

      def install
        bin.install "ged-linux-arm64" => "ged"
        man1.install "ged.1"
      end
    else
      url "https://github.com/colinta/ged/releases/download/1.0.0/ged-linux-amd64.tar.gz"
      sha256 "943a566f65039e275b695805c91399ea5fd072f32cc67dca422598fe2617a99c"

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
