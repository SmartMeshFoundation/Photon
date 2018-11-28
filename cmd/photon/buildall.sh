#!/bin/sh
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
export GOARM=5
go env
./build.sh
mv photon photon_linux_arm_$VERSION
zip -r photon_linux_arm_$VERSION.zip photon_linux_arm_$VERSION
rm -f photon_linux_arm_$VERSION
export GOOS=linux
export GOARCH=amd64
./build.sh
mv photon photon_linux_amd64_$VERSION
zip -r photon_linux_amd64_$VERSION.zip photon_linux_amd64_$VERSION
rm -f photon_linux_amd64_$VERSION
export GOOS=windows
export GOARCH=amd64
./build.sh
mv photon.exe photon_windows_amd64_$VERSION.exe
zip -r photon_windows_amd64_$VERSION.zip photon_windows_amd64_$VERSION.exe
rm -f photon_windows_amd64_$VERSION.exe
export GOOS=darwin
export GOARCH=amd64
./build.sh
mv photon photon_darwin_amd64_$VERSION
zip -r photon_darwin_amd64_$VERSION.zip photon_darwin_amd64_$VERSION
rm -f photon_darwin_amd64_$VERSION

