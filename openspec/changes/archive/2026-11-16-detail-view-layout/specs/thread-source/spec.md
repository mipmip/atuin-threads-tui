# Spec Delta: thread-source

## MODIFIED Requirements

### Requirement: Show a thread

The client SHALL read a thread and its messages by running `atuin threads show
--style json <id>` and decoding the conversation document, including the token
usage, model, and working directory of each message and the files touched by a
patch block.

#### Scenario: Messages and content blocks

- **WHEN** atuin prints a conversation with user, assistant, and tool content
- **THEN** the client returns the session metadata and every message with its
  role, timestamp, and decoded content blocks

#### Scenario: Message usage and patch files

- **WHEN** a message carries a token usage object and a patch block lists a file
- **THEN** the client returns the usage on the message and the file path with its
  additions and deletions on the block
