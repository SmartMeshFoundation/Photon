package photon

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/params"

	"errors"

	"time"

	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network/graph"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/notify"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//run inside loop of photon service
type stateMachineEventHandler struct {
	photon                             *Service
	noEffectiveChainNotifyLoopQuitChan chan *struct{}
}

func newStateMachineEventHandler(photon *Service) *stateMachineEventHandler {
	h := &stateMachineEventHandler{
		photon: photon,
	}
	return h
}

/*
dispatch it to all state managers and log generated events
*/
func (eh *stateMachineEventHandler) dispatchToAllTasks(st transfer.StateChange) {
	for _, mgrs := range eh.photon.Transfer2StateManager {
		eh.dispatch(mgrs, st)
	}
}

/*
dispatch it to the state manager corresponding to `lockSecretHash`
*/
func (eh *stateMachineEventHandler) dispatchBySecretHash(lockSecretHash common.Hash, st transfer.StateChange) {
	for _, mgr := range eh.photon.Transfer2StateManager {
		//todo 这个未必是高效的方式,因为同时进行的 transfer 可能很多,会比较慢.
		if mgr.Identifier == lockSecretHash {
			eh.dispatch(mgr, st)
		}
	}
}

func (eh *stateMachineEventHandler) dispatchByPendingLocksInChannel(channel *channel.Channel, st transfer.StateChange) {
	for lockSecretHash := range channel.OurState.Lock2PendingLocks {
		eh.dispatchBySecretHash(lockSecretHash, st)
	}
}

func (eh *stateMachineEventHandler) dispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) (events []transfer.Event) {
	eh.updateStateManagerFromStateChange(stateManager, stateChange)
	events = stateManager.Dispatch(stateChange)
	for _, e := range events {
		err := eh.OnEvent(e, stateManager)
		if err != nil {
			log.Error(fmt.Sprintf("stateMachineEventHandler dispatch:%v\n", err))
		}
	}
	return
}

/*
我要发送 reveal secret 出去了,应该让每个与密码相关的通道都知道密码.
1.如果我是发送方,多注册一个密码没坏处
2. 如果我参与到的的是 token swap, 那么必须在告诉对方密码的同时,我要知道密码,虽然我自己知道,但必须注册到响应的通道中,
 比如 A-B-C A-D-C ,AC进行 token swap,A 是 maker,C 是 taker,A 给 C 的 expiration 是1000,C 给 A 的 expiration 是500,
那么 C 完全可以不理 A 发出的 secret request,然后在500块以后把 A 给 C 的 token 取走.
*/
/*
 *	eventSendRevealSecret : function to send RevealSecret to all channel participants.
 *
 *	Note that all channel participants in this payment network should refer this secret to their partner once they received the secret.
 *	1. If one participant is sender, repeatedly register the secret has no impact.
 *	2. If he participates in tokenswap, then he must register the secret in responding channel at the time he sends out secret to his partner.
 *		for example, A-B-C A-D-C, AC undergoes token swap, A is the maker, C is the taker, expiration block A sets to C is 1000, expiration block C sets to A is 500.
 *		then C can neglect secret request sent from A then steal those tokens C deposited after 500 block number.
 */
