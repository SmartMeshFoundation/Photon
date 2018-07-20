package smartraiden

import (
	"fmt"

	"errors"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/graph"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

var errSentFailed = errors.New("sent failed")

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
Log a state change, dispatch it to all state managers and log generated events
*/
func (eh *stateMachineEventHandler) logAndDispatchToAllTasks(st transfer.StateChange) {
	for _, mgrs := range eh.raiden.Transfer2StateManager {
		eh.dispatch(mgrs, st)
	}
}

/*
Log a state change, dispatch it to the state manager corresponding to `idenfitier`
        and log generated events
*/
func (eh *stateMachineEventHandler) logAndDispatchBySecretHash(lockSecretHash common.Hash, st transfer.StateChange) {
	for _, mgr := range eh.raiden.Transfer2StateManager {
		//todo 这个未必是高效的方式,因为同时进行的 transfer 可能很多,会比较慢.
		if mgr.Identifier == lockSecretHash {
			eh.dispatch(mgr, st)
		}
	}
}

//Log a state change, dispatch it to the given state manager and log generated events
func (eh *stateMachineEventHandler) logAndDispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) []transfer.Event {
	events := eh.dispatch(stateManager, stateChange)
	return events
}
func (eh *stateMachineEventHandler) dispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) (events []transfer.Event) {
	eh.updateStateManagerFromReceivedMessageOrUserRequest(stateManager, stateChange)
	events = stateManager.Dispatch(stateChange)
	for _, e := range events {
		err := eh.OnEvent(e, stateManager)
		if err != nil {
			log.Error(fmt.Sprintf("stateMachineEventHandler dispatch:%v\n", err))
		}
	}
	return
}
func (eh *stateMachineEventHandler) eventSendMediatedTransfer(event *mediatedtransfer.EventSendMediatedTransfer, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateMediatedTransfer(event.Initiator, event.Target, event.Fee, event.Amount, event.Expiration, event.LockSecretHash)
	if err != nil {
		return
	}
	mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(eh.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	eh.updateStateManagerFromEvent(receiver, mtr, stateManager)
	eh.raiden.conditionQuit("EventSendMediatedTransferBefore")
	err = eh.raiden.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventSendUnlock(event *mediatedtransfer.EventSendBalanceProof, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	tr, err := ch.CreateUnlock(event.LockSecretHash)
	if err != nil {
		return
	}
	tr.Sign(eh.raiden.PrivateKey, tr)
	err = ch.RegisterTransfer(eh.raiden.GetBlockNumber(), tr)
	if err != nil {
		return
	}
	eh.updateStateManagerFromEvent(receiver, tr, stateManager)
	eh.raiden.conditionQuit("EventSendBalanceProofBefore")
	err = eh.raiden.sendAsync(receiver, tr)
	return
}
func (eh *stateMachineEventHandler) eventSendAnnouncedDisposed(event *mediatedtransfer.EventSendAnnounceDisposed, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateAnnouceDisposed(event.LockSecretHash, eh.raiden.GetBlockNumber())
	if err != nil {
		return
	}
	mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterAnnouceDisposed(mtr)
	if err != nil {
		return
	}
	err = eh.raiden.db.MarkLockSecretHashDisposed(event.LockSecretHash, ch.ChannelIdentifier.ChannelIdentifier)
	if err != nil {
		return
	}
	eh.updateStateManagerFromEvent(receiver, mtr, stateManager)
	eh.raiden.conditionQuit("EventSendAnnouncedDisposedBefore")
	err = eh.raiden.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventSendAnnouncedDisposedResponse(event *mediatedtransfer.EventSendAnnounceDisposedResponse, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateAnnounceDisposedResponse(event.LockSecretHash, eh.raiden.GetBlockNumber())
	if err != nil {
		return
	}
	mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterAnnounceDisposedResponse(mtr, eh.raiden.GetBlockNumber())
	if err != nil {
		return
	}
	eh.updateStateManagerFromEvent(receiver, mtr, stateManager)
	eh.raiden.conditionQuit("EventSendAnnouncedDisposedResponseBefore")
	err = eh.raiden.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventContractSendChannelClose(event *mediatedtransfer.EventContractSendChannelClose) (err error) {
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	if graph == nil {
		err = fmt.Errorf("EventContractSendChannelClose but token %s doesn't exist", utils.APex(event.Token))
		return
	}
	ch := graph.ChannelAddress2Channel[event.ChannelIdentifier]
	if ch == nil {
		err = fmt.Errorf("EventContractSendChannelClose  but channel %s doesn't exist,maybe have already settled", utils.HPex(event.ChannelIdentifier))
		return
	}
	balanceProof := ch.OurState.BalanceProofState
	ch.ExternState.Close(balanceProof)
	return
}
func (eh *stateMachineEventHandler) eventWithdrawFailed(e2 *mediatedtransfer.EventWithdrawFailed, manager *transfer.StateManager) (err error) {
	//wait from RemoveExpiredHashlockTransfer from partner.
	//need do nothing ,just wait.
	return nil
}
func (eh *stateMachineEventHandler) eventContractSendWithdraw(e2 *mediatedtransfer.EventContractSendWithdraw, manager *transfer.StateManager) (err error) {
	if manager.Name != target.NameTargetTransition && manager.Name != mediator.NameMediatorTransition {
		panic("EventWithdrawFailed can only comes from a target node or mediated node")
	}
	ch, err := eh.raiden.findChannelByAddress(e2.ChannelIdentifier)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", e2.ChannelIdentifier))
		return
	}
	unlockProofs := ch.PartnerState.GetKnownUnlocks()
	result := ch.ExternState.Unlock(unlockProofs, ch.PartnerState.BalanceProofState.ContractTransferAmount)
	go func() {
		err := <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("withdraw on %s failed, channel is gone, error:%s", ch.ChannelIdentifier, err))
		}
	}()
	return nil
}

/*
the transfer I payed for a payee has expired. give a new balanceproof which doesn't contain this hashlock
*/
func (eh *stateMachineEventHandler) eventUnlockFailed(e2 *mediatedtransfer.EventUnlockFailed, manager *transfer.StateManager) (err error) {
	if manager.Name != mediator.NameMediatorTransition && manager.Name != initiator.NameInitiatorTransition {
		panic("event unlock failed only happen for a mediated node")
	}
	ch, err := eh.raiden.findChannelByAddress(e2.ChannelIdentifier)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", e2.ChannelIdentifier))
		return
	}
	log.Info(fmt.Sprintf("remove expired hashlock channel=%s,hashlock=%s ", e2.ChannelIdentifier, utils.HPex(e2.LockSecretHash)))
	tr, err := ch.CreateRemoveExpiredHashLockTransfer(e2.LockSecretHash, eh.raiden.GetBlockNumber())
	if err != nil {
		log.Warn(fmt.Sprintf("Get Event UnlockFailed ,but hashlock cannot be removed err:%s", err))
		return
	}
	tr.Sign(eh.raiden.PrivateKey, tr)
	err = ch.RegisterRemoveExpiredHashlockTransfer(tr, eh.raiden.GetBlockNumber())
	if err != nil {
		log.Error(fmt.Sprintf("register mine RegisterRemoveExpiredHashlockTransfer err %s", err))
		return
	}
	/*
		save new channel status and sent RemoveExpiredHashlockTransfer must be atomic.
	*/
	eh.updateStateManagerFromEvent(ch.PartnerState.Address, tr, manager)
	eh.raiden.conditionQuit("EventRemoveExpiredHashlockTransferBefore")
	err = eh.raiden.sendAsync(ch.PartnerState.Address, tr)
	return
}
func (eh *stateMachineEventHandler) OnEvent(event transfer.Event, stateManager *transfer.StateManager) (err error) {
	switch e2 := event.(type) {
	case *mediatedtransfer.EventSendMediatedTransfer:
		err = eh.eventSendMediatedTransfer(e2, stateManager)
		eh.raiden.conditionQuit("EventSendMediatedTransferAfter")
	case *mediatedtransfer.EventSendRevealSecret:
		eh.raiden.conditionQuit("EventSendRevealSecretBefore")
		revealMessage := encoding.NewRevealSecret(e2.Secret)
		revealMessage.Sign(eh.raiden.PrivateKey, revealMessage)
		err = eh.raiden.sendAsync(e2.Receiver, revealMessage) //单独处理 reaveal secret
		eh.raiden.conditionQuit("EventSendRevealSecretAfter")
	case *mediatedtransfer.EventSendBalanceProof:
		//unlock and update remotely (send the LockSecretHash message)
		err = eh.eventSendUnlock(e2, stateManager)
		eh.raiden.conditionQuit("EventSendBalanceProofAfter")
	case *mediatedtransfer.EventSendSecretRequest:
		secretRequest := encoding.NewSecretRequest(e2.LockSecretHash, e2.Amount)
		secretRequest.Sign(eh.raiden.PrivateKey, secretRequest)
		eh.updateStateManagerFromEvent(e2.Receiver, secretRequest, stateManager)
		eh.raiden.conditionQuit("EventSendSecretRequestBefore")
		err = eh.raiden.sendAsync(e2.Receiver, secretRequest)
		eh.raiden.conditionQuit("EventSendSecretRequestAfter")
	case *mediatedtransfer.EventSendAnnounceDisposed:
		err = eh.eventSendAnnouncedDisposed(e2, stateManager)
		eh.raiden.conditionQuit("EventSendRefundTransferAfter")
	case *transfer.EventTransferSentSuccess:
		ch := eh.raiden.getChannelWithAddr(e2.ChannelIdentifier)
		if ch == nil {
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
		ch := eh.raiden.getChannelWithAddr(e2.ChannelIdentifier)
		if ch == nil {
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
		//TODO need payer's new signature to remove eh expired lock
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
		//should remove hashlock from channel todo fix bai
		log.Error(fmt.Sprintf("unlockfailed hashlock=%s,reason=%s", utils.HPex(e2.LockSecretHash), e2.Reason))
		err = eh.eventUnlockFailed(e2, stateManager)
		eh.raiden.conditionQuit("EventSendRemoveExpiredHashlockTransferAfter")
	case *mediatedtransfer.EventContractSendChannelClose:
		err = eh.eventContractSendChannelClose(e2)
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
	smkey := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	r := eh.raiden.Transfer2Result[smkey]
	if r == nil { //restart after crash?
		log.Error(fmt.Sprintf("you can ignore this error when this transfer is a direct transfer.\n transfer finished ,but have no relate results :%s", utils.StringInterface(ev, 2)))
		return
	}
	r.Result <- err
	delete(eh.raiden.Transfer2Result, smkey)
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
	graph := graph.NewChannelGraph(eh.raiden.NodeAddress, st.TokenAddress, nil, nil)
	eh.raiden.TokenNetwork2Token[tokenNetworkAddress] = tokenAddress
	eh.raiden.Token2TokenNetwork[tokenAddress] = tokenNetworkAddress
	eh.raiden.Token2ChannelGraph[tokenAddress] = graph
	eh.raiden.Tokens2ConnectionManager[tokenAddress] = NewConnectionManager(eh.raiden, tokenAddress)
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
	graph := eh.raiden.getToken2ChannelGraph(tokenAddress)
	graph.AddPath(participant1, participant2)
	eh.raiden.db.NewNonParticipantChannel(tokenAddress, st.ChannelIdentifier.ChannelIdentifier, participant1, participant2)
	connectionManager, err := eh.raiden.connectionManagerForToken(tokenAddress)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	isParticipant := eh.raiden.NodeAddress == participant2 || eh.raiden.NodeAddress == participant1
	isBootstrap := connectionManager.BootstrapAddr == participant1 || connectionManager.BootstrapAddr == participant2
	partner := st.Participant1
	if partner == eh.raiden.NodeAddress {
		partner = st.Participant2
	}
	if isParticipant {
		eh.raiden.registerNettingChannel(tokenNetworkAddress, partner)
		if !isBootstrap {
			other := participant2
			if other == eh.raiden.NodeAddress {
				other = participant1
			}
			eh.raiden.startHealthCheckFor(other)
		}
	} else if connectionManager.WantsMoreChannels() {
		go func() {
			defer rpanic.PanicRecover("RetryConnect")
			connectionManager.RetryConnect()
		}()
	} else {
		log.Trace("ignoring new channel, this node is not a participant.")
	}
	return nil
}

func (eh *stateMachineEventHandler) handleBalance(st *mediatedtransfer.ContractBalanceStateChange) error {
	tokenAddress := eh.raiden.TokenNetwork2Token[st.TokenNetworkAddress]
	participant := st.ParticipantAddress
	balance := st.Balance
	ch, err := eh.raiden.findChannelByAddress(st.ChannelIdentifier)
	if err != nil {
		//todo 处理这个事件,路由的时候可以考虑节点之间的权重,权重值=双方 deposit 之和
		log.Trace(fmt.Sprintf("ContractBalanceStateChange i'm not a participant,channelAddress=%s", utils.HPex(st.ChannelIdentifier)))
		return nil
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	eh.raiden.db.UpdateChannelContractBalance(channel.NewChannelSerialization(ch))
	if ch.ContractBalance().Cmp(utils.BigInt0) == 0 {
		connectionManager, _ := eh.raiden.connectionManagerForToken(tokenAddress)
		go func() {
			defer rpanic.PanicRecover(fmt.Sprintf("JoinChannel %s", utils.APex(participant)))
			connectionManager.JoinChannel(participant, balance)
		}()
	}
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
func (eh *stateMachineEventHandler) removeSettledChannel(ch *channel.Channel) error {
	graph := eh.raiden.getChannelGraph(ch.ChannelIdentifier.ChannelIdentifier)
	graph.RemoveChannel(ch)
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
	//todo remove channel st.channelAddress ,because eh channel is already settled
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

//如果是对方 unlock 我的锁,那么有可能需要 punish 对方,即使不需要 punish 对方,settle 的时候也需要用到新的 locksroot 和 transferamount
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
	if eh.raiden.NodeAddress == st.Participant {
		ad := eh.raiden.db.GetReceiviedAnnounceDisposed(st.LockHash, ch.ChannelIdentifier.ChannelIdentifier)
		if ad != nil {
			result := ch.ExternState.PunishObsoleteUnlock(common.BytesToHash(ad.LockHash), ad.AdditionalHash, ad.Signature)
			go func() {
				err := <-result.Result
				if err != nil {
					log.Error(fmt.Sprintf("PunishObsoleteUnlock %s ,err %s", utils.BPex(ad.LockHash), err))
				}
				//todo 要不要立即 settle?
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
func (eh *stateMachineEventHandler) handleSecretRegisteredOnChain(st *mediatedtransfer.ContractSecretRevealOnChainStateChange) error {
	eh.raiden.registerRevealedLockSecretHash(st.LockSecretHash, st.BlockNumber)
	//需要 disatch 给相关的 statemanager, 让他们处理未完成的交易.
	eh.logAndDispatchBySecretHash(st.LockSecretHash, st)
	return nil
}

//avoid dead lock
func (eh *stateMachineEventHandler) ChannelStateTransition(c *channel.Channel, st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *transfer.BlockStateChange:
		if c.State == channeltype.StateClosed {
			settlementEnd := c.ExternState.ClosedBlock + int64(c.SettleTimeout)
			if st2.BlockNumber > settlementEnd {
				//should not block todo fix it
				//err = c.ExternState.Settle()
			}
		}
	case *mediatedtransfer.ContractClosedStateChange:
		if st2.ChannelIdentifier == c.ChannelIdentifier.ChannelIdentifier {
			if c.State != channeltype.StateClosed {
				c.ExternState.SetClosed(st2.ClosedBlock)
				c.HandleClosed(st2.ClosedBlock, st2.ClosingAddress)
			} else {
				log.Warn(fmt.Sprintf("channel closed on a different block or close event happened twice channel=%s,closedblock=%d,thisblock=%d",
					c.ChannelIdentifier.String(), c.ExternState.ClosedBlock, st2.ClosedBlock))
			}
		}
	case *mediatedtransfer.ContractSettledStateChange:
		//settled channel should be removed. todo bai fix it
		if st2.ChannelIdentifier == c.ChannelIdentifier.ChannelIdentifier {
			if c.ExternState.SetSettled(st2.SettledBlock) {
				c.HandleSettled(st2.SettledBlock)
			} else {
				log.Warn(fmt.Sprintf("channel is already settled on a different block channeladdress=%s,settleblock=%d,thisblock=%d",
					c.ChannelIdentifier.String(), c.ExternState.SettledBlock, st2.SettledBlock))
			}
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
		channelState.SetContractTransferAmount(st2.TransferAmount)
	case *mediatedtransfer.ContractPunishedStateChange:
		var beneficiaryState, cheaterState *channel.EndState
		if st2.Beneficiary == c.OurState.Address {
			beneficiaryState = c.OurState
			cheaterState = c.PartnerState
		} else if st2.Beneficiary == c.PartnerState.Address {
			beneficiaryState = c.PartnerState
			cheaterState = c.OurState
		} else {
			panic(fmt.Sprintf("channel=%s,but participant =%s",
				st2.ChannelIdentifier.String(),
				st2.Beneficiary.String(),
			))
		}
		beneficiaryState.SetContractTransferAmount(utils.BigInt0)
		beneficiaryState.SetContractLocksroot(utils.EmptyHash)
		beneficiaryState.SetContractNonce(0xfffffff)
		beneficiaryState.ContractBalance = beneficiaryState.ContractBalance.Add(
			beneficiaryState.ContractBalance, cheaterState.ContractBalance,
		)
		cheaterState.ContractBalance = new(big.Int).Set(utils.BigInt0)
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
		//do nothing
	default:
		err = fmt.Errorf("OnBlockchainStateChange unknown statechange :%s", utils.StringInterface1(st))
		log.Error(err.Error())
	}
	return
}

//recive a message and before processed
func (eh *stateMachineEventHandler) updateStateManagerFromReceivedMessageOrUserRequest(mgr *transfer.StateManager, stateChange transfer.StateChange) {
	var msg encoding.Messager
	var quitName string
	switch st2 := stateChange.(type) {
	case *mediatedtransfer.ActionInitTargetStateChange:
		quitName = "ActionInitTargetStateChange"
		msg = st2.Message
		mgr.ChannelAddress = st2.FromRoute.ChannelIdentifier
	case *mediatedtransfer.ReceiveSecretRequestStateChange:
		quitName = "ReceiveSecretRequestStateChange"
		msg = st2.Message
		//todo refund 重新定义了
	case *mediatedtransfer.ReceiveAnnounceDisposedStateChange:
		quitName = "ReceiveAnnounceDisposedStateChange"
		msg = st2.Message
		mgr.ChannelAddresRefund = st2.Message.ChannelIdentifier
	case *mediatedtransfer.ReceiveBalanceProofStateChange:
		quitName = "ReceiveBalanceProofStateChange"
		_, ok := st2.Message.(*encoding.UnLock)
		if ok {
			msg = st2.Message //可能是mediated transfer,direct transfer,refundtransfer,secret 四中情况触发.
		}
	case *mediatedtransfer.ActionInitMediatorStateChange:
		quitName = "ActionInitMediatorStateChange"
		msg = st2.Message
		mgr.ChannelAddress = st2.FromRoute.ChannelIdentifier
	case *mediatedtransfer.ActionInitInitiatorStateChange:
		quitName = "ActionInitInitiatorStateChange"
		mgr.LastReceivedMessage = st2
		//new transfer trigger from user
	case *mediatedtransfer.ReceiveSecretRevealStateChange:
		quitName = "ReceiveSecretRevealStateChange"
		//reveal secret 需要单独处理
	}
	if msg != nil {
		mgr.ManagerState = transfer.StateManagerReceivedMessage
		mgr.LastReceivedMessage = msg
		tag := msg.Tag().(*transfer.MessageTag)
		tag.SetStateManager(mgr)
		msg.SetTag(tag)
		//tx := eh.raiden.db.StartTx()
		//eh.raiden.db.UpdateStateManaer(mgr, tx)
		//if mgr.ChannelIdentifier != utils.EmptyAddress {
		//	ch := eh.raiden.getChannelWithAddr(mgr.ChannelIdentifier)
		//	eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
		//}
		//tx.Commit()
		eh.raiden.conditionQuit(quitName)
	}
}
func (eh *stateMachineEventHandler) updateStateManagerFromEvent(receiver common.Address, msg encoding.Messager, mgr *transfer.StateManager) {
	var msgtoSend encoding.Messager
	switch msg2 := msg.(type) {
	case *encoding.MediatedTransfer:
		msgtoSend = msg2
		if mgr.Name == mediator.NameMediatorTransition {
			mgr.ChannelAddressTo = msg2.ChannelIdentifier
		} else {
			mgr.ChannelAddress = msg2.ChannelIdentifier
		}
	case *encoding.UnLock:
		msgtoSend = msg2
	case *encoding.AnnounceDisposed:
		msgtoSend = msg2 //state manager should be marked as finished? todo
	case *encoding.SecretRequest:
		msgtoSend = msg2
	case *encoding.RemoveExpiredHashlockTransfer:
		msgtoSend = msg2
	case *encoding.AnnounceDisposedResponse:
		msgtoSend = msg2
	default:
		panic(fmt.Sprintf("unknown message updateStateManagerFromEvent :%s", utils.StringInterface(msg, 3)))
	}
	tag := &transfer.MessageTag{
		MessageID:         utils.RandomString(10),
		EchoHash:          utils.Sha3(msg.Pack(), receiver[:]),
		IsASendingMessage: true,
		Receiver:          receiver,
	}
	tag.SetStateManager(mgr)
	msgtoSend.SetTag(tag)
	tx := eh.raiden.db.StartTx()
	mgr.ManagerState = transfer.StateManagerSendMessage
	mgr.LastSendMessage = msgtoSend
	msg3, ok := mgr.LastReceivedMessage.(encoding.Messager) //maybe  ActionInitInitiatorStateChange
	if ok {
		receiveTag := msg3.Tag()
		if receiveTag == nil {
			panic("must not be empty")
		}
		receiveMessageTag := receiveTag.(*transfer.MessageTag)
		if receiveMessageTag.ReceiveProcessComplete == false {
			mgr.ManagerState = transfer.StateManagerReceivedMessageProcessComplete
			log.Trace(fmt.Sprintf("set message %s ReceiveProcessComplete", receiveMessageTag.MessageID))
			receiveMessageTag.ReceiveProcessComplete = true
			ack := eh.raiden.Protocol.CreateAck(receiveMessageTag.EchoHash)
			eh.raiden.db.SaveAck(receiveMessageTag.EchoHash, ack.Pack(), tx)
		}
	} else {
		//user start a transfer.
	}
	eh.raiden.db.UpdateStateManaer(mgr, tx)
	eh.raiden.conditionQuit("InEventTx")
	if mgr.ChannelAddress == utils.EmptyHash {
		panic("channel address must not be empty")
	}
	ch, err := eh.raiden.findChannelByAddress(mgr.ChannelAddress)
	if err != nil {
		panic(fmt.Sprintf("channel %s must exist", utils.HPex(mgr.ChannelAddress)))
	}
	eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	if mgr.ChannelAddressTo != utils.EmptyHash { //for mediated transfer
		ch, err := eh.raiden.findChannelByAddress(mgr.ChannelAddressTo)
		if err != nil {
			panic(fmt.Sprintf("channel %s must exist", utils.HPex(mgr.ChannelAddressTo)))
		}
		eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	}
	if mgr.ChannelAddresRefund != utils.EmptyHash { //for mediated transfer and initiator

	}
	tx.Commit()
}
