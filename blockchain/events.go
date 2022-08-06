package blockchain

import (
	"bytes"
	"context"
	"fmt"

	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/rerr"

	"time"

	"math/big"

	"strings"

	"strconv"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
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
)

var secretRegistryAbi abi.ABI
var tokenNetworkAbi abi.ABI
var topicToEventName map[common.Hash]string

func init() {
	var err error
	secretRegistryAbi, err = abi.JSON(strings.NewReader(contracts.SecretRegistryABI))
	if err != nil {
		panic(fmt.Sprintf("secretRegistryAbi parse err %s", err))
	}
	tokenNetworkAbi, err = abi.JSON(strings.NewReader(contracts.TokensNetworkABI))
	if err != nil {
		panic(fmt.Sprintf("tokenNetworkAbi parse err %s", err))
	}
	topicToEventName = make(map[common.Hash]string)
	topicToEventName[tokenNetworkAbi.Events[params.NameTokenNetworkCreated].Id()] = params.NameTokenNetworkCreated
	topicToEventName[secretRegistryAbi.Events[params.NameSecretRevealed].Id()] = params.NameSecretRevealed
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

type eventID [25]byte //txHash+logIndex
//假定一个tx中事件不可能超过256
func makeEventID(l *types.Log) eventID {
	var e eventID
	copy(e[:], l.TxHash[:])
	e[24] = byte(l.Index)
	return e
}

/*
Events handles all contract events from blockchain
*/
type Events struct {
	StateChangeChannel       chan transfer.StateChange
	lastBlockNumber          int64
	lastBlockNumberTimestamp int64
	isChainEffective         bool
	rpcModuleDependency      RPCModuleDependency
	client                   *helper.SafeEthClient
	stopChan                 chan int                   // has stopped?
	txDone                   map[eventID]uint64         // 该map记录最近30块内处理的events流水,用于事件去重
	firstStart               bool                       //保证ContractHistoryEventCompleteStateChange 只会发送一次
	chainEventRecordDao      models.ChainEventRecordDao // 事件处理记录保存
}

//NewBlockChainEvents create BlockChainEvents
func NewBlockChainEvents(client *helper.SafeEthClient, rpcModuleDependency RPCModuleDependency, chainEventRecordDao models.ChainEventRecordDao) *Events {
	be := &Events{
		StateChangeChannel:  make(chan transfer.StateChange, 10),
		rpcModuleDependency: rpcModuleDependency,
		client:              client,
		txDone:              make(map[eventID]uint64),
		firstStart:          true,
		chainEventRecordDao: chainEventRecordDao,
	}
	return be
}

//Stop event listenging
func (be *Events) Stop() {
	if be.stopChan != nil {
		close(be.stopChan)
	}
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
	/*
		1. start alarm task
	*/
	go be.startAlarmTask()
}

/*
ChangeEthRPCEndpointPort 临时更换连接的公链端口,并重启AlarmTask,for test
*/
func (be *Events) ChangeEthRPCEndpointPort(newPort int) {
	s := strings.Split(be.client.URL, ":")
	s[2] = strconv.Itoa(newPort)
	be.client.URL = strings.Join(s, ":")
	log.Info(fmt.Sprintf("eth-rpc-endpoint change to %s", be.client.URL))
	// 触发重启
	if be.client.IsConnected() {
		be.client.Client.Close()
	}
}

func (be *Events) notifyPhotonStartupCompleteIfNeeded(currentBlock int64) {
	if be.firstStart {
		be.firstStart = false
		//通知photon,历史消息处理完毕,可以进行后续启动了.
		be.StateChangeChannel <- &mediatedtransfer.ContractHistoryEventCompleteStateChange{
			BlockNumber: currentBlock,
		}
	}
}

func (be *Events) startAlarmTask() {
	log.Info(fmt.Sprintf("start getting lasted block number from blocknubmer=%d", be.lastBlockNumber))
	rpanic.PanicRecover("startAlarmTask")
	startUpBlockNumber := be.lastBlockNumber
	currentBlock := be.lastBlockNumber
	currentBlockTimestamp := be.lastBlockNumberTimestamp
	logPeriod := int64(1)
	retryTime := 0
	be.stopChan = make(chan int)
	be.StateChangeChannel <- &transfer.BlockStateChange{BlockNumber: currentBlock}
	/*
			正常处理流程:
			1. 抓取历史事件,排序,发送给photon
			2. 通知photon启动完毕
			其他处理流程:
			通知photon启动完毕
			也就是说无论发生了什么错误,尽快通知photon启动完毕,不要卡主.
		考虑公链的各种错误情况:
		1. 公链网络链接无效的情况
		2. 公链物理连接有效,没有切换节点,但是公链正在同步最新块(已知最新块在3分钟之外)
		3. 公链物理连接有效,但是切换公链,并且新的公链落后于上次已处理块数
		4.  公链物理连接有效,但是切换公链,并且新的公链领先一部分上次已处理块数,但是仍然是正在同步中(已知最新块在3分钟之外)的情况
		5. 连接到公链,并且公链有效(和全网最新块保持同步)
	*/
	for {
		//get the lastest number imediatelly
		ctx, cancelFunc := context.WithTimeout(context.Background(), params.Cfg.ChainRequestTimeout)
		h, err := be.client.HeaderByNumber(ctx, nil)
		if err != nil {
			//无论公链发生什么错误,都应该让photon启动起来,而不是卡主
			be.notifyPhotonStartupCompleteIfNeeded(currentBlock)
			log.Warn(fmt.Sprintf("HeaderByNumber err=%s", err))
			cancelFunc()
			if be.stopChan != nil {
				go be.client.RecoverDisconnect()
			}
			// 连接失败直接通知上层切换到无效公链
			be.isChainEffective = false
			be.StateChangeChannel <- &transfer.EffectiveChainStateChange{
				IsEffective:              false,
				LastBlockNumber:          currentBlock,
				LastBlockNumberTimestamp: currentBlockTimestamp,
			}
			return
		}
		cancelFunc()
		lastedBlock := h.Number.Int64()
		lastedBlockTimestamp := h.Time.Int64()
		// 由于测试环境这个值被修改为了毫秒,必须进行转换. 考虑到如果lastedBlockTimestamp当做秒来解释,将会是几万年以后,因此这么做是合理的
		//todo 为9999999999加上注释,并且给一个合理的解释
		if lastedBlockTimestamp > 9999999999 {
			lastedBlockTimestamp = lastedBlockTimestamp / 1000
		}
		now := time.Now().Unix()
		if lastedBlockTimestamp-now > int64(params.Cfg.BlockPeriod) {
			// 如果本地时间小于最新块的出块时间一个出块间隔以上,说明本地时间服务有问题,这种情况下运行photon是不安全的,直接结束photon
			panic(fmt.Sprintf("local time error local=%d lastedBlockTimestamp=%d, please run photon again after you fix local time server", now, lastedBlockTimestamp))
		}
		//连接到了有效的公链链接,但是公链最新块在三分钟之前,可能公链已经停止同步了,直接切换到无效公链状态
		// 由于启动的时候默认无效公链,所以第一次处循环不可能走到该分支
		//if now-lastedBlockTimestamp >= 180 && startUpBlockNumber == currentBlock {
		if now-lastedBlockTimestamp >= params.Cfg.EffectiveTimeoutSeconds && be.isChainEffective {
			// 最新块的出块时间在3分钟以前,说明连接到了一个无效的公链节点,通知上层切换到无效公链
			be.isChainEffective = false
			be.StateChangeChannel <- &transfer.EffectiveChainStateChange{
				IsEffective:              false,
				LastBlockNumber:          lastedBlock,
				LastBlockNumberTimestamp: lastedBlockTimestamp,
			}
		}
		// 校验公链同步状态,如果尚未同步完成,等待即可,如果是刚启动,也允许启动并等待
		err = be.checkSmcSyncDone()
		if err != nil {
			log.Error(fmt.Sprintf("%s,photon will wait ...", err.Error()))
			be.notifyPhotonStartupCompleteIfNeeded(currentBlock)
			time.Sleep(params.Cfg.PollPeriod)
			continue
		}
		// 这里如果出现切换公链导致获取到的新块比当前块更小的话,只需要等待即可
		if currentBlock >= lastedBlock {
			if startUpBlockNumber >= lastedBlock {
				if startUpBlockNumber > lastedBlock {
					log.Error(fmt.Sprintf("photon last processed number is %d,but spectrum's lastest block number  %d", startUpBlockNumber, lastedBlock))
				}
				// 当启动时获取不到新块,也需要通知photonService,否则会导致api无法启动
				log.Warn(fmt.Sprintf("photon start with blockNumber %d,but lastedBlockNumber on chain also %d", startUpBlockNumber, lastedBlock))
				be.StateChangeChannel <- &transfer.BlockStateChange{BlockNumber: currentBlock}
				startUpBlockNumber = 0
			}
			//在启动的时候连接到了一条无效的公链(不出块)的情况下,photon也应该可以继续启动.
			// 连接到另外一个节点,该节点落后很多,也应该让photon尽快启动
			be.notifyPhotonStartupCompleteIfNeeded(currentBlock)
			time.Sleep(params.Cfg.PollPeriod / 2)
			retryTime++
			if retryTime > 10 {
				log.Warn(fmt.Sprintf("get same block number %d from chain %d times,maybe something wrong with smc ...", lastedBlock, retryTime))
			}
			continue
		}
		retryTime = 0
		if currentBlock != -1 && lastedBlock != currentBlock+1 {
			//log.Warn(fmt.Sprintf("AlarmTask missed %d blocks,currentBlock=%d", lastedBlock-currentBlock-1, currentBlock))
		}
		if lastedBlock%logPeriod == 0 {
			log.Info(fmt.Sprintf("new block :%d", lastedBlock))
		}

		fromBlockNumber := currentBlock - 2*params.Cfg.ForkConfirmNumber
		if fromBlockNumber < 0 {
			fromBlockNumber = 0
		}

		if fromBlockNumber == 0 {
			fromBlockNumber = lastedBlock - 5000
		}

		// get all state change between currentBlock and lastedBlock
		stateChanges, err := be.queryAllStateChange(fromBlockNumber, lastedBlock)
		if err != nil {
			log.Error(fmt.Sprintf("queryAllStateChange err=%s", err))
			//无论公链发生什么错误,都应该让photon启动起来,而不是卡主
			be.notifyPhotonStartupCompleteIfNeeded(currentBlock)
			// 如果这里出现err,不能继续处理该blocknumber,否则会丢事件,直接从该块重新处理即可
			time.Sleep(params.Cfg.PollPeriod / 2)
			continue
		}
		if len(stateChanges) > 0 {
			log.Info(fmt.Sprintf("receive %d events between block %d - %d", len(stateChanges), fromBlockNumber, lastedBlock))
		}

		// refresh block number and notify PhotonService
		currentBlock = lastedBlock
		currentBlockTimestamp = lastedBlockTimestamp
		be.lastBlockNumber = currentBlock
		be.lastBlockNumberTimestamp = currentBlockTimestamp
		var lastSendBlockNumber int64
		// notify Photon service
		//我们需要photon service在处理相关事件的时候知道了对应的块已经发生了,否则可能因为错误的当前块数而出现逻辑错误.
		//同时也需要以下问题得到有效解决
		//A-B交易,A发送RevealSecret以后崩溃,然后很久以后重启
		//如果直接告诉Photon最新块数,那么photon将直接判断该锁过期而发送RemoveExpiredHashLock
		//但是很有可能B已经在链上注册了密码,这个时候A如果发送RemoveExpiredHashLock,将会导致该通道无法使用.
		//因为B会拒绝RemoveExpiredHashLock.为了避免这种情况,一定要在处理最新块之前,处理SerecretRevealOnChain
		for _, sc := range stateChanges {
			if sc.GetBlockNumber() != lastSendBlockNumber {
				be.StateChangeChannel <- &transfer.BlockStateChange{BlockNumber: sc.GetBlockNumber()}
				lastSendBlockNumber = sc.GetBlockNumber()
			}
			be.StateChangeChannel <- sc
		}
		// 先切换有效公链,保证消息处理开始时,
		// 出块时间在3分钟内且大于当前块,被认为是有效最新块,如果当前为无效公链状态,通知上层切换到有效公链状态
		if now-lastedBlockTimestamp < params.Cfg.EffectiveTimeoutSeconds && !be.isChainEffective {
			be.isChainEffective = true
			be.StateChangeChannel <- &transfer.EffectiveChainStateChange{
				IsEffective:              true,
				LastBlockNumber:          lastedBlock,
				LastBlockNumberTimestamp: lastedBlockTimestamp,
			}
		}
		//正常启动流程是,所有历史事件处理完毕,然后再通知photon继续启动
		if lastSendBlockNumber != currentBlock {
			be.StateChangeChannel <- &transfer.BlockStateChange{BlockNumber: currentBlock}
		}
		be.notifyPhotonStartupCompleteIfNeeded(currentBlock)
		// 清除过期流水
		for key, blockNumber := range be.txDone {
			if blockNumber <= uint64(fromBlockNumber) {
				delete(be.txDone, key)
			}
		}
		// wait to next time
		//time.Sleep(be.pollPeriod)
		select {
		case <-time.After(params.Cfg.PollPeriod):
		case <-be.stopChan:
			be.stopChan = nil
			log.Info(fmt.Sprintf("AlarmTask quit complete"))
			return
		}
	}
}

func (be *Events) checkSmcSyncDone() error {
	var err error
	// 2. 校验smc节点同步情况
	sp, err := be.client.SyncProgress(context.Background())
	if err != nil {
		err = rerr.ErrSpectrumSyncError.Errorf("call smc SyncProgress err %s", err)
		return err
	}
	if sp != nil {
		err = rerr.ErrSpectrumBlockError.Errorf("smc is syncing now")
		return err
	}
	return nil
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
	logs, err = rpc.EventsGetInternal(
		rpc.GetQueryConext(), contractAddresses, fromBlock, toBlock, be.client)
	if err != nil {
		return
	}
	return
}

func (be *Events) parseLogsToEvents(logs []types.Log) (stateChanges []mediatedtransfer.ContractStateChange, err error) {
	for _, l := range logs {
		eventName := topicToEventName[l.Topics[0]]
		// 根据已处理流水去重
		if doneBlockNumber, ok := be.txDone[makeEventID(&l)]; ok {
			if doneBlockNumber == l.BlockNumber {
				//log.Trace(fmt.Sprintf("get event txhash=%s repeated,ignore...", l.TxHash.String()))
				continue
			}
			log.Warn(fmt.Sprintf("event tx=%s happened at %d, but now happend at %d ", l.TxHash.String(), doneBlockNumber, l.BlockNumber))
		}
		//chainEventRecordID := be.chainEventRecordDao.MakeChainEventID(&l)
		//// 根据已处理流水去重
		//if doneBlockNumber, delivered := be.chainEventRecordDao.CheckChainEventDelivered(chainEventRecordID); delivered {
		//	if doneBlockNumber == l.BlockNumber {
		//		//log.Trace(fmt.Sprintf("get event txhash=%s repeated,ignore...", l.TxHash.String()))
		//		continue
		//	}
		//	log.Warn(fmt.Sprintf("event tx=%s happened at %d, but now happend at %d ", l.TxHash.String(), doneBlockNumber, l.BlockNumber))
		//}

		// open,deposit,withdraw事件延迟确认,开关默认关闭,方便测试
		if params.Cfg.EnableForkConfirm && needConfirm(eventName) {
			if be.lastBlockNumber-int64(l.BlockNumber) < params.Cfg.ForkConfirmNumber {
				continue
			}
			log.Info(fmt.Sprintf("event %s tx=%s happened at %d, confirmed at %d", eventName, l.TxHash.String(), l.BlockNumber, be.lastBlockNumber))
		}
		// registry secret事件延迟确认,否则在出现恶意分叉的情况下,中间节点有损失资金的风险
		if eventName == params.NameSecretRevealed && params.Cfg.EnableForkConfirm {
			if be.lastBlockNumber-int64(l.BlockNumber) < params.Cfg.ForkConfirmNumber {
				continue
			}
			log.Info(fmt.Sprintf("event %s tx=%s happened at %d, confirmed at %d", eventName, l.TxHash.String(), l.BlockNumber, be.lastBlockNumber))
		}

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
			log.Warn(fmt.Sprintf("receive unkonwn type event from chain : \n%s\n", utils.StringInterface(l, 3)))
		}
		// 记录处理流水
		//be.chainEventRecordDao.NewDeliveredChainEvent(chainEventRecordID, l.BlockNumber)
		be.txDone[makeEventID(&l)] = l.BlockNumber
	}
	return
}

