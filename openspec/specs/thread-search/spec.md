# thread-search Specification

## Purpose
TBD - created by archiving change thread-search-filter. Update Purpose after archive.

## Requirements

### Requirement: Filter the loaded list

The list SHALL filter as the user types after pressing `/`, matching a
case-insensitive substring against title, harness, model, working directory,
preview, last reply, git branch, and thread id.

#### Scenario: Title match

- **WHEN** the filter is `openspec`
- **THEN** only threads whose title contains `openspec` are shown

#### Scenario: Case insensitive field match

- **WHEN** the filter is `DELLIE` and a thread has that text in its cwd
- **THEN** the thread is shown

#### Scenario: No match

- **WHEN** no thread matches
- **THEN** the list is empty, the status says `0 of N threads`, and the app does
  not exit

### Requirement: Filter mode keys

Filter mode SHALL append typed characters, remove the last character on
`backspace`, clear the input on `ctrl+u`, apply with `enter`, and clear the
filter with `esc`.

#### Scenario: Apply then navigate

- **WHEN** `enter` applies a filter and `j` is pressed
- **THEN** the selection moves within the filtered rows only

#### Scenario: Clear

- **WHEN** `esc` is pressed in filter mode
- **THEN** the filter is empty and every thread is listed again

#### Scenario: Backspace to empty

- **WHEN** every typed character is removed with `backspace`
- **THEN** every thread is listed again

### Requirement: Counts reflect the view

The status line SHALL show `N threads` when unfiltered, `N of M threads` when
filtered, and `N results` after a full-text search.

#### Scenario: Filtered count

- **WHEN** three of twelve threads match
- **THEN** the status says `3 of 12 threads`

### Requirement: Full-text search

`ctrl+s` SHALL search atuin with the current filter text, show a loading state
until the result arrives, and replace the list with matching threads.

#### Scenario: Search with results

- **WHEN** the filter is `bootstrap` and `ctrl+s` is pressed
- **THEN** atuin is queried with `bootstrap` and the returned threads are shown

#### Scenario: Empty query

- **WHEN** the filter is empty and `ctrl+s` is pressed
- **THEN** no command is started

#### Scenario: Search fails

- **WHEN** the search request returns an error
- **THEN** the error is shown and the previous list stays visible

### Requirement: Clear a search

`esc` on search results SHALL clear the query and reload the full thread list.

#### Scenario: Back to all threads

- **WHEN** `esc` is pressed while search results are shown
- **THEN** a list request is issued and the status returns to the total count
