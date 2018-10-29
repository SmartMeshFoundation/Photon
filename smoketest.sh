#!/bin/bash
# The script does automatic smoketest on photon api
# Start pwd : $HOME/gopath/src/github.com/SmartMeshFoundation/Photon

# 1. install photon
cd $GOPATH/src/github.com/SmartMeshFoundation/Photon/cmd/photon
./build.sh
cp photon $GOPATH/bin

# 2. build envinit and run
cd ../tools/smoketest/envinit
go build
./envinit --eth-rpc-endpoint=$ETHRPCENDPOINT
if [ $? -ne 0 ]; then
    exit -1
fi

# 3. build smoketest and run
cd ..
rm -rf .photon
if [ ! -d "./log" ];then
    mkdir log
fi
go build
./smoketest
if [ $? -ne 0 ]; then
    echo "failed"
    echo "log/N0.log ----->"
    tail -200 log/N0.log
    echo""
    echo "log/N1.log ----->"
    tail -200 log/N1.log
    echo""
    echo "log/N2.log ----->"
    tail -200 log/N2.log
    echo""
    echo "log/N3.log ----->"
    tail -200 log/N3.log
    echo""
    echo "log/N4.log ----->"
    tail -200 log/N4.log
    echo""
    echo "log/N5.log ----->"
    tail -200 log/N5.log
    echo""
    echo "log/smoketest.log ----->"
    tail -200 log/smoketest.log
    exit -1
else
    echo "succeed"
fi