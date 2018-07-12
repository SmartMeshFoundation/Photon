package models

import (
	"fmt"

	"bytes"
	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

// NewChannel save a just created channel to db
func (model *ModelDB) NewChannel(c *channel.Serialization) error {
	log.Trace(fmt.Sprintf("new channel %s", c.ChannelAddress.String()))
	err := model.db.Save(c)
	//notify new channel added
	model.handleChannelCallback(model.newChannelCallbacks, c)
	if err != nil {
		log.Error(fmt.Sprintf("NewChannel for models err:%s", err))
	}
	return err
}

//UpdateChannelNoTx update channel status without a Tx
func (model *ModelDB) UpdateChannelNoTx(c *channel.Serialization) error {
	log.Trace(fmt.Sprintf("save channel %s", c.ChannelAddress.String()))
	err := model.db.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelNoTx err:%s", err))
	}
	return err
}
func (model *ModelDB) handleChannelCallback(m map[*ChannelCb]bool, c *channel.Serialization) {
	var cbs []*ChannelCb
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
func (model *ModelDB) UpdateChannelContractBalance(c *channel.Serialization) error {
	err := model.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.channelDepositCallbacks, c)
	return nil
}

//UpdateChannel update channel status in a Tx
func (model *ModelDB) UpdateChannel(c *channel.Serialization, tx storm.Node) error {
	//log.Trace(fmt.Sprintf("statemanager save channel status =%s\n", utils.StringInterface(c, 7)))
	err := tx.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannel err=%s", err))
	}
	return err
}

//UpdateChannelState update channel state ,close settle
func (model *ModelDB) UpdateChannelState(c *channel.Serialization) error {
	err := model.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.channelStateCallbacks, c)
	return nil
}

//GetChannel return a channel queried by (token,partner),this channel must not settled
func (model *ModelDB) GetChannel(token, partner common.Address) (c *channel.Serialization, err error) {
	var cs []*channel.Serialization
	if token == utils.EmptyAddress {
		panic("token is empty")
	}
	if partner == utils.EmptyAddress {
		panic("partner is empty")
	}
	err = model.db.Find("TokenAddressString", token.String(), &cs)
	if err != nil {
		return
	}
	for _, c2 := range cs {
		if c2.PartnerAddress == partner && c2.State != transfer.ChannelStateSettled {
			c = c2
			return
		}
	}
	return nil, storm.ErrNotFound
}

//GetChannelByAddress return a channel queried by channel address
func (model *ModelDB) GetChannelByAddress(channelAddress common.Address) (c *channel.Serialization, err error) {
	var c2 channel.Serialization
	err = model.db.One("ChannelAddressString", channelAddress.String(), &c2)
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
		err = model.db.Find("PartnerAddressString", partner.String(), &cs)
	} else if partner == utils.EmptyAddress {
		err = model.db.Find("TokenAddressString", token.String(), &cs)
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
IsThisLockHasWithdraw return ture when  secret has withdrawed on channel?
*/
func (model *ModelDB) IsThisLockHasWithdraw(channel common.Address, secret common.Hash) bool {
	var result bool
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(secret[:])
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
func (model *ModelDB) WithdrawThisLock(channel common.Address, secret common.Hash) {
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(secret[:])
	err := model.db.Set(bucketWithDraw, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("WithdrawThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}

const bucketExpiredHashlock = "expiredHashlock"

/*
IsThisLockRemoved return true when  a expired hashlock has been removed from channel status.
*/
func (model *ModelDB) IsThisLockRemoved(channel common.Address, sender common.Address, secret common.Hash) bool {
	var result bool
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(secret[:])
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
func (model *ModelDB) RemoveLock(channel common.Address, sender common.Address, secret common.Hash) {
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(secret[:])
	key.Write(sender[:])
	err := model.db.Set(bucketExpiredHashlock, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("WithdrawThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}
