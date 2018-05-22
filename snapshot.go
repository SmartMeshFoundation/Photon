package smartraiden

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
)

//save state ,call many times is ok
func (rs *RaidenService) saveSnapshot() {
	log.Info("saveSnapshot...")
	rs.db.IsDbCrashedLastTime()
}

//retore state ,only one time ,just after app start immediately
func (rs *RaidenService) restoreSnapshot() error {
	log.Info("restoreSnapshot...")
	defer func() {
		rs.db.MarkDbOpenedStatus()
		rs.db.SaveRegistryAddress(rs.RegistryAddress)
	}()
	/*
	   When debugging, the registry address may change constantly, Testing is to avoid unnecessary mistakes.
	*/
	registryAddr := rs.db.GetRegistryAddress()
	if registryAddr != rs.RegistryAddress && registryAddr != utils.EmptyAddress {
		err := fmt.Errorf("db registry address not match db=%s,mine=%s", registryAddr.String(), rs.RegistryAddress.String())
		log.Error(err.Error())
		return err
	}
	//never save before
	/*
		The first step    restore the channel state
		The second step  restore the hashlock in channel, so that hashlock changes during subsequent recovery, it can be reflected in the corresponding channel.
		The third step   restore stateManager. This step will send unfinished messages, which may change the channel state.
		The fourth step   recovery processing for unsent revealsecret and the unprocessed revealsecret.
		The fifth step   save the recovered channel state to the database.
	*/
	rs.restoreChannel(rs.db.IsDbCrashedLastTime())
	rs.restoreToken2Hash2Channels()
	rs.restoreStateManager(rs.db.IsDbCrashedLastTime())
	rs.restoreRevealSecret()
	rs.saveChannelStatusAfterStart()
	return nil
}
func (rs *RaidenService) saveChannelStatusAfterStart() {
	for _, g := range rs.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			rs.db.UpdateChannelNoTx(channel.NewChannelSerialization(c))
		}
	}
}

//ds, parameter for validate stored data.
func (rs *RaidenService) restoreChannel(isCrashed bool) error {
	log.Info("restore channel...")
	for _, g := range rs.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			cs, err := rs.db.GetChannelByAddress(c.MyAddress)
			if err != nil {
				if err == storm.ErrNotFound {
					continue //new channel when shutdown
				} else {
					panic(fmt.Sprintf("get channel %s from db err %s", c.MyAddress.String(), err))
				}
			}
			//found a channel,maybe channel settled or new channel opened when i'm down
			if cs.ChannelAddress == c.MyAddress {
				if cs.TokenAddress != c.TokenAddress || cs.OurAddress != c.OurState.Address ||
					cs.PartnerAddress != c.PartnerState.Address ||
					c.RevealTimeout != cs.RevealTimeout {
					log.Error("snapshot data error, channel data error for ", c.MyAddress)
					continue
				} else {
					log.Trace(fmt.Sprintf("retore channel %s\n", utils.StringInterface(cs, 7)))
					c.OurState.BalanceProofState = cs.OurBalanceProof
					c.OurState.TreeState = transfer.NewMerkleTreeStateFromLeaves(cs.OurLeaves)
					c.OurState.Lock2PendingLocks = cs.OurLock2PendingLocks
					c.OurState.Lock2UnclaimedLocks = cs.OurLock2UnclaimedLocks
					c.PartnerState.BalanceProofState = cs.PartnerBalanceProof
					c.PartnerState.TreeState = transfer.NewMerkleTreeStateFromLeaves(cs.PartnerLeaves)
					c.PartnerState.Lock2PendingLocks = cs.PartnerLock2PendingLocks
					c.PartnerState.Lock2UnclaimedLocks = cs.PartnerLock2UnclaimedLocks
				}
			}
		}
	}
	return nil
}

func (rs *RaidenService) restoreDbPointer(state transfer.State) {
	if state == nil {
		return
	}
	switch st2 := state.(type) {
	case *mediatedtransfer.InitiatorState:
		st2.Db = rs.db
	case *mediatedtransfer.TargetState:
		st2.Db = rs.db
	case *mediatedtransfer.MediatorState:
		st2.Db = rs.db
	default:
		panic(fmt.Sprintf("unkown state %s", utils.StringInterface(st2, 3)))
	}
}

