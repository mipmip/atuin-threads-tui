# release-management

Give the project one version number and a repeatable way to publish it.

`cmd/att/VERSION` becomes the only semver source: the Go binary embeds it, the
flake reads it, and `scripts/release.sh patch|minor|major` bumps it, rolls the
changelog, commits, tags `v<version>`, and pushes. A tag drives a goreleaser
build for linux and darwin on amd64 and arm64.
