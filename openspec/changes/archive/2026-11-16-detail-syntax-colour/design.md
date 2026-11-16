# Design

## Context

The detail view builds every message into cached lines once, then draws a
viewport slice. Styling therefore costs nothing per frame as long as it happens
while building. The palette lives in `internal/ui/view.go` as seven lipgloss
styles using 256-colour indexes.

Today only the author label carries a role colour, so the body of a prompt and
the body of a reply look identical. Diffs from patches and edits are drawn
plain, which is the weakest spot because they are the lines people read
closely. Tool calls print their verb, path, and command in one undifferentiated
default colour.

## Palette

One hue per job. Values are 256-colour indexes; five of the six are already in
the app.

| element | style | index |
| --- | --- | --- |
| prompt prose | cyan | 45 |
| reply prose | green | 42 |
| diff addition | green | 42 |
| diff deletion | red | 196 |
| diff hunk header | cyan | 45 |
| command after `$` | yellow | 222 |
| tool verb, thread title | accent bold | 205 |
| path | default bold | none |
| reasoning, summary, rail, chrome | muted | 241 |
| error result | red | 196 |

Why not background tints for diff lines: foreground colours read the same on
light and dark backgrounds, backgrounds do not, and reverse video would fight
the selection style used elsewhere.

## Where colour is applied

Styling happens after wrapping, never before, so `ansi.StringWidth` and the
column arithmetic keep measuring plain text:

```
text block      wrap(block.Text, prose)  ->  styleAll(lines, roleStyle)
patch diff      preserve(diff, width-2)  ->  styleDiff(lines)
edit lines      prefixedLines(...)       ->  styleDiff(lines)
shell command   "  $ " + part            ->  commandStyle
tool head       verb / path / counts     ->  styled pieces via ansi.Truncate
```

`truncate` and `ansi.Truncate` measure visible cells and are ANSI-aware, so a
styled head can still be cut at the column edge.

## Role rule

```
role == "user"        -> cyan
role == "assistant"   -> green
anything else         -> default
```

Reasoning keeps the muted style it has now, and a summary block keeps muted
too, so the role colour applies to `text` blocks only. That is the whole
change for prose; a message that holds both text and a tool call gets the role
colour on its prose and the structural colours on its call.

## Diff rule

A diff line is classified by its first character, in this order:

```
@@ ...   -> hunk header, cyan
+++ / ---   -> file header, muted (not an addition or a deletion)
+ ...   -> addition, green
- ...   -> deletion, red
anything    -> default
```

The `+++` and `---` guards matter: without them a file header reads as a
change and the top of every patch lies about its content.

The `+12 -8` counts in an edit or patch head get the same two colours. Building
them as separate styled fragments keeps the head truncatable: the plain form
stays available for width, the styled fragments are concatenated after.

## Fonts

The app cannot choose a font. It can ask for bold, italic, and faint:

- bold for the tool verb, the path, and the thread title (already the case for
  the title)
- italic is not added: several terminals ignore it and the difference is lost
- faint is not added: support is patchier than the colour fallback

## Risks

- Cyan and green are light colours and can be weak on a light terminal.
  Accepting that keeps the palette the reader already knows; a future theme
  setting (`atuin-threads-tui-3ab2`) can flip it.
- Forty lines of green result would be noise, which is why results are excluded
  from tinting by design.
- Styling per line during the build costs a `lipgloss.Render` per prose line.
  Measured against a 55,527-line thread the build stays within its current
  budget because the rendering is cached, but the frame tests will confirm the
  build time did not regress noticeably.
- `NO_COLOR` and 16-colour terminals are handled by lipgloss's profile
  detection, and every colour here is redundant with a text signal (`-`, `+`,
  `!`, `$`, the author label), so a monochrome terminal loses nothing but
  decoration.

## Testing

- prompt prose carries the cyan SGR, reply prose the green one, tool text none
- a wrapped prompt carries the style on every line
- thinking and summary stay muted
- diff classification: addition, deletion, hunk, `+++`, `---`, plain line
- `+12 -8` counts in an edit head and in a patch head
- shell head: verb, path, command
- result payload has no SGR at all, error result is red
- existing `plain()` assertions still pass, so width and layout are unchanged
