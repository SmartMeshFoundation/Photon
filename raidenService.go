package raiden_network

import (
	"crypto/ecdsa"
	"errors"

	"fmt"

	"os"

	"path/filepath"

	"time"

	"math/rand"

	"sync"

	"github.com/SmartMeshFoundation/raiden-network/blockchain"
	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/network"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/rerr"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/db"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer/initiator"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer/mediator"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer/target"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/theckman/go-flock"
	"go.uber.org/atomic"
)

//simulate async result of python

type SwapKey struct {
	Identifier uint64
	FromToken  common.Address
	FromAmount int64
}
type TokenSwap struct {
	Identifier      uint64
	FromToken       common.Address
	FromAmount      int64
	FromNodeAddress common.Address //the node address of the owner of the `from_token`
	ToToken         common.Address
	ToAmount        int64
	ToNodeAddress   common.Address //the node address of the owner of the `to_token`
}

/*
tokenAddress common.Address, amount int64, target common.Address, identifier uint64
*/
type TransferReq struct {
	TokenAddress common.Address
	Amount       int64
	Target       common.Address
	Identifier   uint64
	ReqId        string
}

// A Raiden node.
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
	Token2ChannelGraph map[common.Address]*network.ChannelGraph //多线程
	//Token2ConnectionsManager todo fix later
	//swapkey_to_tokenswap
	//swapkey_to_greenlettask
	Manager2Token            map[common.Address]common.Address
	Identifier2StateManagers map[uint64][]*transfer.StateManager
	Identifier2Results       map[uint64][]*network.AsyncResult
	SwapKey2TokenSwap        map[SwapKey]*TokenSwap
	SwapKey2Task             map[SwapKey]Tasker
	Tokens2ConnectionManager map[common.Address]*ConnectionManager
	/*
			   This is a map from a hashlock to a list of channels, the same
		         hashlock can be used in more than one token (for tokenswaps), a
		         channel should be removed from this list only when the Lock is
		         released/withdrawn but not when the secret is registered.
	*/
	Token2Hashlock2Channels  map[common.Address]map[common.Hash][]*channel.Channel //多线程
	MessageHandler           *RaidenMessageHandler
	StateMachineEventHandler *StateMachineEventHandler
	BlockChainEvents         *blockchain.BlockChainEvents
	AlarmTask                *blockchain.AlarmTask
	TransactionLog           *db.StateChangeLog
	GreenletTasksDispatcher  *GreenletTasksDispatcher
	FileLocker               *flock.Flock
	SnapshortDir             string
	SerializationFile        string
	BlockNumber              atomic.Int64
	BlockNumberChan          chan int64
	Lock                     sync.RWMutex
	TransferReqChan          chan *TransferReq
	TransferResponseChan     map[string]chan *network.AsyncResult
}

