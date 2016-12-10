#!/bin/bash

echo "[$(date)]: Starting cross-compilation..."

VERSION=$(git tag | tail -n1)

echo "[$(date)]: Version: $VERSION"

printf "Mac OS X... "
GOOS=darwin GOARCH=amd64 go build 
cp hstats hstats-$VERSION-darwin-amd64
mv hstats hstats-$VERSION-osx-amd64
printf "Done\n"

printf "Linux x86_64... "
GOOS=linux GOARCH=amd64 go build
mv hstats hstats-$VERSION-linux-amd64
printf "Done\n"

printf "Linux x86... "
GOOS=linux GOARCH=386 go build
mv hstats hstats-$VERSION-linux-i386
printf "Done\n"

printf "Linux ARM... "
GOOS=linux GOARCH=arm go build
mv hstats hstats-$VERSION-linux-arm
printf "Done\n"

printf "Linux ARM64... "
GOOS=linux GOARCH=arm64 go build
mv hstats hstats-$VERSION-linux-arm64
printf "Done\n"

printf "Linux MIPS64... "
GOOS=linux GOARCH=mips64 go build
mv hstats hstats-$VERSION-linux-mips64
printf "Done\n"

printf "Windows 64-bit... "
GOOS=windows GOARCH=amd64 go build
mv hstats.exe hstats-$VERSION-windows-64bit.exe
printf "Done\n"

printf "Windows 32-bit... "
GOOS=windows GOARCH=386 go build
mv hstats.exe hstats-$VERSION-windows-32bit.exe
printf "Done\n"

echo "[$(date)]: Done"
