package route

import (
	"math/big"

	"encoding/gob"

	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/ethereum/go-ethereum/common"
)

/*
State describes a route state
路由状态我如何收到的或者发送MediatedTransfer
*/
/*
 *	State : describe a route state
 *	route state means how does this participant get or send MediatedTransfer.
 */
type State struct {
	ch                *channel.Channel //don't save pointer
	ChannelIdentifier common.Hash      //崩溃恢复的时候需要	// Needed when handling crash cases.
	IsSend            bool             //用这个 route 来发送还是接收?	// whether this route is used to send or receive.
	Fee               *big.Int         // how much fee to this channel charge charge .
	TotalFee          *big.Int         // how much fee for all path when initiator use this route
	Path              []common.Address // 2019-03消息升级,路由中保存该条路径上所有节点,有序
}

//NewState create route state
func NewState(ch *channel.Channel, path []common.Address) *State {
	return &State{
		ChannelIdentifier: ch.ChannelIdentifier.ChannelIdentifier,
		ch:                ch,
		Path:              path,
	}
}

//CanTransfer can transfer on this hop node
func (rs *State) CanTransfer() bool {
	return rs.ch.CanTransfer()
}

//CanContinueTransfer can continue on this hop node
func (rs *State) CanContinueTransfer() bool {
	return rs.ch.CanContinueTransfer()
}

//SettleTimeout settle timeout of this channel
func (rs *State) SettleTimeout() int {
	return rs.ch.SettleTimeout
}

//RevealTimeout reveal timeout of this channel
func (rs *State) RevealTimeout() int {
	return rs.ch.RevealTimeout
}

//SetClosedBlock set closed block ,for test only
func (rs *State) SetClosedBlock(blockNumbder int64) {
	rs.ch.ExternState.ClosedBlock = blockNumbder
}

//ClosedBlock return closedBlock of this route channel
func (rs *State) ClosedBlock() int64 {
	return rs.ch.ExternState.ClosedBlock
}

//HopNode hop node
func (rs *State) HopNode() common.Address {
	return rs.ch.PartnerState.Address
}

//AvailableBalance avaialabe balance of this route
func (rs *State) AvailableBalance() *big.Int {
	return rs.ch.Distributable()
}

//Channel return Channel
func (rs *State) Channel() *channel.Channel {
	return rs.ch
}

//State of route channel
func (rs *State) State() channeltype.State {
	return rs.ch.State
}

//SetState for test only,
func (rs *State) SetState(state channeltype.State) {
	rs.ch.State = state
}

//StateName return name of the state
func (rs *State) StateName() string {
	return "State"
}

/*
RoutesState is Routing state.
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
	//m := make(map[common.Address]bool)
	for _, r := range availables {
		// 2019-03 消息升级后可能存在多条路径的开头部分节点重复,这里取消校验
		//_, ok := m[r.HopNode()]
		//if ok {
		//	log.Warn("duplicate route for the same address supplied.")
		//	continue
		//}
		rs.AvailableRoutes = append(rs.AvailableRoutes, r)
	}
	return rs
}
func init() {
	gob.Register(&State{})
	gob.Register(&RoutesState{})
}
