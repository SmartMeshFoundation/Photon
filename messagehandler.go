package smartraiden

import (
	"fmt"

	"math/big"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
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
	msg.SetTag(&transfer.MessageTag{
		EchoHash:          hash,
		IsASendingMessage: false,
		MessageId:         utils.RandomString(10),
	})
	switch m2 := msg.(type) {
	case *encoding.SecretRequest:
		f := this.raiden.SecretRequestPredictorMap[m2.HashLock]
		if f != nil {
			ignore := (f)(m2)
			if ignore {
				return errors.New("ignore this secret request")
			}
		}
		err = this.messageSecretRequest(m2)
	case *encoding.RevealSecret:
		this.raiden.db.NewReceivedRevealSecret(models.NewReceivedRevealSecret(m2, hash))
		f := this.raiden.RevealSecretListenerMap[m2.HashLock()]
		if f != nil {
			remove := (f)(m2)
			if remove {
				delete(this.raiden.RevealSecretListenerMap, m2.HashLock())
			}
		}
		err = this.messageRevealSecret(m2) //has no relation with statemanager,duplicate message will be ok
	case *encoding.Secret:
		err = this.messageSecret(m2)
	case *encoding.DirectTransfer:
		err = this.messageDirectTransfer(m2)
	case *encoding.MediatedTransfer:
		for f, _ := range this.raiden.ReceivedMediatedTrasnferListenerMap {
			remove := (*f)(m2)
			if remove {
				delete(this.raiden.ReceivedMediatedTrasnferListenerMap, f)
			}
		}
		err = this.MessageMediatedTransfer(m2)
	case *encoding.RefundTransfer:
		err = this.messageRefundTransfer(m2)
	case *encoding.RemoveExpiredHashlockTransfer:
		err = this.messageRemoveExpiredHashlockTransfer(m2)
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
		Message:      msger,
	}
	this.raiden.StateMachineEventHandler.LogAndDispatchByIdentifier(balanceProof.Identifier, balanceProof)
}
func (this *RaidenMessageHandler) messageRevealSecret(msg *encoding.RevealSecret) error {
	secret := msg.Secret
	sender := msg.Sender
	this.raiden.RegisterSecret(secret)
	stateChange := &mediated_transfer.ReceiveSecretRevealStateChange{secret, sender, msg}
	this.raiden.StateMachineEventHandler.LogAndDispatchToAllTasks(stateChange)
	return nil
}
func (this *RaidenMessageHandler) messageSecretRequest(msg *encoding.SecretRequest) error {
	stateChange := &mediated_transfer.ReceiveSecretRequestStateChange{
		Identifier: msg.Identifier,
		Amount:     new(big.Int).Set(msg.Amount),
		Hashlock:   msg.HashLock,
		Sender:     msg.Sender,
		Message:    msg,
	}
	this.raiden.StateMachineEventHandler.LogAndDispatchByIdentifier(msg.Identifier, stateChange)
	return nil
}
func (this *RaidenMessageHandler) markSecretComplete(msg *encoding.Secret) {
	if msg.Tag() == nil {
		log.Error(fmt.Sprintf("tag must not be nil ,only when token swap %s", utils.StringInterface(msg, 5)))
		return
	}
	tx := this.raiden.db.StartTx()
	msgTag := msg.Tag().(*transfer.MessageTag)
	mgr := msgTag.GetStateManager()

	if msgTag.ReceiveProcessComplete != false {
		/*
				todo must be solved
			When tokenswap is used as an intermediate node, ReceiveProcessComplete is true when it is supposed to be false. for event handler, receiveMessageTag.ReceiveProcessComplete = true
		*/
		//panic(fmt.Sprintf("ReceiveProcessComplete must be false, %s", utils.StringInterface(msg, 6)))
	}

	mgr.ManagerState = transfer.StateManager_ReceivedMessageProcessComplete
	log.Trace(fmt.Sprintf("markSecretComplete set message %s ReceiveProcessComplete", msgTag.MessageId))
	msgTag.ReceiveProcessComplete = true
	ack := this.raiden.Protocol.CreateAck(msgTag.EchoHash)
	this.raiden.db.SaveAck(msgTag.EchoHash, ack.Pack(), tx)
	_, ok := mgr.LastReceivedMessage.(*encoding.Secret)
	if !ok {
		panic("must be a secret message")
	}
	mgr.IsBalanceProofReceived = true
	if mgr.Name == target.NameTargetTransition {
		mgr.ManagerState = transfer.StateManager_TransferComplete
	} else if mgr.Name == initiator.NameInitiatorTransition {
		// initiator should not receive
	} else if mgr.Name == mediator.NameMediatorTransition {
		/*
			how to detect a mediator node is finish or not?
				1. receive prev balanceproof
				2. balanceproof  send to next successfully
			//todo when refund?
		*/
		if mgr.IsBalanceProofSent && mgr.IsBalanceProofReceived {
			mgr.ManagerState = transfer.StateManager_TransferComplete
		}
	}
	this.raiden.db.UpdateStateManaer(mgr, tx)
	if mgr.ChannelAddress == utils.EmptyAddress {
		panic("channeladdress must be valid")
	}
	if mgr.ChannelAddress != msg.Channel {
		log.Info(fmt.Sprintf("this is a secret message from refunded node %s", msg))
	}
	ch := this.raiden.GetChannelWithAddr(msg.Channel)
	this.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	tx.Commit()
	this.raiden.ConditionQuit("SecretSendAck")
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
		fmt.Sprintf("hashlock=%s,identifier=%s,nettingchannel=%s", utils.HPex(hashlock), identifer, nettingChannel)
		if !params.TreatRefundTransferAsNormalMediatedTransfer {
			this.raiden.HandleSecret(identifer, nettingChannel.TokenAddress, secret, msg, hashlock)
		} else {
			err = nettingChannel.RegisterTransfer(this.raiden.GetBlockNumber(), msg)
			if err != nil {
				log.Error(fmt.Sprintf("messageSecret RegisterTransfer err=%s", err))
			}
		}
	}
	//mark balanceproof complete
	this.markSecretComplete(msg)
	/*
		the following seems useless , remove it ,todo fix ,remove
	*/
	/*
		stateChange := &mediated_transfer.ReceiveSecretRevealStateChange{
			Secret:  secret,
			Sender:  msg.Sender,
			Message: nil, //
		}
		this.raiden.StateMachineEventHandler.LogAndDispatchByIdentifier(identifer, stateChange)
	*/
	return nil
}

