# Design briefing: att, thread list and conversation view

## What att is

`att` is a terminal UI for browsing the threads that atuin records from shell
history. Atuin stores every agent session (Claude Code, opencode, Codex) as a
thread with a title, harness, model, working directory, timestamps and a full
message log. `att` reads that log through the `atuin threads` CLI and gives it a
place to be read: a list of threads, a conversation view, a live filter and a
full-text search.

The user is the person who ran those sessions. They come back days later to find
"that session where I fixed the camera" and read what happened, in order,
without opening the agent again.

## Screenshots

All screenshots are 120x36 cells, rendered from the real binary against real
local data (187 threads).

| File                             | Shows                              |
| -------------------------------- | ---------------------------------- |
| `screens/list.png`               | thread list, first row selected    |
| `screens/list-filter.png`        | filter active, "nix", 61 of 187    |
| `screens/detail.png`             | conversation view, tool output     |

The `.ansi` files next to them are the exact captured frames; the `.html` files
are the intermediate renders. Ignore those unless you want to re-render.

## What this session should produce

A design direction for the two main views, concrete enough to implement in
lipgloss:

1. Layout and hierarchy for the thread list. What the eye lands on first, how a
   row is composed, how selection reads, what to do with long titles.
2. Layout and hierarchy for the conversation view. How user, assistant,
   reasoning, tool call and tool result blocks are distinguished, and how a very
   long transcript stays readable.
3. A terminal-safe palette and typographic treatment (bold, dim, color) that
   works on dark and light backgrounds.
4. The filter and search states, including the transition between them.
5. Annotated sketches or annotated screenshots. Prose alone is not enough.

Deliver the design as a document in this folder, for example
`design/DIRECTION.md`, with the reasoning per view and a short list of
implementable changes ordered by impact. No code required in the design
session; implementation happens afterwards in an OpenSpec change.

## Hard constraints

- Terminal only, built with Go, bubbletea and lipgloss. No images, no web view.
- Must work at 80x24 minimum, designed against 120x36.
- Only 256-color ANSI indexes and bold/dim/underline. No background fills that
  assume a specific terminal theme, though tested on black.
- Keyboard only. Mouse support does not exist yet and is out of scope.
- Content comes from atuin's JSON, which the app does not control: titles can be
  empty, very long, or in any language; message counts reach 8000; tool output
  is unbounded.
- The body is padded to the terminal height, with the status line and hint line
  pinned to the bottom. Nothing scrolls the whole screen.

## Current screen anatomy

Both views share four vertical bands:

```
att · atuin threads                                                header
-----------------------------------------------------------------  body
...view content, padded to fill...
-----------------------------------------------------------------  status
187 threads                                                        hints
j/k move · enter read · / filter · ctrl+s full-text · ctrl+r reload · q quit
```

Header: `att` in bold pink 205, then a dim 241 separator and title.

Status line: dim 241, changes per mode. List shows `187 threads`, `61 of 187
threads` while filtering, `N results for "query"` after a search, `searching...`
while waiting. Detail shows `2cf53ec4 · 149 messages · line 133 of 2482`.

Hints line: dim 241, lists every key for the current mode. It is the only
discoverability mechanism.

## Thread list

One line per thread, no wrapping:

```
> 27c12788  Explore skull2-f4m3 change scope  ·  opencode  ·  337 msgs  ·  12s
  2cf53ec4  OpenSpec and Beans project bootstrap  ·  opencode  ·  148 msgs  ·  29s
```

- Marker: `> ` when selected, two spaces otherwise.
- ID: last 8 characters of the atuin id, dim.
- Title: the only content-bearing field, truncated with `...` at the right edge.
- Then harness, message count, and a coarse age (`12s`, `29m`, `2h`, `1d`).
- Selected row: entire line bold pink 212, including the metadata. There is no
  reverse video or background.

The list shows as many rows as fit, with `offset` scrolling. There is no
scrollbar, no grouping by day, no column alignment across rows. Separators are
`  ·  ` in dim gray.

Search results reuse the same list; the status line says how many results and
the hint line stays the same.

## Conversation view

