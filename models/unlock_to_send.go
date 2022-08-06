package models

import "encoding/gob"

/*
UnlockToSend 保存应当发送但由于节点处于无有效公链状态所以没发送的Unlock消息,待切换到有效公链时,再进行发送操作
*/
type UnlockToSend struct {
	Key              []byte `storm:"id"` // 格式为utils.Sha3(LockSecretHash[:], Token[:], Receiver[:]).Bytes()
	LockSecretHash   []byte
	TokenAddress     []byte
	ReceiverAddress  []byte
	SavedTimestamp   int64 // 保存时的时间
	SavedBlockNumber int64 // 保存时的块号
}

func init() {
	gob.Register(&UnlockToSend{})
}
