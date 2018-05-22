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
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
 Class responsible to handle the protocol messages.

    Note:
        This class is not intended to be used standalone, use RaidenService
        instead.
*/
type raidenMessageHandler struct {
	raiden        *RaidenService
	blockedTokens map[common.Address]bool
}

func newRaidenMessageHandler(raiden *RaidenService) *raidenMessageHandler {
	h := &raidenMessageHandler{
		raiden:        raiden,
		blockedTokens: make(map[common.Address]bool),
	}
	return h
}

/*
 Handles `message` and sends an ACK on success.
*/
func (mh *raidenMessageHandler) onMessage(msg encoding.SignedMessager, hash common.Hash) (err error) {
	msg.SetTag(&transfer.MessageTag{
		EchoHash:          hash,
		IsASendingMessage: false,
		MessageId:         utils.RandomString(10),
	})
	switch m2 := msg.(type) {
	case *encoding.SecretRequest:
		f := mh.raiden.SecretRequestPredictorMap[m2.HashLock]
		if f != nil {
			ignore := (f)(m2)
			if ignore {
				return errors.New("ignore mh secret request")
			}
		}
		err = mh.messageSecretRequest(m2)
	case *encoding.RevealSecret:
		mh.raiden.db.NewReceivedRevealSecret(models.NewReceivedRevealSecret(m2, hash))
		f := mh.raiden.RevealSecretListenerMap[m2.HashLock()]
		if f != nil {
			remove := (f)(m2)
			if remove {
				delete(mh.raiden.RevealSecretListenerMap, m2.HashLock())
			}
		}
		err = mh.messageRevealSecret(m2) //has no relation with statemanager,duplicate message will be ok
	case *encoding.Secret:
		err = mh.messageSecret(m2)
	case *encoding.DirectTransfer:
		err = mh.messageDirectTransfer(m2)
	case *encoding.MediatedTransfer:
		err = mh.messageMediatedTransfer(m2)
		if err == nil {
			for f := range mh.raiden.ReceivedMediatedTrasnferListenerMap {
				remove := (*f)(m2)
				if remove {
					delete(mh.raiden.ReceivedMediatedTrasnferListenerMap, f)
				}
			}
		}
	case *encoding.RefundTransfer:
		err = mh.messageRefundTransfer(m2)
	case *encoding.RemoveExpiredHashlockTransfer:
		err = mh.messageRemoveExpiredHashlockTransfer(m2)
	default:
		log.Error(fmt.Sprintf("raidenMessageHandler unknown msg:%s", utils.StringInterface1(msg)))
		return fmt.Errorf("unhandled message cmdid:%d", msg.Cmd())
	}
	return err
}

