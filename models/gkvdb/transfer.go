package gkvdb

import (
	"math"
	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
NewSentTransfer save a new sent transfer to db,this transfer must be success
*/
func (dao *GkvDB) NewSentTransfer(blockNumber int64, channelIdentifier common.Hash, tokenAddr, toAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.SentTransfer {
	if lockSecretHash == utils.EmptyHash {
		// direct transfer, use fakeLockSecretHash
		lockSecretHash = utils.NewRandomHash()
	}
	key := fmt.Sprintf("%s-%d", channelIdentifier.String(), nonce)
	st := &models.SentTransfer{
		Key:               key,
		BlockNumber:       blockNumber,
		ChannelIdentifier: channelIdentifier,
		TokenAddress:      tokenAddr,
		ToAddress:         toAddr,
		Nonce:             nonce,
		Amount:            amount,
		Data:              data,
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
func (dao *GkvDB) NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.ReceivedTransfer {
	if lockSecretHash == utils.EmptyHash {
		// direct transfer, use fakeLockSecretHash
		lockSecretHash = utils.NewRandomHash()
	}
	key := fmt.Sprintf("%s-%d", channelIdentifier.String(), nonce)
	st := &models.ReceivedTransfer{
		Key:               key,
		BlockNumber:       blockNumber,
		ChannelIdentifier: channelIdentifier,
		TokenAddress:      tokenAddr,
		FromAddress:       fromAddr,
		Nonce:             nonce,
		Amount:            amount,
		Data:              data,
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
