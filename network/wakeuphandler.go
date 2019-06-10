package network

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ethereum/go-ethereum/common"
)

type iWakeUpHandler interface {
	registerWakeUpChan(addr common.Address, c chan int)
	unRegisterWakeUpChan(addr common.Address)
	wakeUp(addr common.Address)
}

/*
为3种不同实现的transport提供统一的挂起/唤醒服务
*/
type wakeupHandler struct {
	name                  string // for log
	wakeUpChanListMap     map[common.Address][]chan int
	wakeUpChanListMapLock sync.Mutex
}

func newWakeupHandler(name string) *wakeupHandler {
	return &wakeupHandler{
		name:              name,
		wakeUpChanListMap: make(map[common.Address][]chan int),
	}
}

// registerWakeUpChan 注册唤醒通道,在用户上线时使用
func (h *wakeupHandler) registerWakeUpChan(addr common.Address, c chan int) {
	if c == nil {
		panic("wrong call")
	}
	h.wakeUpChanListMapLock.Lock()
	h.wakeUpChanListMap[addr] = append(h.wakeUpChanListMap[addr], c)
	h.wakeUpChanListMapLock.Unlock()
}

// unRegisterWakeUpChan 移除唤醒通道
func (h *wakeupHandler) unRegisterWakeUpChan(addr common.Address) {
	h.wakeUpChanListMapLock.Lock()
	if _, ok := h.wakeUpChanListMap[addr]; ok {
		delete(h.wakeUpChanListMap, addr)
	}
	h.wakeUpChanListMapLock.Unlock()
}

// wakeUp 当节点上线时调用
func (h *wakeupHandler) wakeUp(addr common.Address) {
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

/*
混合transport
*/
type mixWakeUpHandler struct {
	h1 iWakeUpHandler
	h2 iWakeUpHandler
}

func newMixWakeUpHandler(h1, h2 iWakeUpHandler) *mixWakeUpHandler {
	if h1 == nil || h2 == nil {
		panic("wrong call")
	}
	return &mixWakeUpHandler{
		h1: h1,
		h2: h2,
	}
}

// registerWakeUpChan 注册唤醒通道
func (mh *mixWakeUpHandler) registerWakeUpChan(addr common.Address, c chan int) {
	if c == nil {
		c = make(chan int, 2)
	}
	/*
		一个通道在两个transport中共用,其中任何一个transport探测到对方上线即可
	*/
	mh.h1.registerWakeUpChan(addr, c)
	mh.h2.registerWakeUpChan(addr, c)
}

// unRegisterWakeUpChan 移除唤醒通道
func (mh *mixWakeUpHandler) unRegisterWakeUpChan(addr common.Address) {
	mh.h1.unRegisterWakeUpChan(addr)
	mh.h2.unRegisterWakeUpChan(addr)
}

// wakeUp impl Transporter
func (mh *mixWakeUpHandler) wakeUp(addr common.Address) {
	panic("wrong call")
}
