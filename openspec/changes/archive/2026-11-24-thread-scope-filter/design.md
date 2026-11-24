# Design

## Context

Local data: 257 threads across 84 working directories and exactly two hosts. The
current directory holds 1 thread; the local host owns 161 of the 257.

The host values are opaque device ids:

```
01a18745351c74618f8f9d94eac717e8   161 threads   the local machine
019fed235b887d91b921b2eed7db53c8    96 threads
```

`atuin info` prints the local one:

```
Host info:
host id: 01a18745351c74618f8f9d94eac717e8
```

`~/.local/share/atuin/meta.db` holds the same value in `meta(key, value)`, but
reading it would add a sqlite driver for one row. `threads.db` stores it twice
more (`sessions.host_id`, `messages.host` as an interned id), both behind
sqlite.

## Shape

```
list state
   threads ............ all from atuin
      |
      v  visible()
   scopeDir ?  cwd == thread.CWD
   scopeHost ? localHost == thread.Host
      |
      v
   text filter (/)  --->  rows
```

Both flags live on the model next to `filter` and are applied in `visible()`
before the text filter, so `/`, `ctrl+s` results and the type panel all compose
without special cases.

Scope is session state. It is not persisted, not tied to a thread, and survives
`ctrl+r`.

## Reading the host

`Source` gains `Host(ctx) (string, error)`, implemented with a `Runner` call to
`atuin info` and a scan for a line whose trimmed prefix is `host id:`. The model
asks for it once in `Init` alongside the thread list, and stores the result in
`localHost` with an `hostKnown` flag.

When the read fails, `@` still toggles but the status reads `host: unknown` and
the rows are unchanged. An empty list would look like a bug and hide work.

## Status and hints

- `dir: atuin-threads-tui · 1 of 257 threads`
- `host: local · 161 of 257 threads`
- both: `dir: atuin-threads-tui · host: local · 1 of 257 threads`
- with a filter the count becomes `1 of 3 threads` as it does today
- hints gain `d dir · h host`, placed after `ctrl+r reload`

The directory is labelled with its basename, not the full path, to keep the
status line inside 80 cells; the full path stays visible in the detail header.

## Decisions

- Two toggles instead of one cycling key: the states are independent and a cycle
  would need three presses to get back.
- Compare full paths, not basenames, so two checkouts with the same final
  component do not merge.
- No host column, because a 2 to 6 cell column of hex is noise on a row that is
  already the same value 97% of the time. The column would cost title width for
  no reading value.
- Read the host from `atuin info` instead of adding sqlite, because att already
  speaks to atuin through the CLI and a driver for one row is not worth it.
- Resolve the working directory once at startup; a scope that silently follows a
  later `chdir` would be surprising.

## Risks

- `atuin info` output is not a stable contract, so the parse is a scan for a
  prefix rather than a positional read, and failure degrades to `unknown`.
- A thread captured in a subdirectory of the current project is not matched.
  Exact equality is predictable; prefix matching would need a rule for `/tmp`.
- Both scopes on in a fresh directory shows nothing, which the empty state and
  the status count make explainable.

## Testing

- host parsing from a sample `atuin info` block, a missing block, and empty
  output
- `visible()` with no scope, dir only, host only, both, and each combined with a
  filter query
- unknown host toggling leaves the rows alone
- status text for each combination and the hint line
- keys toggle on and off without touching the cursor beyond clamping
