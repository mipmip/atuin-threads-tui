# Changelog

All notable changes to this project are documented in this file. The format
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- The conversation view is a two column layout: the message body on the left at
  a fixed reading width, and the time, author, kind, and token count in a meta
  column on the right.
- Tool calls show their command, file, or edit instead of raw JSON, and tool
  output keeps its indentation.
- Runs of empty attachment messages collapse into one counted row, and thinking
  placeholders give way to a token count.

### Added

- Patch blocks and session summaries now appear in the conversation view.
