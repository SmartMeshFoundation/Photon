#!/bin/sh
export GIT_COMMIT=$(git rev-list -1 HEAD)
export GO_VERSION=$(go version|sed 's/ //g')
export BUILD_DATE=$(date|sed 's/ //g'|sed 's/://g')
export VERSION=1.0.0
echo $GIT_COMMIT
echo $GO_VERSION
echo $BUILD_DATE
echo $VERSION
gomobile bind -ldflags " -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GitCommit=$GIT_COMMIT -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GoVersion=$GO_VERSION -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.BuildDate=$BUILD_DATE -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.Version=$VERSION "    -target=android

cp mobile.aar /Users/bai/dev/develop-photon/app/libs
