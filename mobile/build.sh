#!/bin/sh
export VERSION=1.0.1
./build_android.sh
zip -r android_$VERSION.zip mobile.aar mobile-sources.jar
rm -f mobile.aar mobile-sources.jar
./build_iOS.sh
zip -r iOS_$VERSION.zip Mobile.framework
rm -rf Mobile.framework
