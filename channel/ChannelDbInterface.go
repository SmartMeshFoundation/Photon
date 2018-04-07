package channel

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-errors/errors"
	"github.com/SmartMeshFoundation/raiden-network/utils"
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
	IsThisLockRemoved(channel common.Address,sender common.Address, secret common.Hash) bool
	/*
		remember this lock has been removed from channel status.
	*/
	RemoveLock(channel common.Address,sender common.Address, secret common.Hash)
	/*
		get the latest channel status
	*/
	GetChannelByAddress(channelAddress common.Address) (c *ChannelSerialization, err error)
}

//for test only
type MockChannelDb struct{
	Keys map[common.Hash]bool
}
func NewMockChannelDb() *MockChannelDb {
	return &MockChannelDb{
		Keys:make(map[common.Hash]bool),
	}
}
/*
	is secret has withdrawed on channel?
*/
func (f*MockChannelDb) IsThisLockHasWithdraw(channel common.Address, secret common.Hash) bool {
	hash:=utils.Sha3(channel[:],secret[:])
	return f.Keys[hash]
}
/*
 I have withdrawed this secret on channel.
*/
func (f*MockChannelDb) WithdrawThisLock(channel common.Address, secret common.Hash) {
	hash:=utils.Sha3(channel[:],secret[:])
	f.Keys[hash]=true
}
/*
	is a expired hashlock has been removed from channel status.
*/
func (f*MockChannelDb) IsThisLockRemoved(channel common.Address,sender common.Address, secret common.Hash) bool {
	hash:=utils.Sha3(channel[:],sender[:],secret[:])
	return f.Keys[hash]
}
/*
	remember this lock has been removed from channel status.
*/
func (f*MockChannelDb) RemoveLock(channel common.Address,sender common.Address, secret common.Hash) {
	hash:=utils.Sha3(channel[:],sender[:],secret[:])
	f.Keys[hash]=true
}
/*
	get the latest channel status
*/
func (f*MockChannelDb) GetChannelByAddress(channelAddress common.Address) (c *ChannelSerialization, err error) {
	return nil,errors.New("not found")
}