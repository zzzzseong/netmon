class Netmon < Formula
  desc "Network monitoring CLI tool"
  homepage "https://github.com/zzzzseong/netmon"
  version "1.0.0"
  license "MIT"
  head "https://github.com/zzzzseong/netmon.git", branch: "main"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-darwin-amd64.tar.gz"
      sha256 "b01de7db39334ca5ae883f01f8afc1303e57fcf20fd74434fb0b1dd9785d3d7b"
    elsif Hardware::CPU.arm?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-darwin-arm64.tar.gz"
      sha256 "b7c76e690d5c898a2ad76842e143d18ff3c85f4f4b755687d101a2b0488e7d1d"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-linux-amd64.tar.gz"
      sha256 "4479b82f440dfe42c1ad36ef3da7c3ed9da06ad3b7191b8db6131893ec5e900f"
    elsif Hardware::CPU.arm?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-linux-arm64.tar.gz"
      sha256 "439d9c687a0ad9e63dd23922e85eb9608abb1267b3da1ed3fc99f68e8f69fcc8"
    end
  end

  def install
    bin.install "netmon"
  end

  test do
    assert_match "Network Monitoring Tool", shell_output("#{bin}/netmon help", 1)
  end
end

