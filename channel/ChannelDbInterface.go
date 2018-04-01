package channel

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChannelDb interface {
	/*
		is secret has withdrawed on channel?
	*/
	IsThisLockHasWithdraw(channel common.Address, secret common.Hash) bool
	/*
	 I have withdrawed this secret on channel.
	*/
	WithdrawThisLock(channel common.Address, secret common.Hash)
	/*
		is a expired hashlock has been removed from channel status.
	*/
	IsThisLockRemoved(channel common.Address, secret common.Hash) bool
	/*
		remember this lock has been removed from channel status.
	*/
	RemoveLock(channel common.Address, secret common.Hash)
	/*
		get the latest channel status
	*/
	GetChannelByAddress(channelAddress common.Address) (c *ChannelSerialization, err error)
}
