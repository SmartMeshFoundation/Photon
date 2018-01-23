package raiden_network

import (
	"fmt"

	"time"

	"sync"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/network"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/rerr"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
	"github.com/kataras/go-errors"
)

type Tasker interface {
	//可能会写入多个想消息,什么时候关闭呢?
	GetResponseChan() chan encoding.SignedMessager
	Start()
	Stop()
}
type GreenletTasksDispatcher struct {
	//use map to simulate set
	lock           sync.Mutex
	hashlock2Tasks map[common.Hash]map[Tasker]bool
}

func NewGreenletTasksDispatcher() *GreenletTasksDispatcher {
	return &GreenletTasksDispatcher{
		hashlock2Tasks: make(map[common.Hash]map[Tasker]bool),
	}
}

/*
Register the task to receive messages based on `hashlock`.

        Registration is required otherwise the task won't receive any messages
        from the protocol, un-registering is done by the `unregister_task`
        function.

        Note:
            Messages are dispatched solely on the hashlock value (being part of
            the message, eg. SecretRequest, or calculated from the message
            content, eg.  RevealSecret), this means the sender needs to be
            checked for the received messages.
*/
func (this *GreenletTasksDispatcher) RegisterTask(task Tasker, hashlock common.Hash) {
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.hashlock2Tasks[hashlock]
	if !ok {
		m = make(map[Tasker]bool)
		this.hashlock2Tasks[hashlock] = m
	}
	m[task] = true
}
func (this *GreenletTasksDispatcher) UnregisterTask(task Tasker, hashlock common.Hash, success bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.hashlock2Tasks[hashlock]
	if !ok {
		log.Warn(fmt.Sprintf("remove unexisted hashlock %s", hashlock.String()))
		return
	}
	delete(m, task)
	close(task.GetResponseChan()) //close channel when task complete.
	if len(m) == 0 {
		delete(this.hashlock2Tasks, hashlock)
	}
}
func (this *GreenletTasksDispatcher) DispatchMessage(msg encoding.SignedMessager, hashlock common.Hash) {
	for task, _ := range this.hashlock2Tasks[hashlock] {
		task.GetResponseChan() <- msg
	}
}
func (this *GreenletTasksDispatcher) Stop() {
	for _, m := range this.hashlock2Tasks {
		for t, _ := range m {
			t.Stop()
		}
	}
}

type BaseMediatedTransferTask struct {
	ResponseChan chan encoding.SignedMessager
	name         string //for debug
}

func (this *BaseMediatedTransferTask) GetResponseChan() chan encoding.SignedMessager {
	return this.ResponseChan
}

/*
 Utility to handle multiple messages for the same hashlock while
        properly handling expiration timeouts.
*/
var errSendAntWaitTimeout = errors.New("sendAndWaitTime timeout")

func (this *BaseMediatedTransferTask) sendAndWaitTime(raiden *RaidenService, recipient common.Address, mtr *encoding.MediatedTransfer, timeout time.Duration) (msg encoding.SignedMessager, err error) {
	var ok bool
	err = raiden.SendAsync(recipient, mtr)
	if err != nil {
		log.Error("raiden.SendAsync", err)
	}
	timeoutCh := time.After(timeout)
	select {
	case <-timeoutCh:
		err = errSendAntWaitTimeout
		log.Debug(fmt.Sprintf("TIMED OUT %s", utils.StringInterface(mtr, 3)))
	case msg, ok = <-this.ResponseChan:
		if !ok {
			err = rerr.GoChannelClosed("ResponseChan")
		}
	}
	return
}

//Utility to handle multiple messages and timeout on a block number.
func (this *BaseMediatedTransferTask) sendAndWaitBlock(raiden *RaidenService, recipient common.Address, mtr *encoding.MediatedTransfer, expirationBlock int64) (msg encoding.SignedMessager, err error) {
	err = raiden.SendAsync(recipient, mtr)
	if err != nil {
		log.Error("raiden.SendAsync", err)
	}
	msg, err = this.messagesUntilBlock(raiden, expirationBlock)
	if err != nil {
		log.Debug(fmt.Sprintf("TIMED OUT ON BLOCK   tr=%s", utils.StringInterface(mtr, 3)))
	}
	return
}

