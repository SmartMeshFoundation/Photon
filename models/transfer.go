package models

import (
	"math/big"

	"fmt"

	"math"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

//SentTransfer transfer's I have sent and success.
type SentTransfer struct {
	Key            string         `storm:"id"`
	BlockNumber    int64          `json:"block_number" storm:"index"`
	ChannelAddress common.Address `json:"channel_address"`
	ToAddress      common.Address `json:"to_address"`
	TokenAddress   common.Address `json:"token_address"`
	Nonce          int64          `json:"nonce"`
	Amount         *big.Int       `json:"amount"`
}

//ReceivedTransfer tokens I have received and where it comes from
type ReceivedTransfer struct {
	Key            string         `storm:"id"`
	BlockNumber    int64          `json:"block_number" storm:"index"`
	ChannelAddress common.Address `json:"channel_address"`
	TokenAddress   common.Address `json:"token_address"`
	FromAddress    common.Address `json:"from_address"`
	Nonce          int64          `json:"nonce"`
	Amount         *big.Int       `json:"amount"`
}

/*
NewSentTransfer save a new sent transfer to db,this trqnsfer must be success
*/
func (model *ModelDB) NewSentTransfer(blockNumber int64, channelAddr, tokenAddr, toAddr common.Address, nonce int64, amount *big.Int) {
	key := fmt.Sprintf("%s-%d", channelAddr.String(), nonce)
	st := &SentTransfer{
		Key:            key,
		BlockNumber:    blockNumber,
		ChannelAddress: channelAddr,
		TokenAddress:   tokenAddr,
		ToAddress:      toAddr,
		Nonce:          nonce,
		Amount:         amount,
	}
	if ost, err := model.GetSentTransfer(key); err == nil {
		log.Error(fmt.Sprintf("NewSentTransfer, but already exist, old=\n%s,new=\n%s",
			utils.StringInterface(ost, 2), utils.StringInterface(st, 2)))
		return
	}
	err := model.db.Save(st)
	if err != nil {
		log.Error(fmt.Sprintf("save SentTransfer err %s", err))
	}
	select {
	case model.SentTransferChan <- st:
	default:
		//nerver block
	}
}

//NewReceivedTransfer save a new received transfer to db
func (model *ModelDB) NewReceivedTransfer(blockNumber int64, channelAddr, tokenAddr, fromAddr common.Address, nonce int64, amount *big.Int) {
	key := fmt.Sprintf("%s-%d", channelAddr.String(), nonce)
	st := &ReceivedTransfer{
		Key:            key,
		BlockNumber:    blockNumber,
		ChannelAddress: channelAddr,
		TokenAddress:   tokenAddr,
		FromAddress:    fromAddr,
		Nonce:          nonce,
		Amount:         amount,
	}
	if ost, err := model.GetReceivedTransfer(key); err == nil {
		log.Error(fmt.Sprintf("NewReceivedTransfer, but already exist, old=\n%s,new=\n%s",
			utils.StringInterface(ost, 2), utils.StringInterface(st, 2)))
		return
	}
	err := model.db.Save(st)
	if err != nil {
		log.Error(fmt.Sprintf("save ReceivedTransfer err %s", err))
	}
	select {
	case model.ReceivedTransferChan <- st:
	default:
		//never block
	}
}

//GetSentTransfer return the sent transfer by key
func (model *ModelDB) GetSentTransfer(key string) (*SentTransfer, error) {
	var s SentTransfer
	err := model.db.One("Key", key, &s)
	return &s, err
}

//GetReceivedTransfer return the received transfer by key
func (model *ModelDB) GetReceivedTransfer(key string) (*ReceivedTransfer, error) {
	var r ReceivedTransfer
	err := model.db.One("Key", key, &r)
	return &r, err
}

//GetSentTransferInBlockRange returns the sent transfer between from and to blocks
func (model *ModelDB) GetSentTransferInBlockRange(fromBlock, toBlock int64) (transfers []*SentTransfer, err error) {
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
func (model *ModelDB) GetReceivedTransferInBlockRange(fromBlock, toBlock int64) (transfers []*ReceivedTransfer, err error) {
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
