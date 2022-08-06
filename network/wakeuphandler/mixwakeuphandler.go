package wakeuphandler

import "github.com/ethereum/go-ethereum/common"

/*
MixWakeUpHandler 供混合Transport层使用的IWakeUpHandler实现,目前photon大部分实现使用的都是混合的Transport层
*/
type MixWakeUpHandler struct {
	h1 IWakeUpHandler
	h2 IWakeUpHandler
}

/*
NewMixWakeUpHandler init
*/
func NewMixWakeUpHandler(h1, h2 IWakeUpHandler) *MixWakeUpHandler {
	if h1 == nil || h2 == nil {
		panic("wrong call")
	}
	return &MixWakeUpHandler{
		h1: h1,
		h2: h2,
	}
}

// RegisterWakeUpChan impl IWakeUpHandler
func (mh *MixWakeUpHandler) RegisterWakeUpChan(addr common.Address, c chan int) {
	if c == nil {
		c = make(chan int, 2)
	}
	/*
		一个通道在两个transport中共用,其中任何一个transport探测到对方上线即可
	*/
	mh.h1.RegisterWakeUpChan(addr, c)
	mh.h2.RegisterWakeUpChan(addr, c)
}

// UnRegisterWakeUpChan impl IWakeUpHandler
func (mh *MixWakeUpHandler) UnRegisterWakeUpChan(addr common.Address) {
	mh.h1.UnRegisterWakeUpChan(addr)
	mh.h2.UnRegisterWakeUpChan(addr)
}

// WakeUp impl IWakeUpHandler, shouldn't call
func (mh *MixWakeUpHandler) WakeUp(addr common.Address) {
	panic("wrong call")
}