//Returns the received messages up to the block `expiration_block`.
func (this *BaseMediatedTransferTask) messagesUntilBlock(raiden *RaidenService, expirationBlock int64) (msg encoding.SignedMessager, err error) {
	var ok bool
	current := raiden.GetBlockNumber()
	for current < expirationBlock {
		timeoutCh := time.After(params.DEFAULT_EVENTS_POLL_TIMEOUT)
		select {
		case <-timeoutCh:
		case msg, ok = <-this.ResponseChan:
			if ok {
				err = nil
			} else {
				err = rerr.GoChannelClosed("ResponseChan")
			}
			return
		}
		current = raiden.GetBlockNumber()
	}
	err = rerr.Timeout(fmt.Sprintf("messagesUntilBlock:%d", current))
	return
}

/*
Wait for a Secret message from our partner to update the local
        state, if the Secret message is not sent within time the channel will
        be closed.

        Note:
            Must be called only once the secret is known.
            Must call `unregister_task` after this function returns.
*/
func (this *BaseMediatedTransferTask) waitForUnlockOrClose(raiden *RaidenService, graph *network.ChannelGraph, channel *channel.Channel, mtr *encoding.MediatedTransfer) {
	if graph.TokenAddress != mtr.Token {
		panic("token addess not equal")
	}
	blockToClose := mtr.Expiration - int64(raiden.Config.RevealTimeout)
	hashlock := mtr.HashLock
	token := mtr.Token
	identifier := mtr.Identifier
	for channel.OurState.IsLocked(hashlock) {
		currentBlock := raiden.GetBlockNumber()
		if currentBlock > blockToClose {
			log.Warn(fmt.Sprintf("close channel (%s,%s) to prevent expiration of Lock %s ", utils.APex(channel.OurState.Address), utils.APex(channel.PartnerState.Address), utils.HPex(hashlock)))
			//为什么给的proof是自己的?
			err := channel.ExternState.Close(channel.OurState.BalanceProofState)
			if err != nil {
				log.Error("close channel err:", err)
			}
			return
		}
		timeoutch := time.After(params.DEFAULT_EVENTS_POLL_TIMEOUT)
		var msg encoding.Messager
		timeout := false
		select {
		case msg = <-this.ResponseChan:
		case <-timeoutch:
			timeout = true
		}
		if timeout {
			continue
		}
		if msg == nil {
			log.Error("ResponseChan channel already closed?")
		} else if msg.Cmd() == encoding.SECRET_CMDID {
			//这里设置断点,看看是没收到消息,还是什么原因
			smsg := msg.(*encoding.Secret)
			secret := smsg.Secret
			hashlock := utils.Sha3(secret[:])
			isValidIdentifier := smsg.Identifier == identifier
			isValidChannel := smsg.Channel == channel.MyAddress
			if isValidIdentifier && isValidChannel {
				err := raiden.HandleSecret(identifier, graph.TokenAddress, secret, smsg, hashlock)
				if err != nil {
					log.Error("handle secret:", err)
				}
			} else {
				//cannot use the message but the secret is okay
				raiden.HandleSecret(identifier, graph.TokenAddress, secret, nil, hashlock)
				log.Error(fmt.Sprintf("Invalid Secret message received, expected message for token=%s,identifier=%d received=%s", token, identifier, utils.StringInterface(msg, 3)))
			}
		} else if msg.Cmd() == encoding.REVEALSECRET_CMDID {
			smsg := msg.(*encoding.RevealSecret)
			secret := smsg.Secret
			hashlock := utils.Sha3(secret[:])
			raiden.HandleSecret(identifier, graph.TokenAddress, secret, nil, hashlock)
		} else {
			log.Error("Invalid message ignoring. ", utils.StringInterface(msg, 3))
		}
	}
}

