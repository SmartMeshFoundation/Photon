package models

import (
	"fmt"

	"bytes"
	"encoding/gob"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/common"
)

const bucketSettledChannel = "settled_channel"

func Unmarshal(b []byte, v interface{}) error {
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}
func (model *ModelDB) NewSettledChannel(c *channeltype.Serialization) error {
	if c.State != channeltype.StateSettled {
		panic("only settled channel can saved to settledChannel")
	}
	key := fmt.Sprintf("%s-%d", c.ChannelIdentifier.ChannelIdentifier.String(), c.ChannelIdentifier.OpenBlockNumber)
	return model.db.Set(bucketSettledChannel, key, c)
}

//如何便利一个 bucket 呢?应该比较容易
func (model *ModelDB) GetAllSettledChannel() (chs []*channeltype.Serialization, err error) {
	model.db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSettledChannel))
		b.ForEach(func(k, v []byte) error {
			if string(k) == "__storm_metadata" {
				return nil
			}
			log.Trace(fmt.Sprintf("GetAllSettledChannel key=%s, value=%s\n", string(k), hex.EncodeToString(v)))
			var c channeltype.Serialization
			err = Unmarshal(v, &c)
			if err != nil {
				return err
			}
			chs = append(chs, &c)
			return nil
		})
		return nil
	})
	return
}

func (model *ModelDB) GetSettledChannel(channelIdentifier common.Hash, openBlockNumber int64) (c *channeltype.Serialization, err error) {
	c = new(channeltype.Serialization)
	key := fmt.Sprintf("%s-%d", channelIdentifier.String(), openBlockNumber)
	err = model.db.Get(bucketSettledChannel, key, c)
	return
}
