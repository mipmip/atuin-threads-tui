# Tasks: att-bootstrap

## 1. Module

- [x] 1.1 Create `go.mod` for `github.com/mipmip/atuin-threads-tui` with Go 1.24
- [x] 1.2 Add `cmd/att/main.go` with `run(args, stdout, stderr) int`

## 2. Command behavior

- [x] 2.1 Implement `--version` printing `att <version>`
- [x] 2.2 Implement `--help` usage text naming `att`
- [x] 2.3 Reject unknown flags with status 2 and a stderr message

## 3. Tests

- [x] 3.1 Table tests for version, help, unknown flag, and no-argument startup

## 4. Build and gate

- [x] 4.1 Point the flake package and app at `att`
- [x] 4.2 Vendor dependencies and confirm `nix flake check` passes
- [x] 4.3 Update `AGENTS.md` commands to name `att`
