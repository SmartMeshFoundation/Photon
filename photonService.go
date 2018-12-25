package photon

import (
	"crypto/ecdsa"
	"errors"

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
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/theckman/go-flock"
)

/*
message sent complete notification
*/
type protocolMessage struct {
	receiver common.Address
	Message  encoding.Messager
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
	Config                   *params.Config
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

	/*
	 */
	PrivateKey            *ecdsa.PrivateKey
	NodeAddress           common.Address
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
		TODO remove this,this design is very weird
	*/
	Token2Hashlock2Channels map[common.Address]map[common.Hash][]*channel.Channel
	FileLocker              *flock.Flock
	BlockNumber             *atomic.Value
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
	isStarting                            bool
	StopCreateNewTransfers                bool // 是否停止接收新交易,默认false,目前仅在用户调用prepare-update接口的时候,会被置为true,直到重启		// boolean to check whether stop receiving new transfers, default to false. Currently it sets to true when clients invoke prepare-update, till it reconnects.
	EthConnectionStatus                   chan netshare.Status
	ChanHistoryContractEventsDealComplete chan struct{}
}

//NewPhotonService create photon service
func NewPhotonService(chain *rpc.BlockChainService, privateKey *ecdsa.PrivateKey, transport network.Transporter, config *params.Config, notifyHandler *notify.Handler, dao models.Dao) (rs *Service, err error) {
	rs = &Service{
		NotifyHandler:      notifyHandler,
		Chain:              chain,
		PrivateKey:         privateKey,
		Config:             config,
		Transport:          transport,
		dao:                dao,
		NodeAddress:        crypto.PubkeyToAddress(privateKey.PublicKey),
		Token2ChannelGraph: make(map[common.Address]*graph.ChannelGraph),
		//todo fixme Token2TokenNetwork 应该是一个token的数组,表示已经注册的token.目前k,v中的v必须是空地址
		Token2TokenNetwork:                    make(map[common.Address]common.Address),
		Transfer2StateManager:                 make(map[common.Hash]*transfer.StateManager),
		Transfer2Result:                       make(map[common.Hash]*utils.AsyncResult),
		Token2Hashlock2Channels:               make(map[common.Address]map[common.Hash][]*channel.Channel),
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
		isStarting:                            true,
		StopCreateNewTransfers:                false,
		EthConnectionStatus:                   make(chan netshare.Status, 10),
		ChanHistoryContractEventsDealComplete: make(chan struct{}),
	}
	rs.BlockNumber.Store(int64(0))
	rs.MessageHandler = newPhotonMessageHandler(rs)
	rs.StateMachineEventHandler = newStateMachineEventHandler(rs)
	rs.Protocol = network.NewPhotonProtocol(transport, privateKey, rs)
	//todo fixme MatrixTransport should have a better contructor function
	mtransport, ok := rs.Transport.(*network.MatrixMixTransport)
	if ok {
		err = mtransport.SetMatrixDB(rs.dao)
		if err != nil {
			return
		}
	}
	rs.Protocol.SetReceivedMessageSaver(NewAckHelper(rs.dao))
	/*
		only one instance for one data directory
	*/
	rs.FileLocker = flock.NewFlock(config.DataBasePath + ".flock.Lock")
	locked, err := rs.FileLocker.TryLock()
	if err != nil || !locked {
		err = fmt.Errorf("another instance already running at %s", config.DataBasePath)
		return
	}
	log.Info(fmt.Sprintf("create photon service registry=%s,node=%s", rs.Chain.GetRegistryAddress().String(), rs.NodeAddress.String()))

	rs.Token2TokenNetwork, err = rs.dao.GetAllTokens()
	if err != nil {
		return
	}
	rs.BlockChainEvents = blockchain.NewBlockChainEvents(chain.Client, chain)
	// pathfinder
	if config.PfsHost != "" {
		rs.PfsProxy = pfsproxy.NewPfsProxy(config.PfsHost, rs.PrivateKey)
	}
	// fee module
	if config.EnableMediationFee {
		rs.FeePolicy, err = NewFeeModule(dao, rs.PfsProxy)
		if err != nil {
			return
		}
	} else {
		rs.FeePolicy = &NoFeePolicy{}
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

	err = rs.registerRegistry()
	if err != nil {
		return
	}
	rs.Protocol.Start()
	rs.restore()

	go func() {
		if rs.Config.ConditionQuit.RandomQuit {
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
	}
	//
	rs.isStarting = false
	rs.startNeighboursHealthCheck()
	// 只有在混合模式下启动时,才订阅其他节点的在线状态
	// Only when starting under MixUDPXMPP, we can subscribe online status of other nodes.
	if rs.Config.NetworkMode == params.MixUDPXMPP || rs.Config.NetworkMode == params.MixUDPMatrix {
		err = rs.startSubscribeNeighborStatus()
		if err != nil {
			err = fmt.Errorf("startSubscribeNeighborStatus err %s", err)
			return
		}
	}
	return nil
}

//Stop the node.
func (rs *Service) Stop() {
	log.Info("photon service stop...")
	close(rs.quitChan)
	rs.Protocol.StopAndWait()
	rs.BlockChainEvents.Stop()
	rs.Chain.Client.Close()
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
					rs.BlockNumber.Store(blockStateChange.BlockNumber)
				} else {
					log.Trace(fmt.Sprintf("statechange received :%s", utils.StringInterface(st, 2)))
				}
				err = rs.StateMachineEventHandler.OnBlockchainStateChange(st)
				if err != nil {
					log.Error(fmt.Sprintf("stateMachineEventHandler.OnBlockchainStateChange %s", err))
				}
				if ok2 {
					rs.dao.SaveLatestBlockNumber(blockStateChange.BlockNumber)
					if rs.ChanHistoryContractEventsDealComplete != nil {
						close(rs.ChanHistoryContractEventsDealComplete)
						rs.ChanHistoryContractEventsDealComplete = nil
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
			select {
			case rs.EthConnectionStatus <- s:
			default:
				//never block
			}
			if s == netshare.Connected {
				rs.handleEthRPCConnectionOK()
			} else {
				rs.NotifyHandler.Notify(notify.LevelWarn, "公链连接失败,正在尝试重连")
			}
		case <-rs.quitChan:
			log.Info(fmt.Sprintf("%s quit now", utils.APex2(rs.NodeAddress)))
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
		err = fmt.Errorf("registerRegistry err:%s", err)
		return
	}
	for token := range token2TokenNetworks {
		err = rs.registerTokenNetwork(token)
		if err != nil {
			err = fmt.Errorf("registerTokenNetwork err:%s", err)
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
	ourState := channel.NewChannelEndState(rs.NodeAddress, big.NewInt(0), nil, mtree.NewMerkleTree(nil))
	partenerState := channel.NewChannelEndState(partnerAddress, big.NewInt(0), nil, mtree.NewMerkleTree(nil))

	externState := channel.NewChannelExternalState(rs.registerChannelForHashlock, tokenNetwork, channelIdentifier, rs.PrivateKey, rs.Chain.Client, rs.dao, 0, rs.NodeAddress, partnerAddress)
	ch, err = channel.NewChannel(ourState, partenerState, externState, tokenAddress, channelIdentifier, rs.Config.RevealTimeout, settleTimeout)
	return
}

/*
block chain tick,
it's the core of HTLC
*/
func (rs *Service) handleBlockNumber(blocknumber int64) {
	statechange := &transfer.BlockStateChange{BlockNumber: blocknumber}
	rs.BlockNumber.Store(blocknumber)
	/*
		todo when to remove statemanager ?
			when currentState==nil && StateManager.ManagerState!=StateManagerStateInit ,should delete rs statemanager.
	*/
	rs.StateMachineEventHandler.dispatchToAllTasks(statechange)
	for _, cg := range rs.Token2ChannelGraph {
		for _, c := range cg.ChannelIdentifier2Channel {
			err := rs.StateMachineEventHandler.ChannelStateTransition(c, statechange)
			if err != nil {
				log.Error(fmt.Sprintf("ChannelStateTransition err %s", err))
			}
		}
	}
	rs.dao.SaveLatestBlockNumber(blocknumber)
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
	if recipient == rs.NodeAddress {
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
	go func() {
		defer rpanic.PanicRecover(fmt.Sprintf("send %s, msg:%s", utils.APex(recipient), msg))
		<-result.Result //always success
		rs.ProtocolMessageSendComplete <- &protocolMessage{
			receiver: recipient,
			Message:  msg,
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
	for _, hashchannel := range rs.Token2Hashlock2Channels {
		for _, ch := range hashchannel[hashlock] {
			err := ch.RegisterSecret(secret)
			err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
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
	for _, hashchannel := range rs.Token2Hashlock2Channels {
		for _, ch := range hashchannel[lockSecretHash] {
			err := ch.RegisterRevealedSecretHash(lockSecretHash, secret, blockNumber)
			err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
			if err != nil {
				log.Error(fmt.Sprintf("RegisterSecret %s to channel %s  err: %s",
					utils.HPex(lockSecretHash), ch.ChannelIdentifier.String(), err))
			}
		}
	}
}
func (rs *Service) registerChannelForHashlock(netchannel *channel.Channel, hashlock common.Hash) {
	tokenAddress := netchannel.TokenAddress
	channelsRegistered := rs.Token2Hashlock2Channels[tokenAddress][hashlock]
	found := false
	for _, c := range channelsRegistered {
		//To determine whether the two channel objects are equal, we simply use the address to identify.
		if c.ExternState.ChannelIdentifier == netchannel.ExternState.ChannelIdentifier {
			found = true
			break
		}
	}
	if !found {
		hashLock2Channels, ok := rs.Token2Hashlock2Channels[tokenAddress]
		if !ok {
			hashLock2Channels = make(map[common.Hash][]*channel.Channel)
			rs.Token2Hashlock2Channels[tokenAddress] = hashLock2Channels
		}
		channelsRegistered = append(channelsRegistered, netchannel)
		rs.Token2Hashlock2Channels[tokenAddress][hashlock] = channelsRegistered
	}
}
func (rs *Service) channelSerilization2Channel(c *channeltype.Serialization, tokenNetwork *rpc.TokenNetworkProxy) (ch *channel.Channel, err error) {
	OurState := channel.NewChannelEndState(c.OurAddress, c.OurContractBalance,
		c.OurBalanceProof, mtree.NewMerkleTree(c.OurLeaves))
	PartnerState := channel.NewChannelEndState(c.PartnerAddress(),
		c.PartnerContractBalance,
		c.PartnerBalanceProof, mtree.NewMerkleTree(c.PartnerLeaves))
	ExternState := channel.NewChannelExternalState(rs.registerChannelForHashlock, tokenNetwork,
		c.ChannelIdentifier, rs.PrivateKey,
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
	tokenNetwork, err = rs.Chain.TokenNetworkWithoutCheck(tokenAddress)
	if err != nil {
		return
	}
	edges, err := rs.dao.GetAllNonParticipantChannelByToken(tokenAddress)
	if err != nil {
		return
	}
	g := graph.NewChannelGraph(rs.NodeAddress, tokenAddress, edges)
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
		result.Result <- errors.New("token not exist")
		return
	}
	directChannel := g.GetPartenerAddress2Channel(target)
	if directChannel == nil || !directChannel.CanTransfer() || directChannel.Distributable().Cmp(amount) < 0 {
		result.Result <- errors.New("no available direct channel")
		return
	}
	tr, err := directChannel.CreateDirectTransfer(amount)
	if err != nil {
		result.Result <- err
		return
	}
	tr.Data = []byte(data)
	err = tr.Sign(rs.PrivateKey, tr)
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
	rs.dao.NewTransferStatus(tokenAddress, tr.FakeLockSecretHash)
	err = rs.sendAsync(directChannel.PartnerState.Address, tr)
	if err != nil {
		result.Result <- err
		return
	}
	rs.dao.UpdateTransferStatusMessage(tokenAddress, tr.FakeLockSecretHash, "DirectTransfer 正在发送")
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
func (rs *Service) startMediatedTransferInternal(tokenAddress, target common.Address, amount *big.Int, fee *big.Int, lockSecretHash common.Hash, expiration int64, secret common.Hash, data string) (result *utils.AsyncResult, stateManager *transfer.StateManager) {
	var availableRoutes []*route.State
	var err error
	targetAmount := new(big.Int).Sub(amount, fee)
	result = utils.NewAsyncResult()
	if rs.PfsProxy != nil {
		availableRoutes, err = rs.getBestRoutesFromPfs(rs.NodeAddress, target, tokenAddress, targetAmount)
		if err != nil {
			result.Result <- errors.New("get route from pathfinder failed")
			return
		}
	} else {
		g := rs.getToken2ChannelGraph(tokenAddress)
		if g == nil {
			result.Result <- errors.New("token not exist")
			return
		}
		availableRoutes = g.GetBestRoutes(rs.Protocol, rs.NodeAddress, target, amount, targetAmount, graph.EmptyExlude, rs)
	}
	if len(availableRoutes) <= 0 {
		result.Result <- errors.New("no available route")
		return
	}
	if rs.Config.IsMeshNetwork {
		result.Result <- errors.New("no mediated transfer on mesh only network")
		return
	}
	/*
		when user specified fee, for test or other purpose.
	*/
	if fee.Cmp(utils.BigInt0) > 0 {
		for _, r := range availableRoutes {
			r.TotalFee = fee //use the user's fee to replace algorithm's
		}
	}
	routesState := route.NewRoutesState(availableRoutes)
	transferState := &mediatedtransfer.LockedTransferState{
		TargetAmount:   new(big.Int).Set(amount),
		Amount:         new(big.Int).Set(amount),
		Token:          tokenAddress,
		Initiator:      rs.NodeAddress,
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
		OurAddress:     rs.NodeAddress,
		Tranfer:        transferState,
		Routes:         routesState,
		BlockNumber:    rs.GetBlockNumber(),
		Secret:         secret,
		LockSecretHash: lockSecretHash,
		Db:             rs.dao,
	}
	stateManager = transfer.NewStateManager(initiator.StateTransition, nil, initiator.NameInitiatorTransition, lockSecretHash, transferState.Token)
	smkey := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := rs.Transfer2StateManager[smkey]
	if manager != nil {
		panic(fmt.Sprintf("manager must be never exist"))
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
func (rs *Service) startMediatedTransfer(tokenAddress, target common.Address, amount *big.Int, fee *big.Int, secret common.Hash, data string) (result *utils.AsyncResult) {
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
	rs.dao.NewTransferStatus(tokenAddress, lockSecretHash)
	result, _ = rs.startMediatedTransferInternal(tokenAddress, target, amount, fee, lockSecretHash, 0, secret, data)
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
	targetAmount := new(big.Int).Sub(msg.PaymentAmount, msg.Fee)
	amount := msg.PaymentAmount
	targetAddr := msg.Target
	fromChannel := ch
	fromRoute := graph.Channel2RouteState(fromChannel, msg.Sender, amount, rs)
	fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg, ch.TokenAddress)
	if stateManager != nil {
		if stateManager.Name != mediator.NameMediatorTransition {
			log.Error(fmt.Sprintf("receive mediator transfer,but i'm not a mediator,msg=%s,stateManager=%s", msg, utils.StringInterface(stateManager, 3)))
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
		ourAddress := rs.NodeAddress
		exclude := graph.MakeExclude(msg.Sender, msg.Initiator)
		var avaiableRoutes []*route.State
		if rs.PfsProxy != nil {
			var err error
			avaiableRoutes, err = rs.getBestRoutesFromPfs(rs.NodeAddress, targetAddr, tokenAddress, targetAmount)
			if err != nil {
				log.Error(fmt.Sprintf("get route from pathfinder failed, err = %s", err.Error()))
			}
		} else {
			g := rs.getToken2ChannelGraph(ch.TokenAddress) //must exist
			avaiableRoutes = g.GetBestRoutes(rs.Protocol, rs.NodeAddress, targetAddr, amount, targetAmount, exclude, rs)
		}
		routesState := route.NewRoutesState(avaiableRoutes)
		blockNumber := rs.GetBlockNumber()
		initMediator := &mediatedtransfer.ActionInitMediatorStateChange{
			OurAddress:  ourAddress,
			FromTranfer: fromTransfer,
			Routes:      routesState,
			FromRoute:   fromRoute,
			BlockNumber: blockNumber,
			Message:     msg,
			Db:          rs.dao,
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
		todo 首先要判断这个密码是否是我声明放弃过的,如果是,就应该谨慎处理.
		锁是有可能重复的,比如 token swap 中.
	*/
	/*
	 *	First time receiving this secret.
	 *	todo We need to check if I have ever abandoned this secret, if so, handle it carefully.
	 * 	Locks might be duplicate, like in toke swap.
	 */
	if rs.dao.IsLockSecretHashChannelIdentifierDisposed(msg.LockSecretHash, ch.ChannelIdentifier.ChannelIdentifier) {
		log.Error(fmt.Sprintf("receive a lock secret hash,and it's my annouce disposed. %s", msg.LockSecretHash.String()))
		//忽略,什么都不做
		// do nothing
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
	fromRoute := graph.Channel2RouteState(fromChannel, msg.Sender, msg.PaymentAmount, rs)
	fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg, ch.TokenAddress)
	initTarget := &mediatedtransfer.ActionInitTargetStateChange{
		OurAddress:  rs.NodeAddress,
		FromRoute:   fromRoute,
		FromTranfer: fromTransfer,
		BlockNumber: rs.GetBlockNumber(),
		Message:     msg,
		Db:          rs.dao,
	}
	stateManager = transfer.NewStateManager(target.StateTransiton, nil, target.NameTargetTransition, fromTransfer.LockSecretHash, fromTransfer.Token)
	//rs.dao.AddStateManager(stateManager)
	rs.Transfer2StateManager[smkey] = stateManager
	rs.StateMachineEventHandler.dispatch(stateManager, initTarget)
	// notify upper
	rs.NotifyHandler.NotifyReceiveMediatedTransfer(msg, ch)
}

func (rs *Service) startHealthCheckFor(address common.Address) {
	if !rs.Config.EnableHealthCheck {
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
func (rs *Service) startSubscribeNeighborStatus() error {
	var err error
	switch t := rs.Transport.(type) {
	case *network.MixTransport:
		err = t.SubscribeNeighbor(rs.dao)
	case *network.MatrixMixTransport:
		err = t.SetMatrixDB(rs.dao)
	default:
		return fmt.Errorf("transport is not mix or matrix transpoter,can't subscribe neighbor status")
	}
	if err != nil {
		log.Warn(fmt.Sprintf("startSubscribeNeighborStatus when mobile mode  err %s ", err))
		return nil
	}
	return err
}
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
func (rs *Service) newChannelAndDeposit(token, partner common.Address, settleTimeout int, amount *big.Int, isNewChannel bool) (result *utils.AsyncResult) {
	if isNewChannel {
		g := rs.Token2ChannelGraph[token]
		if g != nil {
			if g.GetPartenerAddress2Channel(partner) != nil {
				return utils.NewAsyncResultWithError(errors.New("channel already exists"))
			}
		}
	}
	tokenNetwork, err := rs.Chain.TokenNetwork(token)
	if err != nil {
		result = utils.NewAsyncResultWithError(err)
		return
	}
	result = tokenNetwork.NewChannelAndDepositAsync(rs.NodeAddress, partner, settleTimeout, amount)
	return
}

/*
process user's close or settle channel request
*/
func (rs *Service) closeOrSettleChannel(channelIdentifier common.Hash, op string) (result *utils.AsyncResult) {
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result = utils.NewAsyncResultWithError(errors.New("channel not exist"))
		return
	}
	log.Trace(fmt.Sprintf("%s channel %s\n", op, utils.HPex(channelIdentifier)))
	if op == closeChannelReqName {
		result = c.Close()
	} else {
		result = c.Settle()
	}
	return
}
func (rs *Service) cooperativeSettleChannel(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- errors.New("channel not exist")
		return
	}
	_, isOnline := rs.Protocol.GetNetworkStatus(c.PartnerState.Address)
	if !isOnline {
		result.Result <- fmt.Errorf("node %s is not online", c.PartnerState.Address.String())
		return
	}
	log.Trace(fmt.Sprintf("cooperative settle channel %s\n", utils.HPex(channelIdentifier)))
	s, err := c.CreateCooperativeSettleRequest()
	if err != nil {
		result.Result <- err
		return
	}
	c.State = channeltype.StateCooprativeSettle
	err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	if err != nil {
		result.Result <- err
	}
	err = s.Sign(rs.PrivateKey, s)
	err = rs.sendAsync(c.PartnerState.Address, s)
	result.Result <- err
	return
}
func (rs *Service) prepareCooperativeSettleChannel(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- errors.New("channel not exist")
		return
	}
	log.Trace(fmt.Sprintf("prepareCooperativeSettleChannel settle channel %s\n", utils.HPex(channelIdentifier)))
	err = c.PrepareForCooperativeSettle()
	if err != nil {
		result.Result <- err
		return
	}
	err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	result.Result <- err
	return
}
func (rs *Service) cancelPrepareForCooperativeSettleChannelOrWithdraw(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- errors.New("channel not exist")
		return
	}
	log.Trace(fmt.Sprintf("cancelPrepareForCooperativeSettleChannelOrWithdraw   channel %s\n", utils.HPex(channelIdentifier)))
	err = c.CancelWithdrawOrCooperativeSettle()
	if err != nil {
		result.Result <- err
		return
	}
	err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	result.Result <- err
	return
}

func (rs *Service) withdraw(channelIdentifier common.Hash, amount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- errors.New("channel not exist")
		return
	}
	_, isOnline := rs.Protocol.GetNetworkStatus(c.PartnerState.Address)
	if !isOnline {
		result.Result <- fmt.Errorf("node %s is not online", c.PartnerState.Address.String())
		return
	}
	log.Trace(fmt.Sprintf("withdraw channel %s,amount=%s\n", utils.HPex(channelIdentifier), amount))
	s, err := c.CreateWithdrawRequest(amount)
	if err != nil {
		result.Result <- err
		return
	}
	c.State = channeltype.StateWithdraw
	err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	if err != nil {
		result.Result <- err
	}
	err = s.Sign(rs.PrivateKey, s)
	err = rs.sendAsync(c.PartnerState.Address, s)
	result.Result <- err
	return
}
func (rs *Service) prepareForWithdraw(channelIdentifier common.Hash) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByIdentifier(channelIdentifier)
	if err != nil { //settled channel can be queried from dao.
		result.Result <- errors.New("channel not exist")
		return
	}
	log.Trace(fmt.Sprintf("prepareForWithdraw   channel %s\n", utils.HPex(channelIdentifier)))
	err = c.PrepareForWithdraw()
	if err != nil {
		result.Result <- err
		return
	}
	err = rs.dao.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	result.Result <- err
	return
}

/*
process user's token swap maker request
save and restore todo?
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
	result, _ = rs.startMediatedTransferInternal(tokenswap.FromToken, tokenswap.ToNodeAddress, tokenswap.FromAmount, utils.BigInt0, tokenswap.LockSecretHash, 0, tokenswap.Secret, "")
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
	takerExpiration := msg.Expiration - int64(rs.Config.RevealTimeout)
	result, stateManager := rs.startMediatedTransferInternal(tokenswap.ToToken, tokenswap.FromNodeAddress, tokenswap.ToAmount, utils.BigInt0, tokenswap.LockSecretHash, takerExpiration, utils.EmptyHash, "")
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
		result.Result <- errors.New("can not found transfer")
		return
	}
	if manager.Name != initiator.NameInitiatorTransition {
		result.Result <- errors.New("you can only cancel transfers you send")
		return
	}
	transferStatus, err := rs.dao.GetTransferStatus(req.TokenAddress, req.LockSecretHash)
	if err != nil {
		result.Result <- errors.New("can not found transfer status")
		return
	}
	if transferStatus.Status != models.TransferStatusCanCancel {
		result.Result <- errors.New("transfer already can not cancel now")
		return
	}
	stateChange := &transfer.ActionCancelTransferStateChange{
		LockSecretHash: req.LockSecretHash,
	}
	rs.StateMachineEventHandler.dispatch(manager, stateChange)
	rs.dao.UpdateTransferStatus(req.TokenAddress, req.LockSecretHash, models.TransferStatusCanceled, "交易撤销")
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
		}
		smkey := utils.Sha3(msg.FakeLockSecretHash[:], ch.TokenAddress[:])
		if r, ok := rs.Transfer2Result[smkey]; ok {
			r.Result <- nil
		}
		rs.dao.UpdateTransferStatus(ch.TokenAddress, msg.FakeLockSecretHash, models.TransferStatusSuccess, "DirectTransfer 发送成功,交易成功")
	case *encoding.MediatedTransfer:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
		}
		rs.dao.UpdateTransferStatusMessage(ch.TokenAddress, msg.LockSecretHash, "MediatedTransfer 发送成功")
	case *encoding.RevealSecret:
		// save log to dao
		channels := rs.findAllChannelsByLockSecretHash(msg.LockSecretHash())
		for _, c := range channels {
			rs.dao.UpdateTransferStatusMessage(c.TokenAddress, msg.LockSecretHash(), "RevealSecret 发送成功")
		}
	case *encoding.UnLock:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
		}
		rs.dao.UpdateTransferStatus(ch.TokenAddress, msg.LockSecretHash(), models.TransferStatusSuccess, "UnLock 发送成功,交易成功.")
	case *encoding.AnnounceDisposedResponse:
		ch, err := rs.findChannelByIdentifier(msg.ChannelIdentifier)
		if err != nil {
			log.Error(err.Error())
		}
		rs.dao.UpdateTransferStatusMessage(ch.TokenAddress, msg.LockSecretHash, "AnnounceDisposedResponse 发送成功")
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
	if strings.ToLower(eventName) == strings.ToLower(rs.Config.ConditionQuit.QuitEvent) && rs.Config.DebugCrash {
		log.Error(fmt.Sprintf("quitevent=%s\n", eventName))
		debug.PrintStack()
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
	rs.Chain.Registry(rs.Config.RegistryAddress, true)
	// 重连时上传手续费设置给PFS
	if fm, ok := rs.FeePolicy.(*FeeModule); ok {
		err := fm.SubmitFeePolicyToPFS()
		if err != nil {
			log.Error(fmt.Sprintf("set fee policy to pfs err =%s", err.Error()))
		}
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
			result = rs.startMediatedTransfer(r.TokenAddress, r.Target, r.Amount, r.Fee, r.Secret, r.Data)
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
	default:
		panic("unkown req")
	}
	r := req
	r.result <- result
}

func (rs *Service) updateChannelAndSaveAck(c *channel.Channel, tag interface{}) {
	t, ok := tag.(*transfer.MessageTag)
	if !ok || t == nil {
		panic("tag is nil")
	}
	echohash := t.EchoHash
	ack := rs.Protocol.CreateAck(echohash)
	err := rs.dao.UpdateChannelAndSaveAck(channel.NewChannelSerialization(c), echohash, ack.Pack())
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelAndSaveAck %s", err))
	}
}

func (rs *Service) conditionQuitWhenReceiveAck(msg encoding.Messager) {
	var quitName string
	switch msg.(type) {
	case *encoding.SecretRequest:
		quitName = "ReceiveSecretRequestAck"
	case *encoding.RevealSecret:
		quitName = "ReceiveRevealSecretAck"
	case *encoding.UnLock:
	case *encoding.DirectTransfer:
	case *encoding.MediatedTransfer:
		quitName = "ReceiveMediatedTransferAck"
	case *encoding.AnnounceDisposed:
		quitName = "ReceiveAnnounceDisposedAck"
	case *encoding.AnnounceDisposedResponse:
	case *encoding.RemoveExpiredHashlockTransfer:
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
	for _, lockSecretHash2Channels := range rs.Token2Hashlock2Channels {
		chs := lockSecretHash2Channels[lockSecretHash]
		if len(chs) > 0 {
			channels = append(channels, chs...)
		}
	}
	return
}

func (rs *Service) submitBalanceProofToPfs(ch *channel.Channel) {
	if rs.PfsProxy == nil {
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
		bpPartner.Signature,
	)
	if err != nil {
		log.Error(fmt.Sprintf("updateBalanceProofToPfs err = %s", err.Error()))
	}
}

func (rs *Service) getBestRoutesFromPfs(peerFrom, peerTo, token common.Address, amount *big.Int) (routes []*route.State, err error) {
	var paths []pfsproxy.FindPathResponse
	paths, err = rs.PfsProxy.FindPath(peerFrom, peerTo, token, amount)
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
		r := route.NewState(ch)
		r.Fee = rs.FeePolicy.GetNodeChargeFee(partnerAddress, token, amount)
		r.TotalFee = path.Fee
		routes = append(routes, r)
	}
	return
}
