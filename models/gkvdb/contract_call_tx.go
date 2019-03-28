package gkvdb

import (
	"encoding/json"
	"errors"
	"fmt"

	"bytes"

	"time"

	"gitee.com/johng/gkvdb/gkvdb"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewPendingTXInfo 创建pending状态的TXInfo,即自己发起的tx
func (dao *GkvDB) NewPendingTXInfo(tx *types.Transaction, txType models.TXInfoType, channelIdentifier common.Hash, openBlockNumber int64, txParams models.TXParams) (txInfo *models.TXInfo, err error) {
	tokenAddress := utils.EmptyAddress
	if openBlockNumber == 0 && channelIdentifier != utils.EmptyHash {
		c, err2 := dao.GetChannelByAddress(channelIdentifier)
		if err2 != nil {
			log.Error(err2.Error())
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
	tis := txInfo.ToTXInfoSerialization()
	err = dao.saveKeyValueToBucket(models.BucketTXInfo, tis.TXHash, tis)
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
func (dao *GkvDB) SaveEventToTXInfo(event interface{}) (txInfo *models.TXInfo, err error) {
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
func (dao *GkvDB) UpdateTXInfoStatus(txHash common.Hash, status models.TXInfoStatus, pendingBlockNumber int64, gasUsed uint64) (txInfo *models.TXInfo, err error) {
	var tis models.TXInfoSerialization
	err = dao.getKeyValueToBucket(models.BucketTXInfo, txHash[:], &tis)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTXInfoStatus err %s", err))
		err = models.GeneratDBError(err)
		return
	}
	tis.Status = string(status)
	tis.PackBlockNumber = pendingBlockNumber
	tis.PackTime = time.Now().Unix()
	tis.GasUsed = gasUsed
	err = dao.saveKeyValueToBucket(models.BucketTXInfo, tis.TXHash, tis)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateTXInfoStatus err %s", err))
		err = models.GeneratDBError(err)
		return
	}
	log.Trace(fmt.Sprintf("UpdateTXInfoStatus txhash=%s status=%s pendingBlockNumber=%d", txHash.String(), status, pendingBlockNumber))
	txInfo = tis.ToTXInfo()
	return
}

// GetTXInfoList :
// 如果参数不为空,则根据参数查询
func (dao *GkvDB) GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType models.TXInfoType, status models.TXInfoStatus) (list []*models.TXInfo, err error) {
	var tb *gkvdb.Table
	tb, err = dao.db.Table(models.BucketTXInfo)
	if err != nil {
		err = models.GeneratDBError(err)
		return
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var tis models.TXInfoSerialization
		gobDecode(v, &tis)
		appendTXInfoIfMatch(&list, &tis, channelIdentifier, openBlockNumber, tokenAddress, txType, status)
	}
	return
}

func appendTXInfoIfMatch(list *[]*models.TXInfo, tis *models.TXInfoSerialization, channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType models.TXInfoType, status models.TXInfoStatus) {
	var b1, b2, b3, b4, b5 bool
	if channelIdentifier == utils.EmptyHash || bytes.Compare(tis.ChannelIdentifier, channelIdentifier[:]) == 0 {
		b1 = true
	}
	if openBlockNumber <= 0 || tis.OpenBlockNumber == openBlockNumber {
		b2 = true
	}
	if txType == "" || tis.Type == string(txType) {
		b3 = true
	}
	if status == "" || tis.Status == string(status) {
		b4 = true
	}
	if tokenAddress == utils.EmptyAddress || bytes.Compare(tis.TokenAddress, tokenAddress[:]) == 0 {
		b5 = true
	}
	if b1 && b2 && b3 && b4 && b5 {
		*list = append(*list, tis.ToTXInfo())
	}
}
