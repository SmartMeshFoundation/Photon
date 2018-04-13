package smartraiden

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/mediator"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediated_transfer/target"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/log"
	"errors"
)

//save state ,call many times is ok
func (this *RaidenService) SaveSnapshot() {
	log.Info("SaveSnapshot...")
	this.db.IsDbCrashedLastTime()
}

//retore state ,only one time ,just after app start immediately
func (this *RaidenService) RestoreSnapshot() error {
	log.Info("RestoreSnapshot...")
	defer func() {
		this.db.MarkDbOpenedStatus()
		this.db.SaveRegistryAddress(this.RegistryAddress)
	}()
/*
在调试时, registry 地址可能会不断变化,检测一下,以免引起不必要的错误
 */
 registryAddr:=this.db.GetRegistryAddress()
 if registryAddr!=this.RegistryAddress && registryAddr!=utils.EmptyAddress{
 	err:=errors.New(fmt.Sprintf("db registry address not match db=%s,mine=%s",registryAddr.String(),this.RegistryAddress.String()))
 	log.Error(err.Error())
 	return err
 }
	//never save before
	/*
		第一步 恢复channel状态
		第二步  将channel中的hashlock恢复,这样后续恢复过程中,hashlock发生变化的时候可以反映到对应的channel中
		第三步 恢复stateManager,此步恢复以后将会发送未完成的消息,可能会改变channel状态
		第四步 将未完成发送的revealsecret和未处理完的收到的revealsecret 恢复处理,由于revealsecret可以重复收发,所以重复应该没有副作用.
		第五步 将恢复后的channel状态保存到数据库中, 此步骤似乎可以忽略,(在启动过程中,确保没有其他地方会修改channel状态)
	*/
	this.restoreChannel(this.db.IsDbCrashedLastTime())
	this.RestoreToken2Hash2Channels()
	this.restoreStateManager(this.db.IsDbCrashedLastTime())
	this.RestoreRevealSecret()
	this.saveChannelStatusAfterStart()
	return nil
}
func (this *RaidenService) saveChannelStatusAfterStart() {
	for _, g := range this.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			this.db.UpdateChannelNoTx(channel.NewChannelSerialization(c))
		}
	}
}

//ds, parameter for validate stored data.
func (this *RaidenService) restoreChannel(isCrashed bool) error {
	log.Info("restore channel...")
	for _, g := range this.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			cs, err := this.db.GetChannelByAddress(c.MyAddress)
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

func (this*RaidenService) restoreDbPointer(state transfer.State) {
	if state==nil{
		return
	}
	switch st2:=state.(type){
	case *mediated_transfer.InitiatorState:
		st2.Db=this.db
	case *mediated_transfer.TargetState:
		st2.Db=this.db
	case *mediated_transfer.MediatorState:
		st2.Db=this.db
	default:
		panic(fmt.Sprintf("unkown state %s",utils.StringInterface(st2,3)))
	}
}
//function pointer save and restore
func (this *RaidenService) restoreStateManager(isCrashed bool) {
	log.Info(fmt.Sprintf("restore statemanager ,last close correct=%s", !isCrashed))
	mgrs := this.db.GetAllStateManager()
	for _, mgr := range mgrs {
		//log.Trace(fmt.Sprintf("unfinish manager %s", utils.StringInterface(mgr, 7)))
		if mgr.ManagerState == transfer.StateManager_State_Init || mgr.ManagerState == transfer.StateManager_TransferComplete {
			continue
		}
		setStateManagerFuncPointer(mgr)
		idmgrs := this.Identifier2StateManagers[mgr.Identifier]
		idmgrs = append(idmgrs, mgr)
		this.Identifier2StateManagers[mgr.Identifier] = idmgrs
		this.restoreDbPointer(mgr.CurrentState)
	}
	for _, mgrs := range this.Identifier2StateManagers {
		//mannagers for the same channel should be order, otherwise, nonce error.
		for _, mgr := range mgrs {
			log.Trace(fmt.Sprintf("restore state manager:%s\n", utils.StringInterface(mgr, 7)))
			var tag interface{}
			var messageTag *transfer.MessageTag
			switch mgr.ManagerState {
			case transfer.StateManager_TransferComplete:
				//ignore
			case transfer.StateManager_ReceivedMessage:
				st, ok := mgr.LastReceivedMessage.(mediated_transfer.ActionInitInitiatorStateChange)
				if ok {
					st.Db = this.db
					this.StateMachineEventHandler.Dispatch(mgr, st)
				} else {
					//receive a message,and not handled
					//ignore ,partner will try
				}
				if mgr.LastSendMessage == nil {
					break
				}
				fallthrough
			case transfer.StateManager_ReceivedMessageProcessComplete: //there may be message waiting for send
				if mgr.LastSendMessage == nil {
					break
				}
				fallthrough
			case transfer.StateManager_SendMessage:
				/*
				todo fix 应该检测一下是否已经过时了,比如 MediatedTransfer,Secret 这些消息是有时效性的,如果崩溃恢复以后已经过期了,直接丢弃更合理.
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
				this.SendAsync(messageTag.Receiver, mgr.LastSendMessage.(encoding.SignedMessager))
			case transfer.StateManager_SendMessageSuccesss:
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
			state := mgr.CurrentState.(*mediated_transfer.InitiatorState)
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

func (this *RaidenService) RestoreRevealSecret() {
	log.Trace(fmt.Sprintf("RestoreRevealSecret... "))
	receiveSecrets := this.db.GetAllUncompleteReceivedRevealSecret()
	for _, s := range receiveSecrets {
		this.MessageHandler.OnMessage(s.Message, s.EchoHash)
	}
	sentSecrets := this.db.GetAllUncompleteSentRevealSecret()
	for _, s := range sentSecrets {
		this.SendAsync(s.Receiver, s.Message)
	}
}

func (this *RaidenService) RestoreToken2Hash2Channels() {
	log.Trace("RestoreToken2Hash2Channels...")
	for token, g := range this.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			for lock, _ := range c.OurState.Lock2PendingLocks {
				this.RegisterChannelForHashlock(token, c, lock)
			}
			for lock, _ := range c.PartnerState.Lock2PendingLocks {
				this.RegisterChannelForHashlock(token, c, lock)
			}
			for lock, _ := range c.OurState.Lock2UnclaimedLocks {
				this.RegisterChannelForHashlock(token, c, lock)
			}
			for lock, _ := range c.PartnerState.Lock2UnclaimedLocks {
				this.RegisterChannelForHashlock(token, c, lock)
			}
		}
	}
}
