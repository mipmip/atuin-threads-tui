## MODIFIED Requirements

### Requirement: Select messages in a conversation

The conversation view SHALL offer a selection mode over messages, entered with
`v` and left with `esc`. A cursor SHALL move one message with `j` and `k`, ten
messages with `J` and `K`, and to either end with `g` and `G`. `space` SHALL
toggle the message under the cursor in or out of a marked set, so that a few
messages can be picked without including the ones between them. Copying and file
export SHALL act on the marked set in document order. The mode SHALL show its
state by tinting whole messages and by a stripe in the left margin, and SHALL NOT
rebuild the conversation to draw it.

#### Scenario: Enter the selection

- **WHEN** `v` is pressed while a thread is open
- **THEN** the message at the top of the viewport carries the cursor, no message
  is marked yet, and the status reports nothing marked

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

#### Scenario: Space marks a message

- **WHEN** `space` is pressed on an unmarked message
- **THEN** that message is marked, and pressing `space` on it again unmarks it

#### Scenario: Margin marks the selection

- **WHEN** two messages are marked and the cursor sits on a third
- **THEN** every visible line of the marked messages carries a `#` stripe and the
  marked tint, and every visible line of the cursor message carries a `>` stripe
  and the cursor tint

#### Scenario: Marks survive the cursor

- **WHEN** the cursor moves away from a marked message
- **THEN** the message stays marked and keeps its tint and stripe

#### Scenario: Copy takes the marked messages

- **WHEN** `y` is pressed with three messages marked out of five
- **THEN** those three messages are copied in document order, without the two
  messages between or around them

#### Scenario: Nothing is marked

- **WHEN** `y` is pressed in selection mode with nothing marked
- **THEN** a notice says to mark a message first and no clipboard payload is
  written

#### Scenario: Leave without copying

- **WHEN** `esc` is pressed in selection mode
- **THEN** the conversation returns to plain scrolling with the same scroll
  position and the marked set is dropped

#### Scenario: Selection ends with the thread

- **WHEN** a thread is closed, or the message type selection changes
- **THEN** the marked set is cleared

#### Scenario: Empty messages are not selectable

- **WHEN** a collapsed attachment run sits between two messages
- **THEN** the cursor skips it and moving through it does not mark it

## ADDED Requirements

### Requirement: Separate messages with a background

An open conversation SHALL tint the background of each message so that message
boundaries are visible, alternating between two dark tints. The tint SHALL cover
the message body only, SHALL follow a message taller than the viewport, and SHALL
survive the colours already inside a line.

#### Scenario: Alternating tints

- **WHEN** a thread with several messages is open
- **THEN** consecutive messages alternate between two dark background tints

#### Scenario: Tint follows the message

- **WHEN** a message is taller than the viewport
- **THEN** every visible line of that message carries the same tint, and the
  next message starts on the other tint

#### Scenario: Tint stops at the body

- **WHEN** the conversation is painted at 120 columns
- **THEN** the tint covers the body column only, and the gap, the meta rail and
  the scrollbar keep their own background

#### Scenario: Tint survives inner colour

- **WHEN** a line inside a message has its own colour, such as a green diff
  addition
- **THEN** the background tint still covers that line to the end of the body
  column

#### Scenario: Tint in the plain reading view

- **WHEN** a thread is open and selection mode has not been entered
- **THEN** the alternating tints are already visible
