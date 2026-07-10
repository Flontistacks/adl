# typed: false
# frozen_string_literal: true

# Install: brew install Flontistacks/tap/adl
# Tap repo: https://github.com/Flontistacks/homebrew-tap

class Adl < Formula
  desc "Terminal download manager powered by aria2c"
  homepage "https://github.com/Flontistacks/adl"
  url "https://github.com/Flontistacks/adl/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "4518bedf13dc1bd1b93970beffd7269047e2d2496767ed28f327e70d75bf5286"
  license "MIT"
  version "0.1.0"

  head "https://github.com/Flontistacks/adl.git", branch: "main", using: :git

  depends_on "go" => :build
  depends_on "aria2"

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/adl"
    man1.install "man/adl.1"
  end

  test do
    assert_match "Terminal download manager", shell_output("#{bin}/adl --help")
    assert_path_exists man1/"adl.1"
  end
end
