#!/bin/sh
# 0. install all tools


# 1. create a private ethereum
cd  cmd/tools/deploygeth
./deploygeth.sh
cd -
# wait for geth start complete
sleep 1

export TOKEN_NETWORK=0x50839B01D28390048616C8f28dD1A21CF3CacbfF
export KEY1=2ddd679cb0f0754d0e20ef8206ea2210af3b51f159f55cfffbd8550f58daf779
export KEY2=36234555bc087435cf52371f9a0139cb98a4267ba62b722e3f46b90d35f31678
export ISTEST=1
export ETHRPCENDPOINT="http://127.0.0.1:30307"
export KEYSTORE=$GOPATH/src/github.com/SmartMeshFoundation/Photon/testdata/keystore

# 2. deploy contract for test
cd cmd/tools/newtestenv
go install
cd -  
newtestenv --keystore-path ./cmd/tools/deploygeth/privnet/keystore/ --eth-rpc-endpoint $ETHRPCENDPOINT --base 18 --tokennum 2
if [ $? -ne 0 ]; then
   echo "newtestenv failed"
   exit 1
fi
# 3. go test
go test -v  -failfast -timeout 10m -short `go list ./... | grep -v contracts |grep -v casemanager |grep -v gethworkdir`
if [ $? -ne 0 ]; then
    echo "go test failed"
    exit 1
fi

# 4. build photon
cd cmd/photon 
./build.sh
cp photon $GOPATH/bin 
cd -
# 4. smoke test
# chmod +x smoketest.sh
# ./smoketest.sh
# if [ $? -ne 0 ]; then
#     echo "smoketest failed"
#     exit 1
# fi

# 5. casemanager
cd cmd/tools/casemanager
mkdir log
go build
#指定部署私链的rpc
./casemanager --case=all --auto --eth-rpc-endpoint  $ETHRPCENDPOINT
if [ $? -ne 0 ]; then
    echo "casemanager run failed"
    exit 1
fi
cd -

# 6. kill geth
#ps -ef | grep geth  |grep 7888| grep -v grep | awk '{print $2}' |xargs kill -9