/*
Utility to wait until the expiration block.

        For a chain A-B-C, if an attacker controls A and C a mediated transfer
        can be done through B and C will wait for/send a timeout, for that
        reason B must not unregister the hashlock until the Lock has expired,
        otherwise the revealed secret wouldn't be caught.
*/
func (this *BaseMediatedTransferTask) waitExpiration(raiden *RaidenService, mtr *encoding.MediatedTransfer, sleep time.Duration) {
	expiration := mtr.Expiration + 1
	for {
		current := raiden.GetBlockNumber()
		if current > expiration {
			break
		}
		time.Sleep(sleep)
	}
}

func (this *BaseMediatedTransferTask) Stop() {
	//??
}

/*
 Note: send_and_wait_valid methods are used to check the message type and
 sender only, this can be improved by using a encrypted connection between the
 nodes making the signature validation unnecessary
*/
// Implement the swaps as a restartable task (issue #303)

type MakerTokenSwapTask struct {
	BaseMediatedTransferTask
	raiden      *RaidenService
	tokenswap   *TokenSwap
	asyncResult *network.AsyncResult
}

func NewMakerTokenSwapTask(raiden *RaidenService, tokenSwap *TokenSwap, asyncResult *network.AsyncResult) *MakerTokenSwapTask {
	mtt := &MakerTokenSwapTask{
		raiden:                   raiden,
		tokenswap:                tokenSwap,
		asyncResult:              asyncResult,
		BaseMediatedTransferTask: BaseMediatedTransferTask{make(chan encoding.SignedMessager), "MakerTokenSwapTask"},
	}
	return mtt
}

func (mtt *MakerTokenSwapTask) Start() {
	tokenSwap := mtt.tokenswap
	raiden := mtt.raiden
	identifier := tokenSwap.Identifier
	fromToken := tokenSwap.FromToken
	fromAmount := tokenSwap.FromAmount
	toToken := tokenSwap.ToToken
	toAmount := tokenSwap.ToAmount
	toNodeAddress := tokenSwap.ToNodeAddress
	fromGraph := raiden.Token2ChannelGraph[fromToken]
	toGraph := raiden.Token2ChannelGraph[toToken]

	fromRoutes := fromGraph.GetBestRoutes(raiden.Protocol, raiden.NodeAddress, toNodeAddress, fromAmount, utils.EmptyAddress)
	var fee int64 = 0
	for _, route := range fromRoutes {
		//for each new path a new secret must be used
		secret := utils.RandomGenerator()
		hashlock := utils.Sha3(secret[:])

		fromChannel := fromGraph.GetChannelAddress2Channel(route.ChannelAddress)
		raiden.GreenletTasksDispatcher.RegisterTask(mtt, hashlock)
		raiden.RegisterChannelForHashlock(fromToken, fromChannel, hashlock)
		blockNumber := raiden.GetBlockNumber()
		lockExpiration := blockNumber + int64(fromChannel.SettleTimeout)
		fromMtr, _ := fromChannel.CreateMediatedTransfer(raiden.NodeAddress, toNodeAddress, fee, fromAmount, identifier, lockExpiration, hashlock)
		fromMtr.Sign(raiden.PrivateKey, fromMtr)
		//must be the same block number used to compute lock_expiration
		fromChannel.RegisterTransfer(blockNumber, fromMtr)
		//wait for secret request and mediated transfer
		toMtr := mtt.sendAndWaitValidState(raiden, route.HopNode, toNodeAddress, fromMtr, toToken, toAmount)
		if toMtr == nil {
			/*
						 the initiator can unregister right away since it knows the
				                secret wont be revealed
			*/
			raiden.GreenletTasksDispatcher.UnregisterTask(mtt, hashlock, false)
		} else {
			toHop := toMtr.Sender
			toChannel := toGraph.GetPartenerAddress2Channel(toHop)
			err := toChannel.RegisterTransfer(raiden.GetBlockNumber(), toMtr)
			if err != nil {
				log.Error(" toChannel.RegisterTransfer ", err)
			}
			raiden.RegisterChannelForHashlock(toToken, toChannel, hashlock)
			/*
							  A swap is composed of two mediated transfers, we need to
				                 reveal the secret to both, since the maker is one of the ends
				                 we just need to send the reveal secret directly to the taker.
			*/
			revealSecretMsg := encoding.NewRevealSecret(secret)
			revealSecretMsg.Sign(raiden.PrivateKey, revealSecretMsg)
			raiden.SendAsync(toNodeAddress, revealSecretMsg)
			fromChannel.RegisterSecret(secret)

			/*
							  Register the secret with the to_channel and send the
				                 RevealSecret message to the node that is paying the to_token
				                 (this node might, or might not be the same as the taker),
				                 then wait for the withdraw.
			*/
			raiden.HandleSecret(identifier, toToken, secret, nil, hashlock)
			toChannel = toGraph.GetPartenerAddress2Channel(toMtr.Sender)
			mtt.waitForUnlockOrClose(raiden, toGraph, toChannel, toMtr)
			/*
							   unlock the from_token and optimistically reveal the secret
				                 forward
			*/
			raiden.HandleSecret(identifier, fromToken, secret, nil, hashlock)
			raiden.GreenletTasksDispatcher.UnregisterTask(mtt, hashlock, true)
			mtt.asyncResult.Result <- nil
			close(mtt.asyncResult.Result)
			return
		}
	}
	log.Debug(fmt.Sprintf("MAKER TOKEN SWAP FAILED node=%s,to=%s", utils.APex(raiden.NodeAddress), utils.APex(toNodeAddress)))
	//all routes failed
	mtt.asyncResult.Result <- errors.New("all routes failed")
}

