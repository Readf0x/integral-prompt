#!/usr/bin/env bash

rm -rf build/*
mkdir -p build/bin
go generate
go build
cp integral build/bin
cp -r share build
tar -cJf build.tar.xz -C build .

