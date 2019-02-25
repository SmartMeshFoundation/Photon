package stormdb

import (
	"fmt"
	"time"

	"encoding/json"

	"strings"

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
	tokenAddress := utils.EmptyAddress
	if openBlockNumber == 0 && channelIdentifier != utils.EmptyHash {
		c, err2 := model.GetChannelByAddress(channelIdentifier)
		if err2 != nil {
			//log.Error(err2.Error())
		} else {
			openBlockNumber = c.ChannelIdentifier.OpenBlockNumber
			tokenAddress = c.TokenAddress()
		}
	}
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
		if p, ok := txParams.(*models.DepositTXParams); ok && tokenAddress == utils.EmptyAddress {
			tokenAddress = p.TokenAddress
		}
	}
	txInfo = &models.TXInfo{
		TXHash:            tx.Hash(),
		ChannelIdentifier: channelIdentifier,
		OpenBlockNumber:   openBlockNumber,
		TokenAddress:      tokenAddress,
		Type:              txType,
		IsSelfCall:        true,
		TXParams:          txParamsStr,
		Status:            models.TXInfoStatusPending,
		CallTime:          time.Now().Unix(),
		GasPrice:          tx.GasPrice().Uint64(),
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
func (model *StormDB) UpdateTXInfoStatus(txHash common.Hash, status models.TXInfoStatus, packBlockNumber int64, gasUsed uint64) (txInfo *models.TXInfo, err error) {
	var tis models.TXInfoSerialization
	err = model.db.One("TXHash", txHash[:], &tis)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTXInfoStatus err %s", err))
		err = models.GeneratDBError(err)
		return
	}
	tis.Status = string(status)
	tis.PackBlockNumber = packBlockNumber
	tis.PackTime = time.Now().Unix()
	tis.GasUsed = gasUsed
	if tis.OpenBlockNumber == 0 && tis.Type == models.TXInfoTypeDeposit {
		// 通道第一deposit,即通道打开,记录OpenBlockNumber和TokenAddress
		tis.OpenBlockNumber = packBlockNumber
		ch, err2 := model.GetChannelByAddress(common.BytesToHash(tis.ChannelIdentifier))
		if err2 != nil {
			//log.Error(err2.Error())
		} else {
			tis.TokenAddress = ch.TokenAddressBytes
		}
	}
	err = model.db.Save(&tis)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTXInfoStatus err %s", err))
		err = models.GeneratDBError(err)
		return
	}
	log.Trace(fmt.Sprintf("UpdateTXInfoStatus txhash=%s status=%s packBlockNumber=%d", txHash.String(), status, packBlockNumber))
	txInfo = tis.ToTXInfo()
	return
}

// GetTXInfoList :
// 如果参数不为空,则根据参数查询
func (model *StormDB) GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType models.TXInfoType, status models.TXInfoStatus) (list []*models.TXInfo, err error) {
	var selectList []q.Matcher
	if channelIdentifier != utils.EmptyHash {
		selectList = append(selectList, q.Eq("ChannelIdentifier", channelIdentifier[:]))
	}
	if openBlockNumber != 0 {
		selectList = append(selectList, q.Eq("OpenBlockNumber", openBlockNumber))
	}
	if tokenAddress != utils.EmptyAddress {
		selectList = append(selectList, q.Eq("TokenAddress", tokenAddress[:]))
	}
	if txType != "" {
		txTypeStr := string(txType)
		if strings.Contains(txTypeStr, ",") {
			ss := strings.Split(txTypeStr, ",")
			selectList = append(selectList, q.In("Type", ss))
		} else {
			selectList = append(selectList, q.Eq("Type", txType))
		}
	}
	if status != "" {
		txStatusStr := string(status)
		if strings.Contains(txStatusStr, ",") {
			ss := strings.Split(txStatusStr, ",")
			selectList = append(selectList, q.In("Status", ss))
		} else {
			selectList = append(selectList, q.Eq("Status", status))
		}
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
