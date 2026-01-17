# Homebrew formula for backstage-gen-cli
# To use this formula, add the tap:
#   brew tap gautampachnanda101/backstage-gen-cli
#   brew install backstage-gen-cli

class BackstageGenCli < Formula
  desc "CLI tool for generating and validating Backstage catalog-info.yaml files"
  homepage "https://github.com/gautampachnanda101/backstage-gen-cli"
  version "0.1.0-alpha.1"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v#{version}/backstage-gen-cli-darwin-amd64"
      sha256 "PLACEHOLDER_DARWIN_AMD64_SHA256"
    end

    on_arm do
      url "https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v#{version}/backstage-gen-cli-darwin-arm64"
      sha256 "PLACEHOLDER_DARWIN_ARM64_SHA256"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v#{version}/backstage-gen-cli-linux-amd64"
      sha256 "PLACEHOLDER_LINUX_AMD64_SHA256"
    end

    on_arm do
      url "https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v#{version}/backstage-gen-cli-linux-arm64"
      sha256 "PLACEHOLDER_LINUX_ARM64_SHA256"
    end
  end

  def install
    binary_name = "backstage-gen-cli"

    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "backstage-gen-cli-darwin-arm64" => binary_name
      else
        bin.install "backstage-gen-cli-darwin-amd64" => binary_name
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        bin.install "backstage-gen-cli-linux-arm64" => binary_name
      else
        bin.install "backstage-gen-cli-linux-amd64" => binary_name
      end
    end
  end

  test do
    assert_match "backstage-gen-cli", shell_output("#{bin}/backstage-gen-cli --version")

    # Test inspect command
    system "#{bin}/backstage-gen-cli", "inspect", "--json"

    # Test help
    system "#{bin}/backstage-gen-cli", "--help"
  end
end
