package models

import (
	"fmt"

	"bytes"
	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models/cb"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

// NewChannel save a just created channel to db
func (model *ModelDB) NewChannel(c *channeltype.Serialization) error {
	log.Trace(fmt.Sprintf("new channel %s", utils.StringInterface(c, 2)))
	err := model.db.Save(c)
	//notify new channel added
	model.handleChannelCallback(model.newChannelCallbacks, c)
	if err != nil {
		log.Error(fmt.Sprintf("NewChannel for models err:%s", err))
	}
	return err
}

//UpdateChannelNoTx update channel status without a Tx
func (model *ModelDB) UpdateChannelNoTx(c *channeltype.Serialization) error {
	log.Trace(fmt.Sprintf("save channel %s", utils.StringInterface(c, 2)))
	err := model.db.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelNoTx err:%s", err))
	}
	return err
}
func (model *ModelDB) handleChannelCallback(m map[*cb.ChannelCb]bool, c *channeltype.Serialization) {
	var cbs []*cb.ChannelCb
	model.mlock.Lock()
	for f := range m {
		remove := (*f)(c)
		if remove {
			cbs = append(cbs, f)
		}
	}
	for _, f := range cbs {
		delete(m, f)
	}
	model.mlock.Unlock()
}

//UpdateChannelContractBalance update channel balance
func (model *ModelDB) UpdateChannelContractBalance(c *channeltype.Serialization) error {
	err := model.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.channelDepositCallbacks, c)
	return nil
}

//UpdateChannel update channel status in a Tx
func (model *ModelDB) UpdateChannel(c *channeltype.Serialization, tx storm.Node) error {
	log.Trace(fmt.Sprintf("statemanager save channel status =%s\n", utils.StringInterface(c, 2)))
	err := tx.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannel err=%s", err))
	}
	return err
}

//UpdateChannelState update channel state ,close settle
func (model *ModelDB) UpdateChannelState(c *channeltype.Serialization) error {
	err := model.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.channelStateCallbacks, c)
	return nil
}

//RemoveChannel a settled channel from db
func (model *ModelDB) RemoveChannel(c *channeltype.Serialization) error {
	if c.State != channeltype.StateSettled {
		panic("only can remove a settled channel")
	}
	model.handleChannelCallback(model.channelSettledCallbacks, c)
	return model.db.DeleteStruct(c)
}

//GetChannel return a channel queried by (token,partner),this channel must not settled
func (model *ModelDB) GetChannel(token, partner common.Address) (c *channeltype.Serialization, err error) {
	var cs []*channeltype.Serialization
	if token == utils.EmptyAddress {
		panic("token is empty")
	}
	if partner == utils.EmptyAddress {
		panic("partner is empty")
	}
	err = model.db.Find("TokenAddressBytes", token[:], &cs)
	if err != nil {
		return
	}
	for _, c2 := range cs {
		if c2.PartnerAddress() == partner && c2.State != channeltype.StateSettled {
			c = c2
			return
		}
	}
	return nil, storm.ErrNotFound
}

//GetChannelByAddress return a channel queried by channel address
func (model *ModelDB) GetChannelByAddress(channelAddress common.Hash) (c *channeltype.Serialization, err error) {
	var c2 channeltype.Serialization
	err = model.db.One("Key", channelAddress[:], &c2)
	if err == nil {
		c = &c2
	}
	return
}

//GetChannelList returns all related channels
//one of token and partner must be empty
func (model *ModelDB) GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error) {
	if token == utils.EmptyAddress && partner == utils.EmptyAddress {
		err = model.db.All(&cs)
	} else if token == utils.EmptyAddress {
		err = model.db.Find("PartnerAddressBytes", partner[:], &cs)
	} else if partner == utils.EmptyAddress {
		err = model.db.Find("TokenAddressBytes", token[:], &cs)
	} else {
		panic("one of token and partner must be empty")
	}
	if err == storm.ErrNotFound {
		err = nil
	}
	return
}

const bucketWithDraw = "bucketWithdraw"

/*
IsThisLockHasWithdraw return ture when  lockhash has unlocked on channel?
*/
func (model *ModelDB) IsThisLockHasWithdraw(channel common.Hash, lockHash common.Hash) bool {
	var result bool
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(lockHash[:])
	err := model.db.Get(bucketWithDraw, key.Bytes(), &result)
	if err != nil {
		return false
	}
	if result != true {
		panic("withdraw cannot be set to false")
	}
	return result
}

/*
WithdrawThisLock marks that I have withdrawed this secret on channel.
*/
func (model *ModelDB) WithdrawThisLock(channel common.Hash, lockHash common.Hash) {
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(lockHash[:])
	err := model.db.Set(bucketWithDraw, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("WithdrawThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}

const bucketExpiredHashlock = "expiredHashlock"

/*
IsThisLockRemoved return true when  a expired hashlock has been removed from channel status.
*/
func (model *ModelDB) IsThisLockRemoved(channel common.Hash, sender common.Address, lockHash common.Hash) bool {
	var result bool
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(lockHash[:])
	key.Write(sender[:])
	err := model.db.Get(bucketExpiredHashlock, key.Bytes(), &result)
	if err != nil {
		return false
	}
	if result != true {
		panic("expiredHashlock cannot be set to false")
	}
	return result
}

/*
RemoveLock remember this lock has been removed from channel status.
*/
func (model *ModelDB) RemoveLock(channel common.Hash, sender common.Address, lockHash common.Hash) {
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(lockHash[:])
	key.Write(sender[:])
	err := model.db.Set(bucketExpiredHashlock, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("WithdrawThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}
