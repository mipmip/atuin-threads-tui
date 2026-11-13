# Tasks: atuin-threads-client

## 1. Types

- [x] 1.1 Add `Thread`, `Tokens`, `Message`, `Content`, and `Conversation` with JSON tags
- [x] 1.2 Add the `Source` interface used by the UI

## 2. Client

- [x] 2.1 Implement `New` with `WithBinary` and `WithRunner` options
- [x] 2.2 Implement `List` via `atuin threads list --style json`
- [x] 2.3 Implement `Show` via `atuin threads show --style json <id>`
- [x] 2.4 Implement `Search` via `atuin threads search --style json --limit <n> -- <words>`
- [x] 2.5 Wrap command and decode failures with operation context

## 3. Tests

- [x] 3.1 Fixture tests for list, show, and search decoding
- [x] 3.2 Tests for empty list, non-zero exit, malformed JSON, and unknown fields
- [x] 3.3 Test that the default runner reports a missing binary
