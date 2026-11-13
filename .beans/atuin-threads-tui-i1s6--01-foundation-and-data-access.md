---
# atuin-threads-tui-i1s6
title: 01 - Foundation and data access
status: completed
type: milestone
priority: normal
created_at: 2026-11-13T10:26:03Z
updated_at: 2026-11-13T12:10:53Z
openspec-link: openspec/changes/archive/2026-11-13-atuin-threads-client
---

Project skeleton, Go module, Nix flake, and read-only access to the atuin history database so the TUI has something to show.

## Definition of done
- Go module builds and tests run under `nix flake check`.
- atuin history is read from the local SQLite database (or `atuin` CLI) without mutating it.
- A repository abstraction exists that the UI can consume in tests.
