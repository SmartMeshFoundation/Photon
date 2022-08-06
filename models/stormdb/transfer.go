package stormdb

import (
	"math/big"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
)

//NewReceivedTransfer save a new received transfer to db
func (model *StormDB) NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, openBlockNumber int64, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.ReceivedTransfer {
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
		TimeStamp:         time.Now().Unix(),
	}
	if ost, err := model.GetReceivedTransfer(key); err == nil {
		log.Error(fmt.Sprintf("NewReceivedTransfer, but already exist, old=\n%s,new=\n%s",
			utils.StringInterface(ost, 2), utils.StringInterface(st, 2)))
		return nil
	}
	err := model.db.Save(st)
	if err != nil {
		log.Error(fmt.Sprintf("save ReceivedTransfer err %s", err))
	}
	return st
}

//GetReceivedTransfer return the received transfer by key
func (model *StormDB) GetReceivedTransfer(key string) (*models.ReceivedTransfer, error) {
	var r models.ReceivedTransfer
	err := model.db.One("Key", key, &r)
	err = models.GeneratDBError(err)
	return &r, err
}

//GetReceivedTransferList returns the received transfer between from and to blocks
func (model *StormDB) GetReceivedTransferList(tokenAddress common.Address, fromBlock, toBlock, fromTime, toTime int64) (transfers []*models.ReceivedTransfer, err error) {
	var selectList []q.Matcher
	if tokenAddress != utils.EmptyAddress {
		selectList = append(selectList, q.Eq("TokenAddressBytes", tokenAddress[:]))
	}
	if fromBlock > 0 {
		selectList = append(selectList, q.Gte("BlockNumber", fromBlock))
	}
	if toBlock > 0 {
		selectList = append(selectList, q.Lt("BlockNumber", toBlock))
	}
	if fromTime > 0 {
		selectList = append(selectList, q.Gte("TimeStamp", fromTime))
	}
	if toTime > 0 {
		selectList = append(selectList, q.Lt("TimeStamp", toTime))
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
	if err != nil {
		err = rerr.ErrGeneralDBError
	}
	return
}
