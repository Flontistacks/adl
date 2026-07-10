# adl

Terminal download manager for macOS — a Mole-style TUI around [aria2c](https://aria2.github.io/).

## Requirements

- Go 1.22+
- aria2c (`brew install aria2`)

## Build

```bash
go build -o adl ./cmd/adl
# or
make build
```

## Install

### Homebrew (recommended)

```bash
brew tap Flontistacks/adl https://github.com/Flontistacks/adl
brew install --HEAD Flontistacks/adl/adl
```

After a release tag (`v0.1.0`):

```bash
brew install Flontistacks/adl/adl
```

Installs `adl` and `man adl`.

### Manual

```bash
make install PREFIX=/opt/homebrew   # Apple Silicon
make install PREFIX=/usr/local      # Intel Mac
```

## Usage

```bash
adl              # Main menu
adl download     # New download
adl list         # Active downloads
adl settings     # Settings
adl --help       # Quick help
man adl          # Full manual (after install)
```

## Features

- Interactive TUI (Bubble Tea) with English UI
- HTTP/HTTPS URLs, magnet links, and `.torrent` files
- Default download folder `~/Downloads` with folder browser (`b`)
- Live progress bars while the TUI is open
- Pause, resume, cancel controls
- Config at `~/.config/adl/config.yaml`
- aria2c daemon stops when you quit (no background process)

## Config

`~/.config/adl/config.yaml`:

```yaml
download_dir: ~/Downloads
aria2_path: /opt/homebrew/bin/aria2c
rpc_port: 6800
```

## License

MIT
