#!/bin/bash
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. gas        
# 2. golint   
# 3. go vet     
# 4. vetshadow  

#
# gometalinter (github.com/alecthomas/gometalinter) is used to run each static
# checker.
go get -v -u github.com/alecthomas/gometalinter
set -ex

# Make sure gometalinter is installed and $GOPATH/bin is in your path.
# $ go get -v github.com/alecthomas/gometalinter"
# $ gometalinter --install"
if [ ! -x "$(type -p gometalinter)" ]; then
  exit 1
fi
gometalinter --install


# Automatic checks
test -z "$(gometalinter -j 4 --disable-all \
--enable=golint \
--enable=vet \
--enable=gosec \
--enable=vetshadow \
--deadline=10m  \
--vendor \
--skip network/rpc/contracts \
--skip internal/debug \
--skip log \
--skip cmd/tools/test  ./... 2>&1 | grep -v 'ALL_CAPS\|OP_' 2>&1 | tee /dev/stderr)"

