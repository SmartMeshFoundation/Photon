#!/bin/sh
export GIT_COMMIT=$(git rev-list -1 HEAD)
export GO_VERSION=$(go version|sed 's/ //g')
export BUILD_DATE=`date "+%Y-%m-%d-%H:%M:%S"`

