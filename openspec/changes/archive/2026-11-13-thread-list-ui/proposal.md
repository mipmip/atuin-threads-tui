# Proposal: thread-list-ui

## Why

The client can read atuin threads but nothing presents them. This change turns
`att` into a usable browser: a list of threads and a readable view of one
conversation.

## What Changes

- Add `internal/ui` with a bubbletea model, typed messages, and commands.
- Load threads on start, and show loading, empty, and error states.
- Navigate the list with `j`, `k`, arrows, `g`, and `G`; open with `enter`.
- Show thread metadata and every message, with scroll keys and wrapped text.
- Wire `cmd/att` to construct the atuin client and run the program.
- Handle terminal resize by re-truncating rows and re-wrapping messages.

## Capabilities

### New Capabilities

- `thread-browser`: the interactive list and conversation view, including
  navigation, scrolling, and terminal states.

### Modified Capabilities

None.

## Impact

- New dependency on bubbletea and lipgloss, vendored.
- `cmd/att` no longer reports a missing UI.

## Bean

- atuin-threads-tui-qjcn: Thread list and navigation
