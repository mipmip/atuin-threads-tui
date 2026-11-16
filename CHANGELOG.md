# Changelog

All notable changes to this project are documented in this file. The format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Press `v` in a conversation to select a range of messages, `y` to copy it as Markdown, `w` to write it to a file, and `e` to export the whole thread.

- A message type panel, opened with `t` in a conversation, hides prompts,
  replies, thinking, tool calls, tool results, patches, summaries, or notes and
  keeps the selection for the rest of the session.
- Conversation prose takes the colour of its author, diffs mark additions green
  and deletions red, and tool calls colour the verb, the path, and the command.
- A one cell scrollbar in the right margin of the conversation view shows how
  long the thread is and where the visible lines sit.

### Fixed
- Copying a selection no longer freezes the interface when the clipboard tool
  leaves a background process behind.

### Changed

- The conversation view renders a thread once and scrolls from that cache, so
  reading a long conversation no longer re-renders thousands of messages on
  every key press. Threads over 500 messages render in the background.
- The conversation view is a two column layout: the message body on the left at
  a fixed reading width, and the time, author, kind, and token count in a meta
  column on the right.
- Tool calls show their command, file, or edit instead of raw JSON, and tool
  output keeps its indentation.
- Runs of empty attachment messages collapse into one counted row, and thinking
  placeholders give way to a token count.

### Fixed

- Several keys that arrive in one input batch now each count, so a fast scroll
  is no longer swallowed.
- Opening a thread with very long tool output is fast again: each line is cut in
  one pass instead of being rescanned and copied once per chunk.


### Added

- Patch blocks and session summaries now appear in the conversation view.
