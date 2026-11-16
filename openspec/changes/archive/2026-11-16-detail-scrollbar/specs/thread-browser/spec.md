## ADDED Requirements

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
