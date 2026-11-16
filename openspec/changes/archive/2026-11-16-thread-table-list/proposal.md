# Proposal: thread-table-list

## Why

The list view joins its fields with `·` separators, so every column starts at a
different offset and the width of one field pushes the next one around. The
result is a paragraph per thread, not a table: the eye cannot scan down the id,
the date, or the message count, and the relative age hides when a thread
actually ran. With 197 threads on this machine the list is the main surface, so
it has to be readable column by column.

## What Changes

- Replace the `·` separated row with a table: `id`, `date`, `agent`, `title`,
  `msgs`, and `dir`.
- Add a dim header row and a rule under it; both stay fixed while the rows
  scroll.
- Show the update date as `02 Jan 15:04` in local time instead of a relative
  age.
- Map the harness to a two letter agent tag: `opencode` to `oc`,
  `claude-code` to `cc`, and the first two letters of any other harness.
- Show the basename of the working directory in a fixed column, cut from the
  left so the tail of the name survives.
- Drop columns in a fixed order when the terminal gets narrow, so an 80 column
  terminal still has a usable title.
- Count the header and rule as chrome when scrolling, so the selection stays
  visible.
- Pin the flake build to `cmd/att`, so a nested Go module elsewhere in the tree
  cannot break `nix build`.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the list renders as an aligned table with a fixed header,
  and the visible row count follows from the table chrome.

## Impact

- `internal/ui/view.go` gains the table layout and rendering.
- `internal/ui/app.go` scroll math uses the table row count.
- Existing row tests change because `threadRow` moves behind a layout type.
- No new dependency: `charmbracelet/x/ansi` is already vendored through
  lipgloss and provides display width measurement and truncation.

## Bean

- atuin-threads-tui-5q8a: make a UI design first before continuing
