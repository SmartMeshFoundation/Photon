package blockchain

import (
	"time"

	"sync"

	"fmt"

	"context"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/ethereum/go-ethereum/core/types"
)

//stop this call back when return non nil error
type AlarmCallback func(blockNumber int64) error

//Task to notify when a block is mined.
type AlarmTask struct {
	client          *helper.SafeEthClient
	LastBlockNumber int64
	shouldStop      chan struct{}
	waitTime        time.Duration
	callback        []AlarmCallback
	lock            sync.Mutex
}

func NewAlarmTask(client *helper.SafeEthClient) *AlarmTask {
	t := &AlarmTask{
		client:          client,
		waitTime:        time.Second,
		LastBlockNumber: -1,
		shouldStop:      make(chan struct{}), //sync channel
	}
	return t
}

/*
Register a new callback.

        Note:
            The callback will be executed in the AlarmTask context and for
            this reason it should not block, otherwise we can miss block
            changes.
*/
func (this *AlarmTask) RegisterCallback(callback AlarmCallback) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.callback = append(this.callback, callback)
}

//Remove callback from the list of callbacks if it exists
func (this *AlarmTask) RemoveCallback(cb AlarmCallback) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, c := range this.callback {
		addr1 := &c
		addr2 := &cb
		if addr1 == addr2 {
			this.callback = append(this.callback[:k], this.callback[k+1:]...)
		}
	}

}

func (this *AlarmTask) run() {
	log.Debug(fmt.Sprintf("starting block number blocknubmer=%d", this.LastBlockNumber))
	for {
		err := this.waitNewBlock()
		if err != nil {
			time.Sleep(this.waitTime)
		}
	}
}

func (this *AlarmTask) waitNewBlock() error {
	currentBlock := this.LastBlockNumber
	headerCh := make(chan *types.Header, 1)
	//get the lastest number imediatelly
	h, err := this.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	headerCh <- h
	sub, err := this.client.SubscribeNewHead(context.Background(), headerCh)
	if err != nil {
		//reconnect?
		log.Warn("SubscribeNewHead block number err:", err)
		return err
	}
	for {
		select {
		case h, ok := <-headerCh:
			if !ok {
				//client broke?
				return errors.New("SubscribeNewHead channel closed unexpected")
			}
			if currentBlock != -1 && h.Number.Int64() != currentBlock+1 {
				log.Warn(fmt.Sprintf("alarm missed %d blocks", h.Number.Int64()-currentBlock))
			}
			currentBlock = h.Number.Int64()
			if currentBlock%10 == 0 {
				log.Trace(fmt.Sprintf("new block :%d", currentBlock))
			}
			var removes []AlarmCallback
			for _, cb := range this.callback {
				err := cb(currentBlock)
				if err != nil {
					removes = append(removes, cb)
				}
			}
			for _, cb := range removes {
				this.RemoveCallback(cb)
			}
		case <-this.shouldStop:
			sub.Unsubscribe()
			//close(headerCh) //should close by ethclient
			return nil
		case err = <-sub.Err():
			//reconnect here, todo fix ,how to distinguish which error should reconnect
			log.Error(fmt.Sprintf("err=%s", err))
			//spew.Dump(err)
			//if eof try to reconnect
			if err != nil {
				this.client.RecoverDisconnect()
				return errors.New("broken connection")
			}
		}

	}
	return nil
}

func (this *AlarmTask) Start() {
	h, err := this.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(fmt.Sprintf("HeaderByNumber error %s", err))
	}
	this.LastBlockNumber = h.Number.Int64()
	go this.run()
}
func (this *AlarmTask) Stop() {
	log.Info("alarm task stop...")
	this.shouldStop <- struct{}{}
	close(this.shouldStop)
	log.Info("alarm task stop ok...")
}
