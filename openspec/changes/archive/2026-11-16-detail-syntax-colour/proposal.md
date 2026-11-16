# Proposal

## Why

The conversation view draws almost everything in the terminal's default
foreground. A reader scanning a long thread cannot tell at a glance whether a
line is their own prompt, the model's reply, a command, or a changed line in a
diff; only the small author label carries colour. Diffs in particular are the
part people read closely and they are currently monochrome.

## What Changes

- Colour prose by the role that wrote it: prompts in cyan, replies in green,
  everything else neutral.
- Colour diffs by line shape: added lines green, removed lines red, hunk
  headers cyan, and the addition and deletion counts in the edit head in the
  same colours.
- Colour tool calls structurally: the verb in the title accent, the path bold,
  the command after `$` in yellow.
- Keep tool results, file dumps, and other payload in the default colour, and
  keep chrome, rail, hints, and status neutral.
- Keep thinking and summaries muted, and keep errors red.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the conversation view gains colour by role and by block
  kind, with one hue per job.

## Impact

`internal/ui/detail.go` (block and line styling), `internal/ui/view.go` (the
palette), tests in `internal/ui`. No new dependency, no lexer, no change to
`internal/atuin` or to stored data.

## Bean

`atuin-threads-tui-7xhv`
