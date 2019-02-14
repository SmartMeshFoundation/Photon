package gkvdb

import (
	"math"
	"math/big"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
NewSentTransfer save a new sent transfer to db,this transfer must be success
*/
func (dao *GkvDB) NewSentTransfer(blockNumber int64, channelIdentifier common.Hash, openBlockNumber int64, tokenAddr, toAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.SentTransfer {
	if lockSecretHash == utils.EmptyHash {
		// direct transfer, use fakeLockSecretHash
		lockSecretHash = utils.NewRandomHash()
	}
	key := fmt.Sprintf("%s-%d-%d", channelIdentifier.String(), openBlockNumber, nonce)
	st := &models.SentTransfer{
		Key:               key,
		BlockNumber:       blockNumber,
		ChannelIdentifier: channelIdentifier,
		TokenAddress:      tokenAddr,
		ToAddress:         toAddr,
		Nonce:             nonce,
		Amount:            amount,
		Data:              data,
		OpenBlockNumber:   openBlockNumber,
		TimeStamp:         time.Now().Format(time.RFC3339),
	}
	var ost models.SentTransfer
	err := dao.getKeyValueToBucket(models.BucketSentTransfer, key, &ost)
	if err == nil {
		log.Error(fmt.Sprintf("NewSentTransfer, but already exist, old=\n%s,new=\n%s",
			utils.StringInterface(ost, 2), utils.StringInterface(st, 2)))
		return nil
	}
	err = dao.saveKeyValueToBucket(models.BucketSentTransfer, key, st)
	if err != nil {
		log.Error(fmt.Sprintf("save SentTransfer err %s", err))
	}
	return st
}

//NewReceivedTransfer save a new received transfer to db
func (dao *GkvDB) NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, openBlockNumber int64, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.ReceivedTransfer {
	if lockSecretHash == utils.EmptyHash {
		// direct transfer, use fakeLockSecretHash
		lockSecretHash = utils.NewRandomHash()
	}
	key := fmt.Sprintf("%s-%d-%d", channelIdentifier.String(), openBlockNumber, nonce)
	st := &models.ReceivedTransfer{
		Key:               key,
		BlockNumber:       blockNumber,
		ChannelIdentifier: channelIdentifier,
		TokenAddress:      tokenAddr,
		FromAddress:       fromAddr,
		Nonce:             nonce,
		Amount:            amount,
		Data:              data,
		OpenBlockNumber:   openBlockNumber,
		TimeStamp:         time.Now().Format(time.RFC3339),
	}
	var ost models.ReceivedTransfer
	err := dao.getKeyValueToBucket(models.BucketReceivedTransfer, key, &ost)
	if err == nil {
		log.Error(fmt.Sprintf("NewReceivedTransfer, but already exist, old=\n%s,new=\n%s",
			utils.StringInterface(ost, 2), utils.StringInterface(st, 2)))
		return nil
	}
	err = dao.saveKeyValueToBucket(models.BucketReceivedTransfer, key, st)
	if err != nil {
		log.Error(fmt.Sprintf("save ReceivedTransfer err %s", err))
	}
	return st
}

//GetSentTransfer return the sent transfer by key
func (dao *GkvDB) GetSentTransfer(key string) (*models.SentTransfer, error) {
	var s models.SentTransfer
	err := dao.getKeyValueToBucket(models.BucketSentTransfer, key, &s)
	return &s, err
}

//GetReceivedTransfer return the received transfer by key
func (dao *GkvDB) GetReceivedTransfer(key string) (*models.ReceivedTransfer, error) {
	var r models.ReceivedTransfer
	err := dao.getKeyValueToBucket(models.BucketReceivedTransfer, key, &r)
	err = models.GeneratDBError(err)
	return &r, err
}

//GetSentTransferInBlockRange returns the sent transfer between from and to blocks
func (dao *GkvDB) GetSentTransferInBlockRange(fromBlock, toBlock int64) (transfers []*models.SentTransfer, err error) {
	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		toBlock = math.MaxInt64
	}
	tb, err := dao.db.Table(models.BucketSentTransfer)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil && len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var st models.SentTransfer
		gobDecode(v, &st)
		if st.BlockNumber >= fromBlock && st.BlockNumber <= toBlock {
			transfers = append(transfers, &st)
		}
	}
	return
}

//GetReceivedTransferInBlockRange returns the received transfer between from and to blocks
func (dao *GkvDB) GetReceivedTransferInBlockRange(fromBlock, toBlock int64) (transfers []*models.ReceivedTransfer, err error) {
	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		toBlock = math.MaxInt64
	}
	tb, err := dao.db.Table(models.BucketReceivedTransfer)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil && len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var rt models.ReceivedTransfer
		gobDecode(v, &rt)
		if rt.BlockNumber >= fromBlock && rt.BlockNumber <= toBlock {
			transfers = append(transfers, &rt)
		}
	}
	return
}

//GetSentTransferInTimeRange returns the sent transfer between from and to blocks
func (dao *GkvDB) GetSentTransferInTimeRange(from, to time.Time) (transfers []*models.SentTransfer, err error) {

	tb, err := dao.db.Table(models.BucketSentTransfer)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil && len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var st models.SentTransfer
		gobDecode(v, &st)
		t, err := time.Parse(time.RFC3339, st.TimeStamp)
		if err != nil {
			log.Warn(fmt.Sprintf("time parse for %s err %s", utils.StringInterface(st, 3), err))
			continue
		}
		log.Trace(fmt.Sprintf("from=%d,to=%d,t=%d", from.Unix(), to.Unix(), t.Unix()))
		if t.After(from) && t.Before(to) {
			transfers = append(transfers, &st)
		}
	}
	return
}

//GetReceivedTransferInTimeRange returns the received transfer between from and to blocks
func (dao *GkvDB) GetReceivedTransferInTimeRange(from, to time.Time) (transfers []*models.ReceivedTransfer, err error) {

	tb, err := dao.db.Table(models.BucketReceivedTransfer)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil && len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var rt models.ReceivedTransfer
		gobDecode(v, &rt)
		t, err := time.Parse(time.RFC3339, rt.TimeStamp)
		if err != nil {
			log.Warn(fmt.Sprintf("time parse for %s err %s", utils.StringInterface(rt, 3), err))
			continue
		}
		if t.After(from) && t.Before(to) {
			transfers = append(transfers, &rt)
		}
	}
	return
}
