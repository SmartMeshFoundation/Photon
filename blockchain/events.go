package blockchain

import (
	"fmt"
	"sync"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
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
	client                *helper.SafeEthClient
	lock                  sync.RWMutex
	LogChannelMap         map[string]chan types.Log
	RegistryAddress       common.Address //this address is unique
	SecretRegistryAddress common.Address //get from db or from blockchain
	Subscribes            map[string]ethereum.Subscription
	StateChangeChannel    chan transfer.StateChange
	//启动过程中先把收到事件暂存在这个通道中,等启动完毕以后在保存到StateChangeChannel,保证事件被顺序处理.
	startupStateChangeChannel chan mediatedtransfer.ContractStateChange
	stopped                   bool // has stopped?
	quitChan                  chan struct{}
	TokenNetworks             map[common.Address]bool
	historyEventsGot          bool
}

//NewBlockChainEvents create BlockChainEvents
func NewBlockChainEvents(client *helper.SafeEthClient, registryAddress, secretRegistryAddress common.Address, token2TokenNetwork map[common.Address]common.Address) *Events {
	be := &Events{
		client:                    client,
		LogChannelMap:             make(map[string]chan types.Log),
		Subscribes:                make(map[string]ethereum.Subscription),
		RegistryAddress:           registryAddress,
		SecretRegistryAddress:     secretRegistryAddress,
		quitChan:                  make(chan struct{}),
		TokenNetworks:             make(map[common.Address]bool),
		startupStateChangeChannel: make(chan mediatedtransfer.ContractStateChange, 100),
		StateChangeChannel:        make(chan transfer.StateChange, 10),
	}
	for _, tn := range token2TokenNetwork {
		be.TokenNetworks[tn] = true
	}
	for name := range eventAbiMap {
		be.LogChannelMap[name] = make(chan types.Log, 10)
	}
	return be
}

var eventAbiMap = map[string]string{
	params.NameTokenNetworkCreated:       contracts.TokenNetworkRegistryABI,
	params.NameChannelOpened:             contracts.TokenNetworkABI,
	params.NameChannelOpenedAndDeposit:   contracts.TokenNetworkABI,
	params.NameChannelNewDeposit:         contracts.TokenNetworkABI,
	params.NameChannelClosed:             contracts.TokenNetworkABI,
	params.NameChannelSettled:            contracts.TokenNetworkABI,
	params.NameChannelCooperativeSettled: contracts.TokenNetworkABI,
	params.NameChannelWithdraw:           contracts.TokenNetworkABI,
	params.NameChannelUnlocked:           contracts.TokenNetworkABI,
	params.NameBalanceProofUpdated:       contracts.TokenNetworkABI,
	params.NameChannelPunished:           contracts.TokenNetworkABI,
	params.NameSecretRevealed:            contracts.SecretRegistryABI,
	//the following event is for 3rd party
	//params.NameNewDeposit:              monitoringcontracts.MonitoringServiceABI,
	//params.NameNewBalanceProofReceived: monitoringcontracts.MonitoringServiceABI,
	//params.NameRewardClaimed:           monitoringcontracts.MonitoringServiceABI,
	//params.NameWithdrawn:               monitoringcontracts.MonitoringServiceABI,
}

