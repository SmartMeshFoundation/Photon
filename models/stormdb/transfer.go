package stormdb

import (
	"math"
	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

/*
NewSentTransfer save a new sent transfer to db,this transfer must be success
*/
func (model *StormDB) NewSentTransfer(blockNumber int64, channelIdentifier common.Hash, tokenAddr, toAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.SentTransfer {
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
	if ost, err := model.GetSentTransfer(key); err == nil {
		log.Error(fmt.Sprintf("NewSentTransfer, but already exist, old=\n%s,new=\n%s",
			utils.StringInterface(ost, 2), utils.StringInterface(st, 2)))
		return nil
	}
	err := model.db.Save(st)
	if err != nil {
		log.Error(fmt.Sprintf("save SentTransfer err %s", err))
	}
	return st
}

//NewReceivedTransfer save a new received transfer to db
func (model *StormDB) NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *models.ReceivedTransfer {
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

//GetSentTransfer return the sent transfer by key
func (model *StormDB) GetSentTransfer(key string) (*models.SentTransfer, error) {
	var s models.SentTransfer
	err := model.db.One("Key", key, &s)
	return &s, err
}

//GetReceivedTransfer return the received transfer by key
func (model *StormDB) GetReceivedTransfer(key string) (*models.ReceivedTransfer, error) {
	var r models.ReceivedTransfer
	err := model.db.One("Key", key, &r)
	return &r, err
}

//GetSentTransferInBlockRange returns the sent transfer between from and to blocks
func (model *StormDB) GetSentTransferInBlockRange(fromBlock, toBlock int64) (transfers []*models.SentTransfer, err error) {
	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		toBlock = math.MaxInt64
	}
	err = model.db.Range("BlockNumber", fromBlock, toBlock, &transfers)
	if err == storm.ErrNotFound { //ingore not found error
		err = nil
	}
	return
}

//GetReceivedTransferInBlockRange returns the received transfer between from and to blocks
func (model *StormDB) GetReceivedTransferInBlockRange(fromBlock, toBlock int64) (transfers []*models.ReceivedTransfer, err error) {
	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		toBlock = math.MaxInt64
	}
	err = model.db.Range("BlockNumber", fromBlock, toBlock, &transfers)
	if err == storm.ErrNotFound { //ingore not found error
		err = nil
	}
	return
}