/*
 Start the swap by sending the first mediated transfer to the
        taker and wait for mediated transfer for the exchanged token.

        This method will validate the messages received, discard the invalid
        ones, and wait until a valid state is reached. The valid state is
        reached when a mediated transfer for `to_token` with `to_amount` tokens
        and a SecretRequest from the taker are received.

        Returns:
            None: when the timeout was reached.
            MediatedTransfer: when a valid state is reached.
            RefundTransfer: when an invalid state is reached by
                our partner.
*/
func (mtt *MakerTokenSwapTask) sendAndWaitValidState(raiden *RaidenService, nextHop common.Address, targetAddress common.Address,
	fromMtr *encoding.MediatedTransfer, toToken common.Address, toAmount int64) (mtr *encoding.MediatedTransfer) {
	/*
			        a valid state must have a secret request from the maker and a valid
		        mediated transfer for the new token
	*/
	receivedSecretRequest := false
	for {
		msg, err := mtt.sendAndWaitTime(raiden, fromMtr.Recipient, fromMtr, raiden.Config.MsgTimeout)
		if err != nil {
			log.Error("mtt.sendAndWaitTime", err)
			return
		}
		mtr2, ok := msg.(*encoding.MediatedTransfer)
		/*
					      we need a lower expiration because:
			                     - otherwise the previous node is not operating correctly
			                     - we assume that received mediated transfer has a smaller
			                       expiration to properly call close on edge cases
		*/
		transferIsValid := ok && mtr2.Token == toToken && mtr2.Expiration <= fromMtr.Expiration
		/*
					    The MediatedTransfer might be from `next_hop` or most likely from
			             a different node.
		*/
		if transferIsValid {
			if mtr2.Amount.Int64() == toAmount {
				mtr = mtr2
			}
		} else if msg.Cmd() == encoding.SECRETREQUEST_CMDID && msg.GetSender() == targetAddress {
			receivedSecretRequest = true
		} else if msg.Cmd() == encoding.REFUNDTRANSFER_CMDID && msg.GetSender() == nextHop {
			//失败了,应该告知
			log.Info(fmt.Sprintf("receive refund transfer :%s", utils.StringInterface(msg, 3)))
			return nil
		} else {
			/*
							   The other participant must not use a direct transfer to finish
				             the token swap, ignore it
			*/
			log.Error(fmt.Sprintf("invalide message ingoring.. %s", utils.StringInterface(msg, 3)))
		}
		if mtr != nil && receivedSecretRequest {
			return mtr
		}
	}
	return nil
}

