package gkvdb

import (
	"math/big"
	"time"

	"fmt"

	"gitee.com/johng/gkvdb/gkvdb"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

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
		TokenAddressBytes: tokenAddr[:],
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

//GetReceivedTransfer return the received transfer by key
func (dao *GkvDB) GetReceivedTransfer(key string) (*models.ReceivedTransfer, error) {
	var r models.ReceivedTransfer
	err := dao.getKeyValueToBucket(models.BucketReceivedTransfer, key, &r)
	err = models.GeneratDBError(err)
	return &r, err
}

//GetReceivedTransferList returns the received transfer between from and to blocks
func (dao *GkvDB) GetReceivedTransferList(tokenAddress common.Address, fromBlock, toBlock int64) (transfers []*models.ReceivedTransfer, err error) {
	var tb *gkvdb.Table
	tb, err = dao.db.Table(models.BucketReceivedTransfer)
	if err != nil {
		err = models.GeneratDBError(err)
		return
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var st models.ReceivedTransfer
		gobDecode(v, &st)
		appendReceivedTransferIfMatch(&transfers, &st, tokenAddress, fromBlock, toBlock)
	}
	return
}

func appendReceivedTransferIfMatch(list *[]*models.ReceivedTransfer, st *models.ReceivedTransfer, tokenAddress common.Address, fromBlock, toBlock int64) {
	var b1, b2, b3 bool
	if tokenAddress == utils.EmptyAddress || st.TokenAddress == tokenAddress {
		b1 = true
	}
	if fromBlock <= 0 || st.BlockNumber >= fromBlock {
		b2 = true
	}
	if toBlock <= 0 || st.BlockNumber <= toBlock {
		b3 = true
	}
	if b1 && b2 && b3 {
		*list = append(*list, st)
	}
}
