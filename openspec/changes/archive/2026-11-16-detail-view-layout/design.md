# Design

## Context

See `proposal.md` for why. The current renderer flattens a message into
`[role] HH:MM` followed by one indented chunk per content block, then a blank
line. Four properties of the data shape the replacement:

- Token usage lives on the message, not the block. All blocks of one assistant
  turn repeat the same usage object, and `tool` messages carry none.
- `attachment` (425 across 15 threads), `ai-title` (198), and most `system` (141)
  messages decode to empty content.
- opencode writes the literal text `Reasoning · 42 tokens` as its thinking block.
- `patch` and `summary` blocks exist in the data but `blockLines` returns nil for
  both, so 329 blocks never render.

The current `wrap` runs `strings.Fields` over a paragraph, which reflows prose
well but destroys the indentation of YAML, diffs, and code, and measures with
`utf8.RuneCountInString`, which breaks on wide characters.

## Goals / Non-Goals

**Goals:**

- A reading column that stays the same width on every terminal, with message
  annotations parked to its right.
- Every message kind renders something: no silent drops, no raw JSON.
- Line count in a long thread drops because empty and duplicated rows collapse.

**Non-Goals:**

- Changing the list view, the search, or the key map.
- Syntax highlighting, diff colouring, or word-wrap at word boundaries inside
  code blocks.
- A toggle for collapsed runs. If it is ever wanted it becomes its own
  requirement.

## Decisions

### Body left, meta right

```
/clear                                                                            11:15 - user
                                                                                  text

I want a new script that starts the frontend and the backend                      11:53 - user
together, following the conventions in dev-be.sh.                                 text

(12 attachments)                                                                  11:53 - sys


Bash  Inspect dev scripts and README                                              11:54 - asst
  $ cd ~/chatcms-app; ls scripts; cat scripts/dev-fe.sh                           tool - 265 tokens

acceptance-against.sh                                                             11:54 - tool
check-specs.sh                                                                    out

Edit  scripts/dev-fe.sh  +12 -3                                                   11:56 - asst
  - npm run dev                                                                   patch - 91 tokens
  + npm run dev -- --host 0.0.0.0

Summarizing 18 earlier messages                                                   11:57 - asst
                                                                                  summ - 40 tokens
```

Body 80 cells, gap 2, meta 24, total 106. At 120 columns the last 14 cells stay
empty rather than stretching the body.

A left rail was the first alternative and it reads worse: every body line starts
at the same indent either way, but a left rail pushes the text off-centre and
forces the eye to jump back across it at the start of each line. On the right the
annotation is read after the line, which matches how it is used.

### Width arithmetic

`meta` is 24 whenever `width >= 46`, else the meta column is dropped and the
body takes the full width. `body = min(80, width - 26)`.

At 80 columns: body 54, meta 24. At 100: body 74. At 120 and wider: body 80.

The meta needs 24 rather than 20 because a joined kind with a four digit count
reads `think+tool - 2657 tokens`, which is 23 cells.

### Meta lines

Line one is `HH:MM - author`, where author is `user`, `asst`, `sys`, `tool`, or
the raw role for anything unknown. Line two is `kind - N tokens`, or just `kind`.

`kind` is the message's shape, not one block: the set of block types in the
message joined with `+`, in the order they appear, using short names (`text`,
`think`, `tool`, `out`, `patch`, `summ`).

Token count follows the kind: `reasoning` when the message is thinking only,
`output` otherwise, and nothing when the message has no usage object. Showing
input, output, cache, and reasoning side by side would need a table the rail
cannot hold, and `cache_read` in the tens of thousands would dominate the line.

One meta per message rather than per block, because usage is per turn. A
per-block meta would print the same number four times per turn.

Rejected: right-aligning the token number in the meta. It costs alignment of the
`kind` words for a digit column that changes width anyway.

### Blank line between messages

Kept. The time line and the blank line both delimit, which makes a thread of
short messages read as a sequence rather than a wall. A one-line message costs
three rows, which is the price of the rhythm.

### Collapsing empty runs

A message is collapsible when none of its blocks produce a line: empty content,
attachment, ai-title, and empty system messages. A run of two or more such
messages becomes a single row with `(N attachments)` or `(N empty messages)`
when the roles differ. A single empty message is still printed, because a lone
`(1 attachment)` reads as a bug.

The row uses the meta of the first message in the run.

### Tool calls render from fields

```
+------------+-----------------------------------+---------------------------------+
| tool       | first line                        | payload                         |
+------------+-----------------------------------+---------------------------------+
| Bash shell | Bash  <description>               | $ <command, up to 4 lines>      |
| write      | write <path>  +128 lines          | first 3 content lines, then ... |
| edit       | edit <path>  +12 -3               | - old line / + new line, capped |
| read       | read <path>  40..80               | none, the result carries it     |
| Skill      | Skill <name>                      | args, one line                  |
| compress   | compress <topic>                  | none                            |
| other      | <Name>  key=value key=value +3    | values cut at 40 cells          |
+------------+-----------------------------------+---------------------------------+
```

No braces, no quotes, no escapes. Tool names are matched case-insensitively and
mapped to the same renderer, so `Bash` and `shell` behave alike. Unknown tools
fall back to `key=value` pairs rather than JSON, which keeps the guarantee that
no block ever prints a raw document.

`Edit`/`edit` inputs arrive with different field names (`file_path`/`path`,
`old_string`/`oldString`), so the renderer looks for either.

### Whitespace-preserving results

Tool results and patch diffs use a second renderer: split on newlines, keep each
line's leading spaces, expand tabs to the next multiple of 4, and wrap only lines
wider than the body with a hanging indent of two spaces. The existing `... N more
lines` cap stays at 8.

Prose (`text`, `summary`, `thinking`) keeps the reflow that fills the line, but
measures with display width and keeps blank lines as paragraph breaks.

Both paths truncate with `ansi.Truncate`, so a wide character is never split.

### Data

`Message` gains `Tokens` and `Model`; `Content` gains `Files` with `Path`,
`Change`, `Additions`, `Deletions`, and `Diff`. Unknown fields stay tolerated by
the decoder as they are today.

## Risks / Trade-offs

- A 24-cell meta on an 80-column terminal leaves 54 for prose → accepted; the
  alternative is dropping the meta below 100 columns and losing the annotation
  exactly where rows are hardest to scan.
- Collapsing empty runs hides that a thread had 12 attachments → the count is
  shown, and only content-free messages collapse.
- A single very long token (URL, base64 blob) still hard-splits → it splits at the
  body edge instead of the terminal edge, so the damage is contained.
- `output` is 0 for some usage objects while `reasoning` is set → the kind rule
  keeps a `0 tokens` from appearing on a thinking message, but a text message
  with only cache usage would show `0`. Accepted until it is observed.

## Migration Plan

None. The view is read-only and the decoder change is additive.

## Open Questions

None that change the specs. Whether a collapsed run should be expandable is
deliberately deferred; it becomes a new requirement if it is ever wanted.
