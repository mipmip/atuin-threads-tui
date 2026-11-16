---
# atuin-threads-tui-vv7w
title: Release versioning, changelog and tag driven release
status: completed
openspec-link: openspec/changes/archive/2026-11-16-release-management
type: feature
priority: normal
created_at: 2026-11-16T21:24:34Z
updated_at: 2026-11-16T21:34:11Z
parent: atuin-threads-tui-fy6s
---

One semver source in VERSION, a release script that bumps it and rolls the changelog, goreleaser config and a GitHub workflow that builds att on a v* tag.

## Summary of Changes

The version now lives only in `cmd/att/VERSION`, read by `go:embed` for a plain `go build` and by `lib.fileContents` for the flake, so both report the same number. `scripts/release.sh patch|minor|major [--dry-run]` computes the next version, refuses an empty Unreleased section or an existing tag, runs the gate before committing, rolls the changelog into a dated section, commits, tags `v<version>` and pushes both the commit and the tag. `.goreleaser.yaml` builds `att` for linux and darwin on amd64 and arm64 as tar.gz with checksums, and `.github/workflows/release.yml` runs the tests and publishes on a `v*` tag with the changelog section as the release notes. `docs/RELEASING.md` describes the flow and the README points at it; `goreleaser check` validates the config.
