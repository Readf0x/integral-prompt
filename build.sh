#!/usr/bin/env zsh

mkdir -p build
rm -rf build/*
go generate
case "$1" in
  deb)
    TAG="$(git describe --tags --abbrev=0)"
    mkdir build/DEBIAN
    cat > build/DEBIAN/control <<EOF
Package: integral-prompt
Version: ${TAG#v}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Jean <https://github.com/readf0x>
Description: Math themed shell prompt
EOF
    mkdir -p build/usr/local/bin
    cp integral build/usr/local/bin
    chmod 755 build/usr/local/bin/*
    cp -r share build/usr
    dpkg-deb --build build
  ;;
  *)
    mkdir -p build/usr/bin
    CGO_ENABLED=0 go build -ldflags="-extldflags=-static"
    cp integral build/usr/bin
    cp -r share build/usr
    tar -cJf build.tar.xz -C build .
  ;;
esac

