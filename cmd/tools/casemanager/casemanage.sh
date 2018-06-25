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
    tar -cvf /home/gotest/tmp/casemanager.log.tar /home/gotest/casemanage.log /home/gotest/goproj/src/github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/log
    nodemailer wuhan_53@163.com smartraiden@163.com 'casemanager失败,附件为全部日志,请尽快排查问题.场景重现请在服务器193.112.248.133上/home/gotest/goproj/src/github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager目录下执行./casemanager --case=报错case名' -j 'Casenamager场景测试不通过,请尽快排查问题' --atachment '/home/gotest/tmp/casemanager.log.tar' -u smartraiden@163.com -p pass77 -s smtp.163.com
    exit -1
fi