#!/bin/sh
export GIT_COMMIT=$(git rev-list -1 HEAD)
export GO_VERSION=$(go version|sed 's/ //g')
export BUILD_DATE=$(date|sed 's/ //g'|sed 's/://g')
echo $GIT_COMMIT
echo $GO_VERSION
echo $BUILD_DATE
gomobile bind -ldflags " -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GitCommit=$GIT_COMMIT -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GoVersion=$GO_VERSION -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.BuildDate=$BUILD_DATE "  -target=android
gomobile bind -ldflags " -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GitCommit=$GIT_COMMIT -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GoVersion=$GO_VERSION -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.BuildDate=$BUILD_DATE "  -target=ios
#gomobile bind -v -tags=lldb -gcflags "-N -l"  -target=ios
#tar cjf - Mobile.framework/ | split -b 5m - ./out/ios.tar.bz2.
#tar cjf - mobile.aar mobile-sources.jar | split -b 5m - ./out/android.tar.bz2
#cat android.tar.bz2a* | tar xj  
