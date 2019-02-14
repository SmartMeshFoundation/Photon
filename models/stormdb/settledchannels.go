package stormdb

import (
	"fmt"

	"bytes"
	"encoding/gob"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/common"
)

func unmarshal(b []byte, v interface{}) error {
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}

//NewSettledChannel save a settled channel to db
func (model *StormDB) NewSettledChannel(c *channeltype.Serialization) error {
	if c.State != channeltype.StateSettled {
		panic("only settled channel can saved to settledChannel")
	}
	key := fmt.Sprintf("%s-%d", c.ChannelIdentifier.ChannelIdentifier.String(), c.ChannelIdentifier.OpenBlockNumber)
	err := model.db.Set(models.BucketSettledChannel, key, c)
	return models.GeneratDBError(err)
}

//GetAllSettledChannel returns all settled channel
func (model *StormDB) GetAllSettledChannel() (chs []*channeltype.Serialization, err error) {
	err = model.db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(models.BucketSettledChannel))
		err = b.ForEach(func(k, v []byte) error {
			if string(k) == "__storm_metadata" {
				return nil
			}
			//log.Trace(fmt.Sprintf("GetAllSettledChannel key=%s, value=%s\n", string(k), hex.EncodeToString(v)))
			var c channeltype.Serialization
			err = unmarshal(v, &c)
			if err != nil {
				return err
			}
			chs = append(chs, &c)
			return nil
		})
		return nil
	})
	err = models.GeneratDBError(err)
	return
}

//GetSettledChannel 返回某个指定的已经 settle 的 channel
// GetSettledChannel : function to return a specific settled channel.
func (model *StormDB) GetSettledChannel(channelIdentifier common.Hash, openBlockNumber int64) (c *channeltype.Serialization, err error) {
	c = new(channeltype.Serialization)
	key := fmt.Sprintf("%s-%d", channelIdentifier.String(), openBlockNumber)
	err = model.db.Get(models.BucketSettledChannel, key, c)
	err = models.GeneratDBError(err)
	return
}
