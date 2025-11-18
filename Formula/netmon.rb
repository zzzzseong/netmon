class Netmon < Formula
  desc "Network monitoring CLI tool"
  homepage "https://github.com/zzzzseong/netmon"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-darwin-amd64.tar.gz"
      sha256 "48588c0511ae42614c5e6adfeb7868faedf499f0d535337519fd9d5538665bfc"
    elsif Hardware::CPU.arm?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-darwin-arm64.tar.gz"
      sha256 "ed33f3c1c4bca2b6d91f258f881f749f48ec1b4cfaa94bff6e5120648867d2bd"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-linux-amd64.tar.gz"
      sha256 "10962504c38834e3e92840685e6f1b38907e406d28f6a2157116ad7d30cc97e7"
    elsif Hardware::CPU.arm?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-linux-arm64.tar.gz"
      sha256 "79fed0c296e2eb5d2a71a3a86df94c8a7f2bb9430e7ee1f177fa71aebc2197ff"
    end
  end

  def install
    bin.install "netmon"
  end

  test do
    assert_match "Network Monitoring Tool", shell_output("#{bin}/netmon help", 1)
  end
end

