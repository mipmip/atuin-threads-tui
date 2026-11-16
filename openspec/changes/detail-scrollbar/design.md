# Design

## Context

The conversation rendering is cached per thread and width (`detailRendering`),
so the frame already knows `total` lines, the viewport start
(`detailScroll`), and the height it draws (`bodyHeight`). `detailBody` slices
the cache and `View` pads the result up to the body height with empty lines.

After `detail-frame-width` the frame is one left margin, the body, a two-cell
gap, the 24-cell meta rail, and one right margin. That margin is the only cell
in the frame that carries nothing.

See proposal.md for why the bar is wanted.

## Goals / Non-Goals

Goals:

- Show length and position in one cell, without moving any text.
- Keep the body width independent of whether the bar is visible.
- Cost no measurable time per frame.

Non-Goals:

- Mouse dragging or clicking the bar. The app is keyboard only.
- A bar in the list view. It fits the same helper, but it is a separate
  decision.
- Changing the status line. `line X of Y` stays.
- Half-cell thumb precision with `▀` / `▄`. At 55,527 lines one cell already
  covers 1,683 lines.

## Decisions

### The bar occupies the right margin

The body keeps the width `detail-frame-width` gives it, and the bar is drawn in
the last column. Nothing in the layout arithmetic changes. The alternative,
reserving a new column and shrinking the body by one, buys nothing: the margin
is empty anyway.

The column stays reserved when the bar is hidden, so switching between a short
and a long thread never changes the body width or the cache key.

### Track and thumb, two greys and one pink

```
 track  │  U+2502  colour 238
 thumb  █  U+2588  colour 205
```

The track is the dimmest thing on screen and reads as a frame edge rather than
as content. The thumb uses the title pink, the same accent the list already
uses for emphasis, so the app has one accent colour rather than two.

Thumb-only was the alternative: calmest, but it only answers "where", not "how
long", which is half of what the bean asks. A density strip was rejected
because 33 cells over 55,527 lines means every cell averages about 170
messages.

### The maths

```
track = bodyHeight
thumb = clamp(round(track * track / total), 1, track)
top   = round(scroll * (track - thumb) / (total - track))
```

Three details matter:

- `total - track` can be zero or negative; the bar is hidden then, so the
  division never runs.
- Rounding is half-up on integers; with the thumb at one cell the position
  still moves on every page.
- `clampScroll` already keeps `scroll` within `[0, total-track]`, so `top` never
  leaves the track and the thumb touches the bottom row exactly at the last
  line.

The result is two integers per frame, from values the model already holds. No
line is measured, and nothing new is cached.

### Drawing happens after padding

`View` pads the body with empty lines to reach `bodyHeight`, so the bar is drawn
there, over exactly `bodyHeight` rows, rather than inside `detailBody`. Each row
is truncated or padded to `width - 1` cells with `ansi.Truncate` and
`ansi.StringWidth`, then the bar cell is appended. That keeps a wide character
from pushing the bar off the edge and gives the padded empty rows their track.

### No bar while rendering

While the background render runs, `total` is unknown and the body shows
`rendering N messages...`. The bar stays hidden until the first render lands.

## Layout at 120 columns

```
| The script now at `scripts/dev-fe.sh` is the old `dev.sh`,    11:53 - asst          #
| which is the backend script.                                  text - 675 tokens     |
|                                                                                     |
| Bash  Inspect dev scripts and README                          11:53 - asst          |
|   $ cd /home/pim/gh.quiqr/chatcms-app; ls scripts             tool - 265 tokens     |
|                                                                                     |
```

The rail is left-aligned inside its 24 cells, so there are ten empty cells
between `tool - 265 tokens` and the bar in most rows.

## Risks / Trade-offs

- Colour 238 on a very light terminal is a dark line, on a very dark terminal it
  is a faint one. The thumb carries the position in both cases → acceptable, and
  `NO_COLOR` terminals still get two distinct glyphs.
- A one-cell thumb at 55,527 lines jumps by whole screenfuls → correct, and it
  is the honest resolution of a 1,600-screen conversation.
- Two extra width measurements per row cost about a microsecond for 33 rows →
  negligible against the current 8.6 µs frame.

## Migration Plan

None. No data, configuration, or interface changes.

## Open Questions

None.
