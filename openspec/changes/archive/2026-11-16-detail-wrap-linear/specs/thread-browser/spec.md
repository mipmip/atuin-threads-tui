## ADDED Requirements

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