/*
if there is any error, just ignore.
*/
func (this *RaidenMessageHandler) messageRemoveExpiredHashlockTransfer(msg *encoding.RemoveExpiredHashlockTransfer) error {
	this.balanceProof(msg)
	ch := this.raiden.GetChannelWithAddr(msg.Channel)
	if ch == nil {
		log.Warn("received  RemoveExpiredHashlockTransfer ,but relate channel cannot found %s", utils.StringInterface(msg, 7))
		return nil
	}
	err := ch.RegisterRemoveExpiredHashlockTransfer(msg, this.raiden.GetBlockNumber())
	if err != nil {
		log.Warn("RegisterRemoveExpiredHashlockTransfer err %s", err)
	}
	this.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
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
	transferState := &mediated_transfer.LockedTransferState{
		Identifier: msg.Identifier,
		Amount:     new(big.Int).Set(msg.Amount),
		Token:      msg.Token,
		Initiator:  msg.Initiator,
		Target:     msg.Target,
		Expiration: msg.Expiration,
		Hashlock:   msg.HashLock,
		Secret:     utils.EmptyHash,
		Fee:        msg.Fee,
	}
	stateChange := &mediated_transfer.ReceiveTransferRefundStateChange{
		Sender:   msg.Sender,
		Transfer: transferState,
		Message:  msg,
	}
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
		Message:      msg,
	}
	stateChangeId, err := this.raiden.db.LogStateChange(stateChange)
	if err != nil {
		return err
	}
	ch.RegisterTransfer(this.raiden.GetBlockNumber(), msg)
	receiveSuccess := &transfer.EventTransferReceivedSuccess{
		Identifier: msg.Identifier,
		Amount:     amount,
		Initiator:  msg.Sender,
	}
	err = this.raiden.db.LogEvents(stateChangeId, []transfer.Event{receiveSuccess}, this.raiden.GetBlockNumber())
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
	if tokenswap, ok := this.raiden.SwapKey2TokenSwap[key]; ok {
		this.messageTokenSwap(msg, tokenswap)
		//return nil
	}
	graph := this.raiden.GetToken2ChannelGraph(msg.Token)
	if graph == nil {
		return errors.New(fmt.Sprintf("received transfer on unkown token :%s", msg.Token.String()))
	}
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

/*
taker process token swap
*/
func (this *RaidenMessageHandler) messageTokenSwap(msg *encoding.MediatedTransfer, tokenswap *TokenSwap) {
	var hashlock common.Hash = msg.HashLock
	var hasReceiveRevealSecret bool
	var stateManager *transfer.StateManager
	if msg.Identifier != tokenswap.Identifier || msg.Amount.Cmp(tokenswap.FromAmount) != 0 || msg.Initiator != tokenswap.FromNodeAddress || msg.Token != tokenswap.FromToken || msg.Target != tokenswap.ToNodeAddress {
		log.Info("receive a mediated transfer, not match tokenswap condition")
		return
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
		initState, ok := state.(*mediated_transfer.InitiatorState)
		if !ok {
			panic(fmt.Sprintf("must be a InitiatorState"))
		}
		if initState.Transfer.Hashlock != msg.HashLock() {
			panic(fmt.Sprintf("hashlock must be same , state lock=%s,msg lock=%s", utils.HPex(initState.Transfer.Hashlock), utils.HPex(msg.HashLock())))
		}
		initState.Transfer.Secret = msg.Secret
		hasReceiveRevealSecret = true
		delete(this.raiden.SecretRequestPredictorMap, hashlock)
		return true
	}

	result, stateManager := this.raiden.StartTakerMediatedTransfer(tokenswap.ToToken, tokenswap.FromNodeAddress, tokenswap.ToAmount, tokenswap.Identifier, msg.HashLock, msg.Expiration)
	if stateManager == nil {
		log.Error(fmt.Sprintf("taker tokenwap error %s", <-result.Result))
		return
	}
	this.raiden.SecretRequestPredictorMap[hashlock] = secretRequestHook
	this.raiden.RevealSecretListenerMap[hashlock] = receiveRevealSecretHook
}
