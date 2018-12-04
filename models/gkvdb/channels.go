package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// NewChannel save a just created channel to db
func (dao *GkvDB) NewChannel(c *channeltype.Serialization) error {
	//log.Trace(fmt.Sprintf("new channel %s", utils.StringInterface(c, 2)))
	err := dao.saveKeyValueToBucket(models.BucketChannelSerialization, c.GetKey(), c)
	//notify new channel added
	dao.handleChannelCallback(dao.newChannelCallbacks, c)
	if err != nil {
		log.Error(fmt.Sprintf("NewChannel for daos err:%s", err))
	}
	return err
}

//UpdateChannelNoTx update channel status without a Tx
func (dao *GkvDB) UpdateChannelNoTx(c *channeltype.Serialization) error {
	//log.Trace(fmt.Sprintf("save channel %s", utils.StringInterface(c, 2)))
	err := dao.saveKeyValueToBucket(models.BucketChannelSerialization, c.GetKey(), c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannelNoTx err:%s", err))
	}
	return err
}

//UpdateChannelAndSaveAck update channel and save ack, must atomic
func (dao *GkvDB) UpdateChannelAndSaveAck(c *channeltype.Serialization, echohash common.Hash, ack []byte) (err error) {
	tx := dao.StartTx(models.BucketChannelSerialization)
	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()
	err = dao.UpdateChannel(c, tx)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannel err %s", err))
		return
	}
	dao.SaveAck(echohash, ack, tx)
	err = tx.Commit()
	return
}
func (dao *GkvDB) handleChannelCallback(m map[*cb.ChannelCb]bool, c *channeltype.Serialization) {
	var cbs []*cb.ChannelCb
	dao.mlock.Lock()
	for f := range m {
		remove := (*f)(c)
		if remove {
			cbs = append(cbs, f)
		}
	}
	for _, f := range cbs {
		delete(m, f)
	}
	dao.mlock.Unlock()
}

//UpdateChannelContractBalance update channel balance
func (dao *GkvDB) UpdateChannelContractBalance(c *channeltype.Serialization) error {
	err := dao.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	dao.handleChannelCallback(dao.channelDepositCallbacks, c)
	return nil
}

//UpdateChannel update channel status in a Tx
func (dao *GkvDB) UpdateChannel(c *channeltype.Serialization, tx models.TX) error {
	//log.Trace(fmt.Sprintf("statemanager save channel status =%s\n", utils.StringInterface(c, 2)))
	err := tx.Save(c)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateChannel err=%s", err))
	}
	return err
}

//UpdateChannelState update channel state ,close settle
func (dao *GkvDB) UpdateChannelState(c *channeltype.Serialization) error {
	err := dao.UpdateChannelNoTx(c)
	if err != nil {
		return err
	}
	//notify listener
	dao.handleChannelCallback(dao.channelStateCallbacks, c)
	return nil
}

//RemoveChannel a settled channel from db
func (dao *GkvDB) RemoveChannel(c *channeltype.Serialization) error {
	if c.State != channeltype.StateSettled {
		panic("only can remove a settled channel")
	}
	dao.handleChannelCallback(dao.channelSettledCallbacks, c)
	return dao.removeKeyValueFromBucket(models.BucketChannelSerialization, c.GetKey())
}

//GetChannel return a channel queried by (token,partner),this channel must not settled
func (dao *GkvDB) GetChannel(token, partner common.Address) (c *channeltype.Serialization, err error) {
	var cs []*channeltype.Serialization
	if token == utils.EmptyAddress {
		panic("token is empty")
	}
	if partner == utils.EmptyAddress {
		panic("partner is empty")
	}
	tb, err := dao.db.Table(models.BucketChannelSerialization)
	if err != nil {
		return
	}
	buf := tb.Values(-1)
	if len(buf) == 0 {
		err = ErrorNotFound
		return
	}
	for _, v := range buf {
		var ct channeltype.Serialization
		gobDecode(v, &ct)
		cs = append(cs, &ct)
	}
	for _, ct := range cs {
		if ct.TokenAddress() == token && ct.PartnerAddress() == partner {
			c = ct
			return
		}
	}
	return nil, ErrorNotFound
}

//GetChannelByAddress return a channel queried by channel address
func (dao *GkvDB) GetChannelByAddress(ChannelIdentifier common.Hash) (c *channeltype.Serialization, err error) {
	tb, err := dao.db.Table(models.BucketChannelSerialization)
	if err != nil {
		return
	}
	buf := tb.Values(-1)
	if len(buf) == 0 {
		err = ErrorNotFound
		return
	}
	var cs []*channeltype.Serialization
	for _, v := range buf {
		var ct channeltype.Serialization
		gobDecode(v, &ct)
		cs = append(cs, &ct)
	}
	for _, ct := range cs {
		if ct.ChannleAddress() == ChannelIdentifier {
			c = ct
			return
		}
	}
	return nil, ErrorNotFound
}

//GetChannelList returns all related channels
//one of token and partner must be empty
func (dao *GkvDB) GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error) {
	tb, err := dao.db.Table(models.BucketChannelSerialization)
	if err != nil {
		return
	}
	buf := tb.Values(-1)
	if len(buf) == 0 {
		err = ErrorNotFound
		return
	}
	var cst []*channeltype.Serialization
	for _, v := range buf {
		var ct channeltype.Serialization
		gobDecode(v, &ct)
		cst = append(cst, &ct)
	}

	if token == utils.EmptyAddress && partner == utils.EmptyAddress {
		cs = cst
	} else if token == utils.EmptyAddress {
		for _, c := range cst {
			if c.PartnerAddress() == partner {
				cs = append(cs, c)
			}
		}
	} else if partner == utils.EmptyAddress {
		for _, c := range cst {
			if c.TokenAddress() == token {
				cs = append(cs, c)
			}
		}
	} else {
		panic("one of token and partner must be empty")
	}
	if len(cs) == 0 {
		err = ErrorNotFound
	}
	return
}
