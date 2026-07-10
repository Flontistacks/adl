# Homebrew Tap for adl

Install `adl` via Homebrew once this tap is published on GitHub.

## Setup (one-time)

Push this repository to GitHub, then create a separate tap repo:

```bash
# Option A: dedicated tap repository (recommended)
# Create github.com/gertvanduijn/homebrew-tap and copy Formula/adl.rb there

# Option B: use the homebrew-tap/ folder in this repo as a tap
brew tap gertvanduijn/tap https://github.com/gertvanduijn/homebrew-tap
brew install adl
```

## Install from local checkout (before publishing)

From the project root, with git initialized:

```bash
git init
git add .
git commit -m "Initial commit"

brew install --HEAD ./Formula/adl.rb
```

This installs:

- `adl` → `$(brew --prefix)/bin/adl`
- man page → `$(brew --prefix)/share/man/man1/adl.1`

Run `man adl` after install.

## Dependencies

- `aria2` (runtime) — installed automatically
- `go` (build only)

## Release checklist

1. Tag `v0.1.0` on GitHub
2. Compute SHA256: `curl -L https://github.com/gertvanduijn/adl/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256`
3. Update `sha256` in `homebrew-tap/Formula/adl.rb`
4. Push to tap repository
