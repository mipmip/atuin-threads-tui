# Design

## Context

| fact | now | after |
| --- | --- | --- |
| version in Go | `const version = "0.1.0"` in `cmd/att/main.go` | `//go:embed VERSION` |
| version in flake | `version = "0.1.0"` | `lib.fileContents ./VERSION` |
| version in changelog | `## [Unreleased]` only | rolled by the script |
| CI | none | one release workflow |
| release command | none | `scripts/release.sh` |

The repository already ships `CHANGELOG.md` in Keep-a-Changelog form, and
`scripts/ship-change.sh` already shows the jj-first, git-fallback shape a release
script should follow.

## Flow

```
   cmd/att/VERSION (0.2.0)
               |        |
  go:embed ---+           +--- flake.nix version
        |                        |
   cmd/att (att 0.2.0)      nix build (att 0.2.0)

scripts/release.sh minor
   |
   +-- read VERSION -> 0.3.0
   +-- roll CHANGELOG Unreleased -> [0.3.0] - 2026-11-16
   +-- write VERSION
   +-- jj commit "Release 0.3.0"
   +-- jj tag v0.3.0 (git tag fallback)
   +-- push commit and tag
              |
              v
   GitHub Actions on v*
     +-- go test ./...
     +-- goreleaser release
           +-- att_0.3.0_linux_amd64.tar.gz
           +-- att_0.3.0_linux_arm64.tar.gz
           +-- att_0.3.0_darwin_amd64.tar.gz
           +-- att_0.3.0_darwin_arm64.tar.gz
           +-- checksums.txt
           +-- notes from CHANGELOG [0.3.0]
```

## Version plumbing

`VERSION` lives in `cmd/att/`, one line, no `v` prefix, with a trailing newline.
It sits there rather than in the repository root because `go:embed` cannot reach a parent
directory and the spec requires a plain `go build` to report the same version as
`nix build`:

```go
//go:embed VERSION
var versionFile string

var version = strings.TrimSpace(versionFile)
```

The flake and the release script read the same file (`./cmd/att/VERSION`), so
there is one number and three readers, not three numbers.

## Release script

Shaped like `scripts/ship-change.sh`:

```
release.sh <patch|minor|major> [--dry-run]
  guard: on main, clean working copy, changelog has bullets under Unreleased
  compute new version from VERSION
  if --dry-run: print version, changelog section, and the commit and tag it
                would create; exit 0
  rewrite CHANGELOG: insert "## [x.y.z] - YYYY-MM-DD" under a fresh
                     "## [Unreleased]" and move the old bullets into it
  write VERSION
  git add -A
  run scripts/gate.sh
  jj commit -m "Release x.y.z" (git commit fallback)
  jj tag vx.y.z (git tag -a fallback)
  jj bookmark set main / git branch, push commit, push tag
```

The gate runs inside the script for the same reason `ship-change.sh` runs it: a
release is the one commit that should never be pushed untested. `--dry-run`
skips the gate so it stays a read-only preview.

## Goreleaser

```
.goreleaser.yaml
  builds: one entry, id att, main ./cmd/att, env CGO_ENABLED=0,
          ldflags -s -w (version already embedded from the file)
  goos: linux, darwin
  goarch: amd64, arm64
  archives: tar.gz, name att_<version>_<os>_<arch>
  checksum: checksums.txt
  changelog: use github, but release notes come from the changelog section that
             the workflow extracts
```

The binary has no cgo and no runtime assets, so a plain archive is enough. No
Homebrew, no nixpkgs, no container image, no signing.

## Workflow

```
.github/workflows/release.yml
  on: push: tags: ["v*"]
  permissions: contents: write
  steps: checkout (fetch-depth 0 so the changelog section is available),
         setup-go (1.26), go test ./...,
         goreleaser/goreleaser-action@v6 with args release --clean
  env: GITHUB_TOKEN
```

`fetch-depth: 0` is what lets the notes step read history; the notes themselves
come from the `## [x.y.z]` section, extracted with awk, and passed through
`--release-notes`.

## Decisions

- `VERSION` lives in `cmd/att/`, not the root, because `go:embed` cannot reach
  the parent directory and the spec requires a plain `go build` to agree.
- The script rolls the changelog instead of generating it from commits, because
  the project already writes changelog bullets by hand during a ship.
- The gate runs inside the release script, so a tag cannot be created for an
  untested tree.
- The workflow runs the tests again even though the script already did: a tag can
  be pushed from another machine, and the release is the artifact other people
  trust.
- No `v` in the file, `v` on the tag: this matches git and goreleaser defaults.
- The first tag is not cut in this change. Cutting `v0.1.0` is a deliberate
  follow-up after the user has run the app once more.

## Risks

- Two build paths can still disagree if someone edits `cmd/att/VERSION` and
  forgets the flake reads it. The flake reading that exact path removes the
  drift.
- A tag pushed without the script skips the changelog roll; the workflow then
  falls back to an empty notes section rather than failing the release.
- goreleaser changes between majors. Pinning the action major (`@v6`) keeps the
  config compatible.
- `scripts/gate.sh` needs a Go toolchain and network-free vendored deps; that is
  already true in this repo.

## Testing

- a Go test that the version variable is non-empty and matches the file
- a release script dry run against a copy of the repo, asserting the changelog
  and VERSION are untouched and the printed version is right
- a script test for the empty Unreleased guard
- a `nix build` check that `att --version` equals `cat cmd/att/VERSION`
- workflow syntax check by eye; there is no CI in the repo yet to test it