```
Build a testable version                                        title
opencode · deepseek-v4.1-flash · 4b160473 · 90 msgs · 38m ago   meta
/tmp/opensessions                                               cwd

[system] 12:12
[user] 12:12
  can you build this so I can test it
[assistant] 12:13
  Reasoning · 97 tokens
  I'll check the repo state to see what needs building.
  tool shell
    {"command":"rtk git status && rtk git log --oneline -5"}
  result
    [rtk] /!\ No hook installed ...
    * main...origin/main
    ... (8 more lines)
```

- Three header lines (title pink bold, meta dim, cwd dim), then a blank line.
- Each message starts with a role label in brackets, colored: user 45 cyan,
  assistant 42 green, system 241 gray, and some roles (attachment, ai-title)
  fall through to plain text. A `HH:MM` timestamp follows in dim.
- Content is indented two spaces; wrapped at the terminal width.
- Reasoning renders as `Reasoning · N tokens` in dim, no content.
- Tool calls render as `tool <name>` with the raw JSON input below, indented
  four spaces. Results render as `result`, with ` (error)` appended in red 196.
- Tool input and output are capped at 8 lines, then `... (N more lines)`.
- The whole conversation is flattened into a slice of lines and scrolled by
  line; `pgup/pgdn` move a page, `g/G` jump to the ends.
- No scrollbar. No message index or jump-to-message. No way to collapse or
  expand a block.

## Interaction map

| Context | Keys                                                     |
| ------- | -------------------------------------------------------- |
| list    | `j/k` or arrows move, `g/G` first/last, `enter` open      |
| list    | `/` filter, `ctrl+s` search the filter text, `ctrl+r` reload |
| filter  | type to narrow, `enter` apply, `ctrl+u` clear, `esc` clear |
| detail  | `j/k`, `pgup/pgdn`, `g/G`, `esc` or `h` back, `q` quit     |

`/` filters live on the substring across title, harness, model, cwd, preview,
last reply, git branch and IDs. `ctrl+s` sends the text currently in the filter
as an atuin full-text search, limited to 100 results. `esc` in the list clears
the search and reloads everything.

## Where the code lives

```
cmd/att/main.go           entrypoint, --version, --help
internal/atuin/types.go   Thread, Conversation, Message, Content
internal/atuin/client.go  runs atuin threads list|show|search --style json
internal/ui/app.go        model, update loop, modes, navigation, loading
internal/ui/view.go       all rendering, styles, row layout, wrapping
internal/ui/filter.go     filter and search state
```

`internal/ui/view.go` is the file a design change touches first. Styles are
seven lipgloss values at the top of that file.

## Problems worth solving

1. The list spends a lot of width on metadata that barely varies. Harness is
   one of two values, message count is usually irrelevant, age is coarse. Title
   is the thing being scanned. What deserves the space?
2. Selection is bold pink over the whole row. Does that read as selection at a
   glance, and does it survive a light terminal theme?
3. There is no sense of recency beyond the age column. Grouping by day, a
   relative-time gutter, or dimming older rows are all open options.
4. Long titles are truncated at the far right, after the metadata, so the row
   reads as one long sentence. Column alignment is accidental.
5. In the conversation, tool calls and results are at the same indent as
   assistant prose, separated only by the words `tool` and `result`. A long
   transcript becomes a wall.
6. The tool input is raw JSON on one long wrapped line. There is no command
   formatting, no diffs, no syntax color.
7. Reasoning occupies a line but hides its content, which is either useful
   context or noise depending on the reader.
8. A three-line header plus two footer lines on every screen is expensive in
   short terminals. Could the header compress, or move into the status band?
9. `line 133 of 2482` is the only position feedback in a 2500-line transcript.
10. Filter, search and plain list are three states with the same visual
    language. Should they be distinguished more clearly?

## Out of scope for this session

Milestone 03 covers packaging, docs and configuration. Implementation,
tests and the coverage gate come after the design is agreed. Do not redesign
atuin's data model or propose features that need a server.

## Kickoff prompt for the design session

Copy the block below into a fresh session, with the three PNGs attached.

```
I am designing two terminal UI views for att, a Go/bubbletea/lipgloss TUI that
browses atuin agent threads. The attached PNGs are real screenshots: the thread
list, the list with a filter applied, and the conversation view. The briefing is
in design/BRIEFING.md in the repo, read it first.

Propose a design direction for the thread list and the conversation view. Stay
inside terminal constraints: 80x24 minimum, 256-color ANSI, keyboard only, no
images. Show annotated sketches, then a prioritized list of concrete changes I
can implement in lipgloss. Do not write code.
```
