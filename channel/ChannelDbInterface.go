package channel

import (
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-errors/errors"
)

/*
Db is the interface of Database operations about channel
*/
type Db interface {
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
	IsThisLockRemoved(channel common.Address, sender common.Address, secret common.Hash) bool
	/*
		remember this lock has been removed from channel status.
	*/
	RemoveLock(channel common.Address, sender common.Address, secret common.Hash)
	/*
		get the latest channel status
	*/
	GetChannelByAddress(channelAddress common.Address) (c *Serialization, err error)
}

//for test only
type mockChannelDb struct {
	Keys map[common.Hash]bool
}

func newMockChannelDb() *mockChannelDb {
	return &mockChannelDb{
		Keys: make(map[common.Hash]bool),
	}
}

/*
	is secret has withdrawed on channel?
*/
func (f *mockChannelDb) IsThisLockHasWithdraw(channel common.Address, secret common.Hash) bool {
	hash := utils.Sha3(channel[:], secret[:])
	return f.Keys[hash]
}

/*
 I have withdrawed this secret on channel.
*/
func (f *mockChannelDb) WithdrawThisLock(channel common.Address, secret common.Hash) {
	hash := utils.Sha3(channel[:], secret[:])
	f.Keys[hash] = true
}

/*
	is a expired hashlock has been removed from channel status.
*/
func (f *mockChannelDb) IsThisLockRemoved(channel common.Address, sender common.Address, secret common.Hash) bool {
	hash := utils.Sha3(channel[:], sender[:], secret[:])
	return f.Keys[hash]
}

/*
	remember this lock has been removed from channel status.
*/
func (f *mockChannelDb) RemoveLock(channel common.Address, sender common.Address, secret common.Hash) {
	hash := utils.Sha3(channel[:], sender[:], secret[:])
	f.Keys[hash] = true
}

/*
	get the latest channel status
*/
func (f *mockChannelDb) GetChannelByAddress(channelAddress common.Address) (c *Serialization, err error) {
	return nil, errors.New("not found")
}
