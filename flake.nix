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
            fish
            gcc
            go
            zsh
          ];
        };
      };
      packages = rec {
        integral = pkgs.buildGoModule (final: rec {
          pname = "integral";
          version = "v0.5.3";

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
            cp -r share $out

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
          debugEnv = pkgs.stdenv.mkDerivation (finalAttrs: {
            name = "debugEnv";

            src = ./.;

            dontPatch = true;
            dontConfigure = true;
            dontBuild = true;

            installPhase = ''
              mkdir -p $out/.config
              cat << EOF > $out/.bashrc
              [[ \$- == *i* ]] &&
                source -- "${pkgs.blesh}/share/blesh/ble.sh" --attach=none
              $(cat .bashrc)
              [[ ! \''${BLE_VERSION-} ]] || ble-attach
              EOF
              cp .zshrc $out
              cp -r .fishrc/fish $out/.config/fish
            '';
          });
          zsh = pkgs.writeShellScriptBin "debug" ''
            export XDG_DATA_DIRS="${default}/share:$XDG_DATA_DIRS"
            export ZDOTDIR="${test.debugEnv}"
            export PATH="${lib.makeBinPath [default]}:$PATH"
            zsh
          '';
          bash = pkgs.writeShellScriptBin "debug" ''
            export XDG_DATA_DIRS="${default}/share:$XDG_DATA_DIRS"
            export PATH="${lib.makeBinPath [default]}:$PATH"
            bash --rcfile "${test.debugEnv}/.bashrc"
          '';
          fish = pkgs.writeShellScriptBin "debug" ''
            export XDG_DATA_DIRS="${default}/share:$XDG_DATA_DIRS"
            export PATH="${lib.makeBinPath [default]}:$PATH"
            fish
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
