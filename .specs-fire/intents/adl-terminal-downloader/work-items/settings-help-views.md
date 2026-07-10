---
id: settings-help-views
title: Settings and Help Views
intent: adl-terminal-downloader
complexity: low
mode: autopilot
status: pending
depends_on: [tui-shell-menu, config-module]
created: "2026-07-08T21:10:00Z"
---

# Work Item: Settings and Help Views

## Description

Settings screen for default download path and aria2c path. Help screen with keybindings.

## Acceptance Criteria

- [ ] Settings: edit default download directory and aria2c path
- [ ] Settings persist to config file on save
- [ ] Help screen lists all keybindings and subcommands
- [ ] `Esc` returns to main menu

## Technical Notes

(none)

## Dependencies

- tui-shell-menu
- config-module
