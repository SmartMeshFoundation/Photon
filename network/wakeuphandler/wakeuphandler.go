package wakeuphandler

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
WakeUpHandler 为3种不同实现的transport提供统一的挂起/唤醒服务
*/
type WakeUpHandler struct {
	name                  string // for log
	wakeUpChanListMap     map[common.Address][]chan int
	wakeUpChanListMapLock sync.Mutex
}

/*
NewWakeupHandler init
*/
func NewWakeupHandler(name string) *WakeUpHandler {
	return &WakeUpHandler{
		name:              name,
		wakeUpChanListMap: make(map[common.Address][]chan int),
	}
}

// RegisterWakeUpChan impl IWakeUpHandler
func (h *WakeUpHandler) RegisterWakeUpChan(addr common.Address, c chan int) {
	if c == nil {
		panic("wrong call")
	}
	h.wakeUpChanListMapLock.Lock()
	h.wakeUpChanListMap[addr] = append(h.wakeUpChanListMap[addr], c)
	h.wakeUpChanListMapLock.Unlock()
}

// UnRegisterWakeUpChan impl IWakeUpHandler
func (h *WakeUpHandler) UnRegisterWakeUpChan(addr common.Address) {
	h.wakeUpChanListMapLock.Lock()
	if _, ok := h.wakeUpChanListMap[addr]; ok {
		delete(h.wakeUpChanListMap, addr)
	}
	h.wakeUpChanListMapLock.Unlock()
}

// WakeUp impl IWakeUpHandler
func (h *WakeUpHandler) WakeUp(addr common.Address) {
	// 节点上线通知所有已经挂起的通道
	h.wakeUpChanListMapLock.Lock()
	log.Trace(fmt.Sprintf("%s back to online on %s and wakeup %d chan", addr.String(), h.name, len(h.wakeUpChanListMap[addr])))
	if cs, ok := h.wakeUpChanListMap[addr]; ok && len(cs) > 0 {
		for _, c := range cs {
			c <- 1
		}
	}
	h.wakeUpChanListMapLock.Unlock()
}