/*
Taker task, responsible to receive a MediatedTransfer for the
    from_transfer and forward a to_transfer with the same hashlock.
*/
type TakerTokenSwapTask struct {
	BaseMediatedTransferTask
	raiden               *RaidenService
	tokenswap            *TokenSwap
	fromMediatedTransfer *encoding.MediatedTransfer
}

/*
Taker task, responsible to receive a MediatedTransfer for the
    from_transfer and forward a to_transfer with the same hashlock.
*/
func NewTakerTokenSwapTask(raiden *RaidenService, tokenswap *TokenSwap, fromMtr *encoding.MediatedTransfer) *TakerTokenSwapTask {
	ttt := &TakerTokenSwapTask{
		raiden:                   raiden,
		tokenswap:                tokenswap,
		fromMediatedTransfer:     fromMtr,
		BaseMediatedTransferTask: BaseMediatedTransferTask{make(chan encoding.SignedMessager), "TakerTokenSwapTask"},
	}
	return ttt
}

func (this *TakerTokenSwapTask) waitRevealSecret(raiden *RaidenService, takerPayingHop common.Address, expirationBlock int64) *encoding.RevealSecret {
	for {
		msg, err := this.messagesUntilBlock(raiden, expirationBlock)
		if err != nil {
			log.Error("waitRevealSecret :", err)
			return nil
		}
		if msg.Cmd() == encoding.REVEALSECRET_CMDID && msg.GetSender() == takerPayingHop {
			return msg.(*encoding.RevealSecret)
		} else {
			log.Error("Invalid message  supplied to the task, ignoring.", utils.StringInterface(msg, 3))
		}
	}
	return nil
}

/*
Start the second half of the exchange and wait for the SecretReveal
        for it.

        This will send the taker mediated transfer with the maker as a target,
        once the maker receives the transfer he is expected to send a
        RevealSecret backwards.
*/
func (this *TakerTokenSwapTask) sendAndWaitValid(raiden *RaidenService, mtr *encoding.MediatedTransfer, makerPayerHop common.Address) (tr encoding.SignedMessager, secret *encoding.Secret) {
	/*
	   the taker cannot discard the transfer since the secret is controlled
	          by another node (the maker), so we have no option but to wait for a
	          valid response until the Lock expires
	*/
	/*
	   Usually the RevealSecret for the MediatedTransfer from this node to
	        the maker should arrive first, but depending on the number of hops
	        and if the maker-path is optimistically revealing the Secret, then
	        the Secret message might arrive first.
	*/
	var revealSecretMsg *encoding.RevealSecret
	var refundMsg *encoding.RefundTransfer
	var rsOk bool
	for {
		msg, err := this.sendAndWaitBlock(raiden, mtr.Recipient, mtr, mtr.Expiration)
		if err != nil {
			log.Error(fmt.Sprintf("TAKER SWAP TIMED OUT node=%s, hashlock=%s", utils.APex(raiden.NodeAddress), utils.HPex(mtr.HashLock)))
			return
		}
		revealSecretMsg, rsOk = msg.(*encoding.RevealSecret)
		validReveal := rsOk && revealSecretMsg.HashLock() == mtr.HashLock && revealSecretMsg.Sender == makerPayerHop

		refundMsg, rsOk = msg.(*encoding.RefundTransfer)
		validRefund := rsOk && refundMsg.Sender == makerPayerHop && refundMsg.Amount == mtr.Amount && refundMsg.Expiration <= mtr.Expiration && refundMsg.Token == mtr.Token
		if msg.Cmd() == encoding.SECRET_CMDID {
			m2 := msg.(*encoding.Secret)
			if utils.Sha3(m2.Secret[:]) != mtr.HashLock {
				log.Error("Secret doesn't match the hashlock, ignoring.")
				continue
			}
			secret = m2
		} else if validReveal {
			return revealSecretMsg, secret
		} else if validRefund {
			return refundMsg, secret
		} else {
			log.Error(fmt.Sprintf("Invalid message   supplied to the task, ignoring.  %s", utils.StringInterface(msg, 3)))
		}
	}
	return
}

