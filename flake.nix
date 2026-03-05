{
  description = "Integral Prompt for zsh";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }@inputs: 
  {
    devShells."x86_64-linux".default = import ./shell.nix {
      pkgs = nixpkgs.legacyPackages."x86_64-linux";
    };
  };
}
