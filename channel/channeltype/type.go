package channeltype

import (
	"encoding/gob"
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
PendingLock is lock of HTLC
*/
type PendingLock struct {
	Lock       *mtree.Lock
	LockHashed common.Hash
}

/*
UnlockPartialProof is the lock that I have known the secret ,but haven't receive the Balance proof
*/
type UnlockPartialProof struct {
	Lock       *mtree.Lock
	LockHashed common.Hash
	Secret     common.Hash
}

/*
UnlockProof is the info needs withdraw on blockchain
*/
type UnlockProof struct {
	MerkleProof []common.Hash
	Lock        *mtree.Lock
}

// Serialization is the living channel in the database
type Serialization struct {
	ChannelIdentifier      *contracts.ChannelUniqueID
	Key                    []byte `storm:"id"`
	TokenAddressBytes      []byte `storm:"index"`
	PartnerAddressBytes    []byte `storm:"index"`
	OurAddress             common.Address
	RevealTimeout          int
	OurBalanceProof        *transfer.BalanceProofState
	PartnerBalanceProof    *transfer.BalanceProofState
	OurLeaves              []*mtree.Lock
	PartnerLeaves          []*mtree.Lock
	OurKnownSecrets        []common.Hash
	PartnerKnownSecrets    []common.Hash
	State                  State
	OurContractBalance     *big.Int
	PartnerContractBalance *big.Int
	ClosedBlock            int64
	SettledBlock           int64
	SettleTimeout          int
}

func (s *Serialization) ChannleAddress() common.Hash {
	return common.BytesToHash(s.Key)
}
func (s *Serialization) TokenAddress() common.Address {
	return common.BytesToAddress(s.TokenAddressBytes)
}
func (s *Serialization) PartnerAddress() common.Address {
	return common.BytesToAddress(s.PartnerAddressBytes)
}
func (s *Serialization) transferAmount(bp *transfer.BalanceProofState) *big.Int {
	if bp != nil {
		return bp.TransferAmount
	}
	return utils.BigInt0
}
func (s *Serialization) OurBalance() *big.Int {
	x := new(big.Int)
	x.Sub(s.OurContractBalance, s.transferAmount(s.OurBalanceProof))
	x.Add(x, s.transferAmount(s.PartnerBalanceProof))
	return x
}
func (s *Serialization) OurAmountLocked() *big.Int {
	x := new(big.Int)
	for _, l := range s.OurLeaves {
		x = x.Add(x, l.Amount)
	}
	return x
}
func (s *Serialization) PartnerAmountLocked() *big.Int {
	x := new(big.Int)
	for _, l := range s.PartnerLeaves {
		x = x.Add(x, l.Amount)
	}
	return x
}
func (s *Serialization) PartnerBalance() *big.Int {
	x := new(big.Int)
	x.Sub(s.PartnerContractBalance, s.transferAmount(s.PartnerBalanceProof))
	x.Add(x, s.transferAmount(s.OurBalanceProof))
	return x
}
func (s *Serialization) PartnerLock2UnclaimedLocks() map[common.Hash]UnlockPartialProof {
	return nil
}
func (s *Serialization) OurLock2UnclaimedLocks() map[common.Hash]UnlockPartialProof {
	return nil
}
func (s *Serialization) OurLock2PendingLocks() map[common.Hash]PendingLock {
	return nil
}
func (s *Serialization) PartnerLock2PendingLocks() map[common.Hash]PendingLock {
	return nil
}
func init() {
	gob.Register(&PendingLock{})
	gob.Register(&UnlockPartialProof{})
	gob.Register(&Serialization{})

	//make sure don't save this data
	//gob.Register(&UnlockProof{})
}
