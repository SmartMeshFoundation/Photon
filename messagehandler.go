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
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
 Class responsible to handle the protocol messages.

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
		EchoHash: hash,
	})
	switch m2 := msg.(type) {
	case *encoding.SecretRequest:
		f := mh.raiden.SecretRequestPredictorMap[m2.LockSecretHash]
		if f != nil {
			ignore := (f)(m2)
			if ignore {
				return errors.New("ignore this secret request,because of SecretRequestPredictorMap ignores")
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
		err = mh.messageUnlock(m2)
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
	case *encoding.AnnounceDisposedResponse:
		err = mh.messageAnnounceDisposedResponse(m2)
	case *encoding.RemoveExpiredHashlockTransfer:
		err = mh.messageRemoveExpiredHashlockTransfer(m2)
	case *encoding.SettleRequest:
		err = mh.messageSettleRequest(m2)
	case *encoding.SettleResponse:
		err = mh.messageSettleResponse(m2)
	case *encoding.WithdrawRequest:
		err = mh.messageWithdrawRequest(m2)
	case *encoding.WithdrawResponse:
		err = mh.messageWithdrawResponse(m2)
	default:
		log.Error(fmt.Sprintf("raidenMessageHandler unknown msg:%s", utils.StringInterface1(msg)))
		return fmt.Errorf("unhandled message cmdid:%d", msg.Cmd())
	}
	return err
}

//这个到底有什么用啊?看不懂
func (mh *raidenMessageHandler) balanceProof(msg *encoding.UnLock) {
	blanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(msg)
	balanceProof := &mediatedtransfer.ReceiveBalanceProofStateChange{
		LockSecretHash: msg.LockSecretHash(),
		NodeAddress:    msg.Sender,
		BalanceProof:   blanceProof,
		Message:        msg,
	}
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(balanceProof.LockSecretHash, balanceProof)
}
func (mh *raidenMessageHandler) messageRevealSecret(msg *encoding.RevealSecret) error {
	secret := msg.LockSecret
	sender := msg.Sender
	mh.raiden.registerSecret(secret)
	stateChange := &mediatedtransfer.ReceiveSecretRevealStateChange{Secret: secret, Sender: sender, Message: msg}

	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(msg.LockSecretHash(), stateChange)
	return nil
}
func (mh *raidenMessageHandler) messageSecretRequest(msg *encoding.SecretRequest) error {
	stateChange := &mediatedtransfer.ReceiveSecretRequestStateChange{
		Amount:         new(big.Int).Set(msg.PaymentAmount),
		LockSecretHash: msg.LockSecretHash,
		Sender:         msg.Sender,
		Message:        msg,
	}
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(stateChange.LockSecretHash, stateChange)
	return nil
}
func (mh *raidenMessageHandler) messageUnlock(msg *encoding.UnLock) error {
	mh.balanceProof(msg)
	lockSecretHash := msg.LockSecretHash()
	secret := msg.LockSecret
	mh.raiden.registerSecret(secret)
	var ch *channel.Channel
	var err error
	ch, err = mh.raiden.findChannelByAddress(msg.ChannelIdentifier)
	if err != nil {
		log.Info(fmt.Sprintf("Message for unknown channel: %s", err))
		return err
	}
	log.Trace(fmt.Sprintf("lockSecretHash=%s,nettingchannel=%s", utils.HPex(lockSecretHash), ch))
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		log.Error(fmt.Sprintf("messageUnlock RegisterTransfer err=%s", err))
		return err
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}
func (mh *raidenMessageHandler) updateChannelAndSaveAck(c *channel.Channel, tag interface{}) {
	t, ok := tag.(*transfer.MessageTag)
	if !ok || t == nil {
		panic("tag is nil")
	}
	echohash := t.EchoHash
	ack := mh.raiden.Protocol.CreateAck(echohash)
	err := mh.raiden.db.UpdateChannelAndSaveAck(channel.NewChannelSerialization(c), echohash, ack.Pack())
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelAndSaveAck %s", err))
	}
}

