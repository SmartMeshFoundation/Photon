package blockchain

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

/*
Events handles all contract events from blockchain
*/
type Events struct {
	client             *helper.SafeEthClient
	lock               sync.RWMutex
	LogChannelMap      map[string]chan types.Log
	RegistryAddress    common.Address //this address is unique
	Subscribes         map[string]ethereum.Subscription
	StateChangeChannel chan transfer.StateChange
	stopped            bool // has stopped?
}

//NewBlockChainEvents create BlockChainEvents
func NewBlockChainEvents(client *helper.SafeEthClient, registryAddress common.Address) *Events {
	be := &Events{client: client,
		LogChannelMap:   make(map[string]chan types.Log),
		Subscribes:      make(map[string]ethereum.Subscription),
		RegistryAddress: registryAddress,
	}
	for _, name := range eventNames {
		be.LogChannelMap[name] = make(chan types.Log, 10)
	}
	return be
}

var eventNames = []string{params.NameTokenAdded,
	params.NameChannelNew,
	params.NameChannelNewBalance,
	params.NameChannelClosed,
	params.NameChannelSettled,
	params.NameChannelSecretRevealed,
}
var eventAbiMap = map[string]string{
	params.NameChannelNew:            contracts.ChannelManagerContractABI,
	params.NameTokenAdded:            contracts.RegistryABI,
	params.NameChannelNewBalance:     contracts.NettingChannelContractABI,
	params.NameChannelClosed:         contracts.NettingChannelContractABI,
	params.NameChannelSettled:        contracts.NettingChannelContractABI,
	params.NameChannelSecretRevealed: contracts.NettingChannelContractABI,
}

func (be *Events) installEventListener() (err error) {
	var sub ethereum.Subscription
	defer func() {
		//event listener create error,must exit
		if err != nil {
			err = be.uninstallEventListener()
			if err != nil {
				log.Error(fmt.Sprintf("uninstallEventListener err %s", err))
			}
		} else {
			//if ethclient reconnect
			c := be.client.RegisterReConnectNotify("Events")
			go func() {
				_, ok := <-c
				if ok {
					//eventlistener need reinstall
					err = be.installEventListener()
					if err != nil {
						log.Error(fmt.Sprintf("installEventListener err %s", err))
					}
				}
			}()
		}
	}()
	for _, name := range eventNames {
		contractAddr := utils.EmptyAddress
		if name == params.NameTokenAdded { //only registry's contract address is only one
			contractAddr = be.RegistryAddress
		}
		sub, err = rpc.EventSubscribe(contractAddr, name, eventAbiMap[name], be.client, be.LogChannelMap[name])
		if err != nil {
			return
		}
		//ChannelNew
		be.Subscribes[name] = sub
	}
	//try to listen event rightnow
	be.startListenEvent()
	return err
}
func (be *Events) uninstallEventListener() (err error) {
	for _, sub := range be.Subscribes {
		sub.Unsubscribe()
	}
	return nil
}
func (be *Events) startListenEvent() {
	for _, name := range eventNames {
		go func(name string) {
			ch := be.LogChannelMap[name]
			sub := be.Subscribes[name]
			for {
				select {
				case l, ok := <-ch:
					if !ok {
						//channel closed
						return
					}
					debugPrintLog(&l)
					switch name {
					case params.NameTokenAdded:
						ev, err := newEventTokenAdded(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventTokenAdded err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractReceiveTokenAddedStateChange{
							RegistryAddress: ev.ContractAddress,
							TokenAddress:    ev.TokenAddress,
							ManagerAddress:  ev.ChannelManagerAddress,
						})
					case params.NameChannelNew:
						ev, err := newEventEventChannelNew(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventEventChannelNew err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractReceiveNewChannelStateChange{
							ManagerAddress: ev.ContractAddress,
							ChannelAddress: ev.NettingChannelAddress,
							Participant1:   ev.Participant1,
							Participant2:   ev.Participant2,
							SettleTimeout:  ev.SettleTimeout,
						})
					case params.NameChannelNewBalance:
						ev, err := newEventChannelNewBalance(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelNewBalance err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractReceiveBalanceStateChange{
							ChannelAddress:     ev.ContractAddress,
							TokenAddress:       ev.TokenAddress,
							ParticipantAddress: ev.ParticipantAddress,
							Balance:            ev.Balance,
							BlockNumber:        ev.BlockNumber,
						})
					case params.NameChannelClosed:
						ev, err := newEventChannelClosed(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelClosed err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractReceiveClosedStateChange{
							ChannelAddress: ev.ContractAddress,
							ClosingAddress: ev.ClosingAddress,
							ClosedBlock:    ev.BlockNumber,
						})
					case params.NameChannelSettled:
						ev, err := newEventChannelSettled(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelSettled err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractReceiveSettledStateChange{
							ChannelAddress: ev.ContractAddress,
							SettledBlock:   ev.BlockNumber,
						})
					case params.NameChannelSecretRevealed:
						ev, err := newEventChannelSecretRevealed(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelSecretRevealed err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractReceiveWithdrawStateChange{
							ChannelAddress: ev.ContractAddress,
							Secret:         ev.Secret,
							Receiver:       ev.ReceiverAddress,
						})
					case params.NameTransferUpdated:
						ev, err := newEventTransferUpdated(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventTransferUpdated err=%s", err))
							debugPrintLog(&l)
							continue
						}
						be.sendStateChange(&mediatedtransfer.ContractTransferUpdatedStateChange{
							RegistryAddress: ev.RegistryAddress,
							ChannelAddress:  ev.ContractAddress,
							Participant:     ev.NodeAddress,
						})
					default:
						log.Crit(fmt.Sprintf("receive unkown event %s,it must be a bug", name))
					}
					//event to statechange
				case err := <-sub.Err():
					log.Error(fmt.Sprintf("eventlistener %s error:%v", name, err))
					//close(ch)
					return
				}
			}
		}(name)
	}
}

