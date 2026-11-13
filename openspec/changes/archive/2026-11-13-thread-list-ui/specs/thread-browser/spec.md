# thread-browser

## ADDED Requirements

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

The list SHALL keep a selection within bounds and support `j` and down for the
next thread, `k` and up for the previous thread, `g` for the first thread, and
`G` for the last thread.

#### Scenario: Move past the end

- **WHEN** the last thread is selected and `j` is pressed
- **THEN** the selection stays on the last thread

#### Scenario: Empty list navigation

- **WHEN** no threads are loaded and `j` or `k` is pressed
- **THEN** nothing changes and no error is shown

### Requirement: Read a thread

The list SHALL open the selected thread with `enter` and show its title, harness,
model, working directory, update time, and every message.

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

### Requirement: Scroll a conversation

The detail view SHALL scroll by one line with `j` and `k`, by one page with
`pgdown` and `pgup`, and jump with `g` and `G`, without scrolling past either
end.

#### Scenario: Scroll past the last line

- **WHEN** the view is at the last line and `j` is pressed
- **THEN** the offset does not grow

### Requirement: Adapt to terminal size

The program SHALL re-render rows and wrapped message text when the terminal is
resized, and SHALL NOT panic when the terminal is smaller than the fixed
chrome.

#### Scenario: Very short terminal

- **WHEN** the window height is 2
- **THEN** the view renders without panicking

### Requirement: Quit

The program SHALL quit on `q` in the list and on `ctrl+c` anywhere.

#### Scenario: Quit from the list

- **WHEN** `q` is pressed in the list
- **THEN** the program returns a quit command
