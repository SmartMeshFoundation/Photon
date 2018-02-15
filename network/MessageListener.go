package network

import (
	"github.com/SmartMeshFoundation/raiden-network/encoding"
	"github.com/ethereum/go-ethereum/common"
)

type ReceivedMessageSaver interface {
	//call this before message sent
	GetAck(echohash common.Hash) []byte
	SaveAck(echohash common.Hash, msg encoding.Messager, ack []byte)
}
