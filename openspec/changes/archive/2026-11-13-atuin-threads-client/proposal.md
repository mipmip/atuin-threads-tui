# Proposal: atuin-threads-client

## Why

The TUI needs threads, messages, and search results. atuin already stores and
indexes captured agent threads, so the application should read them through
atuin rather than reimplement capture or couple itself to atuin's SQLite
schema.

## What Changes

- Add `internal/atuin` with thread, message, content block, and token types.
- Add a client that runs `atuin threads list`, `show`, and `search` with
  `--style json` and decodes the result.
- Surface a useful error, including atuin's stderr, when the command fails.
- Keep process execution behind a `Runner` seam so tests use fixtures.

## Capabilities

### New Capabilities

- `thread-source`: listing, reading, and searching atuin threads through the
  atuin CLI.

### Modified Capabilities

None.

## Impact

- New dependency-free package `internal/atuin`.
- No change to the flake beyond the vendored module cache.

## Bean

- atuin-threads-tui-w2e2: Read atuin history
