package channeltype

import (
	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
Db is the interface of Database operations about channel
*/
type Db interface {
	/*
		is secret has withdrawed on channel?
	*/
	IsThisLockHasWithdraw(channel common.Hash, secretHash common.Hash) bool
	/*
	 I have withdrawed this secret on channel.
	*/
	WithdrawThisLock(channel common.Hash, secretHash common.Hash)
	/*
		is a expired hashlock has been removed from channel status.
	*/
	IsThisLockRemoved(channel common.Hash, sender common.Address, secretHash common.Hash) bool
	/*
		remember this lock has been removed from channel status.
	*/
	RemoveLock(channel common.Hash, sender common.Address, secretHash common.Hash)
	/*
		get the latest channel status
	*/
	GetChannelByAddress(channelAddress common.Hash) (c *Serialization, err error)

	/*
	 要记录自己放在某个 channel 上放弃了某个锁,到时候一定不能unlock
	*/
	//CanUnlockThisLock(LockSecretHash common.Hash, channelIdentifier common.Hash) bool
}

//MockChannelDb for test only
type MockChannelDb struct {
	Keys map[common.Hash]bool
}

func NewMockChannelDb() Db {
	return &MockChannelDb{
		Keys: make(map[common.Hash]bool),
	}
}

/*
	is secret has withdrawed on channel?
*/
func (f *MockChannelDb) IsThisLockHasWithdraw(channel common.Hash, secretHash common.Hash) bool {
	hash := utils.Sha3(channel[:], secretHash[:])
	return f.Keys[hash]
}

/*
 I have withdrawed this secret on channel.
*/
func (f *MockChannelDb) WithdrawThisLock(channel common.Hash, secretHash common.Hash) {
	hash := utils.Sha3(channel[:], secretHash[:])
	f.Keys[hash] = true
}

/*
	is a expired hashlock has been removed from channel status.
*/
func (f *MockChannelDb) IsThisLockRemoved(channel common.Hash, sender common.Address, secretHash common.Hash) bool {
	hash := utils.Sha3(channel[:], sender[:], secretHash[:])
	return f.Keys[hash]
}

/*
	remember this lock has been removed from channel status.
*/
func (f *MockChannelDb) RemoveLock(channel common.Hash, sender common.Address, secretHash common.Hash) {
	hash := utils.Sha3(channel[:], sender[:], secretHash[:])
	f.Keys[hash] = true
}

/*
	get the latest channel status
*/
func (f *MockChannelDb) GetChannelByAddress(channelAddress common.Hash) (c *Serialization, err error) {
	return nil, errors.New("not found")
}
