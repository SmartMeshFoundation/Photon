package raiden_network

import (
	"fmt"

	"errors"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
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
	stateChangeId, _ := this.raiden.TransactionLog.Log(st)
	for _, mgrs := range this.raiden.Identifier2StateManagers {
		for _, mgr := range mgrs {
			events := this.Dispatch(mgr, st)
			this.raiden.TransactionLog.LogEvents(stateChangeId, events, this.raiden.GetBlockNumber())
		}

	}
}

/*
Log a state change, dispatch it to the state manager corresponding to `idenfitier`
        and log generated events
*/
func (this *StateMachineEventHandler) LogAndDispatchByIdentifier(identifier uint64, st transfer.StateChange) {
	stateChangeId, _ := this.raiden.TransactionLog.Log(st)
	mgrs := this.raiden.Identifier2StateManagers[identifier]
	for _, mgr := range mgrs {
		events := this.Dispatch(mgr, st)
		this.raiden.TransactionLog.LogEvents(stateChangeId, events, this.raiden.GetBlockNumber())
	}
}

//Log a state change, dispatch it to the given state manager and log generated events
func (this *StateMachineEventHandler) LogAndDispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) []transfer.Event {
	stateChangeId, _ := this.raiden.TransactionLog.Log(stateChange)
	events := this.Dispatch(stateManager, stateChange)
	this.raiden.TransactionLog.LogEvents(stateChangeId, events, this.raiden.GetBlockNumber())
	return events
}
func (this *StateMachineEventHandler) Dispatch(stateManager *transfer.StateManager, stateChange transfer.StateChange) (events []transfer.Event) {
	events = stateManager.Dispatch(stateChange)
	for _, e := range events {
		err := this.OnEvent(e)
		if err != nil {
			log.Error(fmt.Sprintf("StateMachineEventHandler Dispatch:%v\n", err))
		}
	}
	return
}
func (this *StateMachineEventHandler) eventSendMediatedTransfer(event *mediated_transfer.EventSendMediatedTransfer) (err error) {
	receiver := event.Receiver
	graph := this.raiden.GetToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateMediatedTransfer(event.Initiator, event.Target, 0, event.Amount, event.Identifier, event.Expiration, event.HashLock)
	if err != nil {
		return
	}
	mtr.Sign(this.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(this.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	err = this.raiden.SendAsync(receiver, mtr)
	return
}
func (this *StateMachineEventHandler) eventSendRefundTransfer(event *mediated_transfer.EventSendRefundTransfer) (err error) {
	receiver := event.Receiver
	graph := this.raiden.GetToken2ChannelGraph(event.Token)
	ch := graph.GetPartenerAddress2Channel(receiver)
	mtr, err := ch.CreateRefundTransfer(event.Initiator, event.Target, 0, event.Amount, event.Identifier, event.Expiration, event.HashLock)
	if err != nil {
		return
	}
	mtr.Sign(this.raiden.PrivateKey, mtr)
	err = ch.RegisterTransfer(this.raiden.GetBlockNumber(), mtr)
	if err != nil {
		return
	}
	err = this.raiden.SendAsync(receiver, mtr)
	return
}
func (this *StateMachineEventHandler) OnEvent(event transfer.Event) (err error) {
	switch e2 := event.(type) {
	case *mediated_transfer.EventSendMediatedTransfer:
		err = this.eventSendMediatedTransfer(e2)
	case *mediated_transfer.EventSendRevealSecret:
		revealMessage := encoding.NewRevealSecret(e2.Secret)
		revealMessage.Sign(this.raiden.PrivateKey, revealMessage)
		err = this.raiden.SendAsync(e2.Receiver, revealMessage)
	case *mediated_transfer.EventSendBalanceProof:
		//unlock and update remotely (send the Secret message)
		err = this.raiden.HandleSecret(e2.Identifier, e2.Token, e2.Secret, nil, utils.Sha3(e2.Secret[:]))
	case *mediated_transfer.EventSendSecretRequest:
		secretRequest := encoding.NewSecretRequest(e2.Identifer, e2.Hashlock, e2.Amount)
		secretRequest.Sign(this.raiden.PrivateKey, secretRequest)
		err = this.raiden.SendAsync(e2.Receiver, secretRequest)
	case *mediated_transfer.EventSendRefundTransfer:
		err = this.eventSendRefundTransfer(e2)
	case *transfer.EventTransferSentSuccess:
		log.Info(fmt.Sprintf("EventTransferSentSuccess for id %d ", e2.Identifier))
		//may receive multi success because of duplicate messages
		/*
			method 1.
			remove this id info after success
			method 2.
			mark success
		*/
		this.raiden.Lock.Lock()
		for _, r := range this.raiden.Identifier2Results[e2.Identifier] {
			r.Result <- nil
			close(r.Result)
		}
		//todo fix this ,when to delete Identifier2StateManager?
		delete(this.raiden.Identifier2Results, e2.Identifier)
		this.raiden.Lock.Unlock()
	case *transfer.EventTransferSentFailed:
		log.Info(fmt.Sprintf("EventTransferSentFailed for id %d", e2.Identifier))
		this.raiden.Lock.Lock()
		for _, r := range this.raiden.Identifier2Results[e2.Identifier] {
			r.Result <- errSentFailed
			close(r.Result)
		}
		delete(this.raiden.Identifier2Results, e2.Identifier)
		this.raiden.Lock.Unlock()
	case *transfer.EventTransferReceivedSuccess:
	case *mediated_transfer.EventUnlockSuccess:
	case *mediated_transfer.EventWithdrawFailed:
	case *mediated_transfer.EventWithdrawSuccess:
		/*
					  The withdraw is currently handled by the netting channel, once the close
			     event is detected all locks will be withdrawn
		*/
	case *mediated_transfer.EventContractSendWithdraw:
		//do nothing for five events above
	case *mediated_transfer.EventUnlockFailed:
		log.Error(fmt.Sprintf("unlockfailed hashlock=%s,reason=%s", e2.Hashlock, e2.Reason))
	case *mediated_transfer.EventContractSendChannelClose:
		graph := this.raiden.GetToken2ChannelGraph(e2.Token)
		ch := graph.GetPartenerAddress2Channel(e2.ChannelAddress)
		balanceProof := ch.OurState.BalanceProofState
		err = ch.ExternState.Close(balanceProof)
	default:
		err = fmt.Errorf("unkown event :%s", utils.StringInterface1(event))
		log.Error(err.Error())
	}
	return
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
	ch.StateTransition(st)
	if ch.ContractBalance() == 0 {
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
	ch.StateTransition(st)
	return nil
}

func (this *StateMachineEventHandler) handleSettled(st *mediated_transfer.ContractReceiveSettledStateChange) error {
	//todo remove channel st.channelAddress ,because this channel is already settled
	ch, err := this.raiden.FindChannelByAddress(st.ChannelAddress)
	if err != nil {
		return err
	}
	ch.StateTransition(st)
	return nil
}
func (this *StateMachineEventHandler) handleWithdraw(st *mediated_transfer.ContractReceiveWithdrawStateChange) error {
	this.raiden.RegisterSecret(st.Secret)
	return nil
}

//only care statechanges about me
func (this *StateMachineEventHandler) filterStateChange(st transfer.StateChange) bool {
	var channelAddress common.Address
	switch st2 := st.(type) { //filter event only about me
	case *mediated_transfer.ContractReceiveTokenAddedStateChange:
		return true
	case *mediated_transfer.ContractReceiveNewChannelStateChange:
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
	this.raiden.Lock.RLock()
	for _, g := range this.raiden.Token2ChannelGraph {
		ch := g.GetChannelAddress2Channel(channelAddress)
		if ch != nil {
			found = true
			break
		}
	}
	this.raiden.Lock.RUnlock()
	return found
}
func (this *StateMachineEventHandler) OnBlockchainStateChange(st transfer.StateChange) (err error) {
	log.Info("statechange received :", utils.StringInterface1(st))
	_, err = this.raiden.TransactionLog.Log(st)
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
