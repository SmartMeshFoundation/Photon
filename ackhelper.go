package smartraiden

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/ethereum/go-ethereum/common"
)

//AckHelper save ack for sent and recevied  message
type AckHelper struct {
	db *models.ModelDB
}

//NewAckHelper create ack
func NewAckHelper(db *models.ModelDB) *AckHelper {
	return &AckHelper{db}
}

//GetAck return a message's ack
func (ah *AckHelper) GetAck(echohash common.Hash) []byte {
	return ah.db.GetAck(echohash)
}

//SaveAck save ack to db
func (ah *AckHelper) SaveAck(echohash common.Hash, msg encoding.Messager, ack []byte) {
	data := ah.GetAck(echohash)
	var ok bool
	switch msg.(type) {
	case *encoding.RevealSecret:
		ok = true
	case *encoding.SecretRequest:
		ok = true
	case *encoding.DirectTransfer:
		ok = true
	case *encoding.AnnounceDisposed:
		ok = true
	case *encoding.RemoveExpiredHashlockTransfer:
		ok = true
	}
	if ok {
		if len(data) > 0 {
			log.Error(fmt.Sprintf("save ack for  %s which is already exist", msg.String()))
		} else {
			tx := ah.db.StartTx()
			ah.db.SaveAck(echohash, ack, tx)
			err := tx.Commit()
			if err != nil {
				log.Error(fmt.Sprintf("SaveAck err %s", err))
			}
		}

	} else {
		if len(data) == 0 {
			log.Error(fmt.Sprintf("save ack for non revealsecret which should be saved before,msg  is %s", msg))
		}
	}
}
