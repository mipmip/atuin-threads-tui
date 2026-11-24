## Why

Selecting a message today gives one `>` in the left margin on the first line of
that message. A message can be a hundred lines tall, so the marker sits far from
where the eye is, and nothing shows which message is current before a range is
committed. The range model (anchor..cursor) also cannot express "mark these
three, skip the one between them", which is what the bean asks for.

The bean reports this as a bug: the select UI is unclear, and rows should be
visible as rows before anything is selected.

## What Changes

- Every message in an open conversation gets a background tint, alternating
  between two dark greys, so message boundaries are visible while scrolling.
- Selection mode marks the whole message instead of one cell: the cursor message
  gets a stronger tint and a pink `>` stripe in the left margin, and every
  marked message gets its own tint and a pink `#` stripe.
- `space` toggles the message under the cursor in or out of the marked set, so a
  few messages can be picked without dragging a range along.
- `y` and `w` copy or write the marked messages in document order; leaving the
  mode with `esc` keeps the view where it is.
- Painting happens on the visible lines only, so moving the cursor never
  rebuilds the conversation.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the selection requirement changes from a range with a
  margin marker to a marked set with tinted rows, and a new requirement covers
  the alternating background of messages.

## Impact

- `internal/ui/select.go`: the marked set replaces the anchor/cursor range.
- `internal/ui/view.go` and `internal/ui/detail.go`: painting of visible lines,
  using the existing `lineInfo`.
- `openspec/specs/thread-browser/spec.md` and the selection tests.
- No new dependency, no client change, no on-disk data.

## Bean

atuin-threads-tui-ds02
