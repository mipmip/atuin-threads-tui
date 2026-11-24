# thread-browser

## ADDED Requirements

### Requirement: Limit the list to one directory or host

The list SHALL offer a directory scope toggled with `.` and a host scope toggled
with `@`. The directory scope SHALL keep only threads whose working directory
equals the process working directory. The host scope SHALL keep only threads
whose host equals the local host id read from `atuin info`. Both scopes SHALL
compose with each other, with the text filter, and with search results.

#### Scenario: Only this directory

- **WHEN** `.` is pressed in a project directory
- **THEN** the list shows only threads whose working directory equals that
  directory and the status line names the directory

#### Scenario: Only this host

- **WHEN** `@` is pressed and `atuin info` reported the local host id
- **THEN** the list shows only threads captured on that host and the status line
  says `host: local`

#### Scenario: Both scopes at once

- **WHEN** `.` and `@` are both on
- **THEN** only threads that match both are shown

#### Scenario: Scope clears

- **WHEN** the active scope key is pressed again
- **THEN** that scope is dropped and the list widens again

#### Scenario: Scope combines with the text filter

- **WHEN** a directory scope is active and `/` is given a query
- **THEN** the rows must match both the scope and the query, and the status
  reports the narrowed count

#### Scenario: Host id is unknown

- **WHEN** `atuin info` could not be read
- **THEN** `@` shows `host: unknown` and leaves the rows unchanged instead of
  emptying the list

#### Scenario: Empty result

- **WHEN** the scope matches no thread
- **THEN** the list says that no thread matches, and the scope key clears it
