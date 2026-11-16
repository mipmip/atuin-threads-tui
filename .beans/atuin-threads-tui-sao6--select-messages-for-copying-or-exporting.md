---
# atuin-threads-tui-sao6
title: select messages for copying or exporting
status: completed
openspec-link: openspec/changes/archive/2026-11-16-detail-message-selection
type: feature
priority: normal
created_at: 2026-11-16T14:51:33Z
updated_at: 2026-11-16T18:42:20Z
parent: atuin-threads-tui-s9te
---

## Summary of Changes

Press `v` in a conversation to select a range of messages; `j`/`k` move one message, `J`/`K` ten, `g`/`G` jump to the ends. Selected messages carry `*` in the left margin and the cursor `>`; cached lines are never restyled. `y` copies the range as Markdown built from the atuin content (fenced commands and diffs, no re-wrapping): a native clipboard tool is tried first, then OSC 52, and the status names the route. `w` writes the selection to `att-<id>-<first>-<last>.md` with O_EXCL so nothing is overwritten, and `e` exports the whole thread. Payloads over 50 KB go to a file instead of the clipboard, and the export follows the active message type filter.
