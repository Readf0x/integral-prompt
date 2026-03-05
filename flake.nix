{
  description = "Integral Prompt for zsh";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, nixpkgs, ... }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      inherit (nixpkgs) lib;
    in {
      devShells = {
        default = pkgs.mkShell {
          packages = with pkgs; [
            bash
            delve
            gcc
            go
            zsh
          ];
        };
      };
      packages = rec {
        integral = pkgs.buildGoModule (final: rec {
          pname = "integral";
          version = "v0.4.0";

          src = ./.;

          nativeBuildInputs = [ pkgs.makeWrapper ];

          vendorHash = "sha256-/gzW1vihul19oMf016fhk32JpuTv3ssSoulo5M05I5E=";

          subPackages = [
            "cmd/integral"
          ];

          preBuild = ''
            go run ./cmd/integral/gen.go "${self.sourceInfo.shortRev or self.sourceInfo.dirtyShortRev}" "${final.version}"
          '';

          postInstall = ''
            mkdir -p $out/share
            cp -r share/integral $out/share/integral

            wrapProgram $out/bin/${pname} \
              --prefix XDG_DATA_DIRS : $out/share
          '';

          meta = {
            description = "Cross shell prompt theme written in Golang";
            homepage = "https://github.com/Readf0x/integral-prompt";
            license = lib.licenses.gpl3Plus;
            mainProgram = pname;
          };
        });
        default = integral;
        test = {
          zsh = pkgs.writeShellScriptBin "debugEnv" ''
            export XDG_DATA_DIRS="${builtins.toString ./.}/share:$XDG_DATA_DIRS"
            export ZDOTDIR="${builtins.toString ./.}"
            export PATH="${lib.makeBinPath [default]}:$PATH"
            zsh
          '';
          bash = let
            rcfile = pkgs.writeTextFile {
              name = "bashrc";
              text = ''
                [[ $- == *i* ]] &&
                  source -- "${pkgs.blesh}/share/blesh/ble.sh" --attach=none
                ${builtins.readFile ./.bashrc}
                [[ ! $\{BLE_VERSION-} ]] || ble-attach
              '';
            };
          in pkgs.writeShellScriptBin "debugEnv" ''
            export XDG_DATA_DIRS="${builtins.toString ./.}/share:$XDG_DATA_DIRS"
            export PATH="${lib.makeBinPath [default]}:$PATH"
            bash --rcfile ${rcfile}
          '';
        };
      };
    }) // {
      homeManagerModules = rec {
        integral = import ./hm.nix self;
        default = integral;
      };
    };
}
