# Design: detail-wrap-linear

## Context

`buildDetail` is measured on two real conversations:

| Thread          | Messages | build at 120 | build at 80 |
| --------------- | -------- | ------------ | ----------- |
| `01a18c27`      | 11,520   | 147 ms       | 200 ms      |
| `2cf53ec4`      | 609      | 1,410 ms     | 2,610 ms    |

The 609-message thread is slower than the 11,520-message one, so the cost is not
driven by message count. A CPU profile attributes 88% of that build to
`preserve` → `wrapPreserving` → `splitWidth`, 0% to decoding and styling.

The offending loop:

```go
for remaining := line; remaining != ""; {
    chunks := splitWidth(remaining, limit)
    out = append(out, prefix+chunks[0])
    remaining = strings.TrimPrefix(remaining, chunks[0])
}
```

`splitWidth` scans the whole remainder and, because it starts with
`ansi.StringWidth(s)`, scans it twice. `chunks[0]` is then removed with a copy of
everything that is left. A 50 KB line at 92 cells therefore costs
500 × 50 KB of scanning and copying instead of 50 KB once.

## Decisions

### Cut with one function

`cutWidth(s, width)` returns the longest prefix that fits and the rest, walking
runes until the width is used up. It is the only place that knows about rune
boundaries and wide characters.

### Walk the line once

`wrapPreserving` keeps a `rest` cursor and calls `cutWidth` on the remainder
only, never on the full line. The first chunk uses the body width, every later
chunk uses the width minus the two-cell hang, which is the same output as before.

### Keep the public shape

`splitWidth` keeps returning `[]string` and `preserve` keeps its early return for
lines that fit, so no caller changes. The guards for a width below one stay.

## Risks

- A rune wider than the limit (a wide CJK character at width 1) must still be
  emitted, or the loop would not advance. `cutWidth` takes at least one rune per
  call.
- The trailing-hang output must stay byte-identical for the existing detail
  tests; they compare rendered frames.

## Measured result

| Measure                | Before   | After  |
| ---------------------- | -------- | ------ |
| `2cf53ec4` at 120      | 1,410 ms | 32 ms  |
| `2cf53ec4` at 80       | 2,610 ms | 32 ms  |
| `01a18c27` at 120      | 147 ms   | 138 ms |

The remaining 32 ms is ordinary work: styling, decoding tool input, and
measuring every block once.

## Out of scope

- Rendering only the visible blocks (the offset table needs every height).
- Markdown-aware prose reflow, which is a separate open question.