func needConfirm(eventName string) bool {

	if eventName == params.NameChannelOpenedAndDeposit ||
		eventName == params.NameChannelNewDeposit ||
		eventName == params.NameChannelWithdraw {
		return true
	}
	return false
}

//eventChannelSettled2StateChange to stateChange
func eventChannelSettled2StateChange(ev *contracts.TokensNetworkChannelSettled) *mediatedtransfer.ContractSettledStateChange {
	return &mediatedtransfer.ContractSettledStateChange{
		ChannelIdentifier: common.Hash(ev.ChannelIdentifier),
		SettledBlock:      int64(ev.Raw.BlockNumber),
	}
}

//eventChannelCooperativeSettled2StateChange to stateChange
func eventChannelCooperativeSettled2StateChange(ev *contracts.TokensNetworkChannelCooperativeSettled) *mediatedtransfer.ContractCooperativeSettledStateChange {
	return &mediatedtransfer.ContractCooperativeSettledStateChange{
		ChannelIdentifier: common.Hash(ev.ChannelIdentifier),
		SettledBlock:      int64(ev.Raw.BlockNumber),
	}
}

//eventChannelPunished2StateChange to stateChange
func eventChannelPunished2StateChange(ev *contracts.TokensNetworkChannelPunished) *mediatedtransfer.ContractPunishedStateChange {
	return &mediatedtransfer.ContractPunishedStateChange{
		ChannelIdentifier: common.Hash(ev.ChannelIdentifier),
		Beneficiary:       ev.Beneficiary,
		BlockNumber:       int64(ev.Raw.BlockNumber),
	}
}

