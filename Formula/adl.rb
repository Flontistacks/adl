# typed: false
# frozen_string_literal: true

# Install: brew install Flontistacks/tap/adl
# Tap repo: https://github.com/Flontistacks/homebrew-tap

class Adl < Formula
  desc "Terminal download manager powered by aria2c"
  homepage "https://github.com/Flontistacks/adl"
  url "https://github.com/Flontistacks/adl/archive/refs/tags/v0.1.7.tar.gz"
  sha256 "f68d82a66be0949c8f185625c62b59791bf2bb6fee42528de4f81439499b5528"
  license "MIT"

  head "https://github.com/Flontistacks/adl.git", branch: "main", using: :git

  depends_on "go" => :build
  depends_on "aria2"

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/adl"
    man1.install "man/adl.1"
  end

  test do
    assert_match "terminal UI for managing downloads", shell_output("#{bin}/adl --help")
    assert_path_exists man1/"adl.1"
  end
end
