---
id: adl-terminal-downloader
title: adl Terminal Downloader
status: in_progress
created: "2026-07-08T21:08:00Z"
---

# Intent: adl Terminal Downloader

## Goal

Build `adl`, a Mole-style terminal download manager for macOS that wraps aria2c. Users get an interactive TUI for starting downloads, watching live progress bars, and managing active transfers — all without leaving the terminal.

## Users

- macOS users who already use Homebrew and the terminal
- Users who want a friendly wrapper around aria2c without learning its flags

## Problem

aria2c is powerful but CLI-heavy. Starting a download requires remembering flags for URL, destination, and torrent/magnet handling. There is no built-in interactive UI for progress and download management.

## Success Criteria

- Running `adl` opens an English TUI main menu (New Download, Active Downloads, Settings, Help, Quit)
- `adl download`, `adl list`, `adl settings` skip directly to the matching view
- New download: single input auto-detects HTTP URL, magnet link, or `.torrent` file path
- Destination defaults to `~/Downloads` with optional folder browser (`b` key)
- Active downloads show live progress bars (%, speed, ETA)
- Controls: `p` pause, `r` resume, `x` cancel, `Enter` for details
- aria2c checked on startup; clear error with `brew install aria2` hint if missing
- `--help` on all commands; Help screen in TUI with keybindings
- Config persists at `~/.config/adl/config.yaml`
- aria2c RPC daemon runs only while TUI is open; stops on exit (no background process)
- Builds as local binary via `go build`; Homebrew formula deferred to later

## Constraints

- Go + Bubble Tea + Cobra stack
- aria2c only (no yt-dlp or other backends in v1)
- English UI only
- macOS primary target
- No download history, scheduling, or notifications in v1
- No `man` page in v1 (deferred to Homebrew release)

## Notes

Design decisions captured via grill-me session on 2026-07-08. Autonomy level: balanced.