func (eh *stateMachineEventHandler) eventSendRevealSecret(event *mediatedtransfer.EventSendRevealSecret, stateManager *transfer.StateManager) (err error) {
	eh.photon.conditionQuit("EventSendRevealSecretBefore")
	/*
			有三种情况发送RevealSecret
			1.我是交易发起方,我需要给target发送密码
			2.我是交易接收方,我需要给上家发送密码,换取unlock
			3.我是中间结点,我需要给上家发送密码,换区unlock
			这时候的崩溃处理就有点麻烦,
			假设1-2-3交易,3发送SecretRequest后崩溃,
			这时候实际上3重启以后,因为没收到reveal secret,所以2-3之间锁是可以移除的.同样1-2之间也是可以移除的.
		但是因为这里使用的是registerSecret,导致1-2之间无法移除,
		之所以采用向所有通道注册密码而不是只向指定通道注册密码,那是考虑到token swap的时候,maker需要在这里通知所有通道知道密码.
		todo 暂时采用这种方式,后续token swap maker应该自行处理通知相关通道密码而不是放在这里.
	*/
	eh.photon.registerSecret(event.Secret)

	revealMessage := encoding.NewRevealSecret(event.Secret)
	// 带上交易附加信息
	revealMessage.Data = []byte(event.Data)
	err = revealMessage.Sign(eh.photon.PrivateKey, revealMessage)
	err = eh.photon.sendAsync(event.Receiver, revealMessage) //单独处理 reaveal secret
	if err == nil {
		std := eh.photon.dao.UpdateSentTransferDetailStatus(event.Token, revealMessage.LockSecretHash(), models.TransferStatusCanNotCancel, fmt.Sprintf("RevealSecret sending target=%s", utils.APex2(event.Receiver)), nil)
		//eh.photon.dao.UpdateTransferStatus(event.Token, revealMessage.LockSecretHash(), models.TransferStatusCanNotCancel, fmt.Sprintf("RevealSecret 正在发送 target=%s", utils.APex2(event.Receiver)))
		//eh.photon.NotifyTransferStatusChange(event.Token, revealMessage.LockSecretHash(), models.TransferStatusCanNotCancel, fmt.Sprintf("RevealSecret 正在发送 target=%s", utils.APex2(event.Receiver)))
		eh.photon.NotifyHandler.NotifySentTransferDetail(std)
	}
	return err
}
func (eh *stateMachineEventHandler) eventSendSecretRequest(event *mediatedtransfer.EventSendSecretRequest, stateManager *transfer.StateManager) (err error) {
	secretRequest := encoding.NewSecretRequest(event.LockSecretHash, event.Amount)
	err = secretRequest.Sign(eh.photon.PrivateKey, secretRequest)
	eh.photon.conditionQuit("EventSendSecretRequestBefore")
	ch := eh.photon.getChannelWithAddr(event.ChannelIdentifier)
	if ch == nil {
		panic("should not found")
	}
	if stateManager.LastReceivedMessage == nil {
		log.Warn(fmt.Sprintf("EventSendSecretRequest %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		eh.photon.UpdateChannelAndSaveAck(ch, stateManager.LastReceivedMessage.Tag())
		stateManager.LastReceivedMessage = nil
	}
	err = eh.photon.sendAsync(event.Receiver, secretRequest)
	return
}
func (eh *stateMachineEventHandler) eventSendMediatedTransfer(event *mediatedtransfer.EventSendMediatedTransfer, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.photon.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	if ch == nil {
		err = fmt.Errorf("receive eventSendMediatedTransfer,but cannot found the channel,there must be error, event=%s,stateManager=%s",
			utils.StringInterface(event, 3), utils.StringInterface(stateManager, 5),
		)
		log.Error(err.Error())
		return
	}
	//log.Trace(fmt.Sprintf("eventSendMediatedTransfer g=%s", utils.StringInterface(g, 3)))
	//log.Trace(fmt.Sprintf("eventSendMediatedTransfer ch=%s", utils.StringInterface(ch, 2)))
	mtr, err := ch.CreateMediatedTransfer(event.Initiator, event.Target, event.Fee, event.Amount, event.Expiration, event.LockSecretHash, event.Path)
	if err != nil {
		return
	}
	//log.Trace(fmt.Sprintf("mtr=%s", utils.StringInterface(mtr, 5)))
	err = mtr.Sign(eh.photon.PrivateKey, mtr)
	err = ch.RegisterTransfer(eh.photon.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	eh.photon.conditionQuit("EventSendMediatedTransferBefore")
	if stateManager.LastReceivedMessage == nil {
		if stateManager.Name != initiator.NameInitiatorTransition {
			log.Warn(fmt.Sprintf("EventSendMediatedTransfer %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		}
		err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		var fromCh *channel.Channel
		fromCh, err = eh.photon.findChannelByIdentifier(event.FromChannel)
		if err != nil {
			return
		}
		t, _ := stateManager.LastReceivedMessage.Tag().(*transfer.MessageTag)
		echohash := t.EchoHash
		ack := eh.photon.Protocol.CreateAck(echohash)
		tx := eh.photon.dao.StartTx()
		eh.photon.dao.SaveAck(echohash, ack.Pack(), tx)
		err = eh.photon.UpdateChannel(channel.NewChannelSerialization(ch), tx)
		if err != nil {
			//数据库保存错误,不可能发生,一旦发生了,程序只能向上层报告错误.
			// database cache fault, impossible to happen.
			// If occurs, then throw this to upper layer.
			//err=tx.Rollback()
			panic(fmt.Sprintf("update channel err %s", err))
		}
		err = eh.photon.UpdateChannel(channel.NewChannelSerialization(fromCh), tx)
		if err != nil {
			//数据库保存错误,不可能发生,一旦发生了,程序只能向上层报告错误.
			// database cache fault, impossible to happen.
			// If occurs, then throw this to upper layer.
			//err=tx.Rollback()
			panic(fmt.Sprintf("update channel err %s", err))
		}
		err = tx.Commit()
		stateManager.LastReceivedMessage = nil
	}
	err = eh.photon.sendAsync(receiver, mtr)
	if err == nil {
		std := eh.photon.dao.UpdateSentTransferDetailStatus(ch.TokenAddress, mtr.LockSecretHash, models.TransferStatusCanCancel, fmt.Sprintf("MediatedTransfer sending target=%s", utils.APex2(receiver)), nil)
		//eh.photon.NotifyTransferStatusChange(ch.TokenAddress, mtr.LockSecretHash, models.TransferStatusCanCancel, fmt.Sprintf("MediatedTransfer 正在发送 target=%s", utils.APex2(receiver)))
		eh.photon.NotifyHandler.NotifySentTransferDetail(std)
	}
	return
}
func (eh *stateMachineEventHandler) eventSendUnlock(event *mediatedtransfer.EventSendBalanceProof) (err error) {
	if !eh.photon.IsChainEffective {
		// 在无有效公链的情况下,阻止本应该发送的unlock消息并保存在数据库,当切换到有效公链的时候,再视情况决定是否发送
		eh.photon.dao.NewUnlockToSend(event.LockSecretHash, event.Token, event.Receiver, eh.photon.GetBlockNumber())
		log.Info(fmt.Sprintf("unlock message [lockSecertHash=%s token=%s receiver=%s] saved in db and wait to send after effective chain",
			event.LockSecretHash.String(), event.Token.String(), event.Receiver.String()))
		return
	}
	receiver := event.Receiver
	g := eh.photon.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	if ch == nil {
		err = fmt.Errorf("receive EventSendBalanceProof,but cannot found the channel,there must be error, event=%s",
			utils.StringInterface(event, 3))
		log.Error(err.Error())
		return
	}
	tr, err := ch.CreateUnlock(event.LockSecretHash)
	if err != nil {
		return
	}
	err = tr.Sign(eh.photon.PrivateKey, tr)
	err = ch.RegisterTransfer(eh.photon.GetBlockNumber(), tr)
	if err != nil {
		return
	}
	eh.photon.conditionQuit("EventSendUnlockBefore")
	err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	err = eh.photon.sendAsync(receiver, tr)
	if err == nil {
		eh.photon.dao.UpdateSentTransferDetailStatusMessage(event.Token, event.LockSecretHash, fmt.Sprintf("Unlock sending target=%s", utils.APex2(receiver)))
	}
	// 清空Token2LockSecretHash2Channels
	eh.photon.removeToken2LockSecretHash2channel(event.LockSecretHash, ch)
	return
}
func (eh *stateMachineEventHandler) eventSendAnnouncedDisposed(event *mediatedtransfer.EventSendAnnounceDisposed, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.photon.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	if ch == nil {
		err = fmt.Errorf("receive eventSendAnnouncedDisposed,but cannot found the channel,there must be error, event=%s,stateManager=%s",
			utils.StringInterface(event, 3), utils.StringInterface(stateManager, 5),
		)
		log.Error(err.Error())
		return
	}
	mtr, err := ch.CreateAnnouceDisposed(event.LockSecretHash, eh.photon.GetBlockNumber(), event.Reason)
	if err != nil {
		return
	}
	err = mtr.Sign(eh.photon.PrivateKey, mtr)
	err = ch.RegisterAnnouceDisposed(mtr)
	if err != nil {
		return
	}
	err = eh.photon.dao.MarkLockSecretHashDisposed(event.LockSecretHash, ch.ChannelIdentifier.ChannelIdentifier)
	if err != nil {
		return
	}
	if stateManager.LastReceivedMessage == nil {
		log.Warn(fmt.Sprintf("EventSendAnnounceDisposed %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		eh.photon.UpdateChannelAndSaveAck(ch, stateManager.LastReceivedMessage.Tag())
		//有可能同一个消息会引发两个 event send, 比如收到 中间节点EventAnnouceDisposed
		// 会触发EventSendAnnounceDisposedResponse 和EventSendMediatedTransfer
		// Maybe a message triggers two event send, when receiving Medaited node's EventAnnounceDisposed
		// which triggers EventSendAnnounceDisposedResponse, and EventSendMediatedTransfer
		stateManager.LastReceivedMessage = nil
	}
	eh.photon.conditionQuit("EventSendAnnouncedDisposedBefore")
	err = eh.photon.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventSendAnnouncedDisposedResponse(event *mediatedtransfer.EventSendAnnounceDisposedResponse, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.photon.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	if ch == nil {
		return fmt.Errorf("GetPartenerAddress2Channel returns nil ,but %s should have channel with %s on token %s",
			utils.APex2(g.OurAddress), utils.APex2(receiver), utils.APex2(g.TokenAddress))
	}
	mtr, err := ch.CreateAnnounceDisposedResponse(event.LockSecretHash, eh.photon.GetBlockNumber())
	if err != nil {
		return
	}
	err = mtr.Sign(eh.photon.PrivateKey, mtr)
	err = ch.RegisterAnnounceDisposedResponse(mtr, eh.photon.GetBlockNumber())
	if err != nil {
		return
	}
	eh.photon.conditionQuit("EventSendAnnouncedDisposedResponseBefore")
	if stateManager.LastReceivedMessage == nil {
		log.Warn(fmt.Sprintf("EventSendAnnounceDisposedResponse %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		eh.photon.UpdateChannelAndSaveAck(ch, stateManager.LastReceivedMessage.Tag())
		stateManager.LastReceivedMessage = nil
	}
	err = eh.photon.sendAsync(receiver, mtr)
	// 清空Token2LockSecretHash2Channels
	eh.photon.removeToken2LockSecretHash2channel(event.LockSecretHash, ch)
	return
}
func (eh *stateMachineEventHandler) eventContractSendRegisterSecret(event *mediatedtransfer.EventContractSendRegisterSecret) (err error) {
	b, err := eh.photon.Chain.SecretRegistryProxy.IsSecretRegistered(event.Secret)
	if err != nil {
		return err
	}
	if b {
		log.Info(fmt.Sprintf("Secret %s already registered", utils.HPex(event.Secret)))
		return
	}
	result := eh.photon.Chain.SecretRegistryProxy.RegisterSecretAsync(event.Secret)
	go func() {
		var err error
		err = <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("register secret on chain err %s,secret=%s you may lose your token because of this error",
				err, event.Secret.String()))
		}
	}()
	return nil
}
func (eh *stateMachineEventHandler) eventWithdrawFailed(e2 *mediatedtransfer.EventWithdrawFailed, manager *transfer.StateManager) (err error) {
	//wait from RemoveExpiredHashlockTransfer from partner.
	//need do nothing ,just wait.
	return nil
}
func (eh *stateMachineEventHandler) eventContractSendUnlock(e2 *mediatedtransfer.EventContractSendUnlock, manager *transfer.StateManager) (err error) {
	ch, err := eh.photon.findChannelByIdentifier(e2.ChannelIdentifier)
	if err != nil {
		log.Error(fmt.Sprintf("EventContractSendUnlock,but cannot find channel %s, eh may happen long later restart after a stop,event=%s",
			e2.ChannelIdentifier, utils.StringInterface(e2, 3),
		))
		return
	}
	var p *channeltype.UnlockProof
	//找到此次链上密码注册对应的锁,肯定应该找到.
	unlockProofs := ch.PartnerState.GetCanUnlockOnChainLocks()
	for _, u := range unlockProofs {
		if u.Lock.LockSecretHash == e2.LockSecretHash {
			p = u
			break
		}
	}
	if p == nil {
		log.Error(fmt.Sprintf("EventContractSendUnlock,but cannot found lock in channel,e2=%s,manager=%s,channel=%s",
			utils.StringInterface(e2, 2), utils.StringInterface(manager, 3), ch,
		))
		return
	}
	result := ch.ExternState.Unlock([]*channeltype.UnlockProof{p}, ch.PartnerState.BalanceProofState.ContractTransferAmount)
	go func() {
		err := <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("contract unlock  on %s failed, channel is gone, error:%s", ch.ChannelIdentifier.String(), err))
		}
	}()
	return nil
}

/*
the transfer I payed for a payee has expired. give a new balanceproof which doesn't contain this hashlock
*/
func (eh *stateMachineEventHandler) eventUnlockFailed(e2 *mediatedtransfer.EventUnlockFailed, manager *transfer.StateManager) (err error) {
	if manager.Name == target.NameTargetTransition {
		panic("event unlock failed can not  happen for a target node")
	}
	ch, err := eh.photon.findChannelByIdentifier(e2.ChannelIdentifier)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", e2.ChannelIdentifier))
		return
	}
	log.Info(fmt.Sprintf("remove expired hashlock channel=%s,hashlock=%s ", utils.HPex(e2.ChannelIdentifier), utils.HPex(e2.LockSecretHash)))
	/*
		unlock 失败,谨慎起见, 只有在对方不知道密码的情况下,才可能成功移除锁.
	*/
	tr, err := ch.CreateRemoveExpiredHashLockTransfer(e2.LockSecretHash, eh.photon.GetBlockNumber())
	if err != nil {
		log.Warn(fmt.Sprintf("Get Event UnlockFailed ,but hashlock cannot be removed err:%s", err))
		return
	}
	err = tr.Sign(eh.photon.PrivateKey, tr)
	err = ch.RegisterRemoveExpiredHashlockTransfer(tr, eh.photon.GetBlockNumber())
	if err != nil {
		log.Error(fmt.Sprintf("register mine RegisterRemoveExpiredHashlockTransfer err %s", err))
		return
	}
	eh.photon.conditionQuit("EventRemoveExpiredHashlockTransferBefore")
	err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	err = eh.photon.sendAsync(ch.PartnerState.Address, tr)
	std := eh.photon.dao.UpdateSentTransferDetailStatus(ch.TokenAddress, e2.LockSecretHash, models.TransferStatusFailed, fmt.Sprintf("transfer timeout err=%s", e2.Reason), nil)
	//eh.photon.NotifyTransferStatusChange(ch.TokenAddress, e2.LockSecretHash, models.TransferStatusFailed, fmt.Sprintf("交易超时失败 err=%s", e2.Reason))
	eh.photon.NotifyHandler.NotifySentTransferDetail(std)
	// 清空Token2LockSecretHash2Channels
	eh.photon.removeToken2LockSecretHash2channel(e2.LockSecretHash, ch)
	return
}

func (eh *stateMachineEventHandler) eventSaveFeeChargeRecord(e *mediatedtransfer.EventSaveFeeChargeRecord) (err error) {
	r := &models.FeeChargeRecord{
		LockSecretHash: e.LockSecretHash,
		TokenAddress:   e.TokenAddress,
		TransferFrom:   e.TransferFrom,
		TransferTo:     e.TransferTo,
		TransferAmount: e.TransferAmount,
		InChannel:      e.InChannel,
		OutChannel:     e.OutChannel,
		Fee:            e.Fee,
		Timestamp:      e.Timestamp,
		Data:           e.Data,
		BlockNumber:    e.BlockNumber,
	}
	return eh.photon.dao.SaveFeeChargeRecord(r)
}

func (eh *stateMachineEventHandler) OnEvent(event transfer.Event, stateManager *transfer.StateManager) (err error) {
	var ch *channel.Channel
	switch e2 := event.(type) {
	case *mediatedtransfer.EventSendMediatedTransfer:
		err = eh.eventSendMediatedTransfer(e2, stateManager)
		eh.photon.conditionQuit("EventSendMediatedTransferAfter")
	case *mediatedtransfer.EventSendRevealSecret:
		err = eh.eventSendRevealSecret(e2, stateManager)
		eh.photon.conditionQuit("EventSendRevealSecretAfter")
	case *mediatedtransfer.EventSendBalanceProof:
		//unlock and update remotely (send the LockSecretHash message)
		err = eh.eventSendUnlock(e2)
		eh.photon.conditionQuit("EventSendUnlockAfter")
	case *mediatedtransfer.EventSendSecretRequest:
		err = eh.eventSendSecretRequest(e2, stateManager)
		eh.photon.conditionQuit("EventSendSecretRequestAfter")
	case *mediatedtransfer.EventSendAnnounceDisposed:
		err = eh.eventSendAnnouncedDisposed(e2, stateManager)
		eh.photon.conditionQuit("EventSendAnnouncedDisposedAfter")
	case *mediatedtransfer.EventSendAnnounceDisposedResponse:
		err = eh.eventSendAnnouncedDisposedResponse(e2, stateManager)
		eh.photon.conditionQuit("EventSendAnnouncedDisposedResponseAfter")
	case *transfer.EventTransferSentSuccess:
		ch, err = eh.photon.findChannelByIdentifier(e2.ChannelIdentifier)
		if err != nil {
			err = fmt.Errorf("receive EventTransferSentSuccess,but channel not exist %s", utils.HPex(e2.ChannelIdentifier))
			return
		}
		err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		if err != nil {
			log.Error(fmt.Sprintf("UpdateChannelNoTx err %s", err))
		}
		//st := eh.photon.dao.NewSentTransfer(eh.photon.GetBlockNumber(), e2.ChannelIdentifier, ch.ChannelIdentifier.OpenBlockNumber, ch.TokenAddress, e2.Target, ch.GetNextNonce(), e2.Amount, e2.LockSecretHash, e2.Data)
		//eh.photon.NotifyHandler.NotifySentTransfer(st)
		eh.finishOneTransfer(event)
	case *transfer.EventTransferSentFailed:
		std := eh.photon.dao.UpdateSentTransferDetailStatus(e2.Token, e2.LockSecretHash, models.TransferStatusFailed, fmt.Sprintf("transfer fail err=%s", e2.Reason), nil)
		//eh.photon.NotifyTransferStatusChange(e2.Token, e2.LockSecretHash, models.TransferStatusFailed, fmt.Sprintf("交易失败 err=%s", e2.Reason))
		eh.photon.NotifyHandler.NotifySentTransferDetail(std)
		eh.finishOneTransfer(event)
	case *transfer.EventTransferReceivedSuccess:
		ch, err = eh.photon.findChannelByIdentifier(e2.ChannelIdentifier)
		if err != nil {
			err = fmt.Errorf("receive EventTransferReceivedSuccess,but channel not exist %s", utils.HPex(e2.ChannelIdentifier))
			return
		}
		err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		if err != nil {
			log.Error(fmt.Sprintf("UpdateChannelNoTx err %s", err))
		}
		rt := eh.photon.dao.NewReceivedTransfer(eh.photon.GetBlockNumber(), e2.ChannelIdentifier, ch.ChannelIdentifier.OpenBlockNumber, ch.TokenAddress, e2.Initiator, ch.PartnerState.BalanceProofState.Nonce, e2.Amount, e2.LockSecretHash, e2.Data)
		eh.photon.NotifyHandler.NotifyReceiveTransfer(rt)
	case *mediatedtransfer.EventUnlockSuccess:
	case *mediatedtransfer.EventWithdrawFailed:
		log.Error(fmt.Sprintf("EventWithdrawFailed hashlock=%s,reason=%s", utils.HPex(e2.LockSecretHash), e2.Reason))
		err = eh.eventWithdrawFailed(e2, stateManager)
	case *mediatedtransfer.EventWithdrawSuccess:
		/*
					  The withdraw is currently handled by the netting channel, once the close
			     event is detected all locks will be withdrawn
		*/
	case *mediatedtransfer.EventContractSendUnlock:
		//do nothing for five events above
		err = eh.eventContractSendUnlock(e2, stateManager)
	case *mediatedtransfer.EventUnlockFailed:
		log.Error(fmt.Sprintf("unlockfailed hashlock=%s,reason=%s", utils.HPex(e2.LockSecretHash), e2.Reason))
		err = eh.eventUnlockFailed(e2, stateManager)
		eh.photon.conditionQuit("EventSendRemoveExpiredHashlockTransferAfter")
	case *mediatedtransfer.EventContractSendRegisterSecret:
		err = eh.eventContractSendRegisterSecret(e2)
	case *mediatedtransfer.EventRemoveStateManager:
		delete(eh.photon.Transfer2StateManager, e2.Key)
	case *mediatedtransfer.EventSaveFeeChargeRecord:
		err = eh.eventSaveFeeChargeRecord(e2)
	default:
		err = fmt.Errorf("unkown event :%s", utils.StringInterface1(event))
		log.Error(err.Error())
	}
	return
}

//remove the successful transfer's state manager
func (eh *stateMachineEventHandler) finishOneTransfer(ev transfer.Event) {
	var err error
	var lockSecretHash common.Hash
	var tokenAddress common.Address
	switch e2 := ev.(type) {
	case *transfer.EventTransferSentSuccess:
		log.Info(fmt.Sprintf("EventTransferSentSuccess for LockSecretHash %s ", e2.LockSecretHash.String()))
		lockSecretHash = e2.LockSecretHash
		tokenAddress = e2.Token
		err = nil
	case *transfer.EventTransferSentFailed:
		log.Warn(fmt.Sprintf("EventTransferSentFailed for LockSecretHash %s,because of %s", e2.LockSecretHash.String(), e2.Reason))
		lockSecretHash = e2.LockSecretHash
		err = errors.New(e2.Reason)
		tokenAddress = e2.Token
	default:
		panic("unknow event")
	}
	if lockSecretHash != utils.EmptyHash {
		smkey := utils.Sha3(lockSecretHash[:], tokenAddress[:])
		r := eh.photon.Transfer2Result[smkey]
		if r == nil { //restart after crash?
			log.Error(fmt.Sprintf("transfer finished ,but have no relate results :%s", utils.StringInterface(ev, 2)))
			return
		}
		r.Result <- err
		delete(eh.photon.Transfer2Result, smkey)
	}
}

//1. 必须能够正确处理重复的ContractTokenAddedStateChange事件
func (eh *stateMachineEventHandler) HandleTokenAdded(st *mediatedtransfer.ContractTokenAddedStateChange) error {
	tokenAddress := st.TokenAddress
	if eh.photon.Token2ChannelGraph[tokenAddress] != nil {
		log.Warn(fmt.Sprintf("receive duplicate ContractTokenAddedStateChange=%s",
			utils.StringInterface(st, 3),
		))
		return nil
	}
	log.Info(fmt.Sprintf("NewTokenAdd token=%s", tokenAddress.String()))
	err := eh.photon.dao.AddToken(st.TokenAddress, utils.EmptyAddress)
	if err != nil {
		return err
	}
	g := graph.NewChannelGraph(eh.photon.NodeAddress, st.TokenAddress, nil)
	eh.photon.Token2TokenNetwork[tokenAddress] = utils.EmptyAddress
	eh.photon.Token2ChannelGraph[tokenAddress] = g
	return nil
}

//1. 必须能够正确处理重复的newchannel 事件.
func (eh *stateMachineEventHandler) handleChannelNew(st *mediatedtransfer.ContractNewChannelStateChange) error {
	// 忽略SettleTimeout小于限定最小值的通道
	minSettleTimeout := eh.photon.getMinSettleTimeout()
	if st.SettleTimeout <= minSettleTimeout {
		log.Warn(fmt.Sprintf("ignore new channel %s because SettleTimeout < %d", st.ChannelIdentifier.String(), minSettleTimeout))
		return nil
	}
	participant1 := st.Participant1
	participant2 := st.Participant2
	tokenAddress := st.TokenAddress
	log.Info(fmt.Sprintf("NewChannel token=%s,participant1=%s,participant2=%s",
		utils.APex2(tokenAddress),
		utils.APex2(participant1),
		utils.APex2(participant2),
	))
	g := eh.photon.getToken2ChannelGraph(tokenAddress)
	g.AddPath(participant1, participant2)
	err := eh.photon.dao.NewNonParticipantChannel(tokenAddress, st.ChannelIdentifier.ChannelIdentifier, participant1, participant2)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	isParticipant := eh.photon.NodeAddress == participant2 || eh.photon.NodeAddress == participant1
	partner := st.Participant1
	if partner == eh.photon.NodeAddress {
		partner = st.Participant2
	}
	if isParticipant {
		c := g.GetPartenerAddress2Channel(partner)
		if c != nil {
			log.Warn(fmt.Sprintf("receive duplicate ContractNewChannelStateChange, c=%s,\n,statechange=%s",
				utils.StringInterface(c, 5), utils.StringInterface(st, 3),
			))
			return nil
		}
		eh.photon.registerChannel(tokenAddress, partner, st.ChannelIdentifier, st.SettleTimeout)
		other := participant2
		if other == eh.photon.NodeAddress {
			other = participant1
		}
		eh.photon.startHealthCheckFor(other)
	} else {
		log.Trace("ignoring new channel, this node is not a participant.")
	}
	return nil
}

//1. 重复的ContractBalanceStateChange没有什么大的影响
func (eh *stateMachineEventHandler) handleBalance(st *mediatedtransfer.ContractBalanceStateChange) error {
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier)
	if err != nil {
		//log.Trace(fmt.Sprintf("ContractBalanceStateChange i'm not a participant,channelIdentifier=%s", utils.HPex(st.ChannelIdentifier)))
		return nil
	}
	if st.GetBlockNumber() < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractBalanceStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = eh.photon.UpdateChannelContractBalance(channel.NewChannelSerialization(ch))
	return err
}

//1. 必须能够正确处理重复的ContractClosedStateChange
func (eh *stateMachineEventHandler) handleClosed(st *mediatedtransfer.ContractClosedStateChange) error {
	channelIdentifier := st.ChannelIdentifier
	ch, err := eh.photon.findChannelByIdentifier(channelIdentifier)
	if err != nil {
		//i'm not a participant
		// 如果不是自己参与的channel,移除路由中的path
		token, p1, p2, err2 := eh.photon.dao.GetNonParticipantChannelByID(st.ChannelIdentifier)
		if err2 != nil {
			log.Warn(fmt.Sprintf("receive ContractClosedStateChange=%s,but channel not found ",
				utils.StringInterface(st, 3),
			))
			return nil
		}
		g := eh.photon.getToken2ChannelGraph(token)
		if g != nil {
			if p1 != utils.EmptyAddress && p2 != utils.EmptyAddress {
				g.RemovePath(p1, p2)
			}
		}
		err = eh.photon.dao.RemoveNonParticipantChannel(st.ChannelIdentifier)
		return err
	}
	if ch.State == channeltype.StateClosed {
		log.Warn(fmt.Sprintf("receive duplicate ContractClosedStateChange=%s,channel already closed",
			utils.StringInterface(st, 3),
		))
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = eh.photon.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}

/*
从内存中将此 channel 所有相关信息都移除
1. channel graph 中的channel 信息
2. 数据库中的 channel 信息
3. 数据库中 non participant 信息
4. statemanager 中有关该 channel 的信息, 会自行移除
*/
/*
 *	Remove this channel infor from local memory.
 *
 *	1. channel information in channel graph
 *	2. channel info in database.
 *	3. non participant info in database
 *	4. channel reference by statemanager
 */
func (eh *stateMachineEventHandler) removeSettledChannel(ch *channel.Channel) error {
	g := eh.photon.getChannelGraph(ch.ChannelIdentifier.ChannelIdentifier)
	g.RemoveChannel(ch)
	cs := channel.NewChannelSerialization(ch)
	err := eh.photon.dao.RemoveChannel(cs)
	if err != nil {
		return err
	}
	err = eh.photon.dao.NewSettledChannel(cs)
	if err != nil {
		return err
	}
	err = eh.photon.dao.RemoveNonParticipantChannel(ch.ChannelIdentifier.ChannelIdentifier)
	/*
		通知上层
	*/
	eh.photon.NotifyHandler.NotifyChannelStatus(channeltype.ChannelSerialization2ChannelDataDetail(cs))
	return err
}
func (eh *stateMachineEventHandler) handleSettled(st *mediatedtransfer.ContractSettledStateChange) error {
	log.Trace(fmt.Sprintf("%s settled event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	// 如果用户在100块合作settle通道,101块又一次打开,然后photon在100+确认块之前崩溃,那么重启后会重复收到该这2次事件,如果这里不验证settle的块号,会导致通道数据被清空
	// 所以这里忽略掉小于OpenBlockNumber的合作Settle事件
	if st.SettledBlock < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractSettledStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
		return err
	}
	return eh.removeSettledChannel(ch)
}

//大部分与 settle 相同,是否可以合并呢?或者合约上干脆合并了?
// Most part of this is same as settle
// can we just combine them?
//1. 必须能够正确处理重复的事件
func (eh *stateMachineEventHandler) handleCooperativeSettled(st *mediatedtransfer.ContractCooperativeSettledStateChange) error {
	log.Trace(fmt.Sprintf("%s cooperative settled event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier)
	if err != nil {
		//i'm not a participant
		// 如果不是自己参与的channel,移除路由中的path
		token, p1, p2, err2 := eh.photon.dao.GetNonParticipantChannelByID(st.ChannelIdentifier)
		if err2 != nil {
			log.Warn(fmt.Sprintf("receive ContractCooperativeSettledStateChange,but channel not found %s",
				utils.StringInterface(st, 3),
			))
			return nil
		}
		g := eh.photon.getToken2ChannelGraph(token)
		if g != nil {
			if p1 != utils.EmptyAddress && p2 != utils.EmptyAddress {
				g.RemovePath(p1, p2)
			}
		}
		return eh.photon.dao.RemoveNonParticipantChannel(st.ChannelIdentifier)
	}
	// 如果用户在100块合作settle通道,101块又一次打开,然后photon在100+确认块之前崩溃,那么重启后会重复收到该这2次事件,如果这里不验证settle的块号,会导致通道数据被清空
	// 所以这里忽略掉小于OpenBlockNumber的合作Settle事件
	if st.SettledBlock < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractCooperativeSettledStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.removeSettledChannel(ch)
	//if true {
	//	g := eh.photon.getChannelGraph(ch.ChannelIdentifier.ChannelIdentifier)
	//	log.Trace(fmt.Sprintf("after settle g=%s", utils.StringInterface(g, 3)))
	//}
	// 通知该通道下所有存在pending lock的state manager,可以放心的announce disposed或者尝试新路由了
	// notify all statemanager with pending locks, then we can send announcedisposed and try another route.
	eh.dispatchByPendingLocksInChannel(ch, st)
	return err
}

//1. 必须能够处理重复的ContractChannelWithdrawStateChange
func (eh *stateMachineEventHandler) handleWithdraw(st *mediatedtransfer.ContractChannelWithdrawStateChange) error {
	log.Trace(fmt.Sprintf("%s withdraw event handle", utils.HPex(st.ChannelIdentifier.ChannelIdentifier)))
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier.ChannelIdentifier)
	if err != nil {
		return nil
	}
	// 考虑到极小状况下会在崩溃重启后收到重复的上一个通道发生的事件,如果这里不验证块号,可能出现上一个channel的withdraw事件在新channel上被处理的BUG,导致新channel失败
	if st.BlockNumber < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractChannelWithdrawStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.photon.UpdateChannelState(channel.NewChannelSerialization(ch))
	// 通知该通道下所有存在pending lock的state manager,可以放心的announce disposed或者尝试新路由了
	// nofity all statemanager with pending locks, and send announce disposed or try new route.
	eh.dispatchByPendingLocksInChannel(ch, st)
	return err
}

//如果是对方 unlock 我的锁,那么有可能需要 punish 对方,即使不需要 punish 对方,settle 的时候也需要用到新的 locksroot 和 transferamount
/*
 *	handleUnlockOnChain : function to handle unlock event.
 *
 *	Note that if case is that my partner unlocks locks of mine, then maybe I need to punish my partner,
 *	even if that's not the case, when channel settle, I still need to use the the locksroot and transferAmount on chain
1. 重复的unlock没什么影响,只要保证后续的事件按序抵达即可
*/
func (eh *stateMachineEventHandler) handleUnlockOnChain(st *mediatedtransfer.ContractUnlockStateChange) error {
	log.Trace(fmt.Sprintf("%s unlock event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	// 考虑到极小状况下会在崩溃重启后收到重复的上一个通道发生的事件,如果这里不验证块号,影响不大,但验证一下肯定没有问题
	if st.BlockNumber < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractUnlockStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handle unlock ChannelStateTransition err=%s", err))
		return err
	}
	//对方解锁我发出去的交易,考虑可否惩罚
	// my partner unlock transfer I sent, consider punish him?
	if eh.photon.NodeAddress == st.Participant {
		ad := eh.photon.dao.GetReceivedAnnounceDisposed(st.LockHash, ch.ChannelIdentifier.ChannelIdentifier)
		if ad != nil {
			result := ch.ExternState.PunishObsoleteUnlock(common.BytesToHash(ad.LockHash), ad.AdditionalHash, ad.Signature)
			go func() {
				var err2 error
				err2 = <-result.Result
				if err2 != nil {
					log.Error(fmt.Sprintf("PunishObsoleteUnlock %s ,err2 %s", utils.BPex(ad.LockHash), err2))
				}
			}()
		}
	}
	err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	return err
}

//必须能够处理重复的punish事件,因为重复的punish只是更新通道状态,所以重复也没什么影响
func (eh *stateMachineEventHandler) handlePunishedOnChain(st *mediatedtransfer.ContractPunishedStateChange) error {
	log.Trace(fmt.Sprintf("%s punished event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier)
	if err != nil {
		log.Warn(fmt.Sprintf("receive ContractPunishedStateChange,but cannot found channel %s",
			utils.StringInterface(st, 3),
		))
		return nil
	}
	// 考虑到极小状况下会在崩溃重启后收到重复的上一个通道发生的事件,如果这里不验证块号,影响不大,但验证一下肯定没有问题
	if st.BlockNumber < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractPunishedStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handle punish ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	return err
}

//1. 必须正确处理重复的ContractBalanceProofUpdatedStateChange,这里只是更新相关参与方的状态,所以重复的事件并不影响
func (eh *stateMachineEventHandler) handleBalanceProofOnChain(st *mediatedtransfer.ContractBalanceProofUpdatedStateChange) error {
	log.Trace(fmt.Sprintf("%s balance proof update event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.photon.findChannelByIdentifier(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	// 考虑到极小状况下会在崩溃重启后收到重复的上一个通道发生的事件,如果这里不验证块号,影响不大,但验证一下肯定没有问题
	if st.BlockNumber < ch.ChannelIdentifier.OpenBlockNumber {
		log.Error("got repeat ContractBalanceProofUpdatedStateChange , ignore ")
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	err = eh.photon.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	return err
}

//1. 必须能够正确处理重复的ContractSecretRevealOnChainStateChange,
//todo 这里有一个潜在的问题A给B发交易,A收到ContractSecretRevealOnChainStateChange,然后会给B发unlock消息,
// 这时候A崩溃,等A立即重启以后,会再次处理ContractSecretRevealOnChainStateChange,从而导致unlock消息发送两次.
// 但是两次unlock消息nonce不同,从而导致通道不可用.
func (eh *stateMachineEventHandler) handleSecretRegisteredOnChain(st *mediatedtransfer.ContractSecretRevealOnChainStateChange) error {
	// 这里需要注册密码,否则unlock消息无法正常发送
	// we need register secret here, otherwise we can not send unlock.
	eh.photon.registerRevealedLockSecretHash(st.LockSecretHash, st.Secret, st.BlockNumber)
	//需要 disatch 给相关的 statemanager, 让他们处理未完成的交易.
	// we need dispatch it to relevant statemanager, and let them handle incomplete transfers.
	eh.dispatchBySecretHash(st.LockSecretHash, st)
	return nil
}

func (eh *stateMachineEventHandler) handleBlockStateChange(st *transfer.BlockStateChange) error {
	eh.dispatchToAllTasks(st)
	//for _, cg := range eh.photon.Token2ChannelGraph {
	//	for _, c := range cg.ChannelIdentifier2Channel {
	//		err := eh.ChannelStateTransition(c, st)
	//		if err != nil {
	//			log.Error(fmt.Sprintf("ChannelStateTransition err %s", err))
	//		}
	//	}
	//}
	return nil
}

/*
	处理有效公链/无效公链状态切换的相关逻辑
*/
func (eh *stateMachineEventHandler) handleEffectiveChainStateChange(st *transfer.EffectiveChainStateChange) (err error) {
	isChainEffective := st.IsEffective
	if isChainEffective == eh.photon.IsChainEffective {
		// 过滤重复
		return
	}
	eh.photon.IsChainEffective = isChainEffective
	eh.photon.EffectiveChangeTimestamp = st.LastBlockNumberTimestamp
	if !isChainEffective {
		// 有效公链切无效公链
		log.Info("photon works without effective chain now...")
		// 1. 启动无有效公链状态下的用户提醒线程
		go eh.startNoEffectiveChainNotifyLoop()
		// 2. 通知上层进入无网
		select {
		case eh.photon.EthConnectionStatus <- netshare.Disconnected:
		default:
			//never block
		}
	} else {
		// 无效公链切有效公链,包含启动时
		log.Info("photon works with effective chain now...")
		// 0. 通知上层进入有网
		select {
		case eh.photon.EthConnectionStatus <- netshare.Connected:
		default:
			//never block
		}
		// 1. 上传手续费设置给PFS
		if fm, ok := eh.photon.FeePolicy.(*FeeModule); ok {
			err2 := fm.SubmitFeePolicyToPFS()
			if err2 != nil {
				log.Error(fmt.Sprintf("set fee policy to pfs err =%s", err2.Error()))
			}
		}
		// 2. 刷新所有通道状态信息到pfs及pms
		for _, cg := range eh.photon.Token2ChannelGraph {
			for _, ch := range cg.ChannelIdentifier2Channel {
				if ch.DelegateState != channeltype.ChannelDelegateStateSuccess {
					// 不管状态,所有尚未settle的通道都需要委托到pms
					eh.photon.submitDelegateToPms(ch)
				}
				if ch.State == channeltype.StateOpened {
					// 仅提交open状态的通道到pfs
					eh.photon.submitBalanceProofToPfs(ch)
				}
			}
		}
		// 3. 获取所有等待发送的Unlock消息并发送,并从db里移除
		unlockToSendList := eh.photon.dao.GetAllUnlockToSend()
		for _, unlockToSend := range unlockToSendList {
			e := &mediatedtransfer.EventSendBalanceProof{
				LockSecretHash: common.BytesToHash(unlockToSend.LockSecretHash),
				Token:          common.BytesToAddress(unlockToSend.TokenAddress),
				Receiver:       common.BytesToAddress(unlockToSend.ReceiverAddress),
			}
			err2 := eh.OnEvent(e, nil)
			if err2 != nil {
				log.Error(fmt.Sprintf("unlock message [lockSecertHash=%s token=%s receiver=%s] saved in db send err : %s",
					e.LockSecretHash.String(), e.Token.String(), e.Receiver.String(), err2.Error()))
			}
			key := utils.Sha3(unlockToSend.LockSecretHash, unlockToSend.TokenAddress, unlockToSend.ReceiverAddress).Bytes()
			eh.photon.dao.RemoveUnlockToSend(key)
		}
		// 4. 关闭提醒线程
		eh.stopNoEffectiveChainNotifyLoop()
	}
	// 下发到所有的stateManager里面,正在进行的交易自行进行对应处理
	eh.dispatchToAllTasks(st)
	return nil
}

func (eh *stateMachineEventHandler) startNoEffectiveChainNotifyLoop() {
	if eh.noEffectiveChainNotifyLoopQuitChan == nil {
		eh.noEffectiveChainNotifyLoopQuitChan = make(chan *struct{})
	}
	periodBlock := eh.photon.getMinSettleTimeout() / 10
	periodSecond := time.Duration(periodBlock) * time.Second // 这里应该取MinSettleTimeout/10 * 出块间隔
	for {
		select {
		case <-eh.noEffectiveChainNotifyLoopQuitChan:
			return
		case <-time.After(periodSecond):
			t := time.Since(time.Unix(eh.photon.EffectiveChangeTimestamp, 0)).Round(time.Second)
			warning := fmt.Sprintf("photon has been worked without effective block chain for about %s", t)
			eh.photon.NotifyHandler.NotifyString(notify.LevelWarn, warning)
			log.Warn(warning)
		}
	}
}

func (eh *stateMachineEventHandler) stopNoEffectiveChainNotifyLoop() {
	if eh.noEffectiveChainNotifyLoopQuitChan != nil {
		close(eh.noEffectiveChainNotifyLoopQuitChan)
		eh.noEffectiveChainNotifyLoopQuitChan = nil
	}
}

//avoid dead lock
func (eh *stateMachineEventHandler) ChannelStateTransition(c *channel.Channel, st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *transfer.BlockStateChange:
		if c.State == channeltype.StateClosed {
			settlementEnd := c.ExternState.ClosedBlock + int64(c.SettleTimeout) + params.PunishBlockNumber
			if st2.BlockNumber > settlementEnd {
				//wait for user call settle
			}
		}
	case *mediatedtransfer.ContractClosedStateChange:
		if c.State != channeltype.StateClosed {
			c.State = channeltype.StateClosed
			c.ExternState.SetClosed(st2.ClosedBlock)
			c.ExternState.SetSettled(st2.ClosedBlock + int64(c.SettleTimeout) + params.PunishBlockNumber)
			c.HandleClosed(st2.ClosingAddress, st2.TransferredAmount, st2.LocksRoot)
		} else {
			log.Warn(fmt.Sprintf("channel closed on a different block or close event happened twice channel=%s,closedblock=%d,thisblock=%d",
				c.ChannelIdentifier.String(), c.ExternState.ClosedBlock, st2.ClosedBlock))
		}
	case *mediatedtransfer.ContractSettledStateChange:
		//settled channel should be removed.
		c.State = channeltype.StateSettled
		if c.ExternState.SetSettled(st2.SettledBlock) {
			c.HandleSettled(st2.SettledBlock)
		} else {
			log.Warn(fmt.Sprintf("channel is already settled on a different block channelIdentifier=%s,settleblock=%d,thisblock=%d",
				c.ChannelIdentifier.String(), c.ExternState.SettledBlock, st2.SettledBlock))
		}
	case *mediatedtransfer.ContractCooperativeSettledStateChange:
		//settled channel should be removed.
		c.State = channeltype.StateSettled
		if c.ExternState.SetSettled(st2.SettledBlock) {
			c.HandleSettled(st2.SettledBlock)
		} else {
			log.Warn(fmt.Sprintf("channel is already settled on a different block channelIdentifier=%s,settleblock=%d,thisblock=%d",
				c.ChannelIdentifier.String(), c.ExternState.SettledBlock, st2.SettledBlock))
		}
	case *mediatedtransfer.ContractChannelWithdrawStateChange:
		if c.ChannelIdentifier.OpenBlockNumber < st2.BlockNumber {
			c.HandleWithdrawed(st2.BlockNumber, st2.Participant1, st2.Participant2, st2.Participant1Balance, st2.Participant2Balance)
		} else {
			log.Warn(fmt.Sprintf("receive withdraw event,but channel's openblocknumber=%d,new openblocknumber=%d",
				c.ChannelIdentifier.OpenBlockNumber, st2.BlockNumber,
			))
		}
	case *mediatedtransfer.ContractBalanceStateChange:
		participant := st2.ParticipantAddress
		balance := st2.Balance
		var channelState *channel.EndState
		channelState, err = c.GetStateFor(participant)
		if err != nil {
			return
		}
		if channelState.ContractBalance.Cmp(balance) != 0 {
			err = channelState.UpdateContractBalance(balance)
		}
	case *mediatedtransfer.ContractUnlockStateChange:
		var channelState *channel.EndState
		channelState, err = c.GetStateFor(st2.Participant)
		if err != nil {
			return
		}
		if c.State == channeltype.StateOpened {
			panic("must closed")
		}
		log.Trace(fmt.Sprintf("channel %s unlocked %s", c.ChannelIdentifier.String(), st2.TransferAmount))
		channelState.SetContractTransferAmount(st2.TransferAmount)
	case *mediatedtransfer.ContractPunishedStateChange:
		c.HandleChannelPunished(st2.Beneficiary)
	case *mediatedtransfer.ContractBalanceProofUpdatedStateChange:
		c.HandleBalanceProofUpdated(st2.Participant, st2.TransferAmount, st2.LocksRoot)
	}
	return

}

func (eh *stateMachineEventHandler) OnBlockchainStateChange(st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *mediatedtransfer.ContractTokenAddedStateChange:
		err = eh.HandleTokenAdded(st2)
	case *mediatedtransfer.ContractNewChannelStateChange:
		eh.photon.conditionQuit("EventNewChannelFromChainBeforeDeal")
		err = eh.handleChannelNew(st2)
		eh.photon.conditionQuit("EventNewChannelFromChainAfterDeal")
	case *mediatedtransfer.ContractBalanceStateChange:
		eh.photon.conditionQuit("EventDepositFromChainBeforeDeal")
		err = eh.handleBalance(st2)
		eh.photon.conditionQuit("EventDepositFromChainAfterDeal")
	case *mediatedtransfer.ContractClosedStateChange:
		eh.photon.conditionQuit("EventChannelCloseFromChainBeforeDeal")
		err = eh.handleClosed(st2)
		eh.photon.conditionQuit("EventChannelCloseFromChainAfterDeal")
	case *mediatedtransfer.ContractSettledStateChange:
		eh.photon.conditionQuit("EventChannelSettleFromChainBeforeDeal")
		err = eh.handleSettled(st2)
		eh.photon.conditionQuit("EventChannelSettleFromChainAfterDeal")
	case *mediatedtransfer.ContractSecretRevealOnChainStateChange:
		err = eh.handleSecretRegisteredOnChain(st2)
	case *mediatedtransfer.ContractUnlockStateChange:
		eh.photon.conditionQuit("EventUnlockFromChainBeforeDeal")
		err = eh.handleUnlockOnChain(st2)
		eh.photon.conditionQuit("EventUnlockFromChainAfterDeal")
	case *mediatedtransfer.ContractPunishedStateChange:
		eh.photon.conditionQuit("EventPunishFromChainBeforeDeal")
		err = eh.handlePunishedOnChain(st2)
		eh.photon.conditionQuit("EventPunishFromChainAfterDeal")
	case *mediatedtransfer.ContractBalanceProofUpdatedStateChange:
		eh.photon.conditionQuit("EventUpdateBalanceProofFromChainBeforeDeal")
		err = eh.handleBalanceProofOnChain(st2)
		eh.photon.conditionQuit("EventUpdateBalanceProofFromChainAfterDeal")
	case *mediatedtransfer.ContractCooperativeSettledStateChange:
		eh.photon.conditionQuit("EventCooperativeSettleFromChainBeforeDeal")
		err = eh.handleCooperativeSettled(st2)
		eh.photon.conditionQuit("EventCooperativeSettleFromChainAfterDeal")
	case *mediatedtransfer.ContractChannelWithdrawStateChange:
		eh.photon.conditionQuit("EventWithdrawFromChainBeforeDeal")
		err = eh.handleWithdraw(st2)
		eh.photon.conditionQuit("EventWithdrawFromChainAfterDeal")
	case *transfer.BlockStateChange:
		err = eh.handleBlockStateChange(st2)
	case *transfer.EffectiveChainStateChange:
		err = eh.handleEffectiveChainStateChange(st2)
	default:
		err = fmt.Errorf("OnBlockchainStateChange unknown statechange :%s", utils.StringInterface1(st))
		log.Error(err.Error())
	}
	return
}

//recive a message and before processed
func (eh *stateMachineEventHandler) updateStateManagerFromStateChange(mgr *transfer.StateManager, stateChange transfer.StateChange) {
	var msg encoding.SignedMessager
	var quitName string
	switch st2 := stateChange.(type) {
	case *mediatedtransfer.ActionInitTargetStateChange:
		quitName = "ActionInitTargetStateChange"
		msg = st2.Message
	case *mediatedtransfer.ReceiveSecretRequestStateChange:
		quitName = "ReceiveSecretRequestStateChange"
		msg = st2.Message
	case *mediatedtransfer.ReceiveAnnounceDisposedStateChange:
		quitName = "ReceiveAnnounceDisposedStateChange"
		msg = st2.Message
	case *mediatedtransfer.ReceiveUnlockStateChange:
		quitName = "ReceiveUnlockStateChange"
	case *mediatedtransfer.ActionInitMediatorStateChange:
		quitName = "ActionInitMediatorStateChange"
		msg = st2.Message
	case *mediatedtransfer.MediatorReReceiveStateChange:
		quitName = "MediatorReReceiveStateChange"
		msg = st2.Message
	case *mediatedtransfer.ActionInitInitiatorStateChange:
		quitName = "ActionInitInitiatorStateChange"
		//new transfer trigger from user
	case *mediatedtransfer.ReceiveSecretRevealStateChange:
		quitName = "ReceiveSecretRevealStateChange"
	}
	if msg != nil {
		mgr.LastReceivedMessage = msg
	}
	if len(quitName) > 0 {
		eh.photon.conditionQuit(quitName)
	}
}
