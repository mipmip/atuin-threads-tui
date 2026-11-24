---
# atuin-threads-tui-ds02
title: the select functionality is not clea
status: completed
openspec-link: openspec/changes/archive/2026-11-24-detail-row-select
type: bug
priority: normal
created_at: 2026-11-24T22:02:10Z
updated_at: 2026-11-24T22:40:09Z
---

Selecting messages inside a session has unclear UI. clicking s give a small > sign at an unclear location. Meesages should appear as full selectable rows and it should be clear before selecting which row is selected. You should have a row-select-mode and then space should select. We want to select a few or we want to use the row position for other purposes

## Summary of Changes

The conversation view now tints messages with alternating backgrounds (234 and 236), so neighbouring messages are separated without borders. In select mode the cursor message is lit with background 240 and a pink `>` stripe in the margin, and every line of it carries the stripe, so a message hundreds of lines tall still reads as one selected block. `space` marks a message: marks are a set, not a range, so any number of messages can be picked, each shown with background 238 and a `#` stripe. `y` copies the marked messages in document order and `w` writes them to a file; with nothing marked the status says so. Selecting fully toward the top of a message is fixed too: a message taller than the viewport now scrolls to its first line instead of its last.
