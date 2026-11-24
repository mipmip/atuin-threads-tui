---
# atuin-threads-tui-yqbs
title: Scope filter for the current directory and host
status: completed
openspec-link: openspec/changes/archive/2026-11-24-thread-scope-filter
type: feature
priority: normal
created_at: 2026-11-24T21:45:55Z
updated_at: 2026-11-24T22:00:01Z
parent: atuin-threads-tui-s9te
---

Press d in the list to keep only threads captured in the working directory and h to keep only threads from this host. Both compose with / and ctrl+s, and the status line names what is active. The local host id is read once from atuin info.

## Summary of Changes

The list gained two scope toggles. `.` keeps only threads captured in the working directory, `@` keeps only threads from this machine, both compose with each other and with `/` and atuin search, and the status line names what is active. The local host id is read once at startup from `atuin info`; when atuin reports no host id, `@` shows `host unknown` instead of an empty list. No column was added.
