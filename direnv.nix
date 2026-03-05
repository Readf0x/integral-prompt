{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  
  packages = with pkgs; [
    coreutils
    delve
    gawk
    git
    go
    openssh
    zsh
  ];

}
