#!/usr/bin/env bash
# Gate for `nix flake check`: build, vet, tests, and coverage thresholds.
set -euo pipefail

CORE_RE="${CORE_RE:-/(internal|core)/}"
MIN_TOTAL="${MIN_TOTAL:-70}"
MIN_CORE="${MIN_CORE:-80}"
PROFILE="$(mktemp -d)/coverage.out"
fail=0

report() {
  awk -v value="$1" -v min="$2" -v label="$3" 'BEGIN {
    printf "%-40s %6.1f%% (min %.0f%%)\n", label, value, min
    exit (value + 0 < min + 0) ? 1 : 0
  }'
}

echo "==> build"
go build -o "$(mktemp -d)" ./...

echo "==> vet"
go vet ./...

echo "==> test"
go test ./... -coverprofile="$PROFILE" -covermode=atomic

overall="$(go tool cover -func="$PROFILE" | awk '/^total:/ { sub(/%/, "", $3); print $3 }')"
echo
echo "==> coverage thresholds"
report "${overall:-0}" "$MIN_TOTAL" "overall" || fail=1

mapfile -t core < <(go list ./... 2>/dev/null | grep -E "$CORE_RE" || true)
if (( ${#core[@]} == 0 )); then
  echo "no core packages match ${CORE_RE}, skipping core gate"
else
  for pkg in "${core[@]}"; do
    out="$(go test -cover "$pkg" 2>&1 || true)"
    pct="$(printf '%s\n' "$out" | sed -n 's/.*coverage: \([0-9.]*\)%.*/\1/p' | head -1)"
    report "${pct:-0}" "$MIN_CORE" "$pkg" || fail=1
  done
fi

echo
if (( fail )); then
  echo "gate: FAILED" >&2
  exit 1
fi
echo "gate: passed"
