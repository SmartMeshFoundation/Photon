package models

import (
	"encoding/gob"

	"encoding/json"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TXInfoStatus tx的状态
type TXInfoStatus string

/* #nosec */
const (
	TXInfoStatusPending = "pending"
	TXInfoStatusSuccess = "success"
	TXInfoStatusFailed  = "failed"
)

// TXInfoType 类型
type TXInfoType string

/* #nosec */
const (
	TXInfoTypeDeposit            = "ChannelDeposit"
	TXInfoTypeClose              = "ChannelClose"
	TXInfoTypeSettle             = "ChannelSettle"
	TXInfoTypeCooperateSettle    = "CooperateSettle"
	TXInfoTypeUpdateBalanceProof = "UpdateBalanceProof"
	TXInfoTypeUnlock             = "Unlock"
	TXInfoTypePunish             = "Punish"
	TXInfoTypeWithdraw           = "Withdraw"
	TXInfoTypeApproveDeposit     = "ApproveDeposit"
	TXInfoTypeRegisterSecret     = "RegisterSecret"
)

// TXInfo 记录已经提交到公链节点的tx信息
type TXInfo struct {
	TXHash            common.Hash    `json:"tx_hash"`
	ChannelIdentifier common.Hash    `json:"channel_identifier"` // 结合OpenBlockNumber唯一确定一个通道
	OpenBlockNumber   int64          `json:"open_block_number"`
	TokenAddress      common.Address `json:"token_address"`
	Type              TXInfoType     `json:"type"`
	IsSelfCall        bool           `json:"is_self_call"` // 是否自己发起的
	TXParams          string         `json:"tx_params"`    // 保存调用tx的参数信息,json格式,内容根据TXType不同而不同,仅自己发起的部分tx里面会带此参数
	Status            TXInfoStatus   `json:"tx_status"`
	Events            []interface{}  `json:"events"`            // 保存这个tx成功之后对应的所有事件
	PackBlockNumber   int64          `json:"pack_block_number"` // tx最终所在的块号
	CallTime          int64          `json:"call_time"`         // tx发起时间戳
	PackTime          int64          `json:"pack_time"`         // tx打包时间戳
	GasPrice          uint64         `json:"gas_price"`
	GasUsed           uint64         `json:"gas_used"` // 消耗的gas
}

// String :
func (ti *TXInfo) String() string {
	buf, err := json.MarshalIndent(ti, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// ToTXInfoSerialization :
func (ti *TXInfo) ToTXInfoSerialization() *TXInfoSerialization {
	return &TXInfoSerialization{
		TXHash:            ti.TXHash[:],
		ChannelIdentifier: ti.ChannelIdentifier[:],
		OpenBlockNumber:   ti.OpenBlockNumber,
		TokenAddress:      ti.TokenAddress[:],
		Type:              string(ti.Type),
		IsSelfCall:        ti.IsSelfCall,
		TXParams:          ti.TXParams,
		Status:            string(ti.Status),
		Events:            ti.Events,
		PackBlockNumber:   ti.PackBlockNumber,
		CallTime:          ti.CallTime,
		PackTime:          ti.PackTime,
		GasPrice:          ti.GasPrice,
		GasUsed:           ti.GasUsed,
	}
}

// TXInfoSerialization :
type TXInfoSerialization struct {
	TXHash            []byte        `storm:"id"`
	ChannelIdentifier []byte        `storm:"index"` // 结合OpenBlockNumber唯一确定一个通道
	OpenBlockNumber   int64         `storm:"index"`
	TokenAddress      []byte        `storm:"index"`
	Type              string        `storm:"index"`
	IsSelfCall        bool          `storm:"index"` // 是否自己发起的
	TXParams          string        // 保存调用tx的参数信息,json格式,内容根据TXType不同而不同,仅自己发起的部分tx里面会带此参数
	Status            string        `storm:"index"`
	Events            []interface{} // 保存这个tx成功之后对应的所有事件
	PackBlockNumber   int64         `storm:"index"` // tx最终所在的块号
	CallTime          int64         `storm:"index"`
	PackTime          int64         `storm:"index"`
	GasPrice          uint64
	GasUsed           uint64
}

// ToTXInfo :
func (tis *TXInfoSerialization) ToTXInfo() *TXInfo {
	return &TXInfo{
		TXHash:            common.BytesToHash(tis.TXHash),
		ChannelIdentifier: common.BytesToHash(tis.ChannelIdentifier),
		OpenBlockNumber:   tis.OpenBlockNumber,
		TokenAddress:      common.BytesToAddress(tis.TokenAddress),
		Type:              TXInfoType(tis.Type),
		IsSelfCall:        tis.IsSelfCall,
		TXParams:          tis.TXParams,
		Status:            TXInfoStatus(tis.Status),
		Events:            tis.Events,
		PackBlockNumber:   tis.PackBlockNumber,
		CallTime:          tis.CallTime,
		PackTime:          tis.PackTime,
		GasPrice:          tis.GasPrice,
		GasUsed:           tis.GasUsed,
	}
}

// TXParams tx的参数,自己发起的tx会带上
type TXParams interface{}

// SecretRegisterTxParams 注册密码的参数
type SecretRegisterTxParams struct {
	Secret common.Hash `json:"secret"`
}

// DepositTXParams :
// 1. 保存在ApproveTX的TXParams中,给崩溃恢复后继续deposit使用
// 2. 保存在DepositTX的TXParams中
type DepositTXParams struct {
	TokenAddress       common.Address `json:"token_address"`
	ParticipantAddress common.Address `json:"participant_address"`
	PartnerAddress     common.Address `json:"partner_address"`
	Amount             *big.Int       `json:"amount"`
	SettleTimeout      uint64         `json:"settle_timeout"`
}

// ChannelCloseOrChannelUpdateBalanceProofTXParams 关闭通道或者UpdateBalanceProof的参数,两种操作复用,根据上层TXInfo中的Type区分
type ChannelCloseOrChannelUpdateBalanceProofTXParams struct {
	TokenAddress       common.Address `json:"token_address"`
	ParticipantAddress common.Address `json:"participant_address"`
	PartnerAddress     common.Address `json:"partner_address"`
	TransferAmount     *big.Int       `json:"transfer_amount"`
	LocksRoot          common.Hash    `json:"locks_root"`
	Nonce              uint64         `json:"nonce"`
	ExtraHash          common.Hash    `json:"extra_hash"`
	Signature          []byte         `json:"signature"`
}

// UnlockTXParams 链上Unlock的参数
type UnlockTXParams struct {
	TokenAddress       common.Address `json:"token_address"`
	ParticipantAddress common.Address `json:"participant_address"`
	PartnerAddress     common.Address `json:"partner_address"`
	TransferAmount     *big.Int       `json:"transfer_amount"`
	Expiration         *big.Int       `json:"expiration"`
	Amount             *big.Int       `json:"amount"`
	LockSecretHash     common.Hash    `json:"lock_secret_hash"`
	Proof              []byte         `json:"proof"`
}

// ChannelSettleTXParams 通道结算的参数
type ChannelSettleTXParams struct {
	TokenAddress     common.Address `json:"token_address"`
	P1Address        common.Address `json:"p1_address"`
	P1TransferAmount *big.Int       `json:"p1_transfer_amount"`
	P1LocksRoot      common.Hash    `json:"p1_locks_root"`
	P2Address        common.Address `json:"p2_address"`
	P2TransferAmount *big.Int       `json:"p2_transfer_amount"`
	P2LocksRoot      common.Hash    `json:"p2_locks_root"`
}

// ChannelWithDrawTXParams 通道取现的参数
type ChannelWithDrawTXParams struct {
	TokenAddress common.Address `json:"token_address"`
	P1Address    common.Address `json:"p1_address"`
	P2Address    common.Address `json:"p2_address"`
	P1Balance    *big.Int       `json:"p1_balance"`
	P1Withdraw   *big.Int       `json:"p1_withdraw"`
	P1Signature  []byte         `json:"p1_signature"`
	P2Signature  []byte         `json:"p2_signature"`
}

// PunishObsoleteUnlockTXParams 通道惩罚的参数
type PunishObsoleteUnlockTXParams struct {
	TokenAddress     common.Address `json:"token_address"`
	Beneficiary      common.Address `json:"beneficiary"`
	Cheater          common.Address `json:"cheater"`
	LockHash         common.Hash    `json:"lock_hash"`
	ExtraHash        common.Hash    `json:"extra_hash"`
	CheaterSignature []byte         `json:"cheater_signature"`
}

// ChannelCooperativeSettleTXParams 通道合作关闭的参数
type ChannelCooperativeSettleTXParams struct {
	TokenAddress common.Address `json:"token_address"`
	P1Address    common.Address `json:"p1_address"`
	P1Balance    *big.Int       `json:"p1_balance"`
	P2Address    common.Address `json:"p2_address"`
	P2Balance    *big.Int       `json:"p2_balance"`
	P1Signature  []byte         `json:"p1_signature"`
	P2Signature  []byte         `json:"p2_signature"`
}

func init() {
	gob.Register(&TXInfoSerialization{})
}
