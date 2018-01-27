package raiden_network

import (
	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/rerr"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

/*
 Class responsible to handle the protocol messages.

    Note:
        This class is not intended to be used standalone, use RaidenService
        instead.
*/
type RaidenMessageHandler struct {
	raiden        *RaidenService
	blockedTokens map[common.Address]bool
}

func NewRaidenMessageHandler(raiden *RaidenService) *RaidenMessageHandler {
	h := &RaidenMessageHandler{
		raiden:        raiden,
		blockedTokens: make(map[common.Address]bool),
	}
	return h
}

/*
 Handles `message` and sends an ACK on success.
*/
func (this *RaidenMessageHandler) OnMessage(msg encoding.SignedMessager, hash common.Hash) (err error) {
	switch m2 := msg.(type) {
	case *encoding.SecretRequest:
		err = this.messageSecretRequest(m2)
	case *encoding.RevealSecret:
		err = this.messageRevealSecret(m2)
	case *encoding.Secret:
		err = this.messageSecret(m2)
	case *encoding.DirectTransfer:
		err = this.messageDirectTransfer(m2)
	case *encoding.MediatedTransfer:
		err = this.MessageMediatedTransfer(m2)
	case *encoding.RefundTransfer:
		err = this.messageRefundTransfer(m2)
	default:
		log.Error(fmt.Sprintf("RaidenMessageHandler unknown msg:%s", utils.StringInterface1(msg)))
		return fmt.Errorf("unhandled message cmdid:%d", msg.Cmd())
	}
	return err
}

func (this *RaidenMessageHandler) balanceProof(msger encoding.EnvelopMessager) {
	//blanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(msger)
	msg := msger.GetEnvelopMessage()
	balanceProof := &mediated_transfer.ReceiveBalanceProofStateChange{
		Identifier:   msg.Identifier,
		NodeAddress:  msg.Sender,
		BalanceProof: transfer.NewBalanceProofStateFromEnvelopMessage(msger),
	}
	this.raiden.StateMachineEventHandler.LogAndDispatchByIdentifier(balanceProof.Identifier, balanceProof)
}
func (this *RaidenMessageHandler) messageRevealSecret(msg *encoding.RevealSecret) error {
	secret := msg.Secret
	sender := msg.Sender
	this.raiden.GreenletTasksDispatcher.DispatchMessage(msg, msg.HashLock())
	this.raiden.RegisterSecret(secret)
	stateChange := &mediated_transfer.ReceiveSecretRevealStateChange{secret, sender}
	this.raiden.StateMachineEventHandler.LogAndDispatchToAllTasks(stateChange)
	return nil
}
func (this *RaidenMessageHandler) messageSecretRequest(msg *encoding.SecretRequest) error {
	this.raiden.GreenletTasksDispatcher.DispatchMessage(msg, msg.HashLock)
	stateChange := &mediated_transfer.ReceiveSecretRequestStateChange{
		Identifier: msg.Identifier,
		Amount:     new(big.Int).Set(msg.Amount),
		Hashlock:   msg.HashLock,
		Sender:     msg.Sender,
	}
	this.raiden.StateMachineEventHandler.LogAndDispatchByIdentifier(msg.Identifier, stateChange)
	return nil
}
func (this *RaidenMessageHandler) messageSecret(msg *encoding.Secret) error {
	this.balanceProof(msg)
	hashlock := msg.HashLock()
	identifer := msg.Identifier
	secret := msg.Secret
	this.raiden.RegisterSecret(secret)
	var nettingChannel *channel.Channel
	var err error
	nettingChannel, err = this.raiden.FindChannelByAddress(msg.Channel)
	if err != nil {
		log.Info(fmt.Sprintf("Message for unknown channel: %s", err))
	} else {
		this.raiden.HandleSecret(identifer, nettingChannel.TokenAddress, secret, msg, hashlock)
	}
	this.raiden.GreenletTasksDispatcher.DispatchMessage(msg, hashlock)
	return nil
}

func (this *RaidenMessageHandler) messageRefundTransfer(msg *encoding.RefundTransfer) (err error) {
	this.balanceProof(msg)
	graph := this.raiden.GetToken2ChannelGraph(msg.Token)
	if !graph.HashChannel(this.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("Direct transfer from node without an existing channel: %s", msg.Sender)
		return
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	err = ch.RegisterTransfer(this.raiden.GetBlockNumber(), msg)
	if err != nil {
		return
	}
	this.raiden.GreenletTasksDispatcher.DispatchMessage(msg, msg.HashLock)
	transferState := &mediated_transfer.LockedTransferState{
		Identifier: msg.Identifier,
		Amount:     new(big.Int).Set(msg.Amount),
		Token:      msg.Token,
		Initiator:  msg.Initiator,
		Target:     msg.Target,
		Expiration: msg.Expiration,
		Hashlock:   msg.HashLock,
		Secret:     utils.EmptyHash}
	stateChange := &mediated_transfer.ReceiveTransferRefundStateChange{msg.Sender, transferState}
	this.raiden.StateMachineEventHandler.LogAndDispatchByIdentifier(msg.Identifier, stateChange)
	return nil
}

func (this *RaidenMessageHandler) messageDirectTransfer(msg *encoding.DirectTransfer) error {
	this.balanceProof(msg)
	if graph := this.raiden.GetToken2ChannelGraph(msg.Token); graph == nil {
		return rerr.UnknownTokenAddress(msg.Token.String())
	}
	if _, ok := this.blockedTokens[msg.Token]; ok {
		return rerr.TransferUnwanted
	}
	graph := this.raiden.GetToken2ChannelGraph(msg.Token)
	if !graph.HashChannel(this.raiden.NodeAddress, msg.Sender) {
		return rerr.UnknownAddress(fmt.Sprintf("Direct transfer from node without an existing channel partner %s  ", msg.Sender))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch.State() != transfer.CHANNEL_STATE_OPENED {
		return rerr.TransferWhenClosed(ch.MyAddress.String())
	}
	var amount = new(big.Int)
	amount = amount.Sub(msg.TransferAmount, ch.PartnerState.TransferAmount())
	stateChange := &transfer.ReceiveTransferDirectStateChange{
		Identifier:   msg.Identifier,
		Amount:       amount,
		TokenAddress: msg.Token,
		Sender:       msg.Sender,
	}
	stateChangeId, err := this.raiden.TransactionLog.Log(stateChange)
	if err != nil {
		return err
	}
	ch.RegisterTransfer(this.raiden.GetBlockNumber(), msg)
	receiveSuccess := &transfer.EventTransferReceivedSuccess{
		Identifier: msg.Identifier,
		Amount:     amount,
		Initiator:  msg.Sender,
	}
	err = this.raiden.TransactionLog.LogEvents(stateChangeId, []transfer.Event{receiveSuccess}, this.raiden.GetBlockNumber())
	return err
}

func (this *RaidenMessageHandler) MessageMediatedTransfer(msg *encoding.MediatedTransfer) error {
	this.balanceProof(msg)
	//  TODO: Reject mediated transfer that the hashlock/identifier is known,
	// this is a downstream bug and the transfer is going in cycles (issue #490)
	key := SwapKey{msg.Identifier, msg.Token, msg.Amount.String()}
	if _, ok := this.blockedTokens[msg.Token]; ok {
		return rerr.TransferUnwanted
	}
	/*
			 TODO: add a separate message for token swaps to simplify message
		     handling (issue #487)
	*/
	if _, ok := this.raiden.SwapKey2TokenSwap[key]; ok {
		this.messageTokenSwap(msg)
		return nil
	}
	graph := this.raiden.GetToken2ChannelGraph(msg.Token)
	if !graph.HashChannel(this.raiden.NodeAddress, msg.Sender) {
		return rerr.ChannelNotFound(fmt.Sprintf("mediated transfer from node without an existing channel %s", msg.Sender))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch.State() != transfer.CHANNEL_STATE_OPENED {
		return rerr.TransferWhenClosed(fmt.Sprintf("Mediated transfer received but the channel is closed %s", ch.MyAddress))
	}
	err := ch.RegisterTransfer(this.raiden.GetBlockNumber(), msg)
	if err != nil {
		return err
	}
	if msg.Target == this.raiden.NodeAddress {
		this.raiden.TargetMediatedTransfer(msg)
	} else {
		this.raiden.MediateMediatedTransfer(msg)
	}
	return nil
}
func (this *RaidenMessageHandler) messageTokenSwap(msg *encoding.MediatedTransfer) {
	key := SwapKey{
		Identifier: msg.Identifier,
		FromToken:  msg.Token,
		FromAmount: msg.Amount.String(),
	}
	/*
			If we are the maker the task is already running and waiting for the
		    taker's MediatedTransfer
	*/
	task := this.raiden.SwapKey2Task[key]
	if task != nil {
		task.GetResponseChan() <- msg
	} else {
		/*
		   If we are the taker we are receiving the maker transfer and should start our new task
		*/
		tokenSwap := this.raiden.SwapKey2TokenSwap[key]
		task := NewTakerTokenSwapTask(this.raiden, tokenSwap, msg)
		go func() {
			task.Start()
		}()
		this.raiden.SwapKey2Task[key] = task
	}
}
