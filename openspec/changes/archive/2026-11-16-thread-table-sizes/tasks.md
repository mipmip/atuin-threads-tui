# Tasks: thread table sizes

## 1. Layout

- [x] 1.1 Measure the agent code, message count, title, and directory widths of
      the threads in view
- [x] 1.2 Clamp those widths to their ranges and derive the table width from
      the columns instead of the terminal width
- [x] 1.3 Order the columns as id, date, agent, message count, title, directory
- [x] 1.4 Shrink the title, then the directory, then drop columns until the
      table fits

## 2. Agent code

- [x] 2.1 Keep the `oc`, `cc`, `co`, and `--` aliases
- [x] 2.2 Show any other harness in lower case, cut to five characters
- [x] 2.3 Grow the header from `ag` to `agent` with the column

## 3. Tests

- [x] 3.1 Cover the column offsets at 120 columns
- [x] 3.2 Cover a title that follows the content and one that hits the cap
- [x] 3.3 Cover the agent column growing for a longer harness name
- [x] 3.4 Cover the drop order at widths from 120 down to 37

## 4. Verify

- [x] 4.1 `bash scripts/gate.sh`
- [x] 4.2 `nix flake check`
- [x] 4.3 Render the list at 120 and 80 columns and compare it to the mockups
