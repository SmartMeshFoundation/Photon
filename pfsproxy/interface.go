package pfsproxy

import (
	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/common"
)

/*
PfsProxy :
api to call pfg server
*/
type PfsProxy interface {
	/*
		submit partner's balance proof to pfg
	*/
	SubmitBalance(nonce uint64, transferAmount, lockAmount *big.Int, openBlockNumber int64, locksroot, channelIdentifier, additionHash common.Hash, signature []byte) error

	/*
		set fee rate of a channel to pfg
	*/
	SetFeeRate(channelIdentifier common.Hash, feeRate *big.Float) error

	/*
		get fee rate of a channel on pfg
	*/
	GetFeeRate(nodeAddress common.Address, channelIdentifier common.Hash) (feeRate *big.Float, effectiveTime time.Time, err error)

	/*
		find path
	*/
	FindPath(peerFrom, peerTo, token common.Address, amount *big.Int) (resp []FindPathResponse, err error)
}
