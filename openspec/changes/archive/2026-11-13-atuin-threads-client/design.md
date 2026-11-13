# Design: atuin-threads-client

## Why the CLI and not the database

`threads.db` is atuin's internal store; its schema changes with atuin releases
and it is written by a running daemon. `atuin threads` is the supported
interface, and `--style json` gives the same fields the pretty output uses.

Costs accepted: one process per request, and atuin must be installed. Both are
fine for an interactive reader, and the process boundary keeps the Go code free
of cgo and SQLite schema assumptions.

## Types

`Thread` mirrors the list document:

```
harness, session_id, atuin_id, host, cwd, git_branch, model,
started_at, updated_at, message_count, tokens{...}, title, preview,
last_reply
```

`Conversation` mirrors `show`: `session` is the same thread document and
`messages` is a list of `{id, role, timestamp, cwd, content[], stop_reason}`.
Content blocks are `text`, `thinking`, `tool_call` (name, input), and
`tool_result` (content, is_error, tool_use_id).

## Execution seam

```go
type Runner func(ctx context.Context, args ...string) ([]byte, error)

func New(opts ...Option) *Client
func WithBinary(path string) Option
func WithRunner(r Runner) Option
```

The default runner is `exec.CommandContext` and includes stderr in its error.
Tests inject a runner that returns fixture JSON, so no atuin installation is
needed in the sandbox.

## Errors

Decode failures and command failures both wrap a context string and remain
errors for the UI to show; no silent fallback to an empty list, because an
empty list and a broken atuin look the same to a user.

## Out of scope

Caching, streaming capture, `tail`, and resume.
