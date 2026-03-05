#!/usr/bin/env sh

sed 's/"v0.*";/"'"$1"'";/' < flake.nix > flake.nix
git add flake.nix
git commit -m "bump version to $1"
git tag -a "$1" -m "$2"
