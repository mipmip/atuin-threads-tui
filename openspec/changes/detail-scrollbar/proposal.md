# Proposal

## Why

The conversation view can be over 1,600 screens long: the largest local thread
renders 55,527 lines, while 33 of them fit on screen. The status line reports
`line 1 of 55527`, but that number says nothing about the shape of the
conversation or where you are inside it. A scrollbar gives the same information
at a glance, and costs one column that is already empty.

## What Changes

- Add a scrollbar to the conversation view: a quiet track and a solid thumb in
  the right margin, one cell wide, spanning the body rows.
- The thumb length follows how much of the thread is visible, and its position
  follows the scroll offset.
- Hide the bar when the whole thread fits on one screen, but keep the column
  reserved so the body width never depends on the content.
- Leave the list view untouched; the status line keeps `line X of Y`.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: a new requirement that the conversation view shows a
  scroll position.

## Impact

`internal/ui/detail.go` (track and thumb maths), `internal/ui/view.go` (drawing
the column), `internal/ui/detail_test.go` and a new test file for the maths. No
client, database, or configuration change.

This change builds on `detail-frame-width`, which gives the detail frame its
one-cell right margin. That change is applied first; the bar then occupies that
margin, so the body width stays whatever `detail-frame-width` decides.

## Bean

- `atuin-threads-tui-rt68`: slider in detailview
