package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
SentRemoveExpiredHashlockTransfer is record of sending a  RemoveExpiredHashlockTransfer
*/
type SentRemoveExpiredHashlockTransfer struct {
	EchoHash       common.Hash
	EchoHashString string `storm:"id"`
	Message        *encoding.RemoveExpiredHashlockTransfer
	Receiver       common.Address
	IsComplete     string `storm:"index"`
}

//IsSentRemoveExpiredHashlockTransferExist returns true when this message has been sent
func (model *ModelDB) IsSentRemoveExpiredHashlockTransferExist(echohash common.Hash) bool {
	var rss SentRemoveExpiredHashlockTransfer
	err := model.db.One("EchoHashString", echohash.String(), &rss)
	return err == nil
}

//NewSentRemoveExpiredHashlockTransfer create a sending RemoveExpiredHashlockTransfer in db
func (model *ModelDB) NewSentRemoveExpiredHashlockTransfer(msg *encoding.RemoveExpiredHashlockTransfer, receiver common.Address, tx storm.Node) {
	echohash := utils.Sha3(msg.Pack(), receiver[:])
	tr := &SentRemoveExpiredHashlockTransfer{
		EchoHash:       echohash,
		EchoHashString: echohash.String(),
		Message:        msg,
		Receiver:       receiver,
		IsComplete:     "false",
	}
	log.Trace(fmt.Sprintf("NewSentRemoveExpiredHashlockTransfer %s", utils.HPex(tr.EchoHash)))
	err := tx.Save(tr)
	if err != nil {
		log.Error(fmt.Sprintf("NewSentRemoveExpiredHashlockTransfer err=%s", err))
	}
}

//UpdateSentRemoveExpiredHashlockTransfer mark message sent complete
func (model *ModelDB) UpdateSentRemoveExpiredHashlockTransfer(echohash common.Hash) {
	var sss SentRemoveExpiredHashlockTransfer
	log.Trace(fmt.Sprintf("UpdateSentRemoveExpiredHashlockTransfer %s", utils.HPex(echohash)))
	err := model.db.One("EchoHashString", echohash.String(), &sss)
	if err != nil {
		panic("UpdateSentRemoveExpiredHashlockTransfer  must exist")
	}
	err = model.db.UpdateField(&sss, "IsComplete", "true")
	if err != nil {
		panic(fmt.Sprintf("UpdateSentRemoveExpiredHashlockTransfer err %s", err))
	}
}

//GetAllUncompleteSentRemoveExpiredHashlockTransfer returns all RemoveExpiredHashlockTransfer message that have not receive ack
func (model *ModelDB) GetAllUncompleteSentRemoveExpiredHashlockTransfer() []*SentRemoveExpiredHashlockTransfer {
	var msgs []*SentRemoveExpiredHashlockTransfer
	err := model.db.Find("IsComplete", "false", &msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllUncompleteSentRemoveExpiredHashlockTransfer err=%s", err))
	}
	log.Trace(fmt.Sprintf("GetAllUncompleteSentRemoveExpiredHashlockTransfer=%s", utils.StringInterface(msgs, 7)))
	return msgs
}
