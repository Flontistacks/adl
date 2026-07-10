# typed: false
# frozen_string_literal: true

# Local development formula — install from this repo checkout:
#   brew install --HEAD ./Formula/adl.rb
#
# Requires a git repository (run `git init && git add . && git commit` first).

class Adl < Formula
  desc "Terminal download manager powered by aria2c"
  homepage "https://github.com/gertvanduijn/adl"
  license "MIT"
  version "0.1.0"

  head Pathname(__FILE__).realpath.dirname.parent.to_s, using: :git

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
