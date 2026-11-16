# Tasks: detail frame width

## 1. Layout arithmetic

- [x] 1.1 Add `detailMargin = 1` and change `newDetailLayout` so it returns a
      `prose` width next to `body` and `meta`: `body = width - 28`,
      `prose = min(80, body)`, and below 48 columns `meta = 0` with
      `body = width - 2`. Verify with a table test that 170 gives body 142,
      120 gives 92, 80 gives 52, 48 gives 20 with a rail, and 47 gives 45
      without one.
- [x] 1.2 Keep `detailBodyMin` from shrinking the body below 1 cell and cover
      width 2 with a test that does not panic.

## 2. Margins and the rail

- [x] 2.1 Prefix every content row with one space and reserve the last cell:
      the title, the harness line, the cwd line, and each message row. Verify
      that a rendered frame at 120 columns has an empty first and last cell on
      every content row.
- [x] 2.2 Place the meta rail at `margin + body + gap` so it ends one cell
      before the window edge, with its text left-aligned. Verify that
      `11:53 - asst` and `tool - 265 tokens` start at the same column.
- [x] 2.3 Truncate the header lines at the body width instead of 80. Verify
      with a 120-column thread whose title is 100 cells long.

## 3. Per-block widths

- [x] 3.1 Give `blockLines` a prose width and a code width, and remove the
      80-cell clamp. Use prose for `text`, `thinking`, and `summary`, and the
      code width for tool calls, tool results, patches, and error blocks.
      Verify that a 92-cell command at 120 columns stays on one row while a
      prose line still stops at 80.
- [x] 3.2 Confirm the prose width follows the body at narrow widths: at 80
      columns prose wraps at 52. Verify with a long paragraph rendered at 80.

## 4. Tests

- [x] 4.1 Update `detail_test.go` expectations that assumed a 120-column body
      of 80 and a rail at cell 82; keep the existing scenario coverage.
- [x] 4.2 Add layout tests at 170, 120, 80, 48, and 47 columns asserting the
      body width, the prose width, the rail position, and the presence of the
      margins.

## 5. Verification

- [x] 5.1 Run `go test ./...` and `bash scripts/gate.sh`; both pass and
      coverage stays above the thresholds.
- [x] 5.2 Run `nix flake check` and `nix build`, then capture the detail view
      at 170x40, 120x36, and 80x24 and compare each against the mockups in
      design.md.
