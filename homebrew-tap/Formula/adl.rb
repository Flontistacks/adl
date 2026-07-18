# typed: false
# frozen_string_literal: true

class Adl < Formula
  desc "Terminal download manager powered by aria2c"
  homepage "https://github.com/Flontistacks/adl"
  url "https://github.com/Flontistacks/adl/archive/refs/tags/v0.1.5.tar.gz"
  version "0.1.5"
  sha256 "3ef1f7dccc19d4a8f3e30abaea01c05c599faa2a73e5e37825c567d0667d9164"
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
