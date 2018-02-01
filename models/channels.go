package models

import (
	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

func (model *ModelDB) AddChannel(c *channel.ChannelSerialization) error {
	var c2 channel.ChannelSerialization
	err := model.db.One("ChannelAddress", c.ChannelAddress, &c2)
	if err == nil {
		err = model.db.Update(c)
	} else {
		err = model.db.Save(c)
		//notify new channel added
		model.handleChannelCallback(model.NewChannelCallbacks, c)

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
	err := model.AddChannel(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.ChannelDepositCallbacks, c)
	return nil
}

//update channel balance? transfer complete?

//update channel state ,close settle
func (model *ModelDB) UpdateChannelState(c *channel.ChannelSerialization) error {
	err := model.AddChannel(c)
	if err != nil {
		return err
	}
	//notify listener
	model.handleChannelCallback(model.ChannelStateCallbacks, c)
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
	err = model.db.Find("TokenAddress", token, &cs)
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
	err = model.db.One("ChannelAddress", channelAddress, &c2)
	if err == nil {
		c = &c2
	}
	return
}

//one of token and partner must be empty
func (model *ModelDB) GetChannelList(token, partner common.Address) (cs []*channel.ChannelSerialization, err error) {
	if token == utils.EmptyAddress && partner == utils.EmptyAddress {
		err = model.db.All(&cs)
		return
	} else if token == utils.EmptyAddress {
		err = model.db.Find("PartnerAddress", partner, &cs)
		return
	} else if partner == utils.EmptyAddress {
		err = model.db.Find("TokenAddress", token, &cs)
		return
	} else {
		panic("one of token and partner must be empty")
	}
	return
}