func NewRaidenService(chain *rpc.BlockChainService, privateKey *ecdsa.PrivateKey, transport network.Transporter,
	discover network.DiscoveryInterface, config *params.Config) (srv *RaidenService) {
	if config.SettleTimeout < params.NETTINGCHANNEL_SETTLE_TIMEOUT_MIN || config.SettleTimeout > params.NETTINGCHANNEL_SETTLE_TIMEOUT_MAX {
		log.Error(fmt.Sprintf("settle timeout must be in range %d-%d",
			params.NETTINGCHANNEL_SETTLE_TIMEOUT_MIN, params.NETTINGCHANNEL_SETTLE_TIMEOUT_MAX))
		os.Exit(1)
	}
	srv = &RaidenService{
		Chain:                    chain,
		Registry:                 chain.Registry(chain.RegistryAddress),
		RegistryAddress:          chain.RegistryAddress,
		PrivateKey:               privateKey,
		Config:                   config,
		NodeAddress:              crypto.PubkeyToAddress(privateKey.PublicKey),
		Token2ChannelGraph:       make(map[common.Address]*network.ChannelGraph),
		Manager2Token:            make(map[common.Address]common.Address),
		Identifier2StateManagers: make(map[uint64][]*transfer.StateManager),
		Identifier2Results:       make(map[uint64][]*network.AsyncResult),
		Token2Hashlock2Channels:  make(map[common.Address]map[common.Hash][]*channel.Channel),
		SwapKey2TokenSwap:        make(map[SwapKey]*TokenSwap),
		SwapKey2Task:             make(map[SwapKey]Tasker),
		Tokens2ConnectionManager: make(map[common.Address]*ConnectionManager),
		AlarmTask:                blockchain.NewAlarmTask(chain.Client),
		BlockChainEvents:         blockchain.NewBlockChainEvents(chain.Client, chain.RegistryAddress),
		GreenletTasksDispatcher:  NewGreenletTasksDispatcher(),
		BlockNumberChan:          make(chan int64, 1),
		TransferReqChan:          make(chan *TransferReq, 1),
		TransferResponseChan:     make(map[string]chan *network.AsyncResult),
	}
	srv.MessageHandler = NewRaidenMessageHandler(srv)
	srv.StateMachineEventHandler = NewStateMachineEventHandler(srv)
	srv.Protocol = network.NewRaidenProtocol(transport, discover, privateKey)
	srv.TransactionLog = db.NewStateChangeLog(config.DataBasePath)
	srv.FileLocker = flock.NewFlock(config.DataBasePath + ".flock.Lock")
	locked, err := srv.FileLocker.TryLock()
	if err != nil || !locked {
		log.Error(fmt.Sprint("another instance already running at %s", config.DataBasePath))
		os.Exit(1)
	}
	srv.SnapshortDir = filepath.Join(config.DataBasePath)
	err = discover.Register(srv.NodeAddress, srv.Config.ExternIp, srv.Config.ExternPort)
	if err != nil {
		log.Error(fmt.Sprintf("register discover endpoint error:%s", err))
		os.Exit(1)
	}
	log.Info("node discovery register complete...")
	//srv.Start()
	return srv
}

// Start the node.
func (this *RaidenService) Start() {
	this.AlarmTask.Start()
	this.AlarmTask.RegisterCallback(func(number int64) error {
		return this.setBlockNumber(number)
	})
	err := this.BlockChainEvents.InstallEventListener()
	if err != nil {
		log.Error(fmt.Sprintf("BlockChainEvents listener error %v", err))
		os.Exit(1)
	}
	/*
			  Registry registration must start *after* the alarm task, this avoid
		         corner cases were the registry is queried in block A, a new block B
		         is mined, and the alarm starts polling at block C.
	*/
	this.RegisterRegistry()
	err = this.RestoreSnapshot()
	if err != nil {
		log.Error(fmt.Sprintf("restore from snapshot error : %v\n you can delete all the database %s to run. but all your trade will lost!!", err, this.Config.DataBasePath))
		os.Exit(1)
	}
	this.Protocol.Start()
	this.loop()
}

//Stop the node.
func (this *RaidenService) Stop() {
	this.AlarmTask.Stop()
	this.Protocol.StopAndWait()
	this.BlockChainEvents.Stop()
	this.SaveSnapshot()
	time.Sleep(time.Second * 3) // let other goroutines quit
	//anther instance cann run now
	this.FileLocker.Unlock()
}
func (this *RaidenService) loop() {
	var err error
	var ok bool
	var m *network.MessageToRaiden
	var st transfer.StateChange
	var blockNumber int64
	var transferReq *TransferReq
	for {
		select {
		case m, ok = <-this.Protocol.ReceivedMessageChannel:
			if ok {
				err = this.MessageHandler.OnMessage(m.Msg, m.EchoHash)
				if err != nil {
					log.Error(fmt.Sprintf("MessageHandler.OnMessage %v", err))
				}
				this.Protocol.ReceivedMessageResultChannel <- err
			} else {
				log.Info("Protocol.ReceivedMessageChannel closed")
				return
			}
		case st, ok = <-this.BlockChainEvents.StateChangeChannel:
			if ok {
				err = this.StateMachineEventHandler.OnBlockchainStateChange(st)
				if err != nil {
					log.Error("StateMachineEventHandler.OnBlockchainStateChange", err)
				}
			} else {
				log.Info("BlockChainEvents.StateChangeChannel closed")
				return
			}
		case blockNumber, ok = <-this.BlockNumberChan:
			if ok {
				this.handleBlockNumber(blockNumber)
			} else {
				log.Info("BlockNumberChan closed")
				return
			}
		case transferReq, ok = <-this.TransferReqChan:
			if ok {
				r := transferReq
				result := this.StartMediatedTransfer(r.TokenAddress, r.Target, r.Amount, r.Identifier)
				this.Lock.Lock()
				//should never block
				this.TransferResponseChan[r.ReqId] <- result
				close(this.TransferResponseChan[r.ReqId])
				delete(this.TransferResponseChan, r.ReqId)
				this.Lock.Unlock()
			} else {
				log.Info("transferReq closed")
			}
		}
	}
}

