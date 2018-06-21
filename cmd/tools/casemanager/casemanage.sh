#!/bin/bash
# The script runs automatic test case for smartraiden in ./cases
# Start pwd : $GOPATH/src/github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager

# for server
export GOPATH=/home/gotest/goproj

echo HOME = $HOME
echo GOPATH = $GOPATH
export PATH=$PATH:$GOPATH/bin
cd $GOPATH/src/github.com/SmartMeshFoundation/SmartRaiden

# get the lasted code
git pull
if [ $? -ne 0 ]; then
    echo "git pull failed"
    exit -1
fi

# build smartraiden
cd cmd/smartraiden
go install
if [ $? -ne 0 ]; then
    echo "smartraiden build failed"
    exit -1
fi

# build casemaneger
cd ../tools/casemanager
go build

# run casemanager
./casemanager --case=all
if [ $? -ne 0 ]; then
    echo "casemanager run failed"
    exit -1
fi