class Netmon < Formula
  desc "Network monitoring CLI tool"
  homepage "https://github.com/zzzzseong/netmon"
  version "1.0.0"
  license "MIT"
  head "https://github.com/zzzzseong/netmon.git", branch: "main"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-darwin-amd64.tar.gz"
      sha256 "" # Update with actual SHA256 for darwin-amd64
    elsif Hardware::CPU.arm?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-darwin-arm64.tar.gz"
      sha256 "" # Update with actual SHA256 for darwin-arm64
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-linux-amd64.tar.gz"
      sha256 "" # Update with actual SHA256 for linux-amd64
    elsif Hardware::CPU.arm?
      url "https://github.com/zzzzseong/netmon/releases/download/v1.0.0/netmon-linux-arm64.tar.gz"
      sha256 "" # Update with actual SHA256 for linux-arm64
    end
  end

  def install
    bin.install "netmon"
  end

  test do
    assert_match "Network Monitoring Tool", shell_output("#{bin}/netmon help", 1)
  end
end

