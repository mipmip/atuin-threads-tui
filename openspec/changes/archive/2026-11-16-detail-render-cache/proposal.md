# Detail render cache

## Why

Reading a large thread is unusable. Opening the biggest thread in the local
database (11,520 messages, 55,527 rendered lines) takes 344 ms per render, and
every keystroke renders the whole conversation three times: once in
`clampScroll`, once in `body`, and once in `status`. That is about one second
per key press, and the freeze happens on the UI goroutine.

A normal thread is not much better. A 224-message thread needs 12.9 ms per
render, so a keystroke costs about 40 ms. Scrolling a conversation should cost
the same no matter how many messages it holds.

## What Changes

- Render the conversation into one block per message, with a running offset
  table, once per thread and terminal width.
- Draw only the blocks that overlap the viewport instead of slicing one huge
  slice of pre-rendered lines.
- Keep the rendering in the model, so scrolling, the status line, and the
  end-of-conversation clamp read a cached total instead of re-rendering.
- Rebuild the rendering when the terminal width changes.
- Build threads over 500 messages in the background while the view shows a
  rendering state, so opening never blocks the UI.
- Measure text faster: per-rune widths from `go-runewidth` in the splitting
  helpers, and wrap widths tracked incrementally instead of measuring the
  growing line for every word.
- Apply a batch of runes that arrives in one key message one key at a time, so
  fast scrolling is not swallowed when several keys arrive together.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the conversation is rendered once and scrolled from a
  cache, and a large thread is rendered in the background.

## Impact

- `internal/ui/detail.go`: the render becomes a per-block structure.
- `internal/ui/app.go`: the model owns the rendering and rebuilds it on open
  and resize.
- `internal/ui/view.go`: the detail body and status line read the cache.
- No change to the atuin client, the list view, or the on-disk data.
- `go.mod`: `github.com/mattn/go-runewidth` becomes a direct dependency (it is
  already vendored through lipgloss).
