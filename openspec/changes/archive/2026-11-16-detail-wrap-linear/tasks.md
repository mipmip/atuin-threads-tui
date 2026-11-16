# Tasks: detail-wrap-linear

## 1. Linear cutter

- [x] 1.1 Add `cutWidth(s, width) (head, rest string)` that walks runes once and
      returns the longest prefix that fits
- [x] 1.2 Test that it consumes at least one rune at width 1 and that CJK
      characters count as two cells

## 2. Rewrite the wrapping

- [x] 2.1 Rewrite `wrapPreserving` to walk the line with `cutWidth`, first chunk
      at the body width, later chunks at the width minus the two-cell hang
- [x] 2.2 Rebuild `splitWidth` on `cutWidth`
- [x] 2.3 Test a 100,000 character line: the same chunks as before, no growth in
      the number of scans

## 3. Verify

- [x] 3.1 `go test ./...` and `bash scripts/gate.sh`
- [x] 3.2 Benchmark `buildDetail` on `2cf53ec4` at 120 and 80 columns and record
      the numbers in `design.md`
- [x] 3.3 `nix flake check` and `nix build`, then open the same thread in the
      binary and scroll
