# Tasks: thread-list-ui

## 1. Model

- [x] 1.1 Add `Model`, `New`, modes, and typed messages
- [x] 1.2 Implement `Init` and commands for listing and opening threads

## 2. Interaction

- [x] 2.1 Implement list navigation with bounds
- [x] 2.2 Implement open and back between list and detail
- [x] 2.3 Implement detail scrolling clamped to rendered lines
- [x] 2.4 Implement quit keys and a retry key for list errors

## 3. Rendering

- [x] 3.1 Render title bar, thread rows, status line, and key hints
- [x] 3.2 Render loading, empty, and error states
- [x] 3.3 Render conversation metadata and content blocks with wrapping
- [x] 3.4 Handle window resize without panics

## 4. Wiring

- [x] 4.1 Run the atuin client from `cmd/att` with bubbletea
- [x] 4.2 Vendor bubbletea and lipgloss
- [x] 4.3 Tests for navigation, open, scroll, states, and rendering helpers
