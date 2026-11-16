# Tasks: detail-type-filter

## 1. Categories and the mask

- [x] 1.1 Add a `typeMask` with bits for prompts, replies, thinking, calls,
      results, patches, summaries, and notes, plus an `allTypes` default and a
      `count()` helper
- [x] 1.2 Add `categoryOf(block, role)` mapping `text` by role into prompts or
      replies, `error` into results, and any unknown type into a category that
      is always shown
- [x] 1.3 Add `selected(mask)` and `without(mask)` helpers for the panel draft
- [x] 1.4 Test every block kind, text on a user and on an assistant message, an
      error block, and an unknown type

## 2. Filter the conversation while building

- [x] 2.1 Pass the mask into `buildDetail` and drop blocks whose category is not
      selected
- [x] 2.2 Skip a message when no block survives, so no hollow meta row is left
- [x] 2.3 Render empty messages only while notes is selected, keeping the
      collapsed run for consecutive ones
- [x] 2.4 Build the second meta line from the visible kinds only
- [x] 2.5 Test a mixed message, a tool-only message that vanishes, a thinking
      only message, notes on and off, and the meta line of a partly hidden
      message

## 3. Panel state and keys

- [x] 3.1 Add a `modeTypes` with a draft mask and the counts captured when the
      panel opens
- [x] 3.2 Open the panel with `t` from the detail view, close it with `esc`
      without applying
- [x] 3.3 Toggle the marked row with space, select all with `a`, and select the
      conversation with `n`
- [x] 3.4 Apply the presets `1` to `4` and apply the draft on `enter`
- [x] 3.5 Keep the applied mask in the model so the next thread uses it
- [x] 3.6 Test opening, toggling, cancelling, applying, each preset, and the
      mask surviving a thread change

## 4. Draw the panel, status and hints

- [x] 4.1 Render the panel with the category rows, their counts, the current
      marks, and the shown-of-total count in its header
- [x] 4.2 Add `t types` to the detail hints and a panel hint line with the keys
- [x] 4.3 Report `N of M msgs` in the detail status and keep the plain count
      when everything is selected
- [x] 4.4 Show `no messages of the selected types` when nothing renders
- [x] 4.5 Test the panel lines, the hint lines, both status forms, and the empty
      body

## 5. Cache key, anchor and background render

- [x] 5.1 Include the mask in the render cache key and reject a background
      result whose mask no longer matches
- [x] 5.2 Record the first message index of each block and scroll to that
      message after a rebuild
- [x] 5.3 Fall back to the next visible block, and to the top, when the anchored
      message is hidden
- [x] 5.4 Test that applying a mask rebuilds once, that the top message stays on
      top, that a hidden anchor falls through, and that a stale mask result is
      dropped

## 6. Verify

- [x] 6.1 Run `gofmt`, `go test ./...` and `bash scripts/gate.sh`
- [x] 6.2 Run `nix flake check` and `nix build` after `git add -A`
- [x] 6.3 Capture the detail view at 120x36 and 80x24 with the conversation
      preset and compare the counts and the meta column