//eventChannelWithdraw2StateChange to stateChange
func eventChannelWithdraw2StateChange(ev *contracts.TokensNetworkChannelWithdraw) *mediatedtransfer.ContractChannelWithdrawStateChange {
	c := &mediatedtransfer.ContractChannelWithdrawStateChange{
		ChannelIdentifier: &contracts.ChannelUniqueID{

			ChannelIdentifier: common.Hash(ev.ChannelIdentifier),
			OpenBlockNumber:   int64(ev.Raw.BlockNumber),
		},
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
func eventTokenNetworkCreated2StateChange(ev *contracts.TokensNetworkTokenNetworkCreated) *mediatedtransfer.ContractTokenAddedStateChange {
	return &mediatedtransfer.ContractTokenAddedStateChange{
		TokenAddress: ev.TokenAddress,
		BlockNumber:  int64(ev.Raw.BlockNumber),
	}
}

//注意与合约上计算方式保持完全一致.
func calcChannelID(token, tokensNetwork, p1, p2 common.Address) common.Hash {
	var channelID common.Hash
	//log.Trace(fmt.Sprintf("p1=%s,p2=%s,tokennetwork=%s", p1.String(), p2.String(), tokenNetwork.String()))
	if bytes.Compare(p1[:], p2[:]) < 0 {
		channelID = utils.Sha3(p1[:], p2[:], token[:], tokensNetwork[:])
	} else {
		channelID = utils.Sha3(p2[:], p1[:], token[:], tokensNetwork[:])
	}
	return channelID
}

//eventChannelOpenAndDeposit2StateChange to statechange
func eventChannelOpenAndDeposit2StateChange(ev *contracts.TokensNetworkChannelOpenedAndDeposit) (ch1 *mediatedtransfer.ContractNewChannelStateChange, ch2 *mediatedtransfer.ContractBalanceStateChange) {
	ch1 = &mediatedtransfer.ContractNewChannelStateChange{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: calcChannelID(ev.Token, ev.Raw.Address, ev.Participant, ev.Partner),
			OpenBlockNumber:   int64(ev.Raw.BlockNumber),
		},
		Participant1:  ev.Participant,
		Participant2:  ev.Partner,
		SettleTimeout: int(ev.SettleTimeout),
		BlockNumber:   int64(ev.Raw.BlockNumber),
		TokenAddress:  ev.Token,
	}
	ch2 = &mediatedtransfer.ContractBalanceStateChange{
		ChannelIdentifier:  ch1.ChannelIdentifier.ChannelIdentifier,
		ParticipantAddress: ev.Participant,
		BlockNumber:        int64(ev.Raw.BlockNumber),
		Balance:            ev.Participant1Deposit,
	}
	return
}

//eventChannelNewDeposit2StateChange to statechange
func eventChannelNewDeposit2StateChange(ev *contracts.TokensNetworkChannelNewDeposit) *mediatedtransfer.ContractBalanceStateChange {
	return &mediatedtransfer.ContractBalanceStateChange{
		ChannelIdentifier:  ev.ChannelIdentifier,
		ParticipantAddress: ev.Participant,
		BlockNumber:        int64(ev.Raw.BlockNumber),
		Balance:            ev.TotalDeposit,
	}
}

//eventChannelClosed2StateChange to statechange
func eventChannelClosed2StateChange(ev *contracts.TokensNetworkChannelClosed) *mediatedtransfer.ContractClosedStateChange {
	c := &mediatedtransfer.ContractClosedStateChange{
		ChannelIdentifier: ev.ChannelIdentifier,
		ClosingAddress:    ev.ClosingParticipant,
		LocksRoot:         ev.Locksroot,
		ClosedBlock:       int64(ev.Raw.BlockNumber),
		TransferredAmount: ev.TransferredAmount,
	}
	if ev.TransferredAmount == nil {
		c.TransferredAmount = new(big.Int)
	}
	return c
}

//eventBalanceProofUpdated2StateChange to statechange
func eventBalanceProofUpdated2StateChange(ev *contracts.TokensNetworkBalanceProofUpdated) *mediatedtransfer.ContractBalanceProofUpdatedStateChange {
	c := &mediatedtransfer.ContractBalanceProofUpdatedStateChange{
		ChannelIdentifier: ev.ChannelIdentifier,
		LocksRoot:         ev.Locksroot,
		TransferAmount:    ev.TransferredAmount,
		Participant:       ev.Participant,
		BlockNumber:       int64(ev.Raw.BlockNumber),
	}
	if c.TransferAmount == nil {
		c.TransferAmount = new(big.Int)
	}
	return c
}

//eventChannelUnlocked2StateChange to statechange
func eventChannelUnlocked2StateChange(ev *contracts.TokensNetworkChannelUnlocked) *mediatedtransfer.ContractUnlockStateChange {
	c := &mediatedtransfer.ContractUnlockStateChange{
		ChannelIdentifier: ev.ChannelIdentifier,
		BlockNumber:       int64(ev.Raw.BlockNumber),
		TransferAmount:    ev.TransferredAmount,
		Participant:       ev.PayerParticipant,
		LockHash:          ev.Lockhash,
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
