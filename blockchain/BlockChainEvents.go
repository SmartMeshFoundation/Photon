package blockchain

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/raiden-network/network/helper"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type BlockChainEvents struct {
	client             *helper.SafeEthClient
	lock               sync.RWMutex
	LogChannelMap      map[string]chan types.Log
	RegistryAddress    common.Address //this address is unique
	Subscribes         map[string]ethereum.Subscription
	StateChangeChannel chan transfer.StateChange
}

func NewBlockChainEvents(client *helper.SafeEthClient, registryAddress common.Address) *BlockChainEvents {
	be := &BlockChainEvents{client: client,
		LogChannelMap:      make(map[string]chan types.Log),
		Subscribes:         make(map[string]ethereum.Subscription),
		StateChangeChannel: make(chan transfer.StateChange, 20),
		RegistryAddress:    registryAddress,
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
	params.NameChannelNew:            rpc.ChannelManagerContractABI,
	params.NameTokenAdded:            rpc.RegistryABI,
	params.NameChannelNewBalance:     rpc.NettingChannelContractABI,
	params.NameChannelClosed:         rpc.NettingChannelContractABI,
	params.NameChannelSettled:        rpc.NettingChannelContractABI,
	params.NameChannelSecretRevealed: rpc.NettingChannelContractABI,
}

func (this *BlockChainEvents) InstallEventListener() (err error) {
	var sub ethereum.Subscription
	defer func() {
		//event listener create error,must exit
		if err != nil {
			this.UninstallEventListener()
		} else {
			//if ethclient reconnect
			c := this.client.RegisterReConnectNotify("BlockChainEvents")
			go func() {
				select {
				case _, ok := <-c:
					if ok {
						//eventlistener need reinstall
						this.InstallEventListener()
					}
				}
			}()
		}
	}()
	for _, name := range eventNames {
		contractAddr := utils.EmptyAddress
		if name == params.NameTokenAdded { //only registry's contract address is only one
			contractAddr = this.RegistryAddress
		}
		sub, err = rpc.EventSubscribe(contractAddr, name, eventAbiMap[name], this.client, this.LogChannelMap[name])
		if err != nil {
			return
		}
		//ChannelNew
		this.Subscribes[name] = sub
	}
	//try to listen event rightnow
	this.startListenEvent()
	return err
}
func (this *BlockChainEvents) UninstallEventListener() (err error) {
	for _, sub := range this.Subscribes {
		sub.Unsubscribe()
	}
	return nil
}
func (this *BlockChainEvents) startListenEvent() {
	for _, name := range eventNames {
		go func(name string) {
			ch := this.LogChannelMap[name]
			sub := this.Subscribes[name]
			for {
				select {
				case l, ok := <-ch:
					if !ok {
						//channel closed
						return
					}
					switch name {
					case params.NameTokenAdded:
						ev, err := NewEventTokenAdded(&l)
						if err != nil {
							continue
						}
						this.sendStateChange(&mediated_transfer.ContractReceiveTokenAddedStateChange{
							RegistryAddress: ev.ContractAddress,
							TokenAddress:    ev.TokenAddress,
							ManagerAddress:  ev.ChannelManagerAddress,
						})
					case params.NameChannelNew:
						ev, err := NewEventEventChannelNew(&l)
						if err != nil {
							continue
						}
						this.sendStateChange(&mediated_transfer.ContractReceiveNewChannelStateChange{
							ManagerAddress: ev.ContractAddress,
							ChannelAddress: ev.NettingChannelAddress,
							Participant1:   ev.Participant1,
							Participant2:   ev.Participant2,
							SettleTimeout:  ev.SettleTimeout,
						})
					case params.NameChannelNewBalance:
						ev, err := NewEventChannelNewBalance(&l)
						if err != nil {
							continue
						}
						this.sendStateChange(&mediated_transfer.ContractReceiveBalanceStateChange{
							ChannelAddress:     ev.ContractAddress,
							TokenAddress:       ev.TokenAddress,
							ParticipantAddress: ev.ParticipantAddress,
							Balance:            ev.Balance,
							BlockNumber:        ev.BlockNumber,
						})
					case params.NameChannelClosed:
						ev, err := NewEventChannelClosed(&l)
						if err != nil {
							continue
						}
						this.sendStateChange(&mediated_transfer.ContractReceiveClosedStateChange{
							ChannelAddress: ev.ContractAddress,
							ClosingAddress: ev.ClosingAddress,
							ClosedBlock:    ev.BlockNumber,
						})
					case params.NameChannelSettled:
						ev, err := NewEventChannelSettled(&l)
						if err != nil {
							continue
						}
						this.sendStateChange(&mediated_transfer.ContractReceiveSettledStateChange{
							ChannelAddress: ev.ContractAddress,
							SettledBlock:   ev.BlockNumber,
						})
					case params.NameChannelSecretRevealed:
						ev, err := NewEventChannelSecretRevealed(&l)
						if err != nil {
							continue
						}
						this.sendStateChange(&mediated_transfer.ContractReceiveWithdrawStateChange{
							ChannelAddress: ev.ContractAddress,
							Secret:         ev.Secret,
							Receiver:       ev.ReceiverAddress,
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
func (this *BlockChainEvents) Stop() {
	close(this.StateChangeChannel)
	//channel close by ethclient
	//for _, ch := range this.LogChannelMap {
	//	close(ch)
	//}
	for _, sub := range this.Subscribes {
		sub.Unsubscribe()
	}
}
func (this *BlockChainEvents) sendStateChange(st transfer.StateChange) {
	this.StateChangeChannel <- st
}
func (this *BlockChainEvents) GetAllRegistryEvents(registryAddress common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	FromBlockNUmber := ethrpc.BlockNumber(fromBlock)
	if FromBlockNUmber < 0 {
		FromBlockNUmber = 0
	}
	ToBlockNumber := ethrpc.BlockNumber(toBlock)
	if toBlock < 0 {
		ToBlockNumber = ethrpc.LatestBlockNumber
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), registryAddress, FromBlockNUmber, ToBlockNumber,
		params.NameTokenAdded, eventAbiMap[params.NameTokenAdded], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventTokenAdded(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	return
}

/*
 These helpers have a better descriptive name and provide the translator for
 the caller.
*/
func (this *BlockChainEvents) GetAllChannelManagerEvents(mgrAddress common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	FromBlockNUmber := ethrpc.BlockNumber(fromBlock)
	if FromBlockNUmber < 0 {
		FromBlockNUmber = 0
	}
	ToBlockNumber := ethrpc.BlockNumber(toBlock)
	if toBlock < 0 {
		ToBlockNumber = ethrpc.LatestBlockNumber
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), mgrAddress, FromBlockNUmber, ToBlockNumber,
		params.NameChannelNew, eventAbiMap[params.NameChannelNew], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventEventChannelNew(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	return

}
func (this *BlockChainEvents) GetAllNettingChannelEvents(chAddr common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
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
		可以通过字形判断事件名字,将四个查询合并成一个,具体就是获取Event Signature
	*/
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelNewBalance, eventAbiMap[params.NameChannelNewBalance], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventChannelNewBalance(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelClosed, eventAbiMap[params.NameChannelClosed], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventChannelClosed(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelSettled, eventAbiMap[params.NameChannelSettled], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventChannelSettled(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameChannelSecretRevealed, eventAbiMap[params.NameChannelSecretRevealed], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventChannelSecretRevealed(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	logs, err = rpc.EventGetInternal(rpc.GetQueryConext(), chAddr, FromBlockNUmber, ToBlockNumber,
		params.NameTransferUpdated, eventAbiMap[params.NameChannelSecretRevealed], this.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := NewEventTransferUpdated(&l)
		if err != nil {
			continue
		}
		events = append(events, e)
	}
	return
}
