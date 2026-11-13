# Proposal: att-bootstrap

## Why

The repository is an empty scaffold: no Go module, no entrypoint, and a flake
that cannot build anything. Until `att` exists as a building binary, no other
milestone can be implemented, tested, or shipped.

## What Changes

- Add the Go module `github.com/mipmip/atuin-threads-tui` targeting Go 1.24.
- Add `cmd/att` with a testable `run` function, `--version`, and a first TUI
  launch that reports a clear error while no data source is wired in.
- Point the flake package, app, and dev shell at the `att` binary, and make the
  gate run `scripts/gate.sh` against a vendored build.
- Vendor dependencies so `nix flake check` builds without network access.

## Capabilities

### New Capabilities

- `cli-entrypoint`: the `att` command, its version output, and its behavior when
  it cannot start.

### Modified Capabilities

None.

## Impact

- `flake.nix`: package name and `mainProgram` become `att`.
- `go.mod`, `go.sum`, `vendor/`: new.
- `AGENTS.md`: commands show `att` where relevant.

## Bean

- atuin-threads-tui-kp3z: Go module and Nix flake

