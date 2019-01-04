#!/bin/sh
export CHANNEL=0x021a641147131481de1be929cd5f2c0c42cc25f40db8ee06d21924e614b49b68
export TOKEN_NETWORK=0xdCE3c72Ab939a49AE393e7cBF08267c1cA2daa42
export KEY1=2ddd679cb0f0754d0e20ef8206ea2210af3b51f159f55cfffbd8550f58daf779
export KEY2=36234555bc087435cf52371f9a0139cb98a4267ba62b722e3f46b90d35f31678
export ISTEST=1
export ETHRPCENDPOINT="ws://127.0.0.1:8546"
export KEYSTORE=$GOPATH/src/github.com/SmartMeshFoundation/Photon/testdata/keystore
go test -v  -failfast -timeout 10m -short `go list ./... | grep -v contracts |grep -v casemanager`