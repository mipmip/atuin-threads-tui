# thread-browser delta

## MODIFIED Requirements

### Requirement: Present threads as a table

The list SHALL render threads as aligned columns under a fixed header row, in
this order: a marker, the short id, the update date, the agent code, the message
count, the title, and the directory. The title column SHALL be the width of the
longest title in view, clamped between 20 and 48 cells. The directory column
SHALL be the width of the longest directory basename in view, clamped between 8
and 24 cells. The agent column SHALL be the width of the longest agent code in
view, clamped between 2 and 5 cells. Every row SHALL occupy the same display
width as the header, and that width SHALL be the sum of the columns and their
gutters rather than the terminal width.

#### Scenario: Columns line up

- **WHEN** the list renders threads with titles of different lengths
- **THEN** the date, agent, message count, title, and directory columns start
  at the same offset on every row and under the header

#### Scenario: Title follows the content

- **WHEN** the longest title in view is 36 cells
- **THEN** the title column is 36 cells, the table is narrower than the
  terminal, and no directory name is cut

#### Scenario: Title is capped

- **WHEN** the longest title in view is longer than 48 cells
- **THEN** the title column is 48 cells and longer titles are cut with an
  ellipsis

#### Scenario: Directory follows the content

- **WHEN** the longest directory basename in view is 22 cells
- **THEN** the directory column is 22 cells and the basename is shown in full

#### Scenario: Header stays visible

- **WHEN** the list is scrolled to the last thread
- **THEN** the header row and its rule are still shown above the rows

#### Scenario: Date is absolute and local

- **WHEN** a thread was updated at 11:06 on 16 November 2026
- **THEN** its date column shows `16 Nov 11:06` in local time, and `-` when the
  thread has no update time

#### Scenario: Agent tag

- **WHEN** the harness is `opencode`, `claude-code`, `codex`, `opencode-fork`,
  or empty
- **THEN** the agent column shows `oc`, `cc`, `co`, `openc`, and `--`

#### Scenario: Directory is the basename

- **WHEN** the working directory is `/home/pim/gh.mipmip/atuin-threads-tui`
- **THEN** the directory column shows `atuin-threads-tui`, cut from the left
  only when the name is longer than the column

#### Scenario: Narrow terminal drops columns

- **WHEN** the table does not fit the terminal
- **THEN** the title column shrinks to 20 cells, then the directory column
  shrinks to 8 cells, and only then the directory, date, agent, and message
  count columns are dropped in that order, while the title keeps at least one
  cell
