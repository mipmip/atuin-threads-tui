# cli-entrypoint

## MODIFIED Requirements

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
