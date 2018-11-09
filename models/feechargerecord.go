package models

import (
	"fmt"
	"math/big"

	"encoding/json"

	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// FeeChargerRecordSerialization :
// 记录中间节点收取手续费的流水
type FeeChargerRecordSerialization struct {
	Key            []byte `storm:"id"`
	LockSecretHash []byte `storm:"index"`
	TokenAddress   []byte
	TransferFrom   []byte
	TransferTo     []byte
	TransferAmount *big.Int
	InChannel      []byte
	OutChannel     []byte
	Fee            *big.Int
	Timestamp      int64
}

// ToFeeChargeRecord :
func (rs *FeeChargerRecordSerialization) ToFeeChargeRecord() *FeeChargeRecord {
	return &FeeChargeRecord{
		Key:            common.BytesToHash(rs.Key),
		LockSecretHash: common.BytesToHash(rs.LockSecretHash),
		TokenAddress:   common.BytesToAddress(rs.TokenAddress),
		TransferFrom:   common.BytesToAddress(rs.TransferFrom),
		TransferTo:     common.BytesToAddress(rs.TransferTo),
		TransferAmount: rs.TransferAmount,
		InChannel:      common.BytesToHash(rs.InChannel),
		OutChannel:     common.BytesToHash(rs.OutChannel),
		Fee:            rs.Fee,
		Timestamp:      rs.Timestamp,
	}
}

// FeeChargeRecord :
// 记录中间节点收取手续费的流水
type FeeChargeRecord struct {
	Key            common.Hash    `json:"key" storm:"id"`
	LockSecretHash common.Hash    `json:"lock_secret_hash"`
	TokenAddress   common.Address `json:"token_address"`
	TransferFrom   common.Address `json:"transfer_from"`
	TransferTo     common.Address `json:"transfer_to"`
	TransferAmount *big.Int       `json:"transfer_amount"`
	InChannel      common.Hash    `json:"in_channel"`  // 我收款的channelID
	OutChannel     common.Hash    `json:"out_channel"` // 我付款的channelID
	Fee            *big.Int       `json:"fee"`
	Timestamp      int64          `json:"timestamp"` // 时间戳,time.Unix()
}

// ToString :
func (r *FeeChargeRecord) ToString() string {
	buf, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// ToSerialized :
func (r *FeeChargeRecord) ToSerialized() *FeeChargerRecordSerialization {
	return &FeeChargerRecordSerialization{
		Key:            r.Key[:],
		LockSecretHash: r.LockSecretHash[:],
		TokenAddress:   r.TokenAddress[:],
		TransferFrom:   r.TransferFrom[:],
		TransferTo:     r.TransferTo[:],
		TransferAmount: r.TransferAmount,
		InChannel:      r.InChannel[:],
		OutChannel:     r.OutChannel[:],
		Fee:            r.Fee,
		Timestamp:      r.Timestamp,
	}
}

// SaveFeeChargeRecord :
func (model *ModelDB) SaveFeeChargeRecord(r *FeeChargerRecordSerialization) (err error) {
	if r.Key == nil {
		key := utils.NewRandomHash()
		r.Key = key[:]
	}
	if r.Timestamp <= 0 {
		r.Timestamp = time.Now().Unix()
	}
	err = model.db.Save(r)
	if err != nil {
		err = fmt.Errorf("SaveFeeChargeRecord err %s", err)
		return
	}
	log.Trace(fmt.Sprintf("charge for transfer:%s", r.ToFeeChargeRecord().ToString()))
	return
}

// GetAllFeeChargeRecord :
func (model *ModelDB) GetAllFeeChargeRecord() (records []*FeeChargeRecord, err error) {
	var rs []*FeeChargerRecordSerialization
	err = model.db.All(&rs)
	if err != nil {
		err = fmt.Errorf("GetAllFeeChargeRecord err %s", err)
		return
	}
	for _, r := range rs {
		records = append(records, r.ToFeeChargeRecord())
	}
	return
}

// GetFeeChargeRecordByLockSecretHash :
func (model *ModelDB) GetFeeChargeRecordByLockSecretHash(lockSecretHash common.Hash) (records []*FeeChargeRecord, err error) {
	var rs []*FeeChargerRecordSerialization
	err = model.db.Find("LockSecretHash", lockSecretHash[:], &rs)
	if err != nil {
		err = fmt.Errorf("GetAllFeeChargeRecordByLockSecretHash err %s", err)
		return
	}
	for _, r := range rs {
		records = append(records, r.ToFeeChargeRecord())
	}
	return
}
