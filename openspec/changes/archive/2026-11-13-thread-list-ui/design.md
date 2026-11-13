# Design: thread-list-ui

## Model

One bubbletea program with a small state machine:

```
modeList   --enter-->  modeDetail
modeDetail --esc----->  modeList
```

All I/O runs in `tea.Cmd` functions and reports back through messages:

```go
type threadsMsg      struct { threads []atuin.Thread; err error }
type conversationMsg struct { conversation atuin.Conversation; err error }
```

The model never calls the source during `Update` or `View`, so tests can drive
it with messages and no goroutines.

## Rendering

`View` composes plain functions that take width and a thread:

- `threadRow(width, thread, selected)` truncates with an ellipsis and keeps a
  fixed column for the short id, harness, message count, and age.
- `wrap(text, width)` folds long message text on word boundaries.
- `renderConversation(conversation, width, scroll, height)` returns the visible
  lines, so the detail scroller clamps against the real wrapped length.

Detail scrolling uses rendered lines, not messages, because one message can be
taller than the terminal.

## Layout

Four bands: title bar, body, status line, key hints. The body gets
`height - 3` rows. If the terminal is too small, the body shrinks to zero
instead of panicking.

## Styling

lipgloss styles for the selected row, roles, keys, and errors. The selected row
also carries a `>` marker so selection survives a monochrome terminal and
plain-text tests.

## Out of scope

Filtering, searching, and persisting state.
