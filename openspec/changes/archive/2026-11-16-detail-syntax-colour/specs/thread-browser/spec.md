# thread-browser

## ADDED Requirements

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
