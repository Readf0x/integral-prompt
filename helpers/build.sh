#!/usr/bin/env zsh

mkdir -p build
if [ build(NF) ]; then
  rm -rf build/*
fi
go generate ./cmd/integral
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
    go build ./cmd/integral
    strip integral -o build/usr/local/bin/integral
    chmod 755 build/usr/local/bin/*
    cp -r share build/usr
    dpkg-deb --build build
  ;;
  debug)
    go build ./cmd/integral -gcflags="all=-N -l"
  ;;
  *)
    mkdir -p build/usr/bin
    CGO_ENABLED=0 go build -ldflags="-extldflags=-static" ./cmd/integral
    strip integral -o build/usr/bin/integral
    cp -r share build/usr
    tar -cJf build.tar.xz -C build .
  ;;
esac

