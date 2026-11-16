# atuin-threads-tui

A terminal user interface for browsing the threads from your atuin shell
history: read, filter, search, and export the conversations that atuin recorded.

```bash
nix run            # start the TUI
nix build          # ./result/bin/att
```

## Development

```bash
nix develop        # go, gopls, sqlite
bash scripts/gate.sh
nix flake check
```

## Releasing

The version lives in `cmd/att/VERSION`. See [docs/RELEASING.md](docs/RELEASING.md)
for the bump, the dry run, and what a `v*` tag publishes.
