# Releasing

The version lives in `cmd/att/VERSION` and nowhere else. The binary embeds it,
`flake.nix` reads it, and `scripts/release.sh` bumps it.

## Before you release

Write the entries under `## [Unreleased]` in `CHANGELOG.md`. The script refuses
to run when that section has no bullet, because a release with no notes is not
worth tagging.

## Dry run

```bash
bash scripts/release.sh patch --dry-run
```

This prints the version it would produce, the changelog section it would write,
and the commit and tag it would create. It writes nothing. Use it to check the
notes before they are frozen.

## Cut the release

```bash
bash scripts/release.sh patch     # 0.1.0 -> 0.1.1
bash scripts/release.sh minor     # 0.1.0 -> 0.2.0
bash scripts/release.sh major     # 0.1.0 -> 1.0.0
```

The script checks that you are on `main`, bumps `cmd/att/VERSION`, rolls
`## [Unreleased]` into a dated `## [x.y.z]` section above a fresh
`## [Unreleased]`, runs `bash scripts/gate.sh`, commits `Release vx.y.z`, moves
the `main` bookmark, tags `vx.y.z`, and pushes the commit and the tag.

It refuses to continue when the tree is not on `main`, when the tag already
exists, or when the gate fails. Nothing is committed before the gate passes.

## What the tag produces

Pushing a tag that starts with `v` starts `.github/workflows/release.yml`, which
tests the code and runs goreleaser. It builds `att` for linux and darwin on
amd64 and arm64, with `CGO_ENABLED=0`, and publishes:

```
att_<version>_linux_amd64.tar.gz
att_<version>_linux_arm64.tar.gz
att_<version>_darwin_amd64.tar.gz
att_<version>_darwin_arm64.tar.gz
checksums.txt
```

The release body is the changelog section for that version, and it is empty when
the section is missing.

## After the workflow

- the release page lists four archives and a checksum file
- `sha256sum -c checksums.txt` passes in the download directory
- `att --version` on a downloaded archive prints the version in the tag

A failed test in the workflow leaves the tag in place without a release. Fix the
problem, then either delete and re-push the tag or cut the next patch.

## Nix

`nix build` reports the same version, because the derivation reads
`cmd/att/VERSION`:

```bash
nix build && ./result/bin/att --version
```
