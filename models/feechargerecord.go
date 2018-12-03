package models

import (
	"math/big"

	"encoding/gob"
	"encoding/json"

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

func init() {
	gob.Register(&FeeChargeRecord{})
	gob.Register(&FeeChargerRecordSerialization{})
}
