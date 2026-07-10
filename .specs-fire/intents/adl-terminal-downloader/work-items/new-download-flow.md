---
id: new-download-flow
title: New Download Flow
intent: adl-terminal-downloader
complexity: medium
mode: confirm
status: pending
depends_on: [aria2-rpc-lifecycle, tui-shell-menu, config-module]
created: "2026-07-08T21:10:00Z"
---

# Work Item: New Download Flow

## Description

Interactive flow to add a new download with auto-detected input type and destination selection.

## Acceptance Criteria

- [ ] Single text input auto-detects HTTP URL, magnet link, or .torrent path
- [ ] Default destination `~/Downloads` from config
- [ ] `b` key opens folder browser
- [ ] Submit starts download via aria2 RPC

## Technical Notes

Use bubbles/textinput for URL field, bubbles/filepicker or custom browser for folders.

## Dependencies

- aria2-rpc-lifecycle
- tui-shell-menu
- config-module
