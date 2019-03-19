#!/bin/sh
export VERSION=1.1.0-${GIT_COMMIT:0-40:4}
./build_android.sh
./build_iOS.sh

