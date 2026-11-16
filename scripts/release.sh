#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
usage: scripts/release.sh <patch|minor|major> [--dry-run]

Bumps cmd/att/VERSION, rolls the Unreleased changelog section into a
dated section, runs the gate, then commits and tags v<version>.
EOF
}

bump="${1:-}"
dry_run=""

case "$bump" in
  patch|minor|major) ;;
  -h|--help|"")
    usage
    exit 2
    ;;
  *)
    echo "unknown bump: $bump" >&2
    usage
    exit 2
    ;;
esac

if [ "${2:-}" = "--dry-run" ]; then
  dry_run=1
elif [ -n "${2:-}" ]; then
  echo "unexpected argument: $2" >&2
  usage
  exit 2
fi

root=$(git rev-parse --show-toplevel)
cd "$root"

version_file="cmd/att/VERSION"
changelog="CHANGELOG.md"

for file in "$version_file" "$changelog"; do
  if [ ! -f "$file" ]; then
    echo "missing $file" >&2
    exit 1
  fi
done

current=$(tr -d ' \t\r\n' < "$version_file")
case "$current" in
  ""|*[!0-9.]*)
    echo "version '$current' is not numeric" >&2
    exit 1
    ;;
esac

IFS=. read -r major minor patch <<<"$current"
major=${major:-0}
minor=${minor:-0}
patch=${patch:-0}

case "$bump" in
  major)
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  minor)
    minor=$((minor + 1))
    patch=0
    ;;
  patch)
    patch=$((patch + 1))
    ;;
esac

next="$major.$minor.$patch"
tag="v$next"
today=$(date +%Y-%m-%d)

if [ -n "$(git rev-parse -q --verify "refs/tags/$tag" || true)" ]; then
  echo "tag $tag already exists" >&2
  exit 1
fi

if ! awk '
  /^## \[Unreleased\]/ { inside = 1; next }
  inside && /^## / { inside = 0 }
  inside && /^- / { found = 1 }
  END { exit found ? 0 : 1 }
' "$changelog"; then
  echo "the Unreleased section has no entries" >&2
  exit 1
fi

rolled=$(
  awk -v ver="$next" -v date="$today" '
    function emit(line) {
      print line
      lastblank = (line == "")
    }
    BEGIN { state = 0; lastblank = 0; seen = 0 }
    /^## \[Unreleased\]/ && state == 0 {
      emit("## [Unreleased]")
      emit("")
      emit("## [" ver "] - " date)
      state = 1
      next
    }
    state == 1 && /^## / {
      if (!lastblank) emit("")
      emit($0)
      state = 2
      next
    }
    state == 1 {
      if (!seen && $0 ~ /^[ \t]*$/) next
      seen = 1
      emit($0)
      next
    }
    { emit($0) }
  ' "$changelog"
)

section=$(
  printf '%s\n' "$rolled" | awk -v ver="$next" '
    /^## \[/ {
      if (on) exit
      if (index($0, "[" ver "]")) on = 1
    }
    on { print }
  '
)

if [ -n "$dry_run" ]; then
  echo "version: $current -> $next"
  echo "tag: $tag"
  echo "commit: Release $tag"
  echo
  echo "changelog section:"
  printf '%s\n' "$section"
  exit 0
fi

if [ -d .jj ]; then
  if [ -z "$(jj log -r 'main & ::@' --no-graph -T commit_id)" ]; then
    echo "not on main" >&2
    exit 1
  fi
else
  branch=$(git rev-parse --abbrev-ref HEAD)
  if [ "$branch" != "main" ]; then
    echo "not on main (on $branch)" >&2
    exit 1
  fi
fi

echo "==> [1/4] bump $current -> $next"
printf '%s\n' "$next" > "$version_file"
printf '%s\n' "$rolled" > "$changelog"
git add -A

echo "==> [2/4] gate"
bash scripts/gate.sh

echo "==> [3/4] commit and tag $tag"
if [ -d .jj ]; then
  jj commit -m "Release $tag"
  jj bookmark set main -r @-
  jj tag set "$tag" -r @-
  jj git push --bookmark main
  git push origin "refs/tags/$tag"
else
  git commit -m "Release $tag"
  git tag -a "$tag" -m "Release $tag"
  git push origin HEAD
  git push origin "refs/tags/$tag"
fi

echo "==> [4/4] released $tag"
