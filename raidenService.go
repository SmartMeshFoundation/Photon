package smartraiden

import (
	"crypto/ecdsa"
	"errors"

	"fmt"

	"path/filepath"

	"time"

	"math/rand"

	"sync"

	"sync/atomic"

	"math/big"

	"strings"

	"os"

	"runtime/debug"

	"github.com/SmartMeshFoundation/SmartRaiden/blockchain"
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/fee"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/target"
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
	Chain              *rpc.BlockChainService
	Registry           *rpc.RegistryProxy
	RegistryAddress    common.Address
	PrivateKey         *ecdsa.PrivateKey
	Transport          network.Transporter
	Discovery          network.Discovery
	Config             *params.Config
	Protocol           *network.RaidenProtocol
	NodeAddress        common.Address
	Token2ChannelGraph map[common.Address]*network.ChannelGraph
	//Token2ConnectionsManager todo fix later
	//swapkey_to_tokenswap
	//swapkey_to_greenlettask
	Manager2Token            map[common.Address]common.Address
	Identifier2StateManagers map[uint64][]*transfer.StateManager
	Identifier2Results       map[uint64][]*network.AsyncResult
	SwapKey2TokenSwap        map[swapKey]*TokenSwap
	Tokens2ConnectionManager map[common.Address]*ConnectionManager //how to save and restore for token swap? todo fix it
	/*
				   This is a map from a hashlock to a list of channels, the same
			         hashlock can be used in more than one token (for tokenswaps), a
			         channel should be removed from this list only when the Lock is
			         released/withdrawn but not when the secret is registered.
		TODO remove this,this design is very weird
	*/
	Token2Hashlock2Channels  map[common.Address]map[common.Hash][]*channel.Channel //Multithread
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
	RoutesTask                  *routesTask
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
}

//NewRaidenService create raiden service
func NewRaidenService(chain *rpc.BlockChainService, privateKey *ecdsa.PrivateKey, transport network.Transporter,
	discover network.DiscoveryInterface, config *params.Config) (srv *RaidenService) {
	if config.SettleTimeout < params.NettingChannelSettleTimeoutMin || config.SettleTimeout > params.NettingChannelSettleTimeoutMax {
		log.Error(fmt.Sprintf("settle timeout must be in range %d-%d",
			params.NettingChannelSettleTimeoutMin, params.NettingChannelSettleTimeoutMax))
		utils.SystemExit(1)
	}
	srv = &RaidenService{
		Chain:                               chain,
		Registry:                            chain.Registry(chain.RegistryAddress),
		RegistryAddress:                     chain.RegistryAddress,
		PrivateKey:                          privateKey,
		Config:                              config,
		NodeAddress:                         crypto.PubkeyToAddress(privateKey.PublicKey),
		Token2ChannelGraph:                  make(map[common.Address]*network.ChannelGraph),
		Manager2Token:                       make(map[common.Address]common.Address),
		Identifier2StateManagers:            make(map[uint64][]*transfer.StateManager),
		Identifier2Results:                  make(map[uint64][]*network.AsyncResult),
		Token2Hashlock2Channels:             make(map[common.Address]map[common.Hash][]*channel.Channel),
		SwapKey2TokenSwap:                   make(map[swapKey]*TokenSwap),
		Tokens2ConnectionManager:            make(map[common.Address]*ConnectionManager),
		AlarmTask:                           blockchain.NewAlarmTask(chain.Client),
		BlockChainEvents:                    blockchain.NewBlockChainEvents(chain.Client, chain.RegistryAddress),
		BlockNumberChan:                     make(chan int64, 1),
		UserReqChan:                         make(chan *apiReq, 10),
		BlockNumber:                         new(atomic.Value),
		ProtocolMessageSendComplete:         make(chan *protocolMessage, 10),
		SecretRequestPredictorMap:           make(map[common.Hash]SecretRequestPredictor),
		RevealSecretListenerMap:             make(map[common.Hash]RevealSecretListener),
		ReceivedMediatedTrasnferListenerMap: make(map[*ReceivedMediatedTrasnferListener]bool),
		SentMediatedTransferListenerMap:     make(map[*SentMediatedTransferListener]bool),
		FeePolicy:                           &ConstantFeePolicy{},
		HealthCheckMap:                      make(map[common.Address]bool),
	}
	var err error
	srv.MessageHandler = newRaidenMessageHandler(srv)
	srv.StateMachineEventHandler = newStateMachineEventHandler(srv)
	srv.Protocol = network.NewRaidenProtocol(transport, discover, privateKey, srv)
	srv.db, err = models.OpenDb(config.DataBasePath)
	if err != nil {
		log.Error("open db error")
		utils.SystemExit(1)
	}
	srv.Protocol.SetReceivedMessageSaver(NewAckHelper(srv.db))
	/*
		only one instance for one data directory
	*/
	srv.FileLocker = flock.NewFlock(config.DataBasePath + ".flock.Lock")
	locked, err := srv.FileLocker.TryLock()
	if err != nil || !locked {
		log.Error(fmt.Sprintf("another instance already running at %s", config.DataBasePath))
		utils.SystemExit(1)
	}
	srv.SnapshortDir = filepath.Join(config.DataBasePath)
	err = discover.Register(srv.NodeAddress, srv.Config.ExternIP, srv.Config.ExternPort)
	if err != nil {
		log.Error(fmt.Sprintf("register discover endpoint error:%s", err))
		utils.SystemExit(1)
	}
	log.Info("node discovery register complete...")
	//srv.Start()
	//start routes detect task
	srv.RoutesTask = newRoutesTask(srv.Protocol, srv.Protocol)
	return srv
}

