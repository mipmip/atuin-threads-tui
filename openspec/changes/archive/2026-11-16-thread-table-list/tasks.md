# Tasks: thread-table-list

## 1. Layout

- [x] 1.1 Add a `listLayout` type that derives column widths and the title
      width from the terminal width and the widest message count
- [x] 1.2 Drop `dir`, `date`, `ag`, then `msgs` when the title would fall below
      20 cells

## 2. Rendering

- [x] 2.1 Render the header row and the rule from the same layout
- [x] 2.2 Render rows with the marker, short id, date, agent tag, title, right
      aligned message count, and directory basename
- [x] 2.3 Map the harness to `oc` and `cc`, with a fallback for other harnesses
- [x] 2.4 Add display width aware truncation and tail truncation helpers
- [x] 2.5 Show the table header and rule in the list body, and keep them out of
      the loading, empty, and error states

## 3. Scrolling

- [x] 3.1 Count the header and rule as chrome in `clampCursor` and in the body
      slice so the selection stays visible

## 4. Tests

- [x] 4.1 Column widths and column dropping at 120, 100, 99, 80, and 40 columns
- [x] 4.2 Agent, date, and directory formatting, including the empty cases
- [x] 4.3 Selection scrolling with the two chrome rows

## 5. Verify

- [x] 5.1 `nix flake check` passes
- [x] 5.2 Capture the list at 120x36 and 80x24 and compare against the mockups
- [x] 5.3 Pin `subPackages` to `cmd/att` so a nested Go tree cannot change the
      build
