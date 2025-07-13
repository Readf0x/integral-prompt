{
  description = "Integral Prompt for zsh";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, nixpkgs, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      inherit (nixpkgs) lib;
    in {
      devShells = {
        default = pkgs.mkShell {
          packages = with pkgs; [
            coreutils
            gawk
            git
            openssh
            zsh
          ];

          shellHook = ''
            exec zsh
          '';

          ZDOTDIR = builtins.toString ./.;
        };
        direnv = import ./direnv.nix { inherit pkgs; };
      };
      packages = rec {
        integral = pkgs.buildGoModule (finalAttrs: rec {
          name = "integral";
          pname = name;
          version = "v0.2.2";

          src = ./.;

          meta = {
            description = "Cross shell prompt theme written in Golang";
            homepage = "https://github.com/Readf0x/integral-prompt";
            license = lib.licenses.gpl3;
          };
        });
        default = integral;
      };
    });
}
