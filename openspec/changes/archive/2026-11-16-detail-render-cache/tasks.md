# Tasks: detail-render-cache

## 1. Render structure

- [x] 1.1 Add `detailBlock` and `detailRendering` with a block start index, a
      running total, and a `slice(from, count)` that copies only the blocks
      overlapping the range
- [x] 1.2 Turn the body of `detailLines()` into a builder that fills a
      `detailRendering`, and keep `detailLines()` as the flattening wrapper for
      existing callers and tests
- [x] 1.3 Keep one block per message, with the trailing blank line inside the
      block, and one block for a collapsed run

## 2. Model ownership

- [x] 2.1 Store the rendering on the model together with a `built` flag, the
      width it was made for, and a render counter
- [x] 2.2 Build synchronously on `conversationMsg` when the thread has at most
      500 messages, otherwise start a background command and set a rendering
      flag
- [x] 2.3 Accept a background result only when its width matches the current
      terminal width, and clamp the scroll offset afterwards
- [x] 2.4 Invalidate and rebuild on a width change while a thread is open
- [x] 2.5 Drop the rendering when the thread is closed

## 3. Drawing from the cache

- [x] 3.1 Draw the detail body from `cache.slice(detailScroll, bodyHeight())`
- [x] 3.2 Clamp the scroll offset against the cached total
- [x] 3.3 Report `line X of Y` in the status line from the cached total, and
      show a rendering state while a background render is pending
- [x] 3.4 Keep a synchronous build as the fallback when no rendering exists and
      none is pending

## 4. Cheaper measurement

- [x] 4.1 Measure runes with `runewidth.RuneWidth` in `splitWidth` and
      `expandTabs`
- [x] 4.2 Track the line width incrementally in `wrap` instead of measuring the
      growing line for every word

## 4b. Input

- [x] 4b.1 Split a key message that carries several runes into one key per
      rune outside the filter, so a fast scroll is not dropped

## 5. Verify

- [x] 5.1 Tests: render counter stays at 1 across 20 scroll keys, cache slice
      matches the flat rendering including a viewport that starts inside a
      block, resize rebuilds, stale result ignored, large thread goes through
      the rendering state
- [x] 5.2 `go test ./...`, `bash scripts/gate.sh`, and `nix flake check`
- [x] 5.3 Time `detailLines()` on an 11,520-message conversation before and
      after, and time a keystroke in a live tmux session
