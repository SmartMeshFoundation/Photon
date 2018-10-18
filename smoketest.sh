#!/bin/bash
# The script does automatic smoketest on smartraiden api
# Start pwd : $HOME/gopath/src/github.com/SmartMeshFoundation/SmartRaiden

# 1. install smartraiden
cd $GOPATH/src/github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden
./build.sh
cp smartraiden $GOPATH/bin

# 2. build envinit and run
cd ../tools/smoketest/envinit
go build
./envinit --eth-rpc-endpoint=$ETHRPCENDPOINT
if [ $? -ne 0 ]; then
    exit -1
fi

# 3. build smoketest and run
cd ..
rm -rf .smartraiden
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