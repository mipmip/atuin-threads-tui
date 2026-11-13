# Tasks: thread-search-filter

## 1. Filter mode

- [x] 1.1 Add `modeFilter`, `filter`, and `filterInput`
- [x] 1.2 Implement `/` to enter, `enter` to apply, `esc` to clear
- [x] 1.3 Implement typing, `backspace`, and `ctrl+u`
- [x] 1.4 Match against title, harness, model, cwd, preview, last reply, branch, and id

## 2. Derived list

- [x] 2.1 Derive visible rows from the filter
- [x] 2.2 Clamp selection and offset to the visible rows
- [x] 2.3 Show `N of M threads` and the filter prompt

## 3. Full-text search

- [x] 3.1 Add `searchMsg` and a command calling `Source.Search`
- [x] 3.2 Run it on `ctrl+s` when the filter is non-empty
- [x] 3.3 Show a searching state and `N results for "query"`
- [x] 3.4 Clear with `esc` and reload the list

## 4. Tests

- [x] 4.1 Filtering, field matching, no match, backspace, and clear
- [x] 4.2 Navigation bounds inside a filter and opening a filtered thread
- [x] 4.3 Search success, empty query, failure, and clearing
