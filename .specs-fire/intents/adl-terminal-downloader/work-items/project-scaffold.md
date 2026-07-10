---
id: project-scaffold
title: Project Scaffold
intent: adl-terminal-downloader
complexity: medium
mode: confirm
status: pending
depends_on: []
created: "2026-07-08T21:10:00Z"
---

# Work Item: Project Scaffold

## Description

Initialize Go module, project structure, Cobra CLI with `adl` root and subcommands (`download`, `list`, `settings`), `--help` on all commands, and aria2c presence check on startup.

## Acceptance Criteria

- [ ] `go.mod` with module path and core dependencies
- [ ] `cmd/adl/main.go` entry point
- [ ] Cobra commands: root, download, list, settings
- [ ] `--help` works on all commands
- [ ] Startup checks for aria2c with `brew install aria2` hint
- [ ] `go build -o adl ./cmd/adl` succeeds

## Technical Notes

Project layout per coding-standards.md. Use charmbracelet and spf13/cobra.

## Dependencies

(none)
