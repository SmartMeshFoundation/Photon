package gkvdb

import (
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//MarkLockSecretHashDisposed mark `locksecrethash` disposed on channel `ChannelIdentifier`
func (dao *GkvDB) MarkLockSecretHashDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) error {
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	sad := &models.SentAnnounceDisposed{
		Key:               key[:],
		LockSecretHash:    lockSecretHash[:],
		ChannelIdentifier: ChannelIdentifier,
	}
	err := dao.saveKeyValueToBucket(models.BucketSentAnnounceDisposed, sad.Key, sad)
	return models.GeneratDBError(err)
}

//IsLockSecretHashDisposed this lockSecretHash has Announced Disposed
func (dao *GkvDB) IsLockSecretHashDisposed(lockSecretHash common.Hash) bool {
	var sads []*models.SentAnnounceDisposed
	tb, err := dao.db.Table(models.BucketSentAnnounceDisposed)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return false
	}
	for _, v := range buf {
		var sad models.SentAnnounceDisposed
		gobDecode(v, &sad)
		sads = append(sads, &sad)
		if common.BytesToHash(sad.LockSecretHash) == lockSecretHash {
			return true
		}
	}
	return false
}

//IsLockSecretHashChannelIdentifierDisposed `lockSecretHash` and `ChannelIdentifier` is the id of AnnounceDisposed
func (dao *GkvDB) IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool {
	sad := new(models.SentAnnounceDisposed)
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := dao.getKeyValueToBucket(models.BucketSentAnnounceDisposed, key[:], sad)
	if err != nil {
		return false
	}
	//log.Trace(fmt.Sprintf("Find SentAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

//MarkLockHashCanPunish 收到了一个放弃声明,需要保存,在收到 unlock 事件的时候进行 punish
/*
 *	MarkLockHashCanPunish : Once receiving an AnnounceDisposed message, we need to store it
 * 	and submit it to enforce punishment procedure while receiving unlock.
 */
func (dao *GkvDB) MarkLockHashCanPunish(r *models.ReceivedAnnounceDisposed) error {
	err := dao.saveKeyValueToBucket(models.BucketReceivedAnnounceDisposed, r.Key, r)
	return models.GeneratDBError(err)
}

//IsLockHashCanPunish can punish this unlock?
func (dao *GkvDB) IsLockHashCanPunish(lockHash, channelIdentifier common.Hash) bool {
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	tb, err := dao.db.Table(models.BucketReceivedAnnounceDisposed)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return false
	}
	for _, v := range buf {
		var rad models.ReceivedAnnounceDisposed
		gobDecode(v, &rad)
		if common.BytesToHash(rad.Key) == key {
			return true
		}
	}
	return false
}

//GetReceivedAnnounceDisposed return a ReceivedAnnounceDisposed ,if not  exist,return nil
func (dao *GkvDB) GetReceivedAnnounceDisposed(lockHash, channelIdentifier common.Hash) *models.ReceivedAnnounceDisposed {
	sad := new(models.ReceivedAnnounceDisposed)
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	err := dao.getKeyValueToBucket(models.BucketReceivedAnnounceDisposed, key[:], sad)
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
func (dao *GkvDB) GetChannelAnnounceDisposed(channelIdentifier common.Hash) (rads []*models.ReceivedAnnounceDisposed) {
	tb, err := dao.db.Table(models.BucketReceivedAnnounceDisposed)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var rad models.ReceivedAnnounceDisposed
		gobDecode(v, &rad)
		if common.BytesToHash(rad.ChannelIdentifier) == channelIdentifier {
			rads = append(rads, &rad)
		}
	}
	return
}
