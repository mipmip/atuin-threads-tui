# detail-message-selection

Select messages in a conversation and copy or export them.

Press `v` in a thread to put a cursor on the message at the top. Move with
`j`/`k`, ten at a time with `J`/`K`, and to the ends with `g`/`G`; selected
messages are marked `*`, the cursor `>`. Press `y` to copy the range as
Markdown, `w` to write it to a file, or `e` to export the whole thread. Copying
uses a native clipboard tool when one exists and falls back to OSC 52, and a
payload above 50 KB is written to a file instead of the clipboard.
