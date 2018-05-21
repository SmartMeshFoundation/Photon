package network

import (
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/ethereum/go-ethereum/common"
)

//ReceivedMessageSaver is designed for ignore duplicated message
type ReceivedMessageSaver interface {
	//GetAck return nil if not found,call this before message sent
	GetAck(echohash common.Hash) []byte
	//SaveAck  marks ack has been sent
	SaveAck(echohash common.Hash, msg encoding.Messager, ack []byte)
}
