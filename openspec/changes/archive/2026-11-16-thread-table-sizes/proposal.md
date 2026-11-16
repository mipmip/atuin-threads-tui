# Thread table sizes

## Why

The table fits the columns to the widest value for the id, date, agent and
message count, but two columns still have a fixed width. The title takes every
cell that is left over, so on a 120 column terminal it grows to 64 cells even
when the longest title is shorter, and the directory stays at 18 cells and cuts
names that would fit. The message count sits at the far right, separated from
the agent tag it belongs with.

## What Changes

- Order the columns as id, date, agent, message count, title, directory.
- Size the title column to the longest title in view, capped at 48 cells and
  never below 20.
- Size the directory column to the longest directory basename in view, capped
  at 24 cells and never below 8.
- Let the agent code grow to five characters; the column and its header follow
  the widest code.
- Let the table be narrower than the terminal instead of stretching the title.
- Degrade by shrinking the title first, then the directory, and only then drop
  the directory, the date, the agent and the message count.
- Keep the header, the rule, the comment-cap in the rule and the selected row
  marker as they are.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the table requirement gains column order, content-driven
  title and directory widths, and a five character agent code.

## Impact

`internal/ui/view.go` for the layout and `internal/ui/table_test.go` and
`internal/ui/app_test.go` for the expectations. No new dependency, no change to
the atuin client.
