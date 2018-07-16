package smartraiden

import (
	"fmt"

	"math/big"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/target"
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
		MessageID:         utils.RandomString(10),
	})
	switch m2 := msg.(type) {
	case *encoding.SecretRequest:
		f := mh.raiden.SecretRequestPredictorMap[m2.LockSecretHash]
		if f != nil {
			ignore := (f)(m2)
			if ignore {
				return errors.New("ignore this secret request")
			}
		}
		err = mh.messageSecretRequest(m2)
	case *encoding.RevealSecret:
		mh.raiden.db.NewReceivedRevealSecret(models.NewReceivedRevealSecret(m2, hash))
		f := mh.raiden.RevealSecretListenerMap[m2.LockSecretHash()]
		if f != nil {
			remove := (f)(m2)
			if remove {
				delete(mh.raiden.RevealSecretListenerMap, m2.LockSecretHash())
			}
		}
		err = mh.messageRevealSecret(m2) //has no relation with statemanager,duplicate message will be ok
	case *encoding.UnLock:
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
	case *encoding.AnnounceDisposed:
		err = mh.messageAnnounceDisposed(m2)
	case *encoding.RemoveExpiredHashlockTransfer:
		err = mh.messageRemoveExpiredHashlockTransfer(m2)
	default:
		log.Error(fmt.Sprintf("raidenMessageHandler unknown msg:%s", utils.StringInterface1(msg)))
		return fmt.Errorf("unhandled message cmdid:%d", msg.Cmd())
	}
	return err
}

//这个到底有什么用啊?看不懂 todo 是否可以移除呢?
func (mh *raidenMessageHandler) balanceProof(msger encoding.EnvelopMessager) {
	//blanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(msger)
	//msg := msger.GetEnvelopMessage()
	//balanceProof := &mediatedtransfer.ReceiveBalanceProofStateChange{
	//	LockSecretHash: msg.L,
	//	NodeAddress:    msg.Sender,
	//	BalanceProof:   transfer.NewBalanceProofStateFromEnvelopMessage(msger),
	//	Message:        msger,
	//}
	//mh.raiden.StateMachineEventHandler.logAndDispatchByIdentifier(balanceProof.LockSecretHash, balanceProof)
}
func (mh *raidenMessageHandler) messageRevealSecret(msg *encoding.RevealSecret) error {
	secret := msg.LockSecret
	sender := msg.Sender
	mh.raiden.registerSecret(secret)
	stateChange := &mediatedtransfer.ReceiveSecretRevealStateChange{Secret: secret, Sender: sender, Message: msg}
	mh.raiden.StateMachineEventHandler.logAndDispatchToAllTasks(stateChange)
	return nil
}
func (mh *raidenMessageHandler) messageSecretRequest(msg *encoding.SecretRequest) error {
	stateChange := &mediatedtransfer.ReceiveSecretRequestStateChange{
		Amount:         new(big.Int).Set(msg.PaymentAmount),
		LockSecretHash: msg.LockSecretHash,
		Sender:         msg.Sender,
		Message:        msg,
	}
	mh.raiden.StateMachineEventHandler.logAndDispatchByIdentifier(stateChange.LockSecretHash, stateChange)
	return nil
}
func (mh *raidenMessageHandler) markSecretComplete(msg *encoding.UnLock) {
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
	log.Trace(fmt.Sprintf("markSecretComplete set message %s ReceiveProcessComplete", msgTag.MessageID))
	msgTag.ReceiveProcessComplete = true
	ack := mh.raiden.Protocol.CreateAck(msgTag.EchoHash)
	mh.raiden.db.SaveAck(msgTag.EchoHash, ack.Pack(), tx)
	_, ok := mgr.LastReceivedMessage.(*encoding.UnLock)
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
	if mgr.ChannelAddress == utils.EmptyHash {
		panic("channeladdress must be valid")
	}
	if mgr.ChannelAddress != msg.ChannelIdentifier {
		log.Info(fmt.Sprintf("mh is a secret message from refunded node %s", msg))
	}
	ch, err := mh.raiden.findChannelByAddress(msg.ChannelIdentifier)
	if err != nil {
		panic(fmt.Sprintf("channel %s must exists", utils.HPex(msg.ChannelIdentifier)))
	}
	mh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	tx.Commit()
	mh.raiden.conditionQuit("SecretSendAck")
}
func (mh *raidenMessageHandler) messageSecret(msg *encoding.UnLock) error {
	mh.balanceProof(msg)
	lockSecretHash := msg.LockSecretHash()
	secret := msg.LockSecret
	mh.raiden.registerSecret(secret)
	var nettingChannel *channel.Channel
	var err error
	nettingChannel, err = mh.raiden.findChannelByAddress(msg.ChannelIdentifier)
	if err != nil {
		log.Info(fmt.Sprintf("Message for unknown channel: %s", err))
	} else {
		log.Trace(fmt.Sprintf("lockSecretHash=%s,nettingchannel=%s", utils.HPex(lockSecretHash), nettingChannel))
		err = nettingChannel.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
		if err != nil {
			log.Error(fmt.Sprintf("messageSecret RegisterTransfer err=%s", err))
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
	ch, err := mh.raiden.findChannelByAddress(msg.ChannelIdentifier)
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
func (mh *raidenMessageHandler) messageAnnounceDisposed(msg *encoding.AnnounceDisposed) (err error) {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unkonwn channel %s", msg.ChannelIdentifier.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("Direct transfer from node without an existing channel: %s", msg.Sender)
		return
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("channel:%s,partner:%s", utils.HPex(msg.ChannelIdentifier), utils.APex2(msg.Sender)))
	}
	err = ch.RegisterAnnouceDisposed(msg)
	if err != nil {
		return
	}
	punish := models.NewReceivedAnnounceDisposed(msg.Lock.Hash(), msg.ChannelIdentifier, msg.GetAdditionalHash(), msg.OpenBlockNumber, msg.Signature)
	err = mh.raiden.db.MarkLockHashCanPunish(punish)
	if err != nil {
		err = fmt.Errorf("MarkLockHashCanPunish %s err %s", utils.StringInterface(punish, 2), err)
		return
	}
	stateChange := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Sender:  msg.Sender,
		Token:   ch.TokenAddress,
		Lock:    msg.Lock,
		Message: msg,
	}
	mh.raiden.StateMachineEventHandler.logAndDispatchByIdentifier(msg.Lock.LockSecretHash, stateChange)
	return nil
}
func (mh *raidenMessageHandler) messageAnnounceDisposedResponse(msg *encoding.AnnounceDisposedResponse) (err error) {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unkonwn channel %s", msg.ChannelIdentifier.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("Direct transfer from node without an existing channel: %s", msg.Sender)
		return
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("channel:%s,partner:%s", utils.HPex(msg.ChannelIdentifier), utils.APex2(msg.Sender)))
	}
	/*
		必须验证我确实发送过这个Dispose
	*/
	b := mh.raiden.db.IsLockSecretHashChannelIdentifierDisposed(msg.LockSecretHash, msg.ChannelIdentifier)
	if !b {
		return fmt.Errorf("maybe a attach, receive a announce disposed response,but i never send announce disposed,msg=%s", msg)
	}
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return
	}
	//保存通道状态即可.
	return nil
}

