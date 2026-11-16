# Tasks: release-management

## 1. One version source

- [x] 1.1 Add `cmd/att/VERSION` holding `0.1.0`
- [x] 1.2 Replace the `const version` in `cmd/att/main.go` with a `//go:embed`
      of `VERSION`, trimmed at init
- [x] 1.3 Drop the ldflags injection from the flake and point its `version` at
      `./cmd/att/VERSION` through `lib.fileContents`
- [x] 1.4 Extend the `cmd/att` version test so it compares against the file
      content and fails when the file is empty

## 2. Release script

- [x] 2.1 Add `scripts/release.sh <patch|minor|major> [--dry-run]` with the
      jj-first, git-fallback shape of `scripts/ship-change.sh`
- [x] 2.2 Compute the next version from `cmd/att/VERSION` and reject an unknown
      bump word
- [x] 2.3 Roll the `## [Unreleased]` bullets into `## [x.y.z] - YYYY-MM-DD`
      above a fresh `## [Unreleased]`, and refuse to continue when Unreleased is
      empty
- [x] 2.4 `--dry-run` prints the version, the changelog section, and the commit
      and tag it would create, and writes nothing
- [x] 2.5 Run `scripts/gate.sh` before committing, then commit, tag `v<version>`,
      move the bookmark, and push the commit and the tag

## 3. Publishing

- [x] 3.1 Add `.goreleaser.yaml` with one build of `./cmd/att`, `CGO_ENABLED=0`,
      linux and darwin on amd64 and arm64, tar.gz archives, and a checksum file
- [x] 3.2 Add `.github/workflows/release.yml` on `push: tags: ["v*"]` with
      `contents: write`, a checkout at full depth, Go setup, `go test ./...`, and
      the goreleaser action reading notes from the changelog section
- [x] 3.3 Make the notes step fall back to an empty body when the changelog has
      no section for the tag

## 4. Docs

- [x] 4.1 Add `docs/RELEASING.md` covering the bump, the dry run, the tag, the
      artifacts, and what to check after the workflow finishes
- [x] 4.2 Add a README section pointing at `docs/RELEASING.md`

## 5. Verify

- [x] 5.1 Run `bash scripts/release.sh patch --dry-run` and confirm no file,
      commit, or tag changed
- [x] 5.2 Run `go test ./...` and `bash scripts/gate.sh`
- [x] 5.3 `git add -A`, then `nix flake check` and `nix build`, and confirm
      `result/bin/att --version` matches `cmd/att/VERSION`
- [x] 5.4 Validate `.goreleaser.yaml` with `goreleaser check` if the binary is
      available, or note the manual check
