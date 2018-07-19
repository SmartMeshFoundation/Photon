package smartraiden

import (
	"crypto/ecdsa"
	"errors"

	"fmt"

	"path/filepath"

	"time"

	"sync/atomic"

	"math/big"

	"strings"

	"os"

	"runtime/debug"

	"github.com/SmartMeshFoundation/SmartRaiden/blockchain"
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/graph"
	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/fee"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/route"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
RaidenService is a raiden node
most of raidenService's member is not thread safe, and should not visit outside the loop method.
*/
type RaidenService struct {
	Chain                 *rpc.BlockChainService
	Registry              *rpc.RegistryProxy
	SecretRegistryAddress common.Address
	RegistryAddress       common.Address
	PrivateKey            *ecdsa.PrivateKey
	Transport             network.Transporter
	Config                *params.Config
	Protocol              *network.RaidenProtocol
	NodeAddress           common.Address
	Token2ChannelGraph    map[common.Address]*graph.ChannelGraph

	TokenNetwork2Token       map[common.Address]common.Address
	Token2TokenNetwork       map[common.Address]common.Address
	Transfer2StateManager    map[common.Hash]*transfer.StateManager
	Transfer2Result          map[common.Hash]*utils.AsyncResult
	SwapKey2TokenSwap        map[swapKey]*TokenSwap
	Tokens2ConnectionManager map[common.Address]*ConnectionManager //how to save and restore for token swap? todo fix it
	/*
				   This is a map from a hashlock to a list of channels, the same
			         hashlock can be used in more than one token (for tokenswaps), a
			         channel should be removed from this list only when the Lock is
			         released/withdrawn but not when the secret is registered.
		TODO remove this,this design is very weird
	*/
	Token2Hashlock2Channels  map[common.Address]map[common.Hash][]*channel.Channel
	MessageHandler           *raidenMessageHandler
	StateMachineEventHandler *stateMachineEventHandler
	BlockChainEvents         *blockchain.Events
	AlarmTask                *blockchain.AlarmTask
	db                       *models.ModelDB
	FileLocker               *flock.Flock
	SnapshortDir             string
	BlockNumber              *atomic.Value
	/*
		new block event
	*/
	BlockNumberChan chan int64
	/*
		chan for user request
	*/
	UserReqChan                 chan *apiReq
	ProtocolMessageSendComplete chan *protocolMessage
	FeePolicy                   fee.Charger //Mediation fee
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
	ReceivedMediatedTrasnferListenerMap map[*ReceivedMediatedTrasnferListener]bool //for tokenswap
	SentMediatedTransferListenerMap     map[*SentMediatedTransferListener]bool     //for tokenswap
	HealthCheckMap                      map[common.Address]bool
	quitChan                            chan struct{} //for quit notification
	ethInited                           bool
	EthConnectionStatus                 chan netshare.Status
}

//NewRaidenService create raiden service
func NewRaidenService(chain *rpc.BlockChainService, privateKey *ecdsa.PrivateKey, transport network.Transporter, config *params.Config) (rs *RaidenService, err error) {
	if config.SettleTimeout < params.ChannelSettleTimeoutMin || config.SettleTimeout > params.ChannelSettleTimeoutMax {
		err = fmt.Errorf("settle timeout must be in range %d-%d",
			params.ChannelSettleTimeoutMin, params.ChannelSettleTimeoutMax)
		return
	}
	rs = &RaidenService{
		Chain:                               chain,
		Registry:                            chain.Registry(chain.RegistryAddress),
		RegistryAddress:                     chain.RegistryAddress,
		PrivateKey:                          privateKey,
		Config:                              config,
		Transport:                           transport,
		NodeAddress:                         crypto.PubkeyToAddress(privateKey.PublicKey),
		Token2ChannelGraph:                  make(map[common.Address]*graph.ChannelGraph),
		TokenNetwork2Token:                  make(map[common.Address]common.Address),
		Token2TokenNetwork:                  make(map[common.Address]common.Address),
		Transfer2StateManager:               make(map[common.Hash]*transfer.StateManager),
		Transfer2Result:                     make(map[common.Hash]*utils.AsyncResult),
		Token2Hashlock2Channels:             make(map[common.Address]map[common.Hash][]*channel.Channel),
		SwapKey2TokenSwap:                   make(map[swapKey]*TokenSwap),
		Tokens2ConnectionManager:            make(map[common.Address]*ConnectionManager),
		AlarmTask:                           blockchain.NewAlarmTask(chain.Client),
		BlockNumberChan:                     make(chan int64, 20), //not block alarm task
		UserReqChan:                         make(chan *apiReq, 10),
		BlockNumber:                         new(atomic.Value),
		ProtocolMessageSendComplete:         make(chan *protocolMessage, 10),
		SecretRequestPredictorMap:           make(map[common.Hash]SecretRequestPredictor),
		RevealSecretListenerMap:             make(map[common.Hash]RevealSecretListener),
		ReceivedMediatedTrasnferListenerMap: make(map[*ReceivedMediatedTrasnferListener]bool),
		SentMediatedTransferListenerMap:     make(map[*SentMediatedTransferListener]bool),
		FeePolicy:                           &ConstantFeePolicy{},
		HealthCheckMap:                      make(map[common.Address]bool),
		quitChan:                            make(chan struct{}),
		EthConnectionStatus:                 make(chan netshare.Status, 10),
	}
	rs.BlockNumber.Store(int64(0))
	rs.MessageHandler = newRaidenMessageHandler(rs)
	rs.StateMachineEventHandler = newStateMachineEventHandler(rs)
	rs.Protocol = network.NewRaidenProtocol(transport, privateKey, rs)
	rs.db, err = models.OpenDb(config.DataBasePath)
	if err != nil {
		err = fmt.Errorf("open db error %s", err)
		return
	}
	rs.Protocol.SetReceivedMessageSaver(NewAckHelper(rs.db))
	/*
		only one instance for one data directory
	*/
	rs.FileLocker = flock.NewFlock(config.DataBasePath + ".flock.Lock")
	locked, err := rs.FileLocker.TryLock()
	if err != nil || !locked {
		err = fmt.Errorf("another instance already running at %s", config.DataBasePath)
		return
	}
	rs.SnapshortDir = filepath.Join(config.DataBasePath)
	log.Info(fmt.Sprintf("create raiden service registry=%s,node=%s", rs.RegistryAddress.String(), rs.NodeAddress.String()))
	if rs.Registry != nil {
		//我已经连接到以太坊全节点
		rs.SecretRegistryAddress, err = rs.Registry.GetContract().Secret_registry_address(nil)
		if err != nil {
			return
		}
		rs.db.SaveSecretRegistryAddress(rs.SecretRegistryAddress)
	} else {
		//读取数据库中存放的 SecretRegistryAddress, 如果没有,说明系统没有初始化过,只能退出.
		rs.SecretRegistryAddress = rs.db.GetSecretRegistryAddress()
		if rs.SecretRegistryAddress == utils.EmptyAddress {
			err = fmt.Errorf("first startup without ethereum rpc connection")
			return
		}
	}
	rs.Token2TokenNetwork, err = rs.db.GetAllTokens()
	if err != nil {
		return
	}
	for t, tn := range rs.Token2TokenNetwork {
		rs.TokenNetwork2Token[tn] = t
	}
	rs.BlockChainEvents = blockchain.NewBlockChainEvents(chain.Client, chain.RegistryAddress, rs.SecretRegistryAddress, rs.Token2TokenNetwork)
	return rs, nil
}

// Start the node.
func (rs *RaidenService) Start() (err error) {

	rs.AlarmTask.RegisterCallback(func(number int64) error {
		rs.db.SaveLatestBlockNumber(number)
		return rs.setBlockNumber(number)
	})
	if rs.Chain.Client.IsConnected() {
		lastHandledBlockNumber := rs.db.GetLatestBlockNumber()
		err = rs.AlarmTask.Start()
		if err != nil {
			log.Error(fmt.Sprintf("alarm task start err %s", err))
			n := rs.db.GetLatestBlockNumber()
			rs.BlockNumber.Store(n)
		} else {
			//must have a valid blocknumber before any transfer operation
			rs.BlockNumber.Store(rs.AlarmTask.LastBlockNumber)
		}
		/*
			events before lastHandledBlockNumber must have been processed, so we start from  lastHandledBlockNumber-1
		*/
		err = rs.BlockChainEvents.Start(lastHandledBlockNumber)
		if err != nil {
			err = fmt.Errorf("Events listener error %v", err)
			return
		}
		/*
			  Registry registration must start *after* the alarm task, rs avoid
				 corner cases were the registry is queried in block A, a new block B
				 is mined, and the alarm starts polling at block C.
		*/
		rs.registerRegistry()
		rs.ethInited = true
	} else {
		log.Warn(fmt.Sprintf("raiden start without eth rpc server"))
		rs.ethInited = false
		err = rs.startWithoutEthRPC()
		if err != nil {

		}
	}

	err = rs.restoreSnapshot()
	if err != nil {
		err = fmt.Errorf("restore from snapshot error : %v\n you can delete all the database %s to run. but all your trade will lost", err, rs.Config.DataBasePath)
		return
	}
	rs.Protocol.Start()
	rs.startNeighboursHealthCheck()
	err = rs.startSubscribeNeighborStatus()
	if err != nil {
		err = fmt.Errorf("startSubscribeNeighborStatus err %s", err)
		return
	}
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
	return nil
}
func (rs *RaidenService) startWithoutEthRPC() (err error) {
	log.Info("start raiden service without eth connection")
	/*
	   从数据库中把 channel 状态恢复了,等连接上以后再重新把链上信息取过来初始化.
	*/
	chmap := make(map[common.Address][]*channeltype.Serialization)
	cs, err := rs.db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		return
	}
	//group by token
	for _, c := range cs {
		if c.State == channeltype.StateSettled {
			continue
		}
		cs2 := chmap[c.TokenAddress()]
		cs2 = append(cs2, c)
		chmap[c.TokenAddress()] = cs2
	}
	for t, cs := range chmap {
		var details []*graph.ChannelDetails
		for _, c := range cs {
			d := rs.makeChannelDetailFromChannelSerialization(c)
			details = append(details, d)
		}
		graph := graph.NewChannelGraph(rs.NodeAddress, t, nil, details)
		rs.Token2ChannelGraph[t] = graph
		rs.Tokens2ConnectionManager[t] = NewConnectionManager(rs, t)
	}
	return nil
}
func (rs *RaidenService) makeChannelDetailFromChannelSerialization(c *channeltype.Serialization) *graph.ChannelDetails {
	return nil
	//tokenAddress := c.TokenAddress
	//addr1, b1, addr2, b2 := c.OurAddress, c.OurBalance, c.PartnerAddress, c.PartnerBalance
	//var ourAddr, partnerAddr common.Address
	//var ourBalance, partnerBalance *big.Int
	//if addr1 == rs.NodeAddress {
	//	ourAddr = addr1
	//	partnerAddr = addr2
	//	ourBalance = b1
	//	partnerBalance = b2
	//} else {
	//	ourAddr = addr2
	//	partnerAddr = addr1
	//	ourBalance = b2
	//	partnerBalance = b1
	//}
	//ourAddr, ourBalance, partnerAddr, partnerBalance = c.OurAddress, c.OurContractBalance, c.PartnerAddress, c.PartnerContractBalance
	//proxy, err := rs.Chain.TokenNetworkWithoutCheck(c.ChannelAddress)
	//if err != nil {
	//	log.Error(fmt.Sprintf("TokenNetworkWithoutCheck err %s", err))
	//}
	//ourState := channel.NewChannelEndState(ourAddr, ourBalance, nil, transfer.EmptyMerkleTreeState)
	//partenerState := channel.NewChannelEndState(partnerAddr, partnerBalance, nil, transfer.EmptyMerkleTreeState)
	//channelAddress := proxy.Address
	//registerChannelForHashlock := func(channel *channel.Channel, hashlock common.Hash) {
	//	rs.registerChannelForHashlock(tokenAddress, channel, hashlock)
	//}
	//externState := channel.NewChannelExternalState(registerChannelForHashlock, proxy, channelAddress, rs.Chain, rs.db, 0, c.ClosedBlock)
	//channelDetail := &graph.ChannelDetails{
	//	ChannelIdentifier: channelAddress,
	//	OurState:          ourState,
	//	PartenerState:     partenerState,
	//	ExternState:       externState,
	//	BlockChainService: rs.Chain,
	//	RevealTimeout:     rs.Config.RevealTimeout,
	//	SettleTimeout:     c.SettleTimeout,
	//}
	//return channelDetail
}

