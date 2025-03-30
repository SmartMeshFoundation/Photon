#!/bin/sh
# must enable cgo,because of plugin
source ../../env.sh
source ../../VERSION

export CGO_ENABLED=1
echo $GIT_COMMIT
go  build  -ldflags "   -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GitCommit=$GIT_COMMIT -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GoVersion=$GO_VERSION -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.BuildDate=$BUILD_DATE -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.Version=$VERSION "
cp photon $GOPATH/bin
