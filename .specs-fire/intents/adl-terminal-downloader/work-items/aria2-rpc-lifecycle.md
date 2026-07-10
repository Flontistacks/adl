---
id: aria2-rpc-lifecycle
title: aria2 RPC Lifecycle
intent: adl-terminal-downloader
complexity: medium
mode: confirm
status: pending
depends_on: [project-scaffold]
created: "2026-07-08T21:10:00Z"
---

# Work Item: aria2 RPC Lifecycle

## Description

Start aria2c as RPC daemon when TUI opens, communicate via JSON-RPC on localhost, stop daemon on TUI exit.

## Acceptance Criteria

- [ ] Start aria2c with `--enable-rpc --rpc-listen-all=false`
- [ ] JSON-RPC client for addUri, tellStatus, pause, unpause, remove
- [ ] Clean shutdown kills daemon process
- [ ] No daemon left running after exit

## Technical Notes

Session-scoped only. Poll interval ~500ms for status updates.

## Dependencies

- project-scaffold
