package models

import (
	"fmt"

	"time"

	"sort"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

/*
SentEnvelopMessager is record of envelop message,that don't received a Ack
*/
type SentEnvelopMessager struct {
	Message  encoding.EnvelopMessager
	Receiver common.Address
	Time     time.Time
	EchoHash []byte `storm:"id"`
}

//NewSentEnvelopMessager create a sending EnvelopMessager in db
func (model *ModelDB) NewSentEnvelopMessager(msg encoding.EnvelopMessager, receiver common.Address) {
	echohash := utils.Sha3(msg.Pack(), receiver[:])
	tr := &SentEnvelopMessager{
		Message:  msg,
		Receiver: receiver,
		Time:     time.Now(),
		EchoHash: echohash[:],
	}
	log.Trace(fmt.Sprintf("NewSentEnvelopMessager %s", utils.BPex(tr.EchoHash)))
	err := model.db.Save(tr)
	if err != nil {
		log.Error(fmt.Sprintf("NewSentEnvelopMessager err=%s", err))
	}
}

//DeleteEnvelopMessager  delete a sending message from db
func (model *ModelDB) DeleteEnvelopMessager(echohash common.Hash) {
	sss := &SentEnvelopMessager{
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
func (model *ModelDB) GetAllOrderedSentEnvelopMessager() []*SentEnvelopMessager {
	var msgs []*SentEnvelopMessager
	err := model.db.All(&msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllOrderedSentEnvelopMessager err=%s", err))
	}
	//log.Trace(fmt.Sprintf("GetAllOrderedSentEnvelopMessager=%s", utils.StringInterface(msgs, 3)))
	sortEnvelopMessager(msgs)
	return msgs
}

type envelopMessageSorter []*SentEnvelopMessager

func (c envelopMessageSorter) Len() int {
	return len(c)
}
func (c envelopMessageSorter) Less(i, j int) bool {
	m1 := c[i].Message.GetEnvelopMessage()
	m2 := c[j].Message.GetEnvelopMessage()
	return m1.Nonce < m2.Nonce
}
func (c envelopMessageSorter) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

/*
must be stable
对于 ChannelOpenedAndDeposit 事件,会产生两个 stateChange,
严格要求有先后顺序
*/
/*
 *	sortEnvelopMessager : function to sort arrays of sent messenger.
 *
 *	Note that for event of ChannelOpenedAndDeposit, two stateChange will be generated.
 *	And they must be in order.
 */
func sortEnvelopMessager(msgs []*SentEnvelopMessager) {
	sort.Stable(envelopMessageSorter(msgs))
}
