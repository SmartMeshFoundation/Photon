#!/bin/bash
# The script runs automatic test case for Photon in ./cases
# Start pwd : $GOPATH/src/github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager

# for server
export GOPATH=/home/gotest/goproj

echo HOME = $HOME
echo GOPATH = $GOPATH
export PATH=$PATH:/home/gotest/go/bin:/home/gotest/goproj/bin
cd $GOPATH/src/github.com/SmartMeshFoundation/Photon

# get the lasted code
git pull
if [ $? -ne 0 ]; then
    echo "git pull failed"
    exit -1
fi

# build Photon
cd cmd/photon
chmod +x ./build.sh
./build.sh
if [ $? -ne 0 ]; then
    echo "Photon build failed"
    exit -1
fi
cp photon $GOPATH/bin/

# run deploygethmatrix.sh
cd ../tools/deploygeth
chmod +x ./deploygethmatrix.sh
./deploygethmatrix.sh

# build casemaneger
cd ../casemanager
go build
rm log/*

# run casemanager
./casemanager --case=all --auto --matrix --slow
if [ $? -ne 0 ]; then
    echo "casemanager run failed"
    tar -czvf /home/gotest/tmp/casemanager-matrix.log.tar.zip /home/gotest/casemanage-matrix.log /home/gotest/goproj/src/github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/log
    nodemailer wuhan_53@163.com,baizhenxuan@qq.com smartraiden@163.com 'casemanager失败,附件为全部日志,请尽快排查问题.场景重现请在服务器193.112.248.133上/home/gotest/goproj/src/github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager目录下执行./casemanager --case=报错case名' -j 'Casenamager-matrix场景测试不通过,请尽快排查问题' --attachment '/home/gotest/tmp/casemanager-matrix.log.tar.zip' -u smartraiden@163.com -p pass77 -s smtp.163.com
    exit -1
fi