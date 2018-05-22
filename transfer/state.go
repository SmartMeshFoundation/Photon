package transfer

import (
	"fmt"

	"encoding/gob"

	"math/big"

	"bytes"
	"encoding/binary"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//ChannelStateClosed closed
const ChannelStateClosed = "closed"

//ChannelStateClosing waiting close
const ChannelStateClosing = "waiting_for_close"

//ChannelStateOpened opened
const ChannelStateOpened = "opened"

//ChannelStateSettled settled
const ChannelStateSettled = "settled"

//ChannelStateSetting waiting settle
const ChannelStateSetting = "waiting_for_settle"

//RouteState describes a route state
type RouteState struct {
	State          string
	HopNode        common.Address
	ChannelAddress common.Address
	AvaibleBalance *big.Int
	Fee            *big.Int // how much fee to this channel charge charge .
	TotalFee       *big.Int // how much fee for all path when initiator use this route
	SettleTimeout  int
	RevealTimeout  int
	ClosedBlock    int64
}

/*
NewRouteState create route state
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
	avaibleBalance *big.Int, settleTimeout int, revealTimeout int, closedBlock int64) *RouteState {
	s := &RouteState{
		State:          state,
		HopNode:        nodeAddress, //hop
		ChannelAddress: channelAddress,
		AvaibleBalance: new(big.Int).Set(avaibleBalance),
		SettleTimeout:  settleTimeout,
		RevealTimeout:  revealTimeout,
		ClosedBlock:    closedBlock,
	}
	return s
}

//StateName return name of the state
func (rs *RouteState) StateName() string {
	return "RouteState"
}

//BalanceProofState is   proof need by contract
type BalanceProofState struct {
	Nonce          int64
	TransferAmount *big.Int
	LocksRoot      common.Hash
	ChannelAddress common.Address
	MessageHash    common.Hash
	//signature is nonce + transferred_amount + locksroot + channel_address + message_hash
	Signature []byte
}

//NewBalanceProofState create BalanceProofState
func NewBalanceProofState(nonce int64, transferAmount *big.Int, locksRoot common.Hash,
	channelAddress common.Address, messageHash common.Hash, signature []byte) *BalanceProofState {
	s := &BalanceProofState{
		Nonce:          nonce,
		TransferAmount: new(big.Int).Set(transferAmount),
		LocksRoot:      locksRoot,
		ChannelAddress: channelAddress,
		MessageHash:    messageHash,
		Signature:      signature,
	}
	return s
}

//NewBalanceProofStateFromEnvelopMessage from locked transfer
func NewBalanceProofStateFromEnvelopMessage(msg encoding.EnvelopMessager) *BalanceProofState {
	envmsg := msg.GetEnvelopMessage()
	msgHash := encoding.HashMessageWithoutSignature(msg)
	return NewBalanceProofState(envmsg.Nonce, envmsg.TransferAmount,
		envmsg.Locksroot, envmsg.Channel,
		msgHash, envmsg.Signature)
}

//IsBalanceProofValid true if valid
func (bpf *BalanceProofState) IsBalanceProofValid() bool {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, bpf.Nonce)
	buf.Write(utils.BigIntTo32Bytes(bpf.TransferAmount))
	buf.Write(bpf.LocksRoot[:])
	buf.Write(bpf.ChannelAddress[:])
	buf.Write(bpf.MessageHash[:])
	dataToSign := buf.Bytes()

	hash := utils.Sha3(dataToSign)
	signature := make([]byte, len(bpf.Signature))
	copy(signature, bpf.Signature)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	//log.Trace(fmt.Sprintf("signer =%s",utils.APex(utils.PubkeyToAddress(pubkey))))
	return err == nil && utils.PubkeyToAddress(pubkey) != utils.EmptyAddress
}

//StateName name of state
func (bpf *BalanceProofState) StateName() string {
	return "BalanceProofState"
}

//MerkleTreeState need by channel
type MerkleTreeState struct {
	Tree *Merkletree
}

//EmptyMerkleTreeState a empty tree
var EmptyMerkleTreeState *MerkleTreeState

func init() {
	tree, _ := NewMerkleTree(nil)
	EmptyMerkleTreeState = NewMerkleTreeState(tree)
}

//NewMerkleTreeState  create MerkleTreeState from MerkleTree
func NewMerkleTreeState(tree *Merkletree) *MerkleTreeState {
	return &MerkleTreeState{
		tree,
	}
}

//NewMerkleTreeStateFromLeaves create MerkleTreeState from leaves
func NewMerkleTreeStateFromLeaves(leaves []common.Hash) *MerkleTreeState {
	tree, _ := NewMerkleTree(leaves)
	return &MerkleTreeState{
		tree,
	}
}

//StateName state  name of MerkleTreeState
func (mt *MerkleTreeState) StateName() string {
	return "MerkleTreeState"
}
func (mt *MerkleTreeState) String() string {
	return fmt.Sprintf("MerkleTreeState{root:%s,layer level:%d}", mt.Tree.MerkleRoot(), len(mt.Tree.Layers))
}

/*
RoutesState is Routing state.
    Args:
        available_routes (list): A list of RouteState instances.
*/
type RoutesState struct {
	AvailableRoutes []*RouteState
	IgnoredRoutes   []*RouteState
	RefundedRoutes  []*RouteState
	CanceledRoutes  []*RouteState
}

//NewRoutesState create routes state from availabes routes
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
