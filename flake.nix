{
  description = "atuin-threads-tui: a terminal UI for browsing atuin shell history threads";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f:
        nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.buildGoModule {
          pname = "att";
          version =
            pkgs.lib.removeSuffix "\n" (pkgs.lib.fileContents ./cmd/att/VERSION);
          src = pkgs.lib.cleanSource ./.;
          subPackages = [ "cmd/att" ];
          vendorHash = null;
          meta.mainProgram = "att";
        };
      });

      apps = forAllSystems (pkgs: {
        default = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.default}/bin/att";
          meta.description = "Browse atuin threads in the terminal";
        };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
            sqlite
            jq
          ];
          shellHook = ''
            echo "atuin-threads-tui dev shell: go $(go version | awk '{print $3}')"
          '';
        };
      });

      checks = forAllSystems (pkgs: {
        gate = pkgs.runCommand "atuin-threads-tui-gate" {
          nativeBuildInputs = with pkgs; [ go bash coreutils gawk gnused ];
          src = pkgs.lib.cleanSource ./.;
        } ''
          cp -r $src src
          chmod -R u+w src
          cd src
          export HOME=$TMPDIR
          export CGO_ENABLED=0
          export GOCACHE=$TMPDIR/go-cache
          export GOPATH=$TMPDIR/go
          export GOFLAGS=-mod=vendor
          export GOPROXY=off
          bash scripts/gate.sh
          touch $out
        '';
      });

      formatter = forAllSystems (pkgs: pkgs.nixpkgs-fmt);
    };
}
