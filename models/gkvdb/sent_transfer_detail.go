package gkvdb

import (
	"fmt"
	"math/big"

	"time"

	"gitee.com/johng/gkvdb/gkvdb"
	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

// NewSentTransferDetail :
func (dao *GkvDB) NewSentTransferDetail(tokenAddress, target common.Address, amount *big.Int, data string, isDirect bool, lockSecretHash common.Hash) {
	std := &models.SentTransferDetail{
		Key:               utils.Sha3(tokenAddress[:], lockSecretHash[:]).String(),
		BlockNumber:       dao.GetLatestBlockNumber(),
		TokenAddress:      tokenAddress,
		TokenAddressBytes: tokenAddress[:],
		TargetAddress:     target,
		Amount:            amount,
		Data:              data,
		IsDirect:          isDirect,
		SendingTime:       time.Now().Unix(),
		FinishTime:        0,
		Status:            models.TransferStatusInit,
		StatusMessage:     "",
		ChannelIdentifier: utils.EmptyHash,
		OpenBlockNumber:   0,
	}
	err := dao.saveKeyValueToBucket(models.BucketSentTransferDetail, std.Key, std)
	if err != nil {
		log.Error(fmt.Sprintf("NewSendTransferDetail key=%s, err %s", std.Key, err))
		return
	}
	log.Trace(fmt.Sprintf("NewSendTransferDetail key=%s lockSecertHash=%s", std.Key, lockSecretHash.String()))
}

// UpdateSentTransferDetailStatus :
func (dao *GkvDB) UpdateSentTransferDetailStatus(tokenAddress common.Address, lockSecretHash common.Hash, status models.TransferStatusCode, statusMessage string, otherParams interface{}) {
	var std models.SentTransferDetail
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := dao.getKeyValueToBucket(models.BucketTransferStatus, key, &std)
	if err == ErrorNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatus err %s", err))
		return
	}
	std.Status = status
	std.StatusMessage = fmt.Sprintf("%s%s\n", std.StatusMessage, statusMessage)
	if status == models.TransferStatusSuccess && otherParams != nil {
		ch, ok := otherParams.(*channel.Channel)
		if ok {
			std.ChannelIdentifier = ch.ChannelIdentifier.ChannelIdentifier
			std.OpenBlockNumber = ch.ChannelIdentifier.OpenBlockNumber
		}
		std.FinishTime = time.Now().Unix()
	}
	if status == models.TransferStatusCanceled || status == models.TransferStatusFailed {
		std.FinishTime = time.Now().Unix()
	}
	err = dao.saveKeyValueToBucket(models.BucketSentTransferDetail, std.Key, &std)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatus err %s", err))
		return
	}
	log.Trace(fmt.Sprintf("UpdateStatus key=%s lockSecretHash=%s %s", key, lockSecretHash.String(), statusMessage))
}

// UpdateSentTransferDetailStatusMessage :
func (dao *GkvDB) UpdateSentTransferDetailStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string) {
	var std models.SentTransferDetail
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := dao.getKeyValueToBucket(models.BucketTransferStatus, key, &std)
	if err == storm.ErrNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatusMessage err %s", err))
		return
	}
	std.StatusMessage = fmt.Sprintf("%s%s\n", std.StatusMessage, statusMessage)
	err = dao.saveKeyValueToBucket(models.BucketSentTransferDetail, std.Key, &std)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatusMessage err %s", err))
		return
	}
	log.Trace(fmt.Sprintf("UpdateStatusMessage key=%s lockSecretHash=%s %s", key, lockSecretHash.String(), statusMessage))
}

// GetSentTransferDetail :
func (dao *GkvDB) GetSentTransferDetail(tokenAddress common.Address, lockSecretHash common.Hash) (*models.SentTransferDetail, error) {
	var std models.SentTransferDetail
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := dao.getKeyValueToBucket(models.BucketTransferStatus, key, &std)
	log.Trace(fmt.Sprintf("GetSentTransferDetail key=%s lockSecretHash=%s err=%s", key, lockSecretHash.String(), err))
	err = models.GeneratDBError(err)
	return &std, err
}

// GetSentTransferDetailList :
// 参数均为查询条件,传空值或负值代表不限制
func (dao *GkvDB) GetSentTransferDetailList(tokenAddress common.Address, fromTime, toTime int64, fromBlock, toBlock int64) (transfers []*models.SentTransferDetail, err error) {
	var tb *gkvdb.Table
	tb, err = dao.db.Table(models.BucketSentTransferDetail)
	if err != nil {
		err = models.GeneratDBError(err)
		return
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var st models.SentTransferDetail
		gobDecode(v, &st)
		appendSentTransferDetailIfMatch(&transfers, &st, tokenAddress, fromTime, toTime, fromBlock, toBlock)
	}
	return
}

func appendSentTransferDetailIfMatch(list *[]*models.SentTransferDetail, st *models.SentTransferDetail, tokenAddress common.Address, fromTime, toTime int64, fromBlock, toBlock int64) {
	var b1, b2, b3, b4, b5 bool
	if tokenAddress == utils.EmptyAddress || st.TokenAddress == tokenAddress {
		b1 = true
	}
	if fromTime <= 0 || st.SendingTime >= fromTime {
		b2 = true
	}
	if toTime <= 0 || st.FinishTime <= toTime {
		b3 = true
	}
	if fromBlock <= 0 || st.BlockNumber >= fromBlock {
		b4 = true
	}
	if toBlock <= 0 || st.BlockNumber <= toBlock {
		b5 = true
	}
	if b1 && b2 && b3 && b4 && b5 {
		*list = append(*list, st)
	}
}
