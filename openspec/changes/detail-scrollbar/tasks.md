# Tasks: detail-scrollbar

## 1. Maths

- [ ] 1.1 Add a `scrollbar` value with `thumb` and `top`, computed from
      `track`, `total`, and `scroll`, with the bar hidden when `total <= track`
- [ ] 1.2 Unit tests: hidden when the thread fits, thumb clamped to one cell,
      thumb about a tenth of the track for a ten-times taller thread, thumb at
      the top row for scroll 0, thumb ending at the bottom row for the last
      line, and no division when `total == track`

## 2. Drawing

- [ ] 2.1 Draw the bar in the right margin of the body in detail mode, over
      exactly `bodyHeight` rows, after padding, using `ansi.Truncate` and
      `ansi.StringWidth` to place the cell at `width - 1`
- [ ] 2.2 Track in colour 238 with `│`, thumb in colour 205 with `█`
- [ ] 2.3 Leave the column empty but reserved when the bar is hidden, and while
      the thread is being rendered
- [ ] 2.4 Confirm the bar does not appear in the list view, filter mode, or the
      rendering message

## 3. Tests

- [ ] 3.1 A long thread shows a thumb on the first body row at the top and a
      track on the other rows
- [ ] 3.2 After `G` the thumb ends on the last body row
- [ ] 3.3 A thread that fits shows no track or thumb, and its body width equals
      the long thread's body width
- [ ] 3.4 A wide character at the end of a row cannot push the bar off the edge
- [ ] 3.5 Existing detail tests still pass with the reserved column

## 4. Verify

- [ ] 4.1 `bash scripts/gate.sh`
- [ ] 4.2 `nix flake check` and `nix build`, then a live frame at 120x36 and
      80x24 with the bar visible mid-thread
