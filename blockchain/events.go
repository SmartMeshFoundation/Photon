package blockchain

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts/monitoringcontracts"
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
	stopped               bool // has stopped?
	quitChan              chan struct{}
	TokenNetworks         map[common.Address]bool
	hasInited             bool
}

//NewBlockChainEvents create BlockChainEvents
func NewBlockChainEvents(client *helper.SafeEthClient, registryAddress, secretRegistryAddress common.Address) *Events {
	be := &Events{
		client:                client,
		LogChannelMap:         make(map[string]chan types.Log),
		Subscribes:            make(map[string]ethereum.Subscription),
		RegistryAddress:       registryAddress,
		SecretRegistryAddress: secretRegistryAddress,
		quitChan:              make(chan struct{}),
		TokenNetworks:         make(map[common.Address]bool),
	}
	for name := range eventAbiMap {
		be.LogChannelMap[name] = make(chan types.Log, 10)
	}
	return be
}

var eventAbiMap = map[string]string{
	params.NameTokenNetworkCreated:       contracts.TokenNetworkRegistryABI,
	params.NameChannelOpened:             contracts.TokenNetworkABI,
	params.NameChannelNewDeposit:         contracts.TokenNetworkABI,
	params.NameChannelClosed:             contracts.TokenNetworkABI,
	params.NameChannelSettled:            contracts.TokenNetworkABI,
	params.NameChannelCooperativeSettled: contracts.TokenNetworkABI,
	params.NameChannelWithdraw:           contracts.TokenNetworkABI,
	params.NameChannelUnlocked:           contracts.TokenNetworkABI,
	params.NameBalanceProofUpdated:       contracts.TokenNetworkABI,
	params.NameSecretRevealed:            contracts.SecretRegistryABI,
	//the following event is for 3rd party
	params.NameNewDeposit:              monitoringcontracts.MonitoringServiceABI,
	params.NameNewBalanceProofReceived: monitoringcontracts.MonitoringServiceABI,
	params.NameRewardClaimed:           monitoringcontracts.MonitoringServiceABI,
	params.NameWithdrawn:               monitoringcontracts.MonitoringServiceABI,
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
				defer rpanic.PanicRecover("installEventListener")
				select {
				case _, ok := <-c:
					if ok {
						//eventlistener need reinstall
						err = be.installEventListener()
						if err != nil {
							log.Error(fmt.Sprintf("installEventListener err %s", err))
						}
					}
				case <-be.quitChan:
					return
				}

			}()
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
func EventChannelSettled2StateChange(ev *contracts.TokenNetworkChannelSettled) *mediatedtransfer.ContractReceiveSettledStateChange {
	return &mediatedtransfer.ContractReceiveSettledStateChange{
		ChannelIdentifier:   common.Hash(ev.Channel_identifier),
		TokenNetworkAddress: ev.Raw.Address,
		SettledBlock:        int64(ev.Raw.BlockNumber),
	}
}

//EventChannelCooperativeSettled2StateChange to stateChange
func EventChannelCooperativeSettled2StateChange(ev *contracts.TokenNetworkChannelCooperativeSettled) *mediatedtransfer.ContractReceiveCooperativeSettledStateChange {
	return &mediatedtransfer.ContractReceiveCooperativeSettledStateChange{
		ChannelIdentifier:   common.Hash(ev.Channel_identifier),
		TokenNetworkAddress: ev.Raw.Address,
		SettledBlock:        int64(ev.Raw.BlockNumber),
	}
}

//EventChannelWithdraw2StateChange to stateChange
func EventChannelWithdraw2StateChange(ev *contracts.TokenNetworkChannelWithdraw) *mediatedtransfer.ContractReceiveChannelWithdrawStateChange {
	return &mediatedtransfer.ContractReceiveChannelWithdrawStateChange{
		ChannelAddress: &contracts.ChannelUniqueID{

			ChannelIdentifier: common.Hash(ev.Channel_identifier),
			OpenBlockNumber:   int64(ev.Raw.BlockNumber),
		},
		TokenNetworkAddress: ev.Raw.Address,
		Participant1:        ev.Participant1,
		Participant2:        ev.Participant2,
		Participant1Balance: ev.Participant1_balance,
		Participant2Balance: ev.Participant2_balance,
	}
}

