---
# atuin-threads-tui-7xhv
title: syntax highlighting
status: completed
openspec-link: openspec/changes/archive/2026-11-16-detail-syntax-colour
type: feature
priority: normal
created_at: 2026-11-16T16:37:27Z
updated_at: 2026-11-16T16:58:28Z
parent: atuin-threads-tui-s9te
---

syntax highlighting and special font settings for the different types

- user color for user prompts
- assistant color for assitant replies
- diff syntax highlighting
- tool call syntax highlighting

## Summary of Changes

Conversation prose now takes the colour of its author: prompts are cyan, assistant
replies green, and other roles stay in the default colour. Reasoning and
summaries keep their muted style.

Diffs are coloured by line: additions green, deletions red, hunks cyan, while
`+++` and `---` file headers stay muted. Tool calls colour the verb and the
thread title in bold pink, paths in bold, the command after `$` in yellow, and
the `+N`/`-M` counts in green and red. Tool results and payloads stay in the
default colour, and errors stay red.

Styling is applied after wrapping, so no width or alignment changed.
