# Design: thread-search-filter

## Two kinds of search

Local filtering is a substring match over fields already in memory. It runs on
every keystroke and never blocks on a process.

Full-text search asks atuin, whose index covers message bodies, reasoning, and
tool calls. It runs once on `ctrl+s` because it costs a process spawn, and it
replaces the list with the matches.

```go
modeList --"/"--> modeFilter --enter--> modeList (filter applied)
                            --esc----> modeList (filter cleared)
modeList --"ctrl+s"--> searchMsg --esc--> reload
```

## State

```go
filter      string  // applied local filter
filterInput string  // text being typed in filter mode
searchQuery string  // non-empty when showing atuin results
```

`visible()` derives the rows from `threads` and `filter`; nothing caches a
second slice, so a filter change cannot go stale. `cursor` and `offset` index the
visible rows, and both clamp whenever the visible length changes.

While filtering, `filter` mirrors `filterInput` so results update on each
keystroke. `enter` keeps the text and leaves the mode; `esc` clears it.

## Clearing search

A search replaces `threads`, so clearing cannot restore the previous list from
memory. `esc` clears `searchQuery` and issues a fresh list command. That is one
process spawn on an explicit user action, which is acceptable.

## Matching

`matches` lowercases both sides and checks each searchable field with
`strings.Contains`. Splitting a query into words would make the filter stricter
than the user expects while typing, so a phrase is a phrase.

## Out of scope

Highlighting matches, fuzzy ranking, and saving queries.