/*
现在channel 所有信息主要是通过事件获取,连接断开过程中发生的事件是无法通过监听获取到,
必须 get 才行.
要保证事件不遗漏,可以重复
*/
func (be *Events) installEventListener() (err error) {
	var sub ethereum.Subscription
	defer func() {
		//event listener create error,must exit
		if err != nil {
			err = be.uninstallEventListener()
			if err != nil {
				log.Error(fmt.Sprintf("uninstallEventListener err %s", err))
			}
		}
	}()
	for name := range eventAbiMap {
		contractAddr := utils.EmptyAddress
		if name == params.NameTokenNetworkCreated { //only registry's contract address is only one
			contractAddr = be.RegistryAddress
		} else if name == params.NameSecretRevealed {
			contractAddr = be.SecretRegistryAddress
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

//EventChannelSettled2StateChange to stateChange
func EventChannelSettled2StateChange(ev *contracts.TokenNetworkChannelSettled) *mediatedtransfer.ContractSettledStateChange {
	return &mediatedtransfer.ContractSettledStateChange{
		ChannelIdentifier:   common.Hash(ev.ChannelIdentifier),
		TokenNetworkAddress: ev.Raw.Address,
		SettledBlock:        int64(ev.Raw.BlockNumber),
	}
}

//EventChannelCooperativeSettled2StateChange to stateChange
func EventChannelCooperativeSettled2StateChange(ev *contracts.TokenNetworkChannelCooperativeSettled) *mediatedtransfer.ContractCooperativeSettledStateChange {
	return &mediatedtransfer.ContractCooperativeSettledStateChange{
		ChannelIdentifier:   common.Hash(ev.ChannelIdentifier),
		TokenNetworkAddress: ev.Raw.Address,
		SettledBlock:        int64(ev.Raw.BlockNumber),
	}
}

//EventChannelPunished2StateChange to stateChange
func EventChannelPunished2StateChange(ev *contracts.TokenNetworkChannelPunished) *mediatedtransfer.ContractPunishedStateChange {
	return &mediatedtransfer.ContractPunishedStateChange{
		ChannelIdentifier:   common.Hash(ev.ChannelIdentifier),
		TokenNetworkAddress: ev.Raw.Address,
		Beneficiary:         ev.Beneficiary,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
}

//EventChannelWithdraw2StateChange to stateChange
func EventChannelWithdraw2StateChange(ev *contracts.TokenNetworkChannelWithdraw) *mediatedtransfer.ContractChannelWithdrawStateChange {
	c := &mediatedtransfer.ContractChannelWithdrawStateChange{
		ChannelIdentifier: &contracts.ChannelUniqueID{

			ChannelIdentifier: common.Hash(ev.ChannelIdentifier),
			OpenBlockNumber:   int64(ev.Raw.BlockNumber),
		},
		TokenNetworkAddress: ev.Raw.Address,
		Participant1:        ev.Participant1,
		Participant2:        ev.Participant2,
		Participant1Balance: ev.Participant1Balance,
		Participant2Balance: ev.Participant2Balance,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
	if c.Participant1Balance == nil {
		c.Participant1Balance = new(big.Int)
	}
	if c.Participant2Balance == nil {
		c.Participant2Balance = new(big.Int)
	}
	return c
}

//EventTokenNetworkCreated2StateChange to statechange
func EventTokenNetworkCreated2StateChange(ev *contracts.TokenNetworkRegistryTokenNetworkCreated) *mediatedtransfer.ContractTokenAddedStateChange {
	return &mediatedtransfer.ContractTokenAddedStateChange{
		RegistryAddress:     ev.Raw.Address,
		TokenAddress:        ev.TokenAddress,
		TokenNetworkAddress: ev.TokenNetworkAddress,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
}

//EventChannelOpen2StateChange to statechange
func EventChannelOpen2StateChange(ev *contracts.TokenNetworkChannelOpened) *mediatedtransfer.ContractNewChannelStateChange {
	return &mediatedtransfer.ContractNewChannelStateChange{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: ev.ChannelIdentifier,
			OpenBlockNumber:   int64(ev.Raw.BlockNumber),
		},
		TokenNetworkAddress: ev.Raw.Address,
		Participant1:        ev.Participant1,
		Participant2:        ev.Participant2,
		SettleTimeout:       int(ev.SettleTimeout),
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
}

//EventChannelOpenAndDeposit2StateChange to statechange
func EventChannelOpenAndDeposit2StateChange(ev *contracts.TokenNetworkChannelOpenedAndDeposit) (ch1 *mediatedtransfer.ContractNewChannelStateChange, ch2 *mediatedtransfer.ContractBalanceStateChange) {
	ch1 = &mediatedtransfer.ContractNewChannelStateChange{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: ev.ChannelIdentifier,
			OpenBlockNumber:   int64(ev.Raw.BlockNumber),
		},
		TokenNetworkAddress: ev.Raw.Address,
		Participant1:        ev.Participant1,
		Participant2:        ev.Participant2,
		SettleTimeout:       int(ev.SettleTimeout),
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
	ch2 = &mediatedtransfer.ContractBalanceStateChange{
		ChannelIdentifier:   ev.ChannelIdentifier,
		ParticipantAddress:  ev.Participant1,
		BlockNumber:         int64(ev.Raw.BlockNumber),
		Balance:             ev.Participant1Deposit,
		TokenNetworkAddress: ev.Raw.Address,
	}
	return
}

//EventChannelNewDeposit2StateChange to statechange
func EventChannelNewDeposit2StateChange(ev *contracts.TokenNetworkChannelNewDeposit) *mediatedtransfer.ContractBalanceStateChange {
	return &mediatedtransfer.ContractBalanceStateChange{
		ChannelIdentifier:   ev.ChannelIdentifier,
		TokenNetworkAddress: ev.Raw.Address,
		ParticipantAddress:  ev.Participant,
		BlockNumber:         int64(ev.Raw.BlockNumber),
		Balance:             ev.TotalDeposit,
	}
}

//EventChannelClosed2StateChange to statechange
func EventChannelClosed2StateChange(ev *contracts.TokenNetworkChannelClosed) *mediatedtransfer.ContractClosedStateChange {
	c := &mediatedtransfer.ContractClosedStateChange{
		TokenNetworkAddress: ev.Raw.Address,
		ChannelIdentifier:   ev.ChannelIdentifier,
		ClosingAddress:      ev.ClosingParticipant,
		LocksRoot:           ev.Locksroot,
		ClosedBlock:         int64(ev.Raw.BlockNumber),
		TransferredAmount:   ev.TransferredAmount,
	}
	if ev.TransferredAmount == nil {
		c.TransferredAmount = new(big.Int)
	}
	return c
}

//EventBalanceProofUpdated2StateChange to statechange
func EventBalanceProofUpdated2StateChange(ev *contracts.TokenNetworkBalanceProofUpdated) *mediatedtransfer.ContractBalanceProofUpdatedStateChange {
	c := &mediatedtransfer.ContractBalanceProofUpdatedStateChange{
		TokenNetworkAddress: ev.Raw.Address,
		ChannelIdentifier:   ev.ChannelIdentifier,
		LocksRoot:           ev.Locksroot,
		TransferAmount:      ev.TransferredAmount,
		Participant:         ev.Participant,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
	if c.TransferAmount == nil {
		c.TransferAmount = new(big.Int)
	}
	return c
}

//EventChannelUnlocked2StateChange to statechange
func EventChannelUnlocked2StateChange(ev *contracts.TokenNetworkChannelUnlocked) *mediatedtransfer.ContractUnlockStateChange {
	c := &mediatedtransfer.ContractUnlockStateChange{
		TokenNetworkAddress: ev.Raw.Address,
		ChannelIdentifier:   ev.ChannelIdentifier,
		BlockNumber:         int64(ev.Raw.BlockNumber),
		TransferAmount:      ev.TransferredAmount,
		Participant:         ev.PayerParticipant,
	}
	if c.TransferAmount == nil {
		c.TransferAmount = new(big.Int)
	}
	return c
}

//EventSecretRevealed2StateChange to statechange
func EventSecretRevealed2StateChange(ev *contracts.SecretRegistrySecretRevealed) *mediatedtransfer.ContractSecretRevealOnChainStateChange {
	return &mediatedtransfer.ContractSecretRevealOnChainStateChange{
		LockSecretHash: ev.Secrethash,
		BlockNumber:    int64(ev.Raw.BlockNumber),
	}
}
func (be *Events) startListenEvent() {
	for name := range eventAbiMap {
		go func(name string) {
			ch := be.LogChannelMap[name]
			sub := be.Subscribes[name]
			defer rpanic.PanicRecover(fmt.Sprintf("startListenEvent %s", name))
			for {
				select {
				case l, ok := <-ch:
					if !ok {
						//channel closed
						return
					}
					switch name {
					case params.NameTokenNetworkCreated:
						ev, err := newEventTokenNetworkCreated(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventTokenNetworkCreated err=%s", err))
							continue
						}
						be.TokenNetworks[ev.TokenNetworkAddress] = true
						be.sendStateChange(EventTokenNetworkCreated2StateChange(ev))
					case params.NameChannelOpened:
						ev, err := newEventChannelOpen(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelOpen err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive event ChannelOpened, but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelOpen2StateChange(ev))
					case params.NameChannelOpenedAndDeposit:
						ev, err := newEventChannelOpenAndDeposit(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelOpen err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive event ChannelOpened, but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						nev, dev := EventChannelOpenAndDeposit2StateChange(ev)
						be.sendStateChange(nev)
						be.sendStateChange(dev)
					case params.NameChannelNewDeposit:
						ev, err := newEventChannelNewDeposit(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelNewDeposit err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive event channel new deposit ,but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelNewDeposit2StateChange(ev))
					case params.NameChannelUnlocked:
						ev, err := newEventChannelUnlocked(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelUnlocked err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("recevie event channel unlocked ,but it's not our contract,ev=\n%s",
								utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelUnlocked2StateChange(ev))
					case params.NameChannelClosed:
						ev, err := newEventChannelClosed(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelClosed err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive NameChannelClosed ,but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelClosed2StateChange(ev))
					case params.NameChannelSettled:
						ev, err := newEventChannelSettled(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelSettled err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive NameChannelSettled,but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelSettled2StateChange(ev))
					case params.NameChannelCooperativeSettled:
						ev, err := newEventChannelCooperativeSettled(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelCooperativeSettled err %s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive channel cooperative settledd,but it's not our contract,ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelCooperativeSettled2StateChange(ev))
					case params.NameChannelPunished:
						ev, err := newEventChannelPunished(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelPunished err %s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive channel punished event,but it's not our contract,ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelPunished2StateChange(ev))
					case params.NameSecretRevealed:
						ev, err := newEventSecretRevealed(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventSecretRevealed err=%s", err))
							continue
						}
						if ev.Raw.Address != be.SecretRegistryAddress {
							log.Info(fmt.Sprintf("receive NameSecretRevealed,but it's not our contract,ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventSecretRevealed2StateChange(ev))
					case params.NameBalanceProofUpdated:
						ev, err := newEventBalanceProofUpdated(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventBalanceProofUpdated err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive channel balance proof updated ,but it's not our contract,ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventBalanceProofUpdated2StateChange(ev))
					case params.NameChannelWithdraw:
						ev, err := newEventChannelWithdraw(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelWithdraw err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive channel withdraw ,but it's not our contract,ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelWithdraw2StateChange(ev))
					default:
						log.Crit(fmt.Sprintf("receive unkown event %s,it must be a bug", name))
					}
					//event to statechange
				case err := <-sub.Err():
					if !be.stopped {
						log.Error(fmt.Sprintf("eventlistener %s error:%v", name, err))
					}
					return
				case <-be.quitChan:
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
	close(be.quitChan)
	for _, sub := range be.Subscribes {
		sub.Unsubscribe()
	}
	log.Info("Events stop ok...")
}
func (be *Events) sendStateChange(st mediatedtransfer.ContractStateChange) {
	if be.stopped {
		return
	}
	//log.Trace(fmt.Sprintf("send statechange %s", utils.StringInterface(st, 2)))
	if be.historyEventsGot {
		be.StateChangeChannel <- st
	} else {
		be.startupStateChangeChannel <- st
	}

}

//GetAllTokenNetworks returns all the token network,events 本身需要知道所有的 tokennetwork, 这样才能处理相关事件.
func (be *Events) GetAllTokenNetworks(fromBlock int64) (events []*contracts.TokenNetworkRegistryTokenNetworkCreated, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), be.RegistryAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameTokenNetworkCreated, eventAbiMap[params.NameTokenNetworkCreated], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventTokenNetworkCreated(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventTokenNetworkCreated err %s", err))
			continue
		}
		events = append(events, e)
	}
	for _, e := range events {
		be.TokenNetworks[e.TokenNetworkAddress] = true
	}
	return
}

/*
GetChannelNew return's a token network's channel since `fromBlock` on tokenNetworkAddress
if tokenNetworkAddress is empty, return all events have this sigature
如果 channel 特别多,比如十万个,怎么办,
为了防止出现这样的情况,应该一个一个 tokennetwork 获取事件,而不要是一起获取.
*/
func (be *Events) GetChannelNew(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelOpened, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelOpened, eventAbiMap[params.NameChannelOpened], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		var e *contracts.TokenNetworkChannelOpened
		e, err = newEventChannelOpen(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelOpen err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

/*
GetChannelNewAndDeposit return's a token network's channel since `fromBlock` on tokenNetworkAddress
if tokenNetworkAddress is empty, return all events have this sigature
如果 channel 特别多,比如十万个,怎么办,
为了防止出现这样的情况,应该一个一个 tokennetwork 获取事件,而不要是一起获取.
*/
func (be *Events) GetChannelNewAndDeposit(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelOpenedAndDeposit, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelOpenedAndDeposit, eventAbiMap[params.NameChannelOpenedAndDeposit], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		var e *contracts.TokenNetworkChannelOpenedAndDeposit
		e, err = newEventChannelOpenAndDeposit(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelOpen err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelClosed return  channel closed events
func (be *Events) GetChannelClosed(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelClosed, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelClosed, eventAbiMap[params.NameChannelClosed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelClosed(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelClosed err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelSettled return all channel settled events since `fromBlock` on tokenNetworkAddress
//if tokenNetworkAddress is empty, return's all events have this signature
func (be *Events) GetChannelSettled(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelSettled, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelSettled, eventAbiMap[params.NameChannelSettled], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelSettled(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelSettled err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelCooperativeSettled return all channel settled events since `fromBlock` on tokenNetworkAddress
//if tokenNetworkAddress is empty, return's all events have this signature
func (be *Events) GetChannelCooperativeSettled(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelCooperativeSettled, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelCooperativeSettled, eventAbiMap[params.NameChannelCooperativeSettled], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelCooperativeSettled(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelSettled err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelPunished punish events of contract
func (be *Events) GetChannelPunished(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelPunished, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelPunished, eventAbiMap[params.NameChannelPunished], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelPunished(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelPunished err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelWithdraw return all channel settled events since `fromBlock` on tokenNetworkAddress
//if tokenNetworkAddress is empty, return's all events have this signature
func (be *Events) GetChannelWithdraw(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelWithdraw, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelWithdraw, eventAbiMap[params.NameChannelWithdraw], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelWithdraw(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelSettled err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelNewDeposit return all channel settled events since `fromBlock` on tokenNetworkAddress
//if tokenNetworkAddress is empty, return's all events have this signature
func (be *Events) GetChannelNewDeposit(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelNewDeposit, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelNewDeposit, eventAbiMap[params.NameChannelNewDeposit], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelNewDeposit(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelSettled err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

//GetChannelUnlocked return all channel settled events since `fromBlock` on tokenNetworkAddress
//if tokenNetworkAddress is empty, return's all events have this signature
func (be *Events) GetChannelUnlocked(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkChannelUnlocked, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelUnlocked, eventAbiMap[params.NameChannelUnlocked], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventChannelUnlocked(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventChannelSettled err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

/*
GetChannelBalanceProofUpdated returns all NonClosing balance proof events since `fromBlock`
*/
func (be *Events) GetChannelBalanceProofUpdated(fromBlock int64, tokenNetworkAddress common.Address) (events []*contracts.TokenNetworkBalanceProofUpdated, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), tokenNetworkAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameBalanceProofUpdated, eventAbiMap[params.NameBalanceProofUpdated], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventBalanceProofUpdated(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventBalanceProofUpdated err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

/*
GetAllSecretRevealed return all secret reveal events
*/
func (be *Events) GetAllSecretRevealed(fromBlock int64) (events []*contracts.SecretRegistrySecretRevealed, err error) {
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), be.SecretRegistryAddress, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameSecretRevealed, eventAbiMap[params.NameSecretRevealed], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventSecretRevealed(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventSecretRevealed err %s", err))
			continue
		}
		events = append(events, e)
	}
	return
}

/*
GetAllStateChangeSince returns all the statechanges that raiden should know when it's offline
tokennetwork合约上发生的所有事情我们都应该按顺序通知使用者
*/
func (be *Events) GetAllStateChangeSince(lastBlockNumber int64) (stateChangs []mediatedtransfer.ContractStateChange, err error) {
	events0, err := be.GetAllTokenNetworks(lastBlockNumber)
	if err != nil {
		return
	}
	for _, e := range events0 {
		stateChangs = append(stateChangs, EventTokenNetworkCreated2StateChange(e))
	}
	var events []*contracts.SecretRegistrySecretRevealed
	events, err = be.GetAllSecretRevealed(lastBlockNumber)
	if err != nil {
		return
	}
	for _, e := range events {
		stateChangs = append(stateChangs, EventSecretRevealed2StateChange(e))
	}
	/*
		把历史发生的事件按照顺序通知给 raidenService,
		如何处理在查询过程中新收到的事件呢?
	*/
	for tokenNetwork := range be.TokenNetworks {
		var err error
		events2, err := be.GetChannelNew(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events2 {
			stateChangs = append(stateChangs, EventChannelOpen2StateChange(e))
		}
		events3, err := be.GetChannelClosed(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events3 {
			stateChangs = append(stateChangs, EventChannelClosed2StateChange(e))
		}
		events4, err := be.GetChannelSettled(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events4 {
			stateChangs = append(stateChangs, EventChannelSettled2StateChange(e))
		}
		events5, err := be.GetChannelCooperativeSettled(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events5 {
			stateChangs = append(stateChangs, EventChannelCooperativeSettled2StateChange(e))
		}
		events6, err := be.GetChannelBalanceProofUpdated(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events6 {
			stateChangs = append(stateChangs, EventBalanceProofUpdated2StateChange(e))
		}
		events7, err := be.GetChannelUnlocked(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events7 {
			stateChangs = append(stateChangs, EventChannelUnlocked2StateChange(e))
		}
		events8, err := be.GetChannelWithdraw(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events8 {
			stateChangs = append(stateChangs, EventChannelWithdraw2StateChange(e))
		}
		events9, err := be.GetChannelNewDeposit(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events9 {
			stateChangs = append(stateChangs, EventChannelNewDeposit2StateChange(e))
		}
		events10, err := be.GetChannelPunished(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events10 {
			stateChangs = append(stateChangs, EventChannelPunished2StateChange(e))
		}
		events11, err := be.GetChannelNewAndDeposit(lastBlockNumber, tokenNetwork)
		if err != nil {
			return nil, err
		}
		for _, e := range events11 {
			st1, st2 := EventChannelOpenAndDeposit2StateChange(e)
			stateChangs = append(stateChangs, st1, st2)
		}
	}
	return
}

/*
Start listening events send to  channel can duplicate but cannot lose.
1. first resend events may lost (duplicat is ok)
2. listen new events on blockchain
有可能启动的时候没联网,等到启动以后某个事件连上了以后在处理.
1.要保证事件按照顺序抵达
2. 保证事件不丢失
3. 事件是可以重复的
*/
func (be *Events) Start(LastBlockNumber int64) error {
	log.Info(fmt.Sprintf("get state change since %d", LastBlockNumber))
	firstStartup := false
	err := be.installEventListener()
	if err != nil {
		return err
	}
	oldstateChanges, err := be.GetAllStateChangeSince(LastBlockNumber)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("get state change since %d complete", LastBlockNumber))
	if !be.historyEventsGot {
		//程序第一次启动,需要等初始化数据完成以后,通知上层
		firstStartup = true
	}
	be.historyEventsGot = true

	go func() {
		var subScribeStateChanges []mediatedtransfer.ContractStateChange
		var hasStateChanges = true
		for {
			if !hasStateChanges {
				break
			}
			select {
			case st := <-be.startupStateChangeChannel:
				subScribeStateChanges = append(subScribeStateChanges, st)
			default:
				hasStateChanges = false
				break
			}
		}
		oldstateChanges = append(oldstateChanges, subScribeStateChanges...)
		//保证按序通知
		sortContractStateChange(oldstateChanges)
		for _, st := range oldstateChanges {
			be.sendStateChange(st)
		}
		if firstStartup {
			be.sendStateChange(new(mediatedtransfer.FakeContractInfoCompleteStateChange))
		}
	}()
	return nil
}
