# Design

## Context

Measured on the largest local thread (`01a18c27`, 11,520 messages, 55,527
rendered lines):

| role | messages | blocks |
| --- | --- | --- |
| assistant | 5,212 | tool_call 2,450, thinking 1,661, text 1,101 |
| tool | 2,449 | tool_result 2,449, patch 114 |
| attachment | 1,954 | none |
| ai-title | 1,164 | none |
| system | 433 | summary 7, none 426 |
| user | 200 | text 200, tool_result 1 |
| away_summary | 108 | summary 108 |

Two facts shape the design:

- Roles and block kinds are different axes. An assistant message usually holds
  thinking, text, and tool calls at once, so filtering must work on blocks and
  then remove messages that end up empty.
- Atuin's own `message_count` is 7,976, which counts messages that have content
  and excludes attachments and titles. The header already uses it while the
  status line uses `len(Messages)`. This change makes both agree.

How empty messages are spread matters too: the 3,544 empty messages form 2,873
runs, and 2,487 of those runs are a single message. Each single one costs a
two-line meta row plus a blank line, about 7,500 lines of the 55,527. That is
why notes is a switch of its own.

## Categories

The mask has one bit per category. A block maps to a category by kind, except
text, which is split by the role that carries it:

| category | blocks | measured |
| --- | --- | --- |
| prompts | `text` on a user message | 200 |
| replies | `text` on any other role | 1,101 |
| thinking | `thinking` | 1,661 |
| calls | `tool_call` | 2,450 |
| results | `tool_result` and `error` | 2,449 |
| patches | `patch` | 114 |
| summaries | `summary` | 115 |
| notes | messages with no content at all | 3,544 |

An unknown block type is always shown, so a new atuin block kind cannot vanish
silently.

Errors join results because they are failed command output, and the existing
renderer already prints them as a result with a `!` prefix.

## Filtering while building

`buildDetail` already walks messages and appends one block per message or one
block per collapsed run. The mask plugs in there:

```
for each message:
    visible := blocks whose category is selected
    if notes is off and message has no content: skip
    if len(visible) == 0: skip          <- no hollow meta rows
    append one block (or one collapsed run when notes is on)
```

The meta second line is built from the visible kinds only. A message whose text
survives but whose tool call does not reads `text - 265 tokens`, not
`text+tool - 265 tokens`.

## Counts

The denominator is the conversation's message count, which is atuin's
`message_count` when present. The numerator is the number of messages that
rendered at least one visible block. Notes are not part of either number,
because atuin does not count them; the panel lists their own count on the notes
row.

The panel counts come from one pass over the conversation when the panel opens,
not per frame: 11,520 messages take about 2 ms to count.

## Panel

```
 att · atuin threads
 show message types                                                   7976 msgs

 > [x] prompts       200   your messages
   [x] replies      1101   assistant prose
   [x] thinking     1661   hidden reasoning
   [x] tool calls   2450   Bash, edit, read
   [x] tool results 2450   command output
   [x] patches       114   diffs
   [x] summaries     115   session wrap-ups
   [x] notes        3544   attachments, titles, empty system
```

This is the frame captured at 80x24 after applying nothing, so the header
carries the plain count; with a filter in place it reads `1301 of 7976 msgs`.
The line stops one cell before the edge, as every other detail row does.

The panel replaces the body for the duration of the selection, which keeps the
draw path the same as the list view instead of compositing an overlay.

Presets:

| key | preset | categories |
| --- | --- | --- |
| 1 | everything | all, notes on |
| 2 | conversation | prompts, replies |
| 3 | tools | calls, results, patches |
| 4 | prompts only | prompts |

`n` selects the conversation preset as a shortcut for the common case, so it
does the same as `2`.

Applying happens on enter rather than on every toggle: rebuilding 11,520
messages takes about 140 ms, and a background rebuild per space press would
flicker. The panel keeps a draft mask and the model keeps the applied one.

## Cache key, anchor and background

The applied mask is part of the render cache key, next to the width. A
background result is accepted only when both its width and its mask still match
the model, which extends the existing stale check by one field.

`detailBlock` gains the index of its first message, so an applied selection can
keep the reading position: take the message index at the top of the viewport
before rebuilding, then scroll to the first block whose range contains or
follows it. When that message is hidden by the new selection, the view falls
back to that next block, and to the top when nothing follows.

## Status and hints

- detail status: `01a18c27 · 1301 of 7976 msgs · line 3 of 4200`; without a
  filter the count collapses to `01a18c27 · 7976 msgs · line 3 of 55527`
- detail hints gain `t types`
- the head band keeps atuin's message count, which now matches the status
- when nothing renders, the body says `no messages of the selected types`

## Decisions

- Blocks, not roles, because "show replies but not tool calls" is the useful
  reading mode and roles cannot express it.
- Hide entirely instead of folding each hidden message into a summary line: the
  point is to make a 55,527-line thread readable, and a fold keeps 5,212
  assistant messages at two lines each.
- Notes as a switch over empty messages rather than a block category, because
  those messages have no blocks.
- Unknown kinds always visible, so a new atuin type degrades to visible rather
  than invisible.
- The selection is a session setting, not a per-thread one; the configuration
  file that makes it permanent is `atuin-threads-tui-3ab2`.
- No persistence in this change, no list view changes, no per-thread override.

## Risks

- Every category unselected is a reachable state. It gets a designed empty body
  and the panel can always be reopened with `t`.
- The anchored message may be hidden by the new selection. The fallback rule
  above keeps the jump small.
- Counting requires a full pass over messages when the panel opens. At 2 ms it
  is not worth caching.
- A filter that hides everything except notes leaves mostly empty rows; the
  status count makes that visible.

## Testing

- category mapping for every block kind, including text by role and an unknown
  kind
- filtering: mixed message keeps its visible part, tool-only message vanishes,
  notes hidden removes runs, notes shown keeps one row per run
- meta names visible kinds only, tokens stay on the message
- counts: numerator, denominator, notes excluded, nothing selected
- panel: open, toggle, apply, cancel, presets, selection across threads
- anchor: top message survives a rebuild, hidden anchor falls through
- cache: mask change rebuilds once, stale result for an old mask is dropped
- status and hints text, empty body
- frames at 120x36 and 80x24 with the conversation preset
