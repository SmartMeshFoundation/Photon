package wakeuphandler

import "github.com/ethereum/go-ethereum/common"

/*
IWakeUpHandler 供Transport层使用,提供消息重发线程的挂起/唤醒服务
*/
type IWakeUpHandler interface {
	/*
		RegisterWakeUpChan 注册一个唤醒通道,Transport会在目标地址上线之后通过该通道通知上层
	*/
	RegisterWakeUpChan(addr common.Address, c chan int)
	/*
		UnRegisterWakeUpChan 取消一个唤醒通道
	*/
	UnRegisterWakeUpChan(addr common.Address)

	/*
		WakeUp Transport层在目标地址上线之后的唤醒方法
	*/
	WakeUp(addr common.Address)
}