// Start the node.
func (rs *RaidenService) Start() {
	lastHandledBlockNumber := rs.db.GetLatestBlockNumber()
	rs.AlarmTask.Start()
	rs.RoutesTask.start()
	//must have a valid blocknumber before any transfer operation
	rs.BlockNumber.Store(rs.AlarmTask.LastBlockNumber)
	rs.AlarmTask.RegisterCallback(func(number int64) error {
		rs.db.SaveLatestBlockNumber(number)
		return rs.setBlockNumber(number)
	})
	/*
		events before lastHandledBlockNumber must have been processed, so we start from  lastHandledBlockNumber-1
	*/
	err := rs.BlockChainEvents.Start(lastHandledBlockNumber)
	if err != nil {
		log.Error(fmt.Sprintf("Events listener error %v", err))
		utils.SystemExit(1)
	}
	/*
			  Registry registration must start *after* the alarm task, rs avoid
		         corner cases were the registry is queried in block A, a new block B
		         is mined, and the alarm starts polling at block C.
	*/
	rs.registerRegistry()
	err = rs.restoreSnapshot()
	if err != nil {
		log.Error(fmt.Sprintf("restore from snapshot error : %v\n you can delete all the database %s to run. but all your trade will lost!!", err, rs.Config.DataBasePath))
		utils.SystemExit(1)
	}
	rs.Protocol.Start()
	rs.startNeighboursHealthCheck()
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
}

//Stop the node.
func (rs *RaidenService) Stop() {
	log.Info("raiden service stop...")
	rs.AlarmTask.Stop()
	rs.RoutesTask.stop()
	rs.Protocol.StopAndWait()
	rs.BlockChainEvents.Stop()
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
	var routestask *routesToDetect
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
					log.Error("stateMachineEventHandler.OnBlockchainStateChange", err)
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
			//before send a transfer, we would better detect if neighbors are online or not.
		case routestask, ok = <-rs.RoutesTask.TaskResult:
			if ok {
				rs.handleRoutesTask(routestask)
			} else {
				log.Info("routesTask.TaskResult closed")
				return
			}
		}
	}
}

//for init
func (rs *RaidenService) registerRegistry() {
	mgrs, err := rs.Chain.GetAllChannelManagers()
	if err != nil {
		log.Error(fmt.Sprintf("registerRegistry err:%s", err))
		utils.SystemExit(1)
	}
	for _, mgr := range mgrs {
		err = rs.registerChannelManager(mgr.Address)
		if err != nil {
			log.Error(fmt.Sprintf("registerChannelManager err:%s", err))
			utils.SystemExit(1)
		}
	}
}

/*
quering my channel details on blockchain
i'm one of the channel participants
*/
func (rs *RaidenService) getChannelDetail(tokenAddress common.Address, proxy *rpc.NettingChannelContractProxy) *network.ChannelDetails {
	addr1, b1, addr2, b2, _ := proxy.AddressAndBalance()
	var ourAddr, partnerAddr common.Address
	var ourBalance, partnerBalance *big.Int
	if addr1 == rs.NodeAddress {
		ourAddr = addr1
		partnerAddr = addr2
		ourBalance = b1
		partnerBalance = b2
	} else {
		ourAddr = addr2
		partnerAddr = addr1
		ourBalance = b2
		partnerBalance = b1
	}
	ourState := channel.NewChannelEndState(ourAddr, ourBalance, nil, transfer.EmptyMerkleTreeState)
	partenerState := channel.NewChannelEndState(partnerAddr, partnerBalance, nil, transfer.EmptyMerkleTreeState)
	channelAddress := proxy.Address
	registerChannelForHashlock := func(channel *channel.Channel, hashlock common.Hash) {
		rs.registerChannelForHashlock(tokenAddress, channel, hashlock)
	}
	externState := channel.NewChannelExternalState(registerChannelForHashlock, proxy, channelAddress, rs.Chain, rs.db)
	channelDetail := &network.ChannelDetails{
		ChannelAddress:    channelAddress,
		OurState:          ourState,
		PartenerState:     partenerState,
		ExternState:       externState,
		BlockChainService: rs.Chain,
		RevealTimeout:     rs.Config.RevealTimeout,
	}
	channelDetail.SettleTimeout, _ = externState.NettingChannel.SettleTimeout()
	return channelDetail
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

func (rs *RaidenService) findChannelByAddress(nettingChannelAddress common.Address) (*channel.Channel, error) {
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
       and ignoring the tokens. Useful for refund transfer, split transfer,
       and token swaps.

       Raises:
           TypeError: If secret is unicode data.
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
					utils.HPex(secret), utils.APex2(ch.MyAddress), err))
			}
			rs.conditionQuit("BeforeSendRevealSecret")
			rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
			//The protocol ignores duplicated messages.
			//make sure not send the same instance multi times.
			rs.sendAsync(ch.PartnerState.Address, encoding.CloneRevealSecret(revealSecretMessage))
		}
	}
}

