package stormdb

import (
	"fmt"

	"encoding/json"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/kataras/go-errors"
)

// NewPendingTXInfo 创建pending状态的TXInfo,即自己发起的tx
func (model *StormDB) NewPendingTXInfo(tx *types.Transaction, txType models.TXInfoType, channelIdentifier common.Hash, openBlockNumber int64, txParams models.TXParams) (txInfo *models.TXInfo, err error) {
	var txParamsStr string
	if txParams != nil {
		if s, ok := txParams.(string); ok {
			txParamsStr = s
		} else {
			var buf []byte
			buf, err = json.Marshal(txParams)
			if err != nil {
				err = models.GeneratDBError(err)
				return
			}
			txParamsStr = string(buf)
		}
	}
	if openBlockNumber == 0 && channelIdentifier != utils.EmptyHash {
		c, err2 := model.GetChannelByAddress(channelIdentifier)
		if err2 != nil {
			log.Error(err.Error())
		} else {
			openBlockNumber = c.ChannelIdentifier.OpenBlockNumber
		}
	}
	txInfo = &models.TXInfo{
		TXHash:            tx.Hash(),
		ChannelIdentifier: channelIdentifier,
		OpenBlockNumber:   openBlockNumber,
		Type:              txType,
		IsSelfCall:        true,
		TXParams:          txParamsStr,
		Status:            models.TXInfoStatusPending,
	}
	err = model.db.Save(txInfo.ToTXInfoSerialization())
	if err != nil {
		log.Error(fmt.Sprintf("NewPendingTXInfo txhash=%s, err %s", txInfo.TXHash.String(), err))
		err = models.GeneratDBError(err)
		return
	}
	log.Trace(fmt.Sprintf("NewPendingTXInfo : \n%s", txInfo))
	return
}

// SaveEventToTXInfo 保存事件到TXInfo里面,当收到链上事件的时候调用
// 如果tx存在,保存事件到tx的事件列表里面
// 如果tx不存在,说明该tx非自己发起,直接创建success状态的tx并保存
// TODO
func (model *StormDB) SaveEventToTXInfo(event interface{}) (txInfo *models.TXInfo, err error) {
	//var txHash, channelIdentifier common.Hash
	//var openBlockNumber int64
	//var txType models.TXInfoType
	//txInfo := &models.TXInfo{
	//	TXHash:            txHash,
	//	ChannelIdentifier: channelIdentifier,
	//	OpenBlockNumber:   openBlockNumber,
	//	Type:              txType,
	//	IsSelfCall:        false,
	//	TXParams:          "",
	//	Status:            models.TXInfoStatusSuccess,
	//}
	return nil, errors.New("TODO")
}

// UpdateTXInfoStatus :
func (model *StormDB) UpdateTXInfoStatus(txHash common.Hash, status models.TXInfoStatus, pendingBlockNumber int64) (err error) {
	var tis models.TXInfoSerialization
	err = model.db.One("TXHash", txHash[:], &tis)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTXInfoStatus err %s", err))
		err = models.GeneratDBError(err)
		return
	}
	tis.Status = string(status)
	tis.PendingBlockNumber = pendingBlockNumber
	err = model.db.Save(&tis)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTXInfoStatus err %s", err))
		err = models.GeneratDBError(err)
		return
	}
	log.Trace(fmt.Sprintf("UpdateTXInfoStatus txhash=%s status=%s pendingBlockNumber=%d", txHash.String(), status, pendingBlockNumber))
	return
}

// GetTXInfoList :
// 如果参数不为空,则根据参数查询
func (model *StormDB) GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, txType models.TXInfoType, status models.TXInfoStatus) (list []*models.TXInfo, err error) {
	var selectList []q.Matcher
	if channelIdentifier != utils.EmptyHash {
		selectList = append(selectList, q.Eq("ChannelIdentifier", channelIdentifier[:]))
	}
	if openBlockNumber != 0 {
		selectList = append(selectList, q.Eq("OpenBlockNumber", openBlockNumber))
	}
	if txType != "" {
		selectList = append(selectList, q.Eq("Type", txType))
	}
	if status != "" {
		selectList = append(selectList, q.Eq("Status", status))
	}
	var l []*models.TXInfoSerialization
	if len(selectList) == 0 {
		err = model.db.All(&l)
	} else {
		q := model.db.Select(selectList...)
		err = q.Find(&l)
	}
	if err == storm.ErrNotFound {
		err = nil
		return
	}
	if err != nil {
		err = fmt.Errorf("GetTXInfoList err %s", err)
		err = models.GeneratDBError(err)
		return
	}
	for _, tis := range l {
		list = append(list, tis.ToTXInfo())
	}
	return
}
