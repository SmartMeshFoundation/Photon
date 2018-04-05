package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//本地注释：布置新场景，注册新Token，注册新雷电网络，启动雷电测试节点
	NewTokenName := NewScene()
	//本地注释：测试交易
	Transfer(NewTokenName, "./../../testdata/TransCase/case1.ini")
	//TransferTestCase(NewTokenName)
	//本地注释：Api逐项测试
	//ApiTest()
}