/*
if there is any error, just ignore.
*/
func (mh *raidenMessageHandler) messageRemoveExpiredHashlockTransfer(msg *encoding.RemoveExpiredHashlockTransfer) error {
	ch, err := mh.raiden.findChannelByAddress(msg.ChannelIdentifier)
	if err != nil {
		log.Warn("received  RemoveExpiredHashlockTransfer ,but relate channel cannot found %s", utils.StringInterface(msg, 7))
		return nil
	}
	if !ch.CanContinueTransfer() {
		log.Warn(fmt.Sprintf("receive msg %s, but channel cannot continue transfer", msg))
		return nil
	}
	err = ch.RegisterRemoveExpiredHashlockTransfer(msg, mh.raiden.GetBlockNumber())
	if err != nil {
		log.Warn("RegisterRemoveExpiredHashlockTransfer err %s", err)
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}
func (mh *raidenMessageHandler) messageAnnounceDisposed(msg *encoding.AnnounceDisposed) (err error) {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unkonwn channel %s", msg.ChannelIdentifier.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("direct transfer from node without an existing channel: %s", msg.Sender)
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
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(msg.Lock.LockSecretHash, stateChange)
	return nil
}
func (mh *raidenMessageHandler) messageAnnounceDisposedResponse(msg *encoding.AnnounceDisposedResponse) (err error) {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unkonwn channel %s", msg.ChannelIdentifier.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("direct transfer from node without an existing channel: %s", msg.Sender)
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
		return fmt.Errorf("maybe a attack, receive a announce disposed response,but i never send announce disposed,msg=%s", msg)
	}
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return
	}
	//保存通道状态即可.
	return mh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
}

