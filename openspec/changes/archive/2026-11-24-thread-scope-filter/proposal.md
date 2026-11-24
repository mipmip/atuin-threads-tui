# Proposal

## Why

The list mixes every project and every machine. Typing a directory into `/`
finds it, but there is no way to say "only this project" or "only this
computer", so a reader scrolling for today's work passes threads from 84
directories and two hosts.

## What Changes

- Add `.` to keep only threads whose working directory is the directory att was
  started in.
- Add `@` to keep only threads captured on this host.
- Both are toggles, compose with each other, with `/` and with `ctrl+s`
  search results, and each is named in the status line while active.
- Read the local host id once at startup from `atuin info`, and treat a failed
  read as "host unknown" rather than as an empty list.
- No new column: the host ids are opaque device ids, so a column would print
  either hex or an invented code and be constant on almost every row.
- No persistence. The scope belongs to the running session, like the message
  type selection.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the list can be limited to the current directory or the
  local host.

## Impact

`internal/atuin/client.go` (read `atuin info`), `internal/ui/app.go` (two
flags and the host id), `internal/ui/view.go` (status and hints),
`internal/ui/filter.go` or `visible()` (apply the scope), tests. No schema
change, no new dependency.

## Bean

`atuin-threads-tui-yqbs`
