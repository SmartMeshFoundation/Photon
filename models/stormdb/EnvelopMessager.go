package stormdb

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

//NewSentEnvelopMessager create a sending EnvelopMessager in db
func (model *StormDB) NewSentEnvelopMessager(msg encoding.EnvelopMessager, receiver common.Address) {
	echohash := utils.Sha3(msg.Pack(), receiver[:])
	tr := &models.SentEnvelopMessager{
		Message:  msg,
		Receiver: receiver,
		Time:     time.Now(),
		EchoHash: echohash[:],
	}
	log.Trace(fmt.Sprintf("NewSentEnvelopMessager EchoHash=%s", utils.BPex(tr.EchoHash)))
	err := model.db.Save(tr)
	if err != nil {
		log.Error(fmt.Sprintf("NewSentEnvelopMessager err=%s", err))
	}
}

//DeleteEnvelopMessager  delete a sending message from db
func (model *StormDB) DeleteEnvelopMessager(echohash common.Hash) {
	sss := &models.SentEnvelopMessager{
		EchoHash: echohash[:],
	}
	err := model.db.DeleteStruct(sss)
	if err != nil {
		//可能这个消息完全不存在
		// this messsage might not exist.
		log.Warn(fmt.Sprintf("try to remove envelop message %s,but err= %s", utils.HPex(echohash), err))
	}
}

//GetAllOrderedSentEnvelopMessager returns all EnvelopMessager message that have not receive ack and order them by nonce
func (model *StormDB) GetAllOrderedSentEnvelopMessager() []*models.SentEnvelopMessager {
	var msgs []*models.SentEnvelopMessager
	err := model.db.All(&msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllOrderedSentEnvelopMessager err=%s", err))
	}
	//log.Trace(fmt.Sprintf("GetAllOrderedSentEnvelopMessager=%s", utils.StringInterface(msgs, 3)))
	models.SortEnvelopMessager(msgs)
	return msgs
}
