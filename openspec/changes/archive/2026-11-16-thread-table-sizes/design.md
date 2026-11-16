# Design: thread table sizes

## Audit of the current table

The table that shipped in `thread-table-list` aligns every column, but two
widths ignore the content:

- The title column is `width - fixed`, so at 120 columns it reserves 64 cells
  whether the longest title needs 36 or 200. Rows with short titles end in a
  wide band of empty cells, and the eye has to cross it to reach the directory.
- The directory column is a fixed 18 cells. A basename longer than 18 is cut
  from the left (`…count-104144963194-main`) while the empty title padding
  beside it goes unused.
- The message count sits after the title, so the number that belongs with the
  agent tag is pushed to the far right, next to a column whose length varies.

The header rule spans the terminal, not the table, which made the empty band
look intentional.

## After

```
  id        date          ag  msgs  title                                             dir
  ────────  ────────────  ──  ────  ────────────────────────────────────────────────  ────────────────────────
> 2cf53ec4  16 Nov 13:08  oc   278  OpenSpec and Beans project bootstrap              atuin-threads-tui
  fd025611  16 Nov 13:03  cc    47  dev-fe.sh wrapper script                          chatcms-app
  3c07dedc  16 Nov 11:46  cc    51  Tuimock tools check                               tui-design-agent-exp2
  b86dc501  16 Nov 11:35  cc     1  (untitled)                                        vve-go
```

At 80 columns the title gives way first, so the directory stays readable:

```
  id        date          ag  msgs  title                 dir
  ────────  ────────────  ──  ────  ────────────────────  ──────────────────────
> 2cf53ec4  16 Nov 13:08  oc   280  OpenSpec and Beans …  atuin-threads-tui
  fd025611  16 Nov 13:03  cc    47  dev-fe.sh wrapper s…  chatcms-app
  3c07dedc  16 Nov 11:46  cc    51  Tuimock tools check   tui-design-agent-exp2
```

## Column widths

| Column | Width                     | Rule                                     |
|--------|---------------------------|------------------------------------------|
| marker | 2                         | `> ` when selected, spaces otherwise     |
| id     | 8                         | last 8 characters of the atuin id        |
| date   | 12                        | `02 Jan 15:04`, local time               |
| agent  | 2 to 5                    | widest code in view, header grows with it|
| msgs   | 4 or more                 | right aligned, digits of the largest     |
| title  | 20 to 48                  | longest title in view, capped            |
| dir    | 8 to 24                   | longest basename in view, capped         |

The table is the sum of its columns plus a two cell gutter between them. It is
not padded to the terminal width: an empty band on the right is more honest
than a title column that is wider than any title.

## Fitting and degradation

1. Take the measured widths, clamped to the ranges above.
2. While the table is too wide, take a cell from the title, down to 20.
3. Then take a cell from the directory, down to 8.
4. Then drop whole columns in this order: directory, date, agent, message count.
5. Whatever remains goes to the title, at least one cell.

The order follows what a reader needs: the title identifies a thread, the
directory only disambiguates threads with the same title. When space runs out
there is nothing left to identify, so the marker and id stay.

## Agent code

`opencode` is `oc`, `claude-code` is `cc`, `codex` is `co`, and an empty
harness is `--`. Any other harness shows its own name in lower case, cut to
five characters, so a fork such as `opencode-fork` reads `openc`. The column
follows the widest code: two cells for the aliases, five when a longer name is
in view, and the header reads `ag` or `agent` accordingly.

## Out of scope

- Sorting, column toggles, and a resizable table.
- A wider palette or a header style beyond the existing muted colour.
- The detail view.
