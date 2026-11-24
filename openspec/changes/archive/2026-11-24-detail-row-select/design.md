## Context

The conversation is rendered once into `detailRendering`: a head (title,
harness line, cwd) and one block per message, each block holding finished lines
that already contain SGR codes. `lineInfo(line)` maps a line to its message, and
`detailBlk` knows where a message starts and ends. Selection today only swaps the
first margin cell of a message.

The bean reports that this is unclear: the marker sits on one line of a message
that may be a hundred lines tall, and a range cannot express "these three".

## Shape

```
 threads -> buildDetail (once) -> blocks of styled lines
                                        |
                          view slot: ~20 visible lines
                                        |
                     paint: stripe + background tint per line
                                        |
                                    terminal
```

Painting never touches the cache. Moving the cursor recolours the visible lines
and nothing else.

## Before and after

Before, at 120 columns (one cell tells you everything):

```
 |  11:53 - user
 |  /openspec-explore ./scripts/dev-be.sh
 |  11:54 - asst
 |> Bash  Inspect dev scripts and README
 |    $ ls scripts
 |  11:54 - tool
 |  1278 .../MINT_FLOW_parametrische_specificatie.txt
```

After:

```
 |> 11:53 - user
 |> /openspec-explore ./scripts/dev-be.sh
 |# 11:54 - asst
 |# Bash  Inspect dev scripts and README
 |#   $ ls scripts
 |  11:54 - tool
 |  1278 .../MINT_FLOW_parametrische_specificatie.txt
```

The first two lines sit on the cursor tint with a `>` stripe on every line. The
next three carry the marked tint and a `#` stripe. The last two are unmarked and
alternate between the even and odd tints. The stripe column is the existing one
cell margin; the tint covers the body cells only.

## Colours

| Row state | Stripe | Background |
| --------- | ------ | ---------- |
| cursor    | `>` pink 205 bold | 240 |
| marked    | `#` pink 205 | 238 |
| even message | space | 234 |
| odd message  | space | 236 |
| thread head  | space | none |

Tints are explicit greys, not the terminal default, because att does not know
what the terminal background is. Light terminals are a known gap and are left
for a later change, as the user asked.

## Painting a line

`paintBackground(line, bg, width string)` pads the line to the body width with
spaces, replaces every inner `\x1b[0m` with `\x1b[0m` plus the background, and
wraps the whole line in the background, closing with a reset. Re-applying after
each inner reset is what keeps a diff line's green or a dim rail from cutting the
tint in half.

Stripes replace the first cell, which is already a space for content rows, so no
column is lost and the body width stays 92 at 120 columns.

## Marked set

`marked map[int]bool` keyed by message index, plus the existing cursor index into
`detail.selectable()`. Space toggles, `esc` clears, closing the thread or changing
the type filter clears. Copying and writing walk `selectable()` in order and take
the marked ones, so document order comes for free.

`y` with nothing marked sets a notice (`mark a message with space first`) and
touches neither clipboard nor terminal.

## Key map in selection mode

| Key | Action |
| --- | ------ |
| `j` `k` | one message |
| `J` `K` | ten messages |
| `g` `G` | first or last selectable message |
| `space` | mark or unmark the message under the cursor |
| `y` | copy the marked messages |
| `w` | write the marked messages to a file |
| `esc` | leave, dropping the marks |

## Decisions

- **Set instead of range.** A range forces every message between anchor and
  cursor to be included; the bean asks for a few. Rejected: range plus a
  subtract mode, which is two mental models.
- **Tint every message, always.** Message boundaries are useful while simply
  reading, not only while selecting.
- **Stripe on every line of a message.** One stripe on the first line is exactly
  the problem the bean reports.
- **Paint, do not rebuild.** A rebuild costs about 150 ms on the largest local
  thread; painting roughly 20 lines costs microseconds.
- **Pink stripes, grey tints.** If the terminal drops the background (no 256
  colours, or `NO_COLOR`), the stripes still say what is marked.
- **No outline mode.** A row-per-message view was considered and rejected: it
  hides the content that the selection exists to export.

## Risks

- Backgrounds may be invisible on terminals that report few colours; the
  stripes carry the meaning there.
- A very long line is padded with spaces to the body width, so a wide thread
  produces slightly larger frames. Bounded by 20 rows.
- `strings.ReplaceAll` over each visible line is cheap, but pathological lines
  with thousands of SGR codes cost more; still bounded by the viewport.

## Test plan

- Painter: plain line, line with inner reset, line shorter than the body, line
  exactly the body width, empty string.
- Zebra: consecutive messages alternate, a message taller than the viewport
  keeps one tint, the head carries no tint.
- Marking: space marks, space again unmarks, marks survive cursor movement,
  leaving drops them, closing the thread drops them, changing types drops them.
- Drawing: `>` on every visible line of the cursor message, `#` on every visible
  line of a marked message, cursor wins over marked, stripe inside the margin.
- Copy: three of five marked copies exactly those three in document order; `y`
  with nothing marked leaves the clipboard alone.
- Performance: moving the cursor does not increase the build counter.
