# typed: false
# frozen_string_literal: true

class Adl < Formula
  desc "Terminal download manager powered by aria2c"
  homepage "https://github.com/Flontistacks/adl"
  url "https://github.com/Flontistacks/adl/archive/refs/tags/v0.1.3.tar.gz"
  sha256 "5ca5b32c805496f79b6ee170c83bd0cf8d9a031ca80a39f87a61ff00c8856f6a"
  license "MIT"
  version "0.1.3"

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
