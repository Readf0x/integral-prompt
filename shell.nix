{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  
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
}
