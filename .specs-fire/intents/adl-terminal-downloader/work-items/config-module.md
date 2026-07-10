---
id: config-module
title: Config Module
intent: adl-terminal-downloader
complexity: low
mode: autopilot
status: pending
depends_on: [project-scaffold]
created: "2026-07-08T21:10:00Z"
---

# Work Item: Config Module

## Description

Load and save user configuration at `~/.config/adl/config.yaml` with sensible defaults.

## Acceptance Criteria

- [ ] Default download dir: `~/Downloads`
- [ ] Configurable aria2c path (auto-detect from PATH)
- [ ] Create config directory on first save
- [ ] Unit tests for load defaults and round-trip save

## Technical Notes

Use gopkg.in/yaml.v3.

## Dependencies

- project-scaffold
