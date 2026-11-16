# Design

## Context

The detail view is line based: `detailScroll` counts rendered lines, and
`detailRendering` holds `head` plus one `detailBlock` per message with its
`start`, its `msg` index, and its `lines`. `messageAt(line)` and
`scrollToMessage(msg)` already map between the two.

That is enough for a message cursor, but it shapes the implementation: the
cursor is an index into the message list, and the marker has to be drawn on
cached lines at draw time.

Two constraints decide most of the rest:

- Rendered lines are hard wrapped at 80 or 92 cells and carry SGR codes, so
  they cannot be the export source.
- OSC 52 has no acknowledgement. A terminal that drops it looks exactly like
  one that accepted it, so the status cannot honestly claim a copy on that
  route.

## Mode and keys

```
 browse                     select (v)                  result
 ------                     ----------                  ------
 j k    line                j k    message              y  clipboard
 space  page                J K    ten messages         w  file (selection)
 t      types               g G    first / last         e  file (whole thread)
 esc    list                y      copy
                            w      write
                            esc    back to browse
```

Selection mode keeps the conversation on screen. `v` marks the message at the
top of the viewport as the anchor and puts the cursor there, so entering and
leaving never moves the reading position.

The cursor is a message index, not a line. Moving it to a message outside the
viewport scrolls with `scrollToMessage`, which already exists for the rebuild
anchor.

## Drawing the cursor

Cached lines are never restyled. The margin cell already exists for every
content row, so the marker goes there:

```
 j/k scroll     ->     v, then j j
                        |
 +--------------------+ |
 | 11:53 - user       | |
 | Read the script.   | |
 |                    | |
 | 11:53 - asst       | |
 | think - 83 tokens  | |
 +--------------------+ |
                        v
 * 11:53 - user
   Read the script for all options.
 * 11:53 - asst
   think - 83 tokens
 > 11:53 - tool
   out
```

`*` marks the first line of a selected message, `>` the cursor, and rows that
belong to no message keep the blank margin. Only the first line of a message is
marked, so a tall selection does not turn into a striped block.

Status in select mode: `fd025611 · 3 selected of 799 msgs · line 120 of 13273`.
Hints replace the browse hints while the mode is active.

## Building the range

The selectable messages are the ones that produced a block:

- a message with no visible block under the current type filter is not in the
  list, so the cursor cannot land on something invisible
- a collapsed `(N attachments)` run is one block but not a selectable message,
  so it is skipped
- the thread head (`head` lines) is not selectable; `messageAt` returns -1 there
  and `v` then anchors on the first block

The list is `detailRendering.selectable()`, the block `msg` values in order,
computed when the rendering is built.

## Markdown payload

One generator, two callers. It takes `[]atuin.Message` plus the hidden mask and
writes to an `io.Writer`.

| block | output |
| --- | --- |
| heading | `### 14:05 assistant` |
| text | the text as it was stored, paragraphs and blank lines kept |
| thinking | a blockquote, one `> ` per line |
| tool_call | `**Bash** - <description>` then a fenced block with the command |
| tool_result | a fenced block with the output |
| error | a fenced block with the output, first line prefixed `! ` |
| patch | `**patch** <path> (+a -d)` then a fenced diff |
| summary | the summary text as a paragraph |

Rules:

- Headings are dropped when exactly one message is copied and it holds only
  prose, because a bare reply pastes cleanly into a chat.
- Prose is not wrapped. Terminals wrap, Markdown renderers wrap, and a stored
  newline in the middle of a sentence is a defect.
- Commands and output keep every original line and get a fence. The fence is
  one backtick longer than the longest run of backticks in the content, so
  pasted code cannot break out of it.
- Hidden categories are skipped, exactly as in the view.
- Messages are separated by one blank line.

## Clipboard ladder

Order, decided at startup with `exec.LookPath`:

```
 y
 |
 +-- payload > 50 KB ------------------> write file, report path
 |
 +-- wl-copy | pbcopy | xclip | xsel ---> run, exit 0 -> "copied N messages"
 |                                          non-zero -> fall through
 |
 +-- OSC 52 ---------------------------> "copied N messages (terminal)"
```

`wl-copy` forks and returns, so the exit status is real and the call does not
block. The OSC 52 sequence is written from a `tea.Cmd` to stdout: it is zero
width, so it cannot disturb the rendered frame, and it reaches the outer
terminal through tmux because `set-clipboard` is `external`.

Naming the route is the point of the ladder: `(terminal)` tells the reader that
the last step is best effort, which no amount of checking can improve.

## Files

- `att-<short id>-<first>-<last>.md`, with 1-based message indices, in the
  working directory: `att-fd025611-12-14.md`
- whole thread: `att-fd025611-1-799.md`, so the name still says what range it
  covers
- created with `O_EXCL`, so an existing file is reported instead of replaced
- the path in the status line is relative, because that is what the user typed
  around

## Decisions

- Message selection, not line selection: lines are a rendering artifact, and
  the unit a reader wants to copy is a message.
- Markers in the margin instead of restyling or rebuilding: a rebuild costs
  149 ms on a large thread, and restyling cached SGR lines is error prone.
- Markdown instead of plain text: fenced code is the whole reason the export
  exists, and Markdown degrades to readable plain text elsewhere.
- One generator for clipboard and file, so the two cannot drift.
- Native tool first, OSC 52 second: the only route with a real success signal
  wins when it exists.
- 50 KB ceiling: terminals and tmux silently truncate large OSC 52 payloads,
  and a truncated paste is worse than a file.
- No filename prompt: a text input mode for one field is not worth the new mode,
  and `-1-799` names are predictable.

## Risks

- A very large selection with a native tool still passes the whole payload on
  stdin; that is a normal exec, not a terminal limit, so it is allowed.
- The export shares block walking with the view but not code, so a new block
  type has to be added in two places. A test asserts that every content type
  either formats or is skipped deliberately.
- Cursor movement relies on the block table, so a rendering must exist. `v` is
  ignored while the background render is still running.
- The selection is by message index, so changing the type filter invalidates it;
  it is cleared rather than remapped.

## Testing

- selectable list, cursor clamping at both ends, ten-step jumps, marker rows
- `v` on a thread whose head is at the top, and while rendering
- Markdown for every block type, heading elision, fence bumping, unwrapped
  command, filter respected
- clipboard ladder with an injected runner: first tool, failing tool falls
  through, no tool uses OSC 52, size ceiling writes a file
- file naming, range in the name, refusal on an existing file, filtered export
- status and hint lines in select mode, and clearing on esc, thread change and
  type change
