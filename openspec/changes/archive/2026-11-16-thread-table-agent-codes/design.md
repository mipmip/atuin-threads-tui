## Context

The list is a table with the columns id, date, agent, message count, title, and
directory. The agent column is measured from the codes in view and clamped
between 2 and 5 cells. The codes are `oc`, `cc`, and `co`.

## Problem

`oc`, `cc`, and `co` are abbreviations of the harness names, invented here.
They save nothing the column does not already save: the column is only as wide
as the code in view, so at the current data the column is 2 cells wide. Two
letters that need a decoder cost more in confusion than four to six letters
cost in cells, and the four-letter code is only paid when an opencode thread is
on screen.

## Decision

| Harness       | Code     |
| ------------- | -------- |
| `claude-code` | `claude` |
| `opencode`    | `ocode`  |
| `codex`       | `codex`  |
| `pi`          | `pi`     |
| other name    | first 6 runes, lowercased |
| empty         | `--`     |

`ocode` rather than `opencode` because the column is six cells at most; the
name is still recognisable without the `pen`. `claude` needs six cells, so the
maximum moves from 5 to 6.

The header reads `agent` once the column is 6 cells, and `ag` when every code in
view is shorter. With `claude` or `ocode` in view it is `agent`; a list with
only `pi` threads keeps `ag`.

## Layout

At 120 columns with a 4-cell message count, a 6-cell agent column, a 48-cell
title, and a 24-cell directory, taken from a live run at 120x24:

```
  id        date          agent   msgs  title                                             dir
  ────────  ────────────  ──────  ────  ────────────────────────────────────────────────  ────────────────────────
> 2cf53ec4  16 Nov 13:25  ocode    325  OpenSpec and Beans project bootstrap              atuin-threads-tui
  fd025611  16 Nov 13:03  claude    47  dev-fe.sh wrapper script                          chatcms-app
  3c07dedc  16 Nov 11:46  claude    51  Tuimock tools check                               tui-design-agent-exp2
  b86dc501  16 Nov 11:35  claude     1  (untitled)                                        vve-go
```

Four cells wider than the previous frame, taken from the slack the table leaves
on a wide terminal.

## Degradation

Unchanged. When the table does not fit, the title shrinks to 20 cells, then the
directory shrinks to 8, and only then are directory, date, agent, and message
count dropped in that order. A 6-cell agent column therefore disappears at the
same point a 5-cell one did, and rows stay aligned in every case because widths
are computed once per frame.

## Out of scope

- A legend or help screen that explains the codes.
- Per-agent colour.
- Reading the harness name from anywhere but the thread itself.