func (this *RaidenService) RegisterRegistry() {
	mgrs, err := this.Chain.GetAllChannelManagers()
	if err != nil {
		log.Error(fmt.Sprintf("RegisterRegistry err:%s", err))
		os.Exit(1)
	}
	for _, mgr := range mgrs {
		err = this.RegisterChannelManager(mgr.Address)
		if err != nil {
			log.Error(fmt.Sprintf("RegisterChannelManager err:%s", err))
			os.Exit(1)
		}
	}
}

func (this *RaidenService) getChannelDetail(tokenAddress common.Address, proxy *rpc.NettingChannelContractProxy) *network.ChannelDetails {
	addr1, b1, addr2, b2, _ := proxy.AddressAndBalance()
	var ourAddr, partnerAddr common.Address
	var ourBalance, partnerBalance int64
	if addr1 == this.NodeAddress {
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
		this.RegisterChannelForHashlock(tokenAddress, channel, hashlock)
	}
	externState := channel.NewChannelExternalState(registerChannelForHashlock, proxy, channelAddress, this.Chain)
	channelDetail := &network.ChannelDetails{
		ChannelAddress:    channelAddress,
		OurState:          ourState,
		PartenerState:     partenerState,
		ExternState:       externState,
		BlockChainService: this.Chain,
		RevealTimeout:     this.Config.RevealTimeout,
	}
	channelDetail.SettleTimeout, _ = externState.NettingChannel.SettleTimeout()
	return channelDetail
}

func (this *RaidenService) setBlockNumber(blocknumber int64) error {
	this.BlockNumberChan <- blocknumber
	return nil
}
func (this *RaidenService) handleBlockNumber(blocknumber int64) error {
	statechange := &transfer.BlockStateChange{blocknumber}
	this.StateMachineEventHandler.LogAndDispatchToAllTasks(statechange)
	for _, cg := range this.CloneToken2ChannelGraph() {
		cg.Lock.Lock() //dead lock!!..
		for _, channel := range cg.ChannelAddress2Channel {
			channel.StateTransition(statechange)
		}
		cg.Lock.Unlock()
	}
	this.BlockNumber.Store(blocknumber)
	return nil
}
func (this *RaidenService) GetBlockNumber() int64 {
	return this.BlockNumber.Load()
}

