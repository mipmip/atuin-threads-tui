# release Specification

## Purpose
TBD - created by archiving change release-management. Update Purpose after archive.

## Requirements

### Requirement: Keep the version in one file

The project SHALL keep its semantic version in a `VERSION` file beside the
command source. The binary, the flake derivation, and the release tag SHALL all derive from
that file, so no second copy of the number exists to drift.

#### Scenario: The binary reports the file

- **WHEN** `att --version` runs on a plain `go build` and on a `nix build`
- **THEN** both print the version held in `VERSION`

#### Scenario: The flake reads the file

- **WHEN** `VERSION` holds `0.2.0` and `nix build` runs
- **THEN** the derivation's version is `0.2.0` and `att --version` prints `0.2.0`

#### Scenario: An unreadable version fails loudly

- **WHEN** `VERSION` is empty or missing
- **THEN** the build fails instead of producing a binary that reports nothing

### Requirement: Cut a release with one command

A release script SHALL accept `patch`, `minor`, or `major`, bump `VERSION`,
move the `## [Unreleased]` section of `CHANGELOG.md` into a dated section for the
new version, leave a fresh empty `## [Unreleased]` above it, commit the change,
tag it `v<version>`, and push the commit and the tag.

#### Scenario: Bump a patch

- **WHEN** `VERSION` holds `0.2.0` and the script runs with `patch`
- **THEN** `VERSION` holds `0.2.1`, the changelog has a `## [0.2.1] - YYYY-MM-DD`
  section carrying the former Unreleased bullets, and a fresh `## [Unreleased]`
  sits above it

#### Scenario: Dry run changes nothing

- **WHEN** the script runs with `--dry-run`
- **THEN** it prints the version it would tag and the changelog it would write,
  and no file, commit, or tag changes

#### Scenario: Nothing to release

- **WHEN** the `## [Unreleased]` section is empty and the script runs without
  `--dry-run`
- **THEN** it stops with a message and changes nothing

#### Scenario: The tag is the version

- **WHEN** a release finishes
- **THEN** the new tag is `v` followed by the content of `VERSION`

### Requirement: Publish binaries on a tag

A GitHub workflow SHALL run on a `v*` tag, build the `att` binary for linux and
darwin on amd64 and arm64, and publish them as release assets with checksums and
notes taken from the changelog section for that version.

#### Scenario: A tag publishes four archives

- **WHEN** a `v0.2.1` tag is pushed
- **THEN** the release carries archives for linux/amd64, linux/arm64,
  darwin/amd64, and darwin/arm64, a `checksums.txt`, and release notes copied
  from the `## [0.2.1]` section

#### Scenario: A failing test blocks the release

- **WHEN** the test step fails on the tagged commit
- **THEN** no release is published

#### Scenario: A non-tag push publishes nothing

- **WHEN** a commit lands on `main` without a tag
- **THEN** the release workflow does not run
