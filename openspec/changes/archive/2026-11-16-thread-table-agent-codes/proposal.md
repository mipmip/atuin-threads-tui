## Why

The agent column shows `oc`, `cc`, and `co`. Those codes are cryptic until you
know that they stand for opencode, claude-code, and codex, and there is no
legend anywhere in the interface. The column is sized to its content, so a
readable code costs only the width it needs, and only when that agent is in
view.

## What Changes

- Show `claude` for `claude-code`, `ocode` for `opencode`, `codex` for `codex`,
  and `pi` for `pi`.
- Raise the agent column maximum from 5 to 6 cells so `claude` fits.
- Cut an unknown harness name to 6 runes instead of 5.
- Keep the `--` placeholder for a thread without a harness.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `thread-browser`: the agent code requirement and the agent column width in
  the table requirement.

## Impact

- `internal/ui/view.go`: `agentTag` and `agentMaxWidth`.
- `internal/ui/table_test.go`: the agent tag test and the column offsets.
