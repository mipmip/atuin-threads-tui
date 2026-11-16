# Detail frame width

## Context

The conversation view renders a body column and a meta rail per message. Today
`newDetailLayout` computes `body = min(80, width - 26)` and puts the rail
directly after a two-cell gap, so the block ends at column 106 in a 120-column
window and the remaining 14 cells stay empty. `blockLines` then clamps every
block, prose and tool output alike, to 80 cells.

Measured on the local thread with 11,520 messages: 84% of command lines and 84%
of result lines fit in 80 cells, 89% in 92, and 94% in 118. The body column
grows to 92 at 120 columns and 142 at 170, which recovers most of that.

## Goals / Non-Goals

Goals:

- Fill the window horizontally, with the meta rail at the right edge.
- Give code the width the window actually has, without changing the prose
  measure.
- Leave a one-cell margin around the detail content.

Non-goals:

- A configurable prose width. Prose stays at 80 cells; the knob is a later
  change.
- Wrapping code at the full window instead of the body column. The rail keeps
  its cells, so code never runs under it.
- Markdown awareness (fences, tables, lists) inside prose blocks. That is a
  separate bug.
- Margins in the list view.

## Decisions

### The rail is pinned to the right edge, its text left-aligned

The rail occupies the 24 cells before the right margin, so `11:53 - asst` and
`tool - 265 tokens` both start at the same column. Right-aligning the text
inside the rail was rejected: two lines of one message would start at
different columns, which reads as ragged.

### The body column is `width - 28`

One left margin, one right margin, a two-cell gap, and the 24-cell rail. At
120: body 92. At 170: body 142. At 80: body 52. The body column no longer has
a ceiling; only prose does.

### Prose keeps its own measure

`prose = min(80, body)`. Wide terminals keep a readable line length, and at 80
columns the prose simply uses the narrowed body.

### Code wraps at the body column

`blockLines` receives two widths instead of one: `prose` for `text`,
`thinking`, and `summary`, and `body` for tool calls, results, patches, and
error blocks. The 80-cell clamp inside `blockLines` goes away.

### The rail drops below 48 columns

The body needs at least 20 cells (`detailBodyMin`), so the rail disappears when
`width - 28 < 20`. The body then takes `width - 2`. Before this change the
threshold was 45 columns; it moves to 47.

### Margins apply to content only

The thread title, the harness line, the cwd, and every message row are
indented by one cell. The app header (`att · atuin threads`), the status line,
and the key hints stay flush, because they are shared chrome with the list
view.

## Layout at 120 columns

```
| The script now at `scripts/dev-fe.sh` is the old `dev.sh`, which is the backend       |
| script.                                                                              |
| So I think the work is:                                             11:53 - asst      |
|                                                                                      |
| Bash  Read config, list specs, scaffold change                      11:53 - asst      |
|   $ cd /home/pim/gh.quiqr/chatcms-app; cat config.yaml 2>/dev/null | head -60         |
|     ls openspec/specs | head -50                                     tool - 265 tokens|
```

At 80 columns the same message uses a 52-cell body:

```
| The script now at `scripts/dev-fe.sh` is the old                  11:53 - asst      |
| `dev.sh`, which is the backend script.                                               |
|                                                                                      |
| Bash  Read config, list specs, scaffold change                    11:53 - asst      |
|   $ cd /home/pim/gh.quiqr/chatcms-app; cat openspec/         tool - 265 tokens       |
|     config.yaml 2>/dev/null | head -60                                               |
```

## Risks / Trade-offs

- [Prose and rail drift apart on very wide terminals] -> Accepted; the rail is
  a page margin at the edge, and prose stays readable at 80 cells.
- [Two cells are lost on narrow terminals] -> At 80 columns prose goes from 54
  to 52 cells; below 47 columns the rail drops and the body takes `width - 2`,
  which is still better than today.
- [The status line and header stay flush while content is indented] -> Visible
  one-cell offset between chrome and content; accepted so both views share the
  same chrome.
- [A rail at the edge makes short meta lines float far from the text] ->
  Mitigated by keeping the rail left-aligned, so it reads as a column.

## Migration Plan

None. The rendering is rebuilt on open and on resize; there is no stored state
and no data format change.
