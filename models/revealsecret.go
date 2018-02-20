package models

import (
	"fmt"

	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ReceivedRevealSecret struct {
	EchoHashString string `storm:"id"`
	EchoHash       common.Hash
	Message        *encoding.RevealSecret
	IsComplete     string `storm:"index"`
}

type SentRevealSecret struct {
	EchoHash       common.Hash
	EchoHashString string `storm:"id"`
	Message        *encoding.RevealSecret
	Receiver       common.Address
	IsComplete     string `storm:"index"`
}

func NewReceivedRevealSecret(msg *encoding.RevealSecret, echohash common.Hash) *ReceivedRevealSecret {
	return &ReceivedRevealSecret{
		EchoHashString: echohash.String(),
		EchoHash:       echohash,
		Message:        msg,
		IsComplete:     "false",
	}
}
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
func (model *ModelDB) IsReceivedRevealSecretExist(echohash common.Hash) bool {
	var rss ReceivedRevealSecret
	err := model.db.One("EchoHashString", echohash.String(), &rss)
	return err == nil
}
func (model *ModelDB) NewReceivedRevealSecret(secret *ReceivedRevealSecret) {
	err := model.db.Save(secret)
	if err != nil {
		panic("ReceivedRevealSecret should not exist ")
	}
}
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

func (model *ModelDB) GetAllUncompleteReceivedRevealSecret() []*ReceivedRevealSecret {
	var msgs []*ReceivedRevealSecret
	err := model.db.Find("IsComplete", "false", &msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllUncompleteReceivedRevealSecret err=%s", err))
	}
	return msgs
}

func (model *ModelDB) IsSentRevealSecretExist(echohash common.Hash) bool {
	var rss SentRevealSecret
	err := model.db.One("EchoHashString", echohash.String(), &rss)
	return err == nil
}

/*
很有可能重复发送reveal secret,直接忽略即可,
*/
func (model *ModelDB) NewSentRevealSecret(secret *SentRevealSecret) {
	log.Trace(fmt.Sprintf("NewSentRevealSecret %s", utils.HPex(secret.EchoHash)))
	err := model.db.Save(secret)
	if err != nil {
		log.Error(fmt.Sprintf("NewSentRevealSecret err=%s", err))
	}
}
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

func (model *ModelDB) GetAllUncompleteSentRevealSecret() []*SentRevealSecret {
	var msgs []*SentRevealSecret
	err := model.db.Find("IsComplete", "false", &msgs)
	if err != nil && err != storm.ErrNotFound {
		panic(fmt.Sprintf("GetAllUncompleteSentRevealSecret err=%s", err))
	}
	log.Trace(fmt.Sprintf("GetAllUncompleteSentRevealSecret=%s", utils.StringInterface(msgs, 7)))
	return msgs
}