//Stop the node.
func (rs *RaidenService) Stop() {
	log.Info("raiden service stop...")
	close(rs.quitChan)
	rs.AlarmTask.Stop()
	rs.Protocol.StopAndWait()
	rs.BlockChainEvents.Stop()
	rs.Chain.Client.Close()
	rs.saveSnapshot()
	time.Sleep(100 * time.Millisecond) // let other goroutines quit
	rs.db.CloseDB()
	//anther instance cann run now
	rs.FileLocker.Unlock()
	log.Info("raiden service stop ok...")
}

/*
main loop of this raiden nodes
process  events below:
1. request from user
2. event from blockchain
3. message from other nodes.
*/
func (rs *RaidenService) loop() {
	var err error
	var ok bool
	var m *network.MessageToRaiden
	var st transfer.StateChange
	var blockNumber int64
	var req *apiReq
	var sentMessage *protocolMessage
	defer rpanic.PanicRecover("raiden service")
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
				err = rs.StateMachineEventHandler.OnBlockchainStateChange(st)
				if err != nil {
					log.Error(fmt.Sprintf("stateMachineEventHandler.OnBlockchainStateChange %s", err))
				}
			} else {
				log.Info("Events.StateChangeChannel closed")
				return
			}
			// new block event, it's the timer of raiden
		case blockNumber, ok = <-rs.BlockNumberChan:
			if ok {
				rs.handleBlockNumber(blockNumber)
			} else {
				log.Info("BlockNumberChan closed")
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
				rs.handleEthRRCConnectionOK()
			}
		case <-rs.quitChan:
			log.Info(fmt.Sprintf("%s quit now", utils.APex2(rs.NodeAddress)))
			return
		}
	}
}

