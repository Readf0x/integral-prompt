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
    in rec {
      devShells = {
        default = pkgs.mkShell {
          packages = with pkgs; [
            coreutils
            gawk
            git
            openssh
            zsh
            packages.default
          ];

          shellHook = ''
            exec zsh
            export XDG_DATA_DIRS="$PWD/share:$XDG_DATA_DIRS"
          '';

          ZDOTDIR = builtins.toString ./.;
        };
        direnv = import ./direnv.nix { inherit pkgs; };
      };
      packages = rec {
        integral = pkgs.buildGoModule rec {
          name = "integral";
          pname = name;
          version = "v0.2.2";

          src = ./.;

          vendorHash = "sha256-nrQ5rkXBAu2prVlJaqHHLIarQsW+/6oH+LMLqpjWvYc=";

          postInstall = ''
            mkdir -p $out/share
            cp -r share/integral $out/share/integral
          '';

          meta = {
            description = "Cross shell prompt theme written in Golang";
            homepage = "https://github.com/Readf0x/integral-prompt";
            license = lib.licenses.gpl3;
          };
        };
        default = integral;
      };
    }) // {
      homeManagerModules = rec {
        integral = import ./hm.nix;
        default = integral;
      };
    };
}
