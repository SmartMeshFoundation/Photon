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
	_, ok := msg.(*encoding.RevealSecret)
	_, ok2 := msg.(*encoding.SecretRequest)
	if ok || ok2 {
		if len(data) > 0 {
			log.Error(fmt.Sprintf("save ack for  RevealSecret which is already exist"))
		} else {
			tx := ah.db.StartTx()
			ah.db.SaveAck(echohash, ack, tx)
			tx.Commit()
		}

	} else {
		if len(data) == 0 {
			log.Error(fmt.Sprintf("save ack for non revealsecret which should be saved before,msg  is %s", msg))
		}
	}
}