//Stop event listenging
func (be *Events) Stop() {
	log.Info("Events stop...")
	be.stopped = true
	close(be.StateChangeChannel)
	for _, sub := range be.Subscribes {
		sub.Unsubscribe()
	}
	log.Info("Events stop ok...")
}
func (be *Events) sendStateChange(st transfer.StateChange) {
	if be.stopped {
		return
	}
	be.StateChangeChannel <- st
}

//GetAllRegistryEvents query all new token events
func (be *Events) GetAllRegistryEvents(registryAddress common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	FromBlockNUmber := ethrpc.BlockNumber(fromBlock)
	if FromBlockNUmber < 0 {
		FromBlockNUmber = 0
	}
	ToBlockNumber := ethrpc.BlockNumber(toBlock)
	if toBlock < 0 {
		ToBlockNumber = ethrpc.LatestBlockNumber
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), registryAddress, FromBlockNUmber, ToBlockNumber,
		params.NameTokenAdded, eventAbiMap[params.NameTokenAdded], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventTokenAdded(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	return
}

/*
GetAllChannelManagerEvents get all new channel events
*/
func (be *Events) GetAllChannelManagerEvents(mgrAddress common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	FromBlockNUmber := ethrpc.BlockNumber(fromBlock)
	if FromBlockNUmber < 0 {
		FromBlockNUmber = 0
	}
	ToBlockNumber := ethrpc.BlockNumber(toBlock)
	if toBlock < 0 {
		ToBlockNumber = ethrpc.LatestBlockNumber
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), mgrAddress, FromBlockNUmber, ToBlockNumber,
		params.NameChannelNew, eventAbiMap[params.NameChannelNew], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventEventChannelNew(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	return

}

//GetAllNettingChannelEvents get channel deposit,close,settle,withdraw,transferupdate events
func (be *Events) GetAllNettingChannelEvents(chAddr common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	FromBlockNUmber := ethrpc.BlockNumber(fromBlock)
	if FromBlockNUmber < 0 {
		FromBlockNUmber = 0
	}
	ToBlockNumber := ethrpc.BlockNumber(toBlock)
	if toBlock < 0 {
		ToBlockNumber = ethrpc.LatestBlockNumber
	}
	/*
				params.NameChannelNewBalance,
				params.NameChannelClosed,
				params.NameChannelSettled,
				params.NameChannelSecretRevealed,
		Use the font to determine the event name and combine the four queries into one, which is to get the Event Signature
	*/
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelNewBalance, eventAbiMap[params.NameChannelNewBalance], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err2 := newEventChannelNewBalance(&l)
		if err2 != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelClosed, eventAbiMap[params.NameChannelClosed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err2 := newEventChannelClosed(&l)
		if err2 != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelSettled, eventAbiMap[params.NameChannelSettled], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err2 := newEventChannelSettled(&l)
		if err2 != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelSecretRevealed, eventAbiMap[params.NameChannelSecretRevealed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err2 := newEventChannelSecretRevealed(&l)
		if err2 != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameTransferUpdated, eventAbiMap[params.NameChannelSecretRevealed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventTransferUpdated(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	return
}

/*
events ChannelClosed and ChannelSecretRevealed must be sent to channels's Channel
*/
func (be *Events) getAllNettingChannelCloseAndWithdrawEvent(fromBlock int64) (stateChanges []transfer.StateChange) {
	FromBlockNUmber := ethrpc.BlockNumber(fromBlock)
	if FromBlockNUmber < 0 {
		FromBlockNUmber = 0
	}

	ToBlockNumber := ethrpc.LatestBlockNumber

	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), utils.EmptyAddress, FromBlockNUmber, ToBlockNumber,
		params.NameChannelClosed, eventAbiMap[params.NameChannelClosed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err2 := newEventChannelClosed(&l)
		if err2 != nil {
			continue
		}
		stateChanges = append(stateChanges, &mediatedtransfer.ContractReceiveClosedStateChange{
			ChannelAddress: e.ContractAddress,
			ClosingAddress: e.ClosingAddress,
			ClosedBlock:    e.BlockNumber,
		})
	}

	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), utils.EmptyAddress, FromBlockNUmber, ToBlockNumber,
		params.NameChannelSecretRevealed, eventAbiMap[params.NameChannelSecretRevealed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelSecretRevealed(&l)
		if err != nil {
			continue
		}
		stateChanges = append(stateChanges, &mediatedtransfer.ContractReceiveWithdrawStateChange{
			ChannelAddress: e.ContractAddress,
			Secret:         e.Secret,
			Receiver:       e.ReceiverAddress,
		})
	}
	return
}

/*
Start listening events send to  channel can duplicate but cannot lose.
1. first resend events may lost (duplicat is ok)
2. listen new events on blockchain
*/
func (be *Events) Start(LastBlockNumber int64) error {
	stateChanges := be.getAllNettingChannelCloseAndWithdrawEvent(LastBlockNumber)
	be.StateChangeChannel = make(chan transfer.StateChange, len(stateChanges)+20)
	for _, st := range stateChanges {
		be.sendStateChange(st)
	}
	return be.installEventListener()
}
