package photon

import (
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/ethereum/go-ethereum/common"
)

//AckHelper save ack for sent and recevied  message
type AckHelper struct {
	dao models.Dao
}

//NewAckHelper create ack
func NewAckHelper(dao models.Dao) *AckHelper {
	return &AckHelper{dao}
}

//GetAck return a message's ack
func (ah *AckHelper) GetAck(echohash common.Hash) []byte {
	return ah.dao.GetAck(echohash)
}

//SaveAck save ack to dao
func (ah *AckHelper) SaveAck(echohash common.Hash, msg encoding.Messager, ack []byte) {
	ah.dao.SaveAckNoTx(echohash, ack)
}
