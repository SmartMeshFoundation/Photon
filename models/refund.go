package models

import (
	"encoding/gob"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
我发出了 AnnonuceDisposed, 那么就要标记这个 channel 上的这个锁我不要去链上兑现了,
如果对方发送过来 AnnounceDisposedResponse, 我要移除这个锁.
*/
type SentAnnounceDisposed struct {
	Key               string      `storm:"id"`
	LockSecretHash    common.Hash `storm:"index"` //假设非恶意的情况下,锁肯定是不会重复的.但是我有可能在多个通道上发送 AnnounceDisposed
	ChannelIdentifier common.Hash
}

/*
收到对方的 disposed, 主要是用来对方unlock 的时候,提交证据,惩罚对方
*/
type ReceivedAnnounceDisposed struct {
	Key               common.Hash `storm:"id"`
	LockHash          common.Hash `storm:"index"` //hash(expiration,locksecrethash,amount)
	ChannelIdentifier common.Hash
	OpenBlockNumber   int64
	AdditionalHash    common.Hash
	Signature         []byte
}

func init() {
	gob.Register(&SentAnnounceDisposed{})
	gob.Register(&ReceivedAnnounceDisposed{})
}
func (model *ModelDB) MarkLockSecretHashDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) error {
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := model.db.Save(&SentAnnounceDisposed{
		Key:               key.String(),
		LockSecretHash:    lockSecretHash,
		ChannelIdentifier: ChannelIdentifier,
	})
	return err
}

func (model *ModelDB) IsLockSecretHashDisposed(lockSecretHash common.Hash) bool {
	sad := new(SentAnnounceDisposed)
	err := model.db.One("LockSecretHash", lockSecretHash, sad)
	if err != nil {
		return false
	}
	log.Trace(fmt.Sprintf("Find SentAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}
func (model *ModelDB) IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool {
	sad := new(SentAnnounceDisposed)
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := model.db.One("Key", key.String(), sad)
	if err != nil {
		return false
	}
	log.Trace(fmt.Sprintf("Find SentAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

func NewReceivedAnnounceDisposed(LockHash, ChannelIdentifier, additionalHash common.Hash, openBlockNumber int64, signature []byte) *ReceivedAnnounceDisposed {
	return &ReceivedAnnounceDisposed{
		Key:               utils.Sha3(LockHash[:], ChannelIdentifier[:]),
		LockHash:          LockHash,
		ChannelIdentifier: ChannelIdentifier,
		OpenBlockNumber:   openBlockNumber,
		AdditionalHash:    additionalHash,
		Signature:         signature,
	}
}

//收到了一个放弃声明,需要保存,在收到 unlock 事件的时候进行 punish
func (model *ModelDB) MarkLockHashCanPunish(r *ReceivedAnnounceDisposed) error {
	return model.db.Save(r)
}
func (model *ModelDB) IsLockHashCanPunish(lockHash, channelIdentifier common.Hash) bool {
	sad := new(ReceivedAnnounceDisposed)
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	err := model.db.One("Key", key, sad)
	if err != nil {
		return false
	}
	//log.Trace(fmt.Sprintf("Find ReceivedAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

func (model *ModelDB) GetReceiviedAnnounceDisposed(lockHash, channelIdentifier common.Hash) *ReceivedAnnounceDisposed {
	sad := new(ReceivedAnnounceDisposed)
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	err := model.db.One("Key", key, sad)
	if err != nil {
		return nil
	}
	//log.Trace(fmt.Sprintf("Find ReceivedAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return sad
}
