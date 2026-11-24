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

### Requirement: Render long lines in linear time

The detail view SHALL cut each rendered line in a single pass, so that the cost
of wrapping one line grows proportionally with its length and does not depend on
how many chunks it produces.

#### Scenario: Very long result line

- **WHEN** a tool result contains a single line of 100,000 characters
- **THEN** the conversation is laid out in one pass over that line
- **AND** the line is split into chunks of the body width with the two-cell
  continuation hang

#### Scenario: Many long lines

- **WHEN** a conversation holds hundreds of long result lines
- **THEN** building the conversation stays proportional to the total number of
  characters

### Requirement: Show the scroll position

The conversation view SHALL show the reader's position inside the thread with a
one-cell scrollbar in the right margin of the body, spanning every body row. It
SHALL be hidden when the whole thread fits on screen, and its column SHALL stay
reserved so the body width does not depend on the content.

#### Scenario: Bar appears for a long thread

- **WHEN** a thread renders more lines than the body holds
- **THEN** the right margin shows a track on every body row and a thumb where the
  visible lines sit

#### Scenario: Bar is hidden for a short thread

- **WHEN** every line of the thread fits on screen
- **THEN** the margin is drawn empty, and the body keeps the width it has for a
  long thread

#### Scenario: Thumb length follows the thread

- **WHEN** the thread is ten times taller than the body
- **THEN** the thumb is about a tenth of the track, never shorter than one cell
  and never longer than the track

#### Scenario: Thumb position follows the scroll

- **WHEN** the reader scrolls to the last line of the thread
- **THEN** the thumb ends at the bottom row, and it sits on the first row while
  the view is at the top

#### Scenario: No bar while rendering

- **WHEN** a thread is still being rendered in the background
- **THEN** no thumb is drawn until the line count is known

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

### Requirement: Colour prose by the role that wrote it

The conversation view SHALL colour the text of a message by its role, using one
hue per voice: prompts in cyan and replies in green. Text from any other role
SHALL stay in the default foreground. Reasoning SHALL stay muted, summaries
SHALL stay muted, and the colour SHALL never be the only sign of a role.

#### Scenario: Prompt prose

- **WHEN** a user message holds a text block
- **THEN** its body prints in cyan, which differs from the default foreground

#### Scenario: Reply prose

- **WHEN** an assistant message holds a text block
- **THEN** its body prints in green

#### Scenario: Other roles stay neutral

- **WHEN** a tool, system, or title message holds a text block
- **THEN** its body prints in the default foreground

#### Scenario: Colour follows wrapping

- **WHEN** a prompt wraps over five lines
- **THEN** every one of those lines carries the cyan style

#### Scenario: Reasoning and summaries stay quiet

- **WHEN** a message holds a thinking block or a summary block
- **THEN** the text prints muted and is not coloured by the role

### Requirement: Colour diffs and tool calls by their structure

The conversation view SHALL colour added diff lines green, removed diff lines
red, and hunk headers cyan, and SHALL colour a tool call by its parts: the verb
in the title accent, the path in bold, and a command in yellow. Tool results and
file payloads SHALL keep the default foreground so that structure stands out and
data does not.

#### Scenario: Diff lines

- **WHEN** a patch or edit block holds an added line, a removed line, and a hunk
  header
- **THEN** the added line prints green, the removed line red, and the hunk
  header cyan

#### Scenario: Change counts

- **WHEN** an edit block reports `+12 -8`
- **THEN** the additions print green and the deletions red, and the other words
  in the head keep their own style

#### Scenario: Tool call parts

- **WHEN** a message holds a shell tool call with a command
- **THEN** the tool name prints in the title accent, the path in bold, and the
  line after `$` prints in yellow

#### Scenario: Payload stays neutral

- **WHEN** a tool result holds twenty lines of output
- **THEN** those lines print in the default foreground

#### Scenario: Errors stay red

- **WHEN** a tool result is marked as an error
- **THEN** its lines and the `!` prefix print in red instead of the default
  foreground

### Requirement: Select messages in a conversation

The conversation view SHALL offer a selection mode over messages, entered with
`v` and left with `esc`. A cursor SHALL move one message with `j` and `k`, ten
messages with `J` and `K`, and to either end with `g` and `G`. The selection
SHALL be the range from the message marked when the mode was entered to the
cursor, and SHALL be shown with a marker in the left margin rather than by
restyling the conversation.

#### Scenario: Enter the selection

- **WHEN** `v` is pressed while a thread is open
- **THEN** the message at the top of the viewport carries the cursor, and the
  status reports one message selected

#### Scenario: Move the cursor

- **WHEN** `j` or `k` is pressed in selection mode
- **THEN** the cursor moves to the next or previous message, and the view
  scrolls so the cursor stays visible

#### Scenario: Move ten messages

- **WHEN** `J` or `K` is pressed in selection mode
- **THEN** the cursor moves ten messages, stopping at the first or last one

#### Scenario: Jump to the ends

