# Design: att-bootstrap

## Module and layout

```
go.mod                      module github.com/mipmip/atuin-threads-tui
cmd/att/main.go             main() delegates to run(args, stdout, stderr)
cmd/att/main_test.go        table tests over run
internal/...                populated by later changes
```

`run` returns an exit code and receives `io.Writer`s. `main` is one line, so the
coverage gate exercises nearly all command behavior.

## Flag handling

Use `flag.NewFlagSet` with `ContinueOnError` and `output` set to the provided
stderr, so unknown flags become a status 2 exit with a useful message.
`--version` and `--help` write to stdout and return 0.

## Gate integration

The flake's `checks.gate` runs `scripts/gate.sh` in a sandbox with `GOPROXY=off`.
Dependencies are committed under `vendor/`, and `buildGoModule` uses
`vendorHash = null`, which makes Nix build from that directory.

## Out of scope

The TUI itself, the atuin data source, and any network behavior.
