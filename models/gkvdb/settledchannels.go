package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/ethereum/go-ethereum/common"
)

//NewSettledChannel save a settled channel to db
func (dao *GkvDB) NewSettledChannel(c *channeltype.Serialization) error {
	if c.State != channeltype.StateSettled {
		panic("only settled channel can saved to settledChannel")
	}
	key := fmt.Sprintf("%s-%d", c.ChannelIdentifier.ChannelIdentifier.String(), c.ChannelIdentifier.OpenBlockNumber)
	return dao.saveKeyValueToBucket(models.BucketSettledChannel, key, c)
}

//GetAllSettledChannel returns all settled channel
func (dao *GkvDB) GetAllSettledChannel() (chs []*channeltype.Serialization, err error) {
	tb, err := dao.db.Table(models.BucketSettledChannel)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var channel channeltype.Serialization
		gobDecode(v, &channel)
		chs = append(chs, &channel)
	}
	return
}

//GetSettledChannel 返回某个指定的已经 settle 的 channel
// GetSettledChannel : function to return a specific settled channel.
func (dao *GkvDB) GetSettledChannel(channelIdentifier common.Hash, openBlockNumber int64) (c *channeltype.Serialization, err error) {
	c = new(channeltype.Serialization)
	key := fmt.Sprintf("%s-%d", channelIdentifier.String(), openBlockNumber)
	err = dao.getKeyValueToBucket(models.BucketSettledChannel, key, c)
	return
}
