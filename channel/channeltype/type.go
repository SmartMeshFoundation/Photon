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
	Lock     *mtree.Lock `json:"lock"`
	LockHash common.Hash `json:"lock_hash"` //hash of this lock
}

/*
UnlockPartialProof is the lock that I have known the secret ,but haven't receive the Balance proof
*/
type UnlockPartialProof struct {
	Lock                *mtree.Lock `json:"lock"`
	LockHash            common.Hash `json:"lock_hash"`
	Secret              common.Hash `json:"secret"`
	IsRegisteredOnChain bool        `json:"is_registered_on_chain"` //该密码是通过链上注册获知的还是通过普通的RevealSecret得知的. 如果是链上注册得知,那么一定是没有过期的.
}

/*
UnlockProof is the info needs withdraw on blockchain
*/
type UnlockProof struct {
	MerkleProof []common.Hash `json:"merkle_proof"`
	Lock        *mtree.Lock   `json:"lock"`
}

//KnownSecret is used to save to db
type KnownSecret struct {
	Secret              common.Hash `json:"secret"`
	IsRegisteredOnChain bool        `json:"is_registered_on_chain"` //该密码是通过链上注册获知的还是通过普通的RevealSecret得知的. 如果是链上注册得知,那么一定是没有过期的.
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
	OurKnownSecrets        []*KnownSecret
	PartnerKnownSecrets    []*KnownSecret
	State                  State
	OurContractBalance     *big.Int
	PartnerContractBalance *big.Int
	ClosedBlock            int64
	SettledBlock           int64
	SettleTimeout          int
}

// GetKey : impl dao.KeyGetter
func (s *Serialization) GetKey() []byte {
	return s.Key
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
		if m2[l.LockSecretHash] != nil {
			//知道密码
			m[l.LockSecretHash] = UnlockPartialProof{
				Lock:                l,
				Secret:              m2[l.LockSecretHash].Secret,
				LockHash:            l.Hash(),
				IsRegisteredOnChain: m2[l.LockSecretHash].IsRegisteredOnChain,
			}
		}
	}
	return m
}

func (s *Serialization) getSecretHashMap(secrets []*KnownSecret) map[common.Hash]*KnownSecret {
	m := make(map[common.Hash]*KnownSecret)
	for _, s := range secrets {
		m[utils.ShaSecret(s.Secret[:])] = s
	}
	return m
}

//OurLock2UnclaimedLocks our lock and know secret
func (s *Serialization) OurLock2UnclaimedLocks() map[common.Hash]UnlockPartialProof {
	m := make(map[common.Hash]UnlockPartialProof)
	m2 := s.getSecretHashMap(s.OurKnownSecrets)
	for _, l := range s.OurLeaves {
		if m2[l.LockSecretHash] != nil {
			//知道密码
			m[l.LockSecretHash] = UnlockPartialProof{
				Lock:                l,
				Secret:              m2[l.LockSecretHash].Secret,
				LockHash:            l.Hash(),
				IsRegisteredOnChain: m2[l.LockSecretHash].IsRegisteredOnChain,
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
		if m2[l.LockSecretHash] == nil {
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
		if m2[l.LockSecretHash] == nil {
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
		if m2[l.LockSecretHash] != nil {
			//知道密码
			if l.Expiration > blockNumber && l.Expiration > min {
				min = l.Expiration
			}
		}
	}
	return min
}

//ChannelDataDetail for user api
type ChannelDataDetail struct {
	ChannelIdentifier   string   `json:"channel_identifier"`
	OpenBlockNumber     int64    `json:"open_block_number"`
	PartnerAddress      string   `json:"partner_address"`
	Balance             *big.Int `json:"balance"`
	PartnerBalance      *big.Int `json:"partner_balance"`
	LockedAmount        *big.Int `json:"locked_amount"`
	PartnerLockedAmount *big.Int `json:"partner_locked_amount"`
	TokenAddress        string   `json:"token_address"`
	State               State    `json:"state"`
	StateString         string   `json:"state_string"`
	SettleTimeout       int      `json:"settle_timeout"`
	RevealTimeout       int      `json:"reveal_timeout"`

	/*
		extended
	*/
	ClosedBlock               int64                              `json:"closed_block"`
	SettledBlock              int64                              `json:"settled_block"`
	OurUnknownSecretLocks     map[common.Hash]PendingLock        `json:"our_unknown_secret_locks,omitempty"`
	OurKnownSecretLocks       map[common.Hash]UnlockPartialProof `json:"our_known_secret_locks,omitempty"`
	PartnerUnknownSecretLocks map[common.Hash]PendingLock        `json:"partner_unknown_secret_locks,omitempty"`
	PartnerKnownSecretLocks   map[common.Hash]UnlockPartialProof `json:"partner_known_secret_locks,omitempty"`
	OurLeaves                 []*mtree.Lock                      `json:"our_leaves,omitempty"`
	PartnerLeaves             []*mtree.Lock                      `json:"partner_leaves,omitempty"`
	OurBalanceProof           *transfer.BalanceProofState        `json:"our_balance_proof,omitempty"`
	PartnerBalanceProof       *transfer.BalanceProofState        `json:"partner_balance_proof,omitempty"`
	Signature                 []byte                             `json:"signature,omitempty"` //my signature of PartnerBalanceProof
}

//ChannelSerialization2ChannelDataDetail 辅助函数
func ChannelSerialization2ChannelDataDetail(c *Serialization) *ChannelDataDetail {
	d := &ChannelDataDetail{
		ChannelIdentifier:         c.ChannelIdentifier.ChannelIdentifier.String(),
		OpenBlockNumber:           c.ChannelIdentifier.OpenBlockNumber,
		PartnerAddress:            c.PartnerAddress().String(),
		Balance:                   c.OurBalance(),
		PartnerBalance:            c.PartnerBalance(),
		State:                     c.State,
		StateString:               c.State.String(),
		SettleTimeout:             c.SettleTimeout,
		RevealTimeout:             c.RevealTimeout,
		TokenAddress:              c.TokenAddress().String(),
		LockedAmount:              c.OurAmountLocked(),
		PartnerLockedAmount:       c.PartnerAmountLocked(),
		ClosedBlock:               c.ClosedBlock,
		SettledBlock:              c.SettledBlock,
		OurLeaves:                 c.OurLeaves,
		PartnerLeaves:             c.PartnerLeaves,
		OurKnownSecretLocks:       c.OurLock2UnclaimedLocks(),
		OurUnknownSecretLocks:     c.OurLock2PendingLocks(),
		PartnerUnknownSecretLocks: c.PartnerLock2PendingLocks(),
		PartnerKnownSecretLocks:   c.PartnerLock2UnclaimedLocks(),
		OurBalanceProof:           c.OurBalanceProof,
		PartnerBalanceProof:       c.PartnerBalanceProof,
	}
	return d
}

func init() {
	gob.Register(&PendingLock{})
	gob.Register(&UnlockPartialProof{})
	gob.Register(&Serialization{})

	//make sure don't save this data
	//gob.Register(&UnlockProof{})
}