func (mh *raidenMessageHandler) messageDirectTransfer(msg *encoding.DirectTransfer) error {
	//mh.balanceProof(msg)
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
	receiveSuccess := &transfer.EventTransferReceivedSuccess{
		Amount:            amount,
		Initiator:         msg.Sender,
		ChannelIdentifier: msg.ChannelIdentifier,
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
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
	if _, ok := mh.blockedTokens[token]; ok {
		return rerr.ErrTransferUnwanted
	}
	graph := mh.raiden.getToken2ChannelGraph(token)
	if graph == nil {
		return fmt.Errorf("received transfer on unkown token :%s", utils.APex2(token))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if !ch.CanTransfer() {
		return rerr.TransferWhenClosed(fmt.Sprintf("Mediated transfer received but the channel is  can not accept any transfer %s", ch.ChannelIdentifier.String()))
	}
	err := ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return err
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	if msg.Target == mh.raiden.NodeAddress {
		mh.raiden.targetMediatedTransfer(msg, ch)
	} else {
		mh.raiden.mediateMediatedTransfer(msg, ch)
	}
	/*
		start  taker's tokenswap ,only if receive a valid mediated transfer
	*/
	key := swapKey{
		LockSecretHash: msg.LockSecretHash,
		FromToken:      token,
		FromAmount:     msg.PaymentAmount.String(),
	}
	if tokenswap, ok := mh.raiden.SwapKey2TokenSwap[key]; ok {
		remove := mh.raiden.messageTokenSwapTaker(msg, tokenswap)
		if remove { //once the swap start,remove mh key immediately. otherwise,maker may repeat mh tokenswap operation.
			delete(mh.raiden.SwapKey2TokenSwap, key)
		}
		//return nil
	}
	return nil
}

func (mh *raidenMessageHandler) messageSettleRequest(msg *encoding.SettleRequest) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return fmt.Errorf("receive settle request but channel state is %s", ch.State)
	}
	err := ch.RegisterCooperativeSettleRequest(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterCooperativeSettleRequest error %s\n", err))
		return err
	}
	settleResponse, err := ch.CreateCooperativeSettleResponse(msg)
	if err != nil {
		//if err, channel can only be closed /settled
		log.Error(fmt.Sprintf("CreateCooperativeSettleResponse err %s", err))
		return err
	}
	if ch.HasAnyUnkonwnSecretTransferOnRoad() {
		//我自己理解 withdraw on channel就可以,防止上一笔交易额外损失
		result := ch.CooperativeSettleChannelOnRequest(msg.Participant1Signature, settleResponse)
		go func() {
			var err2 error
			err2 = <-result.Result
			if err2 != nil {
				log.Error(fmt.Sprintf("CooperativeSettleChannelOnRequest err %s", err2))
			} else {
				log.Info(fmt.Sprintf("CooperativeSettleChannelOnRequest success on channel %s", ch.ChannelIdentifier.String()))
			}
		}()
		return nil
	}
	err = settleResponse.Sign(mh.raiden.PrivateKey, settleResponse)
	if err != nil {
		panic(fmt.Sprintf("sign message for settle response err %s", err))
	}
	err = mh.raiden.sendAsync(msg.Sender, settleResponse)
	if err != nil {
		log.Error(fmt.Sprintf("send message %s, to %s ,err %s", settleResponse, msg.Sender, err))
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}
func (mh *raidenMessageHandler) messageSettleResponse(msg *encoding.SettleResponse) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	/*
			会不会出现碰巧双方都发出了 settle request 这种情况?
			比如一方提出 settle, 另一方提出 withdraw,
			极小概率可能出现,如果出现,就退回到原始的 close/settle 模式
		也要防止另一方主动发出 settle response, 而我并没有发出 settle request 请求这种情况.
	*/
	if ch.State != channeltype.StateCooprativeSettle {
		return fmt.Errorf("receive settle response but channel state is %s", ch.State)
	}
	err := ch.RegisterCooperativeSettleResponse(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterCooperativeSettleResponse error %s\n", err))
		return err
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	result := ch.CooperativeSettleChannel(msg)
	go func() {
		err = <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("CooperativeSettleChannel %s failed, so we can only close/settle this channel", utils.HPex(msg.ChannelIdentifier)))
		}
	}()
	return nil
}
func (mh *raidenMessageHandler) messageWithdrawRequest(msg *encoding.WithdrawRequest) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return fmt.Errorf("receive settle request but channel state is %s", ch.State)
	}
	err := ch.RegisterWithdrawRequest(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterWithdrawRequest error %s\n", err))
		return err
	}
	//这里需要询问用户我想取现多少,现在默认用0来替代.
	withdrawResponse, err := ch.CreateWithdrawResponse(msg, utils.BigInt0)
	if err != nil {
		//if err, channel can only be closed /settled
		log.Error(fmt.Sprintf("CreateWithdrawResponse err %s", err))
		return err
	}
	if ch.HasAnyUnkonwnSecretTransferOnRoad() {
		//我自己理解 withdraw on channel就可以,防止上一笔交易额外损失
		result := ch.WithdrawOnRequest(msg.Participant1Signature, withdrawResponse)
		go func() {
			var err2 error
			err2 = <-result.Result
			if err2 != nil {
				log.Error(fmt.Sprintf("WithdrawOnRequest err %s", err2))
			} else {
				log.Info(fmt.Sprintf("WithdrawOnRequest success on channel %s", ch.ChannelIdentifier.String()))
			}
		}()
		return nil
	}
	err = withdrawResponse.Sign(mh.raiden.PrivateKey, withdrawResponse)
	if err != nil {
		panic(fmt.Sprintf("sign message for withdraw response err %s", err))
	}
	err = mh.raiden.sendAsync(msg.Sender, withdrawResponse)
	if err != nil {
		log.Error(fmt.Sprintf("send message %s, to %s ,err %s", withdrawResponse, msg.Sender, err))
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}
func (mh *raidenMessageHandler) messageWithdrawResponse(msg *encoding.WithdrawResponse) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateWithdraw {
		return fmt.Errorf("receive settle request but channel state is %s", ch.State)
	}
	/*
		要先验证一下我发出去了 withdraw request,并且金额正确,然后才能注册
	*/
	err := ch.RegisterWithdrawResponse(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterTransfer error %s\n", msg))
		return err
	}
	mh.updateChannelAndSaveAck(ch, msg.Tag())
	//如果碰巧崩溃了,如果失败了,都只能回到 close/settle 这种老办法.
	result := ch.Withdraw(msg)
	go func() {
		err = <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("Withdraw %s failed, so we can only close/settle this channel", msg.ChannelIdentifier.String()))
		}
	}()
	return nil
}
