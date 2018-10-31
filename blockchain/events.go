package blockchain

import (
	"context"
	"fmt"

	"time"

	"math/big"

	"strings"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

var secretRegistryAbi abi.ABI
var tokenNetworkRegistryAbi abi.ABI
var tokenNetworkAbi abi.ABI
var topicToEventName map[common.Hash]string

func init() {
	var err error
	tokenNetworkRegistryAbi, err = abi.JSON(strings.NewReader(contracts.TokenNetworkRegistryABI))
	if err != nil {
		panic(fmt.Sprintf("tokenNetworkRegistryAbi parse err %s", err))
	}
	secretRegistryAbi, err = abi.JSON(strings.NewReader(contracts.SecretRegistryABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	tokenNetworkAbi, err = abi.JSON(strings.NewReader(contracts.TokenNetworkABI))
	if err != nil {
		panic(fmt.Sprintf("tokenNetworkAbi parse err %s", err))
	}
	topicToEventName = make(map[common.Hash]string)
	topicToEventName[tokenNetworkRegistryAbi.Events[params.NameTokenNetworkCreated].Id()] = params.NameTokenNetworkCreated
	topicToEventName[secretRegistryAbi.Events[params.NameSecretRevealed].Id()] = params.NameSecretRevealed
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelOpened].Id()] = params.NameChannelOpened
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelOpenedAndDeposit].Id()] = params.NameChannelOpenedAndDeposit
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelNewDeposit].Id()] = params.NameChannelNewDeposit
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelWithdraw].Id()] = params.NameChannelWithdraw
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelClosed].Id()] = params.NameChannelClosed
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelPunished].Id()] = params.NameChannelPunished
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelUnlocked].Id()] = params.NameChannelUnlocked
	topicToEventName[tokenNetworkAbi.Events[params.NameBalanceProofUpdated].Id()] = params.NameBalanceProofUpdated
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelSettled].Id()] = params.NameChannelSettled
	topicToEventName[tokenNetworkAbi.Events[params.NameChannelCooperativeSettled].Id()] = params.NameChannelCooperativeSettled

}

/*
Events handles all contract events from blockchain
*/
type Events struct {
	StateChangeChannel chan transfer.StateChange

	tokenNetworks       map[common.Address]bool
	lastBlockNumber     int64
	rpcModuleDependency RPCModuleDependency
	client              *helper.SafeEthClient
	pollPeriod          time.Duration          // 轮询周期,必须与公链出块间隔一致
	stopped             bool                   // has stopped?
	txDone              map[common.Hash]uint64 // 该map记录最近30块内处理的events流水,用于事件去重
}

//NewBlockChainEvents create BlockChainEvents
func NewBlockChainEvents(client *helper.SafeEthClient, rpcModuleDependency RPCModuleDependency,
	token2TokenNetwork map[common.Address]common.Address) *Events {
	be := &Events{
		StateChangeChannel:  make(chan transfer.StateChange, 10),
		tokenNetworks:       make(map[common.Address]bool),
		rpcModuleDependency: rpcModuleDependency,
		client:              client,
		txDone:              make(map[common.Hash]uint64),
	}
	if token2TokenNetwork != nil {
		for _, tn := range token2TokenNetwork {
			be.tokenNetworks[tn] = true
		}
	}
	return be
}

