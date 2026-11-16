# Design: thread-table-list

## Current row

```
> 27c12788  Explore skull2-f4m3 change scope  ·  opencode  ·  337 msgs  ·  12s
  2cf53ec4  OpenSpec and Beans project bootstrap  ·  opencode  ·  148 msgs  ·  29s
  da8da47e  Testversie op deze computer  ·  claude-code  ·  831 msgs  ·  4m
```

Every field is separated by two spaces and a middle dot, so the start of each
field depends on the length of the one before it. The harness name alone shifts
the message count by six cells between `opencode` and `claude-code`. Sorting
down a column with the eye is impossible, and the relative age (`4m`, `3d`)
says nothing about when a thread ran once it leaves the current day.

## Table

```
  id        date          ag  title                                                             msgs  dir
  ────────  ────────────  ──  ────────────────────────────────────────────────────────────────  ────  ──────────────────
> 27c12788  16 Nov 11:06  oc  Explore skull2-f4m3 change scope                                   337  atuin-threads-tui
  2cf53ec4  16 Nov 11:06  oc  OpenSpec and Beans project bootstrap                               193  atuin-threads-tui
  da8da47e  16 Nov 10:42  cc  Testversie op deze computer                                        831  atuin-threads-tui
  75f570d7  16 Nov 09:58  cc  Eigenaren-synchronisatie en VVE-instellingen                       969  vve-beheer
  4b160473  15 Nov 22:31  oc  (untitled)                                                          90  docs
```

At 80 columns the directory column drops and the title takes the space:

```
  id        date          ag  title                                         msgs
  ────────  ────────────  ──  ────────────────────────────────────────────  ────
> 27c12788  16 Nov 11:06  oc  Explore skull2-f4m3 change scope               337
  2cf53ec4  16 Nov 11:06  oc  OpenSpec and Beans project bootstrap           193
  da8da47e  16 Nov 10:42  cc  Testversie op deze computer                    831
  75f570d7  16 Nov 09:58  cc  Eigenaren-synchronisatie en VVE-instellingen   969
  4b160473  15 Nov 22:31  oc  (untitled)                                      90
```

## Columns

| Column | Width | Alignment | Content                                    |
| ------ | ----- | --------- | ------------------------------------------ |
| marker | 2     | left      | `> ` on the selected row, two spaces else  |
| id     | 8     | left      | last 8 characters of the atuin id           |
| date   | 12    | left      | `UpdatedAt` in local time, `02 Jan 15:04`  |
| ag     | 2     | left      | `oc`, `cc`, or the first two harness letters |
| title  | rest  | left      | title, `(untitled)` when empty             |
| msgs   | 4     | right     | message count                              |
| dir    | 18    | left      | basename of `CWD`, cut from the left       |

Two spaces separate every column. The marker column keeps the selection visible
without colour, so a monochrome terminal and the plain-text tests still show
which row is active.

The last 8 characters of the atuin id are the only useful part: the id is a
ULID, so threads created in the same hour share their leading characters. The
short form is already used by the detail status line, so the two views name a
thread the same way.

`02 Jan 15:04` is 12 cells in every locale and never wraps. The year is left
out because threads older than a year are rare, and the update date is what the
list is ordered by. A missing timestamp renders as `-`.

The message count is right aligned so the digits line up under the `msgs`
header. Width 4 covers 9999 messages; the final row clamp keeps a longer number
from pushing the row past the terminal width.

The directory is shortened to its basename because the parent chain is shared
by most threads and the project name is what tells them apart. Cutting from the
left keeps the end of the name, which is usually the distinctive part. An 18
cell column holds the common project names without cutting.

## Degrading at narrow widths

Columns drop in this order, each time the title would fall below 20 cells:

| Terminal width | Dropped  | Title cells |
| -------------- | -------- | ----------- |
| 100 and up     | none     | width - 56  |
| 56 to 99       | `dir`    | width - 36  |
| 42 to 55       | `date`   | width - 22  |
| 38 to 41       | `ag`     | width - 18  |
| below 38       | `msgs`   | width - 12  |

`id` and `title` survive to the smallest width. The header and the rule always
match the column set of the rows, so nothing shifts when the terminal resizes.

## Header and rule

The header is dim, lowercase, and uses the same column widths as the rows. A
dim rule of `─` cells under it closes the header band. Together they cost two
body rows, which the layout accounts for: rows visible in the list are
`bodyHeight - 2`. Both stay at the top while the rows scroll.

No vertical borders and no per-cell background. The alignment and the dim
header carry the structure; borders at this density would add noise.

## Width arithmetic

```
fixed = 2 (marker) + 8 + 2 + 12 + 2 + 2 + 2 + 18 (dir) + 2 + 4 = 56
title = terminal width - fixed, with dir hidden when width < 100
```

At 120 columns the title gets 64 cells, at 80 columns 44. Both were checked
against the real data: the longest title on this machine is 61 cells.

Cell counts are display cells, not runes. Titles in the history already contain
wide and non-ASCII characters, so truncation and padding go through
`charmbracelet/x/ansi`, which measures grapheme clusters and wide characters.

## Tests

- Column layout at 120, 100, 99, 80, and 40 columns: every row and the header
  have the same display width, and the dropped columns are the ones the order
  predicts.
- Agent mapping for `opencode`, `claude-code`, an unknown harness, and an empty
  harness.
- Date rendering from a fixed `UpdatedAt`, and `-` for a zero time.
- Directory basename, tail truncation with a long path, and `-` for an empty
  `CWD`.
- Selection scrolling: with `bodyHeight - 2` rows visible, moving past the last
  visible row keeps the selection on screen.

## Out of scope

Sorting by column, toggling columns by key, a column picker, changes to the
detail view, and README or configuration work (milestone 03).
