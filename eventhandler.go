package smartraiden

import (
	"fmt"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

var errSentFailed = errors.New("sent failed")

//run inside loop of raiden service
type StateMachineEventHandler struct {
	raiden *RaidenService
}

func NewStateMachineEventHandler(raiden *RaidenService) *StateMachineEventHandler {
	h := &StateMachineEventHandler{
		raiden: raiden,
	}
	return h
}

/*
Log a state change, dispatch it to all state managers and log generated events
*/
func (this *StateMachineEventHandler) LogAndDispatchToAllTasks(st transfer.StateChange) {
	stateChangeId, _ := this.raiden.db.LogStateChange(st)
	for _, mgrs := range this.raiden.Identifier2StateManagers {
		for _, mgr := range mgrs {
			events := this.Dispatch(mgr, st)
			this.raiden.db.LogEvents(stateChangeId, events, this.raiden.GetBlockNumber())
		}

	}
}

/*
Log a state change, dispatch it to the state manager corresponding to `idenfitier`
        and log generated events
*/
func (this *StateMachineEventHandler) LogAndDispatchByIdentifier(identifier uint64, st transfer.StateChange) {
	stateChangeId, _ := this.raiden.db.LogStateChange(st)
	mgrs := this.raiden.Identifier2StateManagers[identifier]
	for _, mgr := range mgrs {
		events := this.Dispatch(mgr, st)
		this.raiden.db.LogEvents(stateChangeId, events, this.raiden.GetBlockNumber())
	}
}

//Log a state change, dispatch it to the given state manager and log generated events
func (this *StateMachineEventHandler) LogAndDispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) []transfer.Event {
	stateChangeId, _ := this.raiden.db.LogStateChange(stateChange)
	events := this.Dispatch(stateManager, stateChange)
	this.raiden.db.LogEvents(stateChangeId, events, this.raiden.GetBlockNumber())
	return events
}
func (this *StateMachineEventHandler) Dispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) (events []transfer.Event) {
	this.updateStateManagerFromReceivedMessageOrUserRequest(stateManager, stateChange)
	events = stateManager.Dispatch(stateChange)
	for _, e := range events {
		err := this.OnEvent(e, stateManager)
		if err != nil {
			log.Error(fmt.Sprintf("StateMachineEventHandler Dispatch:%v\n", err))
		}
	}
	return
}
func (this *StateMachineEventHandler) eventSendMediatedTransfer(event *mediated_transfer.EventSendMediatedTransfer, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := this.raiden.GetToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateMediatedTransfer(event.Initiator, event.Target, event.Fee, event.Amount, event.Identifier, event.Expiration, event.HashLock)
	if err != nil {
		return
	}
	mtr.Sign(this.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(this.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	this.updateStateManagerFromEvent(receiver, mtr, stateManager)
	this.raiden.ConditionQuit("EventSendMediatedTransferBefore")
	err = this.raiden.SendAsync(receiver, mtr)
	return
}
func (this *StateMachineEventHandler) eventSendBalanceProof(event *mediated_transfer.EventSendBalanceProof, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := this.raiden.GetToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	tr, err := ch.CreateSecret(event.Identifier, event.Secret)
	if err != nil {
		return
	}
	tr.Sign(this.raiden.PrivateKey, tr)
	err = ch.RegisterTransfer(this.raiden.GetBlockNumber(), tr)
	if err != nil {
		return
	}
	this.updateStateManagerFromEvent(receiver, tr, stateManager)
	this.raiden.ConditionQuit("EventSendBalanceProofBefore")
	err = this.raiden.SendAsync(receiver, tr)
	return
}
func (this *StateMachineEventHandler) eventSendRefundTransfer(event *mediated_transfer.EventSendRefundTransfer, stateManager *transfer.StateManager) (err error) {
	receiver := event.Receiver
	graph := this.raiden.GetToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateRefundTransfer(event.Initiator, event.Target, utils.BigInt0, event.Amount, event.Identifier, event.Expiration, event.HashLock)
	if err != nil {
		return
	}
	mtr.Sign(this.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(this.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	this.updateStateManagerFromEvent(receiver, mtr, stateManager)
	this.raiden.ConditionQuit("EventSendRefundTransferBefore")
	err = this.raiden.SendAsync(receiver, mtr)
	return
}
func (this *StateMachineEventHandler) eventContractSendChannelClose(event *mediated_transfer.EventContractSendChannelClose) (err error) {
	graph := this.raiden.GetToken2ChannelGraph(event.Token)
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
func (this *StateMachineEventHandler) eventWithdrawFailed(e2 *mediated_transfer.EventWithdrawFailed, manager *transfer.StateManager) (err error) {
	//wait from RemoveExpiredHashlockTransfer from partner.
	return nil
	//if manager.Name != target.NameTargetTransition && manager.Name != mediator.NameMediatorTransition {
	//	panic("EventWithdrawFailed can only comes from a target node or mediated node")
	//}
	//ch, err := this.raiden.FindChannelByAddress(e2.ChannelAddress)
	//if err != nil {
	//	log.Error(fmt.Sprintf("payer's lock expired ,but cannot find channel %s, this may happen long later restart after a stop"))
	//	return
	//}
	//log.Info(fmt.Sprint("remove expired hashlock channel=%s,hashlock=%s", utils.APex(e2.ChannelAddress), utils.HPex(e2.Hashlock)))
	//return ch.RemoveOurExpiredHashlock(e2.Hashlock, this.raiden.GetBlockNumber())
}
func (this *StateMachineEventHandler) eventContractSendWithdraw(e2 *mediated_transfer.EventContractSendWithdraw, manager *transfer.StateManager) (err error) {
	if manager.Name != target.NameTargetTransition && manager.Name != mediator.NameMediatorTransition {
		panic("EventWithdrawFailed can only comes from a target node or mediated node")
	}
	ch, err := this.raiden.FindChannelByAddress(e2.ChannelAddress)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, this may happen long later restart after a stop"))
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
func (this *StateMachineEventHandler) eventUnlockFailed(e2 *mediated_transfer.EventUnlockFailed, manager *transfer.StateManager) (err error) {
	if manager.Name != mediator.NameMediatorTransition && manager.Name != initiator.NameInitiatorTransition {
		panic("event unlock failed only happen for a mediated node")
	}
	ch, err := this.raiden.FindChannelByAddress(e2.ChannelAddress)
	if err != nil {
		log.Error(fmt.Sprintf("payee's lock expired ,but cannot find channel %s, this may happen long later restart after a stop"))
		return
	}
	log.Info(fmt.Sprintf("remove expired hashlock channel=%s,hashlock=%s ", utils.APex(e2.ChannelAddress), utils.HPex(e2.Hashlock)))
	tr, err := ch.CreateRemoveExpiredHashLockTransfer(e2.Hashlock, this.raiden.GetBlockNumber())
	if err != nil {
		log.Warn(fmt.Sprintf("Get Event UnlockFailed ,but hashlock cannot be removed err:%s", err))
		return
	}
	tr.Sign(this.raiden.PrivateKey, tr)
	err = ch.RegisterRemoveExpiredHashlockTransfer(tr, this.raiden.GetBlockNumber())
	if err != nil {
		log.Error(fmt.Sprintf("register mine RegisterRemoveExpiredHashlockTransfer err %s", err))
		return
	}
	/*
		save new channel status and sent RemoveExpiredHashlockTransfer must be atomic.
	*/
	tx := this.raiden.db.StartTx()
	this.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	this.raiden.db.NewSentRemoveExpiredHashlockTransfer(tr, ch.PartnerState.Address, tx)
	tx.Commit()
	err = this.raiden.SendAsync(ch.PartnerState.Address, tr)
	return
}
func (this *StateMachineEventHandler) OnEvent(event transfer.Event, stateManager *transfer.StateManager) (err error) {
	switch e2 := event.(type) {
	case *mediated_transfer.EventSendMediatedTransfer:
		err = this.eventSendMediatedTransfer(e2, stateManager)
		this.raiden.ConditionQuit("EventSendMediatedTransferAfter")
	case *mediated_transfer.EventSendRevealSecret:
		this.raiden.ConditionQuit("EventSendRevealSecretBefore")
		revealMessage := encoding.NewRevealSecret(e2.Secret)
		revealMessage.Sign(this.raiden.PrivateKey, revealMessage)
		err = this.raiden.SendAsync(e2.Receiver, revealMessage) //单独处理 reaveal secret
		this.raiden.ConditionQuit("EventSendRevealSecretAfter")
	case *mediated_transfer.EventSendBalanceProof:
		//unlock and update remotely (send the Secret message)
		err = this.eventSendBalanceProof(e2, stateManager)
		this.raiden.ConditionQuit("EventSendBalanceProofAfter")
	case *mediated_transfer.EventSendSecretRequest:
		secretRequest := encoding.NewSecretRequest(e2.Identifer, e2.Hashlock, e2.Amount)
		secretRequest.Sign(this.raiden.PrivateKey, secretRequest)
		this.updateStateManagerFromEvent(e2.Receiver, secretRequest, stateManager)
		this.raiden.ConditionQuit("EventSendSecretRequestBefore")
		err = this.raiden.SendAsync(e2.Receiver, secretRequest)
		this.raiden.ConditionQuit("EventSendSecretRequestAfter")
	case *mediated_transfer.EventSendRefundTransfer:
		err = this.eventSendRefundTransfer(e2, stateManager)
		this.raiden.ConditionQuit("EventSendRefundTransferAfter")
	case *transfer.EventTransferSentSuccess:
		this.finishOneTransfer(event)
	case *transfer.EventTransferSentFailed:
		this.finishOneTransfer(event)
	case *transfer.EventTransferReceivedSuccess:
	case *mediated_transfer.EventUnlockSuccess:
	case *mediated_transfer.EventWithdrawFailed:
		//TODO need payer's new signature to remove this expired lock
		log.Error(fmt.Sprintf("EventWithdrawFailed hashlock=%s,reason=%s", utils.HPex(e2.Hashlock), e2.Reason))
		err = this.eventWithdrawFailed(e2, stateManager)
	case *mediated_transfer.EventWithdrawSuccess:
		/*
					  The withdraw is currently handled by the netting channel, once the close
			     event is detected all locks will be withdrawn
		*/
	case *mediated_transfer.EventContractSendWithdraw:
		//do nothing for five events above
		err = this.eventContractSendWithdraw(e2, stateManager)
	case *mediated_transfer.EventUnlockFailed:
		//should remove hashlock from channel todo fix bai
		log.Error(fmt.Sprintf("unlockfailed hashlock=%s,reason=%s", utils.HPex(e2.Hashlock), e2.Reason))
		err = this.eventUnlockFailed(e2, stateManager)
	case *mediated_transfer.EventContractSendChannelClose:
		err = this.eventContractSendChannelClose(e2)
	default:
		err = fmt.Errorf("unkown event :%s", utils.StringInterface1(event))
		log.Error(err.Error())
	}
	return
}
func (this *StateMachineEventHandler) finishOneTransfer(ev transfer.Event) {
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
	results := this.raiden.Identifier2Results[identifier]
	if len(results) <= 0 { //restart after crash?
		log.Error(fmt.Sprintf("transfer finished ,but have no relate results :%s", utils.StringInterface(ev, 2)))
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
		delete(this.raiden.Identifier2Results, identifier)
	} else {
		this.raiden.Identifier2Results[identifier] = results
	}
}
func (this *StateMachineEventHandler) HandleTokenAdded(st *mediated_transfer.ContractReceiveTokenAddedStateChange) error {
	managerAddress := st.ManagerAddress
	return this.raiden.RegisterChannelManager(managerAddress)
}
func (this *StateMachineEventHandler) handleChannelNew(st *mediated_transfer.ContractReceiveNewChannelStateChange) error {
	managerAddress := st.ManagerAddress
	ChannelAddres := st.ChannelAddress
	participant1 := st.Participant1
	participant2 := st.Participant2
	tokenAddress := this.raiden.Manager2Token[managerAddress]
	graph := this.raiden.GetToken2ChannelGraph(tokenAddress)
	graph.AddPath(participant1, participant2)
	connectionManager, err := this.raiden.ConnectionManagerForToken(tokenAddress)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	isParticipant := this.raiden.NodeAddress == participant2 || this.raiden.NodeAddress == participant1
	isBootstrap := connectionManager.BOOTSTRAP_ADDR == participant1 || connectionManager.BOOTSTRAP_ADDR == participant2
	if isParticipant {
		this.raiden.RegisterNettingChannel(tokenAddress, ChannelAddres)
		if !isBootstrap {
			other := participant2
			if other == this.raiden.NodeAddress {
				other = participant1
			}
			this.raiden.StartHealthCheckFor(other)
		}
	} else if connectionManager.WantsMoreChannels() {
		go func() {
			connectionManager.RetryConnect()
		}()
	} else {
		log.Info("ignoring new channel, this node is not a participant.")
	}
	return nil
}

func (this *StateMachineEventHandler) handleBalance(st *mediated_transfer.ContractReceiveBalanceStateChange) error {
	channelAddress := st.ChannelAddress
	tokenAddress := st.TokenAddress
	participant := st.ParticipantAddress
	balance := st.Balance
	graph := this.raiden.GetToken2ChannelGraph(tokenAddress)
	ch := graph.GetChannelAddress2Channel(channelAddress)
	err := this.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	this.raiden.db.UpdateChannelContractBalance(channel.NewChannelSerialization(ch))
	if ch.ContractBalance().Cmp(utils.BigInt0) == 0 {
		connectionManager, _ := this.raiden.ConnectionManagerForToken(tokenAddress)
		go func() {
			connectionManager.JoinChannel(participant, balance)
		}()
	}
	return nil
}

func (this *StateMachineEventHandler) handleClosed(st *mediated_transfer.ContractReceiveClosedStateChange) error {
	channelAddress := st.ChannelAddress
	ch, err := this.raiden.FindChannelByAddress(channelAddress)
	if err != nil {
		return err
	}
	err = this.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = this.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}

func (this *StateMachineEventHandler) handleSettled(st *mediated_transfer.ContractReceiveSettledStateChange) error {
	//todo remove channel st.channelAddress ,because this channel is already settled
	log.Trace(fmt.Sprintf("%s settled event handle", st.ChannelAddress.String()))
	ch, err := this.raiden.FindChannelByAddress(st.ChannelAddress)
	if err != nil {
		return err
	}
	err = this.ChannelStateTransition(ch, st)
	if err != nil {
		log.Error(fmt.Sprintf("handleBalance ChannelStateTransition err=%s", err))
	}
	err = this.raiden.db.UpdateChannelState(channel.NewChannelSerialization(ch))
	return err
}
func (this *StateMachineEventHandler) handleWithdraw(st *mediated_transfer.ContractReceiveWithdrawStateChange) error {
	this.raiden.RegisterSecret(st.Secret)
	return nil
}

//avoid dead lock
func (this *StateMachineEventHandler) ChannelStateTransition(c *channel.Channel, st transfer.StateChange) (err error) {
	switch st2 := st.(type) {
	case *transfer.BlockStateChange:
		if c.State() == transfer.CHANNEL_STATE_CLOSED {
			settlementEnd := c.ExternState.ClosedBlock + int64(c.SettleTimeout)
			if st2.BlockNumber > settlementEnd {
				//should not block todo fix it
				//err = c.ExternState.Settle()
			}
		}
	case *mediated_transfer.ContractReceiveClosedStateChange:
		if st2.ChannelAddress == c.MyAddress {
			if !c.IsCloseEventComplete {
				c.ExternState.SetClosed(st2.ClosedBlock)
				//should not block todo fix it
				c.HandleClosed(st2.ClosedBlock, st2.ClosingAddress)
			} else {
				log.Warn(fmt.Sprintf("channel closed on a different block or close event happened twice channel=%s,closedblock=%s,thisblock=%s",
					c.MyAddress.String(), c.ExternState.ClosedBlock, st2.ClosedBlock))
			}
		}
	case *mediated_transfer.ContractReceiveSettledStateChange:
		//settled channel should be removed. todo bai fix it
		if st2.ChannelAddress == c.MyAddress {
			if c.ExternState.SetSettled(st2.SettledBlock) {
				c.HandleSettled(st2.SettledBlock)
			} else {
				log.Warn(fmt.Sprintf("channel is already settled on a different block channeladdress=%s,settleblock=%d,thisblock=%d",
					c.MyAddress.String(), c.ExternState.SettledBlock, st2.SettledBlock))
			}
		}
	case *mediated_transfer.ContractReceiveBalanceStateChange:
		participant := st2.ParticipantAddress
		balance := st2.Balance
		var channelState *channel.ChannelEndState
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
func (this *StateMachineEventHandler) filterStateChange(st transfer.StateChange) bool {
	var channelAddress common.Address
	switch st2 := st.(type) { //filter event only about me
	case *mediated_transfer.ContractReceiveTokenAddedStateChange:
		if st2.RegistryAddress != this.raiden.RegistryAddress { //there maybe many registry contracts on blockchain
			return false
		}
		return true
	case *mediated_transfer.ContractReceiveNewChannelStateChange:
		if this.raiden.Manager2Token[st2.ManagerAddress] == utils.EmptyAddress { //newchannel on another registry contract
			return false
		}
		return true
	case *mediated_transfer.ContractReceiveBalanceStateChange:
		channelAddress = st2.ChannelAddress
	case *mediated_transfer.ContractReceiveClosedStateChange:
		channelAddress = st2.ChannelAddress
	case *mediated_transfer.ContractReceiveSettledStateChange:
		channelAddress = st2.ChannelAddress
	case *mediated_transfer.ContractReceiveWithdrawStateChange:
		channelAddress = st2.ChannelAddress
	default:
		err := fmt.Errorf("OnBlockchainStateChange unknown statechange :%s", utils.StringInterface1(st))
		log.Error(err.Error())
		return false
	}
	found := false
	for _, g := range this.raiden.Token2ChannelGraph {
		ch := g.GetChannelAddress2Channel(channelAddress)
		if ch != nil {
			found = true
			break
		}
	}
	return found
}
func (this *StateMachineEventHandler) OnBlockchainStateChange(st transfer.StateChange) (err error) {
	log.Trace(fmt.Sprintf("statechange received :%s", utils.StringInterface1(st)))
	_, err = this.raiden.db.LogStateChange(st)
	if err != nil {
		return err
	}
	if !this.filterStateChange(st) {
		return nil
	}
	switch st2 := st.(type) {
	case *mediated_transfer.ContractReceiveTokenAddedStateChange:
		err = this.HandleTokenAdded(st2)
	case *mediated_transfer.ContractReceiveNewChannelStateChange:
		err = this.handleChannelNew(st2)
	case *mediated_transfer.ContractReceiveBalanceStateChange:
		err = this.handleBalance(st2)
	case *mediated_transfer.ContractReceiveClosedStateChange:
		err = this.handleClosed(st2)
	case *mediated_transfer.ContractReceiveSettledStateChange:
		err = this.handleSettled(st2)
	case *mediated_transfer.ContractReceiveWithdrawStateChange:
		err = this.handleWithdraw(st2)
	default:
		err = fmt.Errorf("OnBlockchainStateChange unknown statechange :%s", utils.StringInterface1(st))
		log.Error(err.Error())
	}
	return
}

//recive a message and before processed
func (this *StateMachineEventHandler) updateStateManagerFromReceivedMessageOrUserRequest(mgr *transfer.StateManager, stateChange transfer.StateChange) {
	var msg encoding.Messager
	var quitName string
	switch st2 := stateChange.(type) {
	case *mediated_transfer.ActionInitTargetStateChange:
		quitName = "ActionInitTargetStateChange"
		msg = st2.Message
		mgr.ChannelAddress = st2.FromRoute.ChannelAddress
	case *mediated_transfer.ReceiveSecretRequestStateChange:
		quitName = "ReceiveSecretRequestStateChange"
		msg = st2.Message
	case *mediated_transfer.ReceiveTransferRefundStateChange:
		quitName = "ReceiveTransferRefundStateChange"
		msg = st2.Message
		mgr.ChannelAddresRefund = st2.Message.Channel
	case *mediated_transfer.ReceiveBalanceProofStateChange:
		quitName = "ReceiveBalanceProofStateChange"
		_, ok := st2.Message.(*encoding.Secret)
		if ok {
			msg = st2.Message //可能是mediated transfer,direct transfer,refundtransfer,secret 四中情况触发.
		}
	case *mediated_transfer.ActionInitMediatorStateChange:
		quitName = "ActionInitMediatorStateChange"
		msg = st2.Message
		mgr.ChannelAddress = st2.FromRoute.ChannelAddress
	case *mediated_transfer.ActionInitInitiatorStateChange:
		quitName = "ActionInitInitiatorStateChange"
		mgr.LastReceivedMessage = st2
		//new transfer trigger from user
	case *mediated_transfer.ReceiveSecretRevealStateChange:
		quitName = "ReceiveSecretRevealStateChange"
		//reveal secret 需要单独处理
	}
	if msg != nil {
		mgr.ManagerState = transfer.StateManager_ReceivedMessage
		mgr.LastReceivedMessage = msg
		tag := msg.Tag().(*transfer.MessageTag)
		tag.SetStateManager(mgr)
		msg.SetTag(tag)
		//tx := this.raiden.db.StartTx()
		//this.raiden.db.UpdateStateManaer(mgr, tx)
		//if mgr.ChannelAddress != utils.EmptyAddress {
		//	ch := this.raiden.GetChannelWithAddr(mgr.ChannelAddress)
		//	this.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
		//}
		//tx.Commit()
		this.raiden.ConditionQuit(quitName)
	}
}
func (this *StateMachineEventHandler) updateStateManagerFromEvent(receiver common.Address, msg encoding.Messager, mgr *transfer.StateManager) {
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
		MessageId:         utils.RandomString(10),
		EchoHash:          utils.Sha3(msg.Pack(), receiver[:]),
		IsASendingMessage: true,
		Receiver:          receiver,
	}
	tag.SetStateManager(mgr)
	msgtoSend.SetTag(tag)
	tx := this.raiden.db.StartTx()
	mgr.ManagerState = transfer.StateManager_SendMessage
	mgr.LastSendMessage = msgtoSend
	msg3, ok := mgr.LastReceivedMessage.(encoding.Messager) //maybe  ActionInitInitiatorStateChange
	if ok {
		receiveTag := msg3.Tag()
		if receiveTag == nil {
			panic("must not be empty")
		}
		receiveMessageTag := receiveTag.(*transfer.MessageTag)
		if receiveMessageTag.ReceiveProcessComplete == false {
			mgr.ManagerState = transfer.StateManager_ReceivedMessageProcessComplete
			log.Trace(fmt.Sprintf("set message %s ReceiveProcessComplete", receiveMessageTag.MessageId))
			receiveMessageTag.ReceiveProcessComplete = true
			ack := this.raiden.Protocol.CreateAck(receiveMessageTag.EchoHash)
			this.raiden.db.SaveAck(receiveMessageTag.EchoHash, ack.Pack(), tx)
		}
	} else {
		//user start a transfer.
	}
	this.raiden.db.UpdateStateManaer(mgr, tx)
	this.raiden.ConditionQuit("InEventTx")
	if mgr.ChannelAddress == utils.EmptyAddress {
		panic("channel address must not be empty")
	}
	ch, err := this.raiden.FindChannelByAddress(mgr.ChannelAddress)
	if err != nil {
		panic(fmt.Sprintf("channel %s must exist", utils.APex(mgr.ChannelAddress)))
	}
	this.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	if mgr.ChannelAddressTo != utils.EmptyAddress { //for mediated transfer
		ch, err := this.raiden.FindChannelByAddress(mgr.ChannelAddressTo)
		if err != nil {
			panic(fmt.Sprintf("channel %s must exist", utils.APex(mgr.ChannelAddressTo)))
		}
		this.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
	}
	if mgr.ChannelAddresRefund != utils.EmptyAddress { //for mediated transfer and initiator
		_, isrefund := mgr.LastReceivedMessage.(*encoding.RefundTransfer)
		islocktransfer := encoding.IsLockedTransfer(mgr.LastSendMessage) //when receive refund transfer, next message must be a refund transfer or mediated transfer.
		if isrefund && islocktransfer {
			ch, err := this.raiden.FindChannelByAddress(mgr.ChannelAddresRefund)
			if err != nil {
				panic(fmt.Sprintf("channel %s must exist", mgr.ChannelAddresRefund))
			}
			this.raiden.db.UpdateChannel(channel.NewChannelSerialization(ch), tx)
			mgr.ChannelAddresRefund = utils.EmptyAddress
		} else {
			panic("last received message must be a refund transfer and last send must be a mediated transfer")
		}
	}
	tx.Commit()
}