/*
for direct transfer,
we shoud make ack and update channel status ,these two operations atomic
*/
func (mh *raidenMessageHandler) markDirectTransferComplete(msg *encoding.DirectTransfer) {
	if msg.Tag() == nil {
		log.Error(fmt.Sprintf("tag must not be nil ,only when token swap %s", utils.StringInterface(msg, 5)))
		return
	}
	tx := mh.raiden.db.StartTx()
	msgTag := msg.Tag().(*transfer.MessageTag)
	ack := mh.raiden.Protocol.CreateAck(msgTag.EchoHash)
	mh.raiden.db.SaveAck(msgTag.EchoHash, ack.Pack(), tx)
	ch, err := mh.raiden.findChannelByAddress(msg.ChannelIdentifier)
	if err != nil {
		panic(fmt.Sprintf("channel %s must exists", utils.HPex(msg.ChannelIdentifier)))
	}
	mh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	tx.Commit()
	mh.raiden.conditionQuit("DirectTransferSendAck")
}
func (mh *raidenMessageHandler) messageDirectTransfer(msg *encoding.DirectTransfer) error {
	mh.balanceProof(msg)
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	if _, ok := mh.blockedTokens[token]; ok {
		return rerr.ErrTransferUnwanted
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return rerr.TransferWhenClosed(ch.ChannelIdentifier.String())
	}
	var amount = new(big.Int)
	amount = amount.Sub(msg.TransferAmount, ch.PartnerState.TransferAmount())
	err := ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterTransfer error %s\n", msg))
		return err
	}
	mh.markDirectTransferComplete(msg)
	receiveSuccess := &transfer.EventTransferReceivedSuccess{
		Amount:            amount,
		Initiator:         msg.Sender,
		ChannelIdentifier: msg.ChannelIdentifier,
	}
	err = mh.raiden.StateMachineEventHandler.OnEvent(receiveSuccess, nil)
	return err
}

func (mh *raidenMessageHandler) messageMediatedTransfer(msg *encoding.MediatedTransfer) error {
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if mh.raiden.Config.IgnoreMediatedNodeRequest && msg.Target != mh.raiden.NodeAddress {
		return fmt.Errorf("ignored mh mediated transfer, because i don't want to route ")
	}
	if mh.raiden.Config.IsMeshNetwork {
		return fmt.Errorf("deny any mediated transfer when there is no internet connection")
	}
	mh.balanceProof(msg)
	//  TODO: Reject mediated transfer that the hashlock/identifier is known,
	// mh is a downstream bug and the transfer is going in cycles (issue #490)
	if _, ok := mh.blockedTokens[token]; ok {
		return rerr.ErrTransferUnwanted
	}
	graph := mh.raiden.getToken2ChannelGraph(token)
	if graph == nil {
		return fmt.Errorf("received transfer on unkown token :%s", utils.APex2(token))
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		return rerr.ChannelNotFound(fmt.Sprintf("mediated transfer from node without an existing channel %s", msg.Sender))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return rerr.TransferWhenClosed(fmt.Sprintf("Mediated transfer received but the channel is closed %s", ch.ChannelIdentifier))
	}
	err := ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return err
	}
	if msg.Target == mh.raiden.NodeAddress {
		mh.raiden.targetMediatedTransfer(msg, ch)
	} else {
		mh.raiden.mediateMediatedTransfer(msg, ch)
	}
	/*
		start  taker's tokenswap ,only if receive a valid mediated transfer
	*/
	key := swapKey{msg.LockSecretHash, token, msg.PaymentAmount.String()}
	if tokenswap, ok := mh.raiden.SwapKey2TokenSwap[key]; ok {
		remove := mh.raiden.messageTokenSwapTaker(msg, tokenswap)
		if remove { //once the swap start,remove mh key immediately. otherwise,maker may repeat mh tokenswap operation.
			delete(mh.raiden.SwapKey2TokenSwap, key)
		}
		//return nil
	}
	return nil
}
