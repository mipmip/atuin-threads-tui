# Proposal

## Why

The version is written down in three places: the `const` in `cmd/att/main.go`, the
`version` attribute in `flake.nix`, and whatever the tag says. Nothing rejects a
release where they disagree, and there is no way to cut a binary without editing
code. The project is ready for a tag, so the release path should exist before the
first version number is spent.

## What Changes

- Add a `VERSION` file beside the command source as the only place the semver
  lives, starting at `0.1.0`.
- Read it in Go through `//go:embed`, so a plain `go build` reports the same as
  `nix build`, and inject it in the flake through `-ldflags -X main.version`.
- Make the flake take its `version` from `VERSION`, so the derivation and the
  binary cannot drift apart.
- Add `scripts/release.sh patch|minor|major [--dry-run]`, which reads
  `VERSION`, bumps it, rolls the `## [Unreleased]` section of `CHANGELOG.md`
  into a dated section for the new version, commits, tags `v<new>`, and pushes.
  It uses jj like `scripts/ship-change.sh` and falls back to git.
- Add `.goreleaser.yaml` building the single `att` binary for linux and darwin
  on amd64 and arm64, with checksums and release notes taken from the changelog.
- Add `.github/workflows/release.yml` on a `v*` tag, which tests, runs
  goreleaser, and publishes the archive set.
- Add `docs/RELEASING.md` and a README pointer, so the process lives next to the
  code instead of in a commit message.
- Out of scope: Homebrew, nixpkgs, containers, signing, SBOM, and cutting the
  first tag. Those are deliberate follow-ups.

## Capabilities

### New Capabilities

- `release`: where the version number lives, how a release is cut, and what a
  published archive contains.

### Modified Capabilities

- `cli-entrypoint`: `--version` prints the built version instead of a constant.

## Impact

`cmd/att/VERSION`, `cmd/att/main.go` (embed instead of const), `flake.nix` (version and
ldflags), `scripts/release.sh`, `.goreleaser.yaml`, `.github/workflows/release.yml`,
`docs/RELEASING.md`, `README.md`, and the existing `cmd/att` version test. No new
Go dependency; `go:embed` ships with the standard library.

## Bean

`atuin-threads-tui-vv7w`
