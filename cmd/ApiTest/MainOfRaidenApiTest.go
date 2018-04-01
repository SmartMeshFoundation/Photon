package main

import ()

func main() {
	//本地注释：布置新场景，注册新Token，注册新雷电网络，启动雷电测试节点
	NewTokenName := NewScene()
	//本地注释：强化测试正常交易
	TransferTest(NewTokenName)
	//本地注释：Api逐项测试,Api使用范例
	ApiTest()
}