//function pointer save and restore
func (rs *RaidenService) restoreStateManager(isCrashed bool) {
	log.Info(fmt.Sprintf("restore statemanager ,last close correct=%v", !isCrashed))
	mgrs := rs.db.GetAllStateManager()
	for _, mgr := range mgrs {
		//log.Trace(fmt.Sprintf("unfinish manager %s", utils.StringInterface(mgr, 7)))
		if mgr.ManagerState == transfer.StateManagerStateInit || mgr.ManagerState == transfer.StateManagerTransferComplete {
			continue
		}
		setStateManagerFuncPointer(mgr)
		idmgrs := rs.Identifier2StateManagers[mgr.Identifier]
		idmgrs = append(idmgrs, mgr)
		rs.Identifier2StateManagers[mgr.Identifier] = idmgrs
		rs.restoreDbPointer(mgr.CurrentState)
	}
	for _, mgrs := range rs.Identifier2StateManagers {
		//mannagers for the same channel should be order, otherwise, nonce error.
		for _, mgr := range mgrs {
			log.Trace(fmt.Sprintf("restore state manager:%s\n", utils.StringInterface(mgr, 7)))
			var tag interface{}
			var messageTag *transfer.MessageTag
			switch mgr.ManagerState {
			case transfer.StateManagerTransferComplete:
				//ignore
			case transfer.StateManagerReceivedMessage:
				st, ok := mgr.LastReceivedMessage.(mediatedtransfer.ActionInitInitiatorStateChange)
				if ok {
					st.Db = rs.db
					rs.StateMachineEventHandler.dispatch(mgr, st)
				} else {
					//receive a message,and not handled
					//ignore ,partner will try
				}
				if mgr.LastSendMessage == nil {
					break
				}
				fallthrough
			case transfer.StateManagerReceivedMessageProcessComplete: //there may be message waiting for send
				if mgr.LastSendMessage == nil {
					break
				}
				fallthrough
			case transfer.StateManagerSendMessage:
				/*
					todo fix It should be detected whether it is out of date,
					such as ,MediatedTransfer, Secret, which are timeliness, and if it is expired after crash,discarding is more reasonable.
				*/
				tag = mgr.LastSendMessage.Tag()
				if tag == nil {
					panic(fmt.Sprintf("statemanage state error, lastsendmessage has no tag :%s", utils.StringInterface(mgr, 5)))
				}
				messageTag = tag.(*transfer.MessageTag)
				if messageTag.SendingMessageComplete {
					continue // for receive secret message, no need sending any message but ack.
				}
				messageTag.SetStateManager(mgr) //statemanager doesn't save
				rs.sendAsync(messageTag.Receiver, mgr.LastSendMessage.(encoding.SignedMessager))
			case transfer.StateManagerSendMessageSuccesss:
				//do nothing right now.
			}
		}
	}

}
func setStateManagerFuncPointer(mgr *transfer.StateManager) {
	/*
		todo fix tokenswap's randomSecretGenerator
	*/
	switch mgr.Name {
	case initiator.NameInitiatorTransition:
		mgr.FuncStateTransition = initiator.StateTransition
		if mgr.CurrentState != nil {
			state := mgr.CurrentState.(*mediatedtransfer.InitiatorState)
			state.RandomGenerator = utils.RandomSecretGenerator //todo fix for tokenswap
		}
	case mediator.NameMediatorTransition:
		mgr.FuncStateTransition = mediator.StateTransition
	case target.NameTargetTransition:
		mgr.FuncStateTransition = target.StateTransiton
	default:
		log.Error("unkown state manager :", mgr.Name)
	}
}

func (rs *RaidenService) restoreRevealSecret() {
	log.Trace(fmt.Sprintf("restoreRevealSecret... "))
	receiveSecrets := rs.db.GetAllUncompleteReceivedRevealSecret()
	for _, s := range receiveSecrets {
		rs.MessageHandler.onMessage(s.Message, s.EchoHash)
	}
	sentSecrets := rs.db.GetAllUncompleteSentRevealSecret()
	for _, s := range sentSecrets {
		rs.sendAsync(s.Receiver, s.Message)
	}
}

func (rs *RaidenService) restoreToken2Hash2Channels() {
	log.Trace("restoreToken2Hash2Channels...")
	for token, g := range rs.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			for lock := range c.OurState.Lock2PendingLocks {
				rs.registerChannelForHashlock(token, c, lock)
			}
			for lock := range c.PartnerState.Lock2PendingLocks {
				rs.registerChannelForHashlock(token, c, lock)
			}
			for lock := range c.OurState.Lock2UnclaimedLocks {
				rs.registerChannelForHashlock(token, c, lock)
			}
			for lock := range c.PartnerState.Lock2UnclaimedLocks {
				rs.registerChannelForHashlock(token, c, lock)
			}
		}
	}
}
