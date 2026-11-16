## ADDED Requirements

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

## MODIFIED Requirements

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
