package blockchain

import (
	"time"

	"sync"

	"fmt"

	"context"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/ethereum/go-ethereum/core/types"
)

//AlarmCallback stop this call back when return non nil error
type AlarmCallback func(blockNumber int64) error

//AlarmTask notify when a block is mined.
type AlarmTask struct {
	client              *helper.SafeEthClient
	LastBlockNumber     int64
	quitChan            chan struct{}
	stopped             bool
	waitTime            time.Duration
	LastBlockNumberChan chan int64
	lock                sync.Mutex
}

//NewAlarmTask create a alarm task
func NewAlarmTask(client *helper.SafeEthClient) *AlarmTask {
	t := &AlarmTask{
		client:              client,
		waitTime:            time.Second,
		LastBlockNumber:     -1,
		quitChan:            make(chan struct{}), //sync channel
		LastBlockNumberChan: make(chan int64, 10),
	}
	return t
}

func (at *AlarmTask) run() {
	defer rpanic.PanicRecover("alarm task")
	err := at.waitNewBlock()
	if err != nil {
		log.Error("alarm task stopped with err %s", err)
	}
}

func (at *AlarmTask) waitNewBlock() error {
	log.Debug(fmt.Sprintf("start getting lasted block number from blocknubmer=%d", at.LastBlockNumber))
	currentBlock := at.LastBlockNumber
	headerCh := make(chan *types.Header, 1)
	//get the lastest number imediatelly
	h, err := at.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	headerCh <- h
	sub, err := at.client.SubscribeNewHead(context.Background(), headerCh)
	if err != nil {
		//reconnect?
		log.Warn(fmt.Sprintf("SubscribeNewHead block number err: %s", err))
		return err
	}
	for {
		select {
		case h, ok := <-headerCh:
			if at.stopped {
				log.Info(fmt.Sprintf("alarm task quit complete"))
				return nil
			}
			if !ok {
				//client broke?
				return errors.New("SubscribeNewHead channel closed unexpected")
			}
			if currentBlock != -1 && h.Number.Int64() != currentBlock+1 {
				log.Warn(fmt.Sprintf("alarm missed %d blocks", h.Number.Int64()-currentBlock))
			}
			currentBlock = h.Number.Int64()
			at.LastBlockNumber = currentBlock
			if currentBlock%10 == 0 {
				log.Trace(fmt.Sprintf("new block :%d", currentBlock))
			}
			at.LastBlockNumberChan <- currentBlock
		case <-at.quitChan:
			sub.Unsubscribe()
			return nil
		case err = <-sub.Err():
			//reconnect here, todo fix ,how to distinguish which error should reconnect
			log.Error(fmt.Sprintf("err=%s", err))
			//spew.Dump(err)
			//if eof try to reconnect
			if err != nil && !at.stopped {
				at.client.RecoverDisconnect()
				return errors.New("broken connection")
			}
		}
	}
}

//Start this task
func (at *AlarmTask) Start() error {
	h, err := at.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("HeaderByNumber error %s", err)
	}
	at.LastBlockNumber = h.Number.Int64()
	go at.run()
	return nil
}

//Stop this task
func (at *AlarmTask) Stop() {
	log.Info("alarm task stop...")
	at.stopped = true
	close(at.quitChan)
	log.Info("alarm task stop ok...")
}
