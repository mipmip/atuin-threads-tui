## MODIFIED Requirements

### Requirement: Read a thread

The list SHALL open the selected thread with `enter` and show its title,
harness, model, working directory, update time, and every message. The detail
content SHALL keep a margin of one cell on the left and on the right. Each
message SHALL show its time, author, kind, and token usage in a 24-cell meta
column pinned to the right edge of the window, and a body column SHALL fill the
space between the left margin and that meta column. Prose SHALL wrap at 80
cells; commands, diffs, patches, and tool output SHALL wrap at the width of the
body column.

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

- **WHEN** a thread is open in a 120-column terminal
- **THEN** each message prints `HH:MM - author` and `kind - N tokens` in a 24
  cell column that ends one cell before the window edge, and several blocks of
  one message join their kinds

#### Scenario: Tokens only where they exist

- **WHEN** a message reports usage and its kind is thinking
- **THEN** the meta line shows the reasoning count, a message of another kind shows
  the output count, and a message without usage shows only its kind

#### Scenario: Body measured in cells

- **WHEN** the terminal is 120 columns wide
- **THEN** the body column is 92 cells, a prose line stops at 80 cells, and a
  command of 92 cells stays on one row

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

#### Scenario: Frame has margins

- **WHEN** a thread is open
- **THEN** the first and the last cell of the window stay empty on every
  content row, while the app header, status line, and key hints keep their
  existing position

#### Scenario: Patch and summary blocks

- **WHEN** a message contains a patch block or a summary block
- **THEN** the view shows the patch file with its additions and deletions, or the
  summary text

### Requirement: Adapt to terminal size

The program SHALL re-render rows and wrapped message text when the terminal is
resized, and SHALL NOT panic when the terminal is smaller than the fixed
chrome. The width SHALL be part of the render cache key, so a rendering made for
another width is never drawn. The body column SHALL be the window width minus
the two margin cells, the gap, and the 24-cell meta column. Below 48 columns the
meta column SHALL be dropped and the body SHALL take the width minus the two
margin cells.

#### Scenario: Very short terminal

- **WHEN** the window height is 2
- **THEN** the view renders without panicking

#### Scenario: Rendering follows the width

- **WHEN** the terminal is narrowed from 120 to 80 columns
- **THEN** the body column becomes 52 cells, prose wraps at 52 instead of 80,
  and the meta column is redrawn against the right edge
