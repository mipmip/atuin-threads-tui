# Proposal

## Why

The conversation view is a single column of flattened text. Tool calls arrive as
raw JSON wrapped mid-token, runs of empty attachment messages take two rows each
(`fd025611` spends 12% of its 408 lines on 12 empty attachments), opencode's
thinking blocks repeat the placeholder `Reasoning · 42 tokens`, and `patch` and
`summary` blocks are dropped without a trace. Reading a long thread means
scrolling past noise to find the message that mattered.

## What Changes

- Give the conversation view two columns: the body at 80 cells on the left, and a
  meta column of 20 cells on the right, behind a 2-cell gap. Space beyond the
  block stays empty at every terminal width.
- Print the meta column as two lines per message: `HH:MM - author`, then
  `type - N tokens`. Several blocks in one message join their kinds
  (`think+tool - 265 tokens`). Tokens appear only where they mean something:
  `reasoning` for thinking, `output` otherwise, and nothing for roles without
  usage.
- Keep one blank line between messages, so the time line and the gap both mark
  the boundary.
- Collapse a run of messages with no renderable content to one row, for example
  `(12 attachments)`.
- Render tool calls from their input fields instead of raw JSON: `Bash` with the
  description and a `$ command` line, `edit <path> +12 -3` with old and new
  lines, `write <path> +128`, `read <path>`, `Skill <name>`, and `key=value`
  pairs for unknown tools.
- Preserve whitespace in tool results and wrap them by display width, so YAML,
  diffs, and code keep their indentation. Keep the line cap with a `... N more
  lines` marker.
- Start rendering `patch` blocks (path, additions, deletions, diff) and `summary`
  blocks.
- Soft wrap prose at 80 cells and measure text by display width instead of rune
  count.
- Decode `tokens` and `model` on a message, and `files` on a patch block.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the `Read a thread` requirement gains the two column body and
  meta layout, message type and token labels, collapsed runs, tool call and
  result rendering, and the patch and summary blocks.
- `thread-source`: the `Show a thread` requirement now promises token usage and
  model on each message and the file list on patch blocks.

## Impact

- `internal/ui/view.go`: `detailLines`, `blockLines`, `roleLine`, `wrap`,
  `capped`.
- `internal/atuin/types.go`: `Message.Tokens`, `Message.Model`, `Content.Files`.
- Tests in `internal/ui` and `internal/atuin` for the new rows and labels.
- No new dependencies, no flake change.
