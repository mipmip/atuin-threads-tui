# Tasks: detail-row-select

## 1. Paint a background onto a rendered line

- [x] 1.1 Add a helper that pads a line to a given cell width and wraps it in a background, re-applying the background after every inner reset
- [x] 1.2 Test the helper with a plain line, a line that contains its own colour and reset, a line shorter than the width, and an empty line
- [x] 1.3 Confirm the painted width in cells equals the requested width, including a line with a wide character

## 2. Alternate the background per message

- [x] 2.1 Paint consecutive messages with alternating background tints in the open conversation
- [x] 2.2 Leave the thread head, the gap, the meta rail and the scrollbar untouched
- [x] 2.3 Test that two consecutive messages get different tints and that a message taller than the viewport keeps one tint
- [x] 2.4 Test that a coloured line keeps its background to the end of the body column

## 3. Mark a set of messages

- [x] 3.1 Replace the anchor and cursor range with a marked set of message indices and keep the cursor
- [x] 3.2 Bind `space` to mark or unmark the message under the cursor
- [x] 3.3 Enter the mode with nothing marked
- [x] 3.4 Clear the marked set when the thread closes, when the mode is left, and when the message types change
- [x] 3.5 Keep marks when the cursor moves, and clamp the cursor as before
- [x] 3.6 Test marking, unmarking, marks surviving cursor movement, and the clearing cases

## 4. Draw the marks

- [x] 4.1 Paint the cursor message with the cursor tint and a `>` stripe on every visible line
- [x] 4.2 Paint marked messages with the marked tint and a `#` stripe on every visible line, with the cursor winning when both apply
- [x] 4.3 Report the marked count in the status line and keep the line position
- [x] 4.4 Update the selection hints to mention `space`
- [x] 4.5 Test the stripes, the tints and the status text at 120 and 80 columns

## 5. Copy and export the marked set

- [x] 5.1 Make `y` copy the marked messages in document order
- [x] 5.2 Make `w` write the marked messages to a file named after the first and last marked message
- [x] 5.3 Show a notice and touch nothing when `y` is pressed with nothing marked
- [x] 5.4 Test that three marked messages out of five are copied without the ones between them

## 6. Verify

- [x] 6.1 Run `gofmt`, `go test ./...` and `bash scripts/gate.sh`
- [x] 6.2 Run `git add -A`, then `nix flake check` and `nix build`
- [x] 6.3 Capture frames at 120x24 and 80x24 with the cursor on one message and two messages marked, and compare them with the mockup in design.md
