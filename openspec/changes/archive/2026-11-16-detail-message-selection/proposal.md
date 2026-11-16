# Proposal

## Why

Reading a thread produces text that cannot get out of the terminal. The body is
wrapped at the display width, so copying the screen gives hard-wrapped code that
breaks when pasted. There is no way to lift a prompt, a reply, or a whole
exchange into an editor, a chat, or a bug report.

## What Changes

- Add select mode to the conversation view, entered with `v`. The cursor moves
  one message with `j` and `k`, ten with `J` and `K`, and to either end with
  `g` and `G`. The selection is the range from the anchor to the cursor.
- Draw the cursor and the selected messages with a marker in the left margin,
  which is already reserved, so no cached line is restyled and no message is
  rendered again to show the cursor.
- Copy the selection with `y` as Markdown, generated from the atuin message
  content rather than from the rendered lines, so wrapped code comes back out
  unwrapped.
- Copy through a native clipboard tool when one exists (`wl-copy`, `pbcopy`,
  `xclip`, `xsel`), because its exit status is observable, and fall back to an
  OSC 52 sequence otherwise. The status line names the route used.
- Write the selection to a file with `w`, and export the whole thread with `e`
  from browse mode.
- Keep a payload above roughly 50 KB away from the clipboard, where terminals
  drop it, and write it to a file instead.
- Export what the type filter shows, so hidden kinds stay out of the output.
- Skip collapsed attachment runs and messages without content when building the
  range, since they have nothing to export.
- Clear the selection on `esc`, on leaving the thread, and when the type
  selection changes, because message indices move.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the conversation gains message selection, a cursor, and
  clipboard and file export.

## Impact

`internal/ui/app.go` (select mode, keys), `internal/ui/detail.go` (message
indices and margin markers), `internal/ui/view.go` (status and hints), a new
Markdown formatter and a clipboard writer beside them, and tests. No change to
`internal/atuin`, no new module dependency because `ansi.SetSystemClipboard` is
already vendored.

## Bean

`atuin-threads-tui-sao6`
