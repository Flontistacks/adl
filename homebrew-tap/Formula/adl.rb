# typed: false
# frozen_string_literal: true

class Adl < Formula
  desc "Terminal download manager powered by aria2c"
  homepage "https://github.com/Flontistacks/adl"
  license "MIT"
  version "0.1.0"

  depends_on "go" => :build
  depends_on "aria2"

  on_macos do
    if build.head?
      head "https://github.com/Flontistacks/adl.git", branch: "main", using: :git
    else
      url "https://github.com/Flontistacks/adl/archive/refs/tags/v#{version}.tar.gz"
      sha256 "4518bedf13dc1bd1b93970beffd7269047e2d2496767ed28f327e70d75bf5286"
    end
  end

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/adl"
    man1.install "man/adl.1"
  end

  test do
    assert_match "Terminal download manager", shell_output("#{bin}/adl --help")
    assert_path_exists man1/"adl.1"
  end
end
