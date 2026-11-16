# Proposal

## Why

A thread carries far more than the conversation. On the largest local thread,
11,520 messages hold 2,450 tool calls, 2,449 tool results, 1,661 messages that
are nothing but thinking, and 3,544 attachment, title, and empty system
messages. Only 1,301 messages are actual prompts and replies, and reaching them
means scrolling through 55,527 rendered lines. There is no way to hide what the
reader does not want.

## What Changes

- Add a message type mask with one switch per category: prompts, replies,
  thinking, tool calls, tool results, patches, summaries, and notes.
- Filter the conversation while it is built: hidden blocks are dropped, and a
  message with no visible block left is not rendered at all.
- Notes are the empty messages that currently collapse into `(16 attachments)`.
  With notes off they disappear entirely, and with notes on they behave as
  today.
- Show only the visible kinds in the meta column, so the rail never advertises
  content that is not on screen.
- Add a panel, opened with `t`, that lists the categories with their counts and
  toggles them. Space toggles one row, `a` selects all, `n` selects the
  conversation, `1` to `4` select a preset, `enter` applies, `esc` cancels.
- Keep the mask for the rest of the session, so every thread opened afterwards
  uses it. A configuration file is out of scope and belongs to
  `atuin-threads-tui-3ab2`.
- Report the count as `N of M msgs` against atuin's message count, and name the
  active filter in the status line.
- Rebuild the rendering keyed on width and mask, in the background above 500
  messages, and keep the message at the top of the viewport on top.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the conversation shows a chosen subset of message types
  instead of every message, and the render cache is keyed on the mask as well as
  the width.

## Impact

`internal/ui/detail.go` (category of a block, filtering while building),
`internal/ui/app.go` (mask state, panel mode, rebuild), `internal/ui/view.go`
(panel drawing, status, hints), tests in `internal/ui`. No change to
`internal/atuin`, no new dependency, no data migration.

## Bean

`atuin-threads-tui-fyb0`
