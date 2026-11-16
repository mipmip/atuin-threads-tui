# cli-entrypoint Specification

## Purpose
TBD - created by archiving change att-bootstrap. Update Purpose after archive.

## Requirements

### Requirement: Binary identity

The project SHALL build a single executable named `att`. Its version SHALL come
from the repository `VERSION` file, both when built with `go build` and when
built with `nix build`.

#### Scenario: Version flag

- **WHEN** `att --version` is run
- **THEN** it prints `att <version>` to stdout, where `<version>` is the content
  of `VERSION`, and exits with status 0

#### Scenario: Help flag

- **WHEN** `att --help` is run
- **THEN** it prints usage naming `att` and exits with status 0

#### Scenario: Version follows the file

- **WHEN** `VERSION` is changed and the binary is rebuilt
- **THEN** `att --version` prints the new value without editing Go source

### Requirement: Startup failure is reported

The command SHALL report a readable error and exit non-zero when the terminal
UI cannot start.

#### Scenario: Unknown flag

- **WHEN** `att --nope` is run
- **THEN** it writes an error mentioning the unknown flag to stderr and exits
  with status 2

### Requirement: Testable entrypoint

The command SHALL keep argument handling in a function that takes its
arguments and output streams, so the behavior can be tested without a
terminal.

#### Scenario: Run with an explicit stream

- **WHEN** `run([]string{"--version"}, buffer, buffer)` is called
- **THEN** the buffer receives the version line and no process state is
  touched
