# SuperNode
 SuperNode is an off-chain scaling solution for MetaLife.
## Project Status
  This project is still very much a work in progress. It can be used for testing, but it should not be used for real funds. We are doing our best to identify and fix problems, and implement missing features. Any help testing the implementation, reporting bugs, or helping with outstanding issues is very welcome.

## Build and run
```
  go get github.com/MetaLife-Protocol/SuperNode/
  cd $GOPATH/github.com/MetaLife-Protocol/SuperNode
  make 
  ./build/bin/supernode
  
  ./supernode --datadir=./supernode-data \
    --api-address=0.0.0.0:12001 \
    --listen-address=127.0.0.1:12003 \
    --address="0xf413B3187ed510b5b083AB6c5d3BCC259CeF96e9" \
    --keystore-path ./keystore \
    --password-file ***  \
    --eth-rpc-endpoint ws://transport01.smartmesh.cn:33333   \
    --debug   \
    --verbosity 5  \
    --registry-contract-address 0x242e0de2B118279D1479545A131a90A8f67A2512 \
    --pub-address="0xb05Feb81fB4BF6d8B2eB5A5Ae883BAA9E7530cB7" \
    --reward-period 10 \
    --pub-apihost=127.0.0.1:10008 \
    --pfs http://transport01.smartmesh.cn:7000 \
    --xmpp 
```

## mobile support
SuperNode can works on Android and iOS using mobile's API.  it needs [go mobile](https://github.com/golang/mobile) to build mobile library.
### build Android mobile library
```bash
cd mobile
./build_Android.sh 
```
then you can integrate `mobile.aar` into your project.
### build iOS mobile framework
```bash
./build_iOS.sh
```
then you can integrate `Mobile.framework` into your project.
## Requirements
Latest version of SMC

We need go's plugin module.
