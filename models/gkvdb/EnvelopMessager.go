package gkvdb

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//NewSentEnvelopMessager create a sending EnvelopMessager in db
func (dao *GkvDB) NewSentEnvelopMessager(msg encoding.EnvelopMessager, receiver common.Address) {
	echohash := utils.Sha3(msg.Pack(), receiver[:])
	tr := &models.SentEnvelopMessager{
		Message:  msg,
		Receiver: receiver,
		Time:     time.Now(),
		EchoHash: echohash[:],
	}
	log.Trace(fmt.Sprintf("NewSentEnvelopMessager EchoHash=%s", utils.BPex(tr.EchoHash)))
	err := dao.saveKeyValueToBucket(models.BucketEnvelopMessager, tr.EchoHash, tr)
	if err != nil {
		log.Error(fmt.Sprintf("NewSentEnvelopMessager err=%s", err))
	}
}

//DeleteEnvelopMessager  delete a sending message from db
func (dao *GkvDB) DeleteEnvelopMessager(echoHash common.Hash) {
	err := dao.removeKeyValueFromBucket(models.BucketEnvelopMessager, echoHash[:])
	if err != nil {
		//可能这个消息完全不存在
		// this messsage might not exist.
		log.Warn(fmt.Sprintf("try to remove envelop message %s,but err= %s", utils.HPex(echoHash), err))
	}
}

//GetAllOrderedSentEnvelopMessager returns all EnvelopMessager message that have not receive ack and order them by nonce
func (dao *GkvDB) GetAllOrderedSentEnvelopMessager() []*models.SentEnvelopMessager {
	var msgs []*models.SentEnvelopMessager
	tb, err := dao.db.Table(models.BucketEnvelopMessager)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return msgs
	}
	for _, v := range buf {
		var s models.SentEnvelopMessager
		gobDecode(v, &s)
		msgs = append(msgs, &s)
	}
	//log.Trace(fmt.Sprintf("GetAllOrderedSentEnvelopMessager=%s", utils.StringInterface(msgs, 3)))
	models.SortEnvelopMessager(msgs)
	return msgs
}
