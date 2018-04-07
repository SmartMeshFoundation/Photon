package models

import (
	"fmt"

	"bytes"
	"encoding/hex"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func (model *ModelDB) NewChannel(c *channel.ChannelSerialization) error {
	log.Trace(fmt.Sprintf("new channel %s", c.ChannelAddress.String()))
	err := model.db.Save(c)
	//notify new channel added
	model.handleChannelCallback(model.newChannelCallbacks, c)
	if err != nil {
		log.Error(fmt.Sprintf("NewChannel for models err:%s", err))
	}
	return err
}
func (model *ModelDB) UpdateChannelNoTx(c *channel.ChannelSerialization) error {
	log.Trace(fmt.Sprintf("save channel %s", c.ChannelAddress.String()))
	err := model.db.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelNoTx err:%s", err))
	}
	return err
}
func (model *ModelDB) handleChannelCallback(m map[*ChannelCb]bool, c *channel.ChannelSerialization) {
	var cbs []*ChannelCb
	model.mlock.Lock()
	for f, _ := range m {
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

//update channel balance
func (model *ModelDB) UpdateChannelContractBalance(c *channel.ChannelSerialization) error {
	err := model.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.channelDepositCallbacks, c)
	return nil
}

//update channel balance? transfer complete?
func (model *ModelDB) UpdateChannel(c *channel.ChannelSerialization, tx storm.Node) error {
	//log.Trace(fmt.Sprintf("statemanager save channel status =%s\n", utils.StringInterface(c, 7)))
	err := tx.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannel err=%s", err))
	}
	return err
}

//update channel state ,close settle
func (model *ModelDB) UpdateChannelState(c *channel.ChannelSerialization) error {
	err := model.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.channelStateCallbacks, c)
	return nil
}

//channel (token,partner)
func (model *ModelDB) GetChannel(token, partner common.Address) (c *channel.ChannelSerialization, err error) {
	var cs []*channel.ChannelSerialization
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
		if c2.PartnerAddress == partner {
			c = c2
			return
		}
	}
	return nil, storm.ErrNotFound
}

//channel (token,partner)
func (model *ModelDB) GetChannelByAddress(channelAddress common.Address) (c *channel.ChannelSerialization, err error) {
	var c2 channel.ChannelSerialization
	err = model.db.One("ChannelAddressString", channelAddress.String(), &c2)
	if err == nil {
		c = &c2
	}
	return
}

//one of token and partner must be empty
func (model *ModelDB) GetChannelList(token, partner common.Address) (cs []*channel.ChannelSerialization, err error) {
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
	is secret has withdrawed on channel?
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
 I have withdrawed this secret on channel.
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
	is a expired hashlock has been removed from channel status.
*/
func (model *ModelDB) IsThisLockRemoved(channel common.Address,sender common.Address, secret common.Hash) bool {
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
	remember this lock has been removed from channel status.
*/
func (model *ModelDB) RemoveLock(channel common.Address,sender common.Address, secret common.Hash) {
	key := new(bytes.Buffer)
	key.Write(channel[:])
	key.Write(secret[:])
	key.Write(sender[:])
	err := model.db.Set(bucketExpiredHashlock, key.Bytes(), true)
	if err != nil {
		log.Error(fmt.Sprintf("WithdrawThisLock write %s to db err %s", hex.EncodeToString(key.Bytes()), err))
	}
}
