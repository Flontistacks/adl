# Homebrew Tap for adl

Install `adl` via Homebrew from [Flontistacks/adl](https://github.com/Flontistacks/adl).

## Install (HEAD — latest main)

```bash
brew tap Flontistacks/adl https://github.com/Flontistacks/adl
brew install --HEAD Flontistacks/adl/adl
```

## Install (stable — after tagging a release)

```bash
brew tap Flontistacks/adl https://github.com/Flontistacks/adl
brew install Flontistacks/adl/adl
```

## What gets installed

- `adl` → `$(brew --prefix)/bin/adl`
- man page → `$(brew --prefix)/share/man/man1/adl.1`

Run `man adl` after install.

## Dependencies

- `aria2` (runtime) — installed automatically
- `go` (build only)

## Separate tap repository (optional)

To publish as `brew tap Flontistacks/tap`:

1. Create `github.com/Flontistacks/homebrew-tap`
2. Copy `Formula/adl.rb` into that repo
3. Tag `v0.1.0` on the main repo and update `sha256` in the formula

```bash
curl -L https://github.com/Flontistacks/adl/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256
```
