# thread-browser

## ADDED Requirements

### Requirement: Present threads as a table

The list SHALL render threads as aligned columns under a fixed header row, with
a marker, the short id, the update date, an agent tag, the title, the message
count, and the directory. Every row SHALL occupy the same display width, and
the header SHALL use the same column positions as the rows.

#### Scenario: Columns line up

- **WHEN** the list renders threads with titles of different lengths
- **THEN** the id, date, agent, message count, and directory columns start at
  the same offset on every row and under the header

#### Scenario: Header stays visible

- **WHEN** the list is scrolled to the last thread
- **THEN** the header row and its rule are still shown above the rows

#### Scenario: Date is absolute and local

- **WHEN** a thread was updated at 11:06 on 16 November 2026
- **THEN** its date column shows `16 Nov 11:06` in local time, and `-` when the
  thread has no update time

#### Scenario: Agent tag

- **WHEN** the harness is `opencode`, `claude-code`, `codex`, or empty
- **THEN** the agent column shows `oc`, `cc`, `co`, and `--`

#### Scenario: Directory is the basename

- **WHEN** the working directory is `/home/pim/gh.mipmip/atuin-threads-tui`
- **THEN** the directory column shows `atuin-threads-tui`, cut from the left
  when the name is longer than the column

#### Scenario: Narrow terminal drops columns

- **WHEN** the terminal is narrower than 100 columns
- **THEN** the directory column is hidden; below 56 columns the date column is
  hidden, below 42 the agent column, and below 38 the message count column, and
  the title keeps at least 20 cells while the width allows it

## MODIFIED Requirements

### Requirement: Navigate the list

The list SHALL keep a selection within bounds and scroll it into view, and
SHALL support `j` and down for the next thread, `k` and up for the previous
thread, `g` for the first thread, and `G` for the last thread. The header row
and its rule SHALL count as fixed chrome, so the number of thread rows that fit
is the body height minus two.

#### Scenario: Move past the end

- **WHEN** the last thread is selected and `j` is pressed
- **THEN** the selection stays on the last thread

#### Scenario: Empty list navigation

- **WHEN** no threads are loaded and `j` or `k` is pressed
- **THEN** nothing changes and no error is shown

#### Scenario: Selection stays in view

- **WHEN** the selection moves below the last visible row
- **THEN** the list scrolls by one row and the selection remains visible under
  the header