//Stop event listenging
func (be *Events) Stop() {
	be.stopped = true
	log.Info("Events stop ok...")
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
/*
 *  Start listening events send to channel can duplicate but cannot lose.
 *  1. first resend events may lost (duplicate is ok)
 *  2. listen new events on blockchain
 *
 *  It is possible that there is no internet connection when start-up, and missed events have to be regained
 *  after those events starts.
 * 	1. Make sure events sending out with order
 *  2. Make sure events does not get lost.
 *  3. Make sure repeated events are allowed.
 */
func (be *Events) Start(LastBlockNumber int64) {
	log.Info(fmt.Sprintf("get state change since %d", LastBlockNumber))
	be.lastBlockNumber = LastBlockNumber
	be.stopped = false
	/*
		1. start alarm task
	*/
	go be.startAlarmTask()
}

func (be *Events) startAlarmTask() {
	log.Trace(fmt.Sprintf("start getting lasted block number from blocknubmer=%d", be.lastBlockNumber))
	currentBlock := be.lastBlockNumber
	for {
		//get the lastest number imediatelly
		if be.stopped {
			log.Info(fmt.Sprintf("AlarmTask quit complete"))
			return
		}
		if be.pollPeriod == 0 {
			// first time
			if params.ChainID.Int64() == params.TestPrivateChainID {
				be.pollPeriod = params.DefaultEthRPCPollPeriodForTest
			} else {
				be.pollPeriod = params.DefaultEthRPCPollPeriod
			}
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), params.EthRPCTimeout)
		h, err := be.client.HeaderByNumber(ctx, nil)
		if err != nil {
			log.Error(fmt.Sprintf("HeaderByNumber err=%s", err))
			cancelFunc()
			if !be.stopped {
				go be.client.RecoverDisconnect()
			}
			return
		}
		cancelFunc()
		lastedBlock := h.Number.Int64()
		if currentBlock == lastedBlock {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if currentBlock != -1 && lastedBlock != currentBlock+1 {
			log.Warn(fmt.Sprintf("AlarmTask missed %d blocks", lastedBlock-currentBlock-1))
		}
		if lastedBlock%10 == 0 {
			log.Trace(fmt.Sprintf("new block :%d", lastedBlock))
		}

		fromBlockNumber := currentBlock - 2*params.ForkConfirmNumber
		if fromBlockNumber < 0 {
			fromBlockNumber = 0
		}
		// get all state change between currentBlock and lastedBlock
		stateChanges, err := be.queryAllStateChange(fromBlockNumber, lastedBlock)
		if err != nil {
			log.Error(fmt.Sprintf("queryAllStateChange err=%s", err))
			if be.stopped {
				return
			}
			// 如果这里出现err,不能继续处理该blocknumber,否则会丢事件,直接从该块重新处理即可
			continue
		}
		if len(stateChanges) > 0 {
			log.Trace(fmt.Sprintf("receive %d events between block %d - %d", len(stateChanges), currentBlock+1, lastedBlock))
		}

		// refresh block number and notify PhotonService
		currentBlock = lastedBlock
		be.lastBlockNumber = currentBlock

		be.StateChangeChannel <- &transfer.BlockStateChange{BlockNumber: currentBlock}

		// notify Photon service
		for _, sc := range stateChanges {
			be.StateChangeChannel <- sc
		}

		// 清除过期流水
		for key, blockNumber := range be.txDone {
			if blockNumber <= uint64(fromBlockNumber) {
				delete(be.txDone, key)
			}
		}
		// wait to next time
		time.Sleep(be.pollPeriod)
	}
}

func (be *Events) queryAllStateChange(fromBlock int64, toBlock int64) (stateChanges []mediatedtransfer.ContractStateChange, err error) {
	/*
		get all event of contract TokenNetworkRegistry, SecretRegistry , TokenNetwork
	*/
	logs, err := be.getLogsFromChain(fromBlock, toBlock)
	if err != nil {
		return
	}
	stateChanges, err = be.parseLogsToEvents(logs)
	if err != nil {
		return
	}
	// 排序
	sortContractStateChange(stateChanges)
	return
}

func (be *Events) getLogsFromChain(fromBlock int64, toBlock int64) (logs []types.Log, err error) {
	/*
		get all event of contract TokenNetworkRegistry, SecretRegistry , TokenNetwork
	*/
	contractAddresses := []common.Address{
		be.rpcModuleDependency.GetRegistryAddress(),
		be.rpcModuleDependency.GetSecretRegistryAddress(),
	}
	for tokenNetworkAddress := range be.tokenNetworks {
		contractAddresses = append(contractAddresses, tokenNetworkAddress)
	}
	logs, err = rpc.EventsGetInternal(
		rpc.GetQueryConext(), contractAddresses, ethrpc.BlockNumber(fromBlock), ethrpc.BlockNumber(toBlock), be.client)
	if err != nil {
		return
	}
	var newTokenNetworks []common.Address
	for _, l := range logs {
		if topicToEventName[l.Topics[0]] == params.NameTokenNetworkCreated {
			e, err2 := newEventTokenNetworkCreated(&l)
			if err = err2; err != nil {
				return
			}
			newTokenNetworks = append(newTokenNetworks, e.TokenNetworkAddress)
		}
	}
	if len(newTokenNetworks) > 0 {
		for _, tokenNetworkAddress := range newTokenNetworks {
			be.tokenNetworks[tokenNetworkAddress] = true
		}
		var newLogs []types.Log
		newLogs, err = rpc.EventsGetInternal(
			rpc.GetQueryConext(), newTokenNetworks, ethrpc.BlockNumber(fromBlock), ethrpc.BlockNumber(toBlock), be.client)
		if err != nil {
			return
		}
		logs = append(logs, newLogs...)
	}
	return
}

func (be *Events) parseLogsToEvents(logs []types.Log) (stateChanges []mediatedtransfer.ContractStateChange, err error) {
	for _, l := range logs {
		eventName := topicToEventName[l.Topics[0]]

		// 根据已处理流水去重
		if doneBlockNumber, ok := be.txDone[l.TxHash]; ok {
			if doneBlockNumber == l.BlockNumber {
				//log.Trace(fmt.Sprintf("get event txhash=%s repeated,ignore...", l.TxHash.String()))
				continue
			}
			log.Warn(fmt.Sprintf("event tx=%s happened at %d, but now happend at %d ", l.TxHash.String(), doneBlockNumber, l.BlockNumber))
		}
		/*
			if needConfirm {
				if be.lastBlockNumber - l.BlockNumber < 15 {
					// 待确认,暂不处理
					continue
				}
			}
			// 已确认,直接处理上报并记录处理流水
		*/

		switch eventName {
		case params.NameTokenNetworkCreated:
			e, err2 := newEventTokenNetworkCreated(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventTokenNetworkCreated2StateChange(e))
		case params.NameSecretRevealed:
			e, err2 := newEventSecretRevealed(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventSecretRevealed2StateChange(e))
		case params.NameChannelOpened:
			e, err2 := newEventChannelOpen(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelOpen2StateChange(e))
		case params.NameChannelOpenedAndDeposit:
			e, err2 := newEventChannelOpenAndDeposit(&l)
			if err = err2; err != nil {
				return
			}
			oev, dev := eventChannelOpenAndDeposit2StateChange(e)
			stateChanges = append(stateChanges, oev)
			stateChanges = append(stateChanges, dev)
		case params.NameChannelNewDeposit:
			e, err2 := newEventChannelNewDeposit(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelNewDeposit2StateChange(e))
		case params.NameChannelClosed:
			e, err2 := newEventChannelClosed(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelClosed2StateChange(e))
		case params.NameChannelUnlocked:
			e, err2 := newEventChannelUnlocked(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelUnlocked2StateChange(e))
		case params.NameBalanceProofUpdated:
			e, err2 := newEventBalanceProofUpdated(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventBalanceProofUpdated2StateChange(e))
		case params.NameChannelPunished:
			e, err2 := newEventChannelPunished(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelPunished2StateChange(e))
		case params.NameChannelSettled:
			e, err2 := newEventChannelSettled(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelSettled2StateChange(e))
		case params.NameChannelCooperativeSettled:
			e, err2 := newEventChannelCooperativeSettled(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelCooperativeSettled2StateChange(e))
		case params.NameChannelWithdraw:
			e, err2 := newEventChannelWithdraw(&l)
			if err = err2; err != nil {
				return
			}
			stateChanges = append(stateChanges, eventChannelWithdraw2StateChange(e))
		default:
			log.Warn(fmt.Sprintf("receive unkonwn type event from chain : \n%s\n", l.String()))
		}
		// 记录处理流水
		be.txDone[l.TxHash] = l.BlockNumber
	}
	return
}

//eventChannelSettled2StateChange to stateChange
func eventChannelSettled2StateChange(ev *contracts.TokenNetworkChannelSettled) *mediatedtransfer.ContractSettledStateChange {
	return &mediatedtransfer.ContractSettledStateChange{
		ChannelIdentifier:   common.Hash(ev.ChannelIdentifier),
		TokenNetworkAddress: ev.Raw.Address,
		SettledBlock:        int64(ev.Raw.BlockNumber),
	}
}

//eventChannelCooperativeSettled2StateChange to stateChange
func eventChannelCooperativeSettled2StateChange(ev *contracts.TokenNetworkChannelCooperativeSettled) *mediatedtransfer.ContractCooperativeSettledStateChange {
	return &mediatedtransfer.ContractCooperativeSettledStateChange{
		ChannelIdentifier:   common.Hash(ev.ChannelIdentifier),
		TokenNetworkAddress: ev.Raw.Address,
		SettledBlock:        int64(ev.Raw.BlockNumber),
	}
}

//eventChannelPunished2StateChange to stateChange
func eventChannelPunished2StateChange(ev *contracts.TokenNetworkChannelPunished) *mediatedtransfer.ContractPunishedStateChange {
	return &mediatedtransfer.ContractPunishedStateChange{
		ChannelIdentifier:   common.Hash(ev.ChannelIdentifier),
		TokenNetworkAddress: ev.Raw.Address,
		Beneficiary:         ev.Beneficiary,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
}

//eventChannelWithdraw2StateChange to stateChange
func eventChannelWithdraw2StateChange(ev *contracts.TokenNetworkChannelWithdraw) *mediatedtransfer.ContractChannelWithdrawStateChange {
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

//eventTokenNetworkCreated2StateChange to statechange
func eventTokenNetworkCreated2StateChange(ev *contracts.TokenNetworkRegistryTokenNetworkCreated) *mediatedtransfer.ContractTokenAddedStateChange {
	return &mediatedtransfer.ContractTokenAddedStateChange{
		RegistryAddress:     ev.Raw.Address,
		TokenAddress:        ev.TokenAddress,
		TokenNetworkAddress: ev.TokenNetworkAddress,
		BlockNumber:         int64(ev.Raw.BlockNumber),
	}
}

//eventChannelOpen2StateChange to statechange
func eventChannelOpen2StateChange(ev *contracts.TokenNetworkChannelOpened) *mediatedtransfer.ContractNewChannelStateChange {
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

//eventChannelOpenAndDeposit2StateChange to statechange
func eventChannelOpenAndDeposit2StateChange(ev *contracts.TokenNetworkChannelOpenedAndDeposit) (ch1 *mediatedtransfer.ContractNewChannelStateChange, ch2 *mediatedtransfer.ContractBalanceStateChange) {
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

//eventChannelNewDeposit2StateChange to statechange
func eventChannelNewDeposit2StateChange(ev *contracts.TokenNetworkChannelNewDeposit) *mediatedtransfer.ContractBalanceStateChange {
	return &mediatedtransfer.ContractBalanceStateChange{
		ChannelIdentifier:   ev.ChannelIdentifier,
		TokenNetworkAddress: ev.Raw.Address,
		ParticipantAddress:  ev.Participant,
		BlockNumber:         int64(ev.Raw.BlockNumber),
		Balance:             ev.TotalDeposit,
	}
}

//eventChannelClosed2StateChange to statechange
func eventChannelClosed2StateChange(ev *contracts.TokenNetworkChannelClosed) *mediatedtransfer.ContractClosedStateChange {
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

//eventBalanceProofUpdated2StateChange to statechange
func eventBalanceProofUpdated2StateChange(ev *contracts.TokenNetworkBalanceProofUpdated) *mediatedtransfer.ContractBalanceProofUpdatedStateChange {
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

//eventChannelUnlocked2StateChange to statechange
func eventChannelUnlocked2StateChange(ev *contracts.TokenNetworkChannelUnlocked) *mediatedtransfer.ContractUnlockStateChange {
	c := &mediatedtransfer.ContractUnlockStateChange{
		TokenNetworkAddress: ev.Raw.Address,
		ChannelIdentifier:   ev.ChannelIdentifier,
		BlockNumber:         int64(ev.Raw.BlockNumber),
		TransferAmount:      ev.TransferredAmount,
		Participant:         ev.PayerParticipant,
		LockHash:            ev.Lockhash,
	}
	if c.TransferAmount == nil {
		c.TransferAmount = new(big.Int)
	}
	return c
}

//eventSecretRevealed2StateChange to statechange
func eventSecretRevealed2StateChange(ev *contracts.SecretRegistrySecretRevealed) *mediatedtransfer.ContractSecretRevealOnChainStateChange {
	return &mediatedtransfer.ContractSecretRevealOnChainStateChange{
		Secret:         ev.Secret,
		BlockNumber:    int64(ev.Raw.BlockNumber),
		LockSecretHash: utils.ShaSecret(ev.Secret[:]),
	}
}
