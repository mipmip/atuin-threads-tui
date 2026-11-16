# thread-browser

## ADDED Requirements

### Requirement: Choose which message types are shown

The detail view SHALL offer a panel that lists every message category with its
count and lets the reader switch each one on or off. The panel SHALL open with
`t`, toggle the marked row with space, select everything with `a`, select the
conversation with `n`, and select a preset with `1`, `2`, `3`, or `4`. Enter
SHALL apply the selection and `esc` SHALL close the panel without changing what
is shown.

#### Scenario: Panel opens

- **WHEN** `t` is pressed while a thread is open
- **THEN** the view lists prompts, replies, thinking, tool calls, tool results,
  patches, summaries, and notes with a count for each, and the current
  selection is marked

#### Scenario: Toggle a category

- **WHEN** space is pressed on the thinking row
- **THEN** the row changes between marked and unmarked, and the conversation
  stays as it is until the selection is applied

#### Scenario: Apply the selection

- **WHEN** enter is pressed in the panel
- **THEN** the conversation is rebuilt with only the selected categories and the
  panel closes

#### Scenario: Cancel the selection

- **WHEN** `esc` is pressed in the panel
- **THEN** no category changes and the conversation is shown as it was

#### Scenario: Presets

- **WHEN** `1`, `2`, `3`, or `4` is pressed in the panel
- **THEN** the selection becomes everything, the conversation, the tools, or the
  prompts only, in that order

#### Scenario: Selection outlives the thread

- **WHEN** a thread is closed and another one is opened
- **THEN** the same selection applies to the new thread

### Requirement: Show only the selected message types

The conversation SHALL render only the blocks whose category is selected. A
message without a visible block SHALL NOT be rendered, and a message without
content SHALL be rendered only while notes is selected. The meta column SHALL
name only the kinds that are visible, and the status line SHALL report the
number of rendered messages against atuin's message count.

#### Scenario: Hidden blocks are dropped

- **WHEN** thinking is unselected
- **THEN** no reasoning text appears, and a reply that shared a message with
  that reasoning still renders

#### Scenario: Message without visible blocks

- **WHEN** an assistant message holds only a tool call and tool calls are
  unselected
- **THEN** the message does not appear and leaves no empty meta row behind

#### Scenario: Notes hide empty messages

- **WHEN** notes is unselected
- **THEN** attachments, titles, and empty system messages do not appear, and no
  `(16 attachments)` row is shown

#### Scenario: Notes show as one row

- **WHEN** notes is selected and 16 consecutive messages have no content
- **THEN** they appear as one row carrying their count

#### Scenario: Meta names visible kinds only

- **WHEN** a message holds text and a tool call, and tool calls are unselected
- **THEN** the meta column names only the text kind

#### Scenario: Nothing selected

- **WHEN** every category is unselected
- **THEN** the body says that no messages of the selected types are shown and no
  message is drawn

#### Scenario: Status counts shown and total

- **WHEN** atuin reports 7,976 messages for the thread and the selection renders
  1,301 of them
- **THEN** the status line reports `1301 of 7976 msgs`

## MODIFIED Requirements

### Requirement: Read a thread

The list SHALL open the selected thread with `enter` and show its title,
harness, model, working directory, update time, and every visible message. The
detail content SHALL keep a margin of one cell on the left and on the right.
Each message SHALL show its time, author, kind, and token usage in a 24-cell meta
column pinned to the right edge of the window, and a body column SHALL fill the
space between the left margin and that meta column. Prose SHALL wrap at 80
cells; commands, diffs, patches, and tool output SHALL wrap at the width of the
body column.

#### Scenario: Open and return

- **WHEN** `enter` is pressed on a thread and later `esc` is pressed
- **THEN** the conversation is shown and the list is shown again with the same
  selection

#### Scenario: Content blocks

- **WHEN** a message contains text, thinking, a tool call, and a tool result and
  every type is selected
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

- **WHEN** several consecutive messages have no displayable content and notes
  is selected
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

- **WHEN** a message contains a patch block or a summary block and those types
  are selected
- **THEN** the view shows the patch file with its additions and deletions, or the
  summary text

### Requirement: Render a thread once

The detail view SHALL render a conversation into per-message blocks once when it
opens, and SHALL reuse those blocks while the user scrolls. It SHALL rebuild the
rendering when the terminal width or the message type selection changes, and
SHALL render a conversation of more than 500 messages in the background instead
of blocking the interface.

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

#### Scenario: Changing the types rebuilds the rendering

- **WHEN** a different message type selection is applied
- **THEN** the conversation is rendered again with that selection and the
  message at the top of the viewport is still at the top

#### Scenario: Stale render is dropped

- **WHEN** a background render finishes for a width the terminal no longer has
- **THEN** the result is discarded and the current width keeps its own
  rendering
