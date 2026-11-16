# Tasks: thread-table-agent-codes

## 1. Codes

- [x] 1.1 Map `claude-code`, `opencode`, `codex`, and `pi` to `claude`,
      `ocode`, `codex`, and `pi` in `agentTag`.
- [x] 1.2 Raise `agentMaxWidth` to 6 and cut unknown harness names to that
      width.

## 2. Tests

- [x] 2.1 Cover every mapped name, the empty harness, and an unknown long name.
- [x] 2.2 Update the column offsets and the agent column width tests for the
      6-cell maximum.

## 3. Verify

- [x] 3.1 Run `bash scripts/gate.sh` and `nix flake check`.
- [x] 3.2 Render the list at 120x24 and confirm the codes read as intended and
      the columns still line up.
