# adl

Terminal download manager for macOS — an interactive TUI around [aria2c](https://aria2.github.io/).

`adl` gives you an interactive menu in the terminal to add downloads, watch live progress bars, and pause or cancel transfers. No GUI windows, no menu bar icon, no background daemon after you quit.

## Install

```bash
brew install Flontistacks/tap/adl
```

This installs `adl`, the `man adl` page, and the `aria2` dependency.

**Upgrade:**

```bash
brew update && brew upgrade Flontistacks/tap/adl
```

**Latest development build:**

```bash
brew install --HEAD Flontistacks/tap/adl
```

If you previously tapped the old formula from the main repo, remove it:

```bash
brew untap flontistacks/adl
```

## Quick start

```bash
adl              # open main menu
adl download     # add a download directly
adl list         # active downloads with progress bars
adl settings     # default folder & aria2c path
man adl          # full manual
```

**Example download flow:**

1. Run `adl download`
2. Paste a URL, magnet link, or path to a `.torrent` file
3. Press `Enter` — destination defaults to `~/Downloads` (press `b` to browse folders)
4. Press `Enter` again to start
5. Watch progress under **Active Downloads**

Try a small test file:

```text
https://proof.ovh.net/files/1Mb.dat
```

## Features

- **Interactive TUI** — built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **HTTP/HTTPS, magnet links, and `.torrent` files** — auto-detected from a single input field
- **Live progress bars** — speed, percentage, and ETA while the TUI is open
- **Download controls** — pause, resume, cancel, and view details
- **Folder browser** — pick a destination with `b` (default: `~/Downloads`)
- **Session-scoped aria2c** — RPC daemon starts with the app and stops when you quit
- **English UI** — menu, prompts, and help text
- **No telemetry** — `adl` does not phone home; downloads go through aria2c only

## Commands

| Command | Description |
|---------|-------------|
| `adl` | Main menu |
| `adl download` | New download flow |
| `adl list` | Active downloads |
| `adl settings` | Edit settings |
| `adl --help` | Short CLI help |
| `man adl` | Full manual |

## Keybindings

### Main menu

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate |
| `Enter` | Select |
| `?` | Help |
| `Ctrl+c` | Quit |

### New download

| Key | Action |
|-----|--------|
| `Enter` | Next step / start download |
| `b` | Browse destination folder |
| `Esc` | Back to menu |

### Active downloads

| Key | Action |
|-----|--------|
| `j` / `k` | Select download |
| `p` | Pause |
| `r` | Resume |
| `x` | Cancel |
| `Enter` | Details |
| `Esc` | Back to menu |

## Configuration

Settings are stored in `~/.config/adl/config.yaml` (created on first save, mode `0600`).

```yaml
download_dir: ~/Downloads
aria2_path: /opt/homebrew/bin/aria2c
rpc_port: 6800
```

You can also edit these from the in-app **Settings** screen.

The aria2 RPC secret is generated per session and is **not** written to disk.

## Security & privacy

- **Local only** — aria2 RPC listens on `127.0.0.1` while `adl` is running
- **No background process** — quitting the TUI stops the aria2c daemon
- **No analytics** — no usage data sent to Flontistacks or third parties
- **Per-session RPC secret** — random token generated each time you launch `adl`
- **Torrents & magnets** — peer-to-peer traffic is inherent to those protocols; use HTTPS when privacy matters

## Requirements

- macOS (primary target)
- [aria2c](https://aria2.github.io/) — installed automatically via Homebrew

## Build from source

```bash
git clone https://github.com/Flontistacks/adl.git
cd adl
go build -o adl ./cmd/adl
./adl
```

Or with Make:

```bash
make build
make install PREFIX=/opt/homebrew   # Apple Silicon
make install PREFIX=/usr/local      # Intel Mac
```

**Requirements for building:** Go 1.22+

## Project structure

```text
cmd/adl/          CLI entry point (Cobra)
internal/
  aria2/          RPC client & daemon lifecycle
  config/         ~/.config/adl/config.yaml
  download/       URL / magnet / torrent detection
  tui/            Bubble Tea views
man/adl.1         Manual page
Formula/          Homebrew formula (see tap repo)
```

## Related repositories

| Repository | Purpose |
|------------|---------|
| [Flontistacks/adl](https://github.com/Flontistacks/adl) | Source code |
| [Flontistacks/homebrew-tap](https://github.com/Flontistacks/homebrew-tap) | Homebrew formula |

## License

MIT — see [LICENSE](LICENSE).
