# thread-browser Specification

## Purpose
TBD - created by archiving change thread-list-ui. Update Purpose after archive.

## Requirements

### Requirement: Load threads on start

The program SHALL request the thread list when it starts and render a loading
state until the request finishes.

#### Scenario: Threads arrive

- **WHEN** the list request returns three threads
- **THEN** the view shows three rows and the status line says `3 threads`

#### Scenario: No threads

- **WHEN** the list request returns none
- **THEN** the view shows a hint that no threads have been captured

#### Scenario: List request fails

- **WHEN** the list request returns an error
- **THEN** the view shows the error text and the key hint for retrying

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

### Requirement: Read a thread

The list SHALL open the selected thread with `enter` and show its title, harness,
model, working directory, update time, and every message. Each message SHALL show
its time, author, kind, and token usage in a meta column to the right of a body
column no wider than 80 cells.

#### Scenario: Open and return

- **WHEN** `enter` is pressed on a thread and later `esc` is pressed
- **THEN** the conversation is shown and the list is shown again with the same
  selection

#### Scenario: Content blocks

- **WHEN** a message contains text, thinking, a tool call, and a tool result
- **THEN** the view renders each block with a label naming its kind

#### Scenario: Detail request fails

- **WHEN** the conversation request returns an error
- **THEN** the view shows the error and stays in the list

#### Scenario: Message meta column

- **WHEN** a thread is open
- **THEN** each message prints `HH:MM - author` and `kind - N tokens` in a 20 cell
  column to the right of the body, and several blocks of one message join their
  kinds

#### Scenario: Tokens only where they exist

- **WHEN** a message reports usage and its kind is thinking
- **THEN** the meta line shows the reasoning count, a message of another kind shows
  the output count, and a message without usage shows only its kind

#### Scenario: Body measured in cells

- **WHEN** the terminal is wider than the body column plus the meta column
- **THEN** the body stays at 80 cells and the remaining space stays empty

#### Scenario: Repeated empty messages collapse

- **WHEN** several consecutive messages have no displayable content
- **THEN** they appear as one row carrying their count, such as
  `(12 attachments)`

#### Scenario: Tool call rendering

- **WHEN** a message holds a tool call whose input names a command
- **THEN** the body shows the tool name and the command without JSON braces,
  quotes, or escaped characters

#### Scenario: Tool result keeps its layout

- **WHEN** a tool result contains indented lines
- **THEN** the view keeps the indentation and wraps only lines longer than the
  body column

#### Scenario: Patch and summary blocks

- **WHEN** a message contains a patch block or a summary block
- **THEN** the view shows the patch file with its additions and deletions, or the
  summary text

### Requirement: Scroll a conversation

The detail view SHALL scroll by one line with `j` and `k`, by one page with
`pgdown` and `pgup`, and jump with `g` and `G`, without scrolling past either
end. Scrolling SHALL read the rendered blocks from the cache and SHALL NOT
render a message again, and the status line SHALL report the position against
the cached total number of lines.

#### Scenario: Scroll past the last line

- **WHEN** the view is at the last line and `j` is pressed
- **THEN** the offset does not grow

#### Scenario: Fast repeated keys all count

- **WHEN** five `j` keys arrive in one input batch
- **THEN** the offset advances by five lines

#### Scenario: Status reports the last line

- **WHEN** the conversation has 55,527 rendered lines and the view is scrolled
  to the end
- **THEN** the status line reports the last line and no render is triggered

### Requirement: Adapt to terminal size

The program SHALL re-render rows and wrapped message text when the terminal is
resized, and SHALL NOT panic when the terminal is smaller than the fixed
chrome. The width SHALL be part of the render cache key, so a rendering made for
another width is never drawn.

#### Scenario: Very short terminal

- **WHEN** the window height is 2
- **THEN** the view renders without panicking

#### Scenario: Rendering follows the width

- **WHEN** the terminal is narrowed from 120 to 80 columns
- **THEN** the message body is measured again at 80 columns and the meta column
  is redrawn to its right

### Requirement: Quit

The program SHALL quit on `q` in the list and on `ctrl+c` anywhere.

#### Scenario: Quit from the list

- **WHEN** `q` is pressed in the list
- **THEN** the program returns a quit command

### Requirement: Present threads as a table

The list SHALL render threads as aligned columns under a fixed header row, in
this order: a marker, the short id, the update date, the agent code, the message
count, the title, and the directory. The title column SHALL be the width of the
longest title in view, clamped between 20 and 48 cells. The directory column
SHALL be the width of the longest directory basename in view, clamped between 8
and 24 cells. The agent column SHALL be the width of the longest agent code in
view, clamped between 2 and 6 cells. Every row SHALL occupy the same display
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

- **WHEN** the harness is `claude-code`, `opencode`, `codex`, `pi`,
  `opencode-fork`, or empty
- **THEN** the agent column shows `claude`, `ocode`, `codex`, `pi`, `openco`,
  and `--`

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

### Requirement: Render a thread once

The detail view SHALL render a conversation into per-message blocks once when it
opens, and SHALL reuse those blocks while the user scrolls. It SHALL rebuild the
rendering when the terminal width changes, and SHALL render a conversation of
more than 500 messages in the background instead of blocking the interface.

#### Scenario: Cache is reused

- **WHEN** a thread is open and 20 scroll keys are pressed
- **THEN** the conversation is rendered once and every frame reads the blocks
  from the cache

#### Scenario: Large thread renders in the background

- **WHEN** a thread with more than 500 messages is opened
- **THEN** the view shows a rendering state immediately and the conversation
  appears when the render finishes

#### Scenario: Resize rebuilds the rendering

- **WHEN** the terminal width changes while a thread is open
- **THEN** the conversation is rendered again at the new width and the body
  column follows the new width

#### Scenario: Stale render is dropped

- **WHEN** a background render finishes for a width the terminal no longer has
- **THEN** the result is discarded and the current width keeps its own
  rendering