func (rs *RaidenService) registerChannelForHashlock(tokenAddress common.Address,
	netchannel *channel.Channel, hashlock common.Hash) {
	channelsRegistered := rs.Token2Hashlock2Channels[tokenAddress][hashlock]
	found := false
	for _, c := range channelsRegistered {
		//To determine whether the two channel objects are equal, we simply use the address to identify.
		if c.ExternState.ChannelAddress == netchannel.ExternState.ChannelAddress {
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

/*
Unlock/Witdraws locks, register the secret, and send Secret
       messages as necessary.

       This function will:
           - Unlock the locks created by this node and send a Secret message to
           the corresponding partner so that she can withdraw the token.
           - Withdraw the Lock from sender.
           - Register the secret for the locks received and reveal the secret
           to the senders


       Note:
           The channel needs to be registered with
           `raiden.register_channel_for_hashlock`.
//todo 需要再次确认, refund 处理流程有无问题,以及相关细节必须审核.
*/
func (rs *RaidenService) handleSecret(identifier uint64, tokenAddress common.Address, secret common.Hash,
	partnerSecretMessage *encoding.Secret, hashlock common.Hash) (err error) {
	/*
	   handling the secret needs to:
	         - unlock the token for all `forward_channel` (the current one
	           and the ones that failed with a refund)
	         - send a message to each of the forward nodes allowing them
	           to withdraw the token
	         - register the secret for the `originating_channel` so that a
	           proof can be made, if necessary
	         - reveal the secret to the `sender` node (otherwise we
	           cannot withdraw the token)
	*/
	channelsList := rs.Token2Hashlock2Channels[tokenAddress][hashlock]
	var channelsToRemove []*channel.Channel
	revealSecretMessage := encoding.NewRevealSecret(secret)
	revealSecretMessage.Sign(rs.PrivateKey, revealSecretMessage)
	type MsgToSend struct {
		receiver common.Address
		msg      encoding.SignedMessager
	}
	var messagesToSend []*MsgToSend
	log.Trace(fmt.Sprintf("channelsList for %s =%#v", utils.HPex(hashlock), channelsList))
	for _, ch := range channelsList { //Dealing with reused nodes in indirect transactions.
		//unlock a pending Lock
		log.Trace(fmt.Sprintf("process channel %s-%s", utils.APex2(ch.OurState.Address), utils.APex2((ch.PartnerState.Address))))
		if ch.OurState.IsKnown(hashlock) {
			var secretMsg *encoding.Secret
			secretMsg, err = ch.CreateSecret(identifier, secret)
			if err != nil {
				return err
			}
			secretMsg.Sign(rs.PrivateKey, secretMsg)
			//balance proof,Complete rs transaction and receive repeated revealsecret,But secretmsg will only be sent once.
			err = ch.RegisterTransfer(rs.GetBlockNumber(), secretMsg)
			if err != nil {
				return
			}
			messagesToSend = append(messagesToSend, &MsgToSend{ch.PartnerState.Address, secretMsg})
			channelsToRemove = append(channelsToRemove, ch)
		} else if ch.PartnerState.IsKnown(hashlock) {
			//withdraw a pending Lock
			if partnerSecretMessage != nil {
				msg := partnerSecretMessage
				isBalanceProof := msg.Sender == ch.PartnerState.Address && msg.Channel == ch.MyAddress
				if isBalanceProof {
					err = ch.RegisterTransfer(rs.GetBlockNumber(), msg)
					if err != nil {
						return
					}
					channelsToRemove = append(channelsToRemove, ch)
				} else {
					err = ch.RegisterSecret(secret)
					if err != nil {
						return
					}
					rs.conditionQuit("BeforeSendRevealSecret")
					rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
					messagesToSend = append(messagesToSend, &MsgToSend{ch.PartnerState.Address, encoding.CloneRevealSecret(revealSecretMessage)})
				}
			} else {
				err = ch.RegisterSecret(secret)
				if err != nil {
					return
				}
				rs.conditionQuit("BeforeSendRevealSecret")
				rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
				messagesToSend = append(messagesToSend, &MsgToSend{ch.PartnerState.Address, encoding.CloneRevealSecret(revealSecretMessage)})
			}
		} else {
			/*
				todo reimplement handleSecret
			*/
			log.Warn("Channel is registered for a given Lock but the Lock is not contained in it. can be ignored when I'm a mediated node")
		}

	}

	for _, ch := range channelsToRemove {
		//channels_list.remove(channel)
		for k, ch2 := range channelsList {
			if ch2 == ch {
				//to remove
				channelsList = append(channelsList[:k], channelsList[k+1:]...)
				break
			}
		}
	}
	if len(channelsList) == 0 {
		delete(rs.Token2Hashlock2Channels[tokenAddress], hashlock)
	} else {
		rs.Token2Hashlock2Channels[tokenAddress][hashlock] = channelsList
	}
	// send the messages last to avoid races
	for _, msg := range messagesToSend {
		err = rs.sendAsync(msg.receiver, msg.msg)
		if err != nil {
			return
		}
	}
	return
}

func (rs *RaidenService) channelManagerIsRegistered(manager common.Address) bool {
	_, ok := rs.Manager2Token[manager]
	return ok
}
func (rs *RaidenService) registerChannelManager(managerAddress common.Address) (err error) {
	manager := rs.Chain.Manager(managerAddress)
	channels, err := manager.NettingChannelByAddress(rs.NodeAddress)
	if err != nil {
		return
	}
	tokenAddress, _ := manager.TokenAddress()
	edgeList, _ := manager.GetChannelsParticipants()
	var channelsDetails []*network.ChannelDetails
	for _, ch := range channels {
		d := rs.getChannelDetail(tokenAddress, ch)
		channelsDetails = append(channelsDetails, d)
	}
	graph := network.NewChannelGraph(rs.NodeAddress, managerAddress, tokenAddress, edgeList, channelsDetails)
	rs.Manager2Token[managerAddress] = tokenAddress
	rs.Token2ChannelGraph[tokenAddress] = graph
	rs.Tokens2ConnectionManager[tokenAddress] = NewConnectionManager(rs, tokenAddress)
	//new token, save to db
	err = rs.db.AddToken(tokenAddress, managerAddress)
	if err != nil {
		log.Error(err.Error())
	}
	err = rs.db.UpdateTokenNodes(tokenAddress, graph.AllNodes())
	if err != nil {
		log.Error(err.Error())
	}
	// we need restore channel status from database after restart...
	//for _, c := range graph.ChannelAddress2Channel {
	//	err = rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(c))
	//	if err != nil {
	//		log.Info(err.Error())
	//	}
	//}

	return
}

/*
found new channel on blockchain when running...
*/
func (rs *RaidenService) registerNettingChannel(tokenAddress, channelAddress common.Address) {
	nettingChannel, err := rs.Chain.NettingChannel(channelAddress)
	if err != nil {
		log.Error("try to registerNettingChannel not exist channel %s", channelAddress)
	}
	detail := rs.getChannelDetail(tokenAddress, nettingChannel)
	graph := rs.getToken2ChannelGraph(tokenAddress)
	err = graph.AddChannel(detail)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = rs.db.NewChannel(channel.NewChannelSerialization(graph.ChannelAddress2Channel[channelAddress]))
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}
func (rs *RaidenService) connectionManagerForToken(tokenAddress common.Address) (*ConnectionManager, error) {
	mgr, ok := rs.Tokens2ConnectionManager[tokenAddress]
	if ok {
		return mgr, nil
	}
	return nil, rerr.InvalidAddress(fmt.Sprintf("token %s is not registered", utils.APex(tokenAddress)))
}
func (rs *RaidenService) leaveAllTokenNetworksAsync() *network.AsyncResult {
	var leaveResults []*network.AsyncResult
	for token := range rs.Token2ChannelGraph {
		mgr, err := rs.connectionManagerForToken(token)
		if err != nil {
			log.Error(fmt.Sprintf("connectionManagerForToken %s err %s", utils.APex(token), err))
		}
		if mgr != nil {
			leaveResults = append(leaveResults, mgr.LeaveAsync())
		}
	}
	return waitGroupAsyncResult(leaveResults)
}
func waitGroupAsyncResult(results []*network.AsyncResult) *network.AsyncResult {
	totalResult := network.NewAsyncResult()
	wg := sync.WaitGroup{}
	wg.Add(len(results))
	for i := range results {
		go func(i int) {
			<-results[i].Result
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		totalResult.Result <- nil
		close(totalResult.Result)
	}()
	return totalResult
}
func (rs *RaidenService) closeAndSettle() {
	log.Info("raiden will close and settle all channels now")
	var Mgrs []*ConnectionManager
	for t := range rs.Token2ChannelGraph {
		mgr, _ := rs.connectionManagerForToken(t)
		if mgr != nil {
			Mgrs = append(Mgrs, mgr)
		}
	}
	blocksToWait := func() int64 {
		var max int64
		for _, mgr := range Mgrs {
			if max < mgr.minSettleBlocks() {
				max = mgr.minSettleBlocks()
			}
		}
		return max
	}
	var AllChannels []*channel.Channel
	for _, mgr := range Mgrs {
		for _, c := range mgr.openChannels() {
			ch, err := rs.findChannelByAddress(c.ChannelAddress)
			if err != nil {
				panic(fmt.Sprintf("channel %s must exist", utils.APex(c.ChannelAddress)))
			}
			AllChannels = append(AllChannels, ch)
		}
	}
	leavingResult := rs.leaveAllTokenNetworksAsync()
	//using the un-cached block number here
	lastBlock := rs.GetBlockNumber()
	earliestSettlement := lastBlock + blocksToWait()
	/*
			    TODO: estimate and set a `timeout` parameter in seconds
		     based on connection_manager.min_settle_blocks and an average
		     blocktime from the past
	*/
	currentBlock := lastBlock
	for currentBlock < earliestSettlement {
		time.Sleep(time.Second * 10)
		lastBlock := rs.GetBlockNumber()
		if lastBlock != currentBlock {
			currentBlock = lastBlock
			waitBlocksLeft := blocksToWait()
			notSettled := 0
			for _, c := range AllChannels {
				if c.State() != transfer.ChannelStateSettled {
					notSettled++
				}
			}
			if notSettled == 0 {
				log.Debug("nothing left to settle")
				break
			}
			log.Info(fmt.Sprintf("waiting at least %d more blocks for %d channels not yet settled", waitBlocksLeft, notSettled))
		}
		//why  leaving_greenlet.wait
		timeoutch := time.After(time.Second * time.Duration(blocksToWait()))
		select {
		case <-timeoutch:
		case <-leavingResult.Result:
		}
	}
	for _, c := range AllChannels {
		if c.State() != transfer.ChannelStateSettled {
			log.Error("channels were not settled:", utils.APex(c.MyAddress))
		}
	}
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
func (rs *RaidenService) directTransferAsync(tokenAddress, target common.Address, amount *big.Int, identifier uint64) (result *network.AsyncResult) {
	graph := rs.getToken2ChannelGraph(tokenAddress)
	directChannel := graph.GetPartenerAddress2Channel(target)
	result = network.NewAsyncResult()
	if directChannel == nil || !directChannel.CanTransfer() || directChannel.Distributable().Cmp(amount) < 0 {
		result.Result <- errors.New("no available direct channel")
		return
	}
	tr, err := directChannel.CreateDirectTransfer(amount, identifier)
	if err != nil {
		result.Result <- err
		return
	}
	tr.Sign(rs.PrivateKey, tr)
	directChannel.RegisterTransfer(rs.GetBlockNumber(), tr)
	directTransferStateChange := &transfer.ActionTransferDirectStateChange{
		Identifier:   identifier,
		Amount:       amount,
		TokenAddress: tokenAddress,
		NodeAddress:  directChannel.PartnerState.Address,
	}
	// TODO: add the transfer sent event
	stateChangeID, err := rs.db.LogStateChange(directTransferStateChange)
	if err != nil {
		log.Error(fmt.Sprintf("LogStateChange %s err %s", utils.StringInterface(directTransferStateChange, 2), err))
	}
	//This should be set once the direct transfer is acknowledged
	transferSuccess := transfer.EventTransferSentSuccess{
		Identifier: identifier,
		Amount:     amount,
		Target:     target,
	}
	rs.db.LogEvents(stateChangeID, []transfer.Event{transferSuccess}, rs.GetBlockNumber())
	result = rs.Protocol.SendAsync(directChannel.PartnerState.Address, tr)
	return
}

/*
mediated transfer for token swap
we must make sure that taker use the maker's secret.
and taker's lock expiration should be short than maker's todo(fix this)
*/
func (rs *RaidenService) startTakerMediatedTransfer(tokenAddress, target common.Address, amount *big.Int, identifier uint64, hashlock common.Hash, expiration int64) (result *network.AsyncResult, stateManager *transfer.StateManager) {
	return rs.startMediatedTransferInternal(tokenAddress, target, amount, utils.BigInt0, identifier, hashlock, expiration)
}

/*
lauch a new mediated trasfer
Args:
 hashlock: caller can specify a hashlock or use empty ,when empty, will generate a random secret.
 expiration: caller can specify a valid blocknumber or 0, when 0 ,will calculate based on settle timeout of channel.
*/
func (rs *RaidenService) startMediatedTransferInternal(tokenAddress, target common.Address, amount *big.Int, fee *big.Int, identifier uint64, hashlock common.Hash, expiration int64) (result *network.AsyncResult, stateManager *transfer.StateManager) {
	graph := rs.getToken2ChannelGraph(tokenAddress)
	availableRoutes := graph.GetBestRoutes(rs.Protocol, rs.NodeAddress, target, amount, utils.EmptyAddress, rs)
	result = network.NewAsyncResult()
	result.Tag = target //tell the difference when token swap
	if len(availableRoutes) <= 0 {
		result.Result <- errors.New("no available route")
		return
	}
	if identifier == 0 {
		identifier = rand.New(utils.RandSrc).Uint64()
	}
	/*
		when user specified fee, for test or other purpose.
	*/
	if fee.Cmp(utils.BigInt0) > 0 {
		for _, r := range availableRoutes {
			r.TotalFee = fee //use the user's fee to replace algorithm's
		}
	}
	routesState := transfer.NewRoutesState(availableRoutes)
	transferState := &mediatedtransfer.LockedTransferState{
		Identifier:   identifier,
		TargetAmount: new(big.Int).Set(amount),
		Amount:       new(big.Int).Set(amount),
		Token:        tokenAddress,
		Initiator:    rs.NodeAddress,
		Target:       target,
		Expiration:   expiration,
		Hashlock:     utils.EmptyHash,
		Secret:       utils.EmptyHash,
		Fee:          utils.BigInt0,
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
		OurAddress:  rs.NodeAddress,
		Tranfer:     transferState,
		Routes:      routesState,
		BlockNumber: rs.GetBlockNumber(),
		Db:          rs.db,
	}
	if hashlock == utils.EmptyHash {
		initInitiator.RandomGenerator = utils.RandomSecretGenerator
	} else {
		initInitiator.RandomGenerator = utils.NewSepecifiedSecretGenerator(hashlock)
	}
	stateManager = transfer.NewStateManager(initiator.StateTransition, nil, initiator.NameInitiatorTransition, transferState.Identifier, transferState.Token)
	/*
		  TODO: implement the network timeout raiden.config['msg_timeout'] and
			cancel the current transfer if it hapens (issue #374)
	*/
	/*
		first register the rs transfer id, otherwise error may occur imediatelly
	*/
	mgrs := rs.Identifier2StateManagers[identifier]
	mgrs = append(mgrs, stateManager)
	rs.Identifier2StateManagers[identifier] = mgrs
	results := rs.Identifier2Results[identifier]
	results = append(results, result)
	rs.Identifier2Results[identifier] = results
	rs.db.AddStateManager(stateManager)
	//ping before send transfer
	rs.RoutesTask.NewTask <- &routesToDetect{
		RoutesState:     initInitiator.Routes,
		StateManager:    stateManager,
		InitStateChange: initInitiator,
	}
	//rs.stateMachineEventHandler.logAndDispatch(stateManager, initInitiator)
	return
}

/*
1. user start a mediated transfer
2. user start a maker mediated transfer
*/
func (rs *RaidenService) startMediatedTransfer(tokenAddress, target common.Address, amount *big.Int, fee *big.Int, identifier uint64) (result *network.AsyncResult) {
	result, _ = rs.startMediatedTransferInternal(tokenAddress, target, amount, fee, identifier, utils.EmptyHash, 0)
	return
}

//receive a MediatedTransfer, i'm a hop node
func (rs *RaidenService) mediateMediatedTransfer(msg *encoding.MediatedTransfer) {
	amount := msg.Amount
	target := msg.Target
	token := msg.Token
	graph := rs.getToken2ChannelGraph(token)
	avaiableRoutes := graph.GetBestRoutes(rs.Protocol, rs.NodeAddress, target, amount, msg.Sender, rs)
	fromChannel := graph.GetPartenerAddress2Channel(msg.Sender)
	fromRoute := network.Channel2RouteState(fromChannel, msg.Sender, amount, rs)
	ourAddress := rs.NodeAddress
	fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg)
	routesState := transfer.NewRoutesState(avaiableRoutes)
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
	stateManager := transfer.NewStateManager(mediator.StateTransition, nil, mediator.NameMediatorTransition, fromTransfer.Identifier, fromTransfer.Token)
	rs.db.AddStateManager(stateManager)
	mgrs := rs.Identifier2StateManagers[msg.Identifier]
	mgrs = append(mgrs, stateManager)
	rs.Identifier2StateManagers[msg.Identifier] = mgrs //for path A-B-C-F-B-D-E ,node B will have two StateManagers for one identifier
	//ping before send transfer
	rs.RoutesTask.NewTask <- &routesToDetect{
		RoutesState:     initMediator.Routes,
		StateManager:    stateManager,
		InitStateChange: initMediator,
	}
	//rs.stateMachineEventHandler.logAndDispatch(stateManager, initMediator)
}

//receive a MediatedTransfer, i'm the target
func (rs *RaidenService) targetMediatedTransfer(msg *encoding.MediatedTransfer) {
	graph := rs.getToken2ChannelGraph(msg.Token)
	fromChannel := graph.GetPartenerAddress2Channel(msg.Sender)
	fromRoute := network.Channel2RouteState(fromChannel, msg.Sender, msg.Amount, rs)
	fromTransfer := mediatedtransfer.LockedTransferFromMessage(msg)
	initTarget := &mediatedtransfer.ActionInitTargetStateChange{
		OurAddress:  rs.NodeAddress,
		FromRoute:   fromRoute,
		FromTranfer: fromTransfer,
		BlockNumber: rs.GetBlockNumber(),
		Message:     msg,
		Db:          rs.db,
	}
	stateManger := transfer.NewStateManager(target.StateTransiton, nil, target.NameTargetTransition, fromTransfer.Identifier, fromTransfer.Token)
	rs.db.AddStateManager(stateManger)
	identifier := msg.Identifier
	mgrs := rs.Identifier2StateManagers[identifier]
	mgrs = append(mgrs, stateManger)
	rs.Identifier2StateManagers[identifier] = mgrs

	rs.StateMachineEventHandler.logAndDispatch(stateManger, initTarget)
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

func (rs *RaidenService) getToken2ChannelGraph(tokenAddress common.Address) (cg *network.ChannelGraph) {
	cg = rs.Token2ChannelGraph[tokenAddress]
	if cg == nil {
		log.Error(fmt.Sprintf("%s token doesn't exist ", utils.APex(tokenAddress)))
	}
	return
}

//only for test, should call findChannelByAddress
func (rs *RaidenService) getChannelWithAddr(channelAddr common.Address) *channel.Channel {
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
func (rs *RaidenService) newChannel(token, partner common.Address, settleTimeout int) (result *network.AsyncResult) {
	result = network.NewAsyncResult()
	go func() {
		var err error
		defer func() {
			result.Result <- err
			close(result.Result)
		}()
		chMgrAddr, err := rs.Registry.ChannelManagerByToken(token)
		if err != nil {
			return
		}
		chMgr := rs.Chain.Manager(chMgrAddr)
		_, err = chMgr.NewChannel(partner, settleTimeout)
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
func (rs *RaidenService) depositChannel(channelAddress common.Address, amount *big.Int) (result *network.AsyncResult) {
	result = network.NewAsyncResult()
	c, err := rs.findChannelByAddress(channelAddress)
	if err != nil {
		result.Result <- err
		return
	}
	go func() {
		err := c.ExternState.Deposit(amount)
		result.Result <- err
		close(result.Result)
	}()
	return
}

/*
process user's close or settle channel request
*/
func (rs *RaidenService) closeOrSettleChannel(channelAddress common.Address, op string) (result *network.AsyncResult) {
	result = network.NewAsyncResult()
	c, err := rs.findChannelByAddress(channelAddress)
	if err != nil { //settled channel can be queried from db.
		result.Result <- errors.New("channel not exist")
		return
	}
	log.Trace(fmt.Sprintf("%s channel %s\n", op, utils.APex(channelAddress)))
	go func() {
		var err error
		c2, _ := rs.db.GetChannelByAddress(c.MyAddress)
		proof := c2.PartnerBalanceProof
		if op == closeChannelReqName {
			err = c.ExternState.Close(proof)
		} else {
			err = c.ExternState.Settle()
		}
		log.Trace(fmt.Sprintf("%s channel finished err %v", op, err))
		result.Result <- err
	}()
	return
}

/*
process user's token swap maker request
save and restore todo?
*/
func (rs *RaidenService) tokenSwapMaker(tokenswap *TokenSwap) (result *network.AsyncResult) {
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
		if mtr.Identifier == tokenswap.Identifier && mtr.Token == tokenswap.FromToken && mtr.Target == tokenswap.ToNodeAddress && mtr.Amount.Cmp(tokenswap.FromAmount) == 0 {
			if hashlock != utils.EmptyHash {
				log.Info(fmt.Sprintf("tokenswap maker select new path ,because of different hash lock"))
				delete(rs.SecretRequestPredictorMap, hashlock) //old hashlock is invalid,just  remove
			}
			hashlock = mtr.HashLock //hashlock may change when select new route path
			rs.SecretRequestPredictorMap[hashlock] = secretRequestHook
		}
		return false
	}
	receiveMtrHook = func(mtr *encoding.MediatedTransfer) (remove bool) {
		/*
			recevive taker's mediated transfer , the transfer must use argument of tokenswap and have the same hashlock
		*/
		if mtr.Identifier == tokenswap.Identifier && hashlock == mtr.HashLock && mtr.Token == tokenswap.ToToken && mtr.Target == tokenswap.FromNodeAddress && mtr.Amount.Cmp(tokenswap.ToAmount) == 0 {
			hasReceiveTakerMediatedTransfer = true
			delete(rs.SentMediatedTransferListenerMap, &sentMtrHook)
			return true
		}
		return false
	}
	rs.SentMediatedTransferListenerMap[&sentMtrHook] = true
	rs.ReceivedMediatedTrasnferListenerMap[&receiveMtrHook] = true
	result = rs.startMediatedTransfer(tokenswap.FromToken, tokenswap.ToNodeAddress, tokenswap.FromAmount, utils.BigInt0, tokenswap.Identifier)
	return
}

/*
taker process token swap
taker's action is triggered by maker's mediated transfer.
*/
func (rs *RaidenService) messageTokenSwapTaker(msg *encoding.MediatedTransfer, tokenswap *TokenSwap) (remove bool) {
	var hashlock = msg.HashLock
	var hasReceiveRevealSecret bool
	var stateManager *transfer.StateManager
	if msg.Identifier != tokenswap.Identifier || msg.Amount.Cmp(tokenswap.FromAmount) != 0 || msg.Initiator != tokenswap.FromNodeAddress || msg.Token != tokenswap.FromToken || msg.Target != tokenswap.ToNodeAddress {
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
		if msg.HashLock() != hashlock {
			return false
		}
		state := stateManager.CurrentState
		initState, ok := state.(*mediatedtransfer.InitiatorState)
		if !ok {
			panic(fmt.Sprintf("must be a InitiatorState"))
		}
		if initState.Transfer.Hashlock != msg.HashLock() {
			panic(fmt.Sprintf("hashlock must be same , state lock=%s,msg lock=%s", utils.HPex(initState.Transfer.Hashlock), utils.HPex(msg.HashLock())))
		}
		initState.Transfer.Secret = msg.Secret
		hasReceiveRevealSecret = true
		delete(rs.SecretRequestPredictorMap, hashlock)
		return true
	}
	/*
		taker's Expiration must be smaller than maker's ,
		taker and maker may have direct channels on these two tokens.
	*/
	takerExpiration := msg.Expiration - params.DefaultRevealTimeout
	result, stateManager := rs.startTakerMediatedTransfer(tokenswap.ToToken, tokenswap.FromNodeAddress, tokenswap.ToAmount, tokenswap.Identifier, msg.HashLock, takerExpiration)
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
func (rs *RaidenService) tokenSwapTaker(tokenswap *TokenSwap) (result *network.AsyncResult) {
	result = network.NewAsyncResult()
	result.Result <- nil
	key := swapKey{
		Identifier: tokenswap.Identifier,
		FromToken:  tokenswap.FromToken,
		FromAmount: tokenswap.FromAmount.String(),
	}
	rs.SwapKey2TokenSwap[key] = tokenswap
	return
}

//all user's request
func (rs *RaidenService) handleReq(req *apiReq) {
	var result *network.AsyncResult
	switch req.Name {
	case transferReqName: //mediated transfer only
		r := req.Req.(*transferReq)
		result = rs.startMediatedTransfer(r.TokenAddress, r.Target, r.Amount, r.Fee, r.Identifier)
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
			_, ok := sentMessage.Message.(*encoding.Secret)
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
			case *encoding.RemoveExpiredHashlockTransfer:
				rs.db.UpdateSentRemoveExpiredHashlockTransfer(sentMessageTag.EchoHash)
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
func (rs *RaidenService) handleRoutesTask(task *routesToDetect) {
	/*
		no need to modify InitStateChange's Routes, because routesTask share the same instance
	*/
	switch task.InitStateChange.(type) {
	case *mediatedtransfer.ActionInitInitiatorStateChange:
		//do nothing
	case *mediatedtransfer.ActionInitMediatorStateChange:
		//do nothing
	}
	rs.StateMachineEventHandler.logAndDispatch(task.StateManager, task.InitStateChange)
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
