#!/bin/sh
source ../env.sh 
echo $GIT_COMMIT
echo $GO_VERSION
echo $BUILD_DATE
echo $VERSION
gomobile bind -v -ldflags " -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GitCommit=$GIT_COMMIT -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GoVersion=$GO_VERSION -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.BuildDate=$BUILD_DATE -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.Version=$VERSION "    -target=android/arm64,android/arm  

zip -r android_$VERSION.zip mobile.aar mobile-sources.jar
cp mobile.aar /Volumes/dev/develop-photon/smartmesh/app/libs
rm -f mobile.aar mobile-sources.jar