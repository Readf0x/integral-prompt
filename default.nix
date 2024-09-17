{ pkgs ? import <nixpkgs> {} }:
with pkgs;
stdenv.mkDerivation rec {
  pname = "integral-prompt";
  version = "0.2.2";

  src = ./.;

  installPhase = ''
    install -D integral.zsh-theme --target-directory=$out/share/zsh-integral
    install -D lib/* --target-directory=$out/share/zsh-integral/lib
    install -D comp/_integral --target-directory=$out/share/zsh/site-functions
    ln -s $out/share/zsh-integral/integral.zsh-theme $out/share/zsh/site-functions/integral
  '';
}
