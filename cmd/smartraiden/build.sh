#!/bin/sh
export GIT_COMMIT=$(git rev-list -1 HEAD)
export GO_VERSION=$(go version|sed 's/ //g')
export BUILD_DATE=$(date|sed 's/ //g'|sed 's/://g')
echo $GIT_COMMIT
echo $GO_VERSION
echo $BUILD_DATE
go build -ldflags " -X github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl.GitCommit=$GIT_COMMIT -X github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl.GoVersion=$GO_VERSION -X github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl.BuildDate=$BUILD_DATE " 
