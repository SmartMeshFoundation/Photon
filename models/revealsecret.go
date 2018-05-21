package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

//ReceivedRevealSecret represents a receiving reveal secret message.
type ReceivedRevealSecret struct {
	EchoHashString string `storm:"id"`
	EchoHash       common.Hash
	Message        *encoding.RevealSecret
	IsComplete     string `storm:"index"`
}

//SentRevealSecret represents a sending reveal secret message
type SentRevealSecret struct {
	EchoHash       common.Hash
	EchoHashString string `storm:"id"`
	Message        *encoding.RevealSecret
	Receiver       common.Address
	IsComplete     string `storm:"index"`
}

//NewReceivedRevealSecret create a ReceivedRevealSecret
func NewReceivedRevealSecret(msg *encoding.RevealSecret, echohash common.Hash) *ReceivedRevealSecret {
	return &ReceivedRevealSecret{
		EchoHashString: echohash.String(),
		EchoHash:       echohash,
		Message:        msg,
		IsComplete:     "false",
	}
}

//NewSentRevealSecret create a SentRevealSecret
func NewSentRevealSecret(msg *encoding.RevealSecret, receiver common.Address) *SentRevealSecret {
	echohash := utils.Sha3(msg.Pack(), receiver[:])
	return &SentRevealSecret{
		EchoHash:       echohash,
		EchoHashString: echohash.String(),
		Message:        msg,
		Receiver:       receiver,
		IsComplete:     "false",
	}
}

//IsReceivedRevealSecretExist return true when this message has received before.
func (model *ModelDB) IsReceivedRevealSecretExist(echohash common.Hash) bool {
	var rss ReceivedRevealSecret
	err := model.db.One("EchoHashString", echohash.String(), &rss)
	return err == nil
}

//NewReceivedRevealSecret marks receive a reveal secret message
func (model *ModelDB) NewReceivedRevealSecret(secret *ReceivedRevealSecret) {
	err := model.db.Save(secret)
	if err != nil {
		panic("ReceivedRevealSecret should not exist ")
	}
}

//UpdateReceivedRevealSecretComplete marks a revealsecret message has been processed.
func (model *ModelDB) UpdateReceivedRevealSecretComplete(echohash common.Hash) {
	var rss ReceivedRevealSecret
	err := model.db.One("EchoHashString", echohash.String(), &rss)
	if err != nil {
		panic("UpdateReceivedRevealSecretComplete revealsecret must exist")
	}
	err = model.db.UpdateField(&rss, "IsComplete", "true")
	if err != nil {
		panic(fmt.Sprintf("UpdateReceivedRevealSecretComplete err=%s", err))
	}
}

//GetAllUncompleteReceivedRevealSecret return all reveal secret messages that have not been processed before quit.
func (model *ModelDB) GetAllUncompleteReceivedRevealSecret() []*ReceivedRevealSecret {
	var msgs []*ReceivedRevealSecret
	err := model.db.Find("IsComplete", "false", &msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllUncompleteReceivedRevealSecret err=%s", err))
	}
	return msgs
}

//IsSentRevealSecretExist return true when this message can be found in db
func (model *ModelDB) IsSentRevealSecretExist(echohash common.Hash) bool {
	var rss SentRevealSecret
	err := model.db.One("EchoHashString", echohash.String(), &rss)
	return err == nil
}

/*
NewSentRevealSecret It is very likely to send reveal secret repeatedly, which can be ignored directly.
*/
func (model *ModelDB) NewSentRevealSecret(secret *SentRevealSecret) {
	log.Trace(fmt.Sprintf("NewSentRevealSecret %s", utils.HPex(secret.EchoHash)))
	err := model.db.Save(secret)
	if err != nil {
		log.Error(fmt.Sprintf("NewSentRevealSecret err=%s", err))
	}
}

//UpdateSentRevealSecretComplete marks message has been sent complete
func (model *ModelDB) UpdateSentRevealSecretComplete(echohash common.Hash) {
	var sss SentRevealSecret
	log.Trace(fmt.Sprintf("UpdateSentRevealSecretComplete %s", utils.HPex(echohash)))
	err := model.db.One("EchoHashString", echohash.String(), &sss)
	if err != nil {
		panic("UpdateSentRevealSecretComplete revealsecret must exist")
	}
	err = model.db.UpdateField(&sss, "IsComplete", "true")
	if err != nil {
		panic(fmt.Sprintf("UpdateSentRevealSecretComplete err %s", err))
	}
}

//GetAllUncompleteSentRevealSecret get all sending reveal secret messages that have not recevied ack
func (model *ModelDB) GetAllUncompleteSentRevealSecret() []*SentRevealSecret {
	var msgs []*SentRevealSecret
	err := model.db.Find("IsComplete", "false", &msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllUncompleteSentRevealSecret err=%s", err))
	}
	log.Trace(fmt.Sprintf("GetAllUncompleteSentRevealSecret=%s", utils.StringInterface(msgs, 7)))
	return msgs
}
