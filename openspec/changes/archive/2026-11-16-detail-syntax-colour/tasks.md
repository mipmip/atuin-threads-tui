# Tasks: detail-syntax-colour

## 1. Styles and helpers

- [x] 1.1 Add `promptStyle`, `replyStyle`, `commandStyle`, `additionStyle`,
      `deletionStyle`, and `hunkStyle` next to the existing palette in
      `internal/ui/view.go`
- [x] 1.2 Add `roleProseStyle(role)` returning the prompt style, the reply
      style, or none
- [x] 1.3 Add `styleDiffLine(line, style func(...string) string) string` that
      classifies a diff line by its first characters, with `+++` and `---`
      treated as headers and left muted
- [x] 1.4 Test the classifier: addition, deletion, hunk, both file headers, an
      indented plain line, an empty line

## 2. Prose by role

- [x] 2.1 Style the `text` case of `blockLines` with the role style, applied
      after `wrap`
- [x] 2.2 Keep thinking, summary, and unknown blocks on their current style
- [x] 2.3 Test that a prompt is cyan, a reply green, a tool-role text block
      default, and that a five-line prompt carries the style on every line

## 3. Diffs

- [x] 3.1 Colour patch diff lines through the classifier
- [x] 3.2 Colour the `- old` and `+ new` lines of an edit block
- [x] 3.3 Colour the `+N` and `-M` counts in the edit head and in the patch
      head without breaking truncation
- [x] 3.4 Test a patch with an addition, a deletion, and a hunk header, and test
      that the head counts differ in colour from the surrounding words

## 4. Tool calls

- [x] 4.1 Render the verb of a tool call in the accent style and a path in bold
- [x] 4.2 Colour the command after `$` in yellow, including its continuation
      lines
- [x] 4.3 Keep the result payload unstyled and the error result red
- [x] 4.4 Test a shell head, a read head with a path, a plain result with no
      SGR, and an error result

## 5. Verify

- [x] 5.1 Confirm the build time on a large thread did not regress noticeably
- [x] 5.2 Run `gofmt`, `go test ./...` and `bash scripts/gate.sh`
- [x] 5.3 Run `nix flake check` and `nix build` after `git add -A`
- [x] 5.4 Capture the detail view with a diff and a prompt at 120x36 and check
      the colours against this design
