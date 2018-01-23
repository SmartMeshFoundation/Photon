package transfer

import (
	"fmt"

	"encoding/gob"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const CHANNEL_STATE_CLOSED = "closed"
const CHANNEL_STATE_CLOSING = "waiting_for_close"
const CHANNEL_STATE_OPENED = "opened"
const CHANNEL_STATE_SETTLED = "settled"
const CHANNEL_STATE_SETTLING = "waiting_for_settle"

//nothing
//type Stater interface {
//	fmt.Stringer
//	StateName() string
//}

//this describes a route state
type RouteState struct {
	State          string
	HopNode        common.Address
	ChannelAddress common.Address
	AvaibleBalance int64
	SettleTimeout  int
	RevealTimeout  int
	ClosedBlock    int64
}

/*
Args:
        state (string): The current state of the route (opened, closed or settled).
        node_address (address): The address of the next_hop.
        channel_address (address): The address of the on chain netting channel.
        available_balance (int): The current available balance that can be transferred
            through `node_address`.
        settle_timeout (int): The settle_timeout of the channel set in the
            smart contract.
        reveal_timeout (int): The channel configured reveal_timeout.
        closed_block (Nullable[int]): None if the channel is open, otherwise
            the block number at which the channel was closed.
*/
func NewRouteState(state string, nodeAddress common.Address, channelAddress common.Address,
	avaibleBalance int64, settleTimeout int, revealTimeout int, closedBlock int64) *RouteState {
	s := &RouteState{
		State:          state,
		HopNode:        nodeAddress, //hop
		ChannelAddress: channelAddress,
		AvaibleBalance: avaibleBalance,
		SettleTimeout:  settleTimeout,
		RevealTimeout:  revealTimeout,
		ClosedBlock:    closedBlock,
	}
	return s
}

func (this *RouteState) StateName() string {
	return "RouteState"
}

//func (this *RouteState) String() string {
//	return utils.StringInterface(this, 1)
//}

type BalanceProofState struct {
	Nonce          int64
	TransferAmount int64
	LocksRoot      common.Hash
	ChannelAddress common.Address
	MessageHash    common.Hash
	//signature is nonce + transferred_amount + locksroot + channel_address + message_hash
	Signature []byte
}

func NewBalanceProofState(nonce int64, transferAmount int64, locksRoot common.Hash,
	channelAddress common.Address, messageHash common.Hash, signature []byte) *BalanceProofState {
	s := &BalanceProofState{
		Nonce:          nonce,
		TransferAmount: transferAmount,
		LocksRoot:      locksRoot,
		ChannelAddress: channelAddress,
		MessageHash:    messageHash,
		Signature:      signature,
	}
	return s
}
func NewBalanceProofStateFromEnvelopMessage(msg encoding.EnvelopMessager) *BalanceProofState {
	envmsg := msg.GetEnvelopMessage()
	msgHash := encoding.HashMessageWithoutSignature(msg)
	return NewBalanceProofState(envmsg.Nonce, envmsg.TransferAmount.Int64(),
		envmsg.Locksroot, envmsg.Channel,
		msgHash, envmsg.Signature)
}

//func (this *BalanceProofState) String() string {
//	return utils.StringInterface(this, 1)
//}
func (this *BalanceProofState) StateName() string {
	return "BalanceProofState"
}

type MerkleTreeState struct {
	Tree *Merkletree
}

var EmptyMerkleTreeState *MerkleTreeState

func init() {
	tree, _ := NewMerkleTree(nil)
	EmptyMerkleTreeState = NewMerkleTreeState(tree)
}

func NewMerkleTreeState(tree *Merkletree) *MerkleTreeState {
	return &MerkleTreeState{
		tree,
	}
}
func NewMerkleTreeStateFromLeaves(leaves []common.Hash) *MerkleTreeState {
	tree, _ := NewMerkleTree(leaves)
	return &MerkleTreeState{
		tree,
	}
}

func (this *MerkleTreeState) StateName() string {
	return "MerkleTreeState"
}
func (this *MerkleTreeState) String() string {
	return fmt.Sprintf("MerkleTreeState{root:%s,layer level:%d}", this.Tree.MerkleRoot(), len(this.Tree.Layers))
}

/*
Routing state.

    Args:
        available_routes (list): A list of RouteState instances.
*/
type RoutesState struct {
	AvailableRoutes []*RouteState
	IgnoredRoutes   []*RouteState
	RefundedRoutes  []*RouteState
	CanceledRoutes  []*RouteState
}

func NewRoutesState(availables []*RouteState) *RoutesState {
	rs := &RoutesState{}
	m := make(map[common.Address]bool)
	for _, r := range availables {
		_, ok := m[r.HopNode]
		if ok {
			log.Warn("duplicate route for the same address supplied.")
			continue
		}
		rs.AvailableRoutes = append(rs.AvailableRoutes, r)
	}
	return rs
}
func init() {
	gob.Register(&RouteState{})
	gob.Register(&RoutesState{})
	gob.Register(&MerkleTreeState{})
	gob.Register(&BalanceProofState{})
}
