package pfsproxy

import (
	"math/big"

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
		find path
	*/
	FindPath(peerFrom, peerTo, token common.Address, amount *big.Int) (resp []FindPathResponse, err error)

	/*
		set fee rate by account
	*/
	SetAccountFee(feeConstant *big.Int, feePercent int64) (err error)

	/*
		get fee rate by account
	*/
	GetAccountFee() (feeConstant *big.Int, feePercent int64, err error)
	/*
		set fee rate by token
	*/
	SetTokenFee(feeConstant *big.Int, feePercent int64, tokenAddress common.Address) (err error)

	/*
		get fee rate by token
	*/
	GetTokenFee(tokenAddress common.Address) (feeConstant *big.Int, feePercent int64, err error)
	/*
		set fee rate by channel
	*/
	SetChannelFee(feeConstant *big.Int, feePercent int64, channelIdentifier common.Hash) (err error)

	/*
		get fee rate by channel
	*/
	GetChannelFee(channelIdentifier common.Hash) (feeConstant *big.Int, feePercent int64, err error)
}
