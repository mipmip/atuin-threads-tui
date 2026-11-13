# Proposal: thread-search-filter

## Why

A list of every captured thread is only useful if a specific conversation can be
found quickly. The list is already loaded, so narrowing it should be instant, and
looking inside message text should use the full-text index atuin maintains.

## What Changes

- Add a filter mode opened with `/`: case-insensitive substring matching over
  title, harness, model, working directory, preview, last reply, branch, and
  thread id, applied while typing.
- Apply the filter with `enter`, clear it with `esc`.
- Add full-text search with `ctrl+s`, calling `atuin threads search` with the
  current filter and showing the matches.
- Clear an active search with `esc`, which reloads the full list.
- Show counts that distinguish filtered and searched results from the full list.

## Capabilities

### New Capabilities

- `thread-search`: filtering the loaded list and querying atuin's full-text
  index.

### Modified Capabilities

None.

## Impact

- `internal/ui`: a third mode, filter state, and a search message.
- No new dependencies.

## Bean

- atuin-threads-tui-4ojf: Search and filtering