//for init,read db history,只要是我还没处理的链上事件,都还在队列中等着发给我.
func (rs *RaidenService) registerRegistry() {
	token2TokenNetworks, err := rs.db.GetAllTokens()
	if err != nil {
		err = fmt.Errorf("registerRegistry err:%s", err)
		return
	}
	for token, tokenNetwork := range token2TokenNetworks {
		err = rs.registerTokenNetwork(token, tokenNetwork)
		if err != nil {
			err = fmt.Errorf("registerTokenNetwork err:%s", err)
			return
		}
	}
}

/*
quering my channel details on blockchain
i'm one of the channel participants
收到来自链上的事件,新创建了 channel
但是事件有可能重复
*/
func (rs *RaidenService) getChannelDetail(tokenAddress common.Address, proxy *rpc.TokenNetworkProxy, partnerAddress common.Address) (channelDetail *graph.ChannelDetails, err error) {
	/*
		因为有可能在我离线的时候收到一堆事件,所以通道的信息不一定就是新创建时候的状态,
		但是保证后续的事件会继续收到,所以应该按照新通道处理.
	*/

	var ourBalance, partnerBalance *big.Int

	channelID, _, openBlockNumber, _, settleTimeout, err := proxy.GetChannelInfo(rs.NodeAddress, partnerAddress)
	if err != nil {
		return
	}
	ourBalance, _, _, err = proxy.GetChannelParticipantInfo(rs.NodeAddress, partnerAddress)
	if err != nil {
		return
	}
	partnerBalance, _, _, err = proxy.GetChannelParticipantInfo(partnerAddress, rs.NodeAddress)
	if err != nil {
		return
	}
	ourState := channel.NewChannelEndState(rs.NodeAddress, ourBalance, nil, mtree.NewMerkleTree(nil))
	partenerState := channel.NewChannelEndState(partnerAddress, partnerBalance, nil, mtree.NewMerkleTree(nil))
	channelAddress := &contracts.ChannelUniqueID{
		ChannelIdentifier: channelID,
		OpenBlockNumber:   int64(openBlockNumber),
	}

	externState := channel.NewChannelExternalState(rs.registerChannelForHashlock, proxy, channelAddress, rs.PrivateKey, rs.Chain.Client, rs.db, 0, rs.NodeAddress, partnerAddress)
	channelDetail = &graph.ChannelDetails{
		ChannelIdentifier: channelID,
		OpenBlockNumber:   int64(openBlockNumber),
		OurState:          ourState,
		PartenerState:     partenerState,
		ExternState:       externState,
		RevealTimeout:     rs.Config.RevealTimeout,
		SettleTimeout:     int(settleTimeout),
	}
	return
}

func (rs *RaidenService) setBlockNumber(blocknumber int64) error {
	rs.BlockNumberChan <- blocknumber
	return nil
}

