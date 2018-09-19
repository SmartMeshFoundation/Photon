package models

import (
	"encoding/gob"

	"github.com/asdine/storm"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
SentAnnounceDisposed 我发出了 AnnonuceDisposed, 那么就要标记这个 channel 上的这个锁我不要去链上兑现了,
如果对方发送过来 AnnounceDisposedResponse, 我要移除这个锁.
*/
/*
 *	SendAnnounceDisposed : a channel participant sends out an AnnounceDisposed message, then that lock in this channel should be tagged
 *	to inform that this participant should not register it on chain.
 *	If his partner sends AnnounceDisposedResponse, that lock has to be removed.
 */
type SentAnnounceDisposed struct {
	Key               []byte `storm:"id"`
	LockSecretHash    []byte `storm:"index"` //假设非恶意的情况下,锁肯定是不会重复的.但是我有可能在多个通道上发送 AnnounceDisposed,但是肯定不会在同一个通道上发送多次 announce disposed // Assume in honest case, locks are not repeated, but maybe I send AnnounceDisposed in multiple channels.
	ChannelIdentifier common.Hash
}

/*
ReceivedAnnounceDisposed 收到对方的 disposed, 主要是用来对方unlock 的时候,提交证据,惩罚对方
*/
/*
 *	ReceiveAnnounceDisposed : to receive AnnounceDisposed message from channel partner,
 *	mainly to submit proofs and punish fraudulent behaviors while partner submits unlock.
 */
type ReceivedAnnounceDisposed struct {
	Key               []byte `storm:"id"`
	LockHash          []byte `storm:"index"` //hash(expiration,locksecrethash,amount)
	ChannelIdentifier []byte `storm:"index"`
	OpenBlockNumber   int64
	AdditionalHash    common.Hash
	Signature         []byte
}

func init() {
	gob.Register(&SentAnnounceDisposed{})
	gob.Register(&ReceivedAnnounceDisposed{})
}

//MarkLockSecretHashDisposed mark `locksecrethash` disposed on channel `ChannelIdentifier`
func (model *ModelDB) MarkLockSecretHashDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) error {
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := model.db.Save(&SentAnnounceDisposed{
		Key:               key[:],
		LockSecretHash:    lockSecretHash[:],
		ChannelIdentifier: ChannelIdentifier,
	})
	return err
}

//IsLockSecretHashDisposed this lockSecretHash has Announced Disposed
func (model *ModelDB) IsLockSecretHashDisposed(lockSecretHash common.Hash) bool {
	sad := new(SentAnnounceDisposed)
	err := model.db.One("LockSecretHash", lockSecretHash[:], sad)
	if err != nil {
		return false
	}
	log.Trace(fmt.Sprintf("Find SentAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

//IsLockSecretHashChannelIdentifierDisposed `lockSecretHash` and `ChannelIdentifier` is the id of AnnounceDisposed
func (model *ModelDB) IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool {
	sad := new(SentAnnounceDisposed)
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := model.db.One("Key", key[:], sad)
	if err != nil {
		return false
	}
	log.Trace(fmt.Sprintf("Find SentAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

//NewReceivedAnnounceDisposed create ReceivedAnnounceDisposed
func NewReceivedAnnounceDisposed(LockHash, ChannelIdentifier, additionalHash common.Hash, openBlockNumber int64, signature []byte) *ReceivedAnnounceDisposed {
	key := utils.Sha3(LockHash[:], ChannelIdentifier[:])
	return &ReceivedAnnounceDisposed{
		Key:               key[:],
		LockHash:          LockHash[:],
		ChannelIdentifier: ChannelIdentifier[:],
		OpenBlockNumber:   openBlockNumber,
		AdditionalHash:    additionalHash,
		Signature:         signature,
	}
}

//MarkLockHashCanPunish 收到了一个放弃声明,需要保存,在收到 unlock 事件的时候进行 punish
/*
 *	MarkLockHashCanPunish : Once receiving an AnnounceDisposed message, we need to store it
 * 	and submit it to enforce punishment procedure while receiving unlock.
 */
func (model *ModelDB) MarkLockHashCanPunish(r *ReceivedAnnounceDisposed) error {
	return model.db.Save(r)
}

//IsLockHashCanPunish can punish this unlock?
func (model *ModelDB) IsLockHashCanPunish(lockHash, channelIdentifier common.Hash) bool {
	sad := new(ReceivedAnnounceDisposed)
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	err := model.db.One("Key", key[:], sad)
	if err != nil {
		return false
	}
	//log.Trace(fmt.Sprintf("Find ReceivedAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

//GetReceiviedAnnounceDisposed return a ReceivedAnnounceDisposed ,if not  exist,return nil
func (model *ModelDB) GetReceiviedAnnounceDisposed(lockHash, channelIdentifier common.Hash) *ReceivedAnnounceDisposed {
	sad := new(ReceivedAnnounceDisposed)
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	err := model.db.One("Key", key[:], sad)
	if err != nil {
		return nil
	}
	//log.Trace(fmt.Sprintf("Find ReceivedAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return sad
}

/*
GetChannelAnnounceDisposed 获取指定 channel中对方声明放弃的锁,
*/
/*
 *	GetChannelAnnounceDisposed : function to receive disposed locks claimed by channel partner in specific channel
 */
func (model *ModelDB) GetChannelAnnounceDisposed(channelIdentifier common.Hash) []*ReceivedAnnounceDisposed {
	var anns []*ReceivedAnnounceDisposed
	err := model.db.Find("ChannelIdentifier", channelIdentifier[:], &anns)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil
		}
		log.Error(fmt.Sprintf("GetChannelAnnounceDisposed for %s ,err %s", channelIdentifier.String(), err))
		return nil
	}
	return anns
}
