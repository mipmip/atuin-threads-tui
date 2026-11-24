# Tasks: thread-scope-filter

## 1. Read the local host

- [x] 1.1 Add `Host(ctx) (string, error)` to `atuin.Source` and implement it in
      the client by running `atuin info`
- [x] 1.2 Parse the `host id:` line and return an empty string with no error
      when the block is absent
- [x] 1.3 Test the parser with a full block, a block without the line, and empty
      output

## 2. Scope the list

- [x] 2.1 Store `localHost`, `hostKnown`, `scopeDir` and `scopeHost` on the model
      and resolve the working directory once at startup
- [x] 2.2 Apply both scopes in `visible()` before the text filter
- [x] 2.3 Toggle `.` and `@` in the list and clamp the cursor afterwards
- [x] 2.4 Test each scope alone, both together, each with a filter query, and
      an unknown host

## 3. Show the scope

- [x] 3.1 Name the active scopes in the list status with the narrowed count, and
      show `host: unknown` when the host id could not be read
- [x] 3.2 Add `d dir` and `h host` to the list hints
- [x] 3.3 Keep the empty-result hint working with a scope active
- [x] 3.4 Test the status strings, the hint line, and the empty state

## 4. Verify

- [x] 4.1 Run `gofmt`, `go test ./...` and `bash scripts/gate.sh`
- [x] 4.2 Run `nix flake check` and `nix build` after `git add -A`
- [x] 4.3 Capture the list at 120x24 and 80x24 with `.` and with `@` active
