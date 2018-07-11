package channeltype

import (
	"encoding/gob"
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
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
	ChannelIdentifier          *contracts.ChannelUniqueID
	ChannelAddressString       string `storm:"id"` //only for storm, because of save bug
	TokenAddress               common.Address
	PartnerAddress             common.Address
	TokenAddressString         string `storm:"index"`
	PartnerAddressString       string `storm:"index"`
	OurAddress                 common.Address
	RevealTimeout              int
	OurBalanceProof            *transfer.BalanceProofState
	PartnerBalanceProof        *transfer.BalanceProofState
	OurLeaves                  []common.Hash
	PartnerLeaves              []common.Hash
	OurLock2PendingLocks       map[common.Hash]PendingLock
	OurLock2UnclaimedLocks     map[common.Hash]UnlockPartialProof
	PartnerLock2PendingLocks   map[common.Hash]PendingLock
	PartnerLock2UnclaimedLocks map[common.Hash]UnlockPartialProof
	State                      State
	OurBalance                 *big.Int
	PartnerBalance             *big.Int
	OurContractBalance         *big.Int
	PartnerContractBalance     *big.Int
	OurAmountLocked            *big.Int
	PartnerAmountLocked        *big.Int
	ClosedBlock                int64
	SettledBlock               int64
	SettleTimeout              int
}

func init() {
	gob.Register(&PendingLock{})
	gob.Register(&UnlockPartialProof{})
	//make sure don't save this data
	//gob.Register(&UnlockProof{})
}