func (this *TakerTokenSwapTask) Start() {
	var err error
	var fee int64 = 0
	raiden := this.raiden
	tokenSwap := this.tokenswap
	/*
	   this is the MediatedTransfer that wil pay the maker's half of the
	          swap, not necessarily from him
	*/
	makerPayingTransfer := this.fromMediatedTransfer
	/*
	   this is the address of the node that the taker actually has a channel
	          with (might or might not be the maker)
	*/
	makerPayerHop := makerPayingTransfer.Sender
	if tokenSwap.Identifier != makerPayingTransfer.Identifier ||
		tokenSwap.FromToken != makerPayingTransfer.Token ||
		tokenSwap.FromAmount != makerPayingTransfer.GetLock().Amount ||
		tokenSwap.FromNodeAddress != makerPayingTransfer.Initiator {
		log.Error("TakerTokenSwapTask doesn't match , \ntokenswap=%s\n,makerpayingtransfer=%s", utils.StringInterface(tokenSwap, 2), utils.StringInterface(makerPayingTransfer, 3))
		return
	}
	makerReceivingToken := tokenSwap.ToToken
	toAmount := tokenSwap.ToAmount
	identifier := makerPayingTransfer.Identifier
	hashlock := makerPayingTransfer.HashLock
	makerAddress := makerPayingTransfer.Initiator

	TakerReceivingToken := makerPayingTransfer.Token
	takerPayingToken := makerReceivingToken

	fromGraph := raiden.GetToken2ChannelGraph(TakerReceivingToken)

	fromChannel := fromGraph.GetPartenerAddress2Channel(makerPayerHop)
	toGraph := raiden.GetToken2ChannelGraph(makerReceivingToken)
	//update the channel's distributable and merkle tree
	fromChannel.RegisterTransfer(raiden.GetBlockNumber(), makerPayingTransfer)

	//register the task to receive Refund/Secrect/RevealSecret messages
	raiden.GreenletTasksDispatcher.RegisterTask(this, hashlock)
	raiden.RegisterChannelForHashlock(TakerReceivingToken, fromChannel, hashlock)
	/*
	   send to the maker a secret request informing how much the taker will
	   be _paid_, this is used to inform the maker that his part of the
	   mediated transfer is okay
	*/
	secretRequestMsg := encoding.NewSecretRequest(identifier, makerPayingTransfer.HashLock, makerPayingTransfer.Amount.Int64())
	secretRequestMsg.Sign(raiden.PrivateKey, secretRequestMsg)
	raiden.SendAsync(makerAddress, secretRequestMsg)

	/*
	   Note: taker may only try different routes if a RefundTransfer is
	          received, because the maker is the node controlling the secret
	*/
	availableRoutes := toGraph.GetBestRoutes(
		raiden.Protocol,
		raiden.NodeAddress,
		makerAddress,
		makerPayingTransfer.Amount.Int64(), //这个数量错了?感觉应该是toamount呢?
		utils.EmptyAddress,
	)
	if len(availableRoutes) == 0 {
		log.Debug(fmt.Sprintf("TAKER TOKEN SWAP FAILED, NO ROUTES from=%s,to=%s", utils.APex(raiden.NodeAddress), utils.APex(makerAddress)))
		return
	}
	var firstTransfer *encoding.MediatedTransfer
	for _, route := range availableRoutes {
		takerPayingChannel := toGraph.GetChannelAddress2Channel(route.ChannelAddress)
		takerPayingHop := route.HopNode
		lockExpiration := int64(takerPayingChannel.SettleTimeout) + raiden.GetBlockNumber()
		if lockExpiration > makerPayingTransfer.Expiration {
			lockExpiration = makerPayingTransfer.Expiration
		}
		lockExpiration = lockExpiration - int64(raiden.Config.RevealTimeout)
		log.Debug(fmt.Sprintf("TAKER TOKEN SWAP: from=%s,to=%s  hashlock=%s", utils.APex(raiden.NodeAddress), utils.APex(makerAddress), utils.HPex(hashlock)))
		/*
		   make a paying MediatedTransfer with same hashlock/identifier and the
		                taker's paying token/amount
		*/
		takerPayingTransfer, _ := takerPayingChannel.CreateMediatedTransfer(raiden.NodeAddress, makerAddress, fee, toAmount, identifier, lockExpiration, hashlock)
		takerPayingTransfer.Sign(raiden.PrivateKey, takerPayingTransfer)
		err = takerPayingChannel.RegisterTransfer(raiden.GetBlockNumber(), takerPayingTransfer)
		if err != nil {
			log.Error("RegisterTransfer err ", err)
			return
		}
		if firstTransfer == nil {
			firstTransfer = takerPayingTransfer
		}
		log.Debug(fmt.Sprintf("EXCHANGE TRANSFER NEW PATH path=%s,hashlock=%s", utils.APex(takerPayingHop), utils.HPex(hashlock)))
		//register the task to receive Refund/Secrect/RevealSecret messages
		raiden.RegisterChannelForHashlock(makerReceivingToken, takerPayingChannel, hashlock)
		response, secret := this.sendAndWaitValid(raiden, takerPayingTransfer, makerPayerHop)
		/*
			only refunds for `maker_receiving_token` must be considered
			(check send_and_wait_valid)
		*/
		if response == nil {
			log.Debug(fmt.Sprintf("TAKER TOKEN SWAP FAILED from=%s,to=%s", utils.APex(raiden.NodeAddress), utils.APex(makerAddress)))
			//self.async_result.set(False) 这个明显是个错误啊!
			return
		}
		if response.Cmd() == encoding.REFUNDTRANSFER_CMDID {
			tr := response.(*encoding.RefundTransfer)
			if tr.Amount != takerPayingTransfer.Amount {
				log.Info(fmt.Sprintf("partner %s sent an invalid refund message with an invalid amount ", utils.APex(takerPayingHop)))
				raiden.GreenletTasksDispatcher.UnregisterTask(this, hashlock, false)
				return
			} else {
				takerPayingChannel.RegisterTransfer(raiden.GetBlockNumber(), tr)
			}
		} else if response.Cmd() == encoding.REVEALSECRET_CMDID {
			// the secret was registered by the message handler
			/*
				wait for the taker_paying_hop to reveal the secret prior to
				unlocking locally
			*/
			tr := response.(*encoding.RevealSecret)
			if tr.Sender != takerPayingHop {
				response = this.waitRevealSecret(raiden, takerPayingHop, takerPayingTransfer.Expiration)
			}
			//unlock and send the Secret message
			if err = raiden.HandleSecret(identifier, takerPayingToken, tr.Secret, nil, hashlock); err != nil {
				log.Error("taker HandleSecret err ", err.Error())
			}
			/*
				if the secret arrived early, withdraw it, otherwise send the
				RevealSecret forward in the maker-path
			*/
			if secret != nil {
				raiden.HandleSecret(identifier, TakerReceivingToken, tr.Secret, secret, hashlock)
			}
			//wait for the withdraw in case it did not happen yet
			this.waitForUnlockOrClose(raiden, fromGraph, fromChannel, makerPayingTransfer)
			//自行添加,否则会堵死
			raiden.GreenletTasksDispatcher.UnregisterTask(this, hashlock, true)
			return
		} else {
			log.Debug(fmt.Sprintf("TAKER TOKEN SWAP FAILED from=%s,to=%s", utils.APex(raiden.NodeAddress), utils.APex(makerAddress)))
			//self.async_result.set(False) 这个明显是个错误啊!
			return
		}
	}
	//no route is available, wait for the sent mediated transfer to expire
	this.waitExpiration(raiden, firstTransfer, params.DEFAULT_EVENTS_POLL_TIMEOUT)
	log.Debug(fmt.Sprintf("TAKER TOKEN SWAP FAILED from=%s,to=%s", utils.APex(raiden.NodeAddress), utils.APex(makerAddress)))
	// self.async_result.set(False)
	return
}
