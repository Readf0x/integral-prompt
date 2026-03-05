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
    in rec {
      devShells = {
        test = {
          zsh = pkgs.mkShell {
            packages = with pkgs; [
              zsh
              packages.default
            ];

            shellHook = ''
              exec zsh
            '';

            XDG_DATA_DIRS = "${builtins.toString ./.}/share:$XDG_DATA_DIRS";
            ZDOTDIR = builtins.toString ./.;
          };
          bash = pkgs.mkShell {
            packages = with pkgs; [
              bash
              packages.default
            ];

            shellHook = ''
              source .bashrc
            '';

            XDG_DATA_DIRS = "${builtins.toString ./.}/share:$XDG_DATA_DIRS";
          };
        };
        default = pkgs.mkShell {
          packages = with pkgs; [
            bash
            delve
            gcc
            git
            go
            openssh
            zsh
          ];
        };
      };
      packages = rec {
        integral = pkgs.buildGoModule rec {
          pname = "integral";
          version = "v0.3.6";

          src = ./.;

          nativeBuildInputs = [ pkgs.makeWrapper ];

          vendorHash = "sha256-/gzW1vihul19oMf016fhk32JpuTv3ssSoulo5M05I5E=";

          ldflags = [ "-X 'main.VersionString=%s, %s'" ];

          subPackages = [
            "cmd/integral"
          ];

          preBuild = ''
            go run ./cmd/integral/gen.go "Nix build" "${version}"
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
        };
        default = integral;
      };
    }) // {
      homeManagerModules = rec {
        integral = import ./hm.nix self;
        default = integral;
      };
    };
}
