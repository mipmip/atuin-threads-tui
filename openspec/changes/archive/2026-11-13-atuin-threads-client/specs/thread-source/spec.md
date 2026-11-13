# thread-source

## ADDED Requirements

### Requirement: List threads

The client SHALL list captured threads by running `atuin threads list --style
json` and decoding the array it prints.

#### Scenario: Threads are returned in atuin order

- **WHEN** atuin prints two threads, newest first
- **THEN** the client returns both threads in that order with title, harness,
  cwd, message count, token usage, and timestamps populated

#### Scenario: No threads

- **WHEN** atuin prints an empty array
- **THEN** the client returns an empty slice and no error

### Requirement: Show a thread

The client SHALL read a thread and its messages by running `atuin threads show
--style json <id>` and decoding the conversation document.

#### Scenario: Messages and content blocks

- **WHEN** atuin prints a conversation with user, assistant, and tool content
- **THEN** the client returns the session metadata and every message with its
  role, timestamp, and decoded content blocks

### Requirement: Search threads

The client SHALL search threads by running `atuin threads search --style json
--limit <n> <query>` and decoding the resulting thread array.

#### Scenario: Query with several words

- **WHEN** the query is `openspec beans`
- **THEN** the client passes each word as a separate argument and returns the
  matches

### Requirement: Report command failures

The client SHALL return an error naming the failed operation and including the
stderr text from atuin.

#### Scenario: atuin exits non-zero

- **WHEN** the runner reports a failure with stderr `no such command`
- **THEN** listing returns an error whose message contains `no such command`

### Requirement: Report malformed output

The client SHALL return a decode error when atuin prints something that is not
the expected JSON document.

#### Scenario: Truncated document

- **WHEN** atuin prints `[{`
- **THEN** listing returns an error mentioning decoding

### Requirement: Unknown fields are tolerated

The client SHALL ignore JSON fields it does not model, so a newer atuin does
not break the reader.

#### Scenario: Extra field

- **WHEN** a thread document contains `"future_field": true`
- **THEN** listing succeeds and returns the modelled fields
