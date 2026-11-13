# atuin-threads-tui (ATT)

A terminal user interface for browsing the threads from your atuin shell
history. Atuin records every command and groups the history of a session into
threads; This tool gives those threads a place to be read, filtered, and
searched without retyping commands or scrolling through a flat list. Its a
compagnion to Atuin Threads and its purpose it to in depth study the
conversation and also other information stored in threads.

The TUI is written in Go with bubbletea and reads the local atuin SQLite
database read-only. Proposals live in `openspec/`, the issue tracker lives in
`.beans/`, and the build, tests, and coverage gate are defined by `flake.nix`.

## Commands

```bash
att                                           # run the TUI
nix develop                                   # dev shell with go, gopls, sqlite
nix flake check                               # build + vet + tests + coverage gate
nix build                                     # build the binary
go test ./...                                 # tests only, inside the dev shell
openspec list                                 # active changes
openspec validate                             # validate changes and specs
beans list                                    # issues
./scripts/ship-change.sh <change> [subject]   # gate, archive, commit, push
```

Dependencies are vendored into `vendor/` once any are added, so that
`nix flake check` can build offline. Run `go mod vendor` after `go get`.

Coverage gate: `>=70%` overall and `>=80%` for packages matching `/(internal|core)/`.
Override with `MIN_TOTAL` and `MIN_CORE` when running `scripts/gate.sh` directly.

## Beans

When I refer to issues like atuin-threads-tui-rn3b checkout the task
in @.beans/atuin-threads-tui-rn3b-*.md

In this project we will use these tasks as epics for making openspec proposals.

WHEN you create a proposal at a link to this task in the proposal.md.
WHEN a bean is used to create an proposal change the status to "in-progress"
WHEN a proposal is archived add the link to the archived proposal in the frontmatter of this task like this:

```
openspec-link: openspec/changes/archive/....
```

You are allowed to update these statuses in the task frontmatter:

- in-progress
- todo
- draft
- completed
- scrapped

When making changes you are allowed to update the date/time in `updated_at` in the task frontmatter

Besides updating status and openspec-link, you are NOT ALLOWED to modify the contents of the task file.
