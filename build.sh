#!/usr/bin/env bash

rm -rf build/*
mkdir -p build/bin
go generate
go build
cp integral build/bin
cp -r share build
# package step
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
    chmod 755 build/bin/*
    mv build/bin/* build/usr/local/bin/
    rmdir build/bin
    dpkg-deb --build build
  ;;
  *)
    tar -cJf build.tar.xz -C build .
  ;;
esac

