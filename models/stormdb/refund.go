package stormdb

import (
	"github.com/asdine/storm"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//MarkLockSecretHashDisposed mark `locksecrethash` disposed on channel `ChannelIdentifier`
func (model *StormDB) MarkLockSecretHashDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) error {
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := model.db.Save(&models.SentAnnounceDisposed{
		Key:               key[:],
		LockSecretHash:    lockSecretHash[:],
		ChannelIdentifier: ChannelIdentifier,
	})
	err = models.GeneratDBError(err)
	return err
}

//IsLockSecretHashDisposed this lockSecretHash has Announced Disposed
func (model *StormDB) IsLockSecretHashDisposed(lockSecretHash common.Hash) bool {
	sad := new(models.SentAnnounceDisposed)
	err := model.db.One("LockSecretHash", lockSecretHash[:], sad)
	if err != nil {
		return false
	}
	//log.Trace(fmt.Sprintf("Find SentAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

//IsLockSecretHashChannelIdentifierDisposed `lockSecretHash` and `ChannelIdentifier` is the id of AnnounceDisposed
func (model *StormDB) IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool {
	sad := new(models.SentAnnounceDisposed)
	key := utils.Sha3(lockSecretHash[:], ChannelIdentifier[:])
	err := model.db.One("Key", key[:], sad)
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
func (model *StormDB) MarkLockHashCanPunish(r *models.ReceivedAnnounceDisposed) error {
	return model.db.Save(r)
}

//IsLockHashCanPunish can punish this unlock?
func (model *StormDB) IsLockHashCanPunish(lockHash, channelIdentifier common.Hash) bool {
	sad := new(models.ReceivedAnnounceDisposed)
	key := utils.Sha3(lockHash[:], channelIdentifier[:])
	err := model.db.One("Key", key[:], sad)
	if err != nil {
		return false
	}
	//log.Trace(fmt.Sprintf("Find ReceivedAnnounceDisposed=%s", utils.StringInterface(sad, 2)))
	return true
}

//GetReceivedAnnounceDisposed return a ReceivedAnnounceDisposed ,if not  exist,return nil
func (model *StormDB) GetReceivedAnnounceDisposed(lockHash, channelIdentifier common.Hash) *models.ReceivedAnnounceDisposed {
	sad := new(models.ReceivedAnnounceDisposed)
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
func (model *StormDB) GetChannelAnnounceDisposed(channelIdentifier common.Hash) []*models.ReceivedAnnounceDisposed {
	var anns []*models.ReceivedAnnounceDisposed
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

// MarkSubmittedByChannel 将一个通道上收到的所有punish都置为已委托
func (model *StormDB) MarkLockHashCanPunishSubmittedByChannel(channelIdentifier common.Hash) {
	list := model.GetChannelAnnounceDisposed(channelIdentifier)
	if list != nil && len(list) > 0 {
		for _, l := range list {
			err := model.db.UpdateField(l, "IsSubmittedToPms", true)
			if err != nil {
				log.Error(fmt.Sprintf("MarkSubmittedAnnounceDispose failed, channel=%s lockHash=%s err=%s",
					channelIdentifier.String(), common.BytesToHash(l.LockHash).String(), err.Error()))
			}
		}
	}
}
