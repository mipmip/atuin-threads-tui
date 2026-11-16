# thread-browser

## ADDED Requirements

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
