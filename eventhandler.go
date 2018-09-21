package smartraiden

import (
	"fmt"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network/graph"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

//run inside loop of raiden service
type stateMachineEventHandler struct {
	raiden *RaidenService
}

func newStateMachineEventHandler(raiden *RaidenService) *stateMachineEventHandler {
	h := &stateMachineEventHandler{
		raiden: raiden,
	}
	return h
}

/*
dispatch it to all state managers and log generated events
*/
func (eh *stateMachineEventHandler) dispatchToAllTasks(st transfer.StateChange) {
	for _, mgrs := range eh.raiden.Transfer2StateManager {
		eh.dispatch(mgrs, st)
	}
}

/*
dispatch it to the state manager corresponding to `lockSecretHash`
*/
func (eh *stateMachineEventHandler) dispatchBySecretHash(lockSecretHash common.Hash, st transfer.StateChange) {
	for _, mgr := range eh.raiden.Transfer2StateManager {
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
	eh.raiden.conditionQuit("EventSendRevealSecretBefore")
	eh.raiden.registerSecret(event.Secret)
	revealMessage := encoding.NewRevealSecret(event.Secret)
	err = revealMessage.Sign(eh.raiden.PrivateKey, revealMessage)
	err = eh.raiden.sendAsync(event.Receiver, revealMessage) //单独处理 reaveal secret
	if err == nil {
		eh.raiden.db.UpdateTransferStatus(revealMessage.LockSecretHash(), models.TransferStatusCanNotCancel, fmt.Sprintf("RevealSecret 正在发送 target=%s", utils.APex2(event.Receiver)))
	}
	return err
}
func (eh *stateMachineEventHandler) eventSendSecretRequest(event *mediatedtransfer.EventSendSecretRequest, stateManager *transfer.StateManager) (err error) {
	secretRequest := encoding.NewSecretRequest(event.LockSecretHash, event.Amount)
	err = secretRequest.Sign(eh.raiden.PrivateKey, secretRequest)
	eh.raiden.conditionQuit("EventSendSecretRequestBefore")
	ch := eh.raiden.getChannelWithAddr(event.ChannelIdentifier)
	if ch == nil {
		panic("should not found")
	}
	if stateManager.LastReceivedMessage == nil {
		log.Warn(fmt.Sprintf("EventSendSecretRequest %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		eh.raiden.updateChannelAndSaveAck(ch, stateManager.LastReceivedMessage.Tag())
		stateManager.LastReceivedMessage = nil
	}
	err = eh.raiden.sendAsync(event.Receiver, secretRequest)
	return
}
func (eh *stateMachineEventHandler) eventSendMediatedTransfer(event *mediatedtransfer.EventSendMediatedTransfer, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateMediatedTransfer(event.Initiator, event.Target, event.Fee, event.Amount, event.Expiration, event.LockSecretHash)
	if err != nil {
		return
	}
	err = mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(eh.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	eh.raiden.conditionQuit("EventSendMediatedTransferBefore")
	if stateManager.LastReceivedMessage == nil {
		if stateManager.Name != initiator.NameInitiatorTransition {
			log.Warn(fmt.Sprintf("EventSendMediatedTransfer %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		}
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		var fromCh *channel.Channel
		fromCh, err = eh.raiden.findChannelByAddress(event.FromChannel)
		if err != nil {
			return
		}
		t, _ := stateManager.LastReceivedMessage.Tag().(*transfer.MessageTag)
		echohash := t.EchoHash
		ack := eh.raiden.Protocol.CreateAck(echohash)
		tx := eh.raiden.db.StartTx()
		eh.raiden.db.SaveAck(echohash, ack.Pack(), tx)
		err = eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
		err = eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(fromCh), tx)
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
	err = eh.raiden.sendAsync(receiver, mtr)
	if err == nil {
		eh.raiden.db.UpdateTransferStatus(mtr.LockSecretHash, models.TransferStatusCanCancel, fmt.Sprintf("MediatedTransfer 正在发送 target=%s", utils.APex2(receiver)))
	}
	return
}
func (eh *stateMachineEventHandler) eventSendUnlock(event *mediatedtransfer.EventSendBalanceProof, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	tr, err := ch.CreateUnlock(event.LockSecretHash)
	if err != nil {
		return
	}
	err = tr.Sign(eh.raiden.PrivateKey, tr)
	err = ch.RegisterTransfer(eh.raiden.GetBlockNumber(), tr)
	if err != nil {
		return
	}
	eh.raiden.conditionQuit("EventSendUnlockBefore")
	err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	err = eh.raiden.sendAsync(receiver, tr)
	if err == nil {
		eh.raiden.db.UpdateTransferStatusMessage(event.LockSecretHash, fmt.Sprintf("Unlock 正在发送 target=%s", utils.APex2(receiver)))
	}
	return
}
func (eh *stateMachineEventHandler) eventSendAnnouncedDisposed(event *mediatedtransfer.EventSendAnnounceDisposed, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateAnnouceDisposed(event.LockSecretHash, eh.raiden.GetBlockNumber())
	if err != nil {
		return
	}
	err = mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterAnnouceDisposed(mtr)
	if err != nil {
		return
	}
	err = eh.raiden.db.MarkLockSecretHashDisposed(event.LockSecretHash, ch.ChannelIdentifier.ChannelIdentifier)
	if err != nil {
		return
	}
	if stateManager.LastReceivedMessage == nil {
		log.Warn(fmt.Sprintf("EventSendAnnounceDisposed %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		eh.raiden.updateChannelAndSaveAck(ch, stateManager.LastReceivedMessage.Tag())
		//有可能同一个消息会引发两个 event send, 比如收到 中间节点EventAnnouceDisposed
		// 会触发EventSendAnnounceDisposedResponse 和EventSendMediatedTransfer
		// Maybe a message triggers two event send, when receiving Medaited node's EventAnnounceDisposed
		// which triggers EventSendAnnounceDisposedResponse, and EventSendMediatedTransfer
		stateManager.LastReceivedMessage = nil
	}
	eh.raiden.conditionQuit("EventSendAnnouncedDisposedBefore")
	err = eh.raiden.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventSendAnnouncedDisposedResponse(event *mediatedtransfer.EventSendAnnounceDisposedResponse, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	g := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := g.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateAnnounceDisposedResponse(event.LockSecretHash, eh.raiden.GetBlockNumber())
	if err != nil {
		return
	}
	err = mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterAnnounceDisposedResponse(mtr, eh.raiden.GetBlockNumber())
	if err != nil {
		return
	}
	eh.raiden.conditionQuit("EventSendAnnouncedDisposedResponseBefore")
	if stateManager.LastReceivedMessage == nil {
		log.Warn(fmt.Sprintf("EventSendAnnounceDisposedResponse %s,but has no lastReceviedMessage", utils.StringInterface(event, 3)))
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	} else {
		eh.raiden.updateChannelAndSaveAck(ch, stateManager.LastReceivedMessage.Tag())
		stateManager.LastReceivedMessage = nil
	}
	err = eh.raiden.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventContractSendRegisterSecret(event *mediatedtransfer.EventContractSendRegisterSecret) (err error) {
	b, err := eh.raiden.Chain.SecretRegistryProxy.IsSecretRegistered(event.Secret)
	if err != nil {
		return err
	}
	if b {
		log.Info(fmt.Sprintf("Secret %s already registered", utils.HPex(event.Secret)))
		return
	}
	result := eh.raiden.Chain.SecretRegistryProxy.RegisterSecretAsync(event.Secret)
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
func (eh *stateMachineEventHandler) eventContractSendWithdraw(e2 *mediatedtransfer.EventContractSendWithdraw, manager *transfer.StateManager) (err error) {
	//if manager.Name != target.NameTargetTransition && manager.Name != mediator.NameMediatorTransition {
	//	panic("EventWithdrawFailed can only comes from a target node or mediated node")
	//}
	//ch, err := eh.raiden.findChannelByAddress(e2.ChannelIdentifier)
	//if err != nil {
	//	log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", e2.ChannelIdentifier))
	//	return
	//}
	//unlockProofs := ch.PartnerState.GetKnownUnlocks()
	//result := ch.ExternState.Unlock(unlockProofs, ch.PartnerState.BalanceProofState.ContractTransferAmount)
	//go func() {
	//	err := <-result.Result
	//	if err != nil {
	//		log.Error(fmt.Sprintf("withdraw on %s failed, channel is gone, error:%s", ch.ChannelIdentifier.String(), err))
	//	}
	//}()
	return nil
}

/*
the transfer I payed for a payee has expired. give a new balanceproof which doesn't contain this hashlock
*/
func (eh *stateMachineEventHandler) eventUnlockFailed(e2 *mediatedtransfer.EventUnlockFailed, manager *transfer.StateManager) (err error) {
	if manager.Name == target.NameTargetTransition {
		panic("event unlock failed can not  happen for a target node")
	}
	ch, err := eh.raiden.findChannelByAddress(e2.ChannelIdentifier)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", e2.ChannelIdentifier))
		return
	}
	log.Info(fmt.Sprintf("remove expired hashlock channel=%s,hashlock=%s ", utils.HPex(e2.ChannelIdentifier), utils.HPex(e2.LockSecretHash)))
	tr, err := ch.CreateRemoveExpiredHashLockTransfer(e2.LockSecretHash, eh.raiden.GetBlockNumber())
	if err != nil {
		log.Warn(fmt.Sprintf("Get Event UnlockFailed ,but hashlock cannot be removed err:%s", err))
		return
	}
	err = tr.Sign(eh.raiden.PrivateKey, tr)
	err = ch.RegisterRemoveExpiredHashlockTransfer(tr, eh.raiden.GetBlockNumber())
	if err != nil {
		log.Error(fmt.Sprintf("register mine RegisterRemoveExpiredHashlockTransfer err %s", err))
		return
	}
	eh.raiden.conditionQuit("EventRemoveExpiredHashlockTransferBefore")
	err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
	err = eh.raiden.sendAsync(ch.PartnerState.Address, tr)
	return
}
func (eh *stateMachineEventHandler) OnEvent(event transfer.Event, stateManager *transfer.StateManager) (err error) {
	var ch *channel.Channel
	switch e2 := event.(type) {
	case *mediatedtransfer.EventSendMediatedTransfer:
		err = eh.eventSendMediatedTransfer(e2, stateManager)
		eh.raiden.conditionQuit("EventSendMediatedTransferAfter")
	case *mediatedtransfer.EventSendRevealSecret:
		err = eh.eventSendRevealSecret(e2, stateManager)
		eh.raiden.conditionQuit("EventSendRevealSecretAfter")
	case *mediatedtransfer.EventSendBalanceProof:
		//unlock and update remotely (send the LockSecretHash message)
		err = eh.eventSendUnlock(e2, stateManager)
		eh.raiden.conditionQuit("EventSendUnlockAfter")
	case *mediatedtransfer.EventSendSecretRequest:
		err = eh.eventSendSecretRequest(e2, stateManager)
		eh.raiden.conditionQuit("EventSendSecretRequestAfter")
	case *mediatedtransfer.EventSendAnnounceDisposed:
		err = eh.eventSendAnnouncedDisposed(e2, stateManager)
		eh.raiden.conditionQuit("EventSendAnnouncedDisposedAfter")
	case *mediatedtransfer.EventSendAnnounceDisposedResponse:
		err = eh.eventSendAnnouncedDisposedResponse(e2, stateManager)
		eh.raiden.conditionQuit("EventSendAnnouncedDisposedResponseAfter")
	case *transfer.EventTransferSentSuccess:
		ch, err = eh.raiden.findChannelByAddress(e2.ChannelIdentifier)
		if err != nil {
			err = fmt.Errorf("receive EventTransferSentSuccess,but channel not exist %s", utils.HPex(e2.ChannelIdentifier))
			return
		}
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		if err != nil {
			log.Error(fmt.Sprintf("UpdateChannelNoTx err %s", err))
		}
		eh.raiden.db.NewSentTransfer(eh.raiden.GetBlockNumber(), e2.ChannelIdentifier, ch.TokenAddress, e2.Target, ch.GetNextNonce(), e2.Amount)
		eh.finishOneTransfer(event)
	case *transfer.EventTransferSentFailed:
		eh.finishOneTransfer(event)
	case *transfer.EventTransferReceivedSuccess:
		ch, err = eh.raiden.findChannelByAddress(e2.ChannelIdentifier)
		if err != nil {
			err = fmt.Errorf("receive EventTransferReceivedSuccess,but channel not exist %s", utils.HPex(e2.ChannelIdentifier))
			return
		}
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		if err != nil {
			log.Error(fmt.Sprintf("UpdateChannelNoTx err %s", err))
		}
		eh.raiden.db.NewReceivedTransfer(eh.raiden.GetBlockNumber(), e2.ChannelIdentifier, ch.TokenAddress, e2.Initiator, ch.PartnerState.BalanceProofState.Nonce, e2.Amount)
	case *mediatedtransfer.EventUnlockSuccess:
	case *mediatedtransfer.EventWithdrawFailed:
		log.Error(fmt.Sprintf("EventWithdrawFailed hashlock=%s,reason=%s", utils.HPex(e2.LockSecretHash), e2.Reason))
		err = eh.eventWithdrawFailed(e2, stateManager)
	case *mediatedtransfer.EventWithdrawSuccess:
		/*
					  The withdraw is currently handled by the netting channel, once the close
			     event is detected all locks will be withdrawn
		*/
	case *mediatedtransfer.EventContractSendWithdraw:
		//do nothing for five events above
		err = eh.eventContractSendWithdraw(e2, stateManager)
	case *mediatedtransfer.EventUnlockFailed:
		log.Error(fmt.Sprintf("unlockfailed hashlock=%s,reason=%s", utils.HPex(e2.LockSecretHash), e2.Reason))
		err = eh.eventUnlockFailed(e2, stateManager)
		eh.raiden.conditionQuit("EventSendRemoveExpiredHashlockTransferAfter")
	case *mediatedtransfer.EventContractSendRegisterSecret:
		err = eh.eventContractSendRegisterSecret(e2)
	case *mediatedtransfer.EventRemoveStateManager:
		delete(eh.raiden.Transfer2StateManager, e2.Key)
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
		log.Info(fmt.Sprintf("EventTransferSentSuccess for id %d ", e2.LockSecretHash))
		lockSecretHash = e2.LockSecretHash
		tokenAddress = e2.Token
		err = nil
	case *transfer.EventTransferSentFailed:
		log.Warn(fmt.Sprintf("EventTransferSentFailed for id %d,because of %s", e2.LockSecretHash, e2.Reason))
		lockSecretHash = e2.LockSecretHash
		err = errors.New(e2.Reason)
		tokenAddress = e2.Token
	default:
		panic("unknow event")
	}
	if lockSecretHash != utils.EmptyHash {
		smkey := utils.Sha3(lockSecretHash[:], tokenAddress[:])
		r := eh.raiden.Transfer2Result[smkey]
		if r == nil { //restart after crash?
			log.Error(fmt.Sprintf("transfer finished ,but have no relate results :%s", utils.StringInterface(ev, 2)))
			return
		}
		r.Result <- err
		delete(eh.raiden.Transfer2Result, smkey)
	}
}
func (eh *stateMachineEventHandler) HandleTokenAdded(st *mediatedtransfer.ContractTokenAddedStateChange) error {
	if st.RegistryAddress != eh.raiden.RegistryAddress {
		panic("unkown registry")
	}
	tokenAddress := st.TokenAddress
	tokenNetworkAddress := st.TokenNetworkAddress
	log.Info(fmt.Sprintf("NewTokenAdd token=%s,tokennetwork=%s", tokenAddress.String(), tokenNetworkAddress.String()))
	err := eh.raiden.db.AddToken(st.TokenAddress, st.TokenNetworkAddress)
	if err != nil {
		return err
	}
	g := graph.NewChannelGraph(eh.raiden.NodeAddress, st.TokenAddress, nil)
	eh.raiden.TokenNetwork2Token[tokenNetworkAddress] = tokenAddress
	eh.raiden.Token2TokenNetwork[tokenAddress] = tokenNetworkAddress
	eh.raiden.Token2ChannelGraph[tokenAddress] = g
	return nil
}
func (eh *stateMachineEventHandler) handleChannelNew(st *mediatedtransfer.ContractNewChannelStateChange) error {
	tokenNetworkAddress := st.TokenNetworkAddress
	participant1 := st.Participant1
	participant2 := st.Participant2
	tokenAddress := eh.raiden.TokenNetwork2Token[tokenNetworkAddress]
	log.Info(fmt.Sprintf("NewChannel tokenNetwork=%s,token=%s,participant1=%s,participant2=%s",
		utils.APex2(tokenNetworkAddress),
		utils.APex2(tokenAddress),
		utils.APex2(participant1),
		utils.APex2(participant2),
	))
	g := eh.raiden.getToken2ChannelGraph(tokenAddress)
	g.AddPath(participant1, participant2)
	err := eh.raiden.db.NewNonParticipantChannel(tokenAddress, st.ChannelIdentifier.ChannelIdentifier, participant1, participant2)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	isParticipant := eh.raiden.NodeAddress == participant2 || eh.raiden.NodeAddress == participant1
	partner := st.Participant1
	if partner == eh.raiden.NodeAddress {
		partner = st.Participant2
	}
	if isParticipant {
		eh.raiden.registerChannel(tokenNetworkAddress, partner, st.ChannelIdentifier, st.SettleTimeout)
		other := participant2
		if other == eh.raiden.NodeAddress {
			other = participant1
		}
		eh.raiden.startHealthCheckFor(other)
	} else {
		log.Trace("ignoring new channel, this node is not a participant.")
	}
	return nil
}

func (eh *stateMachineEventHandler) handleBalance(st *mediatedtransfer.ContractBalanceStateChange) error {
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
		//todo 处理这个事件,路由的时候可以考虑节点之间的权重,权重值=双方 deposit 之和
		// todo handle this event, when routing we should consider the weight between nodes, weight = sum of deposits between a participant pair.
		log.Trace(fmt.Sprintf("ContractBalanceStateChange i'm not a participant,channelAddress=%s", utils.HPex(st.ChannelIdentifier)))
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = eh.raiden.db.UpdateChannelContractBalance(channel.NewChannelSerialization(ch))
	return nil
}

func (eh *stateMachineEventHandler) handleClosed(st *mediatedtransfer.ContractClosedStateChange) error {
	channelAddress := st.ChannelIdentifier
	ch, err := eh.raiden.findChannelByAddress(channelAddress)
	if err != nil {
		//i'm not a participant
		token := eh.raiden.TokenNetwork2Token[st.TokenNetworkAddress]
		err = eh.raiden.db.RemoveNonParticipantChannel(token, st.ChannelIdentifier)
		return err
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}

/*
从内存中将此 channel 所有相关信息都移除
1. channel graph 中的channel 信息
2. 数据库中的 channel 信息
3. 数据库中 non participant 信息
4. todo statemanager 中有关该 channel 的信息, 是否有?
*/
/*
 *	Remove this channel infor from local memory.
 *
 *	1. channel information in channel graph
 *	2. channel info in database.
 *	3. non participant info in database
 *	4. todo is there any channel info in statemanager,?
 */
func (eh *stateMachineEventHandler) removeSettledChannel(ch *channel.Channel) error {
	g := eh.raiden.getChannelGraph(ch.ChannelIdentifier.ChannelIdentifier)
	g.RemoveChannel(ch)
	cs := channel.NewChannelSerialization(ch)
	err := eh.raiden.db.RemoveChannel(cs)
	if err != nil {
		return err
	}
	err = eh.raiden.db.NewSettledChannel(cs)
	if err != nil {
		return err
	}
	err = eh.raiden.db.RemoveNonParticipantChannel(ch.TokenAddress, ch.ChannelIdentifier.ChannelIdentifier)
	return err
}
func (eh *stateMachineEventHandler) handleSettled(st *mediatedtransfer.ContractSettledStateChange) error {
	log.Trace(fmt.Sprintf("%s settled event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
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
func (eh *stateMachineEventHandler) handleCooperativeSettled(st *mediatedtransfer.ContractCooperativeSettledStateChange) error {
	log.Trace(fmt.Sprintf("%s cooperative settled event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.removeSettledChannel(ch)
	// 通知该通道下所有存在pending lock的state manager,可以放心的announce disposed或者尝试新路由了
	// notify all statemanager with pending locks, then we can send announcedisposed and try another route.
	eh.dispatchByPendingLocksInChannel(ch, st)
	return err
}
func (eh *stateMachineEventHandler) handleWithdraw(st *mediatedtransfer.ContractChannelWithdrawStateChange) error {
	log.Trace(fmt.Sprintf("%s cooperative settled event handle", utils.HPex(st.ChannelIdentifier.ChannelIdentifier)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier.ChannelIdentifier)
	if err != nil {
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
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
 */
func (eh *stateMachineEventHandler) handleUnlockOnChain(st *mediatedtransfer.ContractUnlockStateChange) error {
	log.Trace(fmt.Sprintf("%s unlock event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handle unlock ChannelStateTransition err=%s", err))
		return err
	}
	//对方解锁我发出去的交易,考虑可否惩罚
	// my partner unlock transfer I sent, consider punish him?
	if eh.raiden.NodeAddress == st.Participant {
		ad := eh.raiden.db.GetReceiviedAnnounceDisposed(st.LockHash, ch.ChannelIdentifier.ChannelIdentifier)
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
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}
func (eh *stateMachineEventHandler) handlePunishedOnChain(st *mediatedtransfer.ContractPunishedStateChange) error {
	log.Trace(fmt.Sprintf("%s punished event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handle punish ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}
func (eh *stateMachineEventHandler) handleBalanceProofOnChain(st *mediatedtransfer.ContractBalanceProofUpdatedStateChange) error {
	log.Trace(fmt.Sprintf("%s balance proof update event handle", utils.HPex(st.ChannelIdentifier)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}
func (eh *stateMachineEventHandler) handleSecretRegisteredOnChain(st *mediatedtransfer.ContractSecretRevealOnChainStateChange) error {
	// 这里需要注册密码,否则unlock消息无法正常发送
	// we need register secret here, otherwise we can not send unlock.
	eh.raiden.registerRevealedLockSecretHash(st.LockSecretHash, st.Secret, st.BlockNumber)
	//需要 disatch 给相关的 statemanager, 让他们处理未完成的交易.
	// we need dispatch it to relevant statemanager, and let them handle incomplete transfers.
	eh.dispatchBySecretHash(st.LockSecretHash, st)
	return nil
}

//avoid dead lock
func (eh *stateMachineEventHandler) ChannelStateTransition(c *channel.Channel, st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *transfer.BlockStateChange:
		if c.State == channeltype.StateClosed {
			settlementEnd := c.ExternState.ClosedBlock + int64(c.SettleTimeout) //todo punish time
			if st2.BlockNumber > settlementEnd {
				//should not block todo fix it
				//err = c.ExternState.Settle()
			}
		}
	case *mediatedtransfer.ContractClosedStateChange:
		if c.State != channeltype.StateClosed {
			c.State = channeltype.StateClosed
			c.ExternState.SetClosed(st2.ClosedBlock)
			c.HandleClosed(st2.ClosingAddress, st2.TransferredAmount, st2.LocksRoot)
		} else {
			log.Warn(fmt.Sprintf("channel closed on a different block or close event happened twice channel=%s,closedblock=%d,thisblock=%d",
				c.ChannelIdentifier.String(), c.ExternState.ClosedBlock, st2.ClosedBlock))
		}
	case *mediatedtransfer.ContractSettledStateChange:
		//settled channel should be removed.
		if c.ExternState.SetSettled(st2.SettledBlock) {
			c.HandleSettled(st2.SettledBlock)
		} else {
			log.Warn(fmt.Sprintf("channel is already settled on a different block channeladdress=%s,settleblock=%d,thisblock=%d",
				c.ChannelIdentifier.String(), c.ExternState.SettledBlock, st2.SettledBlock))
		}
	case *mediatedtransfer.ContractCooperativeSettledStateChange:
		//settled channel should be removed.
		if c.ExternState.SetSettled(st2.SettledBlock) {
			c.HandleSettled(st2.SettledBlock)
		} else {
			log.Warn(fmt.Sprintf("channel is already settled on a different block channeladdress=%s,settleblock=%d,thisblock=%d",
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
	log.Trace(fmt.Sprintf("statechange received :%s", utils.StringInterface(st, 2)))
	switch st2 := st.(type) {
	case *mediatedtransfer.ContractTokenAddedStateChange:
		err = eh.HandleTokenAdded(st2)
	case *mediatedtransfer.ContractNewChannelStateChange:
		err = eh.handleChannelNew(st2)
	case *mediatedtransfer.ContractBalanceStateChange:
		err = eh.handleBalance(st2)
	case *mediatedtransfer.ContractClosedStateChange:
		err = eh.handleClosed(st2)
	case *mediatedtransfer.ContractSettledStateChange:
		err = eh.handleSettled(st2)
	case *mediatedtransfer.ContractSecretRevealOnChainStateChange:
		err = eh.handleSecretRegisteredOnChain(st2)
	case *mediatedtransfer.ContractUnlockStateChange:
		err = eh.handleUnlockOnChain(st2)
	case *mediatedtransfer.ContractPunishedStateChange:
		err = eh.handlePunishedOnChain(st2)
	case *mediatedtransfer.ContractBalanceProofUpdatedStateChange:
		err = eh.handleBalanceProofOnChain(st2)
	case *mediatedtransfer.ContractCooperativeSettledStateChange:
		err = eh.handleCooperativeSettled(st2)
	case *mediatedtransfer.ContractChannelWithdrawStateChange:
		err = eh.handleWithdraw(st2)
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
		eh.raiden.conditionQuit(quitName)
	}
}