//EventTokenNetworkCreated2StateChange to statechange
func EventTokenNetworkCreated2StateChange(ev *contracts.TokenNetworkRegistryTokenNetworkCreated) *mediatedtransfer.ContractReceiveTokenAddedStateChange {
	return &mediatedtransfer.ContractReceiveTokenAddedStateChange{
		RegistryAddress:     ev.Raw.Address,
		TokenAddress:        ev.Token_address,
		TokenNetworkAddress: ev.Token_network_address,
	}
}

//EventChannelOpen2StateChange to statechange
func EventChannelOpen2StateChange(ev *contracts.TokenNetworkChannelOpened) *mediatedtransfer.ContractReceiveNewChannelStateChange {
	return &mediatedtransfer.ContractReceiveNewChannelStateChange{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: ev.Channel_identifier,

			OpenBlockNumber: int64(ev.Raw.BlockNumber),
		},
		TokenNetworkAddress: ev.Raw.Address,
		Participant1:        ev.Participant1,
		Participant2:        ev.Participant2,
		SettleTimeout:       int(ev.Settle_timeout.Int64()),
	}
}

//EventChannelNewBalance2StateChange to statechange
func EventChannelNewBalance2StateChange(ev *contracts.TokenNetworkChannelNewDeposit) *mediatedtransfer.ContractReceiveBalanceStateChange {
	return &mediatedtransfer.ContractReceiveBalanceStateChange{

		ChannelIdentifier:   ev.Channel_identifier,
		TokenNetworkAddress: ev.Raw.Address,
		ParticipantAddress:  ev.Participant,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
}

//EventChannelClosed2StateChange to statechange
func EventChannelClosed2StateChange(ev *contracts.TokenNetworkChannelClosed) *mediatedtransfer.ContractReceiveClosedStateChange {
	return &mediatedtransfer.ContractReceiveClosedStateChange{
		TokenNetworkAddress: ev.Raw.Address,
		ChannelIdentifier:   ev.Channel_identifier,
		ClosingAddress:      ev.Closing_participant,
		ClosedBlock:         int64(ev.Raw.BlockNumber),
	}
}

//EventBalanceProofUpdated2StateChange to statechange
func EventBalanceProofUpdated2StateChange(ev *contracts.TokenNetworkBalanceProofUpdated) *mediatedtransfer.ContractBalanceProofUpdatedStateChange {
	return &mediatedtransfer.ContractBalanceProofUpdatedStateChange{
		TokenNetworkAddress: ev.Raw.Address,
		ChannelIdentifier:   ev.Channel_identifier,
		LocksRoot:           ev.Locksroot,
		TransferredAmount:   ev.Transferred_amount,
		Participant:         ev.Participant,
	}
}

//EventSecretRevealed2StateChange to statechange
func EventSecretRevealed2StateChange(ev *contracts.SecretRegistrySecretRevealed) *mediatedtransfer.ContractSecretRevealStateChange {
	return &mediatedtransfer.ContractSecretRevealStateChange{
		Secret: ev.Secrethash,
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
						be.TokenNetworks[ev.Token_network_address] = true
						be.sendStateChange(EventTokenNetworkCreated2StateChange(ev))
					case params.NameChannelOpened:
						ev, err := newEventEventChannelOpen(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventEventChannelOpen err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive event ChannelOpened, but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelOpen2StateChange(ev))
					case params.NameChannelNewDeposit:
						ev, err := newEventChannelNewBalance(&l)
						if err != nil {
							log.Error(fmt.Sprintf("newEventChannelNewBalance err=%s", err))
							continue
						}
						if !be.TokenNetworks[ev.Raw.Address] {
							log.Info(fmt.Sprintf("receive event channel new deposit ,but it's not our contract, ev=\n%s", utils.StringInterface(ev, 3)))
							continue
						}
						be.sendStateChange(EventChannelNewBalance2StateChange(ev))
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
func (be *Events) sendStateChange(st transfer.StateChange) {
	if be.stopped {
		return
	}
	be.StateChangeChannel <- st
}

//GetAllTokenNetworks returns all the token network
func (be *Events) GetAllTokenNetworks(fromBlock int64) (tokens map[common.Address]common.Address, err error) {
	var events []*contracts.TokenNetworkRegistryTokenNetworkCreated
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
	if len(events) > 0 {
		tokens = make(map[common.Address]common.Address)
	}
	for _, e := range events {
		tokens[e.Token_address] = e.Token_network_address
		be.TokenNetworks[e.Token_network_address] = true
	}
	be.hasInited = true
	return
}

/*
GetAllChannels return's a token network's channel since `fromBlock` on tokenNetworkAddress
if tokenNetworkAddress is empty, return all events have this sigature
如果 channel 特别多,比如十万个,怎么办
*/
func (be *Events) GetAllChannels(fromBlock int64, tokenNetworkAddress ...common.Address) (channels []*contracts.ChannelUniqueID, err error) {
	var events []*contracts.TokenNetworkChannelOpened
	addr := utils.EmptyAddress
	if len(tokenNetworkAddress) > 0 {
		addr = tokenNetworkAddress[0]
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), addr, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
		params.NameChannelOpened, eventAbiMap[params.NameChannelOpened], be.client)
	if err != nil {
		return
	}
	for _, l := range logs {
		e, err := newEventEventChannelOpen(&l)
		if err != nil {
			log.Error(fmt.Sprintf("newEventEventChannelOpen err %s", err))
			continue
		}
		events = append(events, e)
	}
	for _, e := range events {
		c := &contracts.ChannelUniqueID{
			ChannelIdentifier: e.Channel_identifier,
			OpenBlockNumber:   int64(e.Raw.BlockNumber),
		}
		channels = append(channels, c)
	}
	return
}

//GetAllChannelClosed return  channel closed events
func (be *Events) GetAllChannelClosed(fromBlock int64, tokenNetworkAddress ...common.Address) (events []*contracts.TokenNetworkChannelClosed, err error) {
	addr := utils.EmptyAddress
	if len(tokenNetworkAddress) > 0 {
		addr = tokenNetworkAddress[0]
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), addr, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
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

//GetAllChannelSettled return all channel settled events since `fromBlock` on tokenNetworkAddress
//if tokenNetworkAddress is empty, return's all events have this signature
func (be *Events) GetAllChannelSettled(fromBlock int64, tokenNetworkAddress ...common.Address) (events []*contracts.TokenNetworkChannelSettled, err error) {
	addr := utils.EmptyAddress
	if len(tokenNetworkAddress) > 0 {
		addr = tokenNetworkAddress[0]
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), addr, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
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

/*
GetAllChannelNonClosingBalanceProofUpdated returns all NonClosing balance proof events since `fromBlock`
*/
func (be *Events) GetAllChannelNonClosingBalanceProofUpdated(fromBlock int64, tokenNetworkAddress ...common.Address) (events []*contracts.TokenNetworkBalanceProofUpdated, err error) {
	addr := utils.EmptyAddress
	if len(tokenNetworkAddress) > 0 {
		addr = tokenNetworkAddress[0]
	}
	logs, err := rpc.EventGetInternal(rpc.GetQueryConext(), addr, ethrpc.BlockNumber(fromBlock), ethrpc.LatestBlockNumber,
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
*/
func (be *Events) GetAllStateChangeSince(lastBlockNumber int64) (stateChangs []transfer.StateChange, err error) {
	if !be.hasInited {
		be.GetAllTokenNetworks(lastBlockNumber)
	}
	var events []*contracts.SecretRegistrySecretRevealed
	events, err = be.GetAllSecretRevealed(lastBlockNumber)
	if err != nil {
		return
	}
	for _, e := range events {
		stateChangs = append(stateChangs, EventSecretRevealed2StateChange(e))
	}
	events2, err := be.GetAllChannelNonClosingBalanceProofUpdated(lastBlockNumber)
	for _, e := range events2 {
		if be.TokenNetworks[e.Raw.Address] {
			stateChangs = append(stateChangs, EventBalanceProofUpdated2StateChange(e))
		}
	}
	return
}

/*
Start listening events send to  channel can duplicate but cannot lose.
1. first resend events may lost (duplicat is ok)
2. listen new events on blockchain
*/
func (be *Events) Start(LastBlockNumber int64) error {
	stateChanges, err := be.GetAllStateChangeSince(LastBlockNumber)
	if err != nil {
		return err
	}
	//maybe too much state changes?
	be.StateChangeChannel = make(chan transfer.StateChange, len(stateChanges)+20)
	for _, st := range stateChanges {
		be.sendStateChange(st)
	}
	return be.installEventListener()
}
