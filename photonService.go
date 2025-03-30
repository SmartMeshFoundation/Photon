package photon

import (
	"context"
	"fmt"

	"time"

	"sync/atomic"

	"math/big"

	"strings"

	"os"

	"runtime/debug"

	"github.com/SmartMeshFoundation/Photon/blockchain"
	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/network/graph"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/fee"
	"github.com/SmartMeshFoundation/Photon/notify"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/SmartMeshFoundation/Photon/pmsproxy"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/theckman/go-flock"
)

/*
message sent complete notification
*/
type protocolMessage struct {
	receiver common.Address
	Message  encoding.Messager
}

// BuildInfo 保存构建信息
type BuildInfo struct {
	GoVersion string `json:"go_version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	Version   string `json:"version"`
}

//SecretRequestPredictor return true to ignore this message,otherwise continue to process
type SecretRequestPredictor func(msg *encoding.SecretRequest) (ignore bool)

//RevealSecretListener return true this listener should not be called next time
type RevealSecretListener func(msg *encoding.RevealSecret) (remove bool)

//ReceivedMediatedTrasnferListener return true this listener should not be called next time
type ReceivedMediatedTrasnferListener func(msg *encoding.MediatedTransfer) (remove bool)

//SentMediatedTransferListener return true this listener should not be called next time
type SentMediatedTransferListener func(msg *encoding.MediatedTransfer) (remove bool)

/*
Service is a photon node
most of Service's member is not thread safe, and should not visit outside the loop method.
*/
type Service struct {
	/*
		module
	*/
	Chain                    *rpc.BlockChainService
	Transport                network.Transporter
	Protocol                 *network.PhotonProtocol
	MessageHandler           *photonMessageHandler
	StateMachineEventHandler *stateMachineEventHandler
	BlockChainEvents         *blockchain.Events
	dao                      models.Dao
	FeePolicy                fee.Charger //Mediation fee
	NotifyHandler            *notify.Handler
	PfsProxy                 pfsproxy.PfsProxy
	PmsProxy                 pmsproxy.PmsProxy

	/*
	 */
	Token2ChannelGraph    map[common.Address]*graph.ChannelGraph
	Token2TokenNetwork    map[common.Address]common.Address
	Transfer2StateManager map[common.Hash]*transfer.StateManager
	Transfer2Result       map[common.Hash]*utils.AsyncResult
	SwapKey2TokenSwap     map[swapKey]*TokenSwap
	/*
		   This is a map from a hashlock to a list of channels, the same
			 hashlock can be used in more than one token (for tokenswaps), a
			 channel should be removed from this list only when the Lock is
			 released/withdrawn but not when the secret is registered.
	*/
	Token2LockSecretHash2Channels map[common.Address]map[common.Hash][]*channel.Channel
	FileLocker                    *flock.Flock
	BlockNumber                   *atomic.Value
	/*
		chan for user request
	*/
	UserReqChan                 chan *apiReq
	ProtocolMessageSendComplete chan *protocolMessage
	/*
		these four maps designed for token swap,but it can be extended for purpose usage.
		for example:
		cross chain.
	*/
	SecretRequestPredictorMap map[common.Hash]SecretRequestPredictor //for tokenswap
	RevealSecretListenerMap   map[common.Hash]RevealSecretListener   //for tokenswap
	/*
		important!:
			we must valid the mediated transfer is valid or not first, then to test  if this mediated transfer matchs any token swap.
	*/
	ReceivedMediatedTrasnferListenerMap   map[*ReceivedMediatedTrasnferListener]bool //for tokenswap
	SentMediatedTransferListenerMap       map[*SentMediatedTransferListener]bool     //for tokenswap
	HealthCheckMap                        map[common.Address]bool
	quitChan                              chan struct{} //for quit notification
	StopCreateNewTransfers                bool          // 是否停止接收新交易,默认false,目前仅在用户调用prepare-update接口的时候,会被置为true,直到重启		// boolean to check whether stop receiving new transfers, default to false. Currently it sets to true when clients invoke prepare-update, till it reconnects.
	EthConnectionStatus                   chan netshare.Status
	ChanHistoryContractEventsDealComplete chan struct{}
	BuildInfo                             *BuildInfo
	ChanSubmitBalanceProofToPFS           chan *channel.Channel // 供submitBalanceProofToPfsLoop线程使用
	ChanSubmitDelegateToPMS               chan *channel.Channel // 供submitDelegateToPmsLoop线程使用

	/*
		取值规则如下:
			1. 启动时默认公链连接无效
			2. 当与公链连接成功且获取到的新块号>本地最新块号时,该标志位应被置为true
			3. 当与公链连接出错时,该标志位应被置为false
			4. 当与公链连接正常,但在3分钟内都没有收到新块时,该标志位被置为false
			5. 当上层应用调用NotifyNetworkDown接口时,该标志位被置为false
	*/
	IsChainEffective         bool  // 当前公链状态是否有效
	EffectiveChangeTimestamp int64 // 公链状态切换时间,即发生状态切换时最后一个有效块的出块时间
}

//NewPhotonService create photon service
func NewPhotonService(chain *rpc.BlockChainService, transport network.Transporter, notifyHandler *notify.Handler, dao models.Dao) (rs *Service, err error) {
	rs = &Service{
		NotifyHandler:      notifyHandler,
		Chain:              chain,
		Transport:          transport,
		dao:                dao,
		Token2ChannelGraph: make(map[common.Address]*graph.ChannelGraph),
		//Token2TokenNetwork 应该是一个token的数组,表示已经注册的token.目前k,v中的v必须是空地址
		Token2TokenNetwork:                    make(map[common.Address]common.Address),
		Transfer2StateManager:                 make(map[common.Hash]*transfer.StateManager),
		Transfer2Result:                       make(map[common.Hash]*utils.AsyncResult),
		Token2LockSecretHash2Channels:         make(map[common.Address]map[common.Hash][]*channel.Channel),
		SwapKey2TokenSwap:                     make(map[swapKey]*TokenSwap),
		UserReqChan:                           make(chan *apiReq, 10),
		BlockNumber:                           new(atomic.Value),
		ProtocolMessageSendComplete:           make(chan *protocolMessage, 10),
		SecretRequestPredictorMap:             make(map[common.Hash]SecretRequestPredictor),
		RevealSecretListenerMap:               make(map[common.Hash]RevealSecretListener),
		ReceivedMediatedTrasnferListenerMap:   make(map[*ReceivedMediatedTrasnferListener]bool),
		SentMediatedTransferListenerMap:       make(map[*SentMediatedTransferListener]bool),
		HealthCheckMap:                        make(map[common.Address]bool),
		quitChan:                              make(chan struct{}),
		StopCreateNewTransfers:                false,
		EthConnectionStatus:                   make(chan netshare.Status, 10),
		ChanHistoryContractEventsDealComplete: make(chan struct{}),
		BuildInfo:                             new(BuildInfo),
		ChanSubmitBalanceProofToPFS:           make(chan *channel.Channel, 100),
		ChanSubmitDelegateToPMS:               make(chan *channel.Channel, 100),
		IsChainEffective:                      false,
	}
	rs.BlockNumber.Store(int64(0))
	rs.MessageHandler = newPhotonMessageHandler(rs)
	rs.StateMachineEventHandler = newStateMachineEventHandler(rs)
	rs.Protocol = network.NewPhotonProtocol(transport, rs)
	////todo fixme MatrixTransport should have a better contructor function
	//mtransport, ok := rs.Transport.(*network.MatrixMixTransport)
	//if ok {
	//	mtransport.SetMatrixDB(rs.dao)
	//}
	rs.Protocol.SetReceivedMessageSaver(NewAckHelper(rs.dao))
	/*
		only one instance for one data directory
	*/
	rs.FileLocker = flock.NewFlock(params.Cfg.DataBasePath + ".flock.Lock")
	locked, err := rs.FileLocker.TryLock()
	if err != nil || !locked {
		err = rerr.ErrPhotonAlreadyRunning.Errorf("another instance already running at %s", params.Cfg.DataBasePath)
		return
	}
	log.Info(fmt.Sprintf("create photon service registry=%s,node=%s", rs.Chain.GetRegistryAddress().String(), params.Cfg.MyAddress.String()))

	rs.Token2TokenNetwork, err = rs.dao.GetAllTokens()
	if err != nil {
		return
	}
	rs.BlockChainEvents = blockchain.NewBlockChainEvents(chain.Client, chain, rs.dao)
	// fee module
	if params.Cfg.EnableMediationFee {
		// pathfinder
		if params.Cfg.PfsHost != "" {
			rs.PfsProxy = pfsproxy.NewPfsProxy(params.Cfg.PfsHost, params.Cfg.PrivateKey)
		}
		rs.FeePolicy = NewFeeModule(dao, rs.PfsProxy)
	} else {
		rs.FeePolicy = &NoFeePolicy{}
	}
	// pms
	if params.Cfg.PmsHost != "" && params.Cfg.PmsAddress != utils.EmptyAddress {
		log.Info(fmt.Sprintf("startup with pms=%s, pms signer=%s", params.Cfg.PmsHost, params.Cfg.PmsAddress.String()))
		rs.PmsProxy = pmsproxy.NewPmsProxy(params.Cfg.PmsHost, params.Cfg.MyAddress, params.Cfg.PmsAddress)
	} else {
		log.Error(fmt.Sprintf("it's unsafe to startup with no pms"))
	}
	return rs, nil
}

// Start the node.
func (rs *Service) Start() (err error) {

	/*
		事先从DB里面获取最后的blocknumber,以免重启后因为超时而拒绝掉之前的MediatedTransfer消息
		Get the last block number from the DB beforehand to avoid rejecting the previous MeditatedTransfer message after restart because of timeout
	*/
	n := rs.dao.GetLatestBlockNumber()
	rs.BlockNumber.Store(n)
	//如果启动的时候就是无网,那么这个时间就完全无效.预置有一个有参考价值的时间
	rs.EffectiveChangeTimestamp = rs.dao.GetLastBlockNumberTime().Unix()
	err = rs.registerRegistry()
	if err != nil {
		return
	}
	//在主循环开启之前,protocol层要准备好,可以发送消息,但是不能接收消息
	rs.Protocol.Start(false)
	//restore 一定要在历史事件处理之前进行,比如链上注册密码事件,需要相应的statemanager发送unlock消息
	rs.restore()
	go func() {
		if params.Cfg.ConditionQuit.RandomQuit {
			go func() {
				isPrime := func(value int) bool {
					if value <= 3 {
						return value >= 2
					}
					if value%2 == 0 || value%3 == 0 {
						return false
					}
					for i := 5; i*i <= value; i += 6 {
						if value%i == 0 || value%(i+2) == 0 {
							return false
						}
					}
					return true
				}
				for {
					/*
							Random sleep is no more than five seconds. If the number of dormancy milliseconds is prime,
						    it will exit directly. This probability is probably 13%.
					*/
					n := utils.NewRandomInt(5000)
					time.Sleep(time.Duration(n) * time.Millisecond)
					if isPrime(n) {
						panic("random quit")
					}
				}
			}()
		}
		rs.loop()
	}()

	// 这里如果状态为connected,则等待积压的block events处理完毕后再启动api以及订阅其他节点的消息
	// 如果状态不为connected,则直接启动api以及订阅其他节点的消息,这样做可能带来的风险:
	// 1. 积压事件处理完毕之前,用户/其他节点通过api/消息对本地数据作出修改,是否会给后续的链上事件同步工作带来问题???
	// Here if status is connected, then after block events completes, we should restart api and subscribe messages from other nodes.
	if rs.Chain.Client.Status == netshare.Connected {
		//wait for start up complete.
		<-rs.ChanHistoryContractEventsDealComplete
		log.Info(fmt.Sprintf("Photon Startup complete and history events process complete."))
	} else {
		log.Info(fmt.Sprintf("Photon Startup complete without effective chain"))
	}

	/*
		将protocol接受消息移到历史事件处理之后,
		保证不在历史事件处理完毕之前进入事件主循环.
		发现问题原因:
		在测试中发现,其他节点启动完毕以后尝试发送消息,如果这时候事件主循环已经启动,那么就会收到消息进行处理.这时候通道状态可能是错的,不匹配的.
		虽然实际中可能会尝试重发,但是会让测试代码进入一种不确定的等待.
		这么做有可能因为接收到过多的消息,而阻塞接受线程,导致消息丢失.但是因为没有处理,对方一定会反复重新发送.
	*/
	rs.Protocol.StartReceive()
	/*
		启动定时提交balance_proof到pfs及pms的线程
	*/
	go rs.submitBalanceProofToPfsLoop()
	go rs.submitDelegateToPmsLoop()
	rs.startNeighboursHealthCheck()
	return nil
}

//Stop the node.
func (rs *Service) Stop() {
	log.Info("photon service stop...")
	close(rs.quitChan)
	rs.Protocol.StopAndWait()
	rs.BlockChainEvents.Stop()
	rs.Chain.Stop()
	rs.NotifyHandler.Stop()
	time.Sleep(100 * time.Millisecond) // let other goroutines quit
	rs.dao.CloseDB()
	//anther instance cann run now
	err := rs.FileLocker.Unlock()
	if err != nil {
		log.Error(fmt.Sprintf("Unlock err %s", err))
	}
	log.Info("photon service stop ok...")
}

/*
main loop of this photon nodes
process  events below:
1. request from user
2. event from blockchain
3. message from other nodes.
*/
func (rs *Service) loop() {
	var err error
	var ok bool
	var m *network.MessageToPhoton
	var st transfer.StateChange
	var req *apiReq
	var sentMessage *protocolMessage

	defer rpanic.PanicRecover("photon service")
	for {
		select {
		//message from other nodes
		case m, ok = <-rs.Protocol.ReceivedMessageChan:
			if ok {
				err = rs.MessageHandler.onMessage(m.Msg, m.EchoHash)
				if err != nil {
					log.Error(fmt.Sprintf("MessageHandler.onMessage %v", err))
				}
				rs.Protocol.ReceivedMessageResultChan <- err
			} else {
				log.Info("Protocol.ReceivedMessageChan closed")
				return
			}
			// contract events from block chain
		case st, ok = <-rs.BlockChainEvents.StateChangeChannel:
			if ok {
				blockStateChange, ok2 := st.(*transfer.BlockStateChange)
				if ok2 {
					rs.handleBlockNumber(blockStateChange)
				} else {
					log.Trace(fmt.Sprintf("statechange received :%s", utils.StringInterface(st, 2)))
					_, isHistoryComplete := st.(*mediatedtransfer.ContractHistoryEventCompleteStateChange)
					if isHistoryComplete {
						if rs.ChanHistoryContractEventsDealComplete != nil {
							close(rs.ChanHistoryContractEventsDealComplete)
							rs.ChanHistoryContractEventsDealComplete = nil
						} else {
							panic("only can receive ContractHistoryEventCompleteStateChange once")
						}
					} else {
						err = rs.StateMachineEventHandler.OnBlockchainStateChange(st)
						if err != nil {
							log.Error(fmt.Sprintf("stateMachineEventHandler.OnBlockchainStateChange %s", err))
						}

					}
				}

			} else {
				log.Info("Events.StateChangeChannel closed")
				return
			}
		//user's request
		case req, ok = <-rs.UserReqChan:
			if ok {
				rs.handleReq(req)
			} else {
				log.Info("req closed")
				return
			}
			//i have sent a message complete
		case sentMessage, ok = <-rs.ProtocolMessageSendComplete:
			if ok {
				rs.handleSentMessage(sentMessage)
			} else {
				log.Info("ProtocolMessageSendComplete closed")
				return
			}
		case s := <-rs.Chain.Client.StatusChan:
			if s == netshare.Connected {
				rs.handleEthRPCConnectionOK()
			} else {
				//通过EthConnectionStatus通知eth连接断开,
				//rs.NotifyHandler.NotifyString(notify.LevelWarn, "公链连接失败,正在尝试重连")
			}
		case <-rs.quitChan:
			log.Info(fmt.Sprintf("%s quit now", utils.APex2(params.Cfg.MyAddress)))
			return
		}
	}
}

//for init,read dao history,只要是我还没处理的链上事件,都还在队列中等着发给我.
// for init, read dao history,
// all on-chain events I have not handled should wait in queue.
func (rs *Service) registerRegistry() (err error) {
	token2TokenNetworks, err := rs.dao.GetAllTokens()
	if err != nil {
		return
	}
	for token := range token2TokenNetworks {
		err = rs.registerTokenNetwork(token)
		if err != nil {
			//err = fmt.Errorf("registerTokenNetwork err:%s", err)
			return
		}
	}
	return nil
}

/*
quering my channel details on blockchain
i'm one of the channel participants
收到来自链上的事件,新创建了 channel
但是事件有可能重复
*/
/*
 *	newChannelFromEvent : function to handle channel query event.
 *
 *	Note that this node is also one of the channle participant, and he receives messages on-chain, to create a new channel.
 *	But those events could be repeated.
 */
func (rs *Service) newChannelFromEvent(tokenNetwork *rpc.TokenNetworkProxy, tokenAddress common.Address, partnerAddress common.Address, channelIdentifier *contracts.ChannelUniqueID, settleTimeout int) (ch *channel.Channel, err error) {
	/*
		因为有可能在我离线的时候收到一堆事件,所以通道的信息不一定就是新创建时候的状态,
		但是保证后续的事件会继续收到,所以应该按照新通道处理.
	*/
	// Because it is possible that I receive a bunch of events when disconnected, so channel states may not be the same as those when first created.
	// But we ensure that following events will be got by me 100%, so we should handle them as creating a new channel.
	ourState := channel.NewChannelEndState(params.Cfg.MyAddress, big.NewInt(0), nil, mtree.NewMerkleTree(nil))
	partenerState := channel.NewChannelEndState(partnerAddress, big.NewInt(0), nil, mtree.NewMerkleTree(nil))

	externState := channel.NewChannelExternalState(rs.registerChannelForHashlock, tokenNetwork, channelIdentifier, params.Cfg.PrivateKey, rs.Chain.Client, rs.dao, 0, params.Cfg.MyAddress, partnerAddress)
	ch, err = channel.NewChannel(ourState, partenerState, externState, tokenAddress, channelIdentifier, params.Cfg.RevealTimeout, settleTimeout)
	return
}

/*
block chain tick,
it's the core of HTLC
*/
func (rs *Service) handleBlockNumber(st *transfer.BlockStateChange) {
	rs.BlockNumber.Store(st.BlockNumber)
	rs.StateMachineEventHandler.dispatchToAllTasks(st)
	for _, cg := range rs.Token2ChannelGraph {
		for _, c := range cg.ChannelIdentifier2Channel {
			err := rs.StateMachineEventHandler.ChannelStateTransition(c, st)
			if err != nil {
				log.Error(fmt.Sprintf("ChannelStateTransition err %s", err))
			}
		}
	}
	rs.dao.SaveLatestBlockNumber(st.BlockNumber)
	return
}

//GetBlockNumber return latest blocknumber of ethereum
func (rs *Service) GetBlockNumber() int64 {
	return rs.BlockNumber.Load().(int64)
}

// GetChannelStatus return status of channel
func (rs *Service) GetChannelStatus(channelIdentifier common.Hash) (int, int64) {
	c := rs.getChannelWithAddr(channelIdentifier)
	if c == nil {
		return channeltype.StateInValid, 0
	}
	return int(c.State), c.ChannelIdentifier.OpenBlockNumber
}

func (rs *Service) findChannelByIdentifier(channelIdentifier common.Hash) (*channel.Channel, error) {
	for _, g := range rs.Token2ChannelGraph {
		ch := g.ChannelIdentifier2Channel[channelIdentifier]
		if ch != nil {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("unknown channel %s", utils.HPex(channelIdentifier))
}

/*
Send `message` to `recipient` using the photon protocol.

       The protocol will take care of resending the message on a given
       interval until an Acknowledgment is received or a given number of
       tries.
*/
func (rs *Service) sendAsync(recipient common.Address, msg encoding.SignedMessager) error {
	if recipient == params.Cfg.MyAddress {
		log.Error(fmt.Sprintf("rs must be a bug ,sending message to it self"))
	}
	mtr, ok := msg.(*encoding.MediatedTransfer)
	if ok && mtr != nil {
		for f := range rs.SentMediatedTransferListenerMap {
			remove := (*f)(mtr)
			if remove {
				delete(rs.SentMediatedTransferListenerMap, f)
			}
		}
	}
	envelopMessager, ok := msg.(encoding.EnvelopMessager)
	if ok && envelopMessager != nil {
		rs.dao.NewSentEnvelopMessager(envelopMessager, recipient)
	}
	result := rs.Protocol.SendAsync(recipient, msg)
	//todo 发送消息量大的时候,会制造大量的goroutine,比较昂贵
	go func() {
		defer rpanic.PanicRecover(fmt.Sprintf("send %s, msg:%s", utils.APex(recipient), msg))
		select {
		case <-rs.quitChan:
			return
		case err := <-result.Result: //如果通道已经settle,那么这个消息是没必要再发送了.这时候会失败
			if err == nil {
				rs.ProtocolMessageSendComplete <- &protocolMessage{
					receiver: recipient,
					Message:  msg,
				}
			} else {
				log.Error(fmt.Sprintf("message %s send finished ,but err=%s", utils.StringInterface(msg, 3), err))
			}
		}
	}()
	return nil
}

/*
SendAndWait Send `message` to `recipient` and wait for the response or `timeout`.

       Args:
           recipient (address): The address of the node that will receive the
               message.
           message: The transfer message.
           timeout (float): How long should we wait for a response from `recipient`.

       Returns:
           None: If the wait timed out
           object: The result from the event
*/
func (rs *Service) SendAndWait(recipient common.Address, message encoding.SignedMessager, timeout time.Duration) error {
	return rs.Protocol.SendAndWait(recipient, message, timeout)
}

/*
Register the secret with any channel that has a hashlock on it.

       This must search through all channels registered for a given hashlock
       and ignoring the tokens.
*/
func (rs *Service) registerSecret(secret common.Hash) {
	hashlock := utils.ShaSecret(secret[:])
	for _, hashchannel := range rs.Token2LockSecretHash2Channels {
		for _, ch := range hashchannel[hashlock] {
			err := ch.RegisterSecret(secret)
			err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
			if err != nil {
				log.Error(fmt.Sprintf("RegisterSecret %s to channel %s  err: %s",
					utils.HPex(secret), ch.ChannelIdentifier.String(), err))
			}
		}
	}
}

/*
链上这个锁对应的密码注册了,
*/
// The secret of this lock has been registered on-chain.
func (rs *Service) registerRevealedLockSecretHash(lockSecretHash, secret common.Hash, blockNumber int64) {
	for _, hashchannel := range rs.Token2LockSecretHash2Channels {
		for _, ch := range hashchannel[lockSecretHash] {
			err := ch.RegisterRevealedSecretHash(lockSecretHash, secret, blockNumber)
			if err != nil {
				log.Error(fmt.Sprintf("RegisterRevealedSecretHash to channel err,locksecrethash=%s,secret=%s,err=%s,ch=%s",
					utils.HPex(lockSecretHash), utils.HPex(secret), err, ch,
				))
				continue
			}
			err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
			if err != nil {
				log.Error(fmt.Sprintf("RegisterSecret %s to channel %s  err: %s",
					utils.HPex(lockSecretHash), ch.ChannelIdentifier.String(), err))
			}
		}
	}
}
func (rs *Service) registerChannelForHashlock(netchannel *channel.Channel, lockSecretHash common.Hash) {
	tokenAddress := netchannel.TokenAddress
	channelsRegistered := rs.Token2LockSecretHash2Channels[tokenAddress][lockSecretHash]
	found := false
	for _, c := range channelsRegistered {
		//To determine whether the two channel objects are equal, we simply use the address to identify.
		if c.ExternState.ChannelIdentifier == netchannel.ExternState.ChannelIdentifier {
			found = true
			break
		}
	}
	if !found {
		hashLock2Channels, ok := rs.Token2LockSecretHash2Channels[tokenAddress]
		if !ok {
			hashLock2Channels = make(map[common.Hash][]*channel.Channel)
			rs.Token2LockSecretHash2Channels[tokenAddress] = hashLock2Channels
		}
		channelsRegistered = append(channelsRegistered, netchannel)
		rs.Token2LockSecretHash2Channels[tokenAddress][lockSecretHash] = channelsRegistered
	}
}
func (rs *Service) channelSerilization2Channel(c *channeltype.Serialization, tokenNetwork *rpc.TokenNetworkProxy) (ch *channel.Channel, err error) {
	OurState := channel.NewChannelEndState(c.OurAddress, c.OurContractBalance,
		c.OurBalanceProof, mtree.NewMerkleTree(c.OurLeaves))
	PartnerState := channel.NewChannelEndState(c.PartnerAddress(),
		c.PartnerContractBalance,
		c.PartnerBalanceProof, mtree.NewMerkleTree(c.PartnerLeaves))
	ExternState := channel.NewChannelExternalState(rs.registerChannelForHashlock, tokenNetwork,
		c.ChannelIdentifier, params.Cfg.PrivateKey,
		rs.Chain.Client, rs.dao, c.ClosedBlock,
		c.OurAddress, c.PartnerAddress())
	ch, err = channel.NewChannel(OurState, PartnerState, ExternState, c.TokenAddress(), c.ChannelIdentifier, c.RevealTimeout, c.SettleTimeout)
	if err != nil {
		return
	}

	ch.OurState.Lock2PendingLocks = c.OurLock2PendingLocks()
	ch.OurState.Lock2UnclaimedLocks = c.OurLock2UnclaimedLocks()
	ch.PartnerState.Lock2PendingLocks = c.PartnerLock2PendingLocks()
	ch.PartnerState.Lock2UnclaimedLocks = c.PartnerLock2UnclaimedLocks()
	ch.State = c.State
	ch.OurState.ContractBalance = c.OurContractBalance
	ch.PartnerState.ContractBalance = c.PartnerContractBalance
	ch.ExternState.ClosedBlock = c.ClosedBlock
	ch.ExternState.SettledBlock = c.SettledBlock
	return
}

//read a token network info from dao
func (rs *Service) registerTokenNetwork(tokenAddress common.Address) (err error) {
	log.Trace(fmt.Sprintf("registerTokenNetwork tokenaddress=%s ", tokenAddress.String()))
	var tokenNetwork *rpc.TokenNetworkProxy
	tokenNetwork, err = rs.Chain.TokenNetwork(tokenAddress)
	if err != nil {
		return
	}
	edges, err := rs.dao.GetAllNonParticipantChannelByToken(tokenAddress)
	if err != nil {
		return
	}
	g := graph.NewChannelGraph(params.Cfg.MyAddress, tokenAddress, edges)
	rs.Token2TokenNetwork[tokenAddress] = utils.EmptyAddress
	rs.Token2ChannelGraph[tokenAddress] = g
	//add channel I participant
	var css []*channeltype.Serialization
	css, err = rs.dao.GetChannelList(tokenAddress, utils.EmptyAddress)
	if err != nil {
		return
	}

	for _, cs := range css {
		//跳过已经 settle 的 channel 加入没有任何意义.
		if cs.State == channeltype.StateSettled {
			continue
		}
		ch, err := rs.channelSerilization2Channel(cs, tokenNetwork)
		if err != nil {
			return err
		}
		err = g.AddChannel(ch)
		if err != nil {
			return err
		}
	}
	return
}

/*
found new channel on blockchain when running...
*/
func (rs *Service) registerChannel(tokenAddress common.Address, partnerAddress common.Address, channelIdentifier *contracts.ChannelUniqueID, settleTimeout int) {
	tokenNetwork, err := rs.Chain.TokenNetwork(tokenAddress)
	if err != nil {
		log.Error(fmt.Sprintf("receive new channel %s-%s,but cannot create tokennetwork err %s",
			utils.APex2(tokenAddress), utils.APex2(partnerAddress), err,
		))
		return
	}
	if rs.getChannel(tokenAddress, partnerAddress) != nil {
		log.Error(fmt.Sprintf("receive new channel %s-%s,but this channel already exist, maybe a duplicate channel event", utils.APex2(tokenAddress), utils.APex2(partnerAddress)))
		return
	}
	ch, err := rs.newChannelFromEvent(tokenNetwork, tokenAddress, partnerAddress, channelIdentifier, settleTimeout)
	if err != nil {
		log.Error(fmt.Sprintf("newChannelFromEvent err %s", err))
		return
	}
	g := rs.getToken2ChannelGraph(tokenAddress)
	err = g.AddChannel(ch)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = rs.dao.NewChannel(channel.NewChannelSerialization(g.ChannelIdentifier2Channel[ch.ChannelIdentifier.ChannelIdentifier]))
	if err != nil {
		log.Error(err.Error())
		return
	}
	//if true {
	//	g := rs.getChannelGraph(ch.ChannelIdentifier.ChannelIdentifier)
	//	log.Trace(fmt.Sprintf("receive new channel g=%s", utils.StringInterface(g, 3)))
	//}
	return
}

/*
Do a direct tranfer with target.

       Direct transfers are non cancellable and non expirable, since these
       transfers are a signed balance proof with the transferred amount
       incremented.

       Because the transfer is non cancellable, there is a level of trust with
       the target. After the message is sent the target is effectively paid
       and then it is not possible to revert.

       The async result will be set to False iff there is no direct channel
       with the target or the payer does not have balance to complete the
       transfer, otherwise because the transfer is non expirable the async
       result *will never be set to False* and if the message is sent it will
       hang until the target node acknowledge the message.

       This transfer should be used as an optimization, since only two packets
       are required to complete the transfer (from the payer's perspective),
       whereas the mediated transfer requires 6 messages.
*/
func (rs *Service) directTransferAsync(tokenAddress, target common.Address, amount *big.Int, data string) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	g := rs.getToken2ChannelGraph(tokenAddress)
	if g == nil {
		result.Result <- rerr.ErrTokenNotFound
		return
	}
	directChannel := g.GetPartenerAddress2Channel(target)
	if directChannel == nil || !directChannel.CanTransfer() {
		result.Result <- rerr.ErrChannelNotFound.Append("no available direct channel")
		return
	}
	if !rs.IsChainEffective && time.Now().Unix()-rs.EffectiveChangeTimestamp >= directChannel.GetHalfSettleTimeoutSeconds() {
		result.Result <- rerr.ErrNotAllowDirectTransfer
		return
	}
	/*
		发之前检测一下,接收方是否在线,如果不在线也不用发了.避免我方发出去,对方收不到这种情况.
	*/
	if rs.Chain.Client != nil && rs.Chain.Client.Status != netshare.Connected {
		//无网情况下的交易,主要考虑移动端,他所连接点公链节点一定是互联网上的而不是局域网的
		mixTransport, ok := rs.Protocol.Transport.(network.MixTranspoter)
		if ok {
			_, isOnline := mixTransport.UDPNodeStatus(target)
			if !isOnline {
				result.Result <- rerr.ErrNodeNotOnline
				return
			}
		} else {
			log.Error(fmt.Sprintf("it'a error when not in test"))
		}

	} else {
		_, isOnline := rs.Protocol.Transport.NodeStatus(target)
		if !isOnline {
			result.Result <- rerr.ErrNodeNotOnline
			return
		}
	}
	tr, err := directChannel.CreateDirectTransfer(amount)
	if err != nil {
		result.Result <- err
		return
	}
	tr.Data = []byte(data)
	err = tr.Sign(params.Cfg.PrivateKey, tr)
	err = directChannel.RegisterTransfer(rs.GetBlockNumber(), tr)
	if err != nil {
		result.Result <- err
		return
	}
	//This should be set once the direct transfer is acknowledged
	transferSuccess := &transfer.EventTransferSentSuccess{
		LockSecretHash:    utils.EmptyHash,
		Amount:            amount,
		Target:            target,
		ChannelIdentifier: directChannel.ChannelIdentifier.ChannelIdentifier,
		Token:             tokenAddress,
		Data:              data,
	}
	/*
		对于DirectTransfer,生成一个假的LockSecretHash,
		用于发起方在这里记录发起的交易状态,后续UpdateTransferStatus会更新DB中的值
	*/
	tr.FakeLockSecretHash = utils.NewRandomHash()
	log.Trace(fmt.Sprintf("send direct transfer, use fake lockSecertHash %s to trace transfer status", tr.FakeLockSecretHash.String()))
	// 构造SentTransferDetail
	rs.dao.NewSentTransferDetail(tokenAddress, target, amount, data, true, tr.FakeLockSecretHash)
	//rs.dao.NewTransferStatus(tokenAddress, tr.FakeLockSecretHash)
	err = rs.sendAsync(directChannel.PartnerState.Address, tr)
	if err != nil {
		result.Result <- err
		return
	}
	rs.dao.UpdateSentTransferDetailStatusMessage(tokenAddress, tr.FakeLockSecretHash, "DirectTransfer sending")
	//rs.dao.UpdateTransferStatusMessage(tokenAddress, tr.FakeLockSecretHash, "DirectTransfer 正在发送")
	/*
		Transfer is success
		whenever partner receive this transfer  or not
	*/
	err = rs.StateMachineEventHandler.OnEvent(transferSuccess, nil)
	if err != nil {
		log.Error(fmt.Sprintf("dispatch transferSuccess err %s", err))
	}
	result.LockSecretHash = tr.FakeLockSecretHash
	//result.Result <- err
	smkey := utils.Sha3(result.LockSecretHash[:], tokenAddress[:])
	rs.Transfer2Result[smkey] = result
	return
}

/*
lauch a new mediated trasfer
Args:
 hashlock: caller can specify a hashlock or use empty ,when empty, will generate a random secret.
 expiration: caller can specify a valid blocknumber or 0, when 0 ,will calculate based on settle timeout of channel.

Calls:
	//1. mediatedTransfer
	//	 1.1 带lockSecretHash和Secret,则为指定密码的交易
	//	 1.2 不带lockSecretHash和Secret则为普通交易,需生成随机密码
	//2. token swap taker 带lockSecretHash,不带secret
	//3. token swap maker 带lockSecretHash,带secret
*/
/*
 *	launch a new mediated transfer
 *
 *	Args :
 *		hashlock : caller can specify a hashlock or use empty, when empty, will generate a random secret.
 *		expiration : caller can specify a valid blocknumber or 0, when 0, will calculate based on settle timeout of channel.
 *	Calls :
 *		1. mediatedTransfer
 *			1.1 if it has lockSecretHash and Secret, then it is a transfer with specific secret.
 *			1.2 If it has no lockSecretHash or Secret, then it is a normal transfer, and we need generate a random secret for it.
 *		2. token swap
 *			2.1 taker should contain lockSecretHash, but no secret.
 *			2.2 maker should contain lockSecretHash and secret.
 */
func (rs *Service) startMediatedTransferInternal(tokenAddress, target common.Address, amount *big.Int, lockSecretHash common.Hash, expiration int64, secret common.Hash, data string, routeInfo []pfsproxy.FindPathResponse) (result *utils.AsyncResult, stateManager *transfer.StateManager) {
	var availableRoutes []*route.State
	//var err error
	//targetAmount := new(big.Int).Sub(amount, fee)
	result = utils.NewAsyncResult()
	g := rs.getToken2ChannelGraph(tokenAddress)
	if g == nil {
		result.Result <- rerr.ErrTokenNotFound
		return
	}
	// 2019-03消息升级过后,如果参数没有RouteInfo,仅支持与target直接拥有通道的情况下发送交易或是在不收费的网络下使用本地路由
	if routeInfo == nil || len(routeInfo) == 0 {
		// 当前为不支持收费的网络下时,使用本地路由
		if rs.PfsProxy == nil {
			log.Trace("get available routes without fee from local channel graph")
			availableRoutes = g.GetBestRoutes(rs.Protocol, params.Cfg.MyAddress, target, amount, amount, graph.EmptyExlude, rs)
		} else {
			log.Trace("get available routes to partner from local channel graph")
			ch := rs.getChannel(tokenAddress, target)
			if ch != nil {
				r := route.NewState(ch, []common.Address{ch.PartnerState.Address})
				r.TotalFee = utils.BigInt0
				availableRoutes = append(availableRoutes, r)
			}
		}
	} else {
		// 用户指定了路由的话,采用用户指定的路由,否则从pfs或者本地查询路由
		log.Trace("get available routes from user req")
		for _, path := range routeInfo {
			if path.Result == nil || len(path.Result) == 0 {
				continue
			}
			partnerAddress := common.HexToAddress(path.Result[0])
			ch := rs.getChannel(tokenAddress, partnerAddress)
			if ch == nil {
				continue
			}
			r := route.NewState(ch, path.GetPath())
			//r.Fee = rs.FeePolicy.GetNodeChargeFee(partnerAddress, tokenAddress, amount) // 发起方不收取手续费
			r.TotalFee = path.Fee
			availableRoutes = append(availableRoutes, r)
		}
	}
	log.Trace(fmt.Sprintf("availableRoutes=%s", utils.StringInterface(availableRoutes, 3)))
	if len(availableRoutes) <= 0 {
		result.Result <- rerr.ErrNoAvailabeRoute
		return
	}
	// 当没有有效公链的时候,不支持发送MediatedTransfer,否则有安全隐患
	if !rs.IsChainEffective {
		result.Result <- rerr.ErrNotAllowMediatedTransfer
		return
	}
	/*
		when user specified fee, for test or other purpose.
	*/
	// 2019-03消息升级过后,手续费以RouteInfo参数中的为准
	//if fee.Cmp(utils.BigInt0) > 0 {
	//	for _, r := range availableRoutes {
	//		r.TotalFee = fee //use the user's fee to replace algorithm's
	//	}
	//}
	routesState := route.NewRoutesState(availableRoutes)
	transferState := &mediatedtransfer.LockedTransferState{
		TargetAmount:   new(big.Int).Set(amount),
		Amount:         new(big.Int).Set(amount),
		Token:          tokenAddress,
		Initiator:      params.Cfg.MyAddress,
		Target:         target,
		Expiration:     expiration,
		LockSecretHash: lockSecretHash,
		Secret:         secret,
		Fee:            utils.BigInt0,
		Data:           data,
	}
	/*
		发起方每次切换路径不再切换密码,不切换依然可以保证安全
	*/
	// Initiator has no need to switch secret, every time he switches the route, and security can be ensured.
	initInitiator := &mediatedtransfer.ActionInitInitiatorStateChange{
		OurAddress:     params.Cfg.MyAddress,
		Tranfer:        transferState,
		Routes:         routesState,
		BlockNumber:    rs.GetBlockNumber(),
		Secret:         secret,
		LockSecretHash: lockSecretHash,
		Db:             rs.dao,
	}
	//log.Trace(fmt.Sprintf("start mediated transfer availableRoutes=%s", utils.StringInterface(availableRoutes, 2)))
	stateManager = transfer.NewStateManager(initiator.StateTransition, nil, initiator.NameInitiatorTransition, lockSecretHash, transferState.Token)
	smkey := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := rs.Transfer2StateManager[smkey]
	if manager != nil {
		result.Result <- rerr.ErrDuplicateTransfer
		return
	}
	rs.Transfer2StateManager[smkey] = stateManager
	rs.Transfer2Result[smkey] = result
	//rs.dao.AddStateManager(stateManager)
	rs.StateMachineEventHandler.dispatch(stateManager, initInitiator)
	return
}

/*
1. user start a mediated transfer
2. user start a mediated transfer with secret
*/
func (rs *Service) startMediatedTransfer(tokenAddress, target common.Address, amount *big.Int, secret common.Hash, data string, routeInfo []pfsproxy.FindPathResponse) (result *utils.AsyncResult) {
	lockSecretHash := utils.EmptyHash
	if secret != utils.EmptyHash {
		lockSecretHash = utils.ShaSecret(secret.Bytes())
		/*用户使用指定的密码来进行交易,那么:
		1. 注册SecretRequestPredictor,防止在用户允许之前发送密码出去
		2. 保证用户在提供密码之后,能移除掉这个predictor
		*/
		/*
		 *	Participants use specified secret to send a transfer, then
		 *		1. Register SecretRequestPredictor, preventing secret is sent out before participants permit.
		 *		2. Ensure that this predictor can be removed once participants provide secret.
		 */
		var secretRequestHook SecretRequestPredictor = func(msg *encoding.SecretRequest) (ignore bool) {
			return true
		}
		rs.SecretRequestPredictorMap[lockSecretHash] = secretRequestHook
		log.Trace(fmt.Sprintf("Register SecretRequestPredictor for secret=[%s] lockSecretHash=[%s]\n", secret.String(), lockSecretHash.String()))
	} else {
		/*
			普通交易，随机生成密码
		*/
		// Normal transfer, generate random secret.
		secret = utils.NewRandomHash()
		lockSecretHash = utils.ShaSecret(secret[:])
	}
	/*
		发起方在这里记录发起的交易状态,后续UpdateTransferStatus会更新DB中的值
	*/
	rs.dao.NewSentTransferDetail(tokenAddress, target, amount, data, false, lockSecretHash)
	//rs.dao.NewTransferStatus(tokenAddress, lockSecretHash)
	result, _ = rs.startMediatedTransferInternal(tokenAddress, target, amount, lockSecretHash, 0, secret, data, routeInfo)
	result.LockSecretHash = lockSecretHash
	return
}

//receive a MediatedTransfer, i'm a hop node
func (rs *Service) mediateMediatedTransfer(msg *encoding.MediatedTransfer, ch *channel.Channel) {
	tokenAddress := ch.TokenAddress
	smkey := utils.Sha3(msg.LockSecretHash[:], tokenAddress[:])
	stateManager := rs.Transfer2StateManager[smkey]
	/*
			第一次收到这个密码,
		首先要判断这个密码是否是我声明放弃过的,如果是,就应该谨慎处理.
			锁是有可能重复的,比如 token swap 中.
	*/
	/*
	 *	First time receiving this secret.
	 *	We need to check if I have ever abandoned this secret, if so, handle it carefully.
	 *	Locks can be duplicated, like in token swap.
	 */
	if rs.dao.IsLockSecretHashChannelIdentifierDisposed(msg.LockSecretHash, ch.ChannelIdentifier.ChannelIdentifier) {
		log.Error(fmt.Sprintf("receive a lock secret hash,and it's my annouce disposed. %s", msg.LockSecretHash.String()))
		//忽略,什么都不做
		// do nothing.
		return
	}
	var avaiableRoutes []*route.State
	amount := msg.PaymentAmount
	//targetAddr := msg.Target
	fromChannel := ch
	fromRoute := graph.Channel2RouteState(fromChannel, msg.Sender, amount, rs, msg.Path)
	fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg, ch.TokenAddress)
	if stateManager != nil {
		if stateManager.Name != mediator.NameMediatorTransition {
			log.Error(fmt.Sprintf("receive mediator transfer,but i'm not a mediator,msg=%s,stateManager=%s", msg, utils.StringInterface(stateManager, 3)))
			return
		}
		// 2019-03 消息升级后,仅在不收费的情况下支持重复交易
		if rs.PfsProxy != nil {
			log.Error(fmt.Sprintf("receive repeate mediator transfer,but i'm not a disable-fee node ,msg=%s,stateManager=%s", msg, utils.StringInterface(stateManager, 3)))
			return
		}
		stateChange := &mediatedtransfer.MediatorReReceiveStateChange{
			Message:      msg,
			FromTransfer: fromTransfer,
			FromRoute:    fromRoute,
			BlockNumber:  rs.GetBlockNumber(),
		}
		rs.StateMachineEventHandler.dispatch(stateManager, stateChange)
	} else {
		// 2019-03 消息升级后,路由以mtr中带有的path为准,有且只有一条,如果在不支持手续费的网络中,则根据本地路由继续交易
		if len(msg.Path) == 0 {
			if rs.PfsProxy != nil {
				log.Error("receive MediatedTransfer without route info,ignore")
				return
			}
			exclude := graph.MakeExclude(msg.Sender, msg.Initiator)
			g := rs.getToken2ChannelGraph(ch.TokenAddress) //must exist
			avaiableRoutes = g.GetBestRoutes(rs.Protocol, params.Cfg.MyAddress, msg.Target, amount, msg.PaymentAmount, exclude, rs)
		} else {
			// 获取下一跳的通道
			myIndexInPath := -1
			for i, addr := range msg.Path {
				if addr == params.Cfg.MyAddress {
					myIndexInPath = i
					break
				}
			}
			if myIndexInPath == -1 {
				log.Error("can not found myself in msg.Path")
				return
			}
			//传递参数有问题,导致没有下一跳
			if myIndexInPath+1 >= len(msg.Path) {
				log.Error(fmt.Sprintf("i'm not target,but cannot find more hop node,msg=%s", utils.StringInterface(msg, 5)))
				return
			}
			nextChan := rs.getChannel(ch.TokenAddress, msg.Path[myIndexInPath+1])
			if nextChan == nil {
				log.Error(fmt.Sprintf("receive path,but channel between me and %s doesn't exist", msg.Path[myIndexInPath+1].String()))
				return
			}
			// 构造路由,手续费根据TargetAmount在下家通道中的费率计算
			availableRoute := route.NewState(nextChan, msg.Path)
			targetAmount := new(big.Int).Sub(msg.PaymentAmount, msg.Fee)
			availableRoute.Fee = rs.FeePolicy.GetNodeChargeFee(nextChan.PartnerState.Address, nextChan.TokenAddress, targetAmount)
			avaiableRoutes = append(avaiableRoutes, availableRoute)
		}
		routesState := route.NewRoutesState(avaiableRoutes)
		blockNumber := rs.GetBlockNumber()
		initMediator := &mediatedtransfer.ActionInitMediatorStateChange{
			OurAddress:               params.Cfg.MyAddress,
			FromTranfer:              fromTransfer,
			Routes:                   routesState,
			FromRoute:                fromRoute,
			BlockNumber:              blockNumber,
			Message:                  msg,
			Db:                       rs.dao,
			IsEffectiveChain:         rs.IsChainEffective,
			EffectiveChangeTimestamp: rs.EffectiveChangeTimestamp,
		}
		stateManager = transfer.NewStateManager(mediator.StateTransition, nil, mediator.NameMediatorTransition, fromTransfer.LockSecretHash, fromTransfer.Token)
		//rs.dao.AddStateManager(stateManager)
		rs.Transfer2StateManager[smkey] = stateManager //for path A-B-C-F-B-D-E ,node B will have two StateManagers for one identifier
		rs.StateMachineEventHandler.dispatch(stateManager, initMediator)
	}
}

//receive a MediatedTransfer, i'm the target
func (rs *Service) targetMediatedTransfer(msg *encoding.MediatedTransfer, ch *channel.Channel) {
	smkey := utils.Sha3(msg.LockSecretHash[:], ch.TokenAddress[:])
	stateManager := rs.Transfer2StateManager[smkey]
	/*
		第一次收到这个密码,
		首先要判断这个密码是否是我声明放弃过的,如果是,就应该谨慎处理.
		锁是有可能重复的,比如 token swap 中.
	*/
	/*
	 *	First time receiving this secret.
	 *	We need to check if I have ever abandoned this secret, if so, handle it carefully.
	 * 	Locks might be duplicate, like in toke swap.
	 */
	if rs.dao.IsLockSecretHashChannelIdentifierDisposed(msg.LockSecretHash, ch.ChannelIdentifier.ChannelIdentifier) {
		//todo 需要通知photon用户
		log.Error(fmt.Sprintf("receive a lock secret hash,and it's my annouce disposed. %s", msg.LockSecretHash.String()))
		return
	}
	if stateManager != nil {
		if stateManager.Name != target.NameTargetTransition {
			log.Error(fmt.Sprintf("receive mediator transfer,but i'm not a target,msg=%s,stateManager=%s", msg, utils.StringInterface(stateManager, 3)))
			return
		}
		log.Error(fmt.Sprintf("receive mediator transfer msg=%s,duplicate? attack?,i'm a target,and has received mediator message. statemanager=%s",
			msg, utils.StringInterface(stateManager, 3)))
		return
	}
	g := rs.getToken2ChannelGraph(ch.TokenAddress)
	fromChannel := g.GetPartenerAddress2Channel(msg.Sender)
	if fromChannel == nil {
		log.Error(fmt.Sprintf("GetPartenerAddress2Channel returns nil ,but %s should have channel with %s on token %s",
			utils.APex2(g.OurAddress), utils.APex2(msg.Sender), utils.APex2(g.TokenAddress)))
		return
	}
	fromRoute := graph.Channel2RouteState(fromChannel, msg.Sender, msg.PaymentAmount, rs, msg.Path)
	fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg, ch.TokenAddress)
	initTarget := &mediatedtransfer.ActionInitTargetStateChange{
		OurAddress:               params.Cfg.MyAddress,
		FromRoute:                fromRoute,
		FromTranfer:              fromTransfer,
		BlockNumber:              rs.GetBlockNumber(),
		Message:                  msg,
		Db:                       rs.dao,
		IsEffectiveChain:         rs.IsChainEffective,
		EffectiveChangeTimestamp: rs.EffectiveChangeTimestamp,
	}
	stateManager = transfer.NewStateManager(target.StateTransiton, nil, target.NameTargetTransition, fromTransfer.LockSecretHash, fromTransfer.Token)
	//rs.dao.AddStateManager(stateManager)
	rs.Transfer2StateManager[smkey] = stateManager
	rs.StateMachineEventHandler.dispatch(stateManager, initTarget)
	// notify upper
	rs.NotifyHandler.NotifyReceiveMediatedTransfer(msg, ch.TokenAddress)
}

func (rs *Service) startHealthCheckFor(address common.Address) {
	if !params.Cfg.EnableHealthCheck {
		return
	}
	if rs.HealthCheckMap[address] {
		log.Info(fmt.Sprintf("addr %s check already start.", utils.APex(address)))
		return
	}
	rs.HealthCheckMap[address] = true
	go func() {
		defer rpanic.PanicRecover(fmt.Sprintf("ping %s", utils.APex(address)))
		log.Trace(fmt.Sprintf("health check for %s started", utils.APex(address)))
		for {
			err := rs.Protocol.SendPing(address)
			if err != nil {
				log.Info(fmt.Sprintf("health check ping %s err %s", utils.APex(address), err))
			}
			time.Sleep(time.Second * 10)
		}
	}()
}

func (rs *Service) startNeighboursHealthCheck() {
	for _, g := range rs.Token2ChannelGraph {
		for addr := range g.PartenerAddress2Channel {
			rs.startHealthCheckFor(addr)
		}
	}
}

//func (rs *Service) startSubscribeNeighborStatus() error {
//	//var err error
//	switch t := rs.Transport.(type) {
//	case *network.MixTransport:
//		//err = t.SubscribeNeighbor(rs.dao)
//		//if err != nil {
//		//	log.Warn(fmt.Sprintf("startSubscribeNeighborStatus when mobile mode  err %s ", err))
//		//}
//	case *network.MatrixMixTransport:
//		t.SetMatrixDB(rs.dao)
//	default:
//		return rerr.ErrTransportTypeUnknown
//	}
//	return nil
//}

func (rs *Service) getToken2ChannelGraph(tokenAddress common.Address) (cg *graph.ChannelGraph) {
	cg = rs.Token2ChannelGraph[tokenAddress]
	if cg == nil {
		log.Error(fmt.Sprintf("%s token doesn't exist ", utils.APex(tokenAddress)))
	}
	return
}
func (rs *Service) getChannelGraph(channelIdentifier common.Hash) (cg *graph.ChannelGraph) {
	ch, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil {
		return
	}
	cg = rs.Token2ChannelGraph[ch.TokenAddress]
	if cg == nil {
		log.Error(fmt.Sprintf("%s token doesn't exist ", utils.APex(ch.TokenAddress)))
	}
	return
}
func (rs *Service) getTokenForChannelIdentifier(channelidentifier common.Hash) (token common.Address) {
	ch, err := rs.findChannelByIdentifier(channelidentifier)
	if err != nil {
		return
	}
	return ch.TokenAddress
}

//only for test, should call findChannelByIdentifier
func (rs *Service) getChannelWithAddr(channelIdentifier common.Hash) *channel.Channel {
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil {

	}

	return c
}

//for test
func (rs *Service) getChannel(tokenAddr, partnerAddr common.Address) *channel.Channel {
	g := rs.getToken2ChannelGraph(tokenAddr)
	if g == nil {
		return nil
	}
	return g.GetPartenerAddress2Channel(partnerAddr)
}

/*
Process user's new channel request
*/
func (rs *Service) newChannelAndDeposit(token, partner common.Address, settleTimeout int, amount *big.Int, isNewChannel bool) *utils.AsyncResult {
	if isNewChannel {
		minSettleTimeout := rs.getMinSettleTimeout()
		if settleTimeout < minSettleTimeout {
			return utils.NewAsyncResultWithError(rerr.ErrArgumentError.Append(fmt.Sprintf("settle_timeout must bigger than %d", minSettleTimeout)))
		}
		g := rs.Token2ChannelGraph[token]
		if g != nil {
			if g.GetPartenerAddress2Channel(partner) != nil {
				return utils.NewAsyncResultWithError(rerr.ErrChannelAlreadExist)
			}
		}
	}
	tokenNetwork, err := rs.Chain.TokenNetwork(token)
	if err != nil {
		panic(err) // never happen
	}
	return utils.NewAsyncResultWithError(tokenNetwork.NewChannelAndDepositAsync(params.Cfg.MyAddress, partner, settleTimeout, amount))
}

/*
process user's close or settle channel request
*/
func (rs *Service) closeOrSettleChannel(channelIdentifier common.Hash, op string) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result = utils.NewAsyncResultWithError(rerr.ErrChannelNotFound)
		return
	}
	log.Trace(fmt.Sprintf("%s channel %s\n", op, utils.HPex(channelIdentifier)))
	if op == closeChannelReqName {
		err = c.Close()
	} else {
		err = c.Settle(rs.GetBlockNumber())
	}
	if err == nil {
		err = rs.UpdateChannelState(channel.NewChannelSerialization(c))
	}
	result.Result <- err
	//通道变化的通知来自于事件,而不是执行结果
	return
}
func (rs *Service) cooperativeSettleChannel(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- rerr.ErrChannelNotFound
		return
	}
	_, isOnline := rs.Protocol.GetNetworkStatus(c.PartnerState.Address)
	if !isOnline {
		result.Result <- rerr.ErrNodeNotOnline.Printf("node %s is not online", c.PartnerState.Address.String())
		return
	}
	// 这里需要先校验下用户余额,防止对方同意了但是自己因为gas不足调用合约失败导致只能强制关闭
	balance, err := rs.Chain.Client.BalanceAt(context.Background(), c.OurState.Address, nil)
	if err != nil {
		result.Result <- rerr.ErrSpectrumNotConnected
		return
	}
	if balance.Cmp(params.Cfg.MinBalance) <= 0 {
		result.Result <- rerr.ErrInsufficientBalanceForGas
		return
	}
	log.Trace(fmt.Sprintf("cooperative settle channel %s\n", utils.HPex(channelIdentifier)))
	s, err := c.CreateCooperativeSettleRequest()
	if err != nil {
		result.Result <- err
		return
	}
	c.State = channeltype.StateCooprativeSettle
	err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	if err != nil {
		result.Result <- err
	}
	err = s.Sign(params.Cfg.PrivateKey, s)
	err = rs.sendAsync(c.PartnerState.Address, s)
	result.Result <- err
	return
}
func (rs *Service) prepareCooperativeSettleChannel(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- rerr.ErrChannelNotFound
		return
	}
	// 查询该通道上是否存在pending状态的deposit,如果有,不允许
	txTypes := fmt.Sprintf("%s,%s", models.TXInfoTypeApproveDeposit, models.TXInfoTypeDeposit)
	pendingDepositList, err := rs.dao.GetTXInfoList(c.ChannelIdentifier.ChannelIdentifier, c.ChannelIdentifier.OpenBlockNumber, utils.EmptyAddress, models.TXInfoType(txTypes), models.TXInfoStatusPending)
	if err != nil {
		result.Result <- err
		return
	}
	if len(pendingDepositList) > 0 {
		result.Result <- rerr.ErrChannelState.Append("can not CooperativeSettle channel when deposit")
		return
	}
	log.Trace(fmt.Sprintf("prepareCooperativeSettleChannel settle channel %s\n", utils.HPex(channelIdentifier)))
	err = c.PrepareForCooperativeSettle()
	if err != nil {
		result.Result <- err
		return
	}
	err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	result.Result <- err
	return
}
func (rs *Service) cancelPrepareForCooperativeSettleChannelOrWithdraw(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- rerr.ErrChannelNotFound
		return
	}
	log.Trace(fmt.Sprintf("cancelPrepareForCooperativeSettleChannelOrWithdraw   channel %s\n", utils.HPex(channelIdentifier)))
	err = c.CancelWithdrawOrCooperativeSettle()
	if err != nil {
		result.Result <- err
		return
	}
	err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	result.Result <- err
	return
}

func (rs *Service) withdraw(channelIdentifier common.Hash, amount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- rerr.ErrChannelNotFound
		return
	}
	if c.State != channeltype.StateOpened && c.State != channeltype.StatePrepareForWithdraw {
		result.Result <- rerr.ErrChannelState.Printf("can not withdraw because state = %s", c.State)
		return
	}
	_, isOnline := rs.Protocol.GetNetworkStatus(c.PartnerState.Address)
	if !isOnline {
		result.Result <- rerr.ErrNodeNotOnline.Printf("node %s is not online", c.PartnerState.Address.String())
		return
	}
	// 查询该通道上是否存在pending状态的deposit,如果有,不允许
	txTypes := fmt.Sprintf("%s,%s", models.TXInfoTypeApproveDeposit, models.TXInfoTypeDeposit)
	pendingDepositList, err := rs.dao.GetTXInfoList(c.ChannelIdentifier.ChannelIdentifier, c.ChannelIdentifier.OpenBlockNumber, utils.EmptyAddress, models.TXInfoType(txTypes), models.TXInfoStatusPending)
	if err != nil {
		result.Result <- err
		return
	}
	if len(pendingDepositList) > 0 {
		result.Result <- rerr.ErrChannelState.Append("can not withdraw on channel when deposit")
		return
	}
	// 这里需要先校验下用户余额,防止对方同意了但是自己因为gas不足调用合约失败导致只能强制关闭
	balance, err := rs.Chain.Client.BalanceAt(context.Background(), c.OurState.Address, nil)
	if err != nil {
		result.Result <- rerr.ErrSpectrumNotConnected
		return
	}
	if balance.Cmp(params.Cfg.MinBalance) <= 0 {
		result.Result <- rerr.ErrInsufficientBalanceForGas
		return
	}
	log.Trace(fmt.Sprintf("withdraw channel %s,amount=%s\n", utils.HPex(channelIdentifier), amount))
	s, err := c.CreateWithdrawRequest(amount)
	if err != nil {
		result.Result <- err
		return
	}
	c.State = channeltype.StateWithdraw
	err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	if err != nil {
		result.Result <- err
	}
	err = s.Sign(params.Cfg.PrivateKey, s)
	err = rs.sendAsync(c.PartnerState.Address, s)
	result.Result <- err
	return
}
func (rs *Service) prepareForWithdraw(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- rerr.ErrChannelNotFound
		return
	}
	log.Trace(fmt.Sprintf("prepareForWithdraw   channel %s,and status=%s\n", utils.HPex(channelIdentifier), c.State))
	err = c.PrepareForWithdraw()
	if err != nil {
		result.Result <- err
		return
	}
	err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	result.Result <- err
	return
}

/*
process user's token swap maker request
todo? save and restore
*/
func (rs *Service) tokenSwapMaker(tokenswap *TokenSwap) (result *utils.AsyncResult) {
	var lockSecretHash common.Hash
	var hasReceiveTakerMediatedTransfer bool
	var sentMtrHook SentMediatedTransferListener
	var receiveMtrHook ReceivedMediatedTrasnferListener
	var secretRequestHook SecretRequestPredictor
	secretRequestHook = func(msg *encoding.SecretRequest) (ignore bool) {
		if !hasReceiveTakerMediatedTransfer {
			/*
				ignore secret request until recieve a valid taker mediated transfer.
				we assume that :
				taker have two independent queue for secret request and mediated transfer
			*/
			return true
		}
		delete(rs.SecretRequestPredictorMap, lockSecretHash) //old hashlock is invalid,just  remove
		return false
	}
	sentMtrHook = func(mtr *encoding.MediatedTransfer) (remove bool) {
		if mtr.LockSecretHash == tokenswap.LockSecretHash && rs.getTokenForChannelIdentifier(mtr.ChannelIdentifier) == tokenswap.FromToken && mtr.Target == tokenswap.ToNodeAddress && mtr.PaymentAmount.Cmp(tokenswap.FromAmount) == 0 {
			if lockSecretHash != utils.EmptyHash {
				log.Info(fmt.Sprintf("tokenswap maker select new path ,because of different hash lock"))
				delete(rs.SecretRequestPredictorMap, lockSecretHash) //old hashlock is invalid,just  remove
			}
			lockSecretHash = mtr.LockSecretHash //hashlock may change when select new route path
			rs.SecretRequestPredictorMap[lockSecretHash] = secretRequestHook
		}
		return false
	}
	receiveMtrHook = func(mtr *encoding.MediatedTransfer) (remove bool) {
		/*
			recevive taker's mediated transfer , the transfer must use argument of tokenswap and have the same hashlock
		*/
		if mtr.LockSecretHash == tokenswap.LockSecretHash && lockSecretHash == mtr.LockSecretHash && rs.getTokenForChannelIdentifier(mtr.ChannelIdentifier) == tokenswap.ToToken && mtr.Target == tokenswap.FromNodeAddress && mtr.PaymentAmount.Cmp(tokenswap.ToAmount) == 0 {
			hasReceiveTakerMediatedTransfer = true
			delete(rs.SentMediatedTransferListenerMap, &sentMtrHook)
			return true
		}
		return false
	}
	rs.SentMediatedTransferListenerMap[&sentMtrHook] = true
	rs.ReceivedMediatedTrasnferListenerMap[&receiveMtrHook] = true
	result, _ = rs.startMediatedTransferInternal(tokenswap.FromToken, tokenswap.ToNodeAddress, tokenswap.FromAmount, tokenswap.LockSecretHash, 0, tokenswap.Secret, "", tokenswap.RouteInfo)
	return
}

/*
taker process token swap
taker's action is triggered by maker's mediated transfer.
*/
func (rs *Service) messageTokenSwapTaker(msg *encoding.MediatedTransfer, tokenswap *TokenSwap) (remove bool) {
	var hashlock = msg.LockSecretHash
	var hasReceiveRevealSecret bool
	var stateManager *transfer.StateManager
	if msg.LockSecretHash != tokenswap.LockSecretHash ||
		msg.PaymentAmount.Cmp(tokenswap.FromAmount) != 0 ||
		msg.Initiator != tokenswap.FromNodeAddress ||
		rs.getTokenForChannelIdentifier(msg.ChannelIdentifier) != tokenswap.FromToken ||
		msg.Target != tokenswap.ToNodeAddress {
		log.Info("receive a mediated transfer, not match tokenswap condition")
		return false
	}
	log.Trace(fmt.Sprintf("begin token swap for %s", msg))
	var secretRequestHook SecretRequestPredictor = func(msg *encoding.SecretRequest) (ignore bool) {
		if !hasReceiveRevealSecret {
			/*
				ignore secret request until recieve a valid reveal secret.
				we assume that :
				maker first send a valid reveal secret and then send secret request, otherwis may deadlock but  taker willnot lose tokens.
			*/
			return true
		}
		return false
	}
	var receiveRevealSecretHook RevealSecretListener = func(msg *encoding.RevealSecret) (remove bool) {
		if msg.LockSecretHash() != hashlock {
			return false
		}
		state := stateManager.CurrentState
		initState, ok := state.(*mediatedtransfer.InitiatorState)
		if !ok {
			panic(fmt.Sprintf("must be a InitiatorState"))
		}
		if initState.Transfer.LockSecretHash != msg.LockSecretHash() {
			panic(fmt.Sprintf("hashlock must be same , state lock=%s,msg lock=%s", utils.HPex(initState.Transfer.LockSecretHash), utils.HPex(msg.LockSecretHash())))
		}
		initState.Transfer.Secret = msg.LockSecret
		hasReceiveRevealSecret = true
		delete(rs.SecretRequestPredictorMap, hashlock)
		return true
	}
	/*
		taker's Expiration must be smaller than maker's ,
		taker and maker may have direct channels on these two tokens.
	*/
	takerExpiration := msg.Expiration - int64(params.Cfg.RevealTimeout)
	result, stateManager := rs.startMediatedTransferInternal(tokenswap.ToToken, tokenswap.FromNodeAddress, tokenswap.ToAmount, tokenswap.LockSecretHash, takerExpiration, utils.EmptyHash, "", tokenswap.RouteInfo)
	if stateManager == nil {
		log.Error(fmt.Sprintf("taker tokenwap error %s", <-result.Result))
		return false
	}
	rs.SecretRequestPredictorMap[hashlock] = secretRequestHook
	rs.RevealSecretListenerMap[hashlock] = receiveRevealSecretHook
	return true
}

/*
process taker's token swap
only mark, if i receive a valid mediated transfer, then start token swap
*/
func (rs *Service) tokenSwapTaker(tokenswap *TokenSwap) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	result.Result <- nil
	key := swapKey{
		LockSecretHash: tokenswap.LockSecretHash,
		FromToken:      tokenswap.FromToken,
		FromAmount:     tokenswap.FromAmount.String(),
	}
	rs.SwapKey2TokenSwap[key] = tokenswap
	return
}

/*
cancel a transfer before secret send
only initiator can call
*/
func (rs *Service) cancelTransfer(req *cancelTransferReq) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	// get transfer info and check
	smKey := utils.Sha3(req.LockSecretHash[:], req.TokenAddress[:])
	manager := rs.Transfer2StateManager[smKey]
	if manager == nil {
		result.Result <- rerr.ErrTransferNotFound
		return
	}
	if manager.Name != initiator.NameInitiatorTransition {
		result.Result <- rerr.ErrTransferCannotCancel.Append("you can only cancel transfers you send")
		return
	}
	transferStatus, err := rs.dao.GetSentTransferDetail(req.TokenAddress, req.LockSecretHash)
	if err != nil {
		result.Result <- rerr.ErrTransferNotFound.Append("can not found transfer status")
		return
	}
	if transferStatus.Status != models.TransferStatusCanCancel {
		result.Result <- rerr.ErrTransferCannotCancel.Printf("status=%d", transferStatus.Status)
		return
	}
	stateChange := &transfer.ActionCancelTransferStateChange{
		LockSecretHash: req.LockSecretHash,
	}
	rs.StateMachineEventHandler.dispatch(manager, stateChange)
	std := rs.dao.UpdateSentTransferDetailStatus(req.TokenAddress, req.LockSecretHash, models.TransferStatusCanceled, "transfer cancel", nil)
	//rs.NotifyTransferStatusChange(req.TokenAddress, req.LockSecretHash, models.TransferStatusCanceled, "交易撤销")
	rs.NotifyHandler.NotifySentTransferDetail(std)
	result.Result <- nil
	return
}

//recieve a ack from
func (rs *Service) handleSentMessage(sentMessage *protocolMessage) {
	data := sentMessage.Message.Pack()
	echohash := utils.Sha3(data, sentMessage.receiver[:])
	_, ok2 := sentMessage.Message.(encoding.EnvelopMessager)
	if ok2 {
		rs.dao.DeleteEnvelopMessager(echohash)
	}
	switch msg := sentMessage.Message.(type) {
	case *encoding.DirectTransfer:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
			return
		}
		smkey := utils.Sha3(msg.FakeLockSecretHash[:], ch.TokenAddress[:])
		if r, ok := rs.Transfer2Result[smkey]; ok {
			r.Result <- nil
		}
		std := rs.dao.UpdateSentTransferDetailStatus(ch.TokenAddress, msg.FakeLockSecretHash, models.TransferStatusSuccess, "DirectTransfer send success,transfer success", ch.ChannelIdentifier)
		//rs.NotifyTransferStatusChange(ch.TokenAddress, msg.FakeLockSecretHash, models.TransferStatusSuccess, "DirectTransfer 发送成功,交易成功")
		rs.NotifyHandler.NotifySentTransferDetail(std)
	case *encoding.MediatedTransfer:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
			return
		}
		rs.dao.UpdateSentTransferDetailStatusMessage(ch.TokenAddress, msg.LockSecretHash, "MediatedTransfer send success")
	case *encoding.RevealSecret:
		// save log to dao
		channels := rs.findAllChannelsByLockSecretHash(msg.LockSecretHash())
		for _, c := range channels {
			rs.dao.UpdateSentTransferDetailStatusMessage(c.TokenAddress, msg.LockSecretHash(), "RevealSecret send success")
		}
	case *encoding.UnLock:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
			return
		}
		std := rs.dao.UpdateSentTransferDetailStatus(ch.TokenAddress, msg.LockSecretHash(), models.TransferStatusSuccess, "UnLock send success,transfer success", ch.ChannelIdentifier)
		//rs.NotifyTransferStatusChange(ch.TokenAddress, msg.LockSecretHash(), models.TransferStatusSuccess, "UnLock 发送成功,交易成功.")
		rs.NotifyHandler.NotifySentTransferDetail(std)
	case *encoding.AnnounceDisposedResponse:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
			return
		}
		rs.dao.UpdateSentTransferDetailStatusMessage(ch.TokenAddress, msg.LockSecretHash, "AnnounceDisposedResponse send success")
	}
	rs.conditionQuitWhenReceiveAck(sentMessage.Message)
	//log.Trace(fmt.Sprintf("msg receive ack :%s", utils.StringInterface(sentMessage, 2)))
}

/*
GetNodeChargeFee implement of FeeCharger
*/
func (rs *Service) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return rs.FeePolicy.GetNodeChargeFee(nodeAddress, tokenAddress, amount)
}

/*
for debug only,quit if eventName exactly match
*/
func (rs *Service) conditionQuit(eventName string) {
	if strings.ToLower(eventName) == strings.ToLower(params.Cfg.ConditionQuit.QuitEvent) && params.Cfg.DebugCrash {
		log.Error(fmt.Sprintf("quitevent=%s\n", eventName))
		//log.Trace(fmt.Sprintf("tokengraph=%s", utils.StringInterface(rs.Token2ChannelGraph, 7)))
		log.Trace(fmt.Sprintf("Transfer2StateManager=%s", utils.StringInterface(rs.Transfer2StateManager, 7)))
		debug.PrintStack()
		//After是发送消息之后,为了确保消息发送成功,主线程sleep100ms
		if strings.Index(eventName, "After") > 0 {
			time.Sleep(time.Millisecond * 100)
		}
		os.Exit(111)
	}
}

/*
GetDao return photon's dao
*/
func (rs *Service) GetDao() models.Dao {
	return rs.dao
}

/*
things to do when Photon connect to eth
*/
func (rs *Service) handleEthRPCConnectionOK() {
	/*
		events before lastHandledBlockNumber must have been processed, so we start from  lastHandledBlockNumber-1
	*/
	rs.BlockChainEvents.Start(rs.dao.GetLatestBlockNumber())
	//启动的时候如果公链 rpc连接有问题,一旦链上,就应该重新初始化 registry, 否则无法进行注册 token 等操作
	// If rpc connection fails in public chain, once reconnecting, we should reinitialize registry,
	// otherwise we can do things like token registry.
	_, err := rs.Chain.Registry(params.Cfg.RegistryAddress, true)
	if err != nil {
		log.Error(fmt.Sprintf("BlockChainService.Registry err =%s", err.Error()))
	}
}

//all user's request
func (rs *Service) handleReq(req *apiReq) {
	var result *utils.AsyncResult
	switch req.Name {
	case transferReqName: //mediated transfer only
		r := req.Req.(*transferReq)
		if r.IsDirectTransfer {
			result = rs.directTransferAsync(r.TokenAddress, r.Target, r.Amount, r.Data)
		} else {
			result = rs.startMediatedTransfer(r.TokenAddress, r.Target, r.Amount, r.Secret, r.Data, r.RouteInfo)
		}
	case newChannelReqName:
		r := req.Req.(*newChannelReq)
		if r.amount != nil && r.amount.Cmp(utils.BigInt0) > 0 {
			result = rs.newChannelAndDeposit(r.tokenAddress, r.partnerAddress, r.settleTimeout, r.amount, r.isNewChannel)
		} else {
			panic("amount must biggner than zero")
		}
	case closeChannelReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.closeOrSettleChannel(r.addr, req.Name)
	case settleChannelReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.closeOrSettleChannel(r.addr, req.Name)
	case tokenSwapMakerReqName:
		r := req.Req.(*tokenSwapMakerReq)
		result = rs.tokenSwapMaker(r.tokenSwap)
	case tokenSwapTakerReqName:
		r := req.Req.(*tokenSwapTakerReq)
		result = rs.tokenSwapTaker(r.tokenSwap)
	case cooperativeSettleChannelReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.cooperativeSettleChannel(r.addr)
	case prepareForCooperativeSettleReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.prepareCooperativeSettleChannel(r.addr)
	case cancelPrepareForCooperativeSettleReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.cancelPrepareForCooperativeSettleChannelOrWithdraw(r.addr)
	case withdrawReqName:
		r := req.Req.(*withdrawReq)
		result = rs.withdraw(r.addr, r.amount)
	case prepareWithdrawReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.prepareForWithdraw(r.addr)
	case cancelPrepareWithdrawReqName:
		r := req.Req.(*closeSettleChannelReq)
		result = rs.cancelPrepareForCooperativeSettleChannelOrWithdraw(r.addr)
	case cancelTransfer:
		r := req.Req.(*cancelTransferReq)
		result = rs.cancelTransfer(r)
	case allowRevealSecretReqName:
		r := req.Req.(*allowRevealSecretReq)
		result = rs.allowRevealSecret(r)
	case registerSecretReqName:
		r := req.Req.(*registerSecretReq)
		result = rs.registerSecretToStateManagerFromUser(r)
	case registerSecretOnChainReqName:
		r := req.Req.(*registerSecretReq)
		result = rs.registerSecretOnChain(r)
	case getUnfinishedReceviedTransferReqName:
		r := req.Req.(*getUnfinishedReceivedTransferReq)
		result = rs.getUnfinishedReceivedTransfer(r)
	case forceUnlockReqName:
		r := req.Req.(*forceUnlockReq)
		result = rs.forceUnlock(r)
	default:
		panic("unkown req")
	}
	r := req
	r.result <- result
}

/*
这一系列update通知没有走callback,而是专门开辟一条道路主要是考虑到这些通知并不是来自用户的请求,
而是
1. photon主动推送
2. 比如交易引起的金额变化,以前是不会通知的,也就没有相应的callback
*/

//UpdateChannelAndSaveAck 保证通道更新和消息确认是一个原子操作
func (rs *Service) UpdateChannelAndSaveAck(c *channel.Channel, tag interface{}) {
	t, ok := tag.(*transfer.MessageTag)
	if !ok || t == nil {
		panic("tag is nil")
	}
	echohash := t.EchoHash
	ack := rs.Protocol.CreateAck(echohash)
	cs := channel.NewChannelSerialization(c)
	err := rs.dao.UpdateChannelAndSaveAck(cs, echohash, ack.Pack())
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelAndSaveAck %s", err))
	}
	rs.NotifyHandler.NotifyChannelStatus(channeltype.ChannelSerialization2ChannelDataDetail(cs))
}

//UpdateChannel 数据库中更新通道状态,同时通知App
func (rs *Service) UpdateChannel(c *channeltype.Serialization, tx models.TX) error {
	err := rs.dao.UpdateChannel(c, tx)
	if err != nil {
		return err
	}
	rs.NotifyHandler.NotifyChannelStatus(channeltype.ChannelSerialization2ChannelDataDetail(c))
	return nil
}

//UpdateChannelNoTx  数据库更新,同时通知App,与updateChannelState的区别就在于回调函数的
func (rs *Service) UpdateChannelNoTx(c *channeltype.Serialization) error {
	err := rs.dao.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	rs.NotifyHandler.NotifyChannelStatus(channeltype.ChannelSerialization2ChannelDataDetail(c))
	return nil
}

//UpdateChannelState 数据库更新,同时通知app
func (rs *Service) UpdateChannelState(c *channeltype.Serialization) error {
	err := rs.dao.UpdateChannelState(c)
	if err != nil {
		return err
	}
	rs.NotifyHandler.NotifyChannelStatus(channeltype.ChannelSerialization2ChannelDataDetail(c))
	return nil
}

//UpdateChannelContractBalance 数据库更新,同时通知app
func (rs *Service) UpdateChannelContractBalance(c *channeltype.Serialization) error {
	err := rs.dao.UpdateChannelContractBalance(c)
	if err != nil {
		return err
	}
	rs.NotifyHandler.NotifyChannelStatus(channeltype.ChannelSerialization2ChannelDataDetail(c))
	return nil
}

func (rs *Service) conditionQuitWhenReceiveAck(msg encoding.Messager) {
	var quitName string
	switch msg.(type) {
	case *encoding.SecretRequest:
		quitName = "ReceiveSecretRequestAck"
	case *encoding.RevealSecret:
		quitName = "ReceiveRevealSecretAck"
	case *encoding.UnLock:
		quitName = "ReceiveUnLockAck"
	case *encoding.DirectTransfer:
		quitName = "ReceiveDirectTransferAck"
	case *encoding.MediatedTransfer:
		quitName = "ReceiveMediatedTransferAck"
	case *encoding.AnnounceDisposed:
		quitName = "ReceiveAnnounceDisposedAck"
	case *encoding.AnnounceDisposedResponse:
		quitName = "ReceiveAnnounceDisposedResponseAck"
	case *encoding.RemoveExpiredHashlockTransfer:
		quitName = "ReceiveRemoveExpiredHashlockTransferAck"
	case *encoding.SettleRequest:
	case *encoding.SettleResponse:
	case *encoding.WithdrawRequest:
	case *encoding.WithdrawResponse:
	default:

	}
	if len(quitName) > 0 {
		rs.conditionQuit(quitName)
	}
}

func (rs *Service) findAllChannelsByLockSecretHash(lockSecretHash common.Hash) (channels []*channel.Channel) {
	for _, lockSecretHash2Channels := range rs.Token2LockSecretHash2Channels {
		chs := lockSecretHash2Channels[lockSecretHash]
		if len(chs) > 0 {
			channels = append(channels, chs...)
		}
	}
	return
}

func (rs *Service) submitDelegateToPms(ch *channel.Channel) {
	select {
	case rs.ChanSubmitDelegateToPMS <- ch:
	default:
		// never block
	}
}

func (rs *Service) submitDelegateToPmsLoop() {
	log.Info("submitDelegateToPmsLoop start...")
	defer rpanic.PanicRecover("submitDelegateToPmsLoop")
	for {
		if rs.PmsProxy == nil {
			log.Info("submitDelegateToPmsLoop stop because PmsProxy is nil")
			return
		}
		select {
		case <-rs.quitChan: //需要有退出机制,否则很快内存耗尽
			return
		case ch, ok := <-rs.ChanSubmitDelegateToPMS:
			if !ok {
				log.Info("submitDelegateToPmsLoop stop because chan close")
				return
			}
			// 1. 获取待提交数据
			data, err := rs.GetDelegateForPms(channel.NewChannelSerialization(ch), utils.EmptyAddress)
			if err != nil {
				log.Error(fmt.Sprintf("submitDelegateToPmsLoop GetDelegateForPms of channel %s err : %s,ignore", ch.ChannelIdentifier.ChannelIdentifier.String(), err.Error()))
				continue
			}
			// 2. 提交至pms,失败最多重试3次
			err = rs.PmsProxy.SubmitDelegate(data)
			if err == pmsproxy.ErrConnect {
				// 网络错误重试3次
				for i := 0; i < 3; i++ {
					err = rs.PmsProxy.SubmitDelegate(data)
					if err == pmsproxy.ErrConnect {
						continue
					} else {
						break
					}
				}
			}
			// 3. 根据提交结果修改db中的channel状态
			if err == nil {
				ch.DelegateState = channeltype.ChannelDelegateStateSuccess
			} else {
				ch.DelegateState = channeltype.ChannelDelegateStateFail
				log.Error(err.Error())
			}
			err = rs.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
			if err != nil {
				log.Error("")
				continue
			}
			// 4. 如果提交成功,修改已提交的punish及SendAnnounceDispose状态
			if ch.DelegateState == channeltype.ChannelDelegateStateSuccess {
				rs.dao.MarkLockHashCanPunishSubmittedByChannel(ch.ChannelIdentifier.ChannelIdentifier)
				rs.dao.MarkSendAnnounceDisposeSubmittedByChannel(ch.ChannelIdentifier.ChannelIdentifier)
			}
		}
	}
}

func (rs *Service) submitBalanceProofToPfs(ch *channel.Channel) {
	select {
	case rs.ChanSubmitBalanceProofToPFS <- ch:
	default:
		// never block
	}
}

func (rs *Service) submitBalanceProofToPfsLoop() {
	defer rpanic.PanicRecover("submitBalanceProofToPfsLoop")
	log.Info("submitBalanceProofToPfsLoop start...")
	for {
		if rs.PfsProxy == nil {
			log.Info("submitBalanceProofToPfsLoop stop because PfsProxy is nil")
			return
		}
		select {
		case <-rs.quitChan: //需要有退出机制,否则很快内存耗尽
			return
		case ch, ok := <-rs.ChanSubmitBalanceProofToPFS:
			if !ok {
				log.Info("submitBalanceProofToPfsLoop stop because chan close")
				return
			}
			bpPartner := ch.PartnerState.BalanceProofState
			err := rs.PfsProxy.SubmitBalance(
				bpPartner.Nonce,
				bpPartner.TransferAmount,
				ch.Outstanding(),
				ch.ChannelIdentifier.OpenBlockNumber,
				bpPartner.LocksRoot,
				ch.ChannelIdentifier.ChannelIdentifier,
				bpPartner.MessageHash,
				ch.PartnerState.Address,
				bpPartner.Signature,
			)
			if err == pfsproxy.ErrConnect {
				log.Warn(fmt.Sprintf("connect to pfs err when submit BalanceProof of channel %s, retry 3 times...", ch.ChannelIdentifier.ChannelIdentifier.String()))
				// 网络错误,重发3次
				for i := 0; i < 3; i++ {
					err = rs.PfsProxy.SubmitBalance(
						bpPartner.Nonce,
						bpPartner.TransferAmount,
						ch.Outstanding(),
						ch.ChannelIdentifier.OpenBlockNumber,
						bpPartner.LocksRoot,
						ch.ChannelIdentifier.ChannelIdentifier,
						bpPartner.MessageHash,
						ch.PartnerState.Address,
						bpPartner.Signature,
					)
					if err == pfsproxy.ErrConnect {
						continue
					} else {
						break
					}
				}
			}
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}

func (rs *Service) getBestRoutesFromPfs(peerFrom, peerTo, token common.Address, amount *big.Int, isInitiator bool) (routes []*route.State, err error) {
	var paths []pfsproxy.FindPathResponse
	paths, err = rs.PfsProxy.FindPath(peerFrom, peerTo, token, amount, isInitiator)
	if err != nil {
		return
	}
	for _, path := range paths {
		if path.Result == nil || path.Result[0] == "" {
			continue
		}
		partnerAddress := common.HexToAddress(path.Result[0])
		ch := rs.getChannel(token, partnerAddress)
		if ch == nil {
			continue
		}
		r := route.NewState(ch, path.GetPath())
		r.Fee = rs.FeePolicy.GetNodeChargeFee(partnerAddress, token, amount)
		r.TotalFee = path.Fee
		routes = append(routes, r)
	}
	return
}
func (rs *Service) forceUnlock(req *forceUnlockReq) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	channelIdentifier := req.ChannelIdentifier
	lockSecretHash := utils.ShaSecret(req.Secret[:])
	secret := req.Secret
	channel := rs.getChannelWithAddr(channelIdentifier)
	if channel == nil {
		result.Result <- rerr.ErrChannelNotFound.Printf("can not find channel %s", channelIdentifier.String())
		return
	}
	tokenNetwork, err := rs.Chain.TokenNetwork(channel.TokenAddress)
	if err != nil {
		result.Result <- err
		return
	}
	// 获取数据
	isSecretRegistered := false
	lock := channel.PartnerState.Lock2PendingLocks[lockSecretHash].Lock
	if lock == nil {
		lock = channel.PartnerState.Lock2UnclaimedLocks[lockSecretHash].Lock
		if lock == nil {
			result.Result <- rerr.ErrTransferNotFound.Printf("can not find lock by lockSecretHash : %s", lockSecretHash.String())
			return
		}
		isSecretRegistered = channel.PartnerState.Lock2UnclaimedLocks[lockSecretHash].IsRegisteredOnChain
	}
	partnerAddress := channel.PartnerState.Address
	transferAmount := channel.PartnerState.BalanceProofState.TransferAmount
	//contractTransferAmout := channel.PartnerState.BalanceProofState.ContractTransferAmount
	locksroot := channel.PartnerState.BalanceProofState.LocksRoot
	nonce := channel.PartnerState.BalanceProofState.Nonce
	signature := channel.PartnerState.BalanceProofState.Signature
	addtionalHash := channel.PartnerState.BalanceProofState.MessageHash
	proof := channel.PartnerState.Tree.MakeProof(lock.Hash())
	//不能阻塞主线程
	go func() {
		//测试接口,不用担心资源释放问题
		defer rpanic.PanicRecover("forceUnlock")
		if channel.State == channeltype.StateOpened {
			// 自己close
			log.Trace(fmt.Sprintf("forceUnlock close : partnerAddress=%s, transferAmount=%d, locksroot=%s nonce=%d, addtionalHash=%s,signature=%s\n",
				partnerAddress.String(), transferAmount, locksroot.String(), nonce, addtionalHash.String(), common.Bytes2Hex(signature)))
			err = tokenNetwork.CloseChannel(partnerAddress, transferAmount, locksroot, nonce, addtionalHash, signature)
			if err != nil {
				result.Result <- rerr.ErrCloseChannel.Printf("forceUnlock : close channel fail %s", err.Error())
				return
			}
		}
		retry := 0
		for {
			channel = rs.getChannelWithAddr(channelIdentifier)
			if channel.State != channeltype.StateClosed {
				time.Sleep(time.Second)
				retry++
				if retry > 10 {
					break
				}
				continue
			}
			break
		}
		if !isSecretRegistered {
			// register
			err = rs.Chain.SecretRegistryProxy.RegisterSecret(secret)
			if err != nil {
				result.Result <- rerr.ErrRegisterSecret.Errorf("ForceUnlock : register secret fail %s", err.Error())
				return
			}
			retry = 0
			for {
				isSecretRegistered, err = rs.Chain.SecretRegistryProxy.IsSecretRegistered(secret)
				if err != nil {
					result.Result <- rerr.ErrRegisterSecret.Errorf("ForceUnlock : register secret fail %s", err.Error())
					return
				}
				if isSecretRegistered {
					break
				}
				retry++
				if retry > 10 {
					break
				}
				time.Sleep(time.Second)
			}

		}
		// unlock
		log.Trace(fmt.Sprintf("forceUnlock unlock : partnerAddress=%s, transferAmount=%d, expiration=%d, amount=%d,lockSecretHash=%s,proof=%s lockHash=%s \n",
			partnerAddress.String(), transferAmount, lock.Expiration, lock.Amount, lock.LockSecretHash.String(), common.Bytes2Hex(mtree.Proof2Bytes(proof)),
			lock.Hash().String()))

		err = tokenNetwork.Unlock(partnerAddress, channel.PartnerState.BalanceProofState.ContractTransferAmount, lock, mtree.Proof2Bytes(proof))
		if err != nil {
			result.Result <- rerr.ErrUnlock.Errorf("forceUnlock : unlock failed %s", err.Error())
			return
		}
		log.Info(fmt.Sprintf("forceUnlock success %s ,partner=%s", lockSecretHash.String(), utils.APex(partnerAddress)))
		result.Result <- nil
	}()
	return
}

func (rs *Service) getUnfinishedReceivedTransfer(req *getUnfinishedReceivedTransferReq) (result *utils.AsyncResult) {
	lockSecretHash := req.LockSecretHash
	tokenAddress := req.TokenAddress
	result = utils.NewAsyncResult()
	//token swap 过滤
	if rs.SecretRequestPredictorMap[lockSecretHash] != nil {
		result.Result <- rerr.ErrTransferNotFound.Errorf("SecretRequestPredictorMap has lockSecretHash")
		return
	}
	key := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := rs.Transfer2StateManager[key]
	if manager == nil {
		result.Result <- rerr.ErrTransferNotFound.Printf("can not find transfer by lock_secret_hash[%s] and token_address[%s]", lockSecretHash.String(), tokenAddress.String())
		return
	}
	state, ok := manager.CurrentState.(*mediatedtransfer.TargetState)
	if !ok {
		// 接收人不是自己
		// I'm not the recipient
		result.Result <- rerr.ErrTransferNotFound.Errorf("i'm not recipient")
		return
	}
	resp := new(TransferDataResponse)
	resp.Initiator = state.FromTransfer.Initiator.String()
	resp.Target = state.FromTransfer.Target.String()
	resp.Token = tokenAddress.String()
	resp.Amount = state.FromTransfer.Amount
	resp.LockSecretHash = state.FromTransfer.LockSecretHash.String()
	resp.Expiration = state.FromTransfer.Expiration - state.BlockNumber
	result.Tag = resp
	result.Result <- nil
	return
}

func (rs *Service) allowRevealSecret(req *allowRevealSecretReq) (result *utils.AsyncResult) {
	lockSecretHash := req.LockSecretHash
	tokenAddress := req.TokenAddress
	result = utils.NewAsyncResult()
	key := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := rs.Transfer2StateManager[key]
	if manager == nil {
		result.Result <- rerr.InvalidState("can not find transfer by lock_secret_hash and token_address")
		return
	}
	state, ok := manager.CurrentState.(*mediatedtransfer.InitiatorState)
	if !ok {
		result.Result <- rerr.InvalidState("wrong state")
		return
	}
	if lockSecretHash != state.LockSecretHash || lockSecretHash != utils.ShaSecret(state.Secret.Bytes()) {
		result.Result <- rerr.InvalidState("wrong lock_secret_hash")
	}
	delete(rs.SecretRequestPredictorMap, lockSecretHash)
	log.Trace(fmt.Sprintf("Remove SecretRequestPredictor for lockSecretHash=%s", lockSecretHash.String()))
	result.Result <- nil
	return
}

func (rs *Service) registerSecretToStateManagerFromUser(req *registerSecretReq) (result *utils.AsyncResult) {
	secret := req.Secret
	tokenAddress := req.TokenAddress
	lockSecretHash := utils.ShaSecret(secret.Bytes())
	result = utils.NewAsyncResult()
	//在channel 中注册密码
	// register secret in channel
	rs.registerSecret(secret)

	key := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := rs.Transfer2StateManager[key]
	if manager == nil {
		result.Result <- rerr.InvalidState("can not find transfer by lock_secret_hash and token_address")
	}
	state, ok := manager.CurrentState.(*mediatedtransfer.TargetState)
	if !ok {
		result.Result <- rerr.InvalidState("wrong state")
		return
	}
	if lockSecretHash != state.FromTransfer.LockSecretHash {
		result.Result <- rerr.InvalidState("wrong secret")
		return
	}
	// 在state manager中注册密码
	// register secret in state manager
	state.FromTransfer.Secret = secret
	state.Secret = secret
	result.Result <- nil
	return
}

func (rs *Service) registerSecretOnChain(req *registerSecretReq) (result *utils.AsyncResult) {
	secret := req.Secret
	return rs.Chain.SecretRegistryProxy.RegisterSecretAsync(secret)
}

// SetBuildInfo 启动时保存构建信息
func (rs *Service) SetBuildInfo(goVersion, gitCommit, buildDate, version string) {
	rs.BuildInfo.GoVersion = goVersion
	rs.BuildInfo.GitCommit = gitCommit
	rs.BuildInfo.BuildDate = buildDate
	rs.BuildInfo.Version = version
}

////NotifyTransferStatusChange notify status change of a sending transfer
//func (rs *Service) NotifyTransferStatusChange(tokenAddress common.Address, lockSecretHash common.Hash, status models.TransferStatusCode, statusMessage string) {
//}

/*
在处理完一个通道中的锁之后,释放该锁在这个map中占据的内存
*/
func (rs *Service) removeToken2LockSecretHash2channel(secretHash common.Hash, ch *channel.Channel) {
	if ch == nil {
		return
	}
	m, ok := rs.Token2LockSecretHash2Channels[ch.TokenAddress]
	if !ok {
		return
	}
	chs, ok := m[secretHash]
	if !ok {
		return
	}
	index := -1
	for i, c := range chs {
		if c.ChannelIdentifier.ChannelIdentifier == ch.ChannelIdentifier.ChannelIdentifier &&
			c.ChannelIdentifier.OpenBlockNumber == ch.ChannelIdentifier.OpenBlockNumber {
			index = i
		}
	}
	if index >= 0 {
		chs = append(chs[:index], chs[index+1:]...)
	}
	if len(chs) == 0 {
		delete(m, secretHash)
	}
}

/*
获取最小SettleTimeout值
*/
func (rs *Service) getMinSettleTimeout() int {
	return params.Cfg.ChannelSettleTimeoutMin
}

// GetDelegateForPms 获取一个通道需要提交给pms的数据
func (rs *Service) GetDelegateForPms(c *channeltype.Serialization, thirdAddr common.Address) (result *pmsproxy.DelegateForPms, err error) {
	if thirdAddr == utils.EmptyAddress {
		thirdAddr = params.Cfg.PmsAddress
	}
	var sig []byte
	// 1. 基础通道数据
	c3 := new(pmsproxy.DelegateForPms)
	c3.ChannelIdentifier = c.ChannelIdentifier.ChannelIdentifier
	c3.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	c3.TokenAddrss = c.TokenAddress()
	c3.PartnerAddress = c.PartnerAddress()
	if c.PartnerBalanceProof == nil {
		result = c3
		return
	}
	// 2. partner的BalanceProof数据
	if c.PartnerBalanceProof.Nonce > 0 {
		c3.UpdateTransfer.Nonce = c.PartnerBalanceProof.Nonce
		c3.UpdateTransfer.TransferAmount = c.PartnerBalanceProof.TransferAmount
		c3.UpdateTransfer.Locksroot = c.PartnerBalanceProof.LocksRoot
		c3.UpdateTransfer.ExtraHash = c.PartnerBalanceProof.MessageHash
		c3.UpdateTransfer.ClosingSignature = c.PartnerBalanceProof.Signature
		sig, err = pmsproxy.SignBalanceProofFor3rd(c, params.Cfg.PrivateKey)
		if err != nil {
			return
		}
		c3.UpdateTransfer.NonClosingSignature = sig
	}
	// 3. unlock数据
	tree := mtree.NewMerkleTree(c.PartnerLeaves)
	var ws []*pmsproxy.DelegateUnlock
	for _, l := range c.PartnerLeaves {
		proof := channel.ComputeProofForLock(l, tree)
		w := &pmsproxy.DelegateUnlock{
			Lock:        l,
			MerkleProof: mtree.Proof2Bytes(proof.MerkleProof),
		}
		w.Signature, err = pmsproxy.SignUnlockFor3rd(c, w, thirdAddr, params.Cfg.PrivateKey)
		ws = append(ws, w)
	}
	c3.Unlocks = ws
	// 4. punish数据
	var ps []*pmsproxy.DelegatePunish
	for _, annouceDisposed := range rs.dao.GetChannelAnnounceDisposed(c.ChannelIdentifier.ChannelIdentifier) {
		//跳过历史 channel
		// omit history channel
		if annouceDisposed.OpenBlockNumber != c.ChannelIdentifier.OpenBlockNumber {
			continue
		}
		// 跳过已经提交过的
		if annouceDisposed.IsSubmittedToPms {
			continue
		}
		p := &pmsproxy.DelegatePunish{
			LockHash:       common.BytesToHash(annouceDisposed.LockHash),
			AdditionalHash: annouceDisposed.AdditionalHash,
			Signature:      annouceDisposed.Signature,
		}
		ps = append(ps, p)
	}
	c3.Punishes = ps
	// 5. SendAnnounceDispose数据
	var sas []*pmsproxy.DelegateAnnounceDisposed
	for _, sendAnnounceDispose := range rs.dao.GetSendAnnounceDisposeByChannel(c.ChannelIdentifier.ChannelIdentifier, false) {
		sas = append(sas, &pmsproxy.DelegateAnnounceDisposed{
			LockSecretHash: common.BytesToHash(sendAnnounceDispose.LockSecretHash),
		})
	}
	c3.AnnouceDisposed = sas
	// 6. 密码注册相关数据,这里只提交当前需要委托注册的,PMS那边采用全量覆盖
	var secrets []*pmsproxy.DelegateSecret
	for _, ourKnownSecret := range c.PartnerKnownSecrets {
		if !ourKnownSecret.IsRegisteredOnChain {
			// 如果没在链上注册过,需要委托
			secret := ourKnownSecret.Secret
			secretHash := utils.ShaSecret(secret[:])
			lock := c.PartnerLock2UnclaimedLocks()[secretHash]
			secrets = append(secrets, &pmsproxy.DelegateSecret{
				Secret:        ourKnownSecret.Secret,
				RegisterBlock: lock.Lock.Expiration - int64(c.RevealTimeout) + 5, // TODO 这里委托注册时间延迟5块,目的是当需要注册密码的时候自己也在线,就自己注册,避免PMS扣费
			})
		}
	}
	c3.Secrets = secrets
	result = c3
	return
}
