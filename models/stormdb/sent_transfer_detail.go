package stormdb

import (
	"fmt"
	"math/big"

	"time"

	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
)

// NewSentTransferDetail :
func (model *StormDB) NewSentTransferDetail(tokenAddress, target common.Address, amount *big.Int, data string, isDirect bool, lockSecretHash common.Hash) {
	std := &models.SentTransferDetail{
		Key:               utils.Sha3(tokenAddress[:], lockSecretHash[:]).String(),
		BlockNumber:       model.GetLatestBlockNumber(),
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
	err := model.db.Save(std)
	if err != nil {
		log.Error(fmt.Sprintf("NewSendTransferDetail key=%s, err %s", std.Key, err))
		return
	}
	log.Trace(fmt.Sprintf("NewSendTransferDetail key=%s lockSecertHash=%s", std.Key, lockSecretHash.String()))
}

// UpdateSentTransferDetailStatus :
func (model *StormDB) UpdateSentTransferDetailStatus(tokenAddress common.Address, lockSecretHash common.Hash, status models.TransferStatusCode, statusMessage string, otherParams interface{}) {
	var std models.SentTransferDetail
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := model.db.One("Key", key, &std)
	if err == storm.ErrNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatus err %s", err))
		return
	}
	std.Status = status
	std.StatusMessage = fmt.Sprintf("%s%s\n", std.StatusMessage, statusMessage)
	if status == models.TransferStatusSuccess {
		if otherParams != nil {
			ch, ok := otherParams.(*channel.Channel)
			if ok {
				std.ChannelIdentifier = ch.ChannelIdentifier.ChannelIdentifier
				std.OpenBlockNumber = ch.ChannelIdentifier.OpenBlockNumber
			}
		}
		std.FinishTime = time.Now().Unix()
	}
	if status == models.TransferStatusCanceled || status == models.TransferStatusFailed {
		std.FinishTime = time.Now().Unix()
	}
	err = model.db.Save(&std)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatus err %s", err))
		return
	}
	log.Trace(fmt.Sprintf("UpdateStatus key=%s lockSecretHash=%s %s", key, lockSecretHash.String(), statusMessage))
}

// UpdateSentTransferDetailStatusMessage :
func (model *StormDB) UpdateSentTransferDetailStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string) {
	var ts models.SentTransferDetail
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := model.db.One("Key", key, &ts)
	if err == storm.ErrNotFound {
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatusMessage err %s", err))
		return
	}
	ts.StatusMessage = fmt.Sprintf("%s%s\n", ts.StatusMessage, statusMessage)
	err = model.db.Save(&ts)
	if err != nil {
		log.Error(fmt.Sprintf("UpdateStatusMessage err %s", err))
		return
	}
	log.Trace(fmt.Sprintf("UpdateStatusMessage key=%s lockSecretHash=%s %s", key, lockSecretHash.String(), statusMessage))
}

// GetSentTransferDetail :
func (model *StormDB) GetSentTransferDetail(tokenAddress common.Address, lockSecretHash common.Hash) (*models.SentTransferDetail, error) {
	var ts models.SentTransferDetail
	key := utils.Sha3(tokenAddress[:], lockSecretHash[:]).String()
	err := model.db.One("Key", key, &ts)
	log.Trace(fmt.Sprintf("GetSentTransferDetail key=%s lockSecretHash=%s err=%s", key, lockSecretHash.String(), err))
	err = models.GeneratDBError(err)
	return &ts, err
}

// GetSentTransferDetailList :
// 参数均为查询条件,传空值或负值代表不限制
func (model *StormDB) GetSentTransferDetailList(tokenAddress common.Address, fromTime, toTime int64, fromBlock, toBlock int64) (transfers []*models.SentTransferDetail, err error) {
	var selectList []q.Matcher
	if tokenAddress != utils.EmptyAddress {
		selectList = append(selectList, q.Eq("TokenAddressBytes", tokenAddress[:]))
	}
	if fromTime > 0 {
		selectList = append(selectList, q.Gte("SendingTime", fromTime))
	}
	if toTime > 0 {
		selectList = append(selectList, q.Lt("FinishTime", toTime))
	}
	if fromBlock > 0 {
		selectList = append(selectList, q.Gte("BlockNumber", fromBlock))
	}
	if toBlock > 0 {
		selectList = append(selectList, q.Lt("BlockNumber", toBlock))
	}
	if len(selectList) == 0 {
		err = model.db.All(&transfers)
	} else {
		q := model.db.Select(selectList...)
		err = q.Find(&transfers)
	}
	if err == storm.ErrNotFound {
		err = nil
	}
	return
}
