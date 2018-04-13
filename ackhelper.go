package smartraiden

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type AckHelper struct {
	db *models.ModelDB
}

func NewAckHelper(db *models.ModelDB) *AckHelper {
	return &AckHelper{db}
}
func (this *AckHelper) GetAck(echohash common.Hash) []byte {
	return this.db.GetAck(echohash)
}
func (this *AckHelper) SaveAck(echohash common.Hash, msg encoding.Messager, ack []byte) {
	data := this.GetAck(echohash)
	_, ok := msg.(*encoding.RevealSecret)
	_, ok2 := msg.(*encoding.SecretRequest)
	if ok || ok2 {
		if len(data) > 0 {
			log.Error(fmt.Sprintf("save ack for  RevealSecret which is already exist"))
		} else {
			tx := this.db.StartTx()
			this.db.SaveAck(echohash, ack, tx)
			tx.Commit()
		}

	} else {
		if len(data) == 0 {
			log.Error(fmt.Sprintf("save ack for non revealsecret which should be saved before,msg  is %s", msg))
		}
	}
}
