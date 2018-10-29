package photon

import (
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/models"
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
	ah.db.SaveAckNoTx(echohash, ack)
}