func (mh *raidenMessageHandler) balanceProof(msger encoding.EnvelopMessager) {
	//blanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(msger)
	msg := msger.GetEnvelopMessage()
	balanceProof := &mediatedtransfer.ReceiveBalanceProofStateChange{
		Identifier:   msg.Identifier,
		NodeAddress:  msg.Sender,
		BalanceProof: transfer.NewBalanceProofStateFromEnvelopMessage(msger),
		Message:      msger,
	}
	mh.raiden.StateMachineEventHandler.logAndDispatchByIdentifier(balanceProof.Identifier, balanceProof)
}
func (mh *raidenMessageHandler) messageRevealSecret(msg *encoding.RevealSecret) error {
	secret := msg.Secret
	sender := msg.Sender
	mh.raiden.registerSecret(secret)
	stateChange := &mediatedtransfer.ReceiveSecretRevealStateChange{Secret: secret, Sender: sender, Message: msg}
	mh.raiden.StateMachineEventHandler.logAndDispatchToAllTasks(stateChange)
	return nil
}
func (mh *raidenMessageHandler) messageSecretRequest(msg *encoding.SecretRequest) error {
	stateChange := &mediatedtransfer.ReceiveSecretRequestStateChange{
		Identifier: msg.Identifier,
		Amount:     new(big.Int).Set(msg.Amount),
		Hashlock:   msg.HashLock,
		Sender:     msg.Sender,
		Message:    msg,
	}
	mh.raiden.StateMachineEventHandler.logAndDispatchByIdentifier(msg.Identifier, stateChange)
	return nil
}
func (mh *raidenMessageHandler) markSecretComplete(msg *encoding.Secret) {
	if msg.Tag() == nil {
		log.Error(fmt.Sprintf("tag must not be nil ,only when token swap %s", utils.StringInterface(msg, 5)))
		return
	}
	tx := mh.raiden.db.StartTx()
	msgTag := msg.Tag().(*transfer.MessageTag)
	mgr := msgTag.GetStateManager()

	if msgTag.ReceiveProcessComplete != false {
		/*
				todo must be solved
			When tokenswap is used as an intermediate node, ReceiveProcessComplete is true when it is supposed to be false. for event handler, receiveMessageTag.ReceiveProcessComplete = true
		*/
		//panic(fmt.Sprintf("ReceiveProcessComplete must be false, %s", utils.StringInterface(msg, 6)))
	}

	mgr.ManagerState = transfer.StateManagerReceivedMessageProcessComplete
	log.Trace(fmt.Sprintf("markSecretComplete set message %s ReceiveProcessComplete", msgTag.MessageId))
	msgTag.ReceiveProcessComplete = true
	ack := mh.raiden.Protocol.CreateAck(msgTag.EchoHash)
	mh.raiden.db.SaveAck(msgTag.EchoHash, ack.Pack(), tx)
	_, ok := mgr.LastReceivedMessage.(*encoding.Secret)
	if !ok {
		panic("must be a secret message")
	}
	mgr.IsBalanceProofReceived = true
	if mgr.Name == target.NameTargetTransition {
		mgr.ManagerState = transfer.StateManagerTransferComplete
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
			mgr.ManagerState = transfer.StateManagerTransferComplete
		}
	}
	mh.raiden.db.UpdateStateManaer(mgr, tx)
	if mgr.ChannelAddress == utils.EmptyAddress {
		panic("channeladdress must be valid")
	}
	if mgr.ChannelAddress != msg.Channel {
		log.Info(fmt.Sprintf("mh is a secret message from refunded node %s", msg))
	}
	ch, err := mh.raiden.findChannelByAddress(msg.Channel)
	if err != nil {
		panic(fmt.Sprintf("channel %s must exists", utils.APex(msg.Channel)))
	}
	mh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	tx.Commit()
	mh.raiden.conditionQuit("SecretSendAck")
}
func (mh *raidenMessageHandler) messageSecret(msg *encoding.Secret) error {
	mh.balanceProof(msg)
	hashlock := msg.HashLock()
	identifer := msg.Identifier
	secret := msg.Secret
	mh.raiden.registerSecret(secret)
	var nettingChannel *channel.Channel
	var err error
	nettingChannel, err = mh.raiden.findChannelByAddress(msg.Channel)
	if err != nil {
		log.Info(fmt.Sprintf("Message for unknown channel: %s", err))
	} else {
		log.Trace(fmt.Sprintf("hashlock=%s,identifier=%d,nettingchannel=%s", utils.HPex(hashlock), identifer, nettingChannel))
		if !params.TreatRefundTransferAsNormalMediatedTransfer {
			mh.raiden.handleSecret(identifer, nettingChannel.TokenAddress, secret, msg, hashlock)
		} else {
			err = nettingChannel.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
			if err != nil {
				log.Error(fmt.Sprintf("messageSecret RegisterTransfer err=%s", err))
			}
		}
	}
	//mark balanceproof complete
	mh.markSecretComplete(msg)
	return nil
}

