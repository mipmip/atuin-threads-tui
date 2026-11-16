# Tasks: detail-view-layout

## 1. Decode the missing fields

- [x] 1.1 Add token usage and model to `Message`, and a file list with path,
      change, additions, deletions, and diff to `Content`, then add a client test
      that decodes a fixture carrying both and verify `go test ./internal/atuin`
      passes

## 2. Two column frame

- [x] 2.1 Add the detail layout that derives body and meta widths from the
      terminal width (body 80 capped, meta 20, gap 2, meta dropped below 44
      usable cells), and verify a test asserts body 80 at 120, body 54 at 80,
      and body taking the full width at 40
- [x] 2.2 Render each message as body lines beside a two line meta (`HH:MM -
      author`, `kind - N tokens`), keeping one blank line between messages, and
      verify a test compares the rendered frame at 120 columns with the design
      mockup line by line

## 3. Message kinds and tokens

- [x] 3.1 Derive the kind label from the block types in the message and pick the
      token count by kind (reasoning for thinking, output otherwise, nothing
      without usage), and verify a table driven test covers text, thinking, tool,
      patch, and a message with no usage
- [x] 3.2 Collapse runs of messages with no displayable content into one counted
      row, and verify tests cover twelve attachments in a row, a mixed
      attachment and ai-title run, and a single empty message that stays printed

## 4. Tool calls and results

- [x] 4.1 Render tool calls from their input fields for Bash, shell, write, edit,
      read, Skill, and compress, with `key=value` for unknown tools, and verify a
      test asserts no rendered tool line contains a brace, a quote, or a
      backslash escape
- [x] 4.2 Render tool results with their indentation and layout preserved,
      expanding tabs and wrapping only over-long lines, and verify a test asserts
      a YAML block keeps its relative indentation and stays inside the body width

## 5. Blocks that never rendered

- [x] 5.1 Render patch blocks as path with additions and deletions plus the diff
      body, and summary blocks as prose, and verify tests assert both appear in
      the rendered frame where today they produce nothing

## 6. Prose that measures correctly

- [x] 6.1 Switch prose wrapping to display width and fill to the body column,
      and verify a test with CJK and combining characters asserts every rendered
      line measures at most the body width

## 7. Integration

- [x] 7.1 Run `bash scripts/gate.sh` and verify the coverage thresholds still
      pass with the new tests in place
- [x] 7.2 Capture the detail view at 120x36 and 80x24 against local atuin data
      and verify the frames match `design.md` for column offsets, meta lines, and
      collapsed runs

## Workflow follow-up

- Archive the change with `./scripts/ship-change.sh detail-view-layout` once the
  gate passes.
