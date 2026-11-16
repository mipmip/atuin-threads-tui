# Proposal: detail-wrap-linear

## Why

Opening a conversation whose tool results contain very long single lines takes
seconds. Measured with the 609-message thread `2cf53ec4`: `buildDetail` takes
1.41 s at 120 columns and 2.61 s at 80 columns, even though the 11,520-message
thread `01a18c27` builds in 147 ms. The cost is quadratic in the length of one
line, so a few large results outweigh thousands of messages.

The profile puts 88% of the build inside `preserve` → `wrapPreserving` →
`splitWidth`. `wrapPreserving` removes the first chunk with
`strings.TrimPrefix(remaining, chunks[0])` and then splits the whole remainder
again, so a 50 KB result line with a 92-cell width is rescanned 500 times and
copied 500 times.

## What Changes

- Add a linear width cutter, `cutWidth`, that returns the head that fits and the
  untouched rest.
- Rewrite `wrapPreserving` to walk the line once with `cutWidth`, keeping the
  two-cell continuation hang.
- Rebuild `splitWidth` on `cutWidth`, so one function owns the rune-boundary
  arithmetic.
- Cover very long lines in the detail tests and record the before/after timing.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: rendering long lines must not cost quadratic time.

## Impact

- `internal/ui/detail.go` only.
- No change to output: the same lines, the same hanging indent.
- No dependency, data, or client change.
