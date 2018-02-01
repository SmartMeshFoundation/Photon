package raiden_network

import (
	"encoding/gob"

	"fmt"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer/initiator"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer/mediator"
	"github.com/SmartMeshFoundation/raiden-network/transfer/mediated_transfer/target"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gotips/log"
)

func init() {
	gob.Register(&Data2Save{})
}

type Data2Save struct {
	Channels        []*channel.ChannelSerialization
	Transfers       map[uint64][]*transfer.StateManager
	RegistryAddress common.Address
}

//save state ,call many times is ok
func (this *RaidenService) SaveSnapshot() {
	if this.Config.Debug {
		return
	}
	log.Info("SaveSnapshot...")
	ds := &Data2Save{
		RegistryAddress: this.RegistryAddress,
		Transfers:       this.Identifier2StateManagers,
	}
	for _, g := range this.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			cs := channel.NewChannelSerialization(c)
			ds.Channels = append(ds.Channels, cs)
		}
	}
	_, err := this.db.Snapshot(1, ds)
	if err != nil {
		log.Error("save snapshot :", err)
	}
}

//retore state ,only one time ,just after app start immediately
func (this *RaidenService) RestoreSnapshot() error {
	if this.Config.Debug {
		return nil
	}
	log.Info("RestoreSnapshot...")
	data, err := this.db.LoadSnapshot()
	if err != nil {
		log.Error("restore snapshot :", err)
		return err
	}
	//never save before
	if data == nil {
		return nil
	}
	ds := data.(*Data2Save)
	if ds.RegistryAddress != this.RegistryAddress {
		err := fmt.Errorf("snapshot data error, registry address error")
		log.Error(err.Error())
		return err
	}
	this.restoreChannel(ds)
	this.restoreStateManager(ds)
	return nil
}
func (this *RaidenService) restoreChannel(ds *Data2Save) error {
	for _, g := range this.Token2ChannelGraph {
		for _, c := range g.ChannelAddress2Channel {
			for _, cs := range ds.Channels {
				//found a channel,maybe channel settled or new channel opened when i'm down
				if cs.ChannelAddress == c.MyAddress {
					if cs.TokenAddress != c.TokenAddress || cs.OurAddress != c.OurState.Address ||
						cs.PartnerAddress != c.PartnerState.Address ||
						c.RevealTimeout != cs.RevealTimeout {
						log.Error("snapshot data error, channel data error for ", c.MyAddress)
						continue
					} else {
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
	}
	return nil
}

//function pointer save and restore
func (this *RaidenService) restoreStateManager(ds *Data2Save) error {
	this.Identifier2StateManagers = ds.Transfers
	for _, ss := range this.Identifier2StateManagers {
		for _, s := range ss {
			switch s.Name {
			case initiator.NameInitiatorTransition:
				s.FuncStateTransition = initiator.StateTransition
			case mediator.NameMediatorTransition:
				s.FuncStateTransition = mediator.StateTransition
			case target.NameTargetTransition:
				s.FuncStateTransition = target.StateTransiton
			default:
				log.Error("unkown state manager :", s.Name)
			}
		}
	}
	return nil
}