- **WHEN** `g` or `G` is pressed in selection mode
- **THEN** the cursor moves to the first or the last selectable message

#### Scenario: Margin marks the selection

- **WHEN** three messages are selected
- **THEN** the first line of each selected message carries `*`, the line of the
  cursor carries `>`, and no other row changes

#### Scenario: Leave without copying

- **WHEN** `esc` is pressed in selection mode
- **THEN** the conversation returns to plain scrolling with the same scroll
  position

#### Scenario: Selection ends with the thread

- **WHEN** a thread is closed, or the message type selection changes
- **THEN** the selection is cleared

#### Scenario: Empty messages are not selectable

- **WHEN** a collapsed attachment run sits between two messages
- **THEN** the cursor skips it and moving through it does not select it

### Requirement: Copy a selection as Markdown

The copy action SHALL build its payload from the atuin message content, never
from the rendered lines, so wrapped text and commands come back unwrapped. It
SHALL prefix each message with a heading unless a single prose message is
copied, and SHALL put the payload on the clipboard through a native tool when
one is available and through an OSC 52 sequence otherwise. A payload larger
than 50 KB SHALL be written to a file instead of the clipboard.

#### Scenario: Copy two messages

- **WHEN** `y` is pressed with two messages selected
- **THEN** the clipboard holds both, each after its own heading with the time
  and the author

#### Scenario: Copy one reply

- **WHEN** one assistant message of prose is selected and `y` is pressed
- **THEN** the clipboard holds its text without a heading

#### Scenario: Long command stays one line

- **WHEN** a selected tool call holds a 300-cell command
- **THEN** the copied command is one line, not the wrapped lines of the view

#### Scenario: Tool call becomes a fence

- **WHEN** a selected message holds a tool call with a description and a
  command
- **THEN** the payload shows the tool name and the description, then a fenced
  block around the command

#### Scenario: Fence collision

- **WHEN** a copied block itself contains three backticks
- **THEN** the fence around it is longer than the backticks inside it

#### Scenario: Native clipboard tool

- **WHEN** a clipboard tool such as `wl-copy` is available and `y` is pressed
- **THEN** that tool receives the payload and the status reports the copy

#### Scenario: No clipboard tool

- **WHEN** no clipboard tool is available and `y` is pressed
- **THEN** the payload is sent as an OSC 52 sequence and the status names the
  terminal as the route

#### Scenario: Copy too large for a clipboard

- **WHEN** the payload is larger than 50 KB
- **THEN** it is written to a file and the status reports the path instead of
  claiming a copy

#### Scenario: Hidden types stay out

- **WHEN** tool calls are hidden and a selected range holds tool calls
- **THEN** the payload contains none of them

### Requirement: Export a selection or a thread to a file

The view SHALL write the selection to a file with `w` and the whole thread with
`e`. The file SHALL be named `att-<short id>-<first>-<last>.md` in the working
directory, SHALL be created only when it does not exist yet, and the status
line SHALL report the path.

#### Scenario: Write the selection

- **WHEN** `w` is pressed with messages 12 to 14 of `fd025611` selected
- **THEN** the file `att-fd025611-12-14.md` is written and its path is shown

#### Scenario: Export the whole thread

- **WHEN** `e` is pressed in the conversation view
- **THEN** every visible message of the thread is written to one file

#### Scenario: Existing file is not overwritten

- **WHEN** the target file already exists
- **THEN** nothing is written and the status says the file exists

#### Scenario: Written file follows the filter

- **WHEN** thinking is hidden and the thread is exported
- **THEN** the file holds no reasoning text

### Requirement: Limit the list to one directory or host

The list SHALL offer a directory scope toggled with `.` and a host scope toggled
with `@`. The directory scope SHALL keep only threads whose working directory
equals the process working directory. The host scope SHALL keep only threads
whose host equals the local host id read from `atuin info`. Both scopes SHALL
compose with each other, with the text filter, and with search results.

#### Scenario: Only this directory

- **WHEN** `.` is pressed in a project directory
- **THEN** the list shows only threads whose working directory equals that
  directory and the status line names the directory

#### Scenario: Only this host

- **WHEN** `@` is pressed and `atuin info` reported the local host id
- **THEN** the list shows only threads captured on that host and the status line
  says `host: local`

#### Scenario: Both scopes at once

- **WHEN** `.` and `@` are both on
- **THEN** only threads that match both are shown

#### Scenario: Scope clears

- **WHEN** the active scope key is pressed again
- **THEN** that scope is dropped and the list widens again

#### Scenario: Scope combines with the text filter

- **WHEN** a directory scope is active and `/` is given a query
- **THEN** the rows must match both the scope and the query, and the status
  reports the narrowed count

#### Scenario: Host id is unknown

- **WHEN** `atuin info` could not be read
- **THEN** `@` shows `host: unknown` and leaves the rows unchanged instead of
  emptying the list

#### Scenario: Empty result

- **WHEN** the scope matches no thread
- **THEN** the list says that no thread matches, and the scope key clears it
