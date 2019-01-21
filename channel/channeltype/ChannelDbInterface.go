package channeltype

import (
	"errors"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
Db is the interface of Database operations about channel
*/
type Db interface {
	/*
		is secret has withdrawed on channel?
	*/
	IsThisLockHasUnlocked(channel common.Hash, lockHash common.Hash) bool
	/*
	 I have withdrawed this secret on channel.
	*/
	UnlockThisLock(channel common.Hash, lockHash common.Hash)
	/*
		is a expired hashlock has been removed from channel status.
	*/
	IsThisLockRemoved(channel common.Hash, sender common.Address, lockHash common.Hash) bool
	/*
		remember this lock has been removed from channel status.
	*/
	RemoveLock(channel common.Hash, sender common.Address, lockHash common.Hash)
	/*
		get the latest channel status
	*/
	GetChannelByAddress(channelIdentifier common.Hash) (c *Serialization, err error)

	/*
	 要记录自己放在某个 channel 上放弃了某个锁,到时候一定不能unlock
	*/
	IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool
}

//MockChannelDb for test only
type MockChannelDb struct {
	Keys map[common.Hash]bool
}

//NewMockChannelDb for test only
func NewMockChannelDb() Db {
	return &MockChannelDb{
		Keys: make(map[common.Hash]bool),
	}
}

//IsThisLockHasUnlocked is secret has withdrawed on channel
func (f *MockChannelDb) IsThisLockHasUnlocked(channel common.Hash, lockhash common.Hash) bool {
	hash := utils.Sha3(channel[:], lockhash[:])
	return f.Keys[hash]
}

//UnlockThisLock I have withdrawed this secret on channel.
func (f *MockChannelDb) UnlockThisLock(channel common.Hash, secretHash common.Hash) {
	hash := utils.Sha3(channel[:], secretHash[:])
	f.Keys[hash] = true
}

//IsThisLockRemoved is a expired hashlock has been removed from channel status.
func (f *MockChannelDb) IsThisLockRemoved(channel common.Hash, sender common.Address, secretHash common.Hash) bool {
	hash := utils.Sha3(channel[:], sender[:], secretHash[:])
	return f.Keys[hash]
}

//RemoveLock remember this lock has been removed from channel status.
func (f *MockChannelDb) RemoveLock(channel common.Hash, sender common.Address, secretHash common.Hash) {
	hash := utils.Sha3(channel[:], sender[:], secretHash[:])
	f.Keys[hash] = true
}

//GetChannelByAddress get the latest channel status
func (f *MockChannelDb) GetChannelByAddress(channelIdentifier common.Hash) (c *Serialization, err error) {
	return nil, errors.New("not found")
}

func (f *MockChannelDb) IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool {
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	return f.Keys[key]
}