func (this *RaidenService) FindChannelByAddress(nettingChannelAddress common.Address) (*channel.Channel, error) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	for _, g := range this.Token2ChannelGraph {
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
func (this *RaidenService) SendAsync(recipient common.Address, msg encoding.SignedMessager) error {
	if recipient == this.NodeAddress {
		log.Error(fmt.Sprintf("this must be a bug ,sending message to it self"))
	}
	_, err := this.Protocol.SendAsync(recipient, msg)
	return err
}

/*
Send `message` to `recipient` and wait for the response or `timeout`.

       Args:
           recipient (address): The address of the node that will receive the
               message.
           message: The transfer message.
           timeout (float): How long should we wait for a response from `recipient`.

       Returns:
           None: If the wait timed out
           object: The result from the event
*/
func (this *RaidenService) SendAndWait(recipient common.Address, message encoding.SignedMessager, timeout time.Duration) error {
	return this.Protocol.SendAndWait(recipient, message, timeout)
}

/*
Register the secret with any channel that has a hashlock on it.

       This must search through all channels registered for a given hashlock
       and ignoring the tokens. Useful for refund transfer, split transfer,
       and token swaps.

       Raises:
           TypeError: If secret is unicode data.
*/
func (this *RaidenService) RegisterSecret(secret common.Hash) {
	hashlock := utils.Sha3(secret[:])
	revealSecretMessage := encoding.NewRevealSecret(secret)
	revealSecretMessage.Sign(this.PrivateKey, revealSecretMessage)
	for _, hashchannel := range this.Token2Hashlock2Channels {
		for _, ch := range hashchannel[hashlock] {
			ch.RegisterSecret(secret)
			//The protocol ignores duplicated messages.
			this.SendAsync(ch.PartnerState.Address, revealSecretMessage)
		}
	}
}

func (this *RaidenService) RegisterChannelForHashlock(tokenAddress common.Address,
	netchannel *channel.Channel, hashlock common.Hash) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	channelsRegistered := this.Token2Hashlock2Channels[tokenAddress][hashlock]
	found := false
	for _, c := range channelsRegistered {
		//判断两个channel对象是否相等,这里只是简单用了地址来判断,可能存在问题 todo
		if c.ExternState.ChannelAddress == netchannel.ExternState.ChannelAddress {
			found = true
			break
		}
	}
	if !found {
		hashLock2Channels, ok := this.Token2Hashlock2Channels[tokenAddress]
		if !ok {
			hashLock2Channels = make(map[common.Hash][]*channel.Channel)
			this.Token2Hashlock2Channels[tokenAddress] = hashLock2Channels
		}
		channelsRegistered = append(channelsRegistered, netchannel)
		this.Token2Hashlock2Channels[tokenAddress][hashlock] = channelsRegistered
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
*/
func (this *RaidenService) HandleSecret(identifier uint64, tokenAddress common.Address, secret common.Hash,
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
	channelsList := this.Token2Hashlock2Channels[tokenAddress][hashlock]
	var channelsToRemove []*channel.Channel
	revealSecretMessage := encoding.NewRevealSecret(secret)
	revealSecretMessage.Sign(this.PrivateKey, revealSecretMessage)
	type MsgToSend struct {
		receiver common.Address
		msg      encoding.SignedMessager
	}
	var messagesToSend []*MsgToSend
	for _, ch := range channelsList { //处理在间接交易过程中重复使用的节点
		//unlock a pending Lock
		if ch.OurState.IsKnown(hashlock) {
			var secretMsg *encoding.Secret
			secretMsg, err = ch.CreateSecret(identifier, secret)
			if err != nil {
				return err
			}
			secretMsg.Sign(this.PrivateKey, secretMsg)
			//balance proof,完成本次交易,收到重复的revealsecret,但是只会发送一次secretmsg,节点故障了,这个channel就不能用了.
			err = ch.RegisterTransfer(this.GetBlockNumber(), secretMsg)
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
					err = ch.RegisterTransfer(this.GetBlockNumber(), msg)
					if err != nil {
						return
					}
					channelsToRemove = append(channelsToRemove, ch)
				} else {
					err = ch.RegisterSecret(secret)
					if err != nil {
						return
					}
					messagesToSend = append(messagesToSend, &MsgToSend{ch.PartnerState.Address, revealSecretMessage})
				}
			} else {
				err = ch.RegisterSecret(secret)
				if err != nil {
					return
				}
				messagesToSend = append(messagesToSend, &MsgToSend{ch.PartnerState.Address, revealSecretMessage})
			}
		} else {
			log.Error("Channel is registered for a given Lock but the Lock is not contained in it.")
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
		delete(this.Token2Hashlock2Channels[tokenAddress], hashlock)
	} else {
		this.Token2Hashlock2Channels[tokenAddress][hashlock] = channelsList
	}
	// send the messages last to avoid races
	for _, msg := range messagesToSend {
		err = this.SendAsync(msg.receiver, msg.msg)
		if err != nil {
			return
		}
	}
	return
}

func (this *RaidenService) ChannelManagerIsRegistered(manager common.Address) bool {
	this.Lock.RLock()
	defer this.Lock.Unlock()
	_, ok := this.Manager2Token[manager]
	return ok
}
func (this *RaidenService) RegisterChannelManager(managerAddress common.Address) (err error) {
	manager := this.Chain.Manager(managerAddress)
	channels, err := manager.NettingChannelByAddress(this.NodeAddress)
	if err != nil {
		return
	}
	tokenAddress, _ := manager.TokenAddress()
	edgeList, _ := manager.GetChannelsParticipants()
	var channelsDetails []*network.ChannelDetails
	for _, ch := range channels {
		d := this.getChannelDetail(tokenAddress, ch)
		channelsDetails = append(channelsDetails, d)
	}
	graph := network.NewChannelGraph(this.NodeAddress, managerAddress, tokenAddress, edgeList, channelsDetails)
	this.Lock.Lock()
	this.Manager2Token[managerAddress] = tokenAddress
	this.Token2ChannelGraph[tokenAddress] = graph
	this.Tokens2ConnectionManager[tokenAddress] = NewConnectionManager(this, tokenAddress, graph)
	this.Lock.Unlock()
	return
}
func (this *RaidenService) RegisterNettingChannel(tokenAddress, channelAddress common.Address) {
	nettingChannel, err := this.Chain.NettingChannel(channelAddress)
	if err != nil {
		log.Error("try to RegisterNettingChannel not exist channel %s", channelAddress)
	}
	detail := this.getChannelDetail(tokenAddress, nettingChannel)
	graph := this.GetToken2ChannelGraph(tokenAddress)
	err = graph.AddChannel(detail)
	if err != nil {
		log.Error(err.Error())
	}
	return
}
func (this *RaidenService) ConnectionManagerForToken(tokenAddress common.Address) (*ConnectionManager, error) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	mgr, ok := this.Tokens2ConnectionManager[tokenAddress]
	if ok {
		return mgr, nil
	}
	return nil, rerr.InvalidAddress(fmt.Sprintf("token %s is not registered", utils.APex(tokenAddress)))
}
func (this *RaidenService) LeaveAllTokenNetworksAsync() *network.AsyncResult {
	var leaveResults []*network.AsyncResult
	for token, _ := range this.CloneToken2ChannelGraph() {
		mgr, _ := this.ConnectionManagerForToken(token)
		if mgr != nil {
			leaveResults = append(leaveResults, mgr.LeaveAsync())
		}
	}
	return WaitGroupAsyncResult(leaveResults)
}
func WaitGroupAsyncResult(results []*network.AsyncResult) *network.AsyncResult {
	totalResult := network.NewAsyncResult()
	wg := sync.WaitGroup{}
	wg.Add(len(results))
	for i, _ := range results {
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
func (this *RaidenService) CloseAndSettle() {
	log.Info("raiden will close and settle all channels now")
	var Mgrs []*ConnectionManager
	this.Lock.RLock()
	for t, _ := range this.Token2ChannelGraph {
		mgr, _ := this.ConnectionManagerForToken(t)
		if mgr != nil {
			Mgrs = append(Mgrs, mgr)
		}
	}
	this.Lock.RUnlock()
	blocksToWait := func() int64 {
		var max int64 = 0
		for _, mgr := range Mgrs {
			if max < mgr.minSettleBlocks() {
				max = mgr.minSettleBlocks()
			}
		}
		return max
	}
	var AllChannels []*channel.Channel
	for _, mgr := range Mgrs {
		AllChannels = append(AllChannels, mgr.openChannels()...)
	}
	leavingResult := this.LeaveAllTokenNetworksAsync()
	//using the un-cached block number here
	lastBlock := this.GetBlockNumber()
	earliestSettlement := lastBlock + blocksToWait()
	/*
			    TODO: estimate and set a `timeout` parameter in seconds
		     based on connection_manager.min_settle_blocks and an average
		     blocktime from the past
	*/
	currentBlock := lastBlock
	for currentBlock < earliestSettlement {
		time.Sleep(time.Second * 10)
		lastBlock := this.GetBlockNumber()
		if lastBlock != currentBlock {
			currentBlock = lastBlock
			waitBlocksLeft := blocksToWait()
			notSettled := 0
			for _, c := range AllChannels {
				if c.State() != transfer.CHANNEL_STATE_SETTLED {
					notSettled++
				}
			}
			if notSettled == 0 {
				log.Debug("nothing left to settle")
				break
			}
			log.Info(fmt.Sprintf("waiting at least %s more blocks for %d channels not yet settled", waitBlocksLeft, notSettled))
		}
		//why  leaving_greenlet.wait
		timeoutch := time.After(time.Second * time.Duration(blocksToWait()))
		select {
		case <-timeoutch:
		case <-leavingResult.Result:
		}
	}
	for _, c := range AllChannels {
		if c.State() != transfer.CHANNEL_STATE_SETTLED {
			log.Error("channels were not settled:", utils.APex(c.MyAddress))
		}
	}
}

/*
Transfer `amount` between this node and `target`.

       This method will start an asyncronous transfer, the transfer might fail
       or succeed depending on a couple of factors:

           - Existence of a path that can be used, through the usage of direct
             or intermediary channels.
           - Network speed, making the transfer sufficiently fast so it doesn't
             expire.
*/
func (this *RaidenService) MediatedTransferAsync(tokenAddress common.Address, amount int64, target common.Address, identifier uint64) *network.AsyncResult {
	req := &TransferReq{
		TokenAddress: tokenAddress,
		Amount:       amount,
		Target:       target,
		Identifier:   identifier,
		ReqId:        utils.RandomString(10),
	}
	result := make(chan *network.AsyncResult, 1)
	this.Lock.Lock()
	this.TransferResponseChan[req.ReqId] = result
	this.Lock.Unlock()
	this.TransferReqChan <- req
	ar := <-result
	return ar
	//return this.StartMediatedTransfer(tokenAddress, target, amount, identifier)
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
func (this *RaidenService) DirectTransferAsync(tokenAddress, target common.Address, amount int64, identifier uint64) (result *network.AsyncResult) {
	graph := this.GetToken2ChannelGraph(tokenAddress)
	directChannel := graph.GetPartenerAddress2Channel(target)
	result = network.NewAsyncResult()
	if directChannel == nil || !directChannel.CanTransfer() || directChannel.Distributable() < amount {
		result.Result <- errors.New("no available direct channel")
		return
	} else {
		tr, err := directChannel.CreateDirectTransfer(amount, identifier)
		if err != nil {
			result.Result <- err
			return
		}
		tr.Sign(this.PrivateKey, tr)
		directChannel.RegisterTransfer(this.GetBlockNumber(), tr)
		directTransferStateChange := &transfer.ActionTransferDirectStateChange{
			Identifier:   identifier,
			Amount:       amount,
			TokenAddress: tokenAddress,
			NodeAddress:  directChannel.PartnerState.Address,
		}
		// TODO: add the transfer sent event
		stateChangeId, _ := this.TransactionLog.Log(directTransferStateChange)
		//This should be set once the direct transfer is acknowledged
		transferSuccess := transfer.EventTransferSentSuccess{
			Identifier: identifier,
			Amount:     amount,
			Target:     target,
		}
		this.TransactionLog.LogEvents(stateChangeId, []transfer.Event{transferSuccess}, this.GetBlockNumber())
		result2, err := this.Protocol.SendAsync(directChannel.PartnerState.Address, tr)
		if err != nil {
			result.Result <- err
			return
		} else {
			result = result2
		}
	}
	return
}
func (this *RaidenService) StartMediatedTransfer(tokenAddress, target common.Address, amount int64, identifier uint64) (result *network.AsyncResult) {
	graph := this.GetToken2ChannelGraph(tokenAddress)
	availableRoutes := graph.GetBestRoutes(this.Protocol, this.NodeAddress, target, amount, utils.EmptyAddress)
	result = network.NewAsyncResult()
	if len(availableRoutes) <= 0 {
		result.Result <- errors.New("no available route")
		return result
	}
	if identifier == 0 {
		identifier = rand.New(utils.RandSrc).Uint64()
	}
	routesState := transfer.NewRoutesState(availableRoutes)
	transferState := &mediated_transfer.LockedTransferState{
		Identifier: identifier,
		Amount:     amount,
		Token:      tokenAddress,
		Initiator:  this.NodeAddress,
		Target:     target,
		Expiration: 0,
		Hashlock:   utils.EmptyHash,
		Secret:     utils.EmptyHash,
	}
	/*
			  Issue #489

		        Raiden may fail after a state change using the random generator is
		        handled but right before the snapshot is taken. If that happens on
		        the next initialization when raiden is recovering and applying the
		        pending state changes a new secret will be generated and the
		        resulting events won't match, this breaks the architecture model,
		        since it's assumed the re-execution of a state change will always
		        produce the same events.

		        TODO: Removed the secret generator from the InitiatorState and add
		        the secret into all state changes that require one, this way the
		        secret will be serialized with the state change and the recovery will
		        use the same /random/ secret.
	*/
	initInitiator := &mediated_transfer.ActionInitInitiatorStateChange{
		OurAddress:      this.NodeAddress,
		Tranfer:         transferState,
		Routes:          routesState,
		RandomGenerator: utils.RandomGenerator,
		BlockNumber:     this.GetBlockNumber(),
	}
	stateManger := transfer.NewStateManager(initiator.StateTransition, nil, initiator.NameInitiatorTransition)
	/*
		  TODO: implement the network timeout raiden.config['msg_timeout'] and
			cancel the current transfer if it hapens (issue #374)
	*/
	/*
		first register the this transfer id, otherwise error may occur imediatelly
	*/
	mgrs := this.Identifier2StateManagers[identifier]
	mgrs = append(mgrs, stateManger)
	this.Identifier2StateManagers[identifier] = mgrs
	results := this.Identifier2Results[identifier]
	results = append(results, result)
	this.Identifier2Results[identifier] = results

	this.StateMachineEventHandler.LogAndDispatch(stateManger, initInitiator)
	return result
}

//收到了MediatedTransfer
func (this *RaidenService) MediateMediatedTransfer(msg *encoding.MediatedTransfer) {
	amount := msg.Amount.Int64()
	target := msg.Target
	token := msg.Token
	graph := this.GetToken2ChannelGraph(token)
	avaiableRoutes := graph.GetBestRoutes(this.Protocol, this.NodeAddress, target, amount, msg.Sender)
	fromChannel := graph.GetPartenerAddress2Channel(msg.Sender)
	fromRoute := network.Channel2RouteState(fromChannel, msg.Sender)
	ourAddress := this.NodeAddress
	fromTransfer := mediated_transfer.LockedTransferFromMessage(msg)
	routesState := transfer.NewRoutesState(avaiableRoutes)
	blockNumber := this.GetBlockNumber()
	initMediator := &mediated_transfer.ActionInitMediatorStateChange{
		OurAddress:  ourAddress,
		FromTranfer: fromTransfer,
		Routes:      routesState,
		FromRoute:   fromRoute,
		BlockNumber: blockNumber,
	}
	stateManager := transfer.NewStateManager(mediator.StateTransition, nil, mediator.NameMediatorTransition)
	this.StateMachineEventHandler.LogAndDispatch(stateManager, initMediator)
	mgrs := this.Identifier2StateManagers[msg.Identifier]
	mgrs = append(mgrs, stateManager)
	this.Identifier2StateManagers[msg.Identifier] = mgrs
}

//收到了MediatedTransfer,我是最终目标
func (this *RaidenService) TargetMediatedTransfer(msg *encoding.MediatedTransfer) {
	graph := this.GetToken2ChannelGraph(msg.Token)
	fromChannel := graph.GetPartenerAddress2Channel(msg.Sender)
	fromRoute := network.Channel2RouteState(fromChannel, msg.Sender)
	fromTransfer := mediated_transfer.LockedTransferFromMessage(msg)
	initTarget := &mediated_transfer.ActionInitTargetStateChange{
		OurAddress:  this.NodeAddress,
		FromRoute:   fromRoute,
		FromTranfer: fromTransfer,
		BlockNumber: this.GetBlockNumber(),
	}
	stateManger := transfer.NewStateManager(target.StateTransiton, nil, target.NameTargetTransition)
	this.StateMachineEventHandler.LogAndDispatch(stateManger, initTarget)
	identifier := msg.Identifier
	mgrs := this.Identifier2StateManagers[identifier]
	mgrs = append(mgrs, stateManger)
	this.Identifier2StateManagers[identifier] = mgrs
}

func (this *RaidenService) StartHealthCheckFor(address common.Address) {

}

func (this *RaidenService) StartNeighboursHealthCheck() {

}

func (this *RaidenService) GetToken2ChannelGraph(tokenAddress common.Address) (cg *network.ChannelGraph) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	cg = this.Token2ChannelGraph[tokenAddress]
	return
}

func (this *RaidenService) CloneToken2ChannelGraph() map[common.Address]*network.ChannelGraph {
	m := make(map[common.Address]*network.ChannelGraph)
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	for k, v := range this.Token2ChannelGraph {
		m[k] = v
	}
	return m
}

func (this *RaidenService) GetChannelWithAddr(channelAddr common.Address) *channel.Channel {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	for _, g := range this.Token2ChannelGraph {
		c := g.GetChannelAddress2Channel(channelAddr)
		if c != nil {
			return c
		}
	}
	return nil
}

func (this *RaidenService) GetChannel(tokenAddr, partnerAddr common.Address) *channel.Channel {
	g := this.GetToken2ChannelGraph(tokenAddr)
	if g == nil {
		return nil
	}
	return g.GetPartenerAddress2Channel(partnerAddr)
}
