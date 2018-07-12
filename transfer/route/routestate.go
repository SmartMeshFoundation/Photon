package route

import (
	"math/big"

	"encoding/gob"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
State describes a route state
路由状态我如何收到的或者发送MediatedTransfer
*/
type State struct {
	ch                *channel.Channel
	ChannelIdentifier common.Hash //崩溃恢复的时候需要
	IsSend            bool        //用这个 route 来发送还是接收?
	Fee               *big.Int    // how much fee to this channel charge charge .
	TotalFee          *big.Int    // how much fee for all path when initiator use this route
}

func (rs *State) CanTransfer() bool {
	return rs.ch.CanTransfer()
}
func (rs *State) CanContinueTransfer() bool {
	return rs.ch.CanContinueTransfer()
}
func (rs *State) SettleTimeout() int {
	return rs.ch.SettleTimeout
}
func (rs *State) RevealTimeout() int {
	return rs.ch.RevealTimeout
}
func (rs *State) ClosedBlock() int64 {
	return rs.ch.ExternState.ClosedBlock
}
func (rs *State) HopNode() common.Address {
	if rs.IsSend {
		return rs.ch.OurState.Address
	}
	return rs.ch.PartnerState.Address
}
func (rs *State) AvailableBalance() *big.Int {
	return rs.ch.Distributable()
}
func (rs *State) Channel() *channel.Channel {
	return rs.ch
}
func (rs *State) State() channeltype.State {
	return rs.ch.State
}

//StateName return name of the state
func (rs *State) StateName() string {
	return "State"
}

/*
RoutesState is Routing state.
    Args:
        available_routes (list): A list of RouteState instances.
*/
type RoutesState struct {
	AvailableRoutes []*State
	IgnoredRoutes   []*State
	RefundedRoutes  []*State
	CanceledRoutes  []*State
}

//NewRoutesState create routes state from availabes routes
func NewRoutesState(availables []*State) *RoutesState {
	rs := &RoutesState{}
	m := make(map[common.Address]bool)
	for _, r := range availables {
		_, ok := m[r.HopNode()]
		if ok {
			log.Warn("duplicate route for the same address supplied.")
			continue
		}
		rs.AvailableRoutes = append(rs.AvailableRoutes, r)
	}
	return rs
}
func init() {
	gob.Register(&State{})
	gob.Register(&RoutesState{})
}
