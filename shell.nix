{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  
  packages = with pkgs; [
    zsh
    git
    openssh
  ];

  shellHook = ''
    exec zsh
  '';

  ZDOTDIR = builtins.toString ./.;
}
