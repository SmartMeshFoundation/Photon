#!/bin/sh
wd=`pwd`
echo $wd 
mkdir -p  gethworkdir/src/github.com/ethereum/
export GOPATH=$wd/gethworkdir
export PATH=$GOPATH/bin:$PATH
#  下载自定义版本geth
cd $GOPATH/src/github.com/ethereum/
git clone https://github.com/nkbai/go-ethereum.git
cd go-ethereum
go get ./...
git pull
cd cmd/geth
go install

cd $wd 

#有可能上一个节点还在运行着没有退出 kill它
ps -ef | grep geth  |grep 7888| grep -v grep | awk '{print $2}' |xargs kill -9

## 准备搭建私链
geth version
rm -rf privnet/geth 
geth --datadir privnet init baipoatestnetmatrix.json

# 尽量避免不必要的log输出,干扰photon信息
geth --datadir=./privnet --unlock 3de45febbd988b6e417e4ebd2c69e42630fefbf0 --password ./privnet/keystore/pass --port 40404 --networkid 7888 --ws --wsaddr 0.0.0.0 --wsorigins "*" --wsport 30306 --rpc --rpccorsdomain "*" --rpcapi eth,admin,web3,net,debug,personal --rpcport 30307 --rpcaddr 127.0.0.1 --mine  --verbosity 1 --nodiscover &
#newtestenv因为总是使用固定的账户,所以合约地址是固定的
##0x50839B01D28390048616C8f28dD1A21CF3CacbfF
