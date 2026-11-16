# Design

## Context

Measured on this machine, 120x40, with the largest thread in the local database:

| Thread | Messages | Rendered lines | One render |
|--------|---------:|---------------:|-----------:|
| `01a18c27…` | 11,520 | 55,527 | 344 ms |
| `fd025611` | 224 | 1,068 | 12.9 ms |

A keystroke in the detail view calls `detailLines()` three times: from
`clampScroll`, from `body`, and from `status`. So 11,520 messages cost about
one second per key, and 224 messages about 40 ms.

The atuin CLI is not the problem: `atuin threads show` returns a 10 MB
conversation in 25 ms, and it already runs in a `tea.Cmd`.

CPU profile of one render: `splitWidth` 32% (an `ansi.StringWidth` call per
rune), JSON decode 20%, lipgloss style resolution 15%, the rest allocation.

## Current flow

```
key j
  Update → handleDetailKey        detailScroll++
  clampScroll                     len(detailLines())      full render
View
  body → detailBody               detailLines()           full render, then slice
  status                          len(detailLines())      full render
```

Only 20 to 40 of the rendered lines are ever drawn.

## New flow

```
key j
  Update → handleDetailKey        detailScroll++
  clampScroll                     cached total           O(1)
View
  body → detailBody               cache.slice(scroll, height)   O(visible blocks)
  status                          cached total           O(1)
```

The conversation is rendered once, when it opens, and again when the width
changes. Scrolling never renders a message.

## Data structures

```go
type detailBlock struct {
	start int      // absolute line of the first line of this block
	lines []string
}

type detailRendering struct {
	width  int
	built  bool
	head   []string        // title, harness line, working directory, blank
	blocks []detailBlock   // one per message, or one per collapsed run
	total  int
}
```

`head` holds the fixed preamble. Each block owns its own trailing blank line,
so a block can be sliced without looking at its neighbours.

`slice(from, count)` binary-searches the first block that covers `from` and
copies only the lines inside `[from, from+count)`, skipping the gaps between
blocks. The cost is proportional to the visible blocks, not the conversation.

The per-block layout is also what a future message selection needs: a row in
the block table maps to a `atuin.Message` index.

## Background rendering

```go
case conversationMsg:
	if len(msg.conversation.Messages) <= detailSyncMessages {
		m.detail = m.buildDetail()        // 500 messages ≈ 30 ms
	} else {
		m.rendering = true
		return m, m.renderDetail()        // tea.Cmd
	}
```

`detailSyncMessages = 500`. Below it, the render happens inside `Update` and
the first frame is already complete, so there is no flicker. Above it, the
view shows a rendering state:

```
att · atuin threads
OpenSpec and Beans project bootstrap
claude-code · claude-opus-5-5 · 01a18c27 · 11520 msgs · 2m ago
/home/pim/gh.mipmip/atuin-threads-tui

rendering 11520 messages...
```

with the status line reading `01a18c27 · 11520 messages · rendering`.

A result is accepted only when its width still matches the current terminal
width; a render that finishes after a resize is dropped, and the resize starts
a new one.

## Resize

`tea.WindowSizeMsg` on a width change invalidates the rendering and rebuilds it
at the new width, synchronously for small threads and in the background for
large ones. The width is part of the cache key, so a stale rendering is never
drawn.

## Faster measurement

Two helpers dominate the render:

- `splitWidth` called `ansi.StringWidth(string(r))` for every rune. It now uses
  `runewidth.RuneWidth(r)`, which is a table lookup. The result is the same for
  CJK, combining marks, and emoji, and a lone rune is exactly what the loop
  measures; `ansi.StringWidth` on a single rune was grapheme segmentation on a
  string of one.
- `wrap` measured the whole growing line for every word, which is quadratic in
  the line length. It now tracks the width of the line and of each chunk it
  appends.

`ansi.StringWidth` stays where a whole string has to be measured at once
(`preserve`, `truncate`, `padRight`), because there grapheme clusters matter.

## Decisions and alternatives

**Per-block cache instead of one flat `[]string`.** A flat slice would need
`strings.Join` over 55,527 lines per frame, and it leaves no place to hang
message identity on. The block table costs one extra allocation per message and
gives both the viewport slice and the future selection index.

**Nothing is rendered lazily per message.** Rendering only the visible messages
would make opening O(visible), but the offset table needs the height of every
message, so the conversation has to be measured once anyway. Rendering once and
caching is simpler and bounded.

**No progressive chunking.** The background render of the largest thread takes
344 ms, which the rendering state covers. Chunking would add ordering rules and
partial-cache invalidation for no visible benefit.

**The header is cached too.** Its `2m ago` fragment is therefore frozen until
the next rebuild (open, resize, or reload). That is a deliberate trade: a
per-frame header would need its own line budget and would shift every block
offset when the harness or model string changes. The list view still shows live
ages.

**Memory.** A 55,527-line rendering is roughly 10 MB of strings, held only
while that thread is open. Dropping the cache on `esc` returns it.

## Measured result

Same conversation, 11,520 messages, 120 columns:

| | before | after |
|---|---:|---:|
| `buildDetail` | 344 ms | 147 ms |
| one `View` frame | 566 ms | 8.6 µs |
| slice of 38 lines | 55,527 line join | 0.9 µs |

Live, in tmux, opening the same thread and sending ten `j` keys 50 ms apart:
the old binary had not passed line 3 after 120 seconds, the new one reached
line 11 in 0.54 s, which is the send interval plus nothing.

## Batched input

The A/B run first showed zero movement: tmux delivered fifty `j` characters in
one read, and bubbletea packs consecutive printable runes into a single
`tea.KeyMsg`. `handleKey` matched `msg.String()` against `"j"`, saw `"jjjj…"`,
and dropped the batch. Outside the filter, a key message with several runes is
now applied one rune at a time. In the filter the runes are still inserted as
text, so pasting a query keeps working.

## Risks

- The `now` clock is captured at build time, so the header age is as old as the
  rendering. Documented above.
- A resize during a background render drops the in-flight result and starts
  again; two rapid resizes render twice in the worst case.
- `runewidth.RuneWidth` and `ansi.StringWidth` could disagree on a rune the
  former does not know. The existing CJK, combining mark, and wide character
  tests cover the cases we can construct.

## Testing

- A render counter on the model: after opening a thread, 20 keystrokes leave the
  counter at 1.
- A cache test that the visible slice equals the matching range of the flat
  rendering, including a viewport that starts inside a block.
- A resize test that the counter grows by one and the body column follows the
  new width.
- A large-thread test that the model reports `rendering` before the command's
  message arrives, and draws the conversation after it.
- A stale-result test: a rendering for another width is ignored.
- Existing detail tests keep calling `detailLines()`, which stays as the pure
  builder used by the cache.