/*
block chain tick,
it's the core of HTLC
*/
func (rs *RaidenService) handleBlockNumber(blocknumber int64) error {
	statechange := &transfer.BlockStateChange{BlockNumber: blocknumber}
	rs.BlockNumber.Store(blocknumber)
	/*
		todo when to remove statemanager ?
			when currentState==nil && StateManager.ManagerState!=StateManagerStateInit ,should delete rs statemanager.
	*/
	rs.StateMachineEventHandler.logAndDispatchToAllTasks(statechange)
	for _, cg := range rs.Token2ChannelGraph {
		for _, channel := range cg.ChannelAddress2Channel {
			rs.StateMachineEventHandler.ChannelStateTransition(channel, statechange)
		}
	}

	return nil
}

//GetBlockNumber return latest blocknumber of ethereum
func (rs *RaidenService) GetBlockNumber() int64 {
	return rs.BlockNumber.Load().(int64)
}

func (rs *RaidenService) findChannelByAddress(nettingChannelAddress common.Hash) (*channel.Channel, error) {
	for _, g := range rs.Token2ChannelGraph {
		ch := g.GetChannelAddress2Channel(nettingChannelAddress)
		if ch != nil {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("unknown channel %s", nettingChannelAddress)
}

/*
Send `message` to `recipient` using the raiden protocol.

       The protocol will take care of resending the message on a given
       interval until an Acknowledgment is received or a given number of
       tries.
*/
func (rs *RaidenService) sendAsync(recipient common.Address, msg encoding.SignedMessager) error {
	if recipient == rs.NodeAddress {
		log.Error(fmt.Sprintf("rs must be a bug ,sending message to it self"))
	}
	revealSecretMessage, ok := msg.(*encoding.RevealSecret)
	if ok && revealSecretMessage != nil {
		srs := models.NewSentRevealSecret(revealSecretMessage, recipient)
		rs.db.NewSentRevealSecret(srs)
		if msg.Tag() == nil {
			msg.SetTag(&transfer.MessageTag{
				EchoHash:          srs.EchoHash,
				IsASendingMessage: true,
			})
		} else {
			messageTag := msg.Tag().(*transfer.MessageTag)
			if messageTag.EchoHash != srs.EchoHash {
				panic("reveal secret's echo hash not equal")
			}
		}
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
func (rs *RaidenService) SendAndWait(recipient common.Address, message encoding.SignedMessager, timeout time.Duration) error {
	return rs.Protocol.SendAndWait(recipient, message, timeout)
}

/*
Register the secret with any channel that has a hashlock on it.

       This must search through all channels registered for a given hashlock
       and ignoring the tokens.
*/
func (rs *RaidenService) registerSecret(secret common.Hash) {
	hashlock := utils.Sha3(secret[:])
	revealSecretMessage := encoding.NewRevealSecret(secret)
	revealSecretMessage.Sign(rs.PrivateKey, revealSecretMessage)
	for _, hashchannel := range rs.Token2Hashlock2Channels {
		for _, ch := range hashchannel[hashlock] {
			err := ch.RegisterSecret(secret)
			if err != nil {
				log.Error(fmt.Sprintf("RegisterSecret %s to channel %s  err: %s",
					utils.HPex(secret), ch.ChannelIdentifier, err))
			}
			rs.conditionQuit("BeforeSendRevealSecret")
			rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
			//The protocol ignores duplicated messages.
			//make sure not send the same instance multi times.
			rs.sendAsync(ch.PartnerState.Address, encoding.CloneRevealSecret(revealSecretMessage))
		}
	}
}

/*
链上这个锁对应的密码注册了,
*/
func (rs *RaidenService) registerRevealedLockSecretHash(lockSecretHash common.Hash, blockNumber int64) {
	for _, hashchannel := range rs.Token2Hashlock2Channels {
		for _, ch := range hashchannel[lockSecretHash] {
			err := ch.RegisterRevealedSecretHash(lockSecretHash, blockNumber)
			if err != nil {
				log.Error(fmt.Sprintf("RegisterSecret %s to channel %s  err: %s",
					utils.HPex(lockSecretHash), ch.ChannelIdentifier, err))
			}
			rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		}
	}
}
func (rs *RaidenService) registerChannelForHashlock(netchannel *channel.Channel, hashlock common.Hash) {
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
func (rs *RaidenService) channelSerilization2ChannelDetail(c *channeltype.Serialization, tokenNetwork *rpc.TokenNetworkProxy) (d *graph.ChannelDetails) {
	d = &graph.ChannelDetails{
		ChannelIdentifier: c.ChannelIdentifier.ChannelIdentifier,
		OpenBlockNumber:   c.ChannelIdentifier.OpenBlockNumber,
		RevealTimeout:     rs.Config.RevealTimeout,
		SettleTimeout:     c.SettleTimeout,
	}
	d.OurState = channel.NewChannelEndState(c.OurAddress, c.OurContractBalance,
		c.OurBalanceProof, mtree.NewMerkleTree(c.OurLeaves))
	d.PartenerState = channel.NewChannelEndState(c.PartnerAddress(),
		c.PartnerContractBalance,
		c.PartnerBalanceProof, mtree.NewMerkleTree(c.PartnerLeaves))
	d.ExternState = channel.NewChannelExternalState(rs.registerChannelForHashlock, tokenNetwork,
		c.ChannelIdentifier, rs.PrivateKey,
		rs.Chain.Client, rs.db, c.ClosedBlock,
		c.OurAddress, c.PartnerAddress())
	return d
}

//read a token network info from db
func (rs *RaidenService) registerTokenNetwork(tokenAddress, tokenNetworkAddress common.Address) (err error) {
	tokenNetwork, err := rs.Chain.TokenNetworkWithoutCheck(tokenNetworkAddress)
	edges, err := rs.db.GetAllNonParticipantChannel(tokenAddress)
	if err != nil {
		return
	}
	channels, err := rs.db.GetChannelList(tokenAddress, utils.EmptyAddress)
	var channelsDetails []*graph.ChannelDetails
	for _, ch := range channels {
		//跳过已经 settle 的 channel 加入没有任何意义.
		if ch.State == channeltype.StateSettled {
			continue
		}
		d := rs.channelSerilization2ChannelDetail(ch, tokenNetwork)
		channelsDetails = append(channelsDetails, d)
	}
	graph := graph.NewChannelGraph(rs.NodeAddress, tokenAddress, edges, channelsDetails)
	rs.TokenNetwork2Token[tokenNetworkAddress] = tokenAddress
	rs.Token2TokenNetwork[tokenAddress] = tokenNetworkAddress
	rs.Token2ChannelGraph[tokenAddress] = graph
	rs.Tokens2ConnectionManager[tokenAddress] = NewConnectionManager(rs, tokenAddress)
	return
}

/*
found new channel on blockchain when running...
*/
func (rs *RaidenService) registerNettingChannel(tokenNetworkAddress common.Address, partnerAddress common.Address) {
	tokenNetwork, err := rs.Chain.TokenNetwork(tokenNetworkAddress)
	if err != nil {
		log.Error(fmt.Sprintf("try to get tokenNetwork %s, err %s", tokenNetworkAddress.String(), err))
		return
	}
	tokenAddress := rs.TokenNetwork2Token[tokenNetworkAddress]
	if rs.getChannel(tokenAddress, partnerAddress) != nil {
		log.Error(fmt.Sprintf("receive new channel %s-%s,but this channel already exist, maybe a duplicate channel event", utils.APex2(tokenAddress), utils.APex2(partnerAddress)))
		return
	}
	detail, err := rs.getChannelDetail(tokenAddress, tokenNetwork, partnerAddress)
	if err != nil {
		log.Error(fmt.Sprintf("getChannelDetail err %s", err))
		return
	}
	graph := rs.getToken2ChannelGraph(tokenAddress)
	err = graph.AddChannel(detail)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = rs.db.NewChannel(channel.NewChannelSerialization(graph.ChannelAddress2Channel[detail.ChannelIdentifier]))
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
func (rs *RaidenService) directTransferAsync(tokenAddress, target common.Address, amount *big.Int) (result *utils.AsyncResult) {
	graph := rs.getToken2ChannelGraph(tokenAddress)
	directChannel := graph.GetPartenerAddress2Channel(target)
	result = utils.NewAsyncResult()
	if directChannel == nil || !directChannel.CanTransfer() || directChannel.Distributable().Cmp(amount) < 0 {
		result.Result <- errors.New("no available direct channel")
		return
	}
	tr, err := directChannel.CreateDirectTransfer(amount)
	if err != nil {
		result.Result <- err
		return
	}
	tr.Sign(rs.PrivateKey, tr)
	directChannel.RegisterTransfer(rs.GetBlockNumber(), tr)
	identifier := utils.NewRandomHash()
	//This should be set once the direct transfer is acknowledged
	transferSuccess := &transfer.EventTransferSentSuccess{
		LockSecretHash:    identifier,
		Amount:            amount,
		Target:            target,
		ChannelIdentifier: directChannel.ChannelIdentifier.ChannelIdentifier,
		Token:             tokenAddress,
	}
	result = rs.Protocol.SendAsync(directChannel.PartnerState.Address, tr)
	err = rs.StateMachineEventHandler.OnEvent(transferSuccess, nil)
	if err != nil {
		log.Error(fmt.Sprintf("dispatch transferSuccess err %s", err))
	}
	return
}

/*
mediated transfer for token swap
we must make sure that taker use the maker's secret.
and taker's lock expiration should be short than maker's todo(fix this)
*/
func (rs *RaidenService) startTakerMediatedTransfer(tokenAddress, target common.Address, amount *big.Int, lockSecretHash common.Hash, hashlock common.Hash, expiration int64) (result *utils.AsyncResult, stateManager *transfer.StateManager) {
	return rs.startMediatedTransferInternal(tokenAddress, target, amount, utils.BigInt0, lockSecretHash, hashlock, expiration)
}

/*
lauch a new mediated trasfer
Args:
 hashlock: caller can specify a hashlock or use empty ,when empty, will generate a random secret.
 expiration: caller can specify a valid blocknumber or 0, when 0 ,will calculate based on settle timeout of channel.
*/
func (rs *RaidenService) startMediatedTransferInternal(tokenAddress, target common.Address, amount *big.Int, fee *big.Int, lockSecretHash common.Hash, hashlock common.Hash, expiration int64) (result *utils.AsyncResult, stateManager *transfer.StateManager) {
	g := rs.getToken2ChannelGraph(tokenAddress)
	availableRoutes := g.GetBestRoutes(rs.Protocol, rs.NodeAddress, target, amount, graph.EmptyExlude, rs)
	result = utils.NewAsyncResult()
	if len(availableRoutes) <= 0 {
		result.Result <- errors.New("no available route")
		return
	}
	if rs.Config.IsMeshNetwork {
		result.Result <- errors.New("no mediated transfer on mesh only network")
		return
	}
	var secret common.Hash
	if lockSecretHash == utils.EmptyHash {
		lockSecretHash = utils.NewRandomHash()
	} else {
		secret = utils.NewRandomHash()
		lockSecretHash = utils.Sha3(secret[:])
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
	}
	/*
			  Issue #489

		        Raiden may fail after a state change using the random generator is
		        handled but right before the snapshot is taken. If that happens on
		        the next initialization when raiden is recovering and applying the
		        pending state changes a new secret will be generated and the
		        resulting events won't match, rs breaks the architecture model,
		        since it's assumed the re-execution of a state change will always
		        produce the same events.

		        TODO: Removed the secret generator from the InitiatorState and add
		        the secret into all state changes that require one, rs way the
		        secret will be serialized with the state change and the recovery will
		        use the same /random/ secret.
	*/
	initInitiator := &mediatedtransfer.ActionInitInitiatorStateChange{
		OurAddress:     rs.NodeAddress,
		Tranfer:        transferState,
		Routes:         routesState,
		BlockNumber:    rs.GetBlockNumber(),
		Secret:         secret,
		LockSecretHash: lockSecretHash,
		Db:             rs.db,
	}
	stateManager = transfer.NewStateManager(initiator.StateTransition, nil, initiator.NameInitiatorTransition, lockSecretHash, transferState.Token)
	smkey := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := rs.Transfer2StateManager[smkey]
	if manager != nil {
		panic(fmt.Sprintf("manager must be never exist"))
	}
	rs.Transfer2StateManager[smkey] = stateManager
	rs.Transfer2Result[smkey] = result
	rs.db.AddStateManager(stateManager)
	rs.StateMachineEventHandler.logAndDispatch(stateManager, initInitiator)
	return
}

/*
1. user start a mediated transfer
2. user start a maker mediated transfer
*/
func (rs *RaidenService) startMediatedTransfer(tokenAddress, target common.Address, amount *big.Int, fee *big.Int, lockSecretHash common.Hash) (result *utils.AsyncResult) {
	result, _ = rs.startMediatedTransferInternal(tokenAddress, target, amount, fee, lockSecretHash, utils.EmptyHash, 0)
	return
}

//receive a MediatedTransfer, i'm a hop node
func (rs *RaidenService) mediateMediatedTransfer(msg *encoding.MediatedTransfer, ch *channel.Channel) {
	tokenAddress := ch.TokenAddress
	smkey := utils.Sha3(msg.LockSecretHash[:], tokenAddress[:])
	stateManager := rs.Transfer2StateManager[smkey]
	if stateManager != nil {
		if stateManager.Name != mediator.NameMediatorTransition {
			log.Error(fmt.Sprintf("receive mediator transfer,but i'm not a mediator,msg=%s,stateManager=%s", msg, utils.StringInterface(stateManager, 3)))
			return
		}
		stateChange := &mediatedtransfer.MediatorReReceiveStateChange{
			Message: msg,
		}
		rs.StateMachineEventHandler.logAndDispatch(stateManager, stateChange)
	} else {
		/*
			第一次收到这个密码,
			todo 首先要判断这个密码是否是我声明放弃过的,如果是,就应该谨慎处理.
		*/
		if rs.db.IsLockSecretHashDisposed(msg.LockSecretHash) {
			//我都声明过这个锁我放弃了,但是仍然收到了这个锁,按道理锁是随机生成,不可能重复的.所以,这一定是攻击
			log.Error(fmt.Sprintf("receive a lock secret hash,and it's my annouce disposed. %s", msg.LockSecretHash.String()))
			//忽略,什么都不做
			return
		}
		amount := msg.PaymentAmount
		target := msg.Target
		exclude := graph.MakeExclude(msg.Sender, msg.Initiator)
		g := rs.getToken2ChannelGraph(ch.TokenAddress) //must exist
		avaiableRoutes := g.GetBestRoutes(rs.Protocol, rs.NodeAddress, target, amount, exclude, rs)
		fromChannel := g.GetPartenerAddress2Channel(msg.Sender)
		fromRoute := graph.Channel2RouteState(fromChannel, msg.Sender, amount, rs)
		ourAddress := rs.NodeAddress
		fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg, ch.TokenAddress)
		routesState := route.NewRoutesState(avaiableRoutes)
		blockNumber := rs.GetBlockNumber()
		initMediator := &mediatedtransfer.ActionInitMediatorStateChange{
			OurAddress:  ourAddress,
			FromTranfer: fromTransfer,
			Routes:      routesState,
			FromRoute:   fromRoute,
			BlockNumber: blockNumber,
			Message:     msg,
			Db:          rs.db,
		}
		stateManager = transfer.NewStateManager(mediator.StateTransition, nil, mediator.NameMediatorTransition, fromTransfer.LockSecretHash, fromTransfer.Token)
		rs.db.AddStateManager(stateManager)
		rs.Transfer2StateManager[smkey] = stateManager //for path A-B-C-F-B-D-E ,node B will have two StateManagers for one identifier
		rs.StateMachineEventHandler.logAndDispatch(stateManager, initMediator)
	}
}

//receive a MediatedTransfer, i'm the target
func (rs *RaidenService) targetMediatedTransfer(msg *encoding.MediatedTransfer, ch *channel.Channel) {
	smkey := utils.Sha3(msg.LockSecretHash[:], ch.TokenAddress[:])
	stateManager := rs.Transfer2StateManager[smkey]
	if stateManager != nil {
		if stateManager.Name != target.NameTargetTransition {
			log.Error(fmt.Sprintf("receive mediator transfer,but i'm not a target,msg=%s,stateManager=%s", msg, utils.StringInterface(stateManager, 3)))
			return
		}
		log.Error(fmt.Sprintf("receive mediator transfer msg=%s,duplicate? attack?,i'm a target,and has received mediator message. statemanager=%s",
			msg, utils.StringInterface(stateManager, 3)))
		return
	} else {
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
			Db:          rs.db,
		}
		stateManger := transfer.NewStateManager(target.StateTransiton, nil, target.NameTargetTransition, fromTransfer.LockSecretHash, fromTransfer.Token)
		rs.db.AddStateManager(stateManger)
		rs.Transfer2StateManager[smkey] = stateManager
		rs.StateMachineEventHandler.logAndDispatch(stateManger, initTarget)
	}
}

func (rs *RaidenService) startHealthCheckFor(address common.Address) {
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

func (rs *RaidenService) startNeighboursHealthCheck() {
	for _, g := range rs.Token2ChannelGraph {
		for addr := range g.PartenerAddress2Channel {
			rs.startHealthCheckFor(addr)
		}
	}
}
func (rs *RaidenService) startSubscribeNeighborStatus() error {
	mt, ok := rs.Transport.(*network.MixTransporter)
	if !ok {
		return fmt.Errorf("transport is not mix transpoter")
	}
	return mt.SubscribeNeighbor(rs.db)
}
func (rs *RaidenService) getToken2ChannelGraph(tokenAddress common.Address) (cg *graph.ChannelGraph) {
	cg = rs.Token2ChannelGraph[tokenAddress]
	if cg == nil {
		log.Error(fmt.Sprintf("%s token doesn't exist ", utils.APex(tokenAddress)))
	}
	return
}
func (rs *RaidenService) getChannelGraph(channelIdentifier common.Hash) (cg *graph.ChannelGraph) {
	ch := rs.getChannelWithAddr(channelIdentifier)
	if ch == nil {
		return
	}
	cg = rs.Token2ChannelGraph[ch.TokenAddress]
	if cg == nil {
		log.Error(fmt.Sprintf("%s token doesn't exist ", utils.APex(ch.TokenAddress)))
	}
	return
}
func (rs *RaidenService) getTokenForChannelIdentifier(channelidentifier common.Hash) (token common.Address) {
	ch := rs.getChannelWithAddr(channelidentifier)
	if ch == nil {
		return
	}
	return ch.TokenAddress
}

//only for test, should call findChannelByAddress
func (rs *RaidenService) getChannelWithAddr(channelAddr common.Hash) *channel.Channel {
	c, _ := rs.findChannelByAddress(channelAddr)
	return c
}

//for test
func (rs *RaidenService) getChannel(tokenAddr, partnerAddr common.Address) *channel.Channel {
	g := rs.getToken2ChannelGraph(tokenAddr)
	if g == nil {
		return nil
	}
	return g.GetPartenerAddress2Channel(partnerAddr)
}

/*
Process user's new channel request
*/
func (rs *RaidenService) newChannel(token, partner common.Address, settleTimeout int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	go func() {
		var err error
		defer rpanic.PanicRecover(fmt.Sprintf("newChannel token:%s,partner:%s", utils.APex(token), utils.APex(partner)))
		defer func() {
			result.Result <- err
			close(result.Result)
		}()
		tokenNetwork, err := rs.Chain.TokenNetwork(rs.Token2TokenNetwork[token])
		if err != nil {
			return
		}
		err = tokenNetwork.NewChannel(partner, settleTimeout)
		if err != nil {
			return
		}
		//defer write result
	}()
	return
}

/*
process user's deposit request
*/
func (rs *RaidenService) depositChannel(channelAddress common.Hash, amount *big.Int) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByAddress(channelAddress)
	if err != nil {
		result.Result <- err
		return
	}
	if c.State != channeltype.StateOpened {
		result.Result <- errors.New("channel can deposit only when at open state")
	}
	result = c.ExternState.Deposit(c.TokenAddress, amount)
	return
}

/*
process user's close or settle channel request
*/
func (rs *RaidenService) closeOrSettleChannel(channelAddress common.Hash, op string) (result *utils.AsyncResult) {
	result = utils.NewAsyncResult()
	c, err := rs.findChannelByAddress(channelAddress)
	if err != nil { //settled channel can be queried from db.
		result.Result <- errors.New("channel not exist")
		return
	}
	log.Trace(fmt.Sprintf("%s channel %s\n", op, utils.HPex(channelAddress)))
	if op == closeChannelReqName {
		result = c.Close()
	} else {
		result = c.Settle()
	}
	return
}

/*
process user's token swap maker request
save and restore todo?
*/
func (rs *RaidenService) tokenSwapMaker(tokenswap *TokenSwap) (result *utils.AsyncResult) {
	var hashlock common.Hash
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
		delete(rs.SecretRequestPredictorMap, hashlock) //old hashlock is invalid,just  remove
		return false
	}
	sentMtrHook = func(mtr *encoding.MediatedTransfer) (remove bool) {
		if mtr.LockSecretHash == tokenswap.LockSecretHash && rs.getTokenForChannelIdentifier(mtr.ChannelIdentifier) == tokenswap.FromToken && mtr.Target == tokenswap.ToNodeAddress && mtr.PaymentAmount.Cmp(tokenswap.FromAmount) == 0 {
			if hashlock != utils.EmptyHash {
				log.Info(fmt.Sprintf("tokenswap maker select new path ,because of different hash lock"))
				delete(rs.SecretRequestPredictorMap, hashlock) //old hashlock is invalid,just  remove
			}
			hashlock = mtr.LockSecretHash //hashlock may change when select new route path
			rs.SecretRequestPredictorMap[hashlock] = secretRequestHook
		}
		return false
	}
	receiveMtrHook = func(mtr *encoding.MediatedTransfer) (remove bool) {
		/*
			recevive taker's mediated transfer , the transfer must use argument of tokenswap and have the same hashlock
		*/
		if mtr.LockSecretHash == tokenswap.LockSecretHash && hashlock == mtr.LockSecretHash && rs.getTokenForChannelIdentifier(mtr.ChannelIdentifier) == tokenswap.ToToken && mtr.Target == tokenswap.FromNodeAddress && mtr.PaymentAmount.Cmp(tokenswap.ToAmount) == 0 {
			hasReceiveTakerMediatedTransfer = true
			delete(rs.SentMediatedTransferListenerMap, &sentMtrHook)
			return true
		}
		return false
	}
	rs.SentMediatedTransferListenerMap[&sentMtrHook] = true
	rs.ReceivedMediatedTrasnferListenerMap[&receiveMtrHook] = true
	result = rs.startMediatedTransfer(tokenswap.FromToken, tokenswap.ToNodeAddress, tokenswap.FromAmount, utils.BigInt0, tokenswap.LockSecretHash)
	return
}

/*
taker process token swap
taker's action is triggered by maker's mediated transfer.
*/
func (rs *RaidenService) messageTokenSwapTaker(msg *encoding.MediatedTransfer, tokenswap *TokenSwap) (remove bool) {
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
	takerExpiration := msg.Expiration - params.DefaultRevealTimeout
	result, stateManager := rs.startTakerMediatedTransfer(tokenswap.ToToken, tokenswap.FromNodeAddress, tokenswap.ToAmount, tokenswap.LockSecretHash, msg.LockSecretHash, takerExpiration)
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
func (rs *RaidenService) tokenSwapTaker(tokenswap *TokenSwap) (result *utils.AsyncResult) {
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

//all user's request
func (rs *RaidenService) handleReq(req *apiReq) {
	var result *utils.AsyncResult
	switch req.Name {
	case transferReqName: //mediated transfer only
		r := req.Req.(*transferReq)
		if r.IsDirectTransfer {
			result = rs.directTransferAsync(r.TokenAddress, r.Target, r.Amount)
		} else {
			result = rs.startMediatedTransfer(r.TokenAddress, r.Target, r.Amount, r.Fee, utils.EmptyHash)
		}
	case newChannelReqName:
		r := req.Req.(*newChannelReq)
		result = rs.newChannel(r.tokenAddress, r.partnerAddress, r.settleTimeout)
	case depositChannelReqName:
		r := req.Req.(*depositChannelReq)
		result = rs.depositChannel(r.addr, r.amount)
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
	default:
		panic("unkown req")
	}
	r := req
	r.result <- result
}

//recieve a ack from
func (rs *RaidenService) handleSentMessage(sentMessage *protocolMessage) {
	log.Trace(fmt.Sprintf("msg receive ack :%s", utils.StringInterface(sentMessage, 2)))
	if sentMessage.Message.Tag() != nil { //
		sentMessageTag := sentMessage.Message.Tag().(*transfer.MessageTag)
		if sentMessageTag.GetStateManager() != nil {
			mgr := sentMessageTag.GetStateManager()
			mgr.ManagerState = transfer.StateManagerSendMessageSuccesss
			sentMessageTag.SendingMessageComplete = true
			tx := rs.db.StartTx()
			_, ok := sentMessage.Message.(*encoding.UnLock)
			if ok {
				mgr.IsBalanceProofSent = true
				if mgr.Name == initiator.NameInitiatorTransition {
					mgr.ManagerState = transfer.StateManagerTransferComplete
				} else if mgr.Name == target.NameTargetTransition {

				} else if mgr.Name == mediator.NameMediatorTransition {
					/*
						how to detect a mediator node is finish or not?
							1. receive prev balanceproof
							2. balanceproof  send to next successfully
						//todo when refund?
					*/
					if mgr.IsBalanceProofSent && mgr.IsBalanceProofReceived {
						mgr.ManagerState = transfer.StateManagerTransferComplete
					}

				}
			}
			rs.db.UpdateStateManaer(mgr, tx)
			tx.Commit()
			rs.conditionQuit(fmt.Sprintf("%sRecevieAck", sentMessage.Message.Name()))
		} else if sentMessageTag.EchoHash != utils.EmptyHash {
			//log.Trace(fmt.Sprintf("reveal sent complete %s", utils.StringInterface(sentMessage.Message, 5)))
			rs.conditionQuit(fmt.Sprintf("%sRecevieAck", sentMessage.Message.Name()))
			switch msg := sentMessage.Message.(type) {
			case *encoding.RevealSecret:
				rs.db.UpdateSentRevealSecretComplete(sentMessageTag.EchoHash)
			default:
				log.Error(fmt.Sprintf("unknown message %s", utils.StringInterface(msg, 7)))
			}

		} else {
			panic(fmt.Sprintf("sent message state unknow :%s", utils.StringInterface(sentMessageTag, 2)))
		}
	} else {
		log.Error(fmt.Sprintf("message must have tag, only when make token swap %s", utils.StringInterface(sentMessage.Message, 3)))
	}
}

/*
GetNodeChargeFee implement of FeeCharger
*/
func (rs *RaidenService) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return rs.FeePolicy.GetNodeChargeFee(nodeAddress, tokenAddress, amount)
}

//SetFeePolicy set fee policy
func (rs *RaidenService) SetFeePolicy(feePolicy fee.Charger) {
	rs.FeePolicy = feePolicy
}

/*
for debug only,quit if eventName exactly match
*/
func (rs *RaidenService) conditionQuit(eventName string) {
	if strings.ToLower(eventName) == strings.ToLower(rs.Config.ConditionQuit.QuitEvent) {
		log.Error(fmt.Sprintf("quitevent=%s\n", eventName))
		debug.PrintStack()
		os.Exit(111)
	}
}

/*
GetDb return raiden's db
*/
func (rs *RaidenService) GetDb() *models.ModelDB {
	return rs.db
}

func (rs *RaidenService) handleEthRRCConnectionOK() {
	if !rs.ethInited {
		log.Info(fmt.Sprintf("eth connection ok, will reinit raiden"))
		rs.ethInited = true
		err := rs.AlarmTask.Start()
		if err != nil {
			log.Error(fmt.Sprintf("alarm task start err %s", err))
			n := rs.db.GetLatestBlockNumber()
			rs.BlockNumber.Store(n)
		} else {
			//must have a valid blocknumber before any transfer operation
			rs.BlockNumber.Store(rs.AlarmTask.LastBlockNumber)
		}
		/*
			events before lastHandledBlockNumber must have been processed, so we start from  lastHandledBlockNumber-1
		*/
		err = rs.BlockChainEvents.Start(rs.db.GetLatestBlockNumber())
		if err != nil {
			err = fmt.Errorf("Events listener error %v", err)
			return
		}
		/*
			  Registry registration must start *after* the alarm task, rs avoid
				 corner cases were the registry is queried in block A, a new block B
				 is mined, and the alarm starts polling at block C.
		*/
		rs.registerRegistry()
		err = rs.restoreChannel(false)
		if err != nil {
			log.Error(fmt.Sprintf("reinit restoreChannel err %s", err))
		}
	}
}
