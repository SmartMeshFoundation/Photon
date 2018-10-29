package channeltype

import (
	"encoding/gob"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
PendingLock is lock of HTLC
*/
type PendingLock struct {
	Lock     *mtree.Lock
	LockHash common.Hash //hash of this lock
}

/*
UnlockPartialProof is the lock that I have known the secret ,but haven't receive the Balance proof
*/
type UnlockPartialProof struct {
	Lock     *mtree.Lock
	LockHash common.Hash
	Secret   common.Hash
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

//ChannleAddress address of channel
func (s *Serialization) ChannleAddress() common.Hash {
	return common.BytesToHash(s.Key)
}

//TokenAddress address of token
func (s *Serialization) TokenAddress() common.Address {
	return common.BytesToAddress(s.TokenAddressBytes)
}

//PartnerAddress partner's address
func (s *Serialization) PartnerAddress() common.Address {
	return common.BytesToAddress(s.PartnerAddressBytes)
}
func (s *Serialization) transferAmount(bp *transfer.BalanceProofState) *big.Int {
	if bp != nil {
		return bp.TransferAmount
	}
	return utils.BigInt0
}

//OurBalance our abalance
func (s *Serialization) OurBalance() *big.Int {
	x := new(big.Int)
	x.Sub(s.OurContractBalance, s.transferAmount(s.OurBalanceProof))
	x.Add(x, s.transferAmount(s.PartnerBalanceProof))
	return x
}

//OurAmountLocked sending token on road
func (s *Serialization) OurAmountLocked() *big.Int {
	x := new(big.Int)
	for _, l := range s.OurLeaves {
		x = x.Add(x, l.Amount)
	}
	return x
}

//PartnerAmountLocked received token on road
func (s *Serialization) PartnerAmountLocked() *big.Int {
	x := new(big.Int)
	for _, l := range s.PartnerLeaves {
		x = x.Add(x, l.Amount)
	}
	return x
}

//PartnerBalance partner's balance
func (s *Serialization) PartnerBalance() *big.Int {
	x := new(big.Int)
	x.Sub(s.PartnerContractBalance, s.transferAmount(s.PartnerBalanceProof))
	x.Add(x, s.transferAmount(s.OurBalanceProof))
	return x
}

//PartnerLock2UnclaimedLocks partner's lock and known secret
func (s *Serialization) PartnerLock2UnclaimedLocks() map[common.Hash]UnlockPartialProof {
	m := make(map[common.Hash]UnlockPartialProof)
	m2 := s.getSecretHashMap(s.PartnerKnownSecrets)
	for _, l := range s.PartnerLeaves {
		if m2[l.LockSecretHash] != utils.EmptyHash {
			//知道密码
			m[l.LockSecretHash] = UnlockPartialProof{
				Lock:     l,
				Secret:   m2[l.LockSecretHash],
				LockHash: l.Hash(),
			}
		}
	}
	return m
}

func (s *Serialization) getSecretHashMap(secrets []common.Hash) map[common.Hash]common.Hash {
	m := make(map[common.Hash]common.Hash)
	for _, s := range secrets {
		m[utils.ShaSecret(s[:])] = s
	}
	return m
}

//OurLock2UnclaimedLocks our lock and know secret
func (s *Serialization) OurLock2UnclaimedLocks() map[common.Hash]UnlockPartialProof {
	m := make(map[common.Hash]UnlockPartialProof)
	m2 := s.getSecretHashMap(s.OurKnownSecrets)
	for _, l := range s.OurLeaves {
		if m2[l.LockSecretHash] != utils.EmptyHash {
			//知道密码
			m[l.LockSecretHash] = UnlockPartialProof{
				Lock:     l,
				Secret:   m2[l.LockSecretHash],
				LockHash: l.Hash(),
			}
		}
	}
	return m
}

//OurLock2PendingLocks our lock and don't know secret
func (s *Serialization) OurLock2PendingLocks() map[common.Hash]PendingLock {
	m := make(map[common.Hash]PendingLock)
	m2 := s.getSecretHashMap(s.OurKnownSecrets)
	for _, l := range s.OurLeaves {
		if m2[l.LockSecretHash] == utils.EmptyHash {
			//不知道密码
			m[l.LockSecretHash] = PendingLock{
				Lock:     l,
				LockHash: l.Hash(),
			}
		}
	}
	return m
}

//PartnerLock2PendingLocks partner's lock and don't know secret
func (s *Serialization) PartnerLock2PendingLocks() map[common.Hash]PendingLock {
	m := make(map[common.Hash]PendingLock)
	m2 := s.getSecretHashMap(s.PartnerKnownSecrets)
	for _, l := range s.PartnerLeaves {
		if m2[l.LockSecretHash] == utils.EmptyHash {
			//不知道密码
			m[l.LockSecretHash] = PendingLock{
				Lock:     l,
				LockHash: l.Hash(),
			}
		}
	}
	return m
}

//MinExpiration 返回对方所有锁中,
// 知道密码过期时间最小的那个,如果已经超过了 expiration,忽略就可.
/*
 *	MinExpiration : return the block number of all locks in channel partner
 * 		which is closest to `expiration` of knownsecrets
 *		if none, return 0.
 */
func (s *Serialization) MinExpiration(blockNumber int64) int64 {
	m2 := s.getSecretHashMap(s.PartnerKnownSecrets)
	var min int64
	for _, l := range s.PartnerLeaves {
		if m2[l.LockSecretHash] != utils.EmptyHash {
			//知道密码
			if l.Expiration > blockNumber && l.Expiration > min {
				min = l.Expiration
			}
		}
	}
	return min
}
func init() {
	gob.Register(&PendingLock{})
	gob.Register(&UnlockPartialProof{})
	gob.Register(&Serialization{})

	//make sure don't save this data
	//gob.Register(&UnlockProof{})
}
