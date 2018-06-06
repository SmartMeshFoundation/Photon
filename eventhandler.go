package smartraiden

import (
	"fmt"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
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
	stateChangeID, _ := eh.raiden.db.LogStateChange(st)
	for _, mgrs := range eh.raiden.Identifier2StateManagers {
		for _, mgr := range mgrs {
			events := eh.dispatch(mgr, st)
			eh.raiden.db.LogEvents(stateChangeID, events, eh.raiden.GetBlockNumber())
		}

	}
}

/*
Log a state change, dispatch it to the state manager corresponding to `idenfitier`
        and log generated events
*/
func (eh *stateMachineEventHandler) logAndDispatchByIdentifier(identifier uint64, st transfer.StateChange) {
	stateChangeID, _ := eh.raiden.db.LogStateChange(st)
	mgrs := eh.raiden.Identifier2StateManagers[identifier]
	for _, mgr := range mgrs {
		events := eh.dispatch(mgr, st)
		eh.raiden.db.LogEvents(stateChangeID, events, eh.raiden.GetBlockNumber())
	}
}

//Log a state change, dispatch it to the given state manager and log generated events
func (eh *stateMachineEventHandler) logAndDispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) []transfer.Event {
	stateChangeID, _ := eh.raiden.db.LogStateChange(stateChange)
	events := eh.dispatch(stateManager, stateChange)
	eh.raiden.db.LogEvents(stateChangeID, events, eh.raiden.GetBlockNumber())
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
	mtr, err := ch.CreateMediatedTransfer(event.Initiator, event.Target, event.Fee, event.Amount, event.Identifier, event.Expiration, event.HashLock)
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
func (eh *stateMachineEventHandler) eventSendBalanceProof(event *mediatedtransfer.EventSendBalanceProof, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	tr, err := ch.CreateSecret(event.Identifier, event.Secret)
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
func (eh *stateMachineEventHandler) eventSendRefundTransfer(event *mediatedtransfer.EventSendRefundTransfer, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateRefundTransfer(event.Initiator, event.Target, utils.BigInt0, event.Amount, event.Identifier, event.Expiration, event.HashLock)
	if err != nil {
		return
	}
	mtr.Sign(eh.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(eh.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	eh.updateStateManagerFromEvent(receiver, mtr, stateManager)
	eh.raiden.conditionQuit("EventSendRefundTransferBefore")
	err = eh.raiden.sendAsync(receiver, mtr)
	return
}
func (eh *stateMachineEventHandler) eventContractSendChannelClose(event *mediatedtransfer.EventContractSendChannelClose) (err error) {
	graph := eh.raiden.getToken2ChannelGraph(event.Token)
	if graph == nil {
		err = fmt.Errorf("EventContractSendChannelClose but token %s doesn't exist", utils.APex(event.Token))
		return
	}
	ch := graph.ChannelAddress2Channel[event.ChannelAddress]
	if ch == nil {
		err = fmt.Errorf("EventContractSendChannelClose  but channel %s doesn't exist,maybe have already settled", utils.APex(event.ChannelAddress))
		return
	}
	balanceProof := ch.OurState.BalanceProofState
	err = ch.ExternState.Close(balanceProof)
	return
}
func (eh *stateMachineEventHandler) eventWithdrawFailed(e2 *mediatedtransfer.EventWithdrawFailed, manager *transfer.StateManager) (err error) {
	//wait from RemoveExpiredHashlockTransfer from partner.
	return nil
	//if manager.Name != target.NameTargetTransition && manager.Name != mediator.NameMediatorTransition {
	//	panic("EventWithdrawFailed can only comes from a target node or mediated node")
	//}
	//ch, err := eh.raiden.findChannelByAddress(e2.ChannelAddress)
	//if err != nil {
	//	log.Error(fmt.Sprintf("payer's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop"))
	//	return
	//}
	//log.Info(fmt.Sprint("remove expired hashlock channel=%s,hashlock=%s", utils.APex(e2.ChannelAddress), utils.HPex(e2.Hashlock)))
	//return ch.RemoveOurExpiredHashlock(e2.Hashlock, eh.raiden.GetBlockNumber())
}
func (eh *stateMachineEventHandler) eventContractSendWithdraw(e2 *mediatedtransfer.EventContractSendWithdraw, manager *transfer.StateManager) (err error) {
	if manager.Name != target.NameTargetTransition && manager.Name != mediator.NameMediatorTransition {
		panic("EventWithdrawFailed can only comes from a target node or mediated node")
	}
	ch, err := eh.raiden.findChannelByAddress(e2.ChannelAddress)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", utils.APex(e2.ChannelAddress)))
		return
	}
	unlockProofs := ch.PartnerState.GetKnownUnlocks()
	err = ch.ExternState.WithDraw(unlockProofs)
	if err != nil {
		log.Error(fmt.Sprintf("withdraw on %s failed, channel is gone, error:%s", utils.APex(ch.MyAddress), err))
	}
	return nil
}

/*
the transfer I payed for a payee has expired. give a new balanceproof which doesn't contain this hashlock
*/
func (eh *stateMachineEventHandler) eventUnlockFailed(e2 *mediatedtransfer.EventUnlockFailed, manager *transfer.StateManager) (err error) {
	if manager.Name != mediator.NameMediatorTransition && manager.Name != initiator.NameInitiatorTransition {
		panic("event unlock failed only happen for a mediated node")
	}
	ch, err := eh.raiden.findChannelByAddress(e2.ChannelAddress)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, eh may happen long later restart after a stop", utils.APex(e2.ChannelAddress)))
		return
	}
	log.Info(fmt.Sprintf("remove expired hashlock channel=%s,hashlock=%s ", utils.APex(e2.ChannelAddress), utils.HPex(e2.Hashlock)))
	tr, err := ch.CreateRemoveExpiredHashLockTransfer(e2.Hashlock, eh.raiden.GetBlockNumber())
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
	tx := eh.raiden.db.StartTx()
	eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	eh.raiden.db.NewSentRemoveExpiredHashlockTransfer(tr, ch.PartnerState.Address, tx)
	tx.Commit()
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
		//unlock and update remotely (send the Secret message)
		err = eh.eventSendBalanceProof(e2, stateManager)
		eh.raiden.conditionQuit("EventSendBalanceProofAfter")
	case *mediatedtransfer.EventSendSecretRequest:
		secretRequest := encoding.NewSecretRequest(e2.Identifer, e2.Hashlock, e2.Amount)
		secretRequest.Sign(eh.raiden.PrivateKey, secretRequest)
		eh.updateStateManagerFromEvent(e2.Receiver, secretRequest, stateManager)
		eh.raiden.conditionQuit("EventSendSecretRequestBefore")
		err = eh.raiden.sendAsync(e2.Receiver, secretRequest)
		eh.raiden.conditionQuit("EventSendSecretRequestAfter")
	case *mediatedtransfer.EventSendRefundTransfer:
		err = eh.eventSendRefundTransfer(e2, stateManager)
		eh.raiden.conditionQuit("EventSendRefundTransferAfter")
	case *transfer.EventTransferSentSuccess:
		ch := eh.raiden.getChannelWithAddr(e2.ChannelAddress)
		if ch == nil {
			err = fmt.Errorf("receive EventTransferSentSuccess,but channel not exist %s", utils.APex(e2.ChannelAddress))
			return
		}
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		if err != nil {
			log.Error(fmt.Sprintf("UpdateChannelNoTx err %s", err))
		}
		eh.raiden.db.NewSentTransfer(eh.raiden.GetBlockNumber(), e2.ChannelAddress, ch.TokenAddress, e2.Target, ch.GetNextNonce(), e2.Amount)
		eh.finishOneTransfer(event)
	case *transfer.EventTransferSentFailed:
		eh.finishOneTransfer(event)
	case *transfer.EventTransferReceivedSuccess:
		ch := eh.raiden.getChannelWithAddr(e2.ChannelAddress)
		if ch == nil {
			err = fmt.Errorf("receive EventTransferReceivedSuccess,but channel not exist %s", utils.APex(e2.ChannelAddress))
			return
		}
		err = eh.raiden.db.UpdateChannelNoTx(channel.NewChannelSerialization(ch))
		if err != nil {
			log.Error(fmt.Sprintf("UpdateChannelNoTx err %s", err))
		}
		eh.raiden.db.NewReceivedTransfer(eh.raiden.GetBlockNumber(), e2.ChannelAddress, ch.TokenAddress, e2.Initiator, ch.PartnerState.BalanceProofState.Nonce, e2.Amount)
	case *mediatedtransfer.EventUnlockSuccess:
	case *mediatedtransfer.EventWithdrawFailed:
		//TODO need payer's new signature to remove eh expired lock
		log.Error(fmt.Sprintf("EventWithdrawFailed hashlock=%s,reason=%s", utils.HPex(e2.Hashlock), e2.Reason))
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
		log.Error(fmt.Sprintf("unlockfailed hashlock=%s,reason=%s", utils.HPex(e2.Hashlock), e2.Reason))
		err = eh.eventUnlockFailed(e2, stateManager)
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
	var identifier uint64
	var target common.Address
	switch e2 := ev.(type) {
	case *transfer.EventTransferSentSuccess:
		log.Info(fmt.Sprintf("EventTransferSentSuccess for id %d ", e2.Identifier))
		identifier = e2.Identifier
		target = e2.Target
		err = nil
	case *transfer.EventTransferSentFailed:
		log.Warn(fmt.Sprintf("EventTransferSentFailed for id %d,because of %s", e2.Identifier, e2.Reason))
		identifier = e2.Identifier
		target = e2.Target
		err = errors.New(e2.Reason)
	default:
		panic("unknow event")
	}
	results := eh.raiden.Identifier2Results[identifier]
	if len(results) <= 0 { //restart after crash?
		log.Error(fmt.Sprintf("you can ignore this error when this transfer is a direct transfer.\n transfer finished ,but have no relate results :%s", utils.StringInterface(ev, 2)))
		return
	}
	for i, r := range results {
		t2, ok := r.Tag.(common.Address)
		if !ok {
			panic("Identifier2Results's tag must be Address")
		}
		if t2 == target {
			r.Result <- err
			results = append(results[:i], results[i+1:]...)
			close(r.Result) //for tokenswap may error todo fix it. 为什么不让tokenswap使用两个有规律的id,而不是完全相同的两个id呢
			break
		}
	}
	if len(results) == 0 {
		delete(eh.raiden.Identifier2Results, identifier)
	} else {
		eh.raiden.Identifier2Results[identifier] = results
	}
}
func (eh *stateMachineEventHandler) HandleTokenAdded(st *mediatedtransfer.ContractReceiveTokenAddedStateChange) error {
	managerAddress := st.ManagerAddress
	return eh.raiden.registerChannelManager(managerAddress)
}
func (eh *stateMachineEventHandler) handleChannelNew(st *mediatedtransfer.ContractReceiveNewChannelStateChange) error {
	managerAddress := st.ManagerAddress
	ChannelAddres := st.ChannelAddress
	participant1 := st.Participant1
	participant2 := st.Participant2
	tokenAddress := eh.raiden.Manager2Token[managerAddress]
	graph := eh.raiden.getToken2ChannelGraph(tokenAddress)
	graph.AddPath(participant1, participant2)
	connectionManager, err := eh.raiden.connectionManagerForToken(tokenAddress)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	isParticipant := eh.raiden.NodeAddress == participant2 || eh.raiden.NodeAddress == participant1
	isBootstrap := connectionManager.BootstrapAddr == participant1 || connectionManager.BootstrapAddr == participant2
	if isParticipant {
		eh.raiden.registerNettingChannel(tokenAddress, ChannelAddres)
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
		log.Info("ignoring new channel, this node is not a participant.")
	}
	return nil
}

func (eh *stateMachineEventHandler) handleBalance(st *mediatedtransfer.ContractReceiveBalanceStateChange) error {
	channelAddress := st.ChannelAddress
	tokenAddress := st.TokenAddress
	participant := st.ParticipantAddress
	balance := st.Balance
	graph := eh.raiden.getToken2ChannelGraph(tokenAddress)
	ch := graph.GetChannelAddress2Channel(channelAddress)
	err := eh.ChannelStateTransition(ch, st)
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

func (eh *stateMachineEventHandler) handleClosed(st *mediatedtransfer.ContractReceiveClosedStateChange) error {
	channelAddress := st.ChannelAddress
	ch, err := eh.raiden.findChannelByAddress(channelAddress)
	if err != nil {
		return err
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}

func (eh *stateMachineEventHandler) handleSettled(st *mediatedtransfer.ContractReceiveSettledStateChange) error {
	//todo remove channel st.channelAddress ,because eh channel is already settled
	log.Trace(fmt.Sprintf("%s settled event handle", utils.APex(st.ChannelAddress)))
	ch, err := eh.raiden.findChannelByAddress(st.ChannelAddress)
	if err != nil {
		return err
	}
	err = eh.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
		return err
	}
	err = eh.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}
func (eh *stateMachineEventHandler) handleWithdraw(st *mediatedtransfer.ContractReceiveWithdrawStateChange) error {
	eh.raiden.registerSecret(st.Secret)
	return nil
}

//avoid dead lock
func (eh *stateMachineEventHandler) ChannelStateTransition(c *channel.Channel, st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *transfer.BlockStateChange:
		if c.State() == transfer.ChannelStateClosed {
			settlementEnd := c.ExternState.ClosedBlock + int64(c.SettleTimeout)
			if st2.BlockNumber > settlementEnd {
				//should not block todo fix it
				//err = c.ExternState.Settle()
			}
		}
	case *mediatedtransfer.ContractReceiveClosedStateChange:
		if st2.ChannelAddress == c.MyAddress {
			if !c.IsCloseEventComplete {
				c.ExternState.SetClosed(st2.ClosedBlock)
				//should not block todo fix it
				c.HandleClosed(st2.ClosedBlock, st2.ClosingAddress)
			} else {
				log.Warn(fmt.Sprintf("channel closed on a different block or close event happened twice channel=%s,closedblock=%d,thisblock=%d",
					c.MyAddress.String(), c.ExternState.ClosedBlock, st2.ClosedBlock))
			}
		}
	case *mediatedtransfer.ContractReceiveSettledStateChange:
		//settled channel should be removed. todo bai fix it
		if st2.ChannelAddress == c.MyAddress {
			if c.ExternState.SetSettled(st2.SettledBlock) {
				c.HandleSettled(st2.SettledBlock)
			} else {
				log.Warn(fmt.Sprintf("channel is already settled on a different block channeladdress=%s,settleblock=%d,thisblock=%d",
					c.MyAddress.String(), c.ExternState.SettledBlock, st2.SettledBlock))
			}
		}
	case *mediatedtransfer.ContractReceiveBalanceStateChange:
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
	}
	return

}

//only care statechanges about me
func (eh *stateMachineEventHandler) filterStateChange(st transfer.StateChange) bool {
	var channelAddress common.Address
	switch st2 := st.(type) { //filter event only about me
	case *mediatedtransfer.ContractReceiveTokenAddedStateChange:
		if st2.RegistryAddress != eh.raiden.RegistryAddress { //there maybe many registry contracts on blockchain
			return false
		}
		return true
	case *mediatedtransfer.ContractReceiveNewChannelStateChange:
		if eh.raiden.Manager2Token[st2.ManagerAddress] == utils.EmptyAddress { //newchannel on another registry contract
			return false
		}
		return true
	case *mediatedtransfer.ContractReceiveBalanceStateChange:
		channelAddress = st2.ChannelAddress
	case *mediatedtransfer.ContractReceiveClosedStateChange:
		channelAddress = st2.ChannelAddress
	case *mediatedtransfer.ContractReceiveSettledStateChange:
		channelAddress = st2.ChannelAddress
	case *mediatedtransfer.ContractReceiveWithdrawStateChange:
		channelAddress = st2.ChannelAddress
	case *mediatedtransfer.ContractTransferUpdatedStateChange:
		channelAddress = st2.ChannelAddress
	default:
		err := fmt.Errorf("OnBlockchainStateChange unknown statechange :%s", utils.StringInterface1(st))
		log.Error(err.Error())
		return false
	}
	found := false
	for _, g := range eh.raiden.Token2ChannelGraph {
		ch := g.GetChannelAddress2Channel(channelAddress)
		if ch != nil {
			found = true
			break
		}
	}
	return found
}
func (eh *stateMachineEventHandler) OnBlockchainStateChange(st transfer.StateChange) (err error) {
	log.Trace(fmt.Sprintf("statechange received :%s", utils.StringInterface(st, 2)))
	_, err = eh.raiden.db.LogStateChange(st)
	if err != nil {
		return err
	}
	if !eh.filterStateChange(st) {
		return nil
	}
	switch st2 := st.(type) {
	case *mediatedtransfer.ContractReceiveTokenAddedStateChange:
		err = eh.HandleTokenAdded(st2)
	case *mediatedtransfer.ContractReceiveNewChannelStateChange:
		err = eh.handleChannelNew(st2)
	case *mediatedtransfer.ContractReceiveBalanceStateChange:
		err = eh.handleBalance(st2)
	case *mediatedtransfer.ContractReceiveClosedStateChange:
		err = eh.handleClosed(st2)
	case *mediatedtransfer.ContractReceiveSettledStateChange:
		err = eh.handleSettled(st2)
	case *mediatedtransfer.ContractReceiveWithdrawStateChange:
		err = eh.handleWithdraw(st2)
	case *mediatedtransfer.ContractTransferUpdatedStateChange:
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
		mgr.ChannelAddress = st2.FromRoute.ChannelAddress
	case *mediatedtransfer.ReceiveSecretRequestStateChange:
		quitName = "ReceiveSecretRequestStateChange"
		msg = st2.Message
	case *mediatedtransfer.ReceiveTransferRefundStateChange:
		quitName = "ReceiveTransferRefundStateChange"
		msg = st2.Message
		mgr.ChannelAddresRefund = st2.Message.Channel
	case *mediatedtransfer.ReceiveBalanceProofStateChange:
		quitName = "ReceiveBalanceProofStateChange"
		_, ok := st2.Message.(*encoding.Secret)
		if ok {
			msg = st2.Message //可能是mediated transfer,direct transfer,refundtransfer,secret 四中情况触发.
		}
	case *mediatedtransfer.ActionInitMediatorStateChange:
		quitName = "ActionInitMediatorStateChange"
		msg = st2.Message
		mgr.ChannelAddress = st2.FromRoute.ChannelAddress
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
		//if mgr.ChannelAddress != utils.EmptyAddress {
		//	ch := eh.raiden.getChannelWithAddr(mgr.ChannelAddress)
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
			mgr.ChannelAddressTo = msg2.Channel
		} else {
			mgr.ChannelAddress = msg2.Channel
		}
	case *encoding.Secret:
		msgtoSend = msg2
	case *encoding.RefundTransfer:
		msgtoSend = msg2 //state manager should be marked as finished? todo
	case *encoding.SecretRequest:
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
	if mgr.ChannelAddress == utils.EmptyAddress {
		panic("channel address must not be empty")
	}
	ch, err := eh.raiden.findChannelByAddress(mgr.ChannelAddress)
	if err != nil {
		panic(fmt.Sprintf("channel %s must exist", utils.APex(mgr.ChannelAddress)))
	}
	eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	if mgr.ChannelAddressTo != utils.EmptyAddress { //for mediated transfer
		ch, err := eh.raiden.findChannelByAddress(mgr.ChannelAddressTo)
		if err != nil {
			panic(fmt.Sprintf("channel %s must exist", utils.APex(mgr.ChannelAddressTo)))
		}
		eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	}
	if mgr.ChannelAddresRefund != utils.EmptyAddress { //for mediated transfer and initiator
		_, isrefund := mgr.LastReceivedMessage.(*encoding.RefundTransfer)
		islocktransfer := encoding.IsLockedTransfer(mgr.LastSendMessage) //when receive refund transfer, next message must be a refund transfer or mediated transfer.
		if isrefund && islocktransfer {
			ch, err := eh.raiden.findChannelByAddress(mgr.ChannelAddresRefund)
			if err != nil {
				panic(fmt.Sprintf("channel %s must exist", mgr.ChannelAddresRefund))
			}
			eh.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
			mgr.ChannelAddresRefund = utils.EmptyAddress
		} else {
			panic("last received message must be a refund transfer and last send must be a mediated transfer")
		}
	}
	tx.Commit()
}
