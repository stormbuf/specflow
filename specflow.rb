# Homebrew formula for specflow
# 安装: brew install ./specflow.rb
# 或提交到 homebrew tap: brew tap stormbuf/specflow && brew install specflow

class Specflow < Formula
  desc "Specflow CLI — spec 驱动的变更生命周期管理工具"
  homepage "https://github.com/stormbuf/specflow"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/stormbuf/specflow/releases/download/v0.1.0/specflow-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64_DARWIN_SHA256"
    end
    on_intel do
      url "https://github.com/stormbuf/specflow/releases/download/v0.1.0/specflow-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64_DARWIN_SHA256"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/stormbuf/specflow/releases/download/v0.1.0/specflow-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64_LINUX_SHA256"
    end
    on_intel do
      url "https://github.com/stormbuf/specflow/releases/download/v0.1.0/specflow-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64_LINUX_SHA256"
    end
  end

  def install
    bin.install "specflow"
  end

  test do
    assert_match "specflow", shell_output("#{bin}/specflow --version")
  end
end