/*
if there is any error, just ignore.
*/
func (mh *raidenMessageHandler) messageRemoveExpiredHashlockTransfer(msg *encoding.RemoveExpiredHashlockTransfer) error {
	mh.balanceProof(msg)
	ch, err := mh.raiden.findChannelByAddress(msg.Channel)
	if err != nil {
		log.Warn("received  RemoveExpiredHashlockTransfer ,but relate channel cannot found %s", utils.StringInterface(msg, 7))
		return nil
	}
	err = ch.RegisterRemoveExpiredHashlockTransfer(msg, mh.raiden.GetBlockNumber())
	if err != nil {
		log.Warn("RegisterRemoveExpiredHashlockTransfer err %s", err)
	}
	mh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	return nil
}
func (mh *raidenMessageHandler) messageRefundTransfer(msg *encoding.RefundTransfer) (err error) {
	mh.balanceProof(msg)
	graph := mh.raiden.getToken2ChannelGraph(msg.Token)
	if graph == nil {
		return rerr.UnknownTokenAddress(msg.Token.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("Direct transfer from node without an existing channel: %s", msg.Sender)
		return
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(msg.Token), utils.APex2(msg.Sender)))
	}
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return
	}
	transferState := &mediatedtransfer.LockedTransferState{
		Identifier:   msg.Identifier,
		TargetAmount: big.NewInt(0).Sub(msg.Amount, msg.Fee),
		Amount:       new(big.Int).Set(msg.Amount),
		Token:        msg.Token,
		Initiator:    msg.Initiator,
		Target:       msg.Target,
		Expiration:   msg.Expiration,
		Hashlock:     msg.HashLock,
		Secret:       utils.EmptyHash,
		Fee:          msg.Fee,
	}
	stateChange := &mediatedtransfer.ReceiveTransferRefundStateChange{
		Sender:   msg.Sender,
		Transfer: transferState,
		Message:  msg,
	}
	mh.raiden.StateMachineEventHandler.logAndDispatchByIdentifier(msg.Identifier, stateChange)
	return nil
}

func (mh *raidenMessageHandler) messageDirectTransfer(msg *encoding.DirectTransfer) error {
	mh.balanceProof(msg)
	if graph := mh.raiden.getToken2ChannelGraph(msg.Token); graph == nil {
		return rerr.UnknownTokenAddress(msg.Token.String())
	}
	if _, ok := mh.blockedTokens[msg.Token]; ok {
		return rerr.TransferUnwanted
	}
	graph := mh.raiden.getToken2ChannelGraph(msg.Token)
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		return rerr.UnknownAddress(fmt.Sprintf("Direct transfer from node without an existing channel partner %s  ", msg.Sender))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(msg.Token), utils.APex2(msg.Sender)))
	}
	if ch.State() != transfer.ChannelStateOpened {
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
	stateChangeID, err := mh.raiden.db.LogStateChange(stateChange)
	if err != nil {
		return err
	}
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		log.Error("RegisterTransfer error %s\n", msg)
		return err
	}
	receiveSuccess := &transfer.EventTransferReceivedSuccess{
		Identifier: msg.Identifier,
		Amount:     amount,
		Initiator:  msg.Sender,
	}
	err = mh.raiden.db.LogEvents(stateChangeID, []transfer.Event{receiveSuccess}, mh.raiden.GetBlockNumber())
	return err
}

func (mh *raidenMessageHandler) messageMediatedTransfer(msg *encoding.MediatedTransfer) error {
	if mh.raiden.Config.IgnoreMediatedNodeRequest && msg.Target != mh.raiden.NodeAddress {
		return fmt.Errorf("ignored mh mediated transfer, because i don't want to route ")
	}
	mh.balanceProof(msg)
	//  TODO: Reject mediated transfer that the hashlock/identifier is known,
	// mh is a downstream bug and the transfer is going in cycles (issue #490)
	if _, ok := mh.blockedTokens[msg.Token]; ok {
		return rerr.TransferUnwanted
	}
	graph := mh.raiden.getToken2ChannelGraph(msg.Token)
	if graph == nil {
		return fmt.Errorf("received transfer on unkown token :%s", msg.Token.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		return rerr.ChannelNotFound(fmt.Sprintf("mediated transfer from node without an existing channel %s", msg.Sender))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(msg.Token), utils.APex2(msg.Sender)))
	}
	if ch.State() != transfer.ChannelStateOpened {
		return rerr.TransferWhenClosed(fmt.Sprintf("Mediated transfer received but the channel is closed %s", ch.MyAddress))
	}
	err := ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return err
	}
	if msg.Target == mh.raiden.NodeAddress {
		mh.raiden.targetMediatedTransfer(msg)
	} else {
		mh.raiden.mediateMediatedTransfer(msg)
	}
	/*
		start  taker's tokenswap ,only if receive a valid mediated transfer
	*/
	key := swapKey{msg.Identifier, msg.Token, msg.Amount.String()}
	if tokenswap, ok := mh.raiden.SwapKey2TokenSwap[key]; ok {
		remove := mh.raiden.messageTokenSwapTaker(msg, tokenswap)
		if remove { //once the swap start,remove mh key immediately. otherwise,maker may repeat mh tokenswap operation.
			delete(mh.raiden.SwapKey2TokenSwap, key)
		}
		//return nil
	}
	return nil
}
