# Spec Delta: thread-browser

## MODIFIED Requirements

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
