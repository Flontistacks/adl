---
id: active-downloads-view
title: Active Downloads View
intent: adl-terminal-downloader
complexity: medium
mode: confirm
status: pending
depends_on: [aria2-rpc-lifecycle, tui-shell-menu]
created: "2026-07-08T21:10:00Z"
---

# Work Item: Active Downloads View

## Description

List active downloads with live progress bars and controls.

## Acceptance Criteria

- [ ] Progress bars showing %, speed, ETA
- [ ] Poll aria2 RPC every ~500ms
- [ ] `p` pause, `r` resume, `x` cancel selected download
- [ ] `Enter` shows download details

## Technical Notes

Use bubbles/progress or custom bar rendering.

## Dependencies

- aria2-rpc-lifecycle
- tui-shell-menu
