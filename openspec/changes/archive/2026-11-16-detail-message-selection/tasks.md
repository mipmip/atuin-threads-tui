# Tasks: detail-message-selection

## 1. Selectable messages

- [x] 1.1 Add `selectable()` to `detailRendering`, returning the message index
      of every block that carries a message, in order
- [x] 1.2 Test that skipped messages and collapsed attachment runs are absent
      and that the list follows the type filter

## 2. Select mode and cursor

- [x] 2.1 Add `modeSelect` with an anchor, a cursor, and the selectable list
      captured when the mode is entered
- [x] 2.2 Enter with `v` from a built rendering, anchored on the message at the
      top of the viewport, and leave with `esc` without moving the scroll
- [x] 2.3 Move with `j`/`k`, `J`/`K`, and `g`/`G`, clamped to the ends, and
      scroll the cursor into view
- [x] 2.4 Clear the selection when the thread closes or the type selection
      changes
- [x] 2.5 Test entering, moving, the ten-message jump, both ends, leaving, and
      the two clearing paths

## 3. Margin markers and chrome

- [x] 3.1 Mark the first row of each selected message with `*` and the cursor
      row with `>`, inside the existing margin cell, without touching cached
      lines
- [x] 3.2 Report `N selected of M msgs` and the line position in the status
      line
- [x] 3.3 Replace the detail hints with the select-mode hints
- [x] 3.4 Test the marked rows, an unmarked conversation, both status forms,
      and the hint lines

## 4. Markdown formatter

- [x] 4.1 Add a formatter that writes selected messages from their content,
      with a `### HH:MM author` heading per message
- [x] 4.2 Format prose without wrapping, and thinking as a blockquote
- [x] 4.3 Format tool calls as the tool name and description followed by a
      fenced command, and results, errors, patches, and summaries as fenced
      blocks
- [x] 4.4 Make the fence one backtick longer than the longest run inside the
      content
- [x] 4.5 Drop the heading when exactly one prose message is exported
- [x] 4.6 Skip categories hidden by the type filter
- [x] 4.7 Test every block type, heading elision, fence bumping, an unwrapped
      300-cell command, and a filtered export

## 5. Clipboard

- [x] 5.1 Detect `wl-copy`, `pbcopy`, `xclip`, and `xsel` once at startup, with
      an injectable runner for tests
- [x] 5.2 Copy with `y`: native tool first, OSC 52 when none is available or
      the tool fails
- [x] 5.3 Send a payload above 50 KB to the file writer instead, and say so
- [x] 5.4 Test the first tool, a failing tool that falls through, the OSC 52
      route, and the size ceiling

## 6. Files

- [x] 6.1 Write the selection with `w` and the whole thread with `e` to
      `att-<short id>-<first>-<last>.md`
- [x] 6.2 Create with `O_EXCL` and report an existing file instead of
      overwriting it
- [x] 6.3 Report the written path in the status line
- [x] 6.4 Test the name, the range in the name, the refusal, and a filtered
      whole-thread export

## 7. Verify

- [x] 7.1 Run `gofmt`, `go test ./...`, and `bash scripts/gate.sh`
- [x] 7.2 Run `nix flake check` and `nix build` after `git add -A`
- [x] 7.3 Capture a 120x36 frame in select mode with a marked range, and check
      the margin markers and the status line
