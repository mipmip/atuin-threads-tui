# Detail frame width

## Why

At 120 columns the conversation view leaves 14 cells empty on the right and
still wraps code at 80, so 16% of command and result lines continue onto a
second row while the window has room. The content also starts in the first
cell of the terminal and the meta rail stops at column 106 instead of reaching
the window edge, which makes the view look like it was cut off rather than
framed.

## What Changes

- Give the detail content a margin of one cell on the left and one on the
  right.
- Pin the meta rail to the right edge of the window: 24 cells wide, one cell
  from the edge, its text left-aligned inside the rail.
- Let the body column fill everything between the left margin and the rail, so
  it is `width - 28` instead of a fixed 80.
- Wrap commands, diffs, patches, and tool output at the body column instead of
  80. Prose keeps its 80-cell measure so long lines stay readable.
- Drop the rail below 48 columns and give the body the width minus the two
  margin cells.
- Correct the spec, which still describes the rail as 20 cells while it has
  been 24 since the two-column layout.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the detail view gains frame margins, a body column that
  follows the window, code that wraps at that column, and a meta rail pinned to
  the right edge.

## Impact

- `internal/ui/detail.go`: the layout, the row assembly, and the width each
  block type is rendered at.
- `internal/ui/detail_test.go`: new width expectations at 170, 120, 80, 48, and
  47 columns.
- No change to the atuin client, the list view, the render cache, or the
  on-disk data.

## Bean

- atuin-threads-tui-hp0w: detailview frame should have margins
- atuin-threads-tui-0ufz: softwrap text, hard wrap code
